import { describe, expect, it } from 'vitest';
import { aiRunScope } from './aiScope';

describe('aiRunScope', () => {
  const base = { id: '1', script: 'X', path: '/x.json', steps: 5 };
  it('cả kịch bản', () => expect(aiRunScope(base)).toBe('cả kịch bản (5 node)'));
  it('từ / đến / chạy tiếp', () => expect(aiRunScope({ ...base, from: 'n3', stopAt: 'n5', resume: true })).toBe('từ n3 đến n5, dùng kết quả lần chạy trước'));
  it('cụm', () => expect(aiRunScope({ ...base, nodes: ['n2', 'n3'] })).toBe('chỉ n2, n3'));
});
