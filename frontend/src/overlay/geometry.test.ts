import { describe, expect, it } from 'vitest';
import { anchorLinesOf, buildMarkers, captureLayout, rectOf, scrollArrowsOf } from './geometry';
import { makeStep, Point } from '../types';

const id = (p: Point) => p;
const key = (s: string, i: number) => `${s}:${i}`;

describe('rectOf', () => {
  it('chuẩn hoá hai góc bất kỳ', () => {
    expect(rectOf({ x: 10, y: 50 }, { x: 4, y: 20 })).toEqual({ x: 4, y: 20, w: 6, h: 30 });
  });
});

describe('buildMarkers', () => {
  const drag = makeStep({ id: 'd', type: 'drag', ref: 'n2', dragFromPrev: true, points: [{ x: 0, y: 0 }, { x: 5, y: 5 }] });
  const click = makeStep({ id: 'c', type: 'click', ref: 'n1', points: [{ x: 1, y: 1 }] });
  const ends = new Map([['c', { start: { x: 1, y: 1 }, end: { x: 1, y: 1 } }], ['d', { start: { x: 1, y: 1 }, end: { x: 9, y: 9 } }]]);

  it('kéo thả có A (ghost khi từ điểm trước, không neo) và B', () => {
    const ms = buildMarkers([click, drag], ends, id, key, null, { x: 0, y: 0 });
    expect(ms.map((m) => [m.label, m.ghost, !!m.pt])).toEqual([['n1', false, true], ['n2A', true, false], ['n2B', false, true]]);
  });

  it('marker đang kéo đi theo con trỏ, toạ độ màn hình giữ nguyên', () => {
    const [m] = buildMarkers([click], ends, id, key, { id: 'c', index: 0 }, { x: 70, y: 80 });
    expect(m.at).toEqual({ x: 70, y: 80 });
    expect(m.scr).toEqual({ x: 1, y: 1 });
  });
});

describe('anchorLinesOf', () => {
  it('đường tới cạnh cuối / giữa khung theo neo', () => {
    const s = makeStep({ id: 'c', type: 'click', ref: 'n1', points: [{ x: 10, y: 0, ax: 'end', ay: 'center' }] });
    const ms = buildMarkers([s], new Map([['c', { start: { x: 50, y: 40 }, end: null }]]), id, key, null, { x: 0, y: 0 });
    const [h, v] = anchorLinesOf(ms, { x: 0, y: 0, w: 100, h: 80 });
    expect([h.x1, h.y1, h.x2]).toEqual([100, 40, 50]);
    expect([v.x1, v.y1, v.y2]).toEqual([50, 40, 40]);
  });
});

describe('scrollArrowsOf', () => {
  it('mũi tên bị cắt khi dài hơn cửa sổ', () => {
    const s = makeStep({ id: 's', type: 'scroll', ref: 'n3', direction: 'down', amount: 1000, points: [{ x: 50, y: 50 }] });
    const [a] = scrollArrowsOf([s], new Map([['s', { start: { x: 50, y: 50 }, end: null }]]), id, { w: 200, h: 200 });
    expect(a.tip).toEqual({ x: 50, y: 188 });
    expect(a.cut).toBe(true);
  });
});

describe('captureLayout', () => {
  const safe = { top: 0, bottom: 800 };
  it('ưu tiên dưới khung', () => {
    expect(captureLayout({ x: 100, y: 100, w: 50, h: 50 }, safe, 1200, true)).toMatchObject({ top: 162, left: 100 });
  });
  it('sát đáy → lên trên khung, nhãn cỡ xuống dưới', () => {
    const l = captureLayout({ x: 100, y: 700, w: 50, h: 80 }, safe, 1200, true);
    expect(l.top).toBe(700 - 12 - 52);
    expect(l.labelTop).toBe(786);
  });
});
