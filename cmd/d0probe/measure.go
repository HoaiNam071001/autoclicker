//go:build probe

package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-vgo/robotgo"
)

type point struct{ X, Y float64 }

type geom struct {
	ScreenX, ScreenY, OuterW, OuterH, InnerW, InnerH, DPR float64
}

// calibration: screen = off + client × scale (mỗi trục).
type calibration struct {
	scale, off point
	resid      float64 // sai lệch ở điểm kiểm (px client)
	geom       geom
	took       time.Duration
}

func (c calibration) String() string {
	return fmt.Sprintf("đo: scale=(%.4f, %.4f) off=(%.2f, %.2f) dư điểm kiểm=%.2f px, DPR=%.2f, inner=%vx%v, mất %v",
		c.scale.X, c.scale.Y, c.off.X, c.off.Y, c.resid, c.geom.DPR, c.geom.InnerW, c.geom.InnerH, c.took.Round(time.Millisecond))
}

func (c calibration) toScreen(p point) point {
	return point{c.off.X + p.X*c.scale.X, c.off.Y + p.Y*c.scale.Y}
}

// calibrate – đưa chuột tới 2 điểm trong vùng trang (đoán từ khung cửa sổ), đọc clientX/Y của pointermove.
// Điểm thứ 3 để kiểm tra công thức tuyến tính.
func calibrate(ctx context.Context) (calibration, error) {
	start := time.Now()
	w, err := windowOf(ctx)
	if err != nil {
		return calibration{}, err
	}
	b := w.Bounds
	at := func(fx, fy float64) point {
		return point{float64(*b.Left) + float64(b.Width)*fx, float64(*b.Top) + float64(b.Height)*fy}
	}
	// x lệch trái (DevTools hay gắn bên phải), y tránh thanh công cụ trên + DevTools gắn dưới.
	s1, s2, s3 := at(0.15, 0.30), at(0.45, 0.55), at(0.30, 0.42)
	c1, err := probeAt(ctx, &s1)
	if err != nil {
		return calibration{}, err
	}
	c2, err := probeAt(ctx, &s2)
	if err != nil {
		return calibration{}, err
	}
	c3, err := probeAt(ctx, &s3)
	if err != nil {
		return calibration{}, err
	}
	var cal calibration
	cal.scale = point{(s2.X - s1.X) / (c2.X - c1.X), (s2.Y - s1.Y) / (c2.Y - c1.Y)}
	cal.off = point{s1.X - c1.X*cal.scale.X, s1.Y - c1.Y*cal.scale.Y}
	pred := cal.toScreen(c3)
	cal.resid = math.Hypot(pred.X-s3.X, pred.Y-s3.Y) / cal.scale.X
	cal.geom, err = readGeom(ctx)
	cal.took = time.Since(start)
	return cal, err
}

// probeAt – robotgo tới s (lệch nhẹ trước để chắc có sự kiện mới), trả clientX/Y trang nhận được.
// s được làm tròn về đúng điểm chuột đã tới (robotgo chỉ nhận số nguyên) – công thức phải dùng điểm này.
func probeAt(ctx context.Context, s *point) (point, error) {
	before, err := evalJS[int](ctx, `__mv.n`)
	if err != nil {
		return point{}, err
	}
	x, y := int(math.Round(s.X)), int(math.Round(s.Y))
	*s = point{float64(x), float64(y)}
	robotgo.Move(x+4, y+4)
	time.Sleep(40 * time.Millisecond)
	robotgo.Move(x, y)
	time.Sleep(120 * time.Millisecond)
	mv, err := evalJS[struct {
		X, Y float64
		N    int
	}](ctx, `({X: __mv.x, Y: __mv.y, N: __mv.n})`)
	if err != nil {
		return point{}, err
	}
	if mv.N == before {
		return point{}, fmt.Errorf("trang không nhận pointermove tại (%d,%d) – cửa sổ bị che / không ở trên cùng?", x, y)
	}
	return point{mv.X, mv.Y}, nil
}

