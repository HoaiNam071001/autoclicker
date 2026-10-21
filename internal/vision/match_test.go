package vision

import (
	"image"
	"image/color"
	"math/rand/v2"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/vision/visiontest"
)

// findTemplate – chỗ tốt nhất, không mặt nạ (test cũ của bước 25).
func findTemplate(s, t Gray) Match { return findPeaks(s, t, Gray{}, candidates)[0] }

func tplOf(g Gray) Template { return Template{Gray: g} }

func TestFindTemplateExact(t *testing.T) {
	img := visiontest.Scene(1200, 800, 1)
	visiontest.Button(img, 700, 420, 120, 40)
	tpl := visiontest.Crop(img, image.Rect(700, 420, 820, 460))
	start := time.Now()
	m := findTemplate(ToGray(img), ToGray(tpl))
	if m.X != 700 || m.Y != 420 || m.Score < 0.99 {
		t.Fatalf("muốn (700, 420) ~1.0, có %+v", m)
	}
	t.Logf("1200×800, mẫu 120×40: %v", time.Since(start))
}

// Mẫu chụp ở Retina (2x) tìm trên màn 1x: thu mẫu về cùng tỷ lệ trước khi tìm (như engine làm).
func TestFindTemplateScaled(t *testing.T) {
	big := visiontest.Scene(1600, 1000, 2)
	visiontest.Button(big, 900, 600, 240, 80)
	tpl2x := ToGray(visiontest.Crop(big, image.Rect(900, 600, 1140, 680)))
	screen1x := ToGray(big).resize(800, 500)
	m := findTemplate(screen1x, tpl2x.resize(tpl2x.w/2, tpl2x.h/2))
	if abs(m.X-450) > 1 || abs(m.Y-300) > 1 || m.Score < 0.9 {
		t.Fatalf("muốn ~(450, 300), có %+v", m)
	}
}

// Không có trên màn hình → độ khớp thấp; có nhiễu / đổi độ sáng nhẹ vẫn khớp cao.
func TestFindTemplateAbsentAndNoise(t *testing.T) {
	img := visiontest.Scene(900, 600, 3)
	other := visiontest.Scene(900, 600, 4)
	visiontest.Button(other, 100, 100, 120, 40)
	tpl := ToGray(visiontest.Crop(other, image.Rect(100, 100, 220, 140)))
	if m := findTemplate(ToGray(img), tpl); m.Score > 0.6 {
		t.Errorf("không có nút mà khớp %.2f tại (%d, %d)", m.Score, m.X, m.Y)
	}
	visiontest.Button(img, 500, 300, 120, 40)
	g := ToGray(img)
	r := rand.New(rand.NewPCG(9, 9))
	for i := range g.p {
		g.p[i] = g.p[i]*0.9 + 10 + float32(r.IntN(9)) - 4 // tối đi 10%, nhiễu ±4
	}
	if m := findTemplate(g, tpl); m.X != 500 || m.Y != 300 || m.Score < 0.95 {
		t.Errorf("nhiễu + đổi sáng: muốn (500, 300) ≥ 0.95, có %+v", m)
	}
}

