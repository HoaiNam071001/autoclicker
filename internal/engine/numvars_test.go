package engine

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
)

// ownOutput – output riêng của node (bỏ CommonOutputs startedAt / durationMs) để so bằng DeepEqual.
func ownOutput(r model.StepResult) map[string]any {
	out := map[string]any{}
	for k, v := range r.Output {
		if k != "startedAt" && k != "durationMs" {
			out[k] = v
		}
	}
	return out
}

// numFlow – Đặt biến n = val → node b (đọc {{vars.n}} trong numVars).
func numFlow(val string, b model.Step) model.Script {
	return model.Script{Loops: 1, Start: "v", Steps: []model.Step{
		modeltest.Node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.Value = "n", val }),
		b,
	}, Edges: []model.Edge{modeltest.Edge("v", model.PortSuccess, "b")}}
}

// Ô số nhận biến (F5c): số lần click / số nấc cuộn / thời gian chờ lấy từ biến lúc chạy; output là giá trị đã dùng.
func TestNumVarsAtRun(t *testing.T) {
	click := modeltest.Node("b", "n2", model.StepClick, func(s *model.Step) {
		s.Points, s.Clicks = []model.Point{{X: 10, Y: 20}}, 1
		s.NumVars = []model.NumVar{{Field: "clicks", Expr: "{{vars.n}} + 1"}}
	})
	in := &devicetest.FakeInput{}
	e := newEngine(in, numFlow("2", click), nil)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(strings.Join(in.Lines(), "\n"), "click left"); got != 3 {
		t.Errorf("click %d lần, muốn 3: %q", got, in.Lines())
	}

	scroll := modeltest.Node("b", "n2", model.StepScroll, func(s *model.Step) {
		s.Direction, s.Amount = "down", 1
		s.NumVars = []model.NumVar{{Field: "amount", Expr: "{{vars.n}}"}}
	})
	in = &devicetest.FakeInput{}
	e = newEngine(in, numFlow("4", scroll), nil)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if want := map[string]any{"direction": "down", "amount": 4}; !reflect.DeepEqual(ownOutput(*e.results["b"]), want) {
		t.Errorf("output cuộn %v, lệnh %q", e.results["b"].Output, in.Lines())
	}

	wait := modeltest.Node("b", "n2", model.StepWait, func(s *model.Step) {
		s.NumVars = []model.NumVar{{Field: "delayAfter", Expr: "{{vars.n}} - 5"}}
	})
	e = newEngine(&devicetest.FakeInput{}, numFlow("25", wait), nil)
	if err := e.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	r := e.results["b"]
	if r.Output["ms"] != 20 || r.Output["durationMs"] == nil || r.Output["startedAt"] == nil {
		t.Errorf("output chờ %v", r.Output)
	}
}

// Biến không ra số / số âm / ngoài khoảng → node fail invalid_input, không thao tác gì.
func TestNumVarsBadValue(t *testing.T) {
	for val, want := range map[string]string{"abc": "clicks: “abc” không phải số", "-3": "không được âm", "0": "số click phải ≥ 1"} {
		click := modeltest.Node("b", "n2", model.StepClick, func(s *model.Step) {
			s.Points, s.Clicks = []model.Point{{X: 10, Y: 20}}, 1
			s.NumVars = []model.NumVar{{Field: "clicks", Expr: "{{vars.n}}"}}
		})
		in := &devicetest.FakeInput{}
		e := newEngine(in, numFlow(val, click), nil)
		err := e.Run(context.Background())
		if err == nil || !strings.Contains(err.Error(), want) || e.results["b"].Error.Code != model.ErrInvalidInput {
			t.Errorf("%s: lỗi %v, muốn %q", val, err, want)
		}
		if len(in.Lines()) != 0 {
			t.Errorf("%s: vẫn thao tác %q", val, in.Lines())
		}
	}
}

// Output riêng của node phím / giữ chuột / cuộn tại điểm (F5c).
func TestActionOutputs(t *testing.T) {
	cases := []struct {
		name string
		node model.Step
		want map[string]any
	}{
		{"phím", modeltest.Node("a", "n1", model.StepKey, func(s *model.Step) { s.Keys, s.Modifiers = []string{"l"}, []string{"cmd", "shift"} }),
			map[string]any{"keys": "cmd+shift+l"}},
		{"giữ", modeltest.Node("a", "n1", model.StepHold, func(s *model.Step) { s.Points, s.HoldMs = []model.Point{{X: 3, Y: 4}}, 5 }),
			map[string]any{"x": 3, "y": 4, "ms": 5}},
		{"cuộn tại điểm", modeltest.Node("a", "n1", model.StepScroll, func(s *model.Step) {
			s.Points, s.Direction, s.Amount = []model.Point{{X: 7, Y: 8}}, "up", 2
		}), map[string]any{"x": 7, "y": 8, "direction": "up", "amount": 2}},
	}
	for _, tc := range cases {
		e := newEngine(&devicetest.FakeInput{}, model.Script{Loops: 1, Start: "a", Steps: []model.Step{tc.node}}, nil)
		if err := e.Run(context.Background()); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got := ownOutput(*e.results["a"]); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: output %v, muốn %v", tc.name, got, tc.want)
		}
	}
}
