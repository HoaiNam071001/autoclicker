package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func pathStep(f func(*model.Step)) model.Step { return modeltest.Step(model.StepPath, f) }

func TestDrawPolyline(t *testing.T) {
	in := &devicetest.FakeInput{}
	e := New(in, model.Script{Region: &model.Rect{X: 100, Y: 100, W: 500, H: 400}}, Options{})
	s := pathStep(func(s *model.Step) {
		s.Points, s.DragMs, s.Closed, s.Button = []model.Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 100, Y: 100}}, 60, true, "left"
	})
	out, err := e.drawPath(context.Background(), s)
	if err != nil {
		t.Fatal(err)
	}
	lines := in.Lines()
	if lines[0] != "move 100,100" || lines[1] != "toggle left true" || lines[len(lines)-1] != "toggle left false" {
		t.Fatalf("đầu / cuối nét: %q", lines)
	}
	for _, want := range []string{"drag 200,100 left", "drag 200,200 left", "drag 100,100 left"} { // qua từng đỉnh, khép về đầu
		if !strings.Contains(strings.Join(lines, "\n"), want) {
			t.Errorf("không đi qua %s: %q", want, lines)
		}
	}
	if out["points"] != 4 || out["x"] != 100 || out["y"] != 100 {
		t.Errorf("output %v", out)
	}
}

func TestDrawShapeAndCancel(t *testing.T) {
	in := &devicetest.FakeInput{}
	e := New(in, model.Script{}, Options{})
	star := pathStep(func(s *model.Step) {
		s.Shape, s.ShapeW, s.DragMs, s.Button, s.Points = model.ShapeStar, 200, 50, "left", []model.Point{{X: 300, Y: 300}}
	})
	out, err := e.drawPath(context.Background(), star)
	if err != nil || out["points"] != 11 || in.Lines()[0] != "move 300,200" { // đỉnh trên cùng: tâm − bán kính
		t.Fatalf("sao: %v %v %q", out, err, in.Lines()[:2])
	}
	// Dừng giữa nét → vẫn nhả nút
	in2 := &devicetest.FakeInput{}
	e2 := New(in2, model.Script{}, Options{})
	star.DragMs = 5000
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	if _, err := e2.drawPath(ctx, star); err == nil {
		t.Fatal("bị dừng mà không báo lỗi")
	}
	if l := in2.Lines(); l[len(l)-1] != "toggle left false" {
		t.Errorf("dừng giữa nét không nhả nút: %q", l[len(l)-3:])
	}
}

func TestDrawHuman(t *testing.T) {
	in := &devicetest.FakeInput{}
	e := New(in, model.Script{Human: model.HumanMedium}, Options{})
	e.setHuman(model.Step{})
	s := pathStep(func(s *model.Step) {
		s.Shape, s.ShapeW, s.ShapeH, s.DragMs, s.Button, s.Points = model.ShapeRect, 200, 100, 200, "left", []model.Point{{X: 300, Y: 300}}
	})
	if _, err := e.drawPath(context.Background(), s); err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(in.Lines(), "\n")
	for _, corner := range []string{"drag 400,250 left", "drag 400,350 left", "drag 200,350 left", "drag 200,250 left"} {
		if !strings.Contains(joined, corner) { // giống người vẫn đi đúng qua các góc
			t.Errorf("thiếu góc %s", corner)
		}
	}
}
