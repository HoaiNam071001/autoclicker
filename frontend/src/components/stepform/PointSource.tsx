// Một điểm của node: công tắc Số / Biến (F5b). Số → children (ô toạ độ + neo); Biến → PointVarFields.
import { ReactNode, useMemo } from 'react';
import { usesVar } from '../../anchor';
import { coordSources } from '../../flow/view';
import { addInput } from '../../scriptInputs';
import { setPointVar, useScript } from '../../store';
import { Step } from '../../types';
import Segmented from '../Segmented';
import PointVarFields from './PointVarFields';
import { SetPoint, useNodeVars } from './shared';

export default function PointSource({ step: s, index, setPoint, label, relative, children }: {
  step: Step; index: number; setPoint: SetPoint; label: ReactNode; relative?: boolean; children: ReactNode;
}) {
  const { steps, start, edges, inputs } = useScript();
  const nodeVars = useNodeVars(s.id);
  const sources = useMemo(() => coordSources({ steps, start, edges }, s.id), [steps, start, edges, s.id]);
  const p = s.points[index];
  const byVar = !!p && usesVar(p);
  return (
    <div className="col-span-2 flex flex-col gap-1.5 text-muted">
      <span className="flex items-center gap-2">
        <span className="min-w-0 flex-1">{label}</span>
        {p && (
          <Segmented value={byVar ? 'var' : 'num'} onChange={(v) => setPointVar(s.id, index, v === 'var')} label="Toạ độ lấy từ"
            options={[{ id: 'num', label: 'Số' }, { id: 'var', label: 'Biến' }]} />
        )}
      </span>
      {byVar ? (
        <PointVarFields point={p} onChange={(q) => setPoint(index, q)} relative={relative} inputs={inputs ?? []} nodeVars={nodeVars}
          sources={sources} onCreateInput={(key) => addInput({ key, label: key, type: 'number' })} />
      ) : children}
    </div>
  );
}
