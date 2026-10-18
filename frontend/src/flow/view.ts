// Node (Step) → dữ liệu hiển thị thẻ trên canvas (FlowNodeView): tên, dòng chi tiết, chip, cổng ra có nhãn, cảnh báo,
// trạng thái chạy. Hàm thuần – FlowEditor (app) gọi với dữ liệu thật; test ở view.test.ts.
import type { FlowNodeView, PortView, RunState } from '../components/flow/model';
import { humanLabel, usesHuman } from './human';
import { Cond, fmtMs, isControl, numVarOf, Script, Step, stepLabel } from '../types';
import { fmtPos } from './detail';
import { reachable } from './edit';
import { CASE_PREFIX, flowOrder, PORT_FAIL, PORT_SUCCESS, PortShape, portsOf } from './model';

// Nhãn + màu cổng ra. Cổng case:<key> lấy nhãn từ case của node.
const PORTS: Record<string, Omit<PortView, 'id'>> = {
  [PORT_SUCCESS]: { label: 'Thành công', kind: 'ok' },
  [PORT_FAIL]: { label: 'Thất bại', kind: 'fail' },
  true: { label: 'Đúng', kind: 'neutral' },
  false: { label: 'Sai', kind: 'neutral' },
  default: { label: 'Khác', kind: 'neutral' },
  none: { label: 'Không thấy', kind: 'neutral' },
  next: { label: 'Còn lượt', kind: 'neutral' },
  done: { label: 'Hết lượt', kind: 'neutral' },
};
export function portView(id: string, s?: Pick<Step, 'type' | 'cases'>): PortView {
  if (id.startsWith(CASE_PREFIX)) {
    const key = id.slice(CASE_PREFIX.length);
    const c = s?.cases?.find((x) => x.key === key);
    return { id, label: c?.label || key, kind: 'neutral' };
  }
  if (id === PORT_SUCCESS && s?.type === 'setVar') return { id, label: 'Tiếp', kind: 'ok' };
  return { id, ...(PORTS[id] ?? { label: id, kind: 'neutral' }) };
}
export const portViews = (s: PortShape) => portsOf(s).map((p) => portView(p, s));

export const nodeTitle = (s: Step) => s.name || `${stepLabel(s)}${s.type === 'click' && !numVarOf(s, 'clicks') && s.clicks > 1 ? ` ×${s.clicks}` : ''}`;

export type ViewOptions = {
  inputs?: (id: string) => number; // số input ghim vào node
  run?: Map<string, { run: RunState; visits?: number }>; // trạng thái lần chạy đang xem
};

export function nodeViews(sc: Pick<Script, 'steps' | 'start' | 'edges'>, opt: ViewOptions = {}): FlowNodeView[] {
  const ok = reachable(sc.edges, sc.start);
  const refs = new Map(sc.steps.map((s) => [s.id, s.ref]));
  const refOf = (id: string) => refs.get(id) ?? '?';
  return sc.steps.map((s) => {
    const n = opt.inputs?.(s.id) ?? 0;
    const r = opt.run?.get(s.id);
    return {
      id: s.id,
      ref: s.ref,
      icon: s.type,
      title: nodeTitle(s),
      detail: isControl(s.type) ? controlDetail(s) : `${fmtPos(s, refOf)} · ${numVarOf(s, 'delayAfter') ?? fmtMs(s.delayAfter)}`,
      ports: portViews(s),
      enabled: s.enabled,
      chips: [...(n ? [`${n} input`] : []), ...(s.retry?.count ? [`thử lại ${s.retry.count}`] : []),
        ...(s.human && usesHuman(s) ? [`giống người: ${humanLabel(s.human)}`] : [])],
      warn: warningOf(s, ok.has(s.id), refs),
      breakpoint: s.breakpoint,
      run: r?.run,
      visits: r?.visits,
    };
  });
}

