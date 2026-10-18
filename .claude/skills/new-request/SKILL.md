---
name: new-request
description: Người dùng đưa yêu cầu mới / đổi yêu cầu / ý tưởng tính năng (vd "thêm tính năng…", "muốn đổi…", "sau này cần…"). Phân loại, đối chiếu plan + quyết định đã chốt, ghi vào docs/tasks đúng chỗ, hỏi những gì cần chốt – trước khi code.
---

# Nhận yêu cầu mới / đổi yêu cầu

Mục tiêu: không yêu cầu nào chỉ nằm trong chat; plan luôn phản ánh đúng ý người dùng.

## 1. Hiểu + đối chiếu

- Viết lại yêu cầu bằng 1–2 câu theo ngôn ngữ người dùng (kết quả họ muốn thấy, không phải cách làm).
- Đối chiếu:
  - Bảng *Tính năng của hệ thống* + bảng thứ tự bước trong `docs/tasks/README.md` – đã có / đã lên plan chưa?
  - `docs/tasks/decisions.md` – có mâu thuẫn quyết định đã chốt không?
  - `docs/ARCHITECTURE.md` – đụng phần nào, ảnh hưởng định dạng file / MCP không?

## 2. Phân loại

| Loại | Làm gì |
|---|---|
| Sửa nhỏ trong bước đang làm / đã xong | Thêm mục con `[ ]` vào file bước đó (ghi *(người dùng yêu cầu YYYY-MM-DD)*) |
| Mở rộng một bước chưa làm | Thêm mục con vào file bước đó, chỉnh mô tả bước |
| Tính năng mới hẳn | File mới `docs/tasks/phase3/<số>-<tên>.md` (theo mẫu file bước khác: mục cha `- [ ] **<số>. <tên>**`, mục con, *Chuẩn bị cho MCP*, *Kiểm tra tay*) + dòng trong bảng thứ tự ở `README.md`, đề xuất vị trí trong thứ tự làm |
| Đổi điều đã chốt | Nêu rõ cái cũ, cái mới, hệ quả (code phải sửa, file kịch bản cũ). Người dùng đồng ý → sửa `decisions.md` (ghi ngày, lý do), không xoá lịch sử mà ghi *thay bởi* |
| Lỗi | Dùng skill `fix-bug` |

## 3. Chốt với người dùng

- Câu hỏi chỉ cho những gì làm đổi hướng làm (hành vi, UI, phạm vi, ưu tiên) – mỗi câu kèm khuyến nghị của bạn và lý do. Gom vào một lần AskUserQuestion.
- Yêu cầu lớn → đề xuất chia mục con theo thứ tự làm được + kiểm tra được; nêu rủi ro (Windows, quyền macOS, hiệu năng chụp màn hình).
- Không bắt đầu code khi người dùng chỉ đang lên ý tưởng / hỏi; nếu họ bảo làm luôn → chuyển sang skill `implement-step`.

## 4. Ghi lại

- Cập nhật file bước / README (bảng thứ tự, *Trạng thái hiện tại* nếu đổi ưu tiên, bảng *Tính năng của hệ thống* thêm dòng ⬜).
- `log.md`: một dòng *Kế hoạch (người dùng yêu cầu): …*.
- `python3 scripts/tasks-index.py`. Commit `docs: …` nếu người dùng muốn commit.
