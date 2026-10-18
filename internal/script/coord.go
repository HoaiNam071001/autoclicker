package script

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"autoclicker/internal/model"
)

// Toạ độ từ biến (bước 20 F5b): ô xVar / yVar sau khi thay biến là số cộng trừ px, vd "512 + 20 - 3" (số thực được
// làm tròn – output w / 2 sau này). Không có phép nhân / chia: chỉ cần lệch px so với chỗ thấy.

// EvalCoord tính biểu thức toạ độ đã thay biến.
func EvalCoord(expr string) (int, error) {
	s := strings.TrimSpace(expr)
	if s == "" {
		return 0, errors.New("chưa nhập toạ độ")
	}
	if strings.Contains(s, "{{") {
		return 0, fmt.Errorf("“%s” còn biến chưa thay", s)
	}
	sum, sign, want := 0.0, 1.0, true // want = đang chờ một số
	for s != "" {
		s = strings.TrimSpace(s)
		switch {
		case !want && (s[0] == '+' || s[0] == '-'):
			sign, want = 1, true
			if s[0] == '-' {
				sign = -1
			}
			s = s[1:]
		case want:
			n := numPrefix(s)
			if n == 0 {
				return 0, fmt.Errorf("“%s” không phải số (cần số hoặc số + / − px)", strings.TrimSpace(expr))
			}
			v, err := strconv.ParseFloat(s[:n], 64)
			if err != nil {
				return 0, fmt.Errorf("“%s” không phải số", s[:n])
			}
			sum, want, s = sum+sign*v, false, s[n:]
		default:
			return 0, fmt.Errorf("“%s” không phải số (cần số hoặc số + / − px)", strings.TrimSpace(expr))
		}
	}
	if want {
		return 0, fmt.Errorf("“%s” thiếu số sau dấu", strings.TrimSpace(expr))
	}
	return int(math.Round(sum)), nil
}

// numPrefix – độ dài số (có thể có dấu − đầu, phần thập phân) ở đầu s.
func numPrefix(s string) int {
	i := 0
	if i < len(s) && s[i] == '-' {
		i++
	}
	digits := 0
	for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.') {
		i++
		digits++
	}
	if digits == 0 {
		return 0
	}
	return i
}

// PointOnScreen – toạ độ màn hình của điểm: dùng biến thì tính xVar / yVar (đã thay biến), không thì theo neo trên b.
func PointOnScreen(p model.Point, b model.Rect) (model.Point, error) {
	if !p.UsesVar() {
		return ResolvePoint(p, b), nil
	}
	x, err := EvalCoord(p.XVar)
	if err != nil {
		return model.Point{}, fmt.Errorf("toạ độ X: %w", err)
	}
	y, err := EvalCoord(p.YVar)
	if err != nil {
		return model.Point{}, fmt.Errorf("toạ độ Y: %w", err)
	}
	return model.Point{X: x, Y: y}, nil
}
