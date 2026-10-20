package vision

import (
	"image"
	"math"
)

// Ảnh xám (float32) cho tìm ảnh: đổi từ ảnh màu, đổi cỡ, ảnh viền (Sobel), mặt nạ.

type Gray struct {
	w, h int
	p    []float32
}

func (g Gray) at(x, y int) float32 { return g.p[y*g.w+x] }

func (g Gray) empty() bool { return g.w == 0 || g.h == 0 }

func ToGray(img image.Image) Gray {
	b := img.Bounds()
	g := Gray{w: b.Dx(), h: b.Dy(), p: make([]float32, b.Dx()*b.Dy())}
	if rgba, ok := img.(*image.RGBA); ok { // ảnh chụp màn hình: đọc thẳng mảng pixel, nhanh hơn At() nhiều lần
		for y := 0; y < g.h; y++ {
			row := rgba.Pix[(b.Min.Y+y-rgba.Rect.Min.Y)*rgba.Stride+(b.Min.X-rgba.Rect.Min.X)*4:]
			for x := 0; x < g.w; x++ {
				g.p[y*g.w+x] = 0.299*float32(row[x*4]) + 0.587*float32(row[x*4+1]) + 0.114*float32(row[x*4+2])
			}
		}
		return g
	}
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			r, gg, bb, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			g.p[y*g.w+x] = float32(0.299*float64(r>>8) + 0.587*float64(gg>>8) + 0.114*float64(bb>>8))
		}
	}
	return g
}

// resize đổi cỡ: thu nhỏ lấy trung bình theo ô (khử nhiễu), phóng to nội suy song tuyến.
func (g Gray) resize(nw, nh int) Gray {
	nw, nh = max(nw, 1), max(nh, 1)
	if nw == g.w && nh == g.h {
		return g
	}
	out := Gray{w: nw, h: nh, p: make([]float32, nw*nh)}
	sx, sy := float64(g.w)/float64(nw), float64(g.h)/float64(nh)
	if sx >= 1 && sy >= 1 {
		for y := 0; y < nh; y++ {
			y0, y1 := int(float64(y)*sy), max(int(float64(y+1)*sy), int(float64(y)*sy)+1)
			for x := 0; x < nw; x++ {
				x0, x1 := int(float64(x)*sx), max(int(float64(x+1)*sx), int(float64(x)*sx)+1)
				var sum float32
				for yy := y0; yy < min(y1, g.h); yy++ {
					for xx := x0; xx < min(x1, g.w); xx++ {
						sum += g.at(xx, yy)
					}
				}
				out.p[y*nw+x] = sum / float32((min(y1, g.h)-y0)*(min(x1, g.w)-x0))
			}
		}
		return out
	}
	for y := 0; y < nh; y++ {
		fy := math.Max(0, (float64(y)+0.5)*sy-0.5)
		y0 := min(int(fy), g.h-1)
		y1, ty := min(y0+1, g.h-1), float32(fy-float64(y0))
		for x := 0; x < nw; x++ {
			fx := math.Max(0, (float64(x)+0.5)*sx-0.5)
			x0 := min(int(fx), g.w-1)
			x1, tx := min(x0+1, g.w-1), float32(fx-float64(x0))
			top := g.at(x0, y0)*(1-tx) + g.at(x1, y0)*tx
			bot := g.at(x0, y1)*(1-tx) + g.at(x1, y1)*tx
			out.p[y*nw+x] = top*(1-ty) + bot*ty
		}
	}
	return out
}

// edges – độ lớn gradient Sobel (thang ~0..255). Đảo sáng / tối (đổi theme) không đổi độ lớn gradient, nên so NCC
// trên ảnh viền chịu được đổi theme / màu nền. Viền ngoài 1 px không có đủ lân cận → 0 (mẫu thì cắt bỏ, xem MatchScaled).
func (g Gray) edges() Gray {
	out := Gray{w: g.w, h: g.h, p: make([]float32, len(g.p))}
	for y := 1; y < g.h-1; y++ {
		for x := 1; x < g.w-1; x++ {
			a, b, c := g.at(x-1, y-1), g.at(x, y-1), g.at(x+1, y-1)
			d, f := g.at(x-1, y), g.at(x+1, y)
			gg, h, i := g.at(x-1, y+1), g.at(x, y+1), g.at(x+1, y+1)
			gx := (c + 2*f + i) - (a + 2*d + gg)
			gy := (gg + 2*h + i) - (a + 2*b + c)
			out.p[y*g.w+x] = float32(math.Sqrt(float64(gx*gx+gy*gy))) / 4
		}
	}
	return out
}

// crop cắt bớt n px mỗi mép.
func (g Gray) crop(n int) Gray {
	if g.empty() || g.w <= 2*n || g.h <= 2*n {
		return g
	}
	out := Gray{w: g.w - 2*n, h: g.h - 2*n}
	out.p = make([]float32, out.w*out.h)
	for y := 0; y < out.h; y++ {
		copy(out.p[y*out.w:(y+1)*out.w], g.p[(y+n)*g.w+n:(y+n)*g.w+n+out.w])
	}
	return out
}

// maskResize – đổi cỡ mặt nạ (0 / 1) rồi làm tròn lại về 0 / 1: điểm ảnh ít nhất nửa thấy được thì tính.
func (g Gray) maskResize(nw, nh int) Gray {
	if g.empty() {
		return g
	}
	if nw == g.w && nh == g.h {
		return g
	}
	out := g.resize(nw, nh)
	for i, v := range out.p {
		if v >= 0.5 {
			out.p[i] = 1
		} else {
			out.p[i] = 0
		}
	}
	return out
}

// erode – mặt nạ thu vào 1 px: điểm cạnh vùng bỏ qua cũng bỏ (so viền: gradient ở mép lỗ lấy cả điểm trong suốt).
func (g Gray) erode() Gray {
	if g.empty() {
		return g
	}
	out := Gray{w: g.w, h: g.h, p: make([]float32, len(g.p))}
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			v := float32(1)
			for dy := -1; dy <= 1 && v > 0; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if xx, yy := x+dx, y+dy; xx >= 0 && yy >= 0 && xx < g.w && yy < g.h && g.at(xx, yy) == 0 {
						v = 0
						break
					}
				}
			}
			out.p[y*g.w+x] = v
		}
	}
	return out
}
