# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** 29 (neo 4 góc) – xong code + test, chờ kiểm tra tay. Bước 17 (phần Windows) và 18 hoãn: chưa có máy Windows
- **Tiếp theo:** 21 → 30 → 21b → 23 → 25 … (bảng thứ tự đầu Phase 3). Khi có máy Windows: quay lại 17–18
- **Cập nhật lần cuối:** 2026-10-16

## Quyết định đã chốt

- Dev trên **macOS** (Apple Silicon); build Windows làm sau (cần MinGW-w64). Trên macOS phải cấp quyền *Accessibility* cho app/terminal thì robotgo mới điều khiển được chuột.
- **Tailwind v4** (`@tailwindcss/vite`), token màu/font khai báo trong `@theme` ở `frontend/src/style.css` – không dùng `tailwind.config.js`.
- Font **Be Vietnam Pro** + **IBM Plex Mono** cài local qua `@fontsource/*` (app phải chạy được offline).
- UI theo **phương án A (MainScreen)** cho MVP. Phương án B (PreviewScreen) giữ làm tham khảo.
- **Design v2** (2026-10-16, `design/DESIGN.md`) thay bản thiết kế cũ `autoclicker-ui/` từ bước 27: theme sáng / tối, một màu nhấn, gradient chỉ ở 4 chỗ, bố cục 3 cột kiểu Windows 11; font Plus Jakarta Sans + JetBrains Mono (vẫn cài local); bỏ màu riêng cho kéo thả và nhóm `dark-*`.
- Thứ tự Phase 3 (2026-10-16, bảng đầu Phase 3): 27 → 27b → 29 → 21 → 30 → 21b → 23 → 25 → 25b → 25c → 26 → 20 → 22 → 24 → 28; chờ ảnh / click theo ảnh làm sớm (cốt lõi cho web / app), nhóm bước làm sau xử lý hết giờ. Không làm bản chạy trên điện thoại.
- **Theme tối = bảng Slate** (`design/newdarkmode.md`, 2026-10-16): theme cũ nhiều nút bị chìm. Có thêm 2 bảng Dimmed / Nord trong file, chưa dùng.
- **Neo điểm theo trục** (29): X neo trái / phải, Y neo trên / dưới (thêm giữa, %), lưu offset dương + tên neo, UI hiện dạng ±; làm ngay sau UI (27, 27b), trước 21, để định dạng file có neo ngay từ đầu (không cần chuyển đổi). Biến đổi nhiều điểm (30) làm sau 27 + 21, ghi lại offset theo neo.
- **UI trước, logic sau** (27b): tính năng chưa có vẫn hiện trong app nhưng tắt + nhãn *Sắp có*; dữ liệu giả chỉ ở màn `#demo` (`frontend/src/demo/`), không thêm loại bước chưa làm vào `types.ts` / `model.go`.
- **Input của kịch bản** áp vào bản sao trước khi chạy (`resolveInputs`), engine không đổi; giá trị lần trước / preset lưu trên máy, không ghi vào file kịch bản; input trong nhóm lặp dùng cùng giá trị mọi vòng (mỗi vòng một giá trị → CSV, bước 28).
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
| 9–11 | 25b → 25c → 26 | Xử lý khi không thấy ảnh, debug chỗ dễ lỗi, focus cửa sổ – xong 26 là chạy được ví dụ Facebook |
| 12 | 20. Nhóm bước | Dùng lại thử lại / bỏ qua của 25b |
| 13–15 | 22 → 24 → 28 | Chạy không cần canh, ghi macro, mở rộng |
| – | 17 (Windows), 18 | Khi có máy Windows; riêng *Kiểm tra tay macOS* của 17 làm lúc nào cũng được (nên làm trước 25 vì chụp ảnh phụ thuộc tỷ lệ màn hình) |

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
  - [x] Debug: nút Chạy từ bước này / Tiếp (F10), breakpoint trên thẻ bước, màn timeline lịch sử chạy *(chưa có trong design; logic: 25c)*
  - [x] Hẹn giờ (chip giờ trong thư viện, hộp đặt giờ), nút Ghi macro trên toolbar *(chưa có trong design; logic: 22, 24)*
  - Xong khi: `#demo` hiện đủ các phần trên ở cả 2 theme, khớp canvas; app thật không có nút nào bấm được mà không chạy; `npm run build` sạch
  - Kết quả: màn mẫu `#demo` (màn soạn đủ loại bước + nhóm, thư viện, inspector từng loại, tab Input, form Chạy, hẹn giờ, ghi macro, debug), `#demo-overlay` (khoanh ảnh, chọn nhiều + biến đổi + ánh xạ), `#demo-minibar` (các trạng thái thanh mini). Code dùng lại khi làm logic: `components/StepCard.tsx`, `components/inputs/*`, `components/steps/*` (model + `StepFields` + `AnchorField`), `ScheduleDialog.tsx`, `AppChrome.tsx`, `MiniBarView`, `LibraryView`. Dữ liệu mẫu chỉ trong `frontend/src/demo/`
  - [x] Người dùng xem `#demo`, `#demo-overlay`, `#demo-minibar` và màn chính (nút tắt + *Sắp có*), góp ý trước khi commit
