// Package modeltest – dựng node / cạnh / điều kiện ngắn gọn cho test (script, engine, root).
package modeltest

import "autoclicker/internal/model"

// Step – node mặc định của loại t (NewStep), không chờ sau khi xong; mod sửa thêm.
func Step(t string, mod func(*model.Step)) model.Step {
	s, err := model.NewStep(t)
	if err != nil {
		panic(err)
	}
	s.DelayAfter = 0
	if mod != nil {
		mod(&s)
	}
	return s
}

// Node – Step với id + ref cố định cho test flow.
func Node(id, ref, typ string, mod func(*model.Step)) model.Step {
	return Step(typ, func(s *model.Step) {
		s.ID, s.Ref = id, ref
		if mod != nil {
			mod(s)
		}
	})
}

// Edge – cạnh from.port → to, id theo from + port.
func Edge(from, port, to string) model.Edge {
	return model.Edge{ID: from + "." + port, From: from, Port: port, To: to}
}

// TextCond – điều kiện "thấy chữ t".
func TextCond(t string) *model.Cond { return &model.Cond{Kind: model.CondText, FindText: t} }
