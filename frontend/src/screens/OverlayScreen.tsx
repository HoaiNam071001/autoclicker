// Overlay chọn điểm – chạy trong cửa sổ trong suốt, full screen, always-on-top
import { MouseEvent, useEffect, useState } from 'react';
import { MousePos } from '../../wailsjs/go/main/App';
import { inWails, setWindowMode } from '../windowMode';
import { Point } from '../types';

const pts = [
  { label: '1', x: 380, y: 260 }, { label: '2', x: 600, y: 260 }, { label: '3', x: 760, y: 450 },
  { label: '4A', x: 330, y: 600, drag: true }, { label: '4B', x: 880, y: 600, drag: true }, { label: '5', x: 1010, y: 520 },
];
const region = { x: 180, y: 150, w: 1040, h: 600 };

export default function OverlayScreen() {
  const [cursor, setCursor] = useState<Point>({ x: 1100, y: 360 });
  // Spike bước 5: so toạ độ trong cửa sổ (client) với toạ độ màn hình thật (robotgo) để đo độ lệch.
  const [probe, setProbe] = useState<{ client: Point; screen: Point } | null>(null);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && setWindowMode('main');
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  const onClick = async (e: MouseEvent) => {
    if ((e.target as HTMLElement).closest('button')) return;
    const client = { x: e.clientX, y: e.clientY };
    const screen = inWails ? await MousePos() : client;
    setProbe({ client, screen });
  };

  return (
    <div className="fixed inset-0 select-none bg-[rgba(20,24,32,.45)]"
      onMouseMove={(e) => setCursor({ x: e.clientX, y: e.clientY })} onClick={onClick}>
      <div className="absolute border-2 border-dashed border-white bg-white/5"
        style={{ left: region.x, top: region.y, width: region.w, height: region.h }} />
      <span className="absolute rounded bg-accent px-2 py-0.5 font-mono text-[13px] text-white"
        style={{ left: region.x, top: region.y - 30 }}>Vùng {region.w} × {region.h}</span>
      {[[region.x, region.y], [region.x + region.w, region.y + region.h]].map(([x, y], i) => (
        <span key={i} className="absolute h-3.5 w-3.5 -translate-x-1/2 -translate-y-1/2 border-2 border-accent bg-white" style={{ left: x, top: y }} />
      ))}

      <svg className="absolute inset-0 h-full w-full" aria-hidden>
        <defs>
          <marker id="ah" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
            <path d="M0 0 10 5 0 10z" fill="#F0913A" />
          </marker>
        </defs>
        <polyline points="380,260 600,260 760,450 330,600" fill="none" stroke="#fff" strokeWidth={2} strokeDasharray="6 6" opacity={0.8} />
        <line x1="330" y1="600" x2="880" y2="600" stroke="#F0913A" strokeWidth={3} markerEnd="url(#ah)" />
        <polyline points="880,600 1010,520" fill="none" stroke="#fff" strokeWidth={2} strokeDasharray="6 6" opacity={0.8} />
      </svg>

      {pts.map((p) => (
        <span key={p.label} style={{ left: p.x, top: p.y }}
          className={`absolute grid h-[34px] w-[34px] -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full border-[3px] border-white text-[13px] font-bold text-white shadow-[0_2px_8px_rgba(0,0,0,.35)] ${p.drag ? 'bg-drag-marker' : 'bg-accent'}`}>
          {p.label}
        </span>
      ))}

      {/* Crosshair + toạ độ theo chuột */}
      <div className="pointer-events-none absolute" style={{ left: cursor.x, top: cursor.y }}>
        <div className="absolute -top-10 h-20 w-px bg-white" />
        <div className="absolute -left-10 h-px w-20 bg-white" />
        <span className="absolute left-3.5 top-3 whitespace-nowrap rounded bg-ink px-2 py-1 font-mono text-[13px] text-white">X {cursor.x}  Y {cursor.y}</span>
      </div>

      <div className="absolute left-1/2 top-6 flex -translate-x-1/2 items-center gap-1.5 rounded-xl bg-white p-1.5 shadow-[0_6px_24px_rgba(0,0,0,.25)]">
        <button className="h-11 rounded-lg bg-accent px-3.5 font-semibold text-white">Đặt điểm</button>
        <button className="h-11 rounded-lg px-3.5">Kéo thả A → B</button>
        <button className="h-11 rounded-lg px-3.5">Chọn vùng</button>
        <span className="mx-1.5 h-7 w-px bg-line" />
        <span className="flex items-center gap-1.5 px-1.5 text-[13px] text-muted">
          <kbd className="rounded border border-line px-1.5 font-mono text-ink">F6</kbd>lưu vị trí chuột
        </span>
        <button onClick={() => setWindowMode('main')} className="h-11 rounded-lg bg-ink px-4 font-semibold text-white">Xong <span className="font-mono text-xs opacity-70">Esc</span></button>
      </div>

      <div className="absolute left-6 top-6 rounded-lg bg-ink/90 px-3 py-2 font-mono text-[13px] text-white">
        {probe ? (
          <>
            <div>Cửa sổ  X {probe.client.x}  Y {probe.client.y}</div>
            <div>Màn hình X {probe.screen.x}  Y {probe.screen.y}</div>
            <div className="text-accent-light">Lệch    X {probe.screen.x - probe.client.x}  Y {probe.screen.y - probe.client.y}</div>
          </>
        ) : 'Click vào vùng tối để đo toạ độ'}
      </div>
    </div>
  );
}
