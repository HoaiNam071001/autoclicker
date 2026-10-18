package web

import (
	"context"
	_ "embed"
	"encoding/json"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

//go:embed find.js
var findJS string

// findResult – kết quả find.js (toạ độ client, CSS px).
type findResult struct {
	Sel     int     `json:"sel"`
	Hidden  int     `json:"hidden"`
	Loading bool    `json:"loading"`
	Geom    geom    `json:"geom"`
	Focus   bool    `json:"focus"`
	Visible bool    `json:"visible"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	W       float64 `json:"w"`
	H       float64 `json:"h"`
	Text    string  `json:"text"`
	InView  bool    `json:"inView"`
	Covered string  `json:"covered"`
}

// found – Found chưa có khung màn hình. Bị che / tâm ngoài viewport khi phải bấm → chưa coi là thấy.
func (r findResult) found(tab string, q Query) Found {
	f := Found{Selector: r.Sel, Hidden: r.Hidden, Tab: tab, Text: r.Text, Covered: r.Covered, Loading: r.Loading}
	f.OK = r.Sel >= 0 && (q.CheckOnly || r.Covered == "" && (r.InView || !q.Scroll))
	return f
}

// cdpTimeout – mỗi lệnh CDP tối đa (trang treo / trình duyệt không trả lời thì báo lỗi, không treo engine).
const cdpTimeout = 5 * time.Second

// run chạy action trên tab tctx; dừng khi ctx (lần chạy) bị huỷ hoặc quá cdpTimeout.
func run[T any](ctx, tctx context.Context, what string, a chromedp.Action[T]) (T, error) {
	cctx, cancel := context.WithTimeout(tctx, cdpTimeout)
	defer cancel()
	stop := context.AfterFunc(ctx, cancel)
	defer stop()
	v, err := chromedp.Run(cctx, a)
	if err != nil && ctx.Err() != nil {
		return v, ctx.Err()
	}
	return v, browserErr(what, err)
}

func runFind(ctx, tctx context.Context, q Query) (findResult, error) {
	sels, err := json.Marshal(q.Selectors)
	if err != nil {
		return findResult{}, err
	}
	js := "(" + strings.TrimSpace(findJS) + ")(" + string(sels) + ", " + map[bool]string{true: "true", false: "false"}[q.Scroll] + ")"
	return run(ctx, tctx, "tìm phần tử", chromedp.Evaluate[findResult](js))
}
