// Đường di chuột giống người cho xem trước (bước 34 H4b) – cùng công thức Path / Fitts / Controls / CurveAt
// ở internal/humanize/path.go, test chung testdata/human_path.json. Số ngẫu nhiên mỗi bên tự bốc nên đường
// không trùng từng điểm với lúc chạy thật, chỉ cùng dáng cong, tăng / giảm tốc, thời gian, đi lố.
import type { HumanLevel, Point } from '../types';

// Phần thông số mức ảnh hưởng đường đi – khớp profiles (internal/humanize/profile.go).
export type HumanProfile = { speed: number; curve: number; overshoot: number };
export const HUMAN_PROFILES: Partial<Record<HumanLevel, HumanProfile>> = {
  low: { speed: 0.6, curve: 0.08, overshoot: 0.05 },
  medium: { speed: 1, curve: 0.15, overshoot: 0.15 },
  high: { speed: 1.3, curve: 0.22, overshoot: 0.25 },
};

// Một đoạn đường: khựng pauseMs rồi đi from → to theo Bezier (c1, c2) trong ms.
export type Leg = { from: Point; to: Point; c1: Point; c2: Point; ms: number; pauseMs: number };
// size: cạnh nhỏ của đích (px, 0 = không rõ); durMs > 0: đi đúng chừng này, không đi lố (kéo thả).
export type HumanTarget = { size?: number; durMs?: number };

// Thời gian (ms) di chuột quãng d tới đích cỡ size: a + b·log2(d / size + 1).
export function fitts(d: number, size: number): number {
  const w = size > 0 ? size : 24;
  return 80 + 110 * Math.log2(d / Math.max(w, 8) + 1);
}

// Hai điểm điều khiển ở 30 % / 70 % quãng, lệch theo pháp tuyến o1 / o2 px.
export function controls(a: Point, b: Point, o1: number, o2: number): [Point, Point] {
  const d = Math.max(Math.hypot(b.x - a.x, b.y - a.y), 1e-9);
  const nx = -(b.y - a.y) / d, ny = (b.x - a.x) / d;
  return [
    { x: a.x + (b.x - a.x) * 0.3 + nx * o1, y: a.y + (b.y - a.y) * 0.3 + ny * o1 },
    { x: a.x + (b.x - a.x) * 0.7 + nx * o2, y: a.y + (b.y - a.y) * 0.7 + ny * o2 },
  ];
}

// Điểm trên đường ở thời điểm t ∈ [0, 1] – tiến độ minimum-jerk (chậm – nhanh – chậm).
export function curveAt(a: Point, c1: Point, c2: Point, b: Point, t: number): Point {
  const s = t * t * t * (10 - 15 * t + 6 * t * t);
  const bz = (p0: number, p1: number, p2: number, p3: number) => {
    const u = 1 - s;
    return u * u * u * p0 + 3 * u * u * s * p1 + 3 * u * s * s * p2 + s * s * s * p3;
  };
  return { x: bz(a.x, c1.x, c2.x, b.x), y: bz(a.y, c1.y, c2.y, b.y) };
}

// Các đoạn đưa chuột from → to (như Path ở Go): một đường cong, hoặc đi lố rồi khựng + chỉnh về đích.
export function humanLegs(from: Point, to: Point, p: HumanProfile, t: HumanTarget = {}, rnd: () => number = Math.random): Leg[] {
  const d = Math.hypot(to.x - from.x, to.y - from.y);
  if (d < 2) return [{ from, to, c1: to, c2: to, ms: 0, pauseMs: 0 }];
  const dur = t.durMs && t.durMs > 0 ? t.durMs : fitts(d, t.size ?? 0) * p.speed * (0.85 + 0.3 * rnd());
  if ((t.durMs ?? 0) > 0 || d < 80 || rnd() >= p.overshoot) return [leg(from, to, p.curve, dur, 0, rnd)];
  const over = Math.min(Math.max(d * (0.03 + 0.05 * rnd()), 6), 40);
  const ux = (to.x - from.x) / d, uy = (to.y - from.y) / d;
  const side = normal(rnd) * over * 0.3;
  const o = { x: Math.round(to.x + ux * over - uy * side), y: Math.round(to.y + uy * over + ux * side) };
  return [
    leg(from, o, p.curve, dur * 0.85, 0, rnd),
    leg(o, to, p.curve / 2, (120 + 100 * rnd()) * p.speed, 40 + Math.floor(rnd() * 80), rnd),
  ];
}

function leg(a: Point, b: Point, curve: number, ms: number, pauseMs: number, rnd: () => number): Leg {
  const bend = Math.hypot(b.x - a.x, b.y - a.y) * curve * (rnd() * 2 - 1);
  const [c1, c2] = controls(a, b, bend * (0.6 + 0.8 * rnd()), bend * (0.6 + 0.8 * rnd()));
  return { from: a, to: b, c1, c2, ms, pauseMs };
}

// Phân bố chuẩn (Box–Muller) từ nguồn đều.
function normal(rnd: () => number): number {
  return Math.sqrt(-2 * Math.log(1 - rnd())) * Math.cos(2 * Math.PI * rnd());
}
