# Cần kiểm tra tay

Những việc **đã xong code** nhưng cần người dùng thử trong app thật (`wails dev`). Mới nhất ở trên.
Thử xong: tick `[x]` ở đây, ghi kết quả vào [log.md](log.md); lỗi thì báo để sửa (thêm mục con trong file của bước).

## 41 – nơi lưu + xoá lịch sử · [phase3/41-noi-luu-lich-su.md](phase3/41-noi-luu-lich-su.md)

- [ ] Cài đặt (bánh răng) → tab **Lưu trữ**: 3 ô hiện *Mặc định: …/Documents/Auto Clicker* và *…/Pictures/Auto Clicker*; *Chọn…* đổi thư mục → ô chữ đậm, có *Về mặc định*; *Mở* mở đúng thư mục trong Finder / Explorer. Tab API / AI (MCP) vẫn như cũ. Xem theme sáng + tối
- [ ] Kịch bản mới → Lưu: hộp thoại mở sẵn thư mục *Kịch bản* đã chọn; *Mở kịch bản* cũng mở ở đó. Kịch bản đã có file → Lưu ghi tại chỗ cũ
- [ ] AI (Claude) tạo kịch bản mới → file vào thư mục *Kịch bản AI tạo* (để trống thì theo *Kịch bản*)
- [ ] Node Chụp màn hình để trống *Thư mục lưu* → ảnh vào thư mục ảnh trong Cài đặt (mặc định giờ là `Pictures/Auto Clicker`); node đã đặt thư mục riêng thì giữ
- [ ] Chọn thư mục rồi xoá / đổi tên thư mục đó trong Finder → Lưu kịch bản mới báo lỗi *không thấy thư mục … Về mặc định*, không lưu lặng lẽ chỗ khác; node Chụp màn hình báo lỗi tương tự
- [ ] Windows có OneDrive chuyển Documents / Pictures → mặc định trỏ vào thư mục trong OneDrive
- [ ] Lịch sử chạy → *Xem tất cả*: *Xoá lần này* xoá lần đang chọn (chọn sang lần cũ hơn kế đó); *Xoá hết…* hỏi lại rồi xoá cả, thanh bên hết lịch sử; đang xem lần đó trên bảng Nhật ký thì bảng đóng
- [ ] Thư viện → Mở gần đây: rê chuột hiện ✕ bỏ một mục; *Xoá danh sách* bỏ hết – file kịch bản vẫn còn trên đĩa

## 40 – đọc chữ tiếng Việt bằng Tesseract (Windows) · [phase3/40-nhan-dang.md](phase3/40-nhan-dang.md)

- [ ] Mẫu thử *Click theo chữ tiếng Việt có dấu* (Windows, 100 % và 125 %) → Bấm đúng “Bản” = 5, Bấm nhầm = 0
- [ ] App thật / web tiếng Việt: Chờ chữ / Click theo chữ có dấu (vd “Đăng nhập”, “Lưu thay đổi”) – tìm đúng, khung click trúng chữ; chữ trắng trên nút màu cũng thấy
- [ ] Node Chờ / Click theo chữ mới: công tắc *Bỏ qua dấu* bật sẵn – tìm “he thong” thấy “hệ thống”; tắt đi thì “hê thông” không khớp “hệ thống”. Kịch bản cũ mở lên vẫn chạy (node trước đây không bật bỏ dấu giờ cũng bỏ dấu)
- [ ] Exe tải từ GitHub Actions chạy trên máy Windows sạch (chưa cài Tesseract): đọc chữ có dấu được ngay, thư mục `%LOCALAPPDATA%\autoclicker\tesseract-…` tự tạo
- [ ] Node Chờ chữ vùng cả màn: không thấy giật chuột / UI chậm rõ khi chờ lâu (mỗi lần đọc cả màn ~0,8 s)
- [ ] AI `describe_screen` trên Windows trả chữ có dấu đúng

## 40 – tìm ảnh: so viền, mặt nạ, kiểm tra màu · [phase3/40-nhan-dang.md](phase3/40-nhan-dang.md)

- [ ] Mẫu thử *Tìm ảnh nâng cao: so viền, màu, mặt nạ* (macOS Retina, zoom 100 %) → chạy qua: Icon theme tối = 2, Nút xanh = 2, Huy hiệu = 2, các ô Nhầm / Bấm trượt = 0. Đổi node 2 về *So theo: Độ sáng* → hết giờ không thấy; tắt *Kiểm tra màu* ở node 4 → có lúc bấm nhầm nút xám
- [ ] App thật: khoanh ảnh mẫu một nút / icon ở theme sáng (vd app có theme hệ thống), *So theo: Đường viền* → đổi hệ điều hành sang theme tối → Thử tìm vẫn thấy, độ khớp ≥ ~80 %
- [ ] Inspector Chờ ảnh: *So theo* + *Kiểm tra màu* (bật → *Lệch màu tối đa* 15 %); Thử tìm trên nút cùng hình khác màu → báo “Thấy hình giống … nhưng màu lệch X %”; thẻ node ghi “· so viền · màu ±15%”; điều kiện Nếu *thấy ảnh* có cùng hai ô – xem theme sáng + tối
- [ ] (Windows 125 %) chạy lại mẫu trên → vẫn qua (ảnh mẫu 2x thu về đúng tỷ lệ cùng mặt nạ)

## Sửa – cửa sổ mất sau khi chạy (Windows 2 màn) / chọn Space (macOS)

- [ ] Windows, 2 màn (màn chính không ở góc 0,0 – vd màn chính bên phải màn phụ): đặt cửa sổ app trên màn chính → Chạy → xong thì màn chính hiện lại đúng chỗ cũ; thử thêm cửa sổ đang phóng to, cửa sổ trên màn phụ, mở / đóng overlay
- [ ] macOS nhiều Space: app mở ở Space 1 → sang Space 2 bấm icon Dock → cửa sổ dời sang Space 2 (không nhảy về Space 1); chạy xong thanh mini → về màn chính vẫn ở Space đang đứng; overlay / thanh mini vẫn nổi trên app toàn màn hình

## Sửa – mở notepad.exe / gõ chữ khi bật bộ gõ tiếng Việt · [phase3/23-cuon-mo-url.md](phase3/23-cuon-mo-url.md)