// clickTargets – bấm từng nút theo số đo, chấm sai số bằng toạ độ pointerdown trang nhận được.
func clickTargets(ctx context.Context, cal calibration, label string) error {
	if label != "" {
		fmt.Printf("-- %s\n", label)
	}
	if _, err := evalJS[bool](ctx, `__downs = []; window.scrollTo(0, 0); true`); err != nil {
		return err
	}
	worst := 0.0
	for _, id := range targetIDs {
		js := fmt.Sprintf(`(async () => { const el = document.getElementById(%q);
			if (%q === "far") { el.scrollIntoView({block: "center"}); await new Promise(r => setTimeout(r, 300)); }
			const r = el.getBoundingClientRect(); return {X: r.left + r.width/2, Y: r.top + r.height/2}; })()`, id, id)
		c, err := evalJS[point](ctx, js)
		if err != nil {
			return err
		}
		s := cal.toScreen(c)
		robotgo.Move(int(math.Round(s.X)), int(math.Round(s.Y)))
		time.Sleep(time.Duration(*settleMs) * time.Millisecond)
		robotgo.Click("left", false)
		time.Sleep(150 * time.Millisecond)
		downs, err := evalJS[[]struct {
			ID      string
			X, Y    float64
			Trusted bool
		}](ctx, `__downs.splice(0).map(d => ({ID: d.id, X: d.x, Y: d.y, Trusted: d.trusted}))`)
		if err != nil {
			return err
		}
		if len(downs) == 0 {
			diag, _ := evalJS[string](ctx, fmt.Sprintf(`(document.elementFromPoint(%f, %f)?.id ?? "-") + " scrollY=" + scrollY + " mv=" + __mv.x + "," + __mv.y`, c.X, c.Y))
			fmt.Printf("  %-4s đích client=(%.1f,%.1f) màn=(%.0f,%.0f) → trang không nhận click (%s)\n", id, c.X, c.Y, s.X, s.Y, diag)
			worst = math.Inf(1)
			continue
		}
		d := downs[len(downs)-1]
		errCSS := math.Hypot(d.X-c.X, d.Y-c.Y)
		errPt := errCSS * cal.scale.X // điểm màn hình = CSS × scale
		worst = math.Max(worst, errPt)
		fmt.Printf("  %-4s đích=(%.1f,%.1f) nhận=(%.1f,%.1f) trúng=%v trusted=%v lệch=%.2f CSS px ≈ %.2f điểm màn hình\n",
			id, c.X, c.Y, d.X, d.Y, d.ID == id, d.Trusted, errCSS, errPt)
	}
	fmt.Printf("  lệch lớn nhất: %.2f điểm màn hình ≈ %.2f px vật lý\n", worst, worst*screenRatio(cal))
	return nil
}

// screenRatio – px vật lý / điểm màn hình: DPR = zoom × tỷ lệ màn, scale (điểm / CSS px) = zoom ⇒ tỷ lệ màn = DPR / scale.
func screenRatio(c calibration) float64 {
	if c.scale.X == 0 || c.geom.DPR == 0 {
		return 1
	}
	return c.geom.DPR / c.scale.X
}

func readGeom(ctx context.Context) (geom, error) {
	return evalJS[geom](ctx, `({ScreenX: screenX, ScreenY: screenY, OuterW: outerWidth, OuterH: outerHeight,
		InnerW: innerWidth, InnerH: innerHeight, DPR: devicePixelRatio})`)
}

func moveWindow(ctx context.Context, left, top int) error {
	err := setWindowPos(ctx, int64(left), int64(top))
	time.Sleep(600 * time.Millisecond)
	return err
}

// bringToFront – thử Page.bringToFront trước; nếu trang chưa có focus thì kích hoạt tiến trình (ghi lại cần cách nào).
func bringToFront(ctx context.Context, pid int) error {
	if err := bringTab(ctx); err != nil {
		return err
	}
	time.Sleep(300 * time.Millisecond)
	focused, err := evalJS[bool](ctx, `document.hasFocus()`)
	if err != nil {
		return err
	}
	if focused {
		fmt.Println("đưa lên trước: Page.bringToFront đủ")
		return nil
	}
	if err := robotgo.ActivePid(pid); err != nil {
		return err
	}
	time.Sleep(400 * time.Millisecond)
	focused, _ = evalJS[bool](ctx, `document.hasFocus()`)
	fmt.Printf("đưa lên trước: Page.bringToFront chưa đủ, thêm ActivePid → hasFocus=%v\n", focused)
	return nil
}
