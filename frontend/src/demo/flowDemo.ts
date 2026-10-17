// Dữ liệu mẫu cho màn #demo-flow (bước 20 F4a): flow "nhắn tin Facebook" dạng node + cạnh, đủ các loại node điều khiển,
// node tắt, node không đi tới được, và 2 lần chạy mẫu (đang chạy / thất bại) để xem trạng thái + nhật ký.
// KHÔNG dùng trong app thật – model thật ở F1 (model.go / types.ts).
import type { StepIconName } from '../components/icons';
import type { FlowEdgeView, FlowNodeView, PortView, RunRow, RunState } from '../components/flow/model';
import type { OutputField, Retry } from '../components/flow/PortField';
import type {
  CounterData, EndData, IfData, SetVarData, SwitchData, SwitchScreenData,
} from '../components/flow/nodes/ControlFields';
import { DemoStep, demoItems, findDemoStep } from './demoScript';

export type NodeBody =
  | { kind: 'step'; step: DemoStep }
  | { kind: 'wait'; ms: number }
  | { kind: 'if'; data: IfData }
  | { kind: 'switch'; data: SwitchData }
  | { kind: 'switchScreen'; data: SwitchScreenData }
  | { kind: 'counter'; data: CounterData }
  | { kind: 'end'; data: EndData }
  | { kind: 'setVar'; data: SetVarData };

export type DemoFlowNode = Omit<FlowNodeView, 'ports'> & { body: NodeBody; retry?: Retry };

const ACTION_PORTS: PortView[] = [
  { id: 'success', label: 'Thành công', kind: 'ok' },
  { id: 'fail', label: 'Thất bại', kind: 'fail' },
];

// Cổng ra theo loại node – switch lấy theo danh sách case (thêm case trong inspector là có cổng mới)
export function portsOf(b: NodeBody): PortView[] {
  switch (b.kind) {
    case 'if': return [{ id: 'true', label: 'Đúng', kind: 'neutral' }, { id: 'false', label: 'Sai', kind: 'neutral' }];
    case 'switchScreen': return [
      ...b.data.cases.map((c): PortView => ({ id: `case:${c.key}`, label: c.label, kind: 'neutral' })),
      { id: 'none', label: 'Không thấy', kind: 'neutral' },
    ];
    case 'switch': return [
      ...b.data.cases.map((c): PortView => ({ id: `case:${c.key}`, label: c.label, kind: 'neutral' })),
      { id: 'default', label: 'Khác', kind: 'neutral' },
    ];
    case 'counter': return [{ id: 'next', label: 'Còn lượt', kind: 'neutral' }, { id: 'done', label: 'Hết lượt', kind: 'neutral' }];
    case 'end': return [];
    case 'setVar': return [{ id: 'success', label: 'Tiếp', kind: 'ok' }];
    default: return ACTION_PORTS;
  }
}

// Output node trả về (hiện trong inspector, chèn bằng {{nX.output.key}})
export function outputsOf(b: NodeBody): OutputField[] {
  const found = [
    { key: 'x', desc: 'tâm chỗ thấy – X' }, { key: 'y', desc: 'tâm chỗ thấy – Y' },
    { key: 'w', desc: 'rộng' }, { key: 'h', desc: 'cao' }, { key: 'score', desc: '% khớp' }, { key: 'text', desc: 'chữ đọc được' },
  ];
  switch (b.kind) {
    case 'switchScreen': return [{ key: 'case', desc: 'case nào hiện' }, ...found];
    case 'switch': return [{ key: 'case', desc: 'case khớp' }, { key: 'value', desc: 'giá trị đã xét' }];
    case 'counter': return [{ key: 'count', desc: 'số lần đã đi qua' }];
    case 'setVar': return [{ key: 'value', desc: 'giá trị sau khi đặt' }];
    case 'if': return [{ key: 'result', desc: 'true / false' }, ...found.slice(0, 2)];
    case 'step': {
      const t = b.step.type;
      if (t === 'waitImage') return found;
      if (t === 'clickImage' || t === 'click') return [{ key: 'x', desc: 'X đã click' }, { key: 'y', desc: 'Y đã click' }];
      if (t === 'open') return [{ key: 'target', desc: 'URL / app đã mở' }];
      if (t === 'type') return [{ key: 'text', desc: 'chữ đã gõ (sau khi thay biến)' }];
      return [];
    }
    default: return [];
  }
}

const step = (id: string) => findDemoStep(id) ?? (demoItems[0] as DemoStep);

export const demoFlowStart = 'a1';

