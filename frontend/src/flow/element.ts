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

// Dòng chi tiết node Đọc phần tử (+ E4: đọc cạnh nhãn, lọc): nhãn “Giá” (OCR) → bên phải · số · chờ đổi · tối đa 10 s
export function fmtElementRead(s: Step): string {
  const near = s.readNear ? ` → ${s.readNear === 'below' ? 'bên dưới' : 'bên phải'}` : '';
  const filter = [s.readPattern ? 'lọc mẫu' : '', s.readNumber ? 'số' : '', s.readChange ? 'chờ đổi' : ''].filter(Boolean).map((t) => ` · ${t}`).join('');
  if (s.findBy !== 'text') return fmtElement(s).replace(' · ', `${near}${filter} · `);
  if (!s.findText?.trim()) return 'chưa nhập chữ của nhãn';
  const nth = (s.occurrence ?? 0) > 1 ? ` (chỗ thứ ${s.occurrence})` : '';
  return `nhãn “${s.findText}”${nth} (OCR)${near}${filter} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
}

// Dòng chi tiết node Đọc chữ trang: cả cửa sổ Google Chrome · gồm phần ngoài khung nhìn · tối đa 10 s
export function fmtPageText(s: Step): string {
  const app = s.elementApp?.trim() || 'app trên cùng';
  const scope = s.selectors?.some((l) => l.trim()) ? `trong ${firstSelector(s.selectors)} · ${app}` : `cả cửa sổ ${app}`;
  return `${scope}${s.pageAll ? ' · gồm phần ngoài khung nhìn' : ''} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
}

// Dòng chi tiết node Đọc bảng: bảng đầu tiên của Google Chrome → ghi ra don-hang.csv · tối đa 10 s
export function fmtReadTable(s: Step): string {
  const app = s.elementApp?.trim() || 'app trên cùng';
  const nth = (s.occurrence ?? 0) > 1 ? `bảng thứ ${s.occurrence}` : 'bảng đầu tiên';
  const scope = s.selectors?.some((l) => l.trim()) ? `bảng trong ${firstSelector(s.selectors)} · ${app}` : `${nth} của ${app}`;
  const file = s.filePath?.trim() ? ` → ghi ${s.filePath.trim().split(/[\\/]/).pop()}` : '';
  return `${scope}${file} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
}

// Dòng cách tìm không dùng được cho phần tử màn hình (css / xpath chỉ có ở node web) – khớp CheckSelector (element/find.go).
// Dòng có {{biến}} chỉ biết lúc chạy. '' = ổn.
export function elementSelectorError(selectors?: string[]): string {
  for (const line of selectorLines(selectors)) {
    if (line.includes('{{')) continue;
    const sel = parseSelector(line);
    if (sel && sel.by !== 'text' && sel.by !== 'aria' && sel.by !== 'role') return `“${line.trim()}”: phần tử màn hình chỉ tìm theo text= / aria= / role= – ${sel.by}= chỉ dùng ở node web`;
  }
  return '';
}

// Mẫu lọc của Đọc phần tử không đọc được (báo ngay dưới ô). Go dùng RE2: cú pháp JS chỉ gần đúng – (?i) JS không có nên
// bỏ trước khi thử; lỗi thật Go báo khi chạy (script/element.go). '' = ổn.
export function patternError(pattern?: string): string {
  if (!pattern) return '';
  try {
    new RegExp(pattern.replace(/^\(\?[imsU]+\)/, ''));
    return '';
  } catch {
    return 'Mẫu lọc sai cú pháp (thiếu ngoặc, ký tự đặc biệt chưa có \\ đứng trước…)';
  }
}
