// Phần dùng chung của inspector node: hàng trường (nhãn trái, ô phải), ô số có menu ⋯, nút chọn lại điểm, tab, hằng style.
import { ReactNode, useMemo } from 'react';
import NumInput, { fieldBox } from '../NumInput';
import { Target } from '../icons';
import PinButton, { KeyChip } from '../inputs/PinButton';
import { FieldMenu, FieldMenuItem, MENU_ICON } from '../FieldMenu';
import { defaultOf } from '../../flow/defaults';
import VarTextInput from '../inputs/VarTextInput';
import { pickPointFor, setNumVarExpr, setNumVarOn, updateStep, useScript } from '../../store';
import { addInput, compatibleInputs, pinField, pinnedInput, pinToInput, setPinnedValue, unpinField } from '../../scriptInputs';
import { MouseButton, numVarOf, NumVarField, ScrollDirection, Step, StepPoint } from '../../types';
import Segmented from '../Segmented';
import { nodeVarOptions } from '../../flow/view';

export const selectBox = `${fieldBox} w-full text-ink`;
export const hint = 'text-xs text-muted';
// Ô chọn hướng cuộn (node Cuộn, Cuộn trang web).
export const DIRECTION_OPTIONS: { id: ScrollDirection; label: string }[] = [{ id: 'up', label: '↑ Lên' }, { id: 'down', label: '↓ Xuống' }, { id: 'left', label: '← Trái' }, { id: 'right', label: '→ Phải' }];
export const fmtSec = (ms?: number) => `${((ms ?? 0) / 1000).toLocaleString('vi-VN', { maximumFractionDigits: 1 })} s`;
// Biến kết quả các node khác ({{n2.output.x}}) cho ô chữ của node id
export function useNodeVars(id: string) {
  const { steps, start, edges } = useScript();
  return useMemo(() => nodeVarOptions({ steps, start, edges }, id), [steps, start, edges, id]);
}

// Sửa một điểm của node (giữ các trường khác của điểm)
export type SetPoint = (i: number, p: Partial<StepPoint>) => void;

export function PickButton({ step: s, title = 'Chọn lại trên màn hình' }: { step: Step; title?: string }) {
  return (
    <button aria-label={title} title={title} onClick={() => pickPointFor(s.id)}
      className="grid h-8 w-9 flex-none place-items-center rounded-control border border-line text-muted hover:bg-hover hover:text-ink"><Target /></button>
  );
}


// Hàng trường design v3 (FieldKit): nhãn 112 px bên trái, ô bên phải. wide: nhãn ở trên, ô rộng hết hàng (ô chữ dài,
// cặp X / Y, danh sách). group: bên trong là nhóm nút (Segmented) – dùng <div>, vì <label> bọc nút thì bấm chữ nhãn
// sẽ bấm luôn nút đầu. extra: chip dưới nhãn (key input đã ghim).
export function Field({ label, children, wide = false, group = false, extra }: {
  label: ReactNode; children: ReactNode; wide?: boolean; group?: boolean; extra?: ReactNode;
}) {
  const Tag = group ? 'div' : 'label';
  const head = <span className="flex min-w-0 flex-col items-start gap-[3px] text-muted"><span className="max-w-full">{label}</span>{extra}</span>;
  if (wide) return <Tag className="col-span-2 flex flex-col gap-1.5">{head}{children}</Tag>;
  return <Tag className="col-span-2 grid grid-cols-[112px_minmax(0,1fr)] items-center gap-2.5">{head}{children}</Tag>;
}

