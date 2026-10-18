package web

import (
	"context"

	"github.com/chromedp/chromedp"
)

// Focused – webType: sau khi click, ô đang focus (xuống cả iframe cùng nguồn) phải là ô nhập và là phần tử đích, con /
// cha của nó hoặc ô của label đích – không thì gõ sẽ rơi vào chỗ khác.
func (c *Chrome) Focused(ctx context.Context) (bool, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.last == nil {
		return false, "chưa tìm phần tử", nil
	}
	r, err := run(ctx, c.last.ctx, "đọc ô đang focus", chromedp.Evaluate[struct {
		OK bool
		At string
	}](pageJS(focusJS)))
	return r.OK, r.At, err
}

const focusJS = `() => {
  let a = document.activeElement;
  while (a && a.contentDocument && a.contentDocument.activeElement) a = a.contentDocument.activeElement; // vào iframe
  const t = window.__acTarget;
  const editable = a && (a.isContentEditable || a.matches('input:not([type=button]):not([type=submit]):not([type=checkbox]):not([type=radio]), textarea, select'));
  const label = t && (t.control || (t.closest && t.closest('label') && t.closest('label').control));
  const ok = !!(editable && t && (a === t || t.contains(a) || a.contains(t) || a === label));
  return { OK: ok, At: ok ? '' : !a || a === document.body ? 'không có ô nào' : describe(a) };
}`
