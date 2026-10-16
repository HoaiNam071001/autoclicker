// Form Chạy kịch bản (artboard RunDialog): chọn bộ giá trị (preset), điền input theo từng bước, ô khác mặc định
// có viền ring-focus + "khôi phục", khung "Sẽ gõ" xem trước câu sau khi thay biến, "Lần sau chạy luôn".
// Bước 21b mới nối Go (resolveInputs) – tới lúc đó dùng ở màn #demo.
import { useEffect, useMemo, useState } from 'react';
import { Play } from '../icons';
import { Switch } from '../RunSettings';
import { fmtValue, InputValue, ScriptInput, splitVars } from './model';

export type RunSection = { steps: string; title: string; keys: string[]; template?: string }; // template: câu có {{biến}} để xem trước
export type Preset = { name: string; values: Record<string, InputValue> };

type Props = {
  scriptName: string;
  inputs: ScriptInput[];
  sections: RunSection[];
  presets: Preset[];
  countdown: boolean;
  onRun: (values: Record<string, InputValue>, remember: boolean) => void;
  onClose: () => void;
};

const defaults = (inputs: ScriptInput[]) => Object.fromEntries(inputs.map((i) => [i.key, i.default ?? (i.type === 'bool' ? false : '')]));
const empty = (v: InputValue | undefined) => v === undefined || v === '';

export default function RunDialog({ scriptName, inputs, sections, presets, countdown, onRun, onClose }: Props) {
  const byKey = useMemo(() => new Map(inputs.map((i) => [i.key, i])), [inputs]);
  const [preset, setPreset] = useState(presets[0]?.name ?? '');
  const [values, setValues] = useState<Record<string, InputValue>>(() => ({ ...defaults(inputs), ...presets[0]?.values }));
  const [remember, setRemember] = useState(false);
  const missing = inputs.filter((i) => i.required && empty(values[i.key]));

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose();
      if (e.key === 'F8' && missing.length === 0) onRun(values, remember);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  });

  const pickPreset = (name: string) => {
    setPreset(name);
    setValues({ ...defaults(inputs), ...presets.find((p) => p.name === name)?.values });
  };
  const set = (k: string, v: InputValue) => {
    setValues((vs) => ({ ...vs, [k]: v }));
    setPreset('');
  };

  return (
    <div className="fixed inset-0 z-30 flex items-start justify-center bg-black/25 px-4 pt-[6vh]" onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div role="dialog" aria-label="Chạy kịch bản"
        className="flex max-h-[88vh] w-[560px] max-w-full flex-col overflow-hidden rounded-panel border border-line bg-surface text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
        <header className="flex items-center gap-3 px-[18px] pt-4 pb-3">
          <span className="grid h-9 w-9 place-items-center rounded-[9px] bg-accent-soft text-accent-text"><Play /></span>
          <div className="flex-1">
            <div className="text-base font-semibold">Chạy kịch bản</div>
            <div className="text-[13px] text-muted">{scriptName}, {inputs.length} input</div>
          </div>
          <button onClick={onClose} aria-label="Đóng" className="grid h-8 w-8 place-items-center rounded-md text-muted hover:bg-hover hover:text-ink">
            <svg width="12" height="12" viewBox="0 0 10 10" aria-hidden><path d="M0 0l10 10M10 0L0 10" stroke="currentColor" strokeWidth="1.3" /></svg>
          </button>
        </header>

        <div className="flex items-center gap-2 border-b border-line px-[18px] pb-3.5">
          <span className="text-[13px] text-muted">Bộ giá trị</span>
          <select value={preset} onChange={(e) => pickPreset(e.target.value)}
            className="h-8 flex-1 rounded-control border border-line bg-field px-2 text-ink outline-none focus:ring-focus">
            {preset === '' && <option value="">(đã sửa – chưa lưu thành bộ)</option>}
            {presets.map((p) => <option key={p.name}>{p.name}</option>)}
          </select>
          <button className="h-8 rounded-control border border-line bg-surface px-2.5 text-[13px] hover:bg-hover" title="Lưu các giá trị đang điền thành một bộ (lưu trên máy)">Lưu bộ mới</button>
        </div>

        <div className="flex min-h-0 flex-1 flex-col gap-3.5 overflow-y-auto px-[18px] py-3.5">
          {sections.map((sec) => (
            <section key={sec.steps} className="flex flex-col gap-2">
              <span className="text-xs font-semibold text-muted"><span className="font-mono">{sec.steps}</span>  {sec.title}</span>
              <div className="grid grid-cols-2 gap-2.5">
                {sec.keys.map((k) => byKey.get(k)).filter((i): i is ScriptInput => !!i).map((inp) => (
                  <Field key={inp.key} inp={inp} value={values[inp.key]} onChange={(v) => set(inp.key, v)} />
                ))}
              </div>
              {sec.template && <Preview template={sec.template} values={values} />}
            </section>
          ))}
        </div>

        <footer className="flex flex-col gap-3 border-t border-line bg-chrome px-[18px] py-3.5">
          <label className="flex items-center gap-2 text-[13px]">
            <input type="checkbox" checked={remember} onChange={(e) => setRemember(e.target.checked)} className="h-4 w-4 accent-accent" />
            Lần sau chạy luôn với giá trị này, không hỏi lại
          </label>
          <div className="flex items-center gap-2">
            <span className={`text-xs ${missing.length ? 'text-warn' : 'text-muted'}`}>
              {missing.length ? `Chưa nhập: ${missing.map((m) => m.label).join(', ')}` : countdown ? 'Đếm ngược 3 giây' : 'Chạy ngay'}
            </span>
            <button onClick={onClose} className="ml-auto h-9 rounded-control border border-line bg-surface px-3.5 hover:bg-hover">Huỷ</button>
            <button onClick={() => onRun(values, remember)} disabled={missing.length > 0}
              className="btn-primary flex h-9 items-center gap-2.5 rounded-control px-[18px] font-semibold disabled:opacity-50">
              <Play className="h-[11px] w-[11px]" />Chạy <span className="font-mono text-[11px] font-normal opacity-80">F8</span>
            </button>
          </div>
        </footer>
      </div>
    </div>
  );
}

