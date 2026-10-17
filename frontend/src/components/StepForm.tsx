// Inspector – cấu hình bước đang chọn (design v2, cột phải), sửa trực tiếp vào store (bước 15)
import { ReactNode, useEffect, useRef, useState } from 'react';
import { Eye, StepIcon, Target } from './icons';
import { openOverlay } from '../windowMode';
import NumInput, { fieldBox } from './NumInput';
import Segmented from './Segmented';
import { Switch } from './RunSettings';
import PinButton, { KeyChip, PinMenu } from './inputs/PinButton';
import VarTextInput from './inputs/VarTextInput';
import InputsPanel from './inputs/InputsPanel';
import AnchorField, { AnchorMap } from './steps/AnchorField';
import ScrollMap from './steps/ScrollMap';
import {
  boundsOf, changeStepType, endpoints, pickPointFor, setAnchor, setDragFromPrev, setDragRelative, setScrollAt, toggleStep, toScreen, updateStep, useScript,
} from '../scriptStore';
import {
  addInput, compatibleInputs, inputViews, moveInput, pinField, pinnedInput, pinToInput, removeInput, setPinnedValue, unpinField, updateInput,
} from '../scriptInputs';
import { loadMemory, stopAutoRun, useRunForm } from '../runForm';
import {
  fmtKey, isValidKey, KEY_NAMES, keyLabel, KeyMode, Modifier, MouseButton, PinField, SCROLL_NOTCH_PX, ScrollDirection, Step, STEP_LABEL, stepLabel,
  StepPoint, StepType, CAN_FAIL,
} from '../types';
import { insideBounds } from '../anchor';
import { checkCapture, pickColor, requestCapture, startCapture, templateUrl, testFind, useImageState } from '../imageSteps';

