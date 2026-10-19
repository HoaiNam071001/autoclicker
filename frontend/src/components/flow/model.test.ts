import { describe, expect, it } from 'vitest';
import { inputY, nodeHeight, portRows, portY } from './model';

describe('cỡ thẻ node v3', () => {
  const action = ['success', 'fail'];
  const counter = ['next', 'done'];

  it('node Thành công / Thất bại không có dòng cổng', () => {
    expect(portRows(action)).toBe(false);
    expect(portRows(['success'])).toBe(false); // Đặt biến
    expect(portRows(counter)).toBe(true);
    expect(nodeHeight(action)).toBe(nodeHeight([])); // Kết thúc cùng cao
  });

  it('chấm Thành công ngang cổng vào, Thất bại sát đáy', () => {
    expect(portY(action, 0)).toBe(inputY);
    expect(portY(action, 1)).toBe(nodeHeight(action) - 11);
  });

  it('node nhiều lối ra: mỗi cổng một dòng, nằm trong thẻ', () => {
    const h = nodeHeight(counter);
    expect(portY(counter, 1) - portY(counter, 0)).toBe(18);
    expect(portY(counter, 1)).toBeLessThan(h);
    expect(h).toBeGreaterThan(nodeHeight(action));
  });
});
