package serve

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kgatilin/archai/internal/worktree"
)

func TestDiscoverDaemon_NoRecord(t *testing.T) {
	root := t.TempDir()
	rec, _, err := DiscoverDaemon(root)
	if err != nil {
		t.Fatalf("DiscoverDaemon: %v", err)
	}
	if rec != nil {
		t.Errorf("expected nil record, got %+v", rec)
	}
}

func TestDiscoverDaemon_LiveRecord(t *testing.T) {
	root := t.TempDir()
	name := worktree.Name(root)
	// Write a serve.json whose PID is this test process — guaranteed alive.
	rec := worktree.ServeRecord{
		PID:       os.Getpid(),
		HTTPAddr:  "127.0.0.1:12345",
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := worktree.WriteServe(root, name, rec); err != nil {
		t.Fatalf("WriteServe: %v", err)
	}

	got, _, err := DiscoverDaemon(root)
	if err != nil {
		t.Fatalf("DiscoverDaemon: %v", err)
	}
	if got == nil {
		t.Fatal("expected record, got nil")
	}
	if got.PID != rec.PID || got.HTTPAddr != rec.HTTPAddr {
		t.Errorf("mismatch: got=%+v want=%+v", got, rec)
	}
}

func TestDiscoverDaemon_StaleRecord(t *testing.T) {
	root := t.TempDir()
	name := worktree.Name(root)
	// PID 0 is never a live process on Unix — PIDAlive returns false.
	rec := worktree.ServeRecord{
		PID:       0,
		HTTPAddr:  "127.0.0.1:1",
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}
	// Write directly so the guard fires without racing real processes.
	if err := os.MkdirAll(filepath.Dir(worktree.ServePath(root, name)), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := worktree.WriteServe(root, name, rec); err != nil {
		t.Fatalf("WriteServe: %v", err)
	}

	got, _, err := DiscoverDaemon(root)
	if err != nil {
		t.Fatalf("DiscoverDaemon: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for stale record, got %+v", got)
	}
}
