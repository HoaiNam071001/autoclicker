# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** 20 – flow kiểu n8n: F0 xong, đang F4a (thiết kế canvas + `#demo-flow`). Chờ kiểm tra tay 25 / 25b / 25e. 30 (sửa nhiều điểm) chờ kiểm tra tay; còn ánh xạ theo 4 điểm mốc. Bước 17 (phần Windows) và 18 hoãn: chưa có máy Windows
- **Tiếp theo:** 20: F4a → F1 → F2 → F4b → F3 → F5 → F6 → F7 → 26 → 25d (web qua DOM) … (bảng thứ tự đầu Phase 3; 25c đã gộp vào 20 F5); MCP cho AI (31) làm sau 22, các bước trước nhớ làm mục *Chuẩn bị cho MCP*. Khi có máy Windows: quay lại 17–18
- **Cập nhật lần cuối:** 2026-10-17 (bỏ cây khối 20a–20c, plan mới cho 20: flow kiểu n8n)

## Quyết định đã chốt

- Dev trên **macOS** (Apple Silicon); build Windows làm sau (cần MinGW-w64). Trên macOS phải cấp quyền *Accessibility* cho app/terminal thì robotgo mới điều khiển được chuột.
- **Tailwind v4** (`@tailwindcss/vite`), token màu/font khai báo trong `@theme` ở `frontend/src/style.css` – không dùng `tailwind.config.js`.
- Font **Be Vietnam Pro** + **IBM Plex Mono** cài local qua `@fontsource/*` (app phải chạy được offline).
- UI theo **phương án A (MainScreen)** cho MVP. Phương án B (PreviewScreen) giữ làm tham khảo.
- **Design v2** (2026-10-16, `design/DESIGN.md`) thay bản thiết kế cũ `autoclicker-ui/` từ bước 27: theme sáng / tối, một màu nhấn, gradient chỉ ở 4 chỗ, bố cục 3 cột kiểu Windows 11; font Plus Jakarta Sans + JetBrains Mono (vẫn cài local); bỏ màu riêng cho kéo thả và nhóm `dark-*`.
- Thứ tự Phase 3 (2026-10-16, sắp lại 2026-10-17, bảng đầu Phase 3): 27 → 27b → 29 → 21 → 30 → 21b → 23 → 25 → 25e → 25b → 20 (flow, gồm cả 25c) → 26 → 25d → 22 → 31a → 31b → 31c → 24 → 28 (thêm 31 ngày 2026-10-17; 25d tách từ 28 cùng ngày); chờ ảnh / click theo ảnh làm sớm (cốt lõi cho web / app), nhóm bước làm sau xử lý hết giờ. Không làm bản chạy trên điện thoại.
- **Tự động hoá qua AI là mục tiêu chính** (31, 2026-10-17): người dùng ra lệnh bằng lời trên Claude / GPT / … (vd "mở bài Tinh Hà Say 2") và AI tự tạo kịch bản, không cần soạn tay. Kết nối qua **MCP** (giao thức chuẩn mà Claude, ChatGPT, Cursor, Gemini CLI… đều hỗ trợ), không viết riêng plugin cho từng AI. App là MCP server; AI tạo / sửa kịch bản qua cùng `checkScript` / `validateScript` như người soạn; kịch bản do AI tạo mở trong app để người dùng duyệt, AI chỉ chạy thật khi người dùng cho phép. **MCP làm sau** (sau 22), nhưng các bước trước phải giữ đường cho nó – mỗi bước có mục *Chuẩn bị cho MCP*:
  - Logic nằm ở **Go** (kiểm tra, `resolveInputs`, chạy bước, chụp / tìm ảnh, đổi toạ độ, danh sách cửa sổ), frontend chỉ gọi – MCP / CLI / API dùng lại, không viết lại
  - Lỗi trả có cấu trúc (`stepId`, `field`, thông điệp) để AI tự sửa được; `id` bước ổn định để sửa theo bước
  - Mỗi loại bước có mô tả + ví dụ JSON trong `model.go` (bắt đầu từ 23) để sau sinh schema cho AI
  - Kịch bản đổi từ ngoài app (file / API / AI) phải nạp lại được vào store, vẫn hoàn tác được; việc chạy từ ngoài đi qua cùng một chỗ hỏi quyền
- **Web qua DOM** (25d, 2026-10-17): tìm phần tử bằng HTML qua CDP (Chrome / Edge / Brave, profile riêng), bấm / cuộn bằng **robotgo** (chuột thật) – CDP `Input.*` chỉ là dự phòng chạy nền. Toạ độ trang → màn hình **đo bằng chuột** (như overlay), không tính từ kích thước cửa sổ. Selector ưu tiên chữ / aria / role, tránh class sinh tự động. Chờ ảnh (25) vẫn giữ cho app desktop và trang không nối CDP được.
- **Theme tối = bảng Slate** (`design/newdarkmode.md`, 2026-10-16): theme cũ nhiều nút bị chìm. Có thêm 2 bảng Dimmed / Nord trong file, chưa dùng.
- **Neo điểm theo trục** (29): X neo trái / phải, Y neo trên / dưới (thêm giữa, %), lưu offset dương + tên neo, UI hiện dạng ±; làm ngay sau UI (27, 27b), trước 21, để định dạng file có neo ngay từ đầu (không cần chuyển đổi). Biến đổi nhiều điểm (30) làm sau 27 + 21, ghi lại offset theo neo.
- **UI trước, logic sau** (27b): tính năng chưa có vẫn hiện trong app nhưng tắt + nhãn *Sắp có*; dữ liệu giả chỉ ở màn `#demo` (`frontend/src/demo/`), không thêm loại bước chưa làm vào `types.ts` / `model.go`.
- **Input của kịch bản** áp vào bản sao trước khi chạy (`resolveInputs`), engine không đổi; giá trị lần trước / preset lưu trên máy, không ghi vào file kịch bản; input trong nhóm lặp dùng cùng giá trị mọi vòng (mỗi vòng một giá trị → CSV, bước 28).
- **File kịch bản** (21, 2026-10-17): `.json` có `version` (hiện 1); mở file chỉ kiểm tra cấu trúc (`checkScript`), không chạy `validateScript` – bản soạn dở vẫn mở được, lúc Chạy mới kiểm tra đủ. *Lưu thành* đổi tên kịch bản theo tên file. Mỗi lần mở một kịch bản; bản nháp tự lưu ghi cả đường dẫn + cờ chưa lưu. Đóng app chọn *Không lưu* thì xoá nháp.
- **Rẽ nhánh = flow kiểu n8n** (20, 2026-10-17 – thay bản cây khối): kịch bản là đồ thị node + cạnh; mỗi node trả `StepResult` (status, cổng ra, output, lỗi, log), cạnh nối theo cổng ra (`success` / `fail` / `true` / `false` / case…). Cổng `fail` không nối = dừng kịch bản báo lỗi. Soạn trên canvas (`@xyflow/react`), bỏ danh sách bước kéo thả. File kịch bản cũ (version < 4) bỏ, không chuyển đổi.
- **Gõ chữ / phím tắt đưa vào MVP** (loại bước `type`, `key`) thay vì Phase 2, để model, engine, form làm một lần.
- Toạ độ lưu theo **toạ độ màn hình thật (robotgo)**. Overlay tự đo độ lệch cửa sổ ↔ màn hình (so `clientX/Y` với `robotgo.Location` khi chuột đứng yên) thay vì tính từ vị trí cửa sổ – Wails macOS đặt cửa sổ dưới menu bar.
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

Loại bước: `click` · `double` · `hold` · `drag` (points[0] → points[1]) · `wait` · `type` (gõ `text`) · `key` (bấm `keys` – một hay nhiều phím chính – kèm `modifiers`: cmd / ctrl / alt / shift; `keyMode`: `chord` nhấn giữ cùng lúc, vd Ctrl+C+T, hoặc `sequence` giữ phím kèm rồi gõ lần lượt).

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
- [x] **9. Dừng khẩn cấp** – Dừng ngay khi chuột chạm góc màn hình.
  - [x] `failsafe.go`: đọc vị trí chuột mỗi 50 ms khi đang chạy; chạm 1 trong 4 góc màn hình chính (≤ 2 px) → dừng. `trackingInput` nhớ điểm engine vừa đưa chuột tới để kịch bản có điểm ở góc không tự kích hoạt
  - [x] Event `engine:stopped` kèm lý do (`f8` / `corner`), màn chính hiện lý do dừng; thay event `hotkey:f8`
  - [x] Test `TestInCorner`, `TestWatchCorners`; `go test -race` 8/8 PASS
  - [x] Kiểm tra tay: đang *Chạy thử engine*, rê chuột vào góc màn hình → dừng, về main, hiện "Dừng khẩn cấp…"
  - Giới hạn: chỉ màn hình chính; nếu kịch bản di chuột liên tục (click dày) thì khó chạm góc kịp giữa 2 lần engine kéo chuột đi – khi đó dùng F8

### Chọn vùng & điểm
- [x] **10. Chọn vùng** – Mặc định toàn màn hình; kéo chọn hình chữ nhật trên overlay. Điểm lưu toạ độ tương đối theo vùng.
  - [x] `scriptStore.ts`: store kịch bản dùng chung MainScreen ↔ Overlay (`useScript`, `setRegion`, `toScreen`); `types.ts` thêm `Rect`, `Script` (khớp Go)
  - [x] Overlay: kéo chuột vẽ vùng (< 10 px bỏ qua, Esc huỷ khi đang kéo), nút *Toàn màn hình*; tự đo lệch toạ độ; crosshair hiện toạ độ màn hình thật; marker + đường nối vẽ từ các bước trong store
  - [x] MainScreen: *Chọn vùng* mở overlay (`openOverlay('region')`), footer hiện vùng + nút về toàn màn hình; danh sách bước đọc từ store
  - [x] Kiểm tra tay: chọn vùng, marker dịch theo vùng, toạ độ crosshair khớp toạ độ thật
  - Giới hạn: độ lệch chỉ là phép tịnh tiến – màn hình scale ≠ 100% (Windows DPI) xử lý ở bước 17
- [x] **11. Chọn điểm** – Click trên overlay để đặt điểm, hoặc F6 lưu vị trí chuột hiện tại.
  - [x] Overlay công cụ *Đặt điểm*: click → thêm bước click mới (toạ độ lấy từ `MousePos`, tương đối theo vùng; ngoài vùng thì báo, không thêm), ở lại overlay để đặt tiếp
  - [x] F6 toàn cục (`hotkey.go`, emit `hotkey:f6`): thêm điểm tại vị trí chuột, không cần mở overlay; bỏ qua khi đang chạy
  - [x] Nút ⌖ trong form: chọn lại điểm cho bước đang chọn (click / F6 trên overlay) rồi về màn chính; Esc huỷ
  - [x] `scriptStore.ts`: `placePoint`, `pickPointFor`, bước đang chọn + thông báo trong store; *Thêm điểm* mở overlay
  - [x] Kiểm tra tay: thêm điểm bằng click và fn+F6 (cả khi app không focus), chọn lại điểm, điểm ngoài vùng bị từ chối
- [x] **12. Marker & đường nối** – Marker đánh số, đường nối theo thứ tự; kéo marker để sửa vị trí.
  - [x] Marker đánh số theo danh sách bước + đường nối + mũi tên kéo thả vẽ từ store (làm ở bước 10)
  - [x] Kéo marker (cả A / B) để dời điểm – toạ độ lấy từ `MousePos` lúc thả, ngoài vùng thì báo và giữ chỗ cũ, Esc huỷ; `movePoint` trong `scriptStore.ts`
  - [x] Click marker (không kéo) → chọn bước, marker bước đang chọn có viền nét đứt; click marker không tạo điểm mới
  - [x] Kiểm tra tay: kéo marker trong app, toạ độ cập nhật ở form màn chính
- [x] **13. Chọn điểm kéo thả** – Chọn A rồi B, hiển thị mũi tên A→B.
  - [x] Overlay công cụ *Kéo thả A → B*: click A (marker A + mũi tên nét đứt theo chuột) rồi click B → thêm bước kéo thả; Esc huỷ A; F6 cũng chọn được A / B
  - [x] Nút *Kéo thả* ở màn chính mở overlay công cụ này; nút ⌖ của bước kéo thả chọn lại cả A và B (tắt cho bước *Chờ*)
  - [x] `scriptStore.ts`: `dragA` trong store, `setDragMode`, `cancelDragA`; `placePoint` gom cặp A, B
  - [x] Bổ sung (người dùng yêu cầu, 2026-10-15): A = vị trí **hoặc từ điểm trước** (chỗ chuột dừng ở bước chuột bật gần nhất phía trước – click / giữ / điểm thả B), B = vị trí **hoặc lệch so với A** (ΔX, ΔY có thể âm). Model `dragFromPrev`, `dragRelative`; engine nhớ vị trí chuột thật (kể cả ngẫu nhiên); kiểm tra báo lỗi nếu không có điểm trước; store `endpoints()` tính toạ độ đầu / cuối như engine cho overlay + kéo marker; đổi chế độ giữ nguyên chỗ trên màn hình; test `TestDragFromPrevAndRelative`, `TestValidateDragFromPrevNeedsPrevPoint`
  - [x] Kiểm tra tay phần bổ sung: kéo từ điểm click trước, kéo nối tiếp từ điểm thả trước, B lệch ΔX/ΔY, kéo marker B khi đang ở chế độ lệch
  - [x] Kiểm tra tay: thêm bước kéo thả, chọn lại A / B, Esc huỷ A

### Giao diện
- [x] **13b. Dựng UI tĩnh từ bản thiết kế** – 4 màn `MainScreen`, `PreviewScreen`, `OverlayScreen`, `MiniBar` trong `frontend/src/screens/`, dữ liệu mẫu, chưa nối Go. Xem thử bằng hash URL: `#preview`, `#overlay`, `#minibar`.
- [x] **14. Danh sách bước** – Sắp xếp bằng kéo thả, bật/tắt, nhân bản, xoá.
  - [x] `components/StepList.tsx` (tách khỏi MainScreen): kéo ⋮⋮ sắp xếp (`@dnd-kit`, kéo ≥ 4 px mới tính), ô tick bật/tắt, nhân bản / xoá hiện khi rê chuột
  - [x] Xoá có *Hoàn tác*; phím Delete / Backspace xoá bước đang chọn (trừ khi đang gõ trong ô nhập); ô *Bật* trong form hoạt động
  - [x] *+ Thêm bước*: menu Click / Kéo thả (mở overlay chọn điểm), Chờ / Gõ chữ / Phím tắt (thêm ngay, mặc định giống `NewStep`)
  - [x] `scriptStore.ts`: `toggleStep`, `duplicateStep`, `removeStep`, `undoDelete`, `moveStep`, `addStep`; màn chính xử lý danh sách trống
  - [x] Kiểm tra tay: sắp xếp, bật/tắt, nhân bản, xoá + hoàn tác, Delete, thêm bước qua menu
- [x] **15. Form cấu hình bước** – Loại hành động, số click, khoảng cách giữa các click, thời gian giữ, thời gian kéo, nội dung gõ, phím + phím kèm, delay sau bước.
  - [x] `components/StepForm.tsx` (tách khỏi MainScreen): form điều khiển sửa thẳng vào store (`updateStep`); ô số giữ bản nháp khi gõ, chỉ ghi khi là số nguyên hợp lệ
  - [x] Trường theo loại: click (nút, vị trí, số click, giữa các click khi > 1), double, giữ (thời gian giữ), kéo thả (A, B, thời gian kéo), chờ, gõ chữ (đếm ký tự, cảnh báo trống), phím tắt (gợi ý tên phím, cảnh báo tên sai, Cmd/Ctrl/Option/Shift – Windows hiện Win/Alt)
  - [x] Đổi loại bước (`changeStepType`) tự bổ sung điểm (kéo thả thêm B lệch A 100 px) và giá trị mặc định
  - [x] Phím tắt nhiều phím chính (theo yêu cầu): model `key` → `keys` + `keyMode` (chord / sequence); engine `pressKeys` dùng `KeyToggle` (phím chính mang cờ phím kèm khi nhấn, thả không kèm để phím kèm vẫn giữ), luôn nhả hết kể cả khi huỷ; form ô phím nhận "c, t" + chọn cách bấm; test `TestKeyChordAndSequence`, `TestKeyCancelReleasesAll`
  - [x] Kiểm tra tay: sửa từng loại bước, đổi loại, giá trị lưu khi chuyển bước / mở overlay
