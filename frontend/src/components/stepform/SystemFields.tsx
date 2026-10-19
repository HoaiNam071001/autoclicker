// Node Hệ thống (bước 35 N5): âm lượng loa – đặt mức / tăng / giảm, tắt / bật tiếng.
import Select from '../Select';
import NumInput from '../NumInput';
import { updateStep } from '../../store';
import { SYS_OPS, sysHasLevel, type SysOp } from '../../flow/system';
import { Step } from '../../types';
import { Field, selectBox } from './shared';

export function SystemFields({ step: s }: { step: Step }) {
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const op = s.sysOp ?? 'volume';
  const level = op === 'volume';
  const v = s.volume ?? 0;
  return (
    <>
      <Field label="Làm gì" help={`Loa mặc định của máy · node sau đọc {{${s.ref}.output.volume}}, muted.`}>
        <Select value={op as string} onChange={(v) => set({ sysOp: v as SysOp, ...(v === 'volume' && !s.volume ? { volume: 50 } : {}) })} className={selectBox} aria-label="Thao tác"
          options={SYS_OPS.map((o) => ({ value: o.id as string, label: o.label }))} />
      </Field>
      {sysHasLevel(op) && (
        <Field label={level ? 'Mức âm lượng' : 'Mỗi lần'} wide group help={level ? 'Đặt mức cũng bật tiếng lại nếu đang tắt.' : 'Giữ trong 0–100 %; đang tắt tiếng thì bật lại.'}>
          <div className="flex items-center gap-3">
            <input type="range" min={level ? 0 : 1} max={level ? 100 : 50} value={level ? v : v || 10} aria-label={level ? 'Mức âm lượng' : 'Mỗi lần'}
              onChange={(e) => set({ volume: Number(e.target.value) })} className="min-w-0 flex-1 accent-[var(--c-accent)]" />
            <NumInput value={level ? v : v || 10} min={level ? 0 : 1} unit="%" onChange={(n) => set({ volume: Math.min(n, 100) })} className="w-[92px] flex-none" />
          </div>
        </Field>
      )}
    </>
  );
}
