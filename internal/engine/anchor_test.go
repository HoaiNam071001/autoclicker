package engine

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script"
	"autoclicker/internal/script/scripttest"
)

// Xong khi của bước 29: 4 điểm cách mỗi cạnh 300 px, đổi vùng 1000×700 → 1400×900 vẫn cách đúng 300 px.
func TestAnchorsFollowRegionResize(t *testing.T) {
	pts := []model.Point{
		{X: 300, Y: 300},
		{X: 300, Y: 300, Ax: script.AnchorEnd},
		{X: 300, Y: 300, Ay: script.AnchorEnd},
		{X: 300, Y: 300, Ax: script.AnchorEnd, Ay: script.AnchorEnd},
	}
	run := func(r model.Rect) []string {
		in := &devicetest.FakeInput{}
		sc := model.Script{Loops: 1, Region: &r}
		for _, p := range pts {
			sc.Steps = append(sc.Steps, modeltest.Step(model.StepClick, func(s *model.Step) { s.Points = []model.Point{p} }))
		}
		if err := newEngine(in, scripttest.Linear(sc), nil).Run(context.Background()); err != nil {
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
		modeltest.Step(model.StepClick, func(s *model.Step) {
			s.Points = []model.Point{{X: 10, Y: 10, Ax: script.AnchorEnd, Ay: script.AnchorEnd}}
		}),
	}}
	e := newEngine(in, scripttest.Linear(sc), nil)
	e.screen = model.Rect{W: 1440, H: 900}
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := in.Lines()[0]; got != "move 1430,890" {
		t.Errorf("got %q", got)
	}
}
