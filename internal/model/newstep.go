package model

import (
	"fmt"
	"sync/atomic"
	"time"
)

// stepSeq – thêm vào id bước mới: gọi NewStep liên tiếp có thể cùng UnixNano (đồng hồ thô) → id trùng.
var stepSeq atomic.Int64

// NewStep trả bước mới với giá trị mặc định theo loại – dùng cho nút "+ Thêm bước".
func NewStep(stepType string) (Step, error) {
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
	case StepScreenshot, StepClipboard, StepFile:
		s.DelayAfter = 300
	case StepSystem:
		s.SysOp, s.Volume, s.DelayAfter = SysVolume, 50, 300
	case StepWindow, StepApp:
		s.TimeoutMs, s.DelayAfter = 10000, 300
	case StepWaitUntil:
		s.At, s.DelayAfter = "12:00:00", 0
	case StepPath:
		s.Points, s.DragMs, s.Clicks = []Point{{}, {X: 200}, {X: 200, Y: 150}}, 1500, 1
	case StepHTTP:
		s.Method, s.URL, s.TimeoutMs, s.DelayAfter = "GET", "https://", 15000, 0
	case StepWebWait:
		s.Selectors, s.TimeoutMs, s.IntervalMs, s.DelayAfter = []string{"text="}, 10000, 200, 300
	case StepWebClick:
		s.Selectors, s.Clicks, s.TimeoutMs, s.IntervalMs = []string{"text="}, 1, 10000, 200
	case StepWebType:
		s.Selectors, s.WebClear, s.TimeoutMs, s.IntervalMs, s.DelayAfter = []string{"text="}, true, 10000, 200, 300
	case StepWebRead:
		s.Selectors, s.TimeoutMs, s.IntervalMs = []string{"text="}, 10000, 200
	case StepElementRead:
		s.Selectors, s.TimeoutMs, s.IntervalMs = []string{"text="}, 10000, 250
	case StepWebScroll:
		s.Selectors, s.Direction, s.Amount, s.TimeoutMs, s.IntervalMs, s.DelayAfter = []string{"text="}, "down", 3, 10000, 200, 300
	case StepIf:
		s.Cond, s.IntervalMs = &Cond{Kind: CondText}, 250
	case StepSwitch:
		s.Cases = []Case{{Key: "c1", Label: "Case 1", Op: "eq"}, {Key: "c2", Label: "Case 2", Op: "eq"}}
	case StepSwitchScreen:
		s.Cases = []Case{{Key: "c1", Label: "Case 1", Cond: &Cond{Kind: CondText}}, {Key: "c2", Label: "Case 2", Cond: &Cond{Kind: CondText}}}
		s.TimeoutMs, s.IntervalMs = 3000, 250
	case StepCounter:
		s.Max = 10
	case StepEnd:
		s.EndStatus = EndSuccess
	case StepSetVar:
		s.VarName, s.VarOp = "ten", VarSet
	default:
		return Step{}, fmt.Errorf("loại bước không hợp lệ: %q", stepType)
	}
	if IsControl(stepType) {
		s.DelayAfter = 0
	}
	return s, nil
}

// NewScript trả kịch bản rỗng mặc định.
func NewScript() Script {
	return Script{Version: ScriptVersion, Name: "Kịch bản 1", Loops: 1, Steps: []Step{}, Edges: []Edge{}, Human: HumanMedium}
}