export const demoFlowNodes: DemoFlowNode[] = [
  { id: 'a1', ref: 'n1', icon: 'open', title: 'Mở trang cá nhân', detail: 'facebook.com/{{profile}} · Chrome', chips: ['1 input'], body: { kind: 'step', step: step('d1') } },
  {
    id: 'a2', ref: 'n2', icon: 'switchScreen', title: 'Đang ở màn nào?', detail: 'chờ tối đa 3 s · 2 case',
    body: { kind: 'switchScreen', data: { timeoutMs: 3000, cases: [
      { key: 'login', label: 'Đăng nhập', by: 'text', text: 'Đăng nhập' },
      { key: 'profile', label: 'Có nút Nhắn tin', by: 'text', text: 'Nhắn tin' },
    ] } },
  },
  { id: 'a3', ref: 'n3', icon: 'clickImage', title: 'Click nút Nhắn tin', detail: 'chỗ thấy ở n2 · giữa chữ', retry: { count: 2, delayMs: 500 }, chips: ['thử lại 2'], body: { kind: 'step', step: step('d4') } },
  { id: 'a4', ref: 'n4', icon: 'wait', title: 'Chờ khung chat', detail: '800 ms', enabled: false, body: { kind: 'wait', ms: 800 } },
  { id: 'a5', ref: 'n5', icon: 'type', title: 'Gõ lời chào', detail: '"xin chào {{M1}} và {{M2}}"', chips: ['2 input'], body: { kind: 'step', step: step('d5') } },
  { id: 'a6', ref: 'n6', icon: 'key', title: 'Enter', detail: 'gửi tin', body: { kind: 'step', step: step('d6') } },
  { id: 'a7', ref: 'n7', icon: 'if', title: 'Đã gửi?', detail: 'thấy chữ "Đã gửi" · chờ 3 s', body: { kind: 'if', data: { kind: 'text', text: 'Đã gửi', waitMs: 3000, left: '', op: 'eq', right: '', negate: false } } },
  { id: 'a8', ref: 'n8', icon: 'end', title: 'Xong', detail: 'thành công · đã gửi tin', body: { kind: 'end', data: { status: 'success', message: 'Đã gửi tin cho {{M1}}' } } },
  { id: 'a9', ref: 'n9', icon: 'end', title: 'Chưa đăng nhập', detail: 'thất bại · cần đăng nhập Facebook', body: { kind: 'end', data: { status: 'fail', message: 'Chưa đăng nhập Facebook trên Chrome' } } },
  { id: 'a10', ref: 'n10', icon: 'counter', title: 'Thử cuộn', detail: 'tối đa 10 lần', body: { kind: 'counter', data: { max: 10 } } },
  { id: 'a11', ref: 'n11', icon: 'scroll', title: 'Cuộn xuống', detail: '3 nấc tại chuột', body: { kind: 'step', step: { ...step('d7'), detail: 'xuống 3 nấc tại chuột' } } },
  { id: 'a12', ref: 'n12', icon: 'end', title: 'Không thấy nút', detail: 'thất bại · sau 10 lần cuộn', body: { kind: 'end', data: { status: 'fail', message: 'Không thấy nút Nhắn tin sau 10 lần cuộn' } } },
  { id: 'a14', ref: 'n13', icon: 'end', title: 'Gửi không được', detail: 'thất bại · không thấy "Đã gửi"', body: { kind: 'end', data: { status: 'fail', message: 'Không thấy chữ “Đã gửi” sau 3 s' } } },
  { id: 'a13', ref: 'n14', icon: 'key', title: 'Esc', detail: 'đóng popup', body: { kind: 'step', step: { ...step('d6'), detail: 'Esc' } } },
];

export const demoFlowEdges: FlowEdgeView[] = [
  { id: 'e1', from: 'a1', port: 'success', to: 'a2' },
  { id: 'e2', from: 'a2', port: 'case:login', to: 'a9' },
  { id: 'e3', from: 'a2', port: 'case:profile', to: 'a3' },
  { id: 'e4', from: 'a2', port: 'none', to: 'a10' },
  { id: 'e5', from: 'a3', port: 'success', to: 'a4' },
  { id: 'e6', from: 'a4', port: 'success', to: 'a5' },
  { id: 'e7', from: 'a5', port: 'success', to: 'a6' },
  { id: 'e8', from: 'a6', port: 'success', to: 'a7' },
  { id: 'e9', from: 'a7', port: 'true', to: 'a8' },
  { id: 'e10', from: 'a7', port: 'false', to: 'a14' },
  { id: 'e11', from: 'a10', port: 'next', to: 'a11' },
  { id: 'e12', from: 'a10', port: 'done', to: 'a12' },
  { id: 'e13', from: 'a11', port: 'success', to: 'a2' },
];

