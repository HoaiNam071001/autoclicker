# 31a. Kết nối AI qua MCP – tạo / sửa kịch bản

- [ ] **31a. Kết nối AI qua MCP – tạo / sửa kịch bản** *(người dùng yêu cầu, 2026-10-17)* – Người dùng ra lệnh trên Claude / GPT / Cursor… và AI tạo ra một hoặc nhiều kịch bản trong app, không cần soạn tay. App làm **MCP server** (Model Context Protocol – chuẩn chung, một server dùng được cho mọi AI hỗ trợ MCP).
  - Kiến trúc:
    - App đang chạy mở MCP server **Streamable HTTP** tại `127.0.0.1:<port>/mcp` (chỉ nghe localhost, token trong `UserConfigDir/autoclicker/`); dùng SDK chính thức `github.com/modelcontextprotocol/go-sdk`. dùng chung server HTTP + token của API local (22)
    - Cùng file chạy có chế độ `autoclicker --mcp` (**stdio**) cho AI chỉ chạy được server dạng lệnh (Claude Desktop, Codex CLI, Gemini CLI…): **luôn chuyển tiếp** sang app đang chạy, app chưa mở thì tự mở – không tự chạy engine (*2026-10-18*: quyền Accessibility / Screen Recording của macOS tính theo tiến trình, chạy trong tiến trình con của Claude Desktop thì quyền thuộc Claude Desktop)
    - Kết nối được: Claude Desktop / Claude Code / Cursor / GitHub Copilot (chế độ agent trong VS Code) / Codex CLI / Gemini CLI (stdio hoặc HTTP localhost). Copilot Studio (Microsoft 365) cũng chỉ nhận MCP từ xa như ChatGPT. ChatGPT bản web / desktop chỉ nhận MCP từ xa qua HTTPS → cần tunnel (cloudflared / ngrok) + token; ghi rõ trong hướng dẫn, không bật mặc định
    - **Dùng được với mọi AI, không riêng Claude** *(người dùng hỏi 2026-10-18)* – thiết kế theo mức chung thấp nhất của các client: tham số tool **phẳng, ít lồng** (Gemini / OpenAI strict kén `oneOf`, lồng sâu) – schema đầy đủ theo từng loại node để ở resource, `validate_script` bắt lỗi; gom **≤ ~20 tool** (client chọn tool kém khi quá nhiều); tool nhận cả `id` lẫn `ref` (n2); kết quả **chữ trước, ảnh sau** (client không xem được ảnh vẫn dùng được – 31b); việc chạy lâu trả `runId` ngay (31c). Luật an toàn (hỏi quyền, công tắc AI, F8) nằm trong app nên AI nào gọi vào cũng như nhau
  - [ ] Schema cho AI: sinh **JSON Schema** của `Script` / `Step` từ `model.go` (mô tả tiếng Việt + tiếng Anh từng trường, giá trị mặc định, loại bước, neo, `{{biến}}`); test kiểm tra schema khớp struct Go (như `types.ts`). Resource `autoclicker://schema`, `autoclicker://guide` (cách viết kịch bản bền: ưu tiên mở URL / phím tắt / chờ ảnh hơn toạ độ cứng, thêm chờ sau khi mở app…), `autoclicker://examples` (vài kịch bản mẫu)
  - [ ] Tool đọc: `list_scripts` (thư viện + mở gần đây), `read_script`, `get_step_types` (loại bước app đang hỗ trợ – AI không dùng loại *Sắp có*), `get_screen_info` (kích thước, tỷ lệ Retina / DPI, số màn hình)
  - [ ] Tool ghi: `validate_script` (trả lỗi rõ theo từng bước để AI tự sửa), `create_script` (một hoặc nhiều kịch bản, lưu vào thư viện), `update_script` (sửa theo `stepId`, không ghi đè cả file; gửi kèm `revision` đọc được lúc `read_script` – kịch bản đã đổi từ lúc đó thì báo xung đột, AI đọc lại rồi sửa; revision lấy từ 22, không ghi vào file), `open_in_app` (mở kịch bản trên MainScreen cho người dùng xem)
  - [ ] Prompt MCP `tao_kich_ban`: khung hỏi đáp mẫu (mục tiêu, app / trang web, giá trị nào là input) để AI hỏi lại khi thiếu thông tin
  - [ ] An toàn: kịch bản AI tạo có `meta.createdBy: "ai"` + chip "AI" trong thư viện, mở ở trạng thái chưa lưu để người dùng duyệt; AI sửa kịch bản đang mở có thay đổi chưa lưu → tạo bản sao, không ghi đè; mọi thay đổi của AI đi qua store nên Ctrl+Z hoàn tác được; bước này **chưa cho AI chạy** kịch bản (31c)
  - [ ] Cài đặt [màn Cài đặt]: bật / tắt MCP, cổng, token (sao chép / tạo lại), nút *Thêm vào Claude Desktop* (ghi `claude_desktop_config.json`, hỏi trước), đoạn cấu hình mẫu cho Claude Code / Cursor / Codex; nhật ký các lệnh AI đã gọi
  - [ ] Test Go: gọi tool qua client MCP trong test (in-memory transport): tạo kịch bản hợp lệ / lỗi, sửa theo `stepId`, schema khớp struct
  - Xong khi: trên Claude Desktop nói "tạo kịch bản mở Spotify, tìm bài Tinh Hà Say 2 và phát" → kịch bản hiện trong thư viện app (chip AI), Xem trước chạy đúng thứ tự, bấm Chạy thì phát nhạc; nói "đổi thành bài khác, cho tên bài thành input" → AI sửa đúng kịch bản, có input `{{bai_hat}}`
