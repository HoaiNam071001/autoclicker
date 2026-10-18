// Phần dùng chung của inspector node: ô có nhãn, ô số ghim được, nút chọn lại điểm, tab, hằng style.
import { ReactNode, useMemo } from 'react';
import NumInput, { fieldBox } from '../NumInput';
import { Target } from '../icons';
import PinButton, { KeyChip, PinMenu, VarButton } from '../inputs/PinButton';
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


// group: bên trong là nhóm nút (Segmented) – dùng <div>, vì <label> bọc nút thì bấm chữ nhãn sẽ bấm luôn nút đầu
export function Field({ label, children, wide = false, group = false }: {
  label: string; children: ReactNode; wide?: boolean; group?: boolean;
}) {
  const Tag = group ? 'div' : 'label';
  return <Tag className={`flex flex-col gap-1.5 text-muted ${wide ? 'col-span-2' : ''}`}>{label}{children}</Tag>;
}

// Ô số ghim được (bước 21b): đã ghim → viền nét đứt + chip key cạnh nhãn, sửa số = sửa mặc định của input
// (mọi trường ghim cùng input đổi theo). Nút {x} (F5c): ô thành ô biến chiếm cả hàng (biểu thức dài) – giá trị lúc chạy
// lấy từ biến, không ghim được.
export function NumField({ label, step: s, field, unit, prefix, min = 0 }: {
  label: string; step: Step; field: NumVarField; unit?: string; prefix?: string; min?: number;
}) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const inp = pinnedInput(script, s.id, field);
  const others = compatibleInputs(script, field).filter((i) => i !== inp);
  const expr = numVarOf(s, field);
  return (
    <div className={`flex flex-col gap-1.5 text-muted ${expr !== undefined ? 'col-span-2' : ''}`}>
      <span className="flex min-w-0 items-center gap-1.5"><span className="truncate">{label}</span>{inp && <KeyChip k={inp.key} />}</span>
      <span className="flex items-start gap-1">
        {expr !== undefined ? (
          <div className="min-w-0 flex-1">
            <VarTextInput value={expr} rows={1} compact inputs={script.inputs ?? []} nodeVars={nodeVars} placeholder="{{n2.output.count}}"
              onCreateInput={(key) => addInput({ key, label: key, type: 'number' })} onChange={(v) => setNumVarExpr(s.id, field, v)} />
          </div>
        ) : (
          <NumInput value={s[field] ?? min} min={min} unit={unit} prefix={prefix} pinned={!!inp} className="min-w-0 flex-1"
            onChange={(v) => (inp ? setPinnedValue(s.id, field, v) : updateStep(s.id, { [field]: v }))} />
        )}
        <VarButton on={expr !== undefined} disabled={!!inp} onToggle={() => setNumVarOn(s.id, field, expr === undefined)} />
        {expr === undefined && (
          <PinMenu pinnedKey={inp?.key} options={others} onPinNew={() => pinField(s.id, field)}
            onPinTo={(key) => pinToInput(s.id, field, key)} onUnpin={() => unpinField(s.id, field)} />
        )}
      </span>
    </div>
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
