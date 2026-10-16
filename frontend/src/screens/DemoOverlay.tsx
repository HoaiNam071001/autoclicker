// Màn #demo-overlay (bước 27b): giao diện overlay cho 2 công cụ chưa có logic, trên một "trang web" giả.
// - Khoanh ảnh mẫu (artboard Overlay, logic ở bước 25): kéo để khoanh, 4 núm, nhãn kích thước, Dùng ảnh này / Khoanh lại.
// - Chọn nhiều điểm (bước 30): click / Shift+click / kéo khung để chọn, khung biến đổi 8 núm + núm xoay,
//   ô nhập ΔX ΔY / % / góc / xiên áp lên điểm mẫu, căn chỉnh / lật, hộp "Ánh xạ sang vùng mới".
import { MouseEvent, ReactNode, useState } from 'react';
import { Region, SelectIcon, StepIcon } from '../components/icons';
import { Point } from '../types';

type Tool = 'select' | 'capture';
type Rect = { x: number; y: number; w: number; h: number };
type Tf = { dx: number; dy: number; sx: number; sy: number; rot: number; kx: number; ky: number }; // %, độ
const ID: Tf = { dx: 0, dy: 0, sx: 100, sy: 100, rot: 0, kx: 0, ky: 0 };

const START: Point[] = [
  { x: 330, y: 250 }, { x: 560, y: 250 }, { x: 790, y: 250 }, { x: 330, y: 470 }, { x: 560, y: 470 }, { x: 1020, y: 560 },
];

const rectOf = (a: Point, b: Point): Rect => ({ x: Math.min(a.x, b.x), y: Math.min(a.y, b.y), w: Math.abs(a.x - b.x), h: Math.abs(a.y - b.y) });
const bbox = (ps: Point[]): Rect => {
  const xs = ps.map((p) => p.x), ys = ps.map((p) => p.y);
  return rectOf({ x: Math.min(...xs), y: Math.min(...ys) }, { x: Math.max(...xs), y: Math.max(...ys) });
};
// Biến đổi quanh tâm khung chọn: co giãn → xiên → xoay → dời (bước 30 viết lại trong transform.ts có test)
function apply(p: Point, c: Point, t: Tf): Point {
  let x = (p.x - c.x) * (t.sx / 100), y = (p.y - c.y) * (t.sy / 100);
  const kx = Math.tan((t.kx * Math.PI) / 180), ky = Math.tan((t.ky * Math.PI) / 180);
  [x, y] = [x + kx * y, y + ky * x];
  const r = (t.rot * Math.PI) / 180;
  [x, y] = [x * Math.cos(r) - y * Math.sin(r), x * Math.sin(r) + y * Math.cos(r)];
  return { x: Math.round(x + c.x + t.dx), y: Math.round(y + c.y + t.dy) };
}

