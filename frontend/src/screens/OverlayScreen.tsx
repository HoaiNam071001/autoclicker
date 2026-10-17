// Overlay chọn vùng / điểm – cửa sổ trong suốt phủ màn hình, always-on-top
import { MouseEvent, ReactNode, useEffect, useRef, useState } from 'react';
import { MousePos } from '../../wailsjs/go/main/App';
import {
  boundsOf, cancelDragA, cancelPick, endpoints, movePoint, pickingStep, placePoint, placingType, selectStep, setDragMode, setFocusLock,
  setRegion, startPick, useEditor,
} from '../scriptStore';
import { Point, Rect, STEP_LABEL, StepPoint } from '../types';
import { fmtAxis } from '../anchor';
import { inWails, OverlayTool, overlayFocus, overlayTool, setWindowMode } from '../windowMode';
import { addSample, currentMapping, isMeasured, resetSamples, toClientPt, toScreenPt } from '../coordMap';
import { SimulationLayer, useSimulation } from '../components/Simulation';
import { Eye, Play, Region, SelectIcon, StepIcon, Stop, Target } from '../components/icons';

const SPEEDS = [1, 2, 4];

const CALIBRATE_MS = 150; // chuột đứng yên bao lâu thì lấy mẫu đo toạ độ (coordMap.ts)
const MIN_REGION = 10; // px – kéo nhỏ hơn coi như click nhầm
const MIN_MOVE = 3; // px – kéo marker ít hơn coi như click (chọn bước)

// Marker đang được kéo: bước, điểm thứ mấy, vị trí bấm xuống (toạ độ cửa sổ), đã kéo đủ xa chưa.
type Moving = { id: string; index: number; from: Point; moved: boolean };

const TOOLS: { id: OverlayTool; label: string; soon?: string; icon: ReactNode }[] = [
  { id: 'point', label: 'Đặt điểm', icon: <StepIcon type="click" /> },
  { id: 'drag', label: 'Kéo thả A → B', icon: <StepIcon type="drag" /> },
  { id: 'region', label: 'Chọn vùng', icon: <Region /> },
  { id: 'select', label: 'Chọn nhiều', soon: '30', icon: <SelectIcon /> },
  { id: 'capture', label: 'Khoanh ảnh', soon: '25', icon: <StepIcon type="waitImage" /> },
];
const toolBtn = 'flex h-11 items-center gap-2 rounded-[7px] px-3.5 font-semibold';

const rectOf = (a: Point, b: Point): Rect => ({
  x: Math.min(a.x, b.x), y: Math.min(a.y, b.y), w: Math.abs(a.x - b.x), h: Math.abs(a.y - b.y),
});

