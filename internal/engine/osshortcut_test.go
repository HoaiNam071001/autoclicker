package engine

import (
	"context"
	"testing"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
	"ghosthand/internal/shortcuts"
)

type fakeShortcuts struct {
	name, input string
	out         string
	err         error
	block       bool
}

func (f *fakeShortcuts) List(context.Context) ([]string, error) { return nil, nil }
func (f *fakeShortcuts) Run(ctx context.Context, name, input string) (string, error) {
	f.name, f.input = name, input
	if f.block {
		<-ctx.Done()
		return "", ctx.Err()
	}
	return f.out, f.err
}

func shortcutRun(t *testing.T, f *fakeShortcuts, patch func(*model.Step)) (map[string]any, *model.StepResult, error) {
	t.Helper()
	s := modeltest.Step(model.StepOSShortcut, func(s *model.Step) {
		s.ShortcutName = "Đổi tên"
		if patch != nil {
			patch(s)
		}
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.shortcuts = f
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, r, err
	}
	return nil, nil, err
}

func TestOSShortcutPassesInputAndReturnsText(t *testing.T) {
	f := &fakeShortcuts{out: "xong\n"}
	out, _, err := shortcutRun(t, f, func(s *model.Step) { s.Text = "a.png" })
	if err != nil || out["text"] != "xong" || f.name != "Đổi tên" || f.input != "a.png" {
		t.Errorf("out=%v err=%v fake=%+v", out, err, f)
	}
}

func TestOSShortcutErrors(t *testing.T) {
	_, r, err := shortcutRun(t, &fakeShortcuts{err: shortcuts.ErrNotFound}, nil)
	if err == nil || r.Error.Code != model.ErrNotFound {
		t.Errorf("not found: %v %+v", err, r)
	}
	_, r, err = shortcutRun(t, &fakeShortcuts{err: shortcuts.ErrUnsupported}, nil)
	if err == nil || r.Error.Code != model.ErrInvalidInput {
		t.Errorf("unsupported: %v %+v", err, r)
	}
	start := time.Now()
	_, r, err = shortcutRun(t, &fakeShortcuts{block: true}, func(s *model.Step) { s.TimeoutMs = 50 })
	if err == nil || r.Error.Code != model.ErrTimeout || time.Since(start) > 2*time.Second {
		t.Errorf("timeout: %v %+v", err, r)
	}
}
