package device

import (
	"errors"

	"ghosthand/internal/model"
)

// Cửa sổ + ứng dụng đang chạy (bước 26: node Cửa sổ / Ứng dụng, tool MCP list_windows / list_apps). Qua interface để
// test engine bằng bản giả (devicetest.FakeWindows / FakeApps). Kiểu dữ liệu ở model/window.go. Windows: Win32
// (window_windows.go).

// Windows – liệt kê / điều khiển cửa sổ. List theo thứ tự chồng (trên cùng trước).
type Windows interface {
	List() ([]model.WindowInfo, error)
	Focus(id uint64) error                   // đưa lên trước, bỏ thu nhỏ
	SetBounds(id uint64, r model.Rect) error // đặt khung nhìn thấy (bỏ thu nhỏ / phóng to trước)
	SetState(id uint64, state string) error  // WinMinimize / WinMaximize / WinRestore
	Close(id uint64) error                   // như bấm × – app có thể hỏi lưu
	WorkArea(r model.Rect) model.Rect        // vùng làm việc (trừ taskbar / thanh menu) của màn chứa r
}

// Trạng thái cho Windows.SetState.
const (
	WinMinimize = "minimize"
	WinMaximize = "maximize"
	WinRestore  = "restore"
)

// Apps – liệt kê / thoát app. Quit thoát bình thường (đóng mọi cửa sổ – app có thể hỏi lưu), ForceQuit buộc thoát.
type Apps interface {
	List() ([]model.AppInfo, error)
	Quit(pid int) error
	ForceQuit(pid int) error
}

type SystemWindows struct{}

type SystemApps struct{}

// ErrWindowGone – cửa sổ đã đóng / không còn (id cũ).
var ErrWindowGone = errors.New("cửa sổ không còn mở")
