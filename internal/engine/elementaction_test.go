package engine

import (
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/element"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

func actionRun(t *testing.T, fake *devicetest.FakeElements, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepElementAction, func(s *model.Step) {
		s.Selectors, s.TimeoutMs, s.IntervalMs = []string{"text=Nhớ tôi"}, 40, 1
		patch(s)
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.elems = fake
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, err
	}
	return nil, err
}

func box(role, title, value string) element.Node {
	return element.Node{Role: role, Title: title, Value: value, Box: elemBox(0, 0, 100, 20)}
}

func TestElementActionPressAndCheck(t *testing.T) {
	win := element.Node{Role: "window", Title: "App", Box: elemBox(0, 0, 300, 200)}
	fake := &devicetest.FakeElements{Frames: [][]element.Node{
		{win, box("checkbox", "Nhớ tôi", "0")}, // trước
		{win, box("checkbox", "Nhớ tôi", "1")}, // sau thao tác
	}}
	out, err := actionRun(t, fake, func(s *model.Step) { s.ElementAct = model.ElActCheck })
	if err != nil || out["value"] != "1" || len(fake.Acts) != 1 || fake.Acts[0] != "1:press:" {
		t.Fatalf("out %v err %v acts %v", out, err, fake.Acts)
	}
	// đã tick sẵn → không bấm
	fake = &devicetest.FakeElements{Frames: [][]element.Node{{win, box("checkbox", "Nhớ tôi", "1")}}}
	if _, err := actionRun(t, fake, func(s *model.Step) { s.ElementAct = model.ElActCheck }); err != nil || len(fake.Acts) != 0 {
		t.Fatalf("đã tick mà vẫn bấm: %v %v", err, fake.Acts)
	}
}

func TestElementActionStateNotChangedFails(t *testing.T) {
	win := element.Node{Role: "window", Box: elemBox(0, 0, 300, 200)}
	fake := &devicetest.FakeElements{Frames: [][]element.Node{{win, box("checkbox", "Nhớ tôi", "0")}}}
	_, err := actionRun(t, fake, func(s *model.Step) { s.ElementAct = model.ElActCheck })
	if err == nil || !strings.Contains(err.Error(), "trạng thái tick không đổi") {
		t.Fatalf("muốn lỗi trạng thái không đổi: %v", err)
	}
}

func TestElementActionSetValue(t *testing.T) {
	win := element.Node{Role: "window", Box: elemBox(0, 0, 300, 200)}
	field := func(v string) element.Node {
		n := box("textbox", "", v)
		n.Placeholder = "Nhớ tôi"
		return n
	}
	fake := &devicetest.FakeElements{Frames: [][]element.Node{{win, field("")}, {win, field("abc")}}}
	out, err := actionRun(t, fake, func(s *model.Step) { s.ElementAct, s.Text = model.ElActSetValue, "abc" })
	if err != nil || out["value"] != "abc" || fake.Acts[0] != "1:setValue:abc" {
		t.Fatalf("out %v err %v acts %v", out, err, fake.Acts)
	}
	fake = &devicetest.FakeElements{Frames: [][]element.Node{{win, field("")}}}
	_, err = actionRun(t, fake, func(s *model.Step) { s.ElementAct, s.Text = model.ElActSetValue, "abc" })
	if err == nil || !strings.Contains(err.Error(), "giá trị vẫn là") {
		t.Fatalf("ô không nhận phải fail: %v", err)
	}
}

func TestElementActionPressDefaultAndAppRefusal(t *testing.T) {
	win := element.Node{Role: "window", Box: elemBox(0, 0, 300, 200)}
	fake := &devicetest.FakeElements{Frames: [][]element.Node{{win, box("button", "Nhớ tôi", "")}}}
	out, err := actionRun(t, fake, func(*model.Step) {})
	if err != nil || out["action"] != "press" || fake.Acts[0] != "1:press:" {
		t.Fatalf("out %v err %v acts %v", out, err, fake.Acts)
	}
	fake = &devicetest.FakeElements{Frames: [][]element.Node{{win, box("button", "Nhớ tôi", "")}}, ActErr: errTest("app không nhận")}
	if _, err := actionRun(t, fake, func(*model.Step) {}); err == nil || !strings.Contains(err.Error(), "app không nhận") {
		t.Fatalf("lỗi của app phải được báo: %v", err)
	}
}

type errTest string

func (e errTest) Error() string { return string(e) }
