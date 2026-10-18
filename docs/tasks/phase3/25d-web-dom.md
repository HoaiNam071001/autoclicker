# 25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật

- [ ] **25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật** *(người dùng yêu cầu, 2026-10-17; tách từ 28)* – Tìm phần tử trên trang theo chữ / selector (vd nút "Xem lịch sử"), tính toạ độ màn hình của nó rồi click / cuộn bằng robotgo như bước thường. Bền hơn ảnh mẫu (không phụ thuộc theme, zoom, hover), nhanh hơn (vài ms), đọc được trạng thái phần tử; AI (31) đọc DOM chính xác hơn nhìn ảnh chụp.
  - Cách làm: **DOM để tìm, robotgo để bấm** → event thật của hệ điều hành (`isTrusted = true`, khó bị coi là bot), dùng lại engine / xem trước / thanh mini. Dự phòng *chạy nền* (cửa sổ thu nhỏ / bị che): `Input.dispatchMouseEvent` của CDP – người dùng chọn "chuột thật / chạy nền"
  - **Ưu tiên chạy chính xác** (người dùng chốt 2026-10-18 – làm ngay sau 32, trước 20 F5b–F7 + 26). Làm theo rủi ro, phần khó nhất trước:
    - D0 → D1 → D2 → D3 → D4 (D0–D3 + D4a xong) (D0 chưa đạt sai số 0–2 px thì chưa làm UI)
    - D1: kết nối CDP + `webWait` / `webClick` (chữ / aria) + engine `Browser` + test; D2: cuộn tới phần tử + không bị che + dự phòng selector; D3: UI (Thử tìm, Chọn phần tử); D4: `webScroll`, `webType`, `webRead`, chạy nền, MCP
  - Phần đã xong (D0–D3, D4a, yêu cầu gốc): [25d-web-dom-xong.md](25d-web-dom-xong.md)
  - [ ] **D4. `webType`, `webRead`, `webScroll`, chạy nền, MCP** – chia 3 phần (người dùng chốt 2026-10-18):
    - [x] **D4a. 3 node mới** `webType` / `webRead` / `webScroll` (2026-10-18) – chi tiết ở file xong
    - [ ] **D4b. Chạy nền** (CDP `Input.dispatchMouseEvent` / `insertText` khi cửa sổ bị che / thu nhỏ) – hỏi người dùng chọn theo node hay theo kịch bản trước khi làm
    - [ ] **D4c. `PageOutline()`** (cây accessibility rút gọn cho AI, 31) – hàm Go độc lập UI
  - [ ] Chuẩn bị cho MCP (31): `FindElement(selector) → {rect màn hình, chữ, hiện / bị che}` và `PageOutline()` (cây accessibility rút gọn cho AI) là hàm Go độc lập UI – Thử tìm, Xem trên màn hình, MCP dùng chung *(`FindElement` xong ở D3, `PageOutline` → D4c)*
  - [x] Kịch bản mẫu (bước 33): mỗi loại `webWait` / `webClick` / `webScroll` / `webType` / `webRead` ≥ 1 mẫu trong `internal/samples/` chạy trên trang sân thử (vd form đăng nhập, feed, bảng dữ liệu bằng DOM thay vì OCR) – `TestEveryStepTypeHasSample` bắt buộc
  - [ ] Kiểm tra tay: trang có nút "Xem lịch sử" nằm dưới màn hình → bước chờ + click chạy đúng khi dời cửa sổ, zoom 125%, mở DevTools, có popup che
