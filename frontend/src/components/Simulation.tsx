// Xem trước kịch bản: con trỏ giả chạy theo các bước trên overlay, không điều khiển chuột thật.
import { useEffect, useRef, useState } from 'react';
import { endpoints } from '../scriptStore';
import { fmtKey, Point, Rect, Step, STEP_LABEL } from '../types';

const TRAVEL_MS = 300; // thời gian con trỏ giả bay tới điểm (engine thật dịch chuyển tức thì)
const TYPE_CHAR_MS = 40;
const KEY_SHOW_MS = 500;
const RIPPLE_MS = 700;
const MAX_TRAIL = 60;
// Khung hình bằng setTimeout thay vì requestAnimationFrame: rAF có thể không chạy khi trình duyệt
// không vẽ khung hình (đã gặp ở Chrome headless) – setTimeout luôn chạy, ~60 khung/giây là đủ mượt.
const FRAME_MS = 16;
const nextFrame = (cb: (now: number) => void) => setTimeout(() => cb(performance.now()), FRAME_MS);

type Segment = { from: Point; to: Point; drag: boolean };
type Ripple = { id: number; at: Point; hold?: number }; // hold: ms giữ chuột (vòng lớn dần)

export type SimState = {
  running: boolean;
  pos: Point | null; // toạ độ màn hình (robotgo) của con trỏ giả
  pressed: boolean;
  step: number | null; // số thứ tự bước đang chạy (từ 1)
  label: string;
  bubble: string; // chữ cạnh con trỏ: chờ / gõ / phím
  trail: Segment[];
  ripples: Ripple[];
};

const idle: SimState = { running: false, pos: null, pressed: false, step: null, label: '', bubble: '', trail: [], ripples: [] };

class Stopped extends Error {}

export function useSimulation(steps: Step[], region: Rect | null, speed: number) {
  const [sim, setSim] = useState<SimState>(idle);
  const token = useRef(0); // tăng lên để huỷ lượt chạy đang dở
  const speedRef = useRef(speed);
  speedRef.current = speed;

  useEffect(() => () => void token.current++, []); // rời overlay → dừng

  const stop = () => {
    token.current++;
    setSim((s) => ({ ...s, running: false, pressed: false, bubble: '', label: s.step ? 'Đã dừng' : '' }));
  };

  const start = async () => {
    const my = ++token.current;
    const alive = () => {
      if (token.current !== my) throw new Stopped();
    };
    const patch = (p: Partial<SimState> | ((s: SimState) => Partial<SimState>)) =>
      setSim((s) => ({ ...s, ...(typeof p === 'function' ? p(s) : p) }));
    // chờ ms (theo tốc độ xem trước), dừng được giữa chừng
    const wait = (ms: number) => new Promise<void>((res, rej) => {
      const t = setTimeout(() => (token.current === my ? res() : rej(new Stopped())), ms / speedRef.current);
      if (token.current !== my) clearTimeout(t);
    });
    let pos: Point | null = null;
    // bay / kéo con trỏ tới p trong ms, để lại vệt
    const moveTo = (p: Point, ms: number, drag = false) => new Promise<void>((res, rej) => {
      const from = pos ?? p;
      const t0 = performance.now();
      const frame = (now: number) => {
        if (token.current !== my) return rej(new Stopped());
        const k = Math.min(Math.max((now - t0) / (ms / speedRef.current), 0), 1);
        const e = drag ? k : 1 - (1 - k) ** 3; // kéo thả đều, bay thì chậm dần
        pos = { x: from.x + (p.x - from.x) * e, y: from.y + (p.y - from.y) * e };
        patch({ pos });
        if (k < 1) return nextFrame(frame);
        if (from.x !== p.x || from.y !== p.y) {
          patch((s) => ({ trail: [...s.trail, { from, to: p, drag }].slice(-MAX_TRAIL) }));
        }
        res();
      };
      nextFrame(frame);
    });
    let rid = 0;
    const ripple = (at: Point, hold?: number) => {
      const r = { id: ++rid, at, hold };
      patch((s) => ({ ripples: [...s.ripples, r] }));
      setTimeout(() => patch((s) => ({ ripples: s.ripples.filter((x) => x.id !== r.id) })),
        ((hold ?? 0) + RIPPLE_MS) / speedRef.current);
    };

    setSim({ ...idle, running: true });
    const ends = endpoints(steps, region);
    try {
      for (const [i, s] of steps.entries()) {
        if (!s.enabled) continue;
        alive();
        patch({ step: i + 1, label: `${STEP_LABEL[s.type]}${s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}`, bubble: '' });
        const e = ends.get(s.id);
        switch (s.type) {
          case 'click':
            for (let n = 0; n < Math.max(s.clicks, 1); n++) {
              if (e?.start) await moveTo(e.start, TRAVEL_MS);
              ripple(pos!);
              if (n < s.clicks - 1) await wait(s.clickInterval);
            }
            break;
          case 'double':
            if (e?.start) await moveTo(e.start, TRAVEL_MS);
            ripple(pos!);
            await wait(100);
            ripple(pos!);
            break;
          case 'hold':
            if (e?.start) await moveTo(e.start, TRAVEL_MS);
            patch({ pressed: true, bubble: `Giữ ${s.holdMs} ms` });
            ripple(pos!, s.holdMs);
            await wait(s.holdMs);
            patch({ pressed: false, bubble: '' });
            break;
          case 'drag':
            if (!e?.start || !e.end) break;
            await moveTo(e.start, TRAVEL_MS);
            patch({ pressed: true });
            await moveTo(e.end, Math.max(s.dragMs, 1), true);
            patch({ pressed: false });
            break;
          case 'type':
            for (let n = 1; n <= [...s.text].length; n++) {
              patch({ bubble: `⌨ ${[...s.text].slice(0, n).join('')}` });
              await wait(TYPE_CHAR_MS);
            }
            break;
          case 'key':
            patch({ bubble: `⌨ ${fmtKey(s)}` });
            await wait(KEY_SHOW_MS);
            break;
          case 'wait':
            break;
        }
        if (s.delayAfter > 0) patch((st) => ({ bubble: st.bubble || `Chờ ${s.delayAfter} ms` }));
        await wait(s.delayAfter);
        patch({ bubble: '' });
      }
      patch({ running: false, pressed: false, label: 'Xong 1 vòng', bubble: '' });
    } catch (err) {
      if (!(err instanceof Stopped)) throw err;
    }
  };

  return { sim, start, stop };
}

