// Biến đổi nhiều điểm cùng lúc (bước 30): ma trận 3×3 (affine: di chuyển / co giãn / xoay / xiên, homography: phối cảnh
// 4 góc), khung bao, căn chỉnh / chia đều / lật. Làm việc trên toạ độ màn hình thật (số thực) – làm tròn khi ghi điểm.
import type { Point, Rect } from './types';

// Ma trận hàng [a b c; d e f; g h i]: x' = (a x + b y + c) / (g x + h y + i), y' = (d x + e y + f) / (g x + h y + i)
export type Mat = readonly [number, number, number, number, number, number, number, number, number];

export const IDENTITY: Mat = [1, 0, 0, 0, 1, 0, 0, 0, 1];

// a ∘ b: áp b trước rồi a.
export function mul(a: Mat, b: Mat): Mat {
  const m = new Array<number>(9);
  for (let r = 0; r < 3; r++) {
    for (let c = 0; c < 3; c++) m[r * 3 + c] = a[r * 3] * b[c] + a[r * 3 + 1] * b[3 + c] + a[r * 3 + 2] * b[6 + c];
  }
  return m as unknown as Mat;
}

// Ghép nhiều phép theo thứ tự áp dụng: chain(t1, t2, t3) = áp t1 rồi t2 rồi t3.
export const chain = (...ms: Mat[]): Mat => ms.reduce((acc, m) => mul(m, acc), IDENTITY);

export function apply(m: Mat, p: Point): Point {
  const w = m[6] * p.x + m[7] * p.y + m[8];
  return { x: (m[0] * p.x + m[1] * p.y + m[2]) / w, y: (m[3] * p.x + m[4] * p.y + m[5]) / w };
}

export const translate = (dx: number, dy: number): Mat => [1, 0, dx, 0, 1, dy, 0, 0, 1];
const about = (m: Mat, p: Point) => chain(translate(-p.x, -p.y), m, translate(p.x, p.y));
export const scaleAbout = (sx: number, sy: number, p: Point) => about([sx, 0, 0, 0, sy, 0, 0, 0, 1], p);
const rad = (deg: number) => (deg * Math.PI) / 180;
export function rotateAbout(deg: number, p: Point) {
  const c = Math.cos(rad(deg)), s = Math.sin(rad(deg));
  return about([c, -s, 0, s, c, 0, 0, 0, 1], p);
}
// Xiên ngang: điểm dời theo x một lượng k × (y − p.y) – đường y = p.y đứng yên. Xiên dọc tương tự.
export const skewXAbout = (k: number, p: Point) => about([1, k, 0, 0, 1, 0, 0, 0, 1], p);
export const skewYAbout = (k: number, p: Point) => about([1, 0, 0, k, 1, 0, 0, 0, 1], p);

// Ô nhập số trên thanh công cụ: ΔX ΔY (px), rộng / cao (%), xoay (°), xiên X / Y (°) – quanh tâm c.
// Thứ tự: co giãn → xiên → xoay → dời.
export type TfParams = { dx: number; dy: number; sx: number; sy: number; rot: number; kx: number; ky: number };
export const NO_TF: TfParams = { dx: 0, dy: 0, sx: 100, sy: 100, rot: 0, kx: 0, ky: 0 };
export const isNoTf = (t: TfParams) => (Object.keys(NO_TF) as (keyof TfParams)[]).every((k) => t[k] === NO_TF[k]);
export const fromParams = (t: TfParams, c: Point): Mat => chain(
  scaleAbout(t.sx / 100, t.sy / 100, c),
  about([1, Math.tan(rad(t.kx)), 0, Math.tan(rad(t.ky)), 1, 0, 0, 0, 1], c),
  rotateAbout(t.rot, c),
  translate(t.dx, t.dy),
);

