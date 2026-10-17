// Flow (bước 20) – phần dùng chung với Go (flow.go): cổng ra theo loại node, thứ tự hiển thị node.
// flowOrder phải cùng công thức với flow.go – test chung dữ liệu testdata/flow_order.json (gốc repo).
import type { Edge, Step } from '../types';

export const PORT_SUCCESS = 'success';
export const PORT_FAIL = 'fail';

// Cổng ra hợp lệ theo thứ tự trên thẻ (trên → dưới). Node điều khiển (F3) thêm cổng riêng.
export const portsOf = (_s: Pick<Step, 'type'>): string[] => [PORT_SUCCESS, PORT_FAIL];

export const edgeOut = (edges: Edge[], from: string, port: string) => edges.find((e) => e.from === from && e.port === port)?.to ?? '';

// Các node theo thứ tự hiển thị (xem trước, danh sách trong form Chạy).
export const orderedSteps = <S extends Pick<Step, 'id' | 'type'>>(sc: { steps: S[]; start: string; edges: Edge[] }): S[] => {
  const byId = new Map(sc.steps.map((s) => [s.id, s]));
  return flowOrder(sc).map((id) => byId.get(id)!);
};

// Thứ tự hiển thị: duyệt sâu từ start, mỗi node đi các cổng theo thứ tự portsOf (nhánh success trước), node đã gặp bỏ qua;
// node không đi tới được xếp cuối theo thứ tự trong mảng.
// Đồ thị tối thiểu để duyệt: node chỉ cần id + loại (cổng ra theo loại).
export type FlowShape = { start: string; steps: Pick<Step, 'id' | 'type'>[]; edges: Edge[] };

export function flowOrder(sc: FlowShape): string[] {
  const byId = new Map(sc.steps.map((s) => [s.id, s]));
  const seen = new Set<string>();
  const order: string[] = [];
  const stack = byId.has(sc.start) ? [sc.start] : [];
  while (stack.length) {
    const id = stack.pop()!;
    if (seen.has(id)) continue;
    seen.add(id);
    order.push(id);
    const ports = portsOf(byId.get(id)!);
    for (let i = ports.length - 1; i >= 0; i--) { // đẩy ngược để cổng đầu được lấy ra trước
      const to = edgeOut(sc.edges, id, ports[i]);
      if (to && !seen.has(to) && byId.has(to)) stack.push(to);
    }
  }
  for (const s of sc.steps) if (!seen.has(s.id)) order.push(s.id);
  return order;
}
