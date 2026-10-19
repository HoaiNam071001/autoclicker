import { describe, expect, it } from 'vitest';
import { makeStep } from '../types';
import { counterText, progressNode } from './progress';

describe('progressNode', () => {
  const steps = [
    makeStep({ id: 'a', ref: 'n2', type: 'waitImage', name: 'Tìm nút Nhắn tin' }),
    makeStep({ id: 'b', ref: 'n5', type: 'counter' }),
    makeStep({ id: 'c', ref: 'n6', type: 'switch', cases: [{ key: 'k1', label: 'Có tin' }] }),
  ];

  it('tới bằng cổng khác Thành công → kèm node trước + nhãn cổng', () => {
    expect(progressNode({ nodeId: 'b', ref: 'n5', from: { ref: 'n2', port: 'fail' } }, steps))
      .toBe('n2 · Tìm nút Nhắn tin · Thất bại → n5 · Đếm');
    expect(progressNode({ nodeId: 'a', ref: 'n2', from: { ref: 'n6', port: 'case:k1' } }, steps))
      .toBe('n6 · Switch theo giá trị · Có tin → n2 · Tìm nút Nhắn tin');
  });

  it('đi thẳng success / node đầu → chỉ node đang chạy', () => {
    expect(progressNode({ nodeId: 'b', ref: 'n5', from: { ref: 'n2', port: 'success' } }, steps)).toBe('n5 · Đếm');
    expect(progressNode({ nodeId: 'a', ref: 'n2' }, steps)).toBe('n2 · Tìm nút Nhắn tin');
    expect(progressNode({}, steps)).toBe('Node');
  });

  it('lượt Đếm', () => {
    expect(counterText({ ref: 'n5', count: 3, max: 10 })).toBe('lượt n5 3 / 10');
  });
});
