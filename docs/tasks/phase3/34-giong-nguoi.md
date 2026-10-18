# 34. Thao tác giống người

- [ ] **34. Thao tác giống người – di chuột, gõ, click, cuộn tự nhiên** *(người dùng yêu cầu, 2026-10-18 – làm sau 20 F5b / F5c)* – Chuột hiện **nhảy thẳng** tới đích (`Move`) + lệch ngẫu nhiên ± px / ± ms (bước 19) – dấu hiệu bot rõ nhất. Bước này làm chuột / phím có nhịp giống người để giảm nguy cơ trang web đánh dấu bot.
  - **Chốt 2026-10-18:** **bật sẵn, mức vừa** cho kịch bản mới; tắt được khi cần nhanh tối đa. Cài đặt theo kịch bản, node ghi đè được (bật / tắt / mức).
  - Giới hạn (ghi vào mô tả cho người dùng): chỉ giảm nguy cơ, **không đảm bảo** – trang lớn (Shopee…) còn xét tần suất, tài khoản, dấu vết trình duyệt; tự động hoá có thể trái điều khoản trang.
  - [ ] H1. Logic thuần Go (`internal/engine/human/` hoặc `internal/humanize`): sinh đường di chuột từ A → B – đường cong ngẫu nhiên (WindMouse / Bezier với điểm điều khiển lệch), vận tốc tăng rồi giảm, thời gian theo khoảng cách + cỡ đích (Fitts), thỉnh thoảng đi lố rồi chỉnh lại, rung nhẹ vài px. Nhận `rand.Source` + mức (`off` / `low` / `medium` / `high`) → test bằng seed (không vượt màn hình, kết thúc đúng đích, thời gian trong khoảng)
  - [ ] H2. Engine: `move` đi theo đường H1 (bước nhỏ qua `device.Input.Move`, chờ giữa các bước, huỷ được khi Dừng / dừng khẩn cấp); giữ nút click 60–150 ms; điểm bấm rơi ngẫu nhiên quanh tâm theo cỡ đích (`w, h` của chờ ảnh / web 25d) thay vì đúng tâm; kéo thả đi theo đường cong
  - [ ] H3. Gõ chữ / phím: nhịp từng ký tự khác nhau, thỉnh thoảng ngừng; tuỳ chọn gõ sai rồi xoá (mức cao). Cuộn: từng nấc, nhịp không đều. Node web (25d D2 chuyển sang): cuộn tới phần tử bằng bánh xe chuột thật giữa trang + đọc lại khung đến khi lọt viewport, thay `scrollIntoView`
  - [ ] H4. Model + UI: `Script.Human` (mức), `Step.Human` ghi đè (trống = theo kịch bản); Cài đặt chạy có mục *Thao tác giống người*; inspector node chuột / phím có ô ghi đè. Xem trước (con trỏ giả) vẽ theo cùng đường – công thức đường dùng chung Go ↔ TS qua `testdata/` hoặc Go sinh đường gửi lên
  - [ ] H5. Mẫu thử: trang sân thử vẽ lại vệt chuột + đo nhịp gõ (so bật / tắt)
  - Chuẩn bị cho MCP (31): mức `human` có trong `StepTypeDocs` / `FlowGuide`
  - Kiểm tra tay: [checks.md](../checks.md) – vệt chuột trên sân thử cong, không nhảy; tắt thì như cũ; Dừng giữa lúc đang di chuột dừng ngay
