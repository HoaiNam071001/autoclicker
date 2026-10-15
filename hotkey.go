package main

import (
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// startHotkeys nghe phím toàn cục (kể cả khi app không được focus). Gọi trong startup.
// macOS: cần quyền Accessibility; phím F8 mặc định là phím media nên phải bấm fn+F8.
func (a *App) startHotkeys() {
	f8 := hook.KeycharToRawcode("f8")
	evs := hook.Start()
	go func() {
		for ev := range evs {
			switch {
			case ev.Kind == hook.HookDisabled:
				runtime.EventsEmit(a.ctx, "hotkey:error", "Không nghe được phím toàn cục – kiểm tra quyền Accessibility")
			case ev.Kind == hook.KeyDown && ev.Rawcode == f8:
				runtime.EventsEmit(a.ctx, "hotkey:f8")
				a.StopScript()
			}
		}
	}()
}

func (a *App) stopHotkeys() {
	hook.End()
}
