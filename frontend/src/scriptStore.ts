import { useSyncExternalStore } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { ScreenBounds } from '../wailsjs/go/main/App';
import { Anchor, makeStep, Point, Rect, Script, sampleSteps, Step, STEP_LABEL, StepPoint, StepType } from './types';
import { anchorPoint, fmtPoint, insideBounds, nearestAnchor, resolvePoint } from './anchor';
import { inWails, openOverlay, setWindowMode } from './windowMode';
import { apply, IDENTITY, Mat, rectToRect } from './transform';

// Store kịch bản đang soạn – dùng chung giữa MainScreen và OverlayScreen (hai màn không cùng mount).
// Lưu / mở file và bản nháp tự lưu: scriptFile.ts (bước 21).
type EditorState = {
  script: Script;
  path: string | null; // file kịch bản đang mở, null = chưa lưu lần nào
  dirty: boolean; // có thay đổi chưa lưu xuống file – mọi thay đổi script qua updateScript đều bật cờ này
  selected: string; // id bước đang chọn trong danh sách
  notice: string; // thông báo ngắn gần nhất (đã thêm điểm, điểm ngoài vùng…)
  dragA: Point | null; // điểm A (toạ độ màn hình) đã chọn, đang chờ điểm B của kéo thả
  deleted: { step: Step; index: number } | null; // bước vừa xoá – để hoàn tác
  countdown: boolean; // đếm ngược 3 giây trước khi chạy
  multi: string[]; // điểm đang chọn nhiều (bước 30), khoá pointKey(id, index)
  past: Script[]; // hoàn tác (cũ nhất trước)
  future: Script[]; // làm lại
  screen: Rect; // khung màn hình chính (robotgo) – "toàn màn hình" neo theo khung này, khớp mainScreen (anchor.go)
};

let state: EditorState = {
  // app thật: kịch bản trống cho tới khi scriptFile.ts nạp bản nháp; trình duyệt: dữ liệu mẫu
  script: { version: 1, name: 'Kịch bản 1', region: null, loops: inWails ? 1 : 10, maxMinutes: 0, steps: inWails ? [] : sampleSteps },
  path: null,
  dirty: false,
  selected: inWails ? '' : '3',
  notice: '',
  dragA: null,
  deleted: null,
  countdown: true,
  multi: [],
  past: [],
  future: [],
  screen: { x: 0, y: 0, w: window.screen.width, h: window.screen.height }, // trình duyệt: tạm lấy cỡ màn hình
};
const listeners = new Set<() => void>();

const set = (patch: Partial<EditorState>) => {
  state = { ...state, ...patch };
  listeners.forEach((l) => l());
};

export const subscribe = (l: () => void) => {
  listeners.add(l);
  return () => listeners.delete(l);
};

export const getEditor = () => state;
export const useEditor = () => useSyncExternalStore(subscribe, () => state);
export const useScript = () => useEditor().script;

// ---- Hoàn tác / làm lại (bước 30) ----
// Mỗi thay đổi kịch bản lưu bản trước đó; các thay đổi sát nhau (gõ số, giữ phím mũi tên) gộp thành một lần hoàn tác.
const HISTORY_MAX = 100;
const MERGE_MS = 500;
let lastEdit = 0;

export const updateScript = (patch: Partial<Script>) => {
  const now = Date.now();
  const past = now - lastEdit < MERGE_MS && state.past.length ? state.past : [...state.past, state.script].slice(-HISTORY_MAX);
  lastEdit = now;
  set({ script: { ...state.script, ...patch }, dirty: true, past, future: [] });
};

// Sau hoàn tác / làm lại: bước đang chọn và điểm chọn nhiều phải còn tồn tại.
const restore = (script: Script, past: Script[], future: Script[], notice: string) => {
  lastEdit = 0;
  const ids = new Set(script.steps.map((s) => s.id));
  set({
    script, past, future, dirty: true, notice, dragA: null,
    selected: ids.has(state.selected) ? state.selected : script.steps[0]?.id ?? '',
    multi: state.multi.filter((k) => ids.has(k.split(':')[0])),
  });
};

export function undo() {
  const prev = state.past.at(-1);
  if (!prev) return setNotice('Không còn gì để hoàn tác');
  restore(prev, state.past.slice(0, -1), [state.script, ...state.future], 'Đã hoàn tác');
}