export default function DemoOverlay() {
  const [tool, setTool] = useState<Tool>('capture');
  const [cap, setCap] = useState<Rect | null>({ x: 1104, y: 366, w: 164, h: 64 });
  const [drag, setDrag] = useState<Point | null>(null);
  const [cursor, setCursor] = useState<Point>({ x: 0, y: 0 });
  const [pts, setPts] = useState(START);
  const [sel, setSel] = useState<number[]>([0, 1, 3, 4]);
  const [tf, setTf] = useState<Tf>(ID);
  const [mapOpen, setMapOpen] = useState(false);
  const [notice, setNotice] = useState('');

  const selPts = sel.map((i) => pts[i]);
  const box = sel.length ? bbox(selPts) : null;
  const center = box ? { x: box.x + box.w / 2, y: box.y + box.h / 2 } : { x: 0, y: 0 };
  const shown = pts.map((p, i) => (sel.includes(i) ? apply(p, center, tf) : p));
  const dirty = JSON.stringify(tf) !== JSON.stringify(ID);
  const commit = () => { setPts(shown); setTf(ID); setNotice(`Đã áp dụng cho ${sel.length} điểm – Ctrl+Z hoàn tác`); };

  const onUI = (e: MouseEvent) => !!(e.target as HTMLElement).closest('[data-ui]');
  const onDown = (e: MouseEvent) => {
    if (onUI(e) || (e.target as HTMLElement).closest('[data-marker]')) return;
    setDrag({ x: e.clientX, y: e.clientY });
  };
  const onUp = (e: MouseEvent) => {
    if (!drag) return;
    const r = rectOf(drag, { x: e.clientX, y: e.clientY });
    setDrag(null);
    if (r.w < 8 || r.h < 8) return tool === 'select' && !e.shiftKey && setSel([]);
    if (tool === 'capture') setCap(r);
    else {
      const inside = pts.flatMap((p, i) => (p.x >= r.x && p.x <= r.x + r.w && p.y >= r.y && p.y <= r.y + r.h ? [i] : []));
      setSel(e.shiftKey ? [...new Set([...sel, ...inside])] : inside);
      setTf(ID);
    }
  };
  const live = drag && rectOf(drag, cursor);

  // căn chỉnh / chia đều / lật – áp ngay lên điểm đã chọn
  const align = (fn: (p: Point, b: Rect, k: number, n: number) => Point) => {
    if (!box) return;
    const order = [...sel].sort((a, b) => pts[a].x - pts[b].x || pts[a].y - pts[b].y);
    setPts(pts.map((p, i) => (sel.includes(i) ? fn(p, box, order.indexOf(i), sel.length) : p)));
  };

  return (
    <div data-theme="dark" className="fixed inset-0 select-none overflow-hidden bg-[#C9CDD4] text-ink"
      onMouseDown={onDown} onMouseUp={onUp} onMouseMove={(e) => setCursor({ x: e.clientX, y: e.clientY })}>
      <FakePage />

      {tool === 'capture' ? (
        <>
          {/* lớp tối có "lỗ" ở vùng khoanh (box-shadow rất lớn) */}
          {(live ?? cap) ? (
            <div className="pointer-events-none absolute rounded-md border-2 border-dashed border-white shadow-[0_0_0_4000px_rgb(14_16_22/.5)]"
              style={pos(live ?? cap!)} />
          ) : <div className="pointer-events-none absolute inset-0 bg-dim" />}
          {cap && !live && (
            <>
              {corners(cap).map(([x, y], i) => (
                <span key={i} className="absolute h-3 w-3 -translate-x-1/2 -translate-y-1/2 rounded-[3px] border-2 border-accent bg-white" style={{ left: x, top: y }} />
              ))}
              <span className="absolute rounded-md bg-accent px-2 py-[3px] font-mono text-xs text-white" style={{ left: cap.x, top: cap.y - 30 }}>
                Ảnh mẫu {Math.round(cap.w)} × {Math.round(cap.h)}
              </span>
              <div data-ui className="absolute flex gap-1.5 rounded-[10px] border border-hud-line bg-hud p-1.5 shadow-[0_8px_24px_rgb(0_0_0/.35)]"
                style={{ left: Math.max(cap.x + cap.w - 228, 8), top: cap.y + cap.h + 14 }}>
                <button onClick={() => setNotice(`Đã dùng ảnh ${Math.round(cap.w)} × ${Math.round(cap.h)} cho bước 03`)} className="btn-primary h-[38px] rounded-lg px-3.5 font-semibold">Dùng ảnh này</button>
                <button onClick={() => setCap(null)} className="h-[38px] rounded-lg bg-white/10 px-3.5 text-hud-text">Khoanh lại</button>
              </div>
            </>
          )}
        </>
      ) : (
        <>
          <div className="pointer-events-none absolute inset-0 bg-dim" />
          <svg className="pointer-events-none absolute inset-0 h-full w-full" aria-hidden>
            {shown.slice(1).map((p, i) => (
              <line key={i} x1={shown[i].x} y1={shown[i].y} x2={p.x} y2={p.y} stroke="#fff" strokeWidth={2} strokeDasharray="6 6" opacity={0.7} />
            ))}
            {/* vị trí cũ (mờ) khi đang xem trước biến đổi */}
            {dirty && sel.map((i) => <circle key={i} cx={pts[i].x} cy={pts[i].y} r={14} fill="none" stroke="#fff" strokeDasharray="3 4" opacity={0.5} />)}
          </svg>
          {box && <TransformBox box={bbox(sel.map((i) => shown[i]))} rot={tf.rot} />}
          {shown.map((p, i) => (
            <span key={i} data-marker onMouseDown={(e) => { e.stopPropagation(); setTf(ID); setSel(e.shiftKey || e.metaKey ? (sel.includes(i) ? sel.filter((s) => s !== i) : [...sel, i]) : [i]); }}
              className={`absolute grid h-[34px] w-[34px] -translate-x-1/2 -translate-y-1/2 cursor-pointer place-items-center rounded-full border-[3px] text-[13px] font-bold text-white shadow-[0_2px_8px_rgb(0_0_0/.35)] ${sel.includes(i) ? 'border-white bg-accent outline-2 outline-offset-2 outline-white' : 'border-white/80 bg-accent/60'}`}
              style={{ left: p.x, top: p.y }}>{i + 1}</span>
          ))}
          {sel.length > 0 && (
            <TransformBar tf={tf} setTf={setTf} count={sel.length} dirty={dirty} onApply={commit} onReset={() => setTf(ID)}
              onMap={() => setMapOpen(true)}
              align={{
                left: () => align((p, b) => ({ ...p, x: b.x })), right: () => align((p, b) => ({ ...p, x: b.x + b.w })),
                top: () => align((p, b) => ({ ...p, y: b.y })), bottom: () => align((p, b) => ({ ...p, y: b.y + b.h })),
                hcenter: () => align((p, b) => ({ ...p, x: b.x + b.w / 2 })), vcenter: () => align((p, b) => ({ ...p, y: b.y + b.h / 2 })),
                hdist: () => align((p, b, k, n) => ({ ...p, x: Math.round(b.x + (b.w * k) / Math.max(n - 1, 1)) })),
                flipH: () => align((p, b) => ({ ...p, x: b.x + b.x + b.w - p.x })), flipV: () => align((p, b) => ({ ...p, y: b.y + b.y + b.h - p.y })),
              }} />
          )}
        </>
      )}
      {live && tool === 'select' && <div className="pointer-events-none absolute border border-accent bg-accent/10" style={pos(live)} />}

      {/* Banner hướng dẫn */}
      <div data-ui className="absolute top-6 left-1/2 flex -translate-x-1/2 flex-col items-center gap-1 rounded-[10px] border border-hud-line bg-hud px-4 py-2.5 text-[13px] text-hud-text">
        <span className="flex items-center gap-2.5">
          <span className="grid h-[22px] w-[22px] place-items-center rounded-[7px] bg-accent-soft text-accent-text">
            {tool === 'capture' ? <StepIcon type="waitImage" className="h-[13px] w-[13px]" /> : <SelectIcon className="h-[13px] w-[13px]" />}
          </span>
          {tool === 'capture'
            ? <span>Kéo để khoanh ảnh mẫu cho <b>bước 03: Chờ ảnh nút Nhắn tin</b></span>
            : <span>Click / Shift+click điểm hoặc kéo khung để chọn · kéo núm để co giãn, núm tròn để xoay · <b>{sel.length}</b> điểm đang chọn</span>}
        </span>
        {notice && <span className="text-ok">{notice}</span>}
      </div>
      <span data-ui className="absolute top-6 left-6 rounded-full bg-accent-soft px-2.5 py-0.5 text-xs font-semibold text-accent-text">Màn mẫu #demo-overlay</span>

      {/* Thanh công cụ đáy */}
      <div data-ui className="absolute bottom-7 left-1/2 flex -translate-x-1/2 items-center gap-1 rounded-xl border border-hud-line bg-hud p-1.5 whitespace-nowrap text-hud-text shadow-[0_12px_32px_rgb(0_0_0/.4)]">
        <HudBtn><StepIcon type="click" />Đặt điểm</HudBtn>
        <HudBtn><StepIcon type="drag" />Kéo thả A → B</HudBtn>
        <HudBtn><Region />Chọn vùng</HudBtn>
        <HudBtn on={tool === 'select'} onClick={() => { setTool('select'); setNotice(''); }}><SelectIcon />Chọn nhiều</HudBtn>
        <HudBtn on={tool === 'capture'} onClick={() => { setTool('capture'); setNotice(''); }}><StepIcon type="waitImage" />Khoanh ảnh</HudBtn>
        <span className="mx-1.5 h-[26px] w-px bg-white/16" />
        <HudBtn plain>Xem trước</HudBtn>
        <button className="btn-primary ml-1 flex h-11 items-center gap-2 rounded-[7px] px-4 font-bold">Xong <span className="font-mono text-[11px] font-normal opacity-70">Esc</span></button>
      </div>

      {mapOpen && <MapDialog onClose={() => setMapOpen(false)} onApply={(m) => { setMapOpen(false); setNotice(`Ánh xạ: ${m}`); }} />}
    </div>
  );
}

