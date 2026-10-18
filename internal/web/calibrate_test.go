package web

import (
	"context"
	"testing"

	"autoclicker/internal/model"
)

// fakePage – trang giả: client = (screen - off) / scale, nhận chuột khi nằm trong vùng trang.
type fakePage struct {
	off, scale fpoint
	page       model.Rect // vùng trang trên màn hình
	n          int
	at         fpoint
}

func (p *fakePage) move(x, y int) error {
	if x >= p.page.X && x < p.page.X+p.page.W && y >= p.page.Y && y < p.page.Y+p.page.H {
		p.n++
		p.at = fpoint{(float64(x) - p.off.X) / p.scale.X, (float64(y) - p.off.Y) / p.scale.Y}
	}
	return nil
}

func (p *fakePage) probe(context.Context) (probeRead, error) { return probeRead{n: p.n, at: p.at}, nil }

func TestMeasure(t *testing.T) {
	win := model.Rect{X: 100, Y: 80, W: 1200, H: 850}
	p := &fakePage{off: fpoint{100, 166}, scale: fpoint{1.25, 1.25}, page: model.Rect{X: 100, Y: 166, W: 1200, H: 764}}
	cal, err := measure(context.Background(), win, p.move, p.probe)
	if err != nil {
		t.Fatal(err)
	}
	// nút client (40, 40, 20×20) → màn hình (150, 216, 25×25)
	if got := cal.rectToScreen(40, 40, 20, 20); got != (model.Rect{X: 150, Y: 216, W: 25, H: 25}) {
		t.Errorf("rectToScreen = %+v", got)
	}
}

func TestMeasureCovered(t *testing.T) {
	win := model.Rect{X: 100, Y: 80, W: 1200, H: 850}
	p := &fakePage{off: fpoint{100, 166}, scale: fpoint{1, 1}, page: model.Rect{}} // không điểm nào tới trang
	if _, err := measure(context.Background(), win, p.move, p.probe); err == nil {
		t.Fatal("trang không nhận chuột phải báo lỗi")
	}
}

func TestFitCalibrationRejects(t *testing.T) {
	s1, c1 := fpoint{100, 100}, fpoint{0, 0}
	s2, c2 := fpoint{300, 300}, fpoint{200, 200}
	if _, err := fitCalibration(s1, c1, s2, c2, fpoint{200, 200}, fpoint{100, 100}); err != nil {
		t.Errorf("đo đúng mà lỗi: %v", err)
	}
	if _, err := fitCalibration(s1, c1, s2, c2, fpoint{200, 200}, fpoint{120, 100}); err == nil {
		t.Error("điểm kiểm lệch 20 px phải báo lỗi")
	}
	if _, err := fitCalibration(s1, c1, s2, c1, s1, c1); err == nil {
		t.Error("hai điểm trùng phải báo lỗi")
	}
}

func TestToClientInverse(t *testing.T) {
	cal := calibration{scale: fpoint{1.25, 2}, off: fpoint{100, 166}}
	if got := cal.toClient(cal.toScreen(fpoint{40, 30})); got != (fpoint{40, 30}) {
		t.Errorf("toClient(toScreen) = %+v", got)
	}
}
