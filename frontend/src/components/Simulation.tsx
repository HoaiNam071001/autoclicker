// Xem trước kịch bản: con trỏ giả chạy theo các bước trên overlay, không điều khiển chuột thật.
import { useEffect, useRef, useState } from 'react';
import { boundsOf, endpoints, refOf } from '../store';
import { orderedSteps } from '../flow/model';
import { fmtSelectors } from '../flow/web';
import { fmtShot } from '../flow/shot';
import { fmtClip } from '../flow/clipboard';
import { fmtHTTP } from '../flow/http';
import { fmtFile } from '../flow/file';
import { fmtSystem } from '../flow/system';
import { fmtKey } from '../keys';
import { effectiveHuman } from '../flow/human';
import { curveAt, HUMAN_PROFILES, humanLegs } from '../overlay/humanPath';
import { Point, Rect, SCROLL_ARROW, scrollDelta, ScrollDirection, Step, stepLabel, Script } from '../types';

const TRAVEL_MS = 300; // thời gian con trỏ giả bay tới điểm khi tắt giống người (engine thật dịch chuyển tức thì)
const TYPE_CHAR_MS = 40;
const KEY_SHOW_MS = 500;
const RIPPLE_MS = 700;
const SCROLL_MIN_MS = 700; // mũi tên cuộn chạy ít nhất chừng này để kịp nhìn (engine thật: 30 ms / nấc)
const MAX_TRAIL = 60;
// Khung hình bằng setTimeout thay vì requestAnimationFrame: rAF có thể không chạy khi trình duyệt
// không vẽ khung hình (đã gặp ở Chrome headless) – setTimeout luôn chạy, ~60 khung/giây là đủ mượt.
const FRAME_MS = 16;
const nextFrame = (cb: (now: number) => void) => setTimeout(() => cb(performance.now()), FRAME_MS);

type Segment = { from: Point; to: Point; drag: boolean; pts?: Point[] }; // pts: các điểm đường cong (giống người)
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
export function useSimulation(sc: Pick<Script, 'steps' | 'start' | 'edges' | 'human'>, region: Rect | null, speed: number) {
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
    // tới p như chuột thật: node có mức giống người → đi theo đường cong + thời gian như engine (humanLegs),
    // tắt → bay thẳng TRAVEL_MS. dragMs: đang kéo thả, đi đúng chừng này
    const travel = async (p: Point, s: Step, dragMs?: number) => {
      const drag = dragMs !== undefined;
      const prof = HUMAN_PROFILES[effectiveHuman(sc.human, s.human)];
      if (!prof || !pos) return moveTo(p, drag ? Math.max(dragMs, 1) : TRAVEL_MS, drag);
      for (const l of humanLegs(pos, p, prof, drag ? { durMs: Math.max(dragMs, 1) } : {})) {
        await wait(l.pauseMs); // khựng trước khi chỉnh lại sau khi đi lố
        const pts: Point[] = [l.from];
        if (l.ms > 0) {
          await animate(l.ms, (k) => {
            pos = curveAt(l.from, l.c1, l.c2, l.to, k);
            pts.push(pos);
            patch({ pos });
          });
        }
        pos = l.to;
        patch((st) => ({ pos: l.to, trail: [...st.trail, { from: l.from, to: l.to, drag, pts: [...pts, l.to] }].slice(-MAX_TRAIL) }));
      }
    };
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
              if (e?.start) await travel(e.start, s);
              ripple(pos!);
              if (n < s.clicks - 1) await wait(s.clickInterval);
            }
            break;
          case 'double':
            if (e?.start) await travel(e.start, s);
            ripple(pos!);
            await wait(100);
            ripple(pos!);
            break;
          case 'hold':
            if (e?.start) await travel(e.start, s);
            patch({ pressed: true, bubble: `Giữ ${s.holdMs} ms` });
            ripple(pos!, s.holdMs);
            await wait(s.holdMs);
            patch({ pressed: false, bubble: '' });
            break;
          case 'drag':
            if (!e?.start || !e.end) break;
            await travel(e.start, s);
            patch({ pressed: true });
            await travel(e.end, s, s.dragMs);
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
            if (e?.start) await travel(e.start, s);
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
          case 'webWait': case 'webClick': case 'webType': case 'webRead': case 'webScroll':
            patch({ bubble: `${stepLabel(s)}: ${fmtSelectors(s.selectors)} – vị trí biết khi chạy thật` });
            await wait(KEY_SHOW_MS);
            break;
          case 'waitPixel':
            if (e?.start) await travel(e.start, s);
            patch({ bubble: `Chờ màu ${s.color} ±${s.tolerance}` });
            await wait(Math.min(s.timeoutMs ?? 0, 800));
            break;
          case 'system':
            patch({ bubble: `🔊 ${fmtSystem(s)}` });
            await wait(KEY_SHOW_MS);
            break;
          case 'file':
            patch({ bubble: `📄 ${fmtFile(s)}` });
            await wait(KEY_SHOW_MS);
            break;
          case 'http':
            patch({ bubble: `🌐 ${fmtHTTP(s)} – kết quả biết khi chạy thật` });
            await wait(KEY_SHOW_MS);
            break;
          case 'clipboard':
            patch({ bubble: `📋 Clipboard: ${fmtClip(s)}` });
            await wait(KEY_SHOW_MS);
            break;
          case 'screenshot':
            patch({ bubble: `📷 Chụp ${fmtShot(s, refOf)}` });
            await wait(KEY_SHOW_MS);
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
