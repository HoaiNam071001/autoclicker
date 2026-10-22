//go:build !windows

package shell

import (
	"os/exec"
	"runtime"
	"syscall"
)

// shellCommand – shell mặc định: zsh trên macOS, sh nơi khác.
func shellCommand(script string) *exec.Cmd {
	if runtime.GOOS == "darwin" {
		return exec.Command("/bin/zsh", "-c", script)
	}
	return exec.Command("/bin/sh", "-c", script)
}

// own cho lệnh một nhóm tiến trình riêng để giết được cả cây.
func own(cmd *exec.Cmd) { cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} }

func killTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
