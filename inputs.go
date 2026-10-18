package main

import (
	"autoclicker/internal/store"
)

// GetInputMemory trả giá trị lần trước / preset của kịch bản id (đường dẫn file, hoặc tên khi chưa lưu).
func (a *App) GetInputMemory(id string) store.InputMemory { return store.InputMemoryFor(id) }

// SaveInputMemory ghi giá trị lần trước / preset của kịch bản id.
func (a *App) SaveInputMemory(id string, m store.InputMemory) error {
	return store.SaveInputMemory(id, m)
}
