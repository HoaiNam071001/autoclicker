// Node Chờ màn hình (bước 42 G1) – cách chờ, dòng chi tiết thẻ. Chạy ở internal/engine/waitscreen.go; vùng như Chụp màn hình.
import { fmtMs, type Step } from '../types';
import { shotAreaLabel } from './shot';

// '' = đứng yên (khớp model.WaitModes)
export type WaitMode = '' | 'change';
export const WAIT_MODES: { id: WaitMode; label: string }[] = [{ id: '', label: 'Đứng yên' }, { id: 'change', label: 'Thay đổi' }];

// Dòng chi tiết thẻ: "đứng yên 1 s · vùng kịch bản · tối đa 15 s". refOf: id node → mã (n2).
export function fmtWaitScreen(s: Step, refOf: (id: string) => string): string {
  const mode = s.waitMode === 'change' ? `thay đổi > ${s.diffPct ?? 0}%` : `đứng yên ${fmtMs(s.stableMs ?? 0)}`;
  return `${mode} · ${shotAreaLabel(s, refOf)} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
}
