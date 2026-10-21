package script

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Bộ tính của node Tính toán (bước 42 G4a) – tự viết, không eval. Cùng cú pháp + kết quả với frontend/src/flow/calc.ts
// (test chung testdata/calc.json):
//   số 12 / 1.5, chuỗi "…" (\" \\), true / false; + - * / % (chỉ số); == != (mọi kiểu) < <= > >= (số);
//   and / && , or / || , not / ! (true-false); round(x[, n]) floor ceil abs min max; ngoặc ().
// Kết quả: float64 (làm tròn 10 chữ số thập phân – bỏ nhiễu 0.1 + 0.2), string hoặc bool.

// EvalCalc tính biểu thức đã thay biến.
func EvalCalc(expr string) (any, error) {
	if strings.TrimSpace(expr) == "" {
		return nil, synErr("chưa nhập biểu thức")
	}
	if strings.Contains(expr, "{{") {
		return nil, fmt.Errorf("“%s” còn biến chưa thay", expr)
	}
	toks, err := calcLex(expr)
	if err != nil {
		return nil, err
	}
	p := &calcParser{toks: toks}
	v, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.i < len(p.toks) {
		return nil, synErr("thừa “%s”", p.toks[p.i].s)
	}
	if f, ok := v.(float64); ok {
		if math.IsInf(f, 0) || math.IsNaN(f) {
			return nil, errors.New("kết quả không phải số hữu hạn")
		}
		if math.Abs(f) < 1e15 {
			v = calcRound(f*1e10) / 1e10
		}
	}
	return v, nil
}

type calcParser struct {
	toks []calcTok
	i    int
}

// eat – token kế là một trong ss (dấu hoặc từ khoá) thì lấy ra.
func (p *calcParser) eat(ss ...string) (string, bool) {
	if p.i >= len(p.toks) {
		return "", false
	}
	t := p.toks[p.i]
	if t.kind != 'o' && t.kind != 'w' {
		return "", false
	}
	for _, s := range ss {
		if t.s == s {
			p.i++
			return s, true
		}
	}
	return "", false
}

func (p *calcParser) or() (any, error) {
	return p.logic(p.and, "or", "||")
}

func (p *calcParser) and() (any, error) {
	return p.logic(p.not, "and", "&&")
}

func (p *calcParser) logic(next func() (any, error), ops ...string) (any, error) {
	a, err := next()
	if err != nil {
		return nil, err
	}
	for {
		op, ok := p.eat(ops...)
		if !ok {
			return a, nil
		}
		b, err := next()
		if err != nil {
			return nil, err
		}
		x, err1 := calcBool(a, op)
		y, err2 := calcBool(b, op)
		if err := errors.Join(err1, err2); err != nil {
			return nil, err
		}
		if ops[0] == "or" {
			a = x || y
		} else {
			a = x && y
		}
	}
}

func (p *calcParser) not() (any, error) {
	if op, ok := p.eat("not", "!"); ok {
		v, err := p.not()
		if err != nil {
			return nil, err
		}
		b, err := calcBool(v, op)
		return !b, err
	}
	return p.compare()
}

func (p *calcParser) compare() (any, error) {
	a, err := p.sum()
	if err != nil {
		return nil, err
	}
	op, ok := p.eat("==", "!=", "<=", ">=", "<", ">")
	if !ok {
		return a, nil
	}
	b, err := p.sum()
	if err != nil {
		return nil, err
	}
	switch op {
	case "==":
		return a == b, nil
	case "!=":
		return a != b, nil
	}
	x, y, err := calcNums(a, b, op)
	if err != nil {
		return nil, err
	}
	return map[string]bool{"<": x < y, "<=": x <= y, ">": x > y, ">=": x >= y}[op], nil
}

func (p *calcParser) sum() (any, error) {
	return p.arith(p.product, "+", "-")
}

func (p *calcParser) product() (any, error) {
	return p.arith(p.unary, "*", "/", "%")
}

func (p *calcParser) arith(next func() (any, error), ops ...string) (any, error) {
	a, err := next()
	if err != nil {
		return nil, err
	}
	for {
		op, ok := p.eat(ops...)
		if !ok {
			return a, nil
		}
		b, err := next()
		if err != nil {
			return nil, err
		}
		x, y, err := calcNums(a, b, op)
		if err != nil {
			return nil, err
		}
		switch op {
		case "+":
			a = x + y
		case "-":
			a = x - y
		case "*":
			a = x * y
		default:
			if y == 0 {
				return nil, errors.New("chia cho 0")
			}
			if op == "/" {
				a = x / y
			} else {
				a = math.Mod(x, y)
			}
		}
	}
}

func (p *calcParser) unary() (any, error) {
	if _, ok := p.eat("-"); ok {
		v, err := p.unary()
		if err != nil {
			return nil, err
		}
		x, _, err := calcNums(v, 0.0, "-")
		return -x, err
	}
	return p.primary()
}

func (p *calcParser) primary() (any, error) {
	if p.i >= len(p.toks) {
		return nil, synErr("biểu thức thiếu số ở cuối")
	}
	t := p.toks[p.i]
	p.i++
	switch {
	case t.kind == 'n':
		return t.n, nil
	case t.kind == 's':
		return t.s, nil
	case t.kind == 'o' && t.s == "(":
		v, err := p.or()
		if err != nil {
			return nil, err
		}
		if _, ok := p.eat(")"); !ok {
			return nil, synErr("thiếu dấu )")
		}
		return v, nil
	case t.kind == 'w' && (t.s == "true" || t.s == "false"):
		return t.s == "true", nil
	case t.kind == 'w':
		return p.call(t.s)
	}
	return nil, synErr("thiếu số trước “%s”", t.s)
}

// call – hàm round / floor / ceil / abs / min / max, đối số là số.
func (p *calcParser) call(name string) (any, error) {
	arity := map[string][2]int{"round": {1, 2}, "floor": {1, 1}, "ceil": {1, 1}, "abs": {1, 1}, "min": {1, 99}, "max": {1, 99}}
	lim, ok := arity[name]
	if !ok {
		return nil, synErr("không có hàm / từ “%s” (có: round floor ceil abs min max, and or not, true false)", name)
	}
	if _, ok := p.eat("("); !ok {
		return nil, synErr("%s cần ( sau tên", name)
	}
	var args []float64
	for {
		v, err := p.or()
		if err != nil {
			return nil, err
		}
		x, _, err := calcNums(v, 0.0, name)
		if err != nil {
			return nil, err
		}
		args = append(args, x)
		if _, ok := p.eat(","); !ok {
			break
		}
	}
	if _, ok := p.eat(")"); !ok {
		return nil, synErr("%s thiếu dấu )", name)
	}
	if len(args) < lim[0] || len(args) > lim[1] {
		return nil, synErr("%s nhận %d–%d số, có %d", name, lim[0], lim[1], len(args))
	}
	x := args[0]
	switch name {
	case "round":
		k := 1.0
		if len(args) == 2 {
			k = math.Pow(10, math.Round(args[1]))
		}
		return calcRound(x*k) / k, nil
	case "floor":
		return math.Floor(x), nil
	case "ceil":
		return math.Ceil(x), nil
	case "abs":
		return math.Abs(x), nil
	}
	for _, a := range args[1:] {
		if name == "min" {
			x = math.Min(x, a)
		} else {
			x = math.Max(x, a)
		}
	}
	return x, nil
}
