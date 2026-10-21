package humanize

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"ghosthand/internal/model"
)

// Cùng dữ liệu với frontend/src/overlay/humanPath.test.ts (testdata/human_path.json) – xem trước vẽ đúng đường chuột thật.
func TestSharedHumanPath(t *testing.T) {
	b, err := os.ReadFile("../../testdata/human_path.json")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Profiles map[string]struct{ Speed, Curve, Overshoot float64 }
		Fitts    []struct {
			D    float64
			Size int
			Ms   float64
		}
		Curves []struct {
			A, B   model.Point
			O1, O2 float64
			C1, C2 FPoint
			At     []struct {
				T float64
				P FPoint
			}
		}
	}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-3 }
	nearP := func(a, b FPoint) bool { return near(a.X, b.X) && near(a.Y, b.Y) }
	if len(data.Profiles) != len(profiles) {
		t.Errorf("testdata có %d mức, profiles có %d", len(data.Profiles), len(profiles))
	}
	for lv, want := range data.Profiles {
		p := profiles[lv]
		if p.Speed != want.Speed || p.Curve != want.Curve || p.Overshoot != want.Overshoot {
			t.Errorf("mức %s: %+v, testdata %+v", lv, p, want)
		}
	}
	for _, f := range data.Fitts {
		if got := Fitts(f.D, f.Size); !near(got, f.Ms) {
			t.Errorf("Fitts(%v, %d) = %v, muốn %v", f.D, f.Size, got, f.Ms)
		}
	}
	for _, c := range data.Curves {
		c1, c2 := Controls(c.A, c.B, c.O1, c.O2)
		if !nearP(c1, c.C1) || !nearP(c2, c.C2) {
			t.Errorf("Controls(%v, %v) = %v %v, muốn %v %v", c.A, c.B, c1, c2, c.C1, c.C2)
		}
		fa, fb := FPoint{float64(c.A.X), float64(c.A.Y)}, FPoint{float64(c.B.X), float64(c.B.Y)}
		for _, a := range c.At {
			if got := CurveAt(fa, c1, c2, fb, a.T); !nearP(got, a.P) {
				t.Errorf("CurveAt(t=%v) = %v, muốn %v", a.T, got, a.P)
			}
		}
	}
}
