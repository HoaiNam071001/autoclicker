export type Point = { x: number; y: number };
export type StepType = 'click' | 'double' | 'hold' | 'drag' | 'wait';
export type MouseButton = 'left' | 'right' | 'center';

export interface Step {
  id: string;
  type: StepType;
  points: Point[]; // drag: points[0] -> points[1]
  button: MouseButton;
  clicks: number;
  clickInterval: number; // ms
  holdMs: number;
  dragMs: number;
  delayAfter: number; // ms
  enabled: boolean;
}

export const STEP_LABEL: Record<StepType, string> = {
  click: 'Click',
  double: 'Double click',
  hold: 'Giữ chuột',
  drag: 'Kéo thả',
  wait: 'Chờ',
};

export const sampleSteps: Step[] = [
  { id: '1', type: 'click', points: [{ x: 412, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, delayAfter: 500, enabled: true },
  { id: '2', type: 'double', points: [{ x: 640, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, delayAfter: 300, enabled: true },
  { id: '3', type: 'click', points: [{ x: 812, y: 436 }], button: 'left', clicks: 5, clickInterval: 120, holdMs: 0, dragMs: 0, delayAfter: 800, enabled: true },
  { id: '4', type: 'drag', points: [{ x: 300, y: 600 }, { x: 900, y: 600 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 400, delayAfter: 1000, enabled: true },
  { id: '5', type: 'wait', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, delayAfter: 2000, enabled: false },
  { id: '6', type: 'click', points: [{ x: 1100, y: 520 }], button: 'right', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, delayAfter: 500, enabled: true },
];

export const fmtMs = (ms: number) => (ms >= 1000 ? `${ms / 1000} s` : `${ms} ms`);
export const fmtPos = (s: Step) =>
  s.type === 'drag'
    ? `A ${s.points[0].x},${s.points[0].y} → B ${s.points[1].x},${s.points[1].y}`
    : s.points[0] ? `X ${s.points[0].x}  Y ${s.points[0].y}` : '—';
