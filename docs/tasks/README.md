# Auto Clicker – Plan

Desktop app tự động click theo kịch bản: chọn vùng / toàn màn hình, chọn điểm, nối theo thứ tự, cấu hình thời gian và số lần click, kéo thả giữa 2 điểm.

## Trạng thái hiện tại

> Cập nhật mục này mỗi khi xong / bắt đầu một bước.

- **Đang làm:** **38** – UI v3 gọn ([phase3/38-ui-v3.md](phase3/38-ui-v3.md)): **xong code** (màu 7 nhóm node, thẻ node 168, khung app + nút Chạy tách đôi, inspector thu gọn + ô nhập có menu ⋯), còn kiểm tra tay. **37** – phần tử màn hình qua Accessibility ([phase3/37-phan-tu-man-hinh.md](phase3/37-phan-tu-man-hinh.md)): **E0** (đọc cây AX macOS – Chrome thường, Finder, Electron) + **E1** (tìm theo text= / aria= / role=, chỗ thứ N, lọc vùng), **E2** (Chờ / Click theo phần tử + điều kiện *thấy phần tử* ở engine; node **Đọc phần tử** có inspector; 3 mẫu thử trên `element.html`), **E3** (UI Tìm theo: Phần tử, Chọn trên màn, Thử tìm), **E4** (đọc chữ cạnh nhãn bằng OCR, lọc mẫu, lấy số, chờ giá trị đổi) xong code; còn E5 (Windows), chạy thật mẫu, chuẩn bị MCP. **35** – node nâng cao ([phase3/35-node-nang-cao.md](phase3/35-node-nang-cao.md)): **xong code** (N0 mô tả node, Chụp màn hình, Clipboard, HTTP, Tệp / thư mục, Âm lượng, Chờ đến giờ, Vẽ đường, mẫu thử), còn kiểm tra tay. **34** – thao tác giống người ([phase3/34-giong-nguoi.md](phase3/34-giong-nguoi.md)): **xong code** (H1–H5, Xem trước đi theo đường cong), còn kiểm tra tay. 20 – flow kiểu n8n ([phase3/20-flow.md](phase3/20-flow.md)): F0, F4a, F1, F2, F4b, F3, F5, F5b (toạ độ nhận biến), **F5c** (output mọi node + ô số nhận biến) xong code – 32 (tách Go) xong code. **33** (kịch bản mẫu để test) xong code – 31 mẫu chạy thật qua 31/31; **25d** (web qua DOM): D0 đo toạ độ 0–2 px trên Retina, **D0–D4 xong code** (node Chờ / Click / Gõ / Đọc / Cuộn phần tử web qua Chrome của app; chờ khung đứng yên + thử hover, iframe, `role=`; inspector có trạng thái kết nối, Chọn phần tử trên trang, Thử tìm; chạy nền theo từng node – Chrome bị che / thu nhỏ vẫn chạy; `PageOutline` cây rút gọn cho AI), còn kiểm tra tay. Việc chờ kiểm tra tay: [checks.md](checks.md). 30 còn ánh xạ theo 4 điểm mốc. Bước 17 (phần Windows) và 18 hoãn: chưa có máy Windows
- **Tiếp theo:** **20 F6** → F7 → 26 … → 36 (việc qua Chrome của app – cuối plan) (bảng thứ tự đầu Phase 3; 25c đã gộp vào 20 F5); MCP cho AI (31) làm sau 22, các bước trước nhớ làm mục *Chuẩn bị cho MCP*. Khi có máy Windows: quay lại 17–18
- **Cập nhật lần cuối:** 2026-10-19 (37 E4 – đọc cạnh nhãn bằng OCR + lọc mẫu)

## Tính năng của hệ thống

> ✅ có, chạy được · 🧪 xong code, chờ bạn kiểm tra tay ([checks.md](checks.md)) · 🟡 có một phần · ⬜ chưa làm. Bấm cột *Bước* để xem chi tiết.
> Cập nhật bảng này khi một tính năng đổi trạng thái (viết tay, khác bảng *Tổng quan* tự sinh ở dưới).

### Soạn kịch bản

