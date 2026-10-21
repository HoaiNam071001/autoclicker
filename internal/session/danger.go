package session

import (
	"errors"

	"ghosthand/internal/model"
	"ghosthand/internal/script"
	"ghosthand/internal/store"
)

// Thao tác nguy hiểm (script.DangerOps – xoá file, buộc thoát app, đóng cửa sổ): chặn ở đây nên hẹn giờ / CLI / API /
// MCP cũng không chạy được khi người dùng chưa đồng ý trong app.

// ErrDangerNotConfirmed – chạy kịch bản có thao tác nguy hiểm chưa được đồng ý.
var ErrDangerNotConfirmed = errors.New("kịch bản có thao tác cần xác nhận (xoá file / thoát app / đóng cửa sổ) – bấm Chạy trong app để xem và đồng ý")

// AllowDanger – lần chạy được làm thao tác nguy hiểm không: không có → false (không cần); có thì phải đã đồng ý, hoặc
// confirm (Options.ConfirmDanger) đúng dấu vân tay lúc này – khi đó ghi nhớ để lần sau không hỏi.
func AllowDanger(sc model.Script, path, confirm string) (bool, error) {
	_, fp := script.DangerOps(sc)
	if fp == "" {
		return false, nil
	}
	key := DangerKey(path, sc.Name)
	if store.DangerConfirmed(key, fp) {
		return true, nil
	}
	if confirm != fp {
		return false, ErrDangerNotConfirmed
	}
	return true, store.ConfirmDanger(key, fp)
}

// DangerKey – khoá ghi nhớ đồng ý: đường dẫn file, chưa lưu thì theo tên.
func DangerKey(path, name string) string {
	if path != "" {
		return path
	}
	return "name:" + name
}