- [x] **16. Cài đặt chạy** – Số vòng / vô hạn / thời lượng, đếm ngược 3 giây, thanh trạng thái khi đang chạy.
  - [x] Nút *Chạy* / F8 (khi chưa chạy, emit `hotkey:f8start`) chạy kịch bản đang soạn: `ValidateScript` trước (lỗi thì ở lại màn chính, hiện ở thanh dưới) → mini bar → `RunScript(sc, countdown)` → về main; đang chạy thì nút thành *Dừng*
  - [x] Đếm ngược 3 giây trong Go (huỷ được bằng F8 / Dừng), emit `engine:countdown`, mini bar hiện "Bắt đầu sau N giây…"
  - [x] `components/RunSettings.tsx`: Số vòng / Vô hạn / Theo thời gian (suy ra từ `loops` + `maxMinutes`), ô số vòng / số phút, bật/tắt đếm ngược; `components/NumInput.tsx` dùng chung
  - [x] Thanh dưới: toạ độ chuột thật (250 ms), lỗi phím toàn cục, lỗi lần chạy trước, lý do dừng, thông báo soạn
  - [x] Dọn Phase 0: bỏ `spike.go`, nút *Thử click / Chạy thử engine / Thử mini bar* và kịch bản demo; `MousePos` trả `Point`
  - [x] Sửa lỗi (người dùng báo): app ở màn A mà overlay chọn vùng / điểm mở sang màn B – `SetWindowMode` phóng to trước khi dời nên cửa sổ tràn sang màn bên cạnh; đổi thành dời về góc trước rồi mới phóng to
  - [x] Kiểm tra tay: chạy kịch bản tự soạn (Số vòng / Vô hạn / Theo thời gian), đếm ngược, F8 bắt đầu / dừng, Ctrl+C+T, capslock

### Hoàn thiện
- [ ] **17. DPI & đa màn hình** – Test toạ độ ở scaling 100%, 125%, 150% và với 2 màn hình.
  - Đã chốt: không làm bộ chọn màn hình – overlay mở trên màn chứa cửa sổ app, toạ độ dùng hệ toạ độ chung (đã thử chọn vùng trên màn thứ 2 OK, 2026-10-15)
  - [x] **Chọn màn hình cho overlay** *(người dùng yêu cầu, 2026-10-17 – đổi quyết định trên)*: `Displays()` / `SetOverlayDisplay(i)` (`window.go`); thanh công cụ overlay có nút Màn 1 / 2… khi > 1 màn, chọn rồi thì các lần mở overlay sau cũng mở ở màn đó (đến khi tắt app); về main thì dời cửa sổ về màn cũ. Wails chỉ đặt vị trí theo màn hiện tại → `moveToDisplay` cộng độ lệch giữa 2 màn; màn hiện tại đoán bằng cỡ màn (Wails) + chuột (`pickDisplay`, test `TestPickDisplay`). Ở màn phụ nút *Toàn màn này* đặt vùng = cả màn đó (vùng nil vẫn là màn chính). Đổi màn → đo lại toạ độ
  - [x] **Tìm ảnh / chờ pixel trên màn phụ** *(người dùng báo lỗi, 2026-10-17: vùng = cả màn MacBook (1920, 406) → "vùng tìm nằm ngoài màn hình chính")*: ScreenCaptureKit chọn màn chứa tâm vùng, `sourceRect` tính từ góc màn đó; `captureArea` chỉ cần vùng nằm gọn trong **một** màn (vắt 2 màn → báo lỗi rõ); *Cả màn* = màn chứa vùng; macOS < 14 (robotgo) vẫn chỉ màn chính. Test `TestCaptureAreaDisplays`; chụp thật màn MacBook ra 3024×1964 OK
  - [x] Sửa (người dùng báo, 2026-10-17): vùng ở màn 1 mà thêm bước Click lại mở overlay ở màn 2 → mỗi lần mở overlay chọn **màn chứa vùng đang chọn** (toàn màn hình = màn chính) – `OverlayDisplayFor` (Go, `displayAt`, test `TestDisplayAt`), frontend gọi trong `openOverlay` (`setEditingBounds` từ scriptStore); nút Màn 1 / 2 trên overlay vẫn đổi được, lần mở sau lại theo vùng
  - [ ] Kiểm tra tay macOS 2 màn: đổi màn trên overlay, đặt điểm / chọn vùng ở màn phụ, Xong → main về đúng chỗ cũ
  - [x] Dừng khẩn cấp nhận góc của mọi màn hình: `displays()` (`input.go`, Windows nhân `robotgo.ScaleF` cho cùng hệ toạ độ với `robotgo.Location`), `hardCorners` bỏ góc tiếp giáp màn khác; test `TestHardCorners`. Máy dev: màn ngoài 1920×1080 + MacBook 1512×982 tại (1920, 406) → 6 góc cứng
  - [x] Overlay đo **tỷ lệ + độ lệch** (`coordMap.ts`, bình phương nhỏ nhất trên các mẫu chuột đứng yên) thay vì chỉ độ lệch – để scale ≠ 100% trên Windows vẽ vùng / marker đúng chỗ; điểm đặt / kéo marker vẫn lấy thẳng từ `MousePos`
  - Ghi chú: robotgo trên Windows tự nhân toạ độ chuột với scale màn chính (`Location`) và chia lại khi `Move` → lấy toạ độ bằng robotgo rồi click lại luôn khớp; chỉ phép đổi sang px trong cửa sổ mới cần tỷ lệ
  - [ ] Kiểm tra tay macOS: dừng khẩn cấp ở góc màn phụ, góc tiếp giáp 2 màn không dừng nhầm; overlay hiện "tỷ lệ 1.00"
  - [ ] Kiểm tra trên Windows ở scale 100% / 125% / 150% và 2 màn khác scale (cần máy Windows – làm cùng bước 18)
- [ ] **18. Build & thử máy sạch** – `wails build` ra .exe, chạy trên máy chưa cài Go/Node.

## Phase 2 – Tiện ích

- [x] **19. Ngẫu nhiên hoá** – Delay ±ms và vị trí ±px cho từng bước.
  - [x] Model: `jitterMs` (± ms cho `delayAfter` và `clickInterval`), `jitterPx` (± px cho mỗi lần di chuột tới điểm, A / B kéo thả lệch riêng)
  - [x] Engine: `jitter` (toạ độ, cho phép âm – màn phụ) / `jitterMs` (không âm), phân bố đều; nguồn ngẫu nhiên `intN` thay được trong test; mỗi click lệch riêng
  - [x] Form: ô "Ngẫu nhiên thời gian" / "Ngẫu nhiên vị trí" kèm giải thích khoảng thực tế; danh sách hiện "±ms"; overlay vẽ ô vùng lệch quanh marker
  - [x] Test `TestJitter`, `TestRunJitterPositionsInBounds`
  - [x] Kiểm tra tay: chạy kịch bản có ± ms / ± px, thấy click lệch trong vùng và thời gian thay đổi
- [x] **19b. Xem trước kịch bản** *(người dùng yêu cầu, 2026-10-15)* – Nút *Xem trước* mở overlay chế độ `preview`: con trỏ giả màu hồng chạy 1 vòng theo thứ tự + thời gian thật của các bước, không điều khiển chuột thật.
  - [x] `components/Simulation.tsx`: `useSimulation` (huỷ bằng token, tốc độ 1× / 2× / 4×) + `SimulationLayer` (vệt đường đi hồng nét đứt / cam khi kéo, gợn click, vòng lớn dần khi giữ, bong bóng chờ / gõ chữ / phím tắt); vị trí lấy từ `endpoints()` nên khớp "từ điểm trước" / "lệch so với A"
  - [x] Overlay: tool `preview` – marker không sửa được, ẩn crosshair, thanh công cụ Chạy thử / Dừng / tốc độ / Sửa điểm, Esc dừng mô phỏng trước; từ chế độ sửa có nút *Xem trước*
  - [x] Khung hình dùng `setTimeout` thay `requestAnimationFrame` (rAF không chạy ở Chrome headless khi thử); đã chụp headless thấy chạy hết 1 vòng
  - [x] Kiểm tra tay trong app
- [x] **19c. Xem / sửa riêng một bước trên overlay** *(người dùng yêu cầu, 2026-10-15)* – Nút *👁 Xem trên màn hình* ở đầu form (bước có vị trí) mở overlay chỉ làm rõ bước đó.
  - [x] `openOverlay(tool, focus)`; overlay chế độ xem riêng: nền tối hơn, vùng + bước khác mờ / xám / không bấm được, đường nối khác mờ, marker bước đang xem nổi lên; thanh công cụ: tên bước, *Chọn lại trên màn hình* (ở lại overlay – `startPick`), *Xem tất cả*, *Xong*; ô thông tin hiện toạ độ A / B
  - [x] Click chỗ trống / F6 không tạo bước mới khi xem riêng (`setFocusLock`), trừ lúc đang chọn lại; Esc huỷ chọn lại trước rồi mới thoát
  - [x] Kiểm tra tay

## Phase 3 – Design v2 & tự động hoá

> Thiết kế: **`design/DESIGN.md`** (token, thành phần, icon, map artboard → file), `design/theme.css` (token Tailwind v4), artboard tham khảo `design/canvas/*.dc.html` + canvas https://claude.ai/artifact/5WdDyFZdB3WKxEXwkticnw. Mỗi việc UI dưới đây ghi artboard tương ứng trong ngoặc vuông.
> Số bước giữ nguyên để khớp nhật ký nên không theo thứ tự số; **làm theo thứ tự từ trên xuống**. Xong 26 là chạy được ví dụ Facebook (cuối phase).

**Thứ tự làm** (từ trên xuống, giống danh sách bên dưới):

| # | Bước | Vì sao ở vị trí này |
|---|---|---|
| 1 | 27. UI v2 cho màn đang có | Người dùng muốn làm giao diện trước; nền cho mọi màn sau |
| 2 | 27b. Khung UI tính năng chưa có | Dựng hết UI theo design một lần cho đồng bộ; bước sau chỉ nối logic |
| 3 | 29. Neo 4 góc | Sửa toạ độ lệch khi đổi vùng; đổi model `Point` nên làm trước khi có file lưu (21) |
| 4 | 21. Lưu / tải | App chưa lưu được gì – tắt là mất kịch bản |
| 5 | 30. Sửa nhiều điểm | Dựa vào neo (29) + UI mới (27); có kịch bản lưu được (21) mới đáng đầu tư chỉnh sửa hàng loạt; kèm hoàn tác / làm lại |
| 6 | 21b. Input | Cần 21 (preset theo từng kịch bản); mở đường cho chuỗi web / app (23, 25) |
| 7–8 | 23 → 25 | Cuộn / mở URL đơn giản, rồi chờ ảnh – cốt lõi tương tác web / app |
| 9 | 25e. Tìm theo chữ (OCR) | Xong 2026-10-17 – tìm ảnh chưa đủ chính xác |
| 10 | 25b. Khi thất bại | Xong 2026-10-17; 20 thay *bỏ qua / nhảy tới* bằng cạnh `fail`, giữ phần thử lại |
| 11 | 20. Flow kiểu n8n (F0–F7) | Người dùng cần luồng thành công / thất bại, if / else, switch (2026-10-17); đổi model lớn nên làm trước các loại bước mới (26, 25d); gồm luôn debug + lịch sử chạy (thay 25c) |
| 12 | 26. Focus cửa sổ | Xong 26 là chạy được ví dụ Facebook |
| 13 | 25d. Web qua DOM | Người dùng yêu cầu (2026-10-17): web bền hơn ảnh mẫu; cần focus (26), cổng `fail` (20); làm trước MCP vì AI dùng DOM tốt hơn ảnh chụp |
| 14 | 22. Hẹn giờ / CLI / API local | Chạy không cần ngồi canh; dựng server HTTP + token mà MCP dùng lại |
| 15–17 | 31a → 31b → 31c. MCP cho AI | Người dùng chọn làm sau (2026-10-17); các bước trước đã có mục *Chuẩn bị cho MCP* nên lúc này chủ yếu là bọc hàm Go sẵn có thành tool |
| 18–19 | 24 → 28 | Ghi macro, mở rộng |
| – | 17 (Windows), 18 | Khi có máy Windows; riêng *Kiểm tra tay macOS* của 17 làm lúc nào cũng được |

- [x] **27. UI v2 cho các màn đang có + theme sáng / tối** – Đổi giao diện mọi thứ đang chạy được sang design v2, không đổi logic.
  - [x] Theme: thay `@theme` trong `style.css` bằng `design/theme.css` (giữ `@import "tailwindcss"`); `data-theme="light|dark"` trên `<html>`, không đặt thì theo hệ điều hành
  - [x] Font: cài `@fontsource/plus-jakarta-sans`, `@fontsource/jetbrains-mono`, bỏ Be Vietnam Pro / IBM Plex Mono (Windows ưu tiên Segoe UI Variable / Cascadia Mono có sẵn)
  - [x] Đổi class theo bảng "Đổi tên token" (DESIGN.md): `ground` → `bg`, `accent-strong` → `accent-text`; bỏ `accent-light`, `drag*` (kéo thả dùng màu trung tính, marker overlay dùng `accent`), bỏ `dark-*` (MiniBar, PreviewScreen dùng token chung)
  - [x] Utility gradient chỉ ở 4 chỗ: `btn-primary` (Chạy, Dùng ảnh này, Xong), `ring-selected` (bước / mục đang chọn), `ring-focus` (ô nhập focus), vòng tiến độ SVG thanh mini
  - [x] `components/icons.tsx`: **toàn bộ** bảng "Icon bước" (SVG 24×24, nét 1.8) – cả loại chưa có (scroll, open, waitImage, clickImage, waitPixel, focusWindow, group, ghim)
  - [x] Thanh tiêu đề tự vẽ kiểu Windows 11 (thu nhỏ / phóng to / đóng bên phải, vẫn kéo được cửa sổ) [`Main`]
  - [x] Bố cục 3 cột [`Main`]: chỗ cho thanh bên thư viện rộng 236 (nội dung ở 27b), danh sách bước co giãn, inspector 320–360 có tab Bước / Input (tab Input ở 27b)
  - [x] Thẻ bước [`Main`]: `bg-surface border-line rounded-card`, padding 9×12, cách 6 px, ô icon 32×32 `bg-soft`; đang chọn `ring-selected` + ô icon `bg-accent-soft text-accent-text`; chip thời gian mono; bước tắt `opacity-50`
  - [x] Inspector (`StepForm.tsx`, `RunSettings.tsx`): ô nhập `bg-field` bo 6–7, focus `ring-focus`, công tắc 40×20, segmented `bg-soft` / `bg-raised`, nút phụ `bg-surface border-line`
  - [x] Thu gọn / mở rộng cột thư viện *(người dùng yêu cầu, 2026-10-16)*: nút ở đầu cột hoặc Ctrl/Cmd+B, thu gọn còn dải 44 px, nhớ trạng thái trên máy (localStorage)
  - [x] Tab Soạn / Xem trước (gạch dưới 3px `bg-accent` rộng 50%) thay nút *Xem trước*
  - [x] Bảng thêm bước Ctrl+K [`AddStep`, `AddStepDark`]: `components/AddStepPalette.tsx` thay menu *+ Thêm bước* trong `StepList.tsx` – ô tìm, 4 nhóm, phím mũi tên + Enter; liệt kê **mọi** loại bước, loại chưa có hiện nhãn *Sắp có* và không chọn được
  - [x] Thanh mini [`MiniBar`, `MiniBarDark`]: `bg-surface border-line rounded-hud`, vòng tiến độ gradient, 2 dòng mono, nút tạm dừng phụ, nút Dừng `btn-danger-soft`
  - [x] Overlay [`Overlay`]: lớp phủ `bg-dim`, khung chọn nét đứt trắng + 4 núm viền accent + nhãn kích thước nền accent; thanh công cụ dưới `bg-hud text-hud-text border-hud-line rounded-xl`, công cụ đang bật nền trắng 14%; marker dùng `accent`
  - [x] Cài đặt theme Sáng / Tối / Theo hệ thống, lưu vào cấu hình app (`UserConfigDir`)
  - [x] Cập nhật mục *Quy ước UI* trong `CLAUDE.md` theo token / font mới
  - Xong khi: mọi màn hiện tại đúng ở cả 2 theme (chụp headless `#`, `#overlay`, `#minibar`, `#preview`), `npm run build` sạch
  - [x] Chụp headless cả 2 theme: màn chính (1280 / 1024), inspector click / kéo thả / gõ chữ / phím tắt, popup lặp, bảng Ctrl+K (tìm không dấu, phím mũi tên, Enter chèn sau bước đang chọn), overlay chọn vùng / xem trước, mini bar
  - [x] Kiểm tra tay trong `wails dev`: kéo cửa sổ bằng thanh tiêu đề, thu nhỏ / phóng to / đóng, phóng to rồi mở overlay → quay về vẫn phóng to; mini bar 540×96 khi chạy; overlay thanh công cụ dưới đáy; đổi theme, tắt app mở lại vẫn giữ; Ctrl/Cmd+K thêm Double click / Giữ chuột qua overlay