func TestFindTemplateFlatAndSmall(t *testing.T) {
	img := visiontest.Scene(600, 400, 5)
	for y := 200; y < 230; y++ {
		for x := 300; x < 360; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	flat := ToGray(visiontest.Crop(img, image.Rect(305, 205, 355, 225))) // mẫu một màu
	if m := findTemplate(ToGray(img), flat); m.X < 300 || m.X > 310 || m.Y < 200 || m.Y > 210 || m.Score < 0.9 {
		t.Errorf("mẫu một màu: %+v", m)
	}
	visiontest.Button(img, 50, 50, 14, 12) // mẫu rất nhỏ
	small := ToGray(visiontest.Crop(img, image.Rect(50, 50, 64, 62)))
	if m := findTemplate(ToGray(img), small); m.X != 50 || m.Y != 50 {
		t.Errorf("mẫu nhỏ: %+v", m)
	}
}

// Cỡ màn Retina 14" (3024×1964 px), nút 240×80 px.
func BenchmarkFindTemplateRetina(b *testing.B) {
	img := visiontest.Scene(3024, 1964, 7)
	visiontest.Button(img, 2000, 1500, 240, 80)
	tpl := ToGray(visiontest.Crop(img, image.Rect(2000, 1500, 2240, 1580)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := findTemplate(ToGray(img), tpl)
		if m.X != 2000 || m.Y != 1500 {
			b.Fatalf("%+v", m)
		}
	}
}

// Nút nhỏ 40×30 px trên màn Retina – trường hợp chậm nhất (ít thu nhỏ được).
func BenchmarkFindTemplateRetinaSmall(b *testing.B) {
	img := visiontest.Scene(3024, 1964, 8)
	visiontest.Button(img, 1200, 900, 40, 30)
	tpl := ToGray(visiontest.Crop(img, image.Rect(1200, 900, 1240, 930)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m := findTemplate(ToGray(img), tpl); m.X != 1200 || m.Y != 900 {
			b.Fatalf("%+v", m)
		}
	}
}

// Chữ xám nhạt mảnh trên nền trắng: thu nhỏ để tìm thô thì mẫu gần như một màu – vẫn phải tìm ra.
func TestFindTemplateLowContrast(t *testing.T) {
	w, h := 1600, 1000
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{250, 250, 250, 255})
		}
	}
	// "chữ": nét dọc 1 px xám nhạt, khoảng cách không đều
	r := rand.New(rand.NewPCG(7, 1))
	for x := 1000; x < 1200; x += 2 + r.IntN(3) {
		for y := 700; y < 740; y++ {
			img.Set(x, y, color.RGBA{243, 243, 243, 255})
		}
	}
	tpl := visiontest.Crop(img, image.Rect(990, 690, 1210, 750))
	m := findTemplate(ToGray(img), ToGray(tpl))
	if m.X != 990 || m.Y != 690 || m.Score < 0.99 {
		t.Fatalf("muốn (990, 690) ~1.0, có %+v", m)
	}
}

// Mẫu chụp ở màn 1x, chạy trên Retina: thu nhỏ ảnh chụp (không phóng to mẫu) – vị trí trả về theo px ảnh chụp.
func TestMatchScaledDownscalesScreen(t *testing.T) {
	big := visiontest.Scene(1600, 1000, 3)
	visiontest.Button(big, 900, 600, 240, 80)
	tpl1x := ToGray(visiontest.Crop(big, image.Rect(900, 600, 1140, 680))).resize(120, 40)
	m, k := MatchScaled(big, model.Rect{W: 800, H: 500}, tplOf(tpl1x), MatchOptions{}, 1, 1)
	if k != 2 || abs(m.X-900) > 2 || abs(m.Y-600) > 2 || abs(m.W-240) > 1 || m.Score < 0.95 {
		t.Fatalf("muốn ~(900, 600) 240 rộng, k = 2; có %+v k = %v", m, k)
	}
}

// Windows: toạ độ = px thật (ảnh / vùng = 1), độ nét nằm ở DPI màn. Mẫu chụp ở Retina (2) tìm trên màn 125 %:
// mẫu thu về 1,25 lần cỡ logic; vị trí vẫn theo px ảnh chụp, tỷ lệ toạ độ vẫn 1.
func TestMatchScaledWindowsDPI(t *testing.T) {
	scr := visiontest.Scene(1000, 800, 5)
	visiontest.Button(scr, 400, 300, 150, 50) // nút 120×40 logic vẽ ở 125 %
	tpl2x := ToGray(visiontest.Crop(scr, image.Rect(400, 300, 550, 350))).resize(240, 80)
	m, k := MatchScaled(scr, model.Rect{W: 1000, H: 800}, tplOf(tpl2x), MatchOptions{}, 2, 1.25)
	if k != 1 || abs(m.X-400) > 2 || abs(m.Y-300) > 2 || abs(m.W-150) > 1 || m.Score < 0.95 {
		t.Fatalf("muốn ~(400, 300) 150 rộng, k = 1; có %+v k = %v", m, k)
	}
}
