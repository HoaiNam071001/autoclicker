//go:build !windows

package cli

import (
	"os/exec"
	"syscall"
)

// AttachConsole – macOS / Linux: tiến trình đã có stdout của terminal.
func AttachConsole() {}

func detach(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} }
