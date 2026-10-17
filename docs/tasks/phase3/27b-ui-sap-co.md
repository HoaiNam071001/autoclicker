# 27b. Dựng sẵn UI cho tính năng chưa có

- [x] **27b. Dựng sẵn UI cho tính năng chưa có** *(người dùng yêu cầu, 2026-10-16)* – Làm hết giao diện theo design v2 trước, logic làm ở bước sau.
  - Cách làm: component nhận dữ liệu qua props, chưa nối Go. Trong app: hiện đúng chỗ nhưng nút / ô bị tắt, nhãn *Sắp có* (tooltip ghi số bước sẽ làm). Màn **`#demo`** nạp kịch bản mẫu (cả loại bước chưa có, dữ liệu giả trong `frontend/src/demo/`, **không** thêm vào `types.ts` / `model.go`) để xem đủ trạng thái, so với canvas
  - [x] Thư viện thanh bên [`Main`] – `components/Library.tsx`: ô tìm, danh sách kịch bản, mục đang mở (`bg-hover` + vạch trái accent), chấm `warn` chưa lưu, "Mở gần đây", nút Kịch bản mới / Mở / Lưu, mục "Lịch sử"; nhãn "Chưa lưu" cạnh tên kịch bản; status bar "Bản nháp tự lưu lúc …" *(logic: 21, 25c)*
  - [x] Input [`StepInputs`, `InputsTab`]: `PinButton.tsx`, trường đã ghim `field-pinned` + chip key; `VarTextInput.tsx` tô `{{biến}}` + popup gợi ý; `InputsPanel.tsx` (tab Input, dòng input, dạng mở rộng để sửa); chip "N input" trên thẻ bước *(logic: 21b)*
  - [x] Form Chạy [`RunDialog`, `RunDialogDark`] – `RunDialog.tsx`: preset, ô theo bước, ô đã đổi + khôi phục, khung "Sẽ gõ", chân với "Lần sau chạy luôn" / Đếm ngược / Huỷ / Chạy *(logic: 21b)*
  - [x] Neo điểm: ô X / Y dạng ±, lưới 3×3 chọn neo + tuỳ chọn %, toạ độ `(−300, +300)` trong danh sách, đường neo trên overlay *(chưa có trong design – vẽ theo thành phần v2; logic: 29)*
  - [x] Sửa nhiều điểm trên overlay: khung biến đổi + núm (co giãn / xoay / xiên), thanh công cụ ΔX / ΔY / % / góc / xiên, nút căn chỉnh / lật, hộp hỏi "Ánh xạ sang vùng mới" *(chưa có trong design; logic: 30)*
  - [x] Inspector cho loại bước mới: cuộn (hướng, số nấc), mở URL / app, **Chờ ảnh** [`Main`] (ảnh mẫu, Chụp lại, Xem trên màn hình, độ khớp, chờ tối đa, kiểm tra mỗi), click theo ảnh (ΔX / ΔY), chờ pixel (điểm, màu, sai số), focus cửa sổ (ô nhập + danh sách cửa sổ) *(logic: 23, 25, 26)*
  - [x] Khi hết giờ [`Main`]: segmented Thử lại / Bỏ qua / Dừng, số lần thử, "Thử hết vẫn không thấy"; chip trên thẻ bước *(logic: 25b)*
  - [x] Overlay tool `capture` [`Overlay`]: khoanh ảnh mẫu, nhãn kích thước, "Dùng ảnh này" / "Khoanh lại" *(logic: 25)*
  - [x] Thanh mini trạng thái chờ [`MiniBar`]: vòng theo thời gian chờ, "đang tìm x / y s", "thử lại x / y", tên preset *(logic: 25, 25b, 21b)*
  - [x] Khung nhóm bước [`Main`]: nét đứt `border-line-strong bg-group rounded-group`, icon + tên + chip `lặp ×N` *(logic: 20)*
  - [x] Debug: nút Chạy từ bước này / Tiếp (F10), breakpoint trên thẻ bước, màn timeline lịch sử chạy *(chưa có trong design; logic: 25c → nay là 20 F5)*
  - [x] Hẹn giờ (chip giờ trong thư viện, hộp đặt giờ), nút Ghi macro trên toolbar *(chưa có trong design; logic: 22, 24)*
  - Xong khi: `#demo` hiện đủ các phần trên ở cả 2 theme, khớp canvas; app thật không có nút nào bấm được mà không chạy; `npm run build` sạch
  - Kết quả: màn mẫu `#demo` (màn soạn đủ loại bước + nhóm, thư viện, inspector từng loại, tab Input, form Chạy, hẹn giờ, ghi macro, debug), `#demo-overlay` (khoanh ảnh, chọn nhiều + biến đổi + ánh xạ), `#demo-minibar` (các trạng thái thanh mini). Code dùng lại khi làm logic: `components/StepCard.tsx`, `components/inputs/*`, `components/steps/*` (model + `StepFields` + `AnchorField`), `ScheduleDialog.tsx`, `AppChrome.tsx`, `MiniBarView`, `LibraryView`. Dữ liệu mẫu chỉ trong `frontend/src/demo/`
  - [x] Người dùng xem `#demo`, `#demo-overlay`, `#demo-minibar` và màn chính (nút tắt + *Sắp có*), góp ý trước khi commit
