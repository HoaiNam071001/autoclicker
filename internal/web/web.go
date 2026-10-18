// Package web – node web (bước 25d): tìm phần tử bằng HTML qua CDP trên Chrome của app (profile riêng), trả khung
// phần tử theo toạ độ màn hình để engine bấm bằng chuột thật. Toạ độ trang → màn hình đo bằng chuột (calibrate.go),
// không tính từ kích thước cửa sổ. Engine dùng qua interface Browser để test bằng trình duyệt giả.
package web

import (
	"context"

	"autoclicker/internal/model"
)

// Browser – trình duyệt engine điều khiển. Một bản cho mỗi lần chạy; Close nhả kết nối (không đóng trình duyệt).
type Browser interface {
	// Open mở url trong tab mới (chưa mở trình duyệt thì bật lên); tab này thành tab mặc định của node web sau.
	Open(ctx context.Context, url string) error
	// Find tìm phần tử một lần (không chờ). Không thấy → Found.OK = false, không lỗi.
	Find(ctx context.Context, q Query) (Found, error)
	// Hover đưa chuột thật tới (x, y) – điểm sắp bấm của lần Find (Scroll) gần nhất – và xem trang có nhận đúng chỗ
	// không (D0: số đo cũ, Dock / thanh menu / cửa sổ khác che mà elementFromPoint không thấy). onTarget = điểm nằm
	// trong khung phần tử → phần tử dưới chuột phải là nó.
	Hover(ctx context.Context, x, y int, onTarget bool) (Hover, error)
	Close()
}

// Query – phần tử cần tìm: thử từng selector theo thứ tự, cái đầu thấy thì dùng.
type Query struct {
	Selectors []model.Selector
	Tab       string // URL / tiêu đề tab chứa chuỗi này ("" = tab đang hiện)
	Scroll    bool   // sắp bấm (webClick): ngoài viewport thì cuộn vào giữa, bị che thì chưa tính là thấy, tab lên trước
	CheckOnly bool   // chỉ xem có hay không (webWait chờ mất) – không đo toạ độ, không di chuột
}

// Found – kết quả một lần tìm.
type Found struct {
	OK       bool
	Rect     model.Rect // khung phần tử theo toạ độ màn hình (robotgo)
	Text     string     // chữ hiển thị (cắt ngắn)
	Selector int        // chỉ số selector đã khớp (-1 = không)
	Covered  string     // thấy nhưng bị phần tử khác che – mô tả phần tử che (OK = false)
	Hidden   int        // số phần tử khớp nhưng đang ẩn / cỡ 0 (gợi ý khi không thấy)
	Tab      string     // URL tab đã tìm
	NoTab    bool       // không có tab nào khớp Query.Tab
	Loading  bool       // trang đang tải – không thấy chưa có nghĩa là đã mất
	Moving   bool       // sắp bấm nhưng khung còn đổi giữa 2 lần đọc (hiệu ứng, trang đang dàn lại) – OK = false
}

// Hover – kết quả thử hover trước khi bấm. OK = false → chưa bấm, tìm lại (Miss nói vì sao, để báo khi hết giờ).
type Hover struct {
	OK   bool
	Miss string
}

// ErrNoBrowser – Chrome của app chưa mở / không còn tab phù hợp (StepError code browser).
var ErrNoBrowser error = &BrowserError{"Chrome của app chưa mở – thêm node Mở URL với tuỳ chọn “Chrome của app” trước node web"}

// BrowserError – lỗi kết nối / điều khiển trình duyệt (StepError code browser).
type BrowserError struct{ Msg string }

func (e *BrowserError) Error() string { return e.Msg }

// QueryError – selector lỗi khi chạy trong trang (CSS / XPath sai cú pháp…) – StepError code invalid_input.
type QueryError struct{ Msg string }

func (e *QueryError) Error() string { return e.Msg }
