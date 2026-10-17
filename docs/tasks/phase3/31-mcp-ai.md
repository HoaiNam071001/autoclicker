# 31a. Kết nối AI qua MCP – tạo / sửa kịch bản

- [ ] **31a. Kết nối AI qua MCP – tạo / sửa kịch bản** *(người dùng yêu cầu, 2026-10-17)* – Người dùng ra lệnh trên Claude / GPT / Cursor… và AI tạo ra một hoặc nhiều kịch bản trong app, không cần soạn tay. App làm **MCP server** (Model Context Protocol – chuẩn chung, một server dùng được cho mọi AI hỗ trợ MCP).
  - Kiến trúc:
    - App đang chạy mở MCP server **Streamable HTTP** tại `127.0.0.1:<port>/mcp` (chỉ nghe localhost, token trong `UserConfigDir/autoclicker/`); dùng SDK chính thức `github.com/modelcontextprotocol/go-sdk`. dùng chung server HTTP + token của API local (22)
    - Cùng file chạy có chế độ `autoclicker --mcp` (**stdio**) cho AI chỉ chạy được server dạng lệnh (Claude Desktop, Codex CLI, Gemini CLI…): chuyển tiếp sang app đang chạy; app chưa mở thì tự mở, hoặc chỉ làm việc với file kịch bản
    - Kết nối được: Claude Desktop / Claude Code / Cursor / GitHub Copilot (chế độ agent trong VS Code) / Codex CLI / Gemini CLI (stdio hoặc HTTP localhost). Copilot Studio (Microsoft 365) cũng chỉ nhận MCP từ xa như ChatGPT. ChatGPT bản web / desktop chỉ nhận MCP từ xa qua HTTPS → cần tunnel (cloudflared / ngrok) + token; ghi rõ trong hướng dẫn, không bật mặc định
  - [ ] Schema cho AI: sinh **JSON Schema** của `Script` / `Step` từ `model.go` (mô tả tiếng Việt + tiếng Anh từng trường, giá trị mặc định, loại bước, neo, `{{biến}}`); test kiểm tra schema khớp struct Go (như `types.ts`). Resource `autoclicker://schema`, `autoclicker://guide` (cách viết kịch bản bền: ưu tiên mở URL / phím tắt / chờ ảnh hơn toạ độ cứng, thêm chờ sau khi mở app…), `autoclicker://examples` (vài kịch bản mẫu)
  - [ ] Tool đọc: `list_scripts` (thư viện + mở gần đây), `read_script`, `get_step_types` (loại bước app đang hỗ trợ – AI không dùng loại *Sắp có*), `get_screen_info` (kích thước, tỷ lệ Retina / DPI, số màn hình)
  - [ ] Tool ghi: `validate_script` (trả lỗi rõ theo từng bước để AI tự sửa), `create_script` (một hoặc nhiều kịch bản, lưu vào thư viện), `update_script` (sửa theo `stepId`, không ghi đè cả file), `open_in_app` (mở kịch bản trên MainScreen cho người dùng xem)
  - [ ] Prompt MCP `tao_kich_ban`: khung hỏi đáp mẫu (mục tiêu, app / trang web, giá trị nào là input) để AI hỏi lại khi thiếu thông tin
  - [ ] An toàn: kịch bản AI tạo có `meta.createdBy: "ai"` + chip "AI" trong thư viện, mở ở trạng thái chưa lưu để người dùng duyệt; AI sửa kịch bản đang mở có thay đổi chưa lưu → tạo bản sao, không ghi đè; mọi thay đổi của AI đi qua store nên Ctrl+Z hoàn tác được; bước này **chưa cho AI chạy** kịch bản (31c)
  - [ ] Cài đặt [màn Cài đặt]: bật / tắt MCP, cổng, token (sao chép / tạo lại), nút *Thêm vào Claude Desktop* (ghi `claude_desktop_config.json`, hỏi trước), đoạn cấu hình mẫu cho Claude Code / Cursor / Codex; nhật ký các lệnh AI đã gọi
  - [ ] Test Go: gọi tool qua client MCP trong test (in-memory transport): tạo kịch bản hợp lệ / lỗi, sửa theo `stepId`, schema khớp struct
  - Xong khi: trên Claude Desktop nói "tạo kịch bản mở Spotify, tìm bài Tinh Hà Say 2 và phát" → kịch bản hiện trong thư viện app (chip AI), Xem trước chạy đúng thứ tự, bấm Chạy thì phát nhạc; nói "đổi thành bài khác, cho tên bài thành input" → AI sửa đúng kịch bản, có input `{{bai_hat}}`
