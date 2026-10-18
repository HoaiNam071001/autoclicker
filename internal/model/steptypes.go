package model

// Mô tả các loại bước cho người đọc file kịch bản và cho AI (MCP, bước 31 – sinh schema / hướng dẫn từ bảng này).
// Thêm loại bước mới: thêm hằng trong model.go, xử lý trong script (checkStep / ValidateScript) / engine, và một dòng ở đây
// (test kiểm tra mọi loại đều có mô tả và ví dụ hợp lệ).

type StepTypeDoc struct {
	Type     string   `json:"type"`
	Label    string   `json:"label"`    // tên tiếng Việt trên UI
	Desc     string   `json:"desc"`     // làm gì, khi nào dùng
	Fields   []string `json:"fields"`   // trường riêng của loại này (ngoài trường chung, xem commonFieldsDoc)
	Pinnable []string `json:"pinnable"` // trường ghim được làm input (script/inputvalue.go)
	Example  string   `json:"example"`  // một node JSON hợp lệ
	Outputs  []string `json:"outputs"`  // khoá output khi thành công – node sau đọc bằng {{ref.output.khoá}}
	Ports    []string `json:"ports"`    // cổng ra (PortsOf) – điền lúc khởi động
}

func init() {
	for i := range stepTypeDocs {
		stepTypeDocs[i].Ports = portsOf(Step{Type: stepTypeDocs[i].Type})
	}
}

// commonFieldsDoc – trường mọi loại node đều có.
const commonFieldsDoc = `id (ổn định – cạnh, start, imageFrom trỏ theo id), ref (mã ngắn duy nhất, vd "n2" – node sau đọc kết quả bằng {{n2.output.x}}), ` +
	`type, name (tên hiện trên thẻ, tuỳ chọn), enabled (false = bỏ qua, đi tiếp cổng success), delayAfter (ms chờ sau khi thành công), ` +
	`jitterMs (± ms ngẫu nhiên), pos {x, y} (vị trí thẻ trên canvas, tuỳ chọn). ` +
	`retry {count, delayMs}: thất bại thì thử lại thêm count lần (0–100), nghỉ delayMs giữa các lần; hết lượt vẫn lỗi → cổng "fail".`

// flowDoc – cấu trúc kịch bản dạng flow (bước 20).
const flowDoc = `Kịch bản là flow: steps = danh sách node (thứ tự trong mảng không quyết định thứ tự chạy), start = id node chạy đầu, ` +
	`edges = [{id, from, port, to}]: node from chạy xong đi ra cổng port thì chạy node to; mỗi cổng tối đa một cạnh, không nối node vào chính nó. ` +
	`Node hành động có cổng "success" / "fail". Cổng fail không nối = lỗi thì dừng kịch bản và báo lỗi; cổng khác không nối = kết thúc vòng. ` +
	`Vòng lặp = cạnh quay lại node trước (tối đa 10000 lần chạy node mỗi vòng). ` +
	`Biến trong ô chữ (text, target, findText): {{key}} = input của kịch bản; {{n2.status}} (success | fail | skipped), {{n2.port}}, ` +
	`{{n2.error.code}} / {{n2.error.message}}, {{n2.output.<khoá>}} = kết quả lần chạy gần nhất của node có ref n2 trong vòng này; ` +
	`{{last.…}} = node vừa chạy trước; {{vars.<tên>}} = biến node setVar đặt. Biến nhận được trong: text, target, findText, ` +
	`cond (findText, left, right), cases (match, cond), value, message. ` +
	`Node điều khiển rẽ nhánh theo cổng: if (true / false), switch / switchScreen (case:<key> theo từng case, rồi default / none), ` +
	`counter (next / done – dùng để giới hạn vòng lặp), end (không có cổng: kết thúc vòng, fail = dừng kịch bản báo lỗi), setVar (success). ` +
	`Điều kiện cond: {kind: text | image | pixel | value | all | any, not}: text = thấy chữ findText (textExact, textNoAccent, occurrence, searchArea); ` +
	`image = thấy ảnh mẫu template (templateScale, threshold, searchArea); pixel = points[0] có màu color ± tolerance; ` +
	`value = so left op right (op: eq | ne | contains (không phân biệt hoa thường) | gt | lt (so số) | empty); all / any = conds[] (và / hoặc). ` +
	`Kết quả mỗi node mỗi node: {nodeId, ref, status, port, output, error {code, message}, logs, attempts, durationMs}; ` +
	`error.code: not_found (chờ ảnh / chữ / màu hết giờ), permission, invalid_input (biến trỏ node chưa chạy…), end (node Kết thúc thất bại), internal.`

