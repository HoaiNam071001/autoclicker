// Overlay chọn vùng / điểm – cửa sổ trong suốt phủ màn hình, always-on-top
import { MouseEvent, ReactNode, useCallback, useEffect, useRef, useState } from 'react';
import { MousePos } from '../../wailsjs/go/main/App';
import {
  boundsOf, cancelDragA, cancelPick, endpoints, getEditor, hasPoints, movePoint, pickingStep, placePoint, placingType, pointKey,
  regionMapping, RegionMapMode, selectableIndexes, selectStep, setDragMode, setFocusLock, setMulti, setRegion, setRegionMapped, startPick,
  transformedSteps, transformSelected, useEditor,
} from '../scriptStore';
import { frameRect, MapDialog, modKeyName, paramsPreview, TransformBar, TransformFrame, useTransformGesture } from '../components/TransformTool';
import { bbox, NO_TF, TfParams, translate } from '../transform';
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
const TOP_EDGE = 12; // px – chuột sát mép trên thì hiện thanh công cụ
const BAR_HIDE_MS = 350; // rời thanh / gợi ý bao lâu thì ẩn
const BAR_INTRO_MS = 2500; // mở overlay: hiện thanh một lúc cho dễ thấy
const NOTICE_MS = 3500; // thông báo trong gợi ý tự tắt

const isMac = navigator.platform.startsWith('Mac');

// Marker đang được kéo: bước, điểm thứ mấy, vị trí bấm xuống (toạ độ cửa sổ), đã kéo đủ xa chưa.
type Moving = { id: string; index: number; from: Point; moved: boolean };

const TOOLS: { id: OverlayTool; label: string; soon?: string; icon: ReactNode }[] = [
  { id: 'point', label: 'Đặt điểm', icon: <StepIcon type="click" /> },
  { id: 'drag', label: 'Kéo thả A → B', icon: <StepIcon type="drag" /> },
  { id: 'region', label: 'Chọn vùng', icon: <Region /> },
  { id: 'select', label: 'Chọn nhiều', icon: <SelectIcon /> },
  { id: 'capture', label: 'Khoanh ảnh', soon: '25', icon: <StepIcon type="waitImage" /> },
];
const toolBtn = 'flex h-11 items-center gap-2 rounded-[7px] px-3.5 font-semibold';

const rectOf = (a: Point, b: Point): Rect => ({
  x: Math.min(a.x, b.x), y: Math.min(a.y, b.y), w: Math.abs(a.x - b.x), h: Math.abs(a.y - b.y),
});

// Điểm chọn được (khoá chọn nhiều → toạ độ màn hình) từ endpoints(): A "từ điểm trước" không chọn được.
function markersOf(ends: ReturnType<typeof endpoints>): [string, Point][] {
  const out: [string, Point][] = [];
  getEditor().script.steps.forEach((s) => {
    const e = ends.get(s.id);
    if (!e) return;
    selectableIndexes(s).forEach((j) => {
      const p = s.type === 'drag' && j === 1 ? e.end : e.start;
      if (p) out.push([pointKey(s.id, j), p]);
    });
  });
  return out;
}

