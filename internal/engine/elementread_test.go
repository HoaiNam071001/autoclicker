package engine

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/element"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script/scripttest"
	"autoclicker/internal/vision"
)

// words – dòng OCR có từng từ (px ảnh): mỗi từ một khung.
func words(ws []string, boxes ...model.Rect) vision.TextLine {
	l, pos := vision.TextLine{}, 0
	for i, w := range ws {
		if i > 0 {
			l.Text += " "
			pos++
		}
		l.Text += w
		l.Words = append(l.Words, vision.TextWord{Start: pos, Len: len([]rune(w)), Box: boxes[i]})
		pos += len([]rune(w))
	}
	l.Box = model.Rect{X: boxes[0].X, Y: boxes[0].Y, W: boxes[len(boxes)-1].X + boxes[len(boxes)-1].W - boxes[0].X, H: boxes[0].H}
	return l
}

// frame – ảnh chụp giả cỡ w×h, điểm (0, 0) mang dấu mark (ảnh khác dấu = điểm ảnh khác).
func frame(w, h int, mark uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(0, 0, color.RGBA{R: mark, A: 255})
	return img
}

// Nhãn là phần tử (có trong cây), giá trị vẽ bằng canvas → OCR vùng bên phải nhãn, chuẩn hoá số.
func TestElementReadNearLabel(t *testing.T) {
	nodes := []element.Node{
		{Role: "window", Box: elemBox(0, 0, 1000, 800)},
		{Role: "text", Raw: "AXStaticText", Value: "Giá", Box: elemBox(100, 200, 40, 20), Depth: 1},
	}
	read := modeltest.Step(model.StepElementRead, func(s *model.Step) {
		s.Selectors, s.IntervalMs, s.ReadNear, s.ReadNumber = []string{`text="Giá"`}, 1, model.ReadNearRight, true
	})
	in := &devicetest.FakeInput{}
	e := newEngine(in, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{read}}), nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{nodes}}
	scr := &fakeScreen{frames: []image.Image{frame(540, 60, 0)}} // vùng NearArea: X 100 Y 180, 540 × 60
	e.scr, e.ocr = scr, &fakeOCR{reads: [][]vision.TextLine{{words([]string{"Giá", "1.234.000₫"}, elemBox(0, 20, 40, 20), elemBox(50, 20, 100, 20))}}}
	if err := e.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	if scr.shots[0] != elemBox(100, 180, 540, 60) {
		t.Errorf("chụp vùng %+v", scr.shots[0])
	}
	out := e.results[e.sc.Steps[0].ID].Output
	if out["text"] != "1234000" || out["raw"] != "1.234.000₫" || out["x"] != 200 || out["role"] != "text" {
		t.Errorf("output %v", out)
	}
}

// Nhãn tìm bằng OCR, giá trị ở dòng dưới; lọc mẫu lấy nhóm.
func TestElementReadTextLabelBelow(t *testing.T) {
	read := modeltest.Step(model.StepElementRead, func(s *model.Step) {
		s.FindBy, s.FindText, s.ReadNear, s.ReadPattern, s.TimeoutMs, s.IntervalMs = model.FindByText, "mã đơn", model.ReadNearBelow, `#(\w+)`, 2000, 1
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Region: &model.Rect{X: 0, Y: 0, W: 400, H: 300}, Steps: []model.Step{read}}), nil)
	e.scr = &fakeScreen{frames: []image.Image{frame(400, 300, 0)}}
	e.ocr = &fakeOCR{reads: [][]vision.TextLine{{line("Mã đơn", 10, 10, 80, 20), line("Đơn #AB12 – đang giao", 10, 40, 200, 20)}}}
	if err := e.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	if out := e.results[e.sc.Steps[0].ID].Output; out["text"] != "AB12" || out["count"] != 1 {
		t.Errorf("output %v", out)
	}
}

