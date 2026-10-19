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

func appStep(op, name string, f func(s *model.Step)) model.Step {
	return modeltest.Step(model.StepApp, func(s *model.Step) {
		s.AppOp, s.AppName, s.TimeoutMs = op, name, 1000
		if f != nil {
			f(s)
		}
	})
}

// appEngine – Chrome 2 tiến trình (một có cửa sổ), Notepad chưa có cửa sổ, explorer.
func appEngine(allow bool) (*Engine, *devicetest.FakeApps) {
	fa := &devicetest.FakeApps{Apps: []model.AppInfo{
		{Name: "chrome", ID: "C:/Chrome/chrome.exe", PID: 10},
		{Name: "chrome", ID: "C:/Chrome/chrome.exe", PID: 11, Windows: 2},
		{Name: "Notepad", PID: 20},
		{Name: "explorer", PID: 30, Windows: 1},
	}}
	e := New(&devicetest.FakeInput{}, model.Script{}, Options{AllowDanger: allow})
	e.apps = fa
	return e, fa
}

func TestAppCheckPorts(t *testing.T) {
	e, _ := appEngine(false)
	for name, port := range map[string]string{"chrome.exe": model.PortYes, "telegram": model.PortNo} {
		res, _ := e.execNode(t.Context(), appStep(model.AppCheck, name, nil), 1)
		if res.Status != model.StatusSuccess || res.Port != port {
			t.Errorf("%s: %s %s, muốn cổng %s", name, res.Status, res.Port, port)
		}
	}
	res, _ := e.execNode(t.Context(), appStep(model.AppCheck, "chrome", nil), 1)
	if o := res.Output; o["pid"] != 11 || o["windows"] != 2 || o["running"] != true {
		t.Errorf("output: %v – muốn tiến trình có cửa sổ, tổng cửa sổ", o)
	}
	// node khác vẫn cổng success
	if res, _ := e.execNode(t.Context(), appStep(model.AppWaitQuit, "telegram", nil), 1); res.Port != model.PortSuccess {
		t.Errorf("chờ thoát: cổng %s", res.Port)
	}
}

func TestAppWaitStart(t *testing.T) {
	e, fa := appEngine(false)
	go func() {
		time.Sleep(300 * time.Millisecond)
		fa.SetWindows(20, 1)
	}()
	start := time.Now()
	out, err := e.appStep(t.Context(), appStep(model.AppWaitStart, "notepad", nil))
	if err != nil || out["windows"] != 1 || time.Since(start) < 250*time.Millisecond {
		t.Fatalf("chờ app có cửa sổ: %v %v sau %v", out, err, time.Since(start))
	}
	var late *timeoutErr
	if _, err := e.appStep(t.Context(), appStep(model.AppWaitStart, "telegram", func(s *model.Step) { s.TimeoutMs = 200 })); !errors.As(err, &late) {
		t.Errorf("app không chạy: %v, muốn timeout", err)
	}
}

func TestAppQuit(t *testing.T) {
	e, fa := appEngine(false)
	fa.SetWindows(20, 1)
	out, err := e.appStep(t.Context(), appStep(model.AppQuit, "notepad", nil))
	if err != nil || out["running"] != false || out["pid"] != 20 {
		t.Fatalf("thoát: %v %v", out, err)
	}
	// chỉ gửi đóng tới tiến trình có cửa sổ; chờ mọi tiến trình của app thoát (tiến trình nền 10 của bản giả không tự
	// thoát như Chrome thật → hết giờ)
	var late *timeoutErr
	if _, err := e.appStep(t.Context(), appStep(model.AppQuit, "chrome", func(s *model.Step) { s.TimeoutMs = 200 })); !errors.As(err, &late) {
		t.Errorf("còn tiến trình nền: %v, muốn timeout", err)
	}
	if !slices.Equal(fa.Log, []string{"quit 20", "quit 11"}) {
		t.Errorf("log %v", fa.Log)
	}
	// không chạy → thành công ngay
	if _, err := e.appStep(t.Context(), appStep(model.AppQuit, "telegram", nil)); err != nil {
		t.Errorf("app không chạy: %v", err)
	}
}

func TestAppQuitAskSave(t *testing.T) {
	e, fa := appEngine(true)
	fa.AskSave = true
	var late *timeoutErr
	_, err := e.appStep(t.Context(), appStep(model.AppQuit, "notepad", func(s *model.Step) { s.TimeoutMs = 200 }))
	if !errors.As(err, &late) {
		t.Fatalf("app hỏi lưu: %v, muốn timeout", err)
	}
	fa.SetWindows(20, 1)
	if _, err := e.appStep(t.Context(), appStep(model.AppQuit, "notepad", func(s *model.Step) { s.TimeoutMs, s.AppForce = 200, true })); err != nil {
		t.Fatalf("quá giờ thì buộc thoát: %v", err)
	}
	if !slices.Equal(fa.Log, []string{"quit 20", "force 20"}) {
		t.Errorf("log %v", fa.Log)
	}
}

func TestAppForceQuitGuards(t *testing.T) {
	e, fa := appEngine(false)
	if _, err := e.appStep(t.Context(), appStep(model.AppForceQuit, "chrome", nil)); !errors.Is(err, errDangerNotAllowed) || errCode(err) != model.ErrPermission {
		t.Errorf("chưa xác nhận: %v", err)
	}
	if _, err := e.appStep(t.Context(), appStep(model.AppQuit, "chrome", func(s *model.Step) { s.AppForce = true })); !errors.Is(err, errDangerNotAllowed) {
		t.Errorf("thoát + buộc thoát chưa xác nhận: %v", err)
	}
	e.allowDanger = true
	var bad *invalidErr
	if _, err := e.appStep(t.Context(), appStep(model.AppForceQuit, "explorer", nil)); !errors.As(err, &bad) {
		t.Errorf("app hệ thống: %v, muốn invalid", err)
	}
	if _, err := e.appStep(t.Context(), appStep(model.AppForceQuit, "chrome", nil)); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(fa.Log, []string{"force 10", "force 11"}) {
		t.Errorf("buộc thoát mọi tiến trình: %v", fa.Log)
	}
}

func TestWindowCloseNeedsConsent(t *testing.T) {
	e, fw := winEngine()
	e.allowDanger = false
	if _, err := e.windowStep(t.Context(), winStep(model.WinClose, nil)); !errors.Is(err, errDangerNotAllowed) || len(fw.Log) != 0 {
		t.Errorf("đóng cửa sổ chưa xác nhận: %v %v", err, fw.Log)
	}
}

func TestAppStopWhileWaiting(t *testing.T) {
	e, _ := appEngine(false)
	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := e.appStep(ctx, appStep(model.AppWaitStart, "telegram", func(s *model.Step) { s.TimeoutMs = 10000 }))
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Errorf("dừng giữa chừng: %v sau %v", err, time.Since(start))
	}
}
