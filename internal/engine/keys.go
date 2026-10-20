package engine

import (
	"context"

	"autoclicker/internal/model"
)

// Bàn phím: node Phím tắt bấm tổ hợp (chord) hoặc lần lượt (sequence) – gõ chữ ở typing.go.

// Thời gian giữ khi bấm phím: chord giữ cả tổ hợp keyHoldMs; sequence nhấn mỗi phím keyHoldMs,
// cách phím sau keyGapMs.
const keyHoldMs, keyGapMs = 50, 50

// pressKeys bấm tổ hợp theo KeyMode; mọi phím đã nhấn luôn được nhả (ngược thứ tự), kể cả khi lỗi / bị huỷ.
// Thời gian giữ / nghỉ theo mức giống người của node (keyTiming).
func (e *Engine) pressKeys(ctx context.Context, s model.Step) error {
	in := e.in
	hold, gap := e.keyTiming()
	var held []string
	defer func() {
		for i := len(held) - 1; i >= 0; i-- {
			in.KeyToggle(held[i], false, nil)
		}
	}()
	press := func(k string, mods []string) error {
		if err := in.KeyToggle(k, true, mods); err != nil {
			return err
		}
		held = append(held, k)
		return nil
	}

	for _, m := range s.Modifiers {
		if err := press(m, nil); err != nil {
			return err
		}
	}
	if s.KeyMode == model.KeySequence {
		for i, k := range s.Keys {
			if err := press(k, s.Modifiers); err != nil {
				return err
			}
			if err := Sleep(ctx, hold()); err != nil {
				return err
			}
			held = held[:len(held)-1]
			if err := in.KeyToggle(k, false, nil); err != nil {
				return err
			}
			if i < len(s.Keys)-1 {
				if err := Sleep(ctx, gap()); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for _, k := range s.Keys {
		if err := press(k, s.Modifiers); err != nil {
			return err
		}
	}
	return Sleep(ctx, hold())
}
