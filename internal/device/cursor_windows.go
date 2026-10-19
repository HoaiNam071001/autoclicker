package device

import (
	"syscall"
	"unsafe"
)

// Chuột trên Windows đọc / đặt thẳng qua Win32, không qua robotgo.Location / Move: app khai báo per-monitor DPI aware
// (build/windows/wails.exe.manifest) nên GetCursorPos / SetCursorPos, khung màn hình và ảnh chụp đều theo pixel thật;
// robotgo vẫn nhân / chia thêm tỷ lệ DPI của *cửa sổ đang ở trên cùng* (GetDpiForWindow) → màn 125 % lúc ghi điểm và
// lúc chạy lệch nhau tuỳ cửa sổ nào đang trước (lỗi người dùng báo 2026-10-19).
var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procSetCursorPos = user32.NewProc("SetCursorPos")
)

func cursorPos() (int, int) {
	var p struct{ X, Y int32 }
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	return int(p.X), int(p.Y)
}

func moveCursor(x, y int) error {
	// int32 → uintptr giữ dấu: màn phụ bên trái / trên màn chính có toạ độ âm
	if r, _, err := procSetCursorPos.Call(uintptr(int32(x)), uintptr(int32(y))); r == 0 {
		return err
	}
	return nil
}

// Kéo trên Windows: nút đang giữ (Toggle) thì đặt chuột là app nhận WM_MOUSEMOVE kèm trạng thái nút – đủ để kéo.
func dragCursor(x, y int, _ string) error { return moveCursor(x, y) }
