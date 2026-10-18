// Node Hệ thống (bước 35 N5) – nhãn thao tác âm lượng, dòng chi tiết thẻ. Chạy ở internal/engine/system.go.
import type { Step } from '../types';

export type SysOp = 'volume' | 'volumeUp' | 'volumeDown' | 'mute' | 'unmute';
export const SYS_OPS: { id: SysOp; label: string }[] = [
  { id: 'volume', label: 'Đặt âm lượng' }, { id: 'volumeUp', label: 'Tăng âm lượng' }, { id: 'volumeDown', label: 'Giảm âm lượng' },
  { id: 'mute', label: 'Tắt tiếng' }, { id: 'unmute', label: 'Bật tiếng' },
];
// Thao tác có ô số %: đặt mức / bước tăng giảm
export const sysHasLevel = (op?: string) => op === 'volume' || op === 'volumeUp' || op === 'volumeDown';

export function fmtSystem(s: Pick<Step, 'sysOp' | 'volume'>): string {
  const step = s.volume || 10;
  switch (s.sysOp) {
    case 'volumeUp': return `âm lượng +${step} %`;
    case 'volumeDown': return `âm lượng −${step} %`;
    case 'mute': return 'tắt tiếng';
    case 'unmute': return 'bật tiếng';
  }
  return `âm lượng = ${s.volume ?? 0} %`;
}
