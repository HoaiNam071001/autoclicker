package main

// Mô tả các loại bước cho người đọc file kịch bản và cho AI (MCP, bước 31 – sinh schema / hướng dẫn từ bảng này).
// Thêm loại bước mới: thêm hằng trong model.go, xử lý trong checkStep / validateScript / engine, và một dòng ở đây
// (test kiểm tra mọi loại đều có mô tả và ví dụ hợp lệ).

type StepTypeDoc struct {
	Type     string   `json:"type"`
	Label    string   `json:"label"`    // tên tiếng Việt trên UI
	Desc     string   `json:"desc"`     // làm gì, khi nào dùng
	Fields   []string `json:"fields"`   // trường riêng của loại này (ngoài trường chung, xem commonFieldsDoc)
	Pinnable []string `json:"pinnable"` // trường ghim được làm input (inputs.go)
	Example  string   `json:"example"`  // một bước JSON hợp lệ
}

// commonFieldsDoc – trường mọi loại bước đều có.
const commonFieldsDoc = `id (ổn định, để sửa / nhảy tới), type, enabled, delayAfter (ms chờ sau bước), jitterMs (± ms ngẫu nhiên). ` +
	`onFail {retries, action, goto}: bước thất bại (chờ ảnh / chữ / màu hết giờ, lỗi) thì thử lại thêm retries lần (0–100), ` +
	`vẫn không được thì action: "stop" dừng kịch bản (mặc định khi không có onFail), "next" bỏ qua sang bước sau, ` +
	`"goto" nhảy tới bước có id = goto (trước hoặc sau; tối đa 1000 lần nhảy mỗi vòng). ` +
	`Vd {"retries":2,"action":"goto","goto":"s5"}.`

// condDoc – mô tả Cond (điều kiện của khối if / group).
const condDoc = `cond.kind: "image" (template, threshold…), "text" (findText, textExact, textNoAccent, occurrence), "pixel" (points[0], color, tolerance) – ` +
	`cùng trường với waitImage / waitPixel, waitMs = chờ tối đa trước khi kết luận không thấy (0 = xem một lần), thấy thì clickImage prev click được; ` +
	`"input": so left op right (op eq | ne | contains | gt | lt | empty; left / right nhận {{biến}}, vd left "{{ten}}"); ` +
	`"result": bước chạy ngay trước thành công (dùng với onFail next); "all" / "any": conds là danh sách điều kiện con (và / hoặc). not: đảo kết quả.`

