import { describe, expect, it } from 'vitest';
import { countOK, countRange, findAllSummary } from './findAll';

describe('findAll', () => {
  it('countRange: min 0 = 1, max 0 = không giới hạn', () => {
    expect(countRange({})).toBe('≥ 1');
    expect(countRange({ minCount: 6 })).toBe('≥ 6');
    expect(countRange({ minCount: 2, maxCount: 5 })).toBe('2–5');
    expect(countRange({ maxCount: 1 })).toBe('1');
  });
  it('countOK khớp Go', () => {
    expect(countOK({}, 0)).toBe(false);
    expect(countOK({}, 3)).toBe(true);
    expect(countOK({ minCount: 2, maxCount: 3 }, 4)).toBe(false);
  });
  it('findAllSummary', () => {
    expect(findAllSummary({ count: 6, ms: 420 }, { minCount: 6 })).toBe('Thấy 6 chỗ · 420 ms');
    expect(findAllSummary({ count: 2 }, { minCount: 6 })).toBe('Thấy 2 chỗ – cần ≥ 6');
    expect(findAllSummary({}, {})).toBe('Thấy 0 chỗ – cần ≥ 1');
  });
});
