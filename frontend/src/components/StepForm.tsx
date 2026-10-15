// Form cấu hình bước đang chọn – sửa trực tiếp vào store (bước 15)
import { ReactNode, useEffect, useState } from 'react';
import { Eye, Target } from './icons';
import { openOverlay } from '../windowMode';
import NumInput from './NumInput';
import {
  changeStepType, endpoints, pickPointFor, setDragFromPrev, setDragRelative, toggleStep, updateStep, useScript,
} from '../scriptStore';
import { fmtKey, isValidKey, KEY_NAMES, KeyMode, Modifier, MouseButton, Point, Step, STEP_LABEL, StepType } from '../types';

const TYPES: StepType[] = ['click', 'double', 'hold', 'drag', 'wait', 'type', 'key'];
const isMac = navigator.platform.startsWith('Mac');
const MODIFIERS: { id: Modifier; label: string }[] = [
  { id: 'cmd', label: isMac ? '⌘ Cmd' : '⊞ Win' },
  { id: 'ctrl', label: 'Ctrl' },
  { id: 'alt', label: isMac ? '⌥ Option' : 'Alt' },
  { id: 'shift', label: '⇧ Shift' },
];
const KEY_MODES: { id: KeyMode; label: string; hint: string }[] = [
  { id: 'chord', label: 'Nhấn giữ cùng lúc', hint: 'Đè phím kèm và mọi phím, rồi nhả – vd Ctrl+C+T' },
  { id: 'sequence', label: 'Giữ phím kèm, gõ lần lượt', hint: 'Đè phím kèm, gõ C rồi T, rồi nhả phím kèm' },
];

type Props = { step: Step; index: number; children?: ReactNode };

