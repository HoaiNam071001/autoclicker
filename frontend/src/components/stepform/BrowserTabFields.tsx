// Node Tab trình duyệt (bước 43 B1): thao tác, trình duyệt, tab cần làm (của node Mở / theo tiêu đề · URL), URL tab mới.
import Segmented from '../Segmented';
import Select from '../Select';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { BROWSERS, TAB_FROM_OPS, TAB_OPS, tabTargeted, type BrowserName, type TabOp } from '../../flow/browserTab';
import { nodeTitle } from '../../flow/view';
import { WIN_MATCHES } from '../../flow/window';
import { Step } from '../../types';
import { errText, Field, fmtSec, hint, NumField, selectBox, useNodeVars } from './shared';

const OP_HINT: Record<TabOp, (t: string) => string> = {
  '': () => 'Chỉ đọc URL + tiêu đề tab, không đổi gì – nối Nếu để rẽ nhánh theo trang đang mở.',
  list: () => 'Mọi tab khớp (trống = mọi tab) vào output items – nối Lặp qua danh sách.',
  switch: (t) => `Chuyển tới tab và đưa cửa sổ lên trước. Chưa thấy tab thì chờ tới ${t}.`,
  new: () => 'Mở tab mới ở cửa sổ trên cùng của trình duyệt. Chưa mở trình duyệt nào → Thất bại (mở bằng node Mở trước).',
  close: (t) => `Đóng một tab (tab đầu tiên khớp). Chưa thấy tab thì chờ tới ${t}.`,
  closeOthers: () => 'Đóng mọi tab khác trong cùng cửa sổ, giữ tab đã chọn.',
  reload: () => 'Tải lại tab. Output đọc ngay sau khi bấm – trang có thể chưa tải xong.',
  back: () => 'Quay lại trang trước trong tab.',
  forward: () => 'Tiến tới trang sau trong tab.',
  waitOpen: (t) => `Chờ có tab khớp. Quá ${t} → cổng Thất bại.`,
  waitClose: (t) => `Chờ hết tab khớp (hoặc tab của node Mở đóng). Quá ${t} → cổng Thất bại.`,
};

export function BrowserTabFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const op = s.tabOp ?? '';
  const opens = script.steps.filter((x) => x.type === 'open' && !x.webBrowser);
  const canFrom = TAB_FROM_OPS.includes(op) && (opens.length > 0 || !!s.tabFrom);
  const from = canFrom ? opens.find((x) => x.id === s.tabFrom) : undefined;
  // thao tác không dùng tab của node Mở → bỏ tabFrom (ô chọn bị ẩn, để lại thì kiểm tra báo lỗi không sửa được)
  const setOp = (v: TabOp) => set({ tabOp: v || undefined, ...(TAB_FROM_OPS.includes(v) ? {} : { tabFrom: undefined }) });
  const varProps = { inputs: script.inputs ?? [], nodeVars, rows: 1, compact: true, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };

  return (
    <>
      <Field label="Làm gì" wide group>
        <Select value={op} onChange={(v) => setOp(v as TabOp)} className={selectBox} aria-label="Thao tác tab"
          options={TAB_OPS.map((o) => ({ value: o.id, label: o.label }))} />
      </Field>
      <Field label="Trình duyệt" wide group help="Đang ở trước = trình duyệt đang chạy có cửa sổ ở trước nhất. Chrome, Edge, Brave, Cốc Cốc thường (đã đăng nhập) – không phải Chrome của app.">
        <Select value={s.browser ?? ''} onChange={(v) => set({ browser: (v || undefined) as BrowserName | undefined })} className={selectBox} aria-label="Trình duyệt"
          options={BROWSERS.map((b) => ({ value: b.id, label: b.label }))} />
      </Field>
      {op === 'new' && (
        <Field label="URL" wide group help="Trống = trang trống. Tên miền tự thêm https://.">
          <VarTextInput {...varProps} value={s.url ?? ''} onChange={(v) => set({ url: v || undefined })} placeholder="https://" />
        </Field>
      )}
      {canFrom && (
        <Field label="Tab node Mở" wide group help="Làm với đúng tab node Mở đã mở (không nhầm tab khác cùng trang). Bỏ chọn = tìm theo tiêu đề / URL.">
          <Select value={from?.id ?? ''} onChange={(v) => set({ tabFrom: v || undefined })} className={selectBox} aria-label="Tab của node Mở"
            options={[{ value: '', label: '— theo tiêu đề / URL —' }, ...opens.map((x) => ({ value: x.id, label: `${x.ref} · ${nodeTitle(x)}` }))]} />
          {s.tabFrom && !from && <span className={errText}>Node Mở đã chọn không còn – chọn lại</span>}
        </Field>
      )}
      {op !== 'new' && !from && (
        <>
          <Field label="Tab" wide group help="Tiêu đề hoặc URL tab chứa chuỗi này (không phân biệt hoa thường, không dấu).">
            <VarTextInput {...varProps} value={s.tab ?? ''} onChange={(v) => set({ tab: v || undefined })}
              placeholder={op === 'list' ? 'trống = mọi tab' : op.startsWith('wait') ? 'vd mail.google.com' : 'trống = tab đang hiện'} />
            {op.startsWith('wait') && !s.tab?.trim() && <span className={errText}>Chưa nhập tab cần chờ</span>}
          </Field>
          {!!s.tab?.trim() && (
            <Field label="So khớp" wide group>
              <Segmented value={s.tabMatch ?? ''} full label="So khớp" onChange={(v) => set({ tabMatch: v || undefined })} options={WIN_MATCHES} />
            </Field>
          )}
        </>
      )}
      {tabTargeted(op) && <NumField label="Chờ tối đa" info="Chưa thấy tab thì chờ tối đa – quá giờ đi cổng Thất bại" step={s} field="timeoutMs" unit="ms" min={0} />}
      <span className={`col-span-2 -mt-1 ${hint}`}>{OP_HINT[op](fmtSec(s.timeoutMs))} Hiện chỉ chạy trên macOS – lần đầu macOS hỏi cho phép điều khiển trình duyệt.</span>
    </>
  );
}
