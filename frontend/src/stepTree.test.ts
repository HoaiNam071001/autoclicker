import { describe, expect, it } from 'vitest';
import { makeStep, Step } from './types';
import { condAt, removeCondAt, updateCondAt, cloneStep, dropLoc, flatten, insertAt, mapStep, moveTo, removeAt, rows } from './stepTree';

const st = (id: string, p: Partial<Step> = {}) => makeStep({ id, type: 'wait', ...p });
// 1: a | 2: Nếu (2.1 t1, 2.2 t2 | không: 2.3 e1) | 3: Nhóm (3.1 g1) | 4: b
const tree = (): Step[] => [
  st('a'),
  st('if', { type: 'if', children: [st('t1'), st('t2')], else: [st('e1')] }),
  st('g', { type: 'group', children: [st('g1')] }),
  st('b'),
];
const ids = (steps: Step[]) => flatten(steps).map((f) => `${f.step.id}=${f.path}`);

describe('cây bước', () => {
  it('đánh số như engine (nhánh không đánh số tiếp)', () => {
    expect(ids(tree())).toEqual(['a=1', 'if=2', 't1=2.1', 't2=2.2', 'e1=2.3', 'g=3', 'g1=3.1', 'b=4']);
  });
  it('bước trong khối đang tắt không active', () => {
    const t = mapStep(tree(), 'g', (s) => ({ ...s, enabled: false }));
    expect(flatten(t).filter((f) => !f.active).map((f) => f.step.id)).toEqual(['g', 'g1']);
  });
  it('sửa / xoá / chèn bước con', () => {
    const base = tree();
    let t = mapStep(base, 'e1', (s) => ({ ...s, delayAfter: 9 }));
    expect(t[1].else![0].delayAfter).toBe(9);
    expect(t[0]).toBe(base[0]); // nhánh không đổi giữ tham chiếu
    expect(t[2]).toBe(base[2]);
    const r = removeAt(t, 't1')!;
    expect(r.loc).toEqual({ parent: 'if', branch: 'children', index: 0 });
    t = insertAt(r.steps, r.loc, r.step);
    expect(ids(t)).toEqual(ids(tree()));
  });
  it('chuyển bước vào / ra khối, không thả khối vào chính nó', () => {
    let t = moveTo(tree(), 'b', { parent: 'g', branch: 'children', index: 0 });
    expect(ids(t)).toEqual(['a=1', 'if=2', 't1=2.1', 't2=2.2', 'e1=2.3', 'g=3', 'b=3.1', 'g1=3.2']);
    t = moveTo(t, 'g1', { parent: null, branch: 'children', index: 0 });
    expect(ids(t)[0]).toBe('g1=1');
    expect(moveTo(tree(), 'if', { parent: 'if', branch: 'else', index: 0 })).toEqual(tree());
    const nested = moveTo(tree(), 'g', { parent: 'if', branch: 'children', index: 0 });
    expect(moveTo(nested, 'if', { parent: 'g', branch: 'children', index: 0 })).toBe(nested);
  });
  it('nhân bản sâu với id mới', () => {
    let n = 0;
    const c = cloneStep(tree()[1], () => `n${++n}`);
    expect(flatten([c]).map((f) => f.step.id)).toEqual(['n1', 'n2', 'n3', 'n4']);
  });
  it('dòng hiển thị + thả theo dòng phía trên', () => {
    const t = tree();
    const r = rows(t, new Set());
    expect(r.map((x) => x.id)).toEqual(['a', 'if', 't1', 't2', 'if:end', 'if:else', 'e1', 'if:elseEnd', 'g', 'g1', 'g:end', 'b']);
    const at = (id: string) => r.find((x) => x.id === id)!;
    expect(dropLoc(t, 'b', at('if'))).toEqual({ parent: 'if', branch: 'children', index: 0 });
    expect(dropLoc(t, 'b', at('t2'))).toEqual({ parent: 'if', branch: 'children', index: 2 });
    expect(dropLoc(t, 'b', at('if:end'))).toEqual({ parent: 'if', branch: 'children', index: 2 });
    expect(dropLoc(t, 'b', at('if:else'))).toEqual({ parent: 'if', branch: 'else', index: 0 });
    expect(dropLoc(t, 'b', at('if:elseEnd'))).toEqual({ parent: null, branch: 'children', index: 2 });
    expect(dropLoc(t, 'a', at('g:end'))).toEqual({ parent: null, branch: 'children', index: 2 }); // a gỡ ra rồi: sau g
    expect(dropLoc(t, 'b', null)).toEqual({ parent: null, branch: 'children', index: 0 });
    // khối gập: thả dưới nó = ngay sau nó
    const folded = rows(t, new Set(['g']));
    expect(folded.map((x) => x.id)).not.toContain('g1');
    expect(dropLoc(t, 'a', folded.find((x) => x.id === 'g')!)).toEqual({ parent: null, branch: 'children', index: 2 });
  });
  it('bước có "khi thất bại → chạy bước khác" hiện nhánh riêng', () => {
    const t = [st('w', { type: 'waitImage', onFail: { action: 'else' }, else: [st('s')] })];
    expect(ids(t)).toEqual(['w=1', 's=1.1']);
    expect(rows(t, new Set()).map((x) => x.id)).toEqual(['w', 'w:else', 's', 'w:elseEnd']);
  });
});

describe('điều kiện lồng nhau', () => {
  const c = { kind: 'all' as const, conds: [{ kind: 'text' as const, findText: 'a' }, { kind: 'any' as const, conds: [{ kind: 'result' as const }, { kind: 'pixel' as const }] }] };
  it('đọc / sửa theo đường dẫn', () => {
    expect(condAt(c, [1, 1])?.kind).toBe('pixel');
    const n = updateCondAt(c, [1, 0], (x) => ({ ...x, not: true }));
    expect(n.conds![1].conds![0].not).toBe(true);
    expect(c.conds[1].conds![0]).toEqual({ kind: 'result' }); // không sửa bản gốc
  });
  it('xoá – còn 1 thì gỡ nhóm', () => {
    expect(removeCondAt(c, [1, 0]).conds![1]).toEqual({ kind: 'pixel' });
    expect(removeCondAt({ kind: 'all', not: true, conds: [{ kind: 'text' }, { kind: 'result' }] }, [1])).toEqual({ kind: 'text', not: true });
  });
});
