package main

import (
	"context"
	"errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type runHandle struct{ cancel context.CancelFunc }

// RunScript chạy kịch bản và chỉ trả về khi chạy xong / bị dừng (StopScript, F8) / lỗi.
// Bị dừng không tính là lỗi. Emit "engine:state" (true khi bắt đầu, false khi kết thúc).
func (a *App) RunScript(sc Script) error {
	a.mu.Lock()
	if a.run != nil {
		a.mu.Unlock()
		return errors.New("kịch bản đang chạy")
	}
	ctx, cancel := context.WithCancel(context.Background())
	h := &runHandle{cancel: cancel}
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

	err := runScript(ctx, robotgoInput{}, sc)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
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
