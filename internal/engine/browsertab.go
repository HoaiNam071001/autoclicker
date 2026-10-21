package engine

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"ghosthand/internal/browser"
	"ghosthand/internal/launch"
	"ghosthand/internal/model"
)

// Node Tab trình duyệt (bước 43 B1): đọc / liệt kê / chuyển / mở / đóng / tải lại tab của trình duyệt thường qua
// browser.Backend (e.tabs). Tab cần làm: của node Mở TabFrom (theo id tab node đó ghi lại) > khớp Tab > tab đang hiện
// của cửa sổ trên cùng. Không thấy thì chờ tới TimeoutMs.

func (e *Engine) browserTab(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	switch s.TabOp {
	case model.TabNew:
		url := strings.TrimSpace(s.URL)
		if url != "" {
			if url = launch.BrowserURL(url); url == "" {
				return nil, &invalidErr{fmt.Sprintf("%q không phải URL / đường dẫn file", s.URL)}
			}
		}
		t, err := e.tabs.NewTab(ctx, s.Browser, url)
		if err != nil {
			return nil, err
		}
		return tabOutput(t), nil
	case model.TabList:
		tabs, err := e.tabs.Tabs(ctx, s.Browser)
		if err != nil {
			return nil, err
		}
		idx, err := model.MatchTabs(tabs, s.Tab, s.TabMatch)
		if err != nil {
			return nil, &invalidErr{err.Error()}
		}
		items := make([]any, len(idx))
		for k, i := range idx {
			items[k] = tabItem(tabs[i])
		}
		return map[string]any{"items": items, "count": len(items)}, nil
	case model.TabWaitClose:
		err := e.pollUntil(ctx, s, "tab", func() (bool, error) {
			_, ok, err := e.findTab(ctx, s)
			return !ok, err
		}, func() error {
			return &timeoutErr{msg: fmt.Sprintf("tab %s vẫn mở sau %s", tabName(s), fmtSeconds(s.TimeoutMs))}
		})
		return map[string]any{}, err
	}
	t, err := e.waitTab(ctx, s)
	if err != nil {
		return nil, err
	}
	return e.tabAction(ctx, s, t)
}

// tabAction làm thao tác của node trên tab t đã tìm thấy.
func (e *Engine) tabAction(ctx context.Context, s model.Step, t model.TabInfo) (map[string]any, error) {
	switch s.TabOp {
	case model.TabSwitch:
		if err := e.tabs.Activate(ctx, t); err != nil {
			return nil, err
		}
		t.Active = true
	case model.TabClose:
		if err := e.tabs.Close(ctx, t); err != nil {
			return nil, err
		}
	case model.TabCloseOthers:
		n, err := e.closeOthers(ctx, s, t)
		out := tabOutput(t)
		out["closed"] = n
		return out, err
	case model.TabReload, model.TabBack, model.TabForward:
		if err := e.tabs.Nav(ctx, t, s.TabOp); err != nil {
			return nil, err
		}
		t = e.rereadTab(ctx, t)
	}
	return tabOutput(t), nil
}

// closeOthers đóng mọi tab khác t trong cùng cửa sổ; tab tự đóng giữa chừng thì bỏ qua.
func (e *Engine) closeOthers(ctx context.Context, s model.Step, t model.TabInfo) (int, error) {
	tabs, err := e.tabs.Tabs(ctx, t.Browser)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, x := range tabs {
		if x.Window != t.Window || x.ID == t.ID {
			continue
		}
		if err := e.gate.wait(ctx); err != nil {
			return n, err
		}
		if err := e.tabs.Close(ctx, x); err != nil && !errors.Is(err, browser.ErrTabGone) {
			return n, err
		}
		n++
	}
	return n, nil
}

// waitTab chờ tới khi có tab cần làm. Hết giờ: chờ mở → timeout, thao tác khác → not_found.
func (e *Engine) waitTab(ctx context.Context, s model.Step) (model.TabInfo, error) {
	var t model.TabInfo
	var tabs []model.TabInfo
	err := e.pollUntil(ctx, s, "tab", func() (bool, error) {
		var ok bool
		var err error
		t, ok, err = e.findTab(ctx, s)
		if !ok && err == nil {
			tabs, _ = e.tabs.Tabs(ctx, s.Browser) // để báo các tab đang mở
		}
		return ok, err
	}, func() error {
		msg := fmt.Sprintf("không thấy tab %s sau %s", tabName(s), fmtSeconds(s.TimeoutMs))
		if s.TabOp == model.TabWaitOpen {
			return &timeoutErr{msg: msg, detail: tabMiss(tabs)}
		}
		return &notFoundErr{msg: msg, detail: tabMiss(tabs)}
	})
	return t, err
}

