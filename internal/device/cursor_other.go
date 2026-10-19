//go:build !windows

package device

import "github.com/go-vgo/robotgo"

// macOS / Linux: robotgo dùng điểm (point) cho cả chuột và khung màn hình – không phải đổi.
func cursorPos() (int, int)                    { return robotgo.Location() }
func moveCursor(x, y int) error                { return robotgo.Move(x, y) }
func dragCursor(x, y int, button string) error { return robotgo.Drag(x, y, button) }
