package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"autoclicker/internal/model"
	"autoclicker/internal/web"
)

// Node web (bước 25d): tìm phần tử qua web.Browser (DOM), bấm bằng chuột thật như node click.

// errNoBrowserOpt – engine không có trình duyệt (Thử node / Thử tìm chưa hỗ trợ node web).
var errNoBrowserOpt = &web.BrowserError{Msg: "lần chạy này không điều khiển trình duyệt"}

// webQuery – selector đã thay biến (dòng trống bỏ qua; đã kiểm tra ở ValidateScript), kèm dòng gốc để báo lại.
func webQuery(s model.Step, click bool) (web.Query, []string, error) {
	q := web.Query{Tab: s.Tab, Scroll: click, CheckOnly: s.Type == model.StepWebWait && s.WebGone}
	var lines []string
	for _, line := range s.Selectors {
		if strings.TrimSpace(line) == "" {
			continue
		}
		sel, err := model.ParseSelector(line) // biến thay vào có thể làm trống giá trị
		if err != nil {
			return q, nil, &invalidErr{err.Error()}
		}
		q.Selectors = append(q.Selectors, sel)
		lines = append(lines, strings.TrimSpace(line))
	}
	return q, lines, nil
}

// webFind chờ tới khi phần tử hiện (gone: mất) hoặc hết TimeoutMs.
func (e *Engine) webFind(ctx context.Context, s model.Step, click bool) (web.Found, []string, error) {
	if e.browser == nil {
		return web.Found{}, nil, errNoBrowserOpt
	}
	q, lines, err := webQuery(s, click)
	if err != nil {
		return web.Found{}, nil, err
	}
	gone := q.CheckOnly
	var last web.Found
	_, err = e.waitLoop(ctx, s, "web", func() (bool, int, error) {
		f, err := e.browser.Find(ctx, q)
		if err != nil {
			return false, 0, err
		}
		last = f
		if gone {
			return !f.OK && !f.Loading && !f.NoTab, 0, nil
		}
		return f.OK, 0, nil
	})
	if errors.Is(err, errTimeout) {
		return last, lines, &notFoundErr{webMissMsg(s, last, lines, gone)}
	}
	return last, lines, err
}

// webMissMsg – vì sao hết giờ: không có tab, bị che, đang ẩn, hay không có.
func webMissMsg(s model.Step, f web.Found, lines []string, gone bool) string {
	after := "sau " + fmtSeconds(s.TimeoutMs)
	if gone {
		return fmt.Sprintf("phần tử “%s” vẫn còn %s", lines[max(f.Selector, 0)], after)
	}
	switch {
	case f.NoTab:
		return fmt.Sprintf("không thấy tab nào có URL / tiêu đề chứa “%s” %s", s.Tab, after)
	case f.Covered != "":
		return fmt.Sprintf("thấy phần tử “%s” nhưng bị %s che %s", lines[f.Selector], f.Covered, after)
	case f.Selector >= 0:
		return fmt.Sprintf("thấy phần tử “%s” nhưng chưa bấm được (nằm ngoài khung trang) %s", lines[f.Selector], after)
	case f.Hidden > 0:
		return fmt.Sprintf("có %d phần tử khớp “%s” nhưng đang ẩn %s", f.Hidden, strings.Join(lines, "”, “"), after)
	}
	return fmt.Sprintf("không thấy phần tử “%s” %s (tab %s)", strings.Join(lines, "”, “"), after, f.Tab)
}

func webOutput(f web.Found, lines []string) map[string]any {
	c := hit{Box: f.Rect}.center()
	return map[string]any{"x": c.X, "y": c.Y, "w": f.Rect.W, "h": f.Rect.H, "text": f.Text, "selector": lines[f.Selector]}
}

func (e *Engine) webWait(ctx context.Context, s model.Step) (map[string]any, error) {
	f, lines, err := e.webFind(ctx, s, false)
	if err != nil || s.WebGone {
		return nil, err
	}
	return webOutput(f, lines), nil
}

// webClick click Clicks lần vào tâm phần tử (+ Dx / Dy).
func (e *Engine) webClick(ctx context.Context, s model.Step) (map[string]any, error) {
	f, lines, err := e.webFind(ctx, s, true)
	if err != nil {
		return nil, err
	}
	out := webOutput(f, lines)
	c := hit{Box: f.Rect}.center()
	n := max(s.Clicks, 1)
	for i := 0; i < n; i++ {
		if err := e.gate.wait(ctx); err != nil {
			return nil, err
		}
		x, y := e.jitter(c.X+s.Dx, s.JitterPx), e.jitter(c.Y+s.Dy, s.JitterPx)
		if err := e.move(x, y); err != nil {
			return nil, err
		}
		if err := e.click(s.Button, false); err != nil {
			return nil, err
		}
		out["clickX"], out["clickY"] = x, y
		if i < n-1 {
			if err := Sleep(ctx, e.jitterMs(s.ClickInterval, s.JitterMs)); err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

// openInBrowser – node Mở với WebBrowser: mở bằng Chrome của app.
func (e *Engine) openInBrowser(ctx context.Context, url string) error {
	if e.browser == nil {
		return errNoBrowserOpt
	}
	return e.browser.Open(ctx, url)
}
