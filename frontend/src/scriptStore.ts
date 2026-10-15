import { useSyncExternalStore } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { makeStep, Point, Rect, Script, sampleSteps, Step } from './types';
import { inWails, openOverlay, setWindowMode } from './windowMode';

// Store kịch bản đang soạn – dùng chung giữa MainScreen và OverlayScreen (hai màn không cùng mount).
// Chưa lưu xuống đĩa (bước 21).
type EditorState = {
  script: Script;
  selected: string; // id bước đang chọn trong danh sách
  notice: string; // thông báo ngắn gần nhất (đã thêm điểm, điểm ngoài vùng…)
  dragA: Point | null; // điểm A (toạ độ màn hình) đã chọn, đang chờ điểm B của kéo thả
};

let state: EditorState = {
  script: { name: 'Kịch bản 1', region: null, loops: 10, maxMinutes: 0, steps: sampleSteps },
  selected: '3',
  notice: '',
  dragA: null,
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

const setSteps = (steps: Step[]) => updateScript({ steps });

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
  const wantsPair = target ? target.type === 'drag' : dragMode;

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
    movePoint(target.id, a ? 1 : 0, screen);
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
  const label = steps[idx].type === 'drag' ? `${idx + 1}${index === 0 ? 'A' : 'B'}` : `${idx + 1}`;
  setSteps(steps.map((s) => (s.id === id ? { ...s, points: s.points.map((p, i) => (i === index ? rel : p)) } : s)));
  set({ selected: id, notice: `Đã đặt điểm ${label}: X ${rel.x} Y ${rel.y}` });
  return true;
}

// F6 (Go nghe phím toàn cục) gửi vị trí chuột hiện tại.
if (inWails) {
  EventsOn('hotkey:f6', (p: Point) => placePoint(p));
}
