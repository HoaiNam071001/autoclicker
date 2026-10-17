package main

import (
	"context"
	"image"
	"image/color"
	"reflect"
	"strings"
	"sync"
	"testing"
)

// fakeScreen trả lần lượt các khung hình (khung cuối lặp lại) và màu pixel; ghi lại vùng đã chụp.
type fakeScreen struct {
	mu     sync.Mutex
	frames []image.Image
	colors []color.RGBA
	shots  []Rect
	pixels int
}

func (f *fakeScreen) Capture(r Rect) (image.Image, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.shots = append(f.shots, r)
	return f.frames[min(len(f.shots)-1, len(f.frames)-1)], nil
}

func (f *fakeScreen) Pixel(x, y int) (color.RGBA, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pixels++
	return f.colors[min(f.pixels-1, len(f.colors)-1)], nil
}

// Màn hình 400×300 có nút tại (250, 120) cỡ 60×20; ảnh mẫu = nút. Vẽ ở Retina (2x); scale 1 = thu nhỏ một nửa
// (như cùng giao diện trên màn 1x).
func screenWithButton(scale int, withButton bool) (*image.RGBA, string) {
	draw := func(btn bool) *image.RGBA {
		img := scene(800, 600, 11)
		if btn {
			button(img, 500, 240, 120, 40)
		}
		if scale == 2 {
			return img
		}
		return half(img)
	}
	ref := draw(true)
	tpl, _ := encodePNG(crop(ref, image.Rect(250*scale, 120*scale, 310*scale, 140*scale)))
	return draw(withButton), tpl
}

// half thu nhỏ một nửa (trung bình 2×2).
func half(img *image.RGBA) *image.RGBA {
	b := img.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx()/2, b.Dy()/2))
	for y := 0; y < b.Dy()/2; y++ {
		for x := 0; x < b.Dx()/2; x++ {
			var s [3]int
			for _, d := range [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}} {
				c := img.RGBAAt(2*x+d[0], 2*y+d[1])
				s[0], s[1], s[2] = s[0]+int(c.R), s[1]+int(c.G), s[2]+int(c.B)
			}
			out.SetRGBA(x, y, color.RGBA{uint8(s[0] / 4), uint8(s[1] / 4), uint8(s[2] / 4), 255})
		}
	}
	return out
}

func imageScript(tpl string, tplScale float64, steps ...Step) Script {
	wait := step(StepWaitImage, func(s *Step) { s.Template, s.TemplateScale, s.TimeoutMs, s.IntervalMs = tpl, tplScale, 2000, 1 })
	return Script{Loops: 1, Region: &Rect{X: 0, Y: 0, W: 400, H: 300}, Steps: append([]Step{wait}, steps...)}
}

func TestWaitImageThenClick(t *testing.T) {
	without, tpl := screenWithButton(1, false)
	with, _ := screenWithButton(1, true)
	scr := &fakeScreen{frames: []image.Image{without, without, with}}
	in := &fakeInput{}
	sc := imageScript(tpl, 1, step(StepClickImage, func(s *Step) { s.Dx = 5 }))
	e := newEngine(in, linear(sc), nil)
	e.scr = scr
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(scr.shots) != 3 || scr.shots[0] != (Rect{0, 0, 400, 300}) {
		t.Errorf("chụp %d lần, vùng %v", len(scr.shots), scr.shots)
	}
	if want := []string{"move 285,130", "click left false"}; !reflect.DeepEqual(in.lines(), want) {
		t.Errorf("có %q, muốn %q (tâm nút 280,130 + dx 5)", in.lines(), want)
	}
}

// Ảnh mẫu chụp trên Retina (2x), chạy trên màn 1x và ngược lại – vẫn thấy đúng chỗ.
func TestWaitImageAcrossScales(t *testing.T) {
	for _, c := range []struct{ tplScale, screenScale int }{{2, 1}, {1, 2}, {2, 2}} {
		_, tpl := screenWithButton(c.tplScale, true)
		screen, _ := screenWithButton(c.screenScale, true)
		in := &fakeInput{}
		e := newEngine(in, linear(imageScript(tpl, float64(c.tplScale), step(StepClickImage, nil))), nil)
		e.scr = &fakeScreen{frames: []image.Image{screen}}
		if err := e.run(context.Background()); err != nil {
			t.Fatalf("mẫu %dx, màn %dx: %v", c.tplScale, c.screenScale, err)
		}
		if got := in.lines()[0]; got != "move 280,130" && got != "move 279,129" && got != "move 280,129" && got != "move 279,130" {
			t.Errorf("mẫu %dx, màn %dx: %s, muốn ~280,130", c.tplScale, c.screenScale, got)
		}
	}
}