- [ ] **29. Neo điểm theo 4 góc vùng** *(người dùng yêu cầu, 2026-10-16)* – Hiện điểm chỉ tính từ góc trên trái của vùng nên điểm ở phía phải / dưới lệch khi đổi vùng. Mỗi trục neo riêng: X neo trái (+) hoặc phải (−), Y neo trên (+) hoặc dưới (−). Ví dụ 4 điểm cách mỗi cạnh 300 px = (+300, +300), (−300, +300), (+300, −300), (−300, −300); đổi vùng thì mỗi điểm bám theo góc của nó.
  - [x] Model: `Point` thêm `ax` (`start` = trái | `end` = phải | `center` | `ratio` = theo % rộng), `ay` (tương tự cho trên / dưới); bỏ trống = `start` (file / kịch bản cũ chạy như cũ). Lưu offset dương + neo riêng thay vì dùng dấu (để phân biệt +0 / −0); UI hiện dạng có dấu ±
  - [x] Toàn màn hình (`region = nil`) cũng neo được: chốt **giữ `nil`**, lúc tính thì khung = màn hình chính (`mainScreen`, Go; frontend lấy qua `ScreenBounds`) – kịch bản mang sang máy khác cỡ màn hình vẫn bám góc. Đa màn hình xử lý ở 17
  - [x] Go: một hàm `resolvePoint(p, region)` dùng chung cho engine (`runStep.pt`) và `validateScript` (điểm phải nằm trong vùng); test 4 góc, center, ratio, đổi cỡ vùng, kịch bản cũ không có neo
  - [x] Frontend: `toScreen` / `toRelative` / `endpoints()` trong `scriptStore.ts` tính theo neo (khớp Go); đổi neo giữ nguyên chỗ trên màn hình (tính lại offset)
  - [x] Neo khi đặt điểm: mặc định **tự chọn góc gần nhất** theo từng trục; kéo marker / chọn lại giữ neo hiện tại
  - [ ] Cài đặt kịch bản *góc gần nhất / luôn trên trái* – chưa làm (cần chỗ lưu cài đặt kịch bản, làm cùng 21)
  - [x] Form: ô X / Y có dấu ± – gõ số âm = neo phải / dưới; nút chọn neo dạng lưới 3×3 nhỏ (4 góc, giữa cạnh, tâm) + tuỳ chọn %; danh sách bước hiện toạ độ dạng `(−300, +300)`
  - [x] Overlay: đường mảnh từ marker tới cạnh / góc đang neo (hiện với bước đang chọn hoặc khi giữ Alt); đổi vùng (kéo vùng mới) thấy marker bám theo góc ngay
  - [x] Kéo thả "lệch so với A" (ΔX, ΔY) giữ nguyên – lệch không phụ thuộc neo
  - [x] Nối UI khung từ 27b (ô ±, lưới chọn neo, đường neo trên overlay)
  - Xong khi: đặt 4 điểm ở 4 góc vùng 1000×700, đổi vùng sang 1400×900 → 4 điểm vẫn cách mỗi cạnh đúng 300 px khi chạy thật và trên overlay
- [ ] **21. Lưu / tải kịch bản** – App chưa lưu gì; thêm `version` trước khi đổi model nhiều.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `Script.version`; file `.json`
  - [ ] Go: `SaveScript`, `SaveScriptAs`, `OpenScript` qua `runtime.SaveFileDialog` / `OpenFileDialog`; `validateScript` khi mở; danh sách mở gần đây (≤ 5) + bản nháp tự lưu trong `os.UserConfigDir()/autoclicker/`
  - [ ] Thư viện ở thanh bên [`Main`]: ô tìm, mục đang mở (`bg-hover` + vạch trái 3×16 `bg-accent`), chấm `warn` khi chưa lưu, "Mở gần đây", nút Kịch bản mới / Mở / Lưu
  - [ ] Nhãn "Chưa lưu" cạnh tên kịch bản; status bar "Bản nháp tự lưu lúc …"; chấm `ok` "Sẵn sàng"
  - [ ] Phím Ctrl+S, Ctrl+Shift+S, Ctrl+O, Ctrl+N (macOS: Cmd); hỏi lưu khi đóng app / đổi kịch bản
  - [ ] Test: ghi rồi đọc lại khớp; file hỏng / sai `version` báo lỗi
  - Xong khi: tắt app mở lại vẫn còn bản nháp; mở / lưu qua hộp thoại trên macOS
