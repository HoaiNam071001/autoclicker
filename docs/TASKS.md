# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** bước 9 – Dừng khẩn cấp. Còn tồn: đo độ lệch toạ độ overlay trên macOS
- **Tiếp theo:** bước 9 (dừng khẩn cấp khi chuột chạm góc màn hình)
- **Cập nhật lần cuối:** 2026-10-15

## Quyết định đã chốt

- Dev trên **macOS** (Apple Silicon); build Windows làm sau (cần MinGW-w64). Trên macOS phải cấp quyền *Accessibility* cho app/terminal thì robotgo mới điều khiển được chuột.
- **Tailwind v4** (`@tailwindcss/vite`), token màu/font khai báo trong `@theme` ở `frontend/src/style.css` – không dùng `tailwind.config.js`.
- Font **Be Vietnam Pro** + **IBM Plex Mono** cài local qua `@fontsource/*` (app phải chạy được offline).
- UI theo **phương án A (MainScreen)** cho MVP. Phương án B (PreviewScreen) giữ làm tham khảo.
- **Gõ chữ / phím tắt đưa vào MVP** (loại bước `type`, `key`) thay vì Phase 2, để model, engine, form làm một lần.
- Cửa sổ luôn **frameless + nền trong suốt** (Wails v2 không đổi được lúc chạy); MainScreen tự vẽ thanh kéo và nút thu nhỏ / đóng.
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

## Data model

Nguồn chuẩn: **`model.go`** (`Point`, `Rect`, `Step`, `Script`, hằng loại bước, `NewStep` / `NewScript` cho giá trị mặc định). Wails sinh class TS vào `frontend/wailsjs/go/models.ts`; `frontend/src/types.ts` giữ interface `Step` với union type chặt hơn và có kiểm tra lúc type-check rằng field khớp struct Go.

Loại bước: `click` · `double` · `hold` · `drag` (points[0] → points[1]) · `wait` · `type` (gõ `text`) · `key` (bấm `key` kèm `modifiers`: cmd / ctrl / alt / shift).

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
- [x] **5. Spike overlay** – Cửa sổ trong suốt, frameless, phủ toàn màn hình, luôn ở trên. Rủi ro nhất; nếu Wails không đáp ứng thì cân nhắc overlay viết riêng bằng Go.
  - [x] Cửa sổ luôn frameless + webview trong suốt (`main.go`); MainScreen có thanh kéo (`.drag`) và nút thu nhỏ / đóng
  - [x] `App.SetWindowMode(main|overlay|minibar)` trong `window.go`: lưu/khôi phục khung main, overlay = cỡ màn hình + always-on-top, mini bar 460×76 giữa trên
  - [x] `windowMode.ts`: gọi Go rồi đổi hash → App đổi màn. Overlay: crosshair theo chuột, Esc / *Xong* thoát, click hiện toạ độ cửa sổ vs robotgo (độ lệch)
  - [x] Kiểm tra tay: nhìn xuyên được, phủ màn hình, nằm trên app khác, độ lệch toạ độ (macOS dự đoán lệch Y = chiều cao menu bar)

## Phase 1 – MVP

### Engine (Go)
- [x] **6. Định nghĩa model** – `Step`, `Script`; Wails tự sinh type TS. Thêm loại `type` (gõ chữ) và `key` (phím / tổ hợp phím).
- [x] **7. Engine chạy bước** – Click trái/phải/giữa, double, giữ N ms, kéo thả A→B, chờ, gõ chữ (`robotgo.TypeStr`, thử với bộ gõ tiếng Việt; dự phòng dán qua clipboard), phím tắt (`robotgo.KeyTap`). Lặp N vòng / vô hạn / X phút; huỷ bằng `context`.
  - [x] `engine.go`: `validateScript` + `runScript(ctx, Input, Script)` – lặp N vòng / vô hạn (`loops = 0`) / `maxMinutes`, toạ độ cộng gốc vùng, giữ / kéo luôn nhả chuột khi bị huỷ, gõ từng ký tự để dừng được giữa chừng
  - [x] `input.go`: interface `Input` + `robotgoInput` (kéo dùng `robotgo.Drag` vì `Move` trên macOS không gửi sự kiện kéo)
  - [x] `runner.go`: `App.RunScript` (chạy tới khi xong / dừng), `StopScript`; F8 gọi `StopScript`; emit `engine:state`
  - [x] `engine_test.go`: 5 test với input giả (click + vùng + vòng lặp, kéo, gõ / phím, huỷ vô hạn có nhả chuột, validate)
  - [x] Kiểm tra tay: nút *Chạy thử engine* với TextEdit – click, gõ tiếng Việt (có / không bật bộ gõ), Enter, double click, kéo chọn, Cmd+A; F8 dừng giữa chừng
- [x] **8. Pause / Resume / Stop** – Emit event tiến độ (bước, vòng, số click) sang frontend.
  - [x] `engine.go`: engine thành struct; `pauseGate` chặn trước mỗi bước / lần click / ký tự; `Progress` (vòng, bước, số click, thời gian không tính lúc tạm dừng), báo ngay khi đổi bước / tạm dừng, còn lại giãn ≥ 100 ms
  - [x] `runner.go`: `PauseScript` / `ResumeScript`, emit `engine:progress`
  - [x] `engineStore.ts`: store chung nghe event (mini bar mở sau vẫn có dữ liệu); MiniBar hiện số thật, ⏸ / ▶ hoạt động, *Dừng* = StopScript + về main
  - [x] Test `TestPauseResumeAndProgress`; `go test -race` 6/6 PASS
  - [x] Kiểm tra tay: *Chạy thử engine* (3 vòng) → cửa sổ thành mini bar, số chạy, tạm dừng / tiếp tục, Dừng hoặc F8 → về main
