# 25. Chờ ảnh, click theo ảnh, chờ pixel

- [x] **25. Chờ ảnh, click theo ảnh, chờ pixel** – Cốt lõi để tương tác web / app (làm sớm).
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Model: `waitImage` (`template` **PNG base64 ngay trong file kịch bản** + `templateScale`, `threshold` %, `timeoutMs`, `intervalMs`, `searchArea` vùng / cả màn chính); `clickImage` (`imageSource` prev = kết quả bước chờ ảnh trước / search = tự tìm, `dx` / `dy`); `waitPixel` (điểm có neo, `color`, `tolerance`, timeout). Ghim được `timeoutMs`, `threshold`
  - [x] Engine (`imagestep.go`): chụp qua interface `Screen` (`capture.go`, robotgo), so khớp `match.go` – **Go thuần, không dùng gocv**: NCC trên ảnh xám, tìm thô trên ảnh thu nhỏ (giữ 6 ứng viên) rồi tìm mịn; mẫu một màu so độ chênh. Màn Retina 3024×1964: ~70 ms (benchmark), máy thật: chụp ~270 ms + tìm ~160 ms cả màn. Chỗ thấy nhớ cho click theo ảnh; tiến độ chờ `Progress.wait`; hết giờ → dừng kèm lỗi "không thấy ảnh mẫu sau 10 s (khớp nhất 46%, cần 90%)" (thử lại / bỏ qua: 25b); thời gian chờ không tính lúc tạm dừng
  - [x] Overlay tool `capture`: kéo khung, tối ngoài khung, nhãn cỡ thật, "Dùng ảnh này" (Go ẩn cửa sổ 0,25 s rồi chụp – `CaptureTemplate`) / "Khoanh lại" / Esc; Ctrl+K "Chờ ảnh" thêm bước rồi mở khoanh luôn; khung kết quả "Thử tìm" vẽ trên overlay
  - [x] Inspector: ảnh mẫu (cỡ theo toạ độ), Khoanh lại, **Thử tìm** (`FindImage` – thấy ở đâu, % khớp, ms, xem trên màn hình), độ khớp, chờ tối đa (ghim được), kiểm tra mỗi, tìm trong vùng / cả màn; click theo ảnh: nguồn "Kết quả bước N" / tự tìm, ΔX ΔY, nút chuột; chờ pixel: điểm, màu + **Lấy màu** (`PixelColorAt`), sai số, chờ tối đa; báo thiếu quyền Screen Recording + nút "Xin quyền" (`capture_darwin.go`, cgo `CGPreflightScreenCaptureAccess`)
  - [x] Bật 3 loại trong bảng Ctrl+K; Xem trước có bong bóng chờ ảnh / click theo ảnh / chờ màu
  - [x] Thanh mini: vòng tiến độ theo thời gian chờ, "đang tìm x / y s · khớp nhất n%"
  - [x] Retina / DPI: ảnh mẫu lưu kèm tỷ lệ lúc chụp, tìm trên màn tỷ lệ khác thì co giãn mẫu (test 2x ↔ 1x)
  - [x] Chuẩn bị cho MCP (31): `CaptureTemplate(rect)`, `FindImage`, `PixelColorAt`, `matchScaled` là hàm Go độc lập với overlay, tạo ảnh mẫu từ toạ độ vùng; `stepTypeDocs` có 3 loại mới kèm ví dụ PNG thật
  - ~~Giới hạn: chỉ tìm trên màn hình chính~~ → tìm được trên mọi màn (2026-10-17, xem 17); test thật `AUTOCLICKER_REAL=1 go test -run RealCapture`
  - Kiểm tra tay: xem [checks.md](../checks.md)
