// Package vision – tìm ảnh mẫu / chữ trên ảnh đã chụp (thuần, không chụp màn hình, không đọc chữ).
package vision

import (
	"image"
	"math"
	"sort"

	"ghosthand/internal/model"
)

// Tìm ảnh mẫu trong ảnh chụp màn hình (bước 25, 40) – Go thuần, không cần OpenCV. Độ khớp NCC (scorer.go) trên ảnh
// xám, hoặc ảnh viền (MatchEdge – chịu đổi theme). Tìm thô trên ảnh thu nhỏ (giữ vài ứng viên tốt nhất) rồi tìm mịn ở
// độ phân giải gốc quanh từng ứng viên; kiểm tra màu chọn ứng viên tốt nhất có màu đúng.

// Match – vị trí (góc trên trái, px của ảnh chụp) và độ khớp 0..1 của chỗ giống mẫu nhất.
type Match struct {
	X, Y, W, H int
	Score      float64
	ColorDiff  int // % lệch màu trung bình so với ảnh mẫu – chỉ tính khi MatchOptions.ColorTol > 0
}

const (
	coarseMaxPixels = 250_000 // ảnh tìm thô tối đa chừng này px
	coarseTpl       = 12      // thu nhỏ tới khi cạnh ngắn của mẫu còn chừng này px (đủ chi tiết để chọn ứng viên)
	coarseMinTpl    = 4       // nhưng không nhỏ hơn
	candidates      = 12      // nhiều chỗ na ná nhau (chữ, ô danh sách) → ứng viên đúng hay rơi khỏi top vài cái
)

type cand struct {
	x, y  int
	score float64
}

// peaks – tối đa n chỗ điểm cao nhất của score trên lưới w×h, cách nhau ít nhất nửa mẫu (khỏi dồn hết vào một chỗ).
func peaks(w, h, tw, th, n int, score func(x, y int) float64) []cand {
	all := make([]cand, 0, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			all = append(all, cand{x, y, score(x, y)})
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].score > all[j].score })
	var out []cand
	for _, c := range all {
		ok := true
		for _, k := range out {
			if abs(c.x-k.x) < max(tw/2, 1) && abs(c.y-k.y) < max(th/2, 1) {
				ok = false
				break
			}
		}
		if ok {
			if out = append(out, c); len(out) == n {
				break
			}
		}
	}
	return out
}

