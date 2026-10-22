package engine

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/element"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script/scripttest"
)

func tableRun(t *testing.T, nodes []element.Node, patch func(*model.Step)) (map[string]any, error) {
	t.Helper()
	s := modeltest.Step(model.StepReadTable, func(s *model.Step) {
		s.TimeoutMs, s.IntervalMs = 20, 1
		patch(s)
	})
	e := newEngine(nil, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{s}}), nil)
	e.elems = &devicetest.FakeElements{Frames: [][]element.Node{nodes}}
	err := e.Run(t.Context())
	if r := e.results[e.sc.Steps[0].ID]; r != nil {
		return r.Output, err
	}
	return nil, err
}

func tableNodes() []element.Node {
	b := elemBox(0, 0, 200, 20)
	cell := func(role, v string, d int) element.Node {
		return element.Node{Role: role, Value: v, Title: v, Box: b, Depth: d}
	}
	return []element.Node{
		{Role: "window", Title: "Trang", Box: b},
		{Role: "table", Desc: "Giá", Box: b, Depth: 1},
		{Role: "row", Box: b, Depth: 2}, cell("columnheader", "Món", 3), cell("columnheader", "Giá vé", 3),
		{Role: "row", Box: b, Depth: 2}, cell("cell", "Phở", 3), cell("cell", "50.000₫", 3),
		{Role: "row", Box: b, Depth: 2}, cell("cell", "Bún", 3), cell("cell", "45.000₫", 3),
		{Role: "row", Box: elemBox(0, 5000, 200, 0), Depth: 2}, cell("cell", "Ẩn", 3), cell("cell", "1", 3),
	}
}

func TestReadTable(t *testing.T) {
	out, err := tableRun(t, tableNodes(), func(*model.Step) {})
	if err != nil {
		t.Fatal(err)
	}
	want := []any{[]any{"Món", "Giá vé"}, []any{"Phở", "50.000₫"}, []any{"Bún", "45.000₫"}}
	if !reflect.DeepEqual(out["rows"], want) || out["count"] != 2 {
		t.Errorf("output %v", out)
	}
	items := out["items"].([]any)
	if items[1].(map[string]any)["gia_ve"] != "45.000₫" || !reflect.DeepEqual(out["columns"], []any{"Món", "Giá vé"}) {
		t.Errorf("items %v columns %v", items, out["columns"])
	}
	out, _ = tableRun(t, tableNodes(), func(s *model.Step) { s.PageAll = true })
	if out["count"] != 3 {
		t.Errorf("pageAll: %v", out["count"])
	}
}

func TestReadTableScopedAndMissing(t *testing.T) {
	if out, err := tableRun(t, tableNodes(), func(s *model.Step) { s.Selectors = []string{`role=table`} }); err != nil || out["count"] != 2 {
		t.Errorf("selectors: %v %v", out, err)
	}
	if _, err := tableRun(t, tableNodes()[:1], func(*model.Step) {}); err == nil || !strings.Contains(err.Error(), "không thấy bảng") {
		t.Errorf("không có bảng: %v", err)
	}
}

func TestReadTableSave(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"b.csv", "b.xlsx"} {
		p := filepath.Join(dir, name)
		out, err := tableRun(t, tableNodes(), func(s *model.Step) { s.FilePath = p })
		if err != nil || out["path"] != p {
			t.Fatalf("%s: %v %v", name, out, err)
		}
		var rows [][]string
		if strings.HasSuffix(name, ".csv") {
			rows, err = readCSV(p)
		} else {
			rows, err = readXLSX(p, "")
		}
		if err != nil || len(rows) != 3 || rows[1][0] != "Phở" {
			t.Errorf("%s đọc lại: %v %v", name, rows, err)
		}
	}
	if _, err := tableRun(t, tableNodes(), func(s *model.Step) { s.FilePath = filepath.Join(dir, "b.txt") }); err == nil {
		t.Error("đuôi .txt phải lỗi")
	}
}