// flowExample – kịch bản mẫu: mở trang → chờ chữ "Nhắn tin" → click → gõ; không thấy thì cuộn rồi chờ lại, tối đa
// 10 lần (node Đếm) rồi kết thúc báo lỗi (test chạy qua script.ValidateScript).
const flowExample = `{"version":4,"name":"Nhắn tin","loops":1,"start":"s1",
"steps":[
 {"id":"s1","ref":"n1","type":"open","target":"facebook.com/{{profile}}","delayAfter":3000,"enabled":true},
 {"id":"s2","ref":"n2","type":"waitImage","findBy":"text","findText":"Nhắn tin","threshold":90,"timeoutMs":3000,"intervalMs":250,"searchArea":"screen","enabled":true},
 {"id":"s3","ref":"n3","type":"clickImage","imageSource":"node","imageFrom":"s2","button":"left","delayAfter":500,"enabled":true},
 {"id":"s4","ref":"n4","type":"scroll","direction":"down","amount":3,"delayAfter":500,"enabled":true},
 {"id":"s5","ref":"n5","type":"type","text":"xin chào – thấy nút ở {{n2.output.x}}, {{n2.output.y}}","enabled":true},
 {"id":"s6","ref":"n6","type":"counter","max":10,"enabled":true},
 {"id":"s7","ref":"n7","type":"end","endStatus":"fail","message":"Không thấy nút Nhắn tin sau {{n6.output.count}} lần","enabled":true}],
"edges":[
 {"id":"e1","from":"s1","port":"success","to":"s2"},
 {"id":"e2","from":"s2","port":"success","to":"s3"},
 {"id":"e3","from":"s2","port":"fail","to":"s6"},
 {"id":"e4","from":"s6","port":"next","to":"s4"},
 {"id":"e5","from":"s6","port":"done","to":"s7"},
 {"id":"e6","from":"s4","port":"success","to":"s2"},
 {"id":"e7","from":"s3","port":"success","to":"s5"}],
"inputs":[{"key":"profile","type":"text","default":"lan.nguyen"}]}`

