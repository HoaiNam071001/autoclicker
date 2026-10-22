package engine

import (
	"context"
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/inputsource"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

type fakeInputs struct {
	cur      int
	list     []inputsource.Source
	selected []string
	stuck    bool // Select không đổi được (cửa sổ không nhận)
}

func (f *fakeInputs) Current() (inputsource.Source, error) { return f.list[f.cur], nil }
func (f *fakeInputs) List() ([]inputsource.Source, error)  { return f.list, nil }
func (f *fakeInputs) Select(id string) error {
	f.selected = append(f.selected, id)
	if f.stuck {
		return nil
	}
	for i, l := range f.list {
		if l.ID == id {
			f.cur = i
		}
	}
	return nil
}

func newFakeInputs() *fakeInputs {
	return &fakeInputs{list: []inputsource.Source{
		{ID: "com.apple.keylayout.ABC", Name: "ABC", Lang: "en"},
		{ID: "com.apple.inputmethod.VietnameseIM.VietnameseTelex", Name: "Telex", Lang: "vi"},
	}}
}

func runInput(t *testing.T, f *fakeInputs, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1}), nil)
	e.inputs = f
	e.apps = &devicetest.FakeApps{Apps: []model.AppInfo{{Name: "UniKey", PID: 7}}}
	return e.inputSource(context.Background(), modeltest.Step(model.StepInputSource, patch))
}

func TestInputSourceReadChangesNothing(t *testing.T) {
	f := newFakeInputs()
	out, err := runInput(t, f, nil)
	if err != nil || out["id"] != "com.apple.keylayout.ABC" || out["vietnamese"] != false || out["changed"] != false || out["imeApp"] != "UniKey" || len(f.selected) != 0 {
		t.Errorf("out=%v err=%v selected=%v", out, err, f.selected)
	}
}

func TestInputSourceSetSwitchesAndReportsPrevious(t *testing.T) {
	f := newFakeInputs()
	out, err := runInput(t, f, func(s *model.Step) { s.InputOp, s.InputLayout = model.InputSet, "telex" })
	if err != nil || out["vietnamese"] != true || out["changed"] != true || out["previousId"] != "com.apple.keylayout.ABC" {
		t.Errorf("out=%v err=%v", out, err)
	}
	// đã đúng bộ gõ rồi thì không gọi Select lần nữa
	if _, err := runInput(t, f, func(s *model.Step) { s.InputOp, s.InputLayout = model.InputSet, "telex" }); err != nil || len(f.selected) != 1 {
		t.Errorf("selected=%v err=%v", f.selected, err)
	}
}

func TestInputSourceSetUnknownListsAvailable(t *testing.T) {
	_, err := runInput(t, newFakeInputs(), func(s *model.Step) { s.InputOp, s.InputLayout = model.InputSet, "dvorak" })
	if errCode(err) != model.ErrNotFound {
		t.Fatalf("err=%v", err)
	}
	if got := err.Error(); !strings.Contains(got, "Telex") || !strings.Contains(got, "ABC") {
		t.Errorf("thông báo thiếu danh sách bộ gõ: %s", got)
	}
}

func TestInputSourceSetStuckFails(t *testing.T) {
	f := newFakeInputs()
	f.stuck = true
	if _, err := runInput(t, f, func(s *model.Step) { s.InputOp, s.InputLayout = model.InputSet, "telex" }); err == nil {
		t.Error("bộ gõ không đổi được phải báo lỗi")
	}
}
