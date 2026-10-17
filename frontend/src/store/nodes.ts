// Node và cạnh của flow (bước 20): thêm / sửa / xoá / nhân bản node, nối / cắt cạnh, node bắt đầu, vị trí trên canvas.
// Phép trên đồ thị là hàm thuần ở flow/edit.ts – ở đây chỉ ghép với state, hoàn tác và thông báo.
import { anchorPoint } from '../anchor';
import { connect, disconnect, InsertAt, insertNode as insertInto, nextRef, pruneEdges, removeEdge as dropEdge, removeNode } from '../flow/edit';
import { flowOrder } from '../flow/model';
import { placeNode, Positions } from '../flow/layout';
import { makeStep, Point, Script, Step, StepType } from '../types';
import { getEditor, patchScriptSilently, setEditor, setNotice, updateScript } from './editor';
import { boundsOf, centerPoint, refOf, toScreen } from './geometry';

export const newId = () => `${Date.now().toString(36)}${Math.random().toString(36).slice(2, 6)}`;

const script = () => getEditor().script;
const setSteps = (steps: Step[]) => updateScript({ steps });

// ---- Chỗ node mới sẽ nối vào ----
// Đặt trước khi mở bảng thêm node / overlay chọn điểm: thả cổng ra ở chỗ trống (kèm vị trí thả), hoặc Ctrl+K khi đang chọn
// node (sau cổng success của nó). null = node mới đứng riêng. Đặt nhiều điểm liên tiếp trên overlay thì node sau nối vào
// node vừa thêm. clearInsertAt khi đóng bảng / overlay.
let insertAt: (InsertAt & { pos?: Point }) | null = null;
export const setInsertAt = (at: (InsertAt & { pos?: Point }) | null) => {
  insertAt = at;
};
export const clearInsertAt = () => setInsertAt(null);
// Mặc định của bảng thêm node khi không thả từ cổng: sau cổng success của node đang chọn.
export const insertAfterSelected = () => {
  const sel = getEditor().selected;
  setInsertAt(script().steps.some((s) => s.id === sel) ? { from: sel, port: 'success' } : null);
};
// Đang chèn vào đâu – dòng ghi trên bảng thêm node.
export const insertTarget = () => insertAt;

// Thêm node vào flow theo insertAt, đặt mã + vị trí, chọn nó và báo. Node sau (đặt điểm liên tiếp) nối vào node này.
export function insertNode(step: Step, notice: (ref: string) => string) {
  const sc = script();
  const node: Step = { ...step, ref: step.ref || nextRef(sc.steps), pos: step.pos ?? insertAt?.pos ?? placeNode(sc.steps, insertAt?.from) };
  updateScript(insertInto(sc, node, insertAt));
  if (insertAt) insertAt = { from: node.id, port: 'success' };
  setEditor({ selected: node.id, notice: notice(node.ref) });
}

// Thêm node không cần chọn điểm (chờ, gõ chữ, phím tắt, cuộn tại chuột, mở, chờ ảnh…) – giá trị mặc định giống NewStep (model.go).
export function addStep(type: Extract<StepType, 'wait' | 'type' | 'key' | 'scroll' | 'open' | 'waitImage' | 'clickImage' | 'waitPixel'>) {
  const defaults: Record<typeof type, Partial<Step>> = {
    wait: { clicks: 0, delayAfter: 1000 },
    type: { clicks: 0, delayAfter: 500 },
    key: { clicks: 0, keys: ['enter'], delayAfter: 500 },
    scroll: { clicks: 0, direction: 'down', amount: 3, delayAfter: 500 },
    open: { clicks: 0, target: 'https://', delayAfter: 2000 },
    waitImage: { clicks: 0, threshold: 90, timeoutMs: 10000, intervalMs: 250, delayAfter: 300 },
    clickImage: { imageSource: 'node', imageFrom: nearestFinder(), threshold: 90, timeoutMs: 10000, intervalMs: 250, delayAfter: 500 },
    waitPixel: { clicks: 0, points: [centerPoint(script().region)], color: '#000000', tolerance: 10, timeoutMs: 5000, intervalMs: 100, delayAfter: 300 },
  };
  insertNode(makeStep({ id: newId(), type, ...defaults[type] }), (ref) => `Đã thêm ${ref}`);
  clearInsertAt();
}

// Node tìm ảnh / chữ (có "chỗ thấy" cho click theo ảnh): chờ ảnh / chữ, click theo ảnh tự tìm.
export const findsImage = (s: Step) => s.type === 'waitImage' || (s.type === 'clickImage' && s.imageSource === 'search');

// Nguồn mặc định cho click theo ảnh mới: node đang nối vào nếu nó tìm ảnh, không thì node tìm ảnh gần nhất theo thứ tự.
function nearestFinder(): string | undefined {
  const sc = script();
  const from = sc.steps.find((s) => s.id === insertAt?.from);
  if (from && findsImage(from)) return from.id;
  const order = flowOrder(sc);
  const upto = from ? order.indexOf(from.id) : order.length - 1;
  const byId = new Map(sc.steps.map((s) => [s.id, s]));
  return order.slice(0, upto + 1).reverse().find((id) => findsImage(byId.get(id)!));
}

