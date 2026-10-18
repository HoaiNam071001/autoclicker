package main

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
)

func TestResolvePointAnchors(t *testing.T) {
	r := model.Rect{X: 100, Y: 50, W: 1000, H: 700}
	cases := []struct {
		p    model.Point
		want model.Point
	}{
		{model.Point{X: 300, Y: 300}, model.Point{X: 400, Y: 350}},                                     // trống = start
		{model.Point{X: 300, Y: 300, Ax: AnchorEnd}, model.Point{X: 800, Y: 350}},                      // phải · trên
		{model.Point{X: 300, Y: 300, Ay: AnchorEnd}, model.Point{X: 400, Y: 450}},                      // trái · dưới
		{model.Point{X: 300, Y: 300, Ax: AnchorEnd, Ay: AnchorEnd}, model.Point{X: 800, Y: 450}},       // phải · dưới
		{model.Point{X: -20, Y: 10, Ax: AnchorCenter, Ay: AnchorCenter}, model.Point{X: 580, Y: 410}},  // lệch so với tâm
		{model.Point{X: 2500, Y: 5000, Ax: AnchorRatio, Ay: AnchorRatio}, model.Point{X: 350, Y: 400}}, // 25% · 50%
	}
	for _, c := range cases {
		if got := resolvePoint(c.p, r); got != c.want {
			t.Errorf("resolvePoint(%+v) = %+v, want %+v", c.p, got, c.want)
		}
	}
}

// Xong khi của bước 29: 4 điểm cách mỗi cạnh 300 px, đổi vùng 1000×700 → 1400×900 vẫn cách đúng 300 px.
func TestAnchorsFollowRegionResize(t *testing.T) {
	pts := []model.Point{
		{X: 300, Y: 300},
		{X: 300, Y: 300, Ax: AnchorEnd},
		{X: 300, Y: 300, Ay: AnchorEnd},
		{X: 300, Y: 300, Ax: AnchorEnd, Ay: AnchorEnd},
	}
	run := func(r model.Rect) []string {
		in := &devicetest.FakeInput{}
		sc := model.Script{Loops: 1, Region: &r}
		for _, p := range pts {
			sc.Steps = append(sc.Steps, step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{p} }))
		}
		if err := newEngine(in, linear(sc), nil).run(context.Background()); err != nil {
			t.Fatal(err)
		}
		var moves []string
		for _, l := range in.Lines() {
			if strings.HasPrefix(l, "move") {
				moves = append(moves, l)
			}
		}
		return moves
	}
	if got, want := run(model.Rect{W: 1000, H: 700}), []string{"move 300,300", "move 700,300", "move 300,400", "move 700,400"}; !reflect.DeepEqual(got, want) {
		t.Errorf("1000×700: got %q, want %q", got, want)
	}
	if got, want := run(model.Rect{X: 10, Y: 20, W: 1400, H: 900}), []string{"move 310,320", "move 1110,320", "move 310,620", "move 1110,620"}; !reflect.DeepEqual(got, want) {
		t.Errorf("1400×900: got %q, want %q", got, want)
	}
}

// Toàn màn hình (Region = nil) neo theo khung màn hình engine được đặt.
func TestAnchorFullScreen(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 10, Y: 10, Ax: AnchorEnd, Ay: AnchorEnd}} }),
	}}
	e := newEngine(in, linear(sc), nil)
	e.screen = model.Rect{W: 1440, H: 900}
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := in.Lines()[0]; got != "move 1430,890" {
		t.Errorf("got %q", got)
	}
}

func TestValidatePointOutsideRegion(t *testing.T) {
	r := model.Rect{W: 1000, H: 700}
	sc := model.Script{Loops: 1, Region: &r, Steps: []model.Step{
		step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{{X: 1200, Y: 10}} }),
	}}
	if err := validateScript(linear(sc), model.Rect{}); err == nil || !strings.Contains(err.Error(), "ngoài vùng") {
		t.Errorf("err = %v", err)
	}
	// B lệch so với A không bị kiểm tra theo vùng
	sc.Steps = []model.Step{step(model.StepDrag, func(s *model.Step) {
		s.Points = []model.Point{{X: 900, Y: 10}, {X: 500, Y: 0}}
		s.DragRelative = true
	})}
	if err := validateScript(linear(sc), model.Rect{}); err != nil {
		t.Errorf("err = %v", err)
	}
}
