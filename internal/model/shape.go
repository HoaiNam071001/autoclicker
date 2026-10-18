package model

import "math"

// FPt – điểm toạ độ thực (đỉnh hình vẽ trước khi làm tròn).
type FPt struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// circleSegments – số đoạn của hình tròn (đủ mượt, ít lệnh chuột).
const circleSegments = 36

// ShapePoints – các đỉnh của hình (node Vẽ đường, bước 35 N7) quanh tâm (cx, cy), cỡ w × h (h ≤ 0 = w), điểm cuối trùng
// điểm đầu (khép kín). Tròn / sao bắt đầu ở đỉnh trên, đi chiều kim đồng hồ; chữ nhật từ góc trên trái. Shape lạ / line → nil.
// Cùng công thức shapePoints (frontend/src/flow/shape.ts), test chung testdata/shapes.json.
func ShapePoints(shape string, cx, cy, w, h float64) []FPt {
	if h <= 0 {
		h = w
	}
	rx, ry := w/2, h/2
	ring := func(n int, r func(i int) float64) []FPt {
		out := make([]FPt, 0, n+1)
		for i := 0; i <= n; i++ {
			a := -math.Pi/2 + 2*math.Pi*float64(i%n)/float64(n)
			k := r(i % n)
			out = append(out, FPt{cx + rx*k*math.Cos(a), cy + ry*k*math.Sin(a)})
		}
		return out
	}
	switch shape {
	case ShapeCircle:
		return ring(circleSegments, func(int) float64 { return 1 })
	case ShapeStar:
		return ring(10, func(i int) float64 { return []float64{1, 0.4}[i%2] })
	case ShapeRect:
		return []FPt{{cx - rx, cy - ry}, {cx + rx, cy - ry}, {cx + rx, cy + ry}, {cx - rx, cy + ry}, {cx - rx, cy - ry}}
	}
	return nil
}
