import { main } from '../wailsjs/go/models';
import { fmtAxis, fmtPoint } from './anchor';

// Neo của một trục (anchor.ts, khớp anchor.go): bỏ trống = start (cách cạnh trái / trên).
export type Anchor = 'start' | 'center' | 'end' | 'ratio';
// Toạ độ màn hình / cửa sổ.
export type Point = { x: number; y: number };
// Điểm lưu trong bước: neo theo cạnh vùng từng trục (khớp Point trong model.go).
export type StepPoint = Point & { ax?: Anchor; ay?: Anchor };
export type Rect = { x: number; y: number; w: number; h: number };
export type StepType = 'click' | 'double' | 'hold' | 'drag' | 'wait' | 'type' | 'key';
export type MouseButton = 'left' | 'right' | 'center';
export type Modifier = 'cmd' | 'ctrl' | 'alt' | 'shift';
// chord: nhấn giữ cùng lúc (Ctrl+C+T) · sequence: giữ phím kèm, gõ lần lượt từng phím chính
export type KeyMode = 'chord' | 'sequence';

export interface Step {
  id: string;
  type: StepType;
  points: StepPoint[]; // drag: points[0] -> points[1]; B lệch so với A (dragRelative) không dùng neo
  button: MouseButton;
  clicks: number;
  clickInterval: number; // ms
  holdMs: number;
  dragMs: number;
  dragFromPrev: boolean; // drag: bắt đầu từ chỗ chuột dừng ở bước chuột trước (bỏ qua points[0])
  dragRelative: boolean; // drag: points[1] là độ lệch (dx, dy) so với điểm bắt đầu
  text: string; // type: chuỗi cần gõ
  keys: string[]; // key: phím chính theo tên robotgo (enter, tab, a, f5…)
  keyMode: KeyMode; // key: chỉ khác nhau khi có ≥ 2 phím chính
  modifiers: Modifier[]; // key: tổ hợp, vd ['cmd'] + 'v'
  delayAfter: number; // ms
  jitterMs: number; // ngẫu nhiên ± ms cho delayAfter và clickInterval
  jitterPx: number; // ngẫu nhiên ± px cho mỗi lần di chuột tới điểm
  enabled: boolean;
}

// Kiểm tra lúc type-check: Step phải khớp struct Go (model.go) – thiếu/thừa/sai kiểu field là tsc báo lỗi.
type Fields<T> = { [K in keyof T as T[K] extends Function ? never : K]: T[K] };
type Same<A, B> = [A] extends [B] ? ([B] extends [A] ? true : false) : false;
const stepMatchesGo: [Step extends Fields<main.Step> ? true : false, Same<keyof Step, keyof Fields<main.Step>>] = [true, true];
void stepMatchesGo;

// Kịch bản – khớp struct Script trong model.go. Toạ độ điểm tương đối theo region (null = toàn màn hình).
export interface Script {
  name: string;
  region: Rect | null; // toạ độ màn hình thật (robotgo)
  loops: number; // 0 = vô hạn
  maxMinutes: number; // 0 = không giới hạn
  steps: Step[];
}
const scriptMatchesGo: Same<keyof Script, keyof Fields<main.Script>> = true;
void scriptMatchesGo;

// Tạo bước với giá trị trống, ghi đè bằng p (dùng cho dữ liệu mẫu / demo).
export const makeStep = (p: Partial<Step> & Pick<Step, 'id' | 'type'>): Step => ({
  points: [], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false,
  text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true, ...p,
});

export const STEP_LABEL: Record<StepType, string> = {
  click: 'Click',
  double: 'Double click',
  hold: 'Giữ chuột',
  drag: 'Kéo thả',
  wait: 'Chờ',
  type: 'Gõ chữ',
  key: 'Phím tắt',
};

// Tên phím robotgo hỗ trợ (ngoài các ký tự đơn như a, 1, /) – gợi ý trong form và kiểm tra.
export const KEY_NAMES = [
  'enter', 'tab', 'esc', 'space', 'backspace', 'delete', 'up', 'down', 'left', 'right',
  'home', 'end', 'pageup', 'pagedown', 'insert', 'capslock', 'printscreen', 'menu',
  ...Array.from({ length: 24 }, (_, i) => `f${i + 1}`),
  'num0', 'num1', 'num2', 'num3', 'num4', 'num5', 'num6', 'num7', 'num8', 'num9', 'num_enter',
  'audio_mute', 'audio_vol_down', 'audio_vol_up', 'audio_play', 'audio_next', 'audio_prev',
];
export const isValidKey = (k: string) => [...k].length === 1 || KEY_NAMES.includes(k);

const cap = (k: string) => k[0].toUpperCase() + k.slice(1);
// chord: "Ctrl + C + T" · sequence: "Ctrl + C → T"
export const fmtKey = (s: Step) => {
  const keys = s.keys.map(cap).join(s.keyMode === 'sequence' ? ' → ' : ' + ');
  return [...s.modifiers.map(cap), keys].filter(Boolean).join(' + ') || '—';
};

export const sampleSteps: Step[] = [
  { id: '1', type: 'click', points: [{ x: 412, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '2', type: 'double', points: [{ x: 640, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '3', type: 'click', points: [{ x: 812, y: 436 }], button: 'left', clicks: 5, clickInterval: 120, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 800, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '4', type: 'drag', points: [{ x: 300, y: 600 }, { x: 900, y: 600 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 400, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 1000, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '5', type: 'wait', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 2000, jitterMs: 0, jitterPx: 0, enabled: false },
  { id: '6', type: 'click', points: [{ x: 1100, y: 520 }], button: 'right', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '7', type: 'type', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: 'xin chao 123', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '8', type: 'key', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: ['enter'], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
];

export const fmtMs = (ms: number) => (ms >= 1000 ? `${ms / 1000} s` : `${ms} ms`);
const fmtSigned = (v: number) => (v > 0 ? `+${v}` : `${v}`);
const fmtPair = (p: StepPoint) => `${fmtAxis(p.x, p.ax)}, ${fmtAxis(p.y, p.ay)}`;
export const fmtPos = (s: Step) =>
  s.type === 'drag'
    ? `${s.dragFromPrev ? 'Từ điểm trước' : `A ${fmtPair(s.points[0])}`} → ${
      s.dragRelative ? `Δ ${fmtSigned(s.points[1].x)},${fmtSigned(s.points[1].y)}` : `B ${fmtPair(s.points[1])}`}`
    : s.type === 'type' ? `"${s.text}"`
    : s.type === 'key' ? fmtKey(s)
    : s.points[0] ? fmtPoint(s.points[0]) : '—';
