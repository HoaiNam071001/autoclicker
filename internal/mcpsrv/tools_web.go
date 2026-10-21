package mcpsrv

import (
	"context"
	"errors"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"ghosthand/internal/store"
)

// Thử trang + ưu tiên node (bước 44): probe_page mở trang trên Chrome của app, dò captcha / chặn bot / canvas → AI
// biết trang dùng được node web không; ưu tiên node (ảnh / chữ mặc định | web) AI đọc / đổi tự do, guide đổi theo.

// nodePrefTips – mục thêm vào guide theo ưu tiên node hiện tại.
var nodePrefTips = map[string]string{
	store.NodePrefWeb: `Ưu tiên node web (Cài đặt / set_node_preference): việc trên trang web → probe_page(url) trước.
- verdict ok: node open (webBrowser true – Chrome của app) + webWait / webClick / webType / webRead / webScroll, selector lấy từ outline (sel) – chính xác hơn tìm chữ / ảnh.
- captcha / blocked / canvas / empty: dùng node ảnh / chữ (open bằng trình duyệt thường, waitImage / clickImage findBy "text") và nói người dùng lý do (signs).
- Chrome của app chỉ để trang không cần đăng nhập – không đăng nhập tài khoản thật ở đó (dễ bị captcha / khoá); trang cần đăng nhập dùng trình duyệt thường + node ảnh / chữ.`,
	store.NodePrefVision: `Ưu tiên node ảnh / chữ (mặc định): trang web cũng làm bằng node open (trình duyệt thường) + tìm chữ / ảnh trên màn hình.
- Chỉ dùng node web khi người dùng xin hoặc sau khi đổi ưu tiên (set_node_preference "web"); muốn biết trang có dùng được node web không → probe_page(url).`,
}

type probeIn struct {
	URL string `json:"url" jsonschema:"trang cần thử (https://… hoặc tên miền)"`
}

type prefIn struct {
	Pref string `json:"pref" jsonschema:"vision (node ảnh / chữ – mặc định) | web (thử trang trên Chrome của app trước, không captcha thì dùng node web)"`
}

func (s *Server) addWebTools() {
	addTool(s, KindRun, &mcp.Tool{Name: "probe_page", Description: "Mở thử trang trên Chrome của app (profile riêng, không đăng nhập), chờ tải xong " +
		"rồi dò captcha / chặn bot / game vẽ canvas: verdict ok (dùng được node web) | captcha | blocked | canvas | empty, useWeb, signs " +
		"(dấu hiệu), elements, outline (≤ 40 phần tử kèm selector sel dùng được cho node web). Gọi trước khi soạn node web cho một trang."}, s.probePage)
	addTool(s, KindRead, &mcp.Tool{Name: "get_preferences", Description: "Cài đặt của app ảnh hưởng cách soạn: nodePref (vision | web), " +
		"mức chạy của AI (run)."}, s.getPrefs)
	addTool(s, KindSettings, &mcp.Tool{Name: "set_node_preference", Description: "Đổi ưu tiên node khi soạn: vision = node ảnh / chữ (mặc định), " +
		"web = thử trang (probe_page) trước, không captcha thì dùng node web. Đổi ngay, app báo người dùng; guide (get_step_types) theo mức mới."}, s.setPref)
}

func (s *Server) probePage(ctx context.Context, _ *mcp.CallToolRequest, in probeIn) (*mcp.CallToolResult, any, error) {
	u := strings.TrimSpace(in.URL)
	if u == "" {
		return fail(errors.New("thiếu url"))
	}
	if !strings.Contains(u, "://") {
		u = "https://" + u
	}
	// chỉ trang web – file:// sẽ cho AI đọc file trên máy qua outline (MCP không có quyền đọc file)
	if l := strings.ToLower(u); !strings.HasPrefix(l, "https://") && !strings.HasPrefix(l, "http://") {
		return fail(errors.New("probe_page chỉ mở trang http:// / https://"))
	}
	pc, err := s.back.CheckPage(ctx, u)
	if err != nil {
		return fail(err)
	}
	return ok(pc)
}

func (s *Server) getPrefs(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	m := s.back.MCP()
	return ok(map[string]any{"nodePref": m.NodePreference(), "run": m.RunMode()})
}

func (s *Server) setPref(_ context.Context, _ *mcp.CallToolRequest, in prefIn) (*mcp.CallToolResult, any, error) {
	p := strings.ToLower(strings.TrimSpace(in.Pref))
	if err := s.back.SetNodePref(p); err != nil {
		return fail(err)
	}
	return ok(map[string]any{"nodePref": p, "tips": nodePrefTips[p]})
}
