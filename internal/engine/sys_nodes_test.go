package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/display"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

// ---- Nguồn (S3)

type fakePower struct{ did []string }

func (f *fakePower) Do(_ context.Context, a string) error { f.did = append(f.did, a); return nil }

func runStep(t *testing.T, typ string, patch func(*model.Step), setup func(*Engine)) (*Engine, map[string]any, string) {
	t.Helper()
	s := modeltest.Step(typ, patch)
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	if setup != nil {
		setup(e)
	}
	if err := e.Run(t.Context()); err != nil {
		t.Log("run:", err)
	}
	r := e.results[e.sc.Steps[0].ID]
	if r == nil {
		return e, nil, ""
	}
	code := ""
	if r.Error != nil {
		code = r.Error.Code
	}
	return e, r.Output, code
}

func TestPowerLockRunsWithoutConsentOrCountdown(t *testing.T) {
	f := &fakePower{}
	_, out, code := runStep(t, model.StepPower, func(s *model.Step) { s.PowerOp = model.PowerLock }, func(e *Engine) { e.powers = f })
	if code != "" || len(f.did) != 1 || f.did[0] != "lock" || out["action"] != "lock" {
		t.Errorf("did=%v out=%v code=%q", f.did, out, code)
	}
}

func TestPowerShutdownNeedsConsent(t *testing.T) {
	f := &fakePower{}
	_, _, code := runStep(t, model.StepPower, func(s *model.Step) { s.PowerOp = model.PowerShutdown }, func(e *Engine) { e.powers = f })
	if code != model.ErrPermission || len(f.did) != 0 {
		t.Errorf("did=%v code=%q", f.did, code)
	}
}

func TestPowerCountdownCanBeCancelled(t *testing.T) {
	f := &fakePower{}
	s := modeltest.Step(model.StepPower, func(s *model.Step) { s.PowerOp, s.CountdownSec = model.PowerRestart, 30 })
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.powers, e.allowDanger = f, true
	ctx, cancel := context.WithCancel(t.Context())
	var sawCountdown bool
	e.onProgress = func(p Progress) {
		if p.Wait != nil && p.Wait.Kind == "power" {
			sawCountdown = true
			cancel() // người dùng bấm Dừng trong lúc đếm
		}
	}
	done := make(chan error, 1)
	go func() { done <- e.Run(ctx) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("đếm ngược không dừng khi huỷ")
	}
	if !sawCountdown || len(f.did) != 0 {
		t.Errorf("sawCountdown=%v did=%v – huỷ lúc đếm phải KHÔNG tắt máy", sawCountdown, f.did)
	}
}

func TestPowerShutdownAfterCountdown(t *testing.T) {
	f := &fakePower{}
	_, out, code := runStep(t, model.StepPower, func(s *model.Step) { s.PowerOp, s.CountdownSec = model.PowerSleep, 1 }, func(e *Engine) { e.powers, e.allowDanger = f, true })
	if code != "" || len(f.did) != 1 || out["countdownSec"] != 1 {
		t.Errorf("did=%v out=%v code=%q", f.did, out, code)
	}
}

// ---- Màn hình / âm thanh (S7)

type fakeDisplay struct {
	dark       bool
	brightness int
	noBright   bool
	devices    []display.Audio
	played     []string
	waited     bool
	setAudioID string
}

func (f *fakeDisplay) Dark() (bool, error)  { return f.dark, nil }
func (f *fakeDisplay) SetDark(d bool) error { f.dark = d; return nil }
func (f *fakeDisplay) Brightness() (int, error) {
	if f.noBright {
		return 0, display.ErrNoBrightness
	}
	return f.brightness, nil
}
func (f *fakeDisplay) SetBrightness(p int) error {
	if f.noBright {
		return display.ErrNoBrightness
	}
	f.brightness = p
	return nil
}
func (f *fakeDisplay) AudioDevices() ([]display.Audio, error) { return f.devices, nil }
func (f *fakeDisplay) SetDefaultAudio(id string, input bool) error {
	f.setAudioID = id
	return nil
}
func (f *fakeDisplay) PlaySound(_ context.Context, target string, wait bool) error {
	f.played, f.waited = append(f.played, target), wait
	return nil
}

func runDisplay(t *testing.T, f *fakeDisplay, patch func(*model.Step)) (map[string]any, string) {
	t.Helper()
	_, out, code := runStep(t, model.StepDisplay, patch, func(e *Engine) { e.disp = f })
	return out, code
}

