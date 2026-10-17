// Phần dùng chung của inspector node: ô có nhãn, ô số ghim được, nút chọn lại điểm, tab, hằng style.
import { ReactNode } from 'react';
import NumInput, { fieldBox } from '../NumInput';
import { Target } from '../icons';
import PinButton, { KeyChip, PinMenu } from '../inputs/PinButton';
import { pickPointFor, updateStep, useScript } from '../../store';
import { compatibleInputs, pinField, pinnedInput, pinToInput, setPinnedValue, unpinField } from '../../scriptInputs';
import { PinField, Step, StepPoint } from '../../types';

export const selectBox = `${fieldBox} w-full text-ink`;
export const hint = 'text-xs text-muted';
export const fmtSec = (ms?: number) => `${((ms ?? 0) / 1000).toLocaleString('vi-VN', { maximumFractionDigits: 1 })} s`;
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
// (mọi trường ghim cùng input đổi theo).
export function NumField({ label, step: s, field, unit, prefix, min = 0 }: {
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
