package vision

import (
	"image"
)

// Chuẩn bị ảnh cho Tesseract (bước 40 O1). Chữ giao diện nhỏ (11–14 px) phải phóng ×2 mới đọc đúng dấu; chữ sáng trên
// nền màu (nút xanh chữ trắng) Tesseract đọc ra rác → đưa mọi chữ về đen trên trắng: điểm càng lệch nền cục bộ càng đen.
// Đo ở 40 O0: phóng ×2 + chuẩn hoá nền đọc đúng 15/15 chuỗi cả theme sáng / tối, DPI 125 %.

const (
	prepWin      = 12  // nền = độ sáng hay gặp nhất trong cửa sổ ±prepWin px (ảnh đã phóng) quanh điểm
	prepContrast = 120 // lệch nền bao nhiêu thì đen hẳn
)

// PrepOCR trả ảnh xám đã phóng scale lần, chữ đen trên nền trắng, thêm lề trắng pad px mỗi bên (Tesseract hay bỏ sót
// chữ sát mép khi vùng khoanh vừa khít).
func PrepOCR(src image.Image, scale, pad int) *image.Gray {
	gs := ToGray(src).resize(src.Bounds().Dx()*scale, src.Bounds().Dy()*scale) // song tuyến – đủ cho OCR, nhanh hơn CatmullRom nhiều
	g := image.NewGray(image.Rect(0, 0, gs.w, gs.h))
	for i, v := range gs.p {
		g.Pix[i] = uint8(min(max(v+0.5, 0), 255))
	}
	w, h := g.Rect.Dx(), g.Rect.Dy()
	bg, bw := backgroundMap(g)
	out := image.NewGray(image.Rect(0, 0, w+2*pad, h+2*pad))
	for i := range out.Pix {
		out.Pix[i] = 255
	}
	for y := 0; y < h; y++ {
		row := bg[(y/2)*bw:]
		dst := out.Pix[(y+pad)*out.Stride+pad:]
		for x := 0; x < w; x++ {
			d := int(g.Pix[y*g.Stride+x]) - int(row[x/2])
			if d < 0 {
				d = -d
			}
			dst[x] = uint8(255 - min(255, d*255/prepContrast))
		}
	}
	return out
}

// backgroundMap – nền quanh từng điểm = bậc hay gặp nhất (histogram 32 bậc) trong cửa sổ ±prepWin: chữ mảnh nên hiếm
// khi chiếm đa số, cửa sổ nhỏ nên mép nút / ô màu vẫn đúng nền từng bên. Tính ở nửa độ phân giải (điểm chẵn, lấy mẫu
// cách 2 px), histogram trượt theo hàng. Trả bản đồ + chiều rộng của nó.
func backgroundMap(g *image.Gray) ([]uint8, int) {
	w, h := g.Rect.Dx(), g.Rect.Dy()
	bw, bh := (w+1)/2, (h+1)/2
	out := make([]uint8, bw*bh)
	col := func(hist *[32]int32, x, y0, y1 int, d int32) {
		if x < 0 || x >= w {
			return
		}
		for y := y0; y <= y1; y += 2 {
			hist[g.Pix[y*g.Stride+x]>>3] += d
		}
	}
	for by := 0; by < bh; by++ {
		y := by * 2
		y0, y1 := max(0, y-prepWin), min(h-1, y+prepWin)
		y0 += (y - y0) % 2 // giữ cùng chẵn lẻ với y – lấy mẫu cách 2
		var hist [32]int32
		for x := -prepWin; x <= prepWin; x += 2 {
			col(&hist, x, y0, y1, 1)
		}
		for bx := 0; bx < bw; bx++ {
			if bx > 0 {
				x := bx * 2
				col(&hist, x-prepWin-2, y0, y1, -1)
				col(&hist, x+prepWin, y0, y1, 1)
			}
			best := 0
			for k := range hist {
				if hist[k] > hist[best] {
					best = k
				}
			}
			out[by*bw+bx] = uint8(best*8 + 4)
		}
	}
	return out, bw
}
