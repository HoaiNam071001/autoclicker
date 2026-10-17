# 23. Cuộn + Mở URL / app

- [x] **23. Cuộn + Mở URL / app** – Loại bước đơn giản.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Model: `scroll` (`direction`, `amount` – ghim được, điểm tuỳ chọn trong `points` nên có neo); `open` (`target` URL / đường dẫn / tên app, `app` tuỳ chọn) – trường `omitempty`, file cũ không đổi, không tăng `version`
  - [x] Engine: `Input.Scroll` (`robotgo.ScrollDir` từng nấc, nghỉ 30 ms – dừng được giữa chừng); `open.go`: `openCommand` (hàm thuần, có test) – macOS `open [-a App] target` (tên app → `open -a`), Windows PowerShell `Start-Process` (không dùng `cmd /c start` vì `&` trong URL bị hiểu là nối lệnh), Linux `xdg-open`; tên miền tự thêm `https://`; biến trong URL tự mã hoá (`resolveInputs`); test với input giả + hàm mở giả
  - [x] UI: bật Cuộn / Mở URL trong bảng Ctrl+K; inspector: hướng + số nấc + cuộn tại chuột / điểm (bản đồ neo, chọn trên màn hình, xem riêng trên overlay); ô URL nhận `{{biến}}` + chọn trình duyệt; Xem trước có bong bóng cuộn / mở
  - [x] Chuẩn bị cho MCP (31): `steptypes.go` – `stepTypeDocs` (tên, mô tả, trường, trường ghim được, ví dụ JSON) cho mọi loại bước + `App.StepTypes()`; test kiểm tra mọi loại có mô tả và ví dụ qua được `checkStep`
  - [x] Thể hiện cuộn tới đâu *(người dùng yêu cầu, 2026-10-17)*: quy ước 1 nấc ≈ 100 px (`SCROLL_NOTCH_PX` / `scrollNotchPx`); inspector có `ScrollMap` (khung vùng / màn hình theo tỷ lệ, điểm cuộn hoặc chuột, mũi tên đúng quãng, dải "hiện thêm ≈ N px", % chiều cao); overlay vẽ mũi tên + nhãn từ marker; Xem trước: mũi tên hồng chạy dần + bong bóng; thẻ bước ghi ≈ px
  - [x] Kiểm tra tay macOS (người dùng xác nhận 2026-10-17); Windows (PowerShell `Start-Process`, cuộn theo nấc) kiểm tra khi có máy – ghi ở 17. Chi tiết đã thử: macOS (mở URL có biến tiếng Việt bằng Chrome / mặc định, mở app theo tên, cuộn tại chuột / điểm); Windows (PowerShell `Start-Process`) khi có máy
