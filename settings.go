package main

import (
	"autoclicker/internal/store"
)

// GetSettings trả cấu hình đã lưu (mặc định nếu chưa có / không đọc được thư mục cấu hình).
func (a *App) GetSettings() store.Settings { return store.CurrentSettings() }

// SetTheme lưu theme: light | dark | system.
func (a *App) SetTheme(theme string) error { return store.SetTheme(theme) }
