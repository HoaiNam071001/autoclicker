package vision

import (
	"testing"

	"ghosthand/internal/model"
)

func rect(x, y, w, h int) model.Rect { return model.Rect{X: x, Y: y, W: w, H: h} }

// wordLine – dòng OCR có từ, mỗi từ một khung (từ cách nhau một khoảng trắng).
func wordLine(words []string, boxes []model.Rect) TextLine {
	l := TextLine{}
	pos := 0
	for i, w := range words {
		if i > 0 {
			l.Text += " "
			pos++
		}
		l.Text += w
		n := len([]rune(w))
		l.Words = append(l.Words, TextWord{Start: pos, Len: n, Box: boxes[i]})
		pos += n
		l.Box = Union(l.Box, boxes[i])
	}
	return l
}

func TestNearText(t *testing.T) {
	label := rect(100, 100, 40, 20) // “Giá:”
	lines := []TextLine{
		wordLine([]string{"Giá:", "1.234.000", "₫"}, []model.Rect{label, rect(150, 100, 90, 20), rect(245, 100, 12, 20)}),
		wordLine([]string{"Tồn", "kho"}, []model.Rect{rect(100, 140, 30, 20), rect(135, 140, 30, 20)}),
		wordLine([]string{"Xa"}, []model.Rect{rect(900, 100, 30, 20)}),
		{Text: "12 cái", Box: rect(100, 170, 60, 20)}, // không có từ
	}
	cases := []struct {
		name  string
		label model.Rect
		dir   string
		span  int
		want  string
		ok    bool
	}{
		{"bên phải – cùng dòng OCR với nhãn", label, model.ReadNearRight, 0, "1.234.000 ₫", true},
		{"bên dưới – dòng gần nhất", label, model.ReadNearBelow, 0, "Tồn kho", true},
		{"bên dưới dòng không có từ", rect(100, 140, 65, 20), model.ReadNearBelow, 0, "12 cái", true},
		{"bên phải ngoài khoảng đọc", rect(100, 140, 65, 20), model.ReadNearRight, 0, "", false},
		{"khoảng đọc ngắn bỏ chữ xa", rect(850, 100, 20, 20), model.ReadNearRight, 10, "", false},
		{"hướng sai", label, "left", 0, "", false},
	}
	for _, c := range cases {
		got, _, ok := NearText(lines, c.label, c.dir, c.span)
		if got != c.want || ok != c.ok {
			t.Errorf("%s: %q %v, muốn %q %v", c.name, got, ok, c.want, c.ok)
		}
	}
	if _, box, _ := NearText(lines, label, model.ReadNearRight, 0); box != rect(150, 100, 107, 20) {
		t.Errorf("khung %+v", box)
	}
}

func TestNearArea(t *testing.T) {
	label := rect(100, 100, 40, 20)
	if got, want := NearArea(label, model.ReadNearRight, 0, model.Rect{}), rect(100, 80, 540, 60); got != want {
		t.Errorf("phải: %+v, muốn %+v", got, want)
	}
	if got, want := NearArea(label, model.ReadNearBelow, 100, model.Rect{}), rect(60, 100, 340, 160); got != want {
		t.Errorf("dưới: %+v, muốn %+v", got, want)
	}
	if got, want := NearArea(label, model.ReadNearRight, 0, rect(0, 90, 400, 300)), rect(100, 90, 300, 50); got != want {
		t.Errorf("kẹp màn: %+v, muốn %+v", got, want)
	}
}
