# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** chưa chọn – xong 19, 19b, 19c. Bước 17 (phần Windows) và 18 hoãn: chưa có máy Windows
- **Tiếp theo:** bước 20 (nhóm bước), 21 (lưu / tải kịch bản JSON). Khi có máy Windows: quay lại 17–18
- **Cập nhật lần cuối:** 2026-10-15

## Quyết định đã chốt

- Dev trên **macOS** (Apple Silicon); build Windows làm sau (cần MinGW-w64). Trên macOS phải cấp quyền *Accessibility* cho app/terminal thì robotgo mới điều khiển được chuột.
- **Tailwind v4** (`@tailwindcss/vite`), token màu/font khai báo trong `@theme` ở `frontend/src/style.css` – không dùng `tailwind.config.js`.
- Font **Be Vietnam Pro** + **IBM Plex Mono** cài local qua `@fontsource/*` (app phải chạy được offline).
- UI theo **phương án A (MainScreen)** cho MVP. Phương án B (PreviewScreen) giữ làm tham khảo.
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
