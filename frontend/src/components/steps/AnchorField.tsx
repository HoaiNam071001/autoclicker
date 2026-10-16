// Ô vị trí có neo theo 4 góc vùng (giao diện cho bước 29): mỗi trục neo riêng.
// X: + = cách cạnh trái, − = cách cạnh phải; Y: + = cách cạnh trên, − = cách cạnh dưới. Gõ số âm là đổi neo.
// Lưới 3×3 chọn neo (4 góc, giữa cạnh, tâm), "%" = theo tỷ lệ cạnh vùng (co giãn theo vùng).
import { ReactNode } from 'react';
import NumInput from '../NumInput';
import { Target } from '../icons';
import { AnchorAxis, AnchoredPoint } from './model';

type Props = { value: AnchoredPoint; onChange: (v: AnchoredPoint) => void; disabled?: boolean; trailing?: ReactNode };

const AXES: AnchorAxis[] = ['start', 'center', 'end'];

// Giá trị hiện trong ô theo neo: start → +n, end → −n, center → lệch so với giữa (có dấu), ratio → %
const shown = (n: number, a: AnchorAxis) => (a === 'end' ? -n : n);
const parse = (n: number, a: AnchorAxis): { n: number; a: AnchorAxis } =>
  a === 'ratio' || a === 'center' ? { n, a } : n < 0 ? { n: -n, a: 'end' } : { n, a: a === 'end' && n > 0 ? 'start' : a };

export function anchorHint({ x, y, ax, ay }: AnchoredPoint) {
  const part = (n: number, a: AnchorAxis, near: string, far: string, mid: string) =>
    a === 'start' ? `cách cạnh ${near} ${n} px` : a === 'end' ? `cách cạnh ${far} ${n} px`
      : a === 'center' ? `${mid} ${n >= 0 ? '+' : ''}${n} px` : `${n}% ${near === 'trái' ? 'chiều rộng' : 'chiều cao'}`;
  return `${part(x, ax, 'trái', 'phải', 'giữa ngang')} · ${part(y, ay, 'trên', 'dưới', 'giữa dọc')}`;
}

export default function AnchorField({ value: v, onChange, disabled, trailing }: Props) {
  const ratio = v.ax === 'ratio';
  const setAxis = (k: 'x' | 'y', n: number) => {
    const a = k === 'x' ? v.ax : v.ay;
    const r = parse(n, a);
    onChange(k === 'x' ? { ...v, x: r.n, ax: r.a } : { ...v, y: r.n, ay: r.a });
  };
  return (
    <div className="flex flex-col gap-1.5">
      <div className="flex gap-2">
        {(['x', 'y'] as const).map((k) => {
          const a = k === 'x' ? v.ax : v.ay;
          return (
            <NumInput key={k} signed prefix={k.toUpperCase()} unit={a === 'ratio' ? '%' : undefined}
              value={shown(v[k], a)} onChange={(n) => !disabled && setAxis(k, n)} className="min-w-0 flex-1" />
          );
        })}
        <AnchorGrid value={v} onChange={onChange} disabled={disabled} />
        <button disabled={disabled} title="Chọn trên màn hình" aria-label="Chọn trên màn hình"
          className="grid h-8 w-8 flex-none place-items-center rounded-control border border-line text-muted enabled:hover:bg-hover enabled:hover:text-ink disabled:opacity-50"><Target /></button>
        {trailing}
      </div>
      <span className="flex items-center gap-2 text-xs text-muted">
        <span className="min-w-0 flex-1">{anchorHint(v)}</span>
        <button disabled={disabled} onClick={() => onChange(ratio ? { ...v, ax: 'start', ay: 'start' } : { ...v, ax: 'ratio', ay: 'ratio', x: 50, y: 50 })}
          title="Theo tỷ lệ: điểm co giãn theo vùng thay vì giữ khoảng cách tới cạnh"
          className={`rounded-full border px-2 py-px font-mono text-[11px] disabled:opacity-50 ${ratio ? 'border-accent bg-accent-soft text-accent-text' : 'border-line enabled:hover:bg-hover'}`}>%</button>
      </span>
    </div>
  );
}

// Lưới 3×3: chấm đang chọn = neo hiện tại của (X, Y)
export function AnchorGrid({ value: v, onChange, disabled }: Omit<Props, 'trailing'>) {
  return (
    <div role="radiogroup" aria-label="Neo theo góc vùng" title={disabled ? 'Neo theo góc vùng – sắp có (bước 29)' : 'Neo theo góc vùng'}
      className={`grid h-8 w-8 flex-none grid-cols-3 place-items-center rounded-control border border-line bg-field p-[3px] ${disabled ? 'opacity-50' : ''}`}>
      {AXES.flatMap((ay) => AXES.map((ax) => {
        const on = v.ax === ax && v.ay === ay;
        return (
          <button key={`${ax}${ay}`} role="radio" aria-checked={on} disabled={disabled} aria-label={`Neo ${ay === 'start' ? 'trên' : ay === 'end' ? 'dưới' : 'giữa'} ${ax === 'start' ? 'trái' : ax === 'end' ? 'phải' : 'giữa'}`}
            onClick={() => onChange({ ...v, ax, ay })}
            className="grid h-full w-full place-items-center">
            <span className={`block rounded-full ${on ? 'h-[7px] w-[7px] bg-accent' : 'h-[3px] w-[3px] bg-muted/60'}`} />
          </button>
        );
      }))}
    </div>
  );
}
