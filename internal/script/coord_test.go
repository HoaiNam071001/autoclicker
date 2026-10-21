package script

import (
	"testing"

	"ghosthand/internal/model"
)

func TestEvalCoord(t *testing.T) {
	ok := map[string]int{"512": 512, " 512 + 20 ": 532, "512+20-3": 529, "-40": -40, "100 - -5": 105, "10.6": 11, "640.5 + 0": 641}
	for in, want := range ok {
		if got, err := EvalCoord(in); err != nil || got != want {
			t.Errorf("EvalCoord(%q) = %d, %v; muốn %d", in, got, err, want)
		}
	}
	for _, in := range []string{"", "  ", "abc", "512 +", "512 20", "{{n2.output.x}}", "1e3", "+ 5", "5 * 2"} {
		if got, err := EvalCoord(in); err == nil {
			t.Errorf("EvalCoord(%q) = %d, muốn lỗi", in, got)
		}
	}
}

func TestPointOnScreen(t *testing.T) {
	b := model.Rect{X: 100, Y: 50, W: 1000, H: 700}
	if p, err := PointOnScreen(model.Point{X: 10, Y: 20}, b); err != nil || p != (model.Point{X: 110, Y: 70}) {
		t.Errorf("điểm số: %+v %v", p, err)
	}
	// dùng biến: toạ độ màn hình, bỏ qua x / y, neo, vùng
	p, err := PointOnScreen(model.Point{X: 10, Y: 20, Ax: AnchorEnd, XVar: "300 + 20", YVar: "40"}, b)
	if err != nil || p != (model.Point{X: 320, Y: 40}) {
		t.Errorf("điểm biến: %+v %v", p, err)
	}
	if _, err := PointOnScreen(model.Point{XVar: "300"}, b); err == nil {
		t.Error("thiếu Y phải lỗi")
	}
}
