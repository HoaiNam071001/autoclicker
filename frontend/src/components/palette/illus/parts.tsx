// Phần tử vẽ dùng chung cho hoạt hoạ minh hoạ node (bước 47): mọi cảnh là SVG 240×135 chạy lặp một vòng DUR,
// thời điểm viết theo phần của vòng (0–1). Màu lấy từ token → theo theme sáng / tối. Dùng SMIL (<animate>) để không cần JS.
import { ReactNode, useEffect, useId, useRef } from 'react';

export const DUR = '2.6s';
const loop = { dur: DUR, repeatCount: 'indefinite' } as const;
const c01 = (n: number) => Math.min(Math.max(n, 0), 1);

// Khung cảnh. Người dùng bật "giảm chuyển động" → dừng ở khoảnh khắc đầy đủ nhất.
export function Stage({ children }: { children: ReactNode }) {
  const ref = useRef<SVGSVGElement>(null);
  useEffect(() => {
    if (window.matchMedia?.('(prefers-reduced-motion: reduce)').matches) {
      ref.current?.pauseAnimations?.();
      ref.current?.setCurrentTime?.(1.8);
    }
  }, []);
  return (
    <svg ref={ref} viewBox="0 0 240 135" className="h-full w-full" aria-hidden>
      <rect width="240" height="135" rx="8" className="fill-bg" />
      {children}
    </svg>
  );
}

// Hiện từ `from` đến `to` của vòng (0–1), ngoài đó ẩn.
export function Show({ from = 0, to = 1, children }: { from?: number; to?: number; children: ReactNode }) {
  return (
    <g opacity={from <= 0 ? 1 : 0}>
      <animate attributeName="opacity" values={from <= 0 ? '1;1;0;0' : '0;1;0;0'} keyTimes={`0;${c01(from)};${c01(to)};1`} calcMode="discrete" {...loop} />
      {children}
    </g>
  );
}

// Di chuyển nhóm qua các mốc [thời điểm, x, y]; trước mốc đầu / sau mốc cuối đứng yên.
export type Key3 = [number, number, number];
export function Move({ pts, children }: { pts: Key3[]; children: ReactNode }) {
  const p: Key3[] = [...pts];
  if (p[0][0] > 0) p.unshift([0, p[0][1], p[0][2]]);
  if (p[p.length - 1][0] < 1) p.push([1, p[p.length - 1][1], p[p.length - 1][2]]);
  return (
    <g transform={`translate(${p[0][1]} ${p[0][2]})`}>
      <animateTransform attributeName="transform" type="translate" values={p.map(([, x, y]) => `${x} ${y}`).join(';')}
        keyTimes={p.map((k) => k[0]).join(';')} {...loop} />
      {children}
    </g>
  );
}

// Giá trị thuộc tính đổi theo thời gian: pts = [thời điểm, giá trị].
export function Tween({ attr, pts }: { attr: string; pts: [number, number][] }) {
  const p = [...pts];
  if (p[0][0] > 0) p.unshift([0, p[0][1]]);
  if (p[p.length - 1][0] < 1) p.push([1, p[p.length - 1][1]]);
  return <animate attributeName={attr} values={p.map((k) => k[1]).join(';')} keyTimes={p.map((k) => k[0]).join(';')} {...loop} />;
}

export const Cur = ({ pts }: { pts: Key3[] }) => (
  <Move pts={pts}>
    <path d="M0 0v13l3.6-3 2.4 5 2.2-1-2.4-5H10z" className="fill-ink stroke-surface" strokeWidth=".8" strokeLinejoin="round" />
  </Move>
);

// Gợn tròn lan ra khi bấm lúc t.
export const Ripple = ({ x, y, t, r = 13 }: { x: number; y: number; t: number; r?: number }) => (
  <circle cx={x} cy={y} r="0" className="fill-none stroke-accent" strokeWidth="2">
    <Tween attr="r" pts={[[0, 0], [t, 0], [t + 0.12, r], [1, r]]} />
    <Tween attr="opacity" pts={[[0, 0], [t, 0.9], [t + 0.12, 0], [1, 0]]} />
  </circle>
);