- [ ] **30. Sửa nhiều điểm cùng lúc: di chuyển, co giãn, xoay, xiên** *(người dùng yêu cầu, 2026-10-16)* – Chọn nhiều marker trên overlay (chế độ sửa và tab Xem trước) rồi biến đổi cả nhóm để khớp vùng mới.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Chọn nhiều: Shift / Cmd+click marker, kéo khung chọn trên chỗ trống (khi đang ở công cụ *Chọn*), Ctrl+A chọn tất cả, Esc bỏ chọn; danh sách bước đồng bộ chọn nhiều
  - [ ] Khung biến đổi quanh các điểm đã chọn: kéo thân = **di chuyển** (phím mũi tên 1 px, Shift 10 px); 8 núm = **co giãn** (Shift giữ tỷ lệ, Alt từ tâm); núm ngoài góc = **xoay** (Shift bước 15°); Cmd/Ctrl + kéo núm cạnh = **xiên (skew)**; Cmd/Ctrl + kéo núm góc = **biến dạng tự do 4 góc (perspective)**; nút **lật ngang / dọc**
  - [ ] Ô nhập số trên thanh công cụ overlay: ΔX, ΔY, rộng / cao %, góc xoay, xiên X / Y – để nhập chính xác
  - [ ] Căn chỉnh: căn trái / phải / trên / dưới / giữa, chia đều khoảng cách ngang / dọc
  - [ ] **Ánh xạ sang vùng mới**: khi đổi vùng hỏi cách xử lý các điểm – *Giữ theo neo* (bước 29, mặc định) / *Co giãn theo vùng* (giữ tỷ lệ vị trí trong vùng) / *Giữ nguyên chỗ trên màn hình*; hoặc chọn 4 điểm mốc cũ → 4 điểm mốc mới để tính phép biến đổi
  - [ ] Toán: một hàm áp ma trận affine (di chuyển / co giãn / xoay / xiên) + homography 4 điểm cho perspective, viết trong `frontend/src/transform.ts`, có test đơn vị; sau khi biến đổi, ghi lại offset theo neo hiện có của từng điểm (bước 29)
  - [ ] Điểm đặc biệt: A "từ điểm trước" bỏ qua (tự tính); B "lệch so với A" biến đổi véc-tơ ΔX/ΔY theo phần tuyến tính (co giãn / xoay / xiên, không cộng di chuyển); điểm bị đẩy ra ngoài vùng → báo và không áp dụng
  - [ ] Hoàn tác / làm lại (Ctrl+Z / Ctrl+Shift+Z) cho mọi thay đổi điểm – hiện chỉ có hoàn tác xoá; lịch sử trong `scriptStore.ts`
  - [ ] Tab Xem trước: bật *Sửa điểm* cho phép chọn nhiều + biến đổi, mô phỏng chạy lại theo vị trí mới
  - Xong khi: chọn 5 điểm, co giãn 150%, xoay 10°, xiên, lật → overlay và lần chạy thật khớp nhau; Ctrl+Z quay lại đúng từng bước; đổi vùng với cả 3 cách xử lý cho kết quả như mô tả
