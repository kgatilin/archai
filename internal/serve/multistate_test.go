package serve

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/kgatilin/archai/internal/worktree"
)

// newGitRepo creates a minimal git repo with a go.mod and one commit,
// so `git worktree list` has a primary worktree to report.
func newGitRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"),
		[]byte("module example.com/multi\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("go.mod: %v", err)
	}
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=t@e",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=t@e",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q", "-b", "main")
	run("add", ".")
	run("commit", "-qm", "init")
	return root
}

// stubLoader returns a StateLoader that increments *count each call
// and yields a bare State rooted at path (no Go extraction). Used to
// assert lazy-load cache behaviour without the reader's overhead.
func stubLoader(count *int64) StateLoader {
	return func(_ context.Context, _, path string) (*State, error) {
		atomic.AddInt64(count, 1)
		return NewState(path), nil
	}
}

func TestMultiState_RefreshAndGet(t *testing.T) {
	root := newGitRepo(t)
	var loadCount int64
	m := NewMultiState(root, stubLoader(&loadCount))

	if err := m.Refresh(); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	names := m.Names()
	if len(names) != 1 {
		t.Fatalf("want 1 worktree, got %d: %v", len(names), names)
	}

	// First Get triggers load.
	if _, err := m.Get(context.Background(), names[0]); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got := atomic.LoadInt64(&loadCount); got != 1 {
		t.Errorf("load count after first Get = %d, want 1", got)
	}

	// Second Get hits cache — loader must not re-run.
	if _, err := m.Get(context.Background(), names[0]); err != nil {
		t.Fatalf("Get (second): %v", err)
	}
	if got := atomic.LoadInt64(&loadCount); got != 1 {
		t.Errorf("load count after second Get = %d, want 1 (cache miss)", got)
	}

	// Unknown worktree returns an error.
	if _, err := m.Get(context.Background(), "nope"); err == nil {
		t.Errorf("Get(nope) expected error")
	}
}

func TestMultiState_RefreshDropsRemoved(t *testing.T) {
	var loadCount int64
	m := NewMultiState(t.TempDir(), stubLoader(&loadCount))

	// Seed state manually for a worktree that "exists".
	m.entries["alpha"] = worktree.Entry{Name: "alpha", Path: "/tmp/alpha"}
	m.order = []string{"alpha"}
	if _, err := m.Get(context.Background(), "alpha"); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if atomic.LoadInt64(&loadCount) != 1 {
		t.Fatalf("load count after first Get = %d, want 1", loadCount)
	}

	// Simulate a Refresh that yields no worktrees. Drop entries to empty
	// and invoke the drop loop directly — the git fallback would reuse
	// the CWD name, which would keep alpha alive.
	m.mu.Lock()
	m.entries = map[string]worktree.Entry{}
	m.order = nil
	for name := range m.states {
		if _, ok := m.entries[name]; !ok {
			delete(m.states, name)
		}
	}
	m.mu.Unlock()

	// Re-seeding and Get should trigger a fresh load (proving the cache
	// was dropped when the worktree disappeared).
	m.mu.Lock()
	m.entries["alpha"] = worktree.Entry{Name: "alpha", Path: "/tmp/alpha"}
	m.order = []string{"alpha"}
	m.mu.Unlock()
	if _, err := m.Get(context.Background(), "alpha"); err != nil {
		t.Fatalf("Get after drop: %v", err)
	}
	if atomic.LoadInt64(&loadCount) != 2 {
		t.Errorf("load count after re-Get = %d, want 2 (cache should have been dropped)", loadCount)
	}
}
