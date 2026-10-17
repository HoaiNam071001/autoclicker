# Phase 2 – Tiện ích

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
