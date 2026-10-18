package web

import (
	"context"
	"fmt"
	"math"
	"time"

	"autoclicker/internal/model"
)

// Đổi toạ độ trang (CSS px, client) → màn hình (robotgo) bằng cách đo (D0: sai số 0–2 px): đưa chuột tới 2 điểm
// trong trang, đọc clientX/Y của pointermove, suy ra screen = off + client × scale mỗi trục; điểm thứ 3 để kiểm.
// Đo lại khi geom đổi (dời / đổi cỡ cửa sổ, zoom, DPI).

type fpoint struct{ X, Y float64 }

// geom – hình học cửa sổ trang báo; đổi thì số đo cũ không còn đúng.
type geom struct {
	ScreenX, ScreenY, OuterW, OuterH, InnerW, InnerH, DPR float64
}

type calibration struct {
	scale, off fpoint
	geom       geom
}

func (c calibration) toScreen(p fpoint) fpoint {
	return fpoint{c.off.X + p.X*c.scale.X, c.off.Y + p.Y*c.scale.Y}
}

// rectToScreen – khung client → khung màn hình (làm tròn về điểm nguyên robotgo).
func (c calibration) rectToScreen(x, y, w, h float64) model.Rect {
	a, b := c.toScreen(fpoint{x, y}), c.toScreen(fpoint{x + w, y + h})
	l, t := int(math.Round(a.X)), int(math.Round(a.Y))
	return model.Rect{X: l, Y: t, W: int(math.Round(b.X)) - l, H: int(math.Round(b.Y)) - t}
}

// maxResidual – điểm kiểm lệch quá số CSS px này thì coi như đo sai (trang chưa nhận chuột kịp, cửa sổ đang dời…).
const maxResidual = 3.0

// fitCalibration – từ 2 cặp (màn hình, client) suy ra công thức, kiểm bằng cặp thứ 3. s là điểm chuột đã làm tròn
// (robotgo chỉ nhận số nguyên – dùng điểm chưa làm tròn lệch 1 px ở mép, D0).
func fitCalibration(s1, c1, s2, c2, s3, c3 fpoint) (calibration, error) {
	if c2.X == c1.X || c2.Y == c1.Y {
		return calibration{}, fmt.Errorf("hai điểm đo trùng toạ độ trang (%.0f,%.0f)", c1.X, c1.Y)
	}
	var cal calibration
	cal.scale = fpoint{(s2.X - s1.X) / (c2.X - c1.X), (s2.Y - s1.Y) / (c2.Y - c1.Y)}
	cal.off = fpoint{s1.X - c1.X*cal.scale.X, s1.Y - c1.Y*cal.scale.Y}
	for _, k := range []float64{cal.scale.X, cal.scale.Y} {
		if k < 0.2 || k > 10 {
			return calibration{}, fmt.Errorf("tỷ lệ trang / màn hình đo được bất thường (%.3f)", k)
		}
	}
	p := cal.toScreen(c3)
	if d := math.Hypot(p.X-s3.X, p.Y-s3.Y) / cal.scale.X; d > maxResidual {
		return calibration{}, fmt.Errorf("điểm kiểm lệch %.1f px – đo lại", d)
	}
	return cal, nil
}

// probeFractions – 3 điểm đo theo tỷ lệ khung cửa sổ: x lệch trái (DevTools hay gắn phải), y tránh thanh công cụ trên
// và DevTools gắn dưới (D0).
var probeFractions = [3]fpoint{{0.15, 0.30}, {0.45, 0.55}, {0.30, 0.42}}

// mover – di chuột thật (device.Input.Move).
type mover func(x, y int) error

// pageProbe – đọc pointermove cuối trang nhận: số lần + clientX/Y.
type pageProbe func(ctx context.Context) (n int, at fpoint, err error)

// measure – đưa chuột tới 3 điểm trong khung cửa sổ win (toạ độ màn hình), ghép công thức.
func measure(ctx context.Context, win model.Rect, move mover, probe pageProbe) (calibration, error) {
	var s, c [3]fpoint
	for i, f := range probeFractions {
		x := int(math.Round(float64(win.X) + float64(win.W)*f.X))
		y := int(math.Round(float64(win.Y) + float64(win.H)*f.Y))
		// lệch nhẹ trước để chắc có sự kiện mới kể cả khi chuột đang đứng đúng chỗ; mốc đếm đọc sau bước lệch
		if err := move(x+4, y+4); err != nil {
			return calibration{}, err
		}
		n0, _, err := probe(ctx)
		if err != nil {
			return calibration{}, err
		}
		if err := move(x, y); err != nil {
			return calibration{}, err
		}
		n, at, err := waitProbe(ctx, probe, n0)
		if err != nil {
			return calibration{}, err
		}
		if n == n0 {
			return calibration{}, &BrowserError{fmt.Sprintf("trang không nhận chuột tại (%d, %d) – cửa sổ trình duyệt bị che?", x, y)}
		}
		s[i], c[i] = fpoint{float64(x), float64(y)}, at
	}
	return fitCalibration(s[0], c[0], s[1], c[1], s[2], c[2])
}

// waitProbe chờ tới khi trang báo pointermove mới (tối đa ~300 ms), rồi thêm một nhịp và đọc lại: sự kiện của bước
// lệch tới muộn thì vẫn lấy được sự kiện cuối (đúng chỗ chuột đứng).
func waitProbe(ctx context.Context, probe pageProbe, n0 int) (int, fpoint, error) {
	var at fpoint
	n := n0
	for i := 0; i < 10 && n == n0; i++ {
		if err := sleepCtx(ctx, 30*time.Millisecond); err != nil {
			return n, at, err
		}
		var err error
		if n, at, err = probe(ctx); err != nil {
			return n, at, err
		}
	}
	if n == n0 {
		return n, at, nil
	}
	if err := sleepCtx(ctx, 40*time.Millisecond); err != nil {
		return n, at, err
	}
	return probe(ctx)
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}
