package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Chế độ cửa sổ. Wails v2 chỉ có một cửa sổ nên overlay / mini bar là cách đổi kích thước,
// vị trí và always-on-top của cửa sổ chính (cửa sổ luôn frameless + nền trong suốt, xem main.go).
const (
	ModeMain    = "main"
	ModeOverlay = "overlay"
	ModeMinibar = "minibar"
)

const miniBarW, miniBarH = 460, 76

type winRect struct{ x, y, w, h int }

// SetWindowMode chuyển cửa sổ sang main | overlay | minibar.
// Lưu vị trí/kích thước của main để khôi phục khi quay lại.
func (a *App) SetWindowMode(mode string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if mode == a.mode {
		return nil
	}

	if a.mode == ModeMain {
		x, y := runtime.WindowGetPosition(a.ctx)
		w, h := runtime.WindowGetSize(a.ctx)
		a.mainRect = winRect{x, y, w, h}
	}

	switch mode {
	case ModeMain:
		r := a.mainRect
		runtime.WindowSetAlwaysOnTop(a.ctx, false)
		runtime.WindowSetSize(a.ctx, r.w, r.h)
		runtime.WindowSetPosition(a.ctx, r.x, r.y)
	case ModeOverlay, ModeMinibar:
		sw, sh, err := a.currentScreenSize()
		if err != nil {
			return err
		}
		runtime.WindowSetAlwaysOnTop(a.ctx, true)
		if mode == ModeOverlay {
			// Dời về góc màn hiện tại TRƯỚC rồi mới phóng to: Wails lấy "màn hiện tại" là màn chứa phần lớn
			// cửa sổ, phóng to trước thì cửa sổ có thể tràn sang màn bên cạnh và overlay mở nhầm màn.
			// SetSize giữ nguyên mép trên nên cửa sổ vẫn ở góc màn này.
			// macOS: vị trí tính từ vùng hiển thị (dưới menu bar) → overlay tự đo độ lệch toạ độ.
			runtime.WindowSetPosition(a.ctx, 0, 0)
			runtime.WindowSetSize(a.ctx, sw, sh)
		} else {
			runtime.WindowSetSize(a.ctx, miniBarW, miniBarH)
			runtime.WindowSetPosition(a.ctx, (sw-miniBarW)/2, 16)
		}
	default:
		return fmt.Errorf("chế độ cửa sổ không hợp lệ: %q", mode)
	}
	a.mode = mode
	return nil
}

// currentScreenSize trả kích thước (logical px) của màn hình chứa cửa sổ.
func (a *App) currentScreenSize() (int, int, error) {
	screens, err := runtime.ScreenGetAll(a.ctx)
	if err != nil {
		return 0, 0, err
	}
	for _, s := range screens {
		if s.IsCurrent {
			return s.Size.Width, s.Size.Height, nil
		}
	}
	for _, s := range screens {
		if s.IsPrimary {
			return s.Size.Width, s.Size.Height, nil
		}
	}
	return 0, 0, fmt.Errorf("không tìm thấy màn hình")
}

// MousePos trả vị trí chuột theo robotgo (toạ độ màn hình thật).
// Overlay so với clientX/Y để đo độ lệch giữa toạ độ trong cửa sổ và toạ độ click.
func (a *App) MousePos() Point {
	x, y := robotgo.Location()
	return Point{X: x, Y: y}
}