// ---- Lần chạy mẫu ----

export type DemoRun = {
  title: string; summary: string; ok: boolean | null;
  nodes: Record<string, { run: RunState; visits?: number }>;
  taken: string[]; active?: string;
  rows: RunRow[];
};

const found = { case: 'profile', x: 1258, y: 406, w: 96, h: 36, score: 100, text: 'Nhắn tin' };

// Đang chạy: không thấy nút → cuộn 1 lần → thấy → click, gõ, Enter → đang xét "Đã gửi?"
export const runningRun: DemoRun = {
  title: 'lần chạy 14:32:05 · vòng 1 / 3', summary: 'Đang chạy n7 · 8,1 s', ok: null,
  nodes: {
    a1: { run: 'ok' }, a2: { run: 'ok', visits: 2 }, a10: { run: 'ok' }, a11: { run: 'ok' }, a3: { run: 'ok' },
    a4: { run: 'skipped' }, a5: { run: 'ok' }, a6: { run: 'ok' }, a7: { run: 'running' },
  },
  taken: ['e1', 'e4', 'e11', 'e13', 'e3', 'e5', 'e6', 'e7', 'e8'], active: 'e8',
  rows: [
    { nodeId: 'a1', at: '14:32:05', status: 'success', port: 'success', next: 'n2', ms: 1210, output: { target: 'https://facebook.com/lan.nguyen.92' } },
    { nodeId: 'a2', at: '14:32:06', status: 'success', port: 'none', next: 'n10', ms: 3020, output: { case: 'none' },
      logs: [{ t: '+0,0', level: 'info', msg: 'chụp màn hình, đọc chữ (0,9 s)' }, { t: '+3,0', level: 'warn', msg: 'hết 3 s, không case nào hiện' }] },
    { nodeId: 'a10', at: '14:32:09', status: 'success', port: 'next', next: 'n11', ms: 0, output: { count: 1 } },
    { nodeId: 'a11', at: '14:32:09', status: 'success', port: 'success', next: 'n2', ms: 320 },
    { nodeId: 'a2', at: '14:32:10', status: 'success', port: 'case:profile', next: 'n3', ms: 940, output: found,
      logs: [{ t: '+0,9', level: 'info', msg: 'thấy "Nhắn tin" ở (1258, 406) – case Có nút Nhắn tin' }] },
    { nodeId: 'a3', at: '14:32:11', status: 'success', port: 'success', next: 'n4', ms: 110, output: { x: 1258, y: 406 } },
    { nodeId: 'a4', at: '14:32:11', status: 'skipped', port: 'success', next: 'n5', logs: [{ t: '+0,0', level: 'info', msg: 'node đang tắt – đi tiếp cổng Thành công' }] },
    { nodeId: 'a5', at: '14:32:11', status: 'success', port: 'success', next: 'n6', ms: 1380, output: { text: 'xin chào Lan và cuối tuần này đi cà phê không?' } },
    { nodeId: 'a6', at: '14:32:13', status: 'success', port: 'success', next: 'n7', ms: 40 },
    { nodeId: 'a7', at: '14:32:13', status: 'running', logs: [{ t: '+0,0', level: 'info', msg: 'đang tìm chữ "Đã gửi" (0,6 / 3 s)' }] },
  ],
};

// Thất bại: 10 lần không thấy nút, cuộn hết lượt → Kết thúc báo lỗi
export const failedRun: DemoRun = (() => {
  const rows: RunRow[] = [{ nodeId: 'a1', at: '14:40:01', status: 'success', port: 'success', next: 'n2', ms: 1180, output: { target: 'https://facebook.com/minh.tran' } }];
  let s = 2;
  const at = () => `14:40:${String(s).padStart(2, '0')}`;
  for (let i = 1; i <= 11; i++) {
    rows.push({ nodeId: 'a2', at: at(), status: 'success', port: 'none', next: 'n10', ms: 3010, output: { case: 'none' } });
    s += 3;
    rows.push({ nodeId: 'a10', at: at(), status: 'success', port: i <= 10 ? 'next' : 'done', next: i <= 10 ? 'n11' : 'n12', ms: 0, output: { count: i } });
    if (i <= 10) rows.push({ nodeId: 'a11', at: at(), status: 'success', port: 'success', next: 'n2', ms: 310 });
    s += 1;
  }
  rows.push({ nodeId: 'a12', at: at(), status: 'fail', ms: 0, error: { code: 'end', message: 'Không thấy nút Nhắn tin sau 10 lần cuộn' },
    logs: [{ t: '+0,0', level: 'error', msg: 'kịch bản dừng: Kết thúc (thất bại)' }] });
  return {
    title: 'lần chạy 14:40:01 · vòng 1 / 1', summary: 'Thất bại ở n12 sau 45,2 s · Không thấy nút Nhắn tin sau 10 lần cuộn', ok: false,
    nodes: { a1: { run: 'ok' }, a2: { run: 'ok', visits: 11 }, a10: { run: 'ok', visits: 11 }, a11: { run: 'ok', visits: 10 }, a12: { run: 'fail' } },
    taken: ['e1', 'e4', 'e11', 'e13', 'e12'], rows,
  };
})();