| | Tính năng | Ghi chú | Bước |
|---|---|---|---|
| 🧪 | Canvas flow kiểu n8n | Kéo node, nối cổng ra → node kế, thả cổng ở chỗ trống để thêm node, tự xếp, zoom, minimap | [20 F4b](phase3/20-flow.md) |
| ✅ | Bảng thêm node Ctrl/Cmd+K | Tìm không dấu, phím mũi tên, nối vào node đang chọn | [20](phase3/20-flow.md) |
| ✅ | Sửa node | Tên, mã node (n2), cổng ra, output, bật / tắt, node bắt đầu; copy / paste, nhân bản, chọn nhiều, xoá | [20 F4b](phase3/20-flow.md) |
| ✅ | Hoàn tác / làm lại | Cmd/Ctrl+Z, Cmd/Ctrl+Shift+Z cho mọi thay đổi | [30](phase3/30-sua-nhieu-diem.md) |
| ✅ | Chọn vùng / toàn màn hình, đặt điểm | Overlay: click hoặc F6, kéo marker, chọn lại điểm, xem riêng một node | [10–13](phase1-mvp.md), [19c](phase2.md) |
| ✅ | Neo điểm theo cạnh vùng | Điểm bám góc / cạnh / giữa / % – đổi cỡ vùng không lệch | [29](phase3/29-neo-diem.md) |
| 🟡 | Sửa nhiều điểm cùng lúc | Di chuyển, co giãn, xoay, xiên, căn, lật, ánh xạ khi đổi vùng – **còn** ánh xạ theo 4 điểm mốc | [30](phase3/30-sua-nhieu-diem.md) |
| ✅ | Input của kịch bản | Ghim trường làm input, `{{biến}}` trong ô chữ, form Chạy, preset, nhớ giá trị lần trước | [21b](phase3/21b-input.md) |
| ✅ | Lưu / mở kịch bản | File `.json`, mở gần đây, nháp tự lưu, hỏi lưu khi đóng | [21](phase3/21-luu-tai.md) |
| ⬜ | Ghi macro | Ghi click / kéo / gõ phím thành node, F9 dừng | [24](phase3/24-ghi-macro.md) |

### Node hành động

