package script

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

func TestValidateFlow(t *testing.T) {
	base := func() model.Script {
		return model.Script{Start: "a", Steps: []model.Step{modeltest.Node("a", "n1", model.StepWait, nil), modeltest.Node("b", "n2", model.StepWait, nil)},
			Edges: []model.Edge{modeltest.Edge("a", model.PortSuccess, "b")}}
	}
	cases := map[string]func(*model.Script){
		"chưa chọn node bắt đầu": func(s *model.Script) { s.Start = "x" },
		"trùng với node khác":    func(s *model.Script) { s.Steps[1].Ref = "n1" },
		"không hợp lệ":           func(s *model.Script) { s.Steps[1].Ref = "vars" },
		"không có cổng":          func(s *model.Script) { s.Edges[0].Port = "true" },
		"nối tới node không còn": func(s *model.Script) { s.Edges[0].To = "x" },
		"nối vào chính node này": func(s *model.Script) { s.Edges[0].To = "a" },
		"có hơn một cạnh": func(s *model.Script) {
			s.Edges = append(s.Edges, model.Edge{ID: "z", From: "a", Port: model.PortSuccess, To: "b"})
		},
		"không có trong kịch bản": func(s *model.Script) {
			s.Steps[1] = modeltest.Node("b", "n2", model.StepText, func(x *model.Step) { x.Text = "{{n9.output.x}}" })
		},
		"chưa chọn node lấy chỗ": func(s *model.Script) { s.Steps[1] = modeltest.Node("b", "n2", model.StepClickImage, nil) },
		"không tìm ảnh / chữ": func(s *model.Script) {
			s.Steps[1] = modeltest.Node("b", "n2", model.StepClickImage, func(x *model.Step) { x.ImageFrom = "a" })
		},
		"kịch bản chưa có node nào": func(s *model.Script) { s.Steps = nil },
	}
	for want, mod := range cases {
		sc := base()
		mod(&sc)
		if err := validateFlow(sc); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi %q, có %v", want, err)
		}
	}
	if err := validateFlow(base()); err != nil {
		t.Errorf("flow hợp lệ báo lỗi: %v", err)
	}
	// biến vars / last không cần node
	sc := base()
	sc.Steps[1] = modeltest.Node("b", "n2", model.StepText, func(x *model.Step) { x.Text = "{{last.status}} {{vars.dem}}" })
	if err := validateFlow(sc); err != nil {
		t.Errorf("vars / last: %v", err)
	}
}

// Lỗi cạnh có id cạnh + node (AI sửa đúng cạnh – 20 F7).
func TestValidateFlowEdgeID(t *testing.T) {
	sc := model.Script{Start: "a", Steps: []model.Step{modeltest.Node("a", "n1", model.StepWait, nil)},
		Edges: []model.Edge{{ID: "e1", From: "a", Port: model.PortSuccess, To: "x"}}}
	var se *model.ScriptError
	if err := validateFlow(sc); !errors.As(err, &se) || se.Edge != "e1" || se.StepID != "a" || se.Ref != "n1" || se.Field != "edges" {
		t.Errorf("lỗi cạnh = %#v", err)
	}
	sc.Edges[0].From = "y"
	if err := validateFlow(sc); !errors.As(err, &se) || se.Edge != "e1" {
		t.Errorf("cạnh từ node không còn = %#v", err)
	}
}

func TestUnreachableNodes(t *testing.T) {
	sc := model.Script{Start: "a", Steps: []model.Step{{ID: "a"}, {ID: "b"}, {ID: "c"}}, Edges: []model.Edge{modeltest.Edge("a", model.PortFail, "c")}}
	if got := unreachableNodes(sc); !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("%q", got)
	}
}

// Ô toạ độ dùng biến (F5b): biến phải là output có thật của node nguồn, cả ô dạng số + / − px.
func TestValidateFlowCoordVars(t *testing.T) {
	flow := func(x, y string) model.Script {
		return model.Script{Start: "a", Steps: []model.Step{
			modeltest.Node("a", "n1", model.StepWaitImage, func(s *model.Step) { s.FindBy, s.FindText = model.FindByText, "Gửi" }),
			modeltest.Node("t", "n2", model.StepText, func(s *model.Step) { s.Text = "x" }),
			modeltest.Node("b", "n3", model.StepClick, func(s *model.Step) { s.Points = []model.Point{{XVar: x, YVar: y}} }),
		}, Edges: []model.Edge{modeltest.Edge("a", model.PortSuccess, "b")}}
	}
	for _, ok := range [][2]string{{"{{n1.output.x}} + 20", "{{n1.output.y}}"}, {"{{vars.x}}", "{{last.output.y}} - 4"}, {"{{diem}}", "300"}} {
		if err := validateFlow(flow(ok[0], ok[1])); err != nil {
			t.Errorf("%q: %v", ok, err)
		}
	}
	bad := map[string][2]string{
		`không có output "x" – có: text`: {"{{n2.output.x}}", "1"},
		"không phải số – ô số dùng":      {"{{n1.status}}", "1"},
		"“{{n1.output.y}} +” sai dạng":   {"1", "{{n1.output.y}} +"},
		"chưa nhập toạ độ Y":             {"{{n1.output.x}}", ""},
		"không có trong kịch bản":        {"{{n9.output.x}}", "1"},
	}
	for want, xy := range bad {
		if err := validateFlow(flow(xy[0], xy[1])); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: muốn lỗi %q, có %v", xy, want, err)
		}
	}
}

func TestVarFieldsSkipUnusedDragA(t *testing.T) {
	s := modeltest.Node("d", "n1", model.StepDrag, func(s *model.Step) {
		s.DragFromPrev = true
		s.Points = []model.Point{{XVar: "{{n9.output.x}}", YVar: "1"}, {XVar: "2", YVar: "3"}}
	})
	var names []string
	for _, f := range VarFields(&s) {
		names = append(names, f.Name)
	}
	if !reflect.DeepEqual(names, []string{"points[1].xVar", "points[1].yVar"}) {
		t.Errorf("%q", names)
	}
}
