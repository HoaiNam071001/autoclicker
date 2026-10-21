// Package browser – điều khiển trình duyệt thường họ Chromium (Chrome / Edge / Brave / Cốc Cốc đã đăng nhập của người
// dùng) ở mức tab: đọc, chuyển, mở, đóng, tải lại (bước 43 B). Không CDP, không extension: macOS qua AppleScript (JXA,
// tabs_darwin.go), Windows qua UI Automation (chưa làm – 43 B1w). Khác internal/web: web điều khiển Chrome của app qua CDP.
package browser

import (
	"context"
	"errors"

	"ghosthand/internal/model"
)

// Backend – thao tác tab. Engine dùng qua interface để test bằng bản giả (browsertest.Fake).
type Backend interface {
	// Tabs – mọi tab của trình duyệt b ("" = mọi trình duyệt đang chạy, trình duyệt đang ở trước trước), trong mỗi
	// trình duyệt cửa sổ trên cùng trước, tab theo thứ tự trên thanh tab. Không có trình duyệt nào chạy → rỗng, không lỗi.
	Tabs(ctx context.Context, b string) ([]model.TabInfo, error)
	Activate(ctx context.Context, t model.TabInfo) error              // chuyển tới tab + đưa cửa sổ, app lên trước
	NewTab(ctx context.Context, b, url string) (model.TabInfo, error) // tab mới ở cửa sổ trên cùng (url "" = trang trống)
	Close(ctx context.Context, t model.TabInfo) error                 // đóng tab
	Nav(ctx context.Context, t model.TabInfo, op string) error        // model.TabReload / TabBack / TabForward
}

var (
	// ErrNotRunning – chưa có trình duyệt nào (đúng loại) đang chạy (mã lỗi browser).
	ErrNotRunning = errors.New("chưa mở trình duyệt – mở bằng node Mở trước")
	// ErrNoAutomation – macOS chưa cho app điều khiển trình duyệt (mã lỗi permission).
	ErrNoAutomation = errors.New("chưa cho phép điều khiển trình duyệt – System Settings → Privacy & Security → Automation, bật trình duyệt cho app / terminal")
	// ErrTabGone – tab đã đóng giữa chừng.
	ErrTabGone = errors.New("tab không còn mở")
)

// System – trình duyệt thật của máy.
type System struct{}
