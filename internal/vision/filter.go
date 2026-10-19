package vision

import (
	"regexp"
	"strings"
)

// Lọc chữ đọc được (bước 37 E4 – node Đọc phần tử, cả khi đọc bằng OCR): lấy phần khớp mẫu regex, chuẩn hoá số.

// ReadFilter – Pattern nil = giữ cả chuỗi; Number = chuẩn hoá thành số.
type ReadFilter struct {
	Pattern *regexp.Regexp
	Number  bool
}

// Apply – chữ sau khi lọc; false = chưa có giá trị dùng được (không khớp mẫu / không có số) – node chờ đọc lại.
// Mẫu có nhóm () thì lấy nhóm đầu, không thì lấy cả phần khớp.
func (f ReadFilter) Apply(raw string) (string, bool) {
	v := raw
	if f.Pattern != nil {
		m := f.Pattern.FindStringSubmatchIndex(raw)
		if m == nil {
			return "", false
		}
		a, b := m[0], m[1]
		if len(m) >= 4 && m[2] >= 0 {
			a, b = m[2], m[3]
		}
		v = raw[a:b]
	}
	v = strings.TrimSpace(v)
	if f.Number {
		return NormalizeNumber(v)
	}
	return v, true
}

// NormalizeNumber – số đầu tiên trong chuỗi, bỏ đơn vị / dấu phân cách nghìn, dấu thập phân thành ".":
// “1.234.000₫” → 1234000, “$1,234.56” → 1234.56, “1.234,5 €” → 1234.5, “1 234” → 1234, “-15” → -15.
// Một dấu . hoặc , duy nhất: theo sau đúng 3 chữ số = phân cách nghìn (cách viết tiền Việt: “1.000” = 1000), khác
// = thập phân (“12,5” → 12.5). Có cả . và , thì dấu sau cùng là thập phân.
func NormalizeNumber(s string) (string, bool) {
	rs := []rune(s)
	i := 0
	for i < len(rs) && !digit(rs[i]) {
		i++
	}
	if i == len(rs) {
		return "", false
	}
	neg := i > 0 && (rs[i-1] == '-' || rs[i-1] == '−')
	var digits []rune
	type sep struct {
		r   rune
		pos int // số chữ số đứng trước dấu
	}
	var seps []sep
	for j := i; j < len(rs); j++ {
		r := rs[j]
		if digit(r) {
			digits = append(digits, r)
			continue
		}
		if strings.ContainsRune(".,\u00a0\u202f '’", r) && j+1 < len(rs) && digit(rs[j+1]) {
			seps = append(seps, sep{r, len(digits)})
			continue
		}
		break
	}
	dec := -1 // vị trí dấu thập phân (số chữ số phần nguyên)
	if n := len(seps); n > 0 && (seps[n-1].r == '.' || seps[n-1].r == ',') {
		last, same, other := seps[n-1], 0, false
		for _, p := range seps {
			if p.r == last.r {
				same++
			} else if p.r == '.' || p.r == ',' {
				other = true
			}
		}
		switch {
		case other && same == 1:
			dec = last.pos
		case !other && same == 1 && len(digits)-last.pos != 3:
			dec = last.pos
		}
	}
	out := string(digits)
	if dec >= 0 {
		out = string(digits[:dec]) + "." + string(digits[dec:])
	}
	if neg {
		out = "-" + out
	}
	return out, true
}

func digit(r rune) bool { return r >= '0' && r <= '9' }