| | Tính năng | Ghi chú | Bước |
|---|---|---|---|
| ✅ | Click / double click / giữ chuột / kéo thả | Nút trái / phải / giữa, số lần click, kéo từ điểm trước, B lệch so với A | [7](phase1-mvp.md), [13](phase1-mvp.md) |
| ✅ | Ngẫu nhiên hoá | ± ms thời gian chờ, ± px vị trí | [19](phase2.md) |
| ✅ | Gõ chữ, phím tắt | Tiếng Việt, tổ hợp phím, gõ lần lượt | [7](phase1-mvp.md) |
| ✅ | Chờ, cuộn, mở URL / app | Cuộn theo nấc tại điểm / tại chuột; mở bằng app chỉ định | [23](phase3/23-cuon-mo-url.md) |
| 🧪 | Chờ ảnh / click theo ảnh / chờ màu pixel | Khoanh ảnh mẫu, Thử tìm, khác tỷ lệ Retina vẫn tìm được, mọi màn hình | [25](phase3/25-cho-anh.md) |
| 🧪 | Chờ chữ / click theo chữ (OCR) | Đọc chữ trên màn hình, bỏ dấu, chỗ thứ N – **chỉ macOS**, Windows chưa có | [25e](phase3/25e-ocr.md) |
| ✅ | Thử lại khi thất bại | Thử lại N lần, nghỉ giữa các lần, hết lượt đi cổng Thất bại | [25b](phase3/25b-khi-that-bai.md), [20](phase3/20-flow.md) |
| 🧪 | Thao tác giống người | Chuột đi đường cong có tăng / giảm tốc, nhịp gõ / click / cuộn tự nhiên – bật sẵn mức vừa | [34](phase3/34-giong-nguoi.md) |
| 🧪 | Chụp màn hình, clipboard, HTTP | Đợt 1 node nâng cao: chụp màn / vùng ra file, đọc / ghi clipboard, gọi API | [35](phase3/35-node-nang-cao.md) |
| 🧪 | Tệp / thư mục, âm lượng, chờ đến giờ, vẽ đường | Đợt 2 node nâng cao (xoá vào Thùng rác); Bluetooth hoãn | [35](phase3/35-node-nang-cao.md) |
| ⬜ | Cửa sổ + Ứng dụng | Node Cửa sổ (focus, chờ mở / đóng, đặt vị trí + cỡ, thu nhỏ / phóng to / đóng, đọc) + node Ứng dụng (đang chạy? cổng yes / no, chờ chạy / thoát, thoát / buộc thoát) – cố định môi trường trước khi click | [26](phase3/26-cua-so-ung-dung.md) |
| 🧪 | Web qua DOM | **Chỉ trang cơ bản, không bắt đăng nhập / captcha** (2026-10-18; việc tiếp theo dời sang [36](phase3/36-chrome-app-sau.md)). Tìm phần tử trang bằng chữ / selector (Chrome), bấm bằng chuột thật – D1: Chờ / Click phần tử web, node Mở bằng Chrome của app; D2: chờ đứng yên + thử hover trước khi bấm, iframe cùng nguồn, `role=`; D3: Chọn phần tử trên trang, Thử tìm, trạng thái kết nối; D4a: Gõ vào ô / Đọc phần tử / Cuộn trang web (tới phần tử, trong khung, tới khi thấy); D4b: chạy nền theo từng node (CDP, không đụng chuột); D4c: `PageOutline` – cây phần tử rút gọn cho AI (MCP 31) | [25d](phase3/25d-web-dom.md) |
| 🧪 | Tìm / đọc phần tử màn hình (Accessibility) | **E0–E4 xong code (macOS)**: đọc cây, Chờ / Click theo phần tử (Tìm theo: Phần tử, Chọn trên màn, Thử tìm), node Đọc phần tử (+ đọc chữ cạnh nhãn bằng OCR, lọc mẫu, lấy số, chờ giá trị đổi), mẫu thử `element.html`. Chrome thường đã đăng nhập + app khác, không CDP: Chờ / Click theo phần tử, Đọc phần tử; dự phòng OCR cạnh nhãn + lọc mẫu – macOS trước, Windows cùng 17–18 | [37](phase3/37-phan-tu-man-hinh.md) |
| 🧪 | Kịch bản mẫu để test | Nút *Mẫu thử* ở Thư viện: mẫu đơn giản cho từng loại bước + nhiều kiểu web app, chạy trên sân thử cục bộ / vài web thật | [33](phase3/33-mau-thu.md) |

### Rẽ nhánh & luồng

| | Tính năng | Ghi chú | Bước |
|---|---|---|---|
| ✅ | Cổng Thành công / Thất bại, vòng lặp bằng cạnh | Thất bại không nối = dừng báo lỗi; chặn lặp vô hạn 10 000 node / vòng | [20 F1–F2](phase3/20-flow.md) |
| ⬜ | Mọi node có output, ô số nhận biến | Gõ chữ / phím / cuộn / chờ / mở cũng trả output; số lần click, thời gian chờ… nhận `{{biến}}` | [20 F5c](phase3/20-flow.md) |
| 🧪 | Toạ độ lấy từ biến | Click / kéo / cuộn / chờ pixel tại `{{n2.output.x}}` (+ lệch px) – vd click cạnh chữ vừa tìm thấy | [20 F5b](phase3/20-flow.md) |
| ✅ | Biến giữa các node | `{{n2.output.x}}`, `{{n2.status}}`, `{{last…}}`, `{{vars…}}`; gợi ý khi gõ `{{` | [20](phase3/20-flow.md) |
| 🧪 | Nếu | Thấy chữ / ảnh / màu, so giá trị, và / hoặc, đảo, chờ tới khi đúng | [20 F3](phase3/20-flow.md) |
| 🧪 | Switch theo giá trị / theo màn hình | Bằng / chứa / khớp mẫu; màn nào hiện trước đi cổng đó | [20 F3](phase3/20-flow.md) |
| 🧪 | Đếm, Kết thúc, Đặt biến | Giới hạn số lần lặp; kết thúc thành công / báo lỗi; gán / cộng biến | [20 F3](phase3/20-flow.md) |
| 🧪 | Thử node | Xét Nếu / Switch trên màn hình lúc này, khung chỗ thấy trên overlay | [20 F3](phase3/20-flow.md) |
| ⬜ | Gọi kịch bản khác như một node (subflow) | | [28](phase3/28-mo-rong.md) |

