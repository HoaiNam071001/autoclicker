// Điều kiện lồng nhau của node Nếu / case Switch theo màn hình (bước 20 F3) – tìm / sửa / xoá theo đường dẫn. Hàm thuần.
import type { Case, Cond, CondKind, Step } from '../types';

// Chỗ một điều kiện trong node: case (switchScreen, bỏ trống = cond của Nếu) + đường dẫn trong conds lồng nhau.
export type CondAt = { case?: number; path: number[] };

export const emptyCond = (kind: CondKind): Cond =>
  kind === 'image' ? { kind, threshold: 90 }
    : kind === 'pixel' ? { kind, color: '#000000', tolerance: 10 }
    : kind === 'value' ? { kind, left: '', op: 'eq', right: '' }
    : kind === 'all' || kind === 'any' ? { kind, conds: [] }
    : { kind, findText: '' };

export function condAt(c: Cond | undefined, path: number[]): Cond | undefined {
  let cur = c;
  for (const i of path) cur = cur?.conds?.[i];
  return cur;
}

// Thay điều kiện tại path bằng f(điều kiện cũ) – trả điều kiện gốc mới.
export function updateCondAt(c: Cond, path: number[], f: (x: Cond) => Cond): Cond {
  if (!path.length) return f(c);
  const [i, ...rest] = path;
  return { ...c, conds: (c.conds ?? []).map((x, j) => (j === i ? updateCondAt(x, rest, f) : x)) };
}

// Xoá điều kiện con tại path; nhóm và / hoặc chỉ còn 1 điều kiện thì gỡ nhóm (giữ cái còn lại, cộng dồn "đảo").
export function removeCondAt(c: Cond, path: number[]): Cond {
  if (!path.length) return c;
  const i = path[path.length - 1];
  return updateCondAt(c, path.slice(0, -1), (g) => {
    const conds = (g.conds ?? []).filter((_, j) => j !== i);
    if (conds.length === 1) return { ...conds[0], not: !!conds[0].not !== !!g.not || undefined };
    return { ...g, conds };
  });
}

// Điều kiện gốc của node tại at (cond của Nếu / cond của case).
export const rootCond = (s: Pick<Step, 'cond' | 'cases'>, at: Pick<CondAt, 'case'>): Cond | undefined =>
  at.case === undefined ? s.cond : s.cases?.[at.case]?.cond;

export const condOf = (s: Pick<Step, 'cond' | 'cases'>, at: CondAt) => condAt(rootCond(s, at), at.path);

// Phần sửa node khi thay điều kiện gốc tại at bằng f(gốc cũ).
export function patchRoot(s: Pick<Step, 'cond' | 'cases'>, at: Pick<CondAt, 'case'>, f: (root: Cond) => Cond): Partial<Step> {
  const root = rootCond(s, at) ?? emptyCond('text');
  if (at.case === undefined) return { cond: f(root) };
  return { cases: (s.cases ?? []).map((c, i) => (i === at.case ? { ...c, cond: f(root) } : c)) };
}

// Phần sửa node khi thay điều kiện tại at bằng f(điều kiện cũ).
export const patchCond = (s: Pick<Step, 'cond' | 'cases'>, at: CondAt, f: (c: Cond) => Cond): Partial<Step> =>
  patchRoot(s, at, (root) => updateCondAt(root, at.path, f));

// Case mới: key c1, c2… chưa dùng (cổng ra theo key – không đổi key của case đã có để giữ cạnh).
export function newCase(cases: Case[], screen: boolean): Case {
  let n = cases.length + 1;
  while (cases.some((c) => c.key === `c${n}`)) n++;
  return { key: `c${n}`, label: `Case ${n}`, ...(screen ? { cond: emptyCond('text') } : { op: 'eq', match: '' }) };
}
