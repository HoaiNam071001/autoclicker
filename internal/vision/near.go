package vision

import (
	"strings"
	"unicode/utf16"

	"autoclicker/internal/model"
)

// Đọc chữ cạnh nhãn (bước 37 E4 – dự phòng khi giá trị không có trong cây phần tử: canvas, app tự vẽ): tìm nhãn
// (phần tử / OCR) rồi lấy dòng chữ gần nhất bên phải / bên dưới nhãn trong các dòng OCR. Nhãn và dòng cùng hệ toạ độ.

// DefaultSpan – khoảng đọc mặc định (px): bên phải = giá trị phải bắt đầu trong chừng này tính từ mép phải nhãn;
// bên dưới = dòng giá trị phải bắt đầu trong chừng này tính từ đáy nhãn.
const DefaultSpan = 300

// unit – một từ (hoặc cả dòng khi bộ đọc không báo từ) kèm vị trí UTF-16 trong dòng.
type unit struct {
	start, end int
	box        model.Rect
}

func units(l TextLine) []unit {
	if len(l.Words) == 0 {
		return []unit{{0, len(utf16.Encode([]rune(l.Text))), l.Box}}
	}
	out := make([]unit, len(l.Words))
	for i, w := range l.Words {
		out[i] = unit{w.Start, w.Start + w.Len, w.Box}
	}
	return out
}

// NearText – chữ gần nhãn nhất theo hướng dir. Bên phải: các từ bắt đầu sau mép phải nhãn (lệch ≤ 2 px), tâm dọc nằm
// trong dải nhãn (± nửa chiều cao) – dòng có từ gần nhãn nhất, lấy từ từ đó tới hết dòng (dòng OCR “Giá: 1.234.000 ₫”
// chứa cả nhãn thì chỉ lấy phần sau nhãn). Bên dưới: dòng đầu tiên có đỉnh dưới tâm nhãn, chồng ngang với nhãn (nới
// sang trái 2 lần chiều cao, sang phải tới span). span ≤ 0 = DefaultSpan.
func NearText(lines []TextLine, label model.Rect, dir string, span int) (string, model.Rect, bool) {
	if span <= 0 {
		span = DefaultSpan
	}
	right, bottom := label.X+label.W, label.Y+label.H
	in := func(b model.Rect) bool {
		cy := b.Y + b.H/2
		if dir == model.ReadNearRight {
			return b.X >= right-2 && b.X <= right+span && cy >= label.Y-label.H/2 && cy <= bottom+label.H/2
		}
		return b.Y >= label.Y+label.H/2 && b.Y <= bottom+span && b.X+b.W >= label.X-2*label.H && b.X <= label.X+max(label.W, span)
	}
	if dir != model.ReadNearRight && dir != model.ReadNearBelow {
		return "", model.Rect{}, false
	}
	best, bestD := -1, 0
	var pick []unit
	for i, l := range lines {
		var us []unit
		for _, u := range units(l) {
			if in(u.box) {
				us = append(us, u)
			}
		}
		if len(us) == 0 {
			continue
		}
		d := us[0].box.X - right // bên phải: từ đầu gần mép nhãn nhất
		if dir == model.ReadNearBelow {
			d = us[0].box.Y - bottom
		}
		if best < 0 || d < bestD {
			best, bestD, pick = i, d, us
		}
	}
	if best < 0 {
		return "", model.Rect{}, false
	}
	u16 := utf16.Encode([]rune(lines[best].Text))
	a, b := pick[0].start, pick[len(pick)-1].end
	var box model.Rect
	for _, u := range pick {
		box = Union(box, u.box)
	}
	text := strings.TrimSpace(strings.TrimLeft(string(utf16.Decode(u16[min(a, len(u16)):min(b, len(u16))])), " :："))
	return text, box, text != ""
}

// NearArea – vùng cần chụp để đọc cạnh nhãn (nhãn lấy từ cây phần tử, chưa có OCR): gồm cả nhãn (bộ đọc thấy đủ ngữ
// cảnh của dòng), rộng thêm 200 px để không cắt ngang chữ bắt đầu ở cuối khoảng đọc; kẹp trong khung màn bound
// (zero = không kẹp).
func NearArea(label model.Rect, dir string, span int, bound model.Rect) model.Rect {
	if span <= 0 {
		span = DefaultSpan
	}
	var r model.Rect
	if dir == model.ReadNearBelow {
		x0 := label.X - 2*label.H
		r = model.Rect{X: x0, Y: label.Y, W: label.X + max(label.W, span) + 200 - x0, H: label.H + span + 2*label.H}
	} else {
		r = model.Rect{X: label.X, Y: label.Y - label.H, W: label.W + span + 200, H: 3 * label.H}
	}
	if bound.W == 0 || bound.H == 0 {
		return r
	}
	x0, y0 := max(r.X, bound.X), max(r.Y, bound.Y)
	x1, y1 := min(r.X+r.W, bound.X+bound.W), min(r.Y+r.H, bound.Y+bound.H)
	return model.Rect{X: x0, Y: y0, W: max(x1-x0, 0), H: max(y1-y0, 0)}
}