### Chạy

| | Tính năng | Ghi chú | Bước |
|---|---|---|---|
| ✅ | Chạy N vòng / vô hạn / theo thời gian | Đếm ngược, F8 chạy / dừng, tạm dừng / tiếp tục, thanh mini | [8](phase1-mvp.md), [16](phase1-mvp.md) |
| ✅ | Dừng khẩn cấp | Rê chuột vào góc màn hình | [9](phase1-mvp.md) |
| 🟡 | Xem trước bằng con trỏ giả | Có – **chưa** đi theo nhánh rẽ (làm ở F6) | [19b](phase2.md), [20 F6](phase3/20-flow.md) |
| 🧪 | Canvas khi chạy + nhật ký từng node | Đường đã đi, node xong / lỗi / ×n, bảng nhật ký (output, log, ảnh lúc lỗi), lời nhắn Kết thúc thành công | [20 F5](phase3/20-flow.md) |
| 🧪 | Lịch sử chạy, debug | Lưu 50 lần chạy gần nhất / kịch bản, mở lại từ thư viện; điểm dừng, từng node (F10), chạy từ node | [20 F5](phase3/20-flow.md) |
| ⬜ | Hẹn giờ, thông báo, CLI / API local | Chạy lúc HH:mm, báo Telegram / Discord, lệnh `autoclicker run` | [22](phase3/22-hen-gio-cli-api.md) |
| ⬜ | Chạy theo từng dòng CSV | | [28](phase3/28-mo-rong.md) |

### AI

| | Tính năng | Ghi chú | Bước |
|---|---|---|---|
| ✅ | Mô tả loại node + hướng dẫn flow cho AI | `StepTypes`, `FlowGuide`, kịch bản mẫu – nền cho MCP | [20 F1](phase3/20-flow.md) |
| ⬜ | AI tạo / sửa kịch bản qua MCP | Claude / ChatGPT / Cursor… ra lệnh bằng lời | [31a](phase3/31-mcp-ai.md) |
| ⬜ | AI nhìn màn hình, AI tự thử và chạy | Đọc chữ / phần tử trên màn trước, chụp ảnh khi cần; chạy có hỏi quyền, công tắc AI, chạy một phần flow / chạy tiếp từ node lỗi, chạy thử khô | [31b, 31c](phase3/31-mcp-ai.md) |
| ⬜ | AI bên trong flow | Node Hỏi AI (nhiều nhà cung cấp, cả AI chạy trên máy), Gọi tool MCP khác – làm sau cùng của 31 | [31d](phase3/31-mcp-ai.md) |

### Giao diện & nền tảng

| | Tính năng | Ghi chú | Bước |
|---|---|---|---|
| ✅ | Giao diện v2, theme sáng / tối | Kiểu Windows 11; UI dựng sẵn cho tính năng *Sắp có* | [27](phase3/27-ui-v2.md), [27b](phase3/27b-ui-sap-co.md) |
| 🧪 | Giao diện v3 – gọn | Màu theo 7 nhóm node, thẻ node nhỏ (168), inspector thu gọn có tóm tắt, nút Chạy tách đôi + menu; theme tối giữ Slate | [38](phase3/38-ui-v3.md) |
| 🟡 | Đa màn hình / DPI | macOS xong (chờ kiểm tra tay 2 màn); Windows chưa thử | [17](phase1-mvp.md) |
| ⬜ | Bản Windows (.exe) | Cần máy Windows | [18](phase1-mvp.md) |

## Tổng quan theo bước

> Bảng tự sinh – sau khi tick chạy `python3 scripts/tasks-index.py`. Bấm tên bước để mở file chi tiết.

<!-- index:start -->
**34 / 49 bước xong code** · 🧪 **[77 việc chờ bạn kiểm tra tay](checks.md)** · ✅ xong · 🧪 xong code, chờ kiểm tra tay · 🟡 đang làm · ⬜ chưa làm

### Phase 3 – Design v2 & tự động hoá (theo thứ tự làm)