export default function StepForm({ step: s, index, children }: Props) {
  const { region } = useScript();
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const setPoint = (i: number, p: Partial<Point>) =>
    set({ points: s.points.map((q, j) => (j === i ? { ...q, ...p } : q)) });
  const mouse = s.type === 'click' || s.type === 'double' || s.type === 'hold' || s.type === 'drag';

  return (
    <section aria-label={`Cấu hình bước ${index}`} className="flex min-h-0 min-w-0 flex-[999_1_560px] flex-col gap-5 overflow-y-auto rounded-[10px] border border-line bg-white px-6 py-5">
      <div className="flex items-center gap-3">
        <span className="grid h-8 w-8 place-items-center rounded-full bg-accent font-bold text-white">{index}</span>
        <h2 className="text-lg font-semibold">Bước {index}</h2>
        {mouse && (
          <button onClick={() => openOverlay('point', s.id)} title="Xem / sửa vị trí riêng bước này trên màn hình – các bước khác mờ đi"
            className="ml-auto flex h-9 items-center gap-2 rounded-lg border border-line px-3 text-[13px] hover:bg-ground">
            <Eye />Xem trên màn hình
          </button>
        )}
        <label className={`${mouse ? 'ml-3' : 'ml-auto'} flex items-center gap-2 text-[#3F4857]`}>
          <input type="checkbox" checked={s.enabled} onChange={() => toggleStep(s.id)} className="h-[18px] w-[18px] accent-accent" />Bật
        </label>
      </div>

      <fieldset className="flex flex-col gap-2">
        <legend className="mb-2 font-semibold">Hành động</legend>
        <div className="flex flex-wrap gap-1.5">
          {TYPES.map((t) => (
            <button key={t} onClick={() => changeStepType(s.id, t)} aria-pressed={t === s.type}
              className={`h-9 rounded-lg px-3 text-[13px] ${t === s.type ? 'border-2 border-accent bg-accent-soft font-semibold text-accent-strong' : 'border border-line bg-white hover:bg-ground'}`}>
              {STEP_LABEL[t]}
            </button>
          ))}
        </div>
      </fieldset>

      <div className="grid grid-cols-2 gap-x-5 gap-y-4">
        {mouse && (
          <label className="flex flex-col gap-1.5 font-semibold">Nút chuột
            <select value={s.button} onChange={(e) => set({ button: e.target.value as MouseButton })}
              className="h-10 rounded-lg border border-line bg-white px-2.5 font-normal">
              <option value="left">Trái</option><option value="right">Phải</option><option value="center">Giữa</option>
            </select>
          </label>
        )}

        {mouse && s.type !== 'drag' && s.points[0] && (
          <div className="flex flex-col gap-1.5">
            <span className="font-semibold">
              Vị trí{region && <span className="ml-1.5 text-xs font-normal text-muted">(trong vùng)</span>}
            </span>
            <div className="flex gap-2">
              {(['x', 'y'] as const).map((k) => (
                <NumInput key={k} prefix={k.toUpperCase()} value={s.points[0][k]} onChange={(v) => setPoint(0, { [k]: v })} className="flex-1" />
              ))}
              <PickButton step={s} />
            </div>
          </div>
        )}
        {s.type === 'drag' && s.points.length >= 2 && <DragPoints step={s} setPoint={setPoint} />}

        {s.type === 'click' && (
          <>
            <Field label="Số click"><NumInput value={s.clicks} min={1} onChange={(v) => set({ clicks: v })} /></Field>
            {s.clicks > 1 && (
              <Field label="Giữa các click"><NumInput value={s.clickInterval} unit="ms" onChange={(v) => set({ clickInterval: v })} /></Field>
            )}
          </>
        )}
        {s.type === 'hold' && (
          <Field label="Giữ chuột trong"><NumInput value={s.holdMs} unit="ms" onChange={(v) => set({ holdMs: v })} /></Field>
        )}
        {s.type === 'drag' && (
          <Field label="Thời gian kéo"><NumInput value={s.dragMs} unit="ms" onChange={(v) => set({ dragMs: v })} /></Field>
        )}

        {s.type === 'type' && (
          <label className="col-span-2 flex flex-col gap-1.5 font-semibold">Nội dung gõ
            <textarea value={s.text} onChange={(e) => set({ text: e.target.value })} rows={3}
              className="rounded-lg border border-line px-2.5 py-2 font-mono font-normal focus:border-accent focus:outline-none" />
            <span className={`text-xs font-normal ${s.text ? 'text-muted' : 'text-danger'}`}>
              {s.text ? `${[...s.text].length} ký tự · ` : 'Chưa nhập nội dung – chạy sẽ báo lỗi · '}
              Gõ vào ô đang focus – thường đặt sau một bước click. Tắt bộ gõ tiếng Việt khi chạy.
            </span>
          </label>
        )}

        {s.type === 'key' && (
          <>
            <KeysField step={s} onChange={(keys) => set({ keys })} />
            <fieldset className="flex flex-col gap-1.5">
              <legend className="mb-1.5 font-semibold">Kèm phím</legend>
              <div className="flex flex-wrap gap-1.5">
                {MODIFIERS.map((m) => (
                  <label key={m.id} className="flex h-10 items-center gap-1.5 rounded-lg border border-line px-3">
                    <input type="checkbox" checked={s.modifiers.includes(m.id)} className="accent-accent"
                      onChange={(e) => set({
                        modifiers: e.target.checked ? [...s.modifiers, m.id] : s.modifiers.filter((x) => x !== m.id),
                      })} />
                    {m.label}
                  </label>
                ))}
              </div>
            </fieldset>
            {s.keys.length > 1 && (
              <fieldset className="col-span-2 flex flex-col gap-1.5">
                <legend className="mb-1.5 font-semibold">Cách bấm</legend>
                <div className="flex flex-wrap gap-1.5">
                  {KEY_MODES.map((m) => (
                    <label key={m.id} className={`flex cursor-pointer flex-col rounded-lg border px-3 py-2 ${s.keyMode === m.id ? 'border-2 border-accent bg-accent-soft' : 'border-line'}`}>
                      <span className="flex items-center gap-2 font-semibold">
                        <input type="radio" name="keyMode" checked={s.keyMode === m.id} onChange={() => set({ keyMode: m.id })} className="accent-accent" />
                        {m.label}
                      </span>
                      <span className="text-xs text-muted">{m.hint}</span>
                    </label>
                  ))}
                </div>
              </fieldset>
            )}
          </>
        )}

        <Field label={s.type === 'wait' ? 'Thời gian chờ' : 'Chờ trước bước tiếp theo'}>
          <NumInput value={s.delayAfter} unit="ms" onChange={(v) => set({ delayAfter: v })} />
        </Field>

        {/* Ngẫu nhiên hoá (bước 19) – thao tác giống người thật hơn */}
        <Field label="Ngẫu nhiên thời gian">
          <NumInput prefix="±" value={s.jitterMs} unit="ms" onChange={(v) => set({ jitterMs: v })} />
          <span className="text-xs font-normal text-muted">
            {s.jitterMs > 0
              ? `Chờ ${Math.max(s.delayAfter - s.jitterMs, 0)}–${s.delayAfter + s.jitterMs} ms${s.type === 'click' && s.clicks > 1 ? ', giữa các click cũng lệch' : ''}`
              : 'Tắt (0)'}
          </span>
        </Field>
        {mouse && (
          <Field label="Ngẫu nhiên vị trí">
            <NumInput prefix="±" value={s.jitterPx} unit="px" onChange={(v) => set({ jitterPx: v })} />
            <span className="text-xs font-normal text-muted">
              {s.jitterPx > 0 ? `Mỗi lần ${s.type === 'drag' ? 'kéo, A và B' : 'click'} lệch trong ô ${2 * s.jitterPx + 1}×${2 * s.jitterPx + 1} px` : 'Tắt (0)'}
            </span>
          </Field>
        )}
      </div>

      {children}
    </section>
  );
}