var stepTypeDocs = []StepTypeDoc{
	{StepClick, "Click", "Click chuột Clicks lần tại points[0] (neo theo cạnh vùng, xem anchor.go).",
		[]string{"points[1]", "button", "clicks", "clickInterval", "jitterPx"}, []string{"clicks", "clickInterval", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"click","points":[{"x":400,"y":300}],"button":"left","clicks":1,"delayAfter":500,"enabled":true}`},
	{StepDouble, "Double click", "Double click tại points[0].",
		[]string{"points[1]", "button", "jitterPx"}, []string{"jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"double","points":[{"x":400,"y":300}],"button":"left","delayAfter":500,"enabled":true}`},
	{StepHold, "Giữ chuột", "Nhấn giữ nút chuột holdMs mili giây tại points[0].",
		[]string{"points[1]", "button", "holdMs", "jitterPx"}, []string{"holdMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"hold","points":[{"x":400,"y":300}],"button":"left","holdMs":1000,"delayAfter":500,"enabled":true}`},
	{StepDrag, "Kéo thả", "Giữ chuột tại points[0] kéo tới points[1] trong dragMs. dragFromPrev: bắt đầu từ chỗ chuột dừng ở bước trước; dragRelative: points[1] là độ lệch so với điểm bắt đầu.",
		[]string{"points[2]", "button", "dragMs", "dragFromPrev", "dragRelative", "jitterPx"}, []string{"dragMs", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"drag","points":[{"x":100,"y":300},{"x":600,"y":300}],"button":"left","dragMs":400,"delayAfter":500,"enabled":true}`},
	{StepWait, "Chờ", "Chỉ chờ delayAfter mili giây.",
		[]string{}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"wait","delayAfter":2000,"enabled":true}`},
	{StepText, "Gõ chữ", "Gõ text vào ô đang focus (thường sau một bước click). Nhận {{biến}}; \\{{ để gõ ra {{.",
		[]string{"text"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"type","text":"xin chào {{ten}}","delayAfter":300,"enabled":true}`},
	{StepKey, "Phím tắt", "Bấm keys (tên phím robotgo: enter, tab, esc, f5, a…) kèm modifiers (cmd, ctrl, alt, shift). keyMode chord: nhấn giữ cùng lúc; sequence: giữ phím kèm, gõ lần lượt.",
		[]string{"keys", "modifiers", "keyMode"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"key","keys":["l"],"modifiers":["cmd"],"keyMode":"chord","delayAfter":300,"enabled":true}`},
	{StepScroll, "Cuộn", "Cuộn amount nấc theo direction (up | down | left | right) tại points[0]; không có điểm thì cuộn ở chỗ chuột đang đứng.",
		[]string{"direction", "amount", "points[0..1]", "jitterPx"}, []string{"amount", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"scroll","direction":"down","amount":5,"delayAfter":500,"enabled":true}`},
	{StepOpen, "Mở URL / app", "Mở target: URL (https://… hoặc tên miền như youtube.com/…), đường dẫn file / app, hoặc tên app (Spotify). app: mở bằng app này, vd Google Chrome; trống = mặc định của máy. Nhận {{biến}} (trong URL tự mã hoá). Không chờ trang tải xong – đặt delayAfter đủ lâu hoặc thêm bước chờ.",
		[]string{"target", "app"}, []string{"delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"open","target":"https://www.youtube.com/results?search_query={{bai_hat}}","app":"Google Chrome","delayAfter":3000,"enabled":true}`},
	{StepWaitImage, "Chờ ảnh", "Chờ tới khi ảnh mẫu template (PNG base64, chụp bằng app – CaptureTemplate) hiện trên màn hình, độ khớp ≥ threshold %. Hết timeoutMs không thấy thì dừng kịch bản. Tìm trong vùng của kịch bản, hoặc searchArea \"screen\" = cả màn chứa vùng. Chỗ thấy được nhớ cho bước clickImage sau. findBy \"text\": tìm chữ findText thay ảnh mẫu (đọc chữ trên màn hình, tiếng Việt / Anh; không phân biệt hoa thường và khoảng trắng; nhận {{biến}}); textExact = cả dòng đúng chuỗi, textNoAccent = bỏ qua dấu, occurrence = chỗ thứ mấy theo thứ tự đọc. Chữ bền hơn ảnh khi đổi theme / zoom – nên dùng khi nút có chữ.",
		[]string{"template", "templateScale", "threshold", "timeoutMs", "intervalMs", "searchArea", "findBy", "findText", "textExact", "textNoAccent", "occurrence"}, []string{"timeoutMs", "threshold", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"waitImage","template":"iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAAAAADhZOFXAAAAHUlEQVR4nGNkYPjPwMjwn4GJAQb+/4dguAgjEWoAyH0MAWBh3h0AAAAASUVORK5CYII=","templateScale":1,"threshold":90,"timeoutMs":10000,"intervalMs":250,"delayAfter":300,"enabled":true}`},
	{StepClickImage, "Click theo ảnh", "Click vào tâm chỗ thấy ảnh (+ dx, dy). imageSource prev: dùng kết quả bước chờ ảnh / click theo ảnh gần nhất phía trước; search: tự tìm ảnh mẫu / chữ riêng (các trường như waitImage, gồm findBy \"text\").",
		[]string{"imageSource", "dx", "dy", "button", "template", "templateScale", "threshold", "timeoutMs", "intervalMs", "searchArea", "findBy", "findText", "textExact", "textNoAccent", "occurrence", "jitterPx"}, []string{"timeoutMs", "threshold", "jitterPx", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s2","type":"clickImage","imageSource":"prev","dx":0,"dy":0,"button":"left","delayAfter":500,"enabled":true}`},
	{StepWaitPixel, "Chờ màu pixel", "Chờ tới khi điểm points[0] có màu color (#RRGGBB), mỗi kênh lệch ≤ tolerance. Hết timeoutMs thì dừng. Không di chuột.",
		[]string{"points[1]", "color", "tolerance", "timeoutMs", "intervalMs"}, []string{"timeoutMs", "delayAfter", "jitterMs", "enabled"},
		`{"id":"s1","type":"waitPixel","points":[{"x":640,"y":300}],"color":"#1877F2","tolerance":10,"timeoutMs":5000,"intervalMs":100,"delayAfter":300,"enabled":true}`},
	{StepIf, "Nếu", "Khối điều kiện: cond đúng → chạy children, sai → chạy else (else có thể chứa một if khác = \"nếu không, nếu\"). " + condDoc,
		[]string{"cond", "children", "else"}, []string{"delayAfter", "enabled"},
		`{"id":"b1","type":"if","cond":{"kind":"text","findText":"Nhắn tin","waitMs":3000,"intervalMs":250},"children":[{"id":"b2","type":"clickImage","imageSource":"prev","button":"left","delayAfter":500,"enabled":true}],"else":[{"id":"b3","type":"scroll","direction":"down","amount":3,"delayAfter":500,"enabled":true}],"delayAfter":0,"enabled":true}`},
	{StepGroup, "Nhóm / Lặp", "Khối chạy children: loopMode \"\" = repeat lần (≥ 1, 1 = chỉ gom nhóm, name là tên); \"until\" = lặp đến khi cond đúng; \"while\" = lặp khi cond còn đúng (cả hai xét cond trước mỗi vòng); maxLoops chặn lặp vô hạn (0 = 100, vượt → bước thất bại). Bước break bên trong thoát khối.",
		[]string{"name", "repeat", "loopMode", "cond", "maxLoops", "children"}, []string{"repeat", "delayAfter", "enabled"},
		`{"id":"g1","type":"group","name":"Cuộn tìm nút","loopMode":"until","cond":{"kind":"text","findText":"Xem lịch sử"},"maxLoops":10,"children":[{"id":"g2","type":"scroll","direction":"down","amount":3,"delayAfter":800,"enabled":true}],"delayAfter":0,"enabled":true}`},
	{StepBreak, "Thoát khối lặp", "Thoát khối group gần nhất chứa nó (thường đặt trong một if). Chỉ hợp lệ bên trong group.",
		[]string{}, []string{"enabled"},
		`{"id":"x1","type":"break","delayAfter":0,"enabled":true}`},
	{StepStop, "Dừng kịch bản", "Dừng cả kịch bản ngay (không chạy các vòng sau). stopError true = báo lỗi kèm message; false = kết thúc bình thường.",
		[]string{"message", "stopError"}, []string{"enabled"},
		`{"id":"x2","type":"stop","message":"Không tìm thấy người nhận","stopError":true,"delayAfter":0,"enabled":true}`},
}

// StepTypes trả mô tả các loại bước app đang hỗ trợ.
func (a *App) StepTypes() []StepTypeDoc { return stepTypeDocs }
