import { describe, expect, it } from 'vitest';
import {
  align, apply, bbox, chain, fromParams, homography, IDENTITY, NO_TF, perspectiveCorner, rectToRect, rotateAbout, scaleAbout,
  skewXAbout, skewYAbout, translate,
} from './transform';
import type { Point } from './types';

const close = (p: Point, q: Point) => {
  expect(p.x).toBeCloseTo(q.x, 6);
  expect(p.y).toBeCloseTo(q.y, 6);
};

describe('affine', () => {
  it('đơn vị và dời', () => {
    close(apply(IDENTITY, { x: 3, y: 4 }), { x: 3, y: 4 });
    close(apply(translate(10, -5), { x: 3, y: 4 }), { x: 13, y: -1 });
  });
  it('co giãn quanh điểm giữ điểm đó đứng yên', () => {
    const c = { x: 100, y: 100 };
    close(apply(scaleAbout(1.5, 2, c), c), c);
    close(apply(scaleAbout(1.5, 2, c), { x: 200, y: 150 }), { x: 250, y: 200 });
  });
  it('xoay 90° theo chiều kim đồng hồ trên màn hình (y hướng xuống)', () => {
    close(apply(rotateAbout(90, { x: 0, y: 0 }), { x: 10, y: 0 }), { x: 0, y: 10 });
    close(apply(rotateAbout(90, { x: 5, y: 5 }), { x: 15, y: 5 }), { x: 5, y: 15 });
  });
  it('xiên giữ đường trục đứng yên', () => {
    close(apply(skewXAbout(0.5, { x: 0, y: 100 }), { x: 10, y: 100 }), { x: 10, y: 100 });
    close(apply(skewXAbout(0.5, { x: 0, y: 100 }), { x: 10, y: 140 }), { x: 30, y: 140 });
    close(apply(skewYAbout(-1, { x: 50, y: 0 }), { x: 60, y: 7 }), { x: 60, y: -3 });
  });
  it('chain áp theo thứ tự', () => {
    // dời rồi co giãn quanh gốc ≠ co giãn rồi dời
    close(apply(chain(translate(10, 0), scaleAbout(2, 2, { x: 0, y: 0 })), { x: 0, y: 0 }), { x: 20, y: 0 });
    close(apply(chain(scaleAbout(2, 2, { x: 0, y: 0 }), translate(10, 0)), { x: 0, y: 0 }), { x: 10, y: 0 });
  });
  it('fromParams: mặc định không đổi gì; 150% + dời quanh tâm', () => {
    const c = { x: 100, y: 100 };
    close(apply(fromParams(NO_TF, c), { x: 37, y: 81 }), { x: 37, y: 81 });
    close(apply(fromParams({ ...NO_TF, sx: 150, sy: 150, dx: 5, dy: -5 }, c), { x: 200, y: 100 }), { x: 255, y: 95 });
  });
});

describe('homography', () => {
  it('4 điểm khớp chính xác, kể cả phối cảnh', () => {
    const src = [{ x: 0, y: 0 }, { x: 100, y: 0 }, { x: 100, y: 100 }, { x: 0, y: 100 }];
    const dst = [{ x: 10, y: 5 }, { x: 120, y: 0 }, { x: 90, y: 130 }, { x: -5, y: 95 }];
    const h = homography(src, dst)!;
    src.forEach((p, i) => close(apply(h, p), dst[i]));
  });
  it('3 điểm thẳng hàng → null', () => {
    const src = [{ x: 0, y: 0 }, { x: 50, y: 0 }, { x: 100, y: 0 }, { x: 0, y: 100 }];
    expect(homography(src, src)).toBeNull();
  });
  it('kéo một góc: 3 góc còn lại đứng yên', () => {
    const r = { x: 0, y: 0, w: 200, h: 100 };
    const h = perspectiveCorner(r, 2, 30, 20)!;
    close(apply(h, { x: 0, y: 0 }), { x: 0, y: 0 });
    close(apply(h, { x: 200, y: 0 }), { x: 200, y: 0 });
    close(apply(h, { x: 200, y: 100 }), { x: 230, y: 120 });
    close(apply(h, { x: 0, y: 100 }), { x: 0, y: 100 });
  });
});

describe('căn chỉnh', () => {
  const ps = [{ x: 10, y: 50 }, { x: 40, y: 10 }, { x: 100, y: 30 }];
  it('khung bao', () => expect(bbox(ps)).toEqual({ x: 10, y: 10, w: 90, h: 40 }));
  it('căn trái / dưới / giữa', () => {
    expect(align(ps, 'left').map((p) => p.x)).toEqual([10, 10, 10]);
    expect(align(ps, 'bottom').map((p) => p.y)).toEqual([50, 50, 50]);
    expect(align(ps, 'hcenter').map((p) => p.x)).toEqual([55, 55, 55]);
  });
  it('chia đều ngang theo thứ tự x, giữ thứ tự mảng', () => {
    const out = align([{ x: 100, y: 0 }, { x: 0, y: 0 }, { x: 20, y: 0 }], 'hdist');
    expect(out.map((p) => p.x)).toEqual([100, 0, 50]);
  });
  it('lật ngang quanh tâm khung', () => {
    expect(align(ps, 'flipH').map((p) => Math.round(p.x))).toEqual([100, 70, 10]);
  });
});

it('vùng cũ → vùng mới giữ tỷ lệ vị trí', () => {
  const m = rectToRect({ x: 0, y: 0, w: 1000, h: 700 }, { x: 100, y: 50, w: 1400, h: 900 });
  close(apply(m, { x: 500, y: 350 }), { x: 800, y: 500 });
  close(apply(m, { x: 1000, y: 700 }), { x: 1500, y: 950 });
});
