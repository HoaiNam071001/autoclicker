package engine

import (
	"reflect"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func normStep(f func(s *model.Step)) model.Step {
	return modeltest.Step(model.StepBrowserNorm, func(s *model.Step) {
		s.TimeoutMs = 500
		s.NormZoom, s.NormFullscreen, s.NormPopups = false, false, false // NewStep bật cả ba
		f(s)
	})
}

// Đưa tab lên trước → thoát toàn màn hình, đặt cỡ giữa màn → zoom → Esc; output là khung sau cùng.
func TestBrowserNormalize(t *testing.T) {
	e, f := tabEngine()
	f.Frames = map[int64]model.TabFrame{1: {Box: model.Rect{W: 1920, H: 1080}, Fullscreen: true}}
	res, _ := e.execNode(t.Context(), normStep(func(s *model.Step) {
		s.Tab = "gmail"
		s.NormZoom, s.NormFullscreen, s.NormPopups, s.NormSize = true, true, true, true
		s.WinPreset, s.WinRect = model.PresetCenter, &model.Rect{W: 1000, H: 700}
	}), 1)
	if res.Status != model.StatusSuccess {
		t.Fatalf("%s %+v", res.Status, res.Error)
	}
	want := []string{"activate 1", "frame 1", "normalize [zoom fullscreen popups bounds=460,190 1000x700] 1", "frame 1"}
	if !reflect.DeepEqual(f.Calls, want) {
		t.Errorf("gọi %q, muốn %q", f.Calls, want)
	}
	o := res.Output
	if o["wasFullscreen"] != true || o["zoomReset"] != true || o["w"] != 1000 || o["h"] != 700 || o["x"] != 960 || o["tabId"] != int64(1) {
		t.Errorf("output %v", o)
	}
}

// Chỉ làm mục được bật; không đặt cỡ thì không đọc khung trước.
func TestBrowserNormalizeOnlyZoom(t *testing.T) {
	e, f := tabEngine()
	res, _ := e.execNode(t.Context(), normStep(func(s *model.Step) { s.NormZoom = true }), 1)
	if res.Status != model.StatusSuccess || res.Output["wasFullscreen"] != false {
		t.Fatalf("%s %+v %v", res.Status, res.Error, res.Output)
	}
	if want := []string{"activate 2", "normalize [zoom] 2", "frame 2"}; !reflect.DeepEqual(f.Calls, want) {
		t.Errorf("gọi %q, muốn %q", f.Calls, want)
	}
}

func TestBrowserNormalizeErrors(t *testing.T) {
	e, f := tabEngine()
	res, _ := e.execNode(t.Context(), normStep(func(s *model.Step) { s.NormZoom, s.Tab = true, "không có tab này" }), 1)
	if res.Error == nil || res.Error.Code != "not_found" || len(f.Calls) != 0 {
		t.Errorf("tab không thấy: %+v %q", res.Error, f.Calls)
	}
	res, _ = e.execNode(t.Context(), normStep(func(s *model.Step) { s.NormSize = true }), 1) // preset toạ độ, chưa nhập cỡ
	if res.Error == nil || res.Error.Code != "invalid_input" {
		t.Errorf("chưa nhập cỡ: %+v", res.Error)
	}
}
