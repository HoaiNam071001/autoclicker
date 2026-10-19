package device

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"autoclicker/internal/element"
	"autoclicker/internal/model"
)

// Đọc cây accessibility thật – chỉ chạy khi AUTOCLICKER_REAL=1 (cần quyền Accessibility). AUTOCLICKER_APP = tên app
// (trống = app trên cùng), AUTOCLICKER_FIND = chữ cần tìm.
func TestRealElements(t *testing.T) {
	if os.Getenv("AUTOCLICKER_REAL") != "1" {
		t.Skip("đặt AUTOCLICKER_REAL=1 để chạy với app thật")
	}
	if runtime.GOOS != "darwin" {
		t.Skip("cây phần tử Windows (UI Automation) làm ở bước 37 E5")
	}
	app := os.Getenv("AUTOCLICKER_APP")
	for round := 1; round <= 2; round++ {
		start := time.Now()
		nodes, err := SystemElements{}.Elements(app)
		if err != nil {
			t.Fatal(err)
		}
		roles := map[string]int{}
		for _, n := range nodes {
			roles[n.Role]++
		}
		t.Logf("lần %d: %d nút trong %v, vai trò: %v", round, len(nodes), time.Since(start), roles)
	}
	nodes, _ := SystemElements{}.Elements(app)
	for i, n := range nodes {
		if i < 40 || n.Role == "button" && i < 400 {
			t.Logf("%s%s(%s) %v title=%q desc=%q value=%q ph=%q", strings.Repeat(" ", min(n.Depth, 12)), n.Role, n.Raw, n.Box, n.Title, n.Desc, n.Value, n.Placeholder)
		}
	}
	if q := os.Getenv("AUTOCLICKER_FIND"); q != "" {
		s, err := model.ParseSelector(q)
		if err != nil {
			t.Fatal(err)
		}
		r := element.Find(nodes, element.Query{Selectors: []model.Selector{s}})
		t.Logf("tìm %q: found=%v count=%d %+v", q, r.Found, r.Count, r.Hit)
	}
}

func TestRealElementsErrors(t *testing.T) {
	if os.Getenv("AUTOCLICKER_REAL") != "1" {
		t.Skip("đặt AUTOCLICKER_REAL=1 để chạy với app thật")
	}
	if _, err := (SystemElements{}).Elements("khong-co-app-nay"); err == nil {
		t.Error("app không chạy phải báo lỗi")
	} else {
		t.Log(err)
	}
}