- [ ] **31b. MCP – AI nhìn màn hình** – Để AI tự biết vị trí nút thay vì đoán toạ độ.
  - [ ] Tool `describe_screen` – **chữ trước ảnh** *(2026-10-18)*: các dòng chữ đọc được (OCR 25e) kèm khung, danh sách cửa sổ (phần tử trang web của Chrome của app → [36](36-chrome-app-sau.md), 2026-10-18) – rẻ token, toạ độ chính xác, dùng được với AI không xem ảnh. Chỉ bọc hàm Go sẵn có
  - [ ] Tool `screenshot` (toàn màn hình / vùng / cửa sổ, thu nhỏ cho nhẹ, trả ảnh + tỷ lệ); toạ độ AI gửi về theo pixel ảnh, Go tự đổi sang toạ độ robotgo (Retina / DPI – dùng chung cách đổi của 17, 25)
  - [ ] Tool `find_image` (dùng so khớp mẫu của 25), `capture_template` (AI khoanh vùng trên ảnh chụp → lưu ảnh mẫu cạnh kịch bản cho bước `waitImage` / `clickImage`), `get_mouse_position`, `list_windows`
  - [ ] Người dùng thấy khi AI chụp màn hình: thông báo nhỏ / chấm đỏ trên thanh mini; tắt được trong Cài đặt
  - Xong khi: AI chụp màn hình YouTube, tự khoanh nút "Phát" làm ảnh mẫu, kịch bản tạo ra dùng `clickImage` thay vì toạ độ cứng và vẫn chạy khi cửa sổ dời chỗ
