package main

import (
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// startHotkeys nghe phím toàn cục (kể cả khi app không được focus). Gọi trong startup.
//   - F8: dừng kịch bản đang chạy
//   - F6: lưu vị trí chuột hiện tại thành điểm mới – emit "hotkey:f6" (Point); bỏ qua khi đang chạy
//
// macOS: cần quyền Accessibility; F6 / F8 mặc định là phím media nên phải bấm kèm fn.
func (a *App) startHotkeys() {
	f6 := hook.KeycharToRawcode("f6")
	f8 := hook.KeycharToRawcode("f8")
	evs := hook.Start()
	go func() {
		for ev := range evs {
			switch {
			case ev.Kind == hook.HookDisabled:
				runtime.EventsEmit(a.ctx, "hotkey:error", "Không nghe được phím toàn cục – kiểm tra quyền Accessibility")
			case ev.Kind == hook.KeyDown && ev.Rawcode == f8:
				a.stopWith(StopByF8)
			case ev.Kind == hook.KeyDown && ev.Rawcode == f6 && a.current() == nil:
				x, y := robotgo.Location()
				runtime.EventsEmit(a.ctx, "hotkey:f6", Point{X: x, Y: y})
			}
		}
	}()
}

func (a *App) stopHotkeys() {
	hook.End()
}
