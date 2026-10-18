// Node web (bước 25d): tìm phần tử bằng HTML qua Chrome của app. Mỗi dòng một cách tìm, dòng đầu ưu tiên.
// Thử tìm / Chọn phần tử trên trang – D3.
import Segmented from '../Segmented';
import NumInput from '../NumInput';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { parseSelector, SELECTOR_LABEL, selectorLines } from '../../flow/web';
import { ButtonField, Field, fmtSec, hint, NumField, useNodeVars } from './shared';

export function WebFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const lines = selectorLines(s.selectors);
  const parsed = lines.map(parseSelector);
  const bad = parsed.findIndex((p) => !p);
  const add = (key: string) => addInput({ key, label: key, type: 'text' });
  return (
    <>
      {s.type === 'webWait' && (
        <Field label="Chờ tới khi" wide group>
          <Segmented value={s.webGone ? 'gone' : 'shown'} onChange={(v) => set({ webGone: v === 'gone' || undefined })} full label="Chờ tới khi"
            options={[{ id: 'shown', label: 'Phần tử hiện' }, { id: 'gone', label: 'Phần tử mất' }]} />
        </Field>
      )}
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span className="flex items-center justify-between">Cách tìm phần tử
          <span className="text-xs">mỗi dòng một cách</span>
        </span>
        <VarTextInput value={(s.selectors ?? []).join('\n')} onChange={(v) => set({ selectors: v.split('\n') })}
          inputs={script.inputs ?? []} nodeVars={nodeVars} onCreateInput={add} rows={3} placeholder={'text=Xem lịch sử\naria=Lịch sử\ncss=#history'} />
        <span className={`text-xs ${lines.length === 0 || bad >= 0 ? 'text-danger-text' : 'text-muted'}`}>
          {lines.length === 0 ? 'Chưa nhập cách tìm – chạy sẽ báo lỗi'
            : bad >= 0 ? `Dòng “${lines[bad].trim()}” chưa có giá trị / sai dạng`
            : `${parsed.map((p) => SELECTOR_LABEL[p!.by] + (p!.exact ? ' (đúng cả chuỗi)' : '')).join(' → ')} · cách đầu thấy thì dùng`}
        </span>
        <span className={hint}>
          <span className="font-mono">text=</span> chữ trên trang (chứa, không phân biệt hoa thường; <span className="font-mono">text="Gửi"</span> = đúng cả chuỗi),{' '}
          <span className="font-mono">aria=</span> nhãn aria-label,{' '}
          <span className="font-mono">role=button[name="Gửi"]</span> vai trò + tên, <span className="font-mono">css=</span>, <span className="font-mono">xpath=</span>.
          Ưu tiên chữ / aria / role – tránh class tự sinh kiểu <span className="font-mono">x1n2onr6</span>.
        </span>
      </div>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">Tab (tuỳ chọn)
        <VarTextInput value={s.tab ?? ''} onChange={(tab) => set({ tab: tab || undefined })} inputs={script.inputs ?? []} nodeVars={nodeVars}
          onCreateInput={add} rows={1} placeholder="bỏ trống = tab đang hiện" />
        <span className={hint}>URL hoặc tiêu đề tab chứa chữ này. Trang phải mở bằng node Mở với “Chrome của app”.</span>
      </div>
      {s.type === 'webClick' && <>
        <Field label="Lệch so với tâm" wide>
          <div className="flex gap-2">
            <NumInput signed prefix="ΔX" value={s.dx ?? 0} onChange={(dx) => set({ dx })} className="flex-1" />
            <NumInput signed prefix="ΔY" value={s.dy ?? 0} onChange={(dy) => set({ dy })} className="flex-1" />
          </div>
        </Field>
        <ButtonField step={s} />
        <NumField label="Số lần click" step={s} field="clicks" min={1} />
        <NumField label="Giữa các click" step={s} field="clickInterval" unit="ms" />
      </>}
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 200} min={20} unit="ms" onChange={(intervalMs) => set({ intervalMs })} />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        {s.webGone ? `Phần tử còn sau ${fmtSec(s.timeoutMs)} → cổng Thất bại.`
          : `Không thấy sau ${fmtSec(s.timeoutMs)} → cổng Thất bại.${s.type === 'webClick' ? ' Phần tử ngoài màn hình thì tự cuộn tới; bị popup che thì chờ tiếp.' : ''}`}
        {' '}Lần đầu chuột đi vào trang vài điểm để đo toạ độ.
      </span>
    </>
  );
}
