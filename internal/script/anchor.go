package script

import (
	"math"

	"ghosthand/internal/model"
)

// Neo của một trục toạ độ (Point.Ax / Point.Ay) – mỗi trục neo riêng vào cạnh của vùng.
const (
	AnchorStart  = "start"  // cách cạnh trái / trên (mặc định – để trống cũng là start, kịch bản cũ chạy như cũ)
	AnchorEnd    = "end"    // cách cạnh phải / dưới (giá trị dương, UI hiện dấu −)
	AnchorCenter = "center" // lệch so với giữa cạnh (có dấu)
	AnchorRatio  = "ratio"  // theo tỷ lệ cạnh vùng, đơn vị 1/10000 (5000 = 50%)
)

// RatioScale – giá trị neo ratio ứng với 100%.
const RatioScale = 10000

// resolveAxis trả toạ độ màn hình của giá trị v neo theo a, trên đoạn [origin, origin+size].
func resolveAxis(v int, a string, origin, size int) int {
	switch a {
	case AnchorEnd:
		return origin + size - v
	case AnchorCenter:
		return origin + size/2 + v
	case AnchorRatio:
		return origin + int(math.Round(float64(size)*float64(v)/RatioScale))
	default:
		return origin + v
	}
}

// ResolvePoint trả toạ độ màn hình (robotgo) của điểm p neo theo khung b (vùng hoặc màn hình).
// Phải khớp resolvePoint trong frontend/src/anchor.ts.
func ResolvePoint(p model.Point, b model.Rect) model.Point {
	return model.Point{X: resolveAxis(p.X, p.Ax, b.X, b.W), Y: resolveAxis(p.Y, p.Ay, b.Y, b.H)}
}

// CoordAbsolute – Step.CoordSpace: điểm là toạ độ màn hình thật, không cộng gốc vùng (khớp ảnh screenshot).
const CoordAbsolute = "absolute"

// StepBounds – khung tính toạ độ điểm của node s: absolute thì gốc (0,0) của màn hình, còn lại như Bounds.
// Phải khớp stepBounds trong frontend/src/anchor.ts.
func StepBounds(sc model.Script, screen model.Rect, s model.Step) model.Rect {
	if s.CoordSpace == CoordAbsolute {
		return model.Rect{W: screen.X + screen.W, H: screen.Y + screen.H}
	}
	return Bounds(sc, screen)
}

// Bounds trả khung tính toạ độ điểm: vùng của kịch bản, không có thì màn hình.
func Bounds(sc model.Script, screen model.Rect) model.Rect {
	if sc.Region != nil {
		return *sc.Region
	}
	return screen
}