const TYPES: StepType[] = ['click', 'double', 'hold', 'drag', 'scroll', 'wait', 'type', 'key', 'open', 'waitImage', 'clickImage', 'waitPixel'];
// Trình duyệt gợi ý cho bước Mở (Windows: open.go đổi sang tên lệnh chrome / msedge / firefox)
const BROWSERS = ['Google Chrome', 'Microsoft Edge', 'Firefox', ...(navigator.platform.startsWith('Mac') ? ['Safari'] : [])];
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
  const script = useScript();
  const { region } = script;
  const [tab, setTab] = useState<'step' | 'inputs'>('step');
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const setPoint = (i: number, p: Partial<StepPoint>) =>
    set({ points: s.points.map((q, j) => (j === i ? { ...q, ...p } : q)) });
  const mouse = s.type === 'click' || s.type === 'double' || s.type === 'hold' || s.type === 'drag';
  const atPoint = mouse || s.type === 'waitPixel' || (s.type === 'scroll' && s.points.length > 0); // có vị trí trên màn hình

  return (
    <aside aria-label={`Thuộc tính bước ${index}`} className="flex w-[340px] flex-none flex-col gap-4 overflow-y-auto border-l border-line bg-surface p-4 max-[1100px]:w-[300px]">
      <div className="flex items-center gap-3">
        <span className="grid h-9 w-9 flex-none place-items-center rounded-[9px] bg-accent-soft text-accent-text"><StepIcon type={s.type} className="h-[18px] w-[18px]" /></span>
        <div className="min-w-0 flex-1">
          <div className="truncate text-[15px] font-semibold">{stepLabel(s)}</div>
          <div className="text-xs text-muted">Bước {index}{s.enabled ? '' : ' · đang tắt'}</div>
        </div>
        {atPoint && (
          <button onClick={() => openOverlay('point', s.id)} aria-label="Xem trên màn hình"
            title="Xem / sửa vị trí riêng bước này trên màn hình – các bước khác mờ đi"
            className="grid h-8 w-8 flex-none place-items-center rounded-control border border-line text-muted hover:bg-hover hover:text-ink"><Eye /></button>
        )}
        <Switch on={s.enabled} onChange={() => toggleStep(s.id)} label="Bật bước" />
      </div>

      <div className="-mt-1 flex gap-[18px] border-b border-line" role="tablist">
        <TabBtn active={tab === 'step'} onClick={() => setTab('step')}>Bước</TabBtn>
        <TabBtn active={tab === 'inputs'} onClick={() => setTab('inputs')} title="Các giá trị hỏi khi bấm Chạy">
          Input {!!script.inputs?.length && <span className="rounded-full bg-accent-soft px-1.5 text-[11px] font-semibold text-accent-text">{script.inputs.length}</span>}
        </TabBtn>
      </div>

      {tab === 'inputs' ? <InputsTab /> : <>

      <label className="flex flex-col gap-1.5 text-muted">Hành động
        <select value={s.type} onChange={(e) => changeStepType(s.id, e.target.value as StepType)} className={selectBox}>
          {TYPES.map((t) => <option key={t} value={t}>{t === s.type ? stepLabel(s) : STEP_LABEL[t]}</option>)}
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
              pick={<PickButton step={s} />} trailing={<PinSoon />} bounds={boundsOf(region)} boundsLabel={region ? 'Vùng' : 'Màn hình'} />
          </div>
        )}
        {s.type === 'drag' && s.points.length >= 2 && <DragPoints step={s} setPoint={setPoint} />}

        {s.type === 'click' && (
          <>
            <NumField label="Số click" step={s} field="clicks" min={1} />
            {s.clicks > 1 && <NumField label="Giữa các click" step={s} field="clickInterval" unit="ms" />}
          </>
        )}
        {s.type === 'hold' && <NumField label="Giữ chuột trong" step={s} field="holdMs" unit="ms" />}
        {s.type === 'drag' && <NumField label="Thời gian kéo" step={s} field="dragMs" unit="ms" />}

        {s.type === 'scroll' && (
          <>
            <Field label="Hướng" wide group>
              <Segmented value={s.direction ?? 'down'} onChange={(direction: ScrollDirection) => set({ direction })} full label="Hướng cuộn"
                options={[{ id: 'up', label: '↑ Lên' }, { id: 'down', label: '↓ Xuống' }, { id: 'left', label: '← Trái' }, { id: 'right', label: '→ Phải' }]} />
            </Field>
            <NumField label="Số nấc" step={s} field="amount" unit="nấc" min={1} />
            <Field label="Quãng cuộn ≈">
              <NumInput value={(s.amount ?? 0) * SCROLL_NOTCH_PX} unit="px" min={SCROLL_NOTCH_PX} pinned={!!pinnedInput(script, s.id, 'amount')}
                onChange={(px) => {
                  const amount = Math.max(1, Math.round(px / SCROLL_NOTCH_PX));
                  const inp = pinnedInput(script, s.id, 'amount');
                  if (inp) setPinnedValue(s.id, 'amount', amount);
                  else set({ amount });
                }} />
              <span className={hint}>Gõ px, tự đổi ra nấc ({SCROLL_NOTCH_PX} px / nấc). Cuộn dài hơn màn hình cũng được.</span>
            </Field>
            <Field label="Cuộn tại" wide group>
              <Segmented value={s.points.length ? 'point' : 'mouse'} onChange={(v) => setScrollAt(s.id, v === 'point')} full label="Cuộn tại"
                options={[{ id: 'mouse', label: 'Chuột' }, { id: 'point', label: 'Điểm' }]} />
            </Field>
            {s.points[0] ? (
              <div className="col-span-2 flex flex-col gap-1.5 text-muted">
                <span>Đặt chuột tại <span className="text-xs">(phải nằm trên {region ? 'vùng' : 'màn hình'} – không phải chỗ cuộn tới)</span></span>
                <AnchorField value={s.points[0]} onChange={(p) => setPoint(0, p)} onAnchor={(ax, ay) => setAnchor(s.id, 0, ax, ay)}
                  pick={<PickButton step={s} />} />
                {!insideBounds(toScreen(s.points[0], region), boundsOf(region)) && (
                  <span className="text-xs text-danger-text">
                    Điểm nằm ngoài {region ? 'vùng' : 'màn hình'} ({boundsOf(region).w} × {boundsOf(region).h}) – chuột không tới được, chạy sẽ báo lỗi.
                    Muốn cuộn xa hơn thì tăng số nấc / quãng cuộn.
                  </span>
                )}
              </div>
            ) : (
              <span className={`col-span-2 ${hint}`}>Cuộn ở chỗ chuột đang đứng – thường sau một bước click vào trang / danh sách cần cuộn.</span>
            )}
            <div className="col-span-2 flex flex-col gap-1.5 text-muted">
              <span>Cuộn tới đâu</span>
              <ScrollMap bounds={boundsOf(region)} label={region ? 'Vùng' : 'Màn hình'} direction={s.direction ?? 'down'} amount={s.amount ?? 0}
                at={endpoints(script.steps, region).get(s.id)?.start ?? null} />
            </div>
          </>
        )}

        {(s.type === 'waitImage' || s.type === 'clickImage' || s.type === 'waitPixel') && <CapturePermission />}
        {s.type === 'waitImage' && <ImageSearchFields step={s} />}
        {s.type === 'clickImage' && <ClickImageFields step={s} />}
        {s.type === 'waitPixel' && <WaitPixelFields step={s} setPoint={setPoint} />}
        {CAN_FAIL.includes(s.type) && <RetryFields step={s} />}

        {s.type === 'open' && (
          <>
            <div className="col-span-2 flex flex-col gap-1.5 text-muted">
              <span className="flex items-center justify-between">URL, đường dẫn hoặc tên app
                <span className="text-xs">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
              </span>
              <VarTextInput value={s.target ?? ''} onChange={(target) => set({ target })} inputs={script.inputs ?? []} rows={2}
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
        )}

        {s.type === 'type' && (
          <div className="col-span-2 flex flex-col gap-1.5 text-muted">
            <span className="flex items-center justify-between">Nội dung gõ
              <span className="text-xs">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
            </span>
            <VarTextInput value={s.text} onChange={(text) => set({ text })} inputs={script.inputs ?? []}
              onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
            <span className={`text-xs ${s.text ? '' : 'text-danger-text'}`}>
              {s.text ? `${[...s.text].length} ký tự · ` : 'Chưa nhập nội dung – chạy sẽ báo lỗi · '}
              Gõ vào ô đang focus – thường đặt sau một bước click. Tắt bộ gõ tiếng Việt khi chạy. <span className="font-mono">{'\\{{'}</span> để gõ ra <span className="font-mono">{'{{'}</span>.
            </span>
          </div>
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
        <NumField label={s.type === 'wait' ? 'Thời gian chờ' : 'Chờ sau bước'} step={s} field="delayAfter" unit="ms" />
        <NumField label="Ngẫu nhiên" step={s} field="jitterMs" unit="ms" prefix="±" />
        <span className={`col-span-2 -mt-1 ${hint}`}>
          {s.jitterMs > 0
            ? `Chờ ${Math.max(s.delayAfter - s.jitterMs, 0)}–${s.delayAfter + s.jitterMs} ms${s.type === 'click' && s.clicks > 1 ? ', giữa các click cũng lệch' : ''}`
            : 'Thời gian chờ cố định (ngẫu nhiên = 0)'}
        </span>
      </section>
      </>}
    </aside>
  );
}

