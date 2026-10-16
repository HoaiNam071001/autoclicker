// Inspector – cấu hình bước đang chọn (design v2, cột phải), sửa trực tiếp vào store (bước 15)
import { ReactNode, useEffect, useState } from 'react';
import { Eye, StepIcon, Target } from './icons';
import { openOverlay } from '../windowMode';
import NumInput, { fieldBox } from './NumInput';
import Segmented from './Segmented';
import { Switch } from './RunSettings';
import PinButton from './inputs/PinButton';
import AnchorField from './steps/AnchorField';
import {
  changeStepType, endpoints, pickPointFor, setAnchor, setDragFromPrev, setDragRelative, toggleStep, updateStep, useScript,
} from '../scriptStore';
import { fmtKey, isValidKey, KEY_NAMES, KeyMode, Modifier, MouseButton, Step, STEP_LABEL, StepPoint, StepType } from '../types';

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
const selectBox = `${fieldBox} w-full text-ink`;
const hint = 'text-xs text-muted';

type Props = { step: Step; index: number };

export default function StepForm({ step: s, index }: Props) {
  const { region } = useScript();
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const setPoint = (i: number, p: Partial<StepPoint>) =>
    set({ points: s.points.map((q, j) => (j === i ? { ...q, ...p } : q)) });
  const mouse = s.type === 'click' || s.type === 'double' || s.type === 'hold' || s.type === 'drag';

  return (
    <aside aria-label={`Thuộc tính bước ${index}`} className="flex w-[340px] flex-none flex-col gap-4 overflow-y-auto border-l border-line bg-surface p-4 max-[1100px]:w-[300px]">
      <div className="flex items-center gap-3">
        <span className="grid h-9 w-9 flex-none place-items-center rounded-[9px] bg-accent-soft text-accent-text"><StepIcon type={s.type} className="h-[18px] w-[18px]" /></span>
        <div className="min-w-0 flex-1">
          <div className="truncate text-[15px] font-semibold">{STEP_LABEL[s.type]}</div>
          <div className="text-xs text-muted">Bước {index}{s.enabled ? '' : ' · đang tắt'}</div>
        </div>
        {mouse && (
          <button onClick={() => openOverlay('point', s.id)} aria-label="Xem trên màn hình"
            title="Xem / sửa vị trí riêng bước này trên màn hình – các bước khác mờ đi"
            className="grid h-8 w-8 flex-none place-items-center rounded-control border border-line text-muted hover:bg-hover hover:text-ink"><Eye /></button>
        )}
        <Switch on={s.enabled} onChange={() => toggleStep(s.id)} label="Bật bước" />
      </div>

      <div className="-mt-1 flex gap-[18px] border-b border-line" role="tablist">
        <button role="tab" aria-selected className="relative h-[34px] px-0.5 font-semibold">
          Bước<span className="absolute inset-x-[20%] -bottom-px h-[3px] rounded-sm bg-accent" />
        </button>
        <button role="tab" aria-selected={false} disabled title="Sắp có (bước 21b)" className="h-[34px] px-0.5 text-muted opacity-60">Input</button>
      </div>

      <label className="flex flex-col gap-1.5 text-muted">Hành động
        <select value={s.type} onChange={(e) => changeStepType(s.id, e.target.value as StepType)} className={selectBox}>
          {TYPES.map((t) => <option key={t} value={t}>{STEP_LABEL[t]}</option>)}
        </select>
      </label>

      <div className="grid grid-cols-2 gap-3">
        {mouse && (
          <Field label="Nút chuột" wide group>
            <Segmented value={s.button} onChange={(button: MouseButton) => set({ button })} full label="Nút chuột"
              options={[{ id: 'left', label: 'Trái' }, { id: 'right', label: 'Phải' }, { id: 'center', label: 'Giữa' }]} />
          </Field>
        )}

        {mouse && s.type !== 'drag' && s.points[0] && (
          <div className="col-span-2 flex flex-col gap-1.5 text-muted">
            <span>Vị trí <span className="text-xs">(neo theo góc {region ? 'vùng' : 'màn hình'})</span></span>
            <AnchorField value={s.points[0]} onChange={(p) => setPoint(0, p)} onAnchor={(ax, ay) => setAnchor(s.id, 0, ax, ay)}
              pick={<PickButton step={s} />} trailing={<PinSoon />} />
          </div>
        )}
        {s.type === 'drag' && s.points.length >= 2 && <DragPoints step={s} setPoint={setPoint} />}

        {s.type === 'click' && (
          <>
            <Field pin label="Số click"><NumInput value={s.clicks} min={1} onChange={(v) => set({ clicks: v })} /></Field>
            {s.clicks > 1 && (
              <Field pin label="Giữa các click"><NumInput value={s.clickInterval} unit="ms" onChange={(v) => set({ clickInterval: v })} /></Field>
            )}
          </>
        )}
        {s.type === 'hold' && (
          <Field pin label="Giữ chuột trong"><NumInput value={s.holdMs} unit="ms" onChange={(v) => set({ holdMs: v })} /></Field>
        )}
        {s.type === 'drag' && (
          <Field pin label="Thời gian kéo"><NumInput value={s.dragMs} unit="ms" onChange={(v) => set({ dragMs: v })} /></Field>
        )}

        {s.type === 'type' && (
          <label className="col-span-2 flex flex-col gap-1.5 text-muted">Nội dung gõ
            <textarea value={s.text} onChange={(e) => set({ text: e.target.value })} rows={3}
              className="rounded-control border border-line bg-field px-2.5 py-2 font-mono text-[13px] text-ink outline-none focus:ring-focus" />
            <span className={`text-xs ${s.text ? '' : 'text-danger-text'}`}>
              {s.text ? `${[...s.text].length} ký tự · ` : 'Chưa nhập nội dung – chạy sẽ báo lỗi · '}
              Gõ vào ô đang focus – thường đặt sau một bước click. Tắt bộ gõ tiếng Việt khi chạy.
            </span>
          </label>
        )}

        {s.type === 'key' && (
          <>
            <KeysField step={s} onChange={(keys) => set({ keys })} />
            <fieldset className="col-span-2 flex flex-col gap-1.5">
              <legend className="mb-1.5 text-muted">Kèm phím</legend>
              <div className="flex flex-wrap gap-1.5">
                {MODIFIERS.map((m) => {
                  const on = s.modifiers.includes(m.id);
                  return (
                    <button key={m.id} aria-pressed={on}
                      onClick={() => set({ modifiers: on ? s.modifiers.filter((x) => x !== m.id) : [...s.modifiers, m.id] })}
                      className={`h-8 rounded-control border px-2.5 text-[13px] ${on ? 'border-accent bg-accent-soft font-semibold text-accent-text' : 'border-line hover:bg-hover'}`}>
                      {m.label}
                    </button>
                  );
                })}
              </div>
            </fieldset>
            {s.keys.length > 1 && (
              <fieldset className="col-span-2 flex flex-col gap-1.5">
                <legend className="mb-1.5 text-muted">Cách bấm</legend>
                <div className="flex flex-col gap-1.5">
                  {KEY_MODES.map((m) => (
                    <label key={m.id} className={`flex cursor-pointer flex-col rounded-card px-3 py-2 ${s.keyMode === m.id ? 'ring-selected' : 'border border-line'}`}>
                      <span className="flex items-center gap-2 font-semibold">
                        <input type="radio" name="keyMode" checked={s.keyMode === m.id} onChange={() => set({ keyMode: m.id })} className="accent-accent" />
                        {m.label}
                      </span>
                      <span className={hint}>{m.hint}</span>
                    </label>
                  ))}
                </div>
              </fieldset>
            )}
          </>
        )}

        {mouse && (
          <Field label="Ngẫu nhiên vị trí" wide>
            <NumInput prefix="±" value={s.jitterPx} unit="px" onChange={(v) => set({ jitterPx: v })} />
            <span className={hint}>
              {s.jitterPx > 0 ? `Mỗi lần ${s.type === 'drag' ? 'kéo, A và B' : 'click'} lệch trong ô ${2 * s.jitterPx + 1}×${2 * s.jitterPx + 1} px` : 'Tắt (0)'}
            </span>
          </Field>
        )}
      </div>

      {/* Chờ sau bước + ngẫu nhiên thời gian (bước 19) ở cuối, giống artboard */}
      <section className="mt-auto grid grid-cols-2 gap-3 border-t border-line pt-3">
        <Field pin label={s.type === 'wait' ? 'Thời gian chờ' : 'Chờ sau bước'}>
          <NumInput value={s.delayAfter} unit="ms" onChange={(v) => set({ delayAfter: v })} />
        </Field>
        <Field pin label="Ngẫu nhiên">
          <NumInput prefix="±" value={s.jitterMs} unit="ms" onChange={(v) => set({ jitterMs: v })} />
        </Field>
        <span className={`col-span-2 -mt-1 ${hint}`}>
          {s.jitterMs > 0
            ? `Chờ ${Math.max(s.delayAfter - s.jitterMs, 0)}–${s.delayAfter + s.jitterMs} ms${s.type === 'click' && s.clicks > 1 ? ', giữa các click cũng lệch' : ''}`
            : 'Thời gian chờ cố định (ngẫu nhiên = 0)'}
        </span>
      </section>
    </aside>
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
    <label className="col-span-2 flex flex-col gap-1.5 text-muted">Phím
      <input value={draft} list="key-names" placeholder="enter · hoặc nhiều phím: c, t"
        onChange={(e) => {
          setDraft(e.target.value);
          onChange(parseKeys(e.target.value));
        }}
        onBlur={() => setDraft(s.keys.join(', '))}
        className={`h-8 rounded-control border bg-field px-2.5 font-mono text-[13px] text-ink outline-none ${ok ? 'border-line focus:ring-focus' : 'border-danger'}`} />
      <datalist id="key-names">{KEY_NAMES.map((k) => <option key={k} value={k} />)}</datalist>
      <span className={`text-xs ${ok ? '' : 'text-danger-text'}`}>
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
      className="grid h-8 w-9 flex-none place-items-center rounded-control border border-line text-muted hover:bg-hover hover:text-ink"><Target /></button>
  );
}

