# 31a. Kết nối AI qua MCP – tạo / sửa kịch bản

- [x] **31a. Kết nối AI qua MCP – tạo / sửa kịch bản** *(người dùng yêu cầu, 2026-10-17)* – Người dùng ra lệnh trên Claude / GPT / Cursor… và AI tạo ra một hoặc nhiều kịch bản trong app, không cần soạn tay. App làm **MCP server** (Model Context Protocol – chuẩn chung, một server dùng được cho mọi AI hỗ trợ MCP).
  - Kiến trúc:
    - App đang chạy mở MCP server **Streamable HTTP** tại `127.0.0.1:<port>/mcp` (chỉ nghe localhost, token trong `UserConfigDir/autoclicker/`); dùng SDK chính thức `github.com/modelcontextprotocol/go-sdk`. dùng chung server HTTP + token của API local (22)
    - Cùng file chạy có chế độ `autoclicker --mcp` (**stdio**) cho AI chỉ chạy được server dạng lệnh (Claude Desktop, Codex CLI, Gemini CLI…): **luôn chuyển tiếp** sang app đang chạy, app chưa mở thì tự mở – không tự chạy engine (*2026-10-18*: quyền Accessibility / Screen Recording của macOS tính theo tiến trình, chạy trong tiến trình con của Claude Desktop thì quyền thuộc Claude Desktop)
    - Kết nối được: Claude Desktop / Claude Code / Cursor / GitHub Copilot (chế độ agent trong VS Code) / Codex CLI / Gemini CLI (stdio hoặc HTTP localhost). Copilot Studio (Microsoft 365) cũng chỉ nhận MCP từ xa như ChatGPT. ChatGPT bản web / desktop chỉ nhận MCP từ xa qua HTTPS → cần tunnel (cloudflared / ngrok) + token; ghi rõ trong hướng dẫn, không bật mặc định
    - **Dùng được với mọi AI, không riêng Claude** *(người dùng hỏi 2026-10-18)* – thiết kế theo mức chung thấp nhất của các client: tham số tool **phẳng, ít lồng** (Gemini / OpenAI strict kén `oneOf`, lồng sâu) – schema đầy đủ theo từng loại node để ở resource, `validate_script` bắt lỗi; gom **≤ ~20 tool** (client chọn tool kém khi quá nhiều); tool nhận cả `id` lẫn `ref` (n2); kết quả **chữ trước, ảnh sau** (client không xem được ảnh vẫn dùng được – 31b); việc chạy lâu trả `runId` ngay (31c). Luật an toàn (hỏi quyền, công tắc AI, F8) nằm trong app nên AI nào gọi vào cũng như nhau
  - [x] Schema cho AI: sinh **JSON Schema** của `Script` / `Step` từ `model.go` (mô tả tiếng Việt + tiếng Anh từng trường, giá trị mặc định, loại bước, neo, `{{biến}}`); test kiểm tra schema khớp struct Go (như `types.ts`). Resource `autoclicker://schema`, `autoclicker://guide` (cách viết kịch bản bền: ưu tiên mở URL / phím tắt / chờ ảnh hơn toạ độ cứng, thêm chờ sau khi mở app…), `autoclicker://examples` (vài kịch bản mẫu)
  - [x] Tool đọc: `list_scripts` (thư viện + mở gần đây), `read_script`, `get_step_types` (loại bước app đang hỗ trợ – AI không dùng loại *Sắp có*), `get_screen_info` (kích thước, tỷ lệ Retina / DPI, số màn hình)
  - [x] Tool ghi: `validate_script` (trả lỗi rõ theo từng bước để AI tự sửa), `create_script` (một hoặc nhiều kịch bản, lưu vào thư viện), `update_script` (sửa theo `stepId`, không ghi đè cả file; gửi kèm `revision` đọc được lúc `read_script` – kịch bản đã đổi từ lúc đó thì báo xung đột, AI đọc lại rồi sửa; revision lấy từ 22, không ghi vào file), `open_in_app` (mở kịch bản trên MainScreen cho người dùng xem)
  - [x] Prompt MCP `tao_kich_ban`: khung hỏi đáp mẫu (mục tiêu, app / trang web, giá trị nào là input) để AI hỏi lại khi thiếu thông tin
  - [x] An toàn: kịch bản AI tạo có `meta.createdBy: "ai"` + chip "AI" trong thư viện, mở ở trạng thái chưa lưu để người dùng duyệt; AI sửa kịch bản đang mở có thay đổi chưa lưu → tạo bản sao, không ghi đè; mọi thay đổi của AI đi qua store nên Ctrl+Z hoàn tác được; bước này **chưa cho AI chạy** kịch bản (31c)
  - [x] Cài đặt [màn Cài đặt]: bật / tắt MCP, cổng, token (sao chép / tạo lại), nút *Thêm vào Claude Desktop* (ghi `claude_desktop_config.json`, hỏi trước), đoạn cấu hình mẫu cho Claude Code / Cursor / Codex; nhật ký các lệnh AI đã gọi
  - [x] Test Go: gọi tool qua client MCP trong test (in-memory transport): tạo kịch bản hợp lệ / lỗi, sửa theo `stepId`, schema khớp struct
  - Đã làm (2026-10-20) – `internal/mcpsrv` (SDK `go-sdk` v1.8, stateless + trả JSON), `mcpapi.go`, `internal/cli/mcp.go`:
    - Người dùng chốt: kịch bản AI tạo lưu vào **Documents/Auto Clicker** (trùng tên thêm số, không ghi đè); MCP **bật sẵn**, tắt trong Cài đặt; app **chỉ mở khi AI gọi `open_in_app`** (thay cho "mở ở trạng thái chưa lưu") – đang dở kịch bản khác thì không mở
    - 8 tool: `list_scripts`, `read_script`, `get_step_types` (kèm hướng dẫn flow + mẹo viết bền), `get_screen_info`, `validate_script`, `create_script`, `update_script`, `open_in_app`; resource `autoclicker://guide|schema|examples`; prompt `tao_kich_ban`
    - AI gửi **dạng gọn** (`script.Compact` – dùng chung với Mẫu thử): node chỉ cần type + trường cần, id / ref tự đặt, không có edges thì nối theo thứ tự; mọi chỗ trỏ node (cạnh, start, imageFrom, đích input) nhận id hoặc ref. `update_script` = `script.Patch` (đè trường gửi lên theo id / ref, thêm / xoá node + cạnh). Cạnh hỏng là lỗi (lúc mở file thì lặng lẽ bỏ)
    - Cấu trúc sai → không lưu; lưu được nhưng chưa chạy được → trả `problems` (vd thiếu ảnh mẫu – AI chưa chụp được tới 31b)
    - Stdio luôn chuyển tiếp: mỗi dòng JSON-RPC = một POST /mcp qua kênh CLI (`api.json`), app chưa mở thì mở lúc nhận request đầu
    - Schema sinh bằng `jsonschema-go` từ struct Go (`[]Cond` trỏ `$defs` để cắt vòng), bỏ `required`; test bắt mọi trường trong `steptypes.go` có trong schema
    - Cài đặt: công tắc MCP, *Thêm vào Claude Desktop…* (hỏi rồi ghi `claude_desktop_config.json`, giữ cài đặt khác), cấu hình mẫu Claude Code / Cursor / Codex / Gemini / HTTP, 5 lệnh AI gần nhất. `meta.client` (tên AI) bỏ – chế độ stateless không biết client
  - Kiểm tra tay: [checks.md](../checks.md) mục 31a
  - Xong khi: trên Claude Desktop nói "tạo kịch bản mở Spotify, tìm bài Tinh Hà Say 2 và phát" → kịch bản hiện trong thư viện app (chip AI), Xem trước chạy đúng thứ tự, bấm Chạy thì phát nhạc; nói "đổi thành bài khác, cho tên bài thành input" → AI sửa đúng kịch bản, có input `{{bai_hat}}`