// Ô phím chính: nhiều phím cách nhau bằng dấu phẩy / khoảng trắng ("c, t"). Giữ bản nháp khi gõ.
// Gợi ý tên phím tự vẽ (datalist của WebView macOS chữ trắng trên nền trắng, không có mô tả): lọc theo phím đang gõ
// (phần sau dấu phẩy cuối) trong cả tên lẫn mô tả; ↑ ↓ chọn, Enter / Tab điền, Esc đóng.
const parseKeys = (text: string) => text.split(/[\s,]+/).map((k) => k.toLowerCase()).filter(Boolean);
const KEY_OPTIONS = KEY_NAMES.map((name) => ({ name, label: keyLabel(name) }));

function KeysField({ step: s, onChange }: { step: Step; onChange: (keys: string[]) => void }) {
  const [draft, setDraft] = useState(s.keys.join(', '));
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState(0);
  const listRef = useRef<HTMLUListElement>(null);
  const joined = s.keys.join(',');
  // Chỉ ghi đè bản nháp khi phím bị đổi từ nơi khác (đổi bước, hoàn tác…), không phải do đang gõ.
  useEffect(() => {
    setDraft((d) => (parseKeys(d).join(',') === joined ? d : s.keys.join(', ')));
  }, [joined]); // eslint-disable-line react-hooks/exhaustive-deps
  const bad = s.keys.filter((k) => !isValidKey(k));
  const ok = s.keys.length > 0 && bad.length === 0;

  // Phím đang gõ = phần sau dấu phân cách cuối.
  const start = draft.search(/[^\s,]*$/);
  const typed = draft.slice(start).toLowerCase();
  const matches = typed
    ? KEY_OPTIONS.filter((o) => o.name.includes(typed) || o.label.toLowerCase().includes(typed))
      .sort((x, y) => Number(!x.name.startsWith(typed)) - Number(!y.name.startsWith(typed)))
    : KEY_OPTIONS;
  const shown = open && matches.length > 0;

  useEffect(() => {
    listRef.current?.children[active]?.scrollIntoView({ block: 'nearest' });
  }, [active, shown]);

  const edit = (text: string) => {
    setDraft(text);
    onChange(parseKeys(text));
  };
  const choose = (name: string) => {
    edit(draft.slice(0, start) + name);
    setOpen(false);
  };
  const onKeyDown = (e: React.KeyboardEvent) => {
    if (!shown) {
      if (e.key === 'ArrowDown') setOpen(true);
      return;
    }
    if (e.key === 'ArrowDown') setActive((active + 1) % matches.length);
    else if (e.key === 'ArrowUp') setActive((active - 1 + matches.length) % matches.length);
    else if ((e.key === 'Enter' || e.key === 'Tab') && typed) choose(matches[active].name);
    else if (e.key === 'Escape') setOpen(false);
    else return;
    e.preventDefault();
    e.stopPropagation();
  };

  return (
    <div className="col-span-2 flex flex-col gap-1.5 text-muted">
      <label htmlFor="key-input">Phím</label>
      <div className="relative">
        <input id="key-input" value={draft} placeholder="enter · hoặc nhiều phím: c, t" autoComplete="off" spellCheck={false}
          role="combobox" aria-expanded={shown} aria-controls="key-options" aria-autocomplete="list"
          onChange={(e) => {
            edit(e.target.value);
            setActive(0);
            setOpen(true);
          }}
          onFocus={() => setOpen(true)} onClick={() => setOpen(true)} onKeyDown={onKeyDown}
          onBlur={() => {
            setOpen(false);
            setDraft(s.keys.join(', '));
          }}
          className={`h-8 w-full rounded-control border bg-field px-2.5 font-mono text-[13px] text-ink outline-none ${ok ? 'border-line focus:ring-focus' : 'border-danger'}`} />
        {shown && (
          <ul id="key-options" ref={listRef} role="listbox" aria-label="Gợi ý tên phím"
            className="absolute inset-x-0 top-full z-30 mt-1 max-h-[240px] overflow-y-auto rounded-panel border border-line bg-surface p-1 shadow-[0_8px_28px_rgb(0_0_0/.16)]">
            {matches.map((o, i) => (
              <li key={o.name} role="option" aria-selected={i === active}
                onMouseDown={(e) => {
                  e.preventDefault(); // giữ focus ô nhập
                  choose(o.name);
                }}
                onMouseEnter={() => setActive(i)}
                className={`flex cursor-pointer items-center gap-3 rounded-md px-2.5 py-1.5 ${i === active ? 'bg-hover' : ''}`}>
                <span className="min-w-[88px] font-mono text-[13px] text-ink">{o.name}</span>
                <span className="min-w-0 flex-1 truncate text-xs text-muted">{o.label}</span>
              </li>
            ))}
          </ul>
        )}
      </div>
      <span className={`text-xs ${ok ? '' : 'text-danger-text'}`}>
        {s.keys.length === 0 ? 'Chưa có phím – chạy sẽ báo lỗi'
          : bad.length ? `Tên phím không hợp lệ: ${bad.join(', ')} – chọn trong gợi ý hoặc gõ 1 ký tự`
          : `Sẽ bấm: ${fmtKey(s)} · phím chữ / số: gõ 1 ký tự (a, 1, /)`}
      </span>
    </div>
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
  // Bản đồ chung A + B: chọn điểm nào để đổi neo (chỉ điểm lưu vị trí – không phải "từ điểm trước" / "lệch so với A").
  const editable = [!s.dragFromPrev && 'A', !s.dragRelative && 'B'].filter(Boolean) as ('A' | 'B')[];
  const [pick, setPick] = useState<'A' | 'B'>('A');
  const cur = editable.includes(pick) ? pick : editable[0];
  const curIdx = cur === 'B' ? 1 : 0;
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
            pick={<PickButton step={s} title="Chọn lại A rồi B trên màn hình" />} grid={false} />
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
            pick={s.dragFromPrev ? <PickButton step={s} title="Chọn lại B trên màn hình" /> : <span />} grid={false} />
        )}
        {s.dragRelative && ends?.end && (
          <span className="text-xs">Thả tại X {ends.end.x} Y {ends.end.y}</span>
        )}
      </div>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span className="flex items-center gap-2">
          <span className="flex-1">Neo trên {region ? 'vùng' : 'màn hình'}</span>
          {editable.length > 1 && (
            <Segmented value={cur} onChange={(v: 'A' | 'B') => setPick(v)} label="Đổi neo cho điểm"
              options={[{ id: 'A', label: 'A' }, { id: 'B', label: 'B' }]} />
          )}
        </span>
        <AnchorMap bounds={boundsOf(region)} label={region ? 'Vùng' : 'Màn hình'}
          value={cur ? s.points[curIdx] : undefined} onAnchor={cur ? (ax, ay) => setAnchor(s.id, curIdx, ax, ay) : undefined}
          marks={[
            ...(ends?.start ? [{ at: ends.start, text: 'A', active: cur === 'A' }] : []),
            ...(ends?.end ? [{ at: ends.end, text: 'B', active: cur === 'B' }] : []),
          ]}
          arrow={ends?.start && ends.end ? [ends.start, ends.end] : null} />
        {!cur && <span className="text-xs">A lấy từ điểm trước, B lệch so với A – không có điểm nào để neo</span>}
      </div>
    </>
  );
}

