# 25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật

- [ ] **25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật** *(người dùng yêu cầu, 2026-10-17; tách từ 28)* – Tìm phần tử trên trang theo chữ / selector (vd nút "Xem lịch sử"), tính toạ độ màn hình của nó rồi click / cuộn bằng robotgo như bước thường. Bền hơn ảnh mẫu (không phụ thuộc theme, zoom, hover), nhanh hơn (vài ms), đọc được trạng thái phần tử; AI (31) đọc DOM chính xác hơn nhìn ảnh chụp.
  - Cách làm: **DOM để tìm, robotgo để bấm** → event thật của hệ điều hành (`isTrusted = true`, khó bị coi là bot), dùng lại engine / xem trước / thanh mini. Dự phòng *chạy nền* (cửa sổ thu nhỏ / bị che): `Input.dispatchMouseEvent` của CDP – người dùng chọn "chuột thật / chạy nền"
  - **Ưu tiên chạy chính xác** (người dùng chốt 2026-10-18 – làm ngay sau 32, trước 20 F5b–F7 + 26). Làm theo rủi ro, phần khó nhất trước:
    - D0 → D1 → D2 → D3 → D4 (D0 chưa đạt sai số 0–2 px thì chưa làm UI)
    - D1: kết nối CDP + `webWait` / `webClick` (chữ / aria) + engine `Browser` + test; D2: cuộn tới phần tử + không bị che + dự phòng selector; D3: UI (Thử tìm, Chọn phần tử); D4: `webScroll`, `webType`, `webRead`, chạy nền, MCP
  - [x] **D0. Thử nghiệm đo toạ độ** (code thử, chưa vào model): trang HTML cục bộ, tìm nút qua CDP → đổi sang toạ độ màn hình bằng cách đo → robotgo bấm; ghi sai số (px) trên màn Retina + màn ngoài, zoom 125%, mở DevTools, có thanh bookmark, dời cửa sổ. Ghi kết quả vào file này
    - Code thử: `cmd/d0probe` (build tag `probe`, ngoài build thường): `go run -tags probe ./cmd/d0probe -scenario all`. Mở Chrome profile tạm, đo bằng 2 điểm `pointermove` (+1 điểm kiểm), bấm 6 nút 20×20 (4 góc, giữa, 1 nút phải cuộn), chấm sai số bằng `clientX/Y` của `pointerdown`. **Xoá khi D1 có test thật `AUTOCLICKER_REAL=1`** (sau khi chạy lại trên Windows + màn ngoài)
    - **Kết quả 2026-10-18** (MacBook Retina 1512×982 điểm, DPR 2, Chrome 155): **đạt 0–2 px** → làm tiếp D1
      | Tình huống | scale đo được | Lệch lớn nhất | Ghi chú |
      |---|---|---|---|
      | Cơ bản | 1.0000 | 0.5 điểm ≈ 1 px vật lý | dư điểm kiểm 0.00; đo mất ~0.5 s |
      | Thanh bookmark | 1.0000 | 1 px | offset Y 132 → 166 tự đúng |
      | Zoom 125% | 1.2500 | 0.71 điểm ≈ 1.4 px | DPR 2 → 2.5 |
      | Mở DevTools (gắn phải) | 1.0000 | 1.4 px | |
      | Nút phải cuộn (`scrollIntoView` center) | – | 1 px | client đổi, offset không đổi |
      | Chrome bị VS Code / Finder che | 1.0000 | 1 px | **`Page.bringToFront` đủ** trên macOS (không cần kích hoạt tiến trình) |
      | Dời cửa sổ, **dùng số đo cũ** | – | ~230 px / trượt | `screenX/Y` đổi → phát hiện được |
      | Dời cửa sổ, đo lại | 0.9997 | 1 px – trừ 1 nút **nằm dưới Dock** (y 960 / màn 982) → không nhận click | |
      | Màn ngoài | – | – | **chưa đo** (chỉ 1 màn) – xem checks.md |
      | `isTrusted` | – | luôn `true` | |
    - Rút ra cho D1 / D2:
      - Cổng debug **không cố định**: máy đã có tiến trình Chrome headless giữ 9333 → dùng `--remote-debugging-port=0` rồi đọc dòng đầu `<profile>/DevToolsActivePort`
      - Toạ độ chuột robotgo là số nguyên → công thức phải dùng **điểm đã làm tròn** (dùng điểm chưa làm tròn lệch 1 px ở mép dưới)
      - Đo lại khi `screenX/Y`, `innerWidth/Height` hoặc `devicePixelRatio` (zoom) đổi
      - `elementFromPoint` không thấy **cửa sổ hệ điều hành che trình duyệt** (Dock, thanh menu, app khác luôn-trên-cùng) → D2 thêm: điểm bấm phải nằm trong vùng nhìn thấy của màn hình (trừ Dock / thanh menu), không thì cuộn cho lọt; **thử hover trước khi bấm**: đưa chuột tới, đọc `pointermove` – trang phải nhận đúng `clientX/Y` và phần tử dưới chuột là đích, không thì đo lại / báo lỗi (thêm ~100 ms, bắt được cả số đo cũ lẫn bị che)
      - chromedp v0.21 đổi API sang generic (`chromedp.Run(ctx, Evaluate[T](…))`, `cdp.Call`, `remote.NewAllocator`) – ví dụ cũ trên mạng không chạy
      - Lần chạy đầu nút phải cuộn 1 lần không nhận click, sau đó 8 lần liền đều trúng – chưa rõ nguyên nhân; D2 *chờ ổn định rect* + thử hover sẽ bắt được
  - [ ] Kết nối Chrome / Edge / Brave qua **CDP** (`chromedp`): app mở trình duyệt với `--remote-debugging-port` + **profile riêng** trong `UserConfigDir/autoclicker/browser-profile` (từ Chrome 136 không bật debug được trên profile mặc định – lần đầu người dùng đăng nhập lại, sau đó profile tự nhớ); nối vào phiên đang mở nếu có. Safari / Firefox không hỗ trợ. Phương án sau (nếu cần profile đang dùng): extension Chrome nói chuyện với app qua WebSocket
  - [ ] Selector, ưu tiên theo độ bền: chữ hiển thị (`text=Xem lịch sử`) / `aria-label` / `role` → `id` / `data-*` → CSS / XPath; **tránh class sinh tự động** (React / Facebook `x1n2onr6…`); cho nhiều selector dự phòng (trang đổi ngôn ngữ); chọn tab theo URL / tiêu đề; phần tử trong iframe cộng toạ độ iframe
  - [ ] **Đổi toạ độ trang → màn hình bằng cách đo** (giống overlay, xem *Quyết định đã chốt*), không tính từ `screenX` + chiều cao thanh tab (sai khi mở DevTools, thanh bookmark, thanh "Chrome đang bị điều khiển", DPI 125%…): robotgo đưa chuột tới 2 điểm trong trang, nghe `mousemove` qua CDP lấy `clientX/Y` → suy ra tỷ lệ (zoom × DPI) + độ lệch; `screen = offset + client × scale`. Đo lại khi cửa sổ dời / đổi cỡ (`screenX/Y`, `outerWidth/Height` đổi) hoặc đổi zoom
  - [ ] Trước khi bấm: phần tử ngoài viewport → cuộn (robotgo 1–2 nấc giữa trang rồi **đọc lại rect** đến khi lọt viewport; không cần cuộn thật thì `scrollIntoView({block:"center"})`); chờ ổn định (rect giống nhau 2 lần cách ~100 ms); **không bị che**: `document.elementFromPoint(tâm)` phải là phần tử hoặc con của nó, không thì chờ / lỗi; cửa sổ trình duyệt phải ở trên cùng: CDP `Page.bringToFront` (không chờ 26); **đo lại rect ngay trước khi di chuột**; điểm bấm trong vùng nhìn thấy của màn hình (không dưới Dock / thanh menu – D0); **thử hover** (trang nhận đúng `pointermove` trên phần tử) rồi mới bấm; nhật ký ghi selector nào đã khớp
  - [ ] Model: `webWait` (chờ phần tử hiện / mất, timeout, ghim được), `webClick` (nút chuột, số click, vị trí trong phần tử: tâm / ±dx dy), `webScroll` (cuộn tới phần tử / theo nấc trong phần tử), `webType` (focus ô rồi gõ bằng bước gõ chữ sẵn có), `webRead` (đọc chữ / thuộc tính vào biến – nối với biến / điều kiện 28); hết giờ dùng `onTimeout` của 25b. Thêm vào `model.go`, `types.ts`, `steptypes.go` (mô tả + ví dụ JSON)
  - [ ] Engine: interface `Browser` (như `Input` / `Screen`), test bằng trình duyệt giả (DOM giả + toạ độ); test thật `AUTOCLICKER_REAL=1` mở trang HTML cục bộ
  - [ ] UI: inspector (selector + dự phòng, tab, chuột thật / chạy nền), **Thử tìm** (khung phần tử trên overlay, toạ độ, chữ đọc được, ms); nút **Chọn phần tử trên trang** (inject script tô viền phần tử dưới chuột, click → tự sinh selector theo thứ tự ưu tiên ở trên); trạng thái kết nối trình duyệt
  - [ ] Chuẩn bị cho MCP (31): `FindElement(selector) → {rect màn hình, chữ, hiện / bị che}` và `PageOutline()` (cây accessibility rút gọn cho AI) là hàm Go độc lập UI – Thử tìm, Xem trên màn hình, MCP dùng chung
  - [ ] Kịch bản mẫu (bước 33): mỗi loại `webWait` / `webClick` / `webScroll` / `webType` / `webRead` ≥ 1 mẫu trong `internal/samples/` chạy trên trang sân thử (vd form đăng nhập, feed, bảng dữ liệu bằng DOM thay vì OCR) – `TestEveryStepTypeHasSample` bắt buộc
  - [ ] Kiểm tra tay: trang có nút "Xem lịch sử" nằm dưới màn hình → bước chờ + click chạy đúng khi dời cửa sổ, zoom 125%, mở DevTools, có popup che
