// Thao tác trên cây bước (bước 20b): khối Nếu / Nhóm có bước con `children`, nhánh `else` (Nếu không, hoặc
// "Khi thất bại → chạy bước khác"). Hàm thuần, không đụng store – scriptStore gọi rồi ghi kết quả (hoàn tác như cũ).
// Đánh số giống Go (blocks.go walkSteps): bước con "3.1", nhánh else đánh số tiếp sau nhánh đúng.
import { Cond, Step } from './types';

export type Branch = 'children' | 'else';
// Vị trí trong cây: danh sách nào (gốc khi parent null), chỗ thứ mấy.
export type Loc = { parent: string | null; branch: Branch; index: number };

export type FlatStep = {
  step: Step;
  path: string; // "3.2"
  depth: number; // 0 = gốc
  loc: Loc;
  active: boolean; // bước và mọi khối chứa nó đều đang bật (engine sẽ chạy tới)
};

export const isBlock = (s: Step) => s.type === 'if' || s.type === 'group';
// Bước có nhánh else hiện trong danh sách: khối Nếu, hoặc bước có "Khi thất bại → chạy bước khác"
export const hasElse = (s: Step) => s.type === 'if' || s.onFail?.action === 'else';

const listOf = (s: Step, b: Branch) => (b === 'children' ? s.children : s.else) ?? [];

// Mọi bước theo thứ tự duyệt (bước cha trước bước con) – cũng là thứ tự engine đi qua khi nhánh nào cũng chạy.
export function flatten(steps: Step[], prefix = '', depth = 0, parent: string | null = null, branch: Branch = 'children', off = 0, active = true): FlatStep[] {
  const out: FlatStep[] = [];
  steps.forEach((s, i) => {
    const path = `${prefix}${off + i + 1}`;
    const on = active && s.enabled;
    out.push({ step: s, path, depth, loc: { parent, branch, index: i }, active: on });
    if (isBlock(s)) out.push(...flatten(listOf(s, 'children'), `${path}.`, depth + 1, s.id, 'children', 0, on));
    if (hasElse(s)) out.push(...flatten(listOf(s, 'else'), `${path}.`, depth + 1, s.id, 'else', s.type === 'if' ? listOf(s, 'children').length : 0, on));
  });
  return out;
}

// Thứ tự Xem trước (con trỏ giả): khối Nếu đi nhánh đúng, khối lặp chạy 1 vòng, bỏ nhánh "Khi thất bại" (20d: chọn nhánh).
export function previewOrder(steps: Step[]): FlatStep[] {
  const byId = new Map(flatten(steps).map((f) => [f.step.id, f]));
  const out: FlatStep[] = [];
  const walk = (list: Step[]) => list.forEach((s) => {
    const f = byId.get(s.id)!;
    if (!f.active) return;
    out.push(f);
    if (isBlock(s)) walk(s.children ?? []);
  });
  walk(steps);
  return out;
}

export const allSteps = (steps: Step[]) => flatten(steps).map((f) => f.step);
export const flatOf = (steps: Step[], id: string) => flatten(steps).find((f) => f.step.id === id);

// Sửa một bước (ở bất kỳ tầng nào) – trả cây mới, nhánh không đổi giữ nguyên tham chiếu.
export function mapStep(steps: Step[], id: string, f: (s: Step) => Step): Step[] {
  let changed = false;
  const out = steps.map((s) => {
    if (s.id === id) {
      changed = true;
      return f(s);
    }
    const children = s.children && mapStep(s.children, id, f);
    const els = s.else && mapStep(s.else, id, f);
    if (children === s.children && els === s.else) return s;
    changed = true;
    return { ...s, children, else: els };
  });
  return changed ? out : steps;
}

// Áp f lên mọi bước (kể cả bước con).
export function mapAll(steps: Step[], f: (s: Step) => Step): Step[] {
  return steps.map((s) => {
    const n = f(s);
    return { ...n, ...(n.children && { children: mapAll(n.children, f) }), ...(n.else && { else: mapAll(n.else, f) }) };
  });
}

// Sửa danh sách tại parent / branch (null = gốc).
function mapList(steps: Step[], parent: string | null, branch: Branch, f: (list: Step[]) => Step[]): Step[] {
  if (parent === null) return f(steps);
  return mapStep(steps, parent, (s) => ({ ...s, [branch]: f(listOf(s, branch)) }));
}

export function removeAt(steps: Step[], id: string): { steps: Step[]; step: Step; loc: Loc } | null {
  const f = flatOf(steps, id);
  if (!f) return null;
  const { parent, branch } = f.loc;
  return { steps: mapList(steps, parent, branch, (l) => l.filter((s) => s.id !== id)), step: f.step, loc: f.loc };
}

export function insertAt(steps: Step[], loc: Loc, step: Step): Step[] {
  return mapList(steps, loc.parent, loc.branch, (l) => {
    const i = Math.max(0, Math.min(loc.index, l.length));
    return [...l.slice(0, i), step, ...l.slice(i)];
  });
}

// Chuyển bước id tới loc (index tính trong danh sách đích SAU khi đã gỡ bước ra). Không cho thả khối vào chính nó.
export function moveTo(steps: Step[], id: string, loc: Loc): Step[] {
  const me = flatOf(steps, id);
  if (!me) return steps;
  if (loc.parent && (loc.parent === id || flatOf([me.step], loc.parent))) return steps; // thả vào chính nó / bước con của nó
  const r = removeAt(steps, id);
  return r ? insertAt(r.steps, loc, r.step) : steps;
}

