package sysinfo

import (
	"context"
	"os/exec"
	"strings"
	"syscall"
)

func readWifi(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "netsh", "wlan", "show", "interfaces")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		k, v, ok := strings.Cut(line, ":")
		// “BSSID” cũng chứa “SSID” → so khớp nguyên khoá
		if ok && strings.TrimSpace(k) == "SSID" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