- [ ] **21b. Input của kịch bản** – Người soạn đánh dấu giá trị nào được đổi khi chạy; bấm Chạy thì điền form. Go áp giá trị vào bản sao kịch bản rồi mới đưa engine, **engine không đổi**. Cùng đổi `version` với 21; CSV (28) dựa vào đây.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model `Script.inputs`: `key`, `label`, `type` (`number`, `ms`, `text`, `multiline`, `bool`, `select`, `point`), `default`, `required`, `min` / `max`, `options`, `hint`, `targets: [{stepId, field}]`. Ví dụ:
    ```json
    "inputs": [
      { "key": "click_delay", "label": "Chờ sau khi click", "type": "ms", "default": 500,
        "targets": [{ "stepId": "s4", "field": "delayAfter" }] },
      { "key": "M1", "label": "Tên người nhận", "type": "text", "required": true, "hint": "Tên hiện trong lời chào" }
    ]
    ```
    - Nguồn: **trường được ghim** (`targets`, một input cho nhiều bước) và **biến `{{key}}`** trong trường chữ (gõ chữ, URL, tên cửa sổ) – tự tạo input khi gõ biến mới
    - Key: chữ không dấu, số, `_`; nhãn tiếng Việt ở `label`; `\{{` để gõ ra `{{`; không ghim được trường tự tính (kéo thả "từ điểm trước")
  - [ ] Go: `resolveInputs(script, values) (Script, error)` – sao chép, ghi trường ghim, thay `{{…}}`, `validateScript`; `RunScript(sc, values, countdown)`; Xem trước dùng chung resolve
  - [ ] Lỗi rõ: thiếu bắt buộc, sai kiểu / ngoài khoảng, biến không có trong `inputs`, input không bước nào dùng, toạ độ ngoài vùng
  - [ ] Giá trị lần trước + preset lưu trong `UserConfigDir` theo từng kịch bản, **không** ghi vào file kịch bản
  - [ ] Inspector [`StepInputs`]: `components/PinButton.tsx` (32×32 cạnh mỗi trường số / chữ; đã ghim nền `accent-soft`); trường đã ghim `field-pinned` + chip key `var-token` cạnh nhãn
  - [ ] Ô chữ có biến [`StepInputs`]: `components/VarTextInput.tsx` – tô `{{biến}}` bằng `var-token`; gõ `{{` mở popup gợi ý (key mono `accent-text` + nhãn `muted`, "+ Tạo input mới…"); dòng phụ "N biến, tự thêm vào tab Input"
  - [ ] Tab Input [`InputsTab`]: `components/InputsPanel.tsx` – chip số input ở tab; mỗi dòng ⋮⋮ kéo, nhãn + `*`, `key · mặc định` mono, chip kiểu, chip bước dùng; input đang sửa mở rộng (`ring-selected`): nhãn, kiểu, mặc định, bắt buộc, gợi ý; bỏ ghim; "+ Thêm input"
  - [ ] Thẻ bước: chip viền `line` "N input"
  - [ ] Form Chạy [`RunDialog`, `RunDialogDark`]: `components/RunDialog.tsx` (popup 560, bo 12) – đầu có tên kịch bản + số input; "Bộ giá trị" (select preset + "Lưu bộ mới"); ô nhóm theo bước, ô đã đổi `ring-focus` + "đã đổi, mặc định … · khôi phục", `*` màu `warn`; khung "Sẽ gõ" xem trước câu sau khi thay biến; chân: "Lần sau chạy luôn với giá trị này", "Đếm ngược 3 giây", Huỷ / Chạy (`btn-primary`, F8)
  - [ ] Nút Chạy / F8 mở form khi kịch bản có input; chạy thẳng nếu đã bật "Lần sau chạy luôn"; thanh mini hiện tên preset
  - [ ] Test: ghi đè trường; thay biến (có tiếng Việt); thiếu biến; sai kiểu; escape; một input nhiều bước; input trong nhóm (sau bước 20)
  - Xong khi: B1 click 500 → chạy với 1000 ms; B2 "xin chào {{M1}} và {{M2}}" gõ đúng giá trị điền
- [ ] **23. Cuộn + Mở URL / app** – Loại bước đơn giản.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `scroll` (`direction`, `amount`, điểm tuỳ chọn); `open` (`target` URL / đường dẫn app, `app` tuỳ chọn)
  - [ ] Engine: `robotgo.Scroll`; `open` → macOS `open [-a App] target`, Windows `cmd /c start "" target`; test với input giả
  - [ ] UI: icon `scroll`, `open` (bảng icon DESIGN.md), thêm vào bảng Ctrl+K; inspector: hướng + số nấc; URL / app + trình duyệt (URL nhận biến, vd `facebook.com/{{profile}}`)
- [ ] **25. Chờ ảnh, click theo ảnh, chờ pixel** – Cốt lõi để tương tác web / app (làm sớm).
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `waitImage` (`template` PNG cạnh file kịch bản hoặc base64, `threshold` %, `timeoutMs`, `intervalMs`, vùng tìm); `clickImage` (dùng kết quả `waitImage` hoặc tự tìm, lệch ΔX/ΔY); `waitPixel` (điểm, màu, sai số, timeout)
  - [ ] Engine: chụp màn hình (`robotgo.CaptureImg`), so khớp mẫu (bắt đầu đơn giản trên ảnh thu nhỏ, cân nhắc `gocv` nếu chậm); lưu vị trí tìm thấy cho bước sau
  - [ ] Overlay tool `capture` [`Overlay`]: khoanh ảnh mẫu, nhãn kích thước, "Dùng ảnh này" / "Khoanh lại"
  - [ ] Inspector "Chờ ảnh" [`Main`]: ảnh mẫu, Chụp lại, Xem trên màn hình, độ khớp, chờ tối đa (ghim được), kiểm tra mỗi
  - [ ] Icon `waitImage`, `clickImage`, `waitPixel`; thêm vào bảng Ctrl+K
  - [ ] Thanh mini [`MiniBar`]: vòng tiến độ theo thời gian chờ, "đang tìm x / y s"
  - [ ] Retina / DPI ≠ 100%: chụp và so cùng tỷ lệ (liên quan 17)