- [ ] **31b. MCP – vòng thử – sửa: AI tự chạy, xem lỗi, sửa, chạy tiếp** *(người dùng yêu cầu 2026-10-20 – gom phần chạy + nhìn khi lỗi của 31b / 31c cũ lên làm trước)* – AI làm trọn vòng không cần người: mở app nền → tạo kịch bản → chạy một phần → node lỗi thì nhận chi tiết lỗi + ảnh chụp → `update_script` → chạy tiếp từ node lỗi với biến cũ, không chạy lại từ đầu.
  - [ ] **Mở app nền**: stdio / tool chạy mà app chưa mở → mở **chỉ thanh mini** (không hiện màn chính; Wails 1 cửa sổ nên là chế độ minibar của cửa sổ chính) – người dùng vẫn thấy AI đang chạy + bấm Dừng. Tool `show_app` mở màn chính khi cần; `open_in_app` giữ như cũ
  - [ ] **Tool chạy** (không chờ, dùng phiên `runId` của 22): `run_script(script, inputs, from, stopAt, only[], seedRunId, shots)` trả `runId` ngay; `wait_run(runId, timeoutS ≤ ~30)` trả trạng thái + vết chạy gọn; `stop_run(runId)`; `resume_run(runId, from)` = chạy tiếp từ node (mặc định node lỗi) với output các node trước lấy từ lần chạy đó (đọc lại kịch bản mới nhất – sau `update_script`)
  - [ ] **Engine chạy một phần**: thêm `engine.Options.StopAt` (chạy đến node rồi dừng), `Only` (chỉ cụm node đã chọn), `Seed` (output node trước lấy từ một lần chạy – node đã bị xoá thì bỏ, báo trong kết quả). Test Go đủ 3 kiểu + kết hợp `from`
  - [ ] **Vết chạy cho AI** (`wait_run`): mỗi node đã chạy một dòng – id, ref, type, status, cổng ra, output (cắt chuỗi dài), attempts, durationMs, vòng; node lỗi kèm `error.code` + `error.message` + `error.detail` + ảnh. Chữ trước ảnh; ảnh thu nhỏ (cạnh dài ~1280) kèm tỷ lệ đổi về toạ độ thật
  - [ ] **`error.detail` có cấu trúc** (`model.StepError.Detail`, map theo loại – UI nhật ký dùng chung; lấy từ `e.miss` / dữ liệu đang có, không tìm lại):
    - chờ / click ảnh: `best` %, `threshold`, khung chỗ khớp nhất, vùng đã tìm
    - chờ / click / đọc chữ: các dòng OCR đọc được trong vùng (chữ + khung, ≤ ~30 dòng), dòng gần giống nhất
    - web: số ứng viên theo từng cách tìm, đoạn `PageOutline` quanh chỗ cần tìm; phần tử (AX / UIA): ứng viên gần đúng (role, tên, khung)
    - màu: màu thực tế; cửa sổ / ứng dụng: các cửa sổ / tiến trình gần khớp
    - **ngữ cảnh chung** mọi node lỗi: cửa sổ đang ở trước (tiêu đề, tiến trình, khung), vị trí chuột, màn hình (cỡ, DPI)
  - [ ] **Output node thêm cho AI** (đối chiếu từng loại khi làm, thiếu thì bổ sung `Outputs` trong `steptypes.go`): click / gõ / phím – cửa sổ focus sau thao tác; chờ ảnh / chữ – điểm khớp + khung thấy; mở app / URL – cửa sổ mới hiện
  - [ ] **Nhìn màn hình khi AI xin** (chữ trước ảnh): `describe_screen(vùng | cửa sổ)` – dòng OCR kèm khung + danh sách cửa sổ; `screenshot(vùng | cửa sổ)` – ảnh thu nhỏ + tỷ lệ, toạ độ AI gửi về theo pixel ảnh → Go đổi sang toạ độ thật (dùng chung cách đổi của 17 / 25). Ảnh chỉ gửi khi node lỗi hoặc AI xin; `shots: "each"` (tuỳ chọn mỗi lượt) chụp sau mỗi node
  - [ ] **Công tắc AI 4 mức** trên thanh mini + Cài đặt: **Tắt** / **Chỉ soạn** / **Soạn + chạy (hỏi)** – hộp hỏi có *cho phép kịch bản này đến khi đóng app* / **Tự chạy** (không hỏi, bật tay). Mọi mức: giới hạn số node + thời gian mỗi lượt (Cài đặt), F8 + góc màn hình dừng khẩn cấp, thao tác nguy hiểm (`danger.go` – buộc thoát, xoá file…) vẫn hỏi trừ khi cho phép riêng; chuyển sang Tắt huỷ ngay lượt đang chạy; chấm báo khi AI đang điều khiển
  - [ ] Người dùng đụng chuột lúc AI đang chạy (vị trí khác chỗ engine vừa đặt) → tạm dừng, `wait_run` báo AI lý do
  - [ ] Lịch sử chạy (20 F5) ghi lần nào do AI chạy; `autoclicker://guide` thêm mục *vòng thử – sửa* (chạy đến node, đọc detail, sửa, `resume_run`)
  - [ ] Test Go: client MCP in-memory + thiết bị giả – lỗi ở node ảnh → `wait_run` có detail + ảnh → `update_script` → `resume_run` chạy tiếp đúng biến cũ, không chạy lại node trước
  - Xong khi: trên Claude Code nói "mở Notepad, gõ chào bạn, lưu ra Desktop/test.txt" khi app chưa mở → app tự mở thanh mini, AI tạo kịch bản, chạy; cố ý để một node lỗi (vd chờ chữ sai) → AI đọc detail / ảnh, sửa node, `resume_run` chạy tiếp từ node đó (Notepad không mở lại), file được lưu