// group: bên trong là nhóm nút (Segmented) – dùng <div>, vì <label> bọc nút thì bấm chữ nhãn sẽ bấm luôn nút đầu
function Field({ label, children, wide = false, group = false }: {
  label: string; children: ReactNode; wide?: boolean; group?: boolean;
}) {
  const Tag = group ? 'div' : 'label';
  return <Tag className={`flex flex-col gap-1.5 text-muted ${wide ? 'col-span-2' : ''}`}>{label}{children}</Tag>;
}

// Ô số ghim được (bước 21b): đã ghim → viền nét đứt + chip key cạnh nhãn, sửa số = sửa mặc định của input
// (mọi trường ghim cùng input đổi theo).
function NumField({ label, step: s, field, unit, prefix, min = 0 }: {
  label: string; step: Step; field: Exclude<PinField, 'enabled'>; unit?: string; prefix?: string; min?: number;
}) {
  const script = useScript();
  const inp = pinnedInput(script, s.id, field);
  const others = compatibleInputs(script, field).filter((i) => i !== inp);
  return (
    <div className="flex flex-col gap-1.5 text-muted">
      <span className="flex min-w-0 items-center gap-1.5"><span className="truncate">{label}</span>{inp && <KeyChip k={inp.key} />}</span>
      <span className="flex gap-1.5">
        <NumInput value={s[field] ?? min} min={min} unit={unit} prefix={prefix} pinned={!!inp} className="min-w-0 flex-1"
          onChange={(v) => (inp ? setPinnedValue(s.id, field, v) : updateStep(s.id, { [field]: v }))} />
        <PinMenu pinnedKey={inp?.key} options={others} onPinNew={() => pinField(s.id, field)}
          onPinTo={(key) => pinToInput(s.id, field, key)} onUnpin={() => unpinField(s.id, field)} />
      </span>
    </div>
  );
}

