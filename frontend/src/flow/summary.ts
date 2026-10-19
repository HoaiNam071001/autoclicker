// Tóm tắt một dòng cho đầu khối thu gọn của inspector (design v3, bước 38): hiện bên phải tên khối khi đóng,
// để không phải mở từng khối mới biết node đang cài gì. Hàm thuần – test ở summary.test.ts.
import type { FlowEdgeView, PortView } from '../components/flow/model';
import { fmtMs, Retry } from '../types';

// Thời gian: chờ sau node (hoặc thời gian chờ của node Chờ) ± ngẫu nhiên. Ô nhận biến thì ghi biến.
export function timingSummary(delay: number | string, jitter: number | string = 0): string {
  if (typeof delay === 'string') return `chờ ${delay}`;
  if (typeof jitter === 'string') return `chờ ${fmtMs(delay)} ± ${jitter}`;
  if (!delay && !jitter) return 'không chờ';
  if (!jitter) return `chờ ${fmtMs(delay)}, cố định`;
  const lo = Math.max(delay - jitter, 0), hi = delay + jitter;
  return hi < 1000 ? `chờ ${lo}–${hi} ms` : `chờ ${fmtMs(lo)}–${fmtMs(hi)}`;
}

export function retrySummary(r: Retry | undefined): string {
  if (!r?.count) return 'không thử lại';
  return `thử lại ${r.count} lần · nghỉ ${fmtMs(r.delayMs ?? 0)}`;
}

const SHORT: Record<string, string> = { success: 'ok', fail: 'lỗi' };
// Cổng ra: "ok → n4 · lỗi → dừng". Nhiều hơn 2 cổng thì chỉ đếm: "3 cổng · 2 đã nối".
export function portsSummary(nodeId: string, ports: PortView[], edges: FlowEdgeView[], refOf: (id: string) => string): string {
  if (!ports.length) return 'không có – kết thúc ở đây';
  const to = (p: PortView) => edges.find((e) => e.from === nodeId && e.port === p.id)?.to;
  if (ports.length > 2) return `${ports.length} cổng · ${ports.filter(to).length} đã nối`;
  return ports.map((p) => {
    const t = to(p);
    return `${SHORT[p.id] ?? p.label} → ${t ? refOf(t) : p.kind === 'fail' ? 'dừng' : 'hết'}`;
  }).join(' · ');
}

export const outputsSummary = (n: number) => `${n} biến`;
