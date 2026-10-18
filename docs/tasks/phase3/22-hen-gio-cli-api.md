# 22. Hẹn giờ + thông báo + CLI / API local

- [ ] **22. Hẹn giờ + thông báo + CLI / API local** – Chạy không cần ngồi canh.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [ ] Hẹn giờ: chạy lúc HH:mm (một lần / hằng ngày) với preset input đã chọn; chip giờ trong thư viện; app thu xuống khay
  - [ ] Thông báo: Telegram / Discord / webhook khi xong hoặc lỗi, kèm ảnh chụp màn hình
  - [ ] CLI: `autoclicker run script.json --preset "Gửi cho Lan" --set M1=Lan`
  - [ ] API local: `POST /run` (kèm values), `POST /stop`, `GET /status`, tiến độ qua SSE; bảo vệ bằng token
  - [ ] Chuẩn bị cho MCP (31): server HTTP chỉ nghe localhost, gắn thêm route được (`/mcp` sau này); kịch bản bị đổi từ ngoài (file / API) → event `script:changed` để frontend nạp lại, giữ hoàn tác; chỗ chặn "hỏi người dùng trước khi chạy" dùng chung cho hẹn giờ / API / AI
  - [ ] **Làm ngay ở bước này, không để tới 31** *(2026-10-18 – để sau thì phải sửa lại CLI / API)*: (1) **phiên chạy** tách khỏi Wails: bắt đầu trả `runId` ngay, trạng thái / dừng / chờ theo `runId` – App, CLI, API, MCP gọi cùng một chỗ (thay `runHandle` trong `runner.go`); (2) kịch bản có **revision** (tăng mỗi lần lưu / đổi từ ngoài, giữ trong bộ nhớ, không ghi file) gửi kèm `script:changed` – API ghi kèm revision cũ thì báo xung đột
