import { describe, expect, it } from 'vitest';
import { fmtPath, pathVertices, shapePoints } from './shape';
// Cùng dữ liệu với internal/model/shape_test.go
import data from '../../../testdata/shapes.json';

describe('shape', () => {
  it.each(data.cases)('shapePoints($shape)', (c) => {
    const got = shapePoints(c.shape, c.cx, c.cy, c.w, c.h);
    expect(got.length).toBe((c.points ?? []).length);
    got.forEach((p, i) => {
      expect(p.x).toBeCloseTo(c.points![i].x, 2);
      expect(p.y).toBeCloseTo(c.points![i].y, 2);
    });
  });
  it('pathVertices', () => {
    const a = { x: 0, y: 0 }, b = { x: 10, y: 0 };
    expect(pathVertices({ closed: true }, [a, b])).toEqual([a, b, a]);
    expect(pathVertices({}, [a, null])).toBeNull();
    expect(pathVertices({ shape: 'rect', shapeW: 10, shapeH: 4 }, [a])).toHaveLength(5);
  });
  it('fmtPath', () => {
    expect(fmtPath({ points: [{ x: 0, y: 0 }, { x: 1, y: 1 }, { x: 2, y: 2 }], closed: true, dragMs: 1500 })).toBe('3 điểm · khép kín · 1,5 s');
    expect(fmtPath({ shape: 'star', shapeW: 300, points: [], dragMs: 2000 })).toBe('sao 300 px · 2 s');
    expect(fmtPath({ shape: 'rect', shapeW: 200, shapeH: 100, points: [], dragMs: 1000 })).toBe('chữ nhật 200 × 100 px · 1 s');
  });
});