- [ ] **31c. MCP – AI tự thử và chạy kịch bản** – AI làm từng bước trên màn hình thật, xem kết quả, bước nào đúng thì ghi vào kịch bản (giống người dùng soạn rồi bấm thử).
  - [ ] Tool `run_step` (chạy một bước, trả ảnh chụp sau bước), `run_script` (kèm giá trị input, dùng `resolveInputs`), `get_run_status` (tiến độ + nhật ký 25c), `stop_script`. Chạy **không chờ**: `run_*` trả `runId` ngay, AI hỏi `get_run_status` / `wait_run(runId, tối đa ~30 s)` – client MCP hay cắt tool chạy lâu (dùng phiên chạy có `runId` của 22)
  - [ ] **Chạy một phần flow** *(2026-10-18)* – để AI sửa node n7 không phải chạy lại từ đầu (mở app, đăng nhập…). Đã có *chạy từ node* (`RunOptions.from`, biến node trước trống); thêm vào `engine.Options`:
    - `StopAt` (chạy đến node) và `Only` (chỉ chạy các node đã chọn – một cụm)
    - `Seed`: output node trước lấy từ một lần chạy trong lịch sử → `resume_run(runId, từ node)` = chạy tiếp từ node lỗi với biến cũ
    - Ghim output node (như pin data n8n): lưu trên máy theo kịch bản, **không ghi vào file**; thẻ node có dấu ghim
    - App dùng chung: chuột phải node → *Chạy đến đây* / *Chạy cụm đã chọn* / *Chạy tiếp từ đây (dùng kết quả lần trước)*
  - [ ] Chạy thử khô `simulate_script`: chạy flow với thiết bị giả (`device` đã là interface), không đụng chuột – AI kiểm tra nhánh / biến trước khi chạy thật. Node tìm ảnh / chữ / web trả theo kịch bản AI đưa (`found` / `notFound`) hoặc mặc định `success`
  - [ ] Quyền: mặc định mỗi lần AI muốn chạy thì app hiện hộp hỏi (tên kịch bản, các bước, Cho phép / Từ chối, "cho phép kịch bản này đến khi đóng app"); F8 và dừng khẩn cấp (góc màn hình) vẫn luôn có tác dụng; giới hạn số bước / thời gian mỗi lượt AI chạy
  - [ ] **Công tắc AI** trên thanh mini + Cài đặt *(2026-10-18)*: 3 mức **Tắt** (từ chối mọi tool) / **Chỉ soạn** / **Soạn + chạy (hỏi quyền)**; chuyển sang Tắt thì huỷ ngay lượt AI đang chạy; chấm báo khi AI đang điều khiển. Mức "tự chạy không hỏi" chưa làm – "cho phép đến khi đóng app" đã đủ, cần thì thêm sau
  - [ ] Người dùng đụng chuột lúc AI đang chạy (vị trí chuột khác chỗ engine vừa đặt) → tạm dừng, báo AI lý do
  - [ ] Chế độ "AI soạn bằng cách thử": các bước AI `run_step` thành công được ghi lại thành kịch bản nháp, kèm ảnh chụp từng bước để người dùng xem lại
  - [ ] Lịch sử chạy (20 F5) ghi rõ lần nào do AI chạy
  - Xong khi: nói "mở bài Tinh Hà Say 2 trên YouTube" → AI hỏi quyền, tự thử (mở URL, chờ, click kết quả đầu), nhạc phát, kịch bản được lưu để lần sau chạy lại không cần AI; AI sửa một node lỗi rồi `resume_run` chạy tiếp từ node đó, không chạy lại từ đầu
- [ ] **31d. AI bên trong flow** *(người dùng hỏi 2026-10-18 – **chỉ làm khi 31a–c đã dùng ổn**, không làm sớm)* – AI xử lý ngay trong luồng (không phải AI bên ngoài điều khiển app). Mặc định tắt; bật thì báo rõ: tốn tiền mỗi lần chạy, ảnh / chữ màn hình gửi ra ngoài (trừ AI chạy trên máy).
  - [ ] Node **Hỏi AI** (`askAI`): đầu vào chữ / ảnh từ output node trước + câu hỏi (nhận biến) → output `text` / `json`, hoặc chọn cổng ra như Switch. Bật / tắt riêng từng node (`enabled` có sẵn): AI tắt / không có key → cổng `skip` với giá trị mặc định, flow vẫn chạy không cần AI
  - [ ] Không gắn một hãng: interface `ai.Provider` (`Ask(prompt, images) → text / json`) trong `internal/ai`, bản giả cho test; bản cài: Anthropic, OpenAI, Gemini, **tương thích OpenAI** (URL + key – phủ luôn AI chạy trên máy Ollama / LM Studio: không tốn tiền, ảnh không ra khỏi máy). Cài đặt chọn nhà cung cấp mặc định, node đổi được; key lưu trên máy, không ghi file kịch bản / log
  - [ ] Tìm phần tử bằng AI (chuyển từ 28): mô tả bằng lời ("nút Gửi") → AI có vision trả khung; dự phòng khi ảnh mẫu / chữ / DOM không thấy. Không tự sửa kịch bản – ghi *đề xuất sửa* cho người dùng duyệt
  - [ ] Node **Gọi tool MCP** (`mcpTool`) – app làm MCP client (SDK Go có cả client): kết nối MCP server khác (Sheets, Notion, server tự viết…), gọi tool như một node, kết quả vào output. Lần đầu kết nối mỗi server phải hỏi quyền; tool ghi / xoá đi qua chỗ hỏi quyền chung (22); kết quả chỉ là dữ liệu, không làm theo như lệnh. Làm sau node Hỏi AI, chỉ khi người dùng cần
  - Xong khi: mẫu *đọc tin nhắn → Hỏi AI "cần trả lời không, soạn câu trả lời" → gõ nháp* chạy được với một nhà cung cấp mạng và một AI chạy trên máy; tắt AI thì flow đi cổng `skip`

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