- [ ] Node Mở target `notepad.exe` (không đường dẫn) → mở Notepad, không mở Chrome; `youtube.com/…`, `parroto.app` vẫn mở trình duyệt
- [ ] Bật UniKey (Telex) → node Gõ chữ *Cách nhập: Tự chọn* gõ đoạn tiếng Việt + tiếng Anh ("too", "dd") vào Notepad / Chrome → đúng nguyên văn; clipboard trước đó (chữ) được trả lại sau khi dán
- [ ] *Gõ phím* + UniKey bật → nhật ký node có cảnh báo bộ gõ; bật giống người → vẫn gõ từng phím; Gõ vào ô web chạy nền không hiện *Cách nhập* – xem inspector theme sáng + tối

## 31b2 – chọn tool MCP cho từng loại AI · [phase3/31-mcp-ai.md](phase3/31-mcp-ai.md)

- [ ] Bánh răng → MCP: có bảng *Tool AI được dùng* (mỗi tool một công tắc, nhãn nhóm); tắt `screenshot` ở *Mặc định* → AI nào chưa có danh sách riêng (kết nối lại) không thấy / không gọi được `screenshot`, bị báo `tool_disabled`
- [ ] Kết nối Claude Code và Cursor (hoặc 2 client khác): sau lần đầu tên từng AI hiện trong ô chọn; tắt `screenshot` riêng cho Cursor, để Claude Code bật → Claude Code vẫn chụp được, Cursor không; *Dùng mặc định* xoá danh sách riêng
- [ ] Khởi động lại app: các lựa chọn còn nguyên

## 31b – AI tự chạy, xem lỗi, sửa, chạy tiếp · [phase3/31-mcp-ai.md](phase3/31-mcp-ai.md)

- [ ] Node Click phần tử web với selector sai nhẹ (vd `text=Xem lich su` khi trang ghi “Xem lịch sử đơn”) → hết giờ, `error.detail` có `perSelector` (0 hiện) + `outline` liệt kê phần tử đúng kèm `sel`; AI sửa selector theo `sel` rồi `resume_run` qua
- [ ] AI chạy kịch bản dài (vd 5 node click có `delayAfter` 1 s) → giữa chừng rê chuột → app tạm dừng trước node kế, thanh mini “Bạn vừa đụng chuột – đã tạm dừng”, AI gọi `wait_run` thấy `pausedBy: touch`; bấm ▶ chạy tiếp được. Khi AI đang tự di chuột (giống người, đường cong) thì không báo nhầm; lượt bạn tự chạy thì rê chuột không tạm dừng
- [ ] AI chạy kịch bản → thanh mini có chip “● AI · Tự chạy / Hỏi” (chấm nhấp nháy); bấm chip → 4 nút mức; chọn *Chỉ soạn* → lượt AI dừng ngay, Cài đặt → MCP cũng đổi theo; chạy tay (nút Chạy) thì không có chip – xem cả theme sáng + tối
- [ ] Bảo AI chạy kịch bản 4–5 node với `shots: "each"` → mỗi node có `shot`, kết quả gửi kèm 3 ảnh mới nhất (đúng cảnh sau khi node chạy xong + chờ `delayAfter`), AI xin ảnh node đầu bằng `wait_run shots: [...]` được; thời gian chạy không bị nhân đôi `delayAfter`; Lịch sử chạy có các ảnh đó
- [ ] (Windows) Node Mở Notepad, `delayAfter` 2 s → output có `window` “Untitled - Notepad” (hoặc tiếng Việt), `newWindow` true, node vẫn chờ đúng ~2 s; mở URL khi Chrome đang mở sẵn → `newWindow` false, `window` = Chrome. Click / gõ / phím sau đó → output `window` / `windowApp` là cửa sổ đang ở trước
- [ ] Chi tiết lỗi gần đúng: bảo AI chạy kịch bản có node sai nhẹ – chờ chữ “Gưi” (thật là “Gửi”), phần tử `text=Mua ngey`, cửa sổ / app gõ sai tên (vd `chromium`) → `wait_run` trả `error.detail` có `closest` / `near` chỉ đúng chữ / phần tử / cửa sổ / app thật (kèm khung), AI tự sửa node rồi `resume_run` qua
- [ ] Bánh răng → *Cho AI điều khiển (MCP)*: có *AI được chạy kịch bản* (Chỉ soạn / Hỏi trước khi chạy / Tự chạy) + *Tối đa mỗi lượt* – xem cả theme sáng + tối
- [ ] Thoát hẳn app → trên Claude Code / Desktop bảo AI làm việc cần chạy (vd "mở Notepad, gõ chào bạn, lưu ra Desktop/test.txt") → app tự mở **không hiện màn chính**; mức *Hỏi*: hiện hộp "AI muốn chạy …" (Enter = cho phép lần này, Esc = từ chối); khi chạy chỉ thấy thanh mini, chạy xong vẫn giữ thanh mini; bấm Dừng trên thanh mini → mở màn chính
- [ ] *Từ chối* → AI báo bị từ chối, không gọi lại liên tục; *Cho phép đến khi đóng app* → lần sau cùng kịch bản không hỏi; mức *Chỉ soạn* → AI báo chỉ được soạn
- [ ] Cố ý cho một node lỗi (vd chờ chữ không có trên màn) → AI nhận ảnh màn hình lúc lỗi + `error.detail` (các dòng chữ đọc được, cửa sổ đang trước), tự sửa node, `resume_run` chạy tiếp từ node đó – các node trước **không** chạy lại (Notepad không mở lần nữa)
- [ ] AI chạy `stopAt` từng đoạn rồi chạy tiếp; `describe_screen` trả chữ đúng chỗ (toạ độ khớp khi click thử), `screenshot` thấy đúng màn – thử thêm màn 125 % / màn phụ
- [ ] Đang AI chạy mà chuyển công tắc MCP sang tắt / *Chỉ soạn* → lượt chạy dừng ngay; F8 / chuột chạm góc vẫn dừng; Lịch sử chạy ghi nguồn **AI**

## 31a – AI tạo / sửa kịch bản qua MCP · [phase3/31-mcp-ai.md](phase3/31-mcp-ai.md)

