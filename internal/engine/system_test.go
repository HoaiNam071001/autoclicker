package engine

import (
	"testing"

	"ghosthand/internal/device/devicetest"
	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func TestSystemVolume(t *testing.T) {
	audio := &devicetest.FakeAudio{Level: 95, Muted: true}
	sys := func(op string, v int) model.Step {
		return modeltest.Step(model.StepSystem, func(s *model.Step) { s.SysOp, s.Volume = op, v })
	}
	e := New(&devicetest.FakeInput{}, model.Script{}, Options{})
	e.audio = audio
	for _, c := range []struct {
		step  model.Step
		level int
		muted bool
	}{
		{sys(model.SysVolumeUp, 0), 100, false}, // +10 nhưng giữ ≤ 100, bật tiếng lại
		{sys(model.SysVolumeDown, 30), 70, false},
		{sys(model.SysMute, 0), 70, true},
		{sys(model.SysUnmute, 0), 70, false},
		{sys(model.SysVolume, 25), 25, false},
		{sys(model.SysVolumeDown, 40), 0, false},
	} {
		out, err := e.systemStep(t.Context(), c.step)
		if err != nil || out["volume"] != c.level || out["muted"] != c.muted {
			t.Errorf("%s %d: %v %v, muốn %d / %v", c.step.SysOp, c.step.Volume, out, err, c.level, c.muted)
		}
	}
}
