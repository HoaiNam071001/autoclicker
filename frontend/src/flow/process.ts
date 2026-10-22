// Node Tiến trình (bước 43 S5) – nhãn thao tác, dòng chi tiết thẻ. Chạy ở internal/engine/process.go.
import type { Step } from '../types';

// '' = đọc (khớp model.ProcOps)
export type ProcOp = '' | 'waitStart' | 'waitQuit' | 'waitIdle';
export const PROC_OPS: { id: ProcOp; label: string }[] = [
  { id: '', label: 'Đọc CPU / RAM, liệt kê' }, { id: 'waitStart', label: 'Chờ tiến trình chạy' },
  { id: 'waitQuit', label: 'Chờ tiến trình thoát' }, { id: 'waitIdle', label: 'Chờ hết bận (CPU thấp)' },
];

export const PROC_IDLE_CPU = 5; // khớp model.ProcIdleCPU

// Dòng chi tiết thẻ: "ffmpeg → chờ hết bận (CPU < 5 %)"
export function fmtProcess(s: Pick<Step, 'procOp' | 'procName' | 'cpuBelow'>): string {
  const op = s.procOp ?? '';
  const name = s.procName?.trim() || (op === '' ? 'mọi tiến trình' : 'chưa nhập tên');
  switch (op) {
    case 'waitStart': return `${name} → chờ chạy`;
    case 'waitQuit': return `${name} → chờ thoát`;
    case 'waitIdle': return `${name} → chờ hết bận (CPU < ${s.cpuBelow || PROC_IDLE_CPU} %)`;
  }
  return `${name} → đọc CPU / RAM`;
}