// Ô số trong inspector (design v3 FieldKit). Nút ⋯ trong ô: dùng biến (F5c – ô thành biểu thức, vd {{n2.output.count}} + 1),
// đặt làm input khi chạy (bước 21b – ghim: viền nét đứt, chip key dưới nhãn, sửa số = sửa mặc định của input; dùng chung
// input cùng kiểu đã có), về mặc định của loại node (flow/defaults). stepper: số nguyên nhỏ có − / +.
export function NumField({ label, step: s, field, unit, prefix, min = 0, stepper = false, help }: {
  label: string; step: Step; field: NumVarField; unit?: string; prefix?: string; min?: number; stepper?: boolean; help?: string;
}) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const inp = pinnedInput(script, s.id, field);
  const others = compatibleInputs(script, field).filter((i) => i !== inp);
  const expr = numVarOf(s, field);
  const def = defaultOf(s.type, field);
  const value = s[field] ?? min;
  const items: FieldMenuItem[] = expr !== undefined ? [] : [
    ...(inp ? [] : [{ label: 'Dùng biến…', icon: MENU_ICON.variable, sub: 'vd {{n2.output.count}} + 1', onClick: () => setNumVarOn(s.id, field, true) }]),
    inp
      ? { label: `Bỏ input ${inp.key}`, icon: MENU_ICON.pin, sub: 'không hỏi khi chạy nữa', onClick: () => unpinField(s.id, field) }
      : { label: 'Đặt làm input khi chạy', icon: MENU_ICON.pin, sub: 'giá trị hiện tại làm mặc định', onClick: () => pinField(s.id, field) },
    ...(inp ? [] : others.map((o) => ({ label: o.key, sub: `dùng chung input: ${o.label}`, mono: true, onClick: () => pinToInput(s.id, field, o.key) }))),
    ...(def !== undefined && def !== value && !inp ? [{ label: `Về mặc định (${def}${unit ? ` ${unit}` : ''})`, icon: MENU_ICON.reset, onClick: () => updateStep(s.id, { [field]: def }) }] : []),
  ];
  return (
    <Field label={label} group extra={inp && <KeyChip k={inp.key} />}>
      {expr !== undefined ? (
        <span className="flex min-w-0 items-center gap-1">
          <span className="min-w-0 flex-1">
            <VarTextInput value={expr} rows={1} compact inputs={script.inputs ?? []} nodeVars={nodeVars} placeholder="{{n2.output.count}}"
              onCreateInput={(key) => addInput({ key, label: key, type: 'number' })} onChange={(v) => setNumVarExpr(s.id, field, v)} />
          </span>
          <button type="button" onClick={() => setNumVarOn(s.id, field, false)} aria-label="Bỏ biến" title="Bỏ biến – dùng số cố định"
            className="grid h-6 w-6 flex-none place-items-center rounded-[5px] text-muted hover:bg-hover hover:text-ink">
            <svg width="10" height="10" viewBox="0 0 10 10" aria-hidden><path d="M1 1l8 8M9 1L1 9" stroke="currentColor" strokeWidth="1.4" /></svg>
          </button>
        </span>
      ) : (
        <NumInput value={value} min={min} unit={unit} prefix={prefix} pinned={!!inp} tag={inp ? 'mặc định' : undefined} stepper={stepper}
          menu={<FieldMenu items={items} help={help} />}
          onChange={(v) => (inp ? setPinnedValue(s.id, field, v) : updateStep(s.id, { [field]: v }))} />
      )}
    </Field>
  );
}

export const PinSoon = () => <PinButton pinned={false} disabled title="Ghim toạ độ làm input – sắp có" />;

export function TabBtn({ active, onClick, title, children }: { active: boolean; onClick: () => void; title?: string; children: ReactNode }) {
  return (
    <button role="tab" aria-selected={active} onClick={onClick} title={title}
      className={`relative flex h-[34px] items-center gap-1.5 px-0.5 ${active ? 'font-semibold' : 'text-muted hover:text-ink'}`}>
      {children}
      {active && <span className="absolute inset-x-[15%] -bottom-px h-[3px] rounded-sm bg-accent" />}
    </button>
  );
}

// Nút chuột của node bấm (click, click theo ảnh, click phần tử web).
export function ButtonField({ step: s }: { step: Step }) {
  return (
    <Field label="Nút chuột" wide group>
      <Segmented value={s.button} onChange={(button: MouseButton) => updateStep(s.id, { button })} full label="Nút chuột"
        options={[{ id: 'left', label: 'Trái' }, { id: 'right', label: 'Phải' }, { id: 'center', label: 'Giữa' }]} />
    </Field>
  );
}
