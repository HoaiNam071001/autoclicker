package web

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"

	"ghosthand/internal/model"
)

// Chạy nền (D4b): bấm / gõ / lăn qua CDP Input.* thay chuột / phím thật – không đưa tab lên trước, không đo toạ độ
// bằng chuột, toạ độ là CSS px trong khung nhìn trang (find.js đã cộng khung iframe).

// findBackground – phần còn lại của Find khi Query.Background: không bringToFront / calibrate, khung theo toạ độ trang.
func (c *Chrome) findBackground(ctx context.Context, t *tabConn, q Query, r findResult, f Found) (Found, error) {
	if q.Scroll { // sắp bấm: trang phải nghĩ mình đang được focus (ô nhận focus, :focus, sự kiện blur) dù cửa sổ ở dưới
		if err := t.emulateFocus(ctx); err != nil {
			return Found{}, err
		}
		var err error
		if r, f, err = steady(ctx, t.ctx, q, r, f.Tab); err != nil || !f.OK {
			return f, err
		}
		c.last = t
	}
	f.Rect = model.Rect{X: int(math.Round(r.X)), Y: int(math.Round(r.Y)), W: int(math.Round(r.W)), H: int(math.Round(r.H))}
	return f, nil
}

func (t *tabConn) emulateFocus(ctx context.Context) error {
	if t.focusEmu {
		return nil
	}
	_, err := run(ctx, t.ctx, "giả lập focus cho tab", func(ctx context.Context, tg *chromedp.Target) (cdp.Empty, error) {
		return cdp.Call(ctx, tg, emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true})
	})
	t.focusEmu = err == nil
	return err
}

// bgTab – tab của lần Find (Scroll) gần nhất; chưa có → lỗi (engine luôn Find trước).
func (c *Chrome) bgTab() (*tabConn, error) {
	if c.last == nil {
		return nil, &BrowserError{"chạy nền: chưa tìm phần tử nào trên tab"}
	}
	return c.last, nil
}

func dispatch[P any](ctx context.Context, t *tabConn, what string, cmd cdp.Command[P, cdp.Empty], p P) error {
	_, err := run(ctx, t.ctx, what, func(ctx context.Context, tg *chromedp.Target) (cdp.Empty, error) {
		return cdp.Call(ctx, tg, cmd, p)
	})
	return err
}

// bgHoldMs – giữ nút giữa nhấn và nhả (trang nghe mousedown / mouseup tách nhau như chuột thật).
const bgHoldMs = 30 * time.Millisecond

func (c *Chrome) BgClick(ctx context.Context, x, y int, button string, count int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, err := c.bgTab()
	if err != nil {
		return err
	}
	btn, mask := mouseButton(button)
	ev := func(typ input.DispatchMouseEventType, b input.MouseButton, buttons int64) input.DispatchMouseEventParams {
		return input.DispatchMouseEventParams{Type: typ, X: float64(x), Y: float64(y), Button: b, Buttons: buttons, ClickCount: int64(count)}
	}
	if err := dispatch(ctx, t, "chạy nền: di chuột", input.DispatchMouseEvent, ev(input.DispatchMouseEventTypeMouseMoved, input.MouseButtonNone, 0)); err != nil {
		return err
	}
	if err := dispatch(ctx, t, "chạy nền: nhấn chuột", input.DispatchMouseEvent, ev(input.DispatchMouseEventTypeMousePressed, btn, mask)); err != nil {
		return err
	}
	waitErr := sleepCtx(ctx, bgHoldMs)
	// luôn nhả, kể cả khi vừa bấm Dừng – không thì trang kẹt ở trạng thái đang nhấn (kéo chọn chữ…)
	err = dispatch(context.WithoutCancel(ctx), t, "chạy nền: nhả chuột", input.DispatchMouseEvent, ev(input.DispatchMouseEventTypeMouseReleased, btn, 0))
	return errors.Join(waitErr, err)
}

