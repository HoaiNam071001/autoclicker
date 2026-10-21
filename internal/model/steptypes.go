package model

import "slices"

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
	Outputs  []string `json:"outputs"`  // khoá output riêng khi thành công (cộng CommonOutputs) – node sau đọc bằng {{ref.output.khoá}}
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
	`retry {count, delayMs}: thất bại thì thử lại thêm count lần (0–100), nghỉ delayMs giữa các lần; hết lượt vẫn lỗi → cổng "fail". ` +
	`human: mức thao tác giống người riêng của node ("" = theo kịch bản | off | low | medium | high).`

// flowDoc – cấu trúc kịch bản dạng flow (bước 20).
const flowDoc = `Kịch bản là flow: steps = danh sách node (thứ tự trong mảng không quyết định thứ tự chạy), start = id node chạy đầu, ` +
	`edges = [{id, from, port, to}]: node from chạy xong đi ra cổng port thì chạy node to; mỗi cổng tối đa một cạnh, không nối node vào chính nó. ` +
	`Node hành động có cổng "success" / "fail". Cổng fail không nối = lỗi thì dừng kịch bản và báo lỗi; cổng khác không nối = kết thúc vòng. ` +
	`Vòng lặp = cạnh quay lại node trước (tối đa 10000 lần chạy node mỗi vòng). ` +
	`Mở đầu chuẩn khi thao tác trên một app: app (kiểm tra) → cổng "no" → open → app waitStart → window setBounds; cổng "yes" nối thẳng window setBounds – lần chạy nào app cũng đã mở, cửa sổ đúng chỗ, đúng cỡ trước khi click toạ độ / tìm ảnh. ` +
	`Script.human (off | low | medium | high, trống = tắt; kịch bản mới: medium): thao tác giống người – chuột đi đường cong có ` +
	`tăng / giảm tốc (lâu hơn khi xa / đích nhỏ, đôi khi đi lố rồi chỉnh), click nhấn – giữ – nhả, điểm bấm lệch quanh tâm chỗ thấy, ` +
	`gõ / phím / cuộn nhịp không đều (high: đôi khi gõ nhầm rồi xoá); chậm hơn nhưng giảm nguy cơ bị nhận là bot (không đảm bảo). ` +

	`Biến trong ô chữ (text, target, findText): {{key}} = input của kịch bản; {{n2.status}} (success | fail | skipped), {{n2.port}}, ` +
	`{{n2.error.code}} / {{n2.error.message}}, {{n2.output.<khoá>}} = kết quả lần chạy gần nhất của node có ref n2 trong vòng này ` +
	`(khoá theo outputs của loại node, cộng startedAt / durationMs mọi node đều có; window / windowApp của click / double / drag / type / key = ` +
	`cửa sổ đang ở trước sau thao tác – kiểm tra có gõ / bấm vào đúng app không; chưa có trên macOS); ` +
	`{{last.…}} = node vừa chạy trước; {{vars.<tên>}} = biến node setVar đặt. Biến nhận được trong: text, target, findText, selectors, tab, ` +
	`elementApp, cond (findText, selectors, elementApp, left, right), cases (match, cond), value, message. ` +
	`Điểm trên màn hình (points của click, double, hold, drag, scroll, waitPixel) nhận toạ độ từ biến: {"xVar":"{{n2.output.x}} + 20","yVar":"{{n2.output.y}}"} – ` +
	`có xVar / yVar thì bỏ qua x, y và neo; mỗi ô là số hoặc biến cộng / trừ px, biến node phải là output của node đó (vd x, y, w, h của waitImage); ` +
	`kết quả là toạ độ màn hình (như output waitImage), không tính theo vùng; drag dragRelative: points[1] từ biến là độ lệch so với A. ` +
	`Trường số (clicks, clickInterval, holdMs, dragMs, delayAfter, jitterMs, jitterPx, amount, timeoutMs, threshold – trong pinnable của loại node) ` +
	`nhận biến qua numVars: [{"field":"amount","expr":"{{n2.output.count}} + 1"}] – lúc chạy giá trị expr (số hoặc biến cộng / trừ, như xVar) ` +
	`thay cho trường đó; biến node phải là output của node đó; ra số sai khoảng (vd clicks < 1, số âm) → node fail invalid_input. ` +
	`Node điều khiển rẽ nhánh theo cổng: if (true / false), switch / switchScreen (case:<key> theo từng case, rồi default / none), ` +
	`counter (next / done – dùng để giới hạn vòng lặp), end (không có cổng: kết thúc vòng, fail = dừng kịch bản báo lỗi), setVar (success). ` +
	`Điều kiện cond: {kind: text | image | element | pixel | value | all | any, not}: text = thấy chữ findText (textExact, textAccent, occurrence, searchArea); ` +
	`element = có phần tử selectors trong cửa sổ app elementApp (như findBy "element" của waitImage; không cuộn); ` +
	`image = thấy ảnh mẫu template (templateScale, threshold, matchMode, colorTol, searchArea – như waitImage); pixel = points[0] có màu color ± tolerance; ` +
	`value = so left op right (op: eq | ne | contains (không phân biệt hoa thường) | gt | lt (so số) | empty); all / any = conds[] (và / hoặc). ` +
	`Kết quả mỗi node: {nodeId, ref, status, port, output, error {code, message}, logs, attempts, durationMs}; ` +
	`error.code: not_found (chờ ảnh / chữ / màu / phần tử hết giờ), browser (chưa mở Chrome của app), permission (Screen Recording / Accessibility), invalid_input (biến trỏ node chưa chạy…), end (node Kết thúc thất bại), internal.`

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
	{StepClick, "Click", "Click chuột Clicks lần tại points[0] (neo theo cạnh vùng, xem script/anchor.go). coordSpace \"absolute\" = toạ độ màn hình thật như ảnh screenshot (mặc định cộng gốc vùng của kịch bản).",
		[]string{"points[1]", "coordSpace", "button", "clicks", "clickInterval", "jitterPx"}, []string{"clicks", "clickInterval", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"click","points":[{"x":400,"y":300}],"button":"left","clicks":1,"delayAfter":500,"enabled":true}`, []string{"x", "y", "window", "windowApp"}, nil},
	{StepDouble, "Double click", "Double click tại points[0].",
		[]string{"points[1]", "coordSpace", "button", "jitterPx"}, []string{"jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"double","points":[{"x":400,"y":300}],"button":"left","delayAfter":500,"enabled":true}`, []string{"x", "y", "window", "windowApp"}, nil},
	{StepHold, "Giữ chuột", "Nhấn giữ nút chuột holdMs mili giây tại points[0].",
		[]string{"points[1]", "coordSpace", "button", "holdMs", "jitterPx"}, []string{"holdMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"hold","points":[{"x":400,"y":300}],"button":"left","holdMs":1000,"delayAfter":500,"enabled":true}`, []string{"x", "y", "ms"}, nil},
	{StepDrag, "Kéo thả", "Giữ chuột tại points[0] kéo tới points[1] trong dragMs. dragFromPrev: bắt đầu từ chỗ chuột dừng ở node chuột chạy trước (chưa có thì lỗi); dragRelative: points[1] là độ lệch so với điểm bắt đầu.",
		[]string{"points[2]", "coordSpace", "button", "dragMs", "dragFromPrev", "dragRelative", "jitterPx"}, []string{"dragMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"drag","points":[{"x":100,"y":300},{"x":600,"y":300}],"button":"left","dragMs":400,"delayAfter":500,"enabled":true}`, []string{"x", "y", "window", "windowApp"}, nil},
	{StepWait, "Chờ", "Chỉ chờ delayAfter mili giây.",
		[]string{}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"wait","delayAfter":2000,"enabled":true}`, []string{"ms"}, nil},
	{StepText, "Gõ chữ", "Gõ text vào ô đang focus (thường sau một bước click). Nhận {{biến}}; \\{{ để gõ ra {{. typeMode: trống = tự chọn (Windows + không giống người: dán, còn lại gõ phím) | \"keys\" = gõ từng phím (nhịp giống người; bộ gõ tiếng Việt UniKey / EVKey đang bật có thể làm méo chữ) | \"paste\" = dán qua clipboard rồi trả chữ cũ (không bị bộ gõ làm méo; ô chặn dán thì không nhập được).",
		[]string{"text", "typeMode"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"type","text":"xin chào {{ten}}","delayAfter":300,"enabled":true}`, []string{"text", "window", "windowApp"}, nil},
	{StepKey, "Phím tắt", "Bấm keys (tên phím robotgo: enter, tab, esc, f5, a…) kèm modifiers (cmd, ctrl, alt, shift). keyMode chord: nhấn giữ cùng lúc; sequence: giữ phím kèm, gõ lần lượt.",
		[]string{"keys", "modifiers", "keyMode"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"key","keys":["l"],"modifiers":["cmd"],"keyMode":"chord","delayAfter":300,"enabled":true}`, []string{"keys", "window", "windowApp"}, nil},
	{StepScroll, "Cuộn", "Cuộn amount nấc theo direction (up | down | left | right) tại points[0]; không có điểm thì cuộn ở chỗ chuột đang đứng.",
		[]string{"direction", "amount", "points[0..1]", "jitterPx"}, []string{"amount", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"scroll","direction":"down","amount":5,"delayAfter":500,"enabled":true}`, []string{"x", "y", "direction", "amount"}, nil},
	{StepOpen, "Mở URL / app", "Mở target: URL (https://… hoặc tên miền như youtube.com/…), đường dẫn file / app, hoặc tên app (Spotify). app: mở bằng app này, vd Google Chrome; trống = mặc định của máy. Nhận {{biến}} (trong URL tự mã hoá). Không chờ trang tải xong – đặt delayAfter đủ lâu hoặc thêm bước chờ; trong lúc chờ delayAfter dò cửa sổ mới hiện: output window / windowApp của nó + newWindow true, không thấy (vd mở tab mới trong Chrome đang mở) → cửa sổ đang ở trước + newWindow false. webBrowser true: mở URL / file bằng Chrome của app (profile riêng của app, Chrome / Edge / Brave – lần đầu phải đăng nhập lại trang) để node web (webWait, webClick, webType, webRead, webScroll) điều khiển được; đang mở thì mở tab mới, bỏ qua app.",
		[]string{"target", "app", "webBrowser"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"open","target":"https://www.youtube.com/results?search_query={{bai_hat}}","app":"Google Chrome","delayAfter":3000,"enabled":true}`, []string{"target", "app", "window", "windowApp", "newWindow"}, nil},
	{StepWaitImage, "Chờ ảnh", "Chờ tới khi ảnh mẫu template (PNG base64, chụp bằng app – CaptureTemplate) hiện trên màn hình, độ khớp ≥ threshold %. Hết timeoutMs không thấy → cổng fail (error.code not_found). Tìm trong vùng của kịch bản, hoặc searchArea \"screen\" = cả màn chứa vùng. Output x, y (tâm), w, h, score – node clickImage chọn imageFrom = id node này để click vào đó. matchMode \"edge\" = so theo viền thay vì độ sáng: ảnh mẫu khoanh ở theme sáng vẫn tìm được ở theme tối / đổi màu nền (kém hơn \"\" khi ảnh ít nét). colorTol (1–100) = sau khi khớp hình, so màu trung bình chỗ thấy với ảnh mẫu, lệch quá colorTol % thì coi như không thấy – phân biệt nút cùng hình khác màu (xám / xanh); 0 = không kiểm tra màu. Ảnh mẫu PNG có điểm trong suốt (alpha) = các điểm đó không tính khi so (số đếm, avatar, chữ hay đổi). findBy \"text\": tìm chữ findText thay ảnh mẫu (đọc chữ trên màn hình, tiếng Việt / Anh; không phân biệt hoa thường và khoảng trắng; nhận {{biến}}); textExact = cả dòng đúng chuỗi, textAccent = so đúng dấu (mặc định bỏ qua dấu), occurrence = chỗ thứ mấy theo thứ tự đọc. Chữ bền hơn ảnh khi đổi theme / zoom – nên dùng khi nút có chữ. findBy \"element\" (macOS): tìm phần tử qua cây accessibility của app elementApp (tên app / bundle id, trống = app có cửa sổ trên cùng) – Chrome thường đã đăng nhập, app khác; selectors như node web nhưng chỉ text= / aria= / role=…[name=]; đọc được chữ chính xác, nội dung ô nhập, nút chỉ có icon; cửa sổ phải đang hiện trên màn (thu nhỏ / Space khác → chờ tới hết giờ); phần tử ngoài khung nhìn thì tự cuộn tới; textAccent, occurrence dùng như tìm chữ; searchArea chỉ lọc khi kịch bản có vùng; output thêm text (chữ / giá trị) và role.",
		[]string{"template", "templateScale", "threshold", "matchMode", "colorTol", "timeoutMs", "intervalMs", "searchArea", "findBy", "findText", "textExact", "textAccent", "occurrence", "selectors", "elementApp"}, []string{"timeoutMs", "threshold", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"waitImage","template":"iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAAAAADhZOFXAAAAHUlEQVR4nGNkYPjPwMjwn4GJAQb+/4dguAgjEWoAyH0MAWBh3h0AAAAASUVORK5CYII=","templateScale":1,"threshold":90,"timeoutMs":10000,"intervalMs":250,"delayAfter":300,"enabled":true}`, []string{"x", "y", "w", "h", "score", "text", "role"}, nil},
	{StepClickImage, "Click theo ảnh", "Click vào tâm chỗ thấy ảnh / chữ (+ dx, dy). imageSource \"node\": chỗ node imageFrom (id node chờ ảnh / chữ, hoặc click theo ảnh tự tìm) thấy lần gần nhất trong vòng – node đó chưa thấy thì lỗi; \"search\": tự tìm ảnh mẫu / chữ / phần tử riêng (các trường như waitImage, gồm findBy \"text\" / \"element\" – phần tử: chờ khung đứng yên rồi mới bấm). Không thấy → cổng fail.",
		[]string{"imageSource", "dx", "dy", "button", "template", "templateScale", "threshold", "matchMode", "colorTol", "timeoutMs", "intervalMs", "searchArea", "findBy", "findText", "textExact", "textAccent", "occurrence", "selectors", "elementApp", "jitterPx"}, []string{"timeoutMs", "threshold", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","ref":"n2","type":"clickImage","imageSource":"node","imageFrom":"s1","dx":0,"dy":0,"button":"left","delayAfter":500,"enabled":true}`, []string{"clickX", "clickY", "x", "y", "w", "h", "score", "text", "role"}, nil},
	{StepWaitPixel, "Chờ màu pixel", "Chờ tới khi điểm points[0] có màu color (#RRGGBB), mỗi kênh lệch ≤ tolerance. Hết timeoutMs → cổng fail. Không di chuột. Output color = màu đọc được.",
		[]string{"points[1]", "color", "tolerance", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"waitPixel","points":[{"x":640,"y":300}],"color":"#1877F2","tolerance":10,"timeoutMs":5000,"intervalMs":100,"delayAfter":300,"enabled":true}`, []string{"color"}, nil},
	{StepScreenshot, "Chụp màn hình", "Chụp màn hình rồi lưu PNG. shotArea: trống = vùng của kịch bản (không có vùng: màn chính); \"screen\" = cả một màn (display: 0 = màn chứa vùng kịch bản, n = màn thứ n theo danh sách màn); \"rect\" = vùng cố định shotRect {x, y, w, h} toạ độ màn hình; \"node\" = khung chỗ thấy của node imageFrom (chờ ảnh / chữ, click theo ảnh tự tìm, if / switchScreen, chờ / click phần tử web không chạy nền) lần gần nhất trong vòng. saveDir: thư mục lưu (trống = thư mục ảnh người dùng chọn trong Cài đặt của app, mặc định Pictures/Ghosthand – nên để trống trừ khi người dùng muốn chỗ riêng; ~ = thư mục người dùng, tạo nếu chưa có); fileName: tên file (trống = ghosthand-<ngày>-<giờ>, tự thêm .png, trùng thì thêm số) – cả hai nhận {{biến}}. toClipboard true: chép ảnh vào clipboard. Cần quyền ghi màn hình (permission). Không di chuột. Output path (đường dẫn file), x, y, w, h (vùng đã chụp, toạ độ màn hình).",
		[]string{"shotArea", "display", "shotRect", "imageFrom", "saveDir", "fileName", "toClipboard"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","ref":"n2","type":"screenshot","shotArea":"node","imageFrom":"s1","fileName":"don-hang-{{ma_don}}","toClipboard":true,"delayAfter":300,"enabled":true}`, []string{"path", "x", "y", "w", "h"}, nil},
	{StepClipboard, "Clipboard", "Đọc / ghi clipboard của máy. clipOp: trống = đọc chữ đang có → output text (clipboard trống / chỉ có ảnh = \"\"); \"write\" = ghi text (nhận {{biến}}, rỗng = xoá clipboard); \"image\" = ghi ảnh PNG từ filePath (nhận {{biến}}, vd {{n2.output.path}} của node screenshot; ~ / đường dẫn tương đối tính từ thư mục người dùng). Copy / paste trọn gói: node key Cmd/Ctrl+C rồi clipboard đọc; clipboard write rồi key Cmd/Ctrl+V. Output text (read / write) hoặc path (image).",
		[]string{"clipOp", "text", "filePath"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s3","ref":"n3","type":"clipboard","clipOp":"write","text":"Mã đơn: {{n2.output.text}}","delayAfter":200,"enabled":true}`, []string{"text", "path"}, nil},
	{StepHTTP, "HTTP", "Gọi API. method GET | POST | PUT | PATCH | DELETE (trống = GET); url http(s)://… – biến chèn nguyên văn ({{base}}/users được), giá trị có dấu cách / tiếng Việt đặt trong query [{key, value}] (tự mã hoá); headers [{key, value}]; bodyType trống (không body) | \"json\" (body là JSON, Content-Type application/json; biến trong body thoát theo chuỗi JSON nên đặt trong dấu nháy: \"{{ten}}\") | \"form\" (form [{key, value}], x-www-form-urlencoded) | \"text\" (body nguyên văn). Mọi ô chữ nhận {{biến}}. timeoutMs: chờ phản hồi tối đa (mặc định 15000). Mã 2xx → cổng success; mã khác → fail error.code http_status (output vẫn có – nhánh fail đọc {{n2.output.status}}); không gửi được / hết giờ → fail network. Output status, body (chữ, tối đa 64 KB), headers (tên chữ thường, set-cookie che), json (phản hồi JSON đã giải mã – đọc sâu {{n2.output.json.data.items.0.id}}, mảng đếm từ 0). Header / tham số có tên kiểu token, key, auth, password không hiện trong lời báo lỗi.",
		[]string{"method", "url", "query", "headers", "bodyType", "body", "form", "timeoutMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","ref":"n2","type":"http","method":"POST","url":"https://api.example.com/orders","headers":[{"key":"Authorization","value":"Bearer {{token}}"}],"bodyType":"json","body":"{\"ma\": \"{{n1.output.text}}\", \"so_luong\": 2}","timeoutMs":15000,"enabled":true}`, []string{"status", "body", "headers", "json"}, nil},
	{StepFile, "Tệp / thư mục", "Thao tác file theo fileOp. Trống: ghi text vào filePath (ghi đè; fileAppend true = thêm vào cuối; tạo thư mục cha nếu thiếu). \"row\": thêm một dòng cells vào file .csv (chưa có thì tạo, có BOM UTF-8 cho Excel) hoặc .xlsx (sheet: tên sheet, trống = sheet đầu; ô trông như số ghi thành số). \"mkdir\": tạo thư mục filePath. \"trash\": chuyển các file (không gồm thư mục con, file ẩn chỉ khi mẫu bắt đầu bằng .) trong thư mục filePath khớp filePattern (vd *.png, bao-cao-*.csv) vào Thùng rác – khôi phục được; không file nào khớp vẫn thành công. \"emptyTrash\": dọn Thùng rác (xoá hẳn). trash / emptyTrash là thao tác xoá: người dùng phải đồng ý trong app trước lần chạy đầu (và khi thư mục / mẫu tên đổi), chưa đồng ý → fail permission. filePath: ~ / đường dẫn tương đối tính từ thư mục người dùng; thư mục gốc, chính thư mục người dùng, ~/Library, thư mục hệ thống bị chặn (invalid_input). filePath, text, cells, sheet, filePattern nhận {{biến}}. Output path (+ bytes khi ghi, rows = số dòng sau khi thêm, count / files khi xoá).",
		[]string{"fileOp", "filePath", "text", "fileAppend", "cells", "sheet", "filePattern"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s3","ref":"n3","type":"file","fileOp":"row","filePath":"~/Documents/don-hang.xlsx","cells":["{{n1.output.text}}","{{n2.output.json.total}}"],"enabled":true}`, []string{"path", "bytes", "rows", "count", "files"}, nil},
	{StepSystem, "Hệ thống", "Âm lượng loa của máy theo sysOp: \"volume\" = đặt mức volume (0–100 %); \"volumeUp\" / \"volumeDown\" = tăng / giảm volume % (0 = 10, giữ trong 0–100); \"mute\" / \"unmute\" = tắt / bật tiếng. Đặt / tăng / giảm cũng bật tiếng lại nếu đang tắt. Không đụng chuột / phím. Output volume (mức sau khi chạy), muted (true / false).",
		[]string{"sysOp", "volume"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"system","sysOp":"volume","volume":30,"delayAfter":300,"enabled":true}`, []string{"volume", "muted"}, nil},
	{StepWindow, "Cửa sổ", "Cố định cửa sổ trước khi thao tác. Chọn cửa sổ: winApp (tên app / tên exe như chrome.exe / bundle id) và / hoặc winTitle (tiêu đề) – không được cả hai trống; winMatch trống = chứa (không phân biệt hoa thường, không dấu) | \"exact\" | \"regex\". Nhiều cửa sổ khớp → cửa sổ trên cùng; bỏ qua cửa sổ của ghosthand. winOp: trống = focus (đưa lên trước, bỏ thu nhỏ); \"waitOpen\" / \"waitClose\" = chờ cửa sổ hiện / mất (hết timeoutMs → fail timeout); \"setBounds\" = đặt khung: winPreset trống = đúng winRect {x, y, w, h} toạ độ màn hình, \"left\" / \"right\" = cỡ winRect.w × h ở góc trên trái / phải (0 = nửa màn), \"center\" = giữa màn, \"full\" = cả vùng làm việc (không đè taskbar / thanh menu) – trên màn display (0 = màn đang chứa cửa sổ, n = màn thứ n); \"minimize\" / \"maximize\" / \"restore\"; \"close\" = như bấm × (app có thể hỏi lưu; thao tác nguy hiểm – người dùng phải đồng ý trong app trước lần chạy đầu, chưa đồng ý → fail permission); \"read\" = chỉ đọc thông tin. Mọi thao tác chờ cửa sổ hiện tới timeoutMs (hết giờ → fail not_found). Đặt cửa sổ bằng setBounds trước khi click toạ độ cứng / tìm ảnh để lần chạy nào cửa sổ cũng ở đúng chỗ, đúng cỡ. Output title, app, pid, x, y (tâm khung), w, h, left, top (góc trên trái), display (màn thứ mấy, từ 1), minimized – khung dùng được cho imageFrom của node Chụp màn hình / Click theo ảnh.",
		[]string{"winOp", "winApp", "winTitle", "winMatch", "winPreset", "winRect", "display", "timeoutMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","ref":"n2","type":"window","winOp":"setBounds","winApp":"chrome.exe","winTitle":"Sân thử","winPreset":"left","winRect":{"x":0,"y":0,"w":1280,"h":800},"timeoutMs":10000,"delayAfter":300,"enabled":true}`, []string{"title", "app", "pid", "x", "y", "w", "h", "left", "top", "display", "minimized"}, nil},
	{StepApp, "Ứng dụng", "Kiểm tra / chờ / thoát một app đang chạy. appName: tên app / tên exe (notepad.exe) / bundle id / đường dẫn exe, nhận {{biến}}; appMatch như winMatch (trống = chứa, không phân biệt hoa thường, không dấu | \"exact\" | \"regex\"). App nhiều tiến trình (Chrome) tính là một. appOp: trống = kiểm tra đang chạy → cổng \"yes\" / \"no\" (không có cổng success / fail); \"waitStart\" = chờ app chạy và có ít nhất một cửa sổ (hết timeoutMs → fail timeout); \"waitQuit\" = chờ app thoát hẳn; \"quit\" = thoát bình thường như đóng mọi cửa sổ (app có thể hỏi lưu – không tự bấm), quá timeoutMs chưa thoát → fail timeout, hoặc appForce true thì buộc thoát; \"forceQuit\" = buộc thoát ngay (mất dữ liệu chưa lưu). forceQuit và quit + appForce là thao tác nguy hiểm: người dùng phải đồng ý trong app trước lần chạy đầu, chưa đồng ý → fail permission. App không chạy thì quit / forceQuit thành công ngay. Không thoát được ghosthand, Finder / Dock, explorer.exe, tiến trình hệ thống (invalid_input). Node này không mở app – dùng node Mở. Mẫu mở đầu: Ứng dụng (kiểm tra) → no → Mở → Ứng dụng waitStart → Cửa sổ setBounds. Output running (true / false), pid, name, id (đường dẫn exe / bundle id), windows (số cửa sổ).",
		[]string{"appOp", "appName", "appMatch", "appForce", "timeoutMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"app","appOp":"waitStart","appName":"notepad.exe","timeoutMs":10000,"delayAfter":300,"enabled":true}`, []string{"running", "pid", "name", "id", "windows"}, nil},
	{StepWaitUntil, "Chờ đến giờ", "Chờ tới đúng giờ at rồi đi tiếp: \"HH:mm\", \"HH:mm:ss\" hoặc \"HH:mm:ss.mmm\" (hôm nay, giờ 24h), có thể kèm ngày phía trước \"2026-11-11 20:00\" / \"11/11/2026 20:00:00.500\". Nhận {{biến}}. ifPassed – giờ đã qua lúc tới node: trống = đi tiếp ngay; \"tomorrow\" = chờ giờ đó ngày mai (giờ không ghi ngày, hợp với vòng lặp hằng ngày); \"fail\" = cổng fail (error.code timeout). useNtp true: canh theo giờ chuẩn internet (NTP) – đồng hồ máy lệch vẫn đúng giờ (săn sale), không lấy được giờ chuẩn thì dùng giờ máy. Chờ dài nghỉ từng đoạn (Dừng / tạm dừng ăn ngay), sát giờ chờ vòng ngắn – lệch vài ms. Output firedAt (giờ chạy tiếp, ISO có ms), target (giờ hẹn), lateMs (trễ so với giờ hẹn, ms), ntpOffsetMs (giờ chuẩn − giờ máy, 0 khi không dùng NTP).",
		[]string{"at", "ifPassed", "useNtp"}, []string{"delayAfter", "enabled"},
		`{"id":"s1","ref":"n1","type":"waitUntil","at":"20:00:00.000","ifPassed":"","useNtp":true,"enabled":true}`, []string{"firedAt", "target", "lateMs", "ntpOffsetMs"}, nil},
	{StepPath, "Vẽ đường", "Giữ nút chuột (button) rồi kéo một nét liền trong dragMs ms. shape trống: đi lần lượt qua points (≥ 2 điểm, neo / biến như click; closed true = nối điểm cuối về đầu). shape \"circle\" / \"rect\" / \"star\": hình quanh tâm points[0], cỡ shapeW × shapeH px (shapeH 0 = bằng shapeW; tròn bắt đầu ở đỉnh, đi chiều kim đồng hồ). Thời gian chia theo độ dài từng đoạn; giống người thì mỗi đoạn hơi cong + nhịp tăng / giảm tốc. Luôn nhả chuột khi xong / bị dừng. Dùng cho app vẽ, chữ ký, mở khoá hình vẽ. Output x, y (chỗ nhả), points (số đỉnh đã đi qua).",
		[]string{"points", "shape", "shapeW", "shapeH", "closed", "button", "dragMs", "jitterPx"}, []string{"dragMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"path","shape":"star","points":[{"x":0,"y":0,"ax":"center","ay":"center"}],"shapeW":300,"button":"left","dragMs":2000,"enabled":true}`, []string{"x", "y", "points"}, nil},
	{StepWebWait, "Chờ phần tử web", "Chờ tới khi phần tử trên trang web hiện (webGone: mất) – đọc HTML qua Chrome của app (mở trang bằng node open có webBrowser true trước). selectors: các cách tìm theo thứ tự ưu tiên, cách đầu thấy thì dùng: \"text=Xem lịch sử\" (chữ hiển thị chứa chuỗi, không phân biệt hoa thường / khoảng trắng; text=\"Gửi\" = đúng cả chuỗi; không tiền tố = text), \"aria=Đóng\" (aria-label), \"role=button[name=\\\"Gửi\\\"]\" (vai trò thẻ / thuộc tính role, name = tên chứa chuỗi; role=link không name = mọi link), \"label=Ngày sinh\" (ô nhập của nhãn: label for / aria-labelledby, không có thì ô đầu tiên sau nhãn trong cùng khối form – nhãn antd; label=\"…\" = nhãn đúng cả chuỗi), \"css=#id\", \"xpath=//button\". Ưu tiên chữ / aria / role / label, tránh class sinh tự động. webScope \"dialog\": chỉ tìm trong hộp thoại đang mở trên cùng (bỏ ô cùng tên ở trang phía sau); occurrence N: chỗ thứ N (thứ tự trong trang) trong số phần tử đang hiện khớp selector (trống = chỗ bấm được đầu tiên). tab: id tab (tab / tabs của probe_page, read_web_outline – chắc chắn đúng tab khi nhiều tab cùng URL, nhưng Chrome của app mở lại thì id đổi: chỉ dùng khi chạy thử trong phiên) hoặc chuỗi có trong URL / tiêu đề tab (bền khi lưu kịch bản); trống = tab đang hiện. Không thấy sau timeoutMs → cổng fail (not_found); Chrome của app chưa mở → fail (browser). Không bấm gì, nhưng lần đầu (và khi cửa sổ dời / đổi cỡ / zoom) di chuột vào trang để đo toạ độ. Output x, y (tâm, toạ độ màn hình), w, h, text, selector (cách đã khớp). webBackground true: chạy nền – không đụng chuột / phím thật, không đưa Chrome lên trước (cửa sổ bị che / thu nhỏ vẫn chạy), không đo toạ độ bằng chuột qua CDP; output x, y theo toạ độ trong trang (CSS px) thay vì màn hình.",
		[]string{"selectors", "tab", "webScope", "occurrence", "webGone", "webBackground", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"webWait","selectors":["text=Xem lịch sử","aria=Lịch sử"],"timeoutMs":10000,"intervalMs":200,"delayAfter":300,"enabled":true}`, []string{"x", "y", "w", "h", "text", "selector"}, nil},
	{StepWebClick, "Click phần tử web", "Chờ phần tử hiện (như webWait) rồi click chuột thật vào tâm (+ dx, dy) clicks lần. Phần tử ngoài màn hình thì cuộn vào giữa trước; bị phần tử khác che thì chờ tiếp tới timeoutMs. Đưa tab lên trước khi bấm. webBackground true: chạy nền – không đụng chuột / phím thật, không đưa Chrome lên trước (cửa sổ bị che / thu nhỏ vẫn chạy), bấm qua CDP; output x, y theo toạ độ trong trang (CSS px) thay vì màn hình. Trang nặng chống bot nên dùng chuột thật.",
		[]string{"selectors", "tab", "webScope", "occurrence", "webBackground", "button", "clicks", "clickInterval", "dx", "dy", "timeoutMs", "intervalMs"}, []string{"clicks", "clickInterval", "timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","ref":"n2","type":"webClick","selectors":["text=\"Gửi\"","css=button[type=submit]"],"button":"left","clicks":1,"timeoutMs":10000,"intervalMs":200,"delayAfter":500,"enabled":true}`, []string{"clickX", "clickY", "x", "y", "w", "h", "text", "selector"}, nil},
	{StepWebType, "Gõ vào ô web", "Tìm ô nhập như webClick, click chuột thật vào ô, kiểm tra ô đã nhận focus (không thì fail not_found, không gõ; skipFocusCheck true = bỏ kiểm tra – ô tuỳ biến focus vào phần tử khác) rồi gõ text bằng bàn phím thật (nhận {{biến}}). webClear true (mặc định): xoá chữ cũ trước (Cmd/Ctrl+A rồi Backspace); webEnter true: bấm Enter sau khi gõ; typeMode như node type (trống = Windows dán, còn lại gõ phím). Selector nên trỏ đúng ô: aria=, role=textbox[name=\"Email\"], css=input[name=email], hoặc text= theo placeholder / nhãn. webBackground true: chạy nền – không đụng chuột / phím thật, không đưa Chrome lên trước (cửa sổ bị che / thu nhỏ vẫn chạy), click + chèn chữ (Input.insertText – cả chuỗi một lần, không gõ từng phím) qua CDP; output x, y theo toạ độ trong trang (CSS px) thay vì màn hình.",
		[]string{"selectors", "tab", "webScope", "occurrence", "webBackground", "text", "typeMode", "webClear", "webEnter", "skipFocusCheck", "dx", "dy", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s3","ref":"n3","type":"webType","selectors":["role=textbox[name=\"Email\"]","css=input[type=email]"],"text":"{{email}}","webClear":true,"timeoutMs":10000,"intervalMs":200,"delayAfter":300,"enabled":true}`, []string{"text", "clickX", "clickY", "x", "y", "w", "h", "selector"}, nil},
	{StepWebRead, "Đọc phần tử web", "Chờ phần tử hiện (như webWait, không di chuột) rồi đọc giá trị vào output text. readFrom: trống = chữ hiển thị (innerText), \"value\" = giá trị ô nhập / select, \"attr\" = thuộc tính tên attr (href / src trả URL đầy đủ). count = số phần tử đang hiện khớp cách tìm đã dùng. Node sau dùng {{n4.output.text}} (vd so trong node if / switch). Không thấy sau timeoutMs → fail not_found.",
		[]string{"selectors", "tab", "webScope", "occurrence", "readFrom", "attr", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s4","ref":"n4","type":"webRead","selectors":["css=#order-status"],"timeoutMs":5000,"intervalMs":200,"enabled":true}`, []string{"text", "count", "selector"}, nil},
	{StepWebScroll, "Cuộn trang web", "Cuộn theo scrollMode. Trống: tới phần tử selectors – đưa vào giữa khung nhìn (không di chuột). \"in\": chuột thật tới tâm phần tử (khung có thanh cuộn riêng) rồi lăn amount nấc theo direction (up | down | left | right). \"until\": lăn amount nấc mỗi lần (mỗi intervalMs) cho tới khi phần tử selectors nằm trong khung nhìn – feed / danh sách tải thêm khi cuộn; lăn trong khung scrollIn (cách tìm như selectors, trống = giữa trang); hết timeoutMs → fail not_found. Lăn bằng bánh xe chuột thật. webBackground true: chạy nền – không đụng chuột / phím thật, không đưa Chrome lên trước (cửa sổ bị che / thu nhỏ vẫn chạy), lăn (mỗi nấc 100 px) qua CDP; output x, y theo toạ độ trong trang (CSS px) thay vì màn hình. Output x, y, w, h, text, selector (until: của phần tử đích, thêm scrolls = số lần đã lăn).",
		[]string{"selectors", "tab", "webScope", "occurrence", "webBackground", "scrollMode", "direction", "amount", "scrollIn", "timeoutMs", "intervalMs"}, []string{"amount", "timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s5","ref":"n5","type":"webScroll","scrollMode":"until","selectors":["text=Bài viết 30"],"direction":"down","amount":3,"timeoutMs":30000,"intervalMs":400,"delayAfter":300,"enabled":true}`, []string{"x", "y", "w", "h", "text", "selector", "scrolls"}, nil},
	{StepWebKey, "Phím trên trang web", "Bấm keys (+ modifiers, keyMode như node key) vào tab của Chrome của app qua CDP – phím vào ô đang focus của trang; không cần Chrome ở trước, không đụng bàn phím thật (khác node key: key gõ vào cửa sổ đang ở trước, dễ trúng app khác). Dùng sau webClick / webType: esc đóng hộp thoại, tab sang ô kế, enter gửi, ctrl+a chọn hết. Phím: enter, tab, esc, backspace, delete, space, up / down / left / right, home, end, pageup, pagedown, f1–f12, một chữ / số / dấu. tab: như node web (id / chuỗi URL; trống = tab node web đang dùng). Không có tab / Chrome chưa mở → fail (browser). Output keys, tab (URL tab đã bấm).",
		[]string{"keys", "modifiers", "keyMode", "tab"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s6","ref":"n6","type":"webKey","keys":["esc"],"delayAfter":300,"enabled":true}`, []string{"keys", "tab"}, nil},
	{StepElementRead, "Đọc phần tử", "macOS: chờ tới khi có phần tử selectors (text= / aria= / role=…[name=]) trong cửa sổ đang hiện của app elementApp (tên app / bundle id, trống = app có cửa sổ trên cùng) – qua cây accessibility, dùng được với Chrome thường đã đăng nhập (không cần Chrome của app) và app khác – rồi đọc vào output text: khối chữ / ô nhập / select → giá trị (chữ chính xác, không qua OCR), phần tử khác → tên (nút chỉ có icon: nhãn). Không di chuột, không cuộn (phần tử ngoài khung nhìn vẫn đọc được). occurrence = chỗ thứ mấy theo thứ tự đọc, textAccent = so đúng dấu (mặc định bỏ dấu khi so); count = số phần tử khớp. Cửa sổ thu nhỏ / ở Space khác → chờ tới hết timeoutMs → fail not_found; thiếu quyền Accessibility → permission. Node sau dùng {{n6.output.text}}. Dự phòng khi giá trị không có trong cây (canvas, app tự vẽ – cần quyền Screen Recording): readNear \"right\" | \"below\" = tìm nhãn rồi OCR chữ bên phải / bên dưới nhãn (dòng gần nhất; readSpan = khoảng đọc px từ mép nhãn, 0 = 300); findBy \"text\" = tìm nhãn bằng OCR (findText, textExact, textAccent, occurrence, searchArea – như waitImage tìm theo chữ; bắt buộc có readNear) thay vì selectors. Lọc: readPattern = regex (có nhóm () thì lấy nhóm đầu), readNumber true = chuẩn hoá số (1.234.000₫ → 1234000, $1,234.56 → 1234.56); chưa khớp mẫu / không có số → chờ đọc lại. readChange true = chờ giá trị đổi: lần đọc đầu làm mốc, xong khi giá trị (sau lọc) khác mốc (output old = mốc); ảnh vùng không đổi thì không OCR lại. Output raw = chữ trước khi lọc.",
		[]string{"selectors", "elementApp", "textAccent", "occurrence", "searchArea", "findBy", "findText", "textExact", "readNear", "readSpan", "readPattern", "readNumber", "readChange", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s6","ref":"n6","type":"elementRead","selectors":["text=₫"],"elementApp":"Google Chrome","readNumber":true,"timeoutMs":10000,"intervalMs":250,"enabled":true}`, []string{"text", "raw", "old", "count", "role", "x", "y", "w", "h"}, nil},
	{StepIf, "Nếu", "Rẽ nhánh theo điều kiện cond (xem flow guide): đúng → cổng true, sai → false. timeoutMs > 0: xem lại mỗi intervalMs tới khi đúng, hết giờ = sai (0 = xem một lần). Thiếu quyền / ảnh hỏng → node lỗi, dừng kịch bản. Output result (true / false) + chỗ thấy x, y, w, h, score, text của điều kiện ảnh / chữ đúng đầu tiên (clickImage chọn imageFrom = node này được).",
		[]string{"cond", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"if","cond":{"kind":"any","conds":[{"kind":"text","findText":"Đã gửi"},{"kind":"value","left":"{{n2.output.text}}","op":"contains","right":"ok"}]},"timeoutMs":3000,"intervalMs":250,"enabled":true}`,
		[]string{"result", "x", "y", "w", "h", "score", "text", "role"}, nil},
	{StepSwitch, "Switch theo giá trị", "So value (thường là {{biến}}) với từng case theo thứ tự: case {key, label, op: eq | contains (không phân biệt hoa thường) | regex, match}; case đầu khớp → cổng case:<key>, không case nào khớp → default.",
		[]string{"value", "cases"}, []string{"enabled"},
		`{"id":"s1","ref":"n1","type":"switch","value":"{{n2.output.text}}","cases":[{"key":"ok","label":"Thành công","op":"contains","match":"đã gửi"},{"key":"err","label":"Lỗi","op":"regex","match":"(?i)lỗi|error"}],"enabled":true}`,
		[]string{"value", "case"}, nil},
	{StepSwitchScreen, "Switch theo màn hình", "Chụp màn hình mỗi intervalMs, xem case nào (cases[].cond: thấy chữ / ảnh / màu) đúng trước – xét theo thứ tự case – thì đi cổng case:<key>; hết timeoutMs không case nào → cổng none. Dùng khi trang có thể ở nhiều trạng thái (đăng nhập / đã vào / popup).",
		[]string{"cases", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "enabled"},
		`{"id":"s1","ref":"n1","type":"switchScreen","cases":[{"key":"login","label":"Đăng nhập","cond":{"kind":"text","findText":"Đăng nhập"}},{"key":"home","label":"Trang chủ","cond":{"kind":"text","findText":"Nhắn tin"}}],"timeoutMs":3000,"intervalMs":250,"enabled":true}`,
		[]string{"case", "x", "y", "w", "h", "score", "text", "role"}, nil},
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

// CommonOutputs – output mọi node có khi thành công / bị bỏ qua (engine điền trong execNode): startedAt (giờ bắt đầu,
// RFC 3339), durationMs (thời gian chạy, kể cả thử lại, chưa tính delayAfter).
var CommonOutputs = []string{"startedAt", "durationMs"}

// OutputsOf – khoá output của loại node t, gồm CommonOutputs (nil = loại không có trong bảng).
func OutputsOf(t string) []string {
	for _, d := range stepTypeDocs {
		if d.Type == t {
			return append(slices.Clone(d.Outputs), CommonOutputs...)
		}
	}
	return nil
}

// NumVarFieldsOf – trường số của loại node t nhận biến (numVars, F5c): các trường ghim được trừ enabled.
func NumVarFieldsOf(t string) []string {
	for _, d := range stepTypeDocs {
		if d.Type == t {
			out := []string{}
			for _, f := range d.Pinnable {
				if f != "enabled" {
					out = append(out, f)
				}
			}
			return out
		}
	}
	return nil
}

// Guide – hướng dẫn cấu trúc flow, trường chung và kịch bản mẫu (cho AI / MCP, bước 31).
func Guide() FlowGuideDoc {
	return FlowGuideDoc{Flow: flowDoc, CommonFields: commonFieldsDoc, Example: flowExample}
}
