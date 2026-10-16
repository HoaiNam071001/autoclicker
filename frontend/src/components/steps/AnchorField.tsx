// Ô vị trí có neo theo cạnh vùng (bước 29): mỗi trục neo riêng.
// X: + = cách cạnh trái, − = cách cạnh phải; Y: + = cách cạnh trên, − = cách cạnh dưới. Gõ số âm là đổi neo.
// Lưới 3×3 chọn neo (4 góc, giữa cạnh, tâm), "%" = theo tỷ lệ cạnh vùng (co giãn theo vùng).
import { ReactNode } from 'react';
import NumInput from '../NumInput';
import { Target } from '../icons';
import { RATIO_SCALE } from '../../anchor';
import { Anchor, StepPoint } from '../../types';

type Props = {
  value: StepPoint;
  onChange: (v: StepPoint) => void; // gõ số: đổi vị trí (số âm ở trục start → neo end)
  // Đổi neo (lưới, nút %): có thì gọi thay onChange để nơi dùng tính lại số, giữ nguyên chỗ trên màn hình.
  onAnchor?: (ax: Anchor, ay: Anchor) => void;
  disabled?: boolean;
  pick?: ReactNode; // nút chọn trên màn hình (mặc định: nút mẫu không làm gì – màn #demo)
  trailing?: ReactNode;
};

const AXES: Anchor[] = ['start', 'center', 'end'];
const axisOf = (a: Anchor | undefined): Anchor => a ?? 'start';

// Giá trị hiện trong ô theo neo: start → +n, end → −n, center → lệch so với giữa (có dấu), ratio → %
const shown = (n: number, a: Anchor) => (a === 'end' ? -n : n);
// Số âm (kể cả "-0") → neo phải / dưới, số dương → trái / trên, 0 giữ neo cũ.
const parse = (n: number, a: Anchor): { n: number; a: Anchor } =>
  a === 'ratio' || a === 'center' ? { n, a }
    : n < 0 || Object.is(n, -0) ? { n: -n, a: 'end' } : n > 0 ? { n, a: 'start' } : { n, a };

export function anchorHint({ x, y, ax, ay }: StepPoint) {
  const part = (n: number, a: Anchor, near: string, far: string, mid: string) =>
    a === 'start' ? `cách cạnh ${near} ${n} px` : a === 'end' ? `cách cạnh ${far} ${n} px`
      : a === 'center' ? `${mid} ${n >= 0 ? '+' : ''}${n} px` : `${Number((n / 100).toFixed(2))}% ${near === 'trái' ? 'chiều rộng' : 'chiều cao'}`;
  return `${part(x, axisOf(ax), 'trái', 'phải', 'giữa ngang')} · ${part(y, axisOf(ay), 'trên', 'dưới', 'giữa dọc')}`;
}

export default function AnchorField({ value: v, onChange, onAnchor, disabled, pick, trailing }: Props) {
  const ratio = v.ax === 'ratio';
  const setAxis = (k: 'x' | 'y', n: number) => {
    const r = parse(n, axisOf(k === 'x' ? v.ax : v.ay));
    onChange(k === 'x' ? { ...v, x: r.n, ax: r.a } : { ...v, y: r.n, ay: r.a });
  };
  const anchor = (ax: Anchor, ay: Anchor, demo: Partial<StepPoint> = {}) =>
    onAnchor ? onAnchor(ax, ay) : onChange({ ...v, ax, ay, ...demo });
  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex gap-2">
        {(['x', 'y'] as const).map((k) => {
          const a = axisOf(k === 'x' ? v.ax : v.ay);
          return (
            <NumInput key={k} signed prefix={k.toUpperCase()} unit={a === 'ratio' ? '%' : undefined} scale={a === 'ratio' ? 100 : 1}
              value={shown(v[k], a)} onChange={(n) => !disabled && setAxis(k, n)} className="min-w-0 flex-1" />
          );
        })}
        <AnchorGrid value={v} onChange={(ax, ay) => anchor(ax, ay)} disabled={disabled} />
        {pick ?? (
          <button disabled={disabled} title="Chọn trên màn hình" aria-label="Chọn trên màn hình"
            className="grid h-8 w-8 flex-none place-items-center rounded-control border border-line text-muted enabled:hover:bg-hover enabled:hover:text-ink disabled:opacity-50"><Target /></button>
        )}
        {trailing}
      </div>
      <span className="flex items-center gap-2 text-xs text-muted">
        <span className="min-w-0 flex-1">{anchorHint(v)}</span>
        <button disabled={disabled} onClick={() => (ratio ? anchor('start', 'start') : anchor('ratio', 'ratio', { x: RATIO_SCALE / 2, y: RATIO_SCALE / 2 }))}
          title="Theo tỷ lệ: điểm co giãn theo vùng thay vì giữ khoảng cách tới cạnh"
          className={`rounded-full border px-2 py-px font-mono text-[11px] disabled:opacity-50 ${ratio ? 'border-accent bg-accent-soft text-accent-text' : 'border-line enabled:hover:bg-hover'}`}>%</button>
      </span>
    </div>
  );
}

// Lưới 3×3: chấm đang chọn = neo hiện tại của (X, Y)
export function AnchorGrid({ value: v, onChange, disabled }: { value: StepPoint; onChange: (ax: Anchor, ay: Anchor) => void; disabled?: boolean }) {
  return (
    <div role="radiogroup" aria-label="Neo theo góc vùng" title="Neo theo góc / cạnh vùng – đổi vùng thì điểm giữ khoảng cách tới góc này"
      className={`grid h-8 w-8 flex-none grid-cols-3 place-items-center rounded-control border border-line bg-field p-[3px] ${disabled ? 'opacity-50' : ''}`}>
      {AXES.flatMap((ay) => AXES.map((ax) => {
        const on = axisOf(v.ax) === ax && axisOf(v.ay) === ay;
        return (
          <button key={`${ax}${ay}`} role="radio" aria-checked={on} disabled={disabled} aria-label={`Neo ${ay === 'start' ? 'trên' : ay === 'end' ? 'dưới' : 'giữa'} ${ax === 'start' ? 'trái' : ax === 'end' ? 'phải' : 'giữa'}`}
            onClick={() => onChange(ax, ay)}
            className="grid h-full w-full place-items-center">
            <span className={`block rounded-full ${on ? 'h-[7px] w-[7px] bg-accent' : 'h-[3px] w-[3px] bg-muted/60'}`} />
          </button>
        );
      }))}
    </div>
  );
}