- [ ] Bánh răng → *Cho AI điều khiển (MCP)* đang bật → *Thêm vào Claude Desktop…* → hỏi đường dẫn → Ghi → thoát hẳn + mở lại Claude Desktop: thấy server `autoclicker` (14 tool). Claude Desktop bản Microsoft Store có thể dùng file cấu hình khác – báo lại nếu không thấy
- [ ] Trên Claude: "tạo kịch bản mở bài Tinh Hà Say 2 trên YouTube, tên bài là input" → file mới trong Documents/Auto Clicker, *Mở gần đây* có chip **AI**; AI gọi `open_in_app` → app hiện lên, mở kịch bản; Xem trước đúng thứ tự; bấm Chạy phát nhạc
- [ ] "đổi thành bài khác / thêm bước…" → AI sửa đúng kịch bản (app nạp lại, Ctrl+Z lấy lại bản trước); đang sửa dở kịch bản đó (chưa lưu) → AI ghi ra bản sao `… (AI).json`, bản của bạn không đổi
- [ ] Đang có thay đổi chưa lưu ở kịch bản khác → AI gọi mở thì app không đổi kịch bản, AI báo lại
- [ ] Claude Code: chép lệnh `claude mcp add …` trong Cài đặt → `/mcp` thấy autoclicker; app chưa mở → tự mở khi Claude gọi
- [ ] Tắt công tắc MCP → AI gọi tool bị từ chối (báo bật trong Cài đặt); bật lại dùng được ngay. Cài đặt hiện 5 lệnh AI gần nhất (✓ / ✗)
- [ ] Bật *API cho chương trình khác* → cấu hình *HTTP (client nối URL)* dùng được với client nối thẳng `http://127.0.0.1:<cổng>/mcp`

## 22 – Hẹn giờ / CLI / API local · [phase3/22-hen-gio-cli-api.md](phase3/22-hen-gio-cli-api.md)

- [ ] Chạy tay như cũ (nút Chạy, F8 dừng, chạy từng node / chạy riêng node, xác nhận thao tác nguy hiểm, lịch sử chạy) – phần chạy đã chuyển sang `internal/session`
- [ ] CLI (PowerShell, app đang mở): `build\bin\autoclicker.exe run <file.json>` → app thu thành thanh mini, terminal in từng node ✓ / ✗ + *Xong sau … s*, `$LASTEXITCODE` = 0; Ctrl+C trong terminal → kịch bản dừng; `--preset "Tên bộ"` / `--set M1=Lan` dùng đúng giá trị; preset sai tên → báo các tên có, mã 1
- [ ] CLI khi app chưa mở: tự mở app rồi chạy; `run` không có file / cờ lạ → in cách dùng, mã 2; `help` in cách dùng
- [ ] Bánh răng trên thanh tiêu đề → bật *API cho chương trình khác*: `curl -X POST http://127.0.0.1:47800/run -H "Authorization: Bearer <token>" -d "{\"path\":\"C:/…/a.json\"}"` chạy được; sai token → 401; tắt công tắc → không nối được; đổi cổng / tạo token mới có tác dụng ngay
- [ ] Hẹn giờ: lưu kịch bản thành file → *N vòng* → Hẹn giờ chạy → Đặt giờ… (chưa lưu file thì nút mờ, có gợi ý); đặt giờ sau 1–2 phút, *Thu nhỏ cửa sổ* bật → app thu nhỏ, tới giờ tự hiện thanh mini và chạy, xong thu nhỏ lại; thư viện hiện chip giờ cạnh tên; *Một lần* chạy xong thì chip mất
- [ ] Thông báo hẹn giờ: Telegram (`botToken@chatId`), Discord webhook, webhook (vd webhook.site) nhận chữ + ảnh màn hình; *Gửi khi: Lỗi* thì lần chạy xong không gửi; tới giờ mà đang chạy kịch bản khác → thông báo lỗi + dòng báo trong app
- [ ] Kịch bản đang mở bị ghi qua API `PUT /script` (kèm `revision` lấy từ `GET /script`) → app nạp lại, có thông báo, Ctrl+Z lấy lại bản trước; gửi revision cũ → 409 conflict

## 26 W3b – Overlay vẽ khung đích node Cửa sổ · [phase3/26-cua-so-ung-dung.md](phase3/26-cua-so-ung-dung.md)

- [ ] Mở mẫu **Cửa sổ – đặt 1280 × 800 bên trái rồi chụp** → mở overlay: thấy khung nét đứt 1280 × 800 ở góc trái trên, nhãn `n… · chrome → 1280 × 800 tại 0, 0 · màn chứa cửa sổ lúc chạy`; chọn node đó → viền liền. Đổi Đặt ở Phải / Giữa / Đầy / Toạ độ → khung đổi theo, Đầy không đè taskbar
- [ ] Xem trước mẫu đó: khung chỉ hiện lúc tới node Cửa sổ (kèm bong bóng *không đổi cửa sổ thật*), cửa sổ thật không bị dời
- [ ] 2 màn hình / màn 125 %: chọn Màn hình 2 → khung vẽ trên màn 2 (mở overlay sang màn 2), khớp chỗ cửa sổ nằm sau khi Chạy

## 26 A1–A3 – Node Ứng dụng + xác nhận thao tác nguy hiểm · [phase3/26-cua-so-ung-dung.md](phase3/26-cua-so-ung-dung.md)

- [ ] Mẫu thử **Ứng dụng – mở Notepad, gõ chữ rồi thoát** (đóng Notepad trước): lần đầu hiện hộp *Kịch bản này sẽ thoát app / đóng cửa sổ* liệt kê n6; Cho phép → Notepad mở, có dòng chữ rồi tự tắt; chạy lại không hỏi nữa. Notepad đang mở sẵn → đi nhánh Đang chạy, dừng báo lỗi, không đụng Notepad
- [ ] Thêm node Ứng dụng (Ctrl+K → Ứng dụng): thẻ có cổng *Đang chạy* / *Không chạy*; đổi sang Thoát → cổng Thành công / Thất bại (cạnh cũ bỏ); ô App hiện ● Đang chạy · n cửa sổ / ○ Không chạy; *Chọn từ app đang chạy…* liệt kê app (Chrome một dòng), gõ lọc không dấu
- [ ] Thoát app đang hỏi lưu (Notepad cũ / Word có chữ chưa lưu) không bật *Quá giờ thì buộc thoát* → hết giờ cổng Thất bại, app vẫn mở hộp hỏi lưu; bật → app bị tắt
- [ ] Kịch bản có node Tệp xoá file (đã đồng ý trước đây): hỏi lại một lần (file đồng ý đổi tên `danger-consent.json`)
- [ ] App ghi `explorer.exe` + Thoát → báo lỗi ngay khi bấm Chạy, không chạy

## 26 W1–W3a – Node Cửa sổ (Windows) · [phase3/26-cua-so-ung-dung.md](phase3/26-cua-so-ung-dung.md)

