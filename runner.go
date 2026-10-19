package main

import (
	"context"
	"errors"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"autoclicker/internal/device"
	"autoclicker/internal/engine"
	"autoclicker/internal/model"
	"autoclicker/internal/script"
	"autoclicker/internal/store"
	"autoclicker/internal/web"
)

// Lý do dừng gửi kèm event "engine:stopped".
const (
	StopByF8     = "f8"
	StopByCorner = "corner"
)

type runHandle struct {
	cancel context.CancelFunc
	eng    *engine.Engine
}

// ValidateScript áp giá trị input rồi kiểm tra như lúc chạy (frontend gọi trước khi thu cửa sổ thành mini bar).
func (a *App) ValidateScript(sc model.Script, values map[string]any) error {
	out, err := script.ResolveInputs(sc, values)
	if err != nil {
		return err
	}
	return script.ValidateScript(out, device.MainScreen())
}

// RunScript áp giá trị input (resolveInputs) vào bản sao kịch bản, đếm ngược opt.Countdown giây rồi chạy; chỉ trả về khi chạy xong / bị dừng
// (StopScript, F8, góc màn hình – kể cả lúc đang đếm ngược) / lỗi. Bị dừng không tính là lỗi.
// Emit "engine:state" (true khi bắt đầu, false khi kết thúc), "engine:countdown" (số giây còn lại)
// "engine:progress" (Progress) và "engine:node" (StepResult mỗi node xong) trong lúc chạy, "engine:done" (RunSummary)
// khi chạy xong – lần chạy lưu vào lịch sử (runlog.go).
func (a *App) RunScript(sc model.Script, values map[string]any, opt RunOptions) error {
	screen := device.MainScreen()
	if opt.From != "" {
		sc = script.RunFrom(sc, opt.From, opt.Only)
	}
	sc, err := script.ResolveInputs(sc, values)
	if err != nil {
		return err
	}
	if err := script.ValidateScript(sc, screen); err != nil {
		return err
	}
	canDanger, err := allowDanger(sc, opt.Path, opt.ConfirmDanger)
	if err != nil {
		return err
	}
	profile, err := store.BrowserProfile()
	if err != nil {
		return err
	}
	a.mu.Lock()
	if a.run != nil {
		a.mu.Unlock()
		return errors.New("kịch bản đang chạy")
	}
	ctx, cancel := context.WithCancel(context.Background())
	in := &device.TrackingInput{Input: device.RobotgoInput{}}
	if opt.Source == "" {
		opt.Source = SourceUser
	}
	rec := model.RunRecord{RunSummary: model.RunSummary{Script: sc.Name, Source: opt.Source}, Inputs: values, Results: []model.StepResult{}}
	types := map[string]string{}
	for _, s := range sc.Steps {
		types[s.ID] = s.Type
	}
	br := web.NewChrome(profile, in.Move) // nối lười – kịch bản không có node web thì không đụng trình duyệt
	defer br.Close()
	eng := engine.New(in, sc, engine.Options{
		Screen: screen, Displays: device.Displays(), Debug: opt.Source == SourceUser, StepMode: opt.Step, KeepShots: true, Browser: br, AllowDanger: canDanger,
		OnProgress: func(p engine.Progress) { runtime.EventsEmit(a.ctx, "engine:progress", p) },
		OnNode: func(r model.StepResult) {
			rec.Add(r, types[r.NodeID])
			runtime.EventsEmit(a.ctx, "engine:node", r)
		},
	})
	h := &runHandle{cancel: cancel, eng: eng}
	a.run = h
	a.mu.Unlock()

	defer func() {
		cancel()
		a.mu.Lock()
		if a.run == h {
			a.run = nil
		}
		a.mu.Unlock()
		runtime.EventsEmit(a.ctx, "engine:state", false)
	}()
	runtime.EventsEmit(a.ctx, "engine:state", true)

	for n := opt.Countdown; n > 0; n-- {
		runtime.EventsEmit(a.ctx, "engine:countdown", n)
		if engine.Sleep(ctx, 1000) != nil {
			return nil // dừng khi đang đếm ngược
		}
	}
	runtime.EventsEmit(a.ctx, "engine:countdown", 0)

	go device.WatchCorners(ctx, device.Location, device.HardCorners(device.Displays()), in, func() { a.stopWith(StopByCorner) })

	rec.StartedAt = time.Now()
	rec.ID = store.NewRunID(rec.StartedAt)
	err = eng.Run(ctx)
	a.finishRun(&rec, eng, opt, err)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if id, miss := eng.Failed(); err != nil && id != "" {
		// frontend chọn bước lỗi; chờ ảnh hết giờ thì kèm chỗ giống nhất để xem lại
		runtime.EventsEmit(a.ctx, "engine:failed", FailedStep{StepID: id, Miss: miss})
	}
	return err
}

