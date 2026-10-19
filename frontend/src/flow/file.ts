// Node Tệp / thư mục (bước 35 N4) – nhãn thao tác, dòng chi tiết thẻ. Chạy ở internal/engine/file.go; thao tác xoá
// cần xác nhận (flow/danger.ts – script.DangerOps).
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
