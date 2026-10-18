package main

import (
	"autoclicker/internal/device"
	"autoclicker/internal/model"
)

// ScreenBounds trả khung màn hình chính – frontend dùng để tính điểm neo khi chọn toàn màn hình.
func (a *App) ScreenBounds() model.Rect { return device.MainScreen() }
