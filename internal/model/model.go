// Package model – kiểu dữ liệu kịch bản (node, cạnh, input, kết quả node) và công thức thuần đi kèm. Không phụ thuộc
// thiết bị / Wails; mọi package khác dùng chung.
package model

// Loại bước (Step.Type).
const (
	StepClick  = "click"  // click Clicks lần, cách nhau ClickInterval
	StepDouble = "double" // double click
	StepHold   = "hold"   // giữ chuột HoldMs
	StepDrag   = "drag"   // kéo Points[0] -> Points[1] trong DragMs
	StepWait   = "wait"   // chỉ chờ DelayAfter
	StepText   = "type"   // gõ chuỗi Text vào ô đang focus
	StepKey    = "key"    // bấm Keys, kèm Modifiers (tổ hợp phím) theo KeyMode
	StepScroll = "scroll" // cuộn Amount nấc theo Direction, tại Points[0] (không có điểm = chỗ chuột đang đứng)
	StepOpen   = "open"   // mở Target (URL / file / app) bằng App hoặc app mặc định của máy (open.go)

	StepWaitImage  = "waitImage"  // chờ tới khi thấy ảnh mẫu Template trên màn hình (imagestep.go)
	StepClickImage = "clickImage" // click vào chỗ thấy ảnh: output của node ImageFrom (ImageSource node) hoặc tự tìm
	StepWaitPixel  = "waitPixel"  // chờ tới khi điểm Points[0] có màu Color (± Tolerance)

	// Node nâng cao (bước 35)
	StepScreenshot = "screenshot" // chụp vùng ShotArea, lưu PNG vào SaveDir / FileName, tuỳ chọn chép vào clipboard (engine/shot.go)
	StepClipboard  = "clipboard"  // ClipOp: đọc chữ vào output / ghi Text / ghi ảnh PNG từ FilePath (engine/clip.go)
	StepHTTP       = "http"       // gọi API: Method URL + Query / Headers / Body; 2xx → success, khác → fail (engine/http.go)
	StepSystem     = "system"     // SysOp: âm lượng đặt / tăng / giảm, tắt / bật tiếng (engine/system.go)
	StepWaitUntil  = "waitUntil"  // chờ đến giờ At (hôm nay / ngày giờ), tuỳ chọn canh giờ chuẩn NTP (engine/waituntil.go)
	StepPath       = "path"       // vẽ: giữ chuột đi qua Points (đường gấp khúc) hoặc hình Shape quanh Points[0] trong DragMs (engine/draw.go)
	StepFile       = "file"       // FileOp: ghi file / thêm dòng csv · xlsx / tạo thư mục / xoá vào Thùng rác / dọn Thùng rác (engine/file.go)

	// Web qua DOM (bước 25d, engine/webstep.go): tìm phần tử bằng Selectors trên tab của Chrome của app, bấm bằng chuột thật.
	StepWebWait   = "webWait"   // chờ tới khi phần tử hiện (WebGone: mất)
	StepWebClick  = "webClick"  // chờ phần tử hiện rồi click vào tâm (+ Dx / Dy)
	StepWebType   = "webType"   // click vào ô, kiểm tra ô nhận focus rồi gõ Text bằng phím thật (WebClear: xoá chữ cũ trước)
	StepWebRead   = "webRead"   // chờ phần tử hiện rồi đọc chữ / value / thuộc tính (ReadFrom) vào output
	StepWebScroll = "webScroll" // cuộn tới phần tử / trong phần tử / tới khi thấy phần tử (ScrollMode)

	// Node điều khiển (bước 20 F3, nodes.go) – không làm gì trên màn hình ngoài nhìn, chọn cổng ra theo kết quả.
	StepIf           = "if"           // Cond đúng → cổng true, sai → false; TimeoutMs > 0 = chờ tới khi đúng
	StepSwitch       = "switch"       // so Value với từng case (Op / Match) → case:<key>, không case nào khớp → default
	StepSwitchScreen = "switchScreen" // chờ tới TimeoutMs, case (Cond ảnh / chữ / màu) nào đúng trước → case:<key>, hết giờ → none
	StepCounter      = "counter"      // đếm số lần đi qua trong vòng: lần 1..Max → next, lần Max+1 → done (rồi đếm lại từ 0)
	StepEnd          = "end"          // kết thúc vòng: EndStatus success (chạy vòng sau) / fail (dừng kịch bản, báo Message)
	StepSetVar       = "setVar"       // đặt biến {{vars.VarName}} = Value (VarOp set) hoặc cộng thêm Value (add)
)

