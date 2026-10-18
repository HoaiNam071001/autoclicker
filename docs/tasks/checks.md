# Cần kiểm tra tay

Những việc **đã xong code** nhưng cần người dùng thử trong app thật (`wails dev`). Mới nhất ở trên.
Thử xong: tick `[x]` ở đây, ghi kết quả vào [log.md](log.md); lỗi thì báo để sửa (thêm mục con trong file của bước).

## 34 – Thao tác giống người · [phase3/34-giong-nguoi.md](phase3/34-giong-nguoi.md)

- [ ] Mẫu thử → **Thao tác giống người – vệt chuột, nhịp gõ** → Chạy: vệt chuột trên trang cong, không nhảy; *Điểm chuột / lần di* hàng chục, *Giữ nút* 60–150 ms, nhịp gõ lệch. Cài đặt chạy → **Tắt** → chạy lại: 1–2 điểm mỗi lần di, cong 0 px
- [ ] Bấm Dừng / F8 / rê chuột vào góc khi chuột đang đi giữa đường → dừng ngay, không bấm thêm; chuột đi gần mép / góc màn không tự kích hoạt dừng khẩn cấp
- [ ] Kịch bản mới: Cài đặt chạy mặc định **Vừa**; inspector node Click / Gõ có ô *Thao tác giống người* (Theo kịch bản / Tắt / …), node web chạy nền thì không có; thẻ node ghi đè có chip – xem cả theme sáng + tối
- [ ] Kéo thả (kanban) và double click vẫn ăn khi bật giống người; Windows (khi có máy): click giữ nút + đường chuột đúng chỗ khi DPI 125 %

## 20 F5c – Output mọi node + ô số nhận biến · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Mẫu thử → **Ô số từ biến – số click, nhịp, thời gian chờ** → Chạy: ô Click trái = 5; lời kết thúc ghi thời gian click + thời gian đã chờ gần bằng nhau; form Chạy hỏi *Nhịp giữa các click*
- [ ] Inspector node Click: bấm **{x}** cạnh *Số click* → ô biến chiếm cả hàng, gõ `{{` gợi ý output node trước (cả `startedAt`, `durationMs`); bấm lại {x} với ô là `4` → về ô số = 4. Ô đang ghim input thì {x} mờ – xem cả theme sáng + tối
- [ ] Gõ sai ô biến (vd `{{n1.status}}`, `abc`) → Chạy báo lỗi rõ ở node đó; biến ra 0 cho Số click → node Thất bại “số click phải ≥ 1”

## 20 F5b – Toạ độ nhận biến · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Mẫu thử → **Kanban – kéo theo toạ độ từ biến** → Chạy: thẻ VIỆC 1 sang cột HOÀN TẤT (HOÀN TẤT = 1); thử với Chrome zoom 125% vẫn đúng
- [ ] Inspector node Click: Vị trí → **Biến** → ô X / Y điền sẵn số hiện tại; chọn *Điền chỗ thấy… → n1* ra `{{n1.output.x}}` / `{{n1.output.y}}`, sửa thành `+ 20`; marker trên overlay của node biến mất, thẻ ghi `X {{n1.output.x}} + 20`; bấm **Số** lại → về điểm cũ. Xem cả theme sáng + tối
- [ ] Gõ sai (vd `{{n1.output.text}}`, `{{n1.output.x}} +`) → Chạy báo lỗi rõ ở node đó trước khi chạy; Kéo thả B *Lệch so với A* + Biến → ô ghi ΔX / ΔY

## 20 F6 – Thanh mini nổi trên app toàn màn hình · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Chạy từng node (F10) → sang Chrome / app khác **toàn màn hình**: thanh mini vẫn hiện, bấm được nút chạy tiếp / dừng; overlay (chọn điểm, Thử tìm) cũng hiện trên app toàn màn hình
- [ ] Trong lúc thanh mini / overlay: icon Dock của app tạm ẩn (đúng như thiết kế); về màn chính icon Dock + menu Edit (Cmd+C / V trong ô nhập) trở lại, Cmd+Tab thấy app
- [ ] Thanh mini không giành focus: đang chạy kịch bản gõ phím vào app khác vẫn gõ đúng chỗ

## 33 – Kịch bản mẫu để test · [phase3/33-mau-thu.md](phase3/33-mau-thu.md)

- [ ] `wails dev` → Thư viện → **Mẫu thử…**: lọc theo kiểu / loại bước, tìm không dấu; *Mở sân thử* mở danh sách trang trong trình duyệt
- [ ] *Dùng mẫu* khi đang có thay đổi chưa lưu → hỏi lưu; mở xong là kịch bản mới chưa có file, Lưu hỏi chỗ lưu
- [ ] Chạy vài mẫu từ app (Click nhiều lần, Kéo thả, Form đăng nhập nhập sai mật khẩu, Kanban, Quiz): trang sân thử hiện đúng kết quả (đã chạy tự động 31 / 31 bằng engine – còn thử qua nút Chạy, thanh mini, F8)
- [ ] Click / bấm theo chữ trong app thật khác (không phải sân thử) không chậm đi rõ (mỗi lần di chuột chờ thêm 40 ms)

