package vision

import (
	"testing"
)

// Dòng do OCR Windows (không có tiếng Việt) đọc thật trên trang sân thử, 2026-10-19.
func TestFindTextApprox(t *testing.T) {
	approx := func(text string, y int) TextLine { l := line(text, 0, y, 200, 20); l.Approx = true; return l }
	lines := []TextLine{
		approx("BÄM vÄo Töl", 0),         // BẤM VÀO TÔI – dấu sai, lệch 1 ký tự
		approx("Sän thir Ghosthand", 30), // Sân thử Ghosthand – “ử” thành “ir”
		approx("Click tråi", 60),         // Click trái
		approx("Töl", 90),
		line("BÄM vÄo Töl", 0, 120, 200, 20), // macOS (không Approx): so như cũ
	}
	cases := []struct {
		q    TextQuery
		want int
	}{
		{TextQuery{Text: "BẤM VÀO TÔI"}, 1},
		{TextQuery{Text: "SAN THU GHOSTHAND", NoAccent: true}, 1},
		{TextQuery{Text: "Sân thử"}, 1},
		{TextQuery{Text: "Click trái", Exact: true}, 1},
		{TextQuery{Text: "Click phải"}, 0}, // lệch 2 ký tự trên 9 – quá ngưỡng
		{TextQuery{Text: "Tôi"}, 0},        // < 5 ký tự phải khớp đúng (“tol” ≠ “toi”)
		{TextQuery{Text: "Töl"}, 3},        // khớp đúng: 2 dòng Approx (so không dấu) + dòng macOS
	}
	for _, c := range cases {
		if got := len(FindText(lines, c.q)); got != c.want {
			t.Errorf("%q: %d chỗ, cần %d", c.q.Text, got, c.want)
		}
	}
	// khung = đoạn khớp: “vÄo Töl” trong dòng không có khung từ → chia theo ký tự, bỏ “BÄM ”
	hits := FindText([]TextLine{approx("BÄM vÄo Töl", 0)}, TextQuery{Text: "vao toi"})
	if len(hits) != 1 || hits[0].Box.X < 60 {
		t.Errorf("khung đoạn khớp: %+v", hits)
	}
}

func TestFuzzyIndex(t *testing.T) {
	s := []rune("abcxdefabcdef")
	if a, b := fuzzyIndex(s, []rune("abcdef"), 0, 1); a != 0 || b != 7 {
		t.Errorf("chỗ đầu (thừa x): [%d, %d)", a, b)
	}
	if a, b := fuzzyIndex(s, []rune("abcdef"), 7, 0); a != 7 || b != 13 {
		t.Errorf("khớp đúng từ 7: [%d, %d)", a, b)
	}
	if a, _ := fuzzyIndex(s, []rune("zzzzzz"), 0, 1); a != -1 {
		t.Errorf("không có: %d", a)
	}
	if d := editDistance([]rune("kitten"), []rune("sitting")); d != 3 {
		t.Errorf("editDistance = %d", d)
	}
}

func TestSimilarity(t *testing.T) {
	for _, c := range []struct {
		got, want string
		min, max  int
	}{
		{"Bấm Gửi ngay", "gửi", 100, 100},
		{"Gui", "Gửi", 100, 100},
		{"Gưi tin", "Gửi tin", 100, 100}, // so không dấu
		{"Gởi tn", "Gửi tin", 60, 70},    // ơ ≠ ư + thiếu i
		{"Đăng nhập", "Gửi", 0, 40},
		{"", "Gửi", 0, 0},
	} {
		if s := Similarity(c.got, c.want); s < c.min || s > c.max {
			t.Errorf("Similarity(%q, %q) = %d, cần %d–%d", c.got, c.want, s, c.min, c.max)
		}
	}
}
