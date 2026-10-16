// Neo điểm theo cạnh vùng (bước 29) – phải khớp anchor.go.
// Mỗi trục neo riêng: start = cách cạnh trái / trên, end = cách cạnh phải / dưới (lưu số dương, hiện dấu −),
// center = lệch so với giữa (có dấu), ratio = theo tỷ lệ cạnh, đơn vị 1/10000 (5000 = 50%).
import type { Anchor, Point, Rect, StepPoint } from './types';

export const RATIO_SCALE = 10000;

const axisOf = (a: Anchor | undefined): Anchor => a ?? 'start';

function resolveAxis(v: number, a: Anchor | undefined, origin: number, size: number) {
  switch (axisOf(a)) {
    case 'end': return origin + size - v;
    case 'center': return origin + Math.trunc(size / 2) + v; // Go: size/2 chia nguyên
    case 'ratio': return origin + Math.round((size * v) / RATIO_SCALE);
    default: return origin + v;
  }
}

// Ngược với resolveAxis: toạ độ màn hình s → giá trị lưu theo neo a.
function unresolveAxis(s: number, a: Anchor | undefined, origin: number, size: number) {
  switch (axisOf(a)) {
    case 'end': return origin + size - s;
    case 'center': return s - origin - Math.trunc(size / 2);
    case 'ratio': return size > 0 ? Math.round(((s - origin) * RATIO_SCALE) / size) : 0;
    default: return s - origin;
  }
}

// Điểm đã lưu (neo theo khung b) → toạ độ màn hình thật.
export const resolvePoint = (p: StepPoint, b: Rect): Point => ({
  x: resolveAxis(p.x, p.ax, b.x, b.w),
  y: resolveAxis(p.y, p.ay, b.y, b.h),
});

// Toạ độ màn hình → điểm lưu với neo cho trước (giữ nguyên chỗ trên màn hình khi đổi neo).
export const anchorPoint = (s: Point, b: Rect, ax: Anchor | undefined, ay: Anchor | undefined): StepPoint => ({
  x: unresolveAxis(s.x, ax, b.x, b.w),
  y: unresolveAxis(s.y, ay, b.y, b.h),
  ...(ax && ax !== 'start' ? { ax } : {}),
  ...(ay && ay !== 'start' ? { ay } : {}),
});

// Neo theo góc gần nhất của từng trục (mặc định khi đặt điểm mới).
export const nearestAnchor = (s: Point, b: Rect) =>
  anchorPoint(s, b, s.x - b.x > b.w / 2 ? 'end' : 'start', s.y - b.y > b.h / 2 ? 'end' : 'start');

export const insideBounds = (s: Point, b: Rect) => s.x >= b.x && s.x <= b.x + b.w && s.y >= b.y && s.y <= b.y + b.h;

// Hiện một trục: start +300, end −300, center giữa+20, ratio 25%.
const pct = (v: number) => `${Number((v / 100).toFixed(2))}%`;
export function fmtAxis(v: number, a: Anchor | undefined) {
  switch (axisOf(a)) {
    case 'end': return `−${v}`;
    case 'center': return `giữa${v >= 0 ? '+' : '−'}${Math.abs(v)}`;
    case 'ratio': return pct(v);
    default: return `+${v}`;
  }
}

export const fmtPoint = (p: StepPoint) => `X ${fmtAxis(p.x, p.ax)}  Y ${fmtAxis(p.y, p.ay)}`;