| Trạng thái | Bước | Mục con xong | Kiểm tra tay |
|---|---|---|---|
| ✅ xong | [27. UI v2 cho các màn đang có + theme sáng / tối](phase3/27-ui-v2.md) | 18 / 18 |  |
| ✅ xong | [27b. Dựng sẵn UI cho tính năng chưa có](phase3/27b-ui-sap-co.md) | 13 / 13 |  |
| ✅ xong | [29. Neo điểm theo 4 góc vùng](phase3/29-neo-diem.md) | 10 / 10 |  |
| ✅ xong | [21. Lưu / tải kịch bản](phase3/21-luu-tai.md) | 8 / 8 |  |
| 🟡 đang làm | [30. Sửa nhiều điểm cùng lúc: di chuyển, co giãn, xoay, xiên](phase3/30-sua-nhieu-diem.md) | 9 / 10 | [1 việc](checks.md) |
| ✅ xong | [21b. Input của kịch bản](phase3/21b-input.md) | 13 / 13 |  |
| ✅ xong | [23. Cuộn + Mở URL / app](phase3/23-cuon-mo-url.md) | 6 / 6 |  |
| 🧪 chờ kiểm tra tay | [25. Chờ ảnh, click theo ảnh, chờ pixel](phase3/25-cho-anh.md) | 8 / 8 | [1 việc](checks.md) |
| 🧪 chờ kiểm tra tay | [25e. Tìm theo chữ (OCR)](phase3/25e-ocr.md) | 5 / 5 | [1 việc](checks.md) |
| ✅ xong | [25b. Khi thất bại](phase3/25b-khi-that-bai.md) | 4 / 4 |  |
| 🟡 đang làm | [20. Flow kiểu n8n – mỗi bước là một module, rẽ nhánh theo kết quả](phase3/20-flow.md) | 58 / 62 | [17 việc](checks.md) |
| 🧪 chờ kiểm tra tay | [32. Tách `package main` ở root thành package trong `internal/`](phase3/32-go-packages.md) | 6 / 6 | [3 việc](checks.md) |
| 🟡 đang làm | [33. Danh sách kịch bản mẫu để test](phase3/33-mau-thu.md) | 8 / 10 | [4 việc](checks.md) |
| 🧪 chờ kiểm tra tay | [25d. Web qua DOM – tìm phần tử bằng HTML, bấm bằng chuột thật](phase3/25d-web-dom.md) | 49 / 49 | [1 việc](checks.md) |
| 🧪 chờ kiểm tra tay | [34. Thao tác giống người – di chuột, gõ, click, cuộn tự nhiên](phase3/34-giong-nguoi.md) | 7 / 7 | [5 việc](checks.md) |
| 🧪 chờ kiểm tra tay | [35. Node nâng cao](phase3/35-node-nang-cao.md) | 9 / 9 | [9 việc](checks.md) |
| 🧪 chờ kiểm tra tay | [38. UI v3 – gọn: màu theo nhóm node, thẻ node nhỏ, inspector thu gọn](phase3/38-ui-v3.md) | 9 / 11 | [7 việc](checks.md) |
| 🟡 đang làm | [37. Tìm / đọc phần tử màn hình qua Accessibility](phase3/37-phan-tu-man-hinh.md) | 7 / 10 | [4 việc](checks.md) |
| ⬜ chưa làm | [26. Node Cửa sổ (`window`) + Ứng dụng (`app`)](phase3/26-cua-so-ung-dung.md) | 0 / 13 |  |
| ⬜ chưa làm | [22. Hẹn giờ + thông báo + CLI / API local](phase3/22-hen-gio-cli-api.md) | 0 / 6 |  |
| ⬜ chưa làm | [31a. Kết nối AI qua MCP – tạo / sửa kịch bản](phase3/31-mcp-ai.md) | 0 / 7 |  |
| ⬜ chưa làm | [31b. MCP – AI nhìn màn hình](phase3/31-mcp-ai.md) | 0 / 4 |  |
| ⬜ chưa làm | [31c. MCP – AI tự thử và chạy kịch bản](phase3/31-mcp-ai.md) | 0 / 8 |  |
| ⬜ chưa làm | [31d. AI bên trong flow](phase3/31-mcp-ai.md) | 0 / 4 |  |
| ⬜ chưa làm | [24. Ghi macro](phase3/24-ghi-macro.md) | – |  |
| ⬜ chưa làm | [28. Mở rộng](phase3/28-mo-rong.md) | 0 / 5 |  |
| ⬜ chưa làm | [36. Web qua Chrome của app – làm sau](phase3/36-chrome-app-sau.md) | 0 / 5 |  |