- [ ] **31b. MCP – AI nhìn màn hình** – Để AI tự biết vị trí nút thay vì đoán toạ độ.
  - [ ] Tool `screenshot` (toàn màn hình / vùng / cửa sổ, thu nhỏ cho nhẹ, trả ảnh + tỷ lệ); toạ độ AI gửi về theo pixel ảnh, Go tự đổi sang toạ độ robotgo (Retina / DPI – dùng chung cách đổi của 17, 25)
  - [ ] Tool `find_image` (dùng so khớp mẫu của 25), `capture_template` (AI khoanh vùng trên ảnh chụp → lưu ảnh mẫu cạnh kịch bản cho bước `waitImage` / `clickImage`), `get_mouse_position`, `list_windows`
  - [ ] Người dùng thấy khi AI chụp màn hình: thông báo nhỏ / chấm đỏ trên thanh mini; tắt được trong Cài đặt
  - Xong khi: AI chụp màn hình YouTube, tự khoanh nút "Phát" làm ảnh mẫu, kịch bản tạo ra dùng `clickImage` thay vì toạ độ cứng và vẫn chạy khi cửa sổ dời chỗ
- [ ] **31c. MCP – AI tự thử và chạy kịch bản** – AI làm từng bước trên màn hình thật, xem kết quả, bước nào đúng thì ghi vào kịch bản (giống người dùng soạn rồi bấm thử).
  - [ ] Tool `run_step` (chạy một bước, trả ảnh chụp sau bước), `run_script` (kèm giá trị input, dùng `resolveInputs`), `get_run_status` (tiến độ + nhật ký 25c), `stop_script`
  - [ ] Quyền: mặc định mỗi lần AI muốn chạy thì app hiện hộp hỏi (tên kịch bản, các bước, Cho phép / Từ chối, "cho phép kịch bản này đến khi đóng app"); F8 và dừng khẩn cấp (góc màn hình) vẫn luôn có tác dụng; giới hạn số bước / thời gian mỗi lượt AI chạy
  - [ ] Chế độ "AI soạn bằng cách thử": các bước AI `run_step` thành công được ghi lại thành kịch bản nháp, kèm ảnh chụp từng bước để người dùng xem lại
  - [ ] Lịch sử chạy (20 F5) ghi rõ lần nào do AI chạy
  - Xong khi: nói "mở bài Tinh Hà Say 2 trên YouTube" → AI hỏi quyền, tự thử (mở URL, chờ, click kết quả đầu), nhạc phát, kịch bản được lưu để lần sau chạy lại không cần AI

### Ví dụ mục tiêu: nhờ AI mở một bài nhạc

Người dùng gõ trên Claude: *"tạo kịch bản mở bài Tinh Hà Say 2 trên YouTube"*. Claude gọi `get_step_types` / đọc `autoclicker://guide`, tạo kịch bản, `validate_script`, rồi `create_script` + `open_in_app`:

| # | Bước | Loại | Input | Cần bước |
|---|---|---|---|---|
| 1 | Mở `https://www.youtube.com/results?search_query={{bai_hat}}` | `open` | `bai_hat` = "Tinh Hà Say 2" | 21b, 23, 31a |
| 2 | Chờ trang tải (chờ ảnh khung kết quả, tối đa 10 s) | `waitImage` | | 25, 31b |
| 3 | Click video đầu tiên | `clickImage` (ảnh mẫu AI tự khoanh) | | 25, 31b |

> Bản đơn giản hơn chỉ cần 31a (không cần 31b): mở URL → chờ 4 s → click toạ độ kết quả đầu (AI hỏi người dùng vị trí hoặc chụp màn hình ở 31b). Kém bền hơn khi cửa sổ đổi cỡ.

### Ví dụ mục tiêu: nhắn tin cho một người trên Facebook

| # | Bước | Loại | Input | Cần bước |
|---|---|---|---|---|
| 1 | Mở `facebook.com/{{profile}}` | `open` | `profile` | 21b, 23 |
| 2 | Đưa Chrome lên trước | `focusWindow` | | 26 |
| 3 | Chờ nút "Nhắn tin" (thử lại 2 lần) | `waitImage` | `wait_s` | 25, 25b |
| 4 | Click nút | `clickImage` | `click_delay` | 25 |
| 5 | Gõ "xin chào {{M1}} và {{M2}}" | `type` | `M1`, `M2` | có sẵn + 21b |
| 6 | Enter | `key` | | có sẵn |
| 7 | Cuộn xuống 5 nấc | `scroll` | | 23 |

> Tự động hoá trên Facebook vi phạm điều khoản của họ và dễ bị khoá tài khoản, nhất là khi nhắn hàng loạt. Thử trên tài khoản phụ hoặc trang của chính bạn.