const PinSoon = () => <PinButton pinned={false} disabled title="Ghim toạ độ làm input – sắp có" />;

function TabBtn({ active, onClick, title, children }: { active: boolean; onClick: () => void; title?: string; children: ReactNode }) {
  return (
    <button role="tab" aria-selected={active} onClick={onClick} title={title}
      className={`relative flex h-[34px] items-center gap-1.5 px-0.5 ${active ? 'font-semibold' : 'text-muted hover:text-ink'}`}>
      {children}
      {active && <span className="absolute inset-x-[15%] -bottom-px h-[3px] rounded-sm bg-accent" />}
    </button>
  );
}

// Tab Input: input của cả kịch bản (không riêng bước đang chọn).
function InputsTab() {
  const script = useScript();
  const { memory } = useRunForm();
  useEffect(() => {
    loadMemory();
  }, [script.name]);
  return (
    <InputsPanel inputs={inputViews(script)} onUpdate={updateInput} onRemove={removeInput} onMove={moveInput}
      onAdd={() => addInput({ label: 'Input mới' })}
      footer={memory.autoRun && !!script.inputs?.length && (
        <div className="flex items-center gap-2 rounded-lg bg-soft px-3 py-2 text-xs text-muted">
          <span className="flex-1">Bấm Chạy đang chạy luôn với giá trị lần trước, không hỏi (Shift + bấm Chạy để sửa)</span>
          <button onClick={stopAutoRun} className="font-semibold text-accent-text">Hỏi lại</button>
        </div>
      )} />
  );
}

// ---- Chờ ảnh / click theo ảnh / chờ pixel (bước 25) ----

// macOS: chưa có quyền ghi màn hình thì chụp chỉ ra hình nền – báo ngay trong inspector.
function CapturePermission() {
  const { allowed } = useImageState();
  useEffect(() => {
    checkCapture();
  }, []);
  if (allowed !== false) return null;
  return (
    <div className="col-span-2 flex flex-col gap-1.5 rounded-lg bg-danger-soft px-3 py-2 text-xs text-danger-text">
      Chưa có quyền ghi màn hình (Screen Recording) – app không chụp được màn hình để tìm ảnh / đọc màu.
      <span>System Settings → Privacy &amp; Security → Screen Recording → bật cho app (khi dev: terminal / VS Code), rồi mở lại app.</span>
      <button onClick={requestCapture} className="self-start rounded-control border border-danger/40 px-2.5 py-1 font-semibold hover:bg-danger/10">Xin quyền</button>
    </div>
  );
}