const pos = (r: Rect) => ({ left: r.x, top: r.y, width: r.w, height: r.h });
const corners = (r: Rect) => [[r.x, r.y], [r.x + r.w, r.y], [r.x, r.y + r.h], [r.x + r.w, r.y + r.h]];

function HudBtn({ on, plain, onClick, children }: { on?: boolean; plain?: boolean; onClick?: () => void; children: ReactNode }) {
  return (
    <button onClick={onClick} className={`flex h-11 items-center gap-2 rounded-[7px] px-3.5 ${plain ? '' : 'font-semibold'} ${on ? 'bg-white/14 text-white' : 'hover:bg-white/6'}`}>{children}</button>
  );
}

// Khung biến đổi quanh các điểm đã chọn: viền accent, 8 núm vuông, núm tròn xoay phía trên
function TransformBox({ box, rot }: { box: Rect; rot: number }) {
  const pad = 22;
  const r = { x: box.x - pad, y: box.y - pad, w: box.w + pad * 2, h: box.h + pad * 2 };
  const hs = [[0, 0], [0.5, 0], [1, 0], [0, 0.5], [1, 0.5], [0, 1], [0.5, 1], [1, 1]];
  return (
    <div className="pointer-events-none absolute" style={pos(r)}>
      <div className="absolute inset-0 border border-accent" />
      {hs.map(([fx, fy], i) => (
        <span key={i} className="absolute h-2.5 w-2.5 -translate-x-1/2 -translate-y-1/2 rounded-[2px] border-[1.5px] border-accent bg-white"
          style={{ left: `${fx * 100}%`, top: `${fy * 100}%` }} />
      ))}
      <span className="absolute -top-7 left-1/2 h-7 w-px -translate-x-1/2 bg-accent" />
      <span className="absolute -top-9 left-1/2 grid h-4 w-4 -translate-x-1/2 place-items-center rounded-full border-[1.5px] border-accent bg-white" title="Kéo để xoay" />
      {rot !== 0 && <span className="absolute -top-9 left-1/2 ml-3 rounded bg-accent px-1.5 font-mono text-[11px] text-white">{rot}°</span>}
    </div>
  );
}

