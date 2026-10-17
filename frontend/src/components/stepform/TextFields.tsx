// Node mở URL / app và node gõ chữ – ô chữ nhận {{biến}}.
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { Field, hint, selectBox, useNodeVars } from './shared';

// Trình duyệt gợi ý cho bước Mở (Windows: open.go đổi sang tên lệnh chrome / msedge / firefox)
const BROWSERS = ['Google Chrome', 'Microsoft Edge', 'Firefox', ...(navigator.platform.startsWith('Mac') ? ['Safari'] : [])];

export function OpenFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
          <>
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">
          <span className="flex items-center justify-between">URL, đường dẫn hoặc tên app
            <span className="text-xs">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
          </span>
          <VarTextInput value={s.target ?? ''} onChange={(target) => set({ target })} inputs={script.inputs ?? []} nodeVars={nodeVars} rows={2}
            placeholder="youtube.com/results?search_query={{bai_hat}}" onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
          <span className={hint}>
            Vd <span className="font-mono">youtube.com/…</span> (tự thêm https://), <span className="font-mono">Spotify</span>, đường dẫn file / app.
            Biến trong URL tự mã hoá (dấu cách, tiếng Việt).
          </span>
        </div>
        <Field label="Mở bằng" wide>
          <select value={s.app ?? ''} onChange={(e) => set({ app: e.target.value })} className={selectBox}>
            <option value="">Mặc định của máy</option>
            {[...BROWSERS, ...(s.app && !BROWSERS.includes(s.app) ? [s.app] : [])].map((b) => <option key={b}>{b}</option>)}
          </select>
          <span className={hint}>Không chờ trang tải xong – đặt “Chờ sau bước” đủ lâu (mặc định 2 s).</span>
        </Field>
      </>
  );
}

export function TypeFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
          <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span className="flex items-center justify-between">Nội dung gõ
          <span className="text-xs">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
        </span>
        <VarTextInput value={s.text} onChange={(text) => set({ text })} inputs={script.inputs ?? []} nodeVars={nodeVars}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
        <span className={`text-xs ${s.text ? '' : 'text-danger-text'}`}>
          {s.text ? `${[...s.text].length} ký tự · ` : 'Chưa nhập nội dung – chạy sẽ báo lỗi · '}
          Gõ vào ô đang focus – thường đặt sau một bước click. Tắt bộ gõ tiếng Việt khi chạy. <span className="font-mono">{'\\{{'}</span> để gõ ra <span className="font-mono">{'{{'}</span>.
        </span>
      </div>
  );
}
