package model

import (
	"fmt"
	"regexp"
	"strings"
)

// Selector của node web (bước 25d): một dòng Step.Selectors = "<cách>=<giá trị>", không tiền tố = text.
// Thứ tự ưu tiên theo độ bền: chữ hiển thị / aria / role → css / xpath (tránh class sinh tự động).
const (
	SelText  = "text"  // chữ hiển thị chứa giá trị (không phân biệt hoa thường / khoảng trắng); "…" trong nháy kép = đúng cả chuỗi
	SelAria  = "aria"  // aria-label / aria-labelledby, cùng cách so như text
	SelRole  = "role"  // vai trò (thẻ hoặc thuộc tính role): role=button, role=link[name="Xem lịch sử"] (tên chứa, như text)
	SelCSS   = "css"   // CSS selector
	SelXPath = "xpath" // XPath
)

// SelectorKinds – các tiền tố hợp lệ (frontend có bản sao trong flow/web.ts).
var SelectorKinds = []string{SelText, SelAria, SelRole, SelCSS, SelXPath}

type Selector struct {
	By    string `json:"by"`
	Value string `json:"value"`
	Exact bool   `json:"exact,omitempty"` // text / aria: giá trị trong nháy kép
	Name  string `json:"name,omitempty"`  // role: tên truy cập phải chứa (chữ / aria-label) – kiểu Playwright [name="…"]
}

// roleRe – "button" hoặc "button[name="Gửi"]" / "button[name=Gửi]".
var roleRe = regexp.MustCompile(`^([a-zA-Z]+)\s*(?:\[\s*name\s*=\s*(?:"([^"]*)"|([^\]"]*?))\s*\])?$`)

// ParseSelector đọc một dòng selector. Giá trị trống → lỗi.
func ParseSelector(line string) (Selector, error) {
	line = strings.TrimSpace(line)
	sel := Selector{By: SelText, Value: line}
	if by, v, ok := strings.Cut(line, "="); ok {
		for _, k := range SelectorKinds {
			if strings.EqualFold(strings.TrimSpace(by), k) {
				sel = Selector{By: k, Value: strings.TrimSpace(v)}
			}
		}
	}
	if sel.By == SelText || sel.By == SelAria {
		if n := len(sel.Value); n >= 2 && sel.Value[0] == '"' && sel.Value[n-1] == '"' {
			sel.Value, sel.Exact = sel.Value[1:n-1], true
		}
	}
	if sel.By == SelRole && sel.Value != "" {
		m := roleRe.FindStringSubmatch(sel.Value)
		if m == nil {
			return sel, fmt.Errorf("selector “%s” sai dạng – viết role=button hoặc role=button[name=\"Gửi\"]", line)
		}
		sel.Value, sel.Name = strings.ToLower(m[1]), strings.TrimSpace(m[2]+m[3])
	}
	if sel.Value == "" {
		return sel, fmt.Errorf("selector “%s” chưa có giá trị", line)
	}
	return sel, nil
}
