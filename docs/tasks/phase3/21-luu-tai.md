# 21. Lưu / tải kịch bản

- [x] **21. Lưu / tải kịch bản** – App chưa lưu gì; thêm `version` trước khi đổi model nhiều.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - [x] Model: `Script.version` (`ScriptVersion = 1`), `Script.anchorMode`; file `.json`
  - [x] Go `storage.go`: `SaveScript` / `SaveScriptAs` / `OpenScript` / `OpenScriptPath` qua `runtime.SaveFileDialog` / `OpenFileDialog`; mở file kiểm tra cấu trúc (`decodeScript` + `checkScript`: version, loại bước, số điểm, neo, số âm; bổ sung slice trống, id trống / trùng); mở gần đây (≤ 5, trong `settings.json`, file mất thì bỏ) + bản nháp `draft.json` (`SaveDraft` / `LoadDraft` / `ClearDraft`) trong `os.UserConfigDir()/autoclicker/`; ghi file tạm rồi đổi tên
  - [x] Frontend `scriptFile.ts`: cờ `dirty` + `path` trong store (mọi `updateScript` bật cờ), tự ghi nháp sau 800 ms, mở app khôi phục nháp (không có thì kịch bản trống – app thật không còn dữ liệu mẫu)
  - [x] Thư viện ở thanh bên [`Main`]: ô tìm (lọc cả mở gần đây), chấm `warn` khi chưa lưu, "Mở gần đây" bấm mở được, nút Kịch bản mới / Mở / Lưu
  - [x] Nhãn "Chưa lưu" theo cờ thật, rê tên kịch bản thấy đường dẫn; status bar "Bản nháp tự lưu lúc …", lỗi lưu / mở; chấm `ok` "Sẵn sàng"
  - [x] Phím Cmd/Ctrl + S, + Shift + S, + O, + N; hộp thoại `SaveChangesDialog` (Lưu / Không lưu / Huỷ, Enter / Esc) khi kịch bản mới, mở khác, đóng app (`OnBeforeClose` → event `app:confirmClose`, `ConfirmQuit`)
  - [x] Test `storage_test.go`: ghi rồi đọc lại khớp, file hỏng / thiếu / sai `version` / giá trị lạ báo lỗi, bổ sung mặc định, nháp, mở gần đây, tên file an toàn
  - [x] Kiểm tra tay trên macOS: lưu / lưu thành / mở / mở gần đây, sửa rồi tắt app mở lại còn nháp, đóng app khi chưa lưu hỏi lưu, file hỏng báo lỗi
  - Xong khi: tắt app mở lại vẫn còn bản nháp; mở / lưu qua hộp thoại trên macOS
