package inputsource

import "testing"

var sources = []Source{
	{ID: "com.apple.keylayout.ABC", Name: "ABC", Lang: "en"},
	{ID: "com.apple.keylayout.US", Name: "U.S.", Lang: "en"},
	{ID: "com.apple.inputmethod.VietnameseIM.VietnameseTelex", Name: "Telex", Lang: "vi"},
}

func TestFind(t *testing.T) {
	cases := map[string]int{
		"abc": 0, "com.apple.keylayout.US": 1, "u.s.": 1, "telex": 2, "vietnamese": 2, "keylayout": 0, "": -1, "dvorak": -1,
	}
	for want, i := range cases {
		if got := Find(sources, want); got != i {
			t.Errorf("Find(%q) = %d, muốn %d", want, got, i)
		}
	}
}

func TestIsVietnamese(t *testing.T) {
	if IsVietnamese(sources[0]) || !IsVietnamese(sources[2]) || !IsVietnamese(Source{ID: "0000042A", Name: "vi-VN", Lang: "vi"}) {
		t.Error("nhận dạng bộ gõ tiếng Việt sai")
	}
}

func TestFindIME(t *testing.T) {
	if got := FindIME([]string{"Finder", "UniKey", "Safari"}); got != "UniKey" {
		t.Errorf("got %q", got)
	}
	if FindIME([]string{"Finder"}) != "" {
		t.Error("không có phần mềm gõ phải trả rỗng")
	}
}

func TestParseHotkey(t *testing.T) {
	for _, c := range []struct {
		in   string
		mods []string
		key  string
		bad  bool
	}{
		{"ctrl+shift", []string{"ctrl", "shift"}, "", false},
		{"Alt + Z", []string{"alt"}, "z", false},
		{"Control+Space", []string{"ctrl"}, "space", false},
		{"option+command+k", []string{"alt", "cmd"}, "k", false},
		{"ctrl+ctrl+shift", []string{"ctrl", "shift"}, "", false},
		{"", nil, "", true}, {"+", nil, "", true}, {"a+b", nil, "", true},
	} {
		mods, key, err := ParseHotkey(c.in)
		if (err != nil) != c.bad || key != c.key || len(mods) != len(c.mods) {
			t.Errorf("ParseHotkey(%q) = %v %q %v", c.in, mods, key, err)
		}
	}
}
