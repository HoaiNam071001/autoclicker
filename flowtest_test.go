package main

import (
	"slices"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Helper test: dựng flow từ danh sách bước theo thứ tự mảng (success → node kế), đặt ref n1, n2…, click theo ảnh chưa
// chọn node nguồn thì lấy node chờ ảnh / chữ gần nhất phía trước. Kịch bản đã có start / cạnh giữ nguyên.
// Test viết gọn theo kiểu danh sách; app thật luôn gửi flow đủ start + cạnh (canvas).
func linear(sc model.Script) model.Script {
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
