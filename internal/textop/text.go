// Package textop – xử lý chữ của các node Lọc mẫu, Tách chữ, Làm sạch chữ, Số & ngày (bước 42 G4b). Hàm thuần, không
// thiết bị; engine/textnodes.go gọi lúc chạy, App.PreviewText gọi cho nút Chạy thử của inspector.
package textop

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"

	"ghosthand/internal/model"
)

// Matches – kết quả Lọc mẫu: Text chỗ khớp đầu (nhóm đầu nếu mẫu có nhóm), Items mọi chỗ (all) hoặc chỉ chỗ đầu,
// Groups mọi nhóm của chỗ đầu.
type Matches struct {
	Text   string
	Items  []string
	Groups []string
}

// Match lọc text theo regex pattern; false = không khớp chỗ nào. Mẫu sai cú pháp → lỗi.
func Match(text, pattern string, all bool) (Matches, bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Matches{}, false, err
	}
	n := 1
	if all {
		n = -1
	}
	found := re.FindAllStringSubmatch(text, n)
	if len(found) == 0 {
		return Matches{}, false, nil
	}
	pick := func(m []string) string {
		if len(m) > 1 {
			return m[1]
		}
		return m[0]
	}
	out := Matches{Text: pick(found[0]), Groups: append([]string{}, found[0][1:]...)}
	for _, m := range found {
		out.Items = append(out.Items, pick(m))
	}
	return out, true, nil
}

// Split tách text theo sep ("" = mỗi dòng; \n \t viết tay = xuống dòng / tab), cắt khoảng trắng từng phần, bỏ phần trống.
func Split(text, sep string) []string {
	sep = strings.NewReplacer(`\n`, "\n", `\t`, "\t").Replace(sep)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if sep == "" {
		sep = "\n"
	}
	items := []string{}
	for _, p := range strings.Split(text, sep) {
		if p = strings.TrimSpace(p); p != "" {
			items = append(items, p)
		}
	}
	return items
}

// Clean làm sạch text theo node Làm sạch chữ, lần lượt: thay → bỏ dấu → hoa / thường → gọn khoảng trắng.
func Clean(text string, s model.Step) (string, error) {
	if s.ReplaceFrom != "" {
		if s.ReplaceRegex {
			re, err := regexp.Compile(s.ReplaceFrom)
			if err != nil {
				return "", err
			}
			text = re.ReplaceAllString(text, s.ReplaceTo)
		} else {
			text = strings.ReplaceAll(text, s.ReplaceFrom, s.ReplaceTo)
		}
	}
	if s.StripAccents {
		text = StripAccents(text)
	}
	switch s.CaseMode {
	case model.CaseLower:
		text = strings.ToLower(text)
	case model.CaseUpper:
		text = strings.ToUpper(text)
	case model.CaseTitle:
		text = titleCase(text)
	}
	if s.Trim {
		text = tidy(text)
	}
	return text, nil
}

// StripAccents bỏ dấu tiếng Việt (giữ hoa / thường): “Đà Nẵng” → “Da Nang”.
func StripAccents(s string) string {
	var b strings.Builder
	for _, r := range norm.NFD.String(s) {
		switch {
		case unicode.Is(unicode.Mn, r):
		case r == 'đ':
			b.WriteRune('d')
		case r == 'Đ':
			b.WriteRune('D')
		default:
			b.WriteRune(r)
		}
	}
	return norm.NFC.String(b.String())
}

// titleCase – chữ đầu mỗi từ hoa, còn lại thường.
func titleCase(s string) string {
	rs := []rune(strings.ToLower(s))
	for i, r := range rs {
		if i == 0 || unicode.IsSpace(rs[i-1]) {
			rs[i] = unicode.ToUpper(r)
		}
	}
	return string(rs)
}

// tidy – cắt khoảng trắng hai đầu từng dòng, gộp khoảng trắng liền nhau thành một dấu cách, bỏ dòng trống.
func tidy(s string) string {
	var lines []string
	for _, l := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		if l = strings.Join(strings.Fields(l), " "); l != "" {
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n")
}
