// Điểm lấy toạ độ từ biến (F5b) – chỉ hiển thị: ô X / Y nhận số / biến cộng trừ px, nút điền nhanh chỗ thấy của node trước.
import Select from '../Select';
import VarTextInput from '../inputs/VarTextInput';
import type { ScriptInput } from '../inputs/model';
import { StepPoint } from '../../types';
import { hint, selectBox } from './shared';

type Props = {
  point: StepPoint;
  onChange: (p: Partial<StepPoint>) => void;
  relative?: boolean; // B "lệch so với A": ô là độ lệch, không phải toạ độ
  inputs: Pick<ScriptInput, 'key' | 'label'>[];
  nodeVars: { key: string; label: string }[];
  sources: { ref: string; label: string }[]; // node phía trước có output x, y (flow/view.ts coordSources)
  onCreateInput?: (key: string) => void;
};

export default function PointVarFields({ point: p, onChange, relative = false, inputs, nodeVars, sources, onCreateInput }: Props) {
  return (
    <div className="flex flex-col gap-1.5">
      {(['x', 'y'] as const).map((k) => (
        <div key={k} className="flex items-start gap-2">
          <span className="w-6 flex-none pt-[7px] font-mono text-xs">{relative ? 'Δ' : ''}{k.toUpperCase()}</span>
          <div className="min-w-0 flex-1">
            <VarTextInput value={(k === 'x' ? p.xVar : p.yVar) ?? ''} rows={1} compact inputs={inputs} nodeVars={nodeVars} onCreateInput={onCreateInput}
              placeholder={`{{n2.output.${k}}} + 20`} onChange={(v) => onChange(k === 'x' ? { xVar: v } : { yVar: v })} />
          </div>
        </div>
      ))}
      {!relative && (sources.length ? (
        <Select value="" onChange={(ref) => onChange({ xVar: `{{${ref}.output.x}}`, yVar: `{{${ref}.output.y}}` })}
          className={selectBox} aria-label="Điền chỗ thấy của node" placeholder="Điền chỗ thấy / chỗ chuột dừng của node…"
          options={sources.map((s) => ({ value: s.ref, label: s.label }))} />
      ) : (
        <span className={hint}>Chưa có node nào phía trước trả toạ độ (Chờ ảnh / chữ, Click, phần tử web…).</span>
      ))}
      <span className={hint}>
        {relative ? 'Độ lệch so với A (px)' : 'Toạ độ màn hình (không tính theo vùng)'} – số hoặc biến, cộng / trừ px được, vd {'{{n2.output.x}} + 20'}.
      </span>
    </div>
  );
}