// OCR Windows (dòng Approx) đọc sai dấu nhãn – đổi toạ độ ở readArea phải giữ Approx để tìm nhãn so gần đúng.
func TestElementReadApproxLabel(t *testing.T) {
	read := modeltest.Step(model.StepElementRead, func(s *model.Step) {
		s.FindBy, s.FindText, s.ReadNear, s.ReadNumber, s.TimeoutMs, s.IntervalMs = model.FindByText, "Giá hiện tại", model.ReadNearRight, true, 500, 1
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Region: &model.Rect{X: 0, Y: 0, W: 400, H: 300}, Steps: []model.Step{read}}), nil)
	e.scr = &fakeScreen{frames: []image.Image{frame(400, 300, 0)}}
	label := words([]string{"Giå", "hién", "tai:", "1.329.000d"}, elemBox(10, 20, 30, 20), elemBox(45, 20, 35, 20), elemBox(85, 20, 30, 20), elemBox(130, 20, 90, 20))
	label.Approx = true
	e.ocr = &fakeOCR{reads: [][]vision.TextLine{{label}}}
	if err := e.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	if out := e.results[e.sc.Steps[0].ID].Output; out["text"] != "1329000" {
		t.Errorf("output %v", out)
	}
}

// Chờ giá trị đổi: lần đọc đầu làm mốc; ảnh y hệt thì không OCR lại.
func TestElementReadWaitsForChange(t *testing.T) {
	read := modeltest.Step(model.StepElementRead, func(s *model.Step) {
		s.FindBy, s.FindText, s.ReadNear, s.ReadNumber, s.ReadChange, s.TimeoutMs, s.IntervalMs = model.FindByText, "Số dư", model.ReadNearRight, true, true, 2000, 1
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Region: &model.Rect{X: 0, Y: 0, W: 400, H: 300}, Steps: []model.Step{read}}), nil)
	e.scr = &fakeScreen{frames: []image.Image{frame(400, 300, 0), frame(400, 300, 0), frame(400, 300, 0), frame(400, 300, 1)}}
	ocr := &fakeOCR{reads: [][]vision.TextLine{
		{line("Số dư", 10, 10, 50, 20), line("100.000 đ", 70, 10, 80, 20)},
		{line("Số dư", 10, 10, 50, 20), line("120.000 đ", 70, 10, 80, 20)},
	}}
	e.ocr = ocr
	if err := e.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	if out := e.results[e.sc.Steps[0].ID].Output; out["text"] != "120000" || out["old"] != "100000" {
		t.Errorf("output %v", out)
	}
	if ocr.n != 2 {
		t.Errorf("OCR %d lần, muốn 2 (ảnh không đổi thì dùng lại)", ocr.n)
	}
}

// Đọc được chữ nhưng không khớp mẫu → chờ tới hết giờ, báo chữ đã đọc.
func TestElementReadPatternMiss(t *testing.T) {
	nodes := append(page(elemBox(100, 200, 80, 30)), element.Node{Role: "text", Value: "Đang tải…", Box: elemBox(100, 300, 80, 20), Depth: 1})
	read := modeltest.Step(model.StepElementRead, func(s *model.Step) {
		s.Selectors, s.ReadNumber, s.TimeoutMs, s.IntervalMs = []string{"role=text"}, true, 30, 1
	})
	_, _, err := elementRun(t, &devicetest.FakeElements{Frames: [][]element.Node{nodes}}, read)
	if err == nil || !strings.Contains(err.Error(), "đọc được “Đang tải…”") {
		t.Errorf("lỗi %v", err)
	}
}

// Thử tìm của Đọc phần tử: tìm theo cây phần tử dù findBy trống.
func TestFindOnceElementRead(t *testing.T) {
	read := modeltest.Step(model.StepElementRead, func(s *model.Step) { s.Selectors = []string{"text=Mua"} })
	e := newEngine(nil, model.Script{Loops: 1}, nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{page(elemBox(100, 200, 80, 30))}}
	res, err := e.FindOnce(read)
	if err != nil || !res.Found || res.X != 100 {
		t.Errorf("%+v %v", res, err)
	}
}