// Lớp vẽ: vệt đường đi, gợn click, con trỏ giả + bong bóng. toClient: toạ độ màn hình → cửa sổ overlay.
export function SimulationLayer({ sim, toClient }: { sim: SimState; toClient: (p: Point) => Point }) {
  const cur = sim.pos && toClient(sim.pos);
  return (
    <div className="pointer-events-none absolute inset-0">
      <svg className="absolute inset-0 h-full w-full" aria-hidden>
        {sim.trail.map((t, i) => {
          const a = toClient(t.from);
          const b = toClient(t.to);
          return (
            <line key={i} x1={a.x} y1={a.y} x2={b.x} y2={b.y}
              stroke={t.drag ? '#F0913A' : '#FF4FA3'} strokeWidth={t.drag ? 4 : 2.5}
              strokeDasharray={t.drag ? undefined : '4 6'} strokeLinecap="round"
              opacity={0.35 + 0.65 * ((i + 1) / sim.trail.length)} />
          );
        })}
      </svg>

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
          <svg width="26" height="32" viewBox="0 0 26 32" className={`drop-shadow-[0_2px_4px_rgba(0,0,0,.4)] transition-transform ${sim.pressed ? 'scale-90' : ''}`}>
            <path d="M2 2 L2 26 L8.5 19.5 L13 30 L17.5 28 L13 18 L22 18 Z" fill="#FF4FA3" stroke="#fff" strokeWidth="2" strokeLinejoin="round" />
          </svg>
          {sim.bubble && (
            <span className="absolute left-7 top-7 w-max max-w-[320px] whitespace-pre-wrap rounded-lg bg-ink px-2.5 py-1.5 font-mono text-[13px] text-white shadow">
              {sim.bubble}
            </span>
          )}
        </div>
      )}
    </div>
  );
}