export function redo() {
  const next = state.future[0];
  if (!next) return setNotice('Không còn gì để làm lại');
  restore(next, [...state.past, state.script], state.future.slice(1), 'Đã làm lại');
}

// Cmd/Ctrl+Z hoàn tác, + Shift + Z (hoặc Ctrl+Y) làm lại – mọi màn; trong ô nhập để trình duyệt tự hoàn tác chữ.
window.addEventListener('keydown', (e) => {
  if (!(navigator.platform.startsWith('Mac') ? e.metaKey : e.ctrlKey) || e.altKey) return;
  if ((e.target as HTMLElement).closest?.('input, textarea, select')) return;
  const k = e.key.toLowerCase();
  if (k === 'z') (e.shiftKey ? redo : undo)();
  else if (k === 'y' && !e.shiftKey) redo();
  else return;
  e.preventDefault();
});

// Thay cả kịch bản (mở file, kịch bản mới, khôi phục nháp): bỏ chọn / hoàn tác / điểm A đang chờ của kịch bản cũ.
export function replaceScript(script: Script, path: string | null, dirty: boolean, notice = '') {
  lastEdit = 0;
  set({ script, path, dirty, selected: script.steps[0]?.id ?? '', deleted: null, dragA: null, notice, multi: [], past: [], future: [] });
}

// Sau khi lưu: đường dẫn + tên mới (Lưu thành đổi tên theo file), giữ nguyên phần còn lại đang soạn.
// stillDirty: có sửa thêm trong lúc hộp thoại lưu đang mở.
export const markSaved = (path: string, name: string, notice: string, stillDirty: boolean) =>
  set({ script: { ...state.script, name }, path, dirty: stillDirty, notice });
export const selectStep = (id: string) => set({ selected: id });
export const setNotice = (notice: string) => set({ notice });
export const setCountdown = (countdown: boolean) => set({ countdown });

const setSteps = (steps: Step[]) => updateScript({ steps });

// ---- Danh sách bước (bước 14) ----

const indexOf = (id: string) => state.script.steps.findIndex((s) => s.id === id);

export function toggleStep(id: string) {
  setSteps(state.script.steps.map((s) => (s.id === id ? { ...s, enabled: !s.enabled } : s)));
}

// Nhân bản ngay sau bước gốc và chọn bản sao.
export function duplicateStep(id: string) {
  const steps = state.script.steps;
  const i = indexOf(id);
  if (i < 0) return;
  const copy: Step = { ...steps[i], id: newId(), points: steps[i].points.map((p) => ({ ...p })), modifiers: [...steps[i].modifiers], keys: [...steps[i].keys] };
  setSteps([...steps.slice(0, i + 1), copy, ...steps.slice(i + 1)]);
  set({ selected: copy.id, notice: `Đã nhân bản bước ${i + 1} thành bước ${i + 2}` });
}

// Xoá và chọn bước kế tiếp (hoặc bước trước nếu xoá bước cuối). Giữ lại để hoàn tác.
export function removeStep(id: string) {
  const steps = state.script.steps;
  const i = indexOf(id);
  if (i < 0) return;
  const rest = steps.filter((s) => s.id !== id);
  setSteps(rest);
  set({
    selected: rest[Math.min(i, rest.length - 1)]?.id ?? '',
    notice: `Đã xoá bước ${i + 1}`,
    deleted: { step: steps[i], index: i },
  });
}

export function undoDelete() {
  const d = state.deleted;
  if (!d) return;
  const steps = state.script.steps;
  setSteps([...steps.slice(0, d.index), d.step, ...steps.slice(d.index)]);
  set({ selected: d.step.id, notice: `Đã khôi phục bước ${d.index + 1}`, deleted: null });
}

export function moveStep(from: number, to: number) {
  const steps = [...state.script.steps];
  const [s] = steps.splice(from, 1);
  steps.splice(to, 0, s);
  setSteps(steps);
}

// ---- Form cấu hình bước (bước 15) ----

export function updateStep(id: string, patch: Partial<Step>) {
  setSteps(state.script.steps.map((s) => (s.id === id ? { ...s, ...patch } : s)));
}

