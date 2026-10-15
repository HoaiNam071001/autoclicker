package main

import (
	"context"
	"errors"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Lý do dừng gửi kèm event "engine:stopped".
const (
	StopByF8     = "f8"
	StopByCorner = "corner"
)

type runHandle struct {
	cancel context.CancelFunc
	eng    *engine
}

// RunScript chạy kịch bản và chỉ trả về khi chạy xong / bị dừng (StopScript, F8) / lỗi.
// Bị dừng không tính là lỗi. Emit "engine:state" (true khi bắt đầu, false khi kết thúc)
// và "engine:progress" (Progress) trong lúc chạy.
func (a *App) RunScript(sc Script) error {
	a.mu.Lock()
	if a.run != nil {
		a.mu.Unlock()
		return errors.New("kịch bản đang chạy")
	}
	ctx, cancel := context.WithCancel(context.Background())
	in := &trackingInput{Input: robotgoInput{}}
	eng := newEngine(in, sc, func(p Progress) {
		runtime.EventsEmit(a.ctx, "engine:progress", p)
	})
	h := &runHandle{cancel: cancel, eng: eng}
	a.run = h
	a.mu.Unlock()

	defer func() {
		cancel()
		a.mu.Lock()
		if a.run == h {
			a.run = nil
		}
		a.mu.Unlock()
		runtime.EventsEmit(a.ctx, "engine:state", false)
	}()
	runtime.EventsEmit(a.ctx, "engine:state", true)

	sw, sh := robotgo.GetScreenSize()
	go watchCorners(ctx, robotgo.Location, sw, sh, in, func() { a.stopWith(StopByCorner) })

	err := eng.run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

// stopWith dừng kịch bản đang chạy và emit "engine:stopped" kèm lý do (F8, góc màn hình).
func (a *App) stopWith(reason string) {
	if a.current() == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "engine:stopped", reason)
	a.StopScript()
}

// StopScript dừng kịch bản đang chạy (không làm gì nếu đang không chạy).
func (a *App) StopScript() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.run != nil {
		a.run.cancel()
		a.run = nil
	}
}

// PauseScript tạm dừng trước bước / lần click kế tiếp. Trả false nếu không có gì để tạm dừng.
func (a *App) PauseScript() bool {
	if h := a.current(); h != nil {
		return h.eng.Pause()
	}
	return false
}

// ResumeScript chạy tiếp sau khi tạm dừng.
func (a *App) ResumeScript() bool {
	if h := a.current(); h != nil {
		return h.eng.Resume()
	}
	return false
}

func (a *App) current() *runHandle {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.run
}