// Cảnh báo khi soạn (không chặn lưu; chạy thì Go kiểm tra đủ): node không đi tới được, click theo ảnh thiếu nguồn,
// node điều khiển thiếu dữ liệu.
function warningOf(s: Step, reachable: boolean, refs: Map<string, string>): string | undefined {
  if (!reachable) return 'Không có đường nào đi tới node này – sẽ không bao giờ chạy';
  if (s.type === 'clickImage' && s.imageSource !== 'search' && !refs.has(s.imageFrom ?? '')) return 'Chưa chọn node lấy chỗ thấy (Chờ ảnh / Chờ chữ)';
  if (s.type === 'if' && s.cond) return condMissing(s.cond);
  if ((s.type === 'switch' || s.type === 'switchScreen') && !s.cases?.length) return 'Chưa có case nào';
  if (s.type === 'switchScreen') {
    for (const c of s.cases ?? []) {
      const m = c.cond ? condMissing(c.cond) : 'chưa có điều kiện';
      if (m) return `Case ${c.label || c.key}: ${m.toLowerCase()}`;
    }
  }
  return undefined;
}

// Điều kiện thiếu dữ liệu (chạy sẽ báo lỗi) – undefined = đủ.
export function condMissing(c: Cond): string | undefined {
  switch (c.kind) {
    case 'text': return c.findText?.trim() ? undefined : 'Chưa nhập chữ cần thấy';
    case 'image': return c.template ? undefined : 'Chưa có ảnh mẫu';
    case 'pixel': return c.points?.[0] ? undefined : 'Chưa chọn điểm màu';
    case 'all': case 'any':
      return c.conds?.length ? c.conds.map(condMissing).find(Boolean) : 'Nhóm chưa có điều kiện nào';
    default: return undefined;
  }
}

const OP_TEXT: Record<string, string> = { eq: '=', ne: '≠', contains: 'chứa', gt: '>', lt: '<' };
// Một dòng mô tả điều kiện: thấy chữ “Gửi”, {{x}} > 3, 2 điều kiện (và)…
export function condText(c: Cond): string {
  const neg = c.not ? 'không ' : '';
  switch (c.kind) {
    case 'text': return `${neg}thấy chữ “${c.findText ?? ''}”`;
    case 'image': return c.template ? `${neg}thấy ảnh (≥ ${c.threshold ?? 90}%)` : 'chưa có ảnh mẫu';
    case 'pixel': return `${c.not ? 'khác' : 'có'} màu ${c.color ?? ''}`;
    case 'value': {
      const t = c.op === 'empty' ? `${c.left ?? ''} trống` : `${c.left ?? ''} ${OP_TEXT[c.op ?? 'eq']} ${c.right ?? ''}`;
      return c.not ? `không (${t})` : t;
    }
    default: return `${c.not ? 'không: ' : ''}${c.conds?.length ?? 0} điều kiện (${c.kind === 'all' ? 'và' : 'hoặc'})`;
  }
}

const sec = (ms: number) => `${(ms / 1000).toLocaleString('vi-VN', { maximumFractionDigits: 1 })} s`;
// Dòng chi tiết trên thẻ node điều khiển.
function controlDetail(s: Step): string {
  const n = s.cases?.length ?? 0;
  switch (s.type) {
    case 'if': return `${s.cond ? condText(s.cond) : '—'}${s.timeoutMs ? ` · chờ ${sec(s.timeoutMs)}` : ''}`;
    case 'switch': return `${s.value || '—'} · ${n} case`;
    case 'switchScreen': return `chờ tối đa ${sec(s.timeoutMs ?? 0)} · ${n} case`;
    case 'counter': return `tối đa ${s.max ?? 0} lần`;
    case 'end': return s.endStatus === 'fail' ? `thất bại${s.message ? ` · ${s.message}` : ''}` : `thành công${s.message ? ` · ${s.message}` : ''}`;
    case 'setVar': return `vars.${s.varName || '?'} ${s.varOp === 'add' ? '+=' : '='} ${s.value ?? ''}`;
    default: return '';
  }
}

