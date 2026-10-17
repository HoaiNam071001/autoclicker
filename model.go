package main

import (
	"fmt"
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

type Step struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`          // click | double | hold | drag | wait | type | key | scroll | open
	Points        []Point  `json:"points"`        // drag dùng points[0] -> points[1]; scroll tuỳ chọn 1 điểm; type/key/wait/open không cần
	Button        string   `json:"button"`        // left | right | center
	Clicks        int      `json:"clicks"`        // click: số lần click
	ClickInterval int      `json:"clickInterval"` // ms giữa các click
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
	Steps      []Step        `json:"steps"`
	AnchorMode string        `json:"anchorMode,omitempty"` // "" = góc gần nhất | start = luôn trên trái
	Inputs     []ScriptInput `json:"inputs,omitempty"`     // giá trị hỏi khi bấm Chạy (inputs.go)
}

// NewStep trả bước mới với giá trị mặc định theo loại – dùng cho nút "+ Thêm bước".
func (a *App) NewStep(stepType string) (Step, error) {
	s := Step{
		ID:         fmt.Sprintf("%d", time.Now().UnixNano()),
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
	default:
		return Step{}, fmt.Errorf("loại bước không hợp lệ: %q", stepType)
	}
	return s, nil
}

// NewScript trả kịch bản rỗng mặc định.
func (a *App) NewScript() Script {
	return Script{Version: ScriptVersion, Name: "Kịch bản 1", Loops: 1, Steps: []Step{}}
}
