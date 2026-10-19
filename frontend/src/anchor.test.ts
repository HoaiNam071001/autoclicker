import { describe, expect, it } from 'vitest';
import { anchorName, anchorSummary } from './anchor';

describe('chữ neo (khối Vị trí)', () => {
  it('tên neo', () => {
    expect(anchorName({})).toBe('neo trên-trái');
    expect(anchorName({ ax: 'end', ay: 'end' })).toBe('neo dưới-phải');
    expect(anchorName({ ax: 'center', ay: 'center' })).toBe('neo tâm');
    expect(anchorName({ ax: 'ratio', ay: 'ratio' })).toBe('theo %');
  });

  it('tóm tắt khoảng cách', () => {
    expect(anchorSummary({ x: 812, y: 436 })).toBe('cách trái 812, trên 436');
    expect(anchorSummary({ x: 20, y: 10, ax: 'center', ay: 'end' })).toBe('giữa ngang +20, cách dưới 10');
    expect(anchorSummary({ x: 5000, y: 2500, ax: 'ratio', ay: 'ratio' })).toBe('50% · 25%');
  });
});
