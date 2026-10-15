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
};

let state: EditorState = {
  script: { name: 'Kịch bản 1', region: null, loops: 10, maxMinutes: 0, steps: sampleSteps },
  selected: '3',
  notice: '',
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

export function pickPointFor(stepId: string) {
  pickTarget = stepId;
  return openOverlay('point');
}

export const pickingStep = () => pickTarget;

export const cancelPick = () => {
  pickTarget = null;
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
  if (pickTarget) {
    const id = pickTarget;
    pickTarget = null;
    const idx = steps.findIndex((s) => s.id === id);
    if (idx >= 0) {
      setSteps(steps.map((s) => (s.id === id ? { ...s, points: [rel, ...s.points.slice(1)] } : s)));
      set({ selected: id, notice: `Đã đặt lại điểm bước ${idx + 1}: X ${rel.x} Y ${rel.y}` });
    }
    setWindowMode('main');
    return true;
  }
  const step = makeStep({ id: newId(), type: 'click', points: [rel], delayAfter: 500 });
  setSteps([...steps, step]);
  set({ selected: step.id, notice: `Đã thêm bước ${steps.length + 1}: click X ${rel.x} Y ${rel.y}` });
  return true;
}

// F6 (Go nghe phím toàn cục) gửi vị trí chuột hiện tại.
if (inWails) {
  EventsOn('hotkey:f6', (p: Point) => placePoint(p));
}