// Đổi loại bước, bổ sung dữ liệu loại mới cần: điểm (click/double/hold 1, kéo thả 2 – B lệch A 100 px), số click ≥ 1.
export function changeStepType(id: string, type: StepType) {
  const s = state.script.steps.find((x) => x.id === id);
  if (!s || s.type === type) return;
  const points = s.points.map((p) => ({ ...p }));
  const need = type === 'drag' ? 2 : type === 'click' || type === 'double' || type === 'hold' ? 1 : 0;
  if (need >= 1 && points.length === 0) points.push({ x: 0, y: 0 });
  if (need === 2 && points.length === 1) {
    const a = toScreen(points[0], state.script.region);
    points.push(anchorPoint({ x: a.x + 100, y: a.y }, boundsOf(state.script.region), points[0].ax, points[0].ay));
  }
  updateStep(id, {
    type,
    points,
    clicks: Math.max(s.clicks, 1),
    holdMs: type === 'hold' && !s.holdMs ? 1000 : s.holdMs,
    dragMs: type === 'drag' && !s.dragMs ? 400 : s.dragMs,
    keys: type === 'key' && s.keys.length === 0 ? ['enter'] : s.keys,
  });
}

// Bước mới từ bảng thêm bước (Ctrl+K) chèn ngay sau bước này; null = thêm vào cuối (F6, đặt điểm thường).
// Đặt nhiều điểm liên tiếp trên overlay thì mỗi điểm chèn sau điểm vừa thêm. cancelPick (đóng overlay) xoá.
let insertAfter: string | null = null;
// Loại bước tạo khi đặt điểm trên overlay (bảng thêm bước chọn Double click / Giữ chuột).
type PointType = Extract<StepType, 'click' | 'double' | 'hold'>;
let pointType: PointType = 'click';
export const placingType = () => pointType; // gợi ý trên overlay

// Thêm bước vào kịch bản (sau insertAfter hoặc cuối), chọn nó và báo. Trả số thứ tự (từ 1).
function insertStep(step: Step, notice: (n: number) => string) {
  const steps = state.script.steps;
  const at = insertAfter ? indexOf(insertAfter) + 1 : 0;
  const i = at > 0 ? at : steps.length;
  setSteps([...steps.slice(0, i), step, ...steps.slice(i)]);
  if (insertAfter) insertAfter = step.id;
  set({ selected: step.id, notice: notice(i + 1) });
}

// Thêm bước không cần chọn điểm (chờ, gõ chữ, phím tắt) ngay sau bước đang chọn – giá trị mặc định giống NewStep (model.go).
export function addStep(type: Extract<StepType, 'wait' | 'type' | 'key'>) {
  const defaults: Record<typeof type, Partial<Step>> = {
    wait: { clicks: 0, delayAfter: 1000 },
    type: { clicks: 0, delayAfter: 500 },
    key: { clicks: 0, keys: ['enter'], delayAfter: 500 },
  };
  insertAfter = state.selected || null;
  insertStep(makeStep({ id: newId(), type, ...defaults[type] }), (n) => `Đã thêm bước ${n}`);
  insertAfter = null;
}

// Bảng thêm bước: bước cần vị trí → mở overlay chọn điểm, bước mới chèn sau bước đang chọn.
export function addPointStep(type: PointType | 'drag') {
  insertAfter = state.selected || null;
  if (type === 'drag') return openOverlay('drag');
  pointType = type;
  return openOverlay('point');
}

// Đổi vùng; điểm lưu theo neo nên mỗi điểm bám theo cạnh / góc nó neo. null = toàn màn hình.
export const setRegion = (region: Rect | null) => updateScript({ region });

// Khung tính toạ độ điểm: vùng, không có thì màn hình chính (bounds trong anchor.go).
export const boundsOf = (region: Rect | null): Rect => region ?? state.screen;

// Điểm đã lưu (neo theo vùng) → toạ độ màn hình thật.
export const toScreen = (p: StepPoint, region: Rect | null): Point => resolvePoint(p, boundsOf(region));

// Toạ độ màn hình → điểm lưu: like = điểm cũ để giữ neo (kéo marker, chọn lại), không có thì neo góc gần nhất.
// null nếu nằm ngoài vùng.
const toRelative = (p: Point, r: Rect | null, like?: StepPoint): StepPoint | null => {
  const b = boundsOf(r);
  if (!insideBounds(p, b)) return null;
  if (like) return anchorPoint(p, b, like.ax, like.ay);
  return state.script.anchorMode === 'start' ? anchorPoint(p, b, 'start', 'start') : nearestAnchor(p, b);
};