- [x] **27b. Dựng sẵn UI cho tính năng chưa có** *(người dùng yêu cầu, 2026-10-16)* – Làm hết giao diện theo design v2 trước, logic làm ở bước sau.
  - Cách làm: component nhận dữ liệu qua props, chưa nối Go. Trong app: hiện đúng chỗ nhưng nút / ô bị tắt, nhãn *Sắp có* (tooltip ghi số bước sẽ làm). Màn **`#demo`** nạp kịch bản mẫu (cả loại bước chưa có, dữ liệu giả trong `frontend/src/demo/`, **không** thêm vào `types.ts` / `model.go`) để xem đủ trạng thái, so với canvas
  - [x] Thư viện thanh bên [`Main`] – `components/Library.tsx`: ô tìm, danh sách kịch bản, mục đang mở (`bg-hover` + vạch trái accent), chấm `warn` chưa lưu, "Mở gần đây", nút Kịch bản mới / Mở / Lưu, mục "Lịch sử"; nhãn "Chưa lưu" cạnh tên kịch bản; status bar "Bản nháp tự lưu lúc …" *(logic: 21, 25c)*
  - [x] Input [`StepInputs`, `InputsTab`]: `PinButton.tsx`, trường đã ghim `field-pinned` + chip key; `VarTextInput.tsx` tô `{{biến}}` + popup gợi ý; `InputsPanel.tsx` (tab Input, dòng input, dạng mở rộng để sửa); chip "N input" trên thẻ bước *(logic: 21b)*
  - [x] Form Chạy [`RunDialog`, `RunDialogDark`] – `RunDialog.tsx`: preset, ô theo bước, ô đã đổi + khôi phục, khung "Sẽ gõ", chân với "Lần sau chạy luôn" / Đếm ngược / Huỷ / Chạy *(logic: 21b)*
  - [x] Neo điểm: ô X / Y dạng ±, lưới 3×3 chọn neo + tuỳ chọn %, toạ độ `(−300, +300)` trong danh sách, đường neo trên overlay *(chưa có trong design – vẽ theo thành phần v2; logic: 29)*
  - [x] Sửa nhiều điểm trên overlay: khung biến đổi + núm (co giãn / xoay / xiên), thanh công cụ ΔX / ΔY / % / góc / xiên, nút căn chỉnh / lật, hộp hỏi "Ánh xạ sang vùng mới" *(chưa có trong design; logic: 30)*
  - [x] Inspector cho loại bước mới: cuộn (hướng, số nấc), mở URL / app, **Chờ ảnh** [`Main`] (ảnh mẫu, Chụp lại, Xem trên màn hình, độ khớp, chờ tối đa, kiểm tra mỗi), click theo ảnh (ΔX / ΔY), chờ pixel (điểm, màu, sai số), focus cửa sổ (ô nhập + danh sách cửa sổ) *(logic: 23, 25, 26)*
  - [x] Khi hết giờ [`Main`]: segmented Thử lại / Bỏ qua / Dừng, số lần thử, "Thử hết vẫn không thấy"; chip trên thẻ bước *(logic: 25b)*
  - [x] Overlay tool `capture` [`Overlay`]: khoanh ảnh mẫu, nhãn kích thước, "Dùng ảnh này" / "Khoanh lại" *(logic: 25)*
  - [x] Thanh mini trạng thái chờ [`MiniBar`]: vòng theo thời gian chờ, "đang tìm x / y s", "thử lại x / y", tên preset *(logic: 25, 25b, 21b)*
  - [x] Khung nhóm bước [`Main`]: nét đứt `border-line-strong bg-group rounded-group`, icon + tên + chip `lặp ×N` *(logic: 20)*
  - [x] Debug: nút Chạy từ bước này / Tiếp (F10), breakpoint trên thẻ bước, màn timeline lịch sử chạy *(chưa có trong design; logic: 25c → nay là 20 F5)*
  - [x] Hẹn giờ (chip giờ trong thư viện, hộp đặt giờ), nút Ghi macro trên toolbar *(chưa có trong design; logic: 22, 24)*
  - Xong khi: `#demo` hiện đủ các phần trên ở cả 2 theme, khớp canvas; app thật không có nút nào bấm được mà không chạy; `npm run build` sạch
  - Kết quả: màn mẫu `#demo` (màn soạn đủ loại bước + nhóm, thư viện, inspector từng loại, tab Input, form Chạy, hẹn giờ, ghi macro, debug), `#demo-overlay` (khoanh ảnh, chọn nhiều + biến đổi + ánh xạ), `#demo-minibar` (các trạng thái thanh mini). Code dùng lại khi làm logic: `components/StepCard.tsx`, `components/inputs/*`, `components/steps/*` (model + `StepFields` + `AnchorField`), `ScheduleDialog.tsx`, `AppChrome.tsx`, `MiniBarView`, `LibraryView`. Dữ liệu mẫu chỉ trong `frontend/src/demo/`
  - [x] Người dùng xem `#demo`, `#demo-overlay`, `#demo-minibar` và màn chính (nút tắt + *Sắp có*), góp ý trước khi commit
- [x] **29. Neo điểm theo 4 góc vùng** *(người dùng yêu cầu, 2026-10-16)* – Hiện điểm chỉ tính từ góc trên trái của vùng nên điểm ở phía phải / dưới lệch khi đổi vùng. Mỗi trục neo riêng: X neo trái (+) hoặc phải (−), Y neo trên (+) hoặc dưới (−). Ví dụ 4 điểm cách mỗi cạnh 300 px = (+300, +300), (−300, +300), (+300, −300), (−300, −300); đổi vùng thì mỗi điểm bám theo góc của nó.
  - [x] Model: `Point` thêm `ax` (`start` = trái | `end` = phải | `center` | `ratio` = theo % rộng), `ay` (tương tự cho trên / dưới); bỏ trống = `start` (file / kịch bản cũ chạy như cũ). Lưu offset dương + neo riêng thay vì dùng dấu (để phân biệt +0 / −0); UI hiện dạng có dấu ±
  - [x] Toàn màn hình (`region = nil`) cũng neo được: chốt **giữ `nil`**, lúc tính thì khung = màn hình chính (`mainScreen`, Go; frontend lấy qua `ScreenBounds`) – kịch bản mang sang máy khác cỡ màn hình vẫn bám góc. Đa màn hình xử lý ở 17
  - [x] Go: một hàm `resolvePoint(p, region)` dùng chung cho engine (`runStep.pt`) và `validateScript` (điểm phải nằm trong vùng); test 4 góc, center, ratio, đổi cỡ vùng, kịch bản cũ không có neo
  - [x] Frontend: `toScreen` / `toRelative` / `endpoints()` trong `scriptStore.ts` tính theo neo (khớp Go); đổi neo giữ nguyên chỗ trên màn hình (tính lại offset)
  - [x] Neo khi đặt điểm: mặc định **tự chọn góc gần nhất** theo từng trục; kéo marker / chọn lại giữ neo hiện tại
  - [x] Cài đặt kịch bản *góc gần nhất / luôn trên trái* – `Script.anchorMode`, chọn trong popup nút "N vòng" (làm cùng 21)
  - [x] Form: ô X / Y có dấu ± – gõ số âm = neo phải / dưới; nút chọn neo dạng lưới 3×3 nhỏ (4 góc, giữa cạnh, tâm) + tuỳ chọn %; danh sách bước hiện toạ độ dạng `(−300, +300)`
  - [x] Overlay: đường mảnh từ marker tới cạnh / góc đang neo (hiện với bước đang chọn hoặc khi giữ Alt); đổi vùng (kéo vùng mới) thấy marker bám theo góc ngay
  - [x] Kéo thả "lệch so với A" (ΔX, ΔY) giữ nguyên – lệch không phụ thuộc neo
  - [x] Nối UI khung từ 27b (ô ±, lưới chọn neo, đường neo trên overlay)
  - Xong khi: đặt 4 điểm ở 4 góc vùng 1000×700, đổi vùng sang 1400×900 → 4 điểm vẫn cách mỗi cạnh đúng 300 px khi chạy thật và trên overlay
- [x] **21. Lưu / tải kịch bản** – App chưa lưu gì; thêm `version` trước khi đổi model nhiều.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Model: `Script.version` (`ScriptVersion = 1`), `Script.anchorMode`; file `.json`
  - [x] Go `storage.go`: `SaveScript` / `SaveScriptAs` / `OpenScript` / `OpenScriptPath` qua `runtime.SaveFileDialog` / `OpenFileDialog`; mở file kiểm tra cấu trúc (`decodeScript` + `checkScript`: version, loại bước, số điểm, neo, số âm; bổ sung slice trống, id trống / trùng); mở gần đây (≤ 5, trong `settings.json`, file mất thì bỏ) + bản nháp `draft.json` (`SaveDraft` / `LoadDraft` / `ClearDraft`) trong `os.UserConfigDir()/autoclicker/`; ghi file tạm rồi đổi tên
  - [x] Frontend `scriptFile.ts`: cờ `dirty` + `path` trong store (mọi `updateScript` bật cờ), tự ghi nháp sau 800 ms, mở app khôi phục nháp (không có thì kịch bản trống – app thật không còn dữ liệu mẫu)
  - [x] Thư viện ở thanh bên [`Main`]: ô tìm (lọc cả mở gần đây), chấm `warn` khi chưa lưu, "Mở gần đây" bấm mở được, nút Kịch bản mới / Mở / Lưu
  - [x] Nhãn "Chưa lưu" theo cờ thật, rê tên kịch bản thấy đường dẫn; status bar "Bản nháp tự lưu lúc …", lỗi lưu / mở; chấm `ok` "Sẵn sàng"
  - [x] Phím Cmd/Ctrl + S, + Shift + S, + O, + N; hộp thoại `SaveChangesDialog` (Lưu / Không lưu / Huỷ, Enter / Esc) khi kịch bản mới, mở khác, đóng app (`OnBeforeClose` → event `app:confirmClose`, `ConfirmQuit`)
  - [x] Test `storage_test.go`: ghi rồi đọc lại khớp, file hỏng / thiếu / sai `version` / giá trị lạ báo lỗi, bổ sung mặc định, nháp, mở gần đây, tên file an toàn
  - [x] Kiểm tra tay trên macOS: lưu / lưu thành / mở / mở gần đây, sửa rồi tắt app mở lại còn nháp, đóng app khi chưa lưu hỏi lưu, file hỏng báo lỗi
  - Xong khi: tắt app mở lại vẫn còn bản nháp; mở / lưu qua hộp thoại trên macOS
- [ ] **30. Sửa nhiều điểm cùng lúc: di chuyển, co giãn, xoay, xiên** *(người dùng yêu cầu, 2026-10-16)* – Chọn nhiều marker trên overlay (chế độ sửa và tab Xem trước) rồi biến đổi cả nhóm để khớp vùng mới.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Chọn nhiều (công cụ *Chọn nhiều* trên overlay thật): Shift / Cmd+click marker, kéo khung chọn trên chỗ trống, Cmd/Ctrl+A chọn tất cả, Esc bỏ chọn (Esc lần nữa mới thoát); danh sách bước viền accent bước có điểm đang chọn, Shift / Cmd+click thẻ thêm / bỏ (`state.multi`, khoá `id:index`)
  - [x] Khung biến đổi (`TransformTool.tsx`): kéo thân / marker đã chọn = **di chuyển** (Shift chỉ ngang / dọc; phím mũi tên 1 px, Shift 10 px); 8 núm = **co giãn** (Shift giữ tỷ lệ, Alt từ tâm); núm tròn = **xoay** (Shift bước 15°); Cmd/Ctrl + núm cạnh = **xiên**; Cmd/Ctrl + núm góc = **phối cảnh**; Esc khi đang kéo = huỷ; xem trước khi kéo, thả mới ghi; khung đỏ khi có điểm ra ngoài vùng
  - [x] Thanh ô nhập số: ΔX, ΔY, rộng / cao %, xoay, xiên X / Y – xem trước ngay, Áp dụng / Enter ghi; thanh tự né khung (dưới → trên → sát đáy), kéo tay nắm "N điểm" để dời
  - [x] Căn trái / giữa / phải / trên / giữa / dưới, chia đều ngang / dọc, lật ngang / dọc (`align` trong `transform.ts`, `moveSelectedTo`)
  - [x] **Ánh xạ khi đổi vùng** (kéo vùng mới / *Toàn màn hình* khi đã có điểm): hộp chọn *Giữ theo neo* (mặc định) / *Co giãn theo vùng* / *Giữ nguyên chỗ trên màn hình*, marker xem trước theo cách đang chọn, báo số điểm sẽ ra ngoài vùng; Enter đổi, Esc giữ vùng cũ
  - [ ] Ánh xạ theo 4 điểm mốc cũ → mới (đã có `homography` – còn thiếu UI chọn mốc). Tạm thời: chọn hết điểm rồi Cmd/Ctrl + kéo góc khung để biến dạng phối cảnh
  - [x] Toán `frontend/src/transform.ts`: ma trận 3×3 (affine + homography), `fromParams`, `perspectiveCorner`, `bbox`, `align`, `rectToRect`; test `transform.test.ts` (vitest, `npm test`) 14 test; sau biến đổi ghi lại offset theo neo hiện có
  - [x] Điểm đặc biệt (`mapPoints` trong `scriptStore.ts`, dùng chung cho biến đổi, xem trước và ánh xạ vùng): A "từ điểm trước" không chọn được; B "lệch so với A" đổi véc-tơ = f(A + Δ) − f(A); điểm ra ngoài vùng → báo, không áp dụng
  - [x] Hoàn tác / làm lại Cmd/Ctrl+Z, Cmd/Ctrl+Shift+Z (Ctrl+Y) cho **mọi** thay đổi kịch bản (`past` / `future` trong store, tối đa 100, gộp thay đổi sát nhau < 500 ms; trong ô nhập để trình duyệt tự hoàn tác chữ); mở / tạo kịch bản xoá lịch sử
  - [x] Tab Xem trước: nút *Sửa điểm* chuyển sang công cụ *Chọn nhiều*, xong bấm *Xem trước* chạy lại theo vị trí mới
  - [ ] Kiểm tra tay trong app
  - Xong khi: chọn 5 điểm, co giãn 150%, xoay 10°, xiên, lật → overlay và lần chạy thật khớp nhau; Ctrl+Z quay lại đúng từng bước; đổi vùng với cả 3 cách xử lý cho kết quả như mô tả