// Output node trả về khi thành công – khớp Outputs trong steptypes.go. Node sau đọc bằng {{ref.output.key}}.
export type OutputField = { key: string; desc: string };
const FOUND: OutputField[] = [
  { key: 'x', desc: 'tâm chỗ thấy – X' }, { key: 'y', desc: 'tâm chỗ thấy – Y' }, { key: 'w', desc: 'rộng' }, { key: 'h', desc: 'cao' },
  { key: 'score', desc: '% khớp' }, { key: 'text', desc: 'chữ tìm thấy (tìm theo chữ)' },
];
const WEB_FOUND: OutputField[] = [
  { key: 'x', desc: 'tâm phần tử – X' }, { key: 'y', desc: 'tâm phần tử – Y' }, { key: 'w', desc: 'rộng' }, { key: 'h', desc: 'cao' },
  { key: 'text', desc: 'chữ của phần tử' }, { key: 'selector', desc: 'cách tìm đã khớp' },
];
// Output mọi node có (CommonOutputs trong steptypes.go, engine điền ở execNode).
const COMMON: OutputField[] = [{ key: 'startedAt', desc: 'giờ bắt đầu' }, { key: 'durationMs', desc: 'thời gian chạy (ms)' }];
export const outputsOf = (s: Step): OutputField[] => [...ownOutputs(s), ...COMMON];
const MOUSE_AT: OutputField[] = [{ key: 'x', desc: 'X chuột dừng' }, { key: 'y', desc: 'Y chuột dừng' }];
function ownOutputs(s: Step): OutputField[] {
  switch (s.type) {
    case 'click': case 'double': case 'drag': return MOUSE_AT;
    case 'hold': return [...MOUSE_AT, { key: 'ms', desc: 'thời gian đã giữ (ms)' }];
    case 'wait': return [{ key: 'ms', desc: 'thời gian đã chờ (ms, gồm ngẫu nhiên)' }];
    case 'key': return [{ key: 'keys', desc: 'tổ hợp đã bấm, vd cmd+l' }];
    case 'scroll': return [...MOUSE_AT, { key: 'direction', desc: 'hướng đã cuộn' }, { key: 'amount', desc: 'số nấc đã cuộn' }];
    case 'type': return [{ key: 'text', desc: 'chữ đã gõ (sau khi thay biến)' }];
    case 'open': return [{ key: 'target', desc: 'URL / app đã mở' }, { key: 'app', desc: 'mở bằng app (trống = mặc định)' }];
    case 'waitImage': return FOUND;
    case 'clickImage': return [{ key: 'clickX', desc: 'X đã click' }, { key: 'clickY', desc: 'Y đã click' }, ...(s.imageSource === 'search' ? FOUND : [])];
    case 'waitPixel': return [{ key: 'color', desc: 'màu đọc được' }];
    case 'system': return [{ key: 'volume', desc: 'âm lượng sau khi chạy (%)' }, { key: 'muted', desc: 'đang tắt tiếng (true / false)' }];
    case 'file': return s.fileOp === 'trash' ? [{ key: 'count', desc: 'số file đã xoá' }, { key: 'files', desc: 'tên các file đã xoá' }]
      : s.fileOp === 'emptyTrash' ? [] : s.fileOp === 'row' ? [{ key: 'path', desc: 'đường dẫn file' }, { key: 'rows', desc: 'số dòng sau khi thêm' }]
      : s.fileOp === 'mkdir' ? [{ key: 'path', desc: 'đường dẫn thư mục' }] : [{ key: 'path', desc: 'đường dẫn file' }, { key: 'bytes', desc: 'số byte đã ghi' }];
    case 'http': return [{ key: 'status', desc: 'mã HTTP, vd 200' }, { key: 'json', desc: 'JSON trả về – đọc sâu json.data.id' }, { key: 'body', desc: 'phản hồi dạng chữ' }, { key: 'headers', desc: 'header trả về (tên chữ thường)' }];
    case 'clipboard': return s.clipOp === 'image' ? [{ key: 'path', desc: 'file ảnh đã chép' }] : [{ key: 'text', desc: s.clipOp === 'write' ? 'chữ đã ghi' : 'chữ đọc được' }];
    case 'screenshot': return [{ key: 'path', desc: 'đường dẫn file ảnh' }, { key: 'x', desc: 'vùng đã chụp – X trái' }, { key: 'y', desc: 'vùng đã chụp – Y trên' },
      { key: 'w', desc: 'rộng' }, { key: 'h', desc: 'cao' }];
    case 'webWait': return s.webGone ? [] : WEB_FOUND;
    case 'webClick': return [{ key: 'clickX', desc: 'X đã click' }, { key: 'clickY', desc: 'Y đã click' }, ...WEB_FOUND];
    case 'webType': return [{ key: 'text', desc: 'chữ đã gõ (sau khi thay biến)' }, { key: 'clickX', desc: 'X đã click' }, { key: 'clickY', desc: 'Y đã click' },
      ...WEB_FOUND.filter((f) => f.key !== 'text')];
    case 'webRead': return [{ key: 'text', desc: 'giá trị đọc được' }, { key: 'count', desc: 'số phần tử khớp' }, { key: 'selector', desc: 'cách tìm đã khớp' }];
    case 'webScroll': return [...WEB_FOUND, ...(s.scrollMode === 'until' ? [{ key: 'scrolls', desc: 'số lần đã lăn' }] : [])];
    case 'if': return [{ key: 'result', desc: 'true / false' }, ...FOUND.map((f) => ({ ...f, desc: `${f.desc} (điều kiện ảnh / chữ đúng)` }))];
    case 'switch': return [{ key: 'case', desc: 'case khớp (default nếu không)' }, { key: 'value', desc: 'giá trị đã xét' }];
    case 'switchScreen': return [{ key: 'case', desc: 'case hiện trước (none nếu không)' }, ...FOUND];
    case 'counter': return [{ key: 'count', desc: 'số lần đã đi qua' }];
    case 'setVar': return [{ key: 'name', desc: 'tên biến' }, { key: 'value', desc: 'giá trị sau khi đặt' }];
    default: return [];
  }
}

