package main

import (
	"context"
	"fmt"
	"sync"
)

// App struct
type App struct {
	ctx context.Context

	mu    sync.Mutex
	spike *spikeRun // khác nil khi vòng lặp click thử đang chạy
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.startHotkeys()
}

// shutdown is called when the app is closing.
func (a *App) shutdown(ctx context.Context) {
	a.StopSpikeLoop()
	a.stopHotkeys()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
