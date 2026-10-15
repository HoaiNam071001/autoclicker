# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** chưa chọn – bước 4 xong
- **Tiếp theo:** bước 5 (overlay) – cũng là nền cho chọn vùng và hiệu ứng click
- **Cập nhật lần cuối:** 2026-10-15

## Quyết định đã chốt

- Dev trên **macOS** (Apple Silicon); build Windows làm sau (cần MinGW-w64). Trên macOS phải cấp quyền *Accessibility* cho app/terminal thì robotgo mới điều khiển được chuột.
- **Tailwind v4** (`@tailwindcss/vite`), token màu/font khai báo trong `@theme` ở `frontend/src/style.css` – không dùng `tailwind.config.js`.
- Font **Be Vietnam Pro** + **IBM Plex Mono** cài local qua `@fontsource/*` (app phải chạy được offline).
- UI theo **phương án A (MainScreen)** cho MVP. Phương án B (PreviewScreen) giữ làm tham khảo.
- Wails v2 chỉ có **một cửa sổ** → overlay và mini bar là các *chế độ* của cửa sổ chính (resize / fullscreen / trong suốt / always-on-top), không phải cửa sổ riêng.

## Stack

- **App shell:** Wails v2 (Go backend + WebView frontend)
- **Frontend:** React + TypeScript + Tailwind
- **Chuột / màn hình:** `github.com/go-vgo/robotgo`
- **Hotkey toàn cục:** `github.com/robotn/gohook`
- **Kéo thả danh sách:** `@dnd-kit/core`, `@dnd-kit/sortable`
- **Nền tảng đầu tiên:** Windows

## Khởi động nhanh

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor
wails init -n autoclicker -t react-ts
cd autoclicker
go get github.com/go-vgo/robotgo github.com/robotn/gohook
wails dev
```

> Windows: robotgo dùng cgo, cần MinGW-w64 (gcc) trong PATH.

## Data model (bản nháp)

```go
type Point struct{ X, Y int }

type Step struct {
    ID            string  `json:"id"`
    Type          string  `json:"type"`          // click | double | hold | drag | wait
    Points        []Point `json:"points"`        // drag dùng points[0] -> points[1]
    Button        string  `json:"button"`        // left | right | center
    Clicks        int     `json:"clicks"`
    ClickInterval int     `json:"clickInterval"` // ms
    HoldMs        int     `json:"holdMs"`
    DragMs        int     `json:"dragMs"`
    DelayAfter    int     `json:"delayAfter"`    // ms, sang bước kế tiếp
    Enabled       bool    `json:"enabled"`
}

type Script struct {
    Name      string `json:"name"`
    Region    *Rect  `json:"region"`    // nil = toàn màn hình; toạ độ điểm tương đối theo vùng
    Loops     int    `json:"loops"`     // 0 = vô hạn
    MaxMinutes int   `json:"maxMinutes"`
    Steps     []Step `json:"steps"`
}
```

---

## Phase 0 – Dựng nền & thử rủi ro

- [x] **1. Cài môi trường** – Go, Node, Wails CLI, MinGW-w64. Xong khi `wails doctor` không báo lỗi. *(làm ngay)*
- [x] **2. Khởi tạo project** – `wails init -n autoclicker -t react-ts`, thêm Tailwind. Xong khi `wails dev` mở được cửa sổ. *(làm ngay)*
- [x] **3. Spike click** – Nút trên UI gọi hàm Go di chuột tới (500, 500) và click. Xong khi click được ra ngoài cửa sổ app. *(làm ngay)*
  - [x] `App.SpikeClick(x, y, delayMs)` trong `spike.go`: chờ → `robotgo.Move` → so vị trí thực tế (lệch = thiếu quyền) → click trái
  - [x] Nút tạm *Thử click 500,500* trên header MainScreen, đếm ngược 3 giây, hiện kết quả / lỗi
  - [x] Kiểm tra tay trong `wails dev`: click trúng cửa sổ khác ở (500, 500)
- [x] **4. Spike hotkey toàn cục** – F8 dừng vòng lặp click, kể cả khi app không được focus.
  - [x] `hotkey.go`: gohook nghe KeyDown F8 từ startup, gọi `StopSpikeLoop`; emit `hotkey:f8` / `hotkey:error`
  - [x] `StartSpikeLoop` / `StopSpikeLoop` trong `spike.go`: click (500, 500) mỗi 500 ms, huỷ bằng `context`, tối đa 50 click; emit `spike:state`
  - [x] Nút *Lặp click 500,500* / *Dừng lặp* trên MainScreen
  - [x] Kiểm tra tay: đang lặp, focus app khác, bấm F8 (macOS: fn+F8) → dừng
- [ ] **5. Spike overlay** – Cửa sổ trong suốt, frameless, phủ toàn màn hình, luôn ở trên. Rủi ro nhất; nếu Wails không đáp ứng thì cân nhắc overlay viết riêng bằng Go.

## Phase 1 – MVP

### Engine (Go)
- [ ] **6. Định nghĩa model** – `Step`, `Script`; Wails tự sinh type TS.
- [ ] **7. Engine chạy bước** – Click trái/phải/giữa, double, giữ N ms, kéo thả A→B, chờ. Lặp N vòng / vô hạn / X phút; huỷ bằng `context`.
- [ ] **8. Pause / Resume / Stop** – Emit event tiến độ (bước, vòng, số click) sang frontend.
- [ ] **9. Dừng khẩn cấp** – Dừng ngay khi chuột chạm góc màn hình.

### Chọn vùng & điểm
- [ ] **10. Chọn vùng** – Mặc định toàn màn hình; kéo chọn hình chữ nhật trên overlay. Điểm lưu toạ độ tương đối theo vùng.
- [ ] **11. Chọn điểm** – Click trên overlay để đặt điểm, hoặc F6 lưu vị trí chuột hiện tại.
- [ ] **12. Marker & đường nối** – Marker đánh số, đường nối theo thứ tự; kéo marker để sửa vị trí.
- [ ] **13. Chọn điểm kéo thả** – Chọn A rồi B, hiển thị mũi tên A→B.

### Giao diện
- [x] **13b. Dựng UI tĩnh từ bản thiết kế** – 4 màn `MainScreen`, `PreviewScreen`, `OverlayScreen`, `MiniBar` trong `frontend/src/screens/`, dữ liệu mẫu, chưa nối Go. Xem thử bằng hash URL: `#preview`, `#overlay`, `#minibar`.
- [ ] **14. Danh sách bước** – Sắp xếp bằng kéo thả, bật/tắt, nhân bản, xoá.
- [ ] **15. Form cấu hình bước** – Loại hành động, số click, khoảng cách giữa các click, thời gian giữ, thời gian kéo, delay sau bước.
- [ ] **16. Cài đặt chạy** – Số vòng / vô hạn / thời lượng, đếm ngược 3 giây, thanh trạng thái khi đang chạy.

