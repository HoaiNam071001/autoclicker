package engine

import (
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/element"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

// uiRun chạy một node waitUIEvent với thiết bị giả; patch chỉnh bước.
func uiRun(t *testing.T, fake *devicetest.FakeElements, wins *seqWindows, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepWaitUIEvent, func(s *model.Step) {
		s.TimeoutMs = 600
		patch(s)
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	if fake != nil {
		e.elems = fake
	}
	if wins != nil {
		e.wins = wins
	}
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, err
	}
	return nil, err
}

// seqWindows – mỗi lần List trả một khung kế tiếp (khung cuối giữ mãi).
type seqWindows struct {
	devicetest.FakeWindows
	frames [][]model.WindowInfo
	n      int
}

func (f *seqWindows) List() ([]model.WindowInfo, error) {
	fr := f.frames[min(f.n, len(f.frames)-1)]
	f.n++
	return fr, nil
}

func win(id uint64, app, title string) model.WindowInfo {
	return model.WindowInfo{ID: id, App: app, Title: title, Box: elemBox(0, 0, 400, 300)}
}

func TestWaitUIEventNewWindow(t *testing.T) {
	a, dlg := win(1, "Notes", "Ghi chú"), win(2, "Notes", "Lưu thành")
	wins := &seqWindows{frames: [][]model.WindowInfo{{a}, {a}, {dlg, a}}}
	out, err := uiRun(t, nil, wins, func(s *model.Step) { s.UIEvent, s.WinTitle = model.UIEventNewWindow, "Lưu" })
	if err != nil || out["title"] != "Lưu thành" || out["event"] != model.UIEventNewWindow {
		t.Fatalf("out %v err %v", out, err)
	}
	// cửa sổ có sẵn lúc đầu không tính là mới
	wins = &seqWindows{frames: [][]model.WindowInfo{{dlg, a}}}
	_, err = uiRun(t, nil, wins, func(s *model.Step) { s.TimeoutMs = 50; s.WinTitle = "Lưu" })
	if err == nil || !strings.Contains(err.Error(), "không thấy cửa sổ mới") {
		t.Fatalf("muốn hết giờ: %v", err)
	}
}

func TestWaitUIEventFocus(t *testing.T) {
	a, b := win(1, "Notes", "Ghi chú"), win(2, "Mail", "Hộp thư")
	wins := &seqWindows{frames: [][]model.WindowInfo{{a, b}, {a, b}, {b, a}}}
	out, err := uiRun(t, nil, wins, func(s *model.Step) { s.UIEvent, s.WinApp = model.UIEventFocus, "Mail" })
	if err != nil || out["app"] != "Mail" {
		t.Fatalf("out %v err %v", out, err)
	}
}

func menuPage(extra ...element.Node) []element.Node {
	return append([]element.Node{{Role: "window", Title: "App", Box: elemBox(0, 0, 300, 200)}, box("button", "Mở", "")}, extra...)
}

func TestWaitUIEventNewElement(t *testing.T) {
	fake := &devicetest.FakeElements{Frames: [][]element.Node{
		menuPage(), menuPage(), menuPage(box("menu", "Tuỳ chọn", "")),
	}}
	out, err := uiRun(t, fake, nil, func(s *model.Step) { s.UIEvent, s.Selectors = model.UIEventNewElement, []string{"role=menu"} })
	if err != nil || out["role"] != "menu" || out["text"] != "Tuỳ chọn" {
		t.Fatalf("out %v err %v", out, err)
	}
	// phần tử có sẵn lúc đầu không tính
	fake = &devicetest.FakeElements{Frames: [][]element.Node{menuPage(box("menu", "Tuỳ chọn", ""))}}
	_, err = uiRun(t, fake, nil, func(s *model.Step) {
		s.TimeoutMs, s.UIEvent, s.Selectors = 50, model.UIEventNewElement, []string{"role=menu"}
	})
	if err == nil || !strings.Contains(err.Error(), "không thấy phần tử mới") {
		t.Fatalf("muốn hết giờ: %v", err)
	}
}

func TestWaitUIEventValueChange(t *testing.T) {
	fake := &devicetest.FakeElements{Frames: [][]element.Node{
		menuPage(box("textbox", "Tiến độ", "10%")), menuPage(box("textbox", "Tiến độ", "10%")), menuPage(box("textbox", "Tiến độ", "100%")),
	}}
	out, err := uiRun(t, fake, nil, func(s *model.Step) { s.UIEvent, s.Selectors = model.UIEventValueChange, []string{"role=textbox"} })
	if err != nil || out["value"] != "100%" {
		t.Fatalf("out %v err %v", out, err)
	}
}
