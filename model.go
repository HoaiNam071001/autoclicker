package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

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
)

// Nguồn vị trí của bước click theo ảnh (Step.ImageSource).
const (
	ImageFromNode   = "node"   // chỗ node ImageFrom (chờ ảnh / chữ, click theo ảnh tự tìm) tìm thấy lần chạy gần nhất
	ImageFromSearch = "search" // tự tìm ảnh mẫu riêng của bước
)

// Hướng cuộn (Step.Direction).
var scrollDirections = []string{"up", "down", "left", "right"}

// Cách bấm khi bước phím có nhiều phím chính (Step.KeyMode).
const (
	KeyChord    = "chord"    // nhấn giữ phím kèm + mọi phím chính cùng lúc rồi nhả (vd Ctrl+C+T)
	KeySequence = "sequence" // giữ phím kèm, gõ lần lượt từng phím chính rồi nhả phím kèm
)

// Point – toạ độ điểm trong vùng, mỗi trục neo riêng (anchor.go): Ax / Ay trống = start (cách cạnh trái / trên).
// Điểm B kéo thả "lệch so với A" (DragRelative) là độ lệch, không dùng neo.
type Point struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Ax string `json:"ax,omitempty"` // start | end | center | ratio
	Ay string `json:"ay,omitempty"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// Step – một node của flow (bước 20): chạy xong trả StepResult (flow.go), cổng ra trong đó chọn cạnh đi tiếp.
type Step struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`           // click | double | hold | drag | wait | type | key | scroll | open | waitImage | clickImage | waitPixel
	Ref           string   `json:"ref"`            // mã ngắn n1, n2… (duy nhất) – node sau đọc kết quả bằng {{n2.output.x}}
	Name          string   `json:"name,omitempty"` // tên hiện trên thẻ ("" = tên theo loại)
	Pos           *Pos     `json:"pos,omitempty"`  // vị trí thẻ trên canvas (nil = chưa đặt, xếp tự động)
	Points        []Point  `json:"points"`         // drag dùng points[0] -> points[1]; scroll tuỳ chọn 1 điểm; type/key/wait/open không cần
	Button        string   `json:"button"`         // left | right | center
	Clicks        int      `json:"clicks"`         // click: số lần click
	ClickInterval int      `json:"clickInterval"`  // ms giữa các click
	HoldMs        int      `json:"holdMs"`
	DragMs        int      `json:"dragMs"`
	DragFromPrev  bool     `json:"dragFromPrev"` // drag: bắt đầu từ chỗ chuột dừng ở bước chuột trước (bỏ qua points[0])
	DragRelative  bool     `json:"dragRelative"` // drag: points[1] là độ lệch (dx, dy) so với điểm bắt đầu
	Text          string   `json:"text"`         // type: chuỗi cần gõ
	Keys          []string `json:"keys"`         // key: phím chính theo tên robotgo (enter, tab, a, f5…)
	KeyMode       string   `json:"keyMode"`      // key: chord | sequence (chỉ khác nhau khi có ≥ 2 phím chính)
	Modifiers     []string `json:"modifiers"`    // key: cmd | ctrl | alt | shift
	DelayAfter    int      `json:"delayAfter"`   // ms, sang bước kế tiếp
	JitterMs      int      `json:"jitterMs"`     // ngẫu nhiên ± ms cho DelayAfter và ClickInterval
	JitterPx      int      `json:"jitterPx"`     // ngẫu nhiên ± px cho mỗi lần di chuột tới điểm
	Enabled       bool     `json:"enabled"`
	Direction     string   `json:"direction,omitempty"` // scroll: up | down | left | right
	Amount        int      `json:"amount,omitempty"`    // scroll: số nấc
	Target        string   `json:"target,omitempty"`    // open: URL, đường dẫn file / app hoặc tên app; nhận {{biến}}
	App           string   `json:"app,omitempty"`       // open: mở bằng app này ("" = mặc định của máy), vd Google Chrome
	// Chờ ảnh / click theo ảnh / chờ pixel (bước 25)
	Template      string  `json:"template,omitempty"`      // ảnh mẫu PNG (base64)
	TemplateScale float64 `json:"templateScale,omitempty"` // px ảnh mẫu trên 1 đơn vị toạ độ lúc chụp (Retina = 2) – tìm trên màn tỷ lệ khác thì co giãn
	Threshold     int     `json:"threshold,omitempty"`     // % khớp tối thiểu (1–100)
	TimeoutMs     int     `json:"timeoutMs,omitempty"`     // chờ tối đa
	IntervalMs    int     `json:"intervalMs,omitempty"`    // kiểm tra mỗi
	SearchArea    string  `json:"searchArea,omitempty"`    // "" = vùng của kịch bản | screen = cả màn chứa vùng (không có vùng: màn chính)
	ImageSource   string  `json:"imageSource,omitempty"`   // clickImage: node | search
	ImageFrom     string  `json:"imageFrom,omitempty"`     // clickImage ImageSource node: id node lấy chỗ thấy
	Dx            int     `json:"dx,omitempty"`            // clickImage: lệch so với tâm ảnh
	Dy            int     `json:"dy,omitempty"`
	FindBy        string  `json:"findBy,omitempty"`       // waitImage / clickImage tự tìm: "" = ảnh mẫu | text = chữ
	FindText      string  `json:"findText,omitempty"`     // FindBy text: chữ cần tìm; nhận {{biến}}
	TextExact     bool    `json:"textExact,omitempty"`    // FindBy text: cả dòng phải đúng chuỗi (mặc định: chứa chuỗi)
	TextNoAccent  bool    `json:"textNoAccent,omitempty"` // FindBy text: bỏ qua dấu
	Occurrence    int     `json:"occurrence,omitempty"`   // FindBy text: chỗ thứ mấy (trên → dưới, trái → phải); 0 = chỗ đầu
	Retry         *Retry  `json:"retry,omitempty"`        // thất bại thì thử lại; hết lượt vẫn lỗi → cổng fail
	Color         string  `json:"color,omitempty"`        // waitPixel: #RRGGBB
	Tolerance     int     `json:"tolerance,omitempty"`    // waitPixel: sai số mỗi kênh màu (0–255)
}

