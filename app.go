package main

import (
	"context"
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
	mainDisp int        // màn chứa cửa sổ main lúc rời main (chỉ số displays()), quay về main thì dời về đây
	winDisp  int        // màn cửa sổ đang nằm (-1 = chưa biết)
	ovDisp   int        // màn người dùng chọn cho overlay (-1 = màn chứa cửa sổ app)
	dirty    bool       // kịch bản đang soạn có thay đổi chưa lưu (storage.go)
	quitting bool       // đã hỏi xong, đóng thật
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{mode: ModeMain, winDisp: -1, ovDisp: -1}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startHotkeys()
	excludeSelfFromCapture() // Windows: thanh mini / overlay không lọt vào ảnh chụp khi tìm ảnh
}

// shutdown is called when the app is closing.
func (a *App) shutdown(ctx context.Context) {
	a.StopScript()
	a.stopHotkeys()
}
