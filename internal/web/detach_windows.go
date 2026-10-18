//go:build windows

package web

import (
	"os/exec"
	"syscall"
)

// detach – nhóm tiến trình riêng, không gắn console của app.
func detach(cmd *exec.Cmd) {
	const createNewProcessGroup, detachedProcess = 0x00000200, 0x00000008
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNewProcessGroup | detachedProcess}
}
