// Thao tác nguy hiểm cần người dùng đồng ý trước lần chạy đầu (script.DangerOps – internal/script/danger.go): xoá file
// vào / dọn Thùng rác, buộc thoát app, thoát app có buộc thoát khi quá giờ, đóng cửa sổ.

// Khớp script.DangerOp.
export type DangerOp = { ref: string; op: string; folder?: string; pattern?: string; target?: string };

// Câu cho hộp xác nhận.
export function dangerOpText(d: DangerOp): string {
  switch (d.op) {
    case 'emptyTrash': return 'Dọn Thùng rác – xoá hẳn mọi thứ trong Thùng rác, không khôi phục được';
    case 'trash': return `Chuyển các file khớp “${d.pattern}” trong ${d.folder} vào Thùng rác (khôi phục được)`;
    case 'forceQuit': return `Buộc thoát ${d.target} – mất dữ liệu chưa lưu`;
    case 'quit': return `Thoát ${d.target}, quá giờ chưa thoát thì buộc thoát – mất dữ liệu chưa lưu`;
    case 'close': return `Đóng cửa sổ ${d.target} (app có thể hỏi lưu)`;
  }
  return d.op;
}

// Thao tác không khôi phục được – tô màu cảnh báo.
export const isPermanent = (d: DangerOp) => d.op !== 'trash' && d.op !== 'close';

// Tiêu đề hộp xác nhận theo nhóm thao tác.
export function dangerTitle(ops: DangerOp[]): string {
  const file = ops.every((o) => o.op === 'trash' || o.op === 'emptyTrash');
  const app = ops.every((o) => o.op === 'quit' || o.op === 'forceQuit' || o.op === 'close');
  return file ? 'Kịch bản này sẽ xoá file' : app ? 'Kịch bản này sẽ thoát app / đóng cửa sổ' : 'Kịch bản có thao tác cần xác nhận';
}
