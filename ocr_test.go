package main

import (
	"context"
	"image"
	"reflect"
	"strings"
	"testing"
)

func line(text string, x, y, w, h int, words ...TextWord) TextLine {
	return TextLine{Text: text, Box: Rect{x, y, w, h}, Words: words}
}

func TestFindText(t *testing.T) {
	lines := []TextLine{
		line("Trang chủ", 10, 100, 80, 20),
		line("Bạn có thể Xem lịch sử ở đây", 10, 10, 280, 20,
			TextWord{0, 3, Rect{10, 10, 30, 20}}, TextWord{11, 3, Rect{110, 10, 35, 20}}, TextWord{15, 4, Rect{150, 10, 40, 20}},
			TextWord{20, 3, Rect{195, 10, 30, 20}}),
		line("Xem lịch sử", 300, 100, 110, 20),
		line("XEMLICH SU", 10, 200, 100, 20),
	}
	cases := []struct {
		q    TextQuery
		want []Rect
	}{
		// thứ tự đọc: dòng trên trước; trong dòng chỉ lấy khung các từ của chuỗi
		{TextQuery{Text: "xem lịch sử"}, []Rect{{110, 10, 115, 20}, {300, 100, 110, 20}}},
		{TextQuery{Text: "Xem lịch sử", Exact: true}, []Rect{{300, 100, 110, 20}}},
		{TextQuery{Text: "xem lich su", NoAccent: true}, []Rect{{110, 10, 115, 20}, {300, 100, 110, 20}, {10, 200, 100, 20}}},
		{TextQuery{Text: "trang  chu"}, nil},
		{TextQuery{Text: "Trang chủ"}, []Rect{{10, 100, 80, 20}}},
		{TextQuery{Text: "   "}, nil},
	}
	for _, c := range cases {
		var got []Rect
		for _, h := range findText(lines, c.q) {
			got = append(got, h.Box)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v: được %v, cần %v", c.q, got, c.want)
		}
	}
	// không có khung từ → chia theo số ký tự
	if b := subBox(line("abcdefghij", 0, 0, 100, 10), 2, 4); b != (Rect{20, 0, 20, 10}) {
		t.Errorf("chia theo tỷ lệ: %v", b)
	}
}

// fakeOCR trả lần lượt các kết quả đọc (cuối lặp lại).
type fakeOCR struct {
	reads [][]TextLine
	n     int
}

func (f *fakeOCR) ReadText(image.Image) ([]TextLine, error) {
	f.n++
	return f.reads[min(f.n-1, len(f.reads)-1)], nil
}

func textScript(text string, steps ...Step) Script {
	wait := step(StepWaitImage, func(s *Step) {
		s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = FindByText, text, 2000, 1
	})
	return Script{Loops: 1, Region: &Rect{X: 100, Y: 50, W: 400, H: 300}, Steps: append([]Step{wait}, steps...)}
}

func TestWaitTextThenClick(t *testing.T) {
	retina := image.NewRGBA(image.Rect(0, 0, 800, 600)) // vùng 400×300 chụp ở 2x
	ocr := &fakeOCR{reads: [][]TextLine{
		{line("Đang tải…", 0, 0, 100, 20)},
		{line("Gửi", 200, 100, 60, 40), line("Gửi", 200, 300, 60, 40)},
	}}
	in := &fakeInput{}
	sc := textScript("gửi", step(StepClickImage, func(s *Step) { s.Dx = 5 }))
	sc.Steps[0].Occurrence = 2
	e := newEngine(in, sc, nil)
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{retina}}, ocr
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	// chỗ thứ 2: px (200, 300) 60×40 → toạ độ 100 + 100 + 15, 50 + 150 + 10
	if want := []string{"move 220,210", "click left false"}; !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("có %q, muốn %q", in.lines(), want)
	}
}

func TestWaitTextTimeout(t *testing.T) {
	e := newEngine(&fakeInput{}, textScript("Xem lịch sử"), nil)
	e.sc.Steps[0].TimeoutMs = 5
	e.scr, e.ocr = &fakeScreen{frames: []image.Image{image.NewRGBA(image.Rect(0, 0, 400, 300))}},
		&fakeOCR{reads: [][]TextLine{{line("Xem lich su", 0, 0, 100, 20)}}}
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
	if err := validateScript(sc, Rect{W: 1920, H: 1080}); err == nil || !strings.Contains(err.Error(), "chưa nhập chữ") {
		t.Errorf("chữ trống: %v", err)
	}
	sc = textScript("Gửi cho {{ten}}")
	sc.Inputs = []ScriptInput{{Key: "ten", Label: "Tên", Type: InputText}}
	out, err := resolveInputs(sc, map[string]any{"ten": "Lan"})
	if err != nil || out.Steps[0].FindText != "Gửi cho Lan" {
		t.Errorf("biến trong chữ cần tìm: %v %q", err, out.Steps[0].FindText)
	}
}
