//go:build !windows

package serve

import (
	"os/exec"
	"syscall"
)

// detachProcess configures the child to run in a new session/process
// group so it survives when the parent exits. Without Setsid the child
// would share the parent's controlling terminal and receive SIGHUP
// when the MCP client closes its stdio.
func detachProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}