- [x] **21b. Input của kịch bản** – Người soạn đánh dấu giá trị nào được đổi khi chạy; bấm Chạy thì điền form. Go áp giá trị vào bản sao kịch bản rồi mới đưa engine, **engine không đổi**. Cùng đổi `version` với 21; CSV (28) dựa vào đây.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - Đã làm (2026-10-17): `inputs.go` (`ScriptInput`, `resolveInputs`, `ScriptError`, `GetInputMemory` / `SaveInputMemory` → `inputs.json`), file kịch bản `version` 2; `RunScript` / `ValidateScript` nhận `values`. Frontend: `scriptInputs.ts` (ghim / dùng chung / bỏ ghim, sửa input, bước dùng input, nhóm form Chạy), `runForm.ts` (Chạy / F8 → form, chạy luôn, preset), `PinMenu`, `NumField` + tab Input trong `StepForm.tsx`
  - Khác plan: kiểu `point` (ghim toạ độ) chưa làm – nút ghim cạnh toạ độ vẫn *Sắp có*; ghim `enabled` (bool) có ở Go nhưng chưa có nút trên UI. Ghim vào bước đã xoá được giữ lại (hoàn tác xoá bước thì ghim quay lại), lúc chạy bỏ qua. Trường đã ghim luôn bằng mặc định của input nên Xem trước dùng mặc định, không gọi resolve
  - [x] Model `Script.inputs`: `key`, `label`, `type` (`number`, `ms`, `text`, `multiline`, `bool`, `select`, `point`), `default`, `required`, `min` / `max`, `options`, `hint`, `targets: [{stepId, field}]`. Ví dụ:
    ```json
    "inputs": [
      { "key": "click_delay", "label": "Chờ sau khi click", "type": "ms", "default": 500,
        "targets": [{ "stepId": "s4", "field": "delayAfter" }] },
      { "key": "M1", "label": "Tên người nhận", "type": "text", "required": true, "hint": "Tên hiện trong lời chào" }
    ]
    ```
    - Nguồn: **trường được ghim** (`targets`, một input cho nhiều bước) và **biến `{{key}}`** trong trường chữ (gõ chữ, URL, tên cửa sổ) – tự tạo input khi gõ biến mới
    - Key: chữ không dấu, số, `_`; nhãn tiếng Việt ở `label`; `\{{` để gõ ra `{{`; không ghim được trường tự tính (kéo thả "từ điểm trước")
  - [x] Go: `resolveInputs(script, values) (Script, error)` – sao chép, ghi trường ghim, thay `{{…}}`, `validateScript`; `RunScript(sc, values, countdown)`; Xem trước dùng chung resolve
  - [x] Lỗi rõ: thiếu bắt buộc, sai kiểu / ngoài khoảng, biến không có trong `inputs`, input không bước nào dùng, toạ độ ngoài vùng
  - [x] Giá trị lần trước + preset lưu trong `UserConfigDir` theo từng kịch bản, **không** ghi vào file kịch bản
  - [x] Inspector [`StepInputs`]: `components/PinButton.tsx` (32×32 cạnh mỗi trường số / chữ; đã ghim nền `accent-soft`); trường đã ghim `field-pinned` + chip key `var-token` cạnh nhãn
  - [x] Ô chữ có biến [`StepInputs`]: `components/VarTextInput.tsx` – tô `{{biến}}` bằng `var-token`; gõ `{{` mở popup gợi ý (key mono `accent-text` + nhãn `muted`, "+ Tạo input mới…"); dòng phụ "N biến, tự thêm vào tab Input"
  - [x] Tab Input [`InputsTab`]: `components/InputsPanel.tsx` – chip số input ở tab; mỗi dòng ⋮⋮ kéo, nhãn + `*`, `key · mặc định` mono, chip kiểu, chip bước dùng; input đang sửa mở rộng (`ring-selected`): nhãn, kiểu, mặc định, bắt buộc, gợi ý; bỏ ghim; "+ Thêm input"
  - [x] Thẻ bước: chip viền `line` "N input"
  - [x] Form Chạy [`RunDialog`, `RunDialogDark`]: `components/RunDialog.tsx` (popup 560, bo 12) – đầu có tên kịch bản + số input; "Bộ giá trị" (select preset + "Lưu bộ mới"); ô nhóm theo bước, ô đã đổi `ring-focus` + "đã đổi, mặc định … · khôi phục", `*` màu `warn`; khung "Sẽ gõ" xem trước câu sau khi thay biến; chân: "Lần sau chạy luôn với giá trị này", "Đếm ngược 3 giây", Huỷ / Chạy (`btn-primary`, F8)
  - [x] Nút Chạy / F8 mở form khi kịch bản có input; chạy thẳng nếu đã bật "Lần sau chạy luôn"; thanh mini hiện tên preset
  - [x] Chuẩn bị cho MCP (31): `resolveInputs` và kiểm tra input viết ở Go (không ở frontend) để CLI / MCP gọi được; `validateScript` trả lỗi có cấu trúc `{stepId, field, message}` thay vì chỉ chuỗi
  - [x] Test: ghi đè trường; thay biến (có tiếng Việt); thiếu biến; sai kiểu; escape; một input nhiều bước; input trong nhóm (sau bước 20) – Go `inputs_test.go` (trừ input trong nhóm – làm ở 20); vitest `scriptInputs.test.ts` (ghim, dùng chung, đổi key, nhóm form Chạy)
  - [x] Kiểm tra tay trong app
  - Xong khi: B1 click 500 → chạy với 1000 ms; B2 "xin chào {{M1}} và {{M2}}" gõ đúng giá trị điền
- [x] **23. Cuộn + Mở URL / app** – Loại bước đơn giản.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Model: `scroll` (`direction`, `amount` – ghim được, điểm tuỳ chọn trong `points` nên có neo); `open` (`target` URL / đường dẫn / tên app, `app` tuỳ chọn) – trường `omitempty`, file cũ không đổi, không tăng `version`
  - [x] Engine: `Input.Scroll` (`robotgo.ScrollDir` từng nấc, nghỉ 30 ms – dừng được giữa chừng); `open.go`: `openCommand` (hàm thuần, có test) – macOS `open [-a App] target` (tên app → `open -a`), Windows PowerShell `Start-Process` (không dùng `cmd /c start` vì `&` trong URL bị hiểu là nối lệnh), Linux `xdg-open`; tên miền tự thêm `https://`; biến trong URL tự mã hoá (`resolveInputs`); test với input giả + hàm mở giả
  - [x] UI: bật Cuộn / Mở URL trong bảng Ctrl+K; inspector: hướng + số nấc + cuộn tại chuột / điểm (bản đồ neo, chọn trên màn hình, xem riêng trên overlay); ô URL nhận `{{biến}}` + chọn trình duyệt; Xem trước có bong bóng cuộn / mở
  - [x] Chuẩn bị cho MCP (31): `steptypes.go` – `stepTypeDocs` (tên, mô tả, trường, trường ghim được, ví dụ JSON) cho mọi loại bước + `App.StepTypes()`; test kiểm tra mọi loại có mô tả và ví dụ qua được `checkStep`
  - [x] Thể hiện cuộn tới đâu *(người dùng yêu cầu, 2026-10-17)*: quy ước 1 nấc ≈ 100 px (`SCROLL_NOTCH_PX` / `scrollNotchPx`); inspector có `ScrollMap` (khung vùng / màn hình theo tỷ lệ, điểm cuộn hoặc chuột, mũi tên đúng quãng, dải "hiện thêm ≈ N px", % chiều cao); overlay vẽ mũi tên + nhãn từ marker; Xem trước: mũi tên hồng chạy dần + bong bóng; thẻ bước ghi ≈ px
  - [x] Kiểm tra tay macOS (người dùng xác nhận 2026-10-17); Windows (PowerShell `Start-Process`, cuộn theo nấc) kiểm tra khi có máy – ghi ở 17. Chi tiết đã thử: macOS (mở URL có biến tiếng Việt bằng Chrome / mặc định, mở app theo tên, cuộn tại chuột / điểm); Windows (PowerShell `Start-Process`) khi có máy
- [ ] **25. Chờ ảnh, click theo ảnh, chờ pixel** – Cốt lõi để tương tác web / app (làm sớm).
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Model: `waitImage` (`template` **PNG base64 ngay trong file kịch bản** + `templateScale`, `threshold` %, `timeoutMs`, `intervalMs`, `searchArea` vùng / cả màn chính); `clickImage` (`imageSource` prev = kết quả bước chờ ảnh trước / search = tự tìm, `dx` / `dy`); `waitPixel` (điểm có neo, `color`, `tolerance`, timeout). Ghim được `timeoutMs`, `threshold`
  - [x] Engine (`imagestep.go`): chụp qua interface `Screen` (`capture.go`, robotgo), so khớp `match.go` – **Go thuần, không dùng gocv**: NCC trên ảnh xám, tìm thô trên ảnh thu nhỏ (giữ 6 ứng viên) rồi tìm mịn; mẫu một màu so độ chênh. Màn Retina 3024×1964: ~70 ms (benchmark), máy thật: chụp ~270 ms + tìm ~160 ms cả màn. Chỗ thấy nhớ cho click theo ảnh; tiến độ chờ `Progress.wait`; hết giờ → dừng kèm lỗi "không thấy ảnh mẫu sau 10 s (khớp nhất 46%, cần 90%)" (thử lại / bỏ qua: 25b); thời gian chờ không tính lúc tạm dừng
  - [x] Overlay tool `capture`: kéo khung, tối ngoài khung, nhãn cỡ thật, "Dùng ảnh này" (Go ẩn cửa sổ 0,25 s rồi chụp – `CaptureTemplate`) / "Khoanh lại" / Esc; Ctrl+K "Chờ ảnh" thêm bước rồi mở khoanh luôn; khung kết quả "Thử tìm" vẽ trên overlay
  - [x] Inspector: ảnh mẫu (cỡ theo toạ độ), Khoanh lại, **Thử tìm** (`FindImage` – thấy ở đâu, % khớp, ms, xem trên màn hình), độ khớp, chờ tối đa (ghim được), kiểm tra mỗi, tìm trong vùng / cả màn; click theo ảnh: nguồn "Kết quả bước N" / tự tìm, ΔX ΔY, nút chuột; chờ pixel: điểm, màu + **Lấy màu** (`PixelColorAt`), sai số, chờ tối đa; báo thiếu quyền Screen Recording + nút "Xin quyền" (`capture_darwin.go`, cgo `CGPreflightScreenCaptureAccess`)
  - [x] Bật 3 loại trong bảng Ctrl+K; Xem trước có bong bóng chờ ảnh / click theo ảnh / chờ màu
  - [x] Thanh mini: vòng tiến độ theo thời gian chờ, "đang tìm x / y s · khớp nhất n%"
  - [x] Retina / DPI: ảnh mẫu lưu kèm tỷ lệ lúc chụp, tìm trên màn tỷ lệ khác thì co giãn mẫu (test 2x ↔ 1x)
  - [x] Chuẩn bị cho MCP (31): `CaptureTemplate(rect)`, `FindImage`, `PixelColorAt`, `matchScaled` là hàm Go độc lập với overlay, tạo ảnh mẫu từ toạ độ vùng; `stepTypeDocs` có 3 loại mới kèm ví dụ PNG thật
  - ~~Giới hạn: chỉ tìm trên màn hình chính~~ → tìm được trên mọi màn (2026-10-17, xem 17); test thật `AUTOCLICKER_REAL=1 go test -run RealCapture`
  - [ ] Kiểm tra tay: khoanh nút trên trang web → Thử tìm → chạy chờ ảnh + click theo ảnh khi cửa sổ dời chỗ; chờ pixel đổi màu; thanh mini khi đang chờ
- [x] **25e. Tìm theo chữ (OCR)** *(người dùng yêu cầu, 2026-10-17 – tìm ảnh chưa đủ chính xác)* – Bước Chờ ảnh / Click theo ảnh (tự tìm) chọn **Tìm theo: Ảnh mẫu / Chữ**; bảng Ctrl+K có *Chờ chữ*, *Click theo chữ*.
  - [x] Go: `ocr.go` (`TextReader`, `findText`: không phân biệt hoa / thường, bỏ khoảng trắng, tuỳ chọn bỏ dấu / cả dòng, chỗ thứ N theo thứ tự đọc; khung = hợp các từ chứa chuỗi), `ocr_darwin.go` (Vision, `vi-VT` + `en-US`, mức accurate; ảnh dẹt chèn nền cho tỷ lệ ≤ 3:1), `ocr_other.go` báo chưa hỗ trợ (Windows.Media.Ocr làm cùng 17 / 18). Model `findBy`, `findText` (nhận `{{biến}}`), `textExact`, `textNoAccent`, `occurrence`; hết giờ báo các dòng đọc được
  - [x] UI: inspector ô chữ + công tắc + chỗ thứ; *Thử tìm* đọc chữ, không thấy thì hiện các dòng đọc được (bấm để dùng); thẻ bước / thanh mini / xem trước ghi chữ
  - [x] Test `TestFindText`, `TestWaitTextThenClick`, `TestWaitTextTimeout`, `TestValidateTextStep`; thật: cả màn 1920×1080 ~0,9 s, thanh menu 1920×40 thấy "Bookmarks" ~0,6 s
  - [x] Tìm ảnh chính xác hơn: khác tỷ lệ (mẫu 1x, màn Retina) thì thu nhỏ ảnh chụp thay vì phóng to mẫu; mẫu ít tương phản (chữ xám nhạt) không còn bị bước tìm thô coi là một màu; 12 ứng viên. Thử thật màn MacBook: đúng 37/127 → 106/127 (16 chỗ giống hệt). Test `TestFindTemplateLowContrast`, `TestMatchScaledDownscalesScreen`
  - [x] Sửa (người dùng báo, 2026-10-17): thêm *Click theo chữ* mà thẻ / inspector / thanh mini ghi "Click theo ảnh" → `stepLabel(s)` (`types.ts`) đặt tên theo cách tìm (Chờ chữ / Click theo chữ), dùng ở mọi chỗ hiện tên một bước; inspector Click bỏ chữ "ảnh" ở nhãn khi tìm theo chữ
  - [ ] Kiểm tra tay: kịch bản Facebook đổi sang Chờ chữ / Click theo chữ
- [x] **25b. Khi thất bại** *(làm sớm 2026-10-17, người dùng yêu cầu if / else – đây là phần nền, phần khối điều kiện ở 20)* – Mỗi bước có cài đặt *Khi không thấy*: thử lại N lần, rồi **Dừng kịch bản / Bỏ qua / Nhảy tới bước**.
  - [x] Model `onFail { retries, action: stop | next | goto, goto }` (nil = dừng như cũ); `checkStep` kiểm action / số lần, `validateScript` kiểm bước nhảy tới còn, đang bật, không phải chính nó
  - [x] Engine (`runWithRetry`): thử lại ngay (không chờ delayAfter), thất bại thì quên chỗ thấy cũ (Click theo kết quả bước trước không click nhầm chỗ cũ); bỏ qua = sang bước sau; nhảy tới bước trước / sau trong cùng vòng, tối đa 1000 lần nhảy mỗi vòng (chặn vòng lặp không dừng); tiến độ `retry` / `retries`
  - [x] UI: inspector *Khi không thấy / không đổi màu / mở không được* (Chờ ảnh, Click theo ảnh, Chờ pixel, Mở) – số lần thử lại, segmented, chọn bước nhảy tới; chip trên thẻ bước (`↻2 → bước 5`); thanh mini "thử lại x / y"
  - [x] Test `TestOnFailRetry`, `TestOnFailNextSkipsStaleClick`, `TestOnFailGoto`, `TestValidateOnFail`; mô tả cho AI ở `commonFieldsDoc` (`steptypes.go`)
  - [ ] Kiểm tra tay: chờ chữ không thấy → thử lại 2 lần → nhảy tới bước cuộn rồi quay lại
  - Không làm: `skipGroup` (bước 20 dùng *Thoát khối* thay)
