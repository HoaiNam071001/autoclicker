package script

import (
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

func TestUnreachableNodes(t *testing.T) {
	sc := model.Script{Start: "a", Steps: []model.Step{{ID: "a"}, {ID: "b"}, {ID: "c"}}, Edges: []model.Edge{modeltest.Edge("a", model.PortFail, "c")}}
	if got := unreachableNodes(sc); !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("%q", got)
	}
}
