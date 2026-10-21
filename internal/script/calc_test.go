package script

import (
	"encoding/json"
	"os"
	"testing"
)

// Cùng dữ liệu với frontend/src/flow/calc.test.ts (testdata/calc.json).
func TestEvalCalc(t *testing.T) {
	b, err := os.ReadFile("../../testdata/calc.json")
	if err != nil {
		t.Fatal(err)
	}
	var data struct {
		Eval []struct {
			Expr  string
			Value any
			Error bool
		}
		Quote []struct{ In, Out string }
	}
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatal(err)
	}
	for _, c := range data.Eval {
		got, err := EvalCalc(c.Expr)
		if (err != nil) != c.Error {
			t.Errorf("EvalCalc(%q): lỗi %v, muốn lỗi = %v", c.Expr, err, c.Error)
			continue
		}
		if !c.Error && got != c.Value {
			t.Errorf("EvalCalc(%q) = %#v, muốn %#v", c.Expr, got, c.Value)
		}
	}
	for _, c := range data.Quote {
		if got := CalcQuote(c.In); got != c.Out {
			t.Errorf("CalcQuote(%q) = %s, muốn %s", c.In, got, c.Out)
		}
	}
}

func TestCheckCalc(t *testing.T) {
	for expr, bad := range map[string]bool{
		"{{n2.output.count}} * 2 + 1":  false,
		"{{n2.output.text}} == \"OK\"": false,
		"10 / {{vars.n}}":              false, // chia 0 chỉ biết khi chạy
		"{{n2.output.count}} *":        true,
		"round({{price}}, 2, 3)":       true,
		"{{n2.output.value}} >> 2":     true,
	} {
		if err := CheckCalc(expr); (err != nil) != bad {
			t.Errorf("CheckCalc(%q) = %v, muốn lỗi = %v", expr, err, bad)
		}
	}
}
