package engine

import (
	"context"
	"strings"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func secretEngine(secrets map[string]string, text string) (*Engine, string, []model.StepResult) {
	sc := model.Script{Start: "v", Steps: []model.Step{
		modeltest.Node("v", "n1", model.StepSetVar, func(s *model.Step) { s.VarName, s.VarOp, s.Value = "pw", model.VarSet, text }),
		modeltest.Node("t", "n2", model.StepText, func(s *model.Step) { s.Text = "{{vars.pw}}" }),
	}, Edges: []model.Edge{modeltest.Edge("v", model.PortSuccess, "t")}}
	e, in := ctlEngine(sc, nil)
	e.rv.secrets, e.mask = secrets, newMasker(secrets)
	var shown []model.StepResult
	e.onNode = func(r model.StepResult) { shown = append(shown, r) }
	if err := e.Run(context.Background()); err != nil {
		println("run err:", err.Error())
	}
	return e, typed(in), shown
}

// Bí mật dùng được ở ô nhận biến (gõ ra đúng giá trị thật) nhưng không lọt vào kết quả node / event / Carry.
func TestSecretUsedButMasked(t *testing.T) {
	e, got, shown := secretEngine(map[string]string{"mk": "Mat-khau-123"}, "{{secret.mk}}")
	if got != "Mat-khau-123" {
		t.Fatalf("gõ ra %q, muốn giá trị thật", got)
	}
	if len(shown) != 2 {
		t.Fatalf("shown=%d", len(shown))
	}
	for _, r := range shown {
		for k, v := range r.Output {
			if s, _ := v.(string); strings.Contains(s, "Mat-khau-123") {
				t.Errorf("%s output %s lộ bí mật: %q", r.Ref, k, s)
			}
		}
	}
	if v := shown[0].Output["value"]; v != maskText {
		t.Errorf("output value=%v, muốn %q", v, maskText)
	}
	if c := e.Carry(); c.Vars["pw"] != maskText {
		t.Errorf("Carry lộ bí mật: %v", c.Vars)
	}
}

func TestSecretMissingFailsWithName(t *testing.T) {
	e, got, _ := secretEngine(nil, "{{secret.khong_co}}")
	r := e.results["v"]
	if got != "" || r == nil || r.Error == nil || r.Error.Code != model.ErrInvalidInput || !strings.Contains(r.Error.Message, "khong_co") {
		t.Errorf("got=%q r=%+v", got, r)
	}
}

func TestMaskerSkipsShortAndPrefersLonger(t *testing.T) {
	m := newMasker(map[string]string{"a": "ab", "b": "abcd", "c": "abcdef"})
	if got := m.str("x ab abcdef abcd"); got != "x ab •••• ••••" {
		t.Errorf("got %q", got)
	}
}