- [ ] Mẫu thử **Cửa sổ – đặt 1280 × 800 bên trái rồi chụp**: cửa sổ trình duyệt nhảy về góc trái 1280 × 800; thu nhỏ / kéo sang chỗ khác / phóng to rồi chạy lại vẫn ra như vậy; lời kết thúc ghi khung + cỡ ảnh khớp
- [ ] Thêm node Cửa sổ (Ctrl+K → Cửa sổ): *Chọn từ cửa sổ đang mở…* liệt kê đúng cửa sổ (không có cửa sổ autoclicker), gõ lọc không dấu, chọn → điền App + Tiêu đề; *Lấy khung hiện tại* điền đúng X, Y, W, H
- [ ] Thao tác Đưa lên trước (cửa sổ đang thu nhỏ / bị che), Thu nhỏ, Phóng to, Đóng (Notepad có chữ chưa lưu → Notepad hỏi lưu, node không tự bấm), Chờ cửa sổ đóng; App / tiêu đề không khớp → hết giờ đi cổng Thất bại
- [ ] Màn 125 % / 150 % và 2 màn hình: Đặt ở Trái / Đầy trên màn 2 – khung đúng vùng làm việc (không đè taskbar)

## 20 F6 / F7 – Xem trước theo nhánh, thanh mini, chạy riêng node · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Chọn một node (vd Click n3) → ▾ cạnh nút Chạy → **Chạy riêng n3**: chỉ node đó chạy một lần (click thật), thẻ node tô kết quả, Lịch sử chạy có lần chạy chỉ một node; node thất bại (vd Chờ chữ không thấy) → báo lỗi dù có nối nhánh Thất bại; F8 dừng được

- [ ] Mẫu có rẽ nhánh (vd ví dụ Facebook: Chờ chữ → Thất bại → Đếm → Cuộn): overlay → Xem trước → **Nhánh** chọn Thất bại cho node Chờ → Chạy thử: con trỏ đi vòng cuộn 1 lần rồi theo Hết lượt; dòng gợi ý ghi đường đã đi; *Về nhánh chính* đi lại nhánh Thành công; đóng / mở lại overlay vẫn nhớ lựa chọn; Esc khi menu mở chỉ đóng menu
- [ ] Chạy thật kịch bản có Đếm + nhánh Thất bại: thanh mini ghi "n2 · … · Thất bại → n4 · Đếm", dòng "lượt n4 2 / 10" tăng dần, hết lượt thì mất

## 37 E5 – Phần tử màn hình Windows · [phase3/37-phan-tu-man-hinh.md](phase3/37-phan-tu-man-hinh.md)

- [ ] Windows, Chrome thường đã đăng nhập (Gmail / Facebook): node Đọc phần tử / Click theo phần tử đọc đúng chữ, bấm đúng nút; App ghi `chrome` hoặc `Google Chrome` đều được; Chọn trên màn sinh dòng tìm đúng phần tử
- [ ] Windows: app khác (Notepad, File Explorer, VS Code) – Thử tìm thấy nút / ô; Chrome thu nhỏ → báo “không có cửa sổ đang hiện”; app chạy quyền admin (vd Task Manager) → báo không đọc được
- [ ] Windows: Đọc phần tử *Lấy số* trên giá có chữ ₫ đọc bằng OCR (số 0 bị đọc thành chữ O vẫn ra số đúng)

## 25e – OCR Windows · [phase3/25e-ocr.md](phase3/25e-ocr.md)

- [ ] Windows: node Chờ chữ / Click theo chữ trên trang tiếng Việt (vd Mẫu thử → **Click theo chữ**, **Bảng dữ liệu**): thấy đúng chữ dù OCR đọc sai dấu; *Thử tìm* không thấy thì danh sách “đọc được” hiện chữ không dấu / sai dấu (bình thường trên Windows). Chữ ngắn < 5 ký tự (vd “HUẾ”) có thể không thấy – dùng chữ dài hơn hoặc ảnh mẫu
- [ ] Windows: Đọc phần tử → đọc chữ cạnh nhãn bằng OCR (37 E4) với nhãn tiếng Việt

## 17 – DPI Windows · [phase1-mvp.md](phase1-mvp.md)

- [ ] Windows 125 % (sửa 2026-10-19 – ảnh mẫu co giãn theo DPI màn): chụp lại ảnh mẫu đã chụp trên Windows trước bản này; chụp ảnh mẫu một nút ở 125 % → đổi scale sang 100 % (Settings → Display) → Thử tìm / chạy vẫn thấy; ảnh mẫu trong inspector hiện đúng cỡ nhìn trên màn. Có 2 màn khác scale: chụp ở màn này, tìm ở màn kia
- [ ] Windows 150 % / 100 %: `AUTOCLICKER_REAL=1 go test ./internal/device/ -run Real -v` qua (chuột đọc lại đúng, DPI in ra đúng scale)

## 18a – Build exe · [phase1-mvp.md](phase1-mvp.md)

- [ ] Windows cửa sổ trong suốt (sửa 2026-10-19): mở overlay → nhìn xuyên thấy màn hình thật (nền tối mờ); Xem trước với *Nền: Không* và *Nền: Trang web* đều thấy màn hình phía sau; màn chính vẫn nền đặc, kéo / đổi cỡ / phóng to cửa sổ như cũ; thanh mini trong suốt quanh viền

- [ ] Windows scale 125 % (sửa 2026-10-19 – chuột qua Win32, pixel thật): tải exe mới nhất ở Actions, xoá điểm cũ rồi ghi lại điểm (F6 / chọn trên overlay) → chạy: click trúng đúng chỗ, cả khi cửa sổ đích là app khác (Chrome, Notepad, Explorer); khung vùng / marker trên overlay nằm đúng chỗ; chờ ảnh / Thử tìm thấy đúng; thử thêm 150 % và 2 màn khác scale nếu có. Kịch bản ghi trên bản cũ ở 125 % có thể lệch – ghi lại điểm
- [ ] GitHub → tab **Actions** → *Build Windows* → lần chạy mới nhất xanh → tải artifact `autoclicker-windows-…` (zip), giải nén, mở `autoclicker.exe` trên Windows (SmartScreen: More info → Run anyway): app mở, tạo kịch bản click 2 điểm chạy đúng, Lưu / mở lại được. Ghi lại lỗi gặp (thiếu WebView2, click lệch khi scale 125 %…)

## 38 – UI v3 – gọn · [phase3/38-ui-v3.md](phase3/38-ui-v3.md)