export const Win = ({ x = 16, y = 12, w = 208, h = 111, bar = true }: { x?: number; y?: number; w?: number; h?: number; bar?: boolean }) => (
  <g>
    <rect x={x} y={y} width={w} height={h} rx="6" className="fill-surface stroke-line-strong" />
    {bar && <>
      <path d={`M${x} ${y + 12}V${y + 6}a6 6 0 0 1 6-6h${w - 12}a6 6 0 0 1 6 6v6z`} className="fill-chrome" />
      {[0, 1, 2].map((i) => <circle key={i} cx={x + 9 + i * 8} cy={y + 6} r="2" className="fill-line-strong" />)}
    </>}
  </g>
);

// Khung trình duyệt: thanh địa chỉ + nội dung.
export const Browser = ({ url = 'https://…', children }: { url?: string; children?: ReactNode }) => (
  <g>
    <Win />
    <rect x="28" y="22" width="184" height="11" rx="5.5" className="fill-field" />
    <text x="36" y="30.2" fontSize="7" className="fill-muted font-mono">{url}</text>
    {children}
  </g>
);

export const Btn = ({ x, y, w = 56, h = 20, label = '', hot = false }: { x: number; y: number; w?: number; h?: number; label?: string; hot?: boolean }) => (
  <g>
    <rect x={x} y={y} width={w} height={h} rx="5" className={hot ? 'fill-accent-soft stroke-accent' : 'fill-field stroke-line-strong'} />
    {label && <text x={x + w / 2} y={y + h / 2 + 3} fontSize="8.5" textAnchor="middle" className="fill-ink font-semibold">{label}</text>}
  </g>
);

export const Ln = ({ x, y, w, cls = 'fill-line-strong' }: { x: number; y: number; w: number; cls?: string }) => <rect x={x} y={y} width={w} height="4" rx="2" className={cls} />;
export const Lines = ({ x, y, ws, gap = 9 }: { x: number; y: number; ws: number[]; gap?: number }) => <>{ws.map((w, i) => <Ln key={i} x={x} y={y + i * gap} w={w} />)}</>;

export const T = ({ x, y, children, s = 9, anchor = 'start', cls = 'fill-ink' }: { x: number; y: number; children: ReactNode; s?: number; anchor?: 'start' | 'middle' | 'end'; cls?: string }) => (
  <text x={x} y={y} fontSize={s} textAnchor={anchor} className={`${cls} font-semibold`}>{children}</text>
);

// Khung nhấn màu chủ đạo quanh thứ node tìm thấy.
export const Ring = ({ x, y, w, h, from = 0.3, to = 1 }: { x: number; y: number; w: number; h: number; from?: number; to?: number }) => (
  <Show from={from} to={to}><rect x={x} y={y} width={w} height={h} rx="4" className="fill-accent-soft stroke-accent" strokeWidth="1.6" strokeDasharray="4 2" /></Show>
);

export const Tick = ({ x, y, from = 0.6 }: { x: number; y: number; from?: number }) => (
  <Show from={from}><g transform={`translate(${x} ${y})`}>
    <circle r="8" className="fill-ok" /><path d="M-3.5 0l2.5 2.8 4.5-5.6" fill="none" className="stroke-on-accent" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </g></Show>
);
export const Cross = ({ x, y, from = 0.6 }: { x: number; y: number; from?: number }) => (
  <Show from={from}><g transform={`translate(${x} ${y})`}>
    <circle r="8" className="fill-danger" /><path d="M-3 -3l6 6M3 -3l-6 6" className="stroke-on-accent" strokeWidth="2" strokeLinecap="round" />
  </g></Show>
);

// Nét vẽ dần từ `from` đến `to`.
export const Draw = ({ d, from = 0.1, to = 0.8, cls = 'stroke-accent' }: { d: string; from?: number; to?: number; cls?: string }) => (
  <path d={d} pathLength="100" fill="none" className={cls} strokeWidth="2.4" strokeLinecap="round" strokeLinejoin="round" strokeDasharray="100" strokeDashoffset="100">
    <Tween attr="stroke-dashoffset" pts={[[0, 100], [from, 100], [to, 0], [1, 0]]} />
  </path>
);

// Chữ hiện dần từng ký tự (cắt bằng clipPath).
export function Typed({ x, y, w, text, from = 0.25, to = 0.65, s = 10 }: { x: number; y: number; w: number; text: string; from?: number; to?: number; s?: number }) {
  const id = useId().replace(/:/g, '');
  return (
    <g>
      <clipPath id={id}><rect x={x} y={y - s} height={s * 1.6} width="0"><Tween attr="width" pts={[[0, 0], [from, 0], [to, w], [1, w]]} /></rect></clipPath>
      <text x={x} y={y} fontSize={s} clipPath={`url(#${id})`} className="fill-ink font-mono">{text}</text>
    </g>
  );
}

