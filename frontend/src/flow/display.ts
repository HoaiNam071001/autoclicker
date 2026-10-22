// Node Màn hình / âm thanh (bước 43 S7) – nhãn thao tác, dòng chi tiết thẻ. Chạy ở internal/engine/display.go.
import type { Step } from '../types';

// '' = đọc (khớp model.DisplayOps)
export type DisplayOp = '' | 'theme' | 'brightness' | 'audioOut' | 'audioIn' | 'sound';
export const DISPLAY_OPS: { id: DisplayOp; label: string }[] = [
  { id: '', label: 'Đọc trạng thái' }, { id: 'theme', label: 'Chế độ sáng / tối' }, { id: 'brightness', label: 'Độ sáng màn hình' },
  { id: 'audioOut', label: 'Đổi loa mặc định' }, { id: 'audioIn', label: 'Đổi micro mặc định' }, { id: 'sound', label: 'Phát âm báo' },
];
export const THEMES: { id: string; label: string }[] = [{ id: 'light', label: 'Sáng' }, { id: 'dark', label: 'Tối' }, { id: 'toggle', label: 'Đảo' }];
export const SYSTEM_SOUNDS = ['Glass', 'Ping', 'Hero', 'Sosumi', 'Funk', 'Submarine', 'Tink']; // macOS; Windows: Asterisk, Beep, Exclamation…

const THEME_LABEL: Record<string, string> = { light: 'sáng', dark: 'tối', toggle: 'đảo sáng ↔ tối' };

// Dòng chi tiết thẻ: "độ sáng 70 %" / "loa → AirPods"
export function fmtDisplay(s: Pick<Step, 'displayOp' | 'displayValue'>): string {
  const v = s.displayValue?.trim() || '';
  switch (s.displayOp ?? '') {
    case 'theme': return `chế độ ${THEME_LABEL[v] ?? (v || 'chưa chọn')}`;
    case 'brightness': return `độ sáng ${v || '?'} %`;
    case 'audioOut': return `loa → ${v || 'chưa nhập'}`;
    case 'audioIn': return `micro → ${v || 'chưa nhập'}`;
    case 'sound': return `phát ${v || 'chưa nhập'}`;
  }
  return 'đọc sáng / tối, độ sáng, loa, micro';
}