- [ ] **20. Flow kiểu n8n – mỗi bước là một module, rẽ nhánh theo kết quả** *(người dùng chốt 2026-10-17, thay bản "cây khối" 20a–20c đã revert)* – Kịch bản thành **đồ thị node + cạnh**: mỗi node chạy xong trả về một `StepResult` chung (status, cổng ra, output, lỗi, log); cạnh nối từ **cổng ra** của node này sang node kế → luồng thành công / thất bại, if / else, switch case, vòng lặp đều là cạnh. Soạn trên **canvas kéo nối như n8n**. Gồm cả debug + lịch sử chạy (25c cũ, ở F5).
  - **Thứ tự làm:** F0 → **F4a** (thiết kế + `#demo-flow`, dữ liệu giả) → F1 (model) → F2 (engine) → **F4b** (nối canvas vào app, bỏ danh sách bước cũ) → F3 (node điều khiển: engine + inspector cùng lúc) → F5 (chạy + nhật ký + debug) → F6 (xem trước, thanh mini) → F7 (MCP). Lý do: duyệt giao diện sớm khi sửa còn rẻ; canvas thay danh sách ngay sau khi engine chạy được để app không bị "nửa cũ nửa mới" lâu
  - **File cũ:** bỏ, không chuyển đổi *(người dùng chốt 2026-10-17)* – `ScriptVersion` 4, file version < 4 và `draft.json` cũ báo "Kịch bản tạo bằng bản cũ, không mở được" (nháp cũ xoá luôn). Bỏ hẳn `onFail.action` / `goto` (thay bằng cạnh `fail`), không giữ code tương thích
  - **Nguyên tắc:**
    - `Step` giữ nguyên làm **node** (mọi trường bước hiện có dùng lại), `Script.steps` = danh sách node (thứ tự trong mảng không có nghĩa chạy); thêm `Script.start` + `Script.edges`
    - Cổng ra không nối: `fail` → **dừng kịch bản, báo lỗi**; `success` / cổng khác → **kết thúc thành công**. Node tắt → `skipped`, đi theo cổng `success`
    - Vòng lặp = cạnh quay lại; chặn lặp vô hạn bằng node *Đếm* + giới hạn 10 000 lần chạy node mỗi vòng
    - Node sau đọc kết quả node trước bằng biến: `{{n2.output.x}}`, `{{n2.status}}`, `{{n2.error.message}}`, `{{last.output.text}}`, `{{vars.ten}}` – `n2` là mã ngắn của node (`ref`, tự đặt n1, n2…, sửa được, duy nhất)
    - *Click theo kết quả bước trước* → chọn node nguồn (`imageFrom: "n2"`); *Kéo từ chỗ chuột dừng* lấy vị trí chuột lúc chạy
    - Lặp N vòng cả kịch bản, tạm dừng / dừng / F8, failsafe giữ nguyên
  - **Chia module** (mỗi file một việc; phần thuần có test riêng; Go và TS cùng tên khái niệm):

    | Module | Go | Frontend | Việc |
    |---|---|---|---|
    | Model flow | `flow.go` | `types.ts` (+ `flow.ts`) | `Edge`, `StepResult`, `StepError`, `portsOf(step)`, hằng cổng / status / mã lỗi |
    | Thứ tự hiển thị | `flow.go` `flowOrder` | `flow.ts` `flowOrder` | đánh số node (marker, Xem trước, form Chạy) – hai bên cùng công thức, test cùng dữ liệu |
    | Kiểm tra | `flowcheck.go` | – | `validateFlow` (cạnh, cổng, `ref`, biến, `imageFrom`, node không tới được) |
    | Engine | `flowrun.go` | – | `runFlow`, `execNode`, thử lại, giới hạn lặp, event |
    | Biến lúc chạy | `vars.go` | `flowVars.ts` | thay `{{nX.output…}}` / `{{vars…}}`, gợi ý biến cho ô nhập |
    | Node điều khiển | `nodes.go` (+ `cond.go`) | `components/flow/nodes/*` | `if`, `switch`, `switchScreen`, `counter`, `end`, `setVar`; Thử node |
    | Sửa đồ thị | – | `flowEdit.ts` | thêm / xoá / nối / cắt cạnh, chèn node giữa cạnh, copy / paste, đặt bắt đầu (thuần, có test) |
    | Xếp tự động | – | `flowLayout.ts` | `elkjs` → vị trí node |
    | Canvas | – | `components/flow/FlowCanvas.tsx`, `NodeCard.tsx`, `FlowEdge.tsx`, `PortField.tsx` | `@xyflow/react`, thẻ node, cạnh có nhãn cổng, mục Cổng ra trong inspector |
    | Nhật ký chạy | `runlog.go` | `components/flow/RunLog.tsx`, `runStore.ts` | lưu / đọc lịch sử chạy, bảng nhật ký |
    | Store | – | `scriptStore.ts` | gọi `flowEdit` + hoàn tác / làm lại (không tự sửa cạnh trong store) |
    | Mẫu | – | `demo/flowDemo.ts`, màn `#demo-flow` | dữ liệu giả đủ mọi trạng thái |

  - **Ví dụ JSON** (đích của F1–F3):
    ```json
    { "version": 4, "start": "a1",
      "steps": [
        { "id": "a1", "ref": "n1", "type": "open", "target": "facebook.com" },
        { "id": "a2", "ref": "n2", "type": "waitImage", "name": "Tìm nút Nhắn tin", "findBy": "text", "findText": "Nhắn tin", "timeoutMs": 3000 },
        { "id": "a3", "ref": "n3", "type": "clickImage", "imageFrom": "n2" },
        { "id": "a4", "ref": "n4", "type": "scroll", "direction": "down", "amount": 3 },
        { "id": "a5", "ref": "n5", "type": "counter", "max": 10 },
        { "id": "a6", "ref": "n6", "type": "end", "endStatus": "fail", "message": "Không thấy nút Nhắn tin" } ],
      "edges": [
        { "from": "a1", "port": "success", "to": "a2" }, { "from": "a2", "port": "success", "to": "a3" },
        { "from": "a2", "port": "fail", "to": "a5" },    { "from": "a5", "port": "next", "to": "a4" },
        { "from": "a5", "port": "done", "to": "a6" },    { "from": "a4", "port": "success", "to": "a2" } ] }
    ```
  - **Loại node và cổng ra:**

    | Nhóm | Node | Cổng ra | Output |
    |---|---|---|---|
    | Hành động | các loại bước hiện có (click, kéo, gõ, phím, cuộn, mở, chờ ảnh / chữ / pixel, click theo ảnh / chữ) | `success` / `fail` | chờ ảnh / chữ: `{x, y, w, h, score, text}`; chờ pixel: `{color}`; click: `{x, y}`; mở: `{target}` |
    | Rẽ nhánh | `if` **Nếu** – điều kiện ảnh / chữ / màu / so giá trị / và-hoặc / đảo | `true` / `false` | từng điều kiện lá đúng / sai + chi tiết |
    | | `switch` **Switch theo giá trị** – `{{biến}}` so với từng case (bằng / chứa / khớp mẫu) | `case:<key>` / `default` | `{value, case}` |
    | | `switchScreen` **Switch theo màn hình** – chờ tối đa X s, cái nào trong N ảnh / chữ hiện trước | `case:<key>` / `none` | `{case, x, y, w, h, score}` |
    | Luồng | `counter` **Đếm** – đếm số lần đi qua (đặt lại mỗi vòng kịch bản) | `next` (còn < max) / `done` | `{count}` |
    | | `end` **Kết thúc** – thành công / thất bại + lời nhắn | – | – |
    | | `setVar` **Đặt biến** – gán / cộng / lấy từ output node khác | `success` | `{name, value}` |
    | Sau này | `subflow` – gọi một kịch bản khác như một node (thay cho Nhóm) | `success` / `fail` | output của node `end` bên trong |

  - [x] **F0. Gỡ bước 20 cũ** *(2026-10-17)*: revert 3 commit 20a / 20b / 20c (`60dcfe2`, `956ec76`, `d5882a7`), giữ 2 commit sửa lỗi sau đó (overlay theo màn chứa vùng, `stepLabel`). Code điều kiện (`evalCond`, `compareValues`, `condtest.go` / `TestCond`, `CondEditor`) **lấy lại từ lịch sử git** ở F3, không viết lại
  - [ ] **F4a. Thiết kế canvas + màn mẫu `#demo-flow`** (dữ liệu giả, chưa nối Go)
    - [ ] `design/DESIGN.md` thêm phần *Canvas flow*: nền canvas (lưới chấm), thẻ node (cỡ, ô icon, tên, mã `n2`, chip tóm tắt), cổng vào / ra + nhãn cổng (`success` màu `ok`, `fail` màu `danger`, case trung tính), cạnh (cong, mũi tên, nhãn, nút × khi rê), node bắt đầu / tắt / đang chọn / cảnh báo / đang chạy / xong / lỗi, minimap + nút zoom, cả 2 theme; token mới (nếu có) thêm cả `theme.css` + `style.css`
    - [ ] Cài `@xyflow/react` + `elkjs`; `components/flow/NodeCard.tsx`, `FlowEdge.tsx`, `FlowCanvas.tsx` nhận dữ liệu qua props (dùng lại được ở F4b)
    - [ ] `demo/flowDemo.ts`: ví dụ Facebook ở trên + 1 node If, 1 Switch theo màn hình, 1 Đếm, 1 Kết thúc, node tắt, node cảnh báo; công tắc xem trạng thái "đang chạy" (node chạy / xong / lỗi, cạnh đã đi) và bảng Nhật ký mẫu
    - [ ] Màn `#demo-flow`: bố cục 3 cột (thư viện | canvas | inspector có mục Cổng ra), Ctrl+K khi kéo cổng ra chỗ trống (giao diện)
    - [ ] Chụp headless cả 2 theme; `npm run build` sạch
    - [ ] Người dùng duyệt `#demo-flow`, góp ý trước khi làm F1
  - [ ] **F1. Model flow + kết quả** (Go + TS, chưa UI – sửa bằng JSON / test)
    - [ ] `model.go`: `Script.start`, `Script.edges []Edge {id, from, port, to}`; `Step.ref`, `Step.name`, `Step.pos {x, y}`, `Step.imageFrom`; `onFail` thành `retry {count, delayMs}`
    - [ ] `flow.go`: `StepResult {nodeId, ref, status: success | fail | skipped | cancelled, port, output, error {code, message}, logs [{t, level, msg}], attempts, startedAt, durationMs}`; mã lỗi `not_found | timeout | permission | invalid_input | cancelled | internal`; `portsOf`; `flowOrder` (duyệt từ `start`, ưu tiên `success` / `true` / case đầu)
    - [ ] `storage.go`: `ScriptVersion` 4, file < 4 báo không mở được; `NewScript` tạo sẵn `start`
    - [ ] `flowcheck.go` `validateFlow`: có `start`; cạnh trỏ node có thật, cổng đúng loại node, mỗi cổng tối đa 1 cạnh; `ref` duy nhất; `{{nX…}}` trỏ node có thật; `imageFrom` là node chờ ảnh / chữ; node không tới được → cảnh báo (không chặn chạy). Lỗi trả `nodeId` + `ref`
    - [ ] `types.ts` khớp Go (tsc kiểm tra), `flow.ts` (`portsOf`, `flowOrder`); `steptypes.go` ghi cổng ra + output mỗi loại, ví dụ qua được `validateFlow`
    - [ ] Test: `validateFlow` từng lỗi, `flowOrder` có vòng lặp (Go + TS cùng kết quả), file cũ bị từ chối
  - [ ] **F2. Engine chạy theo đồ thị** (`flowrun.go`, `vars.go`)
    - [ ] `runFlow`: từ `start` → `execNode` trả `StepResult` → cạnh theo `port` → node kế; hết cạnh kết thúc theo quy tắc. `runStep` bọc lại để trả `StepResult`; thử lại theo `retry`, mỗi lần ghi log; thay `runWithRetry` / `goto` cũ
    - [ ] Ngữ cảnh chạy: `results[nodeId]`, `vars`, `last`; thay `{{nX.output…}}` / `{{vars…}}` lúc chạy từng node (input `{{key}}` vẫn qua `resolveInputs`); click theo ảnh lấy chỗ thấy từ output node `imageFrom`
    - [ ] Giới hạn 10 000 lần chạy node mỗi vòng → lỗi "có thể lặp vô hạn ở n2 → n5 → n4"
    - [ ] Event: `engine:progress` thêm `nodeId`, `ref`, `port`; event mới `engine:node` gửi `StepResult`
    - [ ] Test bằng input giả: `fail` có / không nối, thử lại rồi qua `fail`, vòng lặp, node tắt, biến từ output, dừng / tạm dừng, giới hạn 10 000
  - [ ] **F4b. Nối canvas vào app** (bỏ danh sách bước cũ)
    - [ ] `MainScreen`: cột giữa thành `FlowCanvas`; bỏ `StepList.tsx` + `@dnd-kit/*`; thẻ / màu theo F4a
    - [ ] `flowEdit.ts` (thuần, có test): nối / cắt cạnh, xoá node kèm cạnh, chèn node giữa cạnh (A → B thành A → mới → B), copy / paste node + cạnh bên trong, đặt bắt đầu; `flowLayout.ts` (`elkjs`, nút *Sắp xếp lại*)
    - [ ] Thao tác: kéo node; kéo cổng sang node = nối; **kéo cổng ra chỗ trống → Ctrl+K** thêm node đã nối; Ctrl+K khi đang chọn node = chèn sau cổng `success`; Delete xoá node / cạnh; Shift + kéo chọn nhiều (nối vào chọn nhiều điểm của bước 30); bật / tắt; vừa màn hình, zoom, minimap
    - [ ] Inspector: `StepForm` + *Khi thất bại* (thử lại N lần, nghỉ X ms) + **Cổng ra** (mỗi cổng một ô chọn node đích – thay cho kéo, dùng được bằng bàn phím); sửa `ref`; `VarTextInput` gợi ý cả output các node
    - [ ] `scriptStore.ts`: node + cạnh + vị trí qua `flowEdit`; hoàn tác / làm lại tính cả cạnh, vị trí (kéo node = 1 lần); tự lưu nháp như cũ
    - [ ] Overlay: marker đánh số theo `flowOrder`; chọn node trên canvas → xem riêng điểm của node
    - [ ] Badge cảnh báo trên node lấy từ `validateFlow`
    - [ ] Kiểm tra tay: soạn ví dụ Facebook bằng canvas (chỉ node hành động + cạnh `fail` quay lại), chạy thật
  - [ ] **F3. Node điều khiển** (`nodes.go`, `cond.go`, `components/flow/nodes/*`)
    - [ ] Lấy lại `Cond`, `evalCond`, `compareValues` (từ `60dcfe2`) và `CondEditor` (từ `d5882a7`) cho `if`
    - [ ] `switch`, `switchScreen` (chụp một lần rồi so mọi case), `counter`, `end`, `setVar` – mỗi loại: engine + `checkStep` + `steptypes.go` + inspector + thẻ node + test
    - [ ] **Thử node** (lấy lại `condtest.go` / `TestCond`, mở rộng): chạy riêng `if` / `switch` / `switchScreen` trên màn hình lúc này → cổng sẽ đi + chi tiết từng phần, khung chỗ thấy trên overlay
    - [ ] Ctrl+K nhóm *Rẽ nhánh & luồng*
  - [ ] **F5. Lúc chạy + nhật ký + debug** *(thay bước 25c cũ)* (`runlog.go`, `RunLog.tsx`, `runStore.ts`)
    - [ ] Canvas khi chạy: node đang chạy viền `accent` + vòng tiến độ, xong tô `ok` / `danger`, cạnh vừa đi sáng lên, số lần đi qua trên node
    - [ ] Bảng **Nhật ký** dưới canvas: mỗi dòng một node đã chạy (giờ, `n2`, tên, status, ms, cổng đi) → bấm mở output / lỗi / log / ảnh lúc lỗi; bấm dòng → chọn node
    - [ ] Lưu mỗi lần chạy `UserConfigDir/autoclicker/runs/<kịch bản>/<thời gian>.json` (`StepResult[]`, input đã dùng, nguồn: người dùng / hẹn giờ / CLI / AI, kết quả cuối) + PNG khi node lỗi / hết giờ; giữ 50 lần gần nhất mỗi kịch bản; mục *Lịch sử* trong thư viện mở lại (tô đường đã đi)
    - [ ] Debug: chạy từ node đang chọn, chạy từng node (Tiếp / F10), breakpoint; `RunNode` là hàm Go riêng (Thử node và MCP dùng chung)
  - [ ] **F6. Xem trước, thanh mini**
    - [ ] Xem trước (`Simulation.tsx`) đi theo đồ thị: mặc định `success` / `true` / case đầu; node rẽ nhánh có công tắc "giả sử thất bại" / chọn case; Đếm chạy 1 vòng rồi `done`; đường đang mô phỏng tô trên canvas
    - [ ] Thanh mini: "n2 · Tìm nút Nhắn tin · thất bại → n5 Đếm (3 / 10)"
  - [ ] **F7. Chuẩn bị cho MCP (31):** JSON flow là định dạng duy nhất; `get_step_types` trả cổng ra + output từng loại; `validateFlow` lỗi có `nodeId` / `ref`; kết quả chạy trả `StepResult[]` để AI đọc lỗi và tự sửa kịch bản
  - Sau này (28): đọc chữ trong vùng / clipboard vào biến (`setVar` từ OCR), lặp theo từng dòng CSV, `subflow`
  - Xong khi: ví dụ Facebook ở trên (thấy "Nhắn tin" → click; không thấy → cuộn 3 nấc, lặp tối đa 10 lần rồi kết thúc báo lỗi) soạn hết bằng canvas, chạy thật được, xem được nhật ký từng node; `go test ./...`, `npm test`, `npm run build` sạch
- [ ] **26. Focus cửa sổ theo tên**
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `focusWindow` (`title` hoặc process); Engine: `robotgo.ActivateName`, không thấy → `onTimeout`
  - [ ] UI: icon `focusWindow`; ô nhập (nhận biến) + chọn từ danh sách cửa sổ đang mở
  - [ ] Chuẩn bị cho MCP (31): danh sách cửa sổ là hàm Go (`ListWindows`) dùng chung cho UI và MCP
