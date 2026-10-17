package main

import (
	"context"
	"fmt"
	"sync"
)

// App struct
type App struct {
	ctx context.Context

	mu       sync.Mutex
	run      *runHandle // khác nil khi kịch bản đang chạy
	mode     string     // chế độ cửa sổ hiện tại, xem window.go
	mainRect winRect    // vị trí/kích thước main để khôi phục
	mainMax  bool       // main đang phóng to (nút phóng to trên thanh tiêu đề) – phóng to lại khi quay về
	dirty    bool       // kịch bản đang soạn có thay đổi chưa lưu (storage.go)
	quitting bool       // đã hỏi xong, đóng thật
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{mode: ModeMain}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startHotkeys()
}

// shutdown is called when the app is closing.
func (a *App) shutdown(ctx context.Context) {
	a.StopScript()
	a.stopHotkeys()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
