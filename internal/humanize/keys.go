package humanize

import (
	"strings"
	"unicode"
)

// Backspace – Key.Text của lần xoá lùi (gõ nhầm rồi sửa).
const Backspace = "\b"

// Key – chờ WaitMs rồi gõ Text (một ký tự, hoặc Backspace).
type Key struct {
	Text   string
	WaitMs int
}

// Typing – nhịp gõ chuỗi text: khoảng cách mỗi ký tự khác nhau, thỉnh thoảng ngừng sau dấu cách / dấu câu,
// mức có Typo thì đôi khi gõ nhầm phím bên cạnh, khựng lại, xoá rồi gõ đúng.
func Typing(r R, p Profile, text string) []Key {
	var out []Key
	first := true
	var prev rune
	for _, c := range text {
		wait := 0
		if !first {
			wait = Between(r, p.KeyGap)
			if (prev == ' ' || unicode.IsPunct(prev)) && r.Float64() < p.Pause {
				wait += Between(r, p.PauseMs)
			}
		}
		first, prev = false, c
		if wrong, ok := neighbor(r, c); ok && r.Float64() < p.Typo {
			out = append(out, Key{Text: string(wrong), WaitMs: wait},
				Key{Text: Backspace, WaitMs: 150 + r.IntN(250)}, Key{Text: string(c), WaitMs: Between(r, p.KeyGap)})
			continue
		}
		out = append(out, Key{Text: string(c), WaitMs: wait})
	}
	return out
}

// Hàng phím QWERTY – phím nhầm là phím ngay trái / phải trên cùng hàng.
var keyRows = []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}

func neighbor(r R, c rune) (rune, bool) {
	lower := unicode.ToLower(c)
	for _, row := range keyRows {
		i := strings.IndexRune(row, lower)
		if i < 0 {
			continue
		}
		j := i + 1
		if j >= len(row) || i > 0 && r.IntN(2) == 0 {
			j = i - 1
		}
		n := rune(row[j])
		if unicode.IsUpper(c) {
			n = unicode.ToUpper(n)
		}
		return n, true
	}
	return 0, false
}

// Gaps – khoảng chờ trước mỗi lần trong n lần (nấc cuộn, phím): lần đầu 0, sau đó trong gap, đôi khi ngừng lâu hơn.
func Gaps(r R, p Profile, gap [2]int, n int) []int {
	out := make([]int, n)
	for i := 1; i < n; i++ {
		out[i] = Between(r, gap)
		if r.Float64() < p.Pause {
			out[i] += Between(r, p.PauseMs)
		}
	}
	return out
}