export default function OverlayScreen() {
  const { script: { region, steps }, selected, notice, dragA, multi } = useEditor();
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
  // Chọn nhiều (bước 30): ô nhập số đang xem trước; vùng mới đang chờ chọn cách xử lý điểm (undefined = không chờ).
  const [tf, setTf] = useState<TfParams>(NO_TF);
  const [pending, setPendingState] = useState<{ to: Rect | null } | null>(null);
  const pendingRef = useRef(pending);
  const setPending = (p: { to: Rect | null } | null) => {
    pendingRef.current = p;
    setPendingState(p);
  };
  const [mapMode, setMapMode] = useState<RegionMapMode>('anchor');
  const toolRef = useRef(tool);
  toolRef.current = tool;
  // Thanh công cụ tự ẩn: mở khi chuột sát mép trên / rê vào gợi ý, đóng khi rời. Lúc mở overlay hiện một lúc.
  const [barOpen, setBarOpen] = useState(true);
  const barTimer = useRef<number | undefined>(undefined);
  const openBar = () => {
    clearTimeout(barTimer.current);
    setBarOpen(true);
  };
  const leaveBar = () => {
    clearTimeout(barTimer.current);
    barTimer.current = window.setTimeout(() => setBarOpen(false), BAR_HIDE_MS);
  };
  useEffect(() => {
    barTimer.current = window.setTimeout(() => setBarOpen(false), BAR_INTRO_MS);
    return () => clearTimeout(barTimer.current);
  }, []);
  // Thông báo (đã đặt điểm…) hiện trong gợi ý vài giây rồi tắt.
  const [noticeOn, setNoticeOn] = useState(false);
  useEffect(() => {
    if (!notice) return;
    setNoticeOn(true);
    const t = window.setTimeout(() => setNoticeOn(false), NOTICE_MS);
    return () => clearTimeout(t);
  }, [notice]);
  const skipClick = useRef(false); // click sinh ra sau khi thả marker không được tính là đặt điểm
  const cursorRef = useRef(cursor);
  const timer = useRef<number | undefined>(undefined);

  const toScreenFromClient = useCallback((p: Point): Point => toScreenPt(map, p), [map]);
  const gesture = useTransformGesture(toScreenFromClient);
  const gestureRef = useRef(gesture);
  gestureRef.current = gesture;

  useEffect(() => {
    resetSamples(); // overlay có thể đang mở trên màn khác lần trước
    const onKey = (e: KeyboardEvent) => {
      if (toolRef.current === 'select' && selectKeys(e)) return;
      if (e.key !== 'Escape') return;
      // Đang kéo vùng / marker thì Esc chỉ huỷ thao tác đó.
      if (gestureRef.current.cancel()) return;
      if (simRef.current.running) simRef.current.stop(); // xem trước đang chạy → Esc chỉ dừng
      else if (movingRef.current) setMoving(null);
      else if (dragRef.current) setDragFrom(null);
      else if (cancelDragA()) return;
      else if (repickRef.current) {
        cancelPick();
        setRepicking(false);
      } else if (toolRef.current === 'select' && getEditor().multi.length) setMulti([]); // Esc bỏ chọn trước
      else setWindowMode('main');
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
  useEffect(() => setTf(NO_TF), [multi.join(), tool]); // đổi lựa chọn → bỏ số đang nhập dở trên thanh biến đổi
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
  const onUI = (e: MouseEvent) => !!(e.target as HTMLElement).closest('[data-ui]');

  const onMove = (e: MouseEvent) => {
    const c = { x: e.clientX, y: e.clientY };
    setCursor(c);
    if (e.buttons === 0 && c.y <= TOP_EDGE) openBar();
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
    if (onUI(e) || e.button !== 0 || (tool !== 'region' && tool !== 'select') || focus) return; // kéo vùng / khung chọn
    setDragFrom({ x: e.clientX, y: e.clientY });
  };

  const onMarkerDown = (e: MouseEvent, id: string, index: number) => {
    if (e.button !== 0) return;
    e.stopPropagation(); // không bắt đầu kéo vùng
    skipClick.current = true;
    if (tool === 'select') {
      // Shift / ⌘ / Ctrl + click: thêm / bỏ; click điểm chưa chọn: chỉ chọn điểm đó; kéo = di chuyển cả nhóm.
      const key = pointKey(id, index);
      const has = multi.includes(key);
      if (e.shiftKey || e.metaKey || e.ctrlKey) return setMulti(has ? multi.filter((k) => k !== key) : [...multi, key]);
      const sel = has ? multi : [key];
      if (!has) setMulti(sel);
      selectStep(id);
      return gesture.begin(e, 'move', bbox(markers.filter((m) => sel.includes(m.key)).map((m) => m.scr)));
    }
    setMoving({ id, index, from: { x: e.clientX, y: e.clientY }, moved: false });
  };

  // Đặt điểm: lấy vị trí chuột thật từ Go (không phụ thuộc độ lệch đã đo hay chưa).
  const onClick = async (e: MouseEvent) => {
    if (skipClick.current) {
      skipClick.current = false;
      return;
    }
    if (onUI(e) || preview) return;
    if (focus ? !repickRef.current : tool === 'region' || tool === 'select') return; // xem riêng: chỉ đặt điểm khi đang chọn lại
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
    if (tool === 'select') { // khung chọn: chọn các marker nằm trong; Shift / ⌘ cộng thêm; click chỗ trống = bỏ chọn
      const add = e.shiftKey || e.metaKey || e.ctrlKey;
      if (r.w < MIN_REGION && r.h < MIN_REGION) return add || setMulti([]);
      const inside = markers.filter((m) => m.at.x >= r.x && m.at.x <= r.x + r.w && m.at.y >= r.y && m.at.y <= r.y + r.h).map((m) => m.key);
      return setMulti(add ? [...new Set([...multi, ...inside])] : inside);
    }
    if (r.w < MIN_REGION || r.h < MIN_REGION) return;
    const tl = toScreenFromClient(r);
    requestRegion({ ...tl, w: Math.round(r.w * map.kx), h: Math.round(r.h * map.ky) });
  };

  // Đổi vùng khi đã có điểm: hỏi cách xử lý các điểm (giữ neo / co giãn / giữ chỗ), xem trước trên overlay.
  const requestRegion = (to: Rect | null) => {
    if (!hasPoints()) return setRegion(to);
    setMapMode('anchor');
    setPending({ to });
  };

  // Phím của công cụ Chọn nhiều: ⌘/Ctrl+A chọn tất cả, mũi tên dời 1 px (Shift 10 px). Trả true nếu đã xử lý.
  function selectKeys(e: KeyboardEvent) {
    if ((e.target as HTMLElement).closest('input') || pendingRef.current) return false;
    const { script, multi: sel } = getEditor();
    if (e.key.toLowerCase() === 'a' && (isMac ? e.metaKey : e.ctrlKey)) {
      e.preventDefault();
      setMulti(script.steps.filter((s) => s.enabled).flatMap((s) => selectableIndexes(s).map((j) => pointKey(s.id, j))));
      return true;
    }
    const d = { ArrowLeft: [-1, 0], ArrowRight: [1, 0], ArrowUp: [0, -1], ArrowDown: [0, 1] }[e.key];
    if (!d || !sel.length) return false;
    e.preventDefault();
    const k = e.shiftKey ? 10 : 1;
    transformSelected(translate(d[0] * k, d[1] * k), 'di chuyển');
    return true;
  }

  // Vùng hiển thị: đang kéo → vùng tạm; chờ chọn cách đổi vùng → vùng mới; đã chọn → đổi sang toạ độ cửa sổ;
  // null = toàn màn hình (không vẽ). Khung chọn nhiều (tool select) không phải vùng.
  const shownRegion = pending ? pending.to : region;
  const regionDrag = dragFrom && tool === 'region' ? dragFrom : null;
  const shown = regionDrag ? rectOf(regionDrag, cursor)
    : shownRegion && { ...toClient(shownRegion), w: shownRegion.w / map.kx, h: shownRegion.h / map.ky };
  // Nhãn kích thước vùng theo toạ độ thật (robotgo), không phải px trong cửa sổ.
  const shownSize = regionDrag && shown ? { w: Math.round(shown.w * map.kx), h: Math.round(shown.h * map.ky) } : shownRegion;

  // Marker của các bước đang bật, đánh số giống danh sách bước; kéo thả có A và B.
  // Marker đang kéo vẽ theo con trỏ; đường nối cũng đi theo.
  // ghost: điểm A của bước kéo "từ điểm trước" – trùng marker bước trước nên không vẽ, chỉ dùng cho đường nối.
  // pt: điểm đã lưu (có neo) – không có với A "từ điểm trước" và B "lệch so với A".
  // key: khoá chọn nhiều; scr: toạ độ màn hình (đang xem trước nếu đang biến đổi).
  const markers: { id: string; index: number; key: string; label: string; at: Point; scr: Point; drag: boolean; jitter: number; ghost: boolean; pt?: StepPoint }[] = [];
  // Đang kéo vùng mới: marker tính theo vùng tạm để thấy ngay điểm bám theo góc nó neo.
  const liveRegion = regionDrag && shown && shown.w >= MIN_REGION && shown.h >= MIN_REGION
    ? { ...toScreenFromClient(shown), w: Math.round(shown.w * map.kx), h: Math.round(shown.h * map.ky) } : shownRegion;
  // Xem trước bằng đúng logic ghi thật: chờ đổi vùng → điểm theo cách đang chọn; đang kéo khung / nhập số → điểm sau biến đổi.
  const committed = endpoints(steps, region);
  const selScreen = new Map(markersOf(committed).filter(([k]) => multi.includes(k)));
  const liveMat = tool === 'select' && !pending ? gesture.live ?? paramsPreview(tf, [...selScreen.values()]) : null;
  const view = pending ? regionMapping(pending.to, mapMode) : liveMat ? transformedSteps(liveMat) : { steps, outside: 0 };
  const viewSteps = view.steps;
  const ends = endpoints(viewSteps, liveRegion);
  viewSteps.forEach((s, i) => {
    const e = ends.get(s.id);
    if (!e) return;
    const pts = s.type === 'drag' ? [e.start, e.end] : [e.start];
    pts.forEach((p, j) => {
      if (!p) return;
      const label = s.type === 'drag' ? `${i + 1}${j === 0 ? 'A' : 'B'}` : `${i + 1}`;
      const dragging = moving?.moved && moving.id === s.id && moving.index === j;
      markers.push({
        id: s.id, index: j, key: pointKey(s.id, j), label, at: dragging ? cursor : toClient(p), scr: p, drag: s.type === 'drag',
        jitter: s.jitterPx, ghost: s.type === 'drag' && j === 0 && s.dragFromPrev,
        pt: s.type === 'drag' && (j === 0 ? s.dragFromPrev : s.dragRelative) ? undefined : s.points[j],
      });
    });
  });
  const selPoints = markers.filter((m) => multi.includes(m.key)).map((m) => m.scr);
  const selBox = selPoints.length ? bbox(selPoints) : null;
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
      : 'Kéo marker để sửa · Esc quay lại'
    : preview ? 'Con trỏ giả, không click thật'
    : tool === 'region' ? 'Kéo để khoanh vùng'
    : tool === 'select'
      ? multi.length
        ? `${multi.length} điểm · kéo khung / núm để biến đổi · ${modKeyName}+núm: xiên, phối cảnh`
        : `Click, Shift+click hoặc kéo khung để chọn · ${modKeyName}A tất cả`
    : tool === 'drag'
      ? (dragA ? 'Click điểm B · Esc huỷ A' : `Click điểm A${picking ? ' (chọn lại)' : ''}`)
    : picking ? 'Click để đặt lại điểm'
    : `Click / F6: thêm ${STEP_LABEL[placingType()].toLowerCase()} · kéo marker để dời`;
  const toolName = focus && focusStep ? `Bước ${focusIdx}` : preview ? 'Xem trước' : TOOLS.find((t) => t.id === tool)?.label ?? '';

  return (
    // HUD luôn tối (design v2): data-theme="dark" cho mọi token bên trong
    <div data-theme="dark" className={`fixed inset-0 select-none text-ink ${focus ? 'bg-[rgb(14_16_22/.65)]' : 'bg-dim'} ${preview ? '' : 'cursor-crosshair'}`}
      onMouseMove={onMove} onMouseDown={onDown} onMouseUp={onUp} onClick={onClick}
      // macOS: overlay nằm dưới menu bar – rê nhanh lên trên có thể vượt mép mà không qua vùng TOP_EDGE
      onMouseLeave={(e) => e.buttons === 0 && e.clientY < TOP_EDGE * 3 && openBar()}>
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

      {/* Chọn nhiều: khung biến đổi quanh các điểm đã chọn (theo vị trí đang xem trước) + khung kéo chọn */}
      {tool === 'select' && !pending && selBox && (
        <TransformFrame box={selBox} toClient={toClient} bad={view.outside > 0} single={selPoints.length < 2 || (selBox.w < 1 && selBox.h < 1)}
          begin={(e, kind, _box, fx, fy) => gesture.begin(e, kind, bbox([...selScreen.values()]), fx, fy)} />
      )}
      {tool === 'select' && dragFrom && (
        <div className="pointer-events-none absolute border border-accent bg-accent/10" style={(({ x, y, w, h }) => ({ left: x, top: y, width: w, height: h }))(rectOf(dragFrom, cursor))} />
      )}

      {markers.filter((m) => !m.ghost).map((m) => (
        <span key={m.label} style={{ left: m.at.x, top: m.at.y }}
          title={tool === 'select' ? 'Click chọn · Shift / ⌘ + click thêm bớt · kéo để di chuyển cả nhóm' : 'Kéo để dời điểm, click để chọn bước'}
          onMouseDown={(e) => onMarkerDown(e, m.id, m.index)}
          className={`${tool === 'select' ? (multi.includes(m.key) ? 'z-10 outline-2 outline-offset-2 outline-white' : 'opacity-60') : ''} ${preview ? 'pointer-events-none opacity-70' : ''} ${focus && m.id !== focus ? 'pointer-events-none opacity-25 grayscale' : ''} ${focus === m.id ? 'z-10 scale-110' : ''} absolute grid h-[34px] w-[34px] -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full border-[3px] border-white text-[13px] font-bold text-on-accent shadow-[0_2px_8px_rgb(0_0_0/.35)] bg-accent ${moving ? 'cursor-grabbing' : 'cursor-grab'} ${m.id === selected ? 'outline-[3px] outline-offset-2 outline-white outline-dashed' : ''}`}>
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

      {/* Gợi ý gọn một dòng giữa phía trên – cũng là tab mở thanh công cụ (rê vào / đưa chuột sát mép trên) */}
      <div data-ui onMouseEnter={() => openBar()} onMouseLeave={leaveBar} style={{ top: barOpen ? 76 : 8 }}
        className="absolute left-1/2 flex h-8 max-w-[min(760px,calc(100vw-48px))] -translate-x-1/2 cursor-default items-center gap-2 whitespace-nowrap rounded-full border border-hud-line bg-hud/95 pr-3.5 pl-1.5 text-[12px] text-hud-text shadow-[0_6px_18px_rgb(0_0_0/.3)] transition-[top] duration-150">
        <span className="grid h-[22px] w-[22px] flex-none place-items-center rounded-full bg-accent-soft text-accent-text">
          {preview ? <Eye className="h-3 w-3" /> : tool === 'region' && !focus ? <Region className="h-3 w-3" />
            : tool === 'select' && !focus ? <SelectIcon className="h-3 w-3" />
            : <StepIcon type={tool === 'drag' || focusStep?.type === 'drag' ? 'drag' : 'click'} className="h-3 w-3" />}
        </span>
        <b className="font-semibold text-white">{toolName}</b>
        <span className="min-w-0 truncate text-hud-text/75">{hint}</span>
        {preview && sim.label && (
          <span className="truncate font-mono text-[#FF8CC6]">{sim.step && sim.running ? `${sim.step}/${steps.length} · ` : ''}{sim.label}</span>
        )}
        {notice && noticeOn && <span className="truncate text-ok">· {notice}</span>}
        {!barOpen && <svg width="10" height="10" viewBox="0 0 10 10" className="flex-none text-hud-text/50" aria-hidden><path d="M2 3.5 5 6.5 8 3.5" fill="none" stroke="currentColor" strokeWidth="1.5" /></svg>}
      </div>

      {/* Thông tin toạ độ: góc trên trái, chữ nhỏ */}
      <div data-ui className="absolute top-2 left-2 flex max-w-[240px] flex-col rounded-md border border-hud-line bg-hud/85 px-2 py-1 font-mono text-[11px] leading-[16px] text-hud-text/80">
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

      {/* Thanh công cụ giữa phía trên (đáy bị Dock / taskbar che) – tự ẩn, hiện khi đưa chuột sát mép trên / rê vào gợi ý */}
      <div data-ui onMouseEnter={() => openBar()} onMouseLeave={leaveBar}
        className={`absolute left-1/2 top-2 flex -translate-x-1/2 transition-[translate,opacity] duration-150 ${barOpen ? '' : 'pointer-events-none -translate-y-[calc(100%+16px)] opacity-0'} items-center gap-1 whitespace-nowrap rounded-xl border border-hud-line bg-hud p-1.5 text-hud-text shadow-[0_12px_32px_rgb(0_0_0/.4)]`}>
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
            <button onClick={() => { stopSim(); setTool('select'); }} title="Chọn nhiều điểm để di chuyển / co giãn / xoay, rồi Xem trước lại" className={`${toolBtn} font-normal hover:bg-white/6`}>Sửa điểm</button>
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
              <button onClick={() => requestRegion(null)} className={`${toolBtn} font-normal hover:bg-white/6`}>Toàn màn hình</button>
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

      {tool === 'select' && !pending && selScreen.size > 0 && <TransformBar points={selScreen} tf={tf} setTf={setTf}
        frame={frameRect(bbox([...selScreen.values()]), toClient)} />}
      {pending && (
        <MapDialog from={boundsOf(region)} to={boundsOf(pending.to)} target={pending.to} mode={mapMode} setMode={setMapMode}
          onApply={() => { setRegionMapped(pending.to, mapMode); setPending(null); }}
          onCancel={() => setPending(null)} />
      )}
    </div>
  );
}
