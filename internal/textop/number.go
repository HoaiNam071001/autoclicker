package textop

import (
	"math"
	"strconv"
	"strings"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// ExtractNumber – số đầu tiên trong chữ kiểu Việt / Anh (cùng công thức Lấy số của Đọc phần tử, 37 E4):
// “1.234.000 đ” → 1234000. false = không có số.
func ExtractNumber(s string) (float64, bool) {
	v, ok := vision.NormalizeNumber(s)
	if !ok {
		return 0, false
	}
	f, err := strconv.ParseFloat(v, 64)
	return f, err == nil
}

// FormatNumber – f làm tròn nửa lên decimals chữ số, viết theo style (NumStyles): 1234567.5 → “1.234.567,5”.
func FormatNumber(f float64, decimals int, style string) string {
	decimals = max(0, min(decimals, 10))
	p := math.Pow(10, float64(decimals))
	// Bỏ nhiễu nhị phân trước khi làm tròn: 1.005 × 100 = 100.49999… phải thành 101.
	f = math.Round(math.Round(f*p*1e6)/1e6) / p
	s := strconv.FormatFloat(math.Abs(f), 'f', decimals, 64)
	whole, frac, _ := strings.Cut(s, ".")
	group, point := ".", ","
	switch style {
	case model.NumStyleEN:
		group, point = ",", "."
	case model.NumStylePlain:
		group, point = "", "."
	}
	var b strings.Builder
	if f < 0 {
		b.WriteByte('-')
	}
	for i, r := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			b.WriteString(group)
		}
		b.WriteRune(r)
	}
	if frac != "" {
		b.WriteString(point + frac)
	}
	return b.String()
}
