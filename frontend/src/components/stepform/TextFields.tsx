// Node mở URL / app và node gõ chữ – ô chữ nhận {{biến}}.
import Segmented from '../Segmented';
import Select from '../Select';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step, TypeMode } from '../../types';
import { WEB_BASIC_NOTE } from '../../flow/web';
import { errText, Field, hint, selectBox, useNodeVars, varHint } from './shared';

// Trình duyệt gợi ý cho bước Mở (Windows: open.go đổi sang tên lệnh chrome / msedge / firefox)
const APP_BROWSER = '__app'; // giá trị ô chọn cho Step.webBrowser (bước 25d) – không phải tên app
const BROWSERS = ['Google Chrome', 'Microsoft Edge', 'Firefox', ...(navigator.platform.startsWith('Mac') ? ['Safari'] : [])];

export function OpenFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
          <>
        <Field label="URL / app" wide group side={varHint}
          help="URL, đường dẫn file hoặc tên app. Vd youtube.com/… (tự thêm https://), Spotify, đường dẫn file / app. Biến trong URL tự mã hoá (dấu cách, tiếng Việt).">
          <VarTextInput value={s.target ?? ''} onChange={(target) => set({ target })} inputs={script.inputs ?? []} nodeVars={nodeVars} rows={2}
            placeholder="youtube.com/results?search_query={{bai_hat}}" onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
        </Field>
        <Field label="Mở bằng" group help={s.webBrowser
          ? 'Chrome / Edge / Brave với profile riêng của app (lần đầu phải đăng nhập lại trang, sau đó tự nhớ); đang mở thì mở tab mới. Chỉ mở URL / file. ' + WEB_BASIC_NOTE
          : 'Không chờ trang tải xong – đặt “Chờ sau khi xong” đủ lâu (mặc định 2 s).'}>
          <Select value={s.webBrowser ? APP_BROWSER : s.app ?? ''} className={selectBox} aria-label="Mở bằng"
            onChange={(v) => set(v === APP_BROWSER ? { webBrowser: true } : { webBrowser: undefined, app: v })}
            options={[{ value: '', label: 'Mặc định của máy' }, { value: APP_BROWSER, label: 'Chrome của app – node web điều khiển được' },
              ...[...BROWSERS, ...(s.app && !BROWSERS.includes(s.app) ? [s.app] : [])].map((b) => ({ value: b, label: b }))]} />
        </Field>
      </>
  );
}

export function TypeFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
    <>
      <Field label="Nội dung gõ" wide group side={varHint}
        help="Gõ vào ô đang focus – thường đặt sau một bước click. \{{ để gõ ra {{.">
        <VarTextInput value={s.text} onChange={(text) => set({ text })} inputs={script.inputs ?? []} nodeVars={nodeVars}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
        {s.text ? <span className={hint}>{[...s.text].length} ký tự</span> : <span className={errText}>Chưa nhập nội dung – chạy sẽ báo lỗi</span>}
      </Field>
      <TypeModeField step={s} />
    </>
  );
}

const TYPE_MODE_HELP: Record<TypeMode, string> = {
  '': 'Windows: dán (trừ khi bật giống người); macOS: gõ từng phím.',
  keys: 'Gõ từng phím – có nhịp giống người. Bộ gõ tiếng Việt (UniKey, EVKey) đang bật có thể làm méo chữ.',
  paste: 'Dán qua clipboard rồi trả lại chữ cũ – không bị bộ gõ tiếng Việt làm méo. Ô chặn dán thì không nhập được.',
};

// Cách nhập chữ của node Gõ chữ / Gõ vào ô web (model.TypeModes).
export function TypeModeField({ step: s }: { step: Step }) {
  const mode = s.typeMode ?? '';
  return (
    <Field label="Cách nhập" wide group help={TYPE_MODE_HELP[mode]}>
      <Segmented value={mode} onChange={(v) => updateStep(s.id, { typeMode: (v || undefined) as TypeMode | undefined })} full label="Cách nhập"
        options={[{ id: '', label: 'Tự chọn' }, { id: 'keys', label: 'Gõ phím' }, { id: 'paste', label: 'Dán' }]} />
    </Field>
  );
}
