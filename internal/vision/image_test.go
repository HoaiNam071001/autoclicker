package vision

import (
	"image"
	"image/color"
	"testing"
)

func TestPixelHash(t *testing.T) {
	a := image.NewRGBA(image.Rect(0, 0, 4, 3))
	b := image.NewRGBA(image.Rect(0, 0, 4, 3))
	if PixelHash(a) != PixelHash(b) {
		t.Error("ảnh giống nhau phải cùng hash")
	}
	b.Set(3, 2, color.RGBA{R: 1, A: 255})
	if PixelHash(a) == PixelHash(b) {
		t.Error("đổi một điểm phải đổi hash")
	}
	if PixelHash(a) == PixelHash(image.NewRGBA(image.Rect(0, 0, 3, 4))) {
		t.Error("khác cỡ phải khác hash")
	}
	// ảnh con (SubImage) và ảnh không phải RGBA hash theo điểm ảnh, không theo bộ nhớ
	sub := a.SubImage(image.Rect(1, 1, 3, 3)).(*image.RGBA)
	if PixelHash(sub) != PixelHash(image.NewRGBA(image.Rect(0, 0, 2, 2))) {
		t.Error("ảnh con")
	}
	if PixelHash(image.NewNRGBA(image.Rect(0, 0, 2, 2))) == 0 {
		t.Error("NRGBA")
	}
}