// finishRun ghi kết quả cuối vào lịch sử chạy rồi emit "engine:done". Lưu lỗi thì vẫn báo xong (lịch sử chỉ là phụ).
func (a *App) finishRun(rec *model.RunRecord, eng *engine.Engine, opt RunOptions, err error) {
	rec.DurationMs = time.Since(rec.StartedAt).Milliseconds()
	switch {
	case errors.Is(err, context.Canceled):
		rec.Status = model.RunStopped
	case err != nil:
		rec.Status, rec.Error = model.RunFail, err.Error()
	default:
		rec.Status = model.RunSuccess
	}
	if base, e := store.ConfigDir(); e == nil {
		if e = store.SaveRun(store.RunsDir(base, rec.Script, opt.Path), *rec, eng.Shots()); e != nil {
			runtime.LogErrorf(a.ctx, "lưu lịch sử chạy: %v", e)
		}
	}
	runtime.EventsEmit(a.ctx, "engine:done", rec.RunSummary)
}

// FailedStep – bước làm lần chạy dừng vì lỗi (event engine:failed).
type FailedStep struct {
	StepID string             `json:"stepId"`
	Miss   *engine.FindResult `json:"miss,omitempty"` // chờ ảnh hết giờ: chỗ giống nhất
}

// stopWith dừng kịch bản đang chạy và emit "engine:stopped" kèm lý do (F8, góc màn hình).
func (a *App) stopWith(reason string) {
	if a.current() == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "engine:stopped", reason)
	a.StopScript()
}

// StopScript dừng kịch bản đang chạy (không làm gì nếu đang không chạy).
func (a *App) StopScript() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.run != nil {
		a.run.cancel()
		a.run = nil
	}
}

// PauseScript tạm dừng trước bước / lần click kế tiếp. Trả false nếu không có gì để tạm dừng.
func (a *App) PauseScript() bool {
	if h := a.current(); h != nil {
		return h.eng.Pause()
	}
	return false
}

// StepScript – debug: chạy node kế rồi dừng lại (F10). Trả false nếu không có gì đang chạy.
func (a *App) StepScript() bool {
	if h := a.current(); h != nil {
		return h.eng.Step()
	}
	return false
}

// ResumeScript chạy tiếp sau khi tạm dừng (kể cả đang dừng vì debug – tắt chạy từng node).
func (a *App) ResumeScript() bool {
	if h := a.current(); h != nil {
		return h.eng.Resume()
	}
	return false
}

func (a *App) current() *runHandle {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.run
}

// Lịch sử chạy (bước 20 F5): mỗi lần chạy ghi một model.RunRecord + ảnh lúc node thất bại (engine.Shots) qua
// store.SaveRun; ListRuns / GetRun / RunShot đọc lại cho màn Lịch sử.

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
	Only      bool   `json:"only,omitempty"`   // cùng From: chỉ chạy riêng node From (RunNode cho AI – 31c)
	Step      bool   `json:"step,omitempty"`   // debug: dừng trước mỗi node, Tiếp / F10 chạy node kế
	// dấu vân tay thao tác nguy hiểm người dùng vừa đồng ý (DangerCheck) – chạy được và ghi nhớ; trống = chỉ chạy khi đã đồng ý trước đó
	ConfirmDanger string `json:"confirmDanger,omitempty"`
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
