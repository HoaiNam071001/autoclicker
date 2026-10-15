package main

import "github.com/go-vgo/robotgo"

// Input – thao tác chuột / bàn phím mà engine dùng. Tách interface để test engine bằng input giả.
type Input interface {
	Move(x, y int) error
	Drag(x, y int, button string) error // di chuột khi đang giữ nút
	Toggle(button string, down bool) error
	Click(button string, double bool) error
	TypeStr(s string) error
	KeyTap(key string, modifiers []string) error
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

func (robotgoInput) KeyTap(key string, modifiers []string) error {
	if len(modifiers) == 0 {
		return robotgo.KeyTap(key)
	}
	return robotgo.KeyTap(key, modifiers)
}
