package main

import "slices"

// Helper test: dựng flow từ danh sách bước theo thứ tự mảng (success → node kế), đặt ref n1, n2…, click theo ảnh chưa
// chọn node nguồn thì lấy node chờ ảnh / chữ gần nhất phía trước. Kịch bản đã có start / cạnh giữ nguyên.
// Test viết gọn theo kiểu danh sách; app thật luôn gửi flow đủ start + cạnh (canvas).
func linear(sc Script) Script {
	linearToFlow(&sc)
	return sc
}

func linearToFlow(sc *Script) {
	if sc.Start != "" || len(sc.Edges) > 0 || len(sc.Steps) == 0 {
		return
	}
	sc.Steps = slices.Clone(sc.Steps)
	fixRefs(sc.Steps)
	sc.Edges = nil
	sc.Start = sc.Steps[0].ID
	for i := 0; i+1 < len(sc.Steps); i++ {
		sc.Edges = append(sc.Edges, Edge{ID: "auto" + sc.Steps[i].ID, From: sc.Steps[i].ID, Port: PortSuccess, To: sc.Steps[i+1].ID})
	}
	src := ""
	for i := range sc.Steps {
		s := &sc.Steps[i]
		if s.Type == StepClickImage && s.ImageSource != ImageFromSearch && s.ImageFrom == "" {
			s.ImageFrom = src
		}
		if s.Enabled && (s.Type == StepWaitImage || s.Type == StepClickImage && s.ImageSource == ImageFromSearch) {
			src = s.ID
		}
	}
}
