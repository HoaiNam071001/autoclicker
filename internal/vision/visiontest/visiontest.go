// Package visiontest – ảnh giả dùng chung cho test tìm ảnh (vision, engine).
package visiontest

import (
	"image"
	"image/color"
	"math/rand/v2"
)

// Scene: nền nhiễu nhẹ + vài hình chữ nhật (giống giao diện), trả ảnh và một "nút" (mẫu) lấy từ chính ảnh.
func Scene(w, h int, seed uint64) *image.RGBA {
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

// Button vẽ một "nút Gửi" có chữ giả (sọc) tại (x, y) cỡ w×h, nhân hệ số scale (Retina).
func Button(img *image.RGBA, x, y, w, h int) {
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

// Crop cắt vùng r của img thành ảnh mới.
func Crop(img *image.RGBA, r image.Rectangle) *image.RGBA {
	out := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			out.Set(x, y, img.At(r.Min.X+x, r.Min.Y+y))
		}
	}
	return out
}
