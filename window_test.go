package main

import "testing"

func TestPickDisplay(t *testing.T) {
	ds := []Rect{{0, 0, 1920, 1080}, {1920, 406, 1512, 982}}
	same := []Rect{{0, 0, 1920, 1080}, {1920, 0, 1920, 1080}}
	cases := []struct {
		name  string
		ds    []Rect
		w, h  int
		mouse Point
		want  int
	}{
		{"cỡ khớp màn phụ", ds, 1512, 982, Point{X: 100, Y: 100}, 1},
		{"cỡ khớp màn chính", ds, 1920, 1080, Point{X: 2000, Y: 500}, 0},
		{"cùng cỡ → màn có chuột", same, 1920, 1080, Point{X: 2500, Y: 300}, 1},
		{"cỡ không khớp → màn có chuột", ds, 1280, 720, Point{X: 2000, Y: 500}, 1},
		{"không khớp gì → màn chính", ds, 1280, 720, Point{X: -50, Y: -50}, 0},
	}
	for _, c := range cases {
		if got := pickDisplay(c.ds, c.ds, c.w, c.h, c.mouse); got != c.want {
			t.Errorf("%s: được %d, cần %d", c.name, got, c.want)
		}
	}
}

func TestCaptureAreaDisplays(t *testing.T) {
	ds := []Rect{{0, 0, 1920, 1080}, {1920, 406, 1512, 982}}
	if err := captureArea(Rect{1920, 406, 1512, 982}, ds); err != nil {
		t.Errorf("cả màn phụ: %v", err)
	}
	if err := captureArea(Rect{1800, 500, 300, 200}, ds); err == nil {
		t.Error("vắt qua 2 màn phải báo lỗi")
	}
	if err := captureArea(Rect{100, 2000, 50, 50}, ds); err == nil {
		t.Error("ngoài màn phải báo lỗi")
	}
	e := &engine{sc: Script{Region: &Rect{2000, 500, 300, 200}}, screen: ds[0], displays: ds}
	if got := e.searchArea(Step{SearchArea: "screen"}); got != ds[1] {
		t.Errorf("cả màn với vùng ở màn phụ: %v", got)
	}
}
