# Dev trên Windows

Từ 2026-10-19 người dùng có máy Windows: dev + test phần Windows (17, 18b, 37 E5, OCR Windows, 26 W0) bằng Claude Code chạy **ngay trên máy Windows** – chạy app thật, chuột / màn hình thật. Máy cloud (Linux) chỉ build chéo được, không chạy thử được.

## Cài một lần

1. **Claude Code** bản Windows (desktop app hoặc CLI) + **Git for Windows** (Claude Code dùng Git Bash để chạy lệnh / hook). Không dùng WSL: WSL là Linux, không điều khiển được chuột / cửa sổ Windows.
2. **Go** đúng bản trong `go.mod` (dòng `go 1.xx`), **Node 22**, **Python 3** (lệnh `python3` phải chạy được trong Git Bash – bản cài từ python.org thường chỉ có `python` / `py`: tạo alias hoặc bật *App execution aliases* cho python3).
3. **gcc 64-bit (MinGW-w64)** cho cgo (robotgo / gohook): MSYS2 → `pacman -S mingw-w64-ucrt-x86_64-gcc`, thêm `C:\msys64\ucrt64\bin` vào PATH (hoặc WinLibs qua `winget`). Kiểm tra: `gcc --version`.
4. **Wails CLI** đúng bản trong `go.mod`: `go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0` (`%USERPROFILE%\go\bin` vào PATH) → `wails doctor` không báo thiếu gì (WebView2 có sẵn trên Windows 10/11).
5. Clone repo, `git checkout main`, `cd frontend && npm ci`.

## Lệnh (Git Bash, như trên macOS)

```bash
wails dev                       # chạy app
wails build                     # ra build/bin/ghosthand.exe
go vet ./... && go test ./...   # Go
cd frontend && npm run build && npm test
python3 scripts/check-docs.py
GHOSTHAND_REAL=1 go test ./internal/device/ -run Real   # test thiết bị thật (chuột, chụp, cây phần tử…)
```

## Lưu ý riêng Windows

- **Toạ độ = pixel thật**: app khai báo per-monitor DPI aware (`build/windows/wails.exe.manifest`) – chuột đọc / đặt qua Win32 (`internal/device/cursor_windows.go`), không qua `robotgo.Location` / `Move` (robotgo nhân / chia thêm tỷ lệ DPI – xem bước 17). Code mới đụng toạ độ phải theo `device.Location` / `device.Displays`.
- File riêng hệ điều hành theo hậu tố: `*_darwin.go` (macOS), `*_windows.go`, `*_other.go` (`//go:build !darwin` / `!windows` – báo *chưa hỗ trợ*). Phần Windows mới → thêm `*_windows.go`, sửa build tag của `*_other.go` cho khỏi trùng.
- Đọc chữ: mặc định **Tesseract + `vie`** kèm trong exe (bước 40) – chạy `bash scripts/fetch-tesseract.sh` một lần sau khi clone (cần 7-Zip; tạo `internal/device/tess/tess.zip`, không commit), thiếu thì app tự dùng OCR Windows. OCR Windows (`device/ocr_windows.go`) không có gói tiếng Việt – dòng đọc được mang `Approx`, tìm chữ tự so không dấu + chịu lệch nhẹ (`vision/fuzzy.go`). Gọi WinRT / COM mới → dùng lại `device/winrt_windows.go`.
- Phần tử màn hình Windows (`device/element_windows.go`) qua UI Automation, không cần quyền như macOS; mọi lời gọi COM / WinRT chạy trên một luồng (`onCOM` trong `com_windows.go`) – code COM mới dùng lại.
- Đang chưa có trên Windows: Cửa sổ / Ứng dụng (26 W0 – dùng lại `winlist_windows.go`).
- Test viết theo đường dẫn / hành vi macOS có thể đỏ trên Windows (như `TestGuardPath` đỏ trên Linux): sửa test cho đúng theo hệ điều hành (`runtime.GOOS`), không bỏ test.
- Repo giữ xuống dòng LF (`.gitattributes`) – không đổi sang CRLF (script `.sh` hỏng, gofmt báo khác).
- Hai máy (macOS / cloud và Windows) cùng làm một nhánh: `git pull` trước khi bắt đầu, commit + push khi xong từng phần.
- Bản exe phát hành vẫn do GitHub Actions build (`.github/workflows/build-windows.yml`) – build tại máy chỉ để thử.
