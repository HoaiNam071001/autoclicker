package script

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Phần phụ của bộ tính node Tính toán (calc.go): đóng nháy giá trị biến, tách token, kiểm kiểu giá trị.

var calcNumRe = regexp.MustCompile(`^-?\d+(\.\d+)?$`)

// CalcQuote – giá trị biến chèn vào biểu thức: số / true / false giữ nguyên, còn lại thành chuỗi "…" (chữ OCR như
// "1 + 1" hay "a\"b" không thành phép tính).
func CalcQuote(v string) string {
	t := strings.TrimSpace(v)
	if calcNumRe.MatchString(t) || t == "true" || t == "false" {
		return t
	}
	return `"` + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), `"`, `\"`) + `"`
}

type calcTok struct {
	kind byte // n số | s chuỗi | w từ (true, and, round…) | o dấu
	s    string
	n    float64
}

var calcOps = []string{"==", "!=", "<=", ">=", "&&", "||", "+", "-", "*", "/", "%", "(", ")", ",", "<", ">", "!"}

func calcLex(src string) ([]calcTok, error) {
	var out []calcTok
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c >= '0' && c <= '9' || c == '.':
			j := i
			for j < len(src) && (src[j] >= '0' && src[j] <= '9' || src[j] == '.') {
				j++
			}
			n, err := strconv.ParseFloat(src[i:j], 64)
			if err != nil || src[i] == '.' || src[j-1] == '.' {
				return nil, synErr("“%s” không phải số", src[i:j])
			}
			out, i = append(out, calcTok{kind: 'n', s: src[i:j], n: n}), j
		case c == '"':
			var b strings.Builder
			j := i + 1
			for ; j < len(src) && src[j] != '"'; j++ {
				if src[j] == '\\' && j+1 < len(src) {
					j++
				}
				b.WriteByte(src[j])
			}
			if j >= len(src) {
				return nil, synErr("thiếu dấu \" đóng chuỗi")
			}
			out, i = append(out, calcTok{kind: 's', s: b.String()}), j+1
		case c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_':
			j := i
			for j < len(src) && (src[j] >= 'a' && src[j] <= 'z' || src[j] >= 'A' && src[j] <= 'Z' || src[j] == '_' || src[j] >= '0' && src[j] <= '9') {
				j++
			}
			out, i = append(out, calcTok{kind: 'w', s: strings.ToLower(src[i:j])}), j
		default:
			op := ""
			for _, o := range calcOps {
				if strings.HasPrefix(src[i:], o) {
					op = o
					break
				}
			}
			if op == "" {
				r := []rune(src[i:])[0]
				return nil, synErr("không hiểu ký tự “%c”", r)
			}
			out, i = append(out, calcTok{kind: 'o', s: op}), i+len(op)
		}
	}
	return out, nil
}

// calcRound – làm tròn nửa lên như Math.round của JS (Go math.Round làm tròn xa 0: -2.5 → -3) để hai bên cùng kết quả.
func calcRound(x float64) float64 {
	return math.Floor(x + 0.5)
}

func calcNums(a, b any, op string) (float64, float64, error) {
	x, ok1 := a.(float64)
	y, ok2 := b.(float64)
	if !ok1 {
		return 0, 0, fmt.Errorf("“%s” cần số, có %s", op, calcShow(a))
	}
	if !ok2 {
		return 0, 0, fmt.Errorf("“%s” cần số, có %s", op, calcShow(b))
	}
	return x, y, nil
}

func calcBool(v any, op string) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("“%s” cần true / false, có %s", op, calcShow(v))
	}
	return b, nil
}

func calcShow(v any) string {
	switch x := v.(type) {
	case string:
		return fmt.Sprintf("chữ “%s” (số kiểu 1.234 đ: dùng node Xử lý chữ lấy số trước)", x)
	case bool:
		return strconv.FormatBool(x)
	}
	return fmt.Sprint(v)
}

// calcSyntaxErr – lỗi cú pháp: sai ngay cả khi chưa có giá trị biến (kiểm tra lúc soạn báo được). Lỗi khác (chia 0,
// sai kiểu) chỉ biết khi chạy.
type calcSyntaxErr struct{ msg string }

func (e *calcSyntaxErr) Error() string { return e.msg }

func synErr(format string, a ...any) error { return &calcSyntaxErr{fmt.Sprintf(format, a...)} }

// CheckCalc – lỗi cú pháp của biểu thức còn {{biến}} (biến tạm coi là số 1); nil = chạy được nếu biến đúng kiểu.
func CheckCalc(expr string) error {
	_, err := EvalCalc(NodeVarRe.ReplaceAllString(varRe.ReplaceAllString(expr, "1"), "1"))
	var syn *calcSyntaxErr
	if errors.As(err, &syn) {
		return err
	}
	return nil
}