- [ ] **25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật** *(người dùng yêu cầu, 2026-10-17; tách từ 28)* – Tìm phần tử trên trang theo chữ / selector (vd nút "Xem lịch sử"), tính toạ độ màn hình của nó rồi click / cuộn bằng robotgo như bước thường. Bền hơn ảnh mẫu (không phụ thuộc theme, zoom, hover), nhanh hơn (vài ms), đọc được trạng thái phần tử; AI (31) đọc DOM chính xác hơn nhìn ảnh chụp.
  - Cách làm: **DOM để tìm, robotgo để bấm** → event thật của hệ điều hành (`isTrusted = true`, khó bị coi là bot), dùng lại engine / xem trước / thanh mini. Dự phòng *chạy nền* (cửa sổ thu nhỏ / bị che): `Input.dispatchMouseEvent` của CDP – người dùng chọn "chuột thật / chạy nền"
  - [ ] Kết nối Chrome / Edge / Brave qua **CDP** (`chromedp`): app mở trình duyệt với `--remote-debugging-port` + **profile riêng** trong `UserConfigDir/autoclicker/browser-profile` (từ Chrome 136 không bật debug được trên profile mặc định – lần đầu người dùng đăng nhập lại, sau đó profile tự nhớ); nối vào phiên đang mở nếu có. Safari / Firefox không hỗ trợ. Phương án sau (nếu cần profile đang dùng): extension Chrome nói chuyện với app qua WebSocket
  - [ ] Selector, ưu tiên theo độ bền: chữ hiển thị (`text=Xem lịch sử`) / `aria-label` / `role` → `id` / `data-*` → CSS / XPath; **tránh class sinh tự động** (React / Facebook `x1n2onr6…`); cho nhiều selector dự phòng (trang đổi ngôn ngữ); chọn tab theo URL / tiêu đề; phần tử trong iframe cộng toạ độ iframe
  - [ ] **Đổi toạ độ trang → màn hình bằng cách đo** (giống overlay, xem *Quyết định đã chốt*), không tính từ `screenX` + chiều cao thanh tab (sai khi mở DevTools, thanh bookmark, thanh "Chrome đang bị điều khiển", DPI 125%…): robotgo đưa chuột tới 2 điểm trong trang, nghe `mousemove` qua CDP lấy `clientX/Y` → suy ra tỷ lệ (zoom × DPI) + độ lệch; `screen = offset + client × scale`. Đo lại khi cửa sổ dời / đổi cỡ (`screenX/Y`, `outerWidth/Height` đổi) hoặc đổi zoom
  - [ ] Trước khi bấm: phần tử ngoài viewport → cuộn (robotgo 1–2 nấc giữa trang rồi **đọc lại rect** đến khi lọt viewport; không cần cuộn thật thì `scrollIntoView({block:"center"})`); chờ ổn định (rect giống nhau 2 lần cách ~100 ms); **không bị che**: `document.elementFromPoint(tâm)` phải là phần tử hoặc con của nó, không thì chờ / lỗi; cửa sổ trình duyệt phải ở trên cùng (focus như 26)
  - [ ] Model: `webWait` (chờ phần tử hiện / mất, timeout, ghim được), `webClick` (nút chuột, số click, vị trí trong phần tử: tâm / ±dx dy), `webScroll` (cuộn tới phần tử / theo nấc trong phần tử), `webType` (focus ô rồi gõ bằng bước gõ chữ sẵn có), `webRead` (đọc chữ / thuộc tính vào biến – nối với biến / điều kiện 28); hết giờ dùng `onTimeout` của 25b. Thêm vào `model.go`, `types.ts`, `steptypes.go` (mô tả + ví dụ JSON)
  - [ ] Engine: interface `Browser` (như `Input` / `Screen`), test bằng trình duyệt giả (DOM giả + toạ độ); test thật `AUTOCLICKER_REAL=1` mở trang HTML cục bộ
  - [ ] UI: inspector (selector + dự phòng, tab, chuột thật / chạy nền), **Thử tìm** (khung phần tử trên overlay, toạ độ, chữ đọc được, ms); nút **Chọn phần tử trên trang** (inject script tô viền phần tử dưới chuột, click → tự sinh selector theo thứ tự ưu tiên ở trên); trạng thái kết nối trình duyệt
  - [ ] Chuẩn bị cho MCP (31): `FindElement(selector) → {rect màn hình, chữ, hiện / bị che}` và `PageOutline()` (cây accessibility rút gọn cho AI) là hàm Go độc lập UI – Thử tìm, Xem trên màn hình, MCP dùng chung
  - [ ] Kiểm tra tay: trang có nút "Xem lịch sử" nằm dưới màn hình → bước chờ + click chạy đúng khi dời cửa sổ, zoom 125%, mở DevTools, có popup che
- [ ] **22. Hẹn giờ + thông báo + CLI / API local** – Chạy không cần ngồi canh.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Hẹn giờ: chạy lúc HH:mm (một lần / hằng ngày) với preset input đã chọn; chip giờ trong thư viện; app thu xuống khay
  - [ ] Thông báo: Telegram / Discord / webhook khi xong hoặc lỗi, kèm ảnh chụp màn hình
  - [ ] CLI: `autoclicker run script.json --preset "Gửi cho Lan" --set M1=Lan`
  - [ ] API local: `POST /run` (kèm values), `POST /stop`, `GET /status`, tiến độ qua SSE; bảo vệ bằng token
  - [ ] Chuẩn bị cho MCP (31): server HTTP chỉ nghe localhost, gắn thêm route được (`/mcp` sau này); kịch bản bị đổi từ ngoài (file / API) → event `script:changed` để frontend nạp lại, giữ hoàn tác; chỗ chặn "hỏi người dùng trước khi chạy" dùng chung cho hẹn giờ / API / AI
- [ ] **31a. Kết nối AI qua MCP – tạo / sửa kịch bản** *(người dùng yêu cầu, 2026-10-17)* – Người dùng ra lệnh trên Claude / GPT / Cursor… và AI tạo ra một hoặc nhiều kịch bản trong app, không cần soạn tay. App làm **MCP server** (Model Context Protocol – chuẩn chung, một server dùng được cho mọi AI hỗ trợ MCP).
  - Kiến trúc:
    - App đang chạy mở MCP server **Streamable HTTP** tại `127.0.0.1:<port>/mcp` (chỉ nghe localhost, token trong `UserConfigDir/autoclicker/`); dùng SDK chính thức `github.com/modelcontextprotocol/go-sdk`. dùng chung server HTTP + token của API local (22)
    - Cùng file chạy có chế độ `autoclicker --mcp` (**stdio**) cho AI chỉ chạy được server dạng lệnh (Claude Desktop, Codex CLI, Gemini CLI…): chuyển tiếp sang app đang chạy; app chưa mở thì tự mở, hoặc chỉ làm việc với file kịch bản
    - Kết nối được: Claude Desktop / Claude Code / Cursor / GitHub Copilot (chế độ agent trong VS Code) / Codex CLI / Gemini CLI (stdio hoặc HTTP localhost). Copilot Studio (Microsoft 365) cũng chỉ nhận MCP từ xa như ChatGPT. ChatGPT bản web / desktop chỉ nhận MCP từ xa qua HTTPS → cần tunnel (cloudflared / ngrok) + token; ghi rõ trong hướng dẫn, không bật mặc định
  - [ ] Schema cho AI: sinh **JSON Schema** của `Script` / `Step` từ `model.go` (mô tả tiếng Việt + tiếng Anh từng trường, giá trị mặc định, loại bước, neo, `{{biến}}`); test kiểm tra schema khớp struct Go (như `types.ts`). Resource `autoclicker://schema`, `autoclicker://guide` (cách viết kịch bản bền: ưu tiên mở URL / phím tắt / chờ ảnh hơn toạ độ cứng, thêm chờ sau khi mở app…), `autoclicker://examples` (vài kịch bản mẫu)
  - [ ] Tool đọc: `list_scripts` (thư viện + mở gần đây), `read_script`, `get_step_types` (loại bước app đang hỗ trợ – AI không dùng loại *Sắp có*), `get_screen_info` (kích thước, tỷ lệ Retina / DPI, số màn hình)
  - [ ] Tool ghi: `validate_script` (trả lỗi rõ theo từng bước để AI tự sửa), `create_script` (một hoặc nhiều kịch bản, lưu vào thư viện), `update_script` (sửa theo `stepId`, không ghi đè cả file), `open_in_app` (mở kịch bản trên MainScreen cho người dùng xem)
  - [ ] Prompt MCP `tao_kich_ban`: khung hỏi đáp mẫu (mục tiêu, app / trang web, giá trị nào là input) để AI hỏi lại khi thiếu thông tin
  - [ ] An toàn: kịch bản AI tạo có `meta.createdBy: "ai"` + chip "AI" trong thư viện, mở ở trạng thái chưa lưu để người dùng duyệt; AI sửa kịch bản đang mở có thay đổi chưa lưu → tạo bản sao, không ghi đè; mọi thay đổi của AI đi qua store nên Ctrl+Z hoàn tác được; bước này **chưa cho AI chạy** kịch bản (31c)
  - [ ] Cài đặt [màn Cài đặt]: bật / tắt MCP, cổng, token (sao chép / tạo lại), nút *Thêm vào Claude Desktop* (ghi `claude_desktop_config.json`, hỏi trước), đoạn cấu hình mẫu cho Claude Code / Cursor / Codex; nhật ký các lệnh AI đã gọi
  - [ ] Test Go: gọi tool qua client MCP trong test (in-memory transport): tạo kịch bản hợp lệ / lỗi, sửa theo `stepId`, schema khớp struct
  - Xong khi: trên Claude Desktop nói "tạo kịch bản mở Spotify, tìm bài Tinh Hà Say 2 và phát" → kịch bản hiện trong thư viện app (chip AI), Xem trước chạy đúng thứ tự, bấm Chạy thì phát nhạc; nói "đổi thành bài khác, cho tên bài thành input" → AI sửa đúng kịch bản, có input `{{bai_hat}}`
- [ ] **31b. MCP – AI nhìn màn hình** – Để AI tự biết vị trí nút thay vì đoán toạ độ.
  - [ ] Tool `screenshot` (toàn màn hình / vùng / cửa sổ, thu nhỏ cho nhẹ, trả ảnh + tỷ lệ); toạ độ AI gửi về theo pixel ảnh, Go tự đổi sang toạ độ robotgo (Retina / DPI – dùng chung cách đổi của 17, 25)
  - [ ] Tool `find_image` (dùng so khớp mẫu của 25), `capture_template` (AI khoanh vùng trên ảnh chụp → lưu ảnh mẫu cạnh kịch bản cho bước `waitImage` / `clickImage`), `get_mouse_position`, `list_windows`
  - [ ] Người dùng thấy khi AI chụp màn hình: thông báo nhỏ / chấm đỏ trên thanh mini; tắt được trong Cài đặt
  - Xong khi: AI chụp màn hình YouTube, tự khoanh nút "Phát" làm ảnh mẫu, kịch bản tạo ra dùng `clickImage` thay vì toạ độ cứng và vẫn chạy khi cửa sổ dời chỗ
- [ ] **31c. MCP – AI tự thử và chạy kịch bản** – AI làm từng bước trên màn hình thật, xem kết quả, bước nào đúng thì ghi vào kịch bản (giống người dùng soạn rồi bấm thử).
  - [ ] Tool `run_step` (chạy một bước, trả ảnh chụp sau bước), `run_script` (kèm giá trị input, dùng `resolveInputs`), `get_run_status` (tiến độ + nhật ký 25c), `stop_script`
  - [ ] Quyền: mặc định mỗi lần AI muốn chạy thì app hiện hộp hỏi (tên kịch bản, các bước, Cho phép / Từ chối, "cho phép kịch bản này đến khi đóng app"); F8 và dừng khẩn cấp (góc màn hình) vẫn luôn có tác dụng; giới hạn số bước / thời gian mỗi lượt AI chạy
  - [ ] Chế độ "AI soạn bằng cách thử": các bước AI `run_step` thành công được ghi lại thành kịch bản nháp, kèm ảnh chụp từng bước để người dùng xem lại
  - [ ] Lịch sử chạy (20 F5) ghi rõ lần nào do AI chạy
  - Xong khi: nói "mở bài Tinh Hà Say 2 trên YouTube" → AI hỏi quyền, tự thử (mở URL, chờ, click kết quả đầu), nhạc phát, kịch bản được lưu để lần sau chạy lại không cần AI
- [ ] **24. Ghi macro** – Ghi click / kéo / gõ phím bằng gohook, gộp thành bước, tự tính delay; nút Ghi trên toolbar, F9 dừng.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
- [ ] **28. Mở rộng** – Làm sau cùng.
  - Chrome qua DOM: đã chuyển lên **25d** (2026-10-17)
  - [ ] CSV: mỗi cột là một input key, chạy kịch bản cho từng dòng (mỗi dòng một lần `resolveInputs`)
  - [ ] OCR: đọc chữ trong vùng vào biến, rẽ nhánh theo nội dung (Tesseract)
  - [ ] Tìm phần tử bằng AI: mô tả bằng lời ("nút Gửi"), model có vision trả toạ độ; dự phòng khi ảnh mẫu không khớp
  - [ ] Biến: đọc clipboard / OCR vào biến, đặt / cộng biến – điều kiện `if / else`, "lặp đến khi thấy" đã chuyển lên **20**
  - [ ] Thư viện mẫu, chia sẻ gói kịch bản (kèm ảnh mẫu), thống kê lần chạy (SQLite)


### Ví dụ mục tiêu: nhờ AI mở một bài nhạc

Người dùng gõ trên Claude: *"tạo kịch bản mở bài Tinh Hà Say 2 trên YouTube"*. Claude gọi `get_step_types` / đọc `autoclicker://guide`, tạo kịch bản, `validate_script`, rồi `create_script` + `open_in_app`:

| # | Bước | Loại | Input | Cần bước |
|---|---|---|---|---|
| 1 | Mở `https://www.youtube.com/results?search_query={{bai_hat}}` | `open` | `bai_hat` = "Tinh Hà Say 2" | 21b, 23, 31a |
| 2 | Chờ trang tải (chờ ảnh khung kết quả, tối đa 10 s) | `waitImage` | | 25, 31b |
| 3 | Click video đầu tiên | `clickImage` (ảnh mẫu AI tự khoanh) | | 25, 31b |

> Bản đơn giản hơn chỉ cần 31a (không cần 31b): mở URL → chờ 4 s → click toạ độ kết quả đầu (AI hỏi người dùng vị trí hoặc chụp màn hình ở 31b). Kém bền hơn khi cửa sổ đổi cỡ.

### Ví dụ mục tiêu: nhắn tin cho một người trên Facebook

| # | Bước | Loại | Input | Cần bước |
|---|---|---|---|---|
| 1 | Mở `facebook.com/{{profile}}` | `open` | `profile` | 21b, 23 |
| 2 | Đưa Chrome lên trước | `focusWindow` | | 26 |
| 3 | Chờ nút "Nhắn tin" (thử lại 2 lần) | `waitImage` | `wait_s` | 25, 25b |
| 4 | Click nút | `clickImage` | `click_delay` | 25 |
| 5 | Gõ "xin chào {{M1}} và {{M2}}" | `type` | `M1`, `M2` | có sẵn + 21b |
| 6 | Enter | `key` | | có sẵn |
| 7 | Cuộn xuống 5 nấc | `scroll` | | 23 |

