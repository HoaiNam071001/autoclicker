// Node Nguồn (bước 43 S3) – nhãn hành động, dòng chi tiết thẻ. Chạy ở internal/engine/power.go.
import type { Step } from '../types';

// Khớp model.PowerOps
export type PowerOp = 'lock' | 'displaySleep' | 'sleep' | 'logout' | 'shutdown' | 'restart';
export const POWER_OPS: { id: PowerOp; label: string; danger: boolean }[] = [
  { id: 'lock', label: 'Khoá màn hình', danger: false }, { id: 'displaySleep', label: 'Tắt màn hình', danger: false },
  { id: 'sleep', label: 'Cho máy ngủ', danger: true }, { id: 'logout', label: 'Đăng xuất', danger: true },
  { id: 'shutdown', label: 'Tắt máy', danger: true }, { id: 'restart', label: 'Khởi động lại', danger: true },
];
export const POWER_COUNTDOWN_SEC = 30; // khớp model.PowerCountdownSec

export const powerNeedsConsent = (op?: PowerOp) => !!POWER_OPS.find((o) => o.id === op)?.danger;

// Dòng chi tiết thẻ: "Tắt máy sau 30 giây" / "Khoá màn hình"
export function fmtPower(s: Pick<Step, 'powerOp' | 'countdownSec'>): string {
  const op = POWER_OPS.find((o) => o.id === (s.powerOp ?? 'lock'));
  const label = op?.label ?? 'Khoá màn hình';
  const secs = s.countdownSec || (op?.danger ? POWER_COUNTDOWN_SEC : 0);
  return secs > 0 ? `${label} sau ${secs} giây` : label;
}
