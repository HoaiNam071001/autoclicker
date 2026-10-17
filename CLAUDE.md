# Auto Clicker

Desktop app tự động click theo kịch bản (chọn vùng, chọn điểm, kéo thả A→B, lặp N vòng). Wails v2 + Go + React/TS.

## Bắt đầu mỗi session

1. Đọc **`docs/TASKS.md`**: mục *Trạng thái hiện tại* cho biết đang làm bước nào, *Quyết định đã chốt* là các lựa chọn không cần bàn lại, *Nhật ký* là lịch sử. Thiết kế UI v2 (token, thành phần, icon, artboard) ở **`design/DESIGN.md`** + `design/theme.css`.
2. Xong một bước (hoặc chuyển sang bước mới) → cập nhật `docs/TASKS.md`: tick `[x]`, sửa *Trạng thái hiện tại*, thêm 1 dòng vào *Nhật ký* (ngày tuyệt đối, kết quả kiểm tra).
3. Người dùng trao đổi bằng **tiếng Việt** – trả lời, comment UI và tài liệu bằng tiếng Việt; tên biến/hàm bằng tiếng Anh.

## Lệnh

```bash
wails dev                      # chạy app (hot reload). wails nằm ở ~/go/bin
wails build                    # build ra build/bin/
cd frontend && npm run build   # type-check (tsc) + build frontend – chạy sau mỗi thay đổi TS
cd frontend && npm test        # test đơn vị frontend (vitest, *.test.ts)
```

Xem riêng từng màn trong `wails dev` / trình duyệt: thêm hash `#preview`, `#overlay`, `#minibar` (mặc định là MainScreen). Trên trình duyệt không có Go nên nút Chạy bị tắt.
Màn mẫu giao diện đầy đủ (kể cả tính năng chưa có logic, dữ liệu giả trong `frontend/src/demo/`): `#demo`, `#demo-overlay`, `#demo-minibar`.

## Cấu trúc

- `main.go`, `app.go` – Wails app; method public của `App` được bind sang frontend (`frontend/wailsjs/go/main/App`, file tự sinh – không sửa tay).
- `model.go` – model kịch bản (`Step`, `Script`); `window.go` – chế độ cửa sổ main / overlay / minibar; `hotkey.go` – phím toàn cục (gohook); `engine.go` – chạy kịch bản qua interface `Input` (`input.go` = robotgo, test dùng input giả – `go test ./...`); `runner.go` – `RunScript` / `StopScript` / `PauseScript` / `ResumeScript`; `failsafe.go` – dừng khẩn cấp khi chuột chạm góc màn hình. `open.go` – bước Mở URL / app (`openCommand` theo hệ điều hành). `match.go` – tìm ảnh mẫu (NCC, Go thuần); `capture.go` / `capture_darwin.go` – chụp màn hình qua interface `Screen`, quyền Screen Recording macOS; `imagestep.go` – engine chờ ảnh / click theo ảnh / chờ pixel; `imageapi.go` – `CaptureTemplate`, `FindImage`, `PixelColorAt` cho giao diện (ẩn cửa sổ lúc chụp). `ocr.go` / `ocr_darwin.go` – tìm theo chữ (đọc chữ bằng Vision, so chữ `findText`; Windows chưa có). `steptypes.go` – mô tả + ví dụ JSON từng loại bước (cho AI / MCP; thêm loại bước mới phải thêm dòng, có test). `inputs.go` – input của kịch bản (`ScriptInput`, `resolveInputs` áp giá trị vào bản sao trước khi chạy, `ScriptError` lỗi có cấu trúc, giá trị lần trước / preset trong `inputs.json`). `anchor.go` – neo điểm theo cạnh vùng (`resolvePoint`; `frontend/src/anchor.ts` phải cùng công thức). `frontend/src/transform.ts` – ma trận biến đổi nhiều điểm (bước 30, có test); `components/TransformTool.tsx` – công cụ Chọn nhiều trên overlay (khung biến đổi, thanh số, hộp ánh xạ vùng).
- `frontend/src/scriptStore.ts` – store kịch bản đang soạn (dùng chung các màn; hoàn tác / làm lại, chọn nhiều, `mapPoints`); `stepTree.ts` – thao tác cây bước (khối Nếu / Nhóm: `flatten`, `mapStep`, `moveTo`, `rows` / `dropLoc` cho kéo thả; có test) – mọi chỗ duyệt bước phải đi qua đây, không duyệt `script.steps` phẳng; `blocks.go` – engine khối (Go); `scriptInputs.ts` – ghim trường / sửa input / bước dùng input (bước 21b, có test); `runForm.ts` – bấm Chạy / F8 → form Chạy khi có input, preset, chạy luôn; `imageSteps.ts` – khoanh ảnh mẫu (overlay tool `capture`), thử tìm, lấy màu, quyền chụp màn hình; `components/MockPage.tsx` – trang mẫu phủ nền khi Xem trước cuộn; `coordMap.ts` – đổi toạ độ overlay ↔ robotgo (tỷ lệ + độ lệch đo từ chuột); `engineStore.ts` – store trạng thái engine (event `engine:state` / `engine:progress`); `windowMode.ts` – `setWindowMode` đổi cửa sổ + màn; `theme.ts` – theme sáng / tối / theo hệ thống.
- `settings.go` – cấu hình app (theme, mở gần đây) trong `os.UserConfigDir()/autoclicker/`; `storage.go` – lưu / mở file kịch bản `.json` (`ScriptVersion`, `checkScript`), bản nháp `draft.json`, hỏi lưu khi đóng app. `frontend/src/scriptFile.ts` – Lưu / Mở / Mới / mở gần đây, tự ghi nháp, cờ chưa lưu; `components/SaveChangesDialog.tsx` – hỏi Lưu / Không lưu / Huỷ.
- `frontend/src/screens/` – các màn UI; `components/TitleBar.tsx` – thanh tiêu đề kiểu Windows 11 + chọn theme; `components/Library.tsx` – thanh bên thư viện; `components/StepList.tsx` – thẻ bước (dnd-kit, kéo ô icon); `components/AddStepPalette.tsx` – bảng thêm bước Ctrl/Cmd+K; `components/StepForm.tsx` – inspector sửa bước; `components/RunSettings.tsx` – lặp / đếm ngược (popup nút "N vòng"); `components/Simulation.tsx` – xem trước bằng con trỏ giả; `components/icons.tsx` – icon SVG; `types.ts` – model `Step` (phải khớp struct Go trong `model.go`; tsc tự kiểm tra).
- `frontend/src/style.css` – Tailwind v4 + token thiết kế trong `@theme`.
- `components/inputs/` – ghim (`PinMenu`), ô chữ `{{biến}}`, tab Input, form Chạy (đã nối logic ở 21b). Giao diện dựng sẵn cho tính năng chưa có logic (bước 27b, app thật hiện tắt + *Sắp có*): `components/steps/` (inspector loại bước mới, hết giờ, nhóm, neo ±), `ScheduleDialog.tsx`, `AppChrome.tsx` (`Toolbar`, `StatusBar`, nút debug / ghi). Kiểu dữ liệu UI tạm trong `steps/model.ts` – chuyển vào `types.ts` + `model.go` khi làm logic bước đó.
- `autoclicker-ui/` – bản thiết kế gốc (Tailwind v3, chỉ để tham khảo, không build) – bị thay bởi design v2 từ bước 27.
- `design/` – design v2: `DESIGN.md` (token, thành phần, icon, map artboard → file), `theme.css` (token Tailwind v4 sáng / tối), `canvas/` (mã artboard, chỉ tham khảo – cần runtime canvas, không chạy riêng).

