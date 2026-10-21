package engine

import (
	"fmt"
	"time"

	"ghosthand/internal/model"
	"ghosthand/internal/textop"
)

// Node xử lý chữ (bước 42 G4b): Lọc mẫu, Tách chữ, Làm sạch chữ, Số & ngày – logic ở internal/textop.

// RunTextNode chạy node xử lý chữ s (đã thay biến) → cổng + output. Không khớp / không có số / không đọc được ngày →
// notFoundErr (cổng fail); regex sai → invalidErr. App.PreviewText dùng cho nút Chạy thử của inspector.
func RunTextNode(s model.Step, now time.Time) (string, map[string]any, error) {
	switch s.Type {
	case model.StepTextMatch:
		m, ok, err := textop.Match(s.Value, s.Pattern, s.MatchAll)
		if err != nil {
			return "", nil, &invalidErr{fmt.Sprintf("mẫu regex sai: %v", err)}
		}
		if !ok {
			return "", nil, &notFoundErr{msg: fmt.Sprintf("không có chỗ khớp mẫu “%s” trong “%s”", s.Pattern, clip(s.Value)), detail: map[string]any{"text": clip(s.Value)}}
		}
		return model.PortSuccess, map[string]any{"text": m.Text, "items": anyList(m.Items), "count": len(m.Items), "groups": anyList(m.Groups)}, nil

	case model.StepTextSplit:
		items := textop.Split(s.Value, s.Separator)
		return model.PortSuccess, map[string]any{"items": anyList(items), "count": len(items)}, nil

	case model.StepTextClean:
		t, err := textop.Clean(s.Value, s)
		if err != nil {
			return "", nil, &invalidErr{fmt.Sprintf("mẫu regex sai: %v", err)}
		}
		return model.PortSuccess, map[string]any{"text": t}, nil

	case model.StepTextNumber:
		return textNumber(s, now)
	}
	return "", nil, fmt.Errorf("không phải node xử lý chữ: %s", s.Type)
}

func textNumber(s model.Step, now time.Time) (string, map[string]any, error) {
	if s.NumOp == model.NumDate {
		t, err := textop.ParseDate(s.Value, now)
		if err != nil {
			return "", nil, &notFoundErr{msg: fmt.Sprintf("đọc ngày “%s”: %v", clip(s.Value), err), detail: map[string]any{"text": clip(s.Value)}}
		}
		return model.PortSuccess, map[string]any{"text": textop.FormatDate(t, s.DateLayout), "timestamp": t.Unix()}, nil
	}
	f, ok := textop.ExtractNumber(s.Value)
	if !ok {
		return "", nil, &notFoundErr{msg: fmt.Sprintf("không có số trong “%s”", clip(s.Value)), detail: map[string]any{"text": clip(s.Value)}}
	}
	text := fmt.Sprint(numValue(f))
	if s.NumOp == model.NumFormat {
		text = textop.FormatNumber(f, s.Decimals, s.NumStyle)
	}
	return model.PortSuccess, map[string]any{"value": numValue(f), "text": text}, nil
}

// clip – chữ dài cắt bớt cho lời báo lỗi / nhật ký.
func clip(s string) string {
	if r := []rune(s); len(r) > 80 {
		return string(r[:79]) + "…"
	}
	return s
}

func anyList(items []string) []any {
	out := make([]any, len(items))
	for i, v := range items {
		out[i] = v
	}
	return out
}