// Ảnh mẫu (hoặc chữ) + thử tìm + độ khớp / chờ tối đa / vùng tìm – dùng cho Chờ ảnh và Click theo ảnh (tự tìm).
function ImageSearchFields({ step: s }: { step: Step }) {
  const { finds, busy, error } = useImageState();
  const find = finds[s.id];
  const [size, setSize] = useState<{ w: number; h: number } | null>(null);
  const k = s.templateScale || 1;
  const byText = s.findBy === 'text';
  return (
    <>
      <Field label="Tìm theo" wide group>
        <Segmented value={byText ? 'text' : 'image'} full label="Tìm theo"
          onChange={(v) => updateStep(s.id, { findBy: v === 'text' ? 'text' : undefined })}
          options={[{ id: 'image', label: 'Ảnh mẫu' }, { id: 'text', label: 'Chữ' }]} />
        <span className={hint}>
          {byText ? 'Đọc chữ trên màn hình – không bị ảnh hưởng bởi theme sáng / tối, zoom, đổi màu khi rê chuột. Dùng khi nút / nhãn có chữ.'
            : 'So từng điểm ảnh – dùng cho icon, nút không có chữ. Nút có chữ nên chọn “Chữ”.'}
        </span>
      </Field>
      {byText ? <TextSearchFields step={s} /> : <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span>Ảnh mẫu</span>
        <div className="relative grid min-h-[92px] place-items-center overflow-hidden rounded-card border border-line bg-soft p-2 text-[13px]">
          {s.template ? (
            <>
              <img src={templateUrl(s.template)} alt="Ảnh mẫu" className="max-h-[140px] max-w-full object-contain"
                style={size ? { width: Math.min(size.w / k, 300) } : undefined}
                onLoad={(e) => setSize({ w: e.currentTarget.naturalWidth, h: e.currentTarget.naturalHeight })} />
              {size && (
                <span className="absolute right-1.5 bottom-1.5 rounded bg-hud/80 px-1.5 font-mono text-[11px] text-hud-text">
                  {Math.round(size.w / k)} × {Math.round(size.h / k)}
                </span>
              )}
            </>
          ) : 'Chưa có ảnh mẫu – bấm “Khoanh ảnh mẫu”, kéo khung quanh nút / chữ cần chờ'}
        </div>
        <div className="flex gap-1.5">
          <button onClick={() => startCapture(s.id)} className="h-8 flex-1 rounded-control border border-line bg-surface text-ink hover:bg-hover">
            {s.template ? 'Khoanh lại' : 'Khoanh ảnh mẫu'}
          </button>
          <button onClick={() => testFind(s.id)} disabled={!s.template || busy === s.id}
            title="Tìm ảnh mẫu trên màn hình hiện tại (app ẩn đi một chút) – để biết độ khớp thực tế"
            className="h-8 flex-1 rounded-control border border-line bg-surface text-ink enabled:hover:bg-hover disabled:opacity-50">
            {busy === s.id ? 'Đang tìm…' : 'Thử tìm'}
          </button>
        </div>
        {find && (
          <span className={`flex flex-wrap items-center gap-x-2 text-xs ${find.found ? 'text-ok' : 'text-warn'}`}>
            {find.found
              ? <>Thấy tại X {find.x + Math.round(find.w / 2)} Y {find.y + Math.round(find.h / 2)} · khớp {find.score}% · {find.ms} ms</>
              : find.fromRun
                ? <>Lần chạy trước hết giờ không thấy – chỗ giống nhất khớp {find.score}% (cần {s.threshold}%)</>
                : <>Không thấy trên màn hình hiện tại – khớp nhất {find.score}% (cần {s.threshold}%)</>}
            <button onClick={() => openOverlay('point', s.id)} className="text-muted underline">xem trên màn hình</button>
          </span>
        )}
        {error && busy === '' && <span className="text-xs text-danger-text">{error}</span>}
      </div>}
      {!byText && <NumField label="Độ khớp tối thiểu" step={s} field="threshold" unit="%" min={1} />}
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 250} min={20} unit="ms" onChange={(intervalMs) => updateStep(s.id, { intervalMs })} />
      </Field>
      <span className={`self-end pb-2 ${hint}`}>Không thấy sau {fmtSec(s.timeoutMs)} thì dừng kịch bản</span>
      {!byText && (s.threshold ?? 0) < 70 && (
        <span className="col-span-2 -mt-1 text-xs text-warn">
          Độ khớp {s.threshold}% rất thấp – gần như chỗ nào cũng “khớp”, dễ click nhầm. Ảnh khoanh đúng thường khớp 95–100%;
          nên để 80–90% và bấm Thử tìm để xem độ khớp thật.
        </span>
      )}
      <Field label="Tìm trong" wide group>
        <Segmented value={s.searchArea ?? 'region'} onChange={(v) => updateStep(s.id, { searchArea: v === 'screen' ? 'screen' : undefined })} full
          options={[{ id: 'region', label: 'Vùng của kịch bản' }, { id: 'screen', label: 'Cả màn hình' }]} />
        <span className={hint}>
          {byText ? 'Vùng nhỏ đọc nhanh hơn (cả màn ~1 s mỗi lần) và ít nhầm chữ trùng ở chỗ khác.'
            : 'Vùng nhỏ tìm nhanh hơn và ít nhầm hơn. Ảnh mẫu chụp ở màn Retina vẫn tìm được trên màn thường.'}
          {' '}Cả màn hình = cả màn chứa vùng.
        </span>
      </Field>
    </>
  );
}

