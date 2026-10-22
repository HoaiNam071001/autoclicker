// Node Nén / giải nén (bước 43 S8) – nhãn thao tác, dòng chi tiết thẻ. Chạy ở internal/engine/zip.go.
import type { Step } from '../types';

// '' = nén (khớp model.ZipOps)
export type ZipOp = '' | 'extract';
export const ZIP_OPS: { id: ZipOp; label: string }[] = [{ id: '', label: 'Nén ra .zip' }, { id: 'extract', label: 'Giải nén' }];

const base = (p?: string) => p?.trim().split(/[\\/]/).filter(Boolean).pop() ?? '';

// Dòng chi tiết thẻ: "bao-cao → nén (*.pdf)" / "pack.zip → giải nén"
export function fmtZip(s: Pick<Step, 'zipOp' | 'filePath' | 'filePattern'>): string {
  const src = base(s.filePath) || 'chưa nhập nguồn';
  if (s.zipOp === 'extract') return `${src} → giải nén`;
  const pat = s.filePattern?.trim();
  return `${src} → nén${pat ? ` (${pat})` : ''}`;
}
