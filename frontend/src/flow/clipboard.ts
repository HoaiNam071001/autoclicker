// Node Clipboard (bước 35 N2) – nhãn phép, dòng chi tiết thẻ. Chạy ở internal/engine/clip.go.
import type { Step } from '../types';

// '' = đọc (khớp model.ClipOps)
export type ClipOp = '' | 'write' | 'image';
export const CLIP_OPS: { id: ClipOp; label: string }[] = [{ id: '', label: 'Đọc chữ' }, { id: 'write', label: 'Ghi chữ' }, { id: 'image', label: 'Ghi ảnh' }];

const short = (t: string, n = 40) => (t.length > n ? `${t.slice(0, n)}…` : t);

export function fmtClip(s: Pick<Step, 'clipOp' | 'text' | 'filePath'>): string {
  switch (s.clipOp) {
    case 'write': return s.text ? `ghi “${short(s.text)}”` : 'xoá clipboard';
    case 'image': return s.filePath?.trim() ? `ảnh ${short(s.filePath)}` : 'chưa nhập file ảnh';
  }
  return 'đọc chữ vào output.text';
}
