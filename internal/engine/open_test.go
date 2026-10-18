package engine

import (
	"context"
	"reflect"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script"
	"autoclicker/internal/script/scripttest"
)

// Biến trong URL được mã hoá, trong tên app / văn bản thì giữ nguyên.
func TestResolveOpenTarget(t *testing.T) {
	sc := model.Script{
		Steps: []model.Step{
			{ID: "o", Type: model.StepOpen, Target: "youtube.com/results?search_query={{bai_hat}}", Enabled: true},
			{ID: "t", Type: model.StepText, Text: "{{bai_hat}}", Enabled: true},
		},
		Inputs: []model.ScriptInput{{Key: "bai_hat", Type: model.InputText}},
	}
	out, err := script.ResolveInputs(sc, map[string]any{"bai_hat": "Tinh Hà Say 2 & bạn"})
	if err != nil {
		t.Fatal(err)
	}
	if got := out.Steps[0].Target; got != "youtube.com/results?search_query=Tinh%20H%C3%A0%20Say%202%20%26%20b%E1%BA%A1n" {
		t.Errorf("URL: %q", got)
	}
	if out.Steps[1].Text != "Tinh Hà Say 2 & bạn" {
		t.Errorf("gõ chữ không mã hoá: %q", out.Steps[1].Text)
	}
}

func TestRunScrollAndOpen(t *testing.T) {
	in := &devicetest.FakeInput{}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		modeltest.Step(model.StepOpen, func(s *model.Step) { s.Target = "a.vn"; s.App = "Firefox" }),
		modeltest.Step(model.StepScroll, func(s *model.Step) { s.Amount = 2 }),
		modeltest.Step(model.StepScroll, func(s *model.Step) { s.Direction, s.Amount, s.Points = "up", 1, []model.Point{{X: 10, Y: 20}} }),
	}}
	e := newEngine(in, scripttest.Linear(sc), nil)
	var opened []string
	e.open = func(target, app string) error {
		opened = append(opened, target+" | "+app)
		return nil
	}
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"scroll down", "scroll down", "move 10,20", "scroll up"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q, muốn %q", in.Lines(), want)
	}
	if !reflect.DeepEqual(opened, []string{"a.vn | Firefox"}) {
		t.Errorf("mở: %q", opened)
	}
}
