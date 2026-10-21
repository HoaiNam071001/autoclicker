package mcpsrv

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"ghosthand/internal/model"
	"ghosthand/internal/vision"
)

// Nhìn màn hình khi AI xin (bước 31b) – chữ trước ảnh: describe_screen (chữ OCR kèm khung + cửa sổ) rẻ token, toạ độ
// chính xác, dùng được với AI không xem ảnh; screenshot khi cần nhìn bố cục. Toạ độ trả về theo hệ của chuột.

const (
	maxLines   = 150
	maxWindows = 20
)

type areaIn struct {
	Window string `json:"window,omitempty" jsonschema:"chụp cửa sổ có tiêu đề / tên app chứa chữ này (cửa sổ trên cùng khớp)"`
	X      int    `json:"x,omitempty" jsonschema:"vùng (toạ độ chuột); w = 0 = cả màn hình chính"`
	Y      int    `json:"y,omitempty"`
	W      int    `json:"w,omitempty"`
	H      int    `json:"h,omitempty"`
}

type shotIn struct {
	areaIn
	WebTab string `json:"webTab,omitempty" jsonschema:"chụp riêng một tab của Chrome của app qua trình duyệt (id tab từ probe_page / read_web_outline, hoặc chuỗi trong URL; \"active\" = tab đang hiện) – đúng tab kể cả khi Chrome bị che; bỏ qua window / vùng"`
}

type textLine struct {
	Text string     `json:"text"`
	Box  model.Rect `json:"box"`
}

func (s *Server) addScreenTools() {
	addTool(s, KindScreen, &mcp.Tool{Name: "describe_screen", Description: "Đọc màn hình bằng chữ (rẻ hơn ảnh): các dòng chữ (OCR) kèm khung " +
		"theo toạ độ chuột + danh sách cửa sổ đang mở (trên cùng trước). Dùng để tìm chữ / nút cho node waitImage / clickImage (findBy text) " +
		"hoặc biết cửa sổ nào đang trước."}, s.describeScreen)
	addTool(s, KindScreen, &mcp.Tool{Name: "screenshot", Description: "Chụp màn hình / vùng / cửa sổ, trả ảnh (thu nhỏ) + vùng + tỷ lệ: " +
		"toạ độ chuột = area.x + px ảnh × scale (y tương tự). webTab: chụp một tab Chrome của app (khung nhìn trang, toạ độ trong trang = px ảnh × scale)."}, s.screenshot)
}

// area – vùng cần xem: cửa sổ khớp tên > vùng gửi lên > màn chính.
func (s *Server) area(in areaIn) (model.Rect, error) {
	if q := strings.ToLower(strings.TrimSpace(in.Window)); q != "" {
		ws, err := s.back.Windows()
		if err != nil {
			return model.Rect{}, err
		}
		for _, w := range ws {
			if !w.Minimized && (strings.Contains(strings.ToLower(w.Title), q) || strings.Contains(strings.ToLower(w.App), q)) {
				return w.Box, nil
			}
		}
		return model.Rect{}, fmt.Errorf("không thấy cửa sổ đang hiện nào khớp %q – describe_screen để xem danh sách", in.Window)
	}
	if in.W > 0 && in.H > 0 {
		return model.Rect{X: in.X, Y: in.Y, W: in.W, H: in.H}, nil
	}
	for _, d := range s.back.Displays() {
		if d.Main {
			return d.Rect, nil
		}
	}
	return model.Rect{}, fmt.Errorf("không đọc được màn hình")
}

func (s *Server) screenshot(ctx context.Context, _ *mcp.CallToolRequest, in shotIn) (*mcp.CallToolResult, any, error) {
	if t := strings.TrimSpace(in.WebTab); t != "" {
		return s.tabShot(ctx, t)
	}
	r, err := s.area(in.areaIn)
	if err != nil {
		return fail(err)
	}
	img, err := s.back.Capture(r)
	if err != nil {
		return fail(err)
	}
	c, w := pngContent(img)
	if c == nil {
		return fail(fmt.Errorf("không mã hoá được ảnh"))
	}
	res, _, err := ok(map[string]any{"area": r, "imageWidth": w, "scale": float64(r.W) / float64(w)})
	if err != nil {
		return nil, nil, err
	}
	res.Content = append(res.Content, c)
	return res, nil, nil
}

// tabShot – ảnh một tab Chrome của app (webTab "active" = tab đang hiện); chỉ trang http(s) như read_web_outline.
func (s *Server) tabShot(ctx context.Context, tab string) (*mcp.CallToolResult, any, error) {
	if tab == "active" {
		tab = ""
	}
	sh, err := s.back.TabShot(ctx, tab)
	if err != nil {
		return fail(err)
	}
	if !isWebURL(sh.URL) {
		return fail(errors.New("tab không phải trang http:// / https://"))
	}
	c, w := pngContent(sh.Image)
	if c == nil {
		return fail(fmt.Errorf("không mã hoá được ảnh"))
	}
	res, _, err := ok(map[string]any{"tab": sh.Tab, "url": sh.URL, "imageWidth": w, "scale": float64(sh.CSSWidth) / float64(w),
		"note": "toạ độ trong trang (CSS px, như webBackground / outline) = px ảnh × scale – không phải toạ độ chuột"})
	if err != nil {
		return nil, nil, err
	}
	res.Content = append(res.Content, c)
	return res, nil, nil
}

func (s *Server) describeScreen(_ context.Context, _ *mcp.CallToolRequest, in areaIn) (*mcp.CallToolResult, any, error) {
	r, err := s.area(in)
	if err != nil {
		return fail(err)
	}
	img, err := s.back.Capture(r)
	if err != nil {
		return fail(err)
	}
	read, err := s.back.ReadText(img)
	if err != nil {
		return fail(err)
	}
	k := vision.ImageScale(img, r)
	lines := make([]textLine, 0, min(len(read), maxLines))
	for _, l := range read[:min(len(read), maxLines)] {
		b := l.Box
		lines = append(lines, textLine{Text: l.Text, Box: model.Rect{X: r.X + int(float64(b.X)/k), Y: r.Y + int(float64(b.Y)/k),
			W: int(float64(b.W) / k), H: int(float64(b.H) / k)}})
	}
	out := map[string]any{"area": r, "lines": lines}
	if len(read) > maxLines {
		out["omittedLines"] = len(read) - maxLines
	}
	if ws, err := s.back.Windows(); err == nil {
		out["windows"] = ws[:min(len(ws), maxWindows)]
	}
	return ok(out)
}