// Tìm theo chữ: ô chữ (nhận {{biến}}), cách so, chỗ thứ mấy, thử tìm.
function TextSearchFields({ step: s }: { step: Step }) {
  const script = useScript();
  const { finds, busy, error } = useImageState();
  const find = finds[s.id];
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const empty = !s.findText?.trim();
  return (
    <div className="col-span-2 flex flex-col gap-1.5 text-muted">
      <span className="flex items-center justify-between">Chữ cần tìm
        <span className="text-xs">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
      </span>
      <VarTextInput value={s.findText ?? ''} onChange={(findText) => set({ findText })} inputs={script.inputs ?? []}
        placeholder="Xem lịch sử" onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
      <span className={`text-xs ${empty ? 'text-danger-text' : ''}`}>
        {empty ? 'Chưa nhập chữ – chạy sẽ báo lỗi. ' : ''}Không phân biệt hoa / thường, bỏ qua khoảng trắng. Click vào giữa chữ.
      </span>
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2 py-1 text-[13px] text-ink">
        <label className="flex items-center gap-2"><Switch on={!!s.textExact} onChange={(textExact) => set({ textExact: textExact || undefined })} label="Cả dòng" />Cả dòng đúng chữ này</label>
        <label className="flex items-center gap-2"><Switch on={!!s.textNoAccent} onChange={(textNoAccent) => set({ textNoAccent: textNoAccent || undefined })} label="Bỏ qua dấu" />Bỏ qua dấu</label>
        <span className="flex items-center gap-2">Chỗ thứ
          <NumInput value={Math.max(s.occurrence ?? 1, 1)} min={1} onChange={(n) => set({ occurrence: n > 1 ? n : undefined })} className="w-20" />
        </span>
      </div>
      <span className={hint}>
        “Cả dòng”: chỉ khớp ô / nút chỉ có đúng chữ này (tránh khớp trong câu dài). “Bỏ qua dấu”: khi máy đọc sai dấu.
        Chữ hiện nhiều lần → chỗ thứ mấy tính từ trên xuống, trái sang phải.
      </span>
      <button onClick={() => testFind(s.id)} disabled={empty || busy === s.id}
        title="Đọc chữ trên màn hình hiện tại và tìm (app ẩn đi một chút)"
        className="h-8 rounded-control border border-line bg-surface text-ink enabled:hover:bg-hover disabled:opacity-50">
        {busy === s.id ? 'Đang đọc chữ…' : 'Thử tìm'}
      </button>
      {find && (
        <span className={`flex flex-col gap-1 text-xs ${find.found ? 'text-ok' : 'text-warn'}`}>
          <span className="flex flex-wrap items-center gap-x-2">
            {find.found ? <>Thấy tại X {find.x + Math.round(find.w / 2)} Y {find.y + Math.round(find.h / 2)} · {find.ms} ms</>
              : find.fromRun ? <>Lần chạy trước hết giờ không thấy chữ này</>
              : <>Không thấy trên màn hình hiện tại{find.ms ? ` · ${find.ms} ms` : ''}</>}
            {find.found && <button onClick={() => openOverlay('point', s.id)} className="text-muted underline">xem trên màn hình</button>}
          </span>
          {!find.found && (
            <span className="text-muted">
              {find.read?.length ? <>Đọc được: {find.read.slice(0, 12).map((t, i) => (
                <button key={i} onClick={() => set({ findText: t })} title="Dùng chữ này"
                  className="mr-1 mb-1 inline-block rounded border border-line px-1.5 text-ink hover:bg-hover">{t}</button>
              ))}</> : 'Không đọc được chữ nào trong vùng tìm.'}
            </span>
          )}
        </span>
      )}
      {error && busy === '' && <span className="text-xs text-danger-text">{error}</span>}
    </div>
  );
}

// Khi thất bại: thử lại N lần, nghỉ giữa các lần; hết lượt vẫn lỗi → dừng kịch bản (cổng fail – nối nhánh khác trên canvas, bước 20 F4b).
function RetryFields({ step: s }: { step: Step }) {
  const r = s.retry ?? {};
  const set = (patch: Partial<NonNullable<Step['retry']>>) => {
    const next = { ...r, ...patch };
    updateStep(s.id, { retry: next.count ? next : undefined });
  };
  const what = s.type === 'open' ? 'mở không được' : s.type === 'waitPixel' ? 'không đổi màu' : 'không thấy';
  return (
    <Field label={`Khi ${what}`} wide group>
      <div className="flex flex-wrap items-center gap-2 text-[13px] text-ink">
        <span className="flex items-center gap-2">Thử lại
          <NumInput value={r.count ?? 0} min={0} onChange={(n) => set({ count: Math.min(n, 100) })} className="w-20" unit="lần" />
        </span>
        {!!r.count && (
          <span className="flex items-center gap-2">nghỉ
            <NumInput value={r.delayMs ?? 0} min={0} onChange={(delayMs) => set({ delayMs })} className="w-24" unit="ms" />
          </span>
        )}
      </div>
      <span className="text-xs text-muted">
        {`Thử hết vẫn ${what} thì dừng kịch bản, báo lỗi ở bước này${r.count ? ` (sau ${r.count + 1} lần thử)` : ''}. Rẽ sang nhánh khác khi thất bại: làm trên canvas flow (sắp có).`}
      </span>
    </Field>
  );
}

const fmtSec = (ms?: number) => `${((ms ?? 0) / 1000).toLocaleString('vi-VN', { maximumFractionDigits: 1 })} s`;

