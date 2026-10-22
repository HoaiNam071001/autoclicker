package engine

import (
	"context"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

type fakeSysInfo struct{ path string }

func (f *fakeSysInfo) Read(_ context.Context, p string) (map[string]any, error) {
	f.path = p
	return map[string]any{"batteryPercent": 15, "charging": false, "online": true}, nil
}

func TestSysInfoOutputAndDiskPath(t *testing.T) {
	s := modeltest.Step(model.StepSysInfo, func(s *model.Step) { s.FilePath = "/Volumes/D" })
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	f := &fakeSysInfo{}
	e.sysinfo = f
	if err := e.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	out := e.results[e.sc.Steps[0].ID].Output
	if out["batteryPercent"] != 15 || out["online"] != true || f.path != "/Volumes/D" {
		t.Errorf("out=%v path=%q", out, f.path)
	}
}
