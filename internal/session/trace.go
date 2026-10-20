package session

import (
	"fmt"
	"image"
	"slices"

	"autoclicker/internal/engine"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

// Chạy một phần + vết chạy cho AI (bước 31b): AI chạy đến node / một cụm, đọc kết quả từng node + ảnh lúc lỗi,
// sửa kịch bản rồi chạy tiếp từ node lỗi với biến của lần trước (SeedRun) – không chạy lại từ đầu.

// Trace – một lần chạy kèm kết quả từng node (theo thứ tự chạy) và ảnh lúc node lỗi.
type Trace struct {
	Status
	Inputs  map[string]any         `json:"inputs,omitempty"` // giá trị input đã dùng – resume_run chạy lại cùng giá trị
	Results []model.StepResult     `json:"results"`
	Shots   map[string]image.Image `json:"-"` // theo StepResult.Shot
}

// part – kịch bản sau khi áp From / Only / Nodes / StopAt và kết quả lần chạy SeedRun.
func (m *Manager) part(sc model.Script, opt Options) (model.Script, *engine.Carry, error) {
	for _, id := range append(slices.Clone(opt.Nodes), opt.StopAt) {
		if id != "" && !slices.ContainsFunc(sc.Steps, func(s model.Step) bool { return s.ID == id }) {
			return sc, nil, fmt.Errorf("không có node %q (dùng id node)", id)
		}
	}
	switch {
	case len(opt.Nodes) > 0:
		sc = script.RunPart(sc, opt.From, opt.Nodes)
	case opt.From != "":
		sc = script.RunFrom(sc, opt.From, opt.Only)
	}
	if opt.MaxMinutes > 0 && (sc.MaxMinutes == 0 || sc.MaxMinutes > opt.MaxMinutes) {
		sc.MaxMinutes = opt.MaxMinutes
	}
	if opt.SeedRun == "" {
		return sc, nil, nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.runs[opt.SeedRun]
	switch {
	case r == nil:
		return sc, nil, fmt.Errorf("lần chạy %s: %w", opt.SeedRun, ErrNotFound)
	case r.carry == nil:
		return sc, nil, fmt.Errorf("lần chạy %s chưa xong – chờ xong rồi chạy tiếp", opt.SeedRun)
	}
	return sc, r.carry, nil
}

// Trace – trạng thái + kết quả từng node + ảnh lúc lỗi của lần chạy id ("" = đang chạy / gần nhất).
func (m *Manager) Trace(id string) (Trace, error) {
	st, err := m.Status(id)
	if err != nil {
		return Trace{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	r := m.runs[st.RunID]
	if r == nil { // vừa bị bỏ vì quá cũ
		return Trace{}, ErrNotFound
	}
	return Trace{Status: st, Inputs: r.values, Results: slices.Clone(r.results), Shots: r.shots}, nil
}

// noteNode ghi kết quả node vào vết chạy (OnNode).
func (m *Manager) noteNode(r *run, res model.StepResult) {
	m.mu.Lock()
	r.results = append(r.results, res)
	m.mu.Unlock()
}
