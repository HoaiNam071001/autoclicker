package model

// Node Tab trình duyệt (bước 43 B1, engine/browsertab.go): làm với tab của trình duyệt thường họ Chromium (Chrome đã
// đăng nhập của người dùng – không phải Chrome của app ở node web) qua internal/browser, không CDP / extension.

// Thao tác của node Tab trình duyệt (Step.TabOp).
const (
	TabRead        = ""            // đọc tab (đang hiện / khớp Tab / của node TabFrom) – không đổi gì
	TabList        = "list"        // liệt kê mọi tab khớp Tab → items
	TabSwitch      = "switch"      // chuyển tới tab + đưa cửa sổ lên trước
	TabNew         = "new"         // mở tab mới (URL, trống = trang trống) ở cửa sổ trên cùng
	TabClose       = "close"       // đóng một tab
	TabCloseOthers = "closeOthers" // đóng các tab khác trong cùng cửa sổ
	TabReload      = "reload"
	TabBack        = "back"
	TabForward     = "forward"
	TabWaitOpen    = "waitOpen"  // chờ có tab khớp Tab
	TabWaitClose   = "waitClose" // chờ hết tab khớp Tab (hoặc tab của node TabFrom đóng)
)

var TabOps = []string{TabRead, TabList, TabSwitch, TabNew, TabClose, TabCloseOthers, TabReload, TabBack, TabForward, TabWaitOpen, TabWaitClose}

// Trình duyệt hỗ trợ (Step.Browser) – "" = trình duyệt đang chạy có cửa sổ ở trước nhất.
const (
	BrowserChrome = "chrome"
	BrowserEdge   = "edge"
	BrowserBrave  = "brave"
	BrowserCocCoc = "coccoc"
)

var Browsers = []string{"", BrowserChrome, BrowserEdge, BrowserBrave, BrowserCocCoc}

// TabInfo – một tab. ID / Window do trình duyệt cấp (ổn định tới khi đóng tab / cửa sổ).
type TabInfo struct {
	Browser string `json:"browser"` // Browsers (khác "")
	ID      int64  `json:"id"`
	Window  int64  `json:"window"` // id cửa sổ
	Index   int    `json:"index"`  // thứ tự trong cửa sổ, từ 1
	Title   string `json:"title"`
	URL     string `json:"url"`
	Active  bool   `json:"active"`  // tab đang hiện của cửa sổ đó
	Loading bool   `json:"loading"` // trang đang tải
}

// TabFrame – cửa sổ chứa một tab: khung (cùng hệ toạ độ với điểm click) và có đang toàn màn hình không.
type TabFrame struct {
	Box        Rect `json:"box"`
	Fullscreen bool `json:"fullscreen"`
}

// NormalizeOpts – việc node Chuẩn hoá trình duyệt nhờ trình duyệt làm (browser.Backend.Normalize). Thứ tự: thoát
// toàn màn hình → đặt khung → zoom → đóng popup.
type NormalizeOpts struct {
	Zoom       bool  `json:"zoom"`
	Fullscreen bool  `json:"fullscreen"`
	Popups     bool  `json:"popups"`
	Bounds     *Rect `json:"bounds,omitempty"` // nil = giữ khung
}

// MatchTabs – chỉ số các tab có tiêu đề hoặc URL khớp want theo mode (NameMatches); want "" → mọi tab.
func MatchTabs(tabs []TabInfo, want, mode string) ([]int, error) {
	var out []int
	for i, t := range tabs {
		ok, err := MatchName(t.Title, want, mode)
		if err == nil && !ok {
			ok, err = MatchName(t.URL, want, mode)
		}
		if err != nil {
			return nil, err
		}
		if ok {
			out = append(out, i)
		}
	}
	return out, nil
}

// ActiveTab – chỉ số tab đang hiện của cửa sổ đầu tiên (tabs xếp cửa sổ trên cùng trước); -1 = không có tab nào.
func ActiveTab(tabs []TabInfo) int {
	for i, t := range tabs {
		if t.Window == tabs[0].Window && t.Active {
			return i
		}
	}
	if len(tabs) > 0 {
		return 0
	}
	return -1
}
