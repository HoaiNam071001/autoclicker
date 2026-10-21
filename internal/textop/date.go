package textop

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	isoDate = regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`)
	dmyDate = regexp.MustCompile(`(\d{1,2})[/.-](\d{1,2})[/.-](\d{4})`)
	clock   = regexp.MustCompile(`(\d{1,2}):(\d{2})(?::(\d{2}))?`)
)

// ParseDate đọc ngày giờ trong chữ: ISO 8601 có múi giờ, hoặc ngày dd/mm/yyyy · yyyy-mm-dd và / hoặc giờ HH:mm[:ss]
// ở bất kỳ đâu trong chuỗi (“Hạn 21/10/2026 14:30”). Thiếu ngày = hôm nay, thiếu giờ = 00:00, trống = bây giờ.
func ParseDate(s string, now time.Time) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return now, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.In(now.Location()), nil
	}
	n := func(v string) int { i, _ := strconv.Atoi(v); return i }
	y, mo, d := now.Date()
	dated := false
	if m := isoDate.FindStringSubmatch(s); m != nil {
		y, mo, d, dated = n(m[1]), time.Month(n(m[2])), n(m[3]), true
	} else if m := dmyDate.FindStringSubmatch(s); m != nil {
		y, mo, d, dated = n(m[3]), time.Month(n(m[2])), n(m[1]), true
	}
	h, mi, sec := 0, 0, 0
	m := clock.FindStringSubmatch(s)
	if m != nil {
		h, mi, sec = n(m[1]), n(m[2]), n(m[3])
	}
	if !dated && m == nil {
		return time.Time{}, errors.New("không thấy ngày (dd/mm/yyyy, yyyy-mm-dd) hay giờ (HH:mm)")
	}
	if h > 23 || mi > 59 || sec > 59 || mo < 1 || mo > 12 || d < 1 || d > 31 {
		return time.Time{}, errors.New("ngày / giờ ngoài khoảng")
	}
	t := time.Date(y, mo, d, h, mi, sec, 0, now.Location())
	if t.Day() != d { // 31/02 → Go tự lăn sang tháng sau
		return time.Time{}, errors.New("ngày không có trong tháng")
	}
	return t, nil
}

// dateTokens – mẫu ngày theo thứ tự dài trước (yyyy trước yy, dd trước d).
var dateTokens = []struct {
	tok string
	fmt func(t time.Time) string
}{
	{"yyyy", func(t time.Time) string { return strconv.Itoa(t.Year()) }},
	{"yy", func(t time.Time) string { return pad(t.Year() % 100) }},
	{"MM", func(t time.Time) string { return pad(int(t.Month())) }},
	{"dd", func(t time.Time) string { return pad(t.Day()) }},
	{"HH", func(t time.Time) string { return pad(t.Hour()) }},
	{"mm", func(t time.Time) string { return pad(t.Minute()) }},
	{"ss", func(t time.Time) string { return pad(t.Second()) }},
	{"M", func(t time.Time) string { return strconv.Itoa(int(t.Month())) }},
	{"d", func(t time.Time) string { return strconv.Itoa(t.Day()) }},
	{"H", func(t time.Time) string { return strconv.Itoa(t.Hour()) }},
}

func pad(v int) string { return strconv.Itoa(v/10) + strconv.Itoa(v%10) }

// FormatDate viết t theo mẫu layout (“” = dd/MM/yyyy). Mẫu một chữ (d, M, H) chỉ tính khi không dính chữ cái khác –
// “Hạn: d/M” giữ nguyên chữ “Hạn”.
func FormatDate(t time.Time, layout string) string {
	if layout == "" {
		layout = "dd/MM/yyyy"
	}
	rs := []rune(layout)
	letter := func(i int) bool { return i >= 0 && i < len(rs) && unicode.IsLetter(rs[i]) }
	var b strings.Builder
next:
	for i := 0; i < len(rs); {
		for _, tk := range dateTokens {
			n := len(tk.tok)
			if i+n > len(rs) || string(rs[i:i+n]) != tk.tok {
				continue
			}
			if n == 1 && (letter(i-1) || letter(i+1)) {
				continue
			}
			b.WriteString(tk.fmt(t))
			i += n
			continue next
		}
		b.WriteRune(rs[i])
		i++
	}
	return b.String()
}
