import { describe, expect, it } from 'vitest';
import { effectiveHuman, scriptHuman, usesHuman } from './human';

describe('thao tác giống người (bước 34)', () => {
  it('mức: kịch bản trống = tắt, node trống = theo kịch bản', () => {
    expect(scriptHuman(undefined)).toBe('off');
    expect(effectiveHuman('medium', undefined)).toBe('medium');
    expect(effectiveHuman('medium', 'off')).toBe('off');
    expect(effectiveHuman(undefined, 'high')).toBe('high');
  });
  it('chỉ node dùng chuột / phím thật', () => {
    expect(usesHuman({ type: 'click' })).toBe(true);
    expect(usesHuman({ type: 'wait' })).toBe(false);
    expect(usesHuman({ type: 'webClick' })).toBe(true);
    expect(usesHuman({ type: 'webClick', webBackground: true })).toBe(false);
    expect(usesHuman({ type: 'webRead' })).toBe(false);
    expect(usesHuman({ type: 'webScroll' })).toBe(false); // cuộn tới phần tử bằng trang, không lăn chuột
    expect(usesHuman({ type: 'webScroll', scrollMode: 'in' })).toBe(true);
  });
});
