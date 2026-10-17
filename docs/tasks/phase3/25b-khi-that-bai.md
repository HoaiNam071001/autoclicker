# 25b. Khi thất bại

- [x] **25b. Khi thất bại** *(làm sớm 2026-10-17, người dùng yêu cầu if / else – đây là phần nền, phần khối điều kiện ở 20)* – Mỗi bước có cài đặt *Khi không thấy*: thử lại N lần, rồi **Dừng kịch bản / Bỏ qua / Nhảy tới bước**.
  - [x] Model `onFail { retries, action: stop | next | goto, goto }` (nil = dừng như cũ); `checkStep` kiểm action / số lần, `validateScript` kiểm bước nhảy tới còn, đang bật, không phải chính nó
  - [x] Engine (`runWithRetry`): thử lại ngay (không chờ delayAfter), thất bại thì quên chỗ thấy cũ (Click theo kết quả bước trước không click nhầm chỗ cũ); bỏ qua = sang bước sau; nhảy tới bước trước / sau trong cùng vòng, tối đa 1000 lần nhảy mỗi vòng (chặn vòng lặp không dừng); tiến độ `retry` / `retries`
  - [x] UI: inspector *Khi không thấy / không đổi màu / mở không được* (Chờ ảnh, Click theo ảnh, Chờ pixel, Mở) – số lần thử lại, segmented, chọn bước nhảy tới; chip trên thẻ bước (`↻2 → bước 5`); thanh mini "thử lại x / y"
  - [x] Test `TestOnFailRetry`, `TestOnFailNextSkipsStaleClick`, `TestOnFailGoto`, `TestValidateOnFail`; mô tả cho AI ở `commonFieldsDoc` (`steptypes.go`)
  - Không làm: `skipGroup` (bước 20 dùng *Thoát khối* thay)
