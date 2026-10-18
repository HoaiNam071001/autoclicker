import { describe, expect, it, vi } from 'vitest';
import { chain } from '../flow/edit';
import { makeStep } from '../types';

vi.mock('./editor', () => ({ getEditor: () => ({ screen: { x: 0, y: 0, w: 1000, h: 800 } }) }));
const { endpoints, selectableIndexes } = await import('./geometry');

// Điểm dùng biến (F5b): không có chỗ cố định – không marker, không chọn nhiều, node "từ điểm trước" phía sau cũng không biết.
describe('điểm dùng biến', () => {
  const v = { x: 5, y: 5, xVar: '{{n1.output.x}} + 20', yVar: '{{n1.output.y}}' };
  const steps = [
    makeStep({ id: 'a', ref: 'n1', type: 'click', points: [v] }),
    makeStep({ id: 'b', ref: 'n2', type: 'drag', dragFromPrev: true, dragRelative: true, points: [{ x: 0, y: 0 }, { x: 10, y: 0 }] }),
    makeStep({ id: 'c', ref: 'n3', type: 'click', points: [{ x: 7, y: 8 }] }),
  ];
  const ends = endpoints({ steps, start: 'a', edges: chain(['a', 'b', 'c']) }, null);

  it('endpoints null, node sau vẫn biết nguồn', () => {
    expect(ends.get('a')).toEqual({ start: null, end: null, prevStep: null });
    expect(ends.get('b')).toEqual({ start: null, end: null, prevStep: 'a' });
    expect(ends.get('c')?.start).toEqual({ x: 7, y: 8 });
  });

  it('selectableIndexes bỏ điểm dùng biến', () => {
    expect(selectableIndexes(steps[0])).toEqual([]);
    expect(selectableIndexes(makeStep({ id: 'd', type: 'drag', points: [{ x: 1, y: 1 }, v] }))).toEqual([0]);
  });
});
