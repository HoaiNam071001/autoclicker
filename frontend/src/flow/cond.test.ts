import { describe, expect, it } from 'vitest';
import type { Cond } from '../types';
import { condOf, newCase, patchCond, removeCondAt, updateCondAt } from './cond';

const group: Cond = { kind: 'all', conds: [{ kind: 'text', findText: 'a' }, { kind: 'any', not: true, conds: [{ kind: 'value', left: 'x', op: 'eq' }] }] };

describe('điều kiện lồng nhau', () => {
  it('sửa theo đường dẫn, không đụng bản gốc', () => {
    const out = updateCondAt(group, [1, 0], (c) => ({ ...c, right: 'y' }));
    expect(out.conds![1].conds![0].right).toBe('y');
    expect(group.conds![1].conds![0].right).toBeUndefined();
  });

  it('xoá còn 1 điều kiện thì gỡ nhóm, cộng dồn đảo', () => {
    const g: Cond = { kind: 'any', not: true, conds: [{ kind: 'text', findText: 'a' }, { kind: 'text', findText: 'b', not: true }] };
    expect(removeCondAt(g, [0])).toEqual({ kind: 'text', findText: 'b', not: undefined });
    expect(removeCondAt(g, [1])).toEqual({ kind: 'text', findText: 'a', not: true });
  });

  it('điều kiện của case Switch theo màn hình', () => {
    const s = { cases: [{ key: 'c1' }, { key: 'c2', cond: group }] };
    expect(condOf(s, { case: 1, path: [0] })?.findText).toBe('a');
    const p = patchCond(s, { case: 1, path: [0] }, (c) => ({ ...c, findText: 'b' }));
    expect(p.cases![1].cond!.conds![0].findText).toBe('b');
    expect(p.cases![0]).toBe(s.cases[0]);
  });

  it('case mới không trùng key', () => {
    expect(newCase([{ key: 'c1' }, { key: 'c3' }], false)).toMatchObject({ key: 'c4', op: 'eq' });
    expect(newCase([], true).cond?.kind).toBe('text');
  });
});
