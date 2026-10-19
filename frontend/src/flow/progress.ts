// Chữ thanh mini theo tiến độ engine (Progress – internal/engine/engine.go): node đang chạy, node trước + cổng đã đi,
// lượt node Đếm. Hàm thuần – MiniBar gọi, test ở progress.test.ts.
import type { Progress } from '../engineStore';
import { stepLabel, type Step } from '../types';
import { PORT_SUCCESS } from './model';
import { portView } from './view';

const nameOf = (s: Step) => s.name || stepLabel(s);

// "n5 · Đếm"; tới bằng cổng khác Thành công thì kèm node trước: "n2 · Tìm nút Nhắn tin · Thất bại → n5 · Đếm"
// (đi thẳng success không ghi – đỡ dài, chỗ rẽ nhánh mới cần biết).
export function progressNode(p: Pick<Progress, 'nodeId' | 'ref' | 'from'>, steps: Step[]): string {
  const cur = steps.find((s) => s.id === p.nodeId);
  const node = `${p.ref ?? 'Node'}${cur ? ` · ${nameOf(cur)}` : ''}`;
  const f = p.from;
  if (!f || f.port === PORT_SUCCESS) return node;
  const from = steps.find((s) => s.ref === f.ref);
  return `${f.ref}${from ? ` · ${nameOf(from)}` : ''} · ${portView(f.port, from).label} → ${node}`;
}

export const counterText = (c: NonNullable<Progress['counter']>) => `lượt ${c.ref} ${c.count} / ${c.max}`;