// Điểm A / B của bước kéo thả: A = vị trí hoặc "từ điểm trước"; B = vị trí hoặc lệch so với A.
function DragPoints({ step: s, setPoint }: { step: Step; setPoint: (i: number, p: Partial<StepPoint>) => void }) {
  const { steps, region } = useScript();
  const ends = endpoints(steps, region).get(s.id);
  const inRegion = region && <span className="ml-1.5 text-xs">(trong vùng)</span>;
  return (
    <>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span>Điểm A – bắt đầu kéo{!s.dragFromPrev && inRegion}</span>
        <Segmented value={s.dragFromPrev ? 'prev' : 'pos'} onChange={(v) => setDragFromPrev(s.id, v === 'prev')} full
          options={[{ id: 'pos', label: 'Vị trí' }, { id: 'prev', label: 'Từ điểm trước' }]} />
        {s.dragFromPrev ? (
          <span className={`text-[13px] ${ends?.start ? '' : 'text-danger-text'}`}>
            {ends?.start
              ? `Chỗ chuột dừng ở bước ${ends.prevStep} (X ${ends.start.x} Y ${ends.start.y})`
              : 'Phía trước không có bước chuột nào – chạy sẽ báo lỗi'}
          </span>
        ) : (
          <AnchorField value={s.points[0]} onChange={(p) => setPoint(0, p)} onAnchor={(ax, ay) => setAnchor(s.id, 0, ax, ay)}
            pick={<PickButton step={s} title="Chọn lại A rồi B trên màn hình" />} />
        )}
      </div>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span>Điểm B – thả{!s.dragRelative && inRegion}</span>
        <Segmented value={s.dragRelative ? 'rel' : 'pos'} onChange={(v) => setDragRelative(s.id, v === 'rel')} full
          options={[{ id: 'pos', label: 'Vị trí' }, { id: 'rel', label: 'Lệch so với A' }]} />
        {s.dragRelative ? (
          <div className="flex gap-2">
            {(['x', 'y'] as const).map((k) => (
              <NumInput key={k} signed prefix={`Δ${k.toUpperCase()}`} value={s.points[1][k]} onChange={(v) => setPoint(1, { [k]: v })} className="flex-1" />
            ))}
            {s.dragFromPrev && <PickButton step={s} title="Chọn lại B trên màn hình" />}
          </div>
        ) : (
          <AnchorField value={s.points[1]} onChange={(p) => setPoint(1, p)} onAnchor={(ax, ay) => setAnchor(s.id, 1, ax, ay)}
            pick={s.dragFromPrev ? <PickButton step={s} title="Chọn lại B trên màn hình" /> : <span />} />
        )}
        {s.dragRelative && ends?.end && (
          <span className="text-xs">Thả tại X {ends.end.x} Y {ends.end.y}</span>
        )}
      </div>
    </>
  );
}

// pin: thêm nút ghim (đặt làm input của kịch bản) – giao diện có sẵn, logic ở bước 21b
// group: bên trong là nhóm nút (Segmented) – dùng <div>, vì <label> bọc nút thì bấm chữ nhãn sẽ bấm luôn nút đầu
function Field({ label, children, wide = false, pin = false, group = false }: {
  label: string; children: ReactNode; wide?: boolean; pin?: boolean; group?: boolean;
}) {
  const Tag = group ? 'div' : 'label';
  return (
    <Tag className={`flex flex-col gap-1.5 text-muted ${wide ? 'col-span-2' : ''}`}>{label}
      {pin ? <span className="flex gap-1.5 [&>*:first-child]:min-w-0 [&>*:first-child]:flex-1">{children}<PinSoon /></span> : children}
    </Tag>
  );
}

const PinSoon = () => <PinButton pinned={false} disabled title="Đặt làm input – sắp có (bước 21b)" />;
