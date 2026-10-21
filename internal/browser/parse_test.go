package browser

import (
	"errors"
	"testing"

	"ghosthand/internal/model"
)

// Trình duyệt đang ở trước lên đầu; tab đang hiện theo active của từng cửa sổ.
func TestParseTabs(t *testing.T) {
	out := `[{"browser":"chrome","front":false,"windows":[11,12],"active":[2,1],"ids":[[1,2],[3]],"titles":[["A","B"],["C"]],"urls":[["a","b"],["c"]],"loading":[[false,true],[false]]},
		{"browser":"edge","front":true,"windows":[5],"active":[1],"ids":[[9]],"titles":[["E"]],"urls":[["e"]]}]`
	tabs, err := parseTabs(out)
	if err != nil {
		t.Fatal(err)
	}
	want := []model.TabInfo{
		{Browser: "edge", ID: 9, Window: 5, Index: 1, Title: "E", URL: "e", Active: true},
		{Browser: "chrome", ID: 1, Window: 11, Index: 1, Title: "A", URL: "a"},
		{Browser: "chrome", ID: 2, Window: 11, Index: 2, Title: "B", URL: "b", Active: true, Loading: true},
		{Browser: "chrome", ID: 3, Window: 12, Index: 1, Title: "C", URL: "c", Active: true},
	}
	if len(tabs) != len(want) {
		t.Fatalf("có %d tab: %+v", len(tabs), tabs)
	}
	for i := range want {
		if tabs[i] != want[i] {
			t.Errorf("tab %d: %+v, muốn %+v", i, tabs[i], want[i])
		}
	}
	if tabs, err := parseTabs("[]"); err != nil || len(tabs) != 0 {
		t.Errorf("không trình duyệt nào: %v %v", tabs, err)
	}
}

func TestOsaError(t *testing.T) {
	for msg, want := range map[string]error{
		"execution error: Error: An error occurred. (-1743)": ErrNoAutomation,
		"execution error: Application isn't running. (-600)": ErrNotRunning,
		"execution error: Can't get object. (-1728)":         ErrTabGone,
		"execution error: lỗi khác (-2700)":                  nil,
	} {
		if got := osaError(msg); !errors.Is(got, want) {
			t.Errorf("%s → %v, muốn %v", msg, got, want)
		}
	}
}
