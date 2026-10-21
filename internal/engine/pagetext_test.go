package engine

import (
	"reflect"
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/element"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

func pageTextRun(t *testing.T, nodes []element.Node, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepPageText, func(s *model.Step) {
		s.TimeoutMs, s.IntervalMs = 20, 1
		patch(s)
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{nodes}}
	err := e.Run(t.Context())
	return e.results[e.sc.Steps[0].ID].Output, err
}

func pageNodes() []element.Node {
	b := elemBox(0, 0, 200, 20)
	return []element.Node{
		{Role: "window", Title: "Trang", Box: b},
		{Role: "heading", Title: "Giá vé", Box: b, Depth: 1},
		{Role: "group", Desc: "danh sách", Box: b, Depth: 1},
		{Role: "text", Value: "Hà Nội – Huế", Box: b, Depth: 2},
		{Role: "text", Value: "1.200.000₫", Box: b, Depth: 2},
		{Role: "text", Value: "cuối trang", Box: elemBox(0, 5000, 200, 0), Depth: 1},
	}
}

func TestPageTextWholeWindow(t *testing.T) {
	out, err := pageTextRun(t, pageNodes(), func(*model.Step) {})
	if err != nil {
		t.Fatal(err)
	}
	if out["text"] != "Giá vé\nHà Nội – Huế\n1.200.000₫" || out["count"] != 3 {
		t.Errorf("output %v", out)
	}
	if !reflect.DeepEqual(out["lines"], []any{"Giá vé", "Hà Nội – Huế", "1.200.000₫"}) {
		t.Errorf("lines %v", out["lines"])
	}
	out, _ = pageTextRun(t, pageNodes(), func(s *model.Step) { s.PageAll = true })
	if !strings.HasSuffix(out["text"].(string), "cuối trang") {
		t.Errorf("pageAll: %v", out["text"])
	}
}

func TestPageTextInsideElement(t *testing.T) {
	out, err := pageTextRun(t, pageNodes(), func(s *model.Step) { s.Selectors = []string{`role=group[name="danh sách"]`} })
	if err != nil {
		t.Fatal(err)
	}
	if out["text"] != "Hà Nội – Huế\n1.200.000₫" {
		t.Errorf("output %v", out)
	}
}

func TestPageTextMissing(t *testing.T) {
	_, err := pageTextRun(t, pageNodes(), func(s *model.Step) { s.Selectors = []string{"text=không có"} })
	if err == nil || !strings.Contains(err.Error(), "không đọc được chữ trang") {
		t.Errorf("lỗi %v", err)
	}
	_, err = pageTextRun(t, []element.Node{{Role: "window", Box: elemBox(0, 0, 10, 10)}}, func(*model.Step) {})
	if err == nil || !strings.Contains(err.Error(), "không có chữ nào") {
		t.Errorf("trang trống: %v", err)
	}
}