### Phase 0 – 2

| Trạng thái | Bước | Mục con xong | Kiểm tra tay |
|---|---|---|---|
| ✅ xong | [1. Cài môi trường](phase0.md) | – |  |
| ✅ xong | [2. Khởi tạo project](phase0.md) | – |  |
| ✅ xong | [3. Spike click](phase0.md) | 3 / 3 |  |
| ✅ xong | [4. Spike hotkey toàn cục](phase0.md) | 4 / 4 |  |
| ✅ xong | [5. Spike overlay](phase0.md) | 4 / 4 |  |
| ✅ xong | [6. Định nghĩa model](phase1-mvp.md) | – |  |
| ✅ xong | [7. Engine chạy bước](phase1-mvp.md) | 5 / 5 |  |
| ✅ xong | [8. Pause / Resume / Stop](phase1-mvp.md) | 5 / 5 |  |
| ✅ xong | [9. Dừng khẩn cấp](phase1-mvp.md) | 4 / 4 |  |
| ✅ xong | [10. Chọn vùng](phase1-mvp.md) | 4 / 4 |  |
| ✅ xong | [11. Chọn điểm](phase1-mvp.md) | 5 / 5 |  |
| ✅ xong | [12. Marker & đường nối](phase1-mvp.md) | 4 / 4 |  |
| ✅ xong | [13. Chọn điểm kéo thả](phase1-mvp.md) | 6 / 6 |  |
| ✅ xong | [13b. Dựng UI tĩnh từ bản thiết kế](phase1-mvp.md) | – |  |
| ✅ xong | [14. Danh sách bước](phase1-mvp.md) | 5 / 5 |  |
| ✅ xong | [15. Form cấu hình bước](phase1-mvp.md) | 5 / 5 |  |
| ✅ xong | [16. Cài đặt chạy](phase1-mvp.md) | 7 / 7 |  |
| 🟡 đang làm | [17. DPI & đa màn hình](phase1-mvp.md) | 5 / 6 | [2 việc](checks.md) |
| ⬜ chưa làm | [18. Build & thử máy sạch](phase1-mvp.md) | – |  |
| ✅ xong | [19. Ngẫu nhiên hoá](phase2.md) | 5 / 5 |  |
| ✅ xong | [19b. Xem trước kịch bản](phase2.md) | 4 / 4 |  |
| ✅ xong | [19c. Xem / sửa riêng một bước trên overlay](phase2.md) | 3 / 3 |  |
<!-- index:end -->

## Các file

| File | Nội dung |
|---|---|
| [checks.md](checks.md) | **Cần kiểm tra tay** – việc đã xong code, chờ người dùng thử trong app |
| [decisions.md](decisions.md) | Quyết định đã chốt – không cần bàn lại |
| [overview.md](overview.md) | Stack, cài môi trường máy mới |
| [log.md](log.md) | Nhật ký (mỗi dòng: ngày, kết quả kiểm tra) |
| [phase0.md](phase0.md) · [phase1-mvp.md](phase1-mvp.md) · [phase2.md](phase2.md) | Bước 1–19c (còn 17 phần Windows, 18) |
| `phase3/<số>-<tên>.md` | Mỗi bước Phase 3 một file – thứ tự ở bảng dưới |

**Quy ước:** mục con xong là tick `[x]` ngay; mục cha tick khi hết phần code; việc *Kiểm tra tay* của bước đã xong code ghi ở [checks.md](checks.md) (file bước chỉ để link). Việc mới phát sinh thêm thành mục con trong file của bước. Bước mới → file mới `phase3/<số>-<tên>.md` + một dòng trong bảng dưới.

## Phase 3 – Design v2 & tự động hoá

