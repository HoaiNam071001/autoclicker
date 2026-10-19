package engine

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func winStep(op string, f func(s *model.Step)) model.Step {
	return modeltest.Step(model.StepWindow, func(s *model.Step) {
		s.WinOp, s.WinApp, s.TimeoutMs = op, "zalo", 1000
		if f != nil {
			f(s)
		}
	})
}

func winEngine() (*Engine, *devicetest.FakeWindows) {
	fw := &devicetest.FakeWindows{
		Wins: []model.WindowInfo{
			{ID: 1, Title: "Sân thử – Google Chrome", App: "chrome", PID: 10, Box: model.Rect{X: 100, Y: 100, W: 800, H: 600}},
			{ID: 2, Title: "Zalo", App: "Zalo", PID: 20, Minimized: true},
		},
		Work: model.Rect{X: 0, Y: 0, W: 1920, H: 1040},
	}
	e := New(&devicetest.FakeInput{}, model.Script{}, Options{Displays: []model.Rect{{W: 1920, H: 1080}, {X: 1920, W: 1280, H: 1024}}})
	e.wins = fw
	return e, fw
}

func TestWindowOps(t *testing.T) {
	for _, c := range []struct {
		step model.Step
		log  string
		box  model.Rect
	}{
		{winStep(model.WinFocus, nil), "focus 2", model.Rect{}},
		{winStep(model.WinSetBounds, func(s *model.Step) { s.WinRect = &model.Rect{X: 5, Y: 6, W: 700, H: 500} }), "bounds 2 5,6 700x500", model.Rect{X: 5, Y: 6, W: 700, H: 500}},
		{winStep(model.WinSetBounds, func(s *model.Step) {
			s.WinPreset, s.WinRect = model.PresetRight, &model.Rect{W: 1280, H: 800}
		}), "bounds 2 640,0 1280x800", model.Rect{X: 640, W: 1280, H: 800}},
		{winStep(model.WinMin, nil), "minimize 2", model.Rect{}},
		{winStep(model.WinClose, nil), "close 2", model.Rect{}},
		{winStep(model.WinRead, func(s *model.Step) { s.WinApp, s.WinTitle = "", "sân thử" }), "", model.Rect{X: 100, Y: 100, W: 800, H: 600}},
	} {
		e, fw := winEngine()
		out, err := e.windowStep(t.Context(), c.step)
		if err != nil {
			t.Fatalf("%q: %v", c.step.WinOp, err)
		}
		if got := slices.Concat(fw.Log, []string{""})[0]; got != c.log {
			t.Errorf("%q: log %v, muốn %q", c.step.WinOp, fw.Log, c.log)
		}
		b := c.box
		if out["w"] != b.W || out["h"] != b.H || out["left"] != b.X || out["top"] != b.Y || out["x"] != b.X+b.W/2 {
			t.Errorf("%q: output %v, muốn khung %v", c.step.WinOp, out, b)
		}
	}
}

func TestWindowOutputDisplay(t *testing.T) {
	e, _ := winEngine()
	out := e.windowOutput(model.WindowInfo{Box: model.Rect{X: 2000, Y: 10, W: 400, H: 300}})
	if out["display"] != 2 || out["x"] != 2200 || out["y"] != 160 {
		t.Errorf("output %v", out)
	}
}

func TestWindowPresetOnDisplay(t *testing.T) {
	e, fw := winEngine()
	s := winStep(model.WinSetBounds, func(s *model.Step) { s.WinPreset, s.Display = model.PresetFull, 3 })
	var bad *invalidErr
	if _, err := e.windowStep(t.Context(), s); !errors.As(err, &bad) {
		t.Errorf("màn 3 không có: lỗi %v, muốn invalid", err)
	}
	s.Display = 2
	if _, err := e.windowStep(t.Context(), s); err != nil || fw.Log[0] != "bounds 2 0,0 1920x1040" {
		t.Errorf("full: %v %v", err, fw.Log) // bản giả: mọi màn cùng vùng làm việc Work
	}
}

func TestWindowWaitOpenClose(t *testing.T) {
	e, fw := winEngine()
	go func() {
		time.Sleep(300 * time.Millisecond)
		_ = fw.Close(2)
	}()
	start := time.Now()
	if _, err := e.windowStep(t.Context(), winStep(model.WinWaitClose, nil)); err != nil {
		t.Fatal(err)
	}
	if d := time.Since(start); d < 250*time.Millisecond {
		t.Errorf("chờ đóng xong quá sớm: %v", d)
	}
	// đã đóng: chờ mở hết giờ → timeout, focus hết giờ → not_found
	var late *timeoutErr
	if _, err := e.windowStep(t.Context(), winStep(model.WinWaitOpen, func(s *model.Step) { s.TimeoutMs = 300 })); !errors.As(err, &late) {
		t.Errorf("chờ mở: %v, muốn timeout", err)
	}
	var nf *notFoundErr
	if _, err := e.windowStep(t.Context(), winStep(model.WinFocus, func(s *model.Step) { s.TimeoutMs = 0 })); !errors.As(err, &nf) {
		t.Errorf("focus: %v, muốn not_found", err)
	}
}

func TestWindowStopWhileWaiting(t *testing.T) {
	e, _ := winEngine()
	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := e.windowStep(ctx, winStep(model.WinWaitOpen, func(s *model.Step) { s.WinApp, s.TimeoutMs = "telegram", 10000 }))
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Errorf("dừng giữa chừng: %v sau %v", err, time.Since(start))
	}
}