- [ ] **9. Dừng khẩn cấp** – Dừng ngay khi chuột chạm góc màn hình.

### Chọn vùng & điểm
- [ ] **10. Chọn vùng** – Mặc định toàn màn hình; kéo chọn hình chữ nhật trên overlay. Điểm lưu toạ độ tương đối theo vùng.
- [ ] **11. Chọn điểm** – Click trên overlay để đặt điểm, hoặc F6 lưu vị trí chuột hiện tại.
- [ ] **12. Marker & đường nối** – Marker đánh số, đường nối theo thứ tự; kéo marker để sửa vị trí.
- [ ] **13. Chọn điểm kéo thả** – Chọn A rồi B, hiển thị mũi tên A→B.

### Giao diện
- [x] **13b. Dựng UI tĩnh từ bản thiết kế** – 4 màn `MainScreen`, `PreviewScreen`, `OverlayScreen`, `MiniBar` trong `frontend/src/screens/`, dữ liệu mẫu, chưa nối Go. Xem thử bằng hash URL: `#preview`, `#overlay`, `#minibar`.
- [ ] **14. Danh sách bước** – Sắp xếp bằng kéo thả, bật/tắt, nhân bản, xoá.
- [ ] **15. Form cấu hình bước** – Loại hành động, số click, khoảng cách giữa các click, thời gian giữ, thời gian kéo, nội dung gõ, phím + phím kèm, delay sau bước.
- [ ] **16. Cài đặt chạy** – Số vòng / vô hạn / thời lượng, đếm ngược 3 giây, thanh trạng thái khi đang chạy.

### Hoàn thiện
- [ ] **17. DPI & đa màn hình** – Test toạ độ ở scaling 100%, 125%, 150% và với 2 màn hình.
- [ ] **18. Build & thử máy sạch** – `wails build` ra .exe, chạy trên máy chưa cài Go/Node.

## Phase 2 – Tiện ích

- [ ] **19. Ngẫu nhiên hoá** – Delay ±ms và vị trí ±px cho từng bước.
- [ ] **20. Nhóm bước** – Lặp riêng một nhóm N lần trong kịch bản.
- [ ] **21. Profile** – Lưu/tải kịch bản JSON qua hộp thoại file, danh sách mở gần đây.
- [ ] **22. Hẹn giờ** – Tự chạy lúc HH:mm.
- [ ] **23. Thêm hành động** – Cuộn chuột. *(gõ phím, tổ hợp phím đã chuyển vào MVP – bước 6, 7)*

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
- **2026-10-15** – Commit bước 3–4 (`0fbbb57`). Bước 5: thêm `window.go` (`SetWindowMode`, `MousePos`), cửa sổ frameless trong suốt, `windowMode.ts`, overlay đo toạ độ. Ghi chú: Wails macOS `WindowSetPosition` tính từ `visibleFrame` (dưới menu bar). `go vet` + `go build` + `npm run build` OK; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 5 (kiểm tra tay macOS): overlay trong suốt, nằm trên app khác, crosshair theo chuột OK. Vấn đề: overlay bị đẩy xuống dưới menu bar nên đáy tràn khỏi màn hình (ô đo toạ độ ở đáy bị khuất → chuyển lên góc trên). Cần xử lý độ lệch toạ độ ở bước 10–11 / 17.
- **2026-10-15** – Kiểm tra tay mini bar: hiện giữa phía trên, kéo được, *Dừng* về màn chính OK.
- **2026-10-15** – Xong bước 6: `model.go` (Step, Script, Rect, `NewStep`, `NewScript`), thêm loại bước `type` / `key` theo yêu cầu; `types.ts` khớp struct Go (kiểm tra lúc tsc – đã thử thêm field thừa thì tsc báo lỗi); form MainScreen có ô nội dung gõ / phím + phím kèm. `go vet` + `npm run build` OK.
- **2026-10-15** – Bước 7: engine (`engine.go`, `input.go`, `runner.go`), bỏ vòng lặp click thử của bước 4 (F8 giờ dừng engine), nút *Chạy thử engine* chạy kịch bản demo. `go test` 5/5 PASS, `go vet` + `npm run build` OK; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 7: kiểm tra tay với TextEdit trên macOS – người dùng xác nhận chạy đúng như kỳ vọng.
- **2026-10-15** – Bước 8: pause / resume (`pauseGate`), event tiến độ, `engineStore.ts`, MiniBar dữ liệu thật; demo engine chạy 3 vòng trong mini bar. Ghi chú: `maxMinutes` tính theo đồng hồ thật (kể cả lúc tạm dừng). `go test -race` 6/6, `npm run build` OK; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 8: kiểm tra tay macOS – mini bar hiện tiến độ, tạm dừng / tiếp tục, Dừng / F8 về main OK.
