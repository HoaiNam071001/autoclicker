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