// findPeaks trả các chỗ giống mẫu t nhất trong s (cùng tỷ lệ px), điểm giảm dần; mask rỗng = không mặt nạ.
// Mẫu lớn hơn ảnh → một chỗ Score 0.
func findPeaks(s, t, mask Gray, n int) []Match {
	if t.w > s.w || t.h > s.h || t.empty() {
		return []Match{{}}
	}
	sc := newScorer(s, t, mask)
	// hệ số thu nhỏ cho bước thô: đủ nhỏ để nhanh, mẫu vẫn còn chi tiết
	f := max(1, int(math.Ceil(math.Sqrt(float64(s.w*s.h)/coarseMaxPixels))), min(t.w, t.h)/coarseTpl)
	f = min(f, max(1, min(t.w, t.h)/coarseMinTpl))
	if f == 1 { // ảnh nhỏ: quét cả ảnh ở độ phân giải gốc
		var out []Match
		for _, c := range peaks(s.w-t.w+1, s.h-t.h+1, t.w, t.h, n, sc.at) {
			out = append(out, Match{X: c.x, Y: c.y, W: t.w, H: t.h, Score: math.Max(c.score, 0)})
		}
		return out
	}
	cs, ct := s.resize(s.w/f, s.h/f), t.resize(t.w/f, t.h/f)
	csc := newScorer(cs, ct, mask.maskResize(ct.w, ct.h))
	if !sc.flat && csc.tnorm > 0 {
		// mẫu ít tương phản (chữ xám nhạt…): thu nhỏ xong có thể bị coi là một màu → vẫn so NCC,
		// và hạ ngưỡng cửa sổ phẳng theo độ chi tiết của chính mẫu thu nhỏ
		csc.flat = false
		csc.minVar = math.Min(csc.minVar, csc.tnorm/4)
	}
	// tìm mịn ở độ phân giải gốc quanh từng ứng viên thô
	var out []Match
	for _, c := range peaks(cs.w-ct.w+1, cs.h-ct.h+1, ct.w, ct.h, n, csc.at) {
		best := Match{W: t.w, H: t.h, Score: -1}
		for y := max(c.y*f-f, 0); y <= min(c.y*f+f, s.h-t.h); y++ {
			for x := max(c.x*f-f, 0); x <= min(c.x*f+f, s.w-t.w); x++ {
				if v := sc.at(x, y); v > best.Score {
					best.X, best.Y, best.Score = x, y, v
				}
			}
		}
		best.Score = math.Max(best.Score, 0)
		out = append(out, best)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ImageScale: px ảnh chụp trên 1 đơn vị toạ độ.
func ImageScale(img image.Image, r model.Rect) float64 {
	if r.W == 0 {
		return 1
	}
	return float64(img.Bounds().Dx()) / float64(r.W)
}

// MatchScaled đưa mẫu và ảnh chụp về cùng tỷ lệ px (mẫu chụp ở màn 1x, chạy trên Retina…) rồi tìm. Luôn thu nhỏ
// bên nét hơn (trung bình theo ô ≈ cách màn 1x vẽ) thay vì phóng to bên kia – phóng to làm mẫu nhoè, độ khớp tụt.
// tplScale và độ nét so sánh theo px ảnh / 1 px logic (đơn vị app vẽ giao diện): = ImageScale × dpi, dpi = độ nét màn
// chứa vùng ngoài phần toạ độ đã gồm (Windows: DPI / 96 vì toạ độ là px thật; macOS 1 – toạ độ điểm đã gồm Retina).
// Trả chỗ tốt nhất (có kiểm tra màu: chỗ đúng màu tốt nhất đạt MinScore, không có thì chỗ tốt nhất kèm ColorDiff – để
// báo "giống hình nhưng lệch màu") theo px ảnh chụp và tỷ lệ px / toạ độ của ảnh chụp.
func MatchScaled(img image.Image, area model.Rect, tpl Template, opt MatchOptions, tplScale, dpi float64) (Match, float64) {
	ms, sc := scaledPeaks(img, area, tpl, opt, tplScale, dpi, candidates)
	if opt.ColorTol <= 0 {
		return ms[0], sc
	}
	for i := range ms {
		ms[i].ColorDiff = tpl.colorDiff(img, ms[i].rect())
		if ms[i].Score < opt.MinScore {
			break
		}
		if ms[i].ColorOK(opt) {
			return ms[i], sc
		}
	}
	return ms[0], sc
}

func (m Match) rect() image.Rectangle { return image.Rect(m.X, m.Y, m.X+m.W, m.Y+m.H) }

// scaledPeaks – n chỗ giống mẫu nhất (điểm giảm dần, chưa kiểm tra màu) theo px ảnh chụp, sau khi đưa mẫu và ảnh về
// cùng tỷ lệ như MatchScaled tả; kèm tỷ lệ px / toạ độ của ảnh chụp.
func scaledPeaks(img image.Image, area model.Rect, tpl Template, opt MatchOptions, tplScale, dpi float64, n int) ([]Match, float64) {
	sc := ImageScale(img, area)
	g, t, mask := ToGray(img), tpl.Gray, tpl.mask
	eff := sc * max(dpi, 0.01)
	k := 1.0 // ảnh chụp đã thu nhỏ k lần
	if tplScale > 0 && math.Abs(eff-tplScale) > 0.01 {
		if r := tplScale / eff; r > 1 { // mẫu nét hơn ảnh chụp → thu nhỏ mẫu
			w, h := int(math.Round(float64(t.w)/r)), int(math.Round(float64(t.h)/r))
			t, mask = t.resize(w, h), mask.maskResize(w, h)
		} else { // ảnh chụp nét hơn → thu nhỏ ảnh chụp, đổi kết quả về px ảnh chụp
			k = eff / tplScale
			g = g.resize(int(math.Round(float64(g.w)/k)), int(math.Round(float64(g.h)/k)))
		}
	}
	pad := 0
	if opt.Mode == model.MatchEdge {
		// viền ngoài của mẫu không có đủ lân cận để tính gradient → cắt 1 px, cộng lại vào kết quả; mép lỗ trong suốt cũng vậy
		et := t.edges().crop(1) // mẫu quá nhỏ thì crop không cắt
		g, t, mask, pad = g.edges(), et, mask.erode().crop(1), (t.w-et.w)/2
	}
	ms := findPeaks(g, t, mask, n)
	for i := range ms {
		m := &ms[i]
		m.X, m.Y, m.W, m.H = m.X-pad, m.Y-pad, m.W+2*pad, m.H+2*pad
		if k != 1 {
			m.X, m.Y = int(math.Round(float64(m.X)*k)), int(math.Round(float64(m.Y)*k))
			m.W, m.H = int(math.Round(float64(m.W)*k)), int(math.Round(float64(m.H)*k))
		}
	}
	return ms, sc
}
