package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"ghosthand/internal/element"
	"ghosthand/internal/model"
)

// Node Đọc tài liệu (bước 43 E3): toàn bộ chữ của ô soạn thảo / tài liệu đang mở (TextEdit, Notepad, ô nhiều dòng)
// qua giá trị accessibility – giữ nguyên xuống dòng, không rút gọn như Đọc chữ trang. selectionOnly = chỉ phần bôi đen.

func (e *Engine) readDocument(ctx context.Context, s model.Step) (map[string]any, error) {
	q, scoped, err := pageScope(s)
	if err != nil {
		return nil, err
	}
	var text string
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
		idx := -1
		if scoped {
			r := element.Find(nodes, q)
			if !r.Found {
				why, miss = missWhy(r, nodes, max(s.Occurrence, 1)), func() map[string]any { return elementMiss(s.Selectors, nodes, r, q) }
				return false, 0, nil
			}
			idx = r.Hit.Index
		} else if idx = element.DocumentIndex(nodes, s.SelectionOnly); idx < 0 {
			why = "không thấy ô soạn thảo nào có chữ trong cửa sổ (chỉ selectors để chỉ đúng ô)"
			if s.SelectionOnly {
				why = "chưa có đoạn chữ nào đang được chọn"
			}
			return false, 0, nil
		}
		if text = nodes[idx].Value; s.SelectionOnly {
			text = nodes[idx].Selected
		}
		if text == "" {
			why = "ô tìm được nhưng chưa có chữ"
			if s.SelectionOnly {
				why = "ô tìm được nhưng chưa chọn đoạn chữ nào"
			}
			return false, 0, nil
		}
		return true, 100, nil
	})
	if errors.Is(err, errTimeout) {
		return nil, &notFoundErr{msg: fmt.Sprintf("không đọc được tài liệu sau %s – %s", fmtSeconds(s.TimeoutMs), why), detail: lazyDetail(miss)}
	}
	if err != nil {
		return nil, err
	}
	text = strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	return map[string]any{"text": text, "length": utf8.RuneCountInString(text), "lines": strings.Count(text, "\n") + 1}, nil
}