// Ô phím chính: nhiều phím cách nhau bằng dấu phẩy / khoảng trắng ("c, t"). Giữ bản nháp khi gõ.
const parseKeys = (text: string) => text.split(/[\s,]+/).map((k) => k.toLowerCase()).filter(Boolean);

function KeysField({ step: s, onChange }: { step: Step; onChange: (keys: string[]) => void }) {
  const [draft, setDraft] = useState(s.keys.join(', '));
  const joined = s.keys.join(',');
  // Chỉ ghi đè bản nháp khi phím bị đổi từ nơi khác (đổi bước, hoàn tác…), không phải do đang gõ.
  useEffect(() => {
    setDraft((d) => (parseKeys(d).join(',') === joined ? d : s.keys.join(', ')));
  }, [joined]); // eslint-disable-line react-hooks/exhaustive-deps
  const bad = s.keys.filter((k) => !isValidKey(k));
  const ok = s.keys.length > 0 && bad.length === 0;
  return (
    <label className="flex flex-col gap-1.5 font-semibold">Phím
      <input value={draft} list="key-names" placeholder="enter · hoặc nhiều phím: c, t"
        onChange={(e) => {
          setDraft(e.target.value);
          onChange(parseKeys(e.target.value));
        }}
        onBlur={() => setDraft(s.keys.join(', '))}
        className={`h-10 rounded-lg border px-2.5 font-mono font-normal focus:outline-none ${ok ? 'border-line focus:border-accent' : 'border-danger'}`} />
      <datalist id="key-names">{KEY_NAMES.map((k) => <option key={k} value={k} />)}</datalist>
      <span className={`text-xs font-normal ${ok ? 'text-muted' : 'text-danger'}`}>
        {s.keys.length === 0 ? 'Chưa có phím – chạy sẽ báo lỗi'
          : bad.length ? `Tên phím không hợp lệ: ${bad.join(', ')} – chọn trong gợi ý hoặc gõ 1 ký tự`
          : `Sẽ bấm: ${fmtKey(s)}`}
      </span>
    </label>
  );
}

