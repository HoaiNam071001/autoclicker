//go:build probe

package main

import (
	"context"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// evalJS – chạy JS trong trang, trả giá trị (Promise thì chờ).
func evalJS[T any](ctx context.Context, js string) (T, error) {
	return chromedp.Run(ctx, chromedp.Evaluate[T](js, chromedp.EvalAwaitPromise))
}

func windowOf(ctx context.Context) (browser.GetWindowForTargetResult, error) {
	return chromedp.Run(ctx, func(ctx context.Context, t *chromedp.Target) (browser.GetWindowForTargetResult, error) {
		return cdp.Call(ctx, t, browser.GetWindowForTarget, browser.GetWindowForTargetParams{})
	})
}

func setWindowPos(ctx context.Context, left, top int64) error {
	w, err := windowOf(ctx)
	if err != nil {
		return err
	}
	return chromedp.Do(ctx, chromedp.Func(func(ctx context.Context, t *chromedp.Target) error {
		_, err := cdp.Call(ctx, t, browser.SetWindowBounds, browser.SetWindowBoundsParams{
			WindowID: w.WindowID, Bounds: &browser.Bounds{Left: &left, Top: &top},
		})
		return err
	}))
}

func bringTab(ctx context.Context) error {
	return chromedp.Do(ctx, chromedp.Func(func(ctx context.Context, t *chromedp.Target) error {
		_, err := cdp.Call(ctx, t, page.BringToFront, cdp.Empty{})
		return err
	}))
}
