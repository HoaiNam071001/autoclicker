// Selector của node web (bước 25d) – cùng công thức ParseSelector (internal/model/web.go), test chung testdata/selectors.json.
// Một dòng Step.selectors = "<cách>=<giá trị>", không tiền tố = text; text / aria trong nháy kép = đúng cả chuỗi.
export type SelectorKind = 'text' | 'aria' | 'role' | 'css' | 'xpath';
export const SELECTOR_KINDS: SelectorKind[] = ['text', 'aria', 'role', 'css', 'xpath'];
export const SELECTOR_LABEL: Record<SelectorKind, string> = { text: 'chữ', aria: 'aria', role: 'vai trò', css: 'CSS', xpath: 'XPath' };
// name: role – tên truy cập phải chứa (role=button[name="Gửi"], kiểu Playwright).
export interface Selector { by: SelectorKind; value: string; exact?: boolean; name?: string }

const ROLE_RE = /^([a-zA-Z]+)\s*(?:\[\s*name\s*=\s*(?:"([^"]*)"|([^\]"]*?))\s*\])?$/;

// null = dòng chưa có giá trị / sai dạng (chạy sẽ báo lỗi).
export function parseSelector(line: string): Selector | null {
  const t = line.trim();
  let sel: Selector = { by: 'text', value: t };
  const eq = t.indexOf('=');
  if (eq >= 0) {
    const by = t.slice(0, eq).trim().toLowerCase();
    const kind = SELECTOR_KINDS.find((k) => k === by);
    if (kind) sel = { by: kind, value: t.slice(eq + 1).trim() };
  }
  if ((sel.by === 'text' || sel.by === 'aria') && sel.value.length >= 2 && sel.value.startsWith('"') && sel.value.endsWith('"')) {
    sel = { ...sel, value: sel.value.slice(1, -1), exact: true };
  }
  if (sel.by === 'role' && sel.value !== '') {
    const m = ROLE_RE.exec(sel.value);
    if (!m) return null;
    const name = ((m[2] ?? '') + (m[3] ?? '')).trim();
    sel = { by: 'role', value: m[1].toLowerCase(), ...(name ? { name } : {}) };
  }
  return sel.value === '' ? null : sel;
}

// Các dòng có nội dung của ô selector (dòng trống bỏ qua – như Go).
export const selectorLines = (selectors: string[] | undefined) => (selectors ?? []).filter((l) => l.trim() !== '');

// Dòng chi tiết trên thẻ: cách tìm đầu tiên + số dự phòng.
export function fmtSelectors(selectors: string[] | undefined): string {
  const lines = selectorLines(selectors);
  const first = lines.length ? parseSelector(lines[0]) : null;
  if (!first) return 'chưa nhập cách tìm phần tử';
  const what = first.by === 'text' || first.by === 'aria' ? `“${first.value}”`
    : first.by === 'role' && first.name ? `${first.value} “${first.name}”` : first.value;
  return `${SELECTOR_LABEL[first.by]} ${what}${lines.length > 1 ? ` (+${lines.length - 1} dự phòng)` : ''}`;
}
