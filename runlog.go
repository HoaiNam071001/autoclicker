package main

import (
	"fmt"

	"autoclicker/internal/model"
	"autoclicker/internal/store"
	"autoclicker/internal/vision"
)

// Lịch sử chạy (bước 20 F5): mỗi lần chạy lưu một file runs/<kịch bản>/<id>.json (RunRecord – StepResult từng node,
// input đã dùng, kết quả cuối) + ảnh màn hình <id>-<n>.png lúc node thất bại. Giữ maxRuns lần gần nhất mỗi kịch bản.

const (
	maxShots = 5    // ảnh lúc lỗi tối đa mỗi lần chạy (vòng thử lại qua cổng fail không chụp mãi)
	shotMaxW = 1600 // ảnh rộng hơn thì thu nhỏ (Retina) cho file nhẹ
)

// Nguồn của lần chạy (RunOptions.Source).
const (
	SourceUser     = "user"
	SourceSchedule = "schedule"
	SourceCLI      = "cli"
	SourceAI       = "ai"
)

// RunOptions – cách chạy kịch bản (RunScript).
type RunOptions struct {
	Countdown int    `json:"countdown"`        // giây đếm ngược trước khi chạy
	Path      string `json:"path,omitempty"`   // file kịch bản – gom lịch sử theo kịch bản (chưa lưu: theo tên)
	Source    string `json:"source,omitempty"` // user (mặc định) | schedule | cli | ai
	From      string `json:"from,omitempty"`   // debug: chạy một vòng bắt đầu từ node này thay cho start
	Step      bool   `json:"step,omitempty"`   // debug: dừng trước mỗi node, Tiếp / F10 chạy node kế
}

// shoot chụp màn chứa vùng kịch bản lúc node thất bại (tối đa maxShots mỗi lần chạy); lỗi chụp thì bỏ qua.
func (e *engine) shoot(res *model.StepResult) {
	if e.shots == nil || len(e.shots) >= maxShots || res.Error == nil || res.Error.Code == model.ErrPermission || res.Error.Code == model.ErrEnd {
		return
	}
	img, err := e.scr.Capture(e.searchArea(model.Step{SearchArea: "screen"}))
	if err != nil {
		return
	}
	res.Shot = fmt.Sprintf("%d.png", len(e.shots)+1)
	e.shots[res.Shot] = vision.Shrink(img, shotMaxW)
}

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
