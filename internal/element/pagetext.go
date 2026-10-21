package element

import "strings"

// SubtreeEnd – chỉ số sau phần tử cuối của cây con gốc nodes[i] (danh sách phẳng theo thứ tự duyệt: con cháu là các nút
// liền sau có Depth sâu hơn).
func SubtreeEnd(nodes []Node, i int) int {
	j := i + 1
	for j < len(nodes) && nodes[j].Depth > nodes[i].Depth {
		j++
	}
	return j
}

// rawText – chữ phần tử tự mang: khối chữ / ô nhập → giá trị, còn lại → chữ trên phần tử (nhãn mô tả không tính –
// nút chỉ có icon không phải nội dung trang).
func rawText(n Node) string {
	if n.Role == RoleText || editable(n.Role) {
		if v := strings.TrimSpace(n.Value); v != "" {
			return v
		}
		if n.Role == RoleText {
			return strings.TrimSpace(n.Title)
		}
		return ""
	}
	return strings.TrimSpace(n.Title)
}

// TextLines – chữ trong nodes[from:to] theo thứ tự đọc, mỗi dòng một mục. Phần tử chứa con có chữ thì chỉ lấy chữ của
// con (nút bọc khối chữ “Gửi” không đọc hai lần); cửa sổ bỏ (chỉ là tiêu đề); phần tử ngoài khung nhìn (Shown sai) bỏ
// trừ khi offscreen; dòng liền kề trùng nhau gộp một.
func TextLines(nodes []Node, from, to int, offscreen bool) []string {
	var out []string
	for i := max(from, 0); i < min(to, len(nodes)); i++ {
		n := nodes[i]
		if n.Role == "window" || !offscreen && !n.Shown() {
			continue
		}
		t := rawText(n)
		if t == "" || childHasText(nodes, i, min(to, len(nodes)), offscreen) {
			continue
		}
		for _, l := range strings.Split(t, "\n") {
			l = strings.Join(strings.Fields(l), " ")
			if l != "" && (len(out) == 0 || out[len(out)-1] != l) {
				out = append(out, l)
			}
		}
	}
	return out
}

// childHasText – cây con của nodes[i] (trong [.., to)) có phần tử nào mang chữ (và được đọc).
func childHasText(nodes []Node, i, to int, offscreen bool) bool {
	if nodes[i].Role == RoleText || editable(nodes[i].Role) {
		return false // khối chữ / ô nhập là chữ gốc, con của nó (nếu có) không thay được
	}
	for j := i + 1; j < to && nodes[j].Depth > nodes[i].Depth; j++ {
		if rawText(nodes[j]) != "" && (offscreen || nodes[j].Shown()) {
			return true
		}
	}
	return false
}