func TestWaitImageTimeout(t *testing.T) {
	without, tpl := screenWithButton(1, false)
	sc := imageScript(tpl, 1)
	sc.Steps[0].TimeoutMs, sc.Steps[0].IntervalMs = 60, 20
	var last Progress
	e := newEngine(&fakeInput{}, linear(sc), func(p Progress) {
		if p.Wait != nil {
			last = p
		}
	})
	e.scr = &fakeScreen{frames: []image.Image{without}}
	err := e.run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "không thấy ảnh mẫu sau 0.1 s") && !strings.Contains(err.Error(), "không thấy ảnh mẫu sau") {
		t.Fatalf("muốn lỗi hết giờ, có %v", err)
	}
	if e.failed != sc.Steps[0].ID || e.miss == nil || e.miss.W != 60 || e.miss.H != 20 {
		t.Errorf("bước lỗi %q, chỗ giống nhất %+v", e.failed, e.miss)
	}
	if last.Wait == nil || last.Wait.Kind != "image" || last.Wait.TimeoutMs != 60 {
		t.Errorf("tiến độ chờ: %+v", last.Wait)
	}
}

func TestWaitPixel(t *testing.T) {
	blue := color.RGBA{24, 119, 242, 255}
	scr := &fakeScreen{colors: []color.RGBA{{255, 255, 255, 255}, {30, 115, 240, 255}}}
	sc := Script{Loops: 1, Steps: []Step{step(StepWaitPixel, func(s *Step) {
		s.Points[0] = Point{X: 10, Y: 20}
		s.Color, s.Tolerance, s.IntervalMs = hexColor(blue), 10, 1
	})}}
	e := newEngine(&fakeInput{}, linear(sc), nil)
	e.scr = scr
	if err := e.run(context.Background()); err != nil || scr.pixels != 2 {
		t.Fatalf("muốn thấy ở lần đọc thứ 2: %v, đọc %d lần", err, scr.pixels)
	}
	sc.Steps[0].Tolerance, sc.Steps[0].TimeoutMs, sc.Steps[0].IntervalMs = 2, 30, 10
	e = newEngine(&fakeInput{}, linear(sc), nil)
	e.scr = &fakeScreen{colors: []color.RGBA{{30, 115, 240, 255}}}
	if err := e.run(context.Background()); err == nil || !strings.Contains(err.Error(), "vẫn màu #1E73F0") {
		t.Fatalf("muốn lỗi hết giờ kèm màu hiện tại, có %v", err)
	}
}

func TestValidateImageSteps(t *testing.T) {
	_, tpl := screenWithButton(1, true)
	cases := map[string][]Step{
		"click theo bước trước nhưng không có chờ ảnh": {step(StepClickImage, nil)},
		"chờ ảnh chưa có ảnh mẫu":                      {step(StepWaitImage, nil)},
		"ảnh mẫu hỏng":                                 {step(StepWaitImage, func(s *Step) { s.Template = "abc" })},
		"độ khớp 0":                                    {step(StepWaitImage, func(s *Step) { s.Template, s.Threshold = tpl, 0 })},
		"màu sai":                                      {step(StepWaitPixel, func(s *Step) { s.Color = "xanh" })},
	}
	for name, steps := range cases {
		if err := validateScript(linear(Script{Steps: steps}), Rect{}); err == nil {
			t.Errorf("%s: muốn lỗi", name)
		}
	}
	ok := []Step{step(StepWaitImage, func(s *Step) { s.Template = tpl }), step(StepClickImage, nil),
		step(StepDrag, func(s *Step) { s.DragFromPrev = true })}
	if err := validateScript(linear(Script{Steps: ok}), Rect{}); err != nil {
		t.Errorf("chờ ảnh → click theo ảnh → kéo từ điểm trước: %v", err)
	}
}
