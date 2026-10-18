package engine

import (
	"context"
	"image"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script/scripttest"
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

func TestWaitTextThenClick(t *testing.T) {
	retina := image.NewRGBA(image.Rect(0, 0, 800, 600)) // vùng 400×300 chụp ở 2x
	ocr := &fakeOCR{reads: [][]vision.TextLine{
		{line("Đang tải…", 0, 0, 100, 20)},
		{line("Gửi", 200, 100, 60, 40), line("Gửi", 200, 300, 60, 40)},
	}}
	in := &devicetest.FakeInput{}
	sc := scripttest.TextScript("gửi", modeltest.Step(model.StepClickImage, func(s *model.Step) { s.Dx = 5 }))
	sc.Steps[0].Occurrence = 2
	e := newEngine(in, scripttest.Linear(sc), nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{retina}}, ocr
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	// chỗ thứ 2: px (200, 300) 60×40 → toạ độ 100 + 100 + 15, 50 + 150 + 10
	if want := []string{"move 220,210", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q, muốn %q", in.Lines(), want)
	}
}

func TestWaitTextTimeout(t *testing.T) {
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(scripttest.TextScript("Xem lịch sử")), nil)
	e.sc.Steps[0].TimeoutMs = 5
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}},
		&fakeOCR{reads: [][]vision.TextLine{{line("Xem lich su", 0, 0, 100, 20)}}}
	err := e.Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "không thấy chữ “Xem lịch sử”") || !strings.Contains(err.Error(), "Xem lich su") {
		t.Fatalf("lỗi: %v", err)
	}
	if e.miss == nil || len(e.miss.Read) != 1 {
		t.Errorf("miss: %+v", e.miss)
	}
}