export function Clip({ x, y, w, h, children }: { x: number; y: number; w: number; h: number; children: ReactNode }) {
  const id = useId().replace(/:/g, '');
  return <g><clipPath id={id}><rect x={x} y={y} width={w} height={h} rx="3" /></clipPath><g clipPath={`url(#${id})`}>{children}</g></g>;
}

export const Clock = ({ x, y, r = 22, spin = true }: { x: number; y: number; r?: number; spin?: boolean }) => (
  <g>
    <circle cx={x} cy={y} r={r} className="fill-surface stroke-line-strong" strokeWidth="2" />
    <line x1={x} y1={y} x2={x} y2={y - r * 0.55} className="stroke-ink" strokeWidth="2" strokeLinecap="round" />
    <line x1={x} y1={y} x2={x} y2={y - r * 0.8} className="stroke-accent" strokeWidth="2" strokeLinecap="round">
      {spin && <animateTransform attributeName="transform" type="rotate" values={`0 ${x} ${y};360 ${x} ${y}`} dur="1.4s" repeatCount="indefinite" />}
    </line>
  </g>
);

export const Spinner = ({ x, y, r = 10 }: { x: number; y: number; r?: number }) => (
  <circle cx={x} cy={y} r={r} fill="none" className="stroke-accent" strokeWidth="3" strokeLinecap="round" strokeDasharray={`${r * 2.2} ${r * 4}`}>
    <animateTransform attributeName="transform" type="rotate" values={`0 ${x} ${y};360 ${x} ${y}`} dur="0.9s" repeatCount="indefinite" />
  </circle>
);

export const Arr = ({ x1, y1, x2, y2, cls = 'stroke-muted' }: { x1: number; y1: number; x2: number; y2: number; cls?: string }) => {
  const a = Math.atan2(y2 - y1, x2 - x1);
  const h = (d: number) => `${x2 - 6 * Math.cos(a + d)},${y2 - 6 * Math.sin(a + d)}`;
  return <g className={cls} fill="none" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><path d={`M${x1} ${y1}L${x2} ${y2}`} /><path d={`M${h(0.5)}L${x2} ${y2}L${h(-0.5)}`} /></g>;
};

// Chip chữ nhỏ (biến, tên mẫu…).
export const Chip = ({ x, y, w, label, cls = 'fill-accent-soft stroke-accent' }: { x: number; y: number; w: number; label: string; cls?: string }) => (
  <g><rect x={x} y={y} width={w} height="14" rx="7" className={cls} /><text x={x + w / 2} y={y + 10} fontSize="7.5" textAnchor="middle" className="fill-ink font-mono">{label}</text></g>
);

// Phím trên bàn phím: sáng lên từ `from` đến `to`.
export const KeyCap = ({ x, y, w = 30, label, from, to = 0.8 }: { x: number; y: number; w?: number; label: string; from: number; to?: number }) => (
  <g>
    <rect x={x} y={y} width={w} height="22" rx="5" className="fill-field stroke-line-strong" />
    <Show from={from} to={to}><rect x={x} y={y + 2} width={w} height="20" rx="5" className="fill-accent" /></Show>
    <text x={x + w / 2} y={y + 14} fontSize="8.5" textAnchor="middle" className="fill-ink font-semibold">{label}</text>
    <Show from={from} to={to}><text x={x + w / 2} y={y + 16} fontSize="8.5" textAnchor="middle" className="fill-on-accent font-semibold">{label}</text></Show>
  </g>
);

// Ngôi sao / biểu tượng mẫu để "tìm".
export const Gem = ({ x, y, s = 1 }: { x: number; y: number; s?: number }) => (
  <path transform={`translate(${x} ${y}) scale(${s})`} d="M0-9l2.6 6 6.4.6-4.9 4.2 1.5 6.3L0 4.8l-5.6 3.3 1.5-6.3L-9-2.4l6.4-.6z" className="fill-warn" />
);

export type Scenes = Partial<Record<import('../../../flow/nodeInfo').InfoKey, () => ReactNode>>;
