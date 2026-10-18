// Phần tử màn hình qua Accessibility (bước 37): node Đọc phần tử. ElementQueryFields (cách tìm, app, bỏ dấu, chỗ thứ
// mấy) dùng chung cho Chờ / Click theo phần tử (E3).
import NumInput from '../NumInput';
import { Switch } from '../RunSettings';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { elementSelectorError } from '../../flow/element';
import { SelectorsField } from './WebActionFields';
import { Field, fmtSec, hint, NumField, useNodeVars } from './shared';

export function ElementQueryFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const bad = elementSelectorError(s.selectors);
  return (
    <>
      <SelectorsField step={s} label="Cách tìm phần tử" value={s.selectors} onChange={(selectors) => set({ selectors })}
        placeholder={'text=₫\naria=Thêm vào giỏ\nrole=button[name="Mua ngay"]'} emptyText="Chưa nhập cách tìm – chạy sẽ báo lỗi">
        {bad && <span className="text-xs text-danger-text">{bad}</span>}
        <span className={hint}>
          <span className="font-mono">text=</span> chữ trên phần tử (chứa; <span className="font-mono">text="Gửi"</span> = đúng cả chuỗi),{' '}
          <span className="font-mono">aria=</span> nhãn (nút chỉ có icon),{' '}
          <span className="font-mono">role=button[name="Gửi"]</span> vai trò + tên.
        </span>
      </SelectorsField>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">App (tuỳ chọn)
        <VarTextInput value={s.elementApp ?? ''} onChange={(v) => set({ elementApp: v || undefined })} inputs={script.inputs ?? []} nodeVars={nodeVars}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} rows={1} placeholder="bỏ trống = app có cửa sổ trên cùng" />
        <span className={hint}>Tên app (vd Google Chrome, Safari) hoặc bundle id. Cửa sổ phải đang hiện trên màn – thu nhỏ / ở desktop khác thì không đọc được; trình duyệt chỉ thấy tab đang mở.</span>
      </div>
      <div className="col-span-2 flex flex-wrap items-center gap-x-4 gap-y-2 py-1 text-[13px] text-ink">
        <label className="flex items-center gap-2"><Switch on={!!s.textNoAccent} onChange={(textNoAccent) => set({ textNoAccent: textNoAccent || undefined })} label="Bỏ qua dấu" />Bỏ qua dấu</label>
        <span className="flex items-center gap-2">Chỗ thứ
          <NumInput value={Math.max(s.occurrence ?? 1, 1)} min={1} onChange={(n) => set({ occurrence: n > 1 ? n : undefined })} className="w-20" />
        </span>
      </div>
    </>
  );
}

export function ElementReadFields({ step: s }: { step: Step }) {
  return (
    <>
      <span className={`col-span-2 ${hint}`}>Đọc qua Accessibility của macOS – dùng được với Chrome thường đã đăng nhập và app khác. Cần quyền Accessibility.</span>
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
