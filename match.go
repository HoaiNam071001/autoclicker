package main

import (
	"image"
	"math"
	"sort"
)

// Tìm ảnh mẫu trong ảnh chụp màn hình (bước 25) – Go thuần, không cần OpenCV.
// Độ khớp = tương quan chuẩn hoá (NCC) trên ảnh xám: 1 = giống hệt, không phụ thuộc độ sáng / tương phản chung.
// Tìm thô trên ảnh thu nhỏ (giữ vài ứng viên tốt nhất) rồi tìm mịn ở độ phân giải gốc quanh từng ứng viên.

type gray struct {
	w, h int
	p    []float32
}

func (g gray) at(x, y int) float32 { return g.p[y*g.w+x] }

func toGray(img image.Image) gray {
	b := img.Bounds()
	g := gray{w: b.Dx(), h: b.Dy(), p: make([]float32, b.Dx()*b.Dy())}
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
func (g gray) resize(nw, nh int) gray {
	nw, nh = max(nw, 1), max(nh, 1)
	if nw == g.w && nh == g.h {
		return g
	}
	out := gray{w: nw, h: nh, p: make([]float32, nw*nh)}
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

// Match – vị trí (góc trên trái, px của ảnh chụp) và độ khớp 0..1 của chỗ giống mẫu nhất.
type Match struct {
	X, Y, W, H int
	Score      float64
}

// scorer tính độ khớp mẫu t tại (x, y) của s. Mẫu gần như một màu (không có chi tiết) không dùng được NCC
// → so độ chênh trung bình.
type scorer struct {
	s, t    gray
	tz      []float32 // mẫu trừ trung bình
	tnorm   float64   // Σ tz²
	n       float64
	flat    bool
	tmean   float64
	minVar  float64   // phương sai cửa sổ (× n) dưới mức này = cửa sổ phẳng, không khớp mẫu có chi tiết
	ii, ii2 []float64 // ảnh tích phân của s và s² (cỡ (w+1)×(h+1)) – trung bình / phương sai cửa sổ O(1)
	iw      int
}

const flatVar = 4.0 // phương sai mẫu dưới mức này (thang 0..255) = coi như một màu

func newScorer(s, t gray) *scorer {
	sc := &scorer{s: s, t: t, n: float64(t.w * t.h), iw: s.w + 1}
	var sum float64
	for _, v := range t.p {
		sum += float64(v)
	}
	sc.tmean = sum / sc.n
	sc.tz = make([]float32, len(t.p))
	for i, v := range t.p {
		d := float64(v) - sc.tmean
		sc.tz[i] = float32(d)
		sc.tnorm += d * d
	}
	sc.flat = sc.tnorm/sc.n < flatVar
	sc.minVar = flatVar * sc.n / 4
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

func (sc *scorer) at(x, y int) float64 {
	sum := sc.window(sc.ii, x, y)
	if sc.flat {
		// mẫu một màu: khớp khi cửa sổ cũng phẳng và cùng độ sáng
		mean := sum / sc.n
		v := math.Max(sc.window(sc.ii2, x, y)/sc.n-mean*mean, 0)
		return math.Max(0, 1-math.Abs(mean-sc.tmean)/255-math.Sqrt(v)/128)
	}
	varSum := sc.window(sc.ii2, x, y) - sum*sum/sc.n
	if varSum < sc.minVar {
		return 0 // cửa sổ phẳng, mẫu có chi tiết → không khớp
	}
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
	return dot / math.Sqrt(varSum*sc.tnorm)
}

const (
	coarseMaxPixels = 250_000 // ảnh tìm thô tối đa chừng này px
	coarseTpl       = 12      // thu nhỏ tới khi cạnh ngắn của mẫu còn chừng này px (đủ chi tiết để chọn ứng viên)
	coarseMinTpl    = 4       // nhưng không nhỏ hơn
	candidates      = 12      // nhiều chỗ na ná nhau (chữ, ô danh sách) → ứng viên đúng hay rơi khỏi top vài cái
)

// findTemplate trả chỗ giống mẫu t nhất trong s (cùng tỷ lệ px). Mẫu lớn hơn ảnh → Score 0.
func findTemplate(s, t gray) Match {
	if t.w > s.w || t.h > s.h || t.w == 0 || t.h == 0 {
		return Match{}
	}
	// hệ số thu nhỏ cho bước thô: đủ nhỏ để nhanh, mẫu vẫn còn chi tiết
	f := max(1, int(math.Ceil(math.Sqrt(float64(s.w*s.h)/coarseMaxPixels))), min(t.w, t.h)/coarseTpl)
	f = min(f, max(1, min(t.w, t.h)/coarseMinTpl))
	type cand struct {
		x, y  int
		score float64
	}
	var cands []cand
	if f == 1 {
		cands = []cand{{0, 0, 0}}
	} else {
		cs, ct := s.resize(s.w/f, s.h/f), t.resize(t.w/f, t.h/f)
		sc := newScorer(cs, ct)
		if full := newScorer(t, t); !full.flat && sc.tnorm > 0 {
			// mẫu ít tương phản (chữ xám nhạt…): thu nhỏ xong có thể bị coi là một màu → vẫn so NCC,
			// và hạ ngưỡng cửa sổ phẳng theo độ chi tiết của chính mẫu thu nhỏ
			sc.flat = false
			sc.minVar = math.Min(sc.minVar, sc.tnorm/4)
		}
		all := make([]cand, 0, (cs.w-ct.w+1)*(cs.h-ct.h+1))
		for y := 0; y <= cs.h-ct.h; y++ {
			for x := 0; x <= cs.w-ct.w; x++ {
				all = append(all, cand{x, y, sc.at(x, y)})
			}
		}
		sort.Slice(all, func(i, j int) bool { return all[i].score > all[j].score })
		// giữ vài ứng viên tốt nhất, cách nhau ít nhất nửa mẫu (khỏi dồn hết vào một chỗ)
		for _, c := range all {
			ok := true
			for _, k := range cands {
				if abs(c.x-k.x) < max(ct.w/2, 1) && abs(c.y-k.y) < max(ct.h/2, 1) {
					ok = false
					break
				}
			}
			if ok {
				cands = append(cands, c)
				if len(cands) == candidates {
					break
				}
			}
		}
	}
	// tìm mịn ở độ phân giải gốc quanh từng ứng viên (f == 1: quét cả ảnh)
	sc := newScorer(s, t)
	best := Match{W: t.w, H: t.h, Score: -1}
	for _, c := range cands {
		x0, x1, y0, y1 := 0, s.w-t.w, 0, s.h-t.h
		if f > 1 {
			x0, x1 = max(c.x*f-f, 0), min(c.x*f+f, s.w-t.w)
			y0, y1 = max(c.y*f-f, 0), min(c.y*f+f, s.h-t.h)
		}
		for y := y0; y <= y1; y++ {
			for x := x0; x <= x1; x++ {
				if v := sc.at(x, y); v > best.Score {
					best.X, best.Y, best.Score = x, y, v
				}
			}
		}
	}
	best.Score = math.Max(best.Score, 0)
	return best
}