- [ ] V4 Bảng thêm node (Ctrl+K) và hộp Mẫu thử: icon tô đúng màu nhóm, đọc rõ ở theme tối; lần chạy thất bại – canvas cuộn tới node lỗi, dòng "Flow · N node" không đè chữ lên thẻ
- [ ] V3b3 Inspector mọi loại node (`wails dev`): rê ⓘ cạnh nhãn hiện giải thích, không bị cắt mép; công tắc (Cả dòng, Bỏ qua dấu, Ghi nối tiếp, Vào clipboard, Giờ chuẩn, Khép kín…) bật / tắt đúng; báo lỗi đỏ vẫn hiện khi thiếu dữ liệu; theme tối đọc rõ
- [ ] V3b2 Vị trí (`wails dev`): node Click – ⋯ cạnh X / Y → Lấy toạ độ từ biến, *Dùng số cố định* trở lại; hàng Sơ đồ mở / đóng, bấm chấm đổi neo giữ nguyên chỗ trên màn hình, px / % đổi đúng; chữ neo ở đầu khối đổi theo
- [ ] V3b1 Ô nhập (`wails dev`): rê vào ô số hiện ⋯; Dùng biến… → ô biểu thức, ✕ trở lại số; Đặt làm input khi chạy → viền nét đứt, chip key, form Chạy hỏi giá trị; dùng chung input có sẵn; Về mặc định đúng theo loại node (Chờ 1000 ms, Mở 2000 ms…); − / + ở Số click; menu ⋯ gần đáy inspector mở lên trên
- [ ] V3a Inspector (`wails dev`): mở / đóng các khối, chọn node cùng loại khác vẫn giữ khối đang mở; tóm tắt đúng khi đóng (Thời gian, Khi thất bại, Cổng ra); đổi loại node ở chữ màu dưới tên; nút ⓘ hiện mô tả, bấm ra ngoài đóng
- [ ] V2 Khung app (`wails dev`): nút Chạy (F8 vẫn chạy / dừng), ▾ mở menu – Chạy từ node đang chọn, Chạy từng node (F10 đi tiếp), Xem trước, Sửa input trước khi chạy (kịch bản có input); menu theme đổi Sáng / Tối / Theo hệ thống, lưu sau khi mở lại app; kéo cửa sổ bằng thanh tiêu đề vẫn được (không kéo nhầm khi bấm nút theme); thư viện thu gọn Mở gần đây / Lịch sử chạy nhớ sau khi mở lại
- [ ] V1 Thẻ node mới trên canvas app thật (`wails dev`): kéo nối từ chấm cổng nhỏ (8 px) vẫn dễ trúng; rê chuột vào chấm cổng chưa nối → tooltip "chưa nối: …"; rê vào thẻ có chip ("1 input", "thử lại 2") → tooltip; mở kịch bản cũ rồi bấm *Sắp xếp lại* → dàn đều; xem cả theme sáng + tối (màu nhóm đọc rõ)

## 37 – Phần tử màn hình qua Accessibility · [phase3/37-phan-tu-man-hinh.md](phase3/37-phan-tu-man-hinh.md)

- [ ] E4 Đọc cạnh nhãn (cần thêm quyền Screen Recording): Mẫu thử **Phần tử – đọc chữ cạnh nhãn (OCR)** → lời kết thúc ghi *Số dư 5430000* (số vẽ bằng canvas) và giá cũ → giá mới (chỉ số); inspector Đọc phần tử: *Tìm theo* Chữ (OCR) bỏ được lựa chọn *Chính nó*, Thử tìm nhãn thấy đúng chỗ, gõ mẫu lọc sai `(\d+` → báo đỏ; thử trên app thật có giá trị không đọc được qua cây (game / canvas). Xem cả theme sáng + tối
- [ ] E3 Chờ / Click theo phần tử (Chrome thường mở `element.html` hoặc trang thật): bảng thêm node → Click theo phần tử; **Chọn trên màn** → bấm nút “Mua ngay” / nút chỉ có icon → ô cách tìm tự điền (text="…" / aria="…"), khung hiện đúng chỗ trên overlay; **Thử tìm** thấy đúng, đổi chữ sai → liệt kê phần tử đang có; chạy node bấm trúng. Bấm vào chỗ trống / ảnh không có chữ → báo lỗi rõ
- [ ] E2 Mẫu thử (Chrome thường làm trình duyệt mặc định, cấp quyền Accessibility cho app / terminal): **Phần tử – đọc giá đang đổi** → lời kết thúc ghi đúng giá đang hiện + nam@example.com; **Phần tử – bấm nút chỉ có icon** → trái tim tô đậm, thành công; **Phần tử – nút dưới cuối trang** → trang tự cuộn xuống, bấm trúng XÁC NHẬN ĐƠN, thành công. Lần đầu với một app có thể chậm ~0,5–2 s (app dựng cây). Thử lại khi Chrome ở desktop (Space) khác / thu nhỏ → chờ hết giờ, lỗi ghi *không có cửa sổ đang hiện*. Tắt quyền Accessibility → lỗi quyền ngay
- [ ] E2 Inspector Đọc phần tử: thêm từ bảng (nhóm Ứng dụng – đổi ở 38 V0), gõ `css=.x` → báo đỏ chỉ dùng text= / aria= / role=; ô App nhận `{{biến}}`. Xem cả theme sáng + tối

## 35 – Node nâng cao · [phase3/35-node-nang-cao.md](phase3/35-node-nang-cao.md)