// StepTypes – mọi loại node app hỗ trợ (steptypes.go có mô tả cho từng loại – test kiểm tra).
var StepTypes = []string{StepClick, StepDouble, StepHold, StepDrag, StepWait, StepText, StepKey, StepScroll, StepOpen,
	StepWaitImage, StepClickImage, StepWaitPixel, StepScreenshot, StepClipboard, StepHTTP, StepFile, StepSystem, StepWaitUntil, StepPath, StepWebWait, StepWebClick,
	StepWebType, StepWebRead, StepWebScroll, StepIf, StepSwitch, StepSwitchScreen, StepCounter, StepEnd, StepSetVar}

// IsControl – node điều khiển: không thử lại, không chờ sau khi xong, inspector riêng.
func IsControl(t string) bool {
	switch t {
	case StepIf, StepSwitch, StepSwitchScreen, StepCounter, StepEnd, StepSetVar:
		return true
	}
	return false
}

// Step – một node của flow (bước 20): chạy xong trả StepResult (flow.go), cổng ra trong đó chọn cạnh đi tiếp.
type Step struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`                 // StepTypes: click | … | waitPixel | if | switch | switchScreen | counter | end | setVar
	Ref           string   `json:"ref"`                  // mã ngắn n1, n2… (duy nhất) – node sau đọc kết quả bằng {{n2.output.x}}
	Name          string   `json:"name,omitempty"`       // tên hiện trên thẻ ("" = tên theo loại)
	Pos           *Pos     `json:"pos,omitempty"`        // vị trí thẻ trên canvas (nil = chưa đặt, xếp tự động)
	Breakpoint    bool     `json:"breakpoint,omitempty"` // debug: chạy từ app thì tạm dừng trước node này (bước 20 F5)
	Points        []Point  `json:"points"`               // drag dùng points[0] -> points[1]; scroll tuỳ chọn 1 điểm; type/key/wait/open không cần
	Button        string   `json:"button"`               // left | right | center
	Clicks        int      `json:"clicks"`               // click: số lần click
	ClickInterval int      `json:"clickInterval"`        // ms giữa các click
	HoldMs        int      `json:"holdMs"`
	DragMs        int      `json:"dragMs"`
	DragFromPrev  bool     `json:"dragFromPrev"` // drag: bắt đầu từ chỗ chuột dừng ở bước chuột trước (bỏ qua points[0])
	DragRelative  bool     `json:"dragRelative"` // drag: points[1] là độ lệch (dx, dy) so với điểm bắt đầu
	Text          string   `json:"text"`         // type / webType: chuỗi cần gõ
	Keys          []string `json:"keys"`         // key: phím chính theo tên robotgo (enter, tab, a, f5…)
	KeyMode       string   `json:"keyMode"`      // key: chord | sequence (chỉ khác nhau khi có ≥ 2 phím chính)
	Modifiers     []string `json:"modifiers"`    // key: cmd | ctrl | alt | shift
	DelayAfter    int      `json:"delayAfter"`   // ms, sang bước kế tiếp
	JitterMs      int      `json:"jitterMs"`     // ngẫu nhiên ± ms cho DelayAfter và ClickInterval
	JitterPx      int      `json:"jitterPx"`     // ngẫu nhiên ± px cho mỗi lần di chuột tới điểm
	Enabled       bool     `json:"enabled"`
	Direction     string   `json:"direction,omitempty"`  // scroll: up | down | left | right
	Amount        int      `json:"amount,omitempty"`     // scroll: số nấc
	Target        string   `json:"target,omitempty"`     // open: URL, đường dẫn file / app hoặc tên app; nhận {{biến}}
	App           string   `json:"app,omitempty"`        // open: mở bằng app này ("" = mặc định của máy), vd Google Chrome
	WebBrowser    bool     `json:"webBrowser,omitempty"` // open: mở URL bằng Chrome của app (profile riêng, node web điều khiển được) – bỏ qua App
	// Web qua DOM (bước 25d)
	Selectors []string `json:"selectors,omitempty"` // node web: cách tìm phần tử theo thứ tự ưu tiên (ParseSelector); nhận {{biến}}
	Tab       string   `json:"tab,omitempty"`       // node web: tab có URL / tiêu đề chứa chuỗi này ("" = tab đang hiện); nhận {{biến}}
	WebGone   bool     `json:"webGone,omitempty"`   // webWait: chờ phần tử mất (không còn / bị ẩn)
	WebClear  bool     `json:"webClear,omitempty"`  // webType: xoá chữ cũ trong ô (Cmd/Ctrl+A, Backspace) trước khi gõ
	WebEnter  bool     `json:"webEnter,omitempty"`  // webType: bấm Enter sau khi gõ
	// webWait / webClick / webType / webScroll: chạy nền (D4b) – bấm / gõ / lăn qua CDP, không đụng chuột / phím thật,
	// không đưa Chrome lên trước (cửa sổ bị che / thu nhỏ vẫn chạy); output x, y… theo toạ độ trong trang
	WebBackground bool   `json:"webBackground,omitempty"`
	ReadFrom      string `json:"readFrom,omitempty"` // webRead: "" = chữ hiển thị | value (giá trị ô) | attr (thuộc tính Attr)
	Attr          string `json:"attr,omitempty"`     // webRead ReadFrom attr: tên thuộc tính, vd href
	// webScroll: "" = tới phần tử (vào giữa khung nhìn) | in = lăn Amount nấc Direction trong phần tử | until = lăn tới khi
	// phần tử nằm trong khung nhìn (lăn trong khung ScrollIn, trống = giữa trang)
	ScrollMode string   `json:"scrollMode,omitempty"`
	ScrollIn   []string `json:"scrollIn,omitempty"` // webScroll until: cách tìm khung cuộn (như Selectors); nhận {{biến}}
	// Chờ ảnh / click theo ảnh / chờ pixel (bước 25)
	Template      string  `json:"template,omitempty"`      // ảnh mẫu PNG (base64)
	TemplateScale float64 `json:"templateScale,omitempty"` // px ảnh mẫu trên 1 đơn vị toạ độ lúc chụp (Retina = 2) – tìm trên màn tỷ lệ khác thì co giãn
	Threshold     int     `json:"threshold,omitempty"`     // % khớp tối thiểu (1–100)
	TimeoutMs     int     `json:"timeoutMs,omitempty"`     // chờ tối đa
	IntervalMs    int     `json:"intervalMs,omitempty"`    // kiểm tra mỗi
	SearchArea    string  `json:"searchArea,omitempty"`    // "" = vùng của kịch bản | screen = cả màn chứa vùng (không có vùng: màn chính)
	ImageSource   string  `json:"imageSource,omitempty"`   // clickImage: node | search
	ImageFrom     string  `json:"imageFrom,omitempty"`     // clickImage ImageSource node / screenshot ShotArea node: id node lấy chỗ thấy
	Dx            int     `json:"dx,omitempty"`            // clickImage / webClick: lệch so với tâm ảnh / phần tử
	Dy            int     `json:"dy,omitempty"`
	FindBy        string  `json:"findBy,omitempty"`       // waitImage / clickImage tự tìm: "" = ảnh mẫu | text = chữ | element = phần tử (Selectors)
	ElementApp    string  `json:"elementApp,omitempty"`   // FindBy element: tên app / bundle id ("" = app có cửa sổ trên cùng); nhận {{biến}}
	FindText      string  `json:"findText,omitempty"`     // FindBy text: chữ cần tìm; nhận {{biến}}
	TextExact     bool    `json:"textExact,omitempty"`    // FindBy text: cả dòng phải đúng chuỗi (mặc định: chứa chuỗi)
	TextNoAccent  bool    `json:"textNoAccent,omitempty"` // FindBy text: bỏ qua dấu
	Occurrence    int     `json:"occurrence,omitempty"`   // FindBy text: chỗ thứ mấy (trên → dưới, trái → phải); 0 = chỗ đầu
	Retry         *Retry  `json:"retry,omitempty"`        // thất bại thì thử lại; hết lượt vẫn lỗi → cổng fail
	Human         string  `json:"human,omitempty"`        // mức giống người của node: "" = theo kịch bản | off | low | medium | high
	// Ô số lấy từ biến (F5c): mỗi trường số trong NumVarFieldsOf một dòng – giá trị lúc chạy thay cho trường đó.
	NumVars []NumVar `json:"numVars,omitempty"`
	// Chụp màn hình (bước 35 N1)
	ShotArea    string `json:"shotArea,omitempty"`    // ShotAreas: "" = vùng của kịch bản | screen | rect | node (ImageFrom)
	Display     int    `json:"display,omitempty"`     // shotArea screen: 0 = màn chứa vùng kịch bản (không vùng: màn chính), n = màn thứ n
	ShotRect    *Rect  `json:"shotRect,omitempty"`    // shotArea rect: vùng chụp, toạ độ màn hình
	SaveDir     string `json:"saveDir,omitempty"`     // thư mục lưu ("" = Pictures/autoclicker, ~ = thư mục người dùng); nhận {{biến}}
	FileName    string `json:"fileName,omitempty"`    // tên file ("" = autoclicker-<ngày>-<giờ>; tự thêm .png, trùng thì thêm số); nhận {{biến}}
	ToClipboard bool   `json:"toClipboard,omitempty"` // chép ảnh vào clipboard
	// Clipboard (bước 35 N2)
	ClipOp   string `json:"clipOp,omitempty"`   // ClipOps: "" = đọc chữ | write = ghi Text | image = ghi ảnh từ FilePath
	FilePath string `json:"filePath,omitempty"` // clipboard image: file PNG / file: file hoặc thư mục (~ / tương đối = từ thư mục người dùng); nhận {{biến}}
	// HTTP (bước 35 N3, http.go) – mọi ô chữ nhận {{biến}}; header / tham số tên bí mật (IsSecretKey) che trong nhật ký
	Method   string `json:"method,omitempty"`   // HTTPMethods ("" = GET)
	URL      string `json:"url,omitempty"`      // http(s)://… – biến chèn nguyên văn (giá trị cần mã hoá thì đặt trong Query)
	Query    []KV   `json:"query,omitempty"`    // tham số thêm vào URL (tự mã hoá)
	Headers  []KV   `json:"headers,omitempty"`  // header gửi kèm
	BodyType string `json:"bodyType,omitempty"` // BodyTypes: "" = không body | json | form | text
	Body     string `json:"body,omitempty"`     // bodyType json / text
	Form     []KV   `json:"form,omitempty"`     // bodyType form
	// Tệp / thư mục (bước 35 N4, file.go) – FilePath + Text (ghi) dùng chung với clipboard / gõ chữ
	FileOp      string   `json:"fileOp,omitempty"`      // FileOps: "" = ghi chữ | row = thêm dòng csv / xlsx | mkdir | trash | emptyTrash
	FileAppend  bool     `json:"fileAppend,omitempty"`  // ghi chữ: thêm vào cuối thay vì ghi đè
	Cells       []string `json:"cells,omitempty"`       // row: các ô của dòng mới; nhận {{biến}}
	Sheet       string   `json:"sheet,omitempty"`       // row xlsx: tên sheet ("" = sheet đầu tiên); nhận {{biến}}
	FilePattern string   `json:"filePattern,omitempty"` // trash: mẫu tên file trong thư mục FilePath, vd *.png; nhận {{biến}}
	// Hệ thống (bước 35 N5, system.go)
	SysOp  string `json:"sysOp,omitempty"`  // SysOps: volume = đặt mức | volumeUp | volumeDown | mute | unmute
	Volume int    `json:"volume,omitempty"` // volume: mức 0–100 %; volumeUp / volumeDown: mỗi lần bao nhiêu % (0 = 10)
	// Chờ đến giờ (bước 35 N6, waituntil.go)
	At       string `json:"at,omitempty"`       // HH:mm[:ss[.mmm]] hôm nay, có thể kèm ngày YYYY-MM-DD / DD/MM/YYYY (script.ParseAt); nhận {{biến}}
	IfPassed string `json:"ifPassed,omitempty"` // giờ đã qua: "" = đi tiếp ngay | tomorrow = chờ giờ đó ngày mai (chỉ khi không ghi ngày) | fail
	UseNTP   bool   `json:"useNtp,omitempty"`   // canh theo giờ chuẩn internet (NTP) thay vì đồng hồ máy
	// Vẽ đường (bước 35 N7, draw.go) – Button, DragMs (thời gian vẽ cả nét) dùng chung với kéo thả
	Shape     string `json:"shape,omitempty"`     // Shapes: "" = đi qua Points | circle | rect | star (tâm Points[0], cỡ ShapeW × ShapeH)
	ShapeW    int    `json:"shapeW,omitempty"`    // rộng / đường kính (px)
	ShapeH    int    `json:"shapeH,omitempty"`    // cao (0 = bằng ShapeW)
	Closed    bool   `json:"closed,omitempty"`    // đường gấp khúc: nối điểm cuối về điểm đầu
	Color     string `json:"color,omitempty"`     // waitPixel: #RRGGBB
	Tolerance int    `json:"tolerance,omitempty"` // waitPixel: sai số mỗi kênh màu (0–255)
	// Node điều khiển (bước 20 F3, nodes.go)
	Cond      *Cond  `json:"cond,omitempty"`      // if
	Cases     []Case `json:"cases,omitempty"`     // switch / switchScreen – mỗi case một cổng ra case:<key>
	Value     string `json:"value,omitempty"`     // switch: giá trị cần xét; setVar: giá trị gán / cộng; nhận {{biến}}
	Max       int    `json:"max,omitempty"`       // counter: số lần đi cổng next tối đa
	EndStatus string `json:"endStatus,omitempty"` // end: success | fail
	Message   string `json:"message,omitempty"`   // end: lời nhắn (fail = lời báo lỗi); nhận {{biến}}
	VarName   string `json:"varName,omitempty"`   // setVar: tên biến
	VarOp     string `json:"varOp,omitempty"`     // setVar: set | add
}

// NumVar – trường số Field (tên json, vd clicks, delayAfter) lấy giá trị từ Expr: số / biến cộng trừ, vd
// "{{n2.output.count}} + 1" (như ô toạ độ F5b, script.EvalCoord). Slice thay vì map để thay biến tại chỗ như ô chữ.
type NumVar struct {
	Field string `json:"field"`
	Expr  string `json:"expr"`
}

// NumVarOf – biểu thức biến của trường số field ("" = trường dùng số cố định).
func (s Step) NumVarOf(field string) string {
	for _, v := range s.NumVars {
		if v.Field == field {
			return v.Expr
		}
	}
	return ""
}

// Case – một nhánh của switch / switchScreen. Key duy nhất trong node, cổng ra là "case:" + Key.
type Case struct {
	Key   string `json:"key"`
	Label string `json:"label,omitempty"` // tên hiện trên cổng ("" = Key)
	Op    string `json:"op,omitempty"`    // switch: eq | contains | regex
	Match string `json:"match,omitempty"` // switch: so Value với chuỗi này; nhận {{biến}}
	Cond  *Cond  `json:"cond,omitempty"`  // switchScreen: điều kiện ảnh / chữ / màu
}

// Cond – điều kiện của node Nếu / case của Switch theo màn hình (cond.go). Kind image / text / pixel dùng cùng trường
// với node Chờ ảnh / chữ / pixel; value so 2 vế; all / any ghép nhiều điều kiện.
type Cond struct {
	Kind string `json:"kind"`          // image | text | element | pixel | value | all | any
	Not  bool   `json:"not,omitempty"` // đảo kết quả

	Template      string   `json:"template,omitempty"`
	TemplateScale float64  `json:"templateScale,omitempty"`
	Threshold     int      `json:"threshold,omitempty"`
	SearchArea    string   `json:"searchArea,omitempty"`
	FindText      string   `json:"findText,omitempty"` // nhận {{biến}}
	TextExact     bool     `json:"textExact,omitempty"`
	TextNoAccent  bool     `json:"textNoAccent,omitempty"`
	Occurrence    int      `json:"occurrence,omitempty"`
	Selectors     []string `json:"selectors,omitempty"`  // element: cách tìm phần tử (text= / aria= / role=); nhận {{biến}}
	ElementApp    string   `json:"elementApp,omitempty"` // element: tên app ("" = app có cửa sổ trên cùng); nhận {{biến}}
	Points        []Point  `json:"points,omitempty"`     // pixel: points[0]
	Color         string   `json:"color,omitempty"`
	Tolerance     int      `json:"tolerance,omitempty"`

	Left  string `json:"left,omitempty"`  // value: vế trái, vd {{n2.output.text}}; nhận {{biến}}
	Op    string `json:"op,omitempty"`    // value: eq | ne | contains | gt | lt | empty
	Right string `json:"right,omitempty"` // value: vế phải (không dùng với empty); nhận {{biến}}

	Conds []Cond `json:"conds,omitempty"` // all (và) / any (hoặc)
}

// Step.FindBy
const (
	FindByText    = "text"
	FindByElement = "element" // phần tử màn hình qua Accessibility (bước 37)
)

// Step.ReadFrom (webRead)
const (
	ReadValue = "value"
	ReadAttr  = "attr"
)

// Step.ScrollMode (webScroll)
const (
	ScrollModeIn    = "in"
	ScrollModeUntil = "until"
)

// Retry – node thất bại (không thấy ảnh / chữ / màu, lỗi) thì thử lại thêm Count lần, nghỉ DelayMs giữa các lần.
// Thử hết vẫn lỗi → StepResult fail, đi ra cổng fail (không nối thì dừng kịch bản).
type Retry struct {
	Count   int `json:"count,omitempty"`
	DelayMs int `json:"delayMs,omitempty"`
}

const MaxRetries = 100

// Pos – vị trí thẻ node trên canvas (toạ độ canvas, không phải màn hình).
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Edge – cạnh của flow: node From đi ra cổng Port thì chạy node To. Mỗi cổng tối đa một cạnh.
type Edge struct {
	ID   string `json:"id"`
	From string `json:"from"`
	Port string `json:"port"`
	To   string `json:"to"`
}

// Cách neo điểm mới đặt (Script.AnchorMode).
const (
	AnchorNearest = ""      // góc gần nhất theo từng trục (mặc định)
	AnchorTopLeft = "start" // luôn tính từ góc trên trái
)

// ScriptVersion – định dạng file kịch bản hiện tại. Tăng khi đổi model theo cách file cũ không đọc đúng được.
// 2 (bước 21b): thêm inputs – app cũ không biết inputs nên không được mở file mới.
// 4 (bước 20): kịch bản thành flow node + cạnh. File cũ hơn KHÔNG chuyển đổi (người dùng chốt 2026-10-17) – báo không mở được.
const ScriptVersion = 4

type Script struct {
	Version    int           `json:"version"` // định dạng file, xem ScriptVersion
	Name       string        `json:"name"`
	Region     *Rect         `json:"region"`     // nil = toàn màn hình (màn chính); toạ độ điểm neo theo cạnh vùng
	Loops      int           `json:"loops"`      // 0 = vô hạn
	MaxMinutes int           `json:"maxMinutes"` // 0 = không giới hạn
	Steps      []Step        `json:"steps"`      // các node của flow – thứ tự trong mảng không quyết định thứ tự chạy
	Start      string        `json:"start"`      // id node bắt đầu
	Edges      []Edge        `json:"edges"`
	AnchorMode string        `json:"anchorMode,omitempty"` // "" = góc gần nhất | start = luôn trên trái
	Inputs     []ScriptInput `json:"inputs,omitempty"`     // giá trị hỏi khi bấm Chạy (script/inputs.go)
	Human      string        `json:"human,omitempty"`      // thao tác giống người: "" / off = tắt | low | medium | high
}