function TransformBar({ tf, setTf, count, dirty, onApply, onReset, onMap, align }: {
  tf: Tf; setTf: (t: Tf) => void; count: number; dirty: boolean; onApply: () => void; onReset: () => void; onMap: () => void;
  align: Record<'left' | 'right' | 'top' | 'bottom' | 'hcenter' | 'vcenter' | 'hdist' | 'flipH' | 'flipV', () => void>;
}) {
  const num = (k: keyof Tf, label: string, unit: string, w = 'w-[66px]') => (
    <label className="flex flex-col gap-0.5 text-[11px] text-hud-text/60">{label}
      <span className={`flex h-8 items-center gap-1 rounded-md border border-hud-line bg-black/20 px-2 focus-within:ring-focus ${w}`}>
        <input value={tf[k]} inputMode="numeric" onChange={(e) => /^-?\d*$/.test(e.target.value) && setTf({ ...tf, [k]: Number(e.target.value) || 0 })}
          className="min-w-0 flex-1 bg-transparent font-mono text-[13px] text-white outline-none" />
        <span>{unit}</span>
      </span>
    </label>
  );
  const ico = (d: string, title: string, fn: () => void) => (
    <button title={title} aria-label={title} onClick={fn} className="grid h-8 w-8 place-items-center rounded-md hover:bg-white/10">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" aria-hidden><path d={d} /></svg>
    </button>
  );
  return (
    <div data-ui className="absolute bottom-[96px] left-1/2 flex -translate-x-1/2 items-end gap-3 whitespace-nowrap rounded-xl border border-hud-line bg-hud px-3 py-2.5 text-hud-text shadow-[0_12px_32px_rgb(0_0_0/.4)]">
      <span className="self-center pr-1 text-[13px]"><b className="text-white">{count}</b> điểm</span>
      {num('dx', 'ΔX', 'px')}{num('dy', 'ΔY', 'px')}
      {num('sx', 'Rộng', '%')}{num('sy', 'Cao', '%')}
      {num('rot', 'Xoay', '°', 'w-[58px]')}
      {num('kx', 'Xiên X', '°', 'w-[58px]')}{num('ky', 'Xiên Y', '°', 'w-[58px]')}
      <span className="h-8 w-px self-end bg-white/16" />
      <div className="flex flex-col gap-0.5 text-[11px] text-hud-text/60">Căn / lật
        <div className="flex">
          {ico('M4 4v16M8 8h10M8 16h6', 'Căn trái', align.left)}
          {ico('M12 4v16M6 8h12M8 16h8', 'Căn giữa ngang', align.hcenter)}
          {ico('M20 4v16M6 8h10M10 16h6', 'Căn phải', align.right)}
          {ico('M4 4h16M8 8v10M16 8v6', 'Căn trên', align.top)}
          {ico('M4 12h16M8 6v12M16 8v8', 'Căn giữa dọc', align.vcenter)}
          {ico('M4 20h16M8 6v10M16 10v6', 'Căn dưới', align.bottom)}
          {ico('M4 6v12M12 6v12M20 6v12', 'Chia đều ngang', align.hdist)}
          {ico('M12 3v18M8 7l-4 5 4 5M16 7l4 5-4 5', 'Lật ngang', align.flipH)}
          {ico('M3 12h18M7 8l5-4 5 4M7 16l5 4 5-4', 'Lật dọc', align.flipV)}
        </div>
      </div>
      <span className="h-8 w-px self-end bg-white/16" />
      <div className="flex gap-1.5">
        <button onClick={onMap} title="Ánh xạ các điểm sang vùng mới" className="h-8 rounded-md bg-white/10 px-2.5 text-[13px]">Ánh xạ…</button>
        <button onClick={onReset} disabled={!dirty} className="h-8 rounded-md px-2.5 text-[13px] hover:bg-white/10 disabled:opacity-40">Bỏ</button>
        <button onClick={onApply} disabled={!dirty} className="btn-primary h-8 rounded-md px-3 text-[13px] font-semibold disabled:opacity-50">Áp dụng</button>
      </div>
    </div>
  );
}