- [ ] N8 mẫu: **Tìm và mở app trên máy** (Calculator mở lên; tắt bộ gõ tiếng Việt), **Mở Spotify và phát nhạc** (cần Spotify đăng nhập sẵn), **Săn sale – chờ đúng giờ mở bán** (nhập giờ ~1 phút sau: trang ghi đặt được sau giờ mở bán vài chục ms), **Chụp vùng rồi báo qua HTTP** (HTTP 200)
- [ ] N7 Vẽ đường: Mẫu thử → **Vẽ đường – sao, tròn, chữ nhật, đường gấp khúc** → Chạy: Bảng vẽ có 4 vệt đúng hình, *Số lần kéo* = 4. Overlay (👁): nét viền hiện trước đúng hình; đường gấp khúc có marker n·1, n·2… kéo được, nét đổi theo. Xem trước: con trỏ giả đi đúng nét, giữ chuột. Inspector: thêm / xoá điểm, Khép kín, đổi sang Sao giữ tâm ở giữa. Xem cả theme sáng + tối
- [ ] N6 Chờ đến giờ: Mẫu thử → **Chờ đến giờ – bấm MUA NGAY đúng giờ** → Chạy, nhập giờ ~1 phút sau (vd 20:15:00.000): thanh mini *chờ đến giờ – còn …* đếm ngược, đúng giờ trang ghi *bấm lúc 20:15:00.0xx*; lời kết thúc ghi trễ vài chục ms. Bấm Dừng / tạm dừng giữa lúc chờ ăn ngay. Inspector: xem trước *Còn … / Đã qua …* đổi theo 3 lựa chọn; giờ sai dạng báo đỏ. Xem cả theme sáng + tối
- [ ] N5 Âm lượng: Mẫu thử → **Hệ thống – giảm rồi tăng âm lượng** → Chạy: thanh âm lượng giảm rồi về mức cũ, lời kết thúc ghi 2 mức. Inspector: kéo thanh trượt = ô %; Tắt tiếng / Bật tiếng chạy được. Xem cả theme sáng + tối. Windows (khi có máy): đặt / tăng / giảm / tắt tiếng qua Core Audio
- [ ] N4 Tệp / thư mục: Mẫu thử → **Tệp / thư mục – ghi nhật ký, thêm dòng Excel, dọn file tạm** → Chạy: hiện hộp *Kịch bản này sẽ xoá file* (Huỷ thì không chạy) → Cho phép → `~/Documents/autoclicker-mau-thu` có nhat-ky.txt, don-hang.csv (mở bằng Excel đúng tiếng Việt), don-hang.xlsx (sheet *Đơn hàng*); 2 file tam-*.txt nằm trong Thùng rác. Chạy lại: không hỏi nữa, csv / xlsx thêm dòng. Đổi mẫu tên → hỏi lại. Node *Dọn Thùng rác*: macOS hỏi quyền điều khiển Finder lần đầu. Đặt thư mục `/` hoặc `~` → node báo bị chặn. Xem cả theme sáng + tối. Windows (khi có máy): Recycle Bin, dọn Thùng rác
- [ ] N3 HTTP: Mẫu thử → **HTTP – gọi API, đọc JSON trả về** → Chạy (cần mạng): lời kết thúc ghi HTTP 200, tên “Lan "Hà"”, từ khoá “cà phê & trà”. Inspector: Method + URL một hàng; đổi GET → POST mở sẵn thân JSON; thêm Header `Authorization` → hiện 🔒; JSON sai cú pháp → báo đỏ; Form hiện bảng trường. Đổi URL sang `https://httpbin.org/status/404` → node Thất bại, nhật ký ghi HTTP 404. Xem cả theme sáng + tối
- [ ] N2 Clipboard: Mẫu thử → **Clipboard – ghi, dán, copy lại** → Chạy: ô hiện “Dán từ clipboard: Xin chào Lan!”, lời kết thúc ghi lại đúng chữ. Node Clipboard *Ghi ảnh* với `{{n?.output.path}}` của node Chụp màn hình → dán vào chat ra ảnh. Windows (khi có máy): đọc / ghi tiếng Việt, ghi ảnh
- [ ] N1 Chụp màn hình: Mẫu thử → **Chụp màn hình – chụp chỗ thấy chữ + cả màn** → Chạy: 2 ảnh trong `~/Pictures/Auto Clicker` (thư mục ảnh trong Cài đặt), dán (Cmd+V) vào app chat ra ảnh chữ ĐÃ XONG. Inspector: đổi *Chụp* qua 4 kiểu; *Vùng cố định* → **Khoanh trên màn hình…** → overlay ghi *Dùng vùng này* → về inspector đã có X / Y / W / H; *Cả màn* liệt kê đúng các màn đang cắm. Xem cả theme sáng + tối. Windows (khi có máy): chép ảnh vào clipboard
- [ ] N0 mô tả node: Ctrl/Cmd+K → rê chuột / phím mũi tên qua các mục, khung dưới đổi mô tả theo mục, bảng không giật; inspector có thẻ mô tả dưới tab Node, × ẩn → mở node khác vẫn ẩn, ⓘ hiện lại, tắt app mở lại vẫn nhớ – xem cả theme sáng + tối

## 25d – Ghi chú node web · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Ctrl/Cmd+K: nhóm *Web (trình duyệt)* nằm cuối, dưới tên nhóm có dòng *Chỉ dùng cho trang web cơ bản…*; inspector node web và node Mở chọn *Chrome của app* hiện cùng ghi chú; node web vẫn thêm / chạy như cũ – xem cả theme sáng + tối

## 34 – Thao tác giống người · [phase3/34-giong-nguoi.md](phase3/34-giong-nguoi.md)

- [ ] Mẫu thử → **Thao tác giống người – vệt chuột, nhịp gõ** → Chạy: vệt chuột trên trang cong, không nhảy; *Điểm chuột / lần di* hàng chục, *Giữ nút* 60–150 ms, nhịp gõ lệch. Cài đặt chạy → **Tắt** → chạy lại: 1–2 điểm mỗi lần di, cong 0 px
- [ ] Bấm Dừng / F8 / rê chuột vào góc khi chuột đang đi giữa đường → dừng ngay, không bấm thêm; chuột đi gần mép / góc màn không tự kích hoạt dừng khẩn cấp
- [ ] Kịch bản mới: Cài đặt chạy mặc định **Vừa**; inspector node Click / Gõ có ô *Thao tác giống người* (Theo kịch bản / Tắt / …), node web chạy nền thì không có; thẻ node ghi đè có chip – xem cả theme sáng + tối
- [ ] Overlay → **Xem trước** kịch bản mức Vừa / Cao: con trỏ giả đi đường cong, chậm – nhanh – chậm, đôi khi đi lố rồi khựng + chỉnh; vệt vẽ cong theo đường đã đi; kéo thả cong trong thời gian kéo; node ghi đè Tắt (hoặc kịch bản Tắt) bay thẳng như cũ; bấm Dừng giữa đường dừng ngay
- [ ] Kéo thả (kanban) và double click vẫn ăn khi bật giống người; Windows (khi có máy): click giữ nút + đường chuột đúng chỗ khi DPI 125 %

## 20 F5c – Output mọi node + ô số nhận biến · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Mẫu thử → **Ô số từ biến – số click, nhịp, thời gian chờ** → Chạy: ô Click trái = 5; lời kết thúc ghi thời gian click + thời gian đã chờ gần bằng nhau; form Chạy hỏi *Nhịp giữa các click*
- [ ] Inspector node Click: bấm **{x}** cạnh *Số click* → ô biến chiếm cả hàng, gõ `{{` gợi ý output node trước (cả `startedAt`, `durationMs`); bấm lại {x} với ô là `4` → về ô số = 4. Ô đang ghim input thì {x} mờ – xem cả theme sáng + tối
- [ ] Gõ sai ô biến (vd `{{n1.status}}`, `abc`) → Chạy báo lỗi rõ ở node đó; biến ra 0 cho Số click → node Thất bại “số click phải ≥ 1”

