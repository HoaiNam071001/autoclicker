//go:build !windows

package web

import (
	"os/exec"
	"syscall"
)

// detach – phiên riêng: tín hiệu tắt app (Ctrl+C khi wails dev) không lan sang trình duyệt.
func detach(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} }
