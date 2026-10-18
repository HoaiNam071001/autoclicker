import { describe, expect, it } from 'vitest';
import { makeStep, stepLabel, type Step } from '../types';
import { elementSelectorError, firstSelector, fmtElement } from './element';
import { condMissing, condText } from './view';

const wait = (p: Partial<Step>) =>
  makeStep({ id: 's1', type: 'waitImage', findBy: 'element', timeoutMs: 10000, ...p });

describe('phần tử màn hình', () => {
  it('firstSelector', () => {
    expect(firstSelector()).toBe('(chưa nhập)');
    expect(firstSelector([' ', 'role=button'])).toBe('“role=button”');
    expect(firstSelector(['Mua', 'aria=Mua'])).toBe('“Mua” +1');
  });
  it('fmtElement + tên node', () => {
    expect(fmtElement(wait({ selectors: [] }))).toBe('chưa nhập cách tìm phần tử');
    expect(fmtElement(wait({ selectors: ['Mua ngay'], occurrence: 2, elementApp: 'Google Chrome' })))
      .toBe('phần tử “Mua ngay” (chỗ thứ 2) · Google Chrome · tối đa 10 s');
    expect(stepLabel(wait({}))).toBe('Chờ phần tử');
    expect(stepLabel(wait({ type: 'clickImage', imageSource: 'search' }))).toBe('Click theo phần tử');
  });
  it('điều kiện thấy phần tử', () => {
    expect(condMissing({ kind: 'element', selectors: [' '] })).toBe('Chưa nhập cách tìm phần tử');
    expect(condText({ kind: 'element', selectors: ['role=link'], not: true })).toBe('không thấy phần tử “role=link”');
  });
  it('elementSelectorError', () => {
    expect(elementSelectorError(['text=Mua', 'role=button[name="Gửi"]', '{{cach}}'])).toBe('');
    expect(elementSelectorError(['Mua', 'css=.buy'])).toContain('“css=.buy”');
  });
});
