package engine

import (
	"context"
	"image"
	"image/color"
	"image/draw"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

// spinFrame – màn 400×300 nền trắng, spinner (ô tối) ở vị trí i (-1 = không có).
func spinFrame(i int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	draw.Draw(img, img.Rect, &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	if i >= 0 {
		draw.Draw(img, image.Rect(100+i*40, 100, 140+i*40, 140), &image.Uniform{C: color.Black}, image.Point{}, draw.Src)
	}
	return img
}

func runWaitScreen(t *testing.T, frames []image.Image, f func(*model.Step)) model.StepResult {
	t.Helper()
	s := modeltest.Step(model.StepWaitScreen, func(s *model.Step) {
		s.IntervalMs, s.TimeoutMs, s.DelayAfter = 1, 2000, 0
		f(s)
	})
	sc := model.Script{Loops: 1, Region: &model.Rect{X: 0, Y: 0, W: 400, H: 300}, Steps: []model.Step{s}}
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil)
	e.scr = &fakeScreen{frames: frames}
	var res model.StepResult
	e.onNode = func(r model.StepResult) { res = r }
	_ = e.Run(context.Background())
	return res
}

func TestWaitScreenStable(t *testing.T) {
	// spinner chạy 4 khung rồi tắt → đứng yên
	res := runWaitScreen(t, []image.Image{spinFrame(0), spinFrame(1), spinFrame(2), spinFrame(3), spinFrame(-1)}, func(s *model.Step) { s.StableMs = 30 })
	if res.Status != model.StatusSuccess {
		t.Fatalf("muốn thành công: %+v", res.Error)
	}
	if res.Output["diff"] != 0.0 || res.Output["ms"].(int) < 30 {
		t.Errorf("output %v", res.Output)
	}
}

func TestWaitScreenNeverStable(t *testing.T) {
	var frames []image.Image
	for i := range 2000 {
		frames = append(frames, spinFrame(i%2))
	}
	res := runWaitScreen(t, frames, func(s *model.Step) { s.StableMs, s.TimeoutMs = 1000, 40 })
	if res.Error == nil || res.Error.Code != model.ErrTimeout {
		t.Fatalf("muốn lỗi timeout: %+v", res)
	}
	if d := res.Error.Detail; d["kind"] != "screen" || d["diff"].(float64) < 1 {
		t.Errorf("detail %v", d)
	}
}

func TestWaitScreenChange(t *testing.T) {
	res := runWaitScreen(t, []image.Image{spinFrame(-1), spinFrame(-1), spinFrame(-1), spinFrame(2)}, func(s *model.Step) { s.WaitMode = model.WaitChange })
	if res.Status != model.StatusSuccess || res.Output["diff"].(float64) <= 1 {
		t.Fatalf("muốn thấy thay đổi: %+v %v", res.Error, res.Output)
	}
	// không đổi → hết giờ, diff 0
	res = runWaitScreen(t, []image.Image{spinFrame(-1)}, func(s *model.Step) { s.WaitMode, s.TimeoutMs = model.WaitChange, 30 })
	if res.Error == nil || res.Error.Code != model.ErrTimeout || res.Error.Detail["diff"] != 0.0 {
		t.Fatalf("muốn timeout, diff 0: %+v", res.Error)
	}
}

func TestWaitScreenCancel(t *testing.T) {
	s := modeltest.Step(model.StepWaitScreen, func(s *model.Step) { s.IntervalMs, s.TimeoutMs, s.StableMs = 5, 60000, 60000 })
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.scr = &fakeScreen{frames: []image.Image{spinFrame(-1)}}
	ctx, cancel := context.WithCancel(context.Background())
	e.onNode = func(model.StepResult) {}
	done := make(chan error)
	go func() { done <- e.Run(ctx) }()
	cancel()
	if err := <-done; err == nil {
		t.Error("huỷ giữa chừng phải dừng")
	}
}
