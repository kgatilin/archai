package serve

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/kgatilin/archai/internal/worktree"
)

// StateLoader builds a fresh State for a newly-discovered worktree.
// It is invoked the first time a worktree's State is requested from a
// MultiState; subsequent requests hit the per-worktree cache. The
// default production loader is DefaultStateLoader, which calls
// NewState(path).Load(ctx).
type StateLoader func(ctx context.Context, name, path string) (*State, error)

// DefaultStateLoader is the production StateLoader used when
// NewMultiState is called without one. It builds a fresh State
// anchored at path and loads the full Go + overlay + target model.
func DefaultStateLoader(ctx context.Context, _, path string) (*State, error) {
	state := NewState(path)
	if err := state.Load(ctx); err != nil {
		return nil, err
	}
	return state, nil
}

// MultiState holds one State per discovered worktree and lazy-loads
// each State on first access. The set of worktrees is fixed at
// construction but can be refreshed by calling Refresh() again — this
// is intended for SIGHUP or a periodic poller.
//
// MultiState is safe for concurrent use.
type MultiState struct {
	mu sync.Mutex

	// root is the project root (one of the discovered worktrees, used
	// as the anchor for git worktree list).
	root string

	// entries is the current list of known worktrees, keyed by Name.
	entries map[string]worktree.Entry

	// order is the lexical list of names (matches entries).
	order []string

	// states is the lazy-loaded State for each worktree name.
	states map[string]*State

	// loader is the factory that builds a fresh State for a worktree.
	loader StateLoader
}

// NewMultiState constructs a MultiState rooted at projectRoot, using
// loader to materialize per-worktree States on first access. Pass
// DefaultStateLoader for normal daemon use; lightweight alternatives
// are useful for transport-level tests that want to assert lazy
// behaviour without invoking the Go reader.
//
// The initial worktree list is populated by Refresh(); callers are
// expected to invoke Refresh() once before serving requests.
func NewMultiState(projectRoot string, loader StateLoader) *MultiState {
	if loader == nil {
		loader = DefaultStateLoader
	}
	return &MultiState{
		root:    projectRoot,
		entries: make(map[string]worktree.Entry),
		states:  make(map[string]*State),
		loader:  loader,
	}
}

// Refresh re-discovers worktrees via `git worktree list --porcelain`
// and replaces the internal entry set. Previously-loaded States for
// worktrees that still exist are retained (so lazy caches survive a
// refresh); States for removed worktrees are dropped.
func (m *MultiState) Refresh() error {
	entries, err := worktree.Discover(m.root)
	if err != nil {
		return fmt.Errorf("serve: discover worktrees: %w", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	next := make(map[string]worktree.Entry, len(entries))
	order := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.Name == "" {
			continue
		}
		if _, dup := next[e.Name]; dup {
			// Disambiguate duplicate leaf names by suffixing the branch
			// (rare — two worktrees at paths sharing basename).
			alt := e.Name + "@" + e.Branch
			e.Name = alt
		}
		next[e.Name] = e
		order = append(order, e.Name)
	}
	sort.Strings(order)
	m.entries = next
	m.order = order

	// Drop cached states whose worktrees have disappeared.
	for name := range m.states {
		if _, ok := next[name]; !ok {
			delete(m.states, name)
		}
	}
	return nil
}

// Worktrees returns the discovered entries in lexical order.
func (m *MultiState) Worktrees() []worktree.Entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]worktree.Entry, 0, len(m.order))
	for _, n := range m.order {
		out = append(out, m.entries[n])
	}
	return out
}

// Names returns the discovered worktree names in lexical order.
func (m *MultiState) Names() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]string, len(m.order))
	copy(out, m.order)
	return out
}

// Has reports whether name is a known worktree.
func (m *MultiState) Has(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.entries[name]
	return ok
}

// Default returns the first worktree name in lexical order, or ""
// when no worktrees have been discovered.
func (m *MultiState) Default() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.order) == 0 {
		return ""
	}
	return m.order[0]
}

// Get returns (and lazy-loads) the State for the given worktree name.
// Returns an error when name is unknown or when the underlying load
// fails. Subsequent calls return the cached State.
func (m *MultiState) Get(ctx context.Context, name string) (*State, error) {
	m.mu.Lock()
	entry, ok := m.entries[name]
	if !ok {
		m.mu.Unlock()
		return nil, fmt.Errorf("serve: unknown worktree %q", name)
	}
	if s, ok := m.states[name]; ok {
		m.mu.Unlock()
		return s, nil
	}
	m.mu.Unlock()

	// Load outside the lock so concurrent Get calls for different
	// worktrees don't serialize. We re-check the cache after loading
	// so only one State is ever cached per name.
	loaded, err := m.loader(ctx, name, entry.Path)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if existing, ok := m.states[name]; ok {
		return existing, nil
	}
	m.states[name] = loaded
	return loaded, nil
}

// Entry returns the Entry for the named worktree (path, branch, …).
// Returns false when name is unknown.
func (m *MultiState) Entry(name string) (worktree.Entry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[name]
	return e, ok
}