> Thiết kế: **`design/DESIGN.md`** (token, thành phần, icon, map artboard → file), `design/theme.css` (token Tailwind v4), artboard trên canvas https://claude.ai/artifact/5WdDyFZdB3WKxEXwkticnw. Mỗi việc UI dưới đây ghi artboard tương ứng trong ngoặc vuông.
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
| 11 | [20. Flow kiểu n8n (F0–F5, rồi F5b–F7 sau 25d)](phase3/20-flow.md) | Người dùng cần luồng thành công / thất bại, if / else, switch (2026-10-17); đổi model lớn nên làm trước các loại bước mới (26, 25d); gồm luôn debug + lịch sử chạy (thay 25c) |
| 11b | [32. Tách Go thành package](phase3/32-go-packages.md) | Người dùng chọn (2026-10-18): làm đầu tuần sau, trước F5b – 49 file `package main` ở root không có ranh giới; MCP / API (22, 31) cần dùng lại engine không qua Wails |
| 11b2 | [33. Kịch bản mẫu để test](phase3/33-mau-thu.md) | Người dùng yêu cầu (2026-10-18): làm trước 25d D1 – có mẫu cho mỗi loại bước + sân thử cục bộ để test tay nhanh; bước mới sau này đều thêm mẫu |
| 11c | [25d. Web qua DOM](phase3/25d-web-dom.md) | Người dùng ưu tiên (2026-10-18): web phải chạy chính xác – làm ngay sau 32, chia D0–D4 (đo toạ độ trước); focus bằng CDP `Page.bringToFront` nên không chờ 26. Người dùng yêu cầu (2026-10-17): web bền hơn ảnh mẫu; cần cổng `fail` (20); làm trước MCP vì AI dùng DOM tốt hơn ảnh chụp |
| 11d | [34. Thao tác giống người](phase3/34-giong-nguoi.md) | Người dùng chọn (2026-10-18): sau 20 F5b + F5c – chuột nhảy thẳng là dấu hiệu bot rõ nhất, web (25d) cần nhất |
| 11e | [35. Node nâng cao](phase3/35-node-nang-cao.md) | Người dùng yêu cầu (2026-10-18): sau 34, trước 26 – ít node nền (chụp màn hình, clipboard, HTTP trước) + kịch bản mẫu; node mới dùng output chuẩn F5c |
| 11f | [38. UI v3 – gọn](phase3/38-ui-v3.md) | Người dùng yêu cầu (2026-10-19): làm trước 37 E3 (E3 thêm UI vào inspector – dựng luôn theo kiểu mới); màu theo nhóm node, thẻ node 168, inspector thu gọn, nút Chạy tách đôi; theme tối giữ Slate |
| 11g | [37. Phần tử màn hình qua Accessibility](phase3/37-phan-tu-man-hinh.md) | Người dùng yêu cầu (2026-10-18): sau 35, trước 20 F6 – đọc giá trị web ở Chrome thường đã đăng nhập (không CDP, tránh captcha của 36) + app khác; đã thử chạy trên macOS; phần Windows UIA làm cùng 17–18 |
| 12 | [26. Cửa sổ + Ứng dụng](phase3/26-cua-so-ung-dung.md) | Xong 26 là chạy được ví dụ Facebook. Người dùng mở rộng (2026-10-19): `focusWindow` → node Cửa sổ + node Ứng dụng – cố định app / cửa sổ (vị trí, cỡ) trước khi thao tác |
| 14 | [22. Hẹn giờ / CLI / API local](phase3/22-hen-gio-cli-api.md) | Chạy không cần ngồi canh; dựng server HTTP + token mà MCP dùng lại |
| 15–17 | [31a → 31b → 31c → 31d. MCP cho AI](phase3/31-mcp-ai.md) | Người dùng chọn làm sau (2026-10-17); các bước trước đã có mục *Chuẩn bị cho MCP* nên lúc này chủ yếu là bọc hàm Go sẵn có thành tool. 31d (AI trong flow) chỉ làm khi 31a–c dùng ổn (2026-10-18) |
| 18–19 | [24](phase3/24-ghi-macro.md) → [28](phase3/28-mo-rong.md) | Ghi macro, mở rộng |
| 20 | [36. Web qua Chrome của app – làm sau](phase3/36-chrome-app-sau.md) | Người dùng chốt (2026-10-18): đăng nhập Google / Facebook trong Chrome của app bị captcha, chế độ điều khiển quá rủi ro cho tài khoản thật – dời mọi việc qua Chrome của app xuống cuối, tìm cách sau |
| – | [17 (Windows), 18](phase1-mvp.md) | Khi có máy Windows; riêng *Kiểm tra tay macOS* của 17 làm lúc nào cũng được |
