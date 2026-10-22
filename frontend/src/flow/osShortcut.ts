// Node Shortcuts hệ điều hành (bước 43 S2) – dòng chi tiết thẻ. Chạy ở internal/engine/osshortcut.go.
import type { Step } from '../types';

export function fmtOSShortcut(s: Partial<Pick<Step, 'shortcutName' | 'text'>>): string {
  const name = s.shortcutName?.trim() || '(chưa chọn)';
  return s.text?.trim() ? `${name} · có input` : name;
}