## 20 F5b – Toạ độ nhận biến · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Mẫu thử → **Kanban – kéo theo toạ độ từ biến** → Chạy: thẻ VIỆC 1 sang cột HOÀN TẤT (HOÀN TẤT = 1); thử với Chrome zoom 125% vẫn đúng
- [ ] Inspector node Click: Vị trí → **Biến** → ô X / Y điền sẵn số hiện tại; chọn *Điền chỗ thấy… → n1* ra `{{n1.output.x}}` / `{{n1.output.y}}`, sửa thành `+ 20`; marker trên overlay của node biến mất, thẻ ghi `X {{n1.output.x}} + 20`; bấm **Số** lại → về điểm cũ. Xem cả theme sáng + tối
- [ ] Gõ sai (vd `{{n1.output.text}}`, `{{n1.output.x}} +`) → Chạy báo lỗi rõ ở node đó trước khi chạy; Kéo thả B *Lệch so với A* + Biến → ô ghi ΔX / ΔY

## 20 F6 – Thanh mini nổi trên app toàn màn hình · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Chạy từng node (F10) → sang Chrome / app khác **toàn màn hình**: thanh mini vẫn hiện, bấm được nút chạy tiếp / dừng; overlay (chọn điểm, Thử tìm) cũng hiện trên app toàn màn hình
- [ ] Trong lúc thanh mini / overlay: icon Dock của app tạm ẩn (đúng như thiết kế); về màn chính icon Dock + menu Edit (Cmd+C / V trong ô nhập) trở lại, Cmd+Tab thấy app
- [ ] Thanh mini không giành focus: đang chạy kịch bản gõ phím vào app khác vẫn gõ đúng chỗ

## 33 – Kịch bản mẫu để test · [phase3/33-mau-thu.md](phase3/33-mau-thu.md)

- [ ] `wails dev` → Thư viện → **Mẫu thử…**: lọc theo kiểu / loại bước, tìm không dấu; *Mở sân thử* mở danh sách trang trong trình duyệt
- [ ] *Dùng mẫu* khi đang có thay đổi chưa lưu → hỏi lưu; mở xong là kịch bản mới chưa có file, Lưu hỏi chỗ lưu
- [ ] Chạy vài mẫu từ app (Click nhiều lần, Kéo thả, Form đăng nhập nhập sai mật khẩu, Kanban, Quiz): trang sân thử hiện đúng kết quả (đã chạy tự động 31 / 31 bằng engine – còn thử qua nút Chạy, thanh mini, F8)
- [ ] Click / bấm theo chữ trong app thật khác (không phải sân thử) không chậm đi rõ (mỗi lần di chuột chờ thêm 40 ms)

## 32 – Tách Go thành package · [phase3/32-go-packages.md](phase3/32-go-packages.md)

- [ ] `wails dev`: app mở được, nháp lần trước khôi phục, mở / lưu kịch bản, danh sách mở gần đây, đổi theme
- [ ] Chạy một kịch bản có click + chờ ảnh / chữ + node điều khiển: thanh mini, F8 dừng, chuột chạm góc dừng; lịch sử chạy có lần mới (kèm ảnh khi node lỗi)
- [ ] Thử tìm ảnh / chữ, Thử node, Lấy màu, chụp ảnh mẫu trên overlay; form Chạy nhớ giá trị input lần trước

## 20 F5 – Lúc chạy, nhật ký, lịch sử, debug · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] Chạy ví dụ Facebook (`examples/thu-node-dieu-khien.json`): về màn chính thấy bảng Nhật ký, đường đã đi sáng, ×n trên node lặp; bấm dòng chọn node; node lỗi có ảnh màn hình
- [ ] Node Kết thúc thành công có lời nhắn → "Xong: …" ở thanh trạng thái
- [ ] Lịch sử chạy trong thư viện: chạy vài lần, mở lại lần cũ; mở kịch bản khác thì danh sách đổi theo
- [ ] Debug: đặt Điểm dừng ở một node → Chạy dừng trước node đó (thanh mini), F10 / nút chạy từng node, Tiếp tục; *Từ node này*; *Từng node*

## 20 F3 – Node điều khiển · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] ví dụ Facebook bằng Switch theo màn hình + Đếm + Kết thúc; Nếu thấy ảnh (khoanh ảnh mẫu cho điều kiện), màu (chọn điểm, Lấy màu); Thử node vẽ khung trên overlay – file thử `examples/thu-node-dieu-khien.json`
- [ ] Lời nhắn node Kết thúc thất bại hiện ở "Lỗi lần chạy trước" (thanh trạng thái), node đó được chọn trên canvas

## 20 F4b – Canvas trong app · [phase3/20-flow.md](phase3/20-flow.md)

- [ ] soạn ví dụ Facebook bằng canvas (node hành động + cạnh `fail` quay lại), thả cổng ra chỗ trống thêm node, đặt điểm liên tiếp trên overlay, hoàn tác, chạy thật
- [ ] Thay cho kiểm tra 25b cũ: chờ chữ không thấy → thử lại 2 lần → cạnh `fail` sang node cuộn rồi quay lại chờ chữ

## 25e – Tìm theo chữ · [phase3/25e-ocr.md](phase3/25e-ocr.md)

- [ ] kịch bản Facebook đổi sang Chờ chữ / Click theo chữ

## 25 – Chờ ảnh / click theo ảnh / chờ pixel · [phase3/25-cho-anh.md](phase3/25-cho-anh.md)

- [ ] khoanh nút trên trang web → Thử tìm → chạy chờ ảnh + click theo ảnh khi cửa sổ dời chỗ; chờ pixel đổi màu; thanh mini khi đang chờ

## 30 – Sửa nhiều điểm · [phase3/30-sua-nhieu-diem.md](phase3/30-sua-nhieu-diem.md)

- [ ] Chọn 5 điểm, co giãn 150%, xoay 10°, xiên, lật → overlay và lần chạy thật khớp; hoàn tác / làm lại

## 17 – Đa màn hình (macOS) · [phase1-mvp.md](phase1-mvp.md)

