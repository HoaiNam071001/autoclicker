package main

import (
	"errors"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
	"autoclicker/internal/store"
)

// Xác nhận thao tác xoá (node Tệp / thư mục xoá vào / dọn Thùng rác, bước 35 N4): kịch bản có thao tác xoá thì lần chạy
// đầu (và mỗi khi thư mục / mẫu tên đổi) người dùng phải đồng ý. Chặn ở Go nên hẹn giờ / CLI / MCP (bước 22, 31) cũng
// không xoá được khi chưa đồng ý trong app.

// DeleteCheck – thao tác xoá của kịch bản và đã được đồng ý chưa (frontend hỏi trước khi chạy).
type DeleteCheck struct {
	Ops         []script.DeleteOp `json:"ops"`
	Fingerprint string            `json:"fingerprint"`
	Confirmed   bool              `json:"confirmed"`
}

// errDeleteNotConfirmed – chạy kịch bản có thao tác xoá chưa được đồng ý.
var errDeleteNotConfirmed = errors.New("kịch bản có thao tác xoá file chưa được xác nhận – bấm Chạy trong app để xem và đồng ý")

// DeleteCheck áp input như lúc chạy rồi liệt kê thao tác xoá. path: file kịch bản ("" = chưa lưu, theo tên).
func (a *App) DeleteCheck(sc model.Script, values map[string]any, path string) (DeleteCheck, error) {
	sc, err := script.ResolveInputs(sc, values)
	if err != nil {
		return DeleteCheck{}, err
	}
	ops, fp := script.DeleteOps(sc)
	return DeleteCheck{Ops: ops, Fingerprint: fp, Confirmed: fp == "" || store.DeleteConfirmed(deleteKey(path, sc.Name), fp)}, nil
}

// allowDelete – lần chạy được xoá không: không có thao tác xoá → false (không cần); có thì phải đã đồng ý, hoặc confirm
// (RunOptions.ConfirmDelete) đúng dấu vân tay lúc này – khi đó ghi nhớ để lần sau không hỏi.
func allowDelete(sc model.Script, path, confirm string) (bool, error) {
	_, fp := script.DeleteOps(sc)
	if fp == "" {
		return false, nil
	}
	key := deleteKey(path, sc.Name)
	if store.DeleteConfirmed(key, fp) {
		return true, nil
	}
	if confirm != fp {
		return false, errDeleteNotConfirmed
	}
	return true, store.ConfirmDelete(key, fp)
}

func deleteKey(path, name string) string {
	if path != "" {
		return path
	}
	return "name:" + name
}
