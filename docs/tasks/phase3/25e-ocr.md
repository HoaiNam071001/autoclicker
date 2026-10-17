# 25e. Tìm theo chữ (OCR)

- [x] **25e. Tìm theo chữ (OCR)** *(người dùng yêu cầu, 2026-10-17 – tìm ảnh chưa đủ chính xác)* – Bước Chờ ảnh / Click theo ảnh (tự tìm) chọn **Tìm theo: Ảnh mẫu / Chữ**; bảng Ctrl+K có *Chờ chữ*, *Click theo chữ*.
  - [x] Go: `ocr.go` (`TextReader`, `findText`: không phân biệt hoa / thường, bỏ khoảng trắng, tuỳ chọn bỏ dấu / cả dòng, chỗ thứ N theo thứ tự đọc; khung = hợp các từ chứa chuỗi), `ocr_darwin.go` (Vision, `vi-VT` + `en-US`, mức accurate; ảnh dẹt chèn nền cho tỷ lệ ≤ 3:1), `ocr_other.go` báo chưa hỗ trợ (Windows.Media.Ocr làm cùng 17 / 18). Model `findBy`, `findText` (nhận `{{biến}}`), `textExact`, `textNoAccent`, `occurrence`; hết giờ báo các dòng đọc được
  - [x] UI: inspector ô chữ + công tắc + chỗ thứ; *Thử tìm* đọc chữ, không thấy thì hiện các dòng đọc được (bấm để dùng); thẻ bước / thanh mini / xem trước ghi chữ
  - [x] Test `TestFindText`, `TestWaitTextThenClick`, `TestWaitTextTimeout`, `TestValidateTextStep`; thật: cả màn 1920×1080 ~0,9 s, thanh menu 1920×40 thấy "Bookmarks" ~0,6 s
  - [x] Tìm ảnh chính xác hơn: khác tỷ lệ (mẫu 1x, màn Retina) thì thu nhỏ ảnh chụp thay vì phóng to mẫu; mẫu ít tương phản (chữ xám nhạt) không còn bị bước tìm thô coi là một màu; 12 ứng viên. Thử thật màn MacBook: đúng 37/127 → 106/127 (16 chỗ giống hệt). Test `TestFindTemplateLowContrast`, `TestMatchScaledDownscalesScreen`
  - [x] Sửa (người dùng báo, 2026-10-17): thêm *Click theo chữ* mà thẻ / inspector / thanh mini ghi "Click theo ảnh" → `stepLabel(s)` (`types.ts`) đặt tên theo cách tìm (Chờ chữ / Click theo chữ), dùng ở mọi chỗ hiện tên một bước; inspector Click bỏ chữ "ảnh" ở nhãn khi tìm theo chữ
  - Kiểm tra tay: xem [checks.md](../checks.md)
