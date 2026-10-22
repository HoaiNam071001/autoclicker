package script

import (
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func TestValidateDisplay(t *testing.T) {
	for _, c := range []struct {
		op, val string
		bad     bool
	}{
		{"", "", false}, {"theme", "dark", false}, {"theme", "toggle", false}, {"theme", "blue", true}, {"theme", "", true},
		{"brightness", "70", false}, {"brightness", "101", true}, {"brightness", "abc", true}, {"brightness", "{{vars.b}}", false},
		{"audioOut", "AirPods", false}, {"audioIn", "", true}, {"sound", "Glass", false}, {"sound", "", true},
	} {
		s := modeltest.Step(model.StepDisplay, func(s *model.Step) { s.DisplayOp, s.DisplayValue = c.op, c.val })
		if err := validateDisplay(0, s); (err != nil) != c.bad {
			t.Errorf("op=%q value=%q: err=%v, bad=%v", c.op, c.val, err, c.bad)
		}
	}
}

func TestPowerDangerOnlyForSessionLosingActions(t *testing.T) {
	for _, c := range []struct {
		op        string
		dangerous bool
	}{
		{model.PowerLock, false}, {model.PowerDisplaySleep, false},
		{model.PowerSleep, true}, {model.PowerLogout, true}, {model.PowerShutdown, true}, {model.PowerRestart, true},
	} {
		s := modeltest.Step(model.StepPower, func(s *model.Step) { s.PowerOp = c.op })
		ops, fp := DangerOps(model.Script{Steps: []model.Step{s}})
		if (len(ops) > 0) != c.dangerous || (fp != "") != c.dangerous {
			t.Errorf("%s: ops=%v", c.op, ops)
		}
	}
}
