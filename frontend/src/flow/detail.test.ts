import { describe, expect, it } from 'vitest';
import { matchTags } from './detail';

describe('matchTags', () => {
  it('chỉ ghi cách so khác mặc định', () => {
    expect(matchTags({})).toBe('');
    expect(matchTags({ matchMode: 'edge' })).toBe(' · so viền');
    expect(matchTags({ colorTol: 15 })).toBe(' · màu ±15%');
    expect(matchTags({ matchMode: 'edge', colorTol: 10 })).toBe(' · so viền · màu ±10%');
  });
});
