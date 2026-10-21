package mcpsrv

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"autoclicker/internal/model"
	"autoclicker/internal/samples"
)

// tips – cách viết kịch bản bền (resource autoclicker://guide, get_step_types).
const tips = `Viết kịch bản bền:
- Ưu tiên mở thẳng URL / app (node open, URL có tham số tìm kiếm), phím tắt (key), tìm chữ trên màn hình (waitImage findBy "text" rồi clickImage imageSource "node") hơn toạ độ cứng – toạ độ hỏng khi cửa sổ dời / đổi cỡ.
- Sau khi mở app / trang: chờ thứ cần thấy (waitImage theo chữ, app waitStart) thay vì chỉ delayAfter cố định; delayAfter vài giây chỉ là dự phòng.
- Thao tác trên app desktop: mở đầu chuẩn app (kiểm tra) → open → app waitStart → window setBounds (xem flow) để cửa sổ luôn đúng chỗ.
- Giá trị người dùng muốn đổi mỗi lần chạy (từ khoá, tên bài, người nhận) → khai báo inputs [{key, label, type, default}] và dùng {{key}} trong ô chữ.
- Chờ có thể hết giờ → nối cổng fail (thử lại bằng counter + cạnh quay lại, hoặc end fail kèm message rõ ràng) thay vì để kịch bản dừng không rõ lý do.
- Toạ độ (click theo điểm) theo hệ của get_screen_info; chưa biết vị trí thì hỏi người dùng hoặc tìm theo chữ.
- Không dùng loại node không có trong danh sách. Ảnh mẫu (template) chỉ chụp được bằng app – AI dùng tìm theo chữ (findBy "text").
- Dạng gọn: node chỉ cần type + trường cần (còn lại lấy mặc định), id / ref tự đặt (s1, s2… / n1, n2…) nếu bỏ trống, không có edges thì nối các node theo thứ tự bằng cổng success.

Thử – sửa (chạy thật trên máy người dùng):
- Chưa biết màn hình có gì: describe_screen (chữ + khung + cửa sổ) trước, screenshot khi cần nhìn bố cục.
- Kịch bản dài: run_script với stopAt chạy từng đoạn, xem output từng node rồi resume_run from node kế để chạy tiếp – không chạy lại đoạn đã đúng.
- Người dùng cầm chuột lúc AI đang chạy → app tự tạm dừng (wait_run: current.pausedBy = "touch") – đừng chạy lại hay tranh chuột; chờ người dùng bấm Tiếp tục trên thanh mini hoặc hỏi người dùng.
- Cần xem màn hình sau từng bước (vd kiểm tra click có mở đúng hộp thoại): run_script / resume_run với shots: "each" – chụp sau mỗi node thành công (sau delayAfter), mỗi node có trường shot; kết quả gửi kèm 3 ảnh mới nhất, ảnh khác xin bằng wait_run shots: [tên]. Ảnh cùng vùng với context.shotArea. Chậm hơn – chỉ dùng khi cần.
- Node lỗi: đọc error.code + error.detail (image: best / threshold / bestBox, colorDiff / colorTol khi kiểm tra màu – best ≥ threshold mà vẫn lỗi = đúng hình sai màu; text: read = các dòng đọc được trong vùng kèm khung, closest = dòng giống nhất, approx = máy đang dùng bộ đọc dự phòng không biết dấu tiếng Việt (dấu trong read có thể sai – đừng sửa findText theo dấu đọc được); element: near = phần tử gần giống (role, name, box, score); web: perSelector = số khớp từng selector, outline = phần tử của trang tên gần giống kèm selector dùng được (sel), matched / count / hidden / covered / moving; window / app: near = cửa sổ / app tên gần giống, running = tiến trình khớp; context: cửa sổ đang trước, chuột) và ảnh lúc lỗi → update_script sửa đúng node → resume_run (chạy tiếp từ node lỗi, output / biến node trước giữ nguyên).
- Toạ độ trên ảnh lúc lỗi → toạ độ chuột: shotArea.x + px × shotArea.w / shotWidth (y tương tự). Ưu tiên sửa sang tìm theo chữ hơn sửa toạ độ.
- Bị từ chối quyền chạy / Cài đặt chỉ cho soạn: báo người dùng, không gọi lại liên tục.`

// guide – nội dung get_step_types / autoclicker://guide.
type guide struct {
	model.FlowGuideDoc
	Tips       string              `json:"tips"`
	StepTypes  []model.StepTypeDoc `json:"stepTypes"`
	CommonOuts []string            `json:"commonOutputs"`
	NodePref   string              `json:"nodePref"`     // ưu tiên node hiện tại (44)
	NodeTips   string              `json:"nodePrefTips"` // cách chọn node web / ảnh·chữ theo nodePref
}

// guideDoc – pref: ưu tiên node hiện tại (store.MCPSettings.NodePreference).
func guideDoc(pref string) guide {
	return guide{FlowGuideDoc: model.Guide(), Tips: tips, StepTypes: model.StepTypeDocs(), CommonOuts: model.CommonOutputs,
		NodePref: pref, NodeTips: nodePrefTips[pref]}
}

