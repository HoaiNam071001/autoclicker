package main

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
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
	out, err := resolveInputs(sc, map[string]any{"bai_hat": "Tinh Hà Say 2 & bạn"})
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
		step(model.StepOpen, func(s *model.Step) { s.Target = "a.vn"; s.App = "Firefox" }),
		step(model.StepScroll, func(s *model.Step) { s.Amount = 2 }),
		step(model.StepScroll, func(s *model.Step) { s.Direction, s.Amount, s.Points = "up", 1, []model.Point{{X: 10, Y: 20}} }),
	}}
	e := newEngine(in, linear(sc), nil)
	var opened []string
	e.open = func(target, app string) error {
		opened = append(opened, target+" | "+app)
		return nil
	}
	if err := e.run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := []string{"scroll down", "scroll down", "move 10,20", "scroll up"}; !reflect.DeepEqual(in.Lines(), want) {
		t.Errorf("có %q, muốn %q", in.Lines(), want)
	}
	if !reflect.DeepEqual(opened, []string{"a.vn | Firefox"}) {
		t.Errorf("mở: %q", opened)
	}
}

func TestValidateScrollOpen(t *testing.T) {
	cases := map[string]model.Step{
		"hướng lạ":        step(model.StepScroll, func(s *model.Step) { s.Direction = "xeo" }),
		"0 nấc":           step(model.StepScroll, func(s *model.Step) { s.Amount = 0 }),
		"chưa nhập URL":   step(model.StepOpen, nil),
		"điểm ngoài vùng": step(model.StepScroll, func(s *model.Step) { s.Points = []model.Point{{X: 5000, Y: 5000}} }),
	}
	for name, s := range cases {
		sc := model.Script{Region: &model.Rect{W: 100, H: 100}, Steps: []model.Step{s}}
		if err := validateScript(linear(sc), model.Rect{}); err == nil || !strings.HasPrefix(err.Error(), "n1:") {
			t.Errorf("%s: muốn lỗi ở n1, có %v", name, err)
		}
	}
}
