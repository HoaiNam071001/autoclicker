---
name: senior-review
description: Tự review diff hiện tại theo chuẩn senior của dự án trước khi commit (hoặc khi người dùng nói "review", "rà code", "kiểm tra chất lượng"). Đúng / an toàn trước, rồi kiến trúc, đơn giản, test, UI. Sửa luôn chỗ rõ ràng, hỏi chỗ cần người dùng chọn.
---

# Review kiểu senior

Phạm vi: `git diff` (+ `git diff --cached`, file mới chưa track). Đọc diff trước, chỉ mở file đầy đủ khi cần ngữ cảnh. Báo theo mức: **phải sửa** → **nên sửa** → **gợi ý**; mỗi ý có `file:dòng` + tình huống cụ thể gây lỗi. Không báo ý kiểu "có thể" không có kịch bản lỗi.

## 1. Đúng & an toàn (phải sửa)

- Nhánh lỗi: lỗi bị nuốt, `err` không kiểm, `nil` map / con trỏ (`Retry`, `Cond`, `Region`, `Pos` đều có thể nil), slice rỗng vs nil gửi sang JSON.
- Dừng / tạm dừng: mọi chờ dùng `sleep(ctx, …)` / `gate.wait`, vòng lặp kiểm `ctx`; goroutine có đường thoát; khoá `a.mu` không giữ khi gọi Wails runtime chặn / emit lâu.
- Engine: node tắt, thử lại, cổng không nối, vòng lặp vô hạn, biến chưa có, kết quả đặt lại mỗi vòng.
- Toạ độ: neo, màn phụ toạ độ âm, DPI (Windows), vùng `nil`.
- File / đường dẫn từ frontend: chặn đường dẫn lạ (`safeName`), tên file Windows (`safeFileName`).
- Đa nền tảng: code chỉ macOS phải có build tag + bản `_other.go`; không gọi API chỉ có trên macOS từ code chung.
- Go ↔ TS: struct `model.go` khớp `types.ts` (đã `wails generate module`?); công thức đôi (`portsOf`, `flowOrder`, neo) sửa cả hai + test chung `testdata/`.

## 2. Kiến trúc (nên sửa)

- Đúng lớp: logic thuần không import React / Wails / store; component không dựng `script` mới, không gọi `setEditor`; sửa kịch bản qua `updateScript`.
- Flow: thao tác cạnh / node qua `flow/edit.ts`, thứ tự qua `flowOrder`, tham chiếu bằng `id`, tên / chi tiết node qua `flow/view.ts`.
- Một file một việc (`CLAUDE.md` › *Một file một việc*): diff đưa file vượt ~300 dòng, hoặc thêm tính năng thứ hai vào một file / component → **phải sửa**: tách khối giao diện thành component chỉ hiển thị, phần tính thành hàm thuần có test. Hàm > ~60 dòng nhiều nhánh → tách có tên.
- Dùng lại thành phần sẵn có (`CLAUDE.md` › *Dùng lại trước khi viết mới*); hai chỗ giống nhau → gom.
- Lỗi cho người dùng / AI có cấu trúc (`ScriptError`, `StepError` + mã), chữ tiếng Việt nói cách khắc phục.
- Không code chết, không cờ tương thích ngược, không TODO vô chủ (cầu nối tạm phải ghi bước sẽ xoá trong `docs/tasks/`).

## 3. Đơn giản & hiệu năng (nên sửa)

- Trừu tượng thừa (interface một bản cài, tham số không dùng), state trùng lặp suy ra được.
- Re-render: selector store trả object mới mỗi lần, effect thiếu dọn dẹp (EventsOn trả hàm huỷ), listener đăng ký nhiều lần.
- Chụp màn hình / OCR trong vòng lặp: giới hạn vùng, chụp một lần mỗi lần xét (`snapshot`), có `intervalMs`.

## 4. Test (phải có cho logic mới)

- Logic Go mới / sửa → test bằng thiết bị giả; logic thuần TS → `*.test.ts`. Có ca biên: rỗng, nil, hết giờ, bị huỷ, giá trị ngoài miền.
- Sửa bug → có test đỏ trước khi sửa.

## 5. UI & văn bản (gợi ý / nên sửa)

- Token `design/DESIGN.md` (không màu cứng, `text-on-accent` thay `text-white`, gradient chỉ 4 chỗ), xem cả theme sáng + tối, chữ số dùng `font-mono`.
- Trạng thái rỗng / đang tải / lỗi / bị tắt (*Sắp có*), phím tắt không đụng nhau, tiêu điểm bàn phím.
- Comment tiếng Việt, ngắn, nói *vì sao*; tên biến / hàm tiếng Anh; không comment lại điều code đã nói.

## 6. Kết thúc

Chạy `go vet ./... && go test ./...`, `cd frontend && npm run build && npm test`, `python3 scripts/check-docs.py`. Sửa các ý *phải sửa* / *nên sửa* rõ ràng; ý cần đánh đổi → liệt kê cho người dùng chọn. Báo tóm tắt: đã sửa gì, còn gì.
