package element

import (
	"fmt"

	"autoclicker/internal/model"
	"autoclicker/internal/vision"
)

// Query – cách tìm phần tử. Selector đầu tiên có kết quả thì dùng (như node web 25d).
type Query struct {
	Selectors  []model.Selector
	NoAccent   bool        // bỏ qua dấu tiếng Việt khi so chữ
	Occurrence int         // chỗ thứ mấy theo thứ tự đọc; ≤ 1 = chỗ đầu
	Area       *model.Rect // chỉ lấy phần tử có tâm trong vùng (nil = cả màn)
}

// Hit – phần tử tìm thấy. Index = vị trí trong danh sách nút đã đọc.
type Hit struct {
	Node
	Index int
}

// Result – kết quả Find: Hit có giá trị khi Found; Count = số phần tử khớp selector đã dùng (để biết chỗ thứ N có thật
// không); Selector = selector đã dùng (-1 = không selector nào khớp).
type Result struct {
	Hit      Hit
	Found    bool
	Count    int
	Selector int
}

// CheckSelector – selector dùng được cho phần tử màn hình: css / xpath chỉ có trong trang web của Chrome của app (25d).
func CheckSelector(s model.Selector) error {
	switch s.By {
	case model.SelText, model.SelAria, model.SelRole:
		return nil
	}
	return fmt.Errorf("phần tử màn hình chỉ tìm theo text= / aria= / role= – %s= chỉ dùng được ở node web", s.By)
}

// Find tìm theo q trong nodes (thứ tự đọc). Chỗ thứ N vượt số phần tử khớp → Found = false, Count vẫn báo số đã thấy.
func Find(nodes []Node, q Query) Result {
	for i, s := range q.Selectors {
		hits := match(nodes, s, q)
		if len(hits) == 0 {
			continue
		}
		res := Result{Count: len(hits), Selector: i}
		if k := max(q.Occurrence, 1); k <= len(hits) {
			res.Hit, res.Found = hits[k-1], true
		}
		return res
	}
	return Result{Selector: -1}
}

// match – mọi phần tử khớp một selector, theo thứ tự đọc, đã lọc vùng. text=: bỏ phần tử có con cháu cũng khớp
// (nút chứa khối chữ “Gửi” → lấy khối chữ, khung sát chữ hơn).
func match(nodes []Node, s model.Selector, q Query) []Hit {
	eq := func(got, want string) bool { return vision.MatchText(got, want, s.Exact, q.NoAccent) }
	ok := func(n Node) bool {
		switch s.By {
		case model.SelText:
			return eq(n.ownText(), s.Value)
		case model.SelAria:
			return eq(n.Desc, s.Value)
		case model.SelRole:
			return n.Role == s.Value && (s.Name == "" || vision.MatchText(n.Name(), s.Name, false, q.NoAccent))
		}
		return false
	}
	var hits []Hit
	for i, n := range nodes {
		if !ok(n) || q.Area != nil && !inside(n.Box, *q.Area) {
			continue
		}
		if s.By == model.SelText && childMatches(nodes, i, ok) {
			continue
		}
		hits = append(hits, Hit{Node: n, Index: i})
	}
	return hits
}

// childMatches – có con cháu của nodes[i] (các nút liền sau sâu hơn) cũng khớp.
func childMatches(nodes []Node, i int, ok func(Node) bool) bool {
	for j := i + 1; j < len(nodes) && nodes[j].Depth > nodes[i].Depth; j++ {
		if ok(nodes[j]) {
			return true
		}
	}
	return false
}

// inside – tâm của b nằm trong a.
func inside(b, a model.Rect) bool {
	cx, cy := b.X+b.W/2, b.Y+b.H/2
	return cx >= a.X && cy >= a.Y && cx < a.X+a.W && cy < a.Y+a.H
}

// Settled – khung không đổi quá 1 điểm giữa hai lần đo (trang còn cuộn / chạy hiệu ứng thì chưa bấm).
func Settled(a, b model.Rect) bool {
	d := func(x, y int) bool { return x-y <= 1 && y-x <= 1 }
	return d(a.X, b.X) && d(a.Y, b.Y) && d(a.W, b.W) && d(a.H, b.H)
}
