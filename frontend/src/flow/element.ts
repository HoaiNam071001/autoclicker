import { fmtMs, type Step } from '../types';

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
