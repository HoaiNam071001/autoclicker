package engine

import (
	"context"
	"fmt"
	"strings"

	"ghosthand/internal/humanize"
	"ghosthand/internal/model"
)

// Gõ chữ (node Gõ chữ, Gõ vào ô web): gõ từng phím hoặc dán qua clipboard (Step.TypeMode). Bộ gõ tiếng Việt chạy nền
// (UniKey, EVKey…) xử lý lại cả phím Unicode app gửi qua SendInput → chữ méo ("mình" → "mmmmình"); dán thì không đi qua
// bộ gõ.

// pasteWaitMs – chờ app đọc xong clipboard sau Ctrl/Cmd+V rồi mới trả chữ cũ (Chrome đọc clipboard không đồng bộ).
const pasteWaitMs = 200

// imeApps – tên tiến trình bộ gõ tiếng Việt (chữ thường, so "chứa").
var imeApps = []string{"unikey", "evkey", "openkey", "gotiengviet", "vietkey"}

// usePaste – cách nhập thật của một lần gõ: TypeAuto dán khi pasteAuto (Windows) và không giống người (dán thì mất nhịp gõ).
func usePaste(mode string, pasteAuto, human bool) bool {
	switch mode {
	case model.TypePaste:
		return true
	case model.TypeKeys:
		return false
	}
	return pasteAuto && !human
}

// typeText nhập text vào ô đang focus theo mode (TypeModes).
func (e *Engine) typeText(ctx context.Context, text, mode string) error {
	if text == "" {
		return nil
	}
	if usePaste(mode, e.pasteAuto, e.human != nil) {
		return e.pasteText(ctx, text)
	}
	if ime := e.runningIME(); ime != "" {
		e.note(fmt.Sprintf("bộ gõ %s đang chạy có thể làm méo chữ khi gõ từng phím – tắt tiếng Việt (vd Ctrl+Shift) hoặc chọn cách nhập Dán", ime))
	}
	return e.typeKeys(ctx, text)
}

// typeKeys gõ từng ký tự để dừng / tạm dừng được giữa chừng với chuỗi dài; giống người thì theo nhịp humanize.Typing
// (nhịp không đều, ngừng, gõ nhầm rồi xoá).
func (e *Engine) typeKeys(ctx context.Context, text string) error {
	var keys []humanize.Key
	if e.human != nil {
		keys = humanize.Typing(e.rng, *e.human, text)
	} else {
		for _, r := range text {
			keys = append(keys, humanize.Key{Text: string(r)})
		}
	}
	for _, k := range keys {
		if err := Sleep(ctx, k.WaitMs); err != nil {
			return err
		}
		if err := e.gate.wait(ctx); err != nil {
			return err
		}
		if k.Text == humanize.Backspace {
			if err := e.pressKeys(ctx, model.Step{Keys: []string{"backspace"}}); err != nil {
				return err
			}
			continue
		}
		if err := e.in.TypeStr(k.Text); err != nil {
			return err
		}
	}
	return nil
}

// pasteText ghi text vào clipboard, bấm Ctrl/Cmd+V, rồi trả lại chữ cũ trong clipboard (clipboard cũ không phải chữ –
// ảnh, file – thì không khôi phục được, để nguyên text).
func (e *Engine) pasteText(ctx context.Context, text string) error {
	old, _ := e.clip.ReadText()
	if err := e.clip.WriteText(text); err != nil {
		return fmt.Errorf("không ghi được clipboard để dán: %w", err)
	}
	paste := model.Step{Keys: []string{"v"}, Modifiers: []string{shortcutMod()}, KeyMode: model.KeyChord}
	err := e.pressKeys(ctx, paste)
	if err == nil {
		err = Sleep(ctx, pasteWaitMs)
	}
	if old != "" {
		e.clip.WriteText(old)
	}
	return err
}

// runningIME – tên bộ gõ tiếng Việt đang chạy ("" = không có / không kiểm tra); tra một lần mỗi lần chạy.
func (e *Engine) runningIME() string {
	if !e.imeCheck {
		return ""
	}
	if e.ime == nil {
		name := ""
		if apps, err := e.apps.List(); err == nil {
			for _, a := range apps {
				n := strings.ToLower(a.Name)
				for _, ime := range imeApps {
					if strings.Contains(n, ime) {
						name = a.Name
					}
				}
			}
		}
		e.ime = &name
	}
	return *e.ime
}

// note – lời nhắc của node đang chạy; execNode ghi vào nhật ký node (mức warn).
func (e *Engine) note(msg string) { e.notes = append(e.notes, msg) }
