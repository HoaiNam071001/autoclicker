// Xếp node tự động (bước 20) bằng elkjs, thuật toán layered từ trái sang phải. Cổng ra cố định vị trí theo
// thẻ node (model.ts) nên cạnh ít cắt nhau; vòng lặp (cạnh quay lại) elk tự xử lý.
import type { ELK as ElkApi, ElkNode } from 'elkjs/lib/elk-api';
import { inputY, NODE_W, nodeHeight, portY } from './components/flow/model';
import type { Edge } from './flowEdit';

// elkjs ~1,4 MB – nạp lười lần đầu xếp, không nằm trong bundle chính
let elk: Promise<ElkApi> | null = null;
const getElk = () => (elk ??= import('elkjs/lib/elk.bundled.js').then((m) => new m.default()));

export type LayoutNode = { id: string; ports: string[] };
export type Positions = Record<string, { x: number; y: number }>;

export async function layoutFlow(nodes: LayoutNode[], edges: Edge[]): Promise<Positions> {
  const ids = new Set(nodes.map((n) => n.id));
  const graph: ElkNode = {
    id: 'root',
    layoutOptions: {
      'elk.algorithm': 'layered',
      'elk.direction': 'RIGHT',
      'elk.layered.spacing.nodeNodeBetweenLayers': '72',
      'elk.spacing.nodeNode': '50',
      'elk.layered.nodePlacement.strategy': 'NETWORK_SIMPLEX',
      'elk.layered.crossingMinimization.semiInteractive': 'true',
    },
    children: nodes.map((n) => ({
      id: n.id,
      width: NODE_W,
      height: nodeHeight(n.ports.length),
      layoutOptions: { 'elk.portConstraints': 'FIXED_POS' },
      ports: [
        { id: `${n.id}|in`, x: 0, y: inputY, width: 1, height: 1, layoutOptions: { 'elk.port.side': 'WEST' } },
        ...n.ports.map((p, i) => ({ id: `${n.id}|${p}`, x: NODE_W, y: portY(i), width: 1, height: 1, layoutOptions: { 'elk.port.side': 'EAST' } })),
      ],
    })),
    edges: edges
      .filter((e) => ids.has(e.from) && ids.has(e.to))
      .map((e) => ({ id: e.id, sources: [`${e.from}|${e.port}`], targets: [`${e.to}|in`] })),
  };
  const out = await (await getElk()).layout(graph);
  const pos: Positions = {};
  for (const c of out.children ?? []) pos[c.id] = { x: Math.round(c.x ?? 0), y: Math.round(c.y ?? 0) };
  return pos;
}