## Quy ước UI

Theo **design v2** – chi tiết token, thành phần, icon trong `design/DESIGN.md`.

- Tailwind **v4**: token khai báo trong `style.css` (chép từ `design/theme.css` – sửa token thì sửa cả hai), không tạo `tailwind.config.js`. Body không tô nền (overlay / mini bar cần trong suốt) – mỗi màn tự tô `bg-bg` / `bg-surface`.
- Token: `bg` nền, `chrome` thanh tiêu đề / status bar, `side` thanh bên, `surface` thẻ / inspector / popup, `field` ô nhập, `soft` / `hover` nền pha nhẹ, `line` / `line-strong` viền, `ink` / `muted` chữ, `accent` / `accent-soft` / `accent-text` (chỉ cho trạng thái), `warn`, `ok`, `danger*`; `hud*` / `dim` cho overlay (luôn tối).
- Theme tối dùng bảng **Slate** trong `design/newdarkmode.md` (lớp càng cao càng sáng: `bg` → `surface` → `field` / `soft` → `raised`). Chữ trên nền `accent` dùng `text-on-accent`, không dùng `text-white`.
- Theme: `data-theme="light|dark"` trên `<html>` (`theme.ts`, lưu qua `settings.go`), không đặt = theo hệ điều hành; overlay đặt `data-theme="dark"` riêng. Xem thử trên trình duyệt: `?theme=dark`.
- Gradient chỉ ở 4 chỗ: `btn-primary`, `ring-selected`, `ring-focus`, vòng tiến độ thanh mini. Nút Dừng `btn-danger-soft`.
- Thành phần dùng chung: `NumInput` (+ `fieldBox`), `Segmented`, `Switch` (trong `RunSettings.tsx`), `StepIcon` (`icons.tsx`, đủ icon cả loại bước chưa làm).
- Font `font-sans` (Segoe UI Variable trên Windows → Plus Jakarta Sans) cho chữ, `font-mono` (Cascadia Mono → JetBrains Mono) cho số/toạ độ/thời gian. Font cài local qua `@fontsource`, không dùng Google Fonts.

## Lưu ý nền tảng

- Đang dev trên macOS; mục tiêu phát hành đầu tiên là Windows (robotgo cần cgo + MinGW-w64).
- macOS: cấp quyền Accessibility (System Settings → Privacy & Security) cho terminal / app thì robotgo mới click được.
- Wails v2 chỉ có một cửa sổ: overlay và mini bar là chế độ của cửa sổ chính.
