import { describe, expect, it } from 'vitest';
import { CATEGORIES, categoryOfType, nodeCategory } from './category';

describe('nodeCategory', () => {
  it('theo loại node', () => {
    expect(nodeCategory({ type: 'click' })).toBe('mouse');
    expect(nodeCategory({ type: 'clickImage' })).toBe('mouse'); // cả theo ảnh / chữ / phần tử
    expect(nodeCategory({ type: 'waitImage' })).toBe('wait'); // cả chờ chữ / phần tử
    expect(nodeCategory({ type: 'screenshot' })).toBe('wait');
    expect(nodeCategory({ type: 'setVar' })).toBe('data');
    expect(nodeCategory({ type: 'elementRead' })).toBe('app');
    expect(nodeCategory({ type: 'switchScreen' })).toBe('flow');
    expect(nodeCategory({ type: 'webRead' })).toBe('web');
  });

  it('mở bằng Chrome của app thuộc nhóm Web', () => {
    expect(nodeCategory({ type: 'open' })).toBe('app');
    expect(nodeCategory({ type: 'open', webBrowser: true })).toBe('web');
  });

  it('theo tên loại dạng chuỗi', () => {
    expect(categoryOfType('http')).toBe('data');
    expect(categoryOfType('khongco')).toBeUndefined();
  });

  it('mỗi nhóm có tên, không trùng', () => {
    expect(new Set(CATEGORIES.map((c) => c.id)).size).toBe(7);
  });
});
