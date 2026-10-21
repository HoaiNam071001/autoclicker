// Package scripttest – helper test dùng chung (engine, script_test).
package scripttest

import (
	"slices"

	"ghosthand/internal/model"
	"ghosthand/internal/model/modeltest"
	"ghosthand/internal/script"
)

// Helper test: dựng flow từ danh sách bước theo thứ tự mảng (success → node kế), đặt ref n1, n2…, click theo ảnh chưa
// chọn node nguồn thì lấy node chờ ảnh / chữ gần nhất phía trước. Kịch bản đã có start / cạnh giữ nguyên.
// Test viết gọn theo kiểu danh sách; app thật luôn gửi flow đủ start + cạnh (canvas).
func Linear(sc model.Script) model.Script {
	linearToFlow(&sc)
	return sc
}

func linearToFlow(sc *model.Script) {
	if sc.Start != "" || len(sc.Edges) > 0 || len(sc.Steps) == 0 {
		return
	}
	sc.Steps = slices.Clone(sc.Steps)
	script.FixRefs(sc.Steps)
	sc.Edges = nil
	sc.Start = sc.Steps[0].ID
	for i := 0; i+1 < len(sc.Steps); i++ {
		sc.Edges = append(sc.Edges, model.Edge{ID: "auto" + sc.Steps[i].ID, From: sc.Steps[i].ID, Port: model.PortSuccess, To: sc.Steps[i+1].ID})
	}
	src := ""
	for i := range sc.Steps {
		s := &sc.Steps[i]
		if s.Type == model.StepClickImage && s.ImageSource != model.ImageFromSearch && s.ImageFrom == "" {
			s.ImageFrom = src
		}
		if s.Enabled && (s.Type == model.StepWaitImage || s.Type == model.StepClickImage && s.ImageSource == model.ImageFromSearch) {
			src = s.ID
		}
	}
}

// TextScript – kịch bản: chờ chữ text (trong vùng 400×300) rồi các node steps.
func TextScript(text string, steps ...model.Step) model.Script {
	wait := modeltest.Step(model.StepWaitImage, func(s *model.Step) {
		s.FindBy, s.FindText, s.TimeoutMs, s.IntervalMs = model.FindByText, text, 2000, 1
	})
	return model.Script{Loops: 1, Region: &model.Rect{X: 100, Y: 50, W: 400, H: 300}, Steps: append([]model.Step{wait}, steps...)}
}
