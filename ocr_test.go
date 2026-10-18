package main

import (
	"context"
	"image"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script"
	"autoclicker/internal/vision"
)

func line(text string, x, y, w, h int) vision.TextLine {
	return vision.TextLine{Text: text, Box: model.Rect{X: x, Y: y, W: w, H: h}}
}

// fakeOCR trả lần lượt các kết quả đọc (cuối lặp lại).
type fakeOCR struct {
	reads [][]vision.TextLine
	n     int
}

func (f *fakeOCR) ReadText(image.Image) ([]vision.TextLine, error) {
	f.n++
	return f.reads[min(f.n-1, len(f.reads)-1)], nil
}

func textScript(text string, steps ...model.Step) model.Script {
	wait := modeltest.Step(model.StepWaitImage, func(s *model.Step) {
		s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = model.FindByText, text, 2000, 1
	})
	return model.Script{Loops: 1, Region: &model.Rect{X: 100, Y: 50, W: 400, H: 300}, Steps: append([]model.Step{wait}, steps...)}
}

func TestWaitTextThenClick(t *testing.T) {
	retina := image.NewRGBA(image.Rect(0, 0, 800, 600)) // vùng 400×300 chụp ở 2x
	ocr := &fakeOCR{reads: [][]vision.TextLine{
		{line("Đang tải…", 0, 0, 100, 20)},
		{line("Gửi", 200, 100, 60, 40), line("Gửi", 200, 300, 60, 40)},
	}}
	in := &devicetest.FakeInput{}
	sc := textScript("gửi", modeltest.Step(model.StepClickImage, func(s *model.Step) { s.Dx = 5 }))
	sc.Steps[0].Occurrence = 2
	e := newEngine(in, linear(sc), nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{retina}}, ocr
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	// chỗ thứ 2: px (200, 300) 60×40 → toạ độ 100 + 100 + 15, 50 + 150 + 10
	if want := []string{"move 220,210", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q, muốn %q", in.Lines(), want)
	}
}

func TestWaitTextTimeout(t *testing.T) {
	e := newEngine(&devicetest.FakeInput{}, linear(textScript("Xem lịch sử")), nil)
	e.sc.Steps[0].TimeoutMs = 5
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}},
		&fakeOCR{reads: [][]vision.TextLine{{line("Xem lich su", 0, 0, 100, 20)}}}
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "không thấy chữ “Xem lịch sử”") || !strings.Contains(err.Error(), "Xem lich su") {
		t.Fatalf("lỗi: %v", err)
	}
	if e.miss == nil || len(e.miss.Read) != 1 {
		t.Errorf("miss: %+v", e.miss)
	}
}

func TestValidateTextStep(t *testing.T) {
	sc := textScript("  ")
	if err := script.ValidateScript(linear(sc), model.Rect{W: 1920, H: 1080}); err == nil || !strings.Contains(err.Error(), "chưa nhập chữ") {
		t.Errorf("chữ trống: %v", err)
	}
	sc = textScript("Gửi cho {{ten}}")
	sc.Inputs = []model.ScriptInput{{Key: "ten", Label: "Tên", Type: model.InputText}}
	out, err := script.ResolveInputs(sc, map[string]any{"ten": "Lan"})
	if err != nil || out.Steps[0].FindText != "Gửi cho Lan" {
		t.Errorf("biến trong chữ cần tìm: %v %q", err, out.Steps[0].FindText)
	}
}
