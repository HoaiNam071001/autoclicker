package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"autoclicker/internal/model"
	"autoclicker/internal/session"
	"autoclicker/internal/store"
)

// Chạy kịch bản từ app: mọi lần chạy (nút Chạy, API, CLI, hẹn giờ) đi qua a.sess (internal/session) – ở đây chỉ còn
// method cho frontend. Event của phiên: xem internal/session/exec.go.

// emit – Emitter của session: gửi Wails + SSE của API local.
func (a *App) emit(event string, data any) {
	if event == "log:error" {
		runtime.LogError(a.ctx, data.(string))
		return
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, event, data)
	}
	a.hub.Publish(event, data)
}

// ValidateScript áp giá trị input rồi kiểm tra như lúc chạy (frontend gọi trước khi thu cửa sổ thành mini bar).
func (a *App) ValidateScript(sc model.Script, values map[string]any) error {
	return session.Validate(sc, values)
}

// RunScript chạy kịch bản (session.Start) rồi chờ tới khi xong / bị dừng (StopScript, F8, góc màn hình – kể cả lúc
// đang đếm ngược) / lỗi. Bị dừng không tính là lỗi. Lần chạy lưu vào lịch sử.
func (a *App) RunScript(sc model.Script, values map[string]any, opt session.Options) error {
	id, err := a.sess.Start(sc, values, opt)
	if err != nil {
		return err
	}
	_, err = a.sess.Wait(context.Background(), id)
	return err
}

// StopScript dừng kịch bản đang chạy (không làm gì nếu đang không chạy).
func (a *App) StopScript() { a.sess.Stop("") }

// PauseScript tạm dừng trước bước / lần click kế tiếp. Trả false nếu không có gì để tạm dừng.
func (a *App) PauseScript() bool { return a.sess.Pause() }

// StepScript – debug: chạy node kế rồi dừng lại (F10). Trả false nếu không có gì đang chạy.
func (a *App) StepScript() bool { return a.sess.StepNext() }

// ResumeScript chạy tiếp sau khi tạm dừng (kể cả đang dừng vì debug – tắt chạy từng node).
func (a *App) ResumeScript() bool { return a.sess.Resume() }

// RunStatus – trạng thái lần chạy runId ("" = đang chạy / gần nhất).
func (a *App) RunStatus(runID string) (session.Status, error) { return a.sess.Status(runID) }

// Lịch sử chạy (bước 20 F5): mỗi lần chạy ghi một model.RunRecord + ảnh lúc node thất bại qua store.SaveRun
// (session); ListRuns / GetRun / RunShot đọc lại cho màn Lịch sử.

// ListRuns – lịch sử chạy của kịch bản (tên + đường dẫn file như lúc chạy), mới nhất trước.
func (a *App) ListRuns(name, path string) []model.RunSummary { return store.RunsOf(name, path) }

// GetRun – toàn bộ một lần chạy (nhật ký từng node).
func (a *App) GetRun(name, path, id string) (model.RunRecord, error) {
	return store.GetRun(name, path, id)
}

// RunShot – ảnh màn hình lúc node thất bại (StepResult.Shot) dạng data URL.
func (a *App) RunShot(name, path, id, shot string) (string, error) {
	return store.RunShot(name, path, id, shot)
}