- [ ] macOS 2 màn: đổi màn trên overlay, đặt điểm / chọn vùng ở màn phụ, Xong → main về đúng chỗ cũ
- [ ] macOS: dừng khẩn cấp ở góc màn phụ, góc tiếp giáp 2 màn không dừng nhầm; overlay hiện "tỷ lệ 1.00"

## Chờ máy Windows

- [ ] Bước 17: scale 100% / 125% / 150% và 2 màn khác scale; 23: mở URL / app (PowerShell `Start-Process`), cuộn theo nấc – làm cùng bước 18 ([phase1-mvp.md](phase1-mvp.md))

## Tạm hoãn – node web qua Chrome của app · [phase3/36-chrome-app-sau.md](phase3/36-chrome-app-sau.md)

Người dùng chốt 2026-10-18: việc qua Chrome của app làm sau cùng – thử khi làm 36.

### 25d D4b – Node web: chạy nền · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Tắt hẳn Chrome của app (đã mở từ trước thì chưa có cờ chạy nền) → Mẫu thử → **Web: chạy nền** → Chạy, ngay sau khi trang mở thì chuyển sang app khác / thu nhỏ Chrome: chuột không động, kết thúc thành công “XIN CHÀO DEMO, lăn … lần”
- [ ] Inspector node Click / Gõ / Cuộn / Chờ phần tử web: *Điều khiển* Chuột thật / Chạy nền + dòng giải thích đổi theo; node Đọc không có lựa chọn này; thẻ node hiện “· chạy nền” – xem cả theme sáng + tối
- [ ] Chạy nền khi Chrome zoom 125% và khi tab đích không phải tab đang hiện (tab khác đang mở trên cùng cửa sổ) → vẫn bấm đúng phần tử

### 25d D4a – Node web: gõ, đọc, cuộn · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Ctrl+K nhóm *Web (trình duyệt)* có Gõ vào ô web / Đọc phần tử web / Cuộn trang web; inspector từng loại (công tắc Xoá chữ cũ / Enter, Đọc: chữ / giá trị ô / thuộc tính + ô tên, Cuộn: 3 cách + hướng / số nấc / khung) – xem cả theme sáng + tối
- [ ] Mẫu thử → **Web: đăng nhập, đọc kết quả** → Chạy: gõ vào đúng ô, lần sai hiện chữ đỏ, lần đúng ra XIN CHÀO DEMO, kết thúc thành công
- [ ] Mẫu thử → **Web: cuộn feed tới bài đặc biệt** và **Web: cuộn trong khung + tới phần tử** → Chạy: chuột lăn thật (feed tải thêm), dải ngang lăn mà trang không chạy, dừng ở MỤC 30
- [ ] Node Đọc → **Thử tìm** hiện “đọc được “…”” + số phần tử khớp; node Gõ trỏ nhầm vào nút (không phải ô) → cổng Thất bại “ô nhập không nhận focus”, không gõ bậy

### 25d D3 – Node web: Chọn phần tử, Thử tìm, trạng thái · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Inspector node Chờ / Click phần tử web: dòng trạng thái “Chrome của app: chưa mở” → bấm Mở → Chrome riêng bật, trạng thái thành “đang mở, 1 tab”; xem cả theme sáng + tối
- [ ] Mở trang sân thử web trong Chrome của app → **Chọn phần tử trên trang** → Chrome lên trước, viền cam theo chuột, bấm nút “Đồng ý” → app tự lên lại, ô cách tìm điền 2–3 dòng (`text=Đồng ý`, `role=button[name=…]`…), trang **không** nhận click đó; thử Esc / Huỷ chọn
- [ ] Chọn ô đồng ý trong iframe (mẫu Web: thanh toán trong iframe) → ra `role=checkbox…`
- [ ] **Thử tìm**: lần đầu chuột đi vào trang đo toạ độ rồi app lên lại, báo “Thấy … tại X Y · ms” + chữ; “xem trên màn hình” → khung xanh đúng chỗ phần tử; phần tử bị popup cookie che → báo bị che (khung vàng); phần tử dưới màn hình → “ngoài khung nhìn, lúc click sẽ tự cuộn tới”

### 25d D2 – Node web: trước khi bấm, iframe, role · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Mẫu thử → **Web: thanh toán trong iframe** → Chạy: tick ô đồng ý, chờ nút Thanh toán trượt xong mới bấm, hiện ĐÃ THANH TOÁN, *Click thật* = 2
- [ ] Kéo cửa sổ Chrome của app xuống để nút “Xem lịch sử” / Thanh toán rơi vào chỗ **Dock che** rồi chạy mẫu → node tự cuộn phần tử lên giữa rồi bấm trúng (không bấm vào Dock)
- [ ] Trong lúc node webClick đang chờ: đặt một cửa sổ app khác (luôn trên cùng, vd thanh mini) đè lên nút → không bấm nhầm; hết giờ báo “chưa bấm được: …” rõ lý do

### 25d D1 – Node web (chờ / click phần tử) · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] `wails dev` → Mẫu thử → **Web: mở lịch sử đơn hàng** → Chạy: lần đầu app bật Chrome riêng (profile của app), chạy hết tới LỊCH SỬ ĐƠN HÀNG, ô *Click thật* tăng; chạy lại lần 2 → mở tab mới trong Chrome đó (không bật Chrome thứ hai)
- [ ] Trong lúc chạy mẫu: dời cửa sổ Chrome / zoom 125% (Cmd +) / mở DevTools rồi chạy lại – vẫn bấm trúng (đo lại toạ độ)
- [ ] Tắt Chrome của app rồi chạy một node web (không có node Mở trước) → lỗi rõ “Chrome của app chưa mở…” (mã `browser`), không treo
- [ ] Inspector: thêm node qua Ctrl+K nhóm *Web (trình duyệt)*; nhập nhiều dòng selector, `{{biến}}` trong selector; node Mở chọn “Chrome của app”
- [ ] Trang thật đăng nhập (vd Facebook) trong Chrome của app: đăng nhập một lần, lần sau profile còn nhớ

### 25d D0 – Đo toạ độ web · [phase3/25d-web-dom.md](phase3/25d-web-dom.md)

- [ ] Cắm màn ngoài (khác DPI Retina) rồi chạy `go run -tags probe ./cmd/d0probe -scenario all` – tình huống `display2` lệch ≤ 2 px (chuột sẽ tự di chuyển + bấm ~1 phút, đừng động vào chuột)
- [ ] Khi có máy Windows: chạy lại trên màn DPI 125% (sửa `-chrome` trỏ tới chrome.exe; phím zoom là Ctrl) – dán kết quả vào file bước
