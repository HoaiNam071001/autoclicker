package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ghosthand/internal/element"
	"ghosthand/internal/model"
)

// Node Đọc chữ trang (bước 43 B4): đọc mọi chữ trong cửa sổ app (trang trình duyệt thường, app khác) hoặc trong một
// phần tử (selectors) qua cây accessibility, giữ thứ tự đọc. Chờ tới khi có chữ (trang còn đang tải thì chưa có).

// pageScope – phạm vi đọc: selectors trống = cả cửa sổ; có = cây con của phần tử khớp.
func pageScope(s model.Step) (q element.Query, scoped bool, err error) {
	for _, l := range s.Selectors {
		if strings.TrimSpace(l) != "" {
			q, err = elementQuery(s.Selectors, !s.TextAccent, s.Occurrence, nil)
			return q, true, err
		}
	}
	return q, false, nil
}

func (e *Engine) pageText(ctx context.Context, s model.Step) (map[string]any, error) {
	q, scoped, err := pageScope(s)
	if err != nil {
		return nil, err
	}
	var lines []string
	why := ""
	var miss func() map[string]any
	_, err = e.waitLoop(ctx, s, "element", func() (bool, int, error) {
		nodes, err := e.elems.Elements(s.ElementApp)
		if waitable(err) {
			why = err.Error()
			return false, 0, nil
		}
		if err != nil {
			return false, 0, err
		}
		from, to := 0, len(nodes)
		if scoped {
			r := element.Find(nodes, q)
			if !r.Found {
				why, miss = missWhy(r, nodes, max(s.Occurrence, 1)), func() map[string]any { return elementMiss(s.Selectors, nodes, r, q) }
				return false, 0, nil
			}
			from, to = r.Hit.Index, element.SubtreeEnd(nodes, r.Hit.Index)
		}
		lines = element.TextLines(nodes, from, to, s.PageAll)
		if len(lines) == 0 {
			why = fmt.Sprintf("đã đọc %d phần tử nhưng không có chữ nào đang hiện", to-from)
			return false, 0, nil
		}
		return true, 100, nil
	})
	if errors.Is(err, errTimeout) {
		return nil, &notFoundErr{msg: fmt.Sprintf("không đọc được chữ trang sau %s – %s", fmtSeconds(s.TimeoutMs), why), detail: lazyDetail(miss)}
	}
	if err != nil {
		return nil, err
	}
	items := make([]any, len(lines)) // []any: nối vào forEach / đọc lines.0 trong biến
	for i, l := range lines {
		items[i] = l
	}
	return map[string]any{"text": strings.Join(lines, "\n"), "lines": items, "count": len(lines)}, nil
}
