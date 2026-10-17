package main

import (
	"image"
	"image/color"
	"math/rand/v2"
	"testing"
	"time"
)

// scene: nền nhiễu nhẹ + vài hình chữ nhật (giống giao diện), trả ảnh và một "nút" (mẫu) lấy từ chính ảnh.
func scene(w, h int, seed uint64) *image.RGBA {
	r := rand.New(rand.NewPCG(seed, 1))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := uint8(225 + r.IntN(20))
			img.Set(x, y, color.RGBA{v, v, v, 255})
		}
	}
	for i := 0; i < 40; i++ {
		x0, y0 := r.IntN(w-60), r.IntN(h-30)
		c := color.RGBA{uint8(r.IntN(200)), uint8(r.IntN(200)), uint8(r.IntN(255)), 255}
		for y := y0; y < y0+10+r.IntN(20); y++ {
			for x := x0; x < x0+20+r.IntN(40); x++ {
				img.Set(x, y, c)
			}
		}
	}
	return img
}

// button vẽ một "nút Gửi" có chữ giả (sọc) tại (x, y) cỡ w×h, nhân hệ số scale (Retina).
func button(img *image.RGBA, x, y, w, h int) {
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			c := color.RGBA{24, 119, 242, 255}
			if yy > h/3 && yy < 2*h/3 && xx > w/5 && xx < 4*w/5 && (xx/3)%2 == 0 {
				c = color.RGBA{255, 255, 255, 255}
			}
			img.Set(x+xx, y+yy, c)
		}
	}
}

func crop(img *image.RGBA, r image.Rectangle) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			out.Set(x, y, img.At(r.Min.X+x, r.Min.Y+y))
		}
	}
	return out
}

func TestFindTemplateExact(t *testing.T) {
	img := scene(1200, 800, 1)
	button(img, 700, 420, 120, 40)
	tpl := crop(img, image.Rect(700, 420, 820, 460))
	start := time.Now()
	m := findTemplate(toGray(img), toGray(tpl))
	if m.X != 700 || m.Y != 420 || m.Score < 0.99 {
		t.Fatalf("muốn (700, 420) ~1.0, có %+v", m)
	}
	t.Logf("1200×800, mẫu 120×40: %v", time.Since(start))
}

// Mẫu chụp ở Retina (2x) tìm trên màn 1x: thu mẫu về cùng tỷ lệ trước khi tìm (như engine làm).
func TestFindTemplateScaled(t *testing.T) {
	big := scene(1600, 1000, 2)
	button(big, 900, 600, 240, 80)
	tpl2x := toGray(crop(big, image.Rect(900, 600, 1140, 680)))
	screen1x := toGray(big).resize(800, 500)
	m := findTemplate(screen1x, tpl2x.resize(tpl2x.w/2, tpl2x.h/2))
	if abs(m.X-450) > 1 || abs(m.Y-300) > 1 || m.Score < 0.9 {
		t.Fatalf("muốn ~(450, 300), có %+v", m)
	}
}

// Không có trên màn hình → độ khớp thấp; có nhiễu / đổi độ sáng nhẹ vẫn khớp cao.
func TestFindTemplateAbsentAndNoise(t *testing.T) {
	img := scene(900, 600, 3)
	other := scene(900, 600, 4)
	button(other, 100, 100, 120, 40)
	tpl := toGray(crop(other, image.Rect(100, 100, 220, 140)))
	if m := findTemplate(toGray(img), tpl); m.Score > 0.6 {
		t.Errorf("không có nút mà khớp %.2f tại (%d, %d)", m.Score, m.X, m.Y)
	}
	button(img, 500, 300, 120, 40)
	g := toGray(img)
	r := rand.New(rand.NewPCG(9, 9))
	for i := range g.p {
		g.p[i] = g.p[i]*0.9 + 10 + float32(r.IntN(9)) - 4 // tối đi 10%, nhiễu ±4
	}
	if m := findTemplate(g, tpl); m.X != 500 || m.Y != 300 || m.Score < 0.95 {
		t.Errorf("nhiễu + đổi sáng: muốn (500, 300) ≥ 0.95, có %+v", m)
	}
}

func TestFindTemplateFlatAndSmall(t *testing.T) {
	img := scene(600, 400, 5)
	for y := 200; y < 230; y++ {
		for x := 300; x < 360; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	flat := toGray(crop(img, image.Rect(305, 205, 355, 225))) // mẫu một màu
	if m := findTemplate(toGray(img), flat); m.X < 300 || m.X > 310 || m.Y < 200 || m.Y > 210 || m.Score < 0.9 {
		t.Errorf("mẫu một màu: %+v", m)
	}
	button(img, 50, 50, 14, 12) // mẫu rất nhỏ
	small := toGray(crop(img, image.Rect(50, 50, 64, 62)))
	if m := findTemplate(toGray(img), small); m.X != 50 || m.Y != 50 {
		t.Errorf("mẫu nhỏ: %+v", m)
	}
}

// Cỡ màn Retina 14" (3024×1964 px), nút 240×80 px.
func BenchmarkFindTemplateRetina(b *testing.B) {
	img := scene(3024, 1964, 7)
	button(img, 2000, 1500, 240, 80)
	tpl := toGray(crop(img, image.Rect(2000, 1500, 2240, 1580)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := findTemplate(toGray(img), tpl)
		if m.X != 2000 || m.Y != 1500 {
			b.Fatalf("%+v", m)
		}
	}
}

// Nút nhỏ 40×30 px trên màn Retina – trường hợp chậm nhất (ít thu nhỏ được).
func BenchmarkFindTemplateRetinaSmall(b *testing.B) {
	img := scene(3024, 1964, 8)
	button(img, 1200, 900, 40, 30)
	tpl := toGray(crop(img, image.Rect(1200, 900, 1240, 930)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m := findTemplate(toGray(img), tpl); m.X != 1200 || m.Y != 900 {
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
	tpl := crop(img, image.Rect(990, 690, 1210, 750))
	m := findTemplate(toGray(img), toGray(tpl))
	if m.X != 990 || m.Y != 690 || m.Score < 0.99 {
		t.Fatalf("muốn (990, 690) ~1.0, có %+v", m)
	}
}

// Mẫu chụp ở màn 1x, chạy trên Retina: thu nhỏ ảnh chụp (không phóng to mẫu) – vị trí trả về theo px ảnh chụp.
func TestMatchScaledDownscalesScreen(t *testing.T) {
	big := scene(1600, 1000, 3)
	button(big, 900, 600, 240, 80)
	tpl1x := toGray(crop(big, image.Rect(900, 600, 1140, 680))).resize(120, 40)
	m, k := matchScaled(big, Rect{W: 800, H: 500}, tpl1x, 1)
	if k != 2 || abs(m.X-900) > 2 || abs(m.Y-600) > 2 || abs(m.W-240) > 1 || m.Score < 0.95 {
		t.Fatalf("muốn ~(900, 600) 240 rộng, k = 2; có %+v k = %v", m, k)
	}
}
