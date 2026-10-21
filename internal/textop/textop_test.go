package textop

import (
	"slices"
	"testing"
	"time"

	"ghosthand/internal/model"
)

func TestMatch(t *testing.T) {
	m, ok, err := Match("Mã đơn: AB12, mã phụ: CD34", `[Mm]ã \S+: (\w+)`, true)
	if err != nil || !ok || m.Text != "AB12" || !slices.Equal(m.Items, []string{"AB12", "CD34"}) || !slices.Equal(m.Groups, []string{"AB12"}) {
		t.Errorf("all: %+v %v %v", m, ok, err)
	}
	m, _, _ = Match("a1 b2", `\d`, false)
	if m.Text != "1" || !slices.Equal(m.Items, []string{"1"}) || len(m.Groups) != 0 {
		t.Errorf("chỗ đầu, không nhóm: %+v", m)
	}
	if _, ok, err := Match("abc", `\d+`, false); ok || err != nil {
		t.Errorf("không khớp: %v %v", ok, err)
	}
	if _, _, err := Match("abc", `(`, false); err == nil {
		t.Error("regex sai phải lỗi")
	}
}

func TestSplit(t *testing.T) {
	for _, c := range []struct {
		text, sep string
		want      []string
	}{
		{"Táo\r\n Cam \n\nXoài", "", []string{"Táo", "Cam", "Xoài"}},
		{"a, b,,c ", ",", []string{"a", "b", "c"}},
		{"a\tb", `\t`, []string{"a", "b"}},
		{"  ", "", []string{}},
	} {
		if got := Split(c.text, c.sep); !slices.Equal(got, c.want) {
			t.Errorf("Split(%q, %q) = %q, muốn %q", c.text, c.sep, got, c.want)
		}
	}
}

func TestClean(t *testing.T) {
	for _, c := range []struct {
		in   string
		s    model.Step
		want string
	}{
		{"  Đà   Nẵng \n\n  Việt Nam ", model.Step{Trim: true}, "Đà Nẵng\nViệt Nam"},
		{"Đà Nẵng", model.Step{StripAccents: true, CaseMode: model.CaseLower}, "da nang"},
		{"nguyễn văn a", model.Step{CaseMode: model.CaseTitle}, "Nguyễn Văn A"},
		{"Giá: 100k", model.Step{ReplaceFrom: "k", ReplaceTo: "000"}, "Giá: 100000"},
		{"0912-345-678", model.Step{ReplaceFrom: `\D`, ReplaceRegex: true}, "0912345678"},
		{"Ngày 21/10", model.Step{ReplaceFrom: `(\d+)/(\d+)`, ReplaceTo: "$2-$1", ReplaceRegex: true}, "Ngày 10-21"},
	} {
		if got, err := Clean(c.in, c.s); err != nil || got != c.want {
			t.Errorf("Clean(%q) = %q %v, muốn %q", c.in, got, err, c.want)
		}
	}
	if _, err := Clean("x", model.Step{ReplaceFrom: "(", ReplaceRegex: true}); err == nil {
		t.Error("regex sai phải lỗi")
	}
}

func TestNumber(t *testing.T) {
	if f, ok := ExtractNumber("Giá: 1.234.000 ₫"); !ok || f != 1234000 {
		t.Errorf("lấy số = %v %v", f, ok)
	}
	if _, ok := ExtractNumber("Hết hàng"); ok {
		t.Error("không có số phải false")
	}
	for _, c := range []struct {
		f     float64
		dec   int
		style string
		want  string
	}{
		{1234567.5, 1, "", "1.234.567,5"},
		{1234567.5, 0, "", "1.234.568"},
		{1234567.5, 2, model.NumStyleEN, "1,234,567.50"},
		{1234567.5, 0, model.NumStylePlain, "1234568"},
		{1.005, 2, model.NumStyleEN, "1.01"},
		{-1500, 0, "", "-1.500"},
		{999, 0, "", "999"},
	} {
		if got := FormatNumber(c.f, c.dec, c.style); got != c.want {
			t.Errorf("FormatNumber(%v, %d, %q) = %q, muốn %q", c.f, c.dec, c.style, got, c.want)
		}
	}
}

func TestDate(t *testing.T) {
	now := time.Date(2026, 10, 21, 9, 5, 3, 0, time.Local)
	for _, c := range []struct{ in, layout, want string }{
		{"", "dd/MM/yyyy HH:mm:ss", "21/10/2026 09:05:03"},
		{"Hạn 1/2/2027 14:30", "yyyy-MM-dd HH:mm", "2027-02-01 14:30"},
		{"2026-12-31", "", "31/12/2026"},
		{"14:30", "d/M/yy H:mm", "21/10/26 14:30"},
		{"2026-11-11T20:00:00" + now.Format("-07:00"), "Hạn: dd.MM", "Hạn: 11.11"},
	} {
		d, err := ParseDate(c.in, now)
		if got := FormatDate(d, c.layout); err != nil || got != c.want {
			t.Errorf("%q theo %q = %q %v, muốn %q", c.in, c.layout, got, err, c.want)
		}
	}
	for _, bad := range []string{"Hết hàng", "31/02/2026", "25:00"} {
		if _, err := ParseDate(bad, now); err == nil {
			t.Errorf("%q phải lỗi", bad)
		}
	}
}
