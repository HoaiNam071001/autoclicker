// Chọn trình duyệt + tab cần làm – dùng chung node Tab trình duyệt và Đi tới URL (bước 43 B1 / B2): tab node Mở đã mở
// (tabFrom) hoặc tab có tiêu đề / URL chứa chuỗi (tab + tabMatch).
import Segmented from '../Segmented';
import Select from '../Select';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { BROWSERS, type BrowserName } from '../../flow/browserTab';
import { nodeTitle } from '../../flow/view';
import { WIN_MATCHES } from '../../flow/window';
import { Step } from '../../types';
import { errText, Field, selectBox, useNodeVars } from './shared';

// Ô nhận biến của inspector trình duyệt (URL, tab).
export function useTabVarProps(s: Step) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  return { inputs: script.inputs ?? [], nodeVars, rows: 1, compact: true, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
}

export function BrowserField({ step: s }: { step: Step }) {
  return (
    <Field label="Trình duyệt" wide group help="Đang ở trước = trình duyệt đang chạy có cửa sổ ở trước nhất. Chrome, Edge, Brave, Cốc Cốc thường (đã đăng nhập) – không phải Chrome của app.">
      <Select value={s.browser ?? ''} onChange={(v) => updateStep(s.id, { browser: (v || undefined) as BrowserName | undefined })} className={selectBox} aria-label="Trình duyệt"
        options={BROWSERS.map((b) => ({ value: b.id, label: b.label }))} />
    </Field>
  );
}

// fromOK: thao tác làm được với tab node Mở; placeholder / need: ô Tab trống nghĩa là gì, có bắt buộc không.
export function TabTargetFields({ step: s, fromOK, placeholder, need }: { step: Step; fromOK: boolean; placeholder: string; need: boolean }) {
  const { steps } = useScript();
  const varProps = useTabVarProps(s);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const opens = steps.filter((x) => x.type === 'open' && !x.webBrowser);
  const canFrom = fromOK && (opens.length > 0 || !!s.tabFrom);
  const from = canFrom ? opens.find((x) => x.id === s.tabFrom) : undefined;
  return (
    <>
      {canFrom && (
        <Field label="Tab node Mở" wide group help="Làm với đúng tab node Mở đã mở (không nhầm tab khác cùng trang). Bỏ chọn = tìm theo tiêu đề / URL.">
          <Select value={from?.id ?? ''} onChange={(v) => set({ tabFrom: v || undefined })} className={selectBox} aria-label="Tab của node Mở"
            options={[{ value: '', label: '— theo tiêu đề / URL —' }, ...opens.map((x) => ({ value: x.id, label: `${x.ref} · ${nodeTitle(x)}` }))]} />
          {s.tabFrom && !from && <span className={errText}>Node Mở đã chọn không còn – chọn lại</span>}
        </Field>
      )}
      {!from && (
        <>
          <Field label="Tab" wide group help="Tiêu đề hoặc URL tab chứa chuỗi này (không phân biệt hoa thường, không dấu).">
            <VarTextInput {...varProps} value={s.tab ?? ''} onChange={(v) => set({ tab: v || undefined })} placeholder={placeholder} />
            {need && !s.tab?.trim() && <span className={errText}>Chưa nhập tab cần chờ</span>}
          </Field>
          {!!s.tab?.trim() && (
            <Field label="So khớp" wide group>
              <Segmented value={s.tabMatch ?? ''} full label="So khớp" onChange={(v) => set({ tabMatch: v || undefined })} options={WIN_MATCHES} />
            </Field>
          )}
        </>
      )}
    </>
  );
}
