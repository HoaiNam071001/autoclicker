// Node chuột: nút chuột, vị trí (neo), kéo thả A → B, số click / thời gian giữ / kéo, ngẫu nhiên vị trí.
import { useState } from 'react';
import NumInput from '../NumInput';
import Segmented from '../Segmented';
import AnchorField, { AnchorMap } from '../steps/AnchorField';
import { boundsOf, endpoints, refOf, setAnchor, setDragFromPrev, setDragRelative, updateStep, useScript } from '../../store';
import { MouseButton, Step } from '../../types';
import { Field, hint, NumField, PickButton, PinSoon, SetPoint } from './shared';

export const isMouse = (s: Step) => s.type === 'click' || s.type === 'double' || s.type === 'hold' || s.type === 'drag';

export function MouseFields({ step: s, setPoint }: { step: Step; setPoint: SetPoint }) {
  const { region } = useScript();
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const mouse = isMouse(s);
  return (
    <>
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
    </>
  );
}

// Ngẫu nhiên vị trí – cuối phần chuột
export function JitterField({ step: s }: { step: Step }) {
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
          <Field label="Ngẫu nhiên vị trí" wide>
        <NumInput prefix="±" value={s.jitterPx} unit="px" onChange={(v) => set({ jitterPx: v })} />
        <span className={hint}>
          {s.jitterPx > 0 ? `Mỗi lần ${s.type === 'drag' ? 'kéo, A và B' : 'click'} lệch trong ô ${2 * s.jitterPx + 1}×${2 * s.jitterPx + 1} px` : 'Tắt (0)'}
        </span>
      </Field>
  );
}

// Điểm A / B của bước kéo thả: A = vị trí hoặc "từ điểm trước"; B = vị trí hoặc lệch so với A.
function DragPoints({ step: s, setPoint }: { step: Step; setPoint: SetPoint }) {
  const script = useScript();
  const { region } = script;
  const ends = endpoints(script, region).get(s.id);
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
              ? `Chỗ chuột dừng ở ${refOf(ends.prevStep ?? '')} (X ${ends.start.x} Y ${ends.start.y})`
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
