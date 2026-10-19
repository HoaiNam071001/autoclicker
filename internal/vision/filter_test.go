package vision

import (
	"regexp"
	"testing"
)

func TestNormalizeNumber(t *testing.T) {
	cases := map[string]string{
		"1.234.000₫": "1234000", "Giá: 1.234.000 ₫": "1234000", "$1,234.56": "1234.56", "1.234,5 €": "1234.5",
		"12,5": "12.5", "1,234": "1234", "1.000": "1000", "3.14159": "3.14159", "0.5": "0.5", "-15 điểm": "-15",
		"−3.5": "-3.5", "1 234 567 đ": "1234567", "1\u00a0234": "1234", "Còn 12 sản phẩm, giao 2 ngày": "12", "1.": "1",
		// OCR đọc số 0 thành chữ O (OCR Windows: “1.329.OOOd”)
		"1.329.OOOd": "1329000", "2O phút": "20", "Ok, còn 5": "5", "5 Ohm": "5",
	}
	for in, want := range cases {
		if got, ok := NormalizeNumber(in); got != want || !ok {
			t.Errorf("%q → %q %v, muốn %q", in, got, ok, want)
		}
	}
	if _, ok := NormalizeNumber("hết hàng"); ok {
		t.Error("chữ không có số phải báo false")
	}
}

func TestReadFilter(t *testing.T) {
	cases := []struct {
		name, pattern string
		number        bool
		raw, want     string
		ok            bool
	}{
		{"không lọc", "", false, "  Đã gửi ", "Đã gửi", true},
		{"lấy nhóm", `Mã: (\w+)`, false, "Mã: AB12 – hạn 3 ngày", "AB12", true},
		{"cả phần khớp", `\d+ ngày`, false, "hạn 3 ngày", "3 ngày", true},
		{"không khớp", `Mã: (\w+)`, false, "Đang tải…", "", false},
		{"mẫu + số", `Tổng:(.*)`, true, "Phí 30.000₫ Tổng: 1.234.000₫", "1234000", true},
		{"số – không có số", "", true, "Liên hệ", "", false},
	}
	for _, c := range cases {
		f := ReadFilter{Number: c.number}
		if c.pattern != "" {
			f.Pattern = regexp.MustCompile(c.pattern)
		}
		if got, ok := f.Apply(c.raw); got != c.want || ok != c.ok {
			t.Errorf("%s: %q %v, muốn %q %v", c.name, got, ok, c.want, c.ok)
		}
	}
}