function ClickImageFields({ step: s }: { step: Step }) {
  const { steps } = useScript();
  const i = steps.findIndex((x) => x.id === s.id);
  // bước chờ ảnh / click tự tìm gần nhất phía trước – nguồn vị trí khi "Kết quả bước trước"
  const prev = steps.slice(0, i).map((x, j) => ({ x, n: j + 1 })).reverse()
    .find(({ x }) => x.enabled && (x.type === 'waitImage' || (x.type === 'clickImage' && x.imageSource === 'search')));
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
    <>
      <Field label="Click vào" wide group>
        <Segmented value={s.imageSource ?? 'node'} onChange={(imageSource: 'node' | 'search') => set({ imageSource })} full
          options={[{ id: 'node', label: prev ? `Kết quả bước ${prev.n}` : 'Kết quả bước trước' }, { id: 'search', label: 'Tự tìm (ảnh / chữ)' }]} label="Click vào" />
        <span className={`text-xs ${s.imageSource !== 'search' && !prev ? 'text-danger-text' : 'text-muted'}`}>
          {s.imageSource === 'search' ? `Tự tìm ${s.findBy === 'text' ? 'chữ' : 'ảnh mẫu'} lúc chạy (chờ như bước ${s.findBy === 'text' ? 'Chờ chữ' : 'Chờ ảnh'}) rồi click vào giữa`
            : prev ? `Click vào chỗ bước ${prev.n} (${stepLabel(prev.x)}) tìm thấy – không chụp lại`
            : 'Phía trước chưa có bước Chờ ảnh / Chờ chữ – thêm một bước, hoặc chọn “Tự tìm”'}
        </span>
      </Field>
      {s.imageSource === 'search' && <ImageSearchFields step={s} />}
      <Field label={`Lệch so với tâm ${s.imageSource === 'search' && s.findBy === 'text' ? 'chữ' : 'chỗ thấy'}`} wide group>
        <div className="flex gap-2">
          <NumInput signed prefix="ΔX" value={s.dx ?? 0} onChange={(dx) => set({ dx })} className="flex-1" />
          <NumInput signed prefix="ΔY" value={s.dy ?? 0} onChange={(dy) => set({ dy })} className="flex-1" />
        </div>
        <span className={hint}>Vd tìm nhãn “Họ tên”, ô nhập nằm bên phải 150 px → ΔX 150.</span>
      </Field>
      <Field label="Nút chuột" wide group>
        <Segmented value={s.button} onChange={(button: MouseButton) => set({ button })} full label="Nút chuột"
          options={[{ id: 'left', label: 'Trái' }, { id: 'right', label: 'Phải' }, { id: 'center', label: 'Giữa' }]} />
      </Field>
    </>
  );
}

function WaitPixelFields({ step: s, setPoint }: { step: Step; setPoint: (i: number, p: Partial<StepPoint>) => void }) {
  const { region } = useScript();
  const { busy, error } = useImageState();
  const valid = /^#[0-9a-fA-F]{6}$/.test(s.color ?? '');
  return (
    <>
      {s.points[0] && (
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">
          <span>Điểm</span>
          <AnchorField value={s.points[0]} onChange={(p) => setPoint(0, p)} onAnchor={(ax, ay) => setAnchor(s.id, 0, ax, ay)}
            pick={<PickButton step={s} />} />
        </div>
      )}
      <Field label="Màu cần chờ" wide group>
        <div className="flex gap-2">
          <span className={`${fieldBox} flex-1 gap-2 ${valid ? '' : 'border-danger'}`}>
            <span className="h-4 w-4 flex-none rounded border border-line" style={{ background: valid ? s.color : 'transparent' }} />
            <input value={s.color ?? ''} onChange={(e) => updateStep(s.id, { color: e.target.value.trim() })} spellCheck={false}
              className="min-w-0 flex-1 bg-transparent font-mono text-[13px] text-ink outline-none" />
          </span>
          <button onClick={() => pickColor(s.id)} disabled={busy === s.id}
            title="Lấy màu đang hiện tại điểm này trên màn hình (app ẩn đi một chút) – mở sẵn màn hình có màu cần chờ"
            className="h-8 rounded-control border border-line px-2.5 text-[13px] whitespace-nowrap text-ink enabled:hover:bg-hover disabled:opacity-50">
            {busy === s.id ? 'Đang lấy…' : 'Lấy màu'}
          </button>
        </div>
        {!valid && <span className="text-xs text-danger-text">Màu dạng #RRGGBB, vd #1877F2</span>}
        {error && busy === '' && <span className="text-xs text-danger-text">{error}</span>}
      </Field>
      <Field label="Sai số màu">
        <NumInput prefix="±" value={s.tolerance ?? 0} onChange={(tolerance) => updateStep(s.id, { tolerance: Math.min(tolerance, 255) })} />
      </Field>
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <span className={`col-span-2 ${hint}`}>
        Mỗi {s.intervalMs ?? 100} ms đọc màu tại điểm{region ? ' (trong vùng)' : ''}, mỗi kênh R / G / B lệch ≤ {s.tolerance} là được.
        Không đổi màu sau {fmtSec(s.timeoutMs)} thì dừng kịch bản. Không di chuột.
      </span>
    </>
  );
}
