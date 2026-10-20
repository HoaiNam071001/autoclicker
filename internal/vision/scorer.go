package vision

import "math"

// scorer tính độ khớp mẫu t tại (x, y) của s. Độ khớp = tương quan chuẩn hoá (NCC): 1 = giống hệt, không phụ thuộc
// độ sáng / tương phản chung. Mẫu gần như một màu (không có chi tiết) không dùng được NCC → so độ chênh trung bình.
// Có mặt nạ (ảnh mẫu PNG có điểm trong suốt): chỉ tính các điểm thấy được – chậm hơn (không dùng ảnh tích phân).
type scorer struct {
	s, t    Gray
	tz      []float32 // mẫu trừ trung bình (có mặt nạ: chỉ các điểm thấy được, cùng thứ tự với offs)
	offs    []int     // có mặt nạ: vị trí điểm thấy được, tính theo mảng của s (y*s.w + x)
	tnorm   float64   // Σ tz²
	n       float64
	flat    bool
	tmean   float64
	minVar  float64   // phương sai cửa sổ (× n) dưới mức này = cửa sổ phẳng, không khớp mẫu có chi tiết
	ii, ii2 []float64 // ảnh tích phân của s và s² (cỡ (w+1)×(h+1)) – trung bình / phương sai cửa sổ O(1)
	iw      int
}

const flatVar = 4.0 // phương sai mẫu dưới mức này (thang 0..255) = coi như một màu

// newScorer: mask rỗng = không mặt nạ; có thì cùng cỡ t, 1 = tính, 0 = bỏ qua.
func newScorer(s, t, mask Gray) *scorer {
	sc := &scorer{s: s, t: t, iw: s.w + 1}
	var vals []float32
	if mask.empty() {
		vals = t.p
	} else {
		sc.offs = []int{} // khác nil kể cả khi không còn điểm nào thấy được (mẫu trong suốt hết) – vẫn đi nhánh mặt nạ
		for i, m := range mask.p {
			if m > 0 {
				vals = append(vals, t.p[i])
				sc.offs = append(sc.offs, (i/t.w)*s.w+i%t.w)
			}
		}
	}
	sc.n = float64(max(len(vals), 1))
	var sum float64
	for _, v := range vals {
		sum += float64(v)
	}
	sc.tmean = sum / sc.n
	sc.tz = make([]float32, len(vals))
	for i, v := range vals {
		d := float64(v) - sc.tmean
		sc.tz[i] = float32(d)
		sc.tnorm += d * d
	}
	sc.flat = sc.tnorm/sc.n < flatVar
	sc.minVar = flatVar * sc.n / 4
	if sc.offs != nil {
		return sc
	}
	sc.ii = make([]float64, (s.w+1)*(s.h+1))
	sc.ii2 = make([]float64, (s.w+1)*(s.h+1))
	for y := 0; y < s.h; y++ {
		var row, row2 float64
		for x := 0; x < s.w; x++ {
			v := float64(s.at(x, y))
			row += v
			row2 += v * v
			sc.ii[(y+1)*sc.iw+x+1] = sc.ii[y*sc.iw+x+1] + row
			sc.ii2[(y+1)*sc.iw+x+1] = sc.ii2[y*sc.iw+x+1] + row2
		}
	}
	return sc
}

func (sc *scorer) window(ii []float64, x, y int) float64 {
	w, h := sc.t.w, sc.t.h
	return ii[(y+h)*sc.iw+x+w] - ii[y*sc.iw+x+w] - ii[(y+h)*sc.iw+x] + ii[y*sc.iw+x]
}

// masked – Σ s, Σ s² và Σ s·tz của cửa sổ tại (x, y), chỉ các điểm thấy được của mặt nạ.
func (sc *scorer) masked(x, y int) (s1, s2, dot float64) {
	base := y*sc.s.w + x
	for i, o := range sc.offs { // float64: Σ s² của mẫu lớn vượt độ chính xác float32, trừ nhau ra phương sai sai
		v := float64(sc.s.p[base+o])
		s1 += v
		s2 += v * v
		dot += v * float64(sc.tz[i])
	}
	return s1, s2, dot
}

// dot – Σ s·tz của cửa sổ tại (x, y) (không mặt nạ).
func (sc *scorer) dot(x, y int) float64 {
	var dot float64
	for ty := 0; ty < sc.t.h; ty++ {
		row := sc.s.p[(y+ty)*sc.s.w+x : (y+ty)*sc.s.w+x+sc.t.w]
		tz := sc.tz[ty*sc.t.w : (ty+1)*sc.t.w]
		var acc float32
		for i, v := range row {
			acc += v * tz[i]
		}
		dot += float64(acc)
	}
	return dot
}

func (sc *scorer) at(x, y int) float64 {
	var sum, sum2, dot float64
	if sc.offs != nil {
		sum, sum2, dot = sc.masked(x, y)
	} else {
		sum, sum2 = sc.window(sc.ii, x, y), sc.window(sc.ii2, x, y)
	}
	if sc.flat {
		// mẫu một màu: khớp khi cửa sổ cũng phẳng và cùng độ sáng
		mean := sum / sc.n
		v := math.Max(sum2/sc.n-mean*mean, 0)
		return math.Max(0, 1-math.Abs(mean-sc.tmean)/255-math.Sqrt(v)/128)
	}
	varSum := sum2 - sum*sum/sc.n
	if varSum < sc.minVar {
		return 0 // cửa sổ phẳng, mẫu có chi tiết → không khớp (không mặt nạ: khỏi tính tích chập)
	}
	if sc.offs == nil {
		dot = sc.dot(x, y)
	}
	return dot / math.Sqrt(varSum*sc.tnorm)
}
