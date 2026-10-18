package engine

import (
	"context"
	"errors"
	"time"

	"autoclicker/internal/model"
)

// WebFindResult – kết quả "Thử tìm" của node web (D3; MCP 31 dùng chung): tìm một lần, không cuộn, không bấm.
type WebFindResult struct {
	Found    bool   `json:"found"`
	X        int    `json:"x"` // khung phần tử theo toạ độ màn hình (robotgo)
	Y        int    `json:"y"`
	W        int    `json:"w"`
	H        int    `json:"h"`
	Text     string `json:"text"`
	Selector string `json:"selector"` // dòng selector đã khớp ("" = không)
	Covered  string `json:"covered"`  // phần tử đang che tâm
	OffView  bool   `json:"offView"`  // ngoài vùng nhìn thấy – click sẽ tự cuộn tới
	Hidden   int    `json:"hidden"`   // số phần tử khớp nhưng đang ẩn
	Tab      string `json:"tab"`      // URL tab đã tìm
	NoTab    bool   `json:"noTab"`
	Count    int    `json:"count"` // số phần tử đang hiện khớp
	Value    string `json:"value"` // webRead: giá trị sẽ đọc được
	Ms       int    `json:"ms"`
}

// FindElement tìm phần tử của node web s trên tab hiện tại – một lần, không chờ. Lần đầu chuột đi vào trang để đo
// toạ độ (như lúc chạy). {{biến}} trong s phải thay sẵn (script.WithInputDefaults).
func (e *Engine) FindElement(ctx context.Context, s model.Step) (WebFindResult, error) {
	if e.browser == nil {
		return WebFindResult{}, errNoBrowserOpt
	}
	q, lines, err := webQuery(s, s.Selectors, false)
	if err != nil {
		return WebFindResult{}, err
	}
	if len(lines) == 0 {
		return WebFindResult{}, errors.New("chưa nhập cách tìm phần tử")
	}
	q.Read, q.Attr = readOf(s) // webRead: Thử tìm cho xem luôn giá trị sẽ đọc
	start := time.Now()
	f, err := e.browser.Find(ctx, q)
	if err != nil {
		return WebFindResult{}, err
	}
	r := WebFindResult{
		Found: f.OK, X: f.Rect.X, Y: f.Rect.Y, W: f.Rect.W, H: f.Rect.H, Text: f.Text, Covered: f.Covered, OffView: f.OffView,
		Hidden: f.Hidden, Tab: f.Tab, NoTab: f.NoTab, Count: f.Count, Value: f.Value, Ms: int(time.Since(start).Milliseconds()),
	}
	if f.Selector >= 0 {
		r.Selector = lines[f.Selector]
	}
	return r, nil
}
