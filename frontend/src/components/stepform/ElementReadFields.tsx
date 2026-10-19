// Node Đọc phần tử (bước 37 E2b + E4): tìm (phần tử / nhãn bằng OCR) → đọc (chính nó / chữ cạnh nhãn) → lọc (mẫu, số,
// chờ giá trị đổi).
import NumInput, { fieldBox } from '../NumInput';
import Segmented from '../Segmented';
import { updateStep } from '../../store';
import { Step } from '../../types';
import { patternError } from '../../flow/element';
import { errText, Field, fmtSec, hint, NumField, SwitchField } from './shared';
import { ElementSearchFields } from './ElementFields';
import { CapturePermission, TextSearchFields } from './ImageFields';

const READ_BY_HELP: Record<string, string> = {
  element: 'Tìm trong cây phần tử của app (Accessibility) – chữ chính xác, đọc được cả phần tử ngoài khung nhìn.',
  text: 'Tìm nhãn bằng đọc chữ (OCR) – khi app không có cây phần tử (canvas, game, app tự vẽ). Phải chọn đọc bên phải / bên dưới nhãn.',
};
const READ_NEAR_HELP = 'Chính nó: chữ / giá trị của phần tử. Bên phải / bên dưới: đọc chữ (OCR) cạnh nhãn – khi giá trị không có trong cây (vẽ bằng canvas), cần quyền Screen Recording.';

export default function ElementReadFields({ step: s }: { step: Step }) {
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const byText = s.findBy === 'text';
  const badPattern = patternError(s.readPattern);
  return (
    <>
      <Field label="Tìm theo" wide group help={READ_BY_HELP[byText ? 'text' : 'element']}>
        <Segmented value={byText ? 'text' : 'element'} full label="Tìm theo"
          onChange={(v) => set(v === 'text' ? { findBy: 'text', readNear: s.readNear ?? 'right' } : { findBy: undefined })}
          options={[{ id: 'element', label: 'Phần tử' }, { id: 'text', label: 'Chữ (OCR)' }]} />
      </Field>
      {byText ? <TextSearchFields step={s} label="Chữ của nhãn" placeholder="Số dư:" /> : <ElementSearchFields step={s} />}
      {byText && (
        <Field label="Tìm trong" wide group help="Vùng nhỏ đọc nhanh hơn (cả màn ~1 s mỗi lần) và ít nhầm nhãn trùng ở chỗ khác. Cả màn hình = cả màn chứa vùng.">
          <Segmented value={s.searchArea ?? 'region'} onChange={(v) => set({ searchArea: v === 'screen' ? 'screen' : undefined })} full
            options={[{ id: 'region', label: 'Vùng của kịch bản' }, { id: 'screen', label: 'Cả màn hình' }]} />
        </Field>
      )}
      <Field label="Đọc" wide group help={READ_NEAR_HELP}>
        <Segmented value={s.readNear ?? 'self'} full label="Đọc" onChange={(v) => set({ readNear: v === 'self' ? undefined : v })}
          options={[...(byText ? [] : [{ id: 'self' as const, label: 'Chính nó' }]), { id: 'right' as const, label: 'Bên phải' }, { id: 'below' as const, label: 'Bên dưới' }]} />
      </Field>
      {s.readNear && (
        <Field label="Khoảng đọc" help="Giá trị phải bắt đầu trong chừng này tính từ mép nhãn – lấy dòng chữ gần nhãn nhất">
          <NumInput value={s.readSpan || 300} min={20} unit="px" onChange={(n) => set({ readSpan: n === 300 ? undefined : n })} />
        </Field>
      )}
      <Field label="Lọc mẫu" help={'Regex (cú pháp Go): có nhóm () thì lấy nhóm đầu, vd Mã: (\\w+). Chưa khớp thì chờ đọc lại.'}>
        <input value={s.readPattern ?? ''} onChange={(e) => set({ readPattern: e.target.value || undefined })} placeholder="trống = cả chữ"
          className={`${fieldBox} w-full font-mono text-[12.5px]`} spellCheck={false} />
        {badPattern && <span className={errText}>{badPattern}</span>}
      </Field>
      <SwitchField label="Lấy số" help="Chỉ lấy số đầu tiên, bỏ đơn vị và dấu phân cách: 1.234.000₫ → 1234000, $1,234.56 → 1234.56"
        on={!!s.readNumber} onChange={(readNumber) => set({ readNumber: readNumber || undefined })} />
      <SwitchField label="Chờ giá trị đổi" help="Lần đọc đầu làm mốc, chờ tới khi giá trị khác mốc (output old = giá trị cũ). Màn không đổi thì không đọc chữ lại."
        on={!!s.readChange} onChange={(readChange) => set({ readChange: readChange || undefined })} />
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 250} min={20} unit="ms" onChange={(intervalMs) => updateStep(s.id, { intervalMs })} />
      </Field>
      {s.readNear && <CapturePermission />}
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Không di chuột, không cuộn{s.readNear ? ' – nhãn phải đang hiện trên màn' : '. Khối chữ / ô nhập → giá trị, phần tử khác → tên'}.
        Node sau dùng {`{{${s.ref}.output.text}}`}. {s.readChange ? 'Không đổi' : 'Không đọc được'} sau {fmtSec(s.timeoutMs)} → cổng Thất bại.
      </span>
    </>
  );
}
