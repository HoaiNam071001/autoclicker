// Xem trước kịch bản: con trỏ giả chạy theo các bước trên overlay, không điều khiển chuột thật.
import { useEffect, useRef, useState } from 'react';
import { boundsOf, endpoints } from '../store';
import { orderedSteps } from '../flow/model';
import { fmtSelectors } from '../flow/web';
import { fmtKey, Point, Rect, SCROLL_ARROW, scrollDelta, ScrollDirection, Step, stepLabel, Script } from '../types';
import ScrollProgress from './simulation/ScrollProgress';

const TRAVEL_MS = 300; // thời gian con trỏ giả bay tới điểm (engine thật dịch chuyển tức thì)
const TYPE_CHAR_MS = 40;
const KEY_SHOW_MS = 500;
const RIPPLE_MS = 700;
const SCROLL_MIN_MS = 700; // mũi tên cuộn chạy ít nhất chừng này để kịp nhìn (engine thật: 30 ms / nấc)
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
  // bước cuộn đang chạy: k 0..1 = đã cuộn bao nhiêu phần quãng dist (px); view: cỡ vùng / màn hình đang cuộn
  scroll: { from: Point; dir: ScrollDirection; dist: number; view: { w: number; h: number }; k: number } | null;
  offset: Point; // trang nền đã cuộn bao nhiêu px (cộng dồn các bước cuộn trong lượt xem trước)
};

const idle: SimState = {
  running: false, pos: null, pressed: false, step: null, label: '', bubble: '', trail: [], ripples: [], scroll: null, offset: { x: 0, y: 0 },
};

// Khoảng trang nền cần vẽ: offset nhỏ / lớn nhất qua mọi bước cuộn đang bật (cuộn lên trước thì trang có phần phía trên).
export function scrollRange(steps: Step[]) {
  const r = { minX: 0, maxX: 0, minY: 0, maxY: 0 };
  let x = 0, y = 0;
  for (const s of steps) {
    if (!s.enabled || s.type !== 'scroll') continue;
    const d = scrollDelta(s);
    x += d.x;
    y += d.y;
    r.minX = Math.min(r.minX, x); r.maxX = Math.max(r.maxX, x);
    r.minY = Math.min(r.minY, y); r.maxY = Math.max(r.maxY, y);
  }
  return r;
}

class Stopped extends Error {}

// Xem trước đi theo thứ tự hiển thị của flow (nhánh chính – cổng success trước); rẽ nhánh theo lựa chọn: bước 20 F6.
export function useSimulation(sc: Pick<Script, 'steps' | 'start' | 'edges'>, region: Rect | null, speed: number) {
  const steps = orderedSteps(sc);
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
    let offset: Point = { x: 0, y: 0 };
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
    // chạy f(k) với k 0 → 1 trong ms (theo tốc độ xem trước)
    const animate = (ms: number, f: (k: number) => void) => new Promise<void>((res, rej) => {
      const t0 = performance.now();
      const frame = (now: number) => {
        if (token.current !== my) return rej(new Stopped());
        const k = Math.min((now - t0) / (ms / speedRef.current), 1);
        f(k);
        if (k < 1) nextFrame(frame);
        else res();
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
    const ends = endpoints(sc, region);
    try {
      for (const [i, s] of steps.entries()) {
        if (!s.enabled) continue;
        alive();
        patch({ step: i + 1, label: `${stepLabel(s)}${s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}`, bubble: '' });
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
          case 'scroll': {
            if (e?.start) await moveTo(e.start, TRAVEL_MS);
            // cuộn tại chuột: từ chỗ con trỏ giả đang đứng (chưa có thì giữa vùng / màn hình)
            const b = boundsOf(region);
            const from = e?.start ?? pos ?? { x: b.x + b.w / 2, y: b.y + b.h / 2 };
            const d = scrollDelta(s);
            const dir = s.direction ?? 'down';
            const dist = Math.abs(d.x || d.y);
            patch({ bubble: '' }); // trang thu nhỏ cạnh con trỏ đã ghi quãng cuộn
            // cuộn dài thì xem trước lâu hơn chút để kịp thấy trang trượt (tối đa ~2,5 s)
            const base = offset;
            await animate(Math.min(Math.max(SCROLL_MIN_MS, (s.amount ?? 0) * 60), 2500), (k) => {
              const e2 = 1 - (1 - k) ** 2;
              offset = { x: base.x + d.x * e2, y: base.y + d.y * e2 };
              patch({ scroll: { from, dir, dist, view: { w: b.w, h: b.h }, k: e2 }, offset });
            });
            await wait(600);
            patch({ scroll: null });
            break;
          }
          case 'open':
            patch({ bubble: `Mở ${s.target}${s.app ? ` bằng ${s.app}` : ''}` });
            await wait(KEY_SHOW_MS);
            break;
          case 'waitImage':
            // không biết trước ảnh hiện lúc nào – chờ tượng trưng (≤ 1 s), ghi rõ ngưỡng / thời gian chờ tối đa
            patch({ bubble: s.findBy === 'text' ? `Chờ chữ “${s.findText}” (tối đa ${(s.timeoutMs ?? 0) / 1000} s)` : `Chờ ảnh (khớp ≥ ${s.threshold}%, tối đa ${(s.timeoutMs ?? 0) / 1000} s)` });
            await wait(Math.min(s.timeoutMs ?? 0, 1000));
            break;
          case 'clickImage':
            patch({ bubble: `Click theo ${s.imageSource === 'search' && s.findBy === 'text' ? `chữ “${s.findText}”` : 'ảnh'}${s.imageSource === 'search' ? ' (tự tìm)' : ''} – vị trí biết khi chạy thật` });
            await wait(KEY_SHOW_MS);
            break;
          case 'webWait': case 'webClick':
            patch({ bubble: `${stepLabel(s)}: ${fmtSelectors(s.selectors)} – vị trí biết khi chạy thật` });
            await wait(KEY_SHOW_MS);
            break;
          case 'waitPixel':
            if (e?.start) await moveTo(e.start, TRAVEL_MS);
            patch({ bubble: `Chờ màu ${s.color} ±${s.tolerance}` });
            await wait(Math.min(s.timeoutMs ?? 0, 800));
            break;
          case 'wait':
            break;
          default: // node điều khiển – xem trước chưa đi theo nhánh (bước 20 F6)
            patch({ bubble: `${stepLabel(s)} – rẽ nhánh lúc chạy thật` });
            await wait(KEY_SHOW_MS);
        }
        if (s.delayAfter > 0) patch((st) => ({ bubble: st.bubble || `Chờ ${s.delayAfter} ms` }));
        await wait(s.delayAfter);
        patch({ bubble: '' });
      }
      patch({ running: false, pressed: false, label: 'Xong 1 vòng', bubble: '', scroll: null });
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