const MAP_MODES = [
  { id: 'anchor', title: 'Giữ theo neo', hint: 'Mỗi điểm giữ khoảng cách tới góc / cạnh đang neo (bước 29) – mặc định' },
  { id: 'scale', title: 'Co giãn theo vùng', hint: 'Giữ tỷ lệ vị trí trong vùng: điểm ở giữa vùng cũ vẫn ở giữa vùng mới' },
  { id: 'screen', title: 'Giữ nguyên chỗ trên màn hình', hint: 'Toạ độ thật không đổi, chỉ tính lại theo vùng mới' },
  { id: 'points', title: 'Theo 4 điểm mốc', hint: 'Chọn 4 điểm ở vùng cũ rồi 4 điểm tương ứng ở vùng mới – tự tính phép biến đổi (cả nghiêng / phối cảnh)' },
];

function MapDialog({ onClose, onApply }: { onClose: () => void; onApply: (mode: string) => void }) {
  const [mode, setMode] = useState('anchor');
  return (
    <div data-ui className="absolute inset-0 grid place-items-center bg-black/30">
      <div role="dialog" aria-label="Ánh xạ sang vùng mới" className="flex w-[460px] flex-col gap-3 rounded-panel border border-line bg-surface p-4 text-ink shadow-[0_16px_48px_rgb(0_0_0/.4)]">
        <div>
          <div className="text-base font-semibold">Ánh xạ điểm sang vùng mới</div>
          <div className="text-[13px] text-muted">Vùng 1000 × 700 → 1400 × 900 · 6 điểm</div>
        </div>
        {MAP_MODES.map((m) => (
          <label key={m.id} className={`flex cursor-pointer gap-3 rounded-card p-3 ${mode === m.id ? 'ring-selected' : 'border border-line'}`}>
            <input type="radio" name="map" checked={mode === m.id} onChange={() => setMode(m.id)} className="mt-1 accent-accent" />
            <span className="flex flex-col gap-0.5"><span className="font-semibold">{m.title}</span><span className="text-xs text-muted">{m.hint}</span></span>
          </label>
        ))}
        <div className="flex justify-end gap-2 pt-1">
          <button onClick={onClose} className="h-9 rounded-control border border-line px-3.5 hover:bg-hover">Huỷ</button>
          <button onClick={() => onApply(MAP_MODES.find((m) => m.id === mode)!.title)} className="btn-primary h-9 rounded-control px-4 font-semibold">Áp dụng</button>
        </div>
      </div>
    </div>
  );
}

