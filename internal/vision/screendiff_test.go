package vision

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func fillRect(img *image.RGBA, r image.Rectangle, c color.RGBA) {
	draw.Draw(img, r, &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func TestDiffPercent(t *testing.T) {
	white := color.RGBA{255, 255, 255, 255}
	a := image.NewRGBA(image.Rect(0, 0, 1920, 1080))
	fillRect(a, a.Rect, white)
	b := image.NewRGBA(a.Rect)
	copy(b.Pix, a.Pix)
	if d := DiffPercent(Thumb(a), Thumb(b)); d != 0 {
		t.Fatalf("ảnh y hệt: %v", d)
	}
	if g := Thumb(a); g.w != 320 || g.h != 180 {
		t.Fatalf("cỡ ảnh thu nhỏ %d × %d", g.w, g.h)
	}

	// con trỏ chữ 2 × 20 px nhấp nháy: rất nhỏ so với ngưỡng 1 %
	fillRect(b, image.Rect(500, 500, 502, 520), color.RGBA{A: 255})
	if d := DiffPercent(Thumb(a), Thumb(b)); d <= 0 || d >= 0.1 {
		t.Errorf("con trỏ chữ: %v %%", d)
	}
	// hộp thoại 480 × 270 hiện giữa màn = 1/16 diện tích
	fillRect(b, image.Rect(720, 405, 1200, 675), color.RGBA{40, 40, 40, 255})
	if d := DiffPercent(Thumb(a), Thumb(b)); d < 6 || d > 7 {
		t.Errorf("hộp thoại: %v %%", d)
	}
	// nhiễu nhẹ toàn ảnh (nén, khử răng cưa) không tính
	c := image.NewRGBA(a.Rect)
	fillRect(c, c.Rect, color.RGBA{245, 245, 245, 255})
	if d := DiffPercent(Thumb(a), Thumb(c)); d != 0 {
		t.Errorf("lệch sáng nhẹ: %v", d)
	}
	if d := DiffPercent(Thumb(a), Thumb(image.NewRGBA(image.Rect(0, 0, 100, 100)))); d != 100 {
		t.Errorf("khác cỡ: %v", d)
	}
}
