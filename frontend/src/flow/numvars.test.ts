import { describe, expect, it } from 'vitest';
import { makeStep } from '../types';
import { toggleNumVar, withNumVar } from './numvars';

describe('numVars (F5c)', () => {
  const s = makeStep({ id: 'a', type: 'click', clicks: 2 });
  it('bật điền số hiện tại, sửa giữ vị trí, bỏ dòng cuối thì numVars undefined', () => {
    const on = toggleNumVar(s, 'clicks', true);
    expect(on.numVars).toEqual([{ field: 'clicks', expr: '2' }]);
    const two = withNumVar(withNumVar(on, 'delayAfter', '{{n1.output.ms}}'), 'clicks', '{{vars.n}}');
    expect(two.numVars?.map((v) => v.field)).toEqual(['clicks', 'delayAfter']);
    expect(withNumVar(on, 'clicks', undefined).numVars).toBeUndefined();
  });
  it('tắt: số nguyên giữ lại làm giá trị cố định, biến thì giữ số cũ', () => {
    expect(toggleNumVar(withNumVar(s, 'clicks', ' 5 '), 'clicks', false)).toMatchObject({ clicks: 5, numVars: undefined });
    expect(toggleNumVar(withNumVar(s, 'clicks', '{{vars.n}}'), 'clicks', false).clicks).toBe(2);
  });
});