// Đổi neo điểm thứ index của bước, giữ nguyên chỗ trên màn hình (lưới chọn neo, nút %).
export function setAnchor(id: string, index: number, ax: Anchor, ay: Anchor) {
  const { steps, region } = state.script;
  const s = steps.find((x) => x.id === id);
  const p = s?.points[index];
  if (!s || !p) return;
  const q = anchorPoint(toScreen(p, region), boundsOf(region), ax, ay);
  updateStep(id, { points: s.points.map((x, i) => (i === index ? q : x)) });
}

const newId = () => `${Date.now().toString(36)}${Math.random().toString(36).slice(2, 6)}`;

// Bước đang được chọn lại điểm (nút ⌖ trong form); null = điểm mới sẽ thêm thành bước mới.
let pickTarget: string | null = null;
let pickStay = false; // chọn lại xong thì ở lại overlay (chế độ xem riêng một bước) thay vì về màn chính
// Overlay đang xem riêng một bước: click / F6 không tạo bước mới (trừ khi đang chọn lại điểm).
let focusLock = false;
export const setFocusLock = (on: boolean) => {
  focusLock = on;
};

// Chọn lại điểm ngay trong overlay đang mở (chế độ xem riêng) – không mở lại overlay, xong thì ở lại.
export function startPick(stepId: string) {
  pickTarget = stepId;
  pickStay = true;
  if (state.dragA) set({ dragA: null });
}
// Overlay đang ở công cụ "Kéo thả A → B": điểm mới đi theo cặp A, B thay vì thành bước click.
let dragMode = false;

export function pickPointFor(stepId: string) {
  pickTarget = stepId;
  pickStay = false;
  const step = state.script.steps.find((s) => s.id === stepId);
  return openOverlay(step?.type === 'drag' ? 'drag' : 'point');
}

export const pickingStep = () => pickTarget;

// Gọi khi đổi công cụ / đóng overlay; luôn bỏ điểm A đang chờ.
export const setDragMode = (on: boolean) => {
  dragMode = on;
  if (state.dragA) set({ dragA: null });
};

// Esc trên overlay: huỷ điểm A đang chờ. Trả false nếu không có gì để huỷ.
export const cancelDragA = () => {
  if (!state.dragA) return false;
  set({ dragA: null, notice: 'Đã huỷ điểm A' });
  return true;
};

export const cancelPick = () => {
  pickTarget = null;
  pickStay = false;
  insertAfter = null;
  pointType = 'click';
  setDragMode(false);
};

// Đặt một điểm (toạ độ màn hình thật) từ click trên overlay hoặc F6:
// đang chọn lại → sửa điểm của bước đó rồi về màn chính; không thì thêm bước click mới.
export function placePoint(screen: Point): boolean {
  const { script } = state;
  if (focusLock && !pickTarget) {
    setNotice('Đang xem riêng một bước – bấm “Chọn lại trên màn hình” để đặt điểm, hoặc “Xem tất cả”');
    return false;
  }
  const rel = toRelative(screen, script.region);
  if (!rel) {
    setNotice(`Điểm X ${screen.x} Y ${screen.y} nằm ngoài vùng – chọn trong vùng hoặc đổi về toàn màn hình`);
    return false;
  }
  const steps = script.steps;
  const target = pickTarget && steps.find((s) => s.id === pickTarget);
  // Kéo thả cần 2 điểm – trừ khi bước đang chọn lại bắt đầu "từ điểm trước đó" (chỉ cần B).
  const wantsPair = target ? target.type === 'drag' && !target.dragFromPrev : dragMode;

  // Kéo thả cần 2 điểm: lần đầu chỉ ghi nhớ A.
  if (wantsPair && !state.dragA) {
    set({ dragA: screen, notice: `Điểm A: ${fmtPoint(rel)} – chọn điểm B (Esc để huỷ)` });
    return true;
  }
  const a = wantsPair ? state.dragA! : null;
  if (a) set({ dragA: null });

  if (target) {
    pickTarget = null;
    if (a) movePoint(target.id, 0, a);
    movePoint(target.id, target.type === 'drag' ? 1 : 0, screen);
    if (pickStay) pickStay = false;
    else setWindowMode('main');
    return true;
  }
  if (a) {
    const relA = toRelative(a, script.region)!;
    const step = makeStep({ id: newId(), type: 'drag', points: [relA, rel], dragMs: 400, delayAfter: 500 });
    insertStep(step, (n) => `Đã thêm bước ${n}: kéo A ${fmtPoint(relA)} → B ${fmtPoint(rel)}`);
    return true;
  }
  const step = makeStep({ id: newId(), type: pointType, points: [rel], delayAfter: 500, holdMs: pointType === 'hold' ? 1000 : 0 });
  insertStep(step, (n) => `Đã thêm bước ${n}: ${STEP_LABEL[pointType].toLowerCase()} ${fmtPoint(rel)}`);
  return true;
}

