package device

import (
	"syscall"
	"unsafe"

	"ghosthand/internal/model"
)

var (
	procMonitorFromPoint = user32.NewProc("MonitorFromPoint")
	procGetDpiForMonitor = syscall.NewLazyDLL("shcore.dll").NewProc("GetDpiForMonitor")
)

// Khai báo per-monitor DPI aware ngay khi nạp package (trước khi Wails tạo cửa sổ): app đã có manifest
// (build/windows/wails.exe.manifest) nhưng exe của go test / wails dev / tool thử thì không – thiếu thì Windows ảo hoá
// chuột theo toạ độ logic (125 %: tối đa 1535×863) trong khi khung màn hình và ảnh chụp vẫn là pixel thật → lệch.
// Tiến trình đã đặt (manifest) thì lời gọi thất bại, giữ nguyên – vô hại.
func init() {
	if p := user32.NewProc("SetProcessDpiAwarenessContext"); p.Find() == nil {
		const perMonitorAwareV2 = ^uintptr(3) // DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 = (HANDLE)-4, Windows 10 1703+
		p.Call(perMonitorAwareV2)
	}
}

// DPIScale – tỷ lệ scale (125 % → 1,25) của màn chứa tâm r: toạ độ là px thật nên cùng giao diện ở màn 125 % to hơn
// màn 100 % – ảnh mẫu co giãn theo tỷ lệ này. Không đọc được → 1.
func DPIScale(r model.Rect) float64 {
	if procGetDpiForMonitor.Find() != nil {
		return 1
	}
	// POINT truyền theo giá trị (x64: gói 2 int32 vào một thanh ghi); MONITOR_DEFAULTTONEAREST = 2
	pt := uintptr(uint32(int32(r.X+r.W/2))) | uintptr(uint32(int32(r.Y+r.H/2)))<<32
	mon, _, _ := procMonitorFromPoint.Call(pt, 2)
	var dx, dy uint32
	if hr, _, _ := procGetDpiForMonitor.Call(mon, 0, uintptr(unsafe.Pointer(&dx)), uintptr(unsafe.Pointer(&dy))); hr != 0 || dx == 0 {
		return 1 // MDT_EFFECTIVE_DPI = 0
	}
	return float64(dx) / 96
}
