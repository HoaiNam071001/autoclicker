// Lớp vẽ xem trước trên overlay – chỉ hiển thị (state từ useSimulation).
import type { Point } from '../../types';
import type { SimState } from '../Simulation';
import ScrollProgress from './ScrollProgress';

// Lớp vẽ: vệt đường đi, gợn click, con trỏ giả + bong bóng. toClient: toạ độ màn hình → cửa sổ overlay.
export default function SimulationLayer({ sim, toClient }: { sim: SimState; toClient: (p: Point) => Point }) {
  const cur = sim.pos && toClient(sim.pos);
  return (
    <div className="pointer-events-none absolute inset-0">
      <svg className="absolute inset-0 h-full w-full" aria-hidden>
        {sim.trail.map((t, i) => {
          // đường giống người: vẽ theo các điểm đã đi qua; tắt: đoạn thẳng
          const pts = (t.pts ?? [t.from, t.to]).map(toClient).map((p) => `${p.x},${p.y}`).join(' ');
          return (
            <polyline key={i} points={pts} fill="none"
              stroke={t.drag ? '#F0913A' : '#FF4FA3'} strokeWidth={t.drag ? 4 : 2.5}
              strokeDasharray={t.drag ? undefined : '4 6'} strokeLinecap="round" strokeLinejoin="round"
              opacity={0.35 + 0.65 * ((i + 1) / sim.trail.length)} />
          );
        })}
      </svg>

      {sim.scroll && <ScrollProgress sc={sim.scroll} toClient={toClient} />}

      {sim.ripples.map((r) => {
        const at = toClient(r.at);
        return (
          <span key={r.id} className="absolute -translate-x-1/2 -translate-y-1/2" style={{ left: at.x, top: at.y }}>
            <span className={`block h-12 w-12 rounded-full border-4 border-[#FF4FA3] ${r.hold ? 'sim-hold' : 'sim-ripple'}`}
              style={r.hold ? { animationDuration: `${r.hold}ms` } : undefined} />
          </span>
        );
      })}

      {cur && (
        <div className="absolute" style={{ left: cur.x, top: cur.y }}>
          {/* con trỏ giả: mũi tên hồng viền trắng, đầu mũi tên đúng toạ độ */}
          <svg width="26" height="32" viewBox="0 0 26 32" className={`drop-shadow-[0_2px_4px_rgb(0_0_0/.4)] transition-transform ${sim.pressed ? 'scale-90' : ''}`}>
            <path d="M2 2 L2 26 L8.5 19.5 L13 30 L17.5 28 L13 18 L22 18 Z" fill="#FF4FA3" stroke="#fff" strokeWidth="2" strokeLinejoin="round" />
          </svg>
          {sim.bubble && (
            <span className="absolute left-7 top-7 w-max max-w-[320px] whitespace-pre-wrap rounded-lg bg-hud px-2.5 py-1.5 font-mono text-[13px] text-hud-text shadow">
              {sim.bubble}
            </span>
          )}
        </div>
      )}
    </div>
  );
}
