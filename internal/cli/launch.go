package cli

import (
	"os"
	"os/exec"
	"slices"
)

// BackgroundFlag – app mở nền (bước 31b): khởi động ẩn, AI chạy thì chỉ hiện thanh mini. `--mcp` mở app kiểu này.
const BackgroundFlag = "--background"

// IsBackground – args (không gồm tên chương trình) có cờ mở nền.
func IsBackground(args []string) bool { return slices.Contains(args, BackgroundFlag) }

// launchSelf mở app (chính file exe này) tách khỏi terminal – CLI thoát thì app vẫn mở.
func launchSelf() error { return launch() }

// LaunchBackground – như launchSelf nhưng app mở nền (cho `--mcp`).
func LaunchBackground() error { return launch(BackgroundFlag) }

func launch(args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, args...)
	detach(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
