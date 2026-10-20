package script

import "autoclicker/internal/model"

// RunFrom – bản sao kịch bản để chạy từ node id (debug / AI): một vòng, bắt đầu ở id, các node trước không chạy
// (biến của chúng trống). only: chỉ chạy riêng node đó – bỏ các cạnh ra của nó nên chạy xong là hết vòng
// (thất bại thì báo lỗi như cổng fail không nối). Kiểm tra id có thật ở ValidateScript (start không có node).
func RunFrom(sc model.Script, id string, only bool) model.Script {
	sc.Start, sc.Loops = id, 1
	if !only {
		return sc
	}
	edges := make([]model.Edge, 0, len(sc.Edges))
	for _, e := range sc.Edges {
		if e.From != id {
			edges = append(edges, e)
		}
	}
	sc.Edges = edges
	return sc
}

// RunPart – bản sao kịch bản chỉ chạy cụm node ids (AI sửa một đoạn – bước 31b): một vòng, chỉ giữ cạnh nối hai node
// trong cụm (đi ra ngoài cụm = hết vòng), bắt đầu ở from nếu có, không thì node đầu cụm theo FlowOrder.
func RunPart(sc model.Script, from string, ids []string) model.Script {
	in := make(map[string]bool, len(ids))
	for _, id := range ids {
		in[id] = true
	}
	if from == "" {
		for _, id := range model.FlowOrder(sc) {
			if in[id] {
				from = id
				break
			}
		}
	}
	in[from] = true
	sc.Start, sc.Loops = from, 1
	edges := make([]model.Edge, 0, len(sc.Edges))
	for _, e := range sc.Edges {
		if in[e.From] && in[e.To] {
			edges = append(edges, e)
		}
	}
	sc.Edges = edges
	return sc
}