### Hoàn thiện
- [ ] **17. DPI & đa màn hình** – Test toạ độ ở scaling 100%, 125%, 150% và với 2 màn hình.
- [ ] **18. Build & thử máy sạch** – `wails build` ra .exe, chạy trên máy chưa cài Go/Node.

## Phase 2 – Tiện ích

- [ ] **19. Ngẫu nhiên hoá** – Delay ±ms và vị trí ±px cho từng bước.
- [ ] **20. Nhóm bước** – Lặp riêng một nhóm N lần trong kịch bản.
- [ ] **21. Profile** – Lưu/tải kịch bản JSON qua hộp thoại file, danh sách mở gần đây.
- [ ] **22. Hẹn giờ** – Tự chạy lúc HH:mm.
- [ ] **23. Thêm hành động** – Cuộn chuột, gõ phím, tổ hợp phím.

## Phase 3 – Nâng cao

- [ ] **24. Ghi macro** – Record thao tác chuột thật rồi chuyển thành các bước.
- [ ] **25. Click có điều kiện** – Chỉ click khi pixel đúng màu hoặc tìm thấy ảnh mẫu.
- [ ] **26. Gắn với cửa sổ** – Bám theo một cửa sổ app cụ thể thay vì toạ độ cố định.

---

## Nhật ký

- **2026-10-15** – Xong bước 1–2 (Go 1.27, Wails v2.16, Node 20, macOS). Port UI thiết kế (`autoclicker-ui/`) vào `frontend/`: chuyển Tailwind v3 → v4 `@theme`, font local, bỏ CSS/font template Wails. `npm run build` OK, đã chụp kiểm tra 4 màn.
- **2026-10-15** – Bước 3: thêm `spike.go` (`SpikeClick`), nút thử trên MainScreen, bỏ robotgo khỏi `deps.go`, sinh lại bindings. `go build` + `npm run build` OK. Chạy thử `robotgo.Move` từ terminal: không lỗi nhưng chuột không di chuyển → terminal chưa có quyền Accessibility; chưa kiểm tra click thật.
- **2026-10-15** – Xong bước 3: sau khi cấp quyền Accessibility, bấm *Thử click 500,500* trong `wails dev` → 3 giây sau cửa sổ khác ở (500, 500) được focus, click ra ngoài app OK.
- **2026-10-15** – Bước 4: thêm `hotkey.go` (gohook, F8), vòng lặp click thử trong `spike.go`, xoá `deps.go`, thêm `OnShutdown`. `go vet` + `go build` + `npm run build` OK; chưa kiểm tra tay F8.
- **2026-10-15** – Xong bước 4: kiểm tra tay trong `wails dev` – đang lặp click, app mất focus, bấm fn+F8 → vòng lặp dừng.
