//go:build windows

package serve

import "os/exec"

// detachProcess is a no-op on Windows — exec.Cmd inherits enough of
// its environment that the daemon survives the parent's exit in
// practice. We don't use DETACHED_PROCESS here because it would cost
// us the stderr pipe used during diagnostics.
func detachProcess(cmd *exec.Cmd) {}
