package vision

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf16"

	"golang.org/x/text/unicode/norm"

	"ghosthand/internal/model"
)

// TextWord – một từ trong dòng: vị trí theo UTF-16 trong TextLine.Text (như Vision báo) và khung theo px ảnh.
type TextWord struct {
	Start, Len int
	Box        model.Rect
}

// TextLine – một dòng chữ đọc được, khung theo px ảnh (gốc trên trái).
type TextLine struct {
	Text   string
	Box    model.Rect
	Words  []TextWord
	Approx bool // bộ đọc không biết tiếng Việt (OCR Windows) – FindText so không dấu, chịu lệch nhẹ (fuzzy.go)
}

// TextQuery – cách so chữ. Luôn bỏ qua hoa / thường và khoảng trắng (bộ nhận chữ hay dính / tách chữ).
type TextQuery struct {
	Text       string
	Exact      bool // cả dòng phải đúng chuỗi (không chỉ chứa)
	NoAccent   bool // bỏ qua dấu tiếng Việt (bộ nhận chữ đôi khi đọc sai dấu)
	Occurrence int  // lấy chỗ thứ mấy theo thứ tự đọc (trên → dưới, trái → phải); ≤ 1 = chỗ đầu
}

// textKey chuẩn hoá chữ để so: bỏ khoảng trắng, chữ thường, (tuỳ chọn) bỏ dấu. idx[i] = vị trí rune gốc của key[i].
// Dấu tách rời (NFD – chữ chép từ macOS, một số web / app): bỏ dấu thì bỏ luôn, so đúng dấu thì ghép vào chữ trước
// để khớp chữ dựng sẵn (NFC). Không chuẩn hoá cả chuỗi trước vì idx phải trỏ đúng rune gốc (subBox).
func textKey(s string, noAccent bool) (key []rune, idx []int) {
	for i, r := range []rune(s) {
		if unicode.IsSpace(r) {
			continue
		}
		if unicode.Is(unicode.Mn, r) {
			if !noAccent && len(key) > 0 {
				if c := []rune(norm.NFC.String(string(key[len(key)-1]) + string(r))); len(c) == 1 {
					key[len(key)-1] = c[0]
				}
			}
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

// MatchText: got chứa want (exact: đúng cả chuỗi) theo cùng cách so với tìm chữ – bỏ hoa / thường, khoảng trắng,
// (tuỳ chọn) dấu. Dùng chung cho tìm phần tử màn hình (internal/element). want trống → false.
func MatchText(got, want string, exact, noAccent bool) bool {
	g, _ := textKey(got, noAccent)
	w, _ := textKey(want, noAccent)
	if len(w) == 0 {
		return false
	}
	if exact {
		return string(g) == string(w)
	}
	return runesIndex(g, w, 0) >= 0
}

// TextHit – chỗ thấy chữ: khung theo px ảnh và dòng chứa nó.
type TextHit struct {
	Box  model.Rect
	Line string
}

// FindText trả các chỗ thấy q.Text theo thứ tự đọc.
func FindText(lines []TextLine, q TextQuery) []TextHit {
	if n, _ := textKey(q.Text, q.NoAccent); len(n) == 0 {
		return nil
	}
	var hits []TextHit
	for _, l := range readingOrder(lines) {
		for _, span := range lineMatches(l, q) {
			if q.Exact {
				hits = append(hits, TextHit{Box: l.Box, Line: l.Text})
			} else {
				hits = append(hits, TextHit{Box: subBox(l, span[0], span[1]), Line: l.Text})
			}
		}
	}
	return hits
}

// lineMatches – các khoảng rune gốc [a, b) của dòng khớp q (Exact: cả dòng hoặc không gì). Dòng Approx: so không dấu,
// chịu fuzzyBudget lỗi.
func lineMatches(l TextLine, q TextQuery) [][2]int {
	noAccent := q.NoAccent || l.Approx
	needle, _ := textKey(q.Text, noAccent)
	key, idx := textKey(l.Text, noAccent)
	k := 0
	if l.Approx {
		k = fuzzyBudget(len(needle))
	}
	if q.Exact {
		if editDistance(key, needle) <= k {
			return [][2]int{{0, len([]rune(l.Text))}}
		}
		return nil
	}
	var out [][2]int
	for from := 0; from < len(key); {
		a, b := fuzzyIndex(key, needle, from, k)
		if a < 0 {
			break
		}
		if b > a {
			out = append(out, [2]int{idx[a], idx[b-1] + 1})
		}
		from = max(b, a+1)
	}
	return out
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
			box = Union(box, w.Box)
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

// Union – khung bao a và b; a rỗng → b.
func Union(a, b model.Rect) model.Rect {
	if a.W == 0 || a.H == 0 {
		return b
	}
	x0, y0 := min(a.X, b.X), min(a.Y, b.Y)
	x1, y1 := max(a.X+a.W, b.X+b.W), max(a.Y+a.H, b.Y+b.H)
	return model.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
}

// ReadLines – các dòng theo thứ tự đọc, tối đa n dòng, bỏ dòng trống / trùng chữ (báo "đọc được gì" khi không thấy).
func ReadLines(lines []TextLine, n int) []TextLine {
	var out []TextLine
	seen := map[string]bool{}
	for _, l := range readingOrder(lines) {
		t := strings.TrimSpace(l.Text)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		l.Text = t
		out = append(out, l)
		if len(out) == n {
			break
		}
	}
	return out
}

// TextLines – chữ của ReadLines.
func TextLines(lines []TextLine, n int) []string {
	var out []string
	for _, l := range ReadLines(lines, n) {
		out = append(out, l.Text)
	}
	return out
}
