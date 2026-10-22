package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ghosthand/internal/element"
	"ghosthand/internal/model"
)

// Node Thao tác phần tử (bước 43 E1): tìm phần tử rồi bấm / đặt giá trị / tick / mở rộng / chọn / cuộn thẳng qua cây
// accessibility (AXPress, AXValue… – device.ElementReader.Act), không dùng chuột nên cửa sổ bị che một phần vẫn chạy.
// Thao tác có trạng thái đọc lại được (đặt giá trị, tick) thì đọc lại để kiểm; app không nhận → thử lại tới hết giờ.

func (e *Engine) elementAction(ctx context.Context, s model.Step) (map[string]any, error) {
	q, err := elementQuery(s.Selectors, !s.TextAccent, s.Occurrence, nil)
	if err != nil {
		return nil, err
	}
	act := s.ElementAct
	if act == "" {
		act = model.ElActPress
	}
	var done element.Node
	why := ""
	var miss func() map[string]any
	_, err = e.waitLoop(ctx, s, "element", func() (bool, int, error) {
		nodes, r, err := e.lookElement(s.ElementApp, q)
		switch {
		case err != nil && waitable(err):
			why = err.Error()
			return false, 0, nil
		case err != nil:
			return false, 0, err
		case !r.Found:
			why, miss = missWhy(r, nodes, max(s.Occurrence, 1)), func() map[string]any { return elementMiss(s.Selectors, nodes, r, q) }
			return false, 0, nil
		}
		n := r.Hit.Node
		if (act == model.ElActCheck || act == model.ElActUncheck) && isChecked(n) == (act == model.ElActCheck) {
			done = n // đã đúng trạng thái
			return true, 100, nil
		}
		realAct := act
		if act == model.ElActCheck || act == model.ElActUncheck {
			realAct = model.ElActPress // checkbox / switch: bấm là đảo
		}
		if err := e.elems.Act(s.ElementApp, r.Hit.Index, realAct, s.Text); err != nil {
			why = err.Error()
			return false, 0, nil
		}
		if err := Sleep(ctx, settleMs); err != nil {
			return false, 0, err
		}
		ok, after, bad := e.verifyAction(act, s.Text, s.ElementApp, q)
		if !ok {
			why = bad
			return false, 0, nil
		}
		if after.Role != "" {
			n = after
		}
		done = n
		return true, 100, nil
	})
	if errors.Is(err, errTimeout) {
		return nil, &notFoundErr{msg: fmt.Sprintf("không thao tác được phần tử %s sau %s – %s", strings.Join(s.Selectors, " / "), fmtSeconds(s.TimeoutMs), why), detail: lazyDetail(miss)}
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{"action": act, "role": done.Role, "text": done.Text(), "value": done.Value}, nil
}

// verifyAction đọc lại phần tử sau thao tác: ok = đã đúng ý (hoặc không kiểm được), bad = vì sao chưa.
// Phần tử biến mất / đổi chữ sau thao tác (bấm xong trang đổi) coi là ổn – chỉ kiểm khi còn thấy lại.
func (e *Engine) verifyAction(act, text, app string, q element.Query) (ok bool, n element.Node, bad string) {
	_, r, err := e.lookElement(app, q)
	if err != nil || !r.Found {
		return true, element.Node{}, ""
	}
	n = r.Hit.Node
	switch act {
	case model.ElActSetValue:
		if n.Value != text {
			return false, n, fmt.Sprintf("app nhận lệnh nhưng giá trị vẫn là %q (ô web dựng bằng React có thể không nhận – dùng node Gõ chữ)", n.Value)
		}
	case model.ElActCheck, model.ElActUncheck:
		if isChecked(n) != (act == model.ElActCheck) {
			return false, n, "bấm xong nhưng trạng thái tick không đổi"
		}
	}
	return true, n, ""
}

// isChecked – checkbox / switch / radio đang bật: giá trị AX 1, "true", "on".
func isChecked(n element.Node) bool {
	switch strings.ToLower(strings.TrimSpace(n.Value)) {
	case "1", "true", "on", "checked":
		return true
	}
	return false
}
