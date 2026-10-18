# 25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật

- [ ] **25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật** *(người dùng yêu cầu, 2026-10-17; tách từ 28)* – Tìm phần tử trên trang theo chữ / selector (vd nút "Xem lịch sử"), tính toạ độ màn hình của nó rồi click / cuộn bằng robotgo như bước thường. Bền hơn ảnh mẫu (không phụ thuộc theme, zoom, hover), nhanh hơn (vài ms), đọc được trạng thái phần tử; AI (31) đọc DOM chính xác hơn nhìn ảnh chụp.
  - Cách làm: **DOM để tìm, robotgo để bấm** → event thật của hệ điều hành (`isTrusted = true`, khó bị coi là bot), dùng lại engine / xem trước / thanh mini. Dự phòng *chạy nền* (cửa sổ thu nhỏ / bị che): `Input.dispatchMouseEvent` của CDP – người dùng chọn "chuột thật / chạy nền"
  - **Ưu tiên chạy chính xác** (người dùng chốt 2026-10-18 – làm ngay sau 32, trước 20 F5b–F7 + 26). Làm theo rủi ro, phần khó nhất trước:
    - D0 → D1 → D2 → D3 → D4 (D0, D1, D2 xong) (D0 chưa đạt sai số 0–2 px thì chưa làm UI)
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
  - [x] **D1. CDP + `webWait` / `webClick` + engine `Browser`** – người dùng chốt 2026-10-18:
    - Mở trang: **node Mở có tuỳ chọn `webBrowser`** (“Chrome của app – điều khiển được”): chưa mở thì app bật Chrome / Edge / Brave với profile riêng + cổng debug (`--remote-debugging-port=0`, đọc `DevToolsActivePort`), đang mở thì mở tab mới. Node web nối vào trình duyệt này; chưa mở → lỗi mã `browser`
    - Tab: mặc định **tab đang hiện** (tab vừa mở bằng node Mở, không thì tab dùng gần nhất); trường `tab` tuỳ chọn = URL / tiêu đề chứa chuỗi
    - Selector: **`selectors []string` có tiền tố**, dòng đầu ưu tiên, dòng sau dự phòng: `text=` (chữ hiển thị, chứa – không phân biệt hoa thường / khoảng trắng; `text="…"` = đúng cả chuỗi; không tiền tố = text), `aria=` (aria-label / aria-labelledby), `css=`, `xpath=`. Parse ở `model.ParseSelector` (Go dùng chung kiểm tra + chạy)
    - [x] Model + kiểm tra + steptypes + VarFields (`selectors`, `tab`)
    - [x] Package `internal/web`: `Browser` interface, Chrome thật (mở / nối, chọn tab, `find.js`, đo toạ độ như D0, đo lại khi `screenX/Y` / cỡ / zoom đổi), test phần thuần + test thật `AUTOCLICKER_REAL=1`
    - [x] Engine `webWait` (hiện / mất) + `webClick` (tìm tới timeout, cuộn vào giữa nếu ngoài viewport, chưa bấm khi bị che) + test bằng trình duyệt giả
    - [x] Frontend: loại node, nhãn, icon, palette, view, inspector tối thiểu (selector mỗi dòng một cách tìm, tab, chờ hiện / mất, nút chuột) – Thử tìm / Chọn phần tử để D3
    - [x] Kịch bản mẫu webWait / webClick trên sân thử
    - **Kết quả 2026-10-18** (MacBook Retina, Chrome): `AUTOCLICKER_REAL=1 AUTOCLICKER_SAMPLES=web-lich-su-don,web-bam-nut-chay` qua 2/2 – nút chạy chỗ trúng 5/5; chờ hết “Đang tải…” → Đồng ý popup cookie → tự cuộn tới “Xem lịch sử” nằm dưới màn hình → bấm trúng
    - Chi tiết kỹ thuật: tìm theo chữ quét text node trước (nhanh trên trang lớn), chỉ quét mọi phần tử khi chữ bị tách qua nhiều thẻ; phần tử khớp lấy phần tử sâu nhất; đang tải (tab mới) không tính là “đã mất”; mỗi lệnh CDP tối đa 5 s và dừng được bằng nút Dừng; nối tab bằng `WithDetachOnCancel` (đóng kết nối không đóng tab); trình duyệt bật tách khỏi app (`detach_*.go`) – app tắt thì Chrome vẫn mở
    - Còn cho D2: cuộn bằng chuột thật (giờ là `scrollIntoView`), chờ rect ổn định, thử hover trước khi bấm, vùng bị Dock / thanh menu che, iframe, `role=`; D3: Thử tìm, Chọn phần tử, trạng thái kết nối; `cmd/d0probe` xoá sau khi đo lại trên Windows + màn ngoài
  - [x] **D2. Trước khi bấm cho chắc + iframe + `role=`** (2026-10-18):
    - [x] Chờ khung đứng yên: sắp bấm thì đọc khung 2 lần cách 100 ms (cũng là đo lại ngay trước khi di chuột), lệch > 1 CSS px → `Found.Moving`, tìm lại (`web/aim.go` `steady`)
    - [x] Thử hover (`Browser.Hover`): engine đưa chuột tới điểm bấm đầu (tâm + dx / dy + lệch ngẫu nhiên), trang phải nhận `pointermove` đúng toạ độ (lệch > 3 px → bỏ số đo, đo lại) và phần tử dưới chuột là đích (chỉ đòi khi điểm nằm trong khung – dx / dy ra ngoài thì thôi); không nhận gì (Dock / thanh menu / cửa sổ khác che – `elementFromPoint` không thấy) → cuộn phần tử vào giữa rồi thử lại. Hết giờ báo lý do của lần hover cuối
    - [x] Vùng bị Dock / thanh menu che: bắt bằng thử hover ở trên (không cần đọc vùng làm việc của màn hình – không cần cgo riêng từng hệ điều hành)
    - [x] iframe **cùng nguồn** (kể cả lồng nhau, `srcdoc`): tìm, cuộn, xét bị che từng tầng frame, toạ độ cộng khung nội dung iframe (viền + padding); bộ nghe chuột gắn cả trong frame (đo toạ độ không hỏng khi điểm đo rơi vào iframe). iframe **khác nguồn** (vd khung thanh toán của bên thứ ba) cần target CDP riêng – để sau khi có trang thật cần
    - [x] Selector `role=button`, `role=link[name="Xem lịch sử"]` (kiểu Playwright; vai trò thẻ ngầm định + thuộc tính `role`; name = tên truy cập chứa chuỗi: aria → `<label>` → alt / title / placeholder → chữ). Go `ParseSelector` ↔ TS `parseSelector`, test chung `testdata/selectors.json`
    - [x] Mẫu `web-thanh-toan-iframe` (trang `web-frame.html`): iframe dưới màn hình, ô đồng ý `div role=checkbox`, nút Thanh toán trượt vào 1,5 s
    - **Kết quả 2026-10-18** (MacBook Retina, Chrome): `AUTOCLICKER_REAL=1` 3 mẫu web qua 3/3 (`web-bam-nut-chay` 5/5 trúng, `web-lich-su-don`, `web-thanh-toan-iframe` – bấm Thanh toán sau khi nút trượt xong)
    - Chuyển đi: **cuộn bằng bánh xe chuột thật** (giờ vẫn `scrollIntoView` – trang nhận sự kiện scroll như thường) → 34 H3 (cùng nhóm thao tác giống người)
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
