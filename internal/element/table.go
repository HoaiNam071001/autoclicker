package element

import "strings"

// Đọc bảng (bước 43 B5): bảng trong cây accessibility = nút vai trò table / grid, dòng = row, ô = cell / gridcell /
// columnheader / rowheader (bảng web dựng bằng div có role=grid cũng vào đây). Logic thuần trên danh sách phẳng.

func isTableRole(r string) bool { return r == "table" || r == "grid" || r == "treegrid" }

func isCellRole(r string) bool {
	return r == "cell" || r == "gridcell" || r == "columnheader" || r == "rowheader"
}

// FindTable – chỉ số bảng đầu tiên (nth ≥ 1: bảng thứ nth) trong nodes[from:to]; chính nodes[from] cũng được tính.
// Bảng ngoài khung nhìn bỏ qua trừ khi offscreen. -1 = không có.
func FindTable(nodes []Node, from, to, nth int, offscreen bool) int {
	n := 0
	for i := max(from, 0); i < min(to, len(nodes)); i++ {
		if isTableRole(nodes[i].Role) && (offscreen || nodes[i].Shown()) {
			if n++; n >= max(nth, 1) {
				return i
			}
		}
	}
	return -1
}

// TableRows – nội dung bảng gốc nodes[t]: mỗi dòng một mảng ô (chữ trong ô nối bằng khoảng trắng). Dòng ngoài khung
// nhìn bỏ trừ khi offscreen; bảng lồng nhau: dòng của bảng con không tính vào bảng cha. Dòng không có ô role cell thì
// coi mỗi con trực tiếp là một ô (UIA hay báo ô là text / group).
func TableRows(nodes []Node, t int, offscreen bool) [][]string {
	end := SubtreeEnd(nodes, t)
	var rows [][]string
	for i := t + 1; i < end; {
		n := nodes[i]
		if isTableRole(n.Role) {
			i = SubtreeEnd(nodes, i) // bảng con
			continue
		}
		if n.Role != "row" {
			i++
			continue
		}
		rEnd := SubtreeEnd(nodes, i)
		if offscreen || n.Shown() {
			if cells := rowCells(nodes, i, rEnd, offscreen); len(cells) > 0 {
				rows = append(rows, cells)
			}
		}
		i = rEnd
	}
	return rows
}

func rowCells(nodes []Node, r, end int, offscreen bool) []string {
	var cells []string
	for _, wantRole := range []bool{true, false} {
		for i := r + 1; i < end; {
			n := nodes[i]
			if wantRole && !isCellRole(n.Role) || !wantRole && n.Depth != nodes[r].Depth+1 {
				i++
				continue
			}
			cEnd := SubtreeEnd(nodes, i)
			cells = append(cells, strings.Join(TextLines(nodes, i, cEnd, offscreen), " "))
			i = cEnd
		}
		if len(cells) > 0 {
			break
		}
	}
	return cells
}