// Dời điểm thứ index của bước id tới toạ độ màn hình screen (kéo marker, chọn lại điểm).
// Ngoài vùng → báo và giữ nguyên.
export function movePoint(id: string, index: number, screen: Point): boolean {
  const { steps, region } = state.script;
  const idx = steps.findIndex((s) => s.id === id);
  if (idx < 0) return false;
  const rel = toRelative(screen, region, steps[idx].points[index]);
  if (!rel) {
    setNotice(`Điểm X ${screen.x} Y ${screen.y} nằm ngoài vùng – giữ vị trí cũ`);
    return false;
  }
  const step = steps[idx];
  const label = step.type === 'drag' ? `${idx + 1}${index === 0 ? 'A' : 'B'}` : `${idx + 1}`;
  let value = rel;
  // Điểm B lưu dạng lệch so với điểm bắt đầu: đổi toạ độ màn hình sang (dx, dy).
  if (step.type === 'drag' && index === 1 && step.dragRelative) {
    const start = endpoints(steps, region).get(id)?.start;
    if (!start) {
      setNotice(`Bước ${idx + 1}: chưa xác định được điểm bắt đầu – giữ vị trí cũ`);
      return false;
    }
    value = { x: screen.x - start.x, y: screen.y - start.y };
  }
  setSteps(steps.map((s) => (s.id === id ? { ...s, points: s.points.map((p, i) => (i === index ? value : p)) } : s)));
  set({ selected: id, notice: `Đã đặt điểm ${label}: ${fmtPoint(value)}` });
  return true;
}

// ---- Kéo thả: từ điểm trước / điểm cuối tương đối ----

type Ends = { start: Point | null; end: Point | null; prevStep: number | null };

// Toạ độ màn hình (danh nghĩa, chưa tính ngẫu nhiên) điểm đầu / cuối của các bước chuột đang bật,
// lần theo thứ tự như engine: "điểm trước" = chỗ chuột dừng ở bước chuột bật gần nhất phía trước.
// prevStep: số thứ tự (từ 1) của bước cung cấp điểm trước cho bước kéo "từ điểm trước".
export function endpoints(steps: Step[], region: Rect | null): Map<string, Ends> {
  const out = new Map<string, Ends>();
  let last: Point | null = null;
  let lastStep: number | null = null;
  steps.forEach((s, i) => {
    if (!s.enabled || !s.points.length) return;
    if (s.type === 'click' || s.type === 'double' || s.type === 'hold') {
      const p = toScreen(s.points[0], region);
      out.set(s.id, { start: p, end: p, prevStep: null });
      last = p;
      lastStep = i + 1;
    } else if (s.type === 'drag' && s.points.length >= 2) {
      const start = s.dragFromPrev ? last : toScreen(s.points[0], region);
      const d = s.points[1];
      const end = s.dragRelative ? (start && { x: start.x + d.x, y: start.y + d.y }) : toScreen(d, region);
      out.set(s.id, { start, end, prevStep: s.dragFromPrev ? lastStep : null });
      if (end) {
        last = end;
        lastStep = i + 1;
      }
    }
  });
  return out;
}