// ---- Thêm node mới (bảng thêm node của canvas) ----

const DEFAULTS: Record<string, { icon: StepIconName; title: string; detail: string; body: () => NodeBody }> = {
  click: { icon: 'click', title: 'Click', detail: 'chọn điểm trên overlay', body: () => ({ kind: 'step', step: step('d9') }) },
  scroll: { icon: 'scroll', title: 'Cuộn', detail: 'xuống 3 nấc tại chuột', body: () => ({ kind: 'step', step: step('d7') }) },
  clickImage: { icon: 'clickImage', title: 'Click theo ảnh', detail: 'tự tìm ảnh mẫu', body: () => ({ kind: 'step', step: step('d4') }) },
  type: { icon: 'type', title: 'Gõ chữ', detail: '""', body: () => ({ kind: 'step', step: step('d5') }) },
  key: { icon: 'key', title: 'Phím tắt', detail: 'Enter', body: () => ({ kind: 'step', step: step('d6') }) },
  wait: { icon: 'wait', title: 'Chờ', detail: '500 ms', body: () => ({ kind: 'wait', ms: 500 }) },
  waitImage: { icon: 'waitImage', title: 'Chờ ảnh', detail: 'tối đa 10 s, khớp ≥ 90%', body: () => ({ kind: 'step', step: step('d3') }) },
  waitPixel: { icon: 'waitPixel', title: 'Chờ màu pixel', detail: 'X 640 Y 300 = #1877F2', body: () => ({ kind: 'step', step: step('d8') }) },
  open: { icon: 'open', title: 'Mở URL', detail: 'https://…', body: () => ({ kind: 'step', step: step('d1') }) },
  if: { icon: 'if', title: 'Nếu', detail: 'thấy chữ "…"', body: () => ({ kind: 'if', data: { kind: 'text', text: '', waitMs: 0, left: '', op: 'eq', right: '', negate: false } }) },
  switch: { icon: 'switch', title: 'Switch theo giá trị', detail: '{{…}} · 2 case', body: () => ({ kind: 'switch', data: { value: '', cases: [{ key: 'c1', label: 'Case 1', op: 'eq', match: '' }, { key: 'c2', label: 'Case 2', op: 'eq', match: '' }] } }) },
  switchScreen: { icon: 'switchScreen', title: 'Switch theo màn hình', detail: 'chờ tối đa 3 s · 2 case', body: () => ({ kind: 'switchScreen', data: { timeoutMs: 3000, cases: [{ key: 'c1', label: 'Case 1', by: 'text', text: '' }, { key: 'c2', label: 'Case 2', by: 'text', text: '' }] } }) },
  counter: { icon: 'counter', title: 'Đếm', detail: 'tối đa 10 lần', body: () => ({ kind: 'counter', data: { max: 10 } }) },
  end: { icon: 'end', title: 'Kết thúc', detail: 'thành công', body: () => ({ kind: 'end', data: { status: 'success', message: '' } }) },
  setVar: { icon: 'setVar', title: 'Đặt biến', detail: 'vars.ten = …', body: () => ({ kind: 'setVar', data: { name: 'ten', op: 'set', value: '' } }) },
};

export function newDemoNode(type: keyof typeof DEFAULTS, existing: DemoFlowNode[]): DemoFlowNode {
  const d = DEFAULTS[type];
  let n = existing.length + 1;
  while (existing.some((x) => x.ref === `n${n}`)) n++;
  return { id: `x${Date.now().toString(36)}`, ref: `n${n}`, icon: d.icon, title: d.title, detail: d.detail, body: d.body() };
}
export const NEW_NODE_TYPES = Object.keys(DEFAULTS) as (keyof typeof DEFAULTS)[];
