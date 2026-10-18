// Node web (bước 25d): tìm phần tử bằng HTML qua Chrome của app. Mỗi dòng một cách tìm, dòng đầu ưu tiên.
// Trạng thái Chrome của app, Chọn phần tử trên trang, Thử tìm – D3 (WebTools); trường riêng từng loại – WebActionFields.
import { useEffect } from 'react';
import NumInput from '../NumInput';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { browserStatusText, parseSelector, selectorLines, webFindText } from '../../flow/web';
import { cancelPick, openAppBrowser, pickElement, refreshBrowser, testWebFind, useWebState } from '../../webSteps';
import { inWails, openOverlay } from '../../windowMode';
import { WebTools } from './WebTools';
import { SelectorsField, WebClickFields, webHintText, WebReadFields, WebScrollFields, WebTypeFields, WebWaitFields } from './WebActionFields';
import { Field, fmtSec, hint, NumField, useNodeVars } from './shared';

// Nhãn ô cách tìm theo loại / cách cuộn.
const selectorsLabel = (s: Step) =>
  s.type === 'webType' ? 'Cách tìm ô nhập'
    : s.type === 'webScroll' && s.scrollMode === 'in' ? 'Cách tìm khung cần lăn'
    : s.type === 'webScroll' && s.scrollMode === 'until' ? 'Lăn tới khi thấy phần tử'
    : 'Cách tìm phần tử';

export function WebFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const lines = selectorLines(s.selectors);
  const web = useWebState();
  useEffect(() => { void refreshBrowser(); }, [s.id]);
  const find = web.finds[s.id];
  return (
    <>
      <WebTools status={inWails ? browserStatusText(web.status) : 'Chrome của app: chỉ dùng được trong app'} open={!!web.status?.open}
        enabled={inWails} picking={web.picking === s.id} finding={web.busy === s.id} canFind={lines.length > 0 && lines.every(parseSelector)}
        result={find ? { ...webFindText(find, s.tab, s.type === 'webRead'), elText: find.found ? find.text : undefined } : null}
        error={web.busy === '' && web.picking === '' ? web.error : ''}
        onOpen={openAppBrowser} onRefresh={refreshBrowser} onPick={() => pickElement(s.id)} onCancelPick={cancelPick}
        onFind={() => testWebFind(s.id)} onView={() => openOverlay('point', s.id)} />
      {s.type === 'webWait' && <WebWaitFields step={s} />}
      {s.type === 'webScroll' && <WebScrollFields step={s} />}
      <SelectorsField step={s} label={selectorsLabel(s)} value={s.selectors} onChange={(selectors) => set({ selectors })}
        placeholder={'text=Xem lịch sử\naria=Lịch sử\ncss=#history'} emptyText="Chưa nhập cách tìm – chạy sẽ báo lỗi">
        <span className={hint}>
          <span className="font-mono">text=</span> chữ trên trang (chứa, không phân biệt hoa thường; <span className="font-mono">text="Gửi"</span> = đúng cả chuỗi),{' '}
          <span className="font-mono">aria=</span> nhãn aria-label,{' '}
          <span className="font-mono">role=button[name="Gửi"]</span> vai trò + tên, <span className="font-mono">css=</span>, <span className="font-mono">xpath=</span>.
          Ưu tiên chữ / aria / role – tránh class tự sinh kiểu <span className="font-mono">x1n2onr6</span>.
        </span>
      </SelectorsField>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">Tab (tuỳ chọn)
        <VarTextInput value={s.tab ?? ''} onChange={(tab) => set({ tab: tab || undefined })} inputs={script.inputs ?? []} nodeVars={nodeVars}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} rows={1} placeholder="bỏ trống = tab đang hiện" />
        <span className={hint}>URL hoặc tiêu đề tab chứa chữ này. Trang phải mở bằng node Mở với “Chrome của app”.</span>
      </div>
      {s.type === 'webClick' && <WebClickFields step={s} />}
      {s.type === 'webType' && <WebTypeFields step={s} />}
      {s.type === 'webRead' && <WebReadFields step={s} />}
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 200} min={20} unit="ms" onChange={(intervalMs) => set({ intervalMs })} />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        {webHintText(s, fmtSec(s.timeoutMs))}
        {s.type !== 'webRead' && !(s.type === 'webWait' && s.webGone) && ' Lần đầu chuột đi vào trang vài điểm để đo toạ độ.'}
      </span>
    </>
  );
}