// Homography biến 4 điểm src thành 4 điểm dst (phối cảnh). null nếu 3 điểm thẳng hàng (không giải được).
export function homography(src: Point[], dst: Point[]): Mat | null {
  // 8 phương trình, 8 ẩn (h33 = 1) – khử Gauss có chọn trụ.
  const A: number[][] = [];
  for (let i = 0; i < 4; i++) {
    const { x, y } = src[i], { x: u, y: v } = dst[i];
    A.push([x, y, 1, 0, 0, 0, -u * x, -u * y, u]);
    A.push([0, 0, 0, x, y, 1, -v * x, -v * y, v]);
  }
  for (let c = 0; c < 8; c++) {
    let p = c;
    for (let r = c + 1; r < 8; r++) if (Math.abs(A[r][c]) > Math.abs(A[p][c])) p = r;
    if (Math.abs(A[p][c]) < 1e-9) return null;
    [A[c], A[p]] = [A[p], A[c]];
    for (let r = 0; r < 8; r++) {
      if (r === c) continue;
      const f = A[r][c] / A[c][c];
      for (let k = c; k < 9; k++) A[r][k] -= f * A[c][k];
    }
  }
  const h = A.map((row, i) => row[8] / row[i]);
  return [h[0], h[1], h[2], h[3], h[4], h[5], h[6], h[7], 1];
}

export const corners = (r: Rect): Point[] => [
  { x: r.x, y: r.y }, { x: r.x + r.w, y: r.y }, { x: r.x + r.w, y: r.y + r.h }, { x: r.x, y: r.y + r.h },
];

// Phối cảnh: kéo góc thứ i (thứ tự corners) của khung r thêm (dx, dy), 3 góc còn lại đứng yên.
export function perspectiveCorner(r: Rect, i: number, dx: number, dy: number): Mat | null {
  const src = corners(r);
  const dst = src.map((p, j) => (j === i ? { x: p.x + dx, y: p.y + dy } : p));
  return homography(src, dst);
}

export function bbox(ps: Point[]): Rect {
  if (!ps.length) return { x: 0, y: 0, w: 0, h: 0 };
  const xs = ps.map((p) => p.x), ys = ps.map((p) => p.y);
  const x = Math.min(...xs), y = Math.min(...ys);
  return { x, y, w: Math.max(...xs) - x, h: Math.max(...ys) - y };
}
export const center = (r: Rect): Point => ({ x: r.x + r.w / 2, y: r.y + r.h / 2 });

// Căn chỉnh / chia đều / lật trong khung bao các điểm – trả mảng điểm mới cùng thứ tự.
export type AlignOp = 'left' | 'hcenter' | 'right' | 'top' | 'vcenter' | 'bottom' | 'hdist' | 'vdist' | 'flipH' | 'flipV';
export function align(ps: Point[], op: AlignOp): Point[] {
  const b = bbox(ps), c = center(b);
  const spread = (axis: 'x' | 'y', start: number, size: number) => {
    const order = ps.map((p, i) => i).sort((i, j) => ps[i][axis] - ps[j][axis]);
    const out = ps.map((p) => ({ ...p }));
    order.forEach((i, k) => { out[i][axis] = start + (size * k) / Math.max(ps.length - 1, 1); });
    return out;
  };
  switch (op) {
    case 'left': return ps.map((p) => ({ ...p, x: b.x }));
    case 'hcenter': return ps.map((p) => ({ ...p, x: c.x }));
    case 'right': return ps.map((p) => ({ ...p, x: b.x + b.w }));
    case 'top': return ps.map((p) => ({ ...p, y: b.y }));
    case 'vcenter': return ps.map((p) => ({ ...p, y: c.y }));
    case 'bottom': return ps.map((p) => ({ ...p, y: b.y + b.h }));
    case 'hdist': return spread('x', b.x, b.w);
    case 'vdist': return spread('y', b.y, b.h);
    case 'flipH': return ps.map((p) => apply(scaleAbout(-1, 1, c), p));
    case 'flipV': return ps.map((p) => apply(scaleAbout(1, -1, c), p));
  }
}

// Vùng cũ → vùng mới, giữ tỷ lệ vị trí trong vùng (ánh xạ "Co giãn theo vùng").
export const rectToRect = (from: Rect, to: Rect): Mat => chain(
  translate(-from.x, -from.y),
  [from.w ? to.w / from.w : 1, 0, 0, 0, from.h ? to.h / from.h : 1, 0, 0, 0, 1],
  translate(to.x, to.y),
);
