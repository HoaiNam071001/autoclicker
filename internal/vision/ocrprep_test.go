package vision

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

// Chữ đen trên nền trắng và chữ trắng trên nút xanh đều ra đen trên trắng, ảnh phóng đúng hệ số.
func TestPrepOCR(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 120, 40))
	fill := func(r image.Rectangle, c color.Color) {
		draw.Draw(img, r, image.NewUniform(c), image.Point{}, draw.Src)
	}
	fill(img.Rect, color.White)
	fill(image.Rect(60, 0, 120, 40), color.RGBA{37, 99, 235, 255}) // nút xanh
	fill(image.Rect(20, 10, 23, 30), color.Black)                  // nét chữ tối
	fill(image.Rect(80, 10, 83, 30), color.White)                  // nét chữ trắng trên nút

	g := PrepOCR(img, 2, 5)
	if g.Rect.Dx() != 250 || g.Rect.Dy() != 90 {
		t.Fatalf("cỡ %v", g.Rect)
	}
	at := func(x, y int) uint8 { return g.GrayAt(x*2+1+5, y*2+1+5).Y }
	for _, p := range []image.Point{{21, 20}, {81, 20}} {
		if v := at(p.X, p.Y); v > 60 {
			t.Errorf("nét chữ ở %v phải đen, được %d", p, v)
		}
	}
	for _, p := range []image.Point{{5, 5}, {40, 20}, {70, 5}, {100, 30}} {
		if v := at(p.X, p.Y); v < 200 {
			t.Errorf("nền ở %v phải trắng, được %d", p, v)
		}
	}
}
