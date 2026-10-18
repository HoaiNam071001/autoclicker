// Package engine – chạy kịch bản theo flow (node + cạnh): thao tác chuột / phím, chờ ảnh / chữ / màu, node điều khiển,
// biến, thử lại, tạm dừng / debug. Thiết bị qua device (Input, Screen, TextReader) để test bằng bản giả.
package engine

import (
	"image"
	"math/rand/v2"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"autoclicker/internal/device"
	"autoclicker/internal/humanize"
	"autoclicker/internal/launch"
	"autoclicker/internal/model"
	"autoclicker/internal/vision"
	"autoclicker/internal/web"
)

// Progress – tiến độ gửi sang frontend (event "engine:progress").
type Progress struct {
	Loop      int           `json:"loop"`  // vòng hiện tại, từ 1
	Loops     int           `json:"loops"` // 0 = vô hạn
	Step      int           `json:"step"`  // số thứ tự node đang chạy (theo model.FlowOrder, trong các node đang bật), từ 1
	Total     int           `json:"total"` // số node đang bật
	NodeID    string        `json:"nodeId,omitempty"`
	Ref       string        `json:"ref,omitempty"`
	Clicks    int           `json:"clicks"`
	ElapsedMs int64         `json:"elapsedMs"` // không tính thời gian tạm dừng
	Paused    bool          `json:"paused"`
	Wait      *WaitProgress `json:"wait,omitempty"`    // đang chờ ảnh / pixel
	Retry     int           `json:"retry,omitempty"`   // đang thử lại lần thứ mấy (0 = lần đầu)
	Retries   int           `json:"retries,omitempty"` // số lần thử lại tối đa của bước hiện tại
	Break     string        `json:"break,omitempty"`   // đang dừng trước node NodeID vì debug: breakpoint | step
}

// progressEvery – giãn cách tối thiểu giữa 2 lần báo tiến độ trong cùng một bước (tránh dội event).
const progressEvery = 100 * time.Millisecond

// engine chạy một kịch bản. Tạo mới cho mỗi lần chạy.
type Engine struct {
	in         device.Input
	sc         model.Script
	gate       pauseGate
	screen     model.Rect                     // khung màn hình khi Region = nil (runner đặt = mainScreen); zero = gốc (0, 0)
	displays   []model.Rect                   // khung các màn (tìm ảnh "cả màn", kiểm tra vùng chụp); nil = không kiểm tra
	onProgress func(Progress)                 // có thể nil
	intN       func(n int) int                // số ngẫu nhiên [0, n) – test thay bằng hàm cố định
	open       func(target, app string) error // bước Mở – test thay bằng hàm ghi lại
	scr        device.Screen                  // chụp màn hình / đọc màu (imagestep.go) – test thay bằng màn hình giả
	ocr        device.TextReader              // đọc chữ (tìm theo chữ, ocr.go) – test thay bằng bộ đọc giả
	clip       device.Clipboard               // clipboard (chụp màn hình chép ảnh – shot.go) – test thay bằng bản giả
	httpc      *http.Client                   // node HTTP (http.go) – thời gian chờ theo TimeoutMs của node
	browser    web.Browser                    // node web (webstep.go) – nil = không điều khiển trình duyệt

	onNode  func(model.StepResult)       // mỗi node chạy xong (event engine:node) – có thể nil
	shots   map[string]image.Image       // ảnh màn hình lúc node thất bại theo tên file (nil = không chụp – runlog.go)
	debug   bool                         // dừng ở node có Breakpoint (chạy từ app); không thì bỏ qua breakpoint
	step    atomic.Bool                  // chạy từng node: dừng trước mỗi node (StepScript bật, ResumeScript tắt)
	failed  string                       // id node làm lần chạy dừng vì lỗi (runner gửi engine:failed)
	miss    *FindResult                  // node chờ ảnh hết giờ: chỗ giống nhất (để xem lại vì sao không thấy)
	tpls    map[string]vision.Gray       // ảnh mẫu đã giải mã theo id node
	byID    map[string]model.Step        // node theo id (run đặt)
	order   map[string]int               // số thứ tự hiển thị của node (model.FlowOrder) – thanh mini
	results map[string]*model.StepResult // kết quả lần chạy gần nhất của từng node trong vòng này
	rv      runVars                      // bảng tra biến {{nX…}} / {{last…}} / {{vars…}}
	counts  map[string]int               // node Đếm: số lần đã đi qua trong vòng này

	last    model.Point // chỗ engine vừa đưa chuột tới – cho bước kéo "từ điểm trước đó"
	hasLast bool

	rng   *rand.Rand        // ngẫu nhiên của thao tác giống người (human.go) – test đặt seed
	human *humanize.Profile // mức giống người của node đang chạy (setHuman); nil = tắt

	mu       sync.Mutex
	prog     Progress
	start    time.Time
	lastEmit time.Time
}

func newEngine(in device.Input, sc model.Script, onProgress func(Progress)) *Engine {
	return &Engine{in: in, sc: sc, onProgress: onProgress, intN: rand.IntN, open: launch.OpenTarget, scr: device.RobotgoScreen{}, ocr: device.SystemOCR{}, clip: device.SystemClipboard{}, httpc: &http.Client{},
		rng: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))}
}

// jitter trả v lệch ngẫu nhiên đều trong [v-j, v+j]; j ≤ 0 → v. Dùng cho toạ độ (có thể âm ở màn phụ).
func (e *Engine) jitter(v, j int) int {
	if j <= 0 {
		return v
	}
	return v + e.intN(2*j+1) - j
}

// jitterMs như jitter nhưng không âm – dùng cho thời gian chờ.
func (e *Engine) jitterMs(v, j int) int { return max(e.jitter(v, j), 0) }

// update sửa tiến độ rồi báo sang frontend (force = báo ngay, không giãn cách).
func (e *Engine) update(force bool, f func(*Progress)) {
	e.mu.Lock()
	f(&e.prog)
	now := time.Now()
	if !force && now.Sub(e.lastEmit) < progressEvery {
		e.mu.Unlock()
		return
	}
	e.lastEmit = now
	p := e.snapshotLocked()
	e.mu.Unlock()
	if e.onProgress != nil {
		e.onProgress(p)
	}
}

// Snapshot trả tiến độ hiện tại (kèm thời gian đã chạy và trạng thái tạm dừng).
func (e *Engine) Snapshot() Progress {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.snapshotLocked()
}

func (e *Engine) snapshotLocked() Progress {
	p := e.prog
	p.ElapsedMs = e.gate.activeSince(e.start).Milliseconds()
	p.Paused = e.gate.paused()
	return p
}
