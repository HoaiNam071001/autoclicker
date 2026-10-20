package engine

import (
	"context"
	"image"
	"image/color"
	"reflect"
	"strings"
	"sync"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script/scripttest"
	"autoclicker/internal/vision"
	"autoclicker/internal/vision/visiontest"
)

// fakeScreen trả lần lượt các khung hình (khung cuối lặp lại) và màu pixel; ghi lại vùng đã chụp.
type fakeScreen struct {
	mu     sync.Mutex
	frames []image.Image
	colors []color.RGBA
	shots  []model.Rect
	pixels int
}

func (f *fakeScreen) Capture(r model.Rect) (image.Image, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.shots = append(f.shots, r)
	return f.frames[min(len(f.shots)-1, len(f.frames)-1)], nil
}

func (f *fakeScreen) DPIScale(model.Rect) float64 { return 1 }

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
		img := visiontest.Scene(800, 600, 11)
		if btn {
			visiontest.Button(img, 500, 240, 120, 40)
		}
		if scale == 2 {
			return img
		}
		return half(img)
	}
	ref := draw(true)
	tpl, _ := vision.EncodePNG(visiontest.Crop(ref, image.Rect(250*scale, 120*scale, 310*scale, 140*scale)))
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

func imageScript(tpl string, tplScale float64, steps ...model.Step) model.Script {
	wait := modeltest.Step(model.StepWaitImage, func(s *model.Step) { s.Template, s.TemplateScale, s.TimeoutMs, s.IntervalMs = tpl, tplScale, 2000, 1 })
	return model.Script{Loops: 1, Region: &model.Rect{X: 0, Y: 0, W: 400, H: 300}, Steps: append([]model.Step{wait}, steps...)}
}

func TestWaitImageThenClick(t *testing.T) {
	without, tpl := screenWithButton(1, false)
	with, _ := screenWithButton(1, true)
	scr := &fakeScreen{frames: []image.Image{without, without, with}}
	in := &devicetest.FakeInput{}
	sc := imageScript(tpl, 1, modeltest.Step(model.StepClickImage, func(s *model.Step) { s.Dx = 5 }))
	e := newEngine(in, scripttest.Linear(sc), nil)
	e.scr = scr
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(scr.shots) != 3 || scr.shots[0] != (model.Rect{X: 0, Y: 0, W: 400, H: 300}) {
		t.Errorf("chụp %d lần, vùng %v", len(scr.shots), scr.shots)
	}
	if want := []string{"move 285,130", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q, muốn %q (tâm nút 280,130 + dx 5)", in.Lines(), want)
	}
}

// Ảnh mẫu chụp trên Retina (2x), chạy trên màn 1x và ngược lại – vẫn thấy đúng chỗ.
func TestWaitImageAcrossScales(t *testing.T) {
	for _, c := range []struct{ tplScale, screenScale int }{{2, 1}, {1, 2}, {2, 2}} {
		_, tpl := screenWithButton(c.tplScale, true)
		screen, _ := screenWithButton(c.screenScale, true)
		in := &devicetest.FakeInput{}
		e := newEngine(in, scripttest.Linear(imageScript(tpl, float64(c.tplScale), modeltest.Step(model.StepClickImage, nil))), nil)
		e.scr = &fakeScreen{frames: []image.Image{screen}}
		if err := e.Run(context.Background()); err != nil {
			t.Fatalf("mẫu %dx, màn %dx: %v", c.tplScale, c.screenScale, err)
		}
		if got := in.Lines()[0]; got != "move 280,130" && got != "move 279,129" && got != "move 280,129" && got != "move 279,130" {
			t.Errorf("mẫu %dx, màn %dx: %s, muốn ~280,130", c.tplScale, c.screenScale, got)
		}
	}
}

func TestWaitImageTimeout(t *testing.T) {
	without, tpl := screenWithButton(1, false)
	sc := imageScript(tpl, 1)
	sc.Steps[0].TimeoutMs, sc.Steps[0].IntervalMs = 60, 20
	var last Progress
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), func(p Progress) {
		if p.Wait != nil {
			last = p
		}
	})
	e.scr = &fakeScreen{frames: []image.Image{without}}
	err := e.Run(context.Background())
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

// Kiểm tra màu (bước 40): nút cùng hình nhưng xám → không nhận, lời báo nói rõ lệch màu; tắt kiểm tra thì nhận.
func TestWaitImageColorCheck(t *testing.T) {
	blue := visiontest.Scene(400, 300, 11)
	visiontest.Button(blue, 250, 120, 60, 20)
	tpl, _ := vision.EncodePNG(visiontest.Crop(blue, image.Rect(250, 120, 310, 140)))
	grey := visiontest.Crop(blue, blue.Bounds())
	for y := 120; y < 140; y++ {
		for x := 250; x < 310; x++ {
			if grey.RGBAAt(x, y).B == 242 {
				grey.Set(x, y, color.RGBA{105, 105, 105, 255}) // cùng độ sáng với xanh 24,119,242
			}
		}
	}
	sc := imageScript(tpl, 1)
	s := sc.Steps[0]
	s.TimeoutMs, s.IntervalMs, s.ColorTol = 40, 20, 10
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil)
	e.scr = &fakeScreen{frames: []image.Image{grey}}
	_, err := e.findImage(context.Background(), s)
	if err == nil || !strings.Contains(err.Error(), "nhưng màu lệch") {
		t.Fatalf("muốn lỗi lệch màu, có %v", err)
	}
	wantAll(t, "lệch màu", detailJSON(t, err), `"colorTol":10`, `"colorDiff":`)
	s.ColorTol = 0
	if _, err := e.findImage(context.Background(), s); err != nil {
		t.Fatalf("không kiểm tra màu thì nút xám vẫn khớp: %v", err)
	}
}

func TestWaitPixel(t *testing.T) {
	blue := color.RGBA{24, 119, 242, 255}
	scr := &fakeScreen{colors: []color.RGBA{{255, 255, 255, 255}, {30, 115, 240, 255}}}
	sc := model.Script{Loops: 1, Steps: []model.Step{modeltest.Step(model.StepWaitPixel, func(s *model.Step) {
		s.Points[0] = model.Point{X: 10, Y: 20}
		s.Color, s.Tolerance, s.IntervalMs = vision.HexColor(blue), 10, 1
	})}}
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil)
	e.scr = scr
	if err := e.Run(context.Background()); err != nil || scr.pixels != 2 {
		t.Fatalf("muốn thấy ở lần đọc thứ 2: %v, đọc %d lần", err, scr.pixels)
	}
	sc.Steps[0].Tolerance, sc.Steps[0].TimeoutMs, sc.Steps[0].IntervalMs = 2, 30, 10
	e = newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil)
	e.scr = &fakeScreen{colors: []color.RGBA{{30, 115, 240, 255}}}
	if err := e.Run(context.Background()); err == nil || !strings.Contains(err.Error(), "vẫn màu #1E73F0") {
		t.Fatalf("muốn lỗi hết giờ kèm màu hiện tại, có %v", err)
	}
}
