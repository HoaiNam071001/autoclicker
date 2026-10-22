package sysinfo

import (
	"context"
	"os/exec"
	"strings"
)

// readWifi – tên Wi-Fi đang nối; macOS mới có thể giấu tên nếu app chưa có quyền Vị trí → "".
func readWifi(ctx context.Context) string {
	for _, dev := range []string{"en0", "en1"} {
		out, err := exec.CommandContext(ctx, "networksetup", "-getairportnetwork", dev).Output()
		if err != nil {
			continue
		}
		if _, name, ok := strings.Cut(strings.TrimSpace(string(out)), ": "); ok && !strings.HasPrefix(string(out), "You are not") {
			return strings.TrimSpace(name)
		}
	}
	return ""
}
