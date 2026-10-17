import { describe, expect, it } from 'vitest';
import { makeStep } from '../types';
import { chain } from './edit';
import { nodeViews } from './view';

describe('nodeViews', () => {
  const steps = [
    makeStep({ id: 'a', ref: 'n1', type: 'waitImage', findBy: 'text', findText: 'Gửi', timeoutMs: 3000, retry: { count: 2 } }),
    makeStep({ id: 'b', ref: 'n2', type: 'clickImage', imageSource: 'node', imageFrom: 'a' }),
    makeStep({ id: 'c', ref: 'n3', type: 'clickImage', imageSource: 'node' }),
  ];
  const sc = { steps, start: 'a', edges: chain(['a', 'b']) };

  it('cổng có nhãn + màu, chip thử lại / input, nguồn chỗ thấy theo mã node', () => {
    const [a, b] = nodeViews(sc, { inputs: (id) => (id === 'a' ? 1 : 0) });
    expect(a.ports).toEqual([{ id: 'success', label: 'Thành công', kind: 'ok' }, { id: 'fail', label: 'Thất bại', kind: 'fail' }]);
    expect(a.chips).toEqual(['1 input', 'thử lại 2']);
    expect(a.title).toBe('Chờ chữ');
    expect(b.detail).toContain('chỗ n1 thấy');
  });

  it('cảnh báo node không đi tới được', () => {
    const c = nodeViews(sc)[2];
    expect(c.warn).toContain('Không có đường nào');
    expect(nodeViews({ ...sc, edges: chain(['a', 'c']) })[2].warn).toContain('Chưa chọn node');
  });
});
