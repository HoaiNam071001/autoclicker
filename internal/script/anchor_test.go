package script

import (
	"testing"

	"autoclicker/internal/model"
)

func TestResolvePointAnchors(t *testing.T) {
	r := model.Rect{X: 100, Y: 50, W: 1000, H: 700}
	cases := []struct {
		p    model.Point
		want model.Point
	}{
		{model.Point{X: 300, Y: 300}, model.Point{X: 400, Y: 350}},                                     // trống = start
		{model.Point{X: 300, Y: 300, Ax: AnchorEnd}, model.Point{X: 800, Y: 350}},                      // phải · trên
		{model.Point{X: 300, Y: 300, Ay: AnchorEnd}, model.Point{X: 400, Y: 450}},                      // trái · dưới
		{model.Point{X: 300, Y: 300, Ax: AnchorEnd, Ay: AnchorEnd}, model.Point{X: 800, Y: 450}},       // phải · dưới
		{model.Point{X: -20, Y: 10, Ax: AnchorCenter, Ay: AnchorCenter}, model.Point{X: 580, Y: 410}},  // lệch so với tâm
		{model.Point{X: 2500, Y: 5000, Ax: AnchorRatio, Ay: AnchorRatio}, model.Point{X: 350, Y: 400}}, // 25% · 50%
	}
	for _, c := range cases {
		if got := ResolvePoint(c.p, r); got != c.want {
			t.Errorf("resolvePoint(%+v) = %+v, want %+v", c.p, got, c.want)
		}
	}
}