// "Trang web" giả phía sau overlay (giống artboard Overlay)
function FakePage() {
  return (
    <div className="pointer-events-none absolute top-[70px] left-[120px] h-[760px] w-[1200px] overflow-hidden rounded-lg bg-white text-[#16181D] shadow-[0_20px_60px_rgb(0_0_0/.45)]">
      <div className="flex h-11 items-center bg-[#EEF0F3] px-4">
        <span className="flex h-[26px] w-[520px] items-center rounded-lg bg-white px-3 text-xs text-[#5E6573]">facebook.com/[ten-nguoi-dung]</span>
      </div>
      <div className="grid h-[220px] place-items-center bg-[#E6E9EE] text-[#5E6573]">[Ảnh bìa]</div>
      <div className="-mt-[50px] flex items-end gap-5 px-[60px]">
        <span className="h-[150px] w-[150px] rounded-full border-[5px] border-white bg-[#CBD0D8]" />
        <div className="flex-1 pb-[18px]"><div className="text-[26px] font-bold">[Tên người dùng]</div><div className="text-[#5E6573]">[Thông tin phụ]</div></div>
        <div className="flex gap-2.5 pb-6">
          <span className="h-11 w-[140px] rounded-lg bg-[#E6E9EE]" />
          <span className="grid h-11 w-[148px] place-items-center rounded-lg bg-[#DCE6FF] font-semibold text-[#1E3FAF]">Nhắn tin</span>
        </div>
      </div>
      <div className="mx-[60px] mt-[30px] grid grid-cols-3 gap-4">
        {[0, 1, 2].map((i) => <span key={i} className="h-[220px] rounded-[10px] bg-[#F1F3F6]" />)}
      </div>
    </div>
  );
}
