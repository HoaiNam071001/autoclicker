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
```

Xem riêng từng màn trong `wails dev` / trình duyệt: thêm hash `#preview`, `#overlay`, `#minibar` (mặc định là MainScreen). Trên trình duyệt không có Go nên nút Chạy bị tắt.

## Cấu trúc

- `main.go`, `app.go` – Wails app; method public của `App` được bind sang frontend (`frontend/wailsjs/go/main/App`, file tự sinh – không sửa tay).
- `model.go` – model kịch bản (`Step`, `Script`); `window.go` – chế độ cửa sổ main / overlay / minibar; `hotkey.go` – phím toàn cục (gohook); `engine.go` – chạy kịch bản qua interface `Input` (`input.go` = robotgo, test dùng input giả – `go test ./...`); `runner.go` – `RunScript` / `StopScript` / `PauseScript` / `ResumeScript`; `failsafe.go` – dừng khẩn cấp khi chuột chạm góc màn hình.
- `frontend/src/scriptStore.ts` – store kịch bản đang soạn (dùng chung các màn); `coordMap.ts` – đổi toạ độ overlay ↔ robotgo (tỷ lệ + độ lệch đo từ chuột); `engineStore.ts` – store trạng thái engine (event `engine:state` / `engine:progress`); `windowMode.ts` – `setWindowMode` đổi cửa sổ + màn; `theme.ts` – theme sáng / tối / theo hệ thống.
- `settings.go` – cấu hình app (theme) trong `os.UserConfigDir()/autoclicker/`.
- `frontend/src/screens/` – các màn UI; `components/TitleBar.tsx` – thanh tiêu đề kiểu Windows 11 + chọn theme; `components/Library.tsx` – thanh bên thư viện; `components/StepList.tsx` – thẻ bước (dnd-kit, kéo ô icon); `components/AddStepPalette.tsx` – bảng thêm bước Ctrl/Cmd+K; `components/StepForm.tsx` – inspector sửa bước; `components/RunSettings.tsx` – lặp / đếm ngược (popup nút "N vòng"); `components/Simulation.tsx` – xem trước bằng con trỏ giả; `components/icons.tsx` – icon SVG; `types.ts` – model `Step` (phải khớp struct Go trong `model.go`; tsc tự kiểm tra).
- `frontend/src/style.css` – Tailwind v4 + token thiết kế trong `@theme`.
- `autoclicker-ui/` – bản thiết kế gốc (Tailwind v3, chỉ để tham khảo, không build) – bị thay bởi design v2 từ bước 27.
- `design/` – design v2: `DESIGN.md` (token, thành phần, icon, map artboard → file), `theme.css` (token Tailwind v4 sáng / tối), `canvas/` (mã artboard, chỉ tham khảo – cần runtime canvas, không chạy riêng).

## Quy ước UI

Theo **design v2** – chi tiết token, thành phần, icon trong `design/DESIGN.md`.

- Tailwind **v4**: token khai báo trong `style.css` (chép từ `design/theme.css` – sửa token thì sửa cả hai), không tạo `tailwind.config.js`. Body không tô nền (overlay / mini bar cần trong suốt) – mỗi màn tự tô `bg-bg` / `bg-surface`.
- Token: `bg` nền, `chrome` thanh tiêu đề / status bar, `side` thanh bên, `surface` thẻ / inspector / popup, `field` ô nhập, `soft` / `hover` nền pha nhẹ, `line` / `line-strong` viền, `ink` / `muted` chữ, `accent` / `accent-soft` / `accent-text` (chỉ cho trạng thái), `warn`, `ok`, `danger*`; `hud*` / `dim` cho overlay (luôn tối).
- Theme: `data-theme="light|dark"` trên `<html>` (`theme.ts`, lưu qua `settings.go`), không đặt = theo hệ điều hành; overlay đặt `data-theme="dark"` riêng. Xem thử trên trình duyệt: `?theme=dark`.
- Gradient chỉ ở 4 chỗ: `btn-primary`, `ring-selected`, `ring-focus`, vòng tiến độ thanh mini. Nút Dừng `btn-danger-soft`.
- Thành phần dùng chung: `NumInput` (+ `fieldBox`), `Segmented`, `Switch` (trong `RunSettings.tsx`), `StepIcon` (`icons.tsx`, đủ icon cả loại bước chưa làm).
- Font `font-sans` (Segoe UI Variable trên Windows → Plus Jakarta Sans) cho chữ, `font-mono` (Cascadia Mono → JetBrains Mono) cho số/toạ độ/thời gian. Font cài local qua `@fontsource`, không dùng Google Fonts.

## Lưu ý nền tảng

- Đang dev trên macOS; mục tiêu phát hành đầu tiên là Windows (robotgo cần cgo + MinGW-w64).
- macOS: cấp quyền Accessibility (System Settings → Privacy & Security) cho terminal / app thì robotgo mới click được.
- Wails v2 chỉ có một cửa sổ: overlay và mini bar là chế độ của cửa sổ chính.
