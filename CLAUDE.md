# Auto Clicker

Desktop app tự động click theo kịch bản (chọn vùng, chọn điểm, kéo thả A→B, lặp N vòng). Wails v2 + Go + React/TS.

## Bắt đầu mỗi session

1. Đọc **`docs/TASKS.md`**: mục *Trạng thái hiện tại* cho biết đang làm bước nào, *Quyết định đã chốt* là các lựa chọn không cần bàn lại, *Nhật ký* là lịch sử.
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
- `frontend/src/scriptStore.ts` – store kịch bản đang soạn (dùng chung các màn); `coordMap.ts` – đổi toạ độ overlay ↔ robotgo (tỷ lệ + độ lệch đo từ chuột); `engineStore.ts` – store trạng thái engine (event `engine:state` / `engine:progress`); `windowMode.ts` – `setWindowMode` đổi cửa sổ + màn.
- `frontend/src/screens/` – các màn UI; `components/StepList.tsx` – danh sách bước (dnd-kit); `components/StepForm.tsx` – form sửa bước; `components/RunSettings.tsx` – lặp / đếm ngược; `components/icons.tsx` – icon SVG; `types.ts` – model `Step` (phải khớp struct Go trong `model.go`; tsc tự kiểm tra).
- `frontend/src/style.css` – Tailwind v4 + token thiết kế trong `@theme`.
- `autoclicker-ui/` – bản thiết kế gốc (Tailwind v3, chỉ để tham khảo, không build).

## Quy ước UI

- Tailwind **v4**: thêm màu/font mới vào `@theme` trong `style.css`, không tạo `tailwind.config.js`.
- Token: `ground` nền, `ink` chữ, `muted`, `line` viền; `accent` (xanh) cho điểm click; `drag` (cam) cho kéo thả; `danger`, `ok`; nhóm `dark-*` cho màn nền tối và mini bar.
- Font `font-sans` (Be Vietnam Pro) cho chữ, `font-mono` (IBM Plex Mono) cho số/toạ độ/thời gian. Font cài local qua `@fontsource`, không dùng Google Fonts.

## Lưu ý nền tảng

- Đang dev trên macOS; mục tiêu phát hành đầu tiên là Windows (robotgo cần cgo + MinGW-w64).
- macOS: cấp quyền Accessibility (System Settings → Privacy & Security) cho terminal / app thì robotgo mới click được.
- Wails v2 chỉ có một cửa sổ: overlay và mini bar là chế độ của cửa sổ chính.
