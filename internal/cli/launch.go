package cli

import (
	"os"
	"os/exec"
)

// launchSelf mở app (chính file exe này, không đối số) tách khỏi terminal – CLI thoát thì app vẫn mở.
func launchSelf() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe)
	detach(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
