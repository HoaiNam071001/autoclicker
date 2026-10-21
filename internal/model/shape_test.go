package model

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

// Cùng dữ liệu với frontend/src/flow/shape.test.ts (testdata/shapes.json).
func TestShapePoints(t *testing.T) {
	b, err := os.ReadFile("../../testdata/shapes.json")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Cases []struct {
			Shape        string
			CX, CY, W, H float64
			Points       []FPt
		}
	}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}
	for _, c := range data.Cases {
		got := ShapePoints(c.Shape, c.CX, c.CY, c.W, c.H)
		if len(got) != len(c.Points) {
			t.Errorf("%s: %d đỉnh, muốn %d", c.Shape, len(got), len(c.Points))
			continue
		}
		for i, p := range got {
			if math.Abs(p.X-c.Points[i].X) > 1e-3 || math.Abs(p.Y-c.Points[i].Y) > 1e-3 {
				t.Errorf("%s đỉnh %d: %v, muốn %v", c.Shape, i, p, c.Points[i])
			}
		}
		if n := len(got); n > 0 && (math.Abs(got[0].X-got[n-1].X) > 1e-9 || math.Abs(got[0].Y-got[n-1].Y) > 1e-9) {
			t.Errorf("%s: hình không khép kín", c.Shape)
		}
	}
}

// Mỗi loại node phải có nhóm + gợi ý chọn (AI đọc qua get_step_types), và không còn dòng thừa.
func TestEveryTypeHasUse(t *testing.T) {
	seen := map[string]bool{}
	for _, d := range StepTypeDocs() {
		seen[d.Type] = true
		if u := UseOf(d.Type); u.Group == "" || u.Use == "" {
			t.Errorf("loại %s thiếu gợi ý chọn trong steptypes_use.go", d.Type)
		}
	}
	for k := range stepTypeUses {
		if !seen[k] {
			t.Errorf("steptypes_use.go có loại không tồn tại: %s", k)
		}
	}
}