function PickButton({ step: s, title = 'Chọn lại trên màn hình' }: { step: Step; title?: string }) {
  return (
    <button aria-label={title} title={title} onClick={() => pickPointFor(s.id)}
      className="grid h-10 w-11 flex-none place-items-center rounded-lg border border-line hover:bg-ground"><Target /></button>
  );
}

function Segmented<T extends string>({ value, options, onChange }: {
  value: T; options: { id: T; label: string }[]; onChange: (v: T) => void;
}) {
  return (
    <div className="flex w-fit overflow-hidden rounded-lg border border-line text-[13px] font-normal" role="radiogroup">
      {options.map((o, i) => (
        <button key={o.id} role="radio" aria-checked={value === o.id} onClick={() => onChange(o.id)}
          className={`h-8 px-3 ${i > 0 ? 'border-l border-line' : ''} ${value === o.id ? 'bg-ink text-white' : 'bg-white hover:bg-ground'}`}>
          {o.label}
        </button>
      ))}
    </div>
  );
}

// Điểm A / B của bước kéo thả: A = vị trí hoặc "từ điểm trước"; B = vị trí hoặc lệch so với A.
function DragPoints({ step: s, setPoint }: { step: Step; setPoint: (i: number, p: Partial<Point>) => void }) {
  const { steps, region } = useScript();
  const ends = endpoints(steps, region).get(s.id);
  const inRegion = region && <span className="ml-1.5 text-xs font-normal text-muted">(trong vùng)</span>;
  return (
    <>
      <div className="col-start-1 flex flex-col gap-1.5">
        <span className="font-semibold">Điểm A – bắt đầu kéo{!s.dragFromPrev && inRegion}</span>
        <Segmented value={s.dragFromPrev ? 'prev' : 'pos'} onChange={(v) => setDragFromPrev(s.id, v === 'prev')}
          options={[{ id: 'pos', label: 'Vị trí' }, { id: 'prev', label: 'Từ điểm trước' }]} />
        {s.dragFromPrev ? (
          <span className={`flex h-10 items-center text-[13px] ${ends?.start ? 'text-muted' : 'text-danger'}`}>
            {ends?.start
              ? `Chỗ chuột dừng ở bước ${ends.prevStep} (X ${ends.start.x} Y ${ends.start.y})`
              : 'Phía trước không có bước chuột nào – chạy sẽ báo lỗi'}
          </span>
        ) : (
          <div className="flex gap-2">
            {(['x', 'y'] as const).map((k) => (
              <NumInput key={k} prefix={k.toUpperCase()} value={s.points[0][k]} onChange={(v) => setPoint(0, { [k]: v })} className="flex-1" />
            ))}
            <PickButton step={s} title="Chọn lại A rồi B trên màn hình" />
          </div>
        )}
      </div>
      <div className="flex flex-col gap-1.5">
        <span className="font-semibold">Điểm B – thả{!s.dragRelative && inRegion}</span>
        <Segmented value={s.dragRelative ? 'rel' : 'pos'} onChange={(v) => setDragRelative(s.id, v === 'rel')}
          options={[{ id: 'pos', label: 'Vị trí' }, { id: 'rel', label: 'Lệch so với A' }]} />
        <div className="flex gap-2">
          {(['x', 'y'] as const).map((k) => (
            <NumInput key={k} signed={s.dragRelative} prefix={s.dragRelative ? `Δ${k.toUpperCase()}` : k.toUpperCase()}
              value={s.points[1][k]} onChange={(v) => setPoint(1, { [k]: v })} className="flex-1" />
          ))}
          {s.dragFromPrev && <PickButton step={s} title="Chọn lại B trên màn hình" />}
        </div>
        {s.dragRelative && ends?.end && (
          <span className="text-xs font-normal text-muted">Thả tại X {ends.end.x} Y {ends.end.y}</span>
        )}
      </div>
    </>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return <label className="flex flex-col gap-1.5 font-semibold">{label}{children}</label>;
}
