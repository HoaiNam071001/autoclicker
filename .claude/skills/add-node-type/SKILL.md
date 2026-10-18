---
name: add-node-type
description: Thêm một loại node mới vào flow (vd focus cửa sổ, web DOM, subflow) hoặc thêm trường cho node có sẵn. Danh sách đủ mọi chỗ phải sửa ở Go + TS, theo thứ tự.
---

# Thêm loại node

Đọc trước: `docs/ARCHITECTURE.md` §2 (model, cổng), §4 (engine). Lấy một loại node gần giống nhất làm mẫu và `grep` tên hằng của nó (vd `StepWaitPixel`, `'waitPixel'`) để thấy đủ chỗ đang dùng – danh sách dưới là tối thiểu.

## Go

- [ ] `model.go`: hằng `StepXxx`, thêm vào `stepTypes`, `isControl` nếu là node điều khiển, trường mới của `Step` (json camelCase, `omitempty` cho trường riêng loại), mặc định trong `NewStep`.
- [ ] `~/go/bin/wails generate module`.
- [ ] `storage.go` `checkStep`: miền giá trị hợp lệ khi mở file.
- [ ] `engine.go` `validateScript`: dữ liệu đủ để chạy (lỗi `stepErr(i, s, "<field>", …)` – trường lỗi để UI / AI chỉ đúng chỗ); `flowcheck.go` nếu node trỏ node khác.
- [ ] Chạy: hành động → `runStep` / `doStep` (engine.go) hoặc file riêng nếu lớn (như `imagestep.go`); điều khiển → `runControl` (nodes.go). Trả `output` map. Lỗi không thấy → `notFoundErr`; dữ liệu sai lúc chạy → `invalidErr`. Tôn trọng `ctx` (dùng `sleep(ctx, …)`, không `time.Sleep`).
- [ ] Thiết bị / hệ điều hành mới → interface + bản thật (`_darwin.go` / `_windows.go` / `_other.go`) + bản giả trong test.
- [ ] `flow.go` `portsOf` nếu cổng khác `success` / `fail`.
- [ ] `vars.go` `varFields`: ô chữ nhận `{{biến}}`.
- [ ] `steptypes.go`: mô tả, output, ví dụ JSON (test kiểm đủ loại).
- [ ] Test: `nodes_test.go` / `flowrun_test.go` / file test riêng – chạy đúng, cổng đúng, output đúng, lỗi đúng mã, dừng được giữa chừng.

## Frontend

- [ ] `types.ts`: union loại, trường, nhãn (`STEP_LABEL`…); `components/icons.tsx` icon (có thể đã có từ bước 27b).
- [ ] `flow/model.ts` `portsOf` (cùng công thức Go) + test.
- [ ] `flow/view.ts`: dòng chi tiết, nhãn cổng, `outputsOf` + test.
- [ ] `scriptInputs.ts` `mapVarText`: ô chữ nhận biến (khớp `varFields`).
- [ ] `store/nodes.ts` `addStep` hoặc `store/control.ts` `addControl`.
- [ ] Bảng thêm node (`AddStepPalette`): mục mới, nhóm đúng, từ khoá tìm không dấu.
- [ ] Inspector: hành động → `components/stepform/<Nhóm>Fields.tsx` (+ đăng ký trong `StepForm.tsx`); điều khiển → `components/flow/nodes/`. Dùng `Field` / `NumField` / `PickButton` / `VarTextInput`.
- [ ] Node cần điểm trên màn hình → overlay (`store/points`, `OverlayScreen`), `Simulation` (xem trước).
- [ ] Bỏ phần *Sắp có* tương ứng trong `components/steps/` nếu có.

## Xong

`python3 scripts/check-docs.py` (bắt thiếu nhãn / icon / mô tả), `go test ./...`, `npm run build && npm test`, thêm ví dụ vào `examples/` nếu hữu ích, rồi `wrap-up`. Thêm dòng loại node vào bảng cổng ở `docs/ARCHITECTURE.md` §2 nếu cổng mới.