function Field({ inp, value, onChange }: { inp: ScriptInput; value: InputValue; onChange: (v: InputValue) => void }) {
  const changed = !empty(inp.default) && String(value) !== String(inp.default);
  const wide = inp.type === 'text' || inp.type === 'multiline';
  const box = `flex items-center rounded-control px-2.5 ${changed ? 'ring-focus' : 'border border-line bg-field focus-within:ring-focus'}`;
  const req = inp.required && empty(value);
  return (
    <label className={`flex flex-col gap-[5px] ${wide ? 'col-span-2' : ''}`}>
      <span>
        {inp.label}{inp.required && <span className="text-warn"> *</span>} <span className="font-mono text-[11px] text-muted">{inp.key}</span>
      </span>
      {inp.type === 'multiline' ? (
        <textarea value={String(value)} rows={2} onChange={(e) => onChange(e.target.value)}
          className={`${box} min-h-[58px] py-2 outline-none`} />
      ) : inp.type === 'bool' ? (
        <Switch on={!!value} onChange={onChange} label={inp.label} />
      ) : inp.type === 'select' ? (
        <select value={String(value)} onChange={(e) => onChange(e.target.value)} className={`${box} h-[34px] outline-none`}>
          {inp.options?.map((o) => <option key={o}>{o}</option>)}
        </select>
      ) : (
        <span className={`${box} h-[34px] gap-1`}>
          {inp.prefix && <span className="font-mono text-[13px] text-muted">{inp.prefix}</span>}
          <input value={String(value)} onChange={(e) => onChange(e.target.value)}
            inputMode={inp.type === 'number' || inp.type === 'ms' ? 'numeric' : undefined}
            className={`min-w-0 flex-1 bg-transparent outline-none ${inp.type === 'text' ? '' : 'font-mono text-[13px]'} ${inp.prefix ? 'font-mono text-[13px]' : ''}`} />
          {inp.unit && <span className="text-muted">{inp.unit}</span>}
        </span>
      )}
      {changed ? (
        <span className="flex items-center gap-1.5 text-xs text-accent-text">
          đã đổi, mặc định {fmtValue(inp.default, inp.unit)}
          <button type="button" onClick={() => onChange(inp.default!)} className="text-muted underline">khôi phục</button>
        </span>
      ) : req ? <span className="text-xs text-warn">Bắt buộc – chưa nhập</span>
        : inp.hint && <span className="text-xs text-muted">{inp.hint}</span>}
    </label>
  );
}

// "Sẽ gõ": câu sau khi thay biến, phần thay thế đậm màu accent-text
function Preview({ template, values }: { template: string; values: Record<string, InputValue> }) {
  return (
    <div className="flex flex-col gap-1 rounded-lg bg-soft px-3 py-2.5">
      <span className="text-xs text-muted">Sẽ gõ</span>
      <span className="whitespace-pre-wrap">
        {splitVars(template).map((s, i) => s.varKey
          ? <b key={i} className="font-semibold text-accent-text">{empty(values[s.varKey]) ? `{{${s.varKey}}}` : String(values[s.varKey])}</b>
          : <span key={i}>{s.text === '\\{{' ? '{{' : s.text}</span>)}
      </span>
    </div>
  );
}
