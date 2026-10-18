package vision

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf16"

	"autoclicker/internal/model"
	"golang.org/x/text/unicode/norm"
)

// TextWord – một từ trong dòng: vị trí theo UTF-16 trong TextLine.Text (như Vision báo) và khung theo px ảnh.
type TextWord struct {
	Start, Len int
	Box        model.Rect
}

// TextLine – một dòng chữ đọc được, khung theo px ảnh (gốc trên trái).
type TextLine struct {
	Text  string
	Box   model.Rect
	Words []TextWord
}

// TextQuery – cách so chữ. Luôn bỏ qua hoa / thường và khoảng trắng (bộ nhận chữ hay dính / tách chữ).
type TextQuery struct {
	Text       string
	Exact      bool // cả dòng phải đúng chuỗi (không chỉ chứa)
	NoAccent   bool // bỏ qua dấu tiếng Việt (bộ nhận chữ đôi khi đọc sai dấu)
	Occurrence int  // lấy chỗ thứ mấy theo thứ tự đọc (trên → dưới, trái → phải); ≤ 1 = chỗ đầu
}

// textKey chuẩn hoá chữ để so: bỏ khoảng trắng, chữ thường, (tuỳ chọn) bỏ dấu. idx[i] = vị trí rune gốc của key[i].
func textKey(s string, noAccent bool) (key []rune, idx []int) {
	for i, r := range []rune(s) {
		if unicode.IsSpace(r) {
			continue
		}
		r = unicode.ToLower(r)
		if noAccent {
			if r == 'đ' {
				r = 'd'
			} else if d := []rune(norm.NFD.String(string(r))); len(d) > 0 {
				r = d[0]
			}
		}
		key = append(key, r)
		idx = append(idx, i)
	}
	return key, idx
}

func runesIndex(s, sub []rune, from int) int {
	for i := from; i+len(sub) <= len(s); i++ {
		ok := true
		for j := range sub {
			if s[i+j] != sub[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

// TextHit – chỗ thấy chữ: khung theo px ảnh và dòng chứa nó.
type TextHit struct {
	Box  model.Rect
	Line string
}

// FindText trả các chỗ thấy q.Text theo thứ tự đọc.
func FindText(lines []TextLine, q TextQuery) []TextHit {
	needle, _ := textKey(q.Text, q.NoAccent)
	if len(needle) == 0 {
		return nil
	}
	sorted := readingOrder(lines)
	var hits []TextHit
	for _, l := range sorted {
		key, idx := textKey(l.Text, q.NoAccent)
		if q.Exact {
			if string(key) == string(needle) {
				hits = append(hits, TextHit{Box: l.Box, Line: l.Text})
			}
			continue
		}
		for at := runesIndex(key, needle, 0); at >= 0; at = runesIndex(key, needle, at+len(needle)) {
			a, b := idx[at], idx[at+len(needle)-1]+1 // khoảng rune gốc [a, b)
			hits = append(hits, TextHit{Box: subBox(l, a, b), Line: l.Text})
		}
	}
	return hits
}

// ReadingOrder xếp dòng trên → dưới; các dòng cùng hàng (tâm dọc nằm trong nhau) thì trái → phải.
func readingOrder(lines []TextLine) []TextLine {
	out := append([]TextLine(nil), lines...)
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i].Box, out[j].Box
		ca, cb := a.Y+a.H/2, b.Y+b.H/2
		if cb >= a.Y && cb <= a.Y+a.H || ca >= b.Y && ca <= b.Y+b.H {
			return a.X < b.X
		}
		return ca < cb
	})
	return out
}

// subBox: khung của rune [a, b) trong dòng – hợp các từ chạm khoảng đó; không có từ thì chia theo tỷ lệ số ký tự.
func subBox(l TextLine, a, b int) model.Rect {
	runes := []rune(l.Text)
	if a <= 0 && b >= len(runes) {
		return l.Box
	}
	u16 := func(r int) int { return len(utf16.Encode(runes[:r])) }
	ua, ub := u16(a), u16(b)
	var box model.Rect
	for _, w := range l.Words {
		if w.Start < ub && w.Start+w.Len > ua {
			box = union(box, w.Box)
		}
	}
	if box.W > 0 {
		return box
	}
	n := float64(max(len(runes), 1))
	x0 := l.Box.X + int(float64(l.Box.W)*float64(a)/n)
	x1 := l.Box.X + int(float64(l.Box.W)*float64(b)/n)
	return model.Rect{X: x0, Y: l.Box.Y, W: max(x1-x0, 1), H: l.Box.H}
}

func union(a, b model.Rect) model.Rect {
	if a.W == 0 || a.H == 0 {
		return b
	}
	x0, y0 := min(a.X, b.X), min(a.Y, b.Y)
	x1, y1 := max(a.X+a.W, b.X+b.W), max(a.Y+a.H, b.Y+b.H)
	return model.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
}

// TextLines trả chữ các dòng (để báo "đọc được gì" khi không thấy), tối đa n dòng, bỏ trùng.
func TextLines(lines []TextLine, n int) []string {
	var out []string
	seen := map[string]bool{}
	for _, l := range readingOrder(lines) {
		t := strings.TrimSpace(l.Text)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
		if len(out) == n {
			break
		}
	}
	return out
}