export function updateStep(id: string, patch: Partial<Step>) {
  setSteps(script().steps.map((s) => (s.id === id ? { ...s, ...patch } : s)));
}

export const toggleStep = (id: string) => {
  const s = script().steps.find((x) => x.id === id);
  if (s) updateStep(id, { enabled: !s.enabled });
};

// Nhân bản: id + mã mới, đặt lệch xuống dưới, chưa nối cạnh nào.
export function duplicateStep(id: string) {
  const sc = script();
  const s = sc.steps.find((x) => x.id === id);
  if (!s) return;
  const copy: Step = {
    ...structuredClone(s), id: newId(), ref: nextRef(sc.steps),
    pos: s.pos && { x: s.pos.x + 24, y: s.pos.y + 24 },
  };
  setSteps([...sc.steps, copy]);
  setEditor({ selected: copy.id, notice: `Đã nhân bản ${s.ref} thành ${copy.ref} (chưa nối)` });
}

// Xoá node kèm cạnh vào / ra (hoàn tác bằng ⌘Z / Ctrl+Z).
export function removeStep(id: string) {
  const sc = script();
  const ref = refOf(id);
  if (!sc.steps.some((s) => s.id === id)) return;
  const g = removeNode({ nodes: sc.steps, edges: sc.edges, start: sc.start }, id);
  updateScript({ steps: g.nodes, edges: g.edges, start: g.start });
  setEditor({ selected: '', notice: `Đã xoá ${ref} – ${navigator.platform.startsWith('Mac') ? '⌘' : 'Ctrl+'}Z để hoàn tác` });
}

// Đổi loại node, bổ sung dữ liệu loại mới cần: điểm (click/double/hold 1, kéo thả 2 – B lệch A 100 px), số click ≥ 1.
// Cạnh của cổng loại mới không có bị bỏ.
export function changeStepType(id: string, type: StepType) {
  const sc = script();
  const s = sc.steps.find((x) => x.id === id);
  if (!s || s.type === type) return;
  const points = s.points.slice(0, type === 'scroll' ? 1 : undefined).map((p) => ({ ...p })); // cuộn: điểm tuỳ chọn, tối đa 1
  const need = type === 'drag' ? 2 : type === 'click' || type === 'double' || type === 'hold' || type === 'waitPixel' ? 1 : 0;
  if (need >= 1 && points.length === 0) points.push({ x: 0, y: 0 });
  if (need === 2 && points.length === 1) {
    const a = toScreen(points[0], sc.region);
    points.push(anchorPoint({ x: a.x + 100, y: a.y }, boundsOf(sc.region), points[0].ax, points[0].ay));
  }
  const next: Step = {
    ...s,
    type,
    points,
    clicks: Math.max(s.clicks, 1),
    holdMs: type === 'hold' && !s.holdMs ? 1000 : s.holdMs,
    dragMs: type === 'drag' && !s.dragMs ? 400 : s.dragMs,
    keys: type === 'key' && s.keys.length === 0 ? ['enter'] : s.keys,
    direction: type === 'scroll' ? s.direction ?? 'down' : s.direction,
    amount: type === 'scroll' && !s.amount ? 3 : s.amount,
    target: type === 'open' && !s.target ? 'https://' : s.target,
    threshold: s.threshold ?? 90,
    timeoutMs: s.timeoutMs ?? (type === 'waitPixel' ? 5000 : 10000),
    intervalMs: s.intervalMs ?? (type === 'waitPixel' ? 100 : 250),
    imageSource: type === 'clickImage' ? s.imageSource ?? 'node' : s.imageSource,
    color: type === 'waitPixel' ? s.color ?? '#000000' : s.color,
    tolerance: type === 'waitPixel' ? s.tolerance ?? 10 : s.tolerance,
  };
  const steps = sc.steps.map((x) => (x.id === id ? next : x));
  updateScript({ steps, edges: pruneEdges(steps, sc.edges) });
}

// ---- Cạnh, node bắt đầu, vị trí ----

export function connectPort(from: string, port: string, to: string) {
  const sc = script();
  const edges = connect(sc.edges, from, port, to);
  if (edges !== sc.edges) updateScript({ edges });
}

export const disconnectPort = (from: string, port: string) => updateScript({ edges: disconnect(script().edges, from, port) });
export const removeEdge = (id: string) => updateScript({ edges: dropEdge(script().edges, id) });

export function setStart(id: string) {
  updateScript({ start: id });
  setNotice(`${refOf(id)} là node bắt đầu`);
}

export const moveNode = (id: string, pos: Point) => setSteps(script().steps.map((s) => (s.id === id ? { ...s, pos } : s)));

// Ghi vị trí từ xếp tự động. silent: lần đầu mở kịch bản chưa có vị trí – không tính là sửa.
export function applyLayout(pos: Positions, silent = false) {
  const patch: Partial<Script> = { steps: script().steps.map((s) => (pos[s.id] ? { ...s, pos: pos[s.id] } : s)) };
  if (silent) patchScriptSilently(patch);
  else {
    updateScript(patch);
    setNotice('Đã sắp xếp lại các node');
  }
}
