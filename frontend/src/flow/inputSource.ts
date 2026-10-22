// Node Bộ gõ (bước 43 S6) – nhãn thao tác, dòng chi tiết thẻ. Chạy ở internal/engine/inputsource.go.
import type { Step } from '../types';

// '' = đọc (khớp model.InputOps)
export type InputOp = '' | 'set' | 'toggleIME';
export const INPUT_OPS: { id: InputOp; label: string }[] = [{ id: '', label: 'Đọc bộ gõ' }, { id: 'set', label: 'Đổi sang' }, { id: 'toggleIME', label: 'Bật / tắt tiếng Việt' }];
export const IME_DEFAULT_HOTKEY = 'ctrl+shift'; // khớp model.InputDefaultHotkey

// Dòng chi tiết thẻ: "đổi sang ABC" / "đọc bộ gõ đang dùng"
export function fmtInputSource(s: Pick<Step, 'inputOp' | 'inputLayout' | 'inputHotkey'>): string {
  if (s.inputOp === 'toggleIME') return `bấm ${s.inputHotkey?.trim() || IME_DEFAULT_HOTKEY} (đảo bật / tắt tiếng Việt)`;
  return s.inputOp === 'set' ? `đổi sang ${s.inputLayout?.trim() || 'chưa nhập bộ gõ'}` : 'đọc bộ gõ đang dùng';
}
