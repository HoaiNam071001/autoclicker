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

// ValidateScript áp giá trị input rồi kiểm tra như lúc chạy (frontend gọi trước khi thu cửa sổ thành mini bar).
func (a *App) ValidateScript(sc Script, values map[string]any) error {
	out, err := resolveInputs(sc, values)
	if err != nil {
		return err
	}
	return validateScript(out, mainScreen())
}

// RunScript áp giá trị input (resolveInputs) vào bản sao kịch bản, đếm ngược countdown giây rồi chạy; chỉ trả về khi chạy xong / bị dừng
// (StopScript, F8, góc màn hình – kể cả lúc đang đếm ngược) / lỗi. Bị dừng không tính là lỗi.
// Emit "engine:state" (true khi bắt đầu, false khi kết thúc), "engine:countdown" (số giây còn lại)
// và "engine:progress" (Progress) trong lúc chạy.
func (a *App) RunScript(sc Script, values map[string]any, countdown int) error {
	screen := mainScreen()
	sc, err := resolveInputs(sc, values)
	if err != nil {
		return err
	}
	if err := validateScript(sc, screen); err != nil {
		return err
	}
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
	eng.screen = screen
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

	for n := countdown; n > 0; n-- {
		runtime.EventsEmit(a.ctx, "engine:countdown", n)
		if sleep(ctx, 1000) != nil {
			return nil // dừng khi đang đếm ngược
		}
	}
	runtime.EventsEmit(a.ctx, "engine:countdown", 0)

	go watchCorners(ctx, robotgo.Location, hardCorners(displays()), in, func() { a.stopWith(StopByCorner) })

	err = eng.run(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if err != nil && eng.failed != "" {
		// frontend chọn bước lỗi; chờ ảnh hết giờ thì kèm chỗ giống nhất để xem lại
		runtime.EventsEmit(a.ctx, "engine:failed", FailedStep{StepID: eng.failed, Miss: eng.miss})
	}
	return err
}

// FailedStep – bước làm lần chạy dừng vì lỗi (event engine:failed).
type FailedStep struct {
	StepID string      `json:"stepId"`
	Miss   *FindResult `json:"miss,omitempty"` // chờ ảnh hết giờ: chỗ giống nhất
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
