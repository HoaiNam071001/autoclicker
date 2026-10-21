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
- loading (trang tải lâu): chờ vài giây rồi probe_page không url (dò lại tab, không tải lại) – tối đa 2–3 lần, vẫn loading thì screenshot xem trang rồi quyết.
- captcha / blocked / canvas / empty: dùng node ảnh / chữ (open bằng trình duyệt thường, waitImage / clickImage findBy "text") và nói người dùng lý do (signs).
- Chrome của app chỉ để trang không cần đăng nhập – không đăng nhập tài khoản thật ở đó (dễ bị captcha / khoá); trang cần đăng nhập dùng trình duyệt thường + node ảnh / chữ.`,
	store.NodePrefVision: `Ưu tiên node ảnh / chữ (mặc định): trang web cũng làm bằng node open (trình duyệt thường) + tìm chữ / ảnh trên màn hình.
- Chỉ dùng node web khi người dùng xin hoặc sau khi đổi ưu tiên (set_node_preference "web"); muốn biết trang có dùng được node web không → probe_page(url).`,
}

type probeIn struct {
	URL string `json:"url,omitempty" jsonschema:"trang cần thử (https://… hoặc tên miền); bỏ trống = dò lại tab đang mở, không tải lại"`
}

type prefIn struct {
	Pref string `json:"pref" jsonschema:"vision (node ảnh / chữ – mặc định) | web (thử trang trên Chrome của app trước, không captcha thì dùng node web)"`
}

func (s *Server) addWebTools() {
	addTool(s, KindRun, &mcp.Tool{Name: "probe_page", Description: "Mở thử trang trên Chrome của app (profile riêng, không đăng nhập), chờ trang yên " +
		"(≤ 30 s) rồi dò captcha / chặn bot / game vẽ canvas: verdict ok (dùng được node web) | captcha | blocked | canvas | empty | " +
		"loading (trang tải lâu, chưa hiện gì), useWeb, signs (dấu hiệu), elements, outline (≤ 40 phần tử kèm selector sel dùng được cho node web). " +
		"Gọi trước khi soạn node web cho một trang. Không truyền url = dò lại tab đang mở, không tải lại: verdict loading hoặc loaded false " +
		"→ chờ vài giây rồi gọi lại không url (hoặc screenshot để xem trang)."}, s.probePage)
	addTool(s, KindRead, &mcp.Tool{Name: "get_preferences", Description: "Cài đặt của app ảnh hưởng cách soạn: nodePref (vision | web), " +
		"mức chạy của AI (run)."}, s.getPrefs)
	addTool(s, KindSettings, &mcp.Tool{Name: "set_node_preference", Description: "Đổi ưu tiên node khi soạn: vision = node ảnh / chữ (mặc định), " +
		"web = thử trang (probe_page) trước, không captcha thì dùng node web. Đổi ngay, app báo người dùng; guide (get_step_types) theo mức mới."}, s.setPref)
}

func (s *Server) probePage(ctx context.Context, _ *mcp.CallToolRequest, in probeIn) (*mcp.CallToolResult, any, error) {
	u := strings.TrimSpace(in.URL)
	if u != "" && !strings.Contains(u, "://") {
		u = "https://" + u
	}
	// chỉ trang web – file:// sẽ cho AI đọc file trên máy qua outline (MCP không có quyền đọc file).
	// Dò lại (url rỗng) dùng tab Chrome của app đang mở – tab đó do probe_page / người dùng mở.
	if u != "" && !isWebURL(u) {
		return fail(errors.New("probe_page chỉ mở trang http:// / https://"))
	}
	pc, err := s.back.CheckPage(ctx, u)
	if err != nil {
		return fail(err)
	}
	if !isWebURL(pc.URL) { // dò lại trúng tab file:// / chrome:// người dùng tự mở
		return fail(errors.New("tab đang mở không phải trang http:// / https:// – gọi probe_page kèm url"))
	}
	return ok(pc)
}

func isWebURL(u string) bool {
	l := strings.ToLower(u)
	return strings.HasPrefix(l, "https://") || strings.HasPrefix(l, "http://")
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
