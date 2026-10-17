# Cần kiểm tra tay

Những việc **đã xong code** nhưng cần người dùng thử trong app thật (`wails dev`). Mới nhất ở trên.
Thử xong: tick `[x]` ở đây, ghi kết quả vào [log.md](log.md); lỗi thì báo để sửa (thêm mục con trong file của bước).

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
