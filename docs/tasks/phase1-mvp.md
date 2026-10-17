# Phase 1 – MVP

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
  - Kiểm tra tay: xem [checks.md](checks.md)
  - [x] Dừng khẩn cấp nhận góc của mọi màn hình: `displays()` (`input.go`, Windows nhân `robotgo.ScaleF` cho cùng hệ toạ độ với `robotgo.Location`), `hardCorners` bỏ góc tiếp giáp màn khác; test `TestHardCorners`. Máy dev: màn ngoài 1920×1080 + MacBook 1512×982 tại (1920, 406) → 6 góc cứng
  - [x] Overlay đo **tỷ lệ + độ lệch** (`coordMap.ts`, bình phương nhỏ nhất trên các mẫu chuột đứng yên) thay vì chỉ độ lệch – để scale ≠ 100% trên Windows vẽ vùng / marker đúng chỗ; điểm đặt / kéo marker vẫn lấy thẳng từ `MousePos`
  - Ghi chú: robotgo trên Windows tự nhân toạ độ chuột với scale màn chính (`Location`) và chia lại khi `Move` → lấy toạ độ bằng robotgo rồi click lại luôn khớp; chỉ phép đổi sang px trong cửa sổ mới cần tỷ lệ
  - Kiểm tra tay: xem [checks.md](checks.md)
  - [ ] Kiểm tra trên Windows ở scale 100% / 125% / 150% và 2 màn khác scale (cần máy Windows – làm cùng bước 18)
- [ ] **18. Build & thử máy sạch** – `wails build` ra .exe, chạy trên máy chưa cài Go/Node.
