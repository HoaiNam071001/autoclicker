package engine

import (
	"context"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/procinfo"
	"ghosthand/internal/script/scripttest"
)

// fakeProcs trả lần lượt từng ảnh chụp danh sách tiến trình (ảnh cuối lặp mãi).
type fakeProcs struct {
	snaps [][]procinfo.Info
	n     int
}

func (f *fakeProcs) List(context.Context) ([]procinfo.Info, error) {
	i := min(f.n, len(f.snaps)-1)
	f.n++
	return f.snaps[i], nil
}

func runProc(t *testing.T, f *fakeProcs, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepProcess, patch)
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.procs = f
	err := e.Run(t.Context())
	return e.results[e.sc.Steps[0].ID].Output, err
}

var chrome = []procinfo.Info{{PID: 10, Name: "Google Chrome", CPU: 3, MemMB: 200}, {PID: 11, Name: "Google Chrome Helper", CPU: 40.04, MemMB: 300}, {PID: 12, Name: "Finder", CPU: 1, MemMB: 50}}

func TestProcessReadSumsMatches(t *testing.T) {
	out, err := runProc(t, &fakeProcs{snaps: [][]procinfo.Info{chrome}}, func(s *model.Step) { s.ProcName = "chrome" })
	if err != nil {
		t.Fatal(err)
	}
	items := out["items"].([]any)
	if out["running"] != true || out["count"] != 2 || out["pid"] != 11 || out["cpuPercent"] != 43.0 || out["memMB"] != 500.0 || len(items) != 2 {
		t.Errorf("out=%v", out)
	}
}

func TestProcessReadNotRunningIsNotAnError(t *testing.T) {
	out, err := runProc(t, &fakeProcs{snaps: [][]procinfo.Info{chrome}}, func(s *model.Step) { s.ProcName = "photoshop" })
	if err != nil || out["running"] != false || out["count"] != 0 {
		t.Errorf("out=%v err=%v", out, err)
	}
}

func TestProcessReadAllWhenNameEmpty(t *testing.T) {
	out, _ := runProc(t, &fakeProcs{snaps: [][]procinfo.Info{chrome}}, nil)
	if out["count"] != 3 {
		t.Errorf("out=%v", out)
	}
}

func TestProcessWaitStartAndQuit(t *testing.T) {
	out, err := runProc(t, &fakeProcs{snaps: [][]procinfo.Info{nil, nil, chrome}}, func(s *model.Step) {
		s.ProcOp, s.ProcName = model.ProcWaitStart, "Finder"
		s.ProcMatch = model.MatchExact
	})
	if err != nil || out["running"] != true {
		t.Errorf("waitStart out=%v err=%v", out, err)
	}
	out, err = runProc(t, &fakeProcs{snaps: [][]procinfo.Info{chrome, nil}}, func(s *model.Step) { s.ProcOp, s.ProcName = model.ProcWaitQuit, "Finder" })
	if err != nil || out["running"] != false {
		t.Errorf("waitQuit out=%v err=%v", out, err)
	}
}

func TestProcessWaitIdleResetsWhenBusyAgain(t *testing.T) {
	busy := []procinfo.Info{{PID: 1, Name: "ffmpeg", CPU: 80}}
	quiet := []procinfo.Info{{PID: 1, Name: "ffmpeg", CPU: 1}}
	f := &fakeProcs{snaps: [][]procinfo.Info{busy, quiet, busy, quiet}}
	out, err := runProc(t, f, func(s *model.Step) {
		s.ProcOp, s.ProcName, s.CPUBelow, s.StableMs = model.ProcWaitIdle, "ffmpeg", 10, 0
	})
	if err != nil || out["cpuPercent"] != 1.0 {
		t.Errorf("out=%v err=%v", out, err)
	}
	if f.n != 2 { // StableMs 0: lần yên đầu tiên là xong, không đợi hết ảnh chụp
		t.Errorf("số lần đo = %d", f.n)
	}
}

func procErr(t *testing.T, f *fakeProcs, patch func(*model.Step)) error {
	t.Helper()
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1}), nil)
	e.procs = f
	_, err := e.process(t.Context(), modeltest.Step(model.StepProcess, patch))
	return err
}

func TestProcessWaitIdleNotRunning(t *testing.T) {
	err := procErr(t, &fakeProcs{snaps: [][]procinfo.Info{nil}}, func(s *model.Step) { s.ProcOp, s.ProcName = model.ProcWaitIdle, "ffmpeg" })
	if errCode(err) != model.ErrNotFound {
		t.Errorf("err=%v", err)
	}
}

func TestProcessWaitStartTimeout(t *testing.T) {
	err := procErr(t, &fakeProcs{snaps: [][]procinfo.Info{nil}}, func(s *model.Step) { s.ProcOp, s.ProcName, s.TimeoutMs = model.ProcWaitStart, "x", 1 })
	if errCode(err) != model.ErrTimeout {
		t.Errorf("err=%v", err)
	}
}