func mouseButton(b string) (input.MouseButton, int64) {
	switch b {
	case "right":
		return input.MouseButtonRight, 2
	case "middle":
		return input.MouseButtonMiddle, 4
	}
	return input.MouseButtonLeft, 1
}

// wheelNotch – một nấc lăn chạy nền (CSS px) – như một nấc bánh xe của Chrome trên Windows.
const wheelNotch = 100

func (c *Chrome) BgWheel(ctx context.Context, x, y int, dir string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, err := c.bgTab()
	if err != nil {
		return err
	}
	dx, dy := wheelDelta(dir)
	return dispatch(ctx, t, "chạy nền: lăn chuột", mouseWheel, wheelParams{Type: input.DispatchMouseEventTypeMouseWheel, X: float64(x), Y: float64(y), DeltaX: dx, DeltaY: dy})
}

// wheelParams – Input.dispatchMouseEvent kiểu mouseWheel: Chrome đòi đủ deltaX lẫn deltaY, còn
// input.DispatchMouseEventParams bỏ trường bằng 0 (omitzero).
type wheelParams struct {
	Type   input.DispatchMouseEventType `json:"type"`
	X      float64                      `json:"x"`
	Y      float64                      `json:"y"`
	DeltaX float64                      `json:"deltaX"`
	DeltaY float64                      `json:"deltaY"`
}

var mouseWheel = cdp.Command[wheelParams, cdp.Empty]{Method: input.CommandDispatchMouseEvent}

func wheelDelta(dir string) (float64, float64) {
	switch dir {
	case "up":
		return 0, -wheelNotch
	case "left":
		return -wheelNotch, 0
	case "right":
		return wheelNotch, 0
	}
	return 0, wheelNotch
}

// Phím chạy nền: chọn hết bằng lệnh soạn thảo (không phụ thuộc Cmd / Ctrl theo hệ điều hành), Backspace, Enter.
var (
	keySelectAll = input.DispatchKeyEventParams{Type: input.DispatchKeyEventTypeRawKeyDown, Key: "a", Code: "KeyA", WindowsVirtualKeyCode: 65, Commands: []string{"selectAll"}}
	keyBackspace = input.DispatchKeyEventParams{Type: input.DispatchKeyEventTypeRawKeyDown, Key: "Backspace", Code: "Backspace", WindowsVirtualKeyCode: 8, Commands: []string{"deleteBackward"}}
	keyEnter     = input.DispatchKeyEventParams{Type: input.DispatchKeyEventTypeKeyDown, Key: "Enter", Code: "Enter", WindowsVirtualKeyCode: 13, Text: "\r"}
)

func (c *Chrome) BgType(ctx context.Context, text string, clear, enter bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, err := c.bgTab()
	if err != nil {
		return err
	}
	if clear {
		if err := pressKey(ctx, t, keySelectAll); err != nil {
			return err
		}
		if err := pressKey(ctx, t, keyBackspace); err != nil {
			return err
		}
	}
	if text != "" {
		if err := dispatch(ctx, t, "chạy nền: chèn chữ", input.InsertText, input.InsertTextParams{Text: text}); err != nil {
			return err
		}
	}
	if enter {
		return pressKey(ctx, t, keyEnter)
	}
	return nil
}

// pressKey gửi phím xuống rồi lên (keyUp cùng phím, không lệnh / chữ).
func pressKey(ctx context.Context, t *tabConn, down input.DispatchKeyEventParams) error {
	if err := dispatch(ctx, t, "chạy nền: bấm phím "+down.Key, input.DispatchKeyEvent, down); err != nil {
		return err
	}
	up := input.DispatchKeyEventParams{Type: input.DispatchKeyEventTypeKeyUp, Key: down.Key, Code: down.Code, WindowsVirtualKeyCode: down.WindowsVirtualKeyCode}
	return dispatch(ctx, t, "chạy nền: nhả phím "+down.Key, input.DispatchKeyEvent, up)
}
