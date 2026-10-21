import { describe, expect, it } from 'vitest';
import { fmtWaitScreen } from './waitScreen';
import { outputsOf } from './view';
import { STEP_DEFAULTS } from './defaults';
import { makeStep, type Step } from '../types';

const node = (p: Partial<Step>) => makeStep({ id: 's9', type: 'waitScreen', ...STEP_DEFAULTS.waitScreen, ...p });
const ref = (id: string) => (id === 's1' ? 'n1' : '?');

describe('waitScreen', () => {
  it('fmtWaitScreen', () => {
    expect(fmtWaitScreen(node({}), ref)).toBe('đứng yên 1 s · vùng kịch bản · tối đa 15 s');
    expect(fmtWaitScreen(node({ waitMode: 'change', diffPct: 2, timeoutMs: 5000, shotArea: 'node', imageFrom: 's1' }), ref))
      .toBe('thay đổi > 2% · chỗ n1 thấy · tối đa 5 s');
    expect(fmtWaitScreen(node({ stableMs: 500, shotArea: 'screen' }), ref)).toBe('đứng yên 500 ms · cả màn · tối đa 15 s');
  });
  it('output ms, diff', () => {
    expect(outputsOf(node({})).map((o) => o.key).slice(0, 2)).toEqual(['ms', 'diff']);
  });
});
