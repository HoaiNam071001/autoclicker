// Package browsertest – trình duyệt giả cho test engine (node Tab trình duyệt).
package browsertest

import (
	"context"
	"slices"
	"strconv"
	"sync"

	"ghosthand/internal/browser"
	"ghosthand/internal/model"
)

// Fake – danh sách tab trong bộ nhớ; Calls ghi lại thao tác ("activate 3", "close 3", "reload 3", "new <url>").
type Fake struct {
	mu     sync.Mutex
	List   []model.TabInfo // cửa sổ trên cùng trước, như browser.Backend.Tabs
	Calls  []string
	Err    error         // trả cho mọi thao tác (≠ nil)
	OnList func(f *Fake) // gọi trước mỗi lần Tabs – đổi tab giữa chừng (tab mở / đóng khi đang chờ)
	nextID int64
}

var _ browser.Backend = (*Fake)(nil)

func (f *Fake) Tabs(_ context.Context, b string) ([]model.TabInfo, error) {
	if f.OnList != nil {
		f.OnList(f)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Err != nil {
		return nil, f.Err
	}
	return slices.DeleteFunc(slices.Clone(f.List), func(t model.TabInfo) bool { return b != "" && t.Browser != b }), nil
}

func (f *Fake) Activate(_ context.Context, t model.TabInfo) error {
	return f.do("activate", t, func(i int) {
		for j := range f.List {
			if f.List[j].Window == t.Window {
				f.List[j].Active = j == i
			}
		}
	})
}

func (f *Fake) NewTab(_ context.Context, b, url string) (model.TabInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Calls = append(f.Calls, "new "+url)
	if f.Err != nil {
		return model.TabInfo{}, f.Err
	}
	if len(f.List) == 0 {
		return model.TabInfo{}, browser.ErrNotRunning
	}
	if b == "" {
		b = f.List[0].Browser
	}
	f.nextID++
	t := model.TabInfo{Browser: b, ID: 1000 + f.nextID, Window: f.List[0].Window, URL: url, Title: url, Active: true}
	for j := range f.List {
		if f.List[j].Window == t.Window {
			f.List[j].Active = false
			t.Index = f.List[j].Index + 1
		}
	}
	f.List = append(f.List, t)
	return t, nil
}

func (f *Fake) Close(_ context.Context, t model.TabInfo) error {
	return f.do("close", t, func(i int) { f.List = slices.Delete(f.List, i, i+1) })
}

func (f *Fake) Nav(_ context.Context, t model.TabInfo, op string) error {
	return f.do(op, t, func(int) {})
}

// do ghi thao tác op trên tab t rồi chạy apply với chỉ số của t trong List; tab không còn → browser.ErrTabGone.
func (f *Fake) do(op string, t model.TabInfo, apply func(i int)) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Calls = append(f.Calls, op+" "+strconv.FormatInt(t.ID, 10))
	if f.Err != nil {
		return f.Err
	}
	i := slices.IndexFunc(f.List, func(x model.TabInfo) bool { return x.ID == t.ID })
	if i < 0 {
		return browser.ErrTabGone
	}
	apply(i)
	return nil
}
