# 29. Neo điểm theo 4 góc vùng

- [x] **29. Neo điểm theo 4 góc vùng** *(người dùng yêu cầu, 2026-10-16)* – Hiện điểm chỉ tính từ góc trên trái của vùng nên điểm ở phía phải / dưới lệch khi đổi vùng. Mỗi trục neo riêng: X neo trái (+) hoặc phải (−), Y neo trên (+) hoặc dưới (−). Ví dụ 4 điểm cách mỗi cạnh 300 px = (+300, +300), (−300, +300), (+300, −300), (−300, −300); đổi vùng thì mỗi điểm bám theo góc của nó.
  - [x] Model: `Point` thêm `ax` (`start` = trái | `end` = phải | `center` | `ratio` = theo % rộng), `ay` (tương tự cho trên / dưới); bỏ trống = `start` (file / kịch bản cũ chạy như cũ). Lưu offset dương + neo riêng thay vì dùng dấu (để phân biệt +0 / −0); UI hiện dạng có dấu ±
  - [x] Toàn màn hình (`region = nil`) cũng neo được: chốt **giữ `nil`**, lúc tính thì khung = màn hình chính (`mainScreen`, Go; frontend lấy qua `ScreenBounds`) – kịch bản mang sang máy khác cỡ màn hình vẫn bám góc. Đa màn hình xử lý ở 17
  - [x] Go: một hàm `resolvePoint(p, region)` dùng chung cho engine (`runStep.pt`) và `validateScript` (điểm phải nằm trong vùng); test 4 góc, center, ratio, đổi cỡ vùng, kịch bản cũ không có neo
  - [x] Frontend: `toScreen` / `toRelative` / `endpoints()` trong `scriptStore.ts` tính theo neo (khớp Go); đổi neo giữ nguyên chỗ trên màn hình (tính lại offset)
  - [x] Neo khi đặt điểm: mặc định **tự chọn góc gần nhất** theo từng trục; kéo marker / chọn lại giữ neo hiện tại
  - [x] Cài đặt kịch bản *góc gần nhất / luôn trên trái* – `Script.anchorMode`, chọn trong popup nút "N vòng" (làm cùng 21)
  - [x] Form: ô X / Y có dấu ± – gõ số âm = neo phải / dưới; nút chọn neo dạng lưới 3×3 nhỏ (4 góc, giữa cạnh, tâm) + tuỳ chọn %; danh sách bước hiện toạ độ dạng `(−300, +300)`
  - [x] Overlay: đường mảnh từ marker tới cạnh / góc đang neo (hiện với bước đang chọn hoặc khi giữ Alt); đổi vùng (kéo vùng mới) thấy marker bám theo góc ngay
  - [x] Kéo thả "lệch so với A" (ΔX, ΔY) giữ nguyên – lệch không phụ thuộc neo
  - [x] Nối UI khung từ 27b (ô ±, lưới chọn neo, đường neo trên overlay)
  - Xong khi: đặt 4 điểm ở 4 góc vùng 1000×700, đổi vùng sang 1400×900 → 4 điểm vẫn cách mỗi cạnh đúng 300 px khi chạy thật và trên overlay
