package main

import (
	"errors"

	"autoclicker/internal/model"
	"autoclicker/internal/script"
	"autoclicker/internal/store"
)

// Xác nhận thao tác nguy hiểm (script.DangerOps – xoá file, buộc thoát app, đóng cửa sổ): kịch bản có thao tác nguy hiểm
// thì lần chạy đầu (và mỗi khi đích đổi) người dùng phải đồng ý. Chặn ở Go nên hẹn giờ / CLI / MCP (bước 22, 31) cũng
// không chạy được khi chưa đồng ý trong app.

// DangerCheck – thao tác nguy hiểm của kịch bản và đã được đồng ý chưa (frontend hỏi trước khi chạy).
type DangerCheck struct {
	Ops         []script.DangerOp `json:"ops"`
	Fingerprint string            `json:"fingerprint"`
	Confirmed   bool              `json:"confirmed"`
}

// errDangerNotConfirmed – chạy kịch bản có thao tác nguy hiểm chưa được đồng ý.
var errDangerNotConfirmed = errors.New("kịch bản có thao tác cần xác nhận (xoá file / thoát app / đóng cửa sổ) – bấm Chạy trong app để xem và đồng ý")

// DangerCheck áp input như lúc chạy rồi liệt kê thao tác nguy hiểm. path: file kịch bản ("" = chưa lưu, theo tên).
func (a *App) DangerCheck(sc model.Script, values map[string]any, path string) (DangerCheck, error) {
	sc, err := script.ResolveInputs(sc, values)
	if err != nil {
		return DangerCheck{}, err
	}
	ops, fp := script.DangerOps(sc)
	return DangerCheck{Ops: ops, Fingerprint: fp, Confirmed: fp == "" || store.DangerConfirmed(dangerKey(path, sc.Name), fp)}, nil
}

// allowDanger – lần chạy được làm thao tác nguy hiểm không: không có → false (không cần); có thì phải đã đồng ý, hoặc
// confirm (RunOptions.ConfirmDanger) đúng dấu vân tay lúc này – khi đó ghi nhớ để lần sau không hỏi.
func allowDanger(sc model.Script, path, confirm string) (bool, error) {
	_, fp := script.DangerOps(sc)
	if fp == "" {
		return false, nil
	}
	key := dangerKey(path, sc.Name)
	if store.DangerConfirmed(key, fp) {
		return true, nil
	}
	if confirm != fp {
		return false, errDangerNotConfirmed
	}
	return true, store.ConfirmDanger(key, fp)
}

func dangerKey(path, name string) string {
	if path != "" {
		return path
	}
	return "name:" + name
}
