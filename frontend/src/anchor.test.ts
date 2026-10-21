import { describe, expect, it } from 'vitest';
import { anchorName, anchorSummary, stepBounds } from './anchor';

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

describe('stepBounds', () => {
  const screen = { x: 0, y: 0, w: 1920, h: 1080 };
  const region = { x: 100, y: 50, w: 800, h: 600 };
  it('mặc định / region dùng khung đang tính', () => {
    expect(stepBounds({}, region, screen)).toEqual(region);
    expect(stepBounds({ coordSpace: 'region' }, region, screen)).toEqual(region);
  });
  it('absolute dùng gốc màn hình', () => {
    expect(stepBounds({ coordSpace: 'absolute' }, region, screen)).toEqual({ x: 0, y: 0, w: 1920, h: 1080 });
  });
});
