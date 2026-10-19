package model

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
)

// Cùng dữ liệu với frontend/src/flow/app.test.ts (testdata/apps.json).
func TestMatchAppsShared(t *testing.T) {
	var data struct {
		Apps  []AppInfo
		Cases []struct {
			App, Match string
			Want       any
		}
	}
	b, err := os.ReadFile("../../testdata/apps.json")
	if err == nil {
		err = json.Unmarshal(b, &data)
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range data.Cases {
		got, err := MatchApps(data.Apps, c.App, c.Match)
		if c.Want == "error" {
			if err == nil {
				t.Errorf("%q %s: muốn lỗi", c.App, c.Match)
			}
			continue
		}
		var want []int
		for _, x := range c.Want.([]any) {
			want = append(want, int(x.(float64)))
		}
		if err != nil || !slices.Equal(got, want) {
			t.Errorf("%q %s: %v %v, muốn %v", c.App, c.Match, got, err, want)
		}
	}
}

func TestMainApp(t *testing.T) {
	m, n := MainApp([]AppInfo{{PID: 1}, {PID: 2, Windows: 2}, {PID: 3, Windows: 2}})
	if m.PID != 2 || n != 4 {
		t.Errorf("%d %d", m.PID, n)
	}
	if m, n := MainApp(nil); m.PID != 0 || n != 0 {
		t.Errorf("rỗng: %v %d", m, n)
	}
}

func TestProtectedApp(t *testing.T) {
	for _, n := range []string{"explorer.exe", `C:\Windows\explorer.exe`, "Finder", "svchost", "AutoClicker.exe", "Dock.app"} {
		if !ProtectedApp(n) {
			t.Errorf("%q phải bị chặn", n)
		}
	}
	for _, n := range []string{"notepad", "chrome.exe", "explorer2", "", "Zalo"} {
		if ProtectedApp(n) {
			t.Errorf("%q không bị chặn", n)
		}
	}
}
