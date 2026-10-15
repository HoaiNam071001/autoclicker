package main

import "github.com/go-vgo/robotgo"

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
}

// robotgoInput – Input thật qua robotgo.
type robotgoInput struct{}

func (robotgoInput) Move(x, y int) error { return robotgo.Move(x, y) }

// Drag dùng robotgo.Drag (đã deprecated) vì trên macOS Move chỉ gửi sự kiện "di chuột" kể cả khi
// đang giữ nút – nhiều app không nhận là kéo; Drag gửi đúng sự kiện "kéo". Windows/Linux: như Move.
func (robotgoInput) Drag(x, y int, button string) error { return robotgo.Drag(x, y, button) }

func (robotgoInput) Toggle(button string, down bool) error {
	if down {
		return robotgo.Toggle(button)
	}
	return robotgo.Toggle(button, "up")
}

func (robotgoInput) Click(button string, double bool) error { return robotgo.Click(button, double) }

func (robotgoInput) TypeStr(s string) error { return robotgo.TypeStr(s) }

func (robotgoInput) KeyToggle(key string, down bool, modifiers []string) error {
	if !down {
		// Không truyền modifiers khi thả: robotgo sẽ nhả luôn các phím kèm đang giữ.
		return robotgo.KeyToggle(key, "up")
	}
	if len(modifiers) == 0 {
		return robotgo.KeyToggle(key, "down")
	}
	return robotgo.KeyToggle(key, "down", modifiers)
}
