// Thao tác giống người (bước 34) – chỉ hiển thị: mức của kịch bản (Cài đặt chạy) và ô ghi đè của node (inspector).
import { fieldBox } from './NumInput';
import Segmented from './Segmented';
import { HUMAN_LEVELS, humanHint, humanLabel, scriptHuman } from '../flow/human';
import type { HumanLevel } from '../types';

const hintCls = 'text-xs text-muted';

export function HumanScriptField({ value, onChange }: { value?: HumanLevel; onChange: (lv: HumanLevel) => void }) {
  const lv = scriptHuman(value);
  return (
    <div className="flex flex-col gap-2">
      <span className="flex flex-col">Thao tác giống người
        <span className={hintCls}>Giảm nguy cơ trang web nhận là bot – không đảm bảo; chạy chậm hơn</span>
      </span>
      <Segmented value={lv} options={HUMAN_LEVELS} onChange={onChange} label="Thao tác giống người" full />
      <span className={hintCls}>{humanHint(lv)}</span>
    </div>
  );
}

// Ô ghi đè của node: trống = theo kịch bản (hiện mức đang dùng).
export function HumanNodeField({ value, scriptLevel, onChange, className = '' }: {
  value?: HumanLevel; scriptLevel?: HumanLevel; onChange: (lv: HumanLevel | undefined) => void; className?: string;
}) {
  return (
    <label className={`flex flex-col gap-1.5 text-muted ${className}`}>Thao tác giống người
      <select value={value ?? ''} onChange={(e) => onChange((e.target.value || undefined) as HumanLevel | undefined)}
        className={`${fieldBox} w-full text-ink`}>
        <option value="">Theo kịch bản ({humanLabel(scriptHuman(scriptLevel))})</option>
        {HUMAN_LEVELS.map((l) => <option key={l.id} value={l.id}>{l.label}</option>)}
      </select>
      {value && <span className={hintCls}>{humanHint(value)}</span>}
    </label>
  );
}
