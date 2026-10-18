// Package device – thiết bị thật: chuột / phím (robotgo), chụp màn hình, đọc chữ (OCR), màn hình, dừng khẩn khi chuột
// chạm góc. Mọi thứ qua interface (Input, Screen, TextReader) để test bằng bản giả (devicetest).
package device

import (
	goruntime "runtime"

	"github.com/go-vgo/robotgo"

	"autoclicker/internal/model"
)

// Input – thao tác chuột / bàn phím mà engine dùng. Tách interface để test engine bằng input giả.
type Input interface {
	Move(x, y int) error
	Drag(x, y int, button string) error // di chuột khi đang giữ nút
	Toggle(button string, down bool) error
	Click(button string, double bool) error
	TypeStr(s string) error
	// KeyToggle nhấn (down) hoặc thả một phím. modifiers: cờ phím kèm gửi cùng sự kiện nhấn
	// (macOS cần cờ này mới hiểu là tổ hợp); phím kèm được nhấn / thả riêng bằng chính KeyToggle.
	KeyToggle(key string, down bool, modifiers []string) error
	Scroll(direction string) error // cuộn một nấc: up | down | left | right
}

// RobotgoInput – Input thật qua robotgo.
type RobotgoInput struct{}

func (RobotgoInput) Move(x, y int) error { return robotgo.Move(x, y) }

// Drag dùng robotgo.Drag (đã deprecated) vì trên macOS Move chỉ gửi sự kiện "di chuột" kể cả khi
// đang giữ nút – nhiều app không nhận là kéo; Drag gửi đúng sự kiện "kéo". Windows/Linux: như Move.
func (RobotgoInput) Drag(x, y int, button string) error { return robotgo.Drag(x, y, button) }

func (RobotgoInput) Toggle(button string, down bool) error {
	if down {
		return robotgo.Toggle(button)
	}
	return robotgo.Toggle(button, "up")
}

func (RobotgoInput) Click(button string, double bool) error { return robotgo.Click(button, double) }

func (RobotgoInput) TypeStr(s string) error { return robotgo.TypeStr(s) }

// scrollNotchPx – quãng một nấc cuộn trên macOS. robotgo gửi sự kiện cuộn macOS theo pixel (Windows theo nấc con lăn,
// Chrome / Edge ≈ 100 px một nấc) – gửi 100 px để một nấc hai bên gần bằng nhau. Khớp SCROLL_NOTCH_PX (types.ts).
const scrollNotchPx = 100

func (RobotgoInput) Scroll(direction string) error {
	if goruntime.GOOS == "darwin" {
		return robotgo.ScrollDir(scrollNotchPx, direction)
	}
	return robotgo.ScrollDir(1, direction)
}

func (RobotgoInput) KeyToggle(key string, down bool, modifiers []string) error {
	if !down {
		// Không truyền modifiers khi thả: robotgo sẽ nhả luôn các phím kèm đang giữ.
		return robotgo.KeyToggle(key, "up")
	}
	if len(modifiers) == 0 {
		return robotgo.KeyToggle(key, "down")
	}
	return robotgo.KeyToggle(key, "down", modifiers)
}

// RawDisplays trả khung các màn hình như robotgo báo (chưa nhân tỷ lệ) – cùng đơn vị với vị trí cửa sổ Wails.
func RawDisplays() []model.Rect {
	n := robotgo.DisplaysNum()
	out := make([]model.Rect, 0, n)
	for i := 0; i < n; i++ {
		x, y, w, h := robotgo.GetDisplayBounds(i)
		out = append(out, model.Rect{X: x, Y: y, W: w, H: h})
	}
	return out
}

// Displays trả khung các màn hình theo hệ toạ độ của robotgo.Location / Move.
// Windows: robotgo nhân toạ độ chuột với tỷ lệ DPI màn chính (ScaleF) nên khung màn hình cũng phải nhân theo.
func Displays() []model.Rect {
	f := 1.0
	if goruntime.GOOS == "windows" {
		f = robotgo.ScaleF()
	}
	sc := func(v int) int { return int(float64(v) * f) }
	n := robotgo.DisplaysNum()
	out := make([]model.Rect, 0, n)
	for i := 0; i < n; i++ {
		x, y, w, h := robotgo.GetDisplayBounds(i)
		out = append(out, model.Rect{X: sc(x), Y: sc(y), W: sc(w), H: sc(h)})
	}
	return out
}

// MainScreen – khung màn hình chính theo toạ độ robotgo; "toàn màn hình" (Region = nil) neo theo khung này.
func MainScreen() model.Rect {
	if d := Displays(); len(d) > 0 {
		return d[0]
	}
	w, h := robotgo.GetScreenSize()
	return model.Rect{W: w, H: h}
}
