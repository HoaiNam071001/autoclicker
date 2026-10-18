package humanize

import (
	"math"

	"autoclicker/internal/model"
)

// Đường di chuột: Bezier bậc 3 với 2 điểm điều khiển lệch sang ngang (cong như cổ tay xoay), tốc độ theo đường
// minimum-jerk (tăng rồi giảm), thời gian theo định luật Fitts (xa / đích nhỏ = lâu hơn), thỉnh thoảng đi lố rồi chỉnh.

// stepMs – khoảng giữa hai lần đặt chuột trên đường (~80 lần / giây như chuột thật).
const stepMs = 12

// edgeMargin – điểm giữa đường giữ cách mép màn ≥ px này: không chạm góc dừng khẩn cấp (failsafe cornerSize 3).
const edgeMargin = 6

// Move – chờ WaitMs rồi đưa chuột tới (X, Y). Phần tử cuối luôn là đích.
type Move struct {
	X, Y   int
	WaitMs int
}

// Target – cỡ đích + thời gian cố định.
type Target struct {
	Size  int // cạnh nhỏ của đích (px); 0 = không biết (coi như nút 24 px)
	DurMs int // > 0: đi đúng chừng này (kéo thả theo DragMs), không đi lố
}

// Path – các bước đưa chuột từ from tới to. bounds: khung màn (W = 0 = không giới hạn) – điểm giữa đường nằm trong.
func Path(r R, p Profile, from, to model.Point, t Target, bounds model.Rect) []Move {
	d := dist(from, to)
	if d < 2 {
		return []Move{{X: to.X, Y: to.Y}}
	}
	dur := t.DurMs
	if dur <= 0 {
		dur = int(Fitts(d, t.Size) * p.Speed * (0.85 + 0.3*r.Float64()))
	}
	if t.DurMs > 0 || d < 80 || r.Float64() >= p.Overshoot {
		return curve(r, p, from, to, dur, bounds)
	}
	// đi lố 3–8 % quãng (6–40 px) theo hướng đi, lệch ngang chút, rồi chỉnh về đích chậm hơn
	over := math.Min(math.Max(d*(0.03+0.05*r.Float64()), 6), 40)
	ux, uy := float64(to.X-from.X)/d, float64(to.Y-from.Y)/d
	side := r.NormFloat64() * over * 0.3
	o := clamp(model.Point{X: to.X + round(ux*over-uy*side), Y: to.Y + round(uy*over+ux*side)}, bounds)
	first := curve(r, p, from, o, dur*85/100, bounds)
	fix := curve(r, Profile{Speed: p.Speed, Curve: p.Curve / 2, Tremor: p.Tremor}, o, to, int((120+100*r.Float64())*p.Speed), bounds)
	fix[0].WaitMs += 40 + r.IntN(80) // khựng lại trước khi chỉnh
	return append(first, fix...)
}

// Fitts – thời gian (ms) di chuột quãng d tới đích cỡ size: a + b·log2(d / size + 1).
// Cùng công thức fitts (frontend/src/overlay/humanPath.ts), test chung testdata/human_path.json.
func Fitts(d float64, size int) float64 {
	w := float64(size)
	if w <= 0 {
		w = 24
	}
	return 80 + 110*math.Log2(d/math.Max(w, 8)+1)
}

func curve(r R, p Profile, a, b model.Point, dur int, bounds model.Rect) []Move {
	d := dist(a, b)
	n := max(dur/stepMs, 2)
	bend := d * p.Curve * (r.Float64()*2 - 1)
	c1, c2 := Controls(a, b, bend*(0.6+0.8*r.Float64()), bend*(0.6+0.8*r.Float64()))
	fa, fb := FPoint{float64(a.X), float64(a.Y)}, FPoint{float64(b.X), float64(b.Y)}

	out := make([]Move, 0, n)
	wait := 0
	for i := 1; i <= n; i++ {
		wait += dur / n
		if i == n {
			out = append(out, Move{X: b.X, Y: b.Y, WaitMs: wait})
			break
		}
		c := CurveAt(fa, c1, c2, fb, float64(i)/float64(n))
		q := clamp(model.Point{X: round(c.X + r.NormFloat64()*p.Tremor), Y: round(c.Y + r.NormFloat64()*p.Tremor)}, bounds)
		if len(out) > 0 && out[len(out)-1].X == q.X && out[len(out)-1].Y == q.Y {
			continue // đứng yên: dồn thời gian chờ sang bước sau
		}
		out = append(out, Move{X: q.X, Y: q.Y, WaitMs: wait})
		wait = 0
	}
	return out
}

// FPoint – điểm toạ độ thực (điểm điều khiển, điểm trên đường trước khi làm tròn).
type FPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Controls – hai điểm điều khiển Bezier của đường a → b: ở 30 % / 70 % quãng, lệch theo pháp tuyến o1 / o2 px
// (cùng dấu = cong một chiều). Cùng công thức controls ở humanPath.ts.
func Controls(a, b model.Point, o1, o2 float64) (FPoint, FPoint) {
	ax, ay, bx, by := float64(a.X), float64(a.Y), float64(b.X), float64(b.Y)
	d := math.Max(math.Hypot(bx-ax, by-ay), 1e-9)
	nx, ny := -(by-ay)/d, (bx-ax)/d
	return FPoint{ax + (bx-ax)*0.3 + nx*o1, ay + (by-ay)*0.3 + ny*o1}, FPoint{ax + (bx-ax)*0.7 + nx*o2, ay + (by-ay)*0.7 + ny*o2}
}

// CurveAt – điểm trên đường ở thời điểm t ∈ [0, 1] (tiến độ theo minJerk). Cùng công thức curveAt ở humanPath.ts.
func CurveAt(a, c1, c2, b FPoint, t float64) FPoint {
	s := minJerk(t)
	return FPoint{bezier(a.X, c1.X, c2.X, b.X, s), bezier(a.Y, c1.Y, c2.Y, b.Y, s)}
}

// minJerk – tiến độ theo thời gian t ∈ [0, 1] của chuyển động tay: chậm – nhanh – chậm.
func minJerk(t float64) float64 { return t * t * t * (10 - 15*t + 6*t*t) }

func bezier(p0, p1, p2, p3, t float64) float64 {
	u := 1 - t
	return u*u*u*p0 + 3*u*u*t*p1 + 3*u*t*t*p2 + t*t*t*p3
}

func clamp(q model.Point, b model.Rect) model.Point {
	if b.W <= 0 || b.H <= 0 {
		return q
	}
	q.X = min(max(q.X, b.X+edgeMargin), b.X+b.W-1-edgeMargin)
	q.Y = min(max(q.Y, b.Y+edgeMargin), b.Y+b.H-1-edgeMargin)
	return q
}

func dist(a, b model.Point) float64 { return math.Hypot(float64(b.X-a.X), float64(b.Y-a.Y)) }

func round(v float64) int { return int(math.Round(v)) }

// Aim – điểm bấm quanh tâm (cx, cy) của đích cỡ w × h: phân bố chuẩn (σ = cạnh / 8), giữ trong 1/3 giữa đích.
// Đích không rõ cỡ (w hoặc h ≤ 0) → đúng tâm.
func Aim(r R, cx, cy, w, h int) (int, int) {
	if w <= 0 || h <= 0 {
		return cx, cy
	}
	off := func(side int) int {
		lim := float64(side) / 6
		return round(math.Max(-lim, math.Min(lim, r.NormFloat64()*float64(side)/8)))
	}
	return cx + off(w), cy + off(h)
}
