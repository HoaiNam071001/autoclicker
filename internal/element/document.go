package element

// DocumentIndex – chỉ số ô soạn thảo / tài liệu trong nodes khi người dùng không chỉ selector: ô nhập (textbox,
// searchbox) hoặc document có nội dung, chọn cái có diện tích lớn nhất (thân tài liệu lớn hơn ô tìm kiếm / ô địa chỉ).
// selection = true: chỉ xét phần tử đang có đoạn chữ bôi đen (Selected), cái đầu tiên. Không có → -1.
func DocumentIndex(nodes []Node, selection bool) int {
	best, area := -1, 0
	for i, n := range nodes {
		if n.Role != "textbox" && n.Role != "searchbox" && n.Role != "document" {
			continue
		}
		if selection {
			if n.Selected != "" {
				return i
			}
			continue
		}
		if n.Value == "" {
			continue
		}
		if a := n.Box.W * n.Box.H; a > area || best < 0 {
			best, area = i, a
		}
	}
	return best
}