// Step.FindBy
const FindByText = "text"

// Retry – node thất bại (không thấy ảnh / chữ / màu, lỗi) thì thử lại thêm Count lần, nghỉ DelayMs giữa các lần.
// Thử hết vẫn lỗi → StepResult fail, đi ra cổng fail (không nối thì dừng kịch bản).
type Retry struct {
	Count   int `json:"count,omitempty"`
	DelayMs int `json:"delayMs,omitempty"`
}

const maxRetries = 100

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

type Script struct {
	Version    int           `json:"version"` // định dạng file, xem storage.go
	Name       string        `json:"name"`
	Region     *Rect         `json:"region"`     // nil = toàn màn hình (màn chính); toạ độ điểm neo theo cạnh vùng
	Loops      int           `json:"loops"`      // 0 = vô hạn
	MaxMinutes int           `json:"maxMinutes"` // 0 = không giới hạn
	Steps      []Step        `json:"steps"`      // các node của flow – thứ tự trong mảng không quyết định thứ tự chạy
	Start      string        `json:"start"`      // id node bắt đầu
	Edges      []Edge        `json:"edges"`
	AnchorMode string        `json:"anchorMode,omitempty"` // "" = góc gần nhất | start = luôn trên trái
	Inputs     []ScriptInput `json:"inputs,omitempty"`     // giá trị hỏi khi bấm Chạy (inputs.go)
}

// stepSeq – thêm vào id bước mới: gọi NewStep liên tiếp có thể cùng UnixNano (đồng hồ thô) → id trùng.
var stepSeq atomic.Int64

// NewStep trả bước mới với giá trị mặc định theo loại – dùng cho nút "+ Thêm bước".
func (a *App) NewStep(stepType string) (Step, error) {
	s := Step{
		ID:         fmt.Sprintf("%d-%d", time.Now().UnixNano(), stepSeq.Add(1)),
		Type:       stepType,
		Points:     []Point{},
		Button:     "left",
		Keys:       []string{},
		KeyMode:    KeyChord,
		Modifiers:  []string{},
		DelayAfter: 500,
		Enabled:    true,
	}
	switch stepType {
	case StepClick:
		s.Clicks = 1
		s.Points = []Point{{}}
	case StepDouble:
		s.Points = []Point{{}}
	case StepHold:
		s.HoldMs = 1000
		s.Points = []Point{{}}
	case StepDrag:
		s.DragMs = 400
		s.Points = []Point{{}, {}}
	case StepWait:
		s.DelayAfter = 1000
	case StepText:
	case StepKey:
		s.Keys = []string{"enter"}
	case StepScroll:
		s.Direction, s.Amount = "down", 3
	case StepOpen:
		s.Target = "https://"
		s.DelayAfter = 2000 // chờ trang / app mở
	case StepWaitImage:
		s.Threshold, s.TimeoutMs, s.IntervalMs, s.DelayAfter = 90, 10000, 250, 300
	case StepClickImage:
		s.ImageSource, s.Threshold, s.TimeoutMs, s.IntervalMs = ImageFromNode, 90, 10000, 250
	case StepWaitPixel:
		s.Points = []Point{{}}
		s.Color, s.Tolerance, s.TimeoutMs, s.IntervalMs, s.DelayAfter = "#000000", 10, 5000, 100, 300
	default:
		return Step{}, fmt.Errorf("loại bước không hợp lệ: %q", stepType)
	}
	return s, nil
}

// NewScript trả kịch bản rỗng mặc định.
func (a *App) NewScript() Script {
	return Script{Version: ScriptVersion, Name: "Kịch bản 1", Loops: 1, Steps: []Step{}, Edges: []Edge{}}
}
