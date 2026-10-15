package main

// Giữ các dependency sẽ dùng ở Phase 0 (spike click / hotkey) trong go.mod.
// Xoá file này khi các package đã được import thật.
import (
	_ "github.com/go-vgo/robotgo"
	_ "github.com/robotn/gohook"
)