// Biến đọc kết quả node khác – gợi ý cho ô chữ {{…}} (khớp lookup trong vars.go). Node khác `self`, theo thứ tự hiển thị;
// mỗi node: output khi thành công, rồi status / error.message.
export type NodeVar = { key: string; label: string };
export function nodeVarOptions(sc: Pick<Script, 'steps' | 'start' | 'edges'>, self: string): NodeVar[] {
  const byId = new Map(sc.steps.map((s) => [s.id, s]));
  return flowOrder(sc).flatMap((id) => {
    const s = byId.get(id);
    if (!s || id === self || !s.ref) return [];
    const name = nodeTitle(s);
    return [
      ...outputsOf(s).map((o) => ({ key: `${s.ref}.output.${o.key}`, label: `${name} – ${o.desc}` })),
      { key: `${s.ref}.status`, label: `${name} – success / fail / skipped` },
      { key: `${s.ref}.error.message`, label: `${name} – lời báo lỗi` },
    ];
  });
}

// Node có output x, y (chỗ thấy / chỗ chuột dừng) theo thứ tự flow, trừ self – nút chọn nhanh toạ độ từ biến (F5b).
export function coordSources(sc: Pick<Script, 'steps' | 'start' | 'edges'>, self: string): { ref: string; label: string }[] {
  const byId = new Map(sc.steps.map((s) => [s.id, s]));
  return flowOrder(sc).flatMap((id) => {
    const s = byId.get(id);
    if (!s || id === self || !s.ref) return [];
    const keys = outputsOf(s).map((o) => o.key);
    return keys.includes('x') && keys.includes('y') ? [{ ref: s.ref, label: `${s.ref} · ${nodeTitle(s)}` }] : [];
  });
}
