# 21b. Input của kịch bản

- [x] **21b. Input của kịch bản** – Người soạn đánh dấu giá trị nào được đổi khi chạy; bấm Chạy thì điền form. Go áp giá trị vào bản sao kịch bản rồi mới đưa engine, **engine không đổi**. Cùng đổi `version` với 21; CSV (28) dựa vào đây.
  - UI khung có từ 27b – bước này nối logic, bỏ nhãn *Sắp có*
  - Đã làm (2026-10-17): `inputs.go` (`ScriptInput`, `resolveInputs`, `ScriptError`, `GetInputMemory` / `SaveInputMemory` → `inputs.json`), file kịch bản `version` 2; `RunScript` / `ValidateScript` nhận `values`. Frontend: `scriptInputs.ts` (ghim / dùng chung / bỏ ghim, sửa input, bước dùng input, nhóm form Chạy), `runForm.ts` (Chạy / F8 → form, chạy luôn, preset), `PinMenu`, `NumField` + tab Input trong `StepForm.tsx`
  - Khác plan: kiểu `point` (ghim toạ độ) chưa làm – nút ghim cạnh toạ độ vẫn *Sắp có*; ghim `enabled` (bool) có ở Go nhưng chưa có nút trên UI. Ghim vào bước đã xoá được giữ lại (hoàn tác xoá bước thì ghim quay lại), lúc chạy bỏ qua. Trường đã ghim luôn bằng mặc định của input nên Xem trước dùng mặc định, không gọi resolve
  - [x] Model `Script.inputs`: `key`, `label`, `type` (`number`, `ms`, `text`, `multiline`, `bool`, `select`, `point`), `default`, `required`, `min` / `max`, `options`, `hint`, `targets: [{stepId, field}]`. Ví dụ:
    ```json
    "inputs": [
      { "key": "click_delay", "label": "Chờ sau khi click", "type": "ms", "default": 500,
        "targets": [{ "stepId": "s4", "field": "delayAfter" }] },
      { "key": "M1", "label": "Tên người nhận", "type": "text", "required": true, "hint": "Tên hiện trong lời chào" }
    ]
    ```
    - Nguồn: **trường được ghim** (`targets`, một input cho nhiều bước) và **biến `{{key}}`** trong trường chữ (gõ chữ, URL, tên cửa sổ) – tự tạo input khi gõ biến mới
    - Key: chữ không dấu, số, `_`; nhãn tiếng Việt ở `label`; `\{{` để gõ ra `{{`; không ghim được trường tự tính (kéo thả "từ điểm trước")
  - [x] Go: `resolveInputs(script, values) (Script, error)` – sao chép, ghi trường ghim, thay `{{…}}`, `validateScript`; `RunScript(sc, values, countdown)`; Xem trước dùng chung resolve
  - [x] Lỗi rõ: thiếu bắt buộc, sai kiểu / ngoài khoảng, biến không có trong `inputs`, input không bước nào dùng, toạ độ ngoài vùng
  - [x] Giá trị lần trước + preset lưu trong `UserConfigDir` theo từng kịch bản, **không** ghi vào file kịch bản
  - [x] Inspector [`StepInputs`]: `components/PinButton.tsx` (32×32 cạnh mỗi trường số / chữ; đã ghim nền `accent-soft`); trường đã ghim `field-pinned` + chip key `var-token` cạnh nhãn
  - [x] Ô chữ có biến [`StepInputs`]: `components/VarTextInput.tsx` – tô `{{biến}}` bằng `var-token`; gõ `{{` mở popup gợi ý (key mono `accent-text` + nhãn `muted`, "+ Tạo input mới…"); dòng phụ "N biến, tự thêm vào tab Input"
  - [x] Tab Input [`InputsTab`]: `components/InputsPanel.tsx` – chip số input ở tab; mỗi dòng ⋮⋮ kéo, nhãn + `*`, `key · mặc định` mono, chip kiểu, chip bước dùng; input đang sửa mở rộng (`ring-selected`): nhãn, kiểu, mặc định, bắt buộc, gợi ý; bỏ ghim; "+ Thêm input"
  - [x] Thẻ bước: chip viền `line` "N input"
  - [x] Form Chạy [`RunDialog`, `RunDialogDark`]: `components/RunDialog.tsx` (popup 560, bo 12) – đầu có tên kịch bản + số input; "Bộ giá trị" (select preset + "Lưu bộ mới"); ô nhóm theo bước, ô đã đổi `ring-focus` + "đã đổi, mặc định … · khôi phục", `*` màu `warn`; khung "Sẽ gõ" xem trước câu sau khi thay biến; chân: "Lần sau chạy luôn với giá trị này", "Đếm ngược 3 giây", Huỷ / Chạy (`btn-primary`, F8)
  - [x] Nút Chạy / F8 mở form khi kịch bản có input; chạy thẳng nếu đã bật "Lần sau chạy luôn"; thanh mini hiện tên preset
  - [x] Chuẩn bị cho MCP (31): `resolveInputs` và kiểm tra input viết ở Go (không ở frontend) để CLI / MCP gọi được; `validateScript` trả lỗi có cấu trúc `{stepId, field, message}` thay vì chỉ chuỗi
  - [x] Test: ghi đè trường; thay biến (có tiếng Việt); thiếu biến; sai kiểu; escape; một input nhiều bước; input trong nhóm (sau bước 20) – Go `inputs_test.go` (trừ input trong nhóm – làm ở 20); vitest `scriptInputs.test.ts` (ghim, dùng chung, đổi key, nhóm form Chạy)
  - [x] Kiểm tra tay trong app
  - Xong khi: B1 click 500 → chạy với 1000 ms; B2 "xin chào {{M1}} và {{M2}}" gõ đúng giá trị điền
