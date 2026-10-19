// Phần tử màn hình qua Accessibility (bước 37): node Đọc phần tử. ElementQueryFields (cách tìm, app, bỏ dấu, chỗ thứ
// mấy) dùng chung cho Chờ / Click theo phần tử (E3).
import NumInput from '../NumInput';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { elementSelectorError } from '../../flow/element';
import { SelectorsField } from './WebActionFields';
import { errText, Field, fmtSec, hint, NumField, SwitchField, useNodeVars } from './shared';

export function ElementQueryFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const bad = elementSelectorError(s.selectors);
  return (
    <>
      <SelectorsField step={s} label="Cách tìm phần tử" value={s.selectors} onChange={(selectors) => set({ selectors })}
        placeholder={'text=₫\naria=Thêm vào giỏ\nrole=button[name="Mua ngay"]'} emptyText="Chưa nhập cách tìm – chạy sẽ báo lỗi"
        help={'text= chữ trên phần tử (chứa; text="Gửi" = đúng cả chuỗi), aria= nhãn (nút chỉ có icon), role=button[name="Gửi"] vai trò + tên.'}>
        {bad && <span className={errText}>{bad}</span>}
      </SelectorsField>
      <Field label="App" wide group help="Tuỳ chọn: tên app (vd Google Chrome, Safari) hoặc bundle id. Cửa sổ phải đang hiện trên màn – thu nhỏ / ở desktop khác thì không đọc được; trình duyệt chỉ thấy tab đang mở.">
        <VarTextInput value={s.elementApp ?? ''} onChange={(v) => set({ elementApp: v || undefined })} inputs={script.inputs ?? []} nodeVars={nodeVars}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} rows={1} compact placeholder="trống = cửa sổ trên cùng" />
      </Field>
      <SwitchField label="Bỏ qua dấu" on={!!s.textNoAccent} onChange={(textNoAccent) => set({ textNoAccent: textNoAccent || undefined })} />
      <Field label="Chỗ thứ" help="Nhiều phần tử khớp thì lấy cái thứ mấy (theo thứ tự đọc)">
        <NumInput value={Math.max(s.occurrence ?? 1, 1)} min={1} stepper onChange={(n) => set({ occurrence: n > 1 ? n : undefined })} />
      </Field>
    </>
  );
}

export function ElementReadFields({ step: s }: { step: Step }) {
  return (
    <>
      <ElementQueryFields step={s} />
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 250} min={20} unit="ms" onChange={(intervalMs) => updateStep(s.id, { intervalMs })} />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Không di chuột, không cuộn. Khối chữ / ô nhập → giá trị, phần tử khác → tên. Node sau dùng {`{{${s.ref}.output.text}}`}.
        Không thấy sau {fmtSec(s.timeoutMs)} → cổng Thất bại.
      </span>
    </>
  );
}