> Tự động hoá trên Facebook vi phạm điều khoản của họ và dễ bị khoá tài khoản, nhất là khi nhắn hàng loạt. Thử trên tài khoản phụ hoặc trang của chính bạn.

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
- **2026-10-15** – Commit bước 8 (`a3d3058`). Bước 9: `failsafe.go` (dừng khi chuột chạm góc, bỏ qua khi do engine đưa tới), event `engine:stopped` + lý do. `go test -race` 8/8, `npm run build` OK; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 9: kiểm tra tay macOS – rê chuột vào góc khi đang chạy → dừng, về main OK.
- **2026-10-15** – Commit bước 9 (`c6b8b96`). Bước 10: `scriptStore.ts`, overlay chọn vùng + tự đo lệch toạ độ + marker theo dữ liệu thật. `npm run build` OK, chụp headless kiểm tra overlay / main; chưa kiểm tra tay trong app.
- **2026-10-15** – Xong bước 10: kiểm tra tay macOS – chọn vùng OK, kể cả trên màn hình thứ 2 (không cần bộ chọn màn hình).
- **2026-10-15** – Commit bước 10 (`e61c3b6`). Bước 11: đặt điểm trên overlay, F6 toàn cục, chọn lại điểm cho bước. `npm run build`, `go test` OK; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 11: kiểm tra tay macOS OK. Người dùng hỏi cách xoá điểm → chưa có, thuộc bước 14.
- **2026-10-15** – Bước 12: kéo marker để dời điểm, click marker chọn bước. `npm run build` OK, chụp headless overlay; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 12: kiểm tra tay macOS OK.
- **2026-10-15** – Bước 13: chọn A → B cho kéo thả trên overlay, chọn lại A / B bằng ⌖. `npm run build` OK; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 13: kiểm tra tay macOS OK.
- **2026-10-15** – Bước 14: `StepList.tsx` (dnd-kit sắp xếp, bật/tắt, nhân bản, xoá + hoàn tác, menu thêm bước). `npm run build` OK, chụp headless màn chính; chưa kiểm tra tay.
- **2026-10-15** – Xong bước 14: kiểm tra tay macOS OK.
- **2026-10-15** – Bước 15: `StepForm.tsx` form sửa bước theo loại, `updateStep` / `changeStepType`, `KEY_NAMES` + `isValidKey`. `npm run build` OK, chụp headless form click / kéo thả / phím tắt; chưa kiểm tra tay, chưa commit.
- **2026-10-15** – Bước 15 (tiếp): phím tắt nhiều phím chính + chọn cách bấm (chord / sequence) theo yêu cầu người dùng. `go test -race` 10/10, `npm run build` OK; chưa kiểm tra tay, chưa commit.
- **2026-10-15** – Xong bước 15: người dùng test form OK, commit.
- **2026-10-15** – Commit bước 15 (`2a9e42f`). Bước 16: nút Chạy / F8 chạy kịch bản đang soạn, đếm ngược trong Go, `RunSettings.tsx`, toạ độ chuột ở thanh dưới, dọn code thử Phase 0. `go test`, `npm run build` OK, chụp headless; chưa kiểm tra tay, chưa commit.
- **2026-10-15** – Sửa lỗi overlay mở nhầm màn hình khi có 2 màn (`window.go`: đặt vị trí trước khi phóng to). `go build` OK; chờ người dùng kiểm tra.
- **2026-10-15** – Xong bước 16 (kể cả sửa overlay mở nhầm màn): người dùng kiểm tra tay OK, commit.
- **2026-10-15** – Commit bước 16 (`aadc204`). Bước 17: dừng khẩn cấp mọi màn hình (góc cứng), overlay đo tỷ lệ + độ lệch (`coordMap.ts`). `go test -race` OK (thêm `TestHardCorners`), `npm run build` OK; chưa kiểm tra tay, chưa commit; phần Windows chờ máy.
- **2026-10-15** – Commit phần bước 17 làm được trên macOS (chưa kiểm tra tay); hoãn phần Windows của 17 và bước 18 vì chưa có máy Windows. Chuyển sang bước 19.
- **2026-10-15** – Bước 19: ngẫu nhiên hoá thời gian / vị trí (model, engine, form, overlay). `go test -race` OK, `npm run build` OK, chụp headless form / overlay; chưa kiểm tra tay, chưa commit.
- **2026-10-15** – Kéo thả: thêm “từ điểm trước” và “điểm cuối lệch so với A” theo yêu cầu (model, engine, store `endpoints`, form, overlay). `go test -race` OK, `npm run build` OK, chụp headless; chưa kiểm tra tay, chưa commit (chung với bước 19).
- **2026-10-15** – Thêm *Xem trước kịch bản* (19b) theo yêu cầu: con trỏ giả chạy 1 vòng trên overlay. `npm run build` OK, chụp headless (thời gian ảo) thấy chạy hết vòng; chưa kiểm tra tay, chưa commit.
- **2026-10-15** – Commit theo yêu cầu: bước 19 (ngẫu nhiên hoá), kéo thả từ điểm trước / lệch so với A, 19b (xem trước) – chưa xác nhận kiểm tra tay. Bắt đầu: xem / sửa vị trí riêng một bước trên overlay.
- **2026-10-15** – Commit `f886bd8`. Thêm 19c: xem / sửa riêng một bước trên overlay. `npm run build` OK, chụp headless; chưa kiểm tra tay, chưa commit.
- **2026-10-15** – Đánh dấu xong 19, kéo thả mới, 19b, 19c theo xác nhận của người dùng (yêu cầu commit sau khi thử). Commit + push lên GitHub.
- **2026-10-15** – Sửa giao diện màn chính theo phản hồi (header 2 hàng, thanh dưới bị cắt ở cửa sổ 1024×768): header một hàng gọn (nút h-9, chữ 13px, tên app chỉ hiện chữ từ `xl`), màn chính `h-screen`, danh sách bước và form tự cuộn trong khung, hàng nút loại bước nhỏ lại. Chụp headless 1024×768 và 1440×800 OK; chưa commit.
- **2026-10-15** – Sửa dòng trong danh sách bước bị nhảy khi rê chuột (mô tả dài): nút nhân bản / xoá nổi tuyệt đối ở góc phải, thời gian chờ chỉ `invisible` để giữ chiều rộng. `npm run build` OK; chưa commit.
- **2026-10-16** – Nhận design v2 + kế hoạch mới từ người dùng: chuyển vào `design/` (DESIGN.md, theme.css, canvas) và `docs/ROADMAP.md`. Sắp lại Phase 3 theo ROADMAP: thêm 27 (UI v2), 21b (input), 25b (hết giờ), 25c (debug + lịch sử), 28 (mở rộng); mở rộng 23 (+ mở URL / app), 22 (+ thông báo, CLI / API); bỏ Phase 2 / 3 cũ chưa làm (gộp vào danh sách mới). Chưa đổi code.
- **2026-10-16** – Gộp toàn bộ `docs/ROADMAP.md` + phần design (artboard, thành phần, icon) vào checklist từng bước Phase 3 theo yêu cầu; xoá `docs/ROADMAP.md` (TASKS là nguồn duy nhất), chi tiết thiết kế vẫn ở `design/DESIGN.md`. Chưa đổi code.
- **2026-10-16** – Thêm theo yêu cầu: bước 29 (neo điểm theo 4 góc vùng, ±X / ±Y) và 30 (chọn nhiều điểm: di chuyển, co giãn, xoay, xiên, perspective, căn chỉnh, ánh xạ sang vùng mới, hoàn tác). Xếp sau 27, trước 21. Chưa đổi code.
- **2026-10-16** – Sắp lại thứ tự Phase 3 theo phụ thuộc: 29 → 27 → 21 → 30 → 21b → 23 → 25 → 25b → 25c → 26 → 20 → 22 → 24 → 28; thêm bảng thứ tự + lý do ở đầu Phase 3. Chưa đổi code.
- **2026-10-16** – Theo yêu cầu: đưa UI lên đầu – 27 (đổi giao diện màn đang có) + 27b mới (dựng khung UI cho mọi tính năng chưa có, nhãn *Sắp có*, màn `#demo`); thứ tự mới 27 → 27b → 29 → 21 → 30 → 21b → 23 → 25 → 25b → 25c → 26 → 20 → 22 → 24 → 28. Chưa đổi code.
- **2026-10-16** – Commit kế hoạch (`9ad3e7b`). Bước 27 đợt 1: `style.css` theo `design/theme.css` (body không tô nền để overlay / mini bar trong suốt), font Plus Jakarta Sans + JetBrains Mono (có tiếng Việt), đổi token cũ → mới ở mọi màn, overlay + PreviewScreen `data-theme="dark"` (HUD luôn tối), marker / mũi tên kéo thả dùng `accent`, segmented kiểu `bg-soft` / `bg-raised`, `StepIcon` đủ bộ icon, `theme.ts` (tạm đọc `?theme=`). `npm run build` OK, chụp headless main / overlay / minibar cả 2 theme.
- **2026-10-16** – Bước 27 đợt 2: `TitleBar.tsx` (kiểu Windows 11, double click phóng to), thanh công cụ (tên kịch bản, tab Soạn / Xem trước, nút "N vòng" mở popup `RunSettings`, Chạy `btn-primary`), `Library.tsx` (khung), `StepList.tsx` thành thẻ bước (ô icon = tay kéo, nút bật/tắt / nhân bản / xoá khi rê, chip vùng ở đầu danh sách), `StepForm.tsx` thành inspector 340 px (tab Bước / Input – Input tắt, `Segmented`, `Switch`, `NumInput` cao 32 + `ring-focus`). Go: nhớ trạng thái phóng to khi đổi chế độ cửa sổ (`mainMax`). Chụp headless OK.
- **2026-10-16** – Bước 27 đợt 3: `AddStepPalette.tsx` (Ctrl/Cmd+K, tìm không dấu, 4 nhóm, ←→↑↓ + Enter, loại chưa làm nhãn *Sắp có*); store `addPointStep` / `insertStep` – bước thêm từ bảng chèn sau bước đang chọn, overlay đặt điểm tạo đúng loại (click / double / giữ). Kiểm tra bằng Playwright: tìm "cuon" → Cuộn, Enter thêm Chờ sau bước 3, Esc đóng.
- **2026-10-16** – Bước 27 đợt 4: mini bar 540×96 (vòng tiến độ gradient + icon bước đang chạy, `btn-danger-soft`), overlay theo artboard (thanh công cụ có icon ở đáy, banner hướng dẫn giữa trên, 4 núm góc vùng), `settings.go` (`GetSettings`, `SetTheme`, file `settings.json` trong `UserConfigDir`, test `TestSettingsRoundTrip`, `TestSettingsBadFile`), `theme.ts` + nút Sáng / Tối / Theo hệ thống trên thanh tiêu đề, cập nhật `CLAUDE.md`. `go test -race` OK, `npm run build` OK, chụp headless + Playwright kiểm tra đổi theme; chưa kiểm tra tay, chưa commit.
- **2026-10-16** – Thêm thu gọn / mở rộng cột thư viện theo yêu cầu (`Library.tsx`, icon `SidebarToggle`, Ctrl/Cmd+B, nhớ bằng localStorage). `npm run build` OK, Playwright: 236 → 44 px, tải lại vẫn thu gọn, Cmd+B mở lại.
- **2026-10-16** – Xong bước 27: người dùng kiểm tra tay trong app OK, commit. Bắt đầu 27b.
- **2026-10-16** – Bước 27b đợt 1: tách `StepCard.tsx` (thẻ, khung nhóm, nhãn *Sắp có*), `LibraryView` (tìm, danh sách + giờ hẹn + chấm chưa lưu, mở gần đây, lịch sử chạy, Mới / Mở / Lưu), `AppChrome.tsx` (`Toolbar` có chip "Chưa lưu", `StatusBar`), màn `#demo` với dữ liệu mẫu (`demo/demoScript.ts`). App thật: các mục chưa có logic tắt + *Sắp có*.
- **2026-10-16** – Bước 27b đợt 2: `components/inputs/` – `PinButton` + `KeyChip`, `VarTextInput` (tô `{{biến}}` bằng lớp vẽ dưới textarea, gợi ý khi gõ `{{`, "+ Tạo input mới"), `InputsPanel` (kéo sắp xếp, mở rộng để sửa), `RunDialog` (bộ giá trị, ô đã đổi + khôi phục, bắt buộc, "Sẽ gõ"). Inspector thật có nút ghim tắt. Sửa lỗi con trỏ nhảy sau khi chèn biến (`useLayoutEffect` thay `requestAnimationFrame`). Kiểm tra bằng Playwright.
- **2026-10-16** – Bước 27b đợt 3: `components/steps/` – inspector cuộn, mở URL / app, chờ ảnh, click theo ảnh, chờ pixel, focus cửa sổ, Khi hết giờ, nhóm bước, `AnchorField` (± theo neo, lưới 3×3, %); nút debug trên thanh công cụ; `Field` dùng `<div>` khi bọc nhóm nút (trước đó bấm chữ "Nút chuột" là chọn luôn "Trái").
- **2026-10-16** – Bước 27b đợt 4: overlay thật thêm công cụ Chọn nhiều / Khoanh ảnh (tắt); `#demo-overlay` (khoanh ảnh mẫu; chọn nhiều điểm, khung biến đổi, ô ΔX ΔY % xoay xiên áp thử, căn / chia đều / lật, hộp ánh xạ 4 cách); `MiniBarView` + `#demo-minibar`; `ScheduleDialog` (giờ, lặp theo ngày, bộ giá trị, khay, thông báo); nút Ghi F9. `npm run build` OK, chụp headless cả 2 theme; chưa commit. Người dùng thêm `design/newdarkmode.md` (bảng màu tối Slate / Dimmed / Nord) – chưa áp dụng, chờ hỏi.
- **2026-10-16** – Xong 27b theo xác nhận của người dùng, commit (chưa gồm `design/newdarkmode.md`).
- **2026-10-16** – Theme tối đổi sang bảng Slate (`design/newdarkmode.md`) trong `style.css` + `design/theme.css`: nền / bề mặt / ô nhập tách lớp rõ, accent `#60A5FA`, thêm token `on-accent` (chữ trên nền accent, thay `text-white`), HUD overlay cùng tông. `npm run build` OK, chụp headless màn chính / `#demo` / overlay / mini bar; chưa commit.
- **2026-10-16** – Sửa nút đóng cửa sổ khi rê chuột ở theme sáng bị nền xám + chữ trắng (`hover:bg-hover` đè `hover:bg-danger`) – tách class riêng. Commit theme tối Slate + bản sửa.
- **2026-10-16** – Sửa (người dùng báo): thanh công cụ đáy overlay (Xem trước / Đặt điểm…) bị Dock che – overlay đặt dưới menu bar, cao bằng màn hình nên tràn qua mép dưới. Theo yêu cầu dời thanh công cụ lên giữa phía trên (`OverlayScreen.tsx`), banner hướng dẫn + khung toạ độ xuống ngay dưới. `npm run build` OK, chụp headless `#overlay` 1280 không chồng nhau; chưa kiểm tra tay, chưa commit.
- **2026-10-16** – Bước 29 (neo 4 góc): `anchor.go` (`resolvePoint`, `mainScreen`, `ScreenBounds`) + `anchor.ts` cùng công thức; `Point.ax/ay` (`omitempty`, trống = start), ratio lưu 1/10000; `validateScript` báo điểm ngoài vùng; toàn màn hình neo theo màn chính. Frontend: `StepPoint`, đặt điểm neo góc gần nhất, lưới / % đổi neo giữ chỗ (`setAnchor`), `AnchorField` thật trong inspector (điểm, A, B), danh sách hiện ±, overlay vẽ đường neo (bước đang chọn, giữ Alt = mọi điểm), marker bám vùng ngay khi đang kéo. `go test -race` OK (thêm 4 test neo), `npm run build` OK, kiểm tra `anchor.ts` bằng cùng bộ số với test Go, chụp headless inspector + overlay. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu: ô chọn neo 3×3 quá nhỏ → bản đồ thu nhỏ vùng / màn hình theo đúng tỷ lệ (`AnchorMap` trong `AnchorField.tsx`, tối đa 268×150) dưới ô X / Y: điểm của bước, 9 chấm neo bấm được (vùng bấm 20 px), đường neo + khoảng cách, điểm còn lại của bước kéo thả vẽ mờ, nhãn cỡ vùng. `npm run build` OK, chụp headless sáng / tối; chưa commit.
- **2026-10-17** – Theo yêu cầu: bước kéo thả dùng **một** bản đồ chung cho A và B (marker A / B + mũi tên A → B), nút A | B chọn điểm đổi neo; ẩn A khi "từ điểm trước", ẩn B khi "lệch so với A". `AnchorMap` export, nhận `marks` / `arrow`. `npm run build` OK, chụp headless; chưa commit.
- **2026-10-17** – Người dùng xác nhận kiểm tra tay bước 29 OK → tick. Bước 21 (lưu / tải): `storage.go` (file `.json` có `version`, kiểm tra cấu trúc khi mở, mở gần đây, bản nháp, hỏi lưu khi đóng app qua `OnBeforeClose`), `scriptFile.ts` (cờ chưa lưu, tự ghi nháp, khôi phục khi mở app), `SaveChangesDialog`, thư viện / phím tắt / status bar nối thật; thêm `Script.anchorMode` (mục còn lại của 29). `go test -race` OK (thêm 6 test lưu / đọc), `go build` + `npm run build` OK. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Người dùng xác nhận kiểm tra tay bước 21 OK → tick. Sửa (người dùng báo): gợi ý tên phím (bước Phím tắt) dùng `<datalist>` – WebView macOS vẽ chữ trắng trên nền trắng, không có mô tả → combobox tự vẽ trong `StepForm.tsx` (lọc theo phím đang gõ sau dấu phẩy, ↑ ↓ / Enter / Tab / Esc), mỗi dòng tên phím + mô tả tiếng Việt (`keyLabel` trong `types.ts`). `npm run build` OK, chụp headless sáng / tối thấy rõ chữ.
- **2026-10-17** – Sửa (người dùng báo): Ctrl+K chọn Double click / Giữ chuột nhưng đặt điểm ra Click – StrictMode (`wails dev`) gỡ rồi gắn lại OverlayScreen khi mở, cleanup gọi `cancelPick()` reset loại bước + chỗ chèn + chọn lại điểm. Chỉ huỷ khi hash đã rời `#overlay`. `npm run build` OK; chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Bước 30: `transform.ts` (ma trận 3×3 affine + homography, căn chỉnh, vùng → vùng) + vitest 14 test; store: hoàn tác / làm lại mọi thay đổi, chọn nhiều, `mapPoints` (giữ neo, B lệch, chặn điểm ra ngoài vùng) dùng chung cho biến đổi / xem trước / ánh xạ vùng; overlay: công cụ *Chọn nhiều* (`TransformTool.tsx`: khung 8 núm + xoay, ⌘ xiên / phối cảnh, thanh số né khung + kéo được, căn / chia đều / lật), hộp ánh xạ khi đổi vùng (3 cách, xem trước); danh sách bước đồng bộ. `npm test` + `npm run build` OK; thử headless (CDP) chọn khung, ⌘A, xoay, co giãn, xiên, phối cảnh, ô Rộng 50%, căn trên, ⌘Z ×2, đổi vùng *Co giãn* – toạ độ khớp tính tay. Chưa làm ánh xạ 4 điểm mốc; chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu: overlay gọn hơn – thanh công cụ tự ẩn (hiện khi chuột sát mép trên / rời overlay qua mép trên / rê vào gợi ý, lúc mở hiện 2,5 s); banner hướng dẫn thành viên gợi ý một dòng (tên công cụ + câu ngắn, thông báo tự tắt sau 3,5 s); ô toạ độ thu nhỏ góc trên trái; thanh biến đổi / hộp ánh xạ dời lên. `npm run build` OK, chụp headless trạng thái ẩn / hiện; chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu: thêm bước 31 (kết nối AI qua MCP để AI tạo / sửa / chạy kịch bản từ câu lệnh bằng lời), chia 31a (tạo / sửa kịch bản, sau 21b) · 31b (AI nhìn màn hình, sau 25) · 31c (AI tự thử và chạy, sau 26); cập nhật bảng thứ tự Phase 3, quyết định đã chốt, bước 22 dùng chung server, thêm ví dụ mở bài nhạc. Chỉ sửa plan, chưa code.
- **2026-10-17** – Theo yêu cầu: MCP làm sau – dời 31a / 31b / 31c xuống sau 22 (dùng chung server HTTP + token); thêm quy tắc *giữ đường cho MCP* vào quyết định đã chốt và mục *Chuẩn bị cho MCP* trong 21b, 23, 25, 25c, 26, 22. Chỉ sửa plan.
- **2026-10-17** – Bước 21b (input của kịch bản): Go `inputs.go` – `Script.inputs`, `resolveInputs` (ghi trường ghim, thay `{{biến}}`, `\\{{` thoát, lỗi thiếu / sai kiểu / ngoài khoảng / biến lạ / input không dùng), `ScriptError` có cấu trúc (validateScript cũng trả loại này – chuẩn bị MCP), giá trị lần trước + preset trong `inputs.json`, file `version` 2. Frontend: ghim thật (`PinMenu`: input mới hoặc dùng chung), ô chữ có biến, tab Input (sửa key đổi luôn `{{key}}`), chip "N input", form Chạy thật (preset lưu / xoá, lỗi trong form, chạy luôn + "Hỏi lại"), thanh mini hiện bộ giá trị. `go test` OK (8 test mới), `npm test` 24 test OK, `npm run build` OK, chạy headless trên trình duyệt (ghim, menu, biến, tab Input, form Chạy thật + #demo). Chưa kiểm tra tay trong app, chưa commit.
- **2026-10-17** – Theo câu hỏi của người dùng (đã bật "Lần sau chạy luôn" thì sửa input thế nào): Shift + bấm Chạy luôn mở form; form tick sẵn "Lần sau chạy luôn" khi đang bật (bỏ tick = lần sau hỏi lại); nút Chạy có gợi ý khi đang chạy luôn; vẫn còn nút "Hỏi lại" ở tab Input. `npm run build` + `npm test` OK; chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Người dùng xác nhận kiểm tra tay bước 21b OK → tick, commit. Bắt đầu bước 23.
- **2026-10-17** – Bước 23 (cuộn + mở URL / app): Go – `Step.direction/amount/target/app`, `open.go` (lệnh mở theo hệ điều hành, Windows dùng PowerShell `Start-Process`), engine cuộn từng nấc, `validateScript` / `checkStep` cho 2 loại mới, biến trong URL tự mã hoá, `amount` ghim được; `steptypes.go` (mô tả loại bước cho MCP). Frontend: 2 loại bước trong Ctrl+K, inspector, thẻ bước, Xem trước, overlay (cuộn tại điểm). `go test` OK, `npm test` 25 test OK, `npm run build` OK, headless (thêm Cuộn tại điểm, Mở URL có biến + Chrome). Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu (bước cuộn chưa hình dung được cuộn tới đâu): phát hiện robotgo gửi cuộn macOS theo **pixel** (Windows theo nấc) – `ScrollDir(1)` trên Mac chỉ cuộn 1 px → macOS gửi 100 px / nấc (`scrollNotchPx`, khớp ≈ 1 nấc Chrome / Edge trên Windows). Thêm `ScrollMap` trong inspector, mũi tên cuộn trên overlay, hiệu ứng mũi tên trong Xem trước, ≈ px trên thẻ bước. `go test`, `npm test`, `npm run build` OK; chụp headless inspector / overlay / Xem trước. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu (cuộn dài hơn màn hình; Xem trước cần thứ gì đó cho thấy đang cuộn tới đâu): ô vị trí của bước cuộn đổi tên "Đặt chuột tại" + cảnh báo ngay khi điểm ngoài vùng / màn hình (trước chỉ báo lúc Chạy); thêm ô "Quãng cuộn ≈ px" (tự đổi ra nấc); `ScrollMap` vẽ trang thành từng màn (màn 2, 3…, quá 3 màn vẽ gãy) + "≈ n màn hình"; Xem trước có **trang thu nhỏ** cạnh con trỏ (dòng chữ giả, khung màn hình trượt dọc trang, "đã cuộn x / y px · a / b màn"); bỏ tiền tố "Error:" trong lỗi. `npm test`, `tsc` OK; chụp headless. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu: Xem trước có nút **"Trang: …"** chọn kiểu trang mẫu cho trang thu nhỏ khi cuộn (`components/MockPage.tsx`): Văn bản, Trang web, Danh sách, Lưới, Feed, Bảng – khối vẽ theo kích thước thật (dòng 24 px, dòng danh sách 64 px, hàng ô 300 px, bài 520 px, hàng bảng 36 px) rồi thu nhỏ cùng tỷ lệ; dưới trang ghi "≈ 23 dòng danh sách"…; nhớ lựa chọn trong localStorage. `npm run build`, `npm test` OK; chụp headless menu + 3 kiểu. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu (trang mẫu đè lên nền luôn): Xem trước phủ **trang mẫu kín vùng / màn hình, đúng tỷ lệ px thật** (`PageBackdrop`), mỗi bước cuộn làm trang trượt đúng quãng (cộng dồn trong lượt, `scrollRange` tính độ dài trang); khối được đánh số (dòng #1, Bài 2, Mục 3…) để thấy đã cuộn tới đâu; nút đổi thành "Nền: …", thêm lựa chọn "Không" (nhìn màn hình thật); trang thu nhỏ vẫn giữ. `npm run build`, `npm test` OK; chụp headless trước / giữa / sau khi cuộn. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Theo yêu cầu: nền trang mẫu trong suốt hơn (độ đục 0,88 → 0,62) và theo theme của app (sáng / tối / theo hệ điều hành – `useDarkPage`, bảng màu `LIGHT` / `DARK` trong `MockPage.tsx`; overlay vẫn HUD tối). Áp cho cả nền, trang thu nhỏ và hình trong menu chọn. `npm run build`, `npm test` OK; chụp headless sáng / tối. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Người dùng xác nhận bước 23 OK → tick, commit. Bắt đầu bước 25.
- **2026-10-17** – Theo yêu cầu: logo chính = `frontend/src/assets/logo-app.svg` – component `Logo` (thanh tiêu đề, PreviewScreen; bỏ nền vuông trắng), `build/appicon.png` 1024 (ô bo góc kiểu macOS, vẽ bằng Chrome headless) + `build/windows/icon.ico` 7 cỡ, favicon + tiêu đề "Auto Clicker" trong `index.html`.
- **2026-10-17** – Bước 25 (chờ ảnh, click theo ảnh, chờ pixel): xem các mục trong bước 25. `go test` OK (match + engine với màn hình giả + docs), test thật trên máy dev OK (quyền có, tỷ lệ 2, khớp 1.000), `npm run build`, `npm test` OK; headless: khoanh ảnh, inspector 3 loại, logo. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Sửa (người dùng báo: Thử tìm > 90% nhưng chạy thật chỉ 41%): lúc chạy, thanh mini (giữa trên màn hình) và **con trỏ chuột** lọt vào ảnh chụp (robotgo / ScreenCaptureKit mặc định vẽ con trỏ), Thử tìm thì app đã ẩn. macOS ≥ 14: tự chụp bằng ScreenCaptureKit **bỏ cửa sổ của chính app + không vẽ con trỏ** (`capture_darwin.go`, cgo Objective-C); Windows: `SetWindowDisplayAffinity(WDA_EXCLUDEFROMCAPTURE)` cho cửa sổ app lúc khởi động (`internal/winexclude`, vet GOOS=windows OK, chưa chạy thật); đọc màu pixel cũng chụp 1×1 theo cách này; khoanh ảnh / thử tìm không cần ẩn app nữa (hết nhấp nháy) – chỉ ẩn khi hệ không hỗ trợ. Chuyện đặt 40% vẫn dừng ở 41%: chờ người dùng gửi câu lỗi đầy đủ (bước nào, cần bao nhiêu %). `go test`, test thật OK. Chưa commit.
- **2026-10-17** – Từ ảnh người dùng gửi: bước dừng là **bước 2** (Chờ ảnh khớp ≥ 90%), người dùng hạ ngưỡng ở bước 4 – code đúng, giao diện chưa cho thấy. Thêm: ảnh mẫu thu nhỏ trên thẻ bước; chạy lỗi → Go gửi `engine:failed` (id bước + chỗ giống nhất nếu chờ ảnh hết giờ) → tự chọn bước lỗi, thẻ viền đỏ + nhãn "Dừng ở đây", inspector "Lần chạy trước hết giờ không thấy – chỗ giống nhất khớp n%", xem trên màn hình; cảnh báo khi độ khớp < 70%. `go test`, `npm test`, `npm run build` OK; chụp headless. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Sửa (người dùng báo): khoanh ảnh sát đáy màn hình thì nút "Dùng ảnh này / Khoanh lại" nằm dưới Dock – thêm `safeArea()` (`windowMode.ts`, theo `screen.availTop/availHeight`, trừ Dock / taskbar / thanh menu); nút đặt dưới khung → trên khung → cạnh khung trong vùng an toàn, nhãn cỡ né nút; phím **Enter = Dùng ảnh này**, Esc = khoanh lại; thanh số của Chọn nhiều cũng theo `safeArea`. `npm run build` OK, headless khung sát đáy → nút lên trên. Chưa kiểm tra tay, chưa commit.
- **2026-10-17** – Kế hoạch: thêm bước **25d. Web qua DOM** (tách từ 28, theo yêu cầu người dùng) – tìm phần tử bằng HTML qua CDP, đo toạ độ trang → màn hình bằng chuột, click / cuộn bằng robotgo; xếp sau 26, trước 20 (bảng thứ tự + *Quyết định đã chốt*). Chỉ sửa tài liệu.
- **2026-10-17** – Bước 17: chọn màn hình cho overlay khi có nhiều màn (`Displays`, `SetOverlayDisplay`, nút Màn 1 / 2 trên thanh overlay, *Toàn màn này*). `go test ./...`, `npm run build`, `npm test` OK; chưa kiểm tra tay trên 2 màn.
- **2026-10-17** – Sửa (người dùng báo): chờ ảnh lỗi khi vùng ở màn phụ – chụp / tìm ảnh / lấy màu được trên mọi màn (ScreenCaptureKit theo màn chứa vùng). `go test ./...` OK, chụp thật màn MacBook OK; chưa chạy lại kịch bản của người dùng.
- **2026-10-17** – Bước 25e: tìm theo chữ (Vision macOS) cho Chờ ảnh / Click theo ảnh + sửa độ chính xác tìm ảnh (khác tỷ lệ màn, mẫu ít tương phản). `go test ./...`, `npm run build`, `npm test` OK; thử thật đọc chữ + tìm ảnh trên 2 màn. Chưa kiểm tra tay trong app.
- **2026-10-17** – Bước 25b (Khi thất bại: thử lại / bỏ qua / nhảy tới bước) làm sớm; kế hoạch bước 20 viết lại thành khối điều kiện (Nếu / Lặp / Nhóm / Thoát / Dừng, cây bước, version 2), xếp lên trước 25c. `go test ./...`, `npm run build`, `npm test` OK; chưa kiểm tra tay.
- **2026-10-17** – Bước 20a: khối Nếu / Nhóm-Lặp / Thoát / Dừng + điều kiện (ảnh, chữ, màu, input, bước trước, và / hoặc) ở model + engine + kiểm tra, file version 3. `go test ./...` (8 test mới), `npm run build`, `npm test` OK; thử thật khối Nếu với OCR thanh menu. Chưa có UI soạn khối (20b).
- **2026-10-17** – Bước 20b: danh sách bước dạng cây (thụt lề, đánh số 9.1, Nếu không, thêm vào khối, gập / mở, kéo vào / ra khối), store + input + overlay + xem trước duyệt cây, Ctrl+K thêm khối / bọc bước, inspector khối cơ bản. `npm run build`, `npm test` (32), `go test ./...` OK; chụp headless. Chưa kiểm tra tay kéo thả.
- **2026-10-17** – Bước 20c: điều kiện thấy ảnh / màu pixel / nhiều điều kiện (và / hoặc) trong inspector, khoanh ảnh + chọn điểm cho điều kiện, Thử điều kiện (`TestCond`). `go test ./...`, `npm run build`, `npm test` (34) OK; thử thật TestCond trên thanh menu. Chưa kiểm tra tay trong app.
- **2026-10-17** – Sửa (người dùng báo): overlay mở trên màn chứa vùng đang chọn thay vì màn chứa cửa sổ app / màn chọn lần trước. `go test ./...`, `npm run build`, `npm test` OK; chưa kiểm tra tay 2 màn.
- **2026-10-17** – Sửa (người dùng báo): bước tìm theo chữ hiện đúng tên "Chờ chữ" / "Click theo chữ" (`stepLabel`). `npm run build`, `npm test` OK.
- **2026-10-17** – Người dùng thấy cây khối (20a–20c) khác kỳ vọng → revert 3 commit `60dcfe2`, `956ec76`, `d5882a7` (giữ 2 bản sửa overlay / `stepLabel` sau đó); viết lại plan 20 thành flow kiểu n8n (node + cạnh, `StepResult`, canvas, F1–F7, gộp 25c). `go test ./...`, `npm run build`, `npm test` (25) OK sau revert.
- **2026-10-17** – Sắp lại các mục Phase 3 theo đúng thứ tự làm (20 lên ngay sau 25b, trước 26 / 25d; 25e vào bảng thứ tự); bỏ mục 25c riêng – nội dung nằm trong 20 F5.
- **2026-10-17** – Theo yêu cầu: sắp lại thứ tự trong 20 – F4a (thiết kế + `#demo-flow`) lên đầu, F4b (nối canvas) ngay sau engine, rồi mới F3; bỏ chuyển đổi file kịch bản cũ (version < 4 không mở được); thêm bảng chia module Go / frontend.
