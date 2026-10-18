---
name: fix-bug
description: Sửa lỗi người dùng báo (vd "bị lỗi", "không chạy", "click lệch", "bấm không ăn", kiểm tra tay thất bại trong checks.md). Tìm nguyên nhân gốc, viết test tái hiện trước, sửa, rà chỗ cùng kiểu lỗi, ghi lại.
---

# Sửa bug

## 1. Hiểu đúng lỗi

- Ghi lại: làm gì → mong gì → thấy gì; màn nào (main / overlay / minibar), hệ điều hành, 1 hay nhiều màn hình, theme. Thiếu thông tin quyết định hướng sửa → hỏi người dùng (một lần, gom câu hỏi) – còn lại tự tra.
- Lỗi chạy kịch bản → xin / đọc file kịch bản (`.json`) và lịch sử chạy (`os.UserConfigDir()/autoclicker/runs/…`: `StepResult` có `error.code`, `logs`, ảnh lúc lỗi).
- Định vị bằng `docs/ARCHITECTURE.md` (§3 soạn, §4 chạy, §6 bảng chỗ sửa) rồi mới mở code.

## 2. Tìm nguyên nhân gốc, không vá triệu chứng

- Lần ngược luồng dữ liệu tới chỗ giá trị sai sinh ra. Nghi vấn hay gặp:
  - Toạ độ lệch: neo (`anchor.go` ↔ `anchor.ts`), `coordMap.ts` (DPI, màn phụ, menu bar macOS), vùng `nil` = màn chính.
  - Go và TS tính khác nhau: `portsOf`, `flowOrder`, neo – hai bên phải cùng công thức.
  - State giữa các màn: màn không cùng mount, state phải ở store module.
  - Đồng thời: `a.mu`, `a.run`, goroutine hotkey / failsafe, `ctx` không được tôn trọng (dừng không ăn).
  - Thứ tự node dựa vào mảng `steps` thay vì cạnh.
- Không chắc → thêm log / test nhỏ để chứng minh, không đoán. Nói rõ nguyên nhân cho người dùng bằng 1–2 câu.

## 3. Test tái hiện trước khi sửa

- Go: test với `Input` / `Screen` / `TextReader` giả. TS: logic thuần `*.test.ts`. Test phải **đỏ** trước khi sửa, **xanh** sau.
- Lỗi chỉ thấy với thiết bị thật (quyền, DPI, chuột thật) → không test được thì ghi mục kiểm tra tay vào `checks.md`.

## 4. Sửa

- Sửa nhỏ nhất đúng chỗ gốc; theo quy tắc lớp trong `CLAUDE.md`. Không kèm refactor không liên quan (ghi thành mục con riêng nếu thấy cần).
- `grep` cùng mẫu lỗi ở chỗ khác (bản Go / TS song song, loại node tương tự) và sửa luôn.
- Lỗi cho người dùng phải nói được cách khắc phục (vd thiếu quyền → chỉ chỗ cấp quyền).

## 5. Kiểm tra + ghi lại

- `go vet ./... && go test ./...`, `cd frontend && npm run build && npm test`.
- File bước liên quan: thêm mục con `[x] Sửa: <triệu chứng> – <nguyên nhân>`; lỗi từ `checks.md` thì để mục kiểm tra đó `[ ]` cho người dùng thử lại.
- `log.md` một dòng (ngày tuyệt đối, nguyên nhân, kết quả test). Commit `fix: …` qua skill `wrap-up`.
