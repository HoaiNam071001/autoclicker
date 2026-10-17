// Sửa đồ thị flow (bước 20) – hàm thuần trên danh sách node + cạnh, không đụng store / React.
// Store (F4b) và màn mẫu #demo-flow đều gọi qua đây; hoàn tác / làm lại do nơi gọi lo.
// Quy tắc: mỗi cổng ra tối đa 1 cạnh; không nối node vào chính nó qua cùng một cạnh vô nghĩa (cho phép vòng lặp qua node khác).

import type { Edge } from './types';

export type { Edge };
export type Graph<N extends { id: string }> = { nodes: N[]; edges: Edge[] };

let seq = 0;
export const edgeId = () => `e${Date.now().toString(36)}${(seq++).toString(36)}`;

// Cạnh đang đi ra từ cổng (from, port)
export const edgeFrom = (edges: Edge[], from: string, port: string) => edges.find((e) => e.from === from && e.port === port);

// Nối cổng → node. Cổng đã có cạnh thì thay. Nối vào chính nó bị bỏ qua (vòng lặp phải qua node khác, vd Đếm).
export function connect(edges: Edge[], from: string, port: string, to: string): Edge[] {
  if (from === to) return edges;
  const old = edgeFrom(edges, from, port);
  if (old?.to === to) return edges;
  const rest = edges.filter((e) => e !== old);
  return [...rest, { id: old?.id ?? edgeId(), from, port, to }];
}

export const disconnect = (edges: Edge[], from: string, port: string) => edges.filter((e) => !(e.from === from && e.port === port));

export const removeEdge = (edges: Edge[], id: string) => edges.filter((e) => e.id !== id);

// Xoá node kèm mọi cạnh vào / ra. Node bắt đầu bị xoá → start chuyển sang node đầu còn lại (nếu có).
export function removeNode<N extends { id: string }>(g: Graph<N> & { start: string }, id: string): Graph<N> & { start: string } {
  const nodes = g.nodes.filter((n) => n.id !== id);
  return {
    nodes,
    edges: g.edges.filter((e) => e.from !== id && e.to !== id),
    start: g.start === id ? (nodes[0]?.id ?? '') : g.start,
  };
}

// Thêm node mới nối từ cổng (from, port). Cổng đã có cạnh A → B thì chèn vào giữa: A → mới → B
// (node mới đi tiếp qua cổng `through`, thường là success).
export function addConnected<N extends { id: string }>(g: Graph<N>, node: N, from: string, port: string, through?: string): Graph<N> {
  const old = edgeFrom(g.edges, from, port);
  let edges = connect(g.edges, from, port, node.id);
  if (old && through) edges = connect(edges, node.id, through, old.to);
  return { nodes: [...g.nodes, node], edges };
}

// Các node đi tới được từ start (theo mọi cạnh) – còn lại là node "không đi tới được" (cảnh báo, không chặn chạy)
export function reachable(edges: Edge[], start: string): Set<string> {
  const seen = new Set<string>();
  const stack = start ? [start] : [];
  while (stack.length) {
    const id = stack.pop()!;
    if (seen.has(id)) continue;
    seen.add(id);
    for (const e of edges) if (e.from === id) stack.push(e.to);
  }
  return seen;
}
