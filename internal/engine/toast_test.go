package engine

import (
	"context"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

type fakeToast struct{ title, text []string }

func (f *fakeToast) Show(_ context.Context, title, text string) error {
	f.title, f.text = append(f.title, title), append(f.text, text)
	return nil
}

func runNotify(t *testing.T, f *fakeToast, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepNotify, patch)
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.toasts = f
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, err
	}
	return nil, err
}

func TestNotifyShowsTitleAndText(t *testing.T) {
	f := &fakeToast{}
	out, err := runNotify(t, f, func(s *model.Step) { s.NotifyTitle, s.Text = "Xong", "Đã nén 3 file" })
	if err != nil || len(f.title) != 1 || f.title[0] != "Xong" || f.text[0] != "Đã nén 3 file" || out["title"] != "Xong" {
		t.Errorf("out=%v err=%v got=%v/%v", out, err, f.title, f.text)
	}
}

func TestNotifyDefaultTitleAndEmptyIsInvalid(t *testing.T) {
	f := &fakeToast{}
	if _, err := runNotify(t, f, func(s *model.Step) { s.Text = "hi" }); err != nil || f.title[0] != "Ghosthand" {
		t.Errorf("tiêu đề mặc định: %v %v", f.title, err)
	}
	f = &fakeToast{}
	if _, err := runNotify(t, f, nil); err == nil || len(f.title) != 0 {
		t.Errorf("thông báo trống vẫn hiện: %v", err)
	}
}