// findTab – tab node cần làm lúc này (ok = false: chưa có).
func (e *Engine) findTab(ctx context.Context, s model.Step) (model.TabInfo, bool, error) {
	tabs, err := e.tabs.Tabs(ctx, s.Browser)
	if err != nil {
		return model.TabInfo{}, false, err
	}
	if s.TabFrom != "" {
		return e.openedTab(s, tabs)
	}
	i := model.ActiveTab(tabs)
	if s.Tab != "" || s.TabOp == model.TabWaitOpen || s.TabOp == model.TabWaitClose {
		idx, err := model.MatchTabs(tabs, s.Tab, s.TabMatch)
		if err != nil {
			return model.TabInfo{}, false, &invalidErr{err.Error()}
		}
		i = -1
		if len(idx) > 0 {
			i = idx[0]
		}
	}
	if i < 0 {
		return model.TabInfo{}, false, nil
	}
	return tabs[i], true, nil
}

// openedTab – tab node Mở TabFrom đã mở (output tabId; tab đó đóng rồi → URL tabUrl trong cùng cửa sổ, nhiều tab thì tab
// mới nhất). Node Mở chưa chạy / không ghi được tab → invalid_input (không đoán tab khác – tránh đóng nhầm).
func (e *Engine) openedTab(s model.Step, tabs []model.TabInfo) (model.TabInfo, bool, error) {
	r, ok := e.results[s.TabFrom]
	src := e.byID[s.TabFrom]
	if !ok || r.Output["tabId"] == nil {
		return model.TabInfo{}, false, &invalidErr{fmt.Sprintf("node %s chưa mở tab nào trong vòng này (cần chạy trước node này, mở URL bằng trình duyệt thường)", src.Ref)}
	}
	id, _ := strconv.ParseInt(fmtValue(r.Output["tabId"]), 10, 64)
	win, _ := strconv.ParseInt(fmtValue(r.Output["tabWindow"]), 10, 64)
	url := fmtValue(r.Output["tabUrl"])
	if i := slices.IndexFunc(tabs, func(t model.TabInfo) bool { return t.ID == id }); i >= 0 {
		return tabs[i], true, nil
	}
	best := -1
	for i, t := range tabs {
		if t.Window == win && t.URL == url && (best < 0 || t.ID > tabs[best].ID) {
			best = i
		}
	}
	if best < 0 {
		return model.TabInfo{}, false, nil
	}
	return tabs[best], true, nil
}

// rereadTab – thông tin mới của tab t (sau tải lại / quay lại – URL, tiêu đề đổi); không đọc được thì giữ t.
func (e *Engine) rereadTab(ctx context.Context, t model.TabInfo) model.TabInfo {
	if tabs, err := e.tabs.Tabs(ctx, t.Browser); err == nil {
		if i := slices.IndexFunc(tabs, func(x model.TabInfo) bool { return x.ID == t.ID }); i >= 0 {
			return tabs[i]
		}
	}
	return t
}

// opensTab – có node Tab trình duyệt / Đi tới URL nào làm với tab node Mở id mở không (openStep chỉ đọc tab khi cần – tránh hỏi
// quyền Automation khi không dùng).
func (e *Engine) opensTab(id string) bool {
	return slices.ContainsFunc(e.sc.Steps, func(s model.Step) bool { return (s.Type == model.StepBrowserTab || s.Type == model.StepBrowserGo) && s.TabFrom == id })
}

// putOpenedTab ghi tab đang hiện của trình duyệt ở trước (tab node Mở vừa mở) vào output của node Mở; không đọc được
// thì bỏ – node Tab trình duyệt sau sẽ báo.
func (e *Engine) putOpenedTab(ctx context.Context, out map[string]any) {
	tabs, err := e.tabs.Tabs(ctx, "")
	if i := model.ActiveTab(tabs); err == nil && i >= 0 {
		t := tabs[i]
		out["tabId"], out["tabWindow"], out["tabUrl"], out["tabTitle"], out["tabBrowser"] = t.ID, t.Window, t.URL, t.Title, t.Browser
	}
}

func tabOutput(t model.TabInfo) map[string]any {
	return map[string]any{"url": t.URL, "title": t.Title, "browser": t.Browser, "tabId": t.ID, "window": t.Window, "index": t.Index}
}

func tabItem(t model.TabInfo) map[string]any {
	return map[string]any{"title": t.Title, "url": t.URL, "browser": t.Browser, "window": t.Window, "active": t.Active, "tabId": t.ID}
}

// tabName – tab cần làm để báo lỗi.
func tabName(s model.Step) string {
	switch {
	case s.TabFrom != "":
		return "do node trước mở"
	case s.Tab != "":
		return fmt.Sprintf("“%s”", s.Tab)
	}
	return "đang hiện"
}

// tabMiss – detail khi không thấy tab: tối đa 10 tab đang mở để AI / người dùng sửa chuỗi tìm.
func tabMiss(tabs []model.TabInfo) map[string]any {
	open := []any{}
	for _, t := range tabs[:min(len(tabs), 10)] {
		open = append(open, map[string]any{"title": t.Title, "url": t.URL})
	}
	return map[string]any{"openTabs": open, "tabCount": len(tabs)}
}
