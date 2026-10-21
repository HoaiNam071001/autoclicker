package script

import (
	"strings"
	"testing"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
)

func TestCheckShot(t *testing.T) {
	cases := map[string]func(*model.Step){
		"cách chờ màn hình không hợp lệ": func(s *model.Step) { s.WaitMode = "blink" },
		"0–100 %":                func(s *model.Step) { s.DiffPct = 120 },
		"vùng chụp không hợp lệ": func(s *model.Step) { s.ShotArea = "window" },
	}
	for want, f := range cases {
		s := modeltest.Step(model.StepWaitScreen, f)
		if err := checkStep(&s); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("muốn lỗi chứa %q, nhận %v", want, err)
		}
	}
}
