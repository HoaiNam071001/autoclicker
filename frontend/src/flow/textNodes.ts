// Node xử lý chữ (bước 42 G4b): Lọc mẫu, Tách chữ, Làm sạch chữ, Số & ngày – lựa chọn, nhãn cổng, dòng chi tiết thẻ,
// cảnh báo, output. Chạy ở internal/engine/textnodes.go (logic internal/textop); view.ts gọi các hàm ở đây.
import type { Step, StepType } from '../types';

export type TextType = 'textMatch' | 'textSplit' | 'textClean' | 'textNumber';
export const TEXT_TYPES: TextType[] = ['textMatch', 'textSplit', 'textClean', 'textNumber'];
export const isTextNode = (t: StepType): t is TextType => (TEXT_TYPES as StepType[]).includes(t);

// Khớp model.CaseModes / NumOps / NumStyles ('' = mặc định)
export type CaseMode = '' | 'lower' | 'upper' | 'title';
export const CASE_MODES: { id: CaseMode; label: string }[] = [
  { id: '', label: 'Giữ' }, { id: 'lower', label: 'thường' }, { id: 'upper', label: 'HOA' }, { id: 'title', label: 'Hoa Đầu' },
];
export type NumOp = '' | 'format' | 'date';
export const NUM_OPS: { id: NumOp; label: string }[] = [{ id: '', label: 'Lấy số' }, { id: 'format', label: 'Định dạng số' }, { id: 'date', label: 'Ngày giờ' }];
export type NumStyle = '' | 'en' | 'plain';
export const NUM_STYLES: { id: NumStyle; label: string }[] = [{ id: '', label: '1.234,5' }, { id: 'en', label: '1,234.5' }, { id: 'plain', label: '1234.5' }];

// Nhãn cổng: Lọc mẫu / Số & ngày có cổng fail khi không ra kết quả; Tách / Làm sạch chỉ có Tiếp.
export function textPortLabel(t: TextType, port: string): string {
  if (port !== 'fail') return t === 'textMatch' ? 'Khớp' : t === 'textNumber' ? 'Đọc được' : 'Tiếp';
  return t === 'textMatch' ? 'Không khớp' : 'Không đọc được';
}

const short = (v?: string) => {
  const t = (v ?? '').trim();
  return t.length > 28 ? `${t.slice(0, 27)}…` : t;
};

// Dòng chi tiết thẻ.
export function textDetail(s: Step): string {
  switch (s.type) {
    case 'textMatch': return `/${short(s.pattern) || '?'}/${s.matchAll ? ' · mọi chỗ' : ''}`;
    case 'textSplit': return `tách theo ${s.separator ? `“${s.separator}”` : 'mỗi dòng'}`;
    case 'textClean': {
      const parts = [
        s.replaceFrom && `thay “${short(s.replaceFrom)}”`, s.stripAccents && 'bỏ dấu',
        s.caseMode && CASE_MODES.find((c) => c.id === s.caseMode)?.label.toLowerCase(), s.trim && 'gọn khoảng trắng',
      ].filter(Boolean);
      return parts.join(' · ') || 'chưa chọn việc';
    }
    case 'textNumber':
      if (s.numOp === 'date') return `ngày → ${s.dateLayout || 'dd/MM/yyyy'}`;
      if (s.numOp === 'format') return `định dạng ${NUM_STYLES.find((x) => x.id === (s.numStyle ?? ''))?.label} · ${s.decimals ?? 0} số lẻ`;
      return `lấy số từ ${short(s.value) || '—'}`;
    default: return '';
  }
}

// Thiếu dữ liệu (chạy sẽ báo lỗi) – undefined = đủ. Regex chỉ kiểm ở Go (RE2 khác regex JS).
export function textWarning(s: Step): string | undefined {
  if (!s.value?.trim() && !(s.type === 'textNumber' && s.numOp === 'date')) return 'Chưa có chữ cần xử lý';
  if (s.type === 'textMatch' && !s.pattern) return 'Chưa nhập mẫu';
  return undefined;
}

export type TextOutput = { key: string; desc: string };
export function textOutputs(s: Step): TextOutput[] {
  switch (s.type) {
    case 'textMatch': return [{ key: 'text', desc: 'chỗ khớp đầu (nhóm () đầu nếu có)' }, { key: 'items', desc: s.matchAll ? 'mọi chỗ khớp – mảng' : 'chỗ khớp đầu – mảng 1 phần tử' },
      { key: 'count', desc: 'số chỗ khớp' }, { key: 'groups', desc: 'các nhóm () của chỗ đầu – mảng' }];
    case 'textSplit': return [{ key: 'items', desc: 'các phần – mảng (nối vào Lặp qua danh sách)' }, { key: 'count', desc: 'số phần' }];
    case 'textClean': return [{ key: 'text', desc: 'chữ sau khi làm sạch' }];
    case 'textNumber':
      if (s.numOp === 'date') return [{ key: 'text', desc: 'ngày giờ theo mẫu' }, { key: 'timestamp', desc: 'giây Unix – so ngày bằng Tính toán' }];
      return [{ key: s.numOp === 'format' ? 'text' : 'value', desc: s.numOp === 'format' ? 'số đã định dạng' : 'số – dùng trong Tính toán' },
        { key: s.numOp === 'format' ? 'value' : 'text', desc: s.numOp === 'format' ? 'số gốc' : 'số dạng chữ' }];
    default: return [];
  }
}