// Bật/tắt "từ điểm trước": giữ nguyên chỗ trên màn hình – tắt thì A lấy đúng điểm trước đó.
export function setDragFromPrev(id: string, on: boolean) {
  const { steps, region } = state.script;
  const s = steps.find((x) => x.id === id);
  if (!s || s.type !== 'drag') return;
  const ends = endpoints(steps, region).get(id);
  const patch: Partial<Step> = { dragFromPrev: on };
  if (!on && ends?.start) {
    patch.points = [toRelative(ends.start, region, s.points[0]) ?? s.points[0], s.points[1]];
  }
  updateStep(id, patch);
}

// Bật/tắt "điểm cuối tương đối": đổi points[1] giữa toạ độ và độ lệch, giữ nguyên chỗ trên màn hình.
export function setDragRelative(id: string, on: boolean) {
  const { steps, region } = state.script;
  const s = steps.find((x) => x.id === id);
  if (!s || s.type !== 'drag' || s.dragRelative === on) return;
  const ends = endpoints(steps, region).get(id);
  let b = s.points[1];
  if (ends?.start && ends.end) {
    b = on
      ? { x: ends.end.x - ends.start.x, y: ends.end.y - ends.start.y }
      : toRelative(ends.end, region) ?? anchorPoint(ends.end, boundsOf(region), undefined, undefined);
  }
  updateStep(id, { dragRelative: on, points: [s.points[0], b] });
}

// ---- Chọn nhiều + biến đổi nhiều điểm (bước 30) ----

export const pointKey = (id: string, index: number) => `${id}:${index}`;
export const setMulti = (multi: string[]) => set({ multi });

// Điểm chọn được của một bước: A "từ điểm trước" tự tính nên không chọn / biến đổi.
export const selectableIndexes = (s: Step) =>
  s.points.length === 0 || s.type === 'wait' || s.type === 'type' || s.type === 'key' ? []
    : s.type === 'drag' ? (s.dragFromPrev ? [1] : [0, 1]) : [0];

// Danh sách bước: Cmd/Shift+click thêm / bỏ mọi điểm của bước vào chọn nhiều.
export function toggleMultiStep(id: string) {
  const s = state.script.steps.find((x) => x.id === id);
  if (!s) return;
  const keys = selectableIndexes(s).map((j) => pointKey(id, j));
  const has = keys.some((k) => state.multi.includes(k));
  set({ multi: has ? state.multi.filter((k) => !keys.includes(k)) : [...state.multi, ...keys] });
}

// Áp f (toạ độ màn hình → toạ độ màn hình) lên các điểm được chọn (pick), đổi khung tính từ vùng from sang to.
// Điểm thường: ghi lại theo neo hiện có trong khung mới. B "lệch so với A": véc-tơ đổi theo phần tuyến tính của f
// (f(A + Δ) − f(A), không cộng phần dời). Trả số điểm rơi ra ngoài khung mới – khác 0 thì người gọi không nên áp dụng.
function mapPoints(steps: Step[], from: Rect | null, to: Rect | null, f: (p: Point) => Point, pick: (id: string, j: number) => boolean) {
  const ends = endpoints(steps, from);
  const nb = boundsOf(to);
  let outside = 0;
  const round = (p: Point) => ({ x: Math.round(p.x), y: Math.round(p.y) });
  const out = steps.map((s) => {
    if (!selectableIndexes(s).some((j) => pick(s.id, j))) return s;
    const points = s.points.map((p, j) => {
      if (!pick(s.id, j) || !selectableIndexes(s).includes(j)) return p;
      if (s.type === 'drag' && j === 1 && s.dragRelative) {
        const a = ends.get(s.id)?.start;
        if (!a) return p;
        const fa = f(a), fb = f({ x: a.x + p.x, y: a.y + p.y });
        return round({ x: fb.x - fa.x, y: fb.y - fa.y });
      }
      const q = round(f(toScreen(p, from)));
      if (!insideBounds(q, nb)) outside++;
      return anchorPoint(q, nb, p.ax, p.ay);
    });
    return { ...s, points };
  });
  return { steps: out, outside };
}

// Kịch bản như sẽ thành nếu áp m lên các điểm đang chọn – overlay vẽ xem trước bằng đúng logic ghi thật.
export function transformedSteps(m: Mat) {
  const { steps, region } = state.script;
  const sel = new Set(state.multi);
  return mapPoints(steps, region, region, (p) => apply(m, p), (id, j) => sel.has(pointKey(id, j)));
}