- [ ] **31c. MCP – AI soạn bằng cách thử + ảnh mẫu** – phần còn lại của 31b / 31c cũ.
  - [ ] Tool `find_image` (so khớp mẫu của 25), `capture_template` (AI khoanh vùng trên ảnh chụp → lưu ảnh mẫu cạnh kịch bản cho `waitImage` / `clickImage`), `get_mouse_position`, `list_windows`
  - [ ] Ghim output node (như pin data n8n): lưu trên máy theo kịch bản, **không ghi vào file**; thẻ node có dấu ghim. App dùng chung `StopAt` / `Only` / `Seed` của 31b: chuột phải node → *Chạy đến đây* / *Chạy cụm đã chọn* / *Chạy tiếp từ đây (dùng kết quả lần trước)*
  - [ ] Chạy thử khô `simulate_script`: chạy flow với thiết bị giả, không đụng chuột – AI kiểm tra nhánh / biến trước khi chạy thật. Node tìm ảnh / chữ / web trả theo kịch bản AI đưa (`found` / `notFound`) hoặc mặc định `success`
  - [ ] `run_step` + chế độ "AI soạn bằng cách thử": các bước AI chạy thành công ghi thành kịch bản nháp, kèm ảnh chụp từng bước để người dùng xem lại
  - [ ] Người dùng thấy khi AI chụp màn hình: chấm đỏ trên thanh mini; tắt được trong Cài đặt
  - Xong khi: AI chụp màn hình YouTube, tự khoanh nút "Phát" làm ảnh mẫu, kịch bản dùng `clickImage` thay vì toạ độ cứng và vẫn chạy khi cửa sổ dời chỗ; nói "mở bài Tinh Hà Say 2 trên YouTube" → AI tự thử, nhạc phát, kịch bản được lưu để lần sau chạy không cần AI
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
| 2 | Đưa Chrome lên trước | `window` (focus) | | 26 |
| 3 | Chờ nút "Nhắn tin" (thử lại 2 lần) | `waitImage` | `wait_s` | 25, 25b |
| 4 | Click nút | `clickImage` | `click_delay` | 25 |
| 5 | Gõ "xin chào {{M1}} và {{M2}}" | `type` | `M1`, `M2` | có sẵn + 21b |
| 6 | Enter | `key` | | có sẵn |
| 7 | Cuộn xuống 5 nấc | `scroll` | | 23 |

> Tự động hoá trên Facebook vi phạm điều khoản của họ và dễ bị khoá tài khoản, nhất là khi nhắn hàng loạt. Thử trên tài khoản phụ hoặc trang của chính bạn.
