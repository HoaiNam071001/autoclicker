package vision

import (
	"reflect"
	"testing"

	"autoclicker/internal/model"
)

func line(text string, x, y, w, h int, words ...TextWord) TextLine {
	return TextLine{Text: text, Box: model.Rect{X: x, Y: y, W: w, H: h}, Words: words}
}

func TestFindText(t *testing.T) {
	lines := []TextLine{
		line("Trang chủ", 10, 100, 80, 20),
		line("Bạn có thể Xem lịch sử ở đây", 10, 10, 280, 20,
			TextWord{0, 3, model.Rect{X: 10, Y: 10, W: 30, H: 20}}, TextWord{11, 3, model.Rect{X: 110, Y: 10, W: 35, H: 20}}, TextWord{15, 4, model.Rect{X: 150, Y: 10, W: 40, H: 20}},
			TextWord{20, 3, model.Rect{X: 195, Y: 10, W: 30, H: 20}}),
		line("Xem lịch sử", 300, 100, 110, 20),
		line("XEMLICH SU", 10, 200, 100, 20),
	}
	cases := []struct {
		q    TextQuery
		want []model.Rect
	}{
		// thứ tự đọc: dòng trên trước; trong dòng chỉ lấy khung các từ của chuỗi
		{TextQuery{Text: "xem lịch sử"}, []model.Rect{{X: 110, Y: 10, W: 115, H: 20}, {X: 300, Y: 100, W: 110, H: 20}}},
		{TextQuery{Text: "Xem lịch sử", Exact: true}, []model.Rect{{X: 300, Y: 100, W: 110, H: 20}}},
		{TextQuery{Text: "xem lich su", NoAccent: true}, []model.Rect{{X: 110, Y: 10, W: 115, H: 20}, {X: 300, Y: 100, W: 110, H: 20}, {X: 10, Y: 200, W: 100, H: 20}}},
		{TextQuery{Text: "trang  chu"}, nil},
		{TextQuery{Text: "Trang chủ"}, []model.Rect{{X: 10, Y: 100, W: 80, H: 20}}},
		{TextQuery{Text: "   "}, nil},
	}
	for _, c := range cases {
		var got []model.Rect
		for _, h := range FindText(lines, c.q) {
			got = append(got, h.Box)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%+v: được %v, cần %v", c.q, got, c.want)
		}
	}
	// không có khung từ → chia theo số ký tự
	if b := subBox(line("abcdefghij", 0, 0, 100, 10), 2, 4); b != (model.Rect{X: 20, Y: 0, W: 20, H: 10}) {
		t.Errorf("chia theo tỷ lệ: %v", b)
	}
}