- [ ] **25b. Xử lý khi hết giờ** – Chờ ảnh cần hướng xử lý khi không thấy.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `onTimeout: { action: "retry" | "skip" | "stop", retries, then: "skipGroup" | "stop" }`
  - [ ] Engine: thử lại N lần, bỏ qua bước, bỏ qua phần còn lại của nhóm, hoặc dừng (lý do vào `engine:stopped`)
  - [ ] UI [`Main`]: segmented Thử lại / Bỏ qua / Dừng, số lần thử, "Thử hết vẫn không thấy"; chip trên thẻ bước; thanh mini "thử lại x / y"
- [ ] **25c. Debug từng bước + lịch sử chạy** – Chờ ảnh là chỗ dễ lỗi nhất.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Chạy từ bước đang chọn; chạy từng bước (nút Tiếp / F10); breakpoint trên thẻ bước
  - [ ] Mỗi lần chạy lưu nhật ký: thời gian từng bước, giá trị input đã dùng, ảnh chụp màn hình khi bước lỗi / hết giờ
  - [ ] UI: mục "Lịch sử" trong thư viện; xem lại từng lần chạy dạng timeline
- [ ] **26. Focus cửa sổ theo tên**
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `focusWindow` (`title` hoặc process); Engine: `robotgo.ActivateName`, không thấy → `onTimeout`
  - [ ] UI: icon `focusWindow`; ô nhập (nhận biến) + chọn từ danh sách cửa sổ đang mở
- [ ] **20. Nhóm bước** – Lặp riêng một nhóm N lần; dùng chung thử lại / bỏ qua của 25b.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Model: `group` (`name`, `repeat` – ghim được, `children`); tăng `version` + chuyển đổi file cũ
  - [ ] Engine: chạy lồng, tiến độ theo bước con; `skipGroup` từ 25b
  - [ ] UI [`Main`]: khung `border-dashed border-line-strong bg-group rounded-group p-2`, đầu khung icon `group` + tên + chip `lặp ×N`; kéo bước vào / ra; overlay + Xem trước đánh số theo thứ tự thực tế
- [ ] **22. Hẹn giờ + thông báo + CLI / API local** – Chạy không cần ngồi canh.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Hẹn giờ: chạy lúc HH:mm (một lần / hằng ngày) với preset input đã chọn; chip giờ trong thư viện; app thu xuống khay
  - [ ] Thông báo: Telegram / Discord / webhook khi xong hoặc lỗi, kèm ảnh chụp màn hình
  - [ ] CLI: `autoclicker run script.json --preset "Gửi cho Lan" --set M1=Lan`
  - [ ] API local: `POST /run` (kèm values), `POST /stop`, `GET /status`, tiến độ qua SSE; bảo vệ bằng token
- [ ] **24. Ghi macro** – Ghi click / kéo / gõ phím bằng gohook, gộp thành bước, tự tính delay; nút Ghi trên toolbar, F9 dừng.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
- [ ] **28. Mở rộng** – Làm sau cùng.
  - [ ] Chrome qua DOM (`chromedp`): `web.goto`, `web.waitVisible`, `web.click`, `web.type`, `web.scroll` theo selector; nút "chọn phần tử trên trang"
  - [ ] CSV: mỗi cột là một input key, chạy kịch bản cho từng dòng (mỗi dòng một lần `resolveInputs`)
  - [ ] OCR: đọc chữ trong vùng vào biến, rẽ nhánh theo nội dung (Tesseract)
  - [ ] Tìm phần tử bằng AI: mô tả bằng lời ("nút Gửi"), model có vision trả toạ độ; dự phòng khi ảnh mẫu không khớp
  - [ ] Biến và điều kiện: đọc clipboard / OCR vào biến, `if / else`, "lặp đến khi thấy ảnh"
  - [ ] Thư viện mẫu, chia sẻ gói kịch bản (kèm ảnh mẫu), thống kê lần chạy (SQLite)


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
