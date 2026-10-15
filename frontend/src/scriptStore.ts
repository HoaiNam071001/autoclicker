import { useSyncExternalStore } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { makeStep, Point, Rect, Script, sampleSteps, Step, StepType } from './types';
import { inWails, openOverlay, setWindowMode } from './windowMode';
import { runScript } from './engineStore';

// Store kịch bản đang soạn – dùng chung giữa MainScreen và OverlayScreen (hai màn không cùng mount).
// Chưa lưu xuống đĩa (bước 21).
type EditorState = {
  script: Script;
  selected: string; // id bước đang chọn trong danh sách
  notice: string; // thông báo ngắn gần nhất (đã thêm điểm, điểm ngoài vùng…)
  dragA: Point | null; // điểm A (toạ độ màn hình) đã chọn, đang chờ điểm B của kéo thả
  deleted: { step: Step; index: number } | null; // bước vừa xoá – để hoàn tác
  countdown: boolean; // đếm ngược 3 giây trước khi chạy
};

let state: EditorState = {
  script: { name: 'Kịch bản 1', region: null, loops: 10, maxMinutes: 0, steps: sampleSteps },
  selected: '3',
  notice: '',
  dragA: null,
  deleted: null,
  countdown: true,
};
const listeners = new Set<() => void>();

const set = (patch: Partial<EditorState>) => {
  state = { ...state, ...patch };
  listeners.forEach((l) => l());
};

const subscribe = (l: () => void) => {
  listeners.add(l);
  return () => listeners.delete(l);
};

export const useEditor = () => useSyncExternalStore(subscribe, () => state);
export const useScript = () => useEditor().script;

export const updateScript = (patch: Partial<Script>) => set({ script: { ...state.script, ...patch } });
export const selectStep = (id: string) => set({ selected: id });
export const setNotice = (notice: string) => set({ notice });
export const setCountdown = (countdown: boolean) => set({ countdown });

// Chạy kịch bản đang soạn (nút Chạy, F8 khi chưa chạy).
export const runCurrent = () => runScript(state.script, state.countdown);

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
  if (need === 2 && points.length === 1) points.push({ x: points[0].x + 100, y: points[0].y });
  updateStep(id, {
    type,
    points,
    clicks: Math.max(s.clicks, 1),
    holdMs: type === 'hold' && !s.holdMs ? 1000 : s.holdMs,
    dragMs: type === 'drag' && !s.dragMs ? 400 : s.dragMs,
    keys: type === 'key' && s.keys.length === 0 ? ['enter'] : s.keys,
  });
}

// Thêm bước không cần chọn điểm (chờ, gõ chữ, phím tắt) vào cuối – giá trị mặc định giống NewStep (model.go).
export function addStep(type: Extract<StepType, 'wait' | 'type' | 'key'>) {
  const defaults: Record<typeof type, Partial<Step>> = {
    wait: { clicks: 0, delayAfter: 1000 },
    type: { clicks: 0, delayAfter: 500 },
    key: { clicks: 0, keys: ['enter'], delayAfter: 500 },
  };
  const step = makeStep({ id: newId(), type, ...defaults[type] });
  const steps = state.script.steps;
  setSteps([...steps, step]);
  set({ selected: step.id, notice: `Đã thêm bước ${steps.length + 1}` });
}

// Đổi vùng; toạ độ điểm vẫn tương đối nên các điểm đi theo vùng mới. null = toàn màn hình.
export const setRegion = (region: Rect | null) => updateScript({ region });

// Điểm tương đối theo vùng → toạ độ màn hình thật.
export const toScreen = (p: Point, region: Rect | null): Point =>
  region ? { x: region.x + p.x, y: region.y + p.y } : p;

// Toạ độ màn hình → tương đối theo vùng; null nếu nằm ngoài vùng.
const toRelative = (p: Point, r: Rect | null): Point | null => {
  if (!r) return p;
  const rel = { x: p.x - r.x, y: p.y - r.y };
  return rel.x < 0 || rel.y < 0 || rel.x > r.w || rel.y > r.h ? null : rel;
};

const newId = () => `${Date.now().toString(36)}${Math.random().toString(36).slice(2, 6)}`;

// Bước đang được chọn lại điểm (nút ⌖ trong form); null = điểm mới sẽ thêm thành bước mới.
let pickTarget: string | null = null;
// Overlay đang ở công cụ "Kéo thả A → B": điểm mới đi theo cặp A, B thay vì thành bước click.
let dragMode = false;

export function pickPointFor(stepId: string) {
  pickTarget = stepId;
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
  setDragMode(false);
};

// Đặt một điểm (toạ độ màn hình thật) từ click trên overlay hoặc F6:
// đang chọn lại → sửa điểm của bước đó rồi về màn chính; không thì thêm bước click mới.
export function placePoint(screen: Point): boolean {
  const { script } = state;
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
    set({ dragA: screen, notice: `Điểm A: X ${rel.x} Y ${rel.y} – chọn điểm B (Esc để huỷ)` });
    return true;
  }
  const a = wantsPair ? state.dragA! : null;
  if (a) set({ dragA: null });

  if (target) {
    pickTarget = null;
    if (a) movePoint(target.id, 0, a);
    movePoint(target.id, target.type === 'drag' ? 1 : 0, screen);
    setWindowMode('main');
    return true;
  }
  if (a) {
    const relA = toRelative(a, script.region)!;
    const step = makeStep({ id: newId(), type: 'drag', points: [relA, rel], dragMs: 400, delayAfter: 500 });
    setSteps([...steps, step]);
    set({ selected: step.id, notice: `Đã thêm bước ${steps.length + 1}: kéo A ${relA.x},${relA.y} → B ${rel.x},${rel.y}` });
    return true;
  }
  const step = makeStep({ id: newId(), type: 'click', points: [rel], delayAfter: 500 });
  setSteps([...steps, step]);
  set({ selected: step.id, notice: `Đã thêm bước ${steps.length + 1}: click X ${rel.x} Y ${rel.y}` });
  return true;
}

// Dời điểm thứ index của bước id tới toạ độ màn hình screen (kéo marker, chọn lại điểm).
// Ngoài vùng → báo và giữ nguyên.
export function movePoint(id: string, index: number, screen: Point): boolean {
  const { steps, region } = state.script;
  const idx = steps.findIndex((s) => s.id === id);
  const rel = toRelative(screen, region);
  if (idx < 0) return false;
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
  set({ selected: id, notice: `Đã đặt điểm ${label}: X ${rel.x} Y ${rel.y}` });
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
    patch.points = [toRelative(ends.start, region) ?? s.points[0], s.points[1]];
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
      : region ? { x: ends.end.x - region.x, y: ends.end.y - region.y } : ends.end;
  }
  updateStep(id, { dragRelative: on, points: [s.points[0], b] });
}

// F6 (Go nghe phím toàn cục) gửi vị trí chuột hiện tại.
if (inWails) {
  EventsOn('hotkey:f6', (p: Point) => placePoint(p));
  EventsOn('hotkey:f8start', () => runCurrent());
}
