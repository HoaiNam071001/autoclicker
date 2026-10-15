import { Point } from './types';

// Ánh xạ toạ độ trong cửa sổ overlay (CSS px) ↔ toạ độ robotgo (toạ độ dùng để click), mỗi trục:
//   màn hình = k · cửa sổ + b
// k ≠ 1 khi màn hình scale khác 100% trên Windows (robotgo tính theo DPI), b là vị trí cửa sổ
// (macOS: cửa sổ nằm dưới menu bar; màn phụ: cộng gốc màn đó). Ước lượng từ các mẫu
// (vị trí trong cửa sổ, MousePos) lấy khi chuột đứng yên – không cần biết hệ điều hành / scale.
export type Mapping = { kx: number; bx: number; ky: number; by: number };

const MIN_SPREAD = 100; // px – các mẫu phải cách nhau đủ xa mới ước lượng k, không thì giữ k cũ
const MAX_SAMPLES = 20;

let samples: { c: Point; s: Point }[] = [];
let mapping: Mapping = { kx: 1, bx: 0, ky: 1, by: 0 };
let measured = false;

export const currentMapping = () => mapping;
export const isMeasured = () => measured;

// Gọi mỗi lần mở overlay: có thể đang ở màn khác nên bỏ mẫu cũ, nhưng giữ k (thường giống nhau)
// để mẫu đầu tiên đã cho kết quả đúng.
export const resetSamples = () => {
  samples = [];
};

function fitAxis(cs: number[], ss: number[], k0: number) {
  const n = cs.length;
  const mc = cs.reduce((a, v) => a + v, 0) / n;
  const ms = ss.reduce((a, v) => a + v, 0) / n;
  let k = k0;
  if (n >= 2 && Math.max(...cs) - Math.min(...cs) >= MIN_SPREAD) {
    let cov = 0;
    let varC = 0;
    cs.forEach((c, i) => {
      cov += (c - mc) * (ss[i] - ms);
      varC += (c - mc) ** 2;
    });
    k = cov / varC;
  }
  return { k, b: ms - k * mc };
}

export function addSample(c: Point, s: Point): Mapping {
  samples = [...samples.slice(-(MAX_SAMPLES - 1)), { c, s }];
  const x = fitAxis(samples.map((p) => p.c.x), samples.map((p) => p.s.x), mapping.kx);
  const y = fitAxis(samples.map((p) => p.c.y), samples.map((p) => p.s.y), mapping.ky);
  mapping = { kx: x.k, bx: x.b, ky: y.k, by: y.b };
  measured = true;
  return mapping;
}

export const toScreenPt = (m: Mapping, p: Point): Point => ({
  x: Math.round(m.kx * p.x + m.bx),
  y: Math.round(m.ky * p.y + m.by),
});

export const toClientPt = (m: Mapping, p: Point): Point => ({
  x: (p.x - m.bx) / m.kx,
  y: (p.y - m.by) / m.ky,
});
