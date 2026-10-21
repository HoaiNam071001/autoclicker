package mcpsrv

import (
	"strings"
	"testing"

	"ghosthand/internal/web"
)

func TestProbePage(t *testing.T) {
	cs, back := connect(t)
	back.page = web.PageCheck{Verdict: web.PageCaptcha, URL: "https://game-a.vn/play", Signs: []string{"khung captcha đang hiện"}}
	out, isErr := call(t, cs, "probe_page", map[string]any{"url": "game-a.vn/play"})
	if isErr || out["verdict"] != "captcha" || out["useWeb"] != false {
		t.Fatalf("probe_page: %v", out)
	}
	if back.probed != "https://game-a.vn/play" {
		t.Errorf("url thiếu https:// – %q", back.probed)
	}
	for _, bad := range []string{"file:///C:/Users/a/secret.txt", "chrome://settings"} {
		back.probed = "chưa gọi"
		if out, isErr := call(t, cs, "probe_page", map[string]any{"url": bad}); !isErr || back.probed != "chưa gọi" {
			t.Errorf("url %q phải bị từ chối, không mở: %v", bad, out)
		}
	}
	// không url = dò lại tab đang mở (trang tải lâu)
	back.probed = "chưa gọi"
	if out, isErr := call(t, cs, "probe_page", map[string]any{"url": " "}); isErr || back.probed != "" {
		t.Errorf("dò lại phải gọi CheckPage với url rỗng: %v (probed %q)", out, back.probed)
	}
	// dò lại trúng tab file:// người dùng tự mở → không trả outline cho AI
	back.page = web.PageCheck{Verdict: web.PageOK, URL: "file:///C:/Users/a/secret.html", Outline: []web.OutlineItem{{}}}
	if out, isErr := call(t, cs, "probe_page", nil); !isErr {
		t.Errorf("tab file:// phải bị từ chối: %v", out)
	}
}

func TestNodePreference(t *testing.T) {
	cs, _ := connect(t)
	if out, _ := call(t, cs, "get_preferences", nil); out["nodePref"] != "vision" {
		t.Fatalf("mặc định phải vision: %v", out)
	}
	if g, _ := call(t, cs, "get_step_types", nil); g["nodePref"] != "vision" || g["nodePrefTips"] == "" {
		t.Fatalf("guide mặc định: nodePref %v", g["nodePref"])
	}
	if out, isErr := call(t, cs, "set_node_preference", map[string]any{"pref": "Web"}); isErr || out["nodePref"] != "web" {
		t.Fatalf("set web: %v", out)
	}
	if g, _ := call(t, cs, "get_step_types", nil); g["nodePref"] != "web" || g["nodePrefTips"] != nodePrefTips["web"] {
		t.Errorf("guide sau khi đổi: %v", g["nodePref"])
	}
	if _, isErr := call(t, cs, "set_node_preference", map[string]any{"pref": "dom"}); !isErr {
		t.Error("mức lạ phải lỗi")
	}
}

func TestReadWebOutline(t *testing.T) {
	cs, back := connect(t)
	back.outline = web.Outline{URL: "https://erp.vn/employee", Title: "Nhân viên", Tab: "AB12", Dialogs: 1, Total: 1,
		Items: []web.OutlineItem{{Role: "textbox", Name: "Mã nhân viên", Selector: `role=textbox[name="Mã nhân viên"]`, InView: true}}}
	back.tabs = []web.Tab{{ID: "AB12", URL: "https://erp.vn/employee"}, {ID: "CD34", URL: "file:///C:/Users/a/secret.html"}}
	out, isErr := call(t, cs, "read_web_outline", map[string]any{"tab": " AB12 ", "dialog": true})
	if isErr || back.outlineTab != "AB12" || !back.outlineDialog {
		t.Fatalf("read_web_outline: %v (tab %q dialog %v)", out, back.outlineTab, back.outlineDialog)
	}
	if o, _ := out["outline"].(string); !strings.Contains(o, `role=textbox[name="Mã nhân viên"]`) || out["tab"] != "AB12" {
		t.Errorf("thiếu selector / id tab: %v", out)
	}
	if tabs, _ := out["tabs"].([]any); len(tabs) != 1 {
		t.Errorf("tab file:// không được lộ cho AI: %v", out["tabs"])
	}
	back.outline.URL = "file:///C:/Users/a/secret.html"
	if out, isErr := call(t, cs, "read_web_outline", nil); !isErr {
		t.Errorf("tab file:// phải bị từ chối: %v", out)
	}
}