## 32 – Tách Go thành package · [phase3/32-go-packages.md](phase3/32-go-packages.md)

- [ ] `wails dev`: app mở được, nháp lần trước khôi phục, mở / lưu kịch bản, danh sách mở gần đây, đổi theme
- [ ] Chạy một kịch bản có click + chờ ảnh / chữ + node điều khiển: thanh mini, F8 dừng, chuột chạm góc dừng; lịch sử chạy có lần mới (kèm ảnh khi node lỗi)
- [ ] Thử tìm ảnh / chữ, Thử node, Lấy màu, chụp ảnh mẫu trên overlay; form Chạy nhớ giá trị input lần trước

## 20 F5 – Lúc chạy, nhật ký, lịch sử, debug · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Chạy ví dụ Facebook (`examples/thu-node-dieu-khien.json`): về màn chính thấy bảng Nhật ký, đường đã đi sáng, ×n trên node lặp; bấm dòng chọn node; node lỗi có ảnh màn hình
- [ ] Node Kết thúc thành công có lời nhắn → "Xong: …" ở thanh trạng thái
- [ ] Lịch sử chạy trong thư viện: chạy vài lần, mở lại lần cũ; mở kịch bản khác thì danh sách đổi theo
- [ ] Debug: đặt Điểm dừng ở một node → Chạy dừng trước node đó (thanh mini), F10 / nút chạy từng node, Tiếp tục; *Từ node này*; *Từng node*

## 20 F3 – Node điều khiển · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] ví dụ Facebook bằng Switch theo màn hình + Đếm + Kết thúc; Nếu thấy ảnh (khoanh ảnh mẫu cho điều kiện), màu (chọn điểm, Lấy màu); Thử node vẽ khung trên overlay – file thử `examples/thu-node-dieu-khien.json`
- [ ] Lời nhắn node Kết thúc thất bại hiện ở "Lỗi lần chạy trước" (thanh trạng thái), node đó được chọn trên canvas

## 20 F4b – Canvas trong app · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] soạn ví dụ Facebook bằng canvas (node hành động + cạnh `fail` quay lại), thả cổng ra chỗ trống thêm node, đặt điểm liên tiếp trên overlay, hoàn tác, chạy thật
- [ ] Thay cho kiểm tra 25b cũ: chờ chữ không thấy → thử lại 2 lần → cạnh `fail` sang node cuộn rồi quay lại chờ chữ

## 25e – Tìm theo chữ · [phase3/25e-ocr.md](phase3/25e-ocr.md)

- [ ] kịch bản Facebook đổi sang Chờ chữ / Click theo chữ

## 25 – Chờ ảnh / click theo ảnh / chờ pixel · [phase3/25-cho-anh.md](phase3/25-cho-anh.md)

- [ ] khoanh nút trên trang web → Thử tìm → chạy chờ ảnh + click theo ảnh khi cửa sổ dời chỗ; chờ pixel đổi màu; thanh mini khi đang chờ

## 30 – Sửa nhiều điểm · [phase3/30-sua-nhieu-diem.md](phase3/30-sua-nhieu-diem.md)

- [ ] Chọn 5 điểm, co giãn 150%, xoay 10°, xiên, lật → overlay và lần chạy thật khớp; hoàn tác / làm lại

## 17 – Đa màn hình (macOS) · [phase1-mvp.md](phase1-mvp.md)

- [ ] macOS 2 màn: đổi màn trên overlay, đặt điểm / chọn vùng ở màn phụ, Xong → main về đúng chỗ cũ
- [ ] macOS: dừng khẩn cấp ở góc màn phụ, góc tiếp giáp 2 màn không dừng nhầm; overlay hiện "tỷ lệ 1.00"

## Chờ máy Windows

- [ ] Bước 17: scale 100% / 125% / 150% và 2 màn khác scale; 23: mở URL / app (PowerShell `Start-Process`), cuộn theo nấc – làm cùng bước 18 ([phase1-mvp.md](phase1-mvp.md))

## Tạm hoãn – node web qua Chrome của app · [phase3/36-chrome-app-sau.md](phase3/36-chrome-app-sau.md)

Người dùng chốt 2026-10-18: việc qua Chrome của app làm sau cùng – thử khi làm 36.

### 25d D4b – Node web: chạy nền · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Tắt hẳn Chrome của app (đã mở từ trước thì chưa có cờ chạy nền) → Mẫu thử → **Web: chạy nền** → Chạy, ngay sau khi trang mở thì chuyển sang app khác / thu nhỏ Chrome: chuột không động, kết thúc thành công “XIN CHÀO DEMO, lăn … lần”
- [ ] Inspector node Click / Gõ / Cuộn / Chờ phần tử web: *Điều khiển* Chuột thật / Chạy nền + dòng giải thích đổi theo; node Đọc không có lựa chọn này; thẻ node hiện “· chạy nền” – xem cả theme sáng + tối
- [ ] Chạy nền khi Chrome zoom 125% và khi tab đích không phải tab đang hiện (tab khác đang mở trên cùng cửa sổ) → vẫn bấm đúng phần tử

