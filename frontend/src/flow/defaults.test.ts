import { describe, expect, it } from 'vitest';
import { defaultOf } from './defaults';

describe('defaultOf', () => {
  it('theo loại node (khớp NewStep)', () => {
    expect(defaultOf('wait', 'delayAfter')).toBe(1000);
    expect(defaultOf('open', 'delayAfter')).toBe(2000);
    expect(defaultOf('click', 'clicks')).toBe(1);
    expect(defaultOf('waitImage', 'timeoutMs')).toBe(10000);
    expect(defaultOf('scroll', 'amount')).toBe(3);
  });

  it('ô loại node không có thì không có mặc định', () => {
    expect(defaultOf('click', 'timeoutMs')).toBeUndefined();
  });
});
