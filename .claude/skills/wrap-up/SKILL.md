---
name: wrap-up
description: Chốt việc trước khi commit / kết thúc session (vd "commit", "chốt", "xong rồi", "kết thúc"): đồng bộ docs/tasks, ARCHITECTURE, chạy kiểm tra, commit theo quy ước.
---

# Chốt việc

## 1. Kiểm tra code

```bash
go vet ./... && go test ./...
cd frontend && npm run build && npm test
```

Không xanh thì dừng, báo người dùng (kèm output) – không commit đỏ trừ khi người dùng bảo.

## 2. Đồng bộ tài liệu (theo `CLAUDE.md` › *Bắt đầu mỗi session* bước 2)

- [ ] File bước: mục con xong đã `[x]`; mục cha có mọi mục con code `[x]` → tick mục cha; việc phát sinh đã thành mục con.
- [ ] Việc cần thử tay → `checks.md` (mới nhất ở trên, nhóm `## <số bước> – <tên> · [link]`), file bước chỉ để dòng link.
- [ ] `README.md`: *Trạng thái hiện tại* (đang làm / tiếp theo / cập nhật lần cuối – ngày tuyệt đối); bảng *Tính năng của hệ thống* nếu tính năng đổi trạng thái (⬜ → 🧪 → ✅, 🟡).
- [ ] `log.md`: một dòng `- **YYYY-MM-DD** – <việc>. <kết quả kiểm tra: lệnh nào OK / chưa thử tay>.`
- [ ] `docs/ARCHITECTURE.md`: đổi luồng, thêm method `App` / event / file lõi / cổng → sửa.
- [ ] `CLAUDE.md` › *Cấu trúc*: thêm file / module mới → một cụm mô tả ngắn.
- [ ] `decisions.md`: người dùng vừa chốt điều gì mới → ghi (ngày, lý do).
- [ ] `python3 scripts/tasks-index.py` rồi `python3 scripts/check-docs.py` (phải 0 lỗi).

## 3. Commit (khi người dùng yêu cầu commit / đã cho phép)

- `git status` + `git diff --stat` xem đúng file; không add file rác (`build/bin`, log, ảnh tạm).
- Message tiếng Việt, theo lịch sử: `feat: <tính năng> (bước <số>)`, `fix: <triệu chứng> – <nguyên nhân>`, `docs: …`, `refactor: …`. Thân message (nếu cần) liệt kê ý chính.
- Đang ở `main` và thay đổi lớn → hỏi có muốn tạo nhánh không.

## 4. Báo người dùng

Ngắn: đã xong gì, kiểm tra nào đã chạy (kết quả thật), việc chờ họ thử tay (link `checks.md`), bước tiếp theo.