func TestDisplayTheme(t *testing.T) {
	f := &fakeDisplay{}
	out, code := runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayTheme, "toggle" })
	if code != "" || !f.dark || out["theme"] != "dark" || out["changed"] != true {
		t.Errorf("out=%v code=%q dark=%v", out, code, f.dark)
	}
	out, _ = runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayTheme, "dark" })
	if out["changed"] != false {
		t.Errorf("đã tối rồi mà changed=%v", out["changed"])
	}
	if _, code := runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayTheme, "{{vars.x}}" }); code == "" {
		t.Error("biến chưa có phải lỗi")
	}
}

func TestDisplayBrightness(t *testing.T) {
	f := &fakeDisplay{brightness: 30}
	out, code := runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayBrightness, "80" })
	if code != "" || f.brightness != 80 || out["previous"] != 30 {
		t.Errorf("out=%v code=%q", out, code)
	}
	f = &fakeDisplay{noBright: true}
	if _, code := runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayBrightness, "50" }); code != model.ErrInvalidInput {
		t.Errorf("màn ngoài không chỉnh được phải invalid_input, được %q", code)
	}
}

func TestDisplayAudioSwitchAndNotFound(t *testing.T) {
	f := &fakeDisplay{devices: []display.Audio{
		{ID: "spk", Name: "MacBook Speakers", Default: true}, {ID: "air", Name: "AirPods Pro"}, {ID: "mic", Name: "AirPods Pro", Input: true},
	}}
	out, code := runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayAudioOut, "airpods" })
	if code != "" || f.setAudioID != "air" || out["previous"] != "MacBook Speakers" {
		t.Errorf("out=%v code=%q id=%q", out, code, f.setAudioID)
	}
	_, code = runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue = model.DisplayAudioOut, "sony" })
	if code != model.ErrNotFound {
		t.Errorf("không có thiết bị phải not_found, được %q", code)
	}
}

func TestDisplayReadAndSound(t *testing.T) {
	f := &fakeDisplay{dark: true, brightness: 55, devices: []display.Audio{{ID: "spk", Name: "Loa", Default: true}, {ID: "mic", Name: "Mic", Input: true, Default: true}}}
	out, code := runDisplay(t, f, nil)
	items, _ := out["items"].([]any)
	if code != "" || out["theme"] != "dark" || out["brightness"] != 55 || out["outputName"] != "Loa" || out["inputName"] != "Mic" || len(items) != 2 {
		t.Errorf("out=%v code=%q", out, code)
	}
	out, code = runDisplay(t, f, func(s *model.Step) { s.DisplayOp, s.DisplayValue, s.SoundWait = model.DisplaySound, "Glass", true })
	if code != "" || len(f.played) != 1 || f.played[0] != "Glass" || !f.waited || !strings.Contains(strings.Join(f.played, ""), "Glass") {
		t.Errorf("played=%v waited=%v out=%v", f.played, f.waited, out)
	}
}

// ---- Bộ gõ: toggleIME (S6b)

func toggleIME(t *testing.T, running bool, patch func(*model.Step)) ([]string, string) {
	t.Helper()
	in := &devicetest.FakeInput{}
	_, _, code := runStep(t, model.StepInputSource, func(s *model.Step) {
		s.InputOp = model.InputToggleIME
		if patch != nil {
			patch(s)
		}
	}, func(e *Engine) {
		e.in, e.inputs = in, newFakeInputs()
		e.apps = &devicetest.FakeApps{}
		if running {
			e.apps = &devicetest.FakeApps{Apps: []model.AppInfo{{Name: "UniKey", PID: 7}}}
		}
	})
	return in.Lines(), code
}

func TestToggleIMEPressesDefaultHotkeyOnce(t *testing.T) {
	lines, code := toggleIME(t, true, nil)
	got := strings.Join(lines, "|")
	if code != "" || strings.Count(got, "down ctrl") != 1 || strings.Count(got, "down shift") != 1 || !strings.Contains(got, "up shift") {
		t.Errorf("lines=%v code=%q", lines, code)
	}
}

func TestToggleIMECustomHotkeyWithKey(t *testing.T) {
	lines, code := toggleIME(t, true, func(s *model.Step) { s.InputHotkey = "alt+z" })
	if got := strings.Join(lines, "|"); code != "" || !strings.Contains(got, "down z [alt]") {
		t.Errorf("lines=%v code=%q", lines, code)
	}
}

func TestToggleIMENeedsRunningApp(t *testing.T) {
	lines, code := toggleIME(t, false, nil)
	if code != model.ErrNotFound || len(lines) != 0 {
		t.Errorf("lines=%v code=%q – không có phần mềm gõ thì không được bấm phím", lines, code)
	}
	if _, code := toggleIME(t, true, func(s *model.Step) { s.InputHotkey = "a+b" }); code != model.ErrInvalidInput {
		t.Errorf("phím tắt sai phải invalid_input, được %q", code)
	}
}
