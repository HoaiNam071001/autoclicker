package vision

// So chữ gần đúng cho dòng do bộ đọc không biết tiếng Việt đọc ra (TextLine.Approx – OCR Windows): dấu đọc sai
// (“BẤM” → “BÄM”, bỏ dấu là khớp) và thỉnh thoảng lệch thêm một ký tự (“TÔI” → “Töl”) – chịu tối đa fuzzyBudget lỗi.

// fuzzyBudget – số ký tự được lệch (thêm / bớt / sai): 1 trên 5 ký tự, chuỗi < 5 ký tự phải khớp đúng.
func fuzzyBudget(n int) int { return n / 5 }

// fuzzyIndex tìm từ vị trí from đoạn của s khớp sub với ≤ k lỗi (Sellers: đầu đoạn tự do). Thấy chỗ đầu tiên đạt
// ngưỡng thì đi tiếp khi số lỗi còn giảm (lấy chỗ khớp nhất quanh đó). Trả đoạn [start, end) hoặc -1.
func fuzzyIndex(s, sub []rune, from, k int) (start, end int) {
	m := len(sub)
	cost, st := make([]int, m+1), make([]int, m+1)
	for i := range cost {
		cost[i], st[i] = i, from
	}
	nc, ns := make([]int, m+1), make([]int, m+1)
	best, bestStart, bestEnd := k+1, -1, -1
	for j := from; j < len(s); j++ {
		nc[0], ns[0] = 0, j+1
		for i := 1; i <= m; i++ {
			c, o := cost[i-1], st[i-1] // khớp / sai ký tự
			if s[j] != sub[i-1] {
				c++
			}
			if cost[i]+1 < c { // s thừa ký tự
				c, o = cost[i]+1, st[i]
			}
			if nc[i-1]+1 < c { // s thiếu ký tự
				c, o = nc[i-1]+1, ns[i-1]
			}
			nc[i], ns[i] = c, o
		}
		cost, nc, st, ns = nc, cost, ns, st
		switch {
		case cost[m] < best:
			best, bestStart, bestEnd = cost[m], st[m], j+1
		case bestEnd >= 0:
			return bestStart, bestEnd
		}
	}
	return bestStart, bestEnd
}

// editDistance – số lỗi giữa hai chuỗi (Levenshtein), cho kiểu so cả dòng.
func editDistance(a, b []rune) int {
	prev, cur := make([]int, len(b)+1), make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			c := prev[j-1]
			if a[i-1] != b[j-1] {
				c++
			}
			cur[j] = min(c, prev[j]+1, cur[j-1]+1)
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

// Similarity – got giống want bao nhiêu % (0–100) để xếp ứng viên gần đúng khi không tìm thấy (StepError.Detail):
// so không dấu, bỏ hoa / thường, khoảng trắng; lấy đoạn trong got khớp want ít lỗi nhất (got dài hơn không bị trừ).
func Similarity(got, want string) int {
	g, _ := textKey(got, true)
	w, _ := textKey(want, true)
	m := len(w)
	if m == 0 || len(g) == 0 {
		return 0
	}
	cost, next := make([]int, m+1), make([]int, m+1)
	for i := range cost {
		cost[i] = i
	}
	best := cost[m]
	for _, r := range g {
		next[0] = 0 // đầu đoạn tự do
		for i := 1; i <= m; i++ {
			c := cost[i-1]
			if r != w[i-1] {
				c++
			}
			next[i] = min(c, cost[i]+1, next[i-1]+1)
		}
		cost, next = next, cost
		best = min(best, cost[m])
	}
	return 100 * (m - best) / m
}
