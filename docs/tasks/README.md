# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** 20 – flow kiểu n8n ([phase3/20-flow.md](phase3/20-flow.md)): F0, F4a, F1, F2, F4b, F3 xong code – sang F5. Việc chờ kiểm tra tay: [checks.md](checks.md). 30 còn ánh xạ theo 4 điểm mốc. Bước 17 (phần Windows) và 18 hoãn: chưa có máy Windows
- **Tiếp theo:** 20: F5 → F6 → F7 → 26 → 25d (web qua DOM) … (bảng thứ tự đầu Phase 3; 25c đã gộp vào 20 F5); MCP cho AI (31) làm sau 22, các bước trước nhớ làm mục *Chuẩn bị cho MCP*. Khi có máy Windows: quay lại 17–18
- **Cập nhật lần cuối:** 2026-10-17 (tách docs/TASKS.md thành docs/tasks/)

## Các file

| File | Nội dung |
|---|---|
| [checks.md](checks.md) | **Cần kiểm tra tay** – việc đã xong code, chờ người dùng thử trong app |
| [decisions.md](decisions.md) | Quyết định đã chốt – không cần bàn lại |
| [overview.md](overview.md) | Stack, khởi động nhanh, data model |
| [log.md](log.md) | Nhật ký (mỗi dòng: ngày, kết quả kiểm tra) |
| [phase0.md](phase0.md) · [phase1-mvp.md](phase1-mvp.md) · [phase2.md](phase2.md) | Bước 1–19c (còn 17 phần Windows, 18) |
| `phase3/<số>-<tên>.md` | Mỗi bước Phase 3 một file – thứ tự ở bảng dưới |

**Quy ước:** mục con xong là tick `[x]` ngay; mục cha tick khi hết phần code; việc *Kiểm tra tay* của bước đã xong code ghi ở [checks.md](checks.md) (file bước chỉ để link). Việc mới phát sinh thêm thành mục con trong file của bước. Bước mới → file mới `phase3/<số>-<tên>.md` + một dòng trong bảng dưới.

## Phase 3 – Design v2 & tự động hoá

> Thiết kế: **`design/DESIGN.md`** (token, thành phần, icon, map artboard → file), `design/theme.css` (token Tailwind v4), artboard tham khảo `design/canvas/*.dc.html` + canvas https://claude.ai/artifact/5WdDyFZdB3WKxEXwkticnw. Mỗi việc UI dưới đây ghi artboard tương ứng trong ngoặc vuông.
> Số bước giữ nguyên để khớp nhật ký nên không theo thứ tự số; **làm theo thứ tự từ trên xuống**. Xong 26 là chạy được ví dụ Facebook (cuối phase).

**Thứ tự làm** (từ trên xuống, giống danh sách bên dưới):

| # | Bước | Vì sao ở vị trí này |
|---|---|---|
| 1 | [27. UI v2 cho màn đang có](phase3/27-ui-v2.md) | Người dùng muốn làm giao diện trước; nền cho mọi màn sau |
| 2 | [27b. Khung UI tính năng chưa có](phase3/27b-ui-sap-co.md) | Dựng hết UI theo design một lần cho đồng bộ; bước sau chỉ nối logic |
| 3 | [29. Neo 4 góc](phase3/29-neo-diem.md) | Sửa toạ độ lệch khi đổi vùng; đổi model `Point` nên làm trước khi có file lưu (21) |
| 4 | [21. Lưu / tải](phase3/21-luu-tai.md) | App chưa lưu được gì – tắt là mất kịch bản |
| 5 | [30. Sửa nhiều điểm](phase3/30-sua-nhieu-diem.md) | Dựa vào neo (29) + UI mới (27); có kịch bản lưu được (21) mới đáng đầu tư chỉnh sửa hàng loạt; kèm hoàn tác / làm lại |
| 6 | [21b. Input](phase3/21b-input.md) | Cần 21 (preset theo từng kịch bản); mở đường cho chuỗi web / app (23, 25) |
| 7–8 | [23](phase3/23-cuon-mo-url.md) → [25](phase3/25-cho-anh.md) | Cuộn / mở URL đơn giản, rồi chờ ảnh – cốt lõi tương tác web / app |
| 9 | [25e. Tìm theo chữ (OCR)](phase3/25e-ocr.md) | Xong 2026-10-17 – tìm ảnh chưa đủ chính xác |
| 10 | [25b. Khi thất bại](phase3/25b-khi-that-bai.md) | Xong 2026-10-17; 20 thay *bỏ qua / nhảy tới* bằng cạnh `fail`, giữ phần thử lại |
| 11 | [20. Flow kiểu n8n (F0–F7)](phase3/20-flow.md) | Người dùng cần luồng thành công / thất bại, if / else, switch (2026-10-17); đổi model lớn nên làm trước các loại bước mới (26, 25d); gồm luôn debug + lịch sử chạy (thay 25c) |
| 12 | [26. Focus cửa sổ](phase3/26-focus-cua-so.md) | Xong 26 là chạy được ví dụ Facebook |
| 13 | [25d. Web qua DOM](phase3/25d-web-dom.md) | Người dùng yêu cầu (2026-10-17): web bền hơn ảnh mẫu; cần focus (26), cổng `fail` (20); làm trước MCP vì AI dùng DOM tốt hơn ảnh chụp |
| 14 | [22. Hẹn giờ / CLI / API local](phase3/22-hen-gio-cli-api.md) | Chạy không cần ngồi canh; dựng server HTTP + token mà MCP dùng lại |
| 15–17 | [31a → 31b → 31c. MCP cho AI](phase3/31-mcp-ai.md) | Người dùng chọn làm sau (2026-10-17); các bước trước đã có mục *Chuẩn bị cho MCP* nên lúc này chủ yếu là bọc hàm Go sẵn có thành tool |
| 18–19 | [24](phase3/24-ghi-macro.md) → [28](phase3/28-mo-rong.md) | Ghi macro, mở rộng |
| – | [17 (Windows), 18](phase1-mvp.md) | Khi có máy Windows; riêng *Kiểm tra tay macOS* của 17 làm lúc nào cũng được |
