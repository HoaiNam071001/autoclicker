package engine

import (
	"context"

	"ghosthand/internal/model"
)

// Node Hệ thống (bước 35 N5): âm lượng loa – đặt mức / tăng / giảm (bật tiếng lại), tắt / bật tiếng.
func (e *Engine) systemStep(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	var err error
	switch s.SysOp {
	case model.SysMute, model.SysUnmute:
		err = e.audio.SetMuted(s.SysOp == model.SysMute)
	case model.SysVolumeUp, model.SysVolumeDown:
		step := s.Volume
		if step == 0 {
			step = 10
		}
		if s.SysOp == model.SysVolumeDown {
			step = -step
		}
		var cur int
		if cur, _, err = e.audio.Volume(); err == nil {
			err = e.audio.SetVolume(min(max(cur+step, 0), 100))
		}
	default:
		err = e.audio.SetVolume(min(max(s.Volume, 0), 100))
	}
	if err != nil {
		return nil, err
	}
	level, muted, err := e.audio.Volume()
	if err != nil {
		return nil, err
	}
	return map[string]any{"volume": level, "muted": muted}, nil
}
