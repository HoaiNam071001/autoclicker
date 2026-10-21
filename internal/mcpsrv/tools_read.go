package mcpsrv

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"ghosthand/internal/model"
)

type pathIn struct {
	Path string `json:"path" jsonschema:"đường dẫn file kịch bản (từ list_scripts) hoặc tên kịch bản"`
}

type typesIn struct {
	Type string `json:"type,omitempty" jsonschema:"chỉ lấy một loại node (vd waitImage); trống = mọi loại kèm hướng dẫn flow"`
}

// ScriptEntry – một kịch bản trong list_scripts.
type ScriptEntry struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Steps     int    `json:"steps"`
	CreatedBy string `json:"createdBy,omitempty"` // "ai" = AI tạo
	Open      bool   `json:"open,omitempty"`      // đang mở trên app
	Unsaved   bool   `json:"unsaved,omitempty"`   // đang mở và có thay đổi chưa lưu – update_script sẽ ghi ra bản sao
	Error     string `json:"error,omitempty"`     // không đọc được file
}

func (s *Server) addReadTools() {
	addTool(s, KindRead, &mcp.Tool{Name: "list_scripts", Description: "Danh sách kịch bản: mở gần đây trên app + thư mục kịch bản của AI. " +
		"Mỗi mục: path, name, số node, createdBy (ai), open / unsaved (đang mở trên app)."}, s.listScripts)
	addTool(s, KindRead, &mcp.Tool{Name: "read_script", Description: "Đọc một kịch bản: {path, revision, script}. " +
		"Giữ revision để gửi kèm update_script."}, s.readScript)
	addTool(s, KindRead, &mcp.Tool{Name: "get_step_types", Description: "Hướng dẫn viết kịch bản: cấu trúc flow, biến {{…}}, trường chung, " +
		"từng loại node app hỗ trợ (mô tả, trường, ví dụ JSON, output, cổng ra) và mẹo viết kịch bản bền. Đọc trước khi tạo / sửa."}, s.stepTypes)
	addTool(s, KindRead, &mcp.Tool{Name: "get_screen_info", Description: "Các màn hình: khung theo toạ độ chuột (x, y, w, h), tỷ lệ DPI / Retina, màn chính. " +
		"Toạ độ trong kịch bản dùng cùng hệ này."}, s.screenInfo)
}

// paths – mở gần đây (mới nhất trước) rồi file trong thư mục kịch bản của AI, không trùng.
func (s *Server) paths() []string {
	out := slices.Clone(s.back.Recent())
	if dir, err := s.back.ScriptsDir(); err == nil {
		files, _ := filepath.Glob(filepath.Join(dir, "*.json"))
		for _, f := range files {
			if !slices.ContainsFunc(out, func(p string) bool { return strings.EqualFold(p, f) }) {
				out = append(out, f)
			}
		}
	}
	return out
}

func (s *Server) listScripts(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	ed := s.back.Editor()
	list := []ScriptEntry{}
	for _, p := range s.paths() {
		e := ScriptEntry{Path: p, Name: strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)), Open: samePath(p, ed.Path)}
		e.Unsaved = e.Open && ed.Dirty
		if sc, _, err := s.back.ReadScript(p); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			e.Error = err.Error()
		} else {
			e.Name, e.Steps = sc.Name, len(sc.Steps)
			if sc.Meta != nil {
				e.CreatedBy = sc.Meta.CreatedBy
			}
		}
		list = append(list, e)
	}
	dir, _ := s.back.ScriptsDir()
	return ok(map[string]any{"scripts": list, "aiFolder": dir})
}

func (s *Server) readScript(_ context.Context, _ *mcp.CallToolRequest, in pathIn) (*mcp.CallToolResult, any, error) {
	path, err := s.resolvePath(in.Path)
	if err != nil {
		return fail(err)
	}
	sc, rev, err := s.back.ReadScript(path)
	if err != nil {
		return fail(err)
	}
	ed := s.back.Editor()
	return ok(map[string]any{"path": path, "revision": rev, "script": sc, "unsaved": samePath(path, ed.Path) && ed.Dirty})
}

func (s *Server) stepTypes(_ context.Context, _ *mcp.CallToolRequest, in typesIn) (*mcp.CallToolResult, any, error) {
	if in.Type != "" {
		for _, d := range model.StepTypeDocs() {
			if d.Type == in.Type {
				return ok(d)
			}
		}
		return fail(&model.ScriptError{Field: "type", Msg: "không có loại node " + in.Type + " – gọi get_step_types không tham số để xem danh sách"})
	}
	return ok(guideDoc(s.back.MCP().NodePreference()))
}

func (s *Server) screenInfo(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	return ok(map[string]any{"displays": s.back.Displays(),
		"note": "Toạ độ theo hệ của chuột (robotgo): Windows là pixel thật, macOS là điểm logic (Retina scale 2 = 1 điểm 2 pixel)."})
}

func samePath(a, b string) bool {
	return a != "" && strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
