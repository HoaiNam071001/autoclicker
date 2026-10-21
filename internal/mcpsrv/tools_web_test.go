package mcpsrv

import (
	"testing"

	"ghosthand/internal/web"
)

func TestProbePage(t *testing.T) {
	cs, back := connect(t)
	back.page = web.PageCheck{Verdict: web.PageCaptcha, Signs: []string{"khung captcha đang hiện"}}
	out, isErr := call(t, cs, "probe_page", map[string]any{"url": "game-a.vn/play"})
	if isErr || out["verdict"] != "captcha" || out["useWeb"] != false {
		t.Fatalf("probe_page: %v", out)
	}
	if back.probed != "https://game-a.vn/play" {
		t.Errorf("url thiếu https:// – %q", back.probed)
	}
	for _, bad := range []string{" ", "file:///C:/Users/a/secret.txt", "chrome://settings"} {
		back.probed = ""
		if out, isErr := call(t, cs, "probe_page", map[string]any{"url": bad}); !isErr || back.probed != "" {
			t.Errorf("url %q phải bị từ chối, không mở: %v", bad, out)
		}
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
