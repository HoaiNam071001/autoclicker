// Package samples – kịch bản mẫu để test từng loại bước (bước 33) + trang sân thử cục bộ chúng chạy trên.
// Mẫu là JSON gọn {meta, script}: node chỉ ghi trường cần, còn lại lấy mặc định của model.NewStep; không có
// edges thì nối các node theo thứ tự bằng cổng success. Không phụ thuộc UI / thiết bị – MCP (31) dùng lại được.
package samples

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
)

//go:embed scripts/*.json
var scriptFS embed.FS

// SandboxVar – thay bằng URL file:// của thư mục sân thử lúc mở mẫu (target của node Mở).
const SandboxVar = "{{sandbox}}"

// templatePrefix – "template": "@logo.png" = ảnh mẫu lấy từ file trong pages/ (base64 lúc mở mẫu).
const templatePrefix = "@"

// Info – mô tả một mẫu cho danh sách. Steps tính từ kịch bản (không ghi tay để khỏi lệch).
type Info struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	App         string   `json:"app"`   // kiểu app: Cơ bản, Form, Shop, Chat…
	Level       int      `json:"level"` // 1 dễ – 3 khó
	Page        string   `json:"page"`  // trang sân thử (vd click.html) hoặc URL web thật; "" = mẫu chỉ dùng app của hệ điều hành
	Real        bool     `json:"real"`  // chạy trên web thật – có thể hỏng khi trang đổi giao diện
	Prep        string   `json:"prep,omitempty"`
	Steps       []string `json:"steps"` // các loại bước có trong mẫu, theo thứ tự xuất hiện
}

type sampleFile struct {
	Meta   Info `json:"meta"`
	Script struct {
		Loops  int                 `json:"loops"`
		Region *model.Rect         `json:"region"`
		Start  string              `json:"start"`
		Steps  []json.RawMessage   `json:"steps"`
		Edges  []model.Edge        `json:"edges"`
		Inputs []model.ScriptInput `json:"inputs"`
	} `json:"script"`
}

// List trả mọi mẫu, sắp theo độ khó rồi tên.
func List() ([]Info, error) {
	ids, err := ids()
	if err != nil {
		return nil, err
	}
	out := make([]Info, 0, len(ids))
	for _, id := range ids {
		info, sc, err := load(id)
		if err != nil {
			return nil, err
		}
		info.Steps = stepTypes(sc)
		out = append(out, info)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Level != out[j].Level {
			return out[i].Level < out[j].Level
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}

// Open trả kịch bản của mẫu id với {{sandbox}} trỏ tới sandboxURL (URL file:// của thư mục đã WriteSandbox).
func Open(id, sandboxURL string) (model.Script, error) {
	_, sc, err := load(id)
	if err != nil {
		return model.Script{}, err
	}
	if err := fill(&sc, sandboxURL); err != nil {
		return model.Script{}, fmt.Errorf("mẫu %s: %w", id, err)
	}
	// qua đúng đường mở file để có cùng kiểm tra / sửa ref như kịch bản người dùng
	b, err := script.EncodeScript(sc)
	if err != nil {
		return model.Script{}, err
	}
	return script.DecodeScript(b)
}

func ids() ([]string, error) {
	entries, err := fs.ReadDir(scriptFS, "scripts")
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			out = append(out, strings.TrimSuffix(e.Name(), ".json"))
		}
	}
	return out, nil
}

// load đọc file mẫu → Info + kịch bản đủ trường (chưa thay {{sandbox}} / ảnh mẫu).
func load(id string) (Info, model.Script, error) {
	b, err := scriptFS.ReadFile(path.Join("scripts", id+".json"))
	if err != nil {
		return Info{}, model.Script{}, fmt.Errorf("không có mẫu %q", id)
	}
	var f sampleFile
	if err := json.Unmarshal(b, &f); err != nil {
		return Info{}, model.Script{}, fmt.Errorf("mẫu %s hỏng: %w", id, err)
	}
	sc := model.NewScript()
	sc.Name, sc.Region, sc.Inputs = f.Meta.Name, f.Script.Region, f.Script.Inputs
	if f.Script.Loops > 0 {
		sc.Loops = f.Script.Loops
	}
	for i, raw := range f.Script.Steps {
		st, err := decodeStep(raw, i)
		if err != nil {
			return Info{}, model.Script{}, fmt.Errorf("mẫu %s node %d: %w", id, i+1, err)
		}
		sc.Steps = append(sc.Steps, st)
	}
	sc.Edges = f.Script.Edges
	if len(sc.Edges) == 0 {
		for i := 1; i < len(sc.Steps); i++ {
			sc.Edges = append(sc.Edges, model.Edge{From: sc.Steps[i-1].ID, Port: model.PortSuccess, To: sc.Steps[i].ID})
		}
	}
	for i := range sc.Edges {
		if sc.Edges[i].ID == "" {
			sc.Edges[i].ID = fmt.Sprintf("e%d", i+1)
		}
	}
	sc.Start = f.Script.Start
	if sc.Start == "" && len(sc.Steps) > 0 {
		sc.Start = sc.Steps[0].ID
	}
	f.Meta.ID = id
	return f.Meta, sc, nil
}

// decodeStep: mặc định theo loại (NewStep) rồi đè trường ghi trong mẫu; id / ref mặc định s<i> / n<i>.
func decodeStep(raw json.RawMessage, i int) (model.Step, error) {
	var head struct{ Type string }
	if err := json.Unmarshal(raw, &head); err != nil {
		return model.Step{}, err
	}
	st, err := model.NewStep(head.Type)
	if err != nil {
		return model.Step{}, err
	}
	st.ID, st.Ref = "", ""
	if err := json.Unmarshal(raw, &st); err != nil {
		return model.Step{}, err
	}
	if st.ID == "" {
		st.ID = fmt.Sprintf("s%d", i+1)
	}
	if st.Ref == "" {
		st.Ref = fmt.Sprintf("n%d", i+1)
	}
	return st, nil
}

// fill thay {{sandbox}} trong target và "@file" trong ảnh mẫu (node, điều kiện, case).
func fill(sc *model.Script, sandboxURL string) error {
	var err error
	tpl := func(t *string) {
		if err == nil && strings.HasPrefix(*t, templatePrefix) {
			*t, err = pageBase64(strings.TrimPrefix(*t, templatePrefix))
		}
	}
	var cond func(c *model.Cond)
	cond = func(c *model.Cond) {
		if c == nil {
			return
		}
		tpl(&c.Template)
		for i := range c.Conds {
			cond(&c.Conds[i])
		}
	}
	for i := range sc.Steps {
		st := &sc.Steps[i]
		st.Target = strings.ReplaceAll(st.Target, SandboxVar, sandboxURL)
		tpl(&st.Template)
		cond(st.Cond)
		for j := range st.Cases {
			cond(st.Cases[j].Cond)
		}
	}
	return err
}

func pageBase64(name string) (string, error) {
	b, err := pageFS.ReadFile(path.Join("pages", name))
	if err != nil {
		return "", fmt.Errorf("không có ảnh mẫu %s trong sân thử", name)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// stepTypes – các loại bước trong kịch bản, không trùng, theo thứ tự xuất hiện.
func stepTypes(sc model.Script) []string {
	seen := map[string]bool{}
	var out []string
	for _, st := range sc.Steps {
		if !seen[st.Type] {
			seen[st.Type] = true
			out = append(out, st.Type)
		}
	}
	return out
}
