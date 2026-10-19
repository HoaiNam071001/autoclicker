package devicetest

import (
	"errors"
	"fmt"
	"sync"

	"autoclicker/internal/model"
)

// FakeWindows – cửa sổ giả (bước 26): thao tác đổi Wins và ghi Log ("focus 2", "bounds 2 0,0 800x600"…).
// Work = vùng làm việc trả cho mọi màn.
type FakeWindows struct {
	mu   sync.Mutex
	Wins []model.WindowInfo
	Work model.Rect
	Log  []string
}

var errGone = errors.New("cửa sổ không còn mở")

func (f *FakeWindows) List() ([]model.WindowInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]model.WindowInfo(nil), f.Wins...), nil
}

// with chạy fn trên cửa sổ id (đã khoá) và ghi log; không có → lỗi như bản thật.
func (f *FakeWindows) with(id uint64, log string, fn func(w *model.WindowInfo)) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := range f.Wins {
		if f.Wins[i].ID == id {
			f.Log = append(f.Log, log)
			fn(&f.Wins[i])
			return nil
		}
	}
	return errGone
}

func (f *FakeWindows) Focus(id uint64) error {
	return f.with(id, fmt.Sprintf("focus %d", id), func(w *model.WindowInfo) { w.Minimized = false })
}

func (f *FakeWindows) SetBounds(id uint64, r model.Rect) error {
	return f.with(id, fmt.Sprintf("bounds %d %d,%d %dx%d", id, r.X, r.Y, r.W, r.H), func(w *model.WindowInfo) {
		w.Box, w.Minimized = r, false
	})
}

func (f *FakeWindows) SetState(id uint64, state string) error {
	return f.with(id, fmt.Sprintf("%s %d", state, id), func(w *model.WindowInfo) { w.Minimized = state == "minimize" })
}

func (f *FakeWindows) Close(id uint64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, w := range f.Wins {
		if w.ID == id {
			f.Log = append(f.Log, fmt.Sprintf("close %d", id))
			f.Wins = append(f.Wins[:i], f.Wins[i+1:]...)
			return nil
		}
	}
	return errGone
}

func (f *FakeWindows) WorkArea(model.Rect) model.Rect { return f.Work }

// FakeApps – app giả: Quit / ForceQuit bỏ app khỏi Apps và ghi Log ("quit 10", "force 10").
type FakeApps struct {
	mu   sync.Mutex
	Apps []model.AppInfo
	Log  []string
}

func (f *FakeApps) List() ([]model.AppInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]model.AppInfo(nil), f.Apps...), nil
}

func (f *FakeApps) Quit(pid int) error      { return f.remove(pid, "quit") }
func (f *FakeApps) ForceQuit(pid int) error { return f.remove(pid, "force") }

func (f *FakeApps) remove(pid int, op string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, a := range f.Apps {
		if a.PID == pid {
			f.Log = append(f.Log, fmt.Sprintf("%s %d", op, pid))
			f.Apps = append(f.Apps[:i], f.Apps[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("không có tiến trình %d", pid)
}
