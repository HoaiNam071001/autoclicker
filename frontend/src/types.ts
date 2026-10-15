import { main } from '../wailsjs/go/models';

export type Point = { x: number; y: number };
export type StepType = 'click' | 'double' | 'hold' | 'drag' | 'wait' | 'type' | 'key';
export type MouseButton = 'left' | 'right' | 'center';
export type Modifier = 'cmd' | 'ctrl' | 'alt' | 'shift';

export interface Step {
  id: string;
  type: StepType;
  points: Point[]; // drag: points[0] -> points[1]
  button: MouseButton;
  clicks: number;
  clickInterval: number; // ms
  holdMs: number;
  dragMs: number;
  text: string; // type: chuỗi cần gõ
  key: string; // key: tên phím robotgo (enter, tab, a, f5…)
  modifiers: Modifier[]; // key: tổ hợp, vd ['cmd'] + 'v'
  delayAfter: number; // ms
  enabled: boolean;
}

// Kiểm tra lúc type-check: Step phải khớp struct Go (model.go) – thiếu/thừa/sai kiểu field là tsc báo lỗi.
type Fields<T> = { [K in keyof T as T[K] extends Function ? never : K]: T[K] };
type Same<A, B> = [A] extends [B] ? ([B] extends [A] ? true : false) : false;
const stepMatchesGo: [Step extends Fields<main.Step> ? true : false, Same<keyof Step, keyof Fields<main.Step>>] = [true, true];
void stepMatchesGo;

// Tạo bước với giá trị trống, ghi đè bằng p (dùng cho dữ liệu mẫu / demo).
export const makeStep = (p: Partial<Step> & Pick<Step, 'id' | 'type'>): Step => ({
  points: [], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0,
  text: '', key: '', modifiers: [], delayAfter: 300, enabled: true, ...p,
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

export const fmtKey = (s: Step) => [...s.modifiers, s.key].map((k) => k[0].toUpperCase() + k.slice(1)).join(' + ');

export const sampleSteps: Step[] = [
  { id: '1', type: 'click', points: [{ x: 412, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, text: '', key: '', modifiers: [], delayAfter: 500, enabled: true },
  { id: '2', type: 'double', points: [{ x: 640, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, text: '', key: '', modifiers: [], delayAfter: 300, enabled: true },
  { id: '3', type: 'click', points: [{ x: 812, y: 436 }], button: 'left', clicks: 5, clickInterval: 120, holdMs: 0, dragMs: 0, text: '', key: '', modifiers: [], delayAfter: 800, enabled: true },
  { id: '4', type: 'drag', points: [{ x: 300, y: 600 }, { x: 900, y: 600 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 400, text: '', key: '', modifiers: [], delayAfter: 1000, enabled: true },
  { id: '5', type: 'wait', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, text: '', key: '', modifiers: [], delayAfter: 2000, enabled: false },
  { id: '6', type: 'click', points: [{ x: 1100, y: 520 }], button: 'right', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, text: '', key: '', modifiers: [], delayAfter: 500, enabled: true },
  { id: '7', type: 'type', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, text: 'xin chao 123', key: '', modifiers: [], delayAfter: 300, enabled: true },
  { id: '8', type: 'key', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, text: '', key: 'enter', modifiers: [], delayAfter: 500, enabled: true },
];

export const fmtMs = (ms: number) => (ms >= 1000 ? `${ms / 1000} s` : `${ms} ms`);
export const fmtPos = (s: Step) =>
  s.type === 'drag'
    ? `A ${s.points[0].x},${s.points[0].y} → B ${s.points[1].x},${s.points[1].y}`
    : s.type === 'type' ? `"${s.text}"`
    : s.type === 'key' ? fmtKey(s)
    : s.points[0] ? `X ${s.points[0].x}  Y ${s.points[0].y}` : '—';