export default function OverlayScreen() {
  const { script: { region, steps }, selected, notice, dragA } = useEditor();
  const [picking] = useState(pickingStep); // id bước đang chọn lại điểm (nút ⌖), null = thêm bước mới
  const [tool, setTool] = useState<OverlayTool>(overlayTool);
  const [cursor, setCursor] = useState<Point>({ x: -100, y: -100 });
  const [map, setMap] = useState(currentMapping);
  const [measured, setMeasured] = useState(isMeasured);
  const [dragFrom, setDragFromState] = useState<Point | null>(null); // toạ độ cửa sổ
  const dragRef = useRef<Point | null>(null);
  const setDragFrom = (p: Point | null) => {
    dragRef.current = p;
    setDragFromState(p);
  };
  const [moving, setMovingState] = useState<Moving | null>(null);
  const movingRef = useRef<Moving | null>(null);
  const setMoving = (m: Moving | null) => {
    movingRef.current = m;
    setMovingState(m);
  };
  const preview = tool === 'preview';
  const [speed, setSpeed] = useState(1);
  const { sim, start: startSim, stop: stopSim } = useSimulation(steps, region, speed);
  const simRef = useRef({ running: false, stop: stopSim });
  simRef.current = { running: sim.running, stop: stopSim };
  // Xem riêng một bước (nút 👁 trong form): bước khác + vùng mờ đi, chỉ sửa được bước này.
  const [focus, setFocus] = useState(overlayFocus);
  const [repicking, setRepickingState] = useState(false); // đang "Chọn lại trên màn hình" trong chế độ xem riêng
  const repickRef = useRef(false);
  const setRepicking = (on: boolean) => {
    repickRef.current = on;
    setRepickingState(on);
  };
  const [alt, setAlt] = useState(false); // giữ Alt: hiện đường neo của mọi điểm
  const skipClick = useRef(false); // click sinh ra sau khi thả marker không được tính là đặt điểm
  const cursorRef = useRef(cursor);
  const timer = useRef<number | undefined>(undefined);

  useEffect(() => {
    resetSamples(); // overlay có thể đang mở trên màn khác lần trước
    const onKey = (e: KeyboardEvent) => {
      if (e.key !== 'Escape') return;
      // Đang kéo vùng / marker thì Esc chỉ huỷ thao tác đó.
      if (simRef.current.running) simRef.current.stop(); // xem trước đang chạy → Esc chỉ dừng
      else if (movingRef.current) setMoving(null);
      else if (dragRef.current) setDragFrom(null);
      else if (cancelDragA()) return;
      else if (repickRef.current) {
        cancelPick();
        setRepicking(false);
      } else setWindowMode('main');
    };
    const onAlt = (e: KeyboardEvent) => setAlt(e.altKey);
    window.addEventListener('keydown', onKey);
    window.addEventListener('keydown', onAlt);
    window.addEventListener('keyup', onAlt);
    return () => {
      window.removeEventListener('keydown', onKey);
      window.removeEventListener('keydown', onAlt);
      window.removeEventListener('keyup', onAlt);
      clearTimeout(timer.current);
      // Thoát overlay (Esc / Xong) mà chưa chọn lại điểm → huỷ. Chỉ khi đã rời overlay thật: StrictMode (wails dev)
      // gỡ rồi gắn lại màn ngay khi mở – huỷ lúc đó sẽ mất loại bước (Double click / Giữ chuột) và chỗ chèn từ Ctrl+K.
      if (location.hash !== '#overlay') cancelPick();
    };
  }, []);

  useEffect(() => setDragMode(tool === 'drag'), [tool]);
  // Chọn lại xong bằng F6 (không qua onClick) → tắt trạng thái "đang chọn lại".
  useEffect(() => {
    if (repickRef.current && pickingStep() === null) setRepicking(false);
  }, [steps, dragA]);
  useEffect(() => {
    setFocusLock(!!focus);
    if (focus) selectStep(focus);
    return () => setFocusLock(false);
  }, [focus]);

  const toClient = (p: Point): Point => toClientPt(map, p);
  const toScreenFromClient = (p: Point): Point => toScreenPt(map, p);
  const onUI = (e: MouseEvent) => !!(e.target as HTMLElement).closest('[data-ui]');

  const onMove = (e: MouseEvent) => {
    const c = { x: e.clientX, y: e.clientY };
    setCursor(c);
    cursorRef.current = c;
    const m = movingRef.current;
    if (m && !m.moved && Math.hypot(c.x - m.from.x, c.y - m.from.y) >= MIN_MOVE) setMoving({ ...m, moved: true });
    if (!inWails) return;
    clearTimeout(timer.current);
    timer.current = window.setTimeout(async () => {
      const s = await MousePos();
      const now = cursorRef.current;
      if (now.x !== c.x || now.y !== c.y) return; // chuột đã đi tiếp → bỏ mẫu này
      setMap(addSample(c, s));
      setMeasured(true);
    }, CALIBRATE_MS);
  };

  const onDown = (e: MouseEvent) => {
    skipClick.current = false; // bấm ngoài marker (marker chặn sự kiện này)
    if (onUI(e) || e.button !== 0 || tool !== 'region' || focus) return;
    setDragFrom({ x: e.clientX, y: e.clientY });
  };

  const onMarkerDown = (e: MouseEvent, id: string, index: number) => {
    if (e.button !== 0) return;
    e.stopPropagation(); // không bắt đầu kéo vùng
    skipClick.current = true;
    setMoving({ id, index, from: { x: e.clientX, y: e.clientY }, moved: false });
  };

  // Đặt điểm: lấy vị trí chuột thật từ Go (không phụ thuộc độ lệch đã đo hay chưa).
  const onClick = async (e: MouseEvent) => {
    if (skipClick.current) {
      skipClick.current = false;
      return;
    }
    if (onUI(e) || preview) return;
    if (focus ? !repickRef.current : tool === 'region') return; // xem riêng: chỉ đặt điểm khi đang chọn lại
    const screen = inWails ? await MousePos() : toScreenFromClient({ x: e.clientX, y: e.clientY });
    placePoint(screen);
    if (focus) setRepicking(pickingStep() !== null); // kéo thả: còn chờ B thì vẫn đang chọn
  };

  const onUp = async (e: MouseEvent) => {
    const m = movingRef.current;
    if (m) {
      setMoving(null);
      if (!m.moved) return selectStep(m.id);
      const screen = inWails ? await MousePos() : toScreenFromClient({ x: e.clientX, y: e.clientY });
      movePoint(m.id, m.index, screen);
      return;
    }
    if (!dragFrom) return;
    const r = rectOf(dragFrom, { x: e.clientX, y: e.clientY });
    setDragFrom(null);
    if (r.w < MIN_REGION || r.h < MIN_REGION) return;
    const tl = toScreenFromClient(r);
    setRegion({ ...tl, w: Math.round(r.w * map.kx), h: Math.round(r.h * map.ky) });
  };

  // Vùng hiển thị: đang kéo → vùng tạm; đã chọn → đổi sang toạ độ cửa sổ; null = toàn màn hình (không vẽ).
  const shown = dragFrom ? rectOf(dragFrom, cursor)
    : region && { ...toClient(region), w: region.w / map.kx, h: region.h / map.ky };
  // Nhãn kích thước vùng theo toạ độ thật (robotgo), không phải px trong cửa sổ.
  const shownSize = dragFrom && shown ? { w: Math.round(shown.w * map.kx), h: Math.round(shown.h * map.ky) } : region;

  // Marker của các bước đang bật, đánh số giống danh sách bước; kéo thả có A và B.
  // Marker đang kéo vẽ theo con trỏ; đường nối cũng đi theo.
  // ghost: điểm A của bước kéo "từ điểm trước" – trùng marker bước trước nên không vẽ, chỉ dùng cho đường nối.
  // pt: điểm đã lưu (có neo) – không có với A "từ điểm trước" và B "lệch so với A".
  const markers: { id: string; index: number; label: string; at: Point; drag: boolean; jitter: number; ghost: boolean; pt?: StepPoint }[] = [];
  // Đang kéo vùng mới: marker tính theo vùng tạm để thấy ngay điểm bám theo góc nó neo.
  const liveRegion = dragFrom && shown && shown.w >= MIN_REGION && shown.h >= MIN_REGION
    ? { ...toScreenFromClient(shown), w: Math.round(shown.w * map.kx), h: Math.round(shown.h * map.ky) } : region;
  const ends = endpoints(steps, liveRegion);
  steps.forEach((s, i) => {
    const e = ends.get(s.id);
    if (!e) return;
    const pts = s.type === 'drag' ? [e.start, e.end] : [e.start];
    pts.forEach((p, j) => {
      if (!p) return;
      const label = s.type === 'drag' ? `${i + 1}${j === 0 ? 'A' : 'B'}` : `${i + 1}`;
      const dragging = moving?.moved && moving.id === s.id && moving.index === j;
      markers.push({
        id: s.id, index: j, label, at: dragging ? cursor : toClient(p), drag: s.type === 'drag',
        jitter: s.jitterPx, ghost: s.type === 'drag' && j === 0 && s.dragFromPrev,
        pt: s.type === 'drag' && (j === 0 ? s.dragFromPrev : s.dragRelative) ? undefined : s.points[j],
      });
    });
  });
  // Đường neo (bước 29): từ marker tới cạnh / giữa đang neo – bước đang chọn / xem riêng, giữ Alt thì mọi điểm.
  const b = boundsOf(liveRegion);
  const bc = { ...toClient(b), w: b.w / map.kx, h: b.h / map.ky };
  const anchorLines = preview ? [] : markers.filter((m) => m.pt && (alt || m.id === (focus ?? selected))).flatMap((m) => {
    const p = m.pt!;
    const edge = (a: StepPoint['ax'], o: number, size: number) => (a === 'end' ? o + size : a === 'center' ? o + size / 2 : o);
    const ex = edge(p.ax, bc.x, bc.w);
    const ey = edge(p.ay, bc.y, bc.h);
    return [
      { key: `${m.label}x`, x1: ex, y1: m.at.y, x2: m.at.x, y2: m.at.y, text: fmtAxis(p.x, p.ax) },
      { key: `${m.label}y`, x1: m.at.x, y1: ey, x2: m.at.x, y2: m.at.y, text: fmtAxis(p.y, p.ay) },
    ];
  });
  const screenCursor = toScreenFromClient(cursor);
  const focusIdx = focus ? steps.findIndex((x) => x.id === focus) + 1 : 0;
  const focusStep = focusIdx > 0 ? steps[focusIdx - 1] : null;
  const focusEnds = focus ? ends.get(focus) : undefined;
  const hint = focus && focusStep
    ? repicking
      ? focusStep.type === 'drag' && !focusStep.dragFromPrev
        ? (dragA ? 'Click điểm B (nơi thả)' : 'Click điểm A (nơi bắt đầu kéo)')
        : 'Click (hoặc fn+F6) để đặt lại điểm'
      : `Đang xem riêng bước ${focusIdx} – kéo marker để sửa vị trí, Esc / Xong để quay lại`
    : preview ? 'Xem trước 1 vòng – con trỏ giả, không điều khiển chuột thật'
    : tool === 'region' ? 'Kéo chuột để chọn vùng – toạ độ các điểm tính theo vùng'
    : tool === 'drag'
      ? (dragA ? 'Click điểm B (nơi thả) · Esc huỷ A' : `Click điểm A (nơi bắt đầu kéo)${picking ? ' – chọn lại cho bước đang chọn' : ''}`)
    : picking ? 'Click để đặt lại điểm cho bước đang chọn'
    : `Click (hoặc fn+F6) để thêm bước ${STEP_LABEL[placingType()].toLowerCase()} mới · kéo marker để dời điểm`;

  return (
    // HUD luôn tối (design v2): data-theme="dark" cho mọi token bên trong
    <div data-theme="dark" className={`fixed inset-0 select-none text-ink ${focus ? 'bg-[rgb(14_16_22/.65)]' : 'bg-dim'} ${preview ? '' : 'cursor-crosshair'}`}
      onMouseMove={onMove} onMouseDown={onDown} onMouseUp={onUp} onClick={onClick}>
      {shown && (
        <div className={focus ? 'opacity-30' : ''}>
          <div className="absolute rounded-md border-2 border-dashed border-white bg-white/5"
            style={{ left: shown.x, top: shown.y, width: shown.w, height: shown.h }} />
          <span className="absolute rounded-md bg-accent px-2 py-[3px] font-mono text-xs text-on-accent"
            style={{ left: shown.x, top: shown.y - 30 }}>Vùng {shownSize?.w} × {shownSize?.h}</span>
          {[[shown.x, shown.y], [shown.x + shown.w, shown.y], [shown.x, shown.y + shown.h], [shown.x + shown.w, shown.y + shown.h]].map(([x, y], i) => (
            <span key={i} className="absolute h-3 w-3 -translate-x-1/2 -translate-y-1/2 rounded-[3px] border-2 border-accent bg-white" style={{ left: x, top: y }} />
          ))}
        </div>
      )}

      <svg className="pointer-events-none absolute inset-0 h-full w-full" aria-hidden>
        <defs>
          <marker id="ah" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
            <path d="M0 0 10 5 0 10z" className="fill-accent" />
          </marker>
        </defs>
        {anchorLines.map((l) => (
          <g key={l.key}>
            <line x1={l.x1} y1={l.y1} x2={l.x2} y2={l.y2} stroke="#fff" strokeWidth={1} strokeDasharray="3 4" opacity={0.7} />
            <circle cx={l.x1} cy={l.y1} r={3} fill="#fff" opacity={0.8} />
            <text x={(l.x1 + l.x2) / 2} y={(l.y1 + l.y2) / 2} dx={l.y1 === l.y2 ? 0 : 6} dy={l.y1 === l.y2 ? -6 : 4}
              textAnchor={l.y1 === l.y2 ? 'middle' : 'start'} className="fill-white font-mono text-[11px]" opacity={0.85}>{l.text}</text>
          </g>
        ))}
        {dragA && (
          <line x1={toClient(dragA).x} y1={toClient(dragA).y} x2={cursor.x} y2={cursor.y}
            className="stroke-accent" strokeWidth={3} strokeDasharray="8 5" markerEnd="url(#ah)" />
        )}
        {markers.slice(1).map((m, i) => {
          const prev = markers[i];
          const isDrag = prev.drag && m.drag && prev.label.endsWith('A');
          // xem riêng: chỉ đường A → B của bước đang xem giữ rõ
          const dim = focus && !(prev.id === focus && m.id === focus);
          return isDrag
            ? <line key={i} x1={prev.at.x} y1={prev.at.y} x2={m.at.x} y2={m.at.y} className="stroke-accent" strokeWidth={3} markerEnd="url(#ah)" opacity={dim ? 0.2 : 1} />
            : <line key={i} x1={prev.at.x} y1={prev.at.y} x2={m.at.x} y2={m.at.y} stroke="#fff" strokeWidth={2} strokeDasharray="6 6" opacity={dim ? 0.15 : 0.8} />;
        })}
      </svg>

      {/* Ô ngẫu nhiên vị trí (± px, bước 19) quanh marker */}
      {markers.filter((m) => m.jitter > 0 && !m.ghost).map((m) => (
        <span key={`j${m.label}`} aria-hidden
          className={`pointer-events-none absolute -translate-x-1/2 -translate-y-1/2 border border-dashed ${focus && m.id !== focus ? 'opacity-20' : ''} border-accent bg-accent/15`}
          style={{ left: m.at.x, top: m.at.y, width: (2 * m.jitter + 1) / map.kx, height: (2 * m.jitter + 1) / map.ky }} />
      ))}

      {markers.filter((m) => !m.ghost).map((m) => (
        <span key={m.label} style={{ left: m.at.x, top: m.at.y }} title="Kéo để dời điểm, click để chọn bước"
          onMouseDown={(e) => onMarkerDown(e, m.id, m.index)}
          className={`${preview ? 'pointer-events-none opacity-70' : ''} ${focus && m.id !== focus ? 'pointer-events-none opacity-25 grayscale' : ''} ${focus === m.id ? 'z-10 scale-110' : ''} absolute grid h-[34px] w-[34px] -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full border-[3px] border-white text-[13px] font-bold text-on-accent shadow-[0_2px_8px_rgb(0_0_0/.35)] bg-accent ${moving ? 'cursor-grabbing' : 'cursor-grab'} ${m.id === selected ? 'outline-[3px] outline-offset-2 outline-white outline-dashed' : ''}`}>
          {m.label}
        </span>
      ))}

      {dragA && (
        <span style={{ left: toClient(dragA).x, top: toClient(dragA).y }}
          className="pointer-events-none absolute grid h-[34px] w-[34px] -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full border-[3px] border-white bg-accent text-[13px] font-bold text-on-accent shadow-[0_2px_8px_rgb(0_0_0/.35)]">
          A
        </span>
      )}

      {preview && <SimulationLayer sim={sim} toClient={toClient} />}

      {/* Crosshair + toạ độ màn hình thật theo chuột (ẩn khi xem trước) */}
      <div className={`pointer-events-none absolute ${preview ? 'hidden' : ''}`} style={{ left: cursor.x, top: cursor.y }}>
        <div className="absolute -top-10 h-20 w-px bg-white" />
        <div className="absolute -left-10 h-px w-20 bg-white" />
        <span className="absolute left-3.5 top-3 whitespace-nowrap rounded bg-hud px-2 py-1 font-mono text-[13px] text-hud-text">X {screenCursor.x}  Y {screenCursor.y}</span>
      </div>

      {/* Banner hướng dẫn giữa phía trên, ngay dưới thanh công cụ */}
      <div data-ui className="absolute left-1/2 top-[96px] flex max-w-[min(720px,calc(100vw-600px))] -translate-x-1/2 flex-col items-center gap-1 rounded-[10px] border border-hud-line bg-hud px-4 py-2.5 text-center text-[13px] text-hud-text">
        <span className="flex items-center gap-2.5">
          <span className="grid h-[22px] w-[22px] flex-none place-items-center rounded-[7px] bg-accent-soft text-accent-text">
            {preview ? <Eye className="h-[13px] w-[13px]" /> : tool === 'region' && !focus ? <Region className="h-[13px] w-[13px]" />
              : <StepIcon type={tool === 'drag' || focusStep?.type === 'drag' ? 'drag' : 'click'} className="h-[13px] w-[13px]" />}
          </span>
          <span>{hint}</span>
        </span>
        {preview && sim.label && (
          <span className="font-mono text-[#FF8CC6]">{sim.step && sim.running ? `Bước ${sim.step}/${steps.length} · ` : ''}{sim.label}</span>
        )}
        {notice && <span className="text-ok">{notice}</span>}
      </div>

      {/* Thông tin toạ độ bên trái, ngang banner hướng dẫn */}
      <div data-ui className="absolute left-6 top-[96px] flex max-w-[260px] flex-col gap-0.5 rounded-lg border border-hud-line bg-hud/95 px-3 py-2 font-mono text-xs text-hud-text">
        <span>{region ? `Vùng X ${region.x} Y ${region.y} · ${region.w} × ${region.h}` : 'Vùng: toàn màn hình'}</span>
        {focus && focusStep && (focusStep.type === 'drag' ? [['A', focusEnds?.start], ['B', focusEnds?.end]] as const : [['Điểm', focusEnds?.start]] as const).map(([k, p]) => (
          <span key={k}>{k}: {p ? `X ${p.x} Y ${p.y}` : '—'}</span>
        ))}
        {inWails && (
          <span className="text-hud-text/60">
            {measured
              ? `Đo toạ độ: tỷ lệ ${map.kx.toFixed(2)}${Math.abs(map.kx - map.ky) > 0.01 ? `/${map.ky.toFixed(2)}` : ''} · lệch X ${Math.round(map.bx)} Y ${Math.round(map.by)}`
              : 'Đang đo toạ độ – rê chuột rồi dừng lại'}
          </span>
        )}
      </div>

      {/* Thanh công cụ giữa phía trên (artboard để ở đáy – dời lên vì đáy bị Dock / taskbar che) */}
      <div data-ui className="absolute left-1/2 top-6 flex -translate-x-1/2 items-center gap-1 whitespace-nowrap rounded-xl border border-hud-line bg-hud p-1.5 text-hud-text shadow-[0_12px_32px_rgb(0_0_0/.4)]">
        {focus && focusStep ? (
          <>
            <span className="flex items-center gap-2 px-2.5 font-semibold">
              <span className="grid h-7 w-7 place-items-center rounded-full bg-accent text-[13px] text-on-accent">{focusIdx}</span>
              Bước {focusIdx} · {STEP_LABEL[focusStep.type]}
            </span>
            <button onClick={() => { if (repicking) { cancelPick(); setRepicking(false); } else { startPick(focus); setRepicking(true); } }}
              className={`${toolBtn} ${repicking ? 'bg-white/14 text-white' : 'hover:bg-white/6'}`}>
              <Target />{repicking ? 'Huỷ chọn lại' : 'Chọn lại trên màn hình'}
            </button>
            <span className="mx-1.5 h-[26px] w-px bg-white/16" />
            <button onClick={() => { cancelPick(); setRepicking(false); setFocus(null); setTool('point'); }} className={`${toolBtn} font-normal hover:bg-white/6`}>Xem tất cả</button>
          </>
        ) : preview ? (
          <>
            <button onClick={() => (sim.running ? stopSim() : startSim())}
              className={`${toolBtn} ${sim.running ? 'btn-danger-soft' : 'btn-primary'}`}>
              {sim.running ? <><Stop className="h-2.5 w-2.5" />Dừng</> : <><Play className="h-3 w-3" />Chạy thử</>}
            </button>
            <div className="flex rounded-lg bg-white/5 p-[3px]" role="radiogroup" aria-label="Tốc độ xem trước">
              {SPEEDS.map((v) => (
                <button key={v} role="radio" aria-checked={speed === v} onClick={() => setSpeed(v)}
                  className={`h-[38px] rounded-md px-3 font-mono ${speed === v ? 'bg-white/14 text-white' : 'hover:bg-white/6'}`}>
                  {v}×
                </button>
              ))}
            </div>
            <span className="mx-1.5 h-[26px] w-px bg-white/16" />
            <button onClick={() => { stopSim(); setTool('point'); }} className={`${toolBtn} font-normal hover:bg-white/6`}>Sửa điểm</button>
          </>
        ) : (
          <>
            {TOOLS.map((t) => (
              <button key={t.id} disabled={!!t.soon} title={t.soon && `Sắp có (bước ${t.soon})`} onClick={() => setTool(t.id)}
                className={`${toolBtn} disabled:opacity-40 ${tool === t.id ? 'bg-white/14 text-white' : 'hover:bg-white/6'}`}>
                {t.icon}{t.label}
              </button>
            ))}
            {tool === 'region' && region && (
              <button onClick={() => setRegion(null)} className={`${toolBtn} font-normal hover:bg-white/6`}>Toàn màn hình</button>
            )}
            <span className="mx-1.5 h-[26px] w-px bg-white/16" />
            <span className="flex items-center gap-1.5 px-1.5 text-[13px] text-hud-text/60">
              <kbd className="rounded border border-hud-line px-1.5 font-mono text-hud-text">F6</kbd>lưu vị trí chuột
            </span>
            <button onClick={() => setTool('preview')} className={`${toolBtn} font-normal hover:bg-white/6`}>Xem trước</button>
          </>
        )}
        <button onClick={() => setWindowMode('main')} className={`${toolBtn} btn-primary ml-1 font-bold`}>Xong <span className="font-mono text-[11px] font-normal opacity-70">Esc</span></button>
      </div>
    </div>
  );
}