var stepTypeDocs = []StepTypeDoc{
	{StepClick, "Click", "Click chuột Clicks lần tại points[0] (neo theo cạnh vùng, xem script/anchor.go).",
		[]string{"points[1]", "button", "clicks", "clickInterval", "jitterPx"}, []string{"clicks", "clickInterval", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"click","points":[{"x":400,"y":300}],"button":"left","clicks":1,"delayAfter":500,"enabled":true}`, []string{"x", "y"}, nil},
	{StepDouble, "Double click", "Double click tại points[0].",
		[]string{"points[1]", "button", "jitterPx"}, []string{"jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"double","points":[{"x":400,"y":300}],"button":"left","delayAfter":500,"enabled":true}`, []string{"x", "y"}, nil},
	{StepHold, "Giữ chuột", "Nhấn giữ nút chuột holdMs mili giây tại points[0].",
		[]string{"points[1]", "button", "holdMs", "jitterPx"}, []string{"holdMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"hold","points":[{"x":400,"y":300}],"button":"left","holdMs":1000,"delayAfter":500,"enabled":true}`, []string{"x", "y"}, nil},
	{StepDrag, "Kéo thả", "Giữ chuột tại points[0] kéo tới points[1] trong dragMs. dragFromPrev: bắt đầu từ chỗ chuột dừng ở node chuột chạy trước (chưa có thì lỗi); dragRelative: points[1] là độ lệch so với điểm bắt đầu.",
		[]string{"points[2]", "button", "dragMs", "dragFromPrev", "dragRelative", "jitterPx"}, []string{"dragMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"drag","points":[{"x":100,"y":300},{"x":600,"y":300}],"button":"left","dragMs":400,"delayAfter":500,"enabled":true}`, []string{"x", "y"}, nil},
	{StepWait, "Chờ", "Chỉ chờ delayAfter mili giây.",
		[]string{}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"wait","delayAfter":2000,"enabled":true}`, []string{}, nil},
	{StepText, "Gõ chữ", "Gõ text vào ô đang focus (thường sau một bước click). Nhận {{biến}}; \\{{ để gõ ra {{.",
		[]string{"text"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"type","text":"xin chào {{ten}}","delayAfter":300,"enabled":true}`, []string{"text"}, nil},
	{StepKey, "Phím tắt", "Bấm keys (tên phím robotgo: enter, tab, esc, f5, a…) kèm modifiers (cmd, ctrl, alt, shift). keyMode chord: nhấn giữ cùng lúc; sequence: giữ phím kèm, gõ lần lượt.",
		[]string{"keys", "modifiers", "keyMode"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"key","keys":["l"],"modifiers":["cmd"],"keyMode":"chord","delayAfter":300,"enabled":true}`, []string{}, nil},
	{StepScroll, "Cuộn", "Cuộn amount nấc theo direction (up | down | left | right) tại points[0]; không có điểm thì cuộn ở chỗ chuột đang đứng.",
		[]string{"direction", "amount", "points[0..1]", "jitterPx"}, []string{"amount", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"scroll","direction":"down","amount":5,"delayAfter":500,"enabled":true}`, []string{}, nil},
	{StepOpen, "Mở URL / app", "Mở target: URL (https://… hoặc tên miền như youtube.com/…), đường dẫn file / app, hoặc tên app (Spotify). app: mở bằng app này, vd Google Chrome; trống = mặc định của máy. Nhận {{biến}} (trong URL tự mã hoá). Không chờ trang tải xong – đặt delayAfter đủ lâu hoặc thêm bước chờ.",
		[]string{"target", "app"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"open","target":"https://www.youtube.com/results?search_query={{bai_hat}}","app":"Google Chrome","delayAfter":3000,"enabled":true}`, []string{"target"}, nil},
	{StepWaitImage, "Chờ ảnh", "Chờ tới khi ảnh mẫu template (PNG base64, chụp bằng app – CaptureTemplate) hiện trên màn hình, độ khớp ≥ threshold %. Hết timeoutMs không thấy → cổng fail (error.code not_found). Tìm trong vùng của kịch bản, hoặc searchArea \"screen\" = cả màn chứa vùng. Output x, y (tâm), w, h, score – node clickImage chọn imageFrom = id node này để click vào đó. findBy \"text\": tìm chữ findText thay ảnh mẫu (đọc chữ trên màn hình, tiếng Việt / Anh; không phân biệt hoa thường và khoảng trắng; nhận {{biến}}); textExact = cả dòng đúng chuỗi, textNoAccent = bỏ qua dấu, occurrence = chỗ thứ mấy theo thứ tự đọc. Chữ bền hơn ảnh khi đổi theme / zoom – nên dùng khi nút có chữ.",
		[]string{"template", "templateScale", "threshold", "timeoutMs", "intervalMs", "searchArea", "findBy", "findText", "textExact", "textNoAccent", "occurrence"}, []string{"timeoutMs", "threshold", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"waitImage","template":"iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAAAAADhZOFXAAAAHUlEQVR4nGNkYPjPwMjwn4GJAQb+/4dguAgjEWoAyH0MAWBh3h0AAAAASUVORK5CYII=","templateScale":1,"threshold":90,"timeoutMs":10000,"intervalMs":250,"delayAfter":300,"enabled":true}`, []string{"x", "y", "w", "h", "score", "text"}, nil},
	{StepClickImage, "Click theo ảnh", "Click vào tâm chỗ thấy ảnh / chữ (+ dx, dy). imageSource \"node\": chỗ node imageFrom (id node chờ ảnh / chữ, hoặc click theo ảnh tự tìm) thấy lần gần nhất trong vòng – node đó chưa thấy thì lỗi; \"search\": tự tìm ảnh mẫu / chữ riêng (các trường như waitImage, gồm findBy \"text\"). Không thấy → cổng fail.",
		[]string{"imageSource", "dx", "dy", "button", "template", "templateScale", "threshold", "timeoutMs", "intervalMs", "searchArea", "findBy", "findText", "textExact", "textNoAccent", "occurrence", "jitterPx"}, []string{"timeoutMs", "threshold", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","ref":"n2","type":"clickImage","imageSource":"node","imageFrom":"s1","dx":0,"dy":0,"button":"left","delayAfter":500,"enabled":true}`, []string{"clickX", "clickY", "x", "y", "w", "h", "score", "text"}, nil},
	{StepWaitPixel, "Chờ màu pixel", "Chờ tới khi điểm points[0] có màu color (#RRGGBB), mỗi kênh lệch ≤ tolerance. Hết timeoutMs → cổng fail. Không di chuột. Output color = màu đọc được.",
		[]string{"points[1]", "color", "tolerance", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"waitPixel","points":[{"x":640,"y":300}],"color":"#1877F2","tolerance":10,"timeoutMs":5000,"intervalMs":100,"delayAfter":300,"enabled":true}`, []string{"color"}, nil},
	{StepIf, "Nếu", "Rẽ nhánh theo điều kiện cond (xem flow guide): đúng → cổng true, sai → false. timeoutMs > 0: xem lại mỗi intervalMs tới khi đúng, hết giờ = sai (0 = xem một lần). Thiếu quyền / ảnh hỏng → node lỗi, dừng kịch bản. Output result (true / false) + chỗ thấy x, y, w, h, score, text của điều kiện ảnh / chữ đúng đầu tiên (clickImage chọn imageFrom = node này được).",
		[]string{"cond", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"if","cond":{"kind":"any","conds":[{"kind":"text","findText":"Đã gửi"},{"kind":"value","left":"{{n2.output.text}}","op":"contains","right":"ok"}]},"timeoutMs":3000,"intervalMs":250,"enabled":true}`,
		[]string{"result", "x", "y", "w", "h", "score", "text"}, nil},
	{StepSwitch, "Switch theo giá trị", "So value (thường là {{biến}}) với từng case theo thứ tự: case {key, label, op: eq | contains (không phân biệt hoa thường) | regex, match}; case đầu khớp → cổng case:<key>, không case nào khớp → default.",
		[]string{"value", "cases"}, []string{"enabled"},
		`{"id":"s1","ref":"n1","type":"switch","value":"{{n2.output.text}}","cases":[{"key":"ok","label":"Thành công","op":"contains","match":"đã gửi"},{"key":"err","label":"Lỗi","op":"regex","match":"(?i)lỗi|error"}],"enabled":true}`,
		[]string{"value", "case"}, nil},
	{StepSwitchScreen, "Switch theo màn hình", "Chụp màn hình mỗi intervalMs, xem case nào (cases[].cond: thấy chữ / ảnh / màu) đúng trước – xét theo thứ tự case – thì đi cổng case:<key>; hết timeoutMs không case nào → cổng none. Dùng khi trang có thể ở nhiều trạng thái (đăng nhập / đã vào / popup).",
		[]string{"cases", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"switchScreen","cases":[{"key":"login","label":"Đăng nhập","cond":{"kind":"text","findText":"Đăng nhập"}},{"key":"home","label":"Trang chủ","cond":{"kind":"text","findText":"Nhắn tin"}}],"timeoutMs":3000,"intervalMs":250,"enabled":true}`,
		[]string{"case", "x", "y", "w", "h", "score", "text"}, nil},
	{StepCounter, "Đếm", "Đếm số lần đi qua trong vòng này: lần 1..max → cổng next, lần max+1 → done (rồi đếm lại từ 0). Đặt trên đường quay lại của vòng lặp để giới hạn số lần thử (vd không thấy nút → Đếm → next: cuộn rồi tìm lại; done: Kết thúc báo lỗi).",
		[]string{"max"}, []string{"enabled"},
		`{"id":"s1","ref":"n1","type":"counter","max":10,"enabled":true}`, []string{"count"}, nil},
	{StepEnd, "Kết thúc", "Kết thúc vòng chạy: endStatus success = xong vòng bình thường (còn vòng lặp thì chạy vòng sau); fail = dừng kịch bản, báo lỗi message (error.code end). Không có cổng ra; thành công có message thì output {message} – hiện cho người dùng khi chạy xong. Nhận {{biến}} trong message.",
		[]string{"endStatus", "message"}, []string{"enabled"},
		`{"id":"s1","ref":"n1","type":"end","endStatus":"fail","message":"Không thấy nút Nhắn tin","enabled":true}`, []string{"message"}, nil},
	{StepSetVar, "Đặt biến", "Đặt biến {{vars.varName}} dùng ở node sau trong vòng: varOp set = gán value (nhận {{biến}}, vd {{n2.output.text}}); add = cộng value (số) vào giá trị hiện có (chưa có = 0). Biến đặt lại mỗi vòng.",
		[]string{"varName", "varOp", "value"}, []string{"enabled"},
		`{"id":"s1","ref":"n1","type":"setVar","varName":"so_lan","varOp":"add","value":"1","enabled":true}`, []string{"name", "value"}, nil},
}

type FlowGuideDoc struct {
	Flow         string `json:"flow"`
	CommonFields string `json:"commonFields"`
	Example      string `json:"example"`
}

// StepTypeDocs – mô tả các loại bước app đang hỗ trợ.
func StepTypeDocs() []StepTypeDoc { return stepTypeDocs }

// Guide – hướng dẫn cấu trúc flow, trường chung và kịch bản mẫu (cho AI / MCP, bước 31).
func Guide() FlowGuideDoc {
	return FlowGuideDoc{Flow: flowDoc, CommonFields: commonFieldsDoc, Example: flowExample}
}
