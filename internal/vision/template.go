package vision

import (
	"image"
	"math"
)

// Ảnh mẫu đã giải mã cho tìm ảnh (bước 40): ảnh xám, mặt nạ (điểm trong suốt bỏ qua) và màu trung bình để kiểm tra màu.

type Template struct {
	Gray  Gray
	mask  Gray       // rỗng = ảnh không có điểm trong suốt; có thì cùng cỡ Gray, 1 = tính, 0 = bỏ qua
	color [3]float64 // RGB trung bình các điểm thấy được
}

// alphaCut – alpha dưới mức này (0..255) coi là trong suốt.
const alphaCut = 128

func NewTemplate(img image.Image) Template {
	t := Template{Gray: ToGray(img)}
	b := img.Bounds()
	mask := Gray{w: b.Dx(), h: b.Dy(), p: make([]float32, b.Dx()*b.Dy())}
	var sum [3]float64
	n, holes := 0, false
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, bb, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			if a>>8 < alphaCut {
				holes = true
				continue
			}
			mask.p[y*mask.w+x] = 1
			// RGBA() đã nhân alpha – chia lại để ra màu thật của điểm hơi trong
			sum[0] += float64(r) / float64(a) * 255
			sum[1] += float64(g) / float64(a) * 255
			sum[2] += float64(bb) / float64(a) * 255
			n++
		}
	}
	if holes {
		t.mask = mask
	}
	if n > 0 {
		t.color = [3]float64{sum[0] / float64(n), sum[1] / float64(n), sum[2] / float64(n)}
	}
	return t
}

// Masked – ảnh mẫu có vùng trong suốt (bỏ qua khi so).
func (t Template) Masked() bool { return !t.mask.empty() }

// MatchOptions – cách so (Step.MatchMode), độ lệch màu tối đa (Step.ColorTol, 0 = không kiểm tra màu) và độ khớp
// tối thiểu 0..1 (Step.Threshold – chỗ đúng màu mà khớp thấp hơn thì không chọn thay chỗ giống hình nhất).
type MatchOptions struct {
	Mode     string
	ColorTol int
	MinScore float64
}

// colorDiff – độ lệch (0..100 %) giữa màu trung bình chỗ r (px của img) với màu trung bình ảnh mẫu; điểm trong suốt
// của mẫu không tính. Lấy mẫu thưa (≤ ~64×64 điểm) – chỉ cần trung bình.
func (t Template) colorDiff(img image.Image, r image.Rectangle) int {
	r = r.Add(img.Bounds().Min).Intersect(img.Bounds())
	if r.Empty() {
		return 100
	}
	step := max(1, max(r.Dx(), r.Dy())/64)
	var sum [3]float64
	n := 0
	for y := r.Min.Y; y < r.Max.Y; y += step {
		for x := r.Min.X; x < r.Max.X; x += step {
			if t.Masked() {
				tx := min((x-r.Min.X)*t.mask.w/r.Dx(), t.mask.w-1)
				ty := min((y-r.Min.Y)*t.mask.h/r.Dy(), t.mask.h-1)
				if t.mask.at(tx, ty) == 0 {
					continue
				}
			}
			cr, cg, cb, _ := img.At(x, y).RGBA()
			sum[0] += float64(cr >> 8)
			sum[1] += float64(cg >> 8)
			sum[2] += float64(cb >> 8)
			n++
		}
	}
	if n == 0 {
		return 100
	}
	var d2 float64
	for i := range sum {
		d := sum[i]/float64(n) - t.color[i]
		d2 += d * d
	}
	return int(math.Round(math.Sqrt(d2) / (255 * math.Sqrt(3)) * 100))
}

// ColorOK – chỗ thấy đạt kiểm tra màu (luôn đúng khi không kiểm tra).
func (m Match) ColorOK(opt MatchOptions) bool {
	return opt.ColorTol <= 0 || m.ColorDiff <= opt.ColorTol
}
