package script

import (
	"fmt"
	"regexp"

	"autoclicker/internal/model"
)

// NodeVarRe – biến đọc kết quả node: {{n2.output.x}}, {{n2.status}}, {{n2.error.message}}, {{vars.ten}}, {{last.output.text}}.
// Khác {{key}} của input (không có dấu chấm, inputs.go varRe).
var NodeVarRe = regexp.MustCompile(`\{\{\s*([A-Za-z][A-Za-z0-9_]*)\.([A-Za-z0-9_.]+)\s*\}\}`)

// ValidateFlow kiểm tra cấu trúc flow trước khi chạy: node bắt đầu, cạnh, mã node, nguồn chỗ thấy của click theo ảnh,
// biến trỏ node. Dữ liệu từng node (điểm, ảnh mẫu, phím…) do validateScript kiểm tra.
// Node không đi tới được từ start không phải lỗi (unreachableNodes – UI cảnh báo).
func validateFlow(sc model.Script) error {
	if len(sc.Steps) == 0 {
		return &model.ScriptError{Msg: "kịch bản chưa có node nào"}
	}
	idx := make(map[string]int, len(sc.Steps))
	refs := make(map[string]int, len(sc.Steps))
	for i, s := range sc.Steps {
		idx[s.ID] = i
		if !model.ValidRef(s.Ref) {
			return stepErr(i, s, "ref", "mã node %q không hợp lệ – chữ cái đầu, sau đó chữ / số / _ (không dùng vars, last)", s.Ref)
		}
		if _, dup := refs[s.Ref]; dup {
			return stepErr(i, s, "ref", "mã node %q trùng với node khác", s.Ref)
		}
		refs[s.Ref] = i
	}
	if _, ok := idx[sc.Start]; !ok {
		return &model.ScriptError{Field: "start", Msg: "chưa chọn node bắt đầu"}
	}

	ports := map[[2]string]bool{}
	for _, e := range sc.Edges {
		i, ok := idx[e.From]
		if !ok {
			return &model.ScriptError{Field: "edges", Msg: fmt.Sprintf("cạnh %s đi ra từ node không còn trong kịch bản", e.ID)}
		}
		s := sc.Steps[i]
		_, toOK := idx[e.To]
		switch {
		case !model.HasPort(s, e.Port):
			return stepErr(i, s, "edges", "node không có cổng %q", e.Port)
		case !toOK:
			return stepErr(i, s, "edges", "cổng %s nối tới node không còn trong kịch bản", e.Port)
		case e.From == e.To:
			return stepErr(i, s, "edges", "cổng %s nối vào chính node này – vòng lặp phải đi qua node khác (vd Đếm)", e.Port)
		case ports[[2]string{e.From, e.Port}]:
			return stepErr(i, s, "edges", "cổng %s có hơn một cạnh", e.Port)
		}
		ports[[2]string{e.From, e.Port}] = true
	}

	for i, s := range sc.Steps {
		if !s.Enabled {
			continue
		}
		if s.Type == model.StepClickImage && s.ImageSource != model.ImageFromSearch {
			j, ok := idx[s.ImageFrom]
			switch {
			case s.ImageFrom == "":
				return stepErr(i, s, "imageFrom", "chưa chọn node lấy chỗ thấy (Chờ ảnh / Chờ chữ)")
			case !ok:
				return stepErr(i, s, "imageFrom", "node lấy chỗ thấy không còn trong kịch bản – chọn lại")
			case !findsImage(sc.Steps[j]):
				return stepErr(i, s, "imageFrom", "node %s không tìm ảnh / chữ nên không có chỗ thấy", sc.Steps[j].Ref)
			}
		}
		for _, f := range VarFields(&s) {
			for _, m := range NodeVarRe.FindAllStringSubmatch(*f.Ptr, -1) {
				if name := m[1]; name != "vars" && name != "last" {
					if _, ok := refs[name]; !ok {
						return stepErr(i, s, f.Name, "biến %s trỏ tới node %q không có trong kịch bản", m[0], name)
					}
				}
			}
		}
	}
	return nil
}

// findsImage: node có "chỗ thấy" để click theo – chờ ảnh / chữ, click theo ảnh tự tìm, Nếu / Switch theo màn hình
// (chỗ thấy của điều kiện đúng đầu tiên).
func findsImage(s model.Step) bool {
	return s.Type == model.StepWaitImage || s.Type == model.StepClickImage && s.ImageSource == model.ImageFromSearch || s.Type == model.StepIf || s.Type == model.StepSwitchScreen
}

// UnreachableNodes – id các node không đi tới được từ start (sẽ không bao giờ chạy), theo thứ tự trong mảng.
func unreachableNodes(sc model.Script) []string {
	ok := model.ReachableFrom(sc)
	out := []string{}
	for _, s := range sc.Steps {
		if !ok[s.ID] {
			out = append(out, s.ID)
		}
	}
	return out
}
