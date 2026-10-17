package main

// Mô tả các loại bước cho người đọc file kịch bản và cho AI (MCP, bước 31 – sinh schema / hướng dẫn từ bảng này).
// Thêm loại bước mới: thêm hằng trong model.go, xử lý trong checkStep / validateScript / engine, và một dòng ở đây
// (test kiểm tra mọi loại đều có mô tả và ví dụ hợp lệ).

type StepTypeDoc struct {
	Type     string   `json:"type"`
	Label    string   `json:"label"`    // tên tiếng Việt trên UI
	Desc     string   `json:"desc"`     // làm gì, khi nào dùng
	Fields   []string `json:"fields"`   // trường riêng của loại này (ngoài id, type, enabled, delayAfter, jitterMs)
	Pinnable []string `json:"pinnable"` // trường ghim được làm input (inputs.go)
	Example  string   `json:"example"`  // một bước JSON hợp lệ
}

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
}

// StepTypes trả mô tả các loại bước app đang hỗ trợ.
func (a *App) StepTypes() []StepTypeDoc { return stepTypeDocs }
