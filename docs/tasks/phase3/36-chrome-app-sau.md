# 36. Web qua Chrome của app – làm sau

- [ ] **36. Web qua Chrome của app – làm sau** *(người dùng chốt 2026-10-18 – dời mọi việc tương tác qua Chrome của app xuống cuối plan, tìm cách sau)* – Node web (25d) vẫn giữ và dùng được cho trang cơ bản, không bắt đăng nhập / captcha; trang lớn (Google, Facebook…) người dùng dùng Chrome thường + node màn hình (ảnh / chữ / chuột / phím).
  - Lý do: đăng nhập Gmail / Facebook trong Chrome của app (profile riêng, mở với cổng điều khiển `--remote-debugging-port`) bị captcha nhiều lần, sợ bị khoá tài khoản; chạy trình duyệt ở chế độ điều khiển (dev) bị coi là quá rủi ro cho tài khoản thật
  - Kết quả thử 2026-10-18 (macOS, Chrome 155): tắt Chrome của app → mở cùng profile **không** có cổng điều khiển → đăng nhập Gmail → tắt → mở lại chế độ điều khiển: Gmail vẫn vào thẳng hộp thư. Facebook đăng nhập ngay trong chế độ điều khiển thì hỏi captcha
  - [ ] Tìm cách đăng nhập / dùng trang lớn an toàn: nút *Đăng nhập tài khoản web* (mở profile không cổng điều khiển → người dùng đăng nhập, bấm Xong → mở lại chế độ điều khiển – đã thử được với Gmail); xem lại dấu vết CDP trang có thể phát hiện (`Runtime.enable`…); phương án extension người dùng thấy quá phức tạp (2026-10-18)
  - [ ] Kiểm tra tay 25d: trang có nút "Xem lịch sử" nằm dưới màn hình → bước chờ + click chạy đúng khi dời cửa sổ, zoom 125%, mở DevTools, có popup che *(chuyển từ 25d)* – cùng các mục 25d trong [checks.md](../checks.md) (nhóm *Tạm hoãn*)
  - [ ] Node web: cuộn tới phần tử bằng bánh xe chuột thật giữa trang + đọc lại khung đến khi lọt viewport, thay `scrollIntoView` *(chuyển từ 34 H3b)*
  - [ ] Mẫu thử web cho 35: vào trang cá nhân FB (Mở URL + Click phần tử web), săn sale Shopee (web DOM + 34) *(chuyển từ 35)*
  - [ ] MCP 31b `describe_screen`: thêm phần tử trang web khi là Chrome của app (`PageOutline` đã có ở 25d D4c) *(chuyển từ 31)*

## Lỗi riêng Windows (thử 2026-10-19, `TestRealSamples`, màn 1920×1080 scale 125 %)

- [ ] `web-cuon-trong-khung`: “đã lăn 17 lần, thấy CỘT 12 nhưng chưa vào khung nhìn” – lăn ngang trong khung không tới (robotgo Windows lăn ngang bằng `MOUSEEVENTF_HWHEEL`; xem hướng / Chrome có nhận không)
- [ ] `web-thanh-toan-iframe`: tick ô đồng ý trong iframe được, nút “Thanh toán” báo “nằm ngoài khung trang” sau 5 s – có thể do khung nhìn thấp (màn 864 px logic) mà phần tử trong iframe không được cuộn vào
