package engine

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ghosthand/internal/display"
	"ghosthand/internal/model"
)

// Node Màn hình / âm thanh (bước 43 S7): chế độ sáng / tối, độ sáng, loa + micro mặc định, phát âm báo. Âm lượng ở node
// Hệ thống. Đổi xong đọc lại để chắc đã đổi.

func (e *Engine) displayStep(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	v := strings.TrimSpace(s.DisplayValue)
	var out map[string]any
	var err error
	switch s.DisplayOp {
	case model.DisplayTheme:
		out, err = e.setTheme(v)
	case model.DisplayBrightness:
		out, err = e.setBrightness(v)
	case model.DisplayAudioOut, model.DisplayAudioIn:
		out, err = e.setAudio(v, s.DisplayOp == model.DisplayAudioIn)
	case model.DisplaySound:
		err = e.disp.PlaySound(ctx, e.soundTarget(v), s.SoundWait)
		out = map[string]any{"sound": v, "waited": s.SoundWait}
	default:
		out, err = e.readDisplay()
	}
	switch {
	case ctx.Err() != nil:
		return nil, ctx.Err()
	case errors.Is(err, display.ErrUnsupported):
		return nil, &invalidErr{err.Error()}
	case errors.Is(err, display.ErrNoBrightness):
		return nil, &invalidErr{err.Error()}
	}
	return out, err
}

func (e *Engine) soundTarget(v string) string {
	if strings.ContainsAny(v, `/\`) || strings.HasPrefix(v, "~") {
		if p, err := userPath(v); err == nil {
			return p
		}
	}
	return v
}

func (e *Engine) setTheme(v string) (map[string]any, error) {
	prev, err := e.disp.Dark()
	if err != nil {
		return nil, err
	}
	dark := prev
	switch v {
	case "light":
		dark = false
	case "dark":
		dark = true
	case "toggle":
		dark = !prev
	default:
		return nil, &invalidErr{fmt.Sprintf("chế độ %q không hợp lệ – light, dark hoặc toggle", v)}
	}
	if dark != prev {
		if err := e.disp.SetDark(dark); err != nil {
			return nil, err
		}
	}
	now, err := e.disp.Dark()
	if err != nil {
		return nil, err
	}
	if now != dark {
		return nil, fmt.Errorf("đã yêu cầu đổi chế độ nhưng hệ điều hành vẫn chưa đổi")
	}
	return map[string]any{"dark": now, "theme": themeName(now), "changed": prev != now}, nil
}

func themeName(dark bool) string {
	if dark {
		return "dark"
	}
	return "light"
}

func (e *Engine) setBrightness(v string) (map[string]any, error) {
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 || n > 100 {
		return nil, &invalidErr{fmt.Sprintf("độ sáng %q không hợp lệ – số nguyên 0–100", v)}
	}
	prev, _ := e.disp.Brightness()
	if err := e.disp.SetBrightness(n); err != nil {
		return nil, err
	}
	now, err := e.disp.Brightness()
	if err != nil {
		now = n
	}
	return map[string]any{"brightness": now, "previous": prev}, nil
}

func (e *Engine) setAudio(query string, input bool) (map[string]any, error) {
	list, err := e.disp.AudioDevices()
	if err != nil {
		return nil, err
	}
	i := display.FindAudio(list, query, input)
	if i < 0 {
		var names []string
		for _, a := range list {
			if a.Input == input {
				names = append(names, a.Name)
			}
		}
		kind := "loa"
		if input {
			kind = "micro"
		}
		return nil, &notFoundErr{
			msg:    fmt.Sprintf("không có %s %q – máy có: %s", kind, query, strings.Join(names, ", ")),
			detail: map[string]any{"device": query, "available": names},
		}
	}
	prev := ""
	for _, a := range list {
		if a.Input == input && a.Default {
			prev = a.Name
		}
	}
	want := list[i]
	if err := e.disp.SetDefaultAudio(want.ID, input); err != nil {
		return nil, err
	}
	return map[string]any{"name": want.Name, "id": want.ID, "previous": prev}, nil
}

func (e *Engine) readDisplay() (map[string]any, error) {
	dark, err := e.disp.Dark()
	if err != nil {
		return nil, err
	}
	bright, berr := e.disp.Brightness()
	if berr != nil {
		bright = -1
	}
	out := map[string]any{"dark": dark, "theme": themeName(dark), "brightness": bright, "outputName": "", "inputName": ""}
	if list, err := e.disp.AudioDevices(); err == nil {
		items := []any{}
		for _, a := range list {
			kind := "loa"
			if a.Input {
				kind = "micro"
			}
			items = append(items, kind+": "+a.Name)
			if a.Default {
				if a.Input {
					out["inputName"] = a.Name
				} else {
					out["outputName"] = a.Name
				}
			}
		}
		out["items"] = items
	}
	return out, nil
}
