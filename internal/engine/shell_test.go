package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
	"ghosthand/internal/shell"
)

// fakeShell ghi lại Spec và trả kết quả định sẵn; block = chờ ctx huỷ (giả lệnh chạy mãi).
type fakeShell struct {
	got   []shell.Spec
	res   shell.Result
	err   error
	block bool
}

func (f *fakeShell) Run(ctx context.Context, sp shell.Spec) (shell.Result, error) {
	f.got = append(f.got, sp)
	if f.block {
		<-ctx.Done()
		return shell.Result{}, ctx.Err()
	}
	return f.res, f.err
}

func shellRun(t *testing.T, f *fakeShell, allow bool, patch func(*model.Step)) (map[string]any, *model.StepResult, error) {
	t.Helper()
	s := modeltest.Step(model.StepShell, func(s *model.Step) {
		s.Script = "echo hi"
		if patch != nil {
			patch(s)
		}
	})
	e := New(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), Options{AllowDanger: allow})
	e.sh = f
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, r, err
	}
	return nil, nil, err
}

func TestShellOutputAndEnv(t *testing.T) {
	f := &fakeShell{res: shell.Result{Stdout: "xong\r\n", Stderr: "\n"}}
	out, _, err := shellRun(t, f, true, func(s *model.Step) {
		s.ShellEnv = []model.KV{{Key: "TEN", Value: "Nam"}}
	})
	if err != nil || out["stdout"] != "xong" || out["stderr"] != "" || out["exitCode"] != 0 {
		t.Fatalf("output sai: %v %v", out, err)
	}
	if len(f.got) != 1 || f.got[0].Script != "echo hi" || len(f.got[0].Env) != 1 || f.got[0].Env[0] != "TEN=Nam" {
		t.Errorf("Spec sai: %+v", f.got)
	}
}

func TestShellScriptIsNotSubstituted(t *testing.T) {
	f := &fakeShell{}
	_, _, err := shellRun(t, f, true, func(s *model.Step) { s.Script = `echo "{{vars.x}}"` })
	if err != nil || f.got[0].Script != `echo "{{vars.x}}"` {
		t.Errorf("script bị thay biến: %+v %v", f.got, err)
	}
}

func TestShellNeedsConsent(t *testing.T) {
	f := &fakeShell{}
	_, r, err := shellRun(t, f, false, nil)
	if err == nil || len(f.got) != 0 {
		t.Fatalf("chưa đồng ý mà vẫn chạy: %v", f.got)
	}
	if r == nil || r.Error == nil || r.Error.Code != model.ErrPermission {
		t.Errorf("muốn lỗi permission: %+v", r)
	}
}

func TestShellExitCodeFails(t *testing.T) {
	f := &fakeShell{res: shell.Result{Stdout: "một phần", Stderr: "hỏng rồi\nchi tiết", ExitCode: 2}}
	out, r, err := shellRun(t, f, true, nil)
	if err == nil || out["exitCode"] != 2 || out["stdout"] != "một phần" {
		t.Fatalf("out=%v err=%v", out, err)
	}
	if r.Error == nil || r.Error.Code != model.ErrExitCode || !strings.Contains(r.Error.Message, "mã 2: hỏng rồi") {
		t.Errorf("lỗi sai: %+v", r.Error)
	}
}

func TestShellTimeoutKills(t *testing.T) {
	f := &fakeShell{block: true}
	start := time.Now()
	_, r, err := shellRun(t, f, true, func(s *model.Step) { s.TimeoutMs = 50 })
	if err == nil || r.Error == nil || r.Error.Code != model.ErrTimeout {
		t.Fatalf("muốn timeout: %v %+v", err, r)
	}
	if time.Since(start) > 2*time.Second {
		t.Error("hết giờ mà không dừng")
	}
}

func TestShellUnsupportedAndBadDir(t *testing.T) {
	_, r, err := shellRun(t, &fakeShell{err: shell.ErrUnsupported}, true, func(s *model.Step) { s.ShellLang = model.ShellJXA })
	if err == nil || r.Error.Code != model.ErrInvalidInput {
		t.Errorf("unsupported: %v %+v", err, r)
	}
	_, r, err = shellRun(t, &fakeShell{}, true, func(s *model.Step) { s.FilePath = "/khong/co/thu/muc/nay" })
	if err == nil || r.Error.Code != model.ErrInvalidInput || !strings.Contains(r.Error.Message, "thư mục") {
		t.Errorf("thư mục: %v %+v", err, r)
	}
}
