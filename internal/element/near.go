package element

import (
	"sort"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// Near – phần tử gần giống cách tìm (StepError.Detail khi không thấy, bước 31b): AI / người xem để sửa selector.
type Near struct {
	Role  string     `json:"role"`
	Name  string     `json:"name"`
	Box   model.Rect `json:"box"`
	Score int        `json:"score"` // % giống (vision.Similarity), cao nhất trong các selector
}

// nearMin – dưới mức này coi như không liên quan.
const nearMin = 50

// NearMatches – tối đa n phần tử giống q nhất (đã lọc vùng), giống nhiều trước; phần tử khớp hẳn cũng có (100 %) –
// khi thiếu chỗ thứ N, chúng là các chỗ đã thấy.
func NearMatches(nodes []Node, q Query, n int) []Near {
	var out []Near
	for _, nd := range nodes {
		if nd.Role == "window" || q.Area != nil && !inside(nd.Box, *q.Area) {
			continue
		}
		best := 0
		for _, s := range q.Selectors {
			best = max(best, nearScore(nd, s))
		}
		if best >= nearMin {
			out = append(out, Near{Role: nd.Role, Name: nd.Text(), Box: nd.Box, Score: best})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out[:min(len(out), n)]
}

// nearScore – một phần tử giống một selector bao nhiêu %. role=: sai vai trò thì chỉ tính nửa phần tên.
func nearScore(nd Node, s model.Selector) int {
	switch s.By {
	case model.SelText:
		return max(vision.Similarity(nd.ownText(), s.Value), vision.Similarity(nd.Name(), s.Value))
	case model.SelAria:
		return max(vision.Similarity(nd.Desc, s.Value), vision.Similarity(nd.Name(), s.Value))
	case model.SelRole:
		name := 100
		if s.Name != "" {
			name = vision.Similarity(nd.Name(), s.Name)
		}
		if nd.Role != s.Value {
			if s.Name == "" {
				return 0
			}
			return name / 2
		}
		return name
	}
	return 0
}
