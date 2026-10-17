// Chọn nhiều điểm + biến đổi hàng loạt (bước 30) và đổi vùng kèm cách xử lý các điểm.
import { anchorPoint, insideBounds } from '../anchor';
import { apply, IDENTITY, Mat, rectToRect } from '../transform';
import { Point, Rect, Script, Step } from '../types';
import { getEditor, setEditor, setNotice, updateScript } from './editor';
import { boundsOf, endpoints, pointKey, selectableIndexes, toScreen } from './geometry';

const script = () => getEditor().script;

export const setMulti = (multi: string[]) => setEditor({ multi });

// Thẻ node: Cmd/Shift+click thêm / bỏ mọi điểm của node vào chọn nhiều.
export function toggleMultiStep(id: string) {
  const s = script().steps.find((x) => x.id === id);
  if (!s) return;
  const { multi } = getEditor();
  const keys = selectableIndexes(s).map((j) => pointKey(id, j));
  const has = keys.some((k) => multi.includes(k));
  setEditor({ multi: has ? multi.filter((k) => !keys.includes(k)) : [...multi, ...keys] });
}

// Áp f (toạ độ màn hình → toạ độ màn hình) lên các điểm được chọn (pick), đổi khung tính từ vùng from sang to.
// Điểm thường: ghi lại theo neo hiện có trong khung mới. B "lệch so với A": véc-tơ đổi theo phần tuyến tính của f
// (f(A + Δ) − f(A), không cộng phần dời). Trả số điểm rơi ra ngoài khung mới – khác 0 thì người gọi không nên áp dụng.
function mapPoints(sc: Script, from: Rect | null, to: Rect | null, f: (p: Point) => Point, pick: (id: string, j: number) => boolean) {
  const ends = endpoints(sc, from);
  const nb = boundsOf(to);
  let outside = 0;
  const round = (p: Point) => ({ x: Math.round(p.x), y: Math.round(p.y) });
  const steps = sc.steps.map((s) => {
    if (!selectableIndexes(s).some((j) => pick(s.id, j))) return s;
    const points = s.points.map((p, j) => {
      if (!pick(s.id, j) || !selectableIndexes(s).includes(j)) return p;
      if (s.type === 'drag' && j === 1 && s.dragRelative) {
        const a = ends.get(s.id)?.start;
        if (!a) return p;
        const fa = f(a), fb = f({ x: a.x + p.x, y: a.y + p.y });
        return round({ x: fb.x - fa.x, y: fb.y - fa.y });
      }
      const q = round(f(toScreen(p, from)));
      if (!insideBounds(q, nb)) outside++;
      return anchorPoint(q, nb, p.ax, p.ay);
    });
    return { ...s, points };
  });
  return { steps, outside };
}

// Kịch bản như sẽ thành nếu áp m lên các điểm đang chọn – overlay vẽ xem trước bằng đúng logic ghi thật.
export function transformedSteps(m: Mat) {
  const sc = script();
  const sel = new Set(getEditor().multi);
  return mapPoints(sc, sc.region, sc.region, (p) => apply(m, p), (id, j) => sel.has(pointKey(id, j)));
}

// Biến đổi các điểm đang chọn nhiều bằng ma trận m (toạ độ màn hình). Điểm bị đẩy ra ngoài vùng → báo, không áp dụng.
export function transformSelected(m: Mat, what: string): boolean {
  const n = getEditor().multi.length;
  if (!n) return false;
  const r = transformedSteps(m);
  if (r.outside) {
    setNotice(`${r.outside} điểm sẽ ra ngoài vùng – không áp dụng ${what}`);
    return false;
  }
  updateScript({ steps: r.steps });
  setNotice(`Đã ${what} ${n} điểm – ${navigator.platform.startsWith('Mac') ? '⌘' : 'Ctrl+'}Z để hoàn tác`);
  return true;
}

// Đặt toạ độ màn hình mới cho từng điểm đang chọn (căn chỉnh / chia đều / lật): next[khoá] → điểm mới.
// Hai lượt: ghi điểm thường trước, rồi tính lại B "lệch so với A" theo điểm bắt đầu mới (A có thể vừa dời).
export function moveSelectedTo(next: Map<string, Point>, what: string) {
  const sc = script();
  const b = boundsOf(sc.region);
  const round = (p: Point) => ({ x: Math.round(p.x), y: Math.round(p.y) });
  const isRel = (s: Step, j: number) => s.type === 'drag' && j === 1 && s.dragRelative;
  let outside = 0;
  const first = sc.steps.map((s) => ({
    ...s,
    points: s.points.map((p, j) => {
      const q = next.get(pointKey(s.id, j));
      if (!q || isRel(s, j)) return p;
      if (!insideBounds(round(q), b)) outside++;
      return anchorPoint(round(q), b, p.ax, p.ay);
    }),
  }));
  const ends = endpoints({ ...sc, steps: first }, sc.region);
  const steps = first.map((s) => ({
    ...s,
    points: s.points.map((p, j) => {
      const q = next.get(pointKey(s.id, j));
      const start = ends.get(s.id)?.start;
      if (!q || !isRel(s, j) || !start) return p;
      if (!insideBounds(round(q), b)) outside++;
      return { x: Math.round(q.x) - start.x, y: Math.round(q.y) - start.y };
    }),
  }));
  if (outside) return setNotice(`${outside} điểm sẽ ra ngoài vùng – không ${what}`);
  updateScript({ steps });
  setNotice(`Đã ${what} ${next.size} điểm`);
}

// ---- Đổi vùng: cách xử lý các điểm (bước 30) ----
// anchor: giữ theo neo (bước 29) · scale: giữ tỷ lệ vị trí trong vùng · screen: giữ nguyên chỗ trên màn hình
export type RegionMapMode = 'anchor' | 'scale' | 'screen';

export function regionMapping(to: Rect | null, mode: RegionMapMode) {
  const sc = script();
  const f = mode === 'scale' ? rectToRect(boundsOf(sc.region), boundsOf(to)) : IDENTITY;
  return mapPoints(sc, sc.region, to, (p) => apply(f, p), () => mode !== 'anchor');
}

// Số điểm rơi ra ngoài vùng mới với từng cách – hiện trong hộp chọn.
export function mappingOutside(to: Rect | null, mode: RegionMapMode) {
  if (mode !== 'anchor') return regionMapping(to, mode).outside;
  const b = boundsOf(to);
  let n = 0;
  endpoints(script(), to).forEach((e) => [e.start, e.end].forEach((p, j) => {
    if (p && (j === 0 || e.end !== e.start) && !insideBounds(p, b)) n++;
  }));
  return n;
}

export function setRegionMapped(to: Rect | null, mode: RegionMapMode) {
  const { steps } = regionMapping(to, mode);
  updateScript({ region: to, steps });
  setNotice(to ? `Đã đổi vùng ${to.w} × ${to.h}` : 'Đã về toàn màn hình');
}

// Đổi vùng; điểm lưu theo neo nên mỗi điểm bám theo cạnh / góc nó neo. null = toàn màn hình.
export const setRegion = (region: Rect | null) => updateScript({ region });

// Kịch bản có điểm nào phụ thuộc vùng không (không có thì đổi vùng khỏi hỏi).
export const hasPoints = () => script().steps.some((s) => selectableIndexes(s).length > 0);
