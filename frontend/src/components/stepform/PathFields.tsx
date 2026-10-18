// Node Vẽ đường (bước 35 N7): hình (qua các điểm / tròn / chữ nhật / sao), cỡ hình, danh sách điểm, thời gian vẽ.
// Tâm hình / vị trí điểm sửa ở overlay (nút 👁) hoặc ô Vị trí (MouseFields).
import NumInput from '../NumInput';
import Segmented from '../Segmented';
import { Switch } from '../RunSettings';
import { updateStep } from '../../store';
import { fmtPoint, usesVar } from '../../anchor';
import { SHAPES, type Shape } from '../../flow/shape';
import { Step, StepPoint } from '../../types';
import { Field, hint, NumField } from './shared';

export function PathFields({ step: s }: { step: Step }) {
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const shape = s.shape ?? '';
  const setShape = (v: Shape) => set({
    shape: v || undefined,
    ...(v && !s.shapeW ? { shapeW: 200 } : {}),
    // từ "qua các điểm" sang hình: tâm = giữa các điểm (cùng neo, không biến) – hình không lệch sang điểm đầu
    ...(v && !shape && s.points.length > 1 ? { points: [centerOf(s.points)] } : {}),
    // sang "qua các điểm" mà chỉ còn tâm → thêm điểm thứ hai lệch 150 px
    ...(!v && s.points.length < 2 ? { points: [...s.points, shifted(s.points[0], 150, 0)] } : {}),
  });
  return (
    <>
      <Field label="Hình" wide group>
        <Segmented value={shape} options={SHAPES.map((x) => (x.id ? x : { ...x, label: 'Các điểm' }))} onChange={setShape} full label="Hình" />
      </Field>
      {shape ? <>
        <Field label={shape === 'circle' ? 'Đường kính / rộng' : 'Rộng'}>
          <NumInput value={s.shapeW ?? 0} min={4} unit="px" onChange={(shapeW) => set({ shapeW })} />
        </Field>
        {shape !== 'star' && <Field label="Cao">
          <NumInput value={s.shapeH || s.shapeW || 0} min={4} unit="px" onChange={(shapeH) => set({ shapeH: shapeH === s.shapeW ? undefined : shapeH })} />
        </Field>}
      </> : <PathPoints step={s} />}
      <NumField label="Vẽ cả nét trong" step={s} field="dragMs" unit="ms" min={1} />
      <span className={`col-span-2 -mt-1 ${hint}`}>Giữ chuột từ đầu tới cuối – một nét liền. Bấm 👁 ở đầu để xem / kéo {shape ? 'tâm hình' : 'các điểm'} trên màn hình.</span>
    </>
  );
}

// Trung bình các điểm khi cùng neo và không dùng biến; không thì giữ điểm đầu.
function centerOf(ps: StepPoint[]): StepPoint {
  const [a] = ps;
  if (ps.some((p) => usesVar(p) || p.ax !== a.ax || p.ay !== a.ay)) return a;
  const avg = (k: 'x' | 'y') => Math.round(ps.reduce((n, p) => n + p[k], 0) / ps.length);
  return { ...a, x: avg('x'), y: avg('y') };
}

const shifted = (p: StepPoint | undefined, dx: number, dy: number): StepPoint =>
  p && !usesVar(p) ? { ...p, x: p.x + dx, y: p.y + dy } : { x: dx, y: dy, ax: 'center', ay: 'center' };

function PathPoints({ step: s }: { step: Step }) {
  const set = (points: StepPoint[]) => updateStep(s.id, { points });
  return (
    <div className="col-span-2 flex flex-col gap-1.5 text-muted">
      <span className="flex items-center justify-between">Các điểm (đi lần lượt)<span className="text-xs">{s.points.length} điểm</span></span>
      <ol className="flex flex-col gap-1">
        {s.points.map((p, i) => (
          <li key={i} className="flex h-8 items-center gap-2 rounded-control border border-line bg-field px-2.5">
            <span className="w-4 font-mono text-xs">{i + 1}</span>
            <span className="min-w-0 flex-1 truncate font-mono text-[13px] text-ink">{usesVar(p) ? `${p.xVar}, ${p.yVar}` : fmtPoint(p)}</span>
            <button onClick={() => set(s.points.filter((_, j) => j !== i))} disabled={s.points.length <= 2} aria-label={`Xoá điểm ${i + 1}`}
              className="grid h-6 w-6 place-items-center rounded-control hover:bg-hover hover:text-ink disabled:opacity-30">
              <svg viewBox="0 0 24 24" className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" aria-hidden><path d="M6 6l12 12M18 6L6 18" /></svg>
            </button>
          </li>
        ))}
      </ol>
      <div className="flex items-center gap-3">
        <button onClick={() => set([...s.points, shifted(s.points[s.points.length - 1], 80, 40)])}
          className="h-8 rounded-control border border-dashed border-line-strong px-3 text-[13px] text-ink hover:bg-hover">+ Thêm điểm</button>
        <label className="ml-auto flex items-center gap-2 text-ink">
          <Switch on={!!s.closed} onChange={(on) => updateStep(s.id, { closed: on || undefined })} label="Khép kín" />Khép kín
        </label>
      </div>
    </div>
  );
}
