// Package sysinfo đọc thông tin máy cho node Thông tin máy (bước 43 S4): pin, CPU / RAM, ổ đĩa, màn hình, mạng,
// tên máy. Phần dùng chung viết bằng gopsutil; pin và Wi-Fi riêng từng hệ điều hành (battery_*.go, wifi_*.go).
package sysinfo

import (
	"context"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"

	"ghosthand/internal/device"
)

// Backend – đọc toàn bộ thông tin máy một lần. diskPath trống = thư mục người dùng. Test thay bản giả.
type Backend interface {
	Read(ctx context.Context, diskPath string) (map[string]any, error)
}

// System – Backend thật.
type System struct{}

// Battery – kết quả đọc pin; Has = false khi máy không có pin (máy bàn).
type Battery struct {
	Has      bool
	Percent  int
	Charging bool
}

func (System) Read(ctx context.Context, diskPath string) (map[string]any, error) {
	out := map[string]any{"os": runtime.GOOS}
	out["hostname"], _ = os.Hostname()
	if u, err := user.Current(); err == nil {
		out["user"] = u.Username
	} else {
		out["user"] = ""
	}

	b := readBattery(ctx)
	out["hasBattery"], out["charging"] = b.Has, b.Charging
	out["batteryPercent"] = 100 // máy bàn: coi như đầy, để điều kiện “pin < 20 %” không dừng nhầm
	if b.Has {
		out["batteryPercent"] = b.Percent
	}

	if p, err := cpu.PercentWithContext(ctx, 300*time.Millisecond, false); err == nil && len(p) > 0 {
		out["cpuPercent"] = round1(p[0])
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	} else {
		out["cpuPercent"] = 0.0
	}
	if m, err := mem.VirtualMemoryWithContext(ctx); err == nil {
		out["memUsedPercent"] = round1(m.UsedPercent)
		out["memTotalMB"] = int(m.Total >> 20)
		out["memFreeMB"] = int(m.Available >> 20)
	}
	if u, err := disk.UsageWithContext(ctx, resolveDiskPath(diskPath)); err == nil {
		out["diskFreeGB"] = round1(float64(u.Free) / (1 << 30))
		out["diskTotalGB"] = round1(float64(u.Total) / (1 << 30))
	}

	d := device.Displays()
	out["displays"] = len(d)
	out["screenWidth"], out["screenHeight"], out["scale"] = 0, 0, 1.0
	if len(d) > 0 {
		out["screenWidth"], out["screenHeight"], out["scale"] = d[0].W, d[0].H, screenScale(d[0])
	}

	ip := localIP()
	out["online"], out["ip"] = ip != "", ip
	out["wifi"] = readWifi(ctx)
	return out, nil
}

// resolveDiskPath – trống / ~ = thư mục người dùng (cùng quy ước đường dẫn của node Tệp).
func resolveDiskPath(p string) string {
	p = strings.TrimSpace(p)
	home, _ := os.UserHomeDir()
	switch {
	case p == "" || p == "~":
		return home
	case strings.HasPrefix(p, "~/") || strings.HasPrefix(p, `~\`):
		return filepath.Join(home, p[2:])
	}
	return p
}

// localIP – IPv4 của giao diện mạng đang dùng (không gửi gói nào: UDP “connect” chỉ chọn tuyến). Không có mạng → "".
func localIP() string {
	c, err := net.Dial("udp4", "192.0.2.1:9")
	if err != nil {
		return ""
	}
	defer c.Close()
	if a, ok := c.LocalAddr().(*net.UDPAddr); ok && !a.IP.IsLoopback() {
		return a.IP.String()
	}
	return ""
}

func round1(f float64) float64 { return float64(int(f*10+0.5)) / 10 }