### 25d D4a – Node web: gõ, đọc, cuộn · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Ctrl+K nhóm *Web (trình duyệt)* có Gõ vào ô web / Đọc phần tử web / Cuộn trang web; inspector từng loại (công tắc Xoá chữ cũ / Enter, Đọc: chữ / giá trị ô / thuộc tính + ô tên, Cuộn: 3 cách + hướng / số nấc / khung) – xem cả theme sáng + tối
- [ ] Mẫu thử → **Web: đăng nhập, đọc kết quả** → Chạy: gõ vào đúng ô, lần sai hiện chữ đỏ, lần đúng ra XIN CHÀO DEMO, kết thúc thành công
- [ ] Mẫu thử → **Web: cuộn feed tới bài đặc biệt** và **Web: cuộn trong khung + tới phần tử** → Chạy: chuột lăn thật (feed tải thêm), dải ngang lăn mà trang không chạy, dừng ở MỤC 30
- [ ] Node Đọc → **Thử tìm** hiện “đọc được “…”” + số phần tử khớp; node Gõ trỏ nhầm vào nút (không phải ô) → cổng Thất bại “ô nhập không nhận focus”, không gõ bậy

### 25d D3 – Node web: Chọn phần tử, Thử tìm, trạng thái · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Inspector node Chờ / Click phần tử web: dòng trạng thái “Chrome của app: chưa mở” → bấm Mở → Chrome riêng bật, trạng thái thành “đang mở, 1 tab”; xem cả theme sáng + tối
- [ ] Mở trang sân thử web trong Chrome của app → **Chọn phần tử trên trang** → Chrome lên trước, viền cam theo chuột, bấm nút “Đồng ý” → app tự lên lại, ô cách tìm điền 2–3 dòng (`text=Đồng ý`, `role=button[name=…]`…), trang **không** nhận click đó; thử Esc / Huỷ chọn
- [ ] Chọn ô đồng ý trong iframe (mẫu Web: thanh toán trong iframe) → ra `role=checkbox…`
- [ ] **Thử tìm**: lần đầu chuột đi vào trang đo toạ độ rồi app lên lại, báo “Thấy … tại X Y · ms” + chữ; “xem trên màn hình” → khung xanh đúng chỗ phần tử; phần tử bị popup cookie che → báo bị che (khung vàng); phần tử dưới màn hình → “ngoài khung nhìn, lúc click sẽ tự cuộn tới”

### 25d D2 – Node web: trước khi bấm, iframe, role · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Mẫu thử → **Web: thanh toán trong iframe** → Chạy: tick ô đồng ý, chờ nút Thanh toán trượt xong mới bấm, hiện ĐÃ THANH TOÁN, *Click thật* = 2
- [ ] Kéo cửa sổ Chrome của app xuống để nút “Xem lịch sử” / Thanh toán rơi vào chỗ **Dock che** rồi chạy mẫu → node tự cuộn phần tử lên giữa rồi bấm trúng (không bấm vào Dock)
- [ ] Trong lúc node webClick đang chờ: đặt một cửa sổ app khác (luôn trên cùng, vd thanh mini) đè lên nút → không bấm nhầm; hết giờ báo “chưa bấm được: …” rõ lý do

### 25d D1 – Node web (chờ / click phần tử) · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] `wails dev` → Mẫu thử → **Web: mở lịch sử đơn hàng** → Chạy: lần đầu app bật Chrome riêng (profile của app), chạy hết tới LỊCH SỬ ĐƠN HÀNG, ô *Click thật* tăng; chạy lại lần 2 → mở tab mới trong Chrome đó (không bật Chrome thứ hai)
- [ ] Trong lúc chạy mẫu: dời cửa sổ Chrome / zoom 125% (Cmd +) / mở DevTools rồi chạy lại – vẫn bấm trúng (đo lại toạ độ)
- [ ] Tắt Chrome của app rồi chạy một node web (không có node Mở trước) → lỗi rõ “Chrome của app chưa mở…” (mã `browser`), không treo
- [ ] Inspector: thêm node qua Ctrl+K nhóm *Web (trình duyệt)*; nhập nhiều dòng selector, `{{biến}}` trong selector; node Mở chọn “Chrome của app”
- [ ] Trang thật đăng nhập (vd Facebook) trong Chrome của app: đăng nhập một lần, lần sau profile còn nhớ

### 25d D0 – Đo toạ độ web · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Cắm màn ngoài (khác DPI Retina) rồi chạy `go run -tags probe ./cmd/d0probe -scenario all` – tình huống `display2` lệch ≤ 2 px (chuột sẽ tự di chuyển + bấm ~1 phút, đừng động vào chuột)
- [ ] Khi có máy Windows: chạy lại trên màn DPI 125% (sửa `-chrome` trỏ tới chrome.exe; phím zoom là Ctrl) – dán kết quả vào file bước
