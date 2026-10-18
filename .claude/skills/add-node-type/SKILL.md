---
name: add-node-type
description: Thêm một loại node mới vào flow (vd focus cửa sổ, web DOM, subflow) hoặc thêm trường cho node có sẵn. Danh sách đủ mọi chỗ phải sửa ở Go + TS, theo thứ tự.
---

# Thêm loại node

Đọc trước: `docs/ARCHITECTURE.md` §2 (model, cổng), §4 (engine). Lấy một loại node gần giống nhất làm mẫu và `grep` tên hằng của nó (vd `StepWaitPixel`, `'waitPixel'`) để thấy đủ chỗ đang dùng – danh sách dưới là tối thiểu.

## Go

- [ ] `internal/model/model.go`: hằng `StepXxx`, thêm vào `StepTypes`, `IsControl` nếu là node điều khiển, trường mới của `Step` (json camelCase, `omitempty` cho trường riêng loại), mặc định trong `NewStep`.
- [ ] `~/go/bin/wails generate module`.
- [ ] `internal/script/check.go` `checkStep` (node điều khiển: `control.go` `checkControl`): miền giá trị hợp lệ khi mở file.
- [ ] `internal/script/validate.go` `ValidateScript` (điều khiển: `control.go` `validateControl`): dữ liệu đủ để chạy (lỗi `stepErr(i, s, "<field>", …)` – trường lỗi để UI / AI chỉ đúng chỗ); `script/flowcheck.go` nếu node trỏ node khác.
- [ ] Chạy: hành động → `runStep` / `doStep` (`internal/engine/actions.go`) hoặc file riêng nếu lớn (như `engine/imagestep.go`); điều khiển → `runControl` (`engine/nodes.go`). Trả `output` map. Lỗi không thấy → `notFoundErr`; dữ liệu sai lúc chạy → `invalidErr`. Tôn trọng `ctx` (dùng `Sleep(ctx, …)` của engine, không `time.Sleep`).
- [ ] Thiết bị / hệ điều hành mới → interface + bản thật (`_darwin.go` / `_windows.go` / `_other.go`) + bản giả trong test.
- [ ] `internal/model/flow.go` `PortsOf` nếu cổng khác `success` / `fail`.
- [ ] `internal/script/varfields.go` `VarFields`: ô chữ nhận `{{biến}}`.
- [ ] `internal/model/steptypes.go`: mô tả, output, ví dụ JSON (test kiểm đủ loại).
- [ ] Test: `nodes_test.go` / `flowrun_test.go` / file test riêng – chạy đúng, cổng đúng, output đúng, lỗi đúng mã, dừng được giữa chừng.
- [ ] **Kịch bản mẫu** (bước 33, bắt buộc – `TestEveryStepTypeHasSample` đỏ nếu thiếu): ≥ 1 mẫu trong `internal/samples/scripts/<id>.json` (JSON gọn `{meta, script}`, node chỉ ghi trường cần) chạy trên trang sân thử `internal/samples/pages/` – cần trang mới thì thêm trang (tự hiện kết quả, chữ to cho OCR, không chép chữ đích vào nhật ký trang) + dòng trong `pages/index.html`. Chạy thật: `AUTOCLICKER_REAL=1 AUTOCLICKER_SAMPLES=<id> go test ./internal/samples/ -run TestRealSamples -v`.

## Frontend

- [ ] `types.ts`: union loại, trường, nhãn (`STEP_LABEL`…); `components/icons.tsx` icon (có thể đã có từ bước 27b).
- [ ] `flow/model.ts` `portsOf` (cùng công thức Go) + test.
- [ ] `flow/view.ts`: dòng chi tiết, nhãn cổng, `outputsOf` + test.
- [ ] `scriptInputs.ts` `mapVarText`: ô chữ nhận biến (khớp `script.VarFields`).
- [ ] `store/nodes.ts` `addStep` hoặc `store/control.ts` `addControl`.
- [ ] Bảng thêm node (`AddStepPalette`): mục mới, nhóm đúng, từ khoá tìm không dấu.
- [ ] Inspector: hành động → `components/stepform/<Nhóm>Fields.tsx` (+ đăng ký trong `StepForm.tsx`); điều khiển → `components/flow/nodes/`. Dùng `Field` / `NumField` / `PickButton` / `VarTextInput`.
- [ ] Node cần điểm trên màn hình → overlay (`store/points`, `OverlayScreen`), `Simulation` (xem trước).
- [ ] Bỏ phần *Sắp có* tương ứng trong `components/steps/` nếu có.

## Xong

`python3 scripts/check-docs.py` (bắt thiếu nhãn / icon / mô tả), `go test ./...`, `npm run build && npm test`, thêm ví dụ vào `examples/` nếu hữu ích, rồi `wrap-up`. Thêm dòng loại node vào bảng cổng ở `docs/ARCHITECTURE.md` §2 nếu cổng mới.
