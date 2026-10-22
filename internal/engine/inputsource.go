package engine

import (
	"context"
	"fmt"
	"strings"

	"ghosthand/internal/inputsource"
	"ghosthand/internal/model"
)

// Node Bộ gõ (bước 43 S6, S6b): đọc bộ gõ hiện tại / đổi sang bộ gõ khớp InputLayout (inputsource.Find). Phần mềm gõ tiếng
// Việt bên thứ ba (Unikey, EVKey…) được *phát hiện* (output imeApp) và đảo bật / tắt bằng phím tắt riêng của chúng
// (toggleIME) – chúng không có API đọc trạng thái nên không biết đang bật hay tắt.

const inputVerifyTries = 10 // Windows đổi bộ gõ bất đồng bộ: chờ tối đa ~1 s cho tới khi thấy đã đổi

func (e *Engine) inputSource(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	if s.InputOp == model.InputToggleIME {
		return e.toggleIME(ctx, s)
	}
	list, err := e.inputs.List()
	if err != nil {
		return nil, err
	}
	cur, err := e.inputs.Current()
	if err != nil {
		return nil, err
	}
	prev := cur
	if s.InputOp == model.InputSet {
		if cur, err = e.switchInput(ctx, s, list, cur); err != nil {
			return nil, err
		}
	}
	names := make([]any, len(list))
	for i, l := range list {
		names[i] = l.Name
	}
	return map[string]any{
		"id": cur.ID, "name": cur.Name, "lang": cur.Lang, "vietnamese": inputsource.IsVietnamese(cur),
		"previousId": prev.ID, "changed": prev.ID != cur.ID, "imeApp": e.findIME(), "items": names,
	}, nil
}

func (e *Engine) switchInput(ctx context.Context, s model.Step, list []inputsource.Source, cur inputsource.Source) (inputsource.Source, error) {
	i := inputsource.Find(list, s.InputLayout)
	if i < 0 {
		ids := make([]string, len(list))
		for j, l := range list {
			ids[j] = l.Name + " (" + l.ID + ")"
		}
		return cur, &notFoundErr{
			msg:    fmt.Sprintf("không có bộ gõ %q – máy có: %s", s.InputLayout, strings.Join(ids, ", ")),
			detail: map[string]any{"layout": s.InputLayout, "available": ids},
		}
	}
	want := list[i]
	if want.ID == cur.ID {
		return cur, nil
	}
	if err := e.inputs.Select(want.ID); err != nil {
		return cur, err
	}
	for range inputVerifyTries {
		if now, err := e.inputs.Current(); err == nil && now.ID == want.ID {
			return now, nil
		}
		if err := Sleep(ctx, 100); err != nil {
			return cur, err
		}
	}
	return cur, fmt.Errorf("đã yêu cầu đổi sang %s nhưng bộ gõ vẫn chưa đổi (cửa sổ đang ở trên cùng có thể không nhận)", want.Name)
}

// findIME – tên phần mềm gõ tiếng Việt đang chạy ("" = không có hoặc không đọc được danh sách app).
func (e *Engine) findIME() string {
	apps, _ := e.apps.List()
	names := make([]string, len(apps))
	for i, a := range apps {
		names[i] = a.Name
	}
	return inputsource.FindIME(names)
}

// toggleIME bấm phím tắt đổi chế độ gõ của phần mềm gõ đang chạy – đúng MỘT lần.
func (e *Engine) toggleIME(ctx context.Context, s model.Step) (map[string]any, error) {
	ime := e.findIME()
	if ime == "" {
		return nil, &notFoundErr{msg: "không thấy Unikey / EVKey / OpenKey đang chạy – chỉ bật / tắt được phần mềm gõ đang mở", detail: map[string]any{"kind": "ime"}}
	}
	hotkey := strings.TrimSpace(s.InputHotkey)
	if hotkey == "" {
		hotkey = model.InputDefaultHotkey
	}
	mods, key, err := inputsource.ParseHotkey(hotkey)
	if err != nil {
		return nil, &invalidErr{err.Error()}
	}
	press := model.Step{KeyMode: model.KeyChord, Modifiers: mods}
	if key != "" {
		press.Keys = []string{key}
	}
	if err := e.pressKeys(ctx, press); err != nil {
		return nil, err
	}
	return map[string]any{"imeApp": ime, "hotkey": hotkey, "toggled": true}, nil
}
