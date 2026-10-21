// Node Tab trình duyệt (bước 43 B1) + Đi tới URL (B2, internal/engine/browsergo.go) – nhãn thao tác / trình duyệt, dòng chi tiết thẻ, cảnh báo thiếu, output. Chạy ở
// internal/engine/browsertab.go (kiểm tra: internal/script/tabcheck.go).
import type { Step } from '../types';
import type { OutputField } from './view';

export type TabOp = '' | 'list' | 'switch' | 'new' | 'close' | 'closeOthers' | 'reload' | 'back' | 'forward' | 'waitOpen' | 'waitClose';
export type BrowserName = '' | 'chrome' | 'edge' | 'brave' | 'coccoc';

export const TAB_OPS: { id: TabOp; label: string }[] = [
  { id: '', label: 'Đọc tab' }, { id: 'list', label: 'Liệt kê tab' }, { id: 'switch', label: 'Chuyển tới tab' },
  { id: 'new', label: 'Mở tab mới' }, { id: 'close', label: 'Đóng tab' }, { id: 'closeOthers', label: 'Đóng các tab khác' },
  { id: 'reload', label: 'Tải lại' }, { id: 'back', label: 'Quay lại' }, { id: 'forward', label: 'Tiến tới' },
  { id: 'waitOpen', label: 'Chờ tab mở' }, { id: 'waitClose', label: 'Chờ tab đóng' },
];

export const BROWSERS: { id: BrowserName; label: string }[] = [
  { id: '', label: 'Đang ở trước' }, { id: 'chrome', label: 'Chrome' }, { id: 'edge', label: 'Edge' }, { id: 'brave', label: 'Brave' }, { id: 'coccoc', label: 'Cốc Cốc' },
];

// Thao tác làm được với tab của node Mở (khớp TabFromOps – internal/script/tabcheck.go).
export const TAB_FROM_OPS: TabOp[] = ['', 'switch', 'close', 'closeOthers', 'reload', 'back', 'forward', 'waitClose'];
// Thao tác cần chọn tab (mở tab mới / liệt kê thì không).
export const tabTargeted = (op: TabOp) => op !== 'new' && op !== 'list';

type TabStep = Pick<Step, 'tabOp' | 'tab' | 'tabFrom' | 'url' | 'browser'> & { type?: Step['type'] };

// Tab node làm: "tab của n1" / “Gmail” / "tab đang hiện".
function tabTarget(s: TabStep, refOf: (id: string) => string): string {
  if (s.tabFrom && TAB_FROM_OPS.includes(s.tabOp ?? '')) return `tab của ${refOf(s.tabFrom)}`;
  return s.tab?.trim() ? `“${s.tab.trim()}”` : 'tab đang hiện';
}

// Dòng chi tiết thẻ: "Đóng tab · tab của n1", "Mở tab mới · example.com", "Liệt kê tab · “gmail” · Edge".
export function fmtBrowserTab(s: TabStep, refOf: (id: string) => string): string {
  const op = s.tabOp ?? '';
  const label = TAB_OPS.find((o) => o.id === op)?.label ?? op;
  const what = op === 'new' ? s.url?.trim() || 'trang trống' : op === 'list' ? (s.tab?.trim() ? `“${s.tab.trim()}”` : 'mọi tab') : tabTarget(s, refOf);
  const browser = s.browser ? ` · ${BROWSERS.find((b) => b.id === s.browser)?.label ?? s.browser}` : '';
  return `${label} · ${what}${browser}`;
}

// Dòng chi tiết node Đi tới URL: "example.com · tab của n1".
export const fmtBrowserGo = (s: TabStep, refOf: (id: string) => string): string =>
  `${s.url?.trim() || 'chưa nhập URL'} · ${tabTarget(s, refOf)}${s.browser ? ` · ${BROWSERS.find((b) => b.id === s.browser)?.label ?? s.browser}` : ''}`;

// Cảnh báo khi soạn (khớp validateTab / validateTabFrom); refs: id các node còn trong kịch bản.
export function browserTabWarning(s: TabStep, refs: Map<string, string>): string | undefined {
  const op = s.tabOp ?? '';
  if (s.type === 'browserGo' && !s.url?.trim()) return 'Chưa nhập URL';
  if (s.tabFrom && TAB_FROM_OPS.includes(op) && !refs.has(s.tabFrom)) return 'Node Mở đã chọn không còn – chọn lại';
  if (op === 'waitOpen' && !s.tab?.trim()) return 'Chưa nhập tiêu đề / URL tab cần chờ';
  if (op === 'waitClose' && !s.tab?.trim() && !s.tabFrom) return 'Chưa nhập tab cần chờ đóng';
  return undefined;
}

const TAB_INFO: OutputField[] = [{ key: 'url', desc: 'URL tab' }, { key: 'title', desc: 'tiêu đề tab' }, { key: 'browser', desc: 'trình duyệt (chrome, edge…)' },
  { key: 'index', desc: 'thứ tự tab trong cửa sổ (từ 1)' }, { key: 'tabId', desc: 'mã tab' }, { key: 'window', desc: 'mã cửa sổ' }];

// Output theo thao tác (khớp tabOutput / items – internal/engine/browsertab.go).
export function browserTabOutputs(s: Pick<Step, 'tabOp'> & { type?: Step['type'] }): OutputField[] {
  if (s.type === 'browserGo') return [...TAB_INFO, { key: 'ms', desc: 'thời gian tới khi tải xong (ms)' }];
  switch (s.tabOp ?? '') {
    case 'list': return [{ key: 'items', desc: 'danh sách tab [{title, url, browser, window, active, tabId}] – nối vào Lặp qua danh sách' }, { key: 'count', desc: 'số tab' }];
    case 'waitClose': return [];
    case 'closeOthers': return [...TAB_INFO, { key: 'closed', desc: 'số tab đã đóng' }];
  }
  return TAB_INFO;
}
