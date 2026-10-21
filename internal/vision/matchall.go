package vision

import (
	"image"
	"sort"

	"ghosthand/internal/model"
)

// Tìm mọi chỗ khớp ảnh mẫu (node Tìm tất cả, bước 42 G3): nhiều ứng viên hơn MatchScaled, gộp chỗ chồng lấn, thứ tự đọc.

// MatchAllMax – tối đa chừng này chỗ mỗi lần tìm (cũng là số ứng viên thô giữ lại).
const MatchAllMax = 200

// MatchAll trả mọi chỗ khớp ≥ opt.MinScore (và đúng màu khi opt.ColorTol > 0) theo px ảnh chụp, xếp theo thứ tự đọc,
// kèm tỷ lệ px / toạ độ của ảnh chụp. Hai chỗ chồng nhau (tâm chỗ này nằm trong khung chỗ kia) giữ chỗ khớp cao hơn.
func MatchAll(img image.Image, area model.Rect, tpl Template, opt MatchOptions, tplScale, dpi float64) ([]Match, float64) {
	ms, sc := scaledPeaks(img, area, tpl, opt, tplScale, dpi, MatchAllMax)
	sort.SliceStable(ms, func(i, j int) bool { return ms[i].Score > ms[j].Score }) // tìm mịn có thể đổi thứ tự điểm
	var out []Match
	for _, m := range ms {
		if m.Score < opt.MinScore || m.Score <= 0 || overlaps(m, out) {
			continue
		}
		if opt.ColorTol > 0 {
			if m.ColorDiff = tpl.colorDiff(img, m.rect()); !m.ColorOK(opt) {
				continue
			}
		}
		out = append(out, m)
	}
	sort.SliceStable(out, func(i, j int) bool { return ReadingLess(out[i].box(), out[j].box()) })
	return out, sc
}

func (m Match) box() model.Rect { return model.Rect{X: m.X, Y: m.Y, W: m.W, H: m.H} }

// overlaps – tâm m nằm trong khung một chỗ đã giữ (hoặc ngược lại): cùng một chỗ thấy hai lần.
func overlaps(m Match, kept []Match) bool {
	in := func(a, b Match) bool {
		cx, cy := a.X+a.W/2, a.Y+a.H/2
		return cx >= b.X && cx < b.X+b.W && cy >= b.Y && cy < b.Y+b.H
	}
	for _, k := range kept {
		if in(m, k) || in(k, m) {
			return true
		}
	}
	return false
}
