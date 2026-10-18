import { fmtMs, type Step } from '../types';
import { parseSelector, selectorLines } from './web';

// Phần tử màn hình qua Accessibility (bước 37, internal/engine/elementstep.go): chữ hiển thị cho node Chờ / Click theo
// phần tử và điều kiện "thấy phần tử".

// Cách tìm đầu tiên (+ "+n" nếu còn dòng dự phòng); chưa có → "(chưa nhập)".
export function firstSelector(selectors?: string[]): string {
  const lines = (selectors ?? []).map((l) => l.trim()).filter(Boolean);
  if (!lines.length) return '(chưa nhập)';
  return `“${lines[0]}”${lines.length > 1 ? ` +${lines.length - 1}` : ''}`;
}

// Dòng chi tiết node: phần tử “role=button[name=Mua]” · Google Chrome (chỗ thứ 2) · tối đa 10 s
export function fmtElement(s: Step): string {
  if (!s.selectors?.some((l) => l.trim())) return 'chưa nhập cách tìm phần tử';
  const nth = (s.occurrence ?? 0) > 1 ? ` (chỗ thứ ${s.occurrence})` : '';
  return `phần tử ${firstSelector(s.selectors)}${nth} · ${s.elementApp?.trim() || 'app trên cùng'} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
}

// Dòng cách tìm không dùng được cho phần tử màn hình (css / xpath chỉ có ở node web) – khớp CheckSelector (element/find.go).
// Dòng có {{biến}} chỉ biết lúc chạy. '' = ổn.
export function elementSelectorError(selectors?: string[]): string {
  for (const line of selectorLines(selectors)) {
    if (line.includes('{{')) continue;
    const sel = parseSelector(line);
    if (sel && (sel.by === 'css' || sel.by === 'xpath')) return `“${line.trim()}”: phần tử màn hình chỉ tìm theo text= / aria= / role= – css / xpath chỉ dùng ở node web`;
  }
  return '';
}