// Sao chép sâu, đặt id mới cho bước và mọi bước con.
export function cloneStep(s: Step, newId: () => string): Step {
  return {
    ...s, id: newId(), points: s.points.map((p) => ({ ...p })), keys: [...s.keys], modifiers: [...s.modifiers],
    ...(s.children && { children: s.children.map((c) => cloneStep(c, newId)) }),
    ...(s.else && { else: s.else.map((c) => cloneStep(c, newId)) }),
    ...(s.cond && { cond: structuredClone(s.cond) }),
    ...(s.onFail && { onFail: { ...s.onFail } }),
  };
}

// ---- Kéo thả trong danh sách phẳng có thụt lề ----
// Danh sách hiển thị = các thẻ bước + dòng phụ (tiêu đề "Nếu không", cuối khối "+ Thêm bước"). Thả xuống thì vị trí
// mới suy ra từ dòng ngay phía trên chỗ thả: tiêu đề khối đang mở → bước đầu của khối; "Nếu không" → bước đầu nhánh
// else; cuối khối → ngay sau khối đó; bước thường → ngay sau bước đó, cùng danh sách.
export type Row =
  | { kind: 'step'; id: string; f: FlatStep; open: boolean } // open: khối đang mở (hiện bước con)
  | { kind: 'else'; id: string; parent: string; depth: number; label: string }
  | { kind: 'end'; id: string; parent: string; branch: Branch; depth: number };

// collapsed: id các khối đang gập; hide: id bước đang kéo (ẩn bước con của nó).
export function rows(steps: Step[], collapsed: Set<string>, hide: string | null = null): Row[] {
  const out: Row[] = [];
  const walk = (list: Step[], prefix: string, depth: number, parent: string | null, branch: Branch, off: number, active: boolean) => {
    list.forEach((s, i) => {
      const path = `${prefix}${off + i + 1}`;
      const on = active && s.enabled;
      const open = (isBlock(s) || hasElse(s)) && !collapsed.has(s.id) && s.id !== hide;
      out.push({ kind: 'step', id: s.id, f: { step: s, path, depth, loc: { parent, branch, index: i }, active: on }, open });
      if (!open) return;
      if (isBlock(s)) {
        walk(listOf(s, 'children'), `${path}.`, depth + 1, s.id, 'children', 0, on);
        if (s.type !== 'if') out.push({ kind: 'end', id: `${s.id}:end`, parent: s.id, branch: 'children', depth: depth + 1 });
      }
      if (hasElse(s)) {
        if (s.type === 'if') out.push({ kind: 'end', id: `${s.id}:end`, parent: s.id, branch: 'children', depth: depth + 1 });
        out.push({ kind: 'else', id: `${s.id}:else`, parent: s.id, depth, label: s.type === 'if' ? 'Nếu không' : 'Khi thất bại' });
        walk(listOf(s, 'else'), `${path}.`, depth + 1, s.id, 'else', s.type === 'if' ? listOf(s, 'children').length : 0, on);
        out.push({ kind: 'end', id: `${s.id}:elseEnd`, parent: s.id, branch: 'else', depth: depth + 1 });
      }
    });
  };
  walk(steps, '', 0, null, 'children', 0, true);
  return out;
}

// Vị trí mới khi thả bước id ngay dưới dòng `above` (null = đầu danh sách gốc).
export function dropLoc(steps: Step[], id: string, above: Row | null): Loc {
  if (!above) return { parent: null, branch: 'children', index: 0 };
  const idx = (parent: string | null, branch: Branch, afterId: string | null) => {
    const list = (parent === null ? steps : listOf(flatOf(steps, parent)!.step, branch)).filter((s) => s.id !== id);
    return afterId === null ? 0 : list.findIndex((s) => s.id === afterId) + 1;
  };
  switch (above.kind) {
    case 'else':
      return { parent: above.parent, branch: 'else', index: 0 };
    case 'end': {
      // cuối nhánh đúng của Nếu → đầu "Nếu không" không hợp lý; coi như cuối nhánh đó
      const p = flatOf(steps, above.parent)!;
      if (above.branch === 'children' && p.step.type === 'if') {
        const list = listOf(p.step, 'children').filter((s) => s.id !== id);
        return { parent: above.parent, branch: 'children', index: list.length };
      }
      if (above.branch === 'else' || !hasElse(p.step)) return { ...p.loc, index: idx(p.loc.parent, p.loc.branch, p.step.id) };
      const list = listOf(p.step, 'children').filter((s) => s.id !== id);
      return { parent: above.parent, branch: 'children', index: list.length };
    }
    case 'step': {
      const s = above.f.step;
      if (above.open && isBlock(s)) return { parent: s.id, branch: 'children', index: 0 };
      if (above.open && hasElse(s)) return { parent: s.id, branch: 'else', index: 0 };
      return { ...above.f.loc, index: idx(above.f.loc.parent, above.f.loc.branch, s.id) };
    }
  }
}

// ---- Điều kiện lồng nhau (bước 20c): path = chỉ số trong conds từ điều kiện gốc ([] = gốc) ----

export function condAt(c: Cond, path: number[]): Cond | undefined {
  let cur: Cond | undefined = c;
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
  const parent = path.slice(0, -1);
  const i = path[path.length - 1];
  return updateCondAt(c, parent, (g) => {
    const conds = (g.conds ?? []).filter((_, j) => j !== i);
    if (conds.length === 1) return { ...conds[0], not: !!conds[0].not !== !!g.not || undefined };
    return { ...g, conds };
  });
}
