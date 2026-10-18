// Node Tệp / thư mục (bước 35 N4) – nhãn thao tác, dòng chi tiết thẻ, chữ thao tác xoá cho hộp xác nhận.
// Chạy ở internal/engine/file.go; thao tác xoá liệt kê ở script.DeleteOps.
import type { Step } from '../types';

// '' = ghi chữ (khớp model.FileOps)
export type FileOp = '' | 'row' | 'mkdir' | 'trash' | 'emptyTrash';
export const FILE_OPS: { id: FileOp; label: string }[] = [
  { id: '', label: 'Ghi chữ vào file' }, { id: 'row', label: 'Thêm dòng vào bảng (.csv / .xlsx)' }, { id: 'mkdir', label: 'Tạo thư mục' },
  { id: 'trash', label: 'Xoá file vào Thùng rác' }, { id: 'emptyTrash', label: 'Dọn Thùng rác' },
];
export const isDeleteOp = (op?: string) => op === 'trash' || op === 'emptyTrash';

// Đường dẫn trông như file (có đuôi .xxx ở phần cuối) → thư mục chứa nó; còn lại giữ nguyên.
export const folderOf = (p: string) => (/[^\\/]\.[A-Za-z0-9]{1,5}$/.test(p.trim()) ? p.trim().replace(/[\\/][^\\/]*$/, '') : p);

const base = (p?: string) => (p ?? '').trim().replace(/[\\/]+$/, '').split(/[\\/]/).pop() || '';

export function fmtFile(s: Pick<Step, 'fileOp' | 'filePath' | 'fileAppend' | 'cells' | 'filePattern'>): string {
  const name = base(s.filePath) || '(chưa chọn)';
  switch (s.fileOp) {
    case 'row': return `thêm dòng ${(s.cells ?? []).length} ô → ${name}`;
    case 'mkdir': return `tạo thư mục ${name}`;
    case 'trash': return `xoá ${s.filePattern?.trim() || '(chưa có mẫu)'} trong ${name} → Thùng rác`;
    case 'emptyTrash': return 'dọn Thùng rác (xoá hẳn)';
  }
  return `${s.fileAppend ? 'thêm vào cuối' : 'ghi'} ${name}`;
}

// Một thao tác xoá (khớp script.DeleteOp) → câu cho hộp xác nhận.
export type DeleteOp = { ref: string; op: string; folder: string; pattern: string };
export const deleteOpText = (d: DeleteOp) => (d.op === 'emptyTrash'
  ? 'Dọn Thùng rác – xoá hẳn mọi thứ trong Thùng rác, không khôi phục được'
  : `Chuyển các file khớp “${d.pattern}” trong ${d.folder} vào Thùng rác (khôi phục được)`);
