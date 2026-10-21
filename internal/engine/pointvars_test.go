package engine

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

// Toạ độ nhận biến (F5b): click lệch 30 px bên phải chỗ thấy chữ "Gửi" – toạ độ màn hình, không cộng vùng.
func TestPointFromWaitOutput(t *testing.T) {
	e, in := sendFlow(0, nil)
	e.sc.Steps[1] = modeltest.Node("c", "n2", model.StepClick, func(s *model.Step) {
		s.Points = []model.Point{{X: 5, Y: 5, Ax: "end", XVar: "{{n1.output.x}} + 30", YVar: "{{n1.output.y}} - 2"}}
	})
	e.sc.Edges = e.sc.Edges[:1] // w → c, không gõ tiếp
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"move 250,158", "click left false"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("lệnh %q, muốn %q", in.Lines(), want)
	}
	if r := e.results["c"]; r.Output["x"] != 250 || r.Output["y"] != 158 {
		t.Errorf("output click: %+v", r.Output)
	}
}

func TestPointVarNotNumber(t *testing.T) {
	sc := model.Script{Loops: 1, Start: "a", Steps: []model.Step{
		modeltest.Node("a", "n1", model.StepText, func(s *model.Step) { s.Text = "chào" }),
		modeltest.Node("b", "n2", model.StepClick, func(s *model.Step) {
			s.Points = []model.Point{{XVar: "{{n1.output.text}}", YVar: "10"}}
		}),
	}, Edges: []model.Edge{modeltest.Edge("a", model.PortSuccess, "b")}}
	in := &devicetest.FakeInput{}
	e := newEngine(in, sc, nil)
	err := e.Run(context.Background())
	if err == nil || !strings.Contains(err.Error(), "toạ độ X: “chào” không phải số") || e.results["b"].Error.Code != model.ErrInvalidInput {
		t.Fatalf("lỗi %v, kết quả %+v", err, e.results["b"].Error)
	}
	if want := []string{"type c", "type h", "type à", "type o"}; !reflect.DeepEqual(in.Lines(), want) { // không di chuột
		t.Errorf("lệnh %q", in.Lines())
	}
}

// Kéo thả: A từ biến (setVar), B lệch so với A – độ lệch cộng vào A đã tính; B dùng biến cũng là độ lệch.
func TestDragFromVarPoint(t *testing.T) {
	for _, tc := range []struct {
		name string
		b    model.Point
	}{{"B số", model.Point{X: 100, Y: -10}}, {"B biến", model.Point{XVar: "{{vars.dx}}", YVar: "-10"}}} {
		t.Run(tc.name, func(t *testing.T) {
			sc := model.Script{Loops: 1, Start: "v", Region: &model.Rect{X: 1000, Y: 1000, W: 10, H: 10}, Steps: []model.Step{
				modeltest.Node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.Value = "dx", "100" }),
				modeltest.Node("d", "n2", model.StepDrag, func(s *model.Step) {
					s.DragRelative, s.DragMs = true, 10
					s.Points = []model.Point{{XVar: "{{vars.dx}} + 200", YVar: "300"}, tc.b}
				}),
			}, Edges: []model.Edge{modeltest.Edge("v", model.PortSuccess, "d")}}
			in := &devicetest.FakeInput{}
			e := newEngine(in, sc, nil)
			if err := e.Run(context.Background()); err != nil {
				t.Fatal(err)
			}
			if want := []string{"move 300,300", "toggle left true", "drag 400,290 left", "toggle left false"}; !reflect.DeepEqual(in.Lines(), want) {
				t.Errorf("lệnh %q, muốn %q", in.Lines(), want)
			}
		})
	}
}