// Biến đổi các điểm đang chọn nhiều bằng ma trận m (toạ độ màn hình). Điểm bị đẩy ra ngoài vùng → báo, không áp dụng.
export function transformSelected(m: Mat, what: string): boolean {
  const n = state.multi.length;
  if (!n) return false;
  const r = transformedSteps(m);
  if (r.outside) {
    setNotice(`${r.outside} điểm sẽ ra ngoài vùng – không áp dụng ${what}`);
    return false;
  }
  setSteps(r.steps);
  setNotice(`Đã ${what} ${n} điểm – ${navigator.platform.startsWith('Mac') ? '⌘' : 'Ctrl+'}Z để hoàn tác`);
  return true;
}

// Đặt toạ độ màn hình mới cho từng điểm đang chọn (căn chỉnh / chia đều / lật): next[khoá] → điểm mới.
// Hai lượt: ghi điểm thường trước, rồi tính lại B "lệch so với A" theo điểm bắt đầu mới (A có thể vừa dời).
export function moveSelectedTo(next: Map<string, Point>, what: string) {
  const { steps, region } = state.script;
  const b = boundsOf(region);
  const round = (p: Point) => ({ x: Math.round(p.x), y: Math.round(p.y) });
  const isRel = (s: Step, j: number) => s.type === 'drag' && j === 1 && s.dragRelative;
  let outside = 0;
  const first = steps.map((s) => ({
    ...s,
    points: s.points.map((p, j) => {
      const q = next.get(pointKey(s.id, j));
      if (!q || isRel(s, j)) return p;
      if (!insideBounds(round(q), b)) outside++;
      return anchorPoint(round(q), b, p.ax, p.ay);
    }),
  }));
  const ends = endpoints(first, region);
  const out = first.map((s) => ({
    ...s,
    points: s.points.map((p, j) => {
      const q = next.get(pointKey(s.id, j));
      const start = ends.get(s.id)?.start;
      if (!q || !isRel(s, j) || !start) return p;
      if (!insideBounds(round(q), b)) outside++;
      return { x: Math.round(q.x) - start.x, y: Math.round(q.y) - start.y };
    }),
  }));
  if (outside) return setNotice(`${outside} điểm sẽ ra ngoài vùng – không ${what}`);
  setSteps(out);
  setNotice(`Đã ${what} ${next.size} điểm`);
}

// ---- Đổi vùng: cách xử lý các điểm (bước 30) ----
// anchor: giữ theo neo (bước 29) · scale: giữ tỷ lệ vị trí trong vùng · screen: giữ nguyên chỗ trên màn hình
export type RegionMapMode = 'anchor' | 'scale' | 'screen';

export function regionMapping(to: Rect | null, mode: RegionMapMode) {
  const { steps, region } = state.script;
  const f = mode === 'scale' ? rectToRect(boundsOf(region), boundsOf(to)) : IDENTITY;
  return mapPoints(steps, region, to, (p) => apply(f, p), () => mode !== 'anchor');
}

// Số điểm rơi ra ngoài vùng mới với từng cách – hiện trong hộp chọn.
export function mappingOutside(to: Rect | null, mode: RegionMapMode) {
  if (mode !== 'anchor') return regionMapping(to, mode).outside;
  const b = boundsOf(to);
  let n = 0;
  endpoints(state.script.steps, to).forEach((e) => [e.start, e.end].forEach((p, j) => {
    if (p && (j === 0 || e.end !== e.start) && !insideBounds(p, b)) n++;
  }));
  return n;
}

export function setRegionMapped(to: Rect | null, mode: RegionMapMode) {
  const { steps } = regionMapping(to, mode);
  updateScript({ region: to, steps });
  setNotice(to ? `Đã đổi vùng ${to.w} × ${to.h}` : 'Đã về toàn màn hình');
}

// Kịch bản có điểm nào phụ thuộc vùng không (không có thì đổi vùng khỏi hỏi).
export const hasPoints = () => state.script.steps.some((s) => selectableIndexes(s).length > 0);

// F6 (Go nghe phím toàn cục) gửi vị trí chuột hiện tại.
if (inWails) {
  ScreenBounds().then((screen) => set({ screen })).catch(() => {});
  EventsOn('hotkey:f6', (p: Point) => placePoint(p));
}
