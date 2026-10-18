# Cần kiểm tra tay

Những việc **đã xong code** nhưng cần người dùng thử trong app thật (`wails dev`). Mới nhất ở trên.
Thử xong: tick `[x]` ở đây, ghi kết quả vào [log.md](log.md); lỗi thì báo để sửa (thêm mục con trong file của bước).

## 20 F6 – Thanh mini nổi trên app toàn màn hình · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Chạy từng node (F10) → sang Chrome / app khác **toàn màn hình**: thanh mini vẫn hiện, bấm được nút chạy tiếp / dừng; overlay (chọn điểm, Thử tìm) cũng hiện trên app toàn màn hình
- [ ] Trong lúc thanh mini / overlay: icon Dock của app tạm ẩn (đúng như thiết kế); về màn chính icon Dock + menu Edit (Cmd+C / V trong ô nhập) trở lại, Cmd+Tab thấy app
- [ ] Thanh mini không giành focus: đang chạy kịch bản gõ phím vào app khác vẫn gõ đúng chỗ

## 33 – Kịch bản mẫu để test · [phase3/33-mau-thu.md](phase3/33-mau-thu.md)

- [ ] `wails dev` → Thư viện → **Mẫu thử…**: lọc theo kiểu / loại bước, tìm không dấu; *Mở sân thử* mở danh sách trang trong trình duyệt
- [ ] *Dùng mẫu* khi đang có thay đổi chưa lưu → hỏi lưu; mở xong là kịch bản mới chưa có file, Lưu hỏi chỗ lưu
- [ ] Chạy vài mẫu từ app (Click nhiều lần, Kéo thả, Form đăng nhập nhập sai mật khẩu, Kanban, Quiz): trang sân thử hiện đúng kết quả (đã chạy tự động 31 / 31 bằng engine – còn thử qua nút Chạy, thanh mini, F8)
- [ ] Click / bấm theo chữ trong app thật khác (không phải sân thử) không chậm đi rõ (mỗi lần di chuột chờ thêm 40 ms)

## 25d D0 – Đo toạ độ web · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Cắm màn ngoài (khác DPI Retina) rồi chạy `go run -tags probe ./cmd/d0probe -scenario all` – tình huống `display2` lệch ≤ 2 px (chuột sẽ tự di chuyển + bấm ~1 phút, đừng động vào chuột)
- [ ] Khi có máy Windows: chạy lại trên màn DPI 125% (sửa `-chrome` trỏ tới chrome.exe; phím zoom là Ctrl) – dán kết quả vào file bước

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
