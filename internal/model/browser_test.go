package model

import (
	"slices"
	"testing"
)

func TestMatchTabsAndActive(t *testing.T) {
	tabs := []TabInfo{
		{ID: 1, Window: 1, Title: "Hộp thư – Gmail", URL: "https://mail.google.com/"},
		{ID: 2, Window: 1, Title: "Bảng giá", URL: "http://x/price-list.html", Active: true},
		{ID: 3, Window: 2, Title: "YouTube", URL: "https://youtube.com/", Active: true},
	}
	for want, wantIdx := range map[string][]int{"": {0, 1, 2}, "bang gia": {1}, "GOOGLE.com": {0}, "price": {1}, "zalo": nil} {
		if got, err := MatchTabs(tabs, want, MatchContains); err != nil || !slices.Equal(got, wantIdx) {
			t.Errorf("%q: %v %v, muốn %v", want, got, err, wantIdx)
		}
	}
	if _, err := MatchTabs(tabs, "(", MatchRegex); err == nil {
		t.Error("regex sai phải báo lỗi")
	}
	if i := ActiveTab(tabs); i != 1 {
		t.Errorf("tab đang hiện của cửa sổ trên cùng: %d", i)
	}
	if i := ActiveTab(nil); i != -1 {
		t.Errorf("không tab: %d", i)
	}
}