// exampleIDs – mẫu thử (bước 33) làm ví dụ: chạy được trên máy thật, không phụ thuộc sân thử / ảnh mẫu.
var exampleIDs = []string{"mo-spotify-phat", "google-tim-kiem", "ung-dung-mo-dong", "http-goi-api"}

// ScriptSchema – JSON Schema của file kịch bản, sinh từ struct Go (khớp model.go). Mọi trường đều tuỳ chọn khi gửi
// cho create_script (dạng gọn), nên bỏ "required".
func ScriptSchema() (*jsonschema.Schema, error) {
	opts := &jsonschema.ForOptions{TypeSchemas: map[reflect.Type]*jsonschema.Schema{
		// Cond.Conds lồng chính Cond – trỏ $defs để không vòng
		reflect.TypeFor[[]model.Cond](): {Type: "array", Items: &jsonschema.Schema{Ref: "#/$defs/cond"}},
	}}
	sc, err := jsonschema.For[model.Script](opts)
	if err != nil {
		return nil, err
	}
	cond, err := jsonschema.For[model.Cond](opts)
	if err != nil {
		return nil, err
	}
	sc.Defs = map[string]*jsonschema.Schema{"cond": cond}
	sc.Title = "Kịch bản Auto Clicker"
	sc.Description = "Flow node + cạnh. Chi tiết từng loại node: get_step_types / autoclicker://guide."
	sc.Properties["version"].Description = "định dạng file – app tự ghi"
	step := sc.Properties["steps"].Items
	types := []any{}
	for _, d := range model.StepTypeDocs() {
		types = append(types, d.Type)
	}
	step.Properties["type"].Enum = types
	dropRequired(sc)
	return sc, nil
}

func dropRequired(s *jsonschema.Schema) {
	if s == nil {
		return
	}
	s.Required = nil
	s.AdditionalProperties = nil
	for _, p := range s.Properties {
		dropRequired(p)
	}
	for _, d := range s.Defs {
		dropRequired(d)
	}
	dropRequired(s.Items)
}

func (s *Server) addDocs() {
	res := func(uri, name, desc, mime string, body func() (any, error)) {
		s.mcp.AddResource(&mcp.Resource{URI: uri, Name: name, Description: desc, MIMEType: mime},
			func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
				v, err := body()
				if err != nil {
					return nil, err
				}
				b, err := json.MarshalIndent(v, "", " ")
				if err != nil {
					return nil, err
				}
				return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: uri, MIMEType: mime, Text: string(b)}}}, nil
			})
	}
	res("autoclicker://guide", "guide", "Hướng dẫn viết kịch bản: flow, biến, từng loại node, mẹo viết bền", "application/json",
		func() (any, error) { return guideDoc(s.back.MCP().NodePreference()), nil })
	res("autoclicker://schema", "schema", "JSON Schema của file kịch bản", "application/schema+json",
		func() (any, error) { return ScriptSchema() })
	res("autoclicker://examples", "examples", "Kịch bản mẫu dạng gọn (chạy được với create_script)", "application/json", examples)

	s.mcp.AddPrompt(&mcp.Prompt{Name: "tao_kich_ban", Title: "Tạo kịch bản Auto Clicker",
		Description: "Hỏi đủ thông tin rồi tạo kịch bản tự động thao tác trên máy",
		Arguments:   []*mcp.PromptArgument{{Name: "muc_tieu", Description: "muốn tự động làm gì, vd: mở bài Tinh Hà Say 2 trên YouTube"}}},
		func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return &mcp.GetPromptResult{Messages: []*mcp.PromptMessage{{Role: "user",
				Content: &mcp.TextContent{Text: promptText(req.Params.Arguments["muc_tieu"])}}}}, nil
		})
}

func examples() (any, error) {
	var out []map[string]any
	for _, id := range exampleIDs {
		raw, err := samples.Raw(id)
		if err != nil {
			return nil, err
		}
		out = append(out, map[string]any{"id": id, "script": raw})
	}
	return map[string]any{"full": json.RawMessage(model.Guide().Example), "compact": out}, nil
}

func promptText(goal string) string {
	if goal == "" {
		goal = "(chưa nói – hỏi người dùng)"
	}
	return `Tôi muốn tạo kịch bản Auto Clicker để: ` + goal + `

Trước khi tạo, hỏi lại tôi những gì còn thiếu (mỗi lần vài câu ngắn):
1. Mục tiêu cụ thể: kết thúc thế nào thì coi là xong?
2. Làm trên app / trang web nào? App đã mở sẵn chưa, cửa sổ nằm đâu?
3. Giá trị nào đổi mỗi lần chạy (thành input {{key}}), giá trị nào cố định?
4. Lặp bao nhiêu lần, lỗi thì thử lại hay dừng?

Sau đó: get_step_types → viết kịch bản (ưu tiên mở URL / phím tắt / tìm chữ hơn toạ độ) → validate_script tới khi ok → create_script → open_in_app, rồi tóm tắt các bước cho tôi.`
}
