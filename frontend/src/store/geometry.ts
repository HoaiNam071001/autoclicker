// Toạ độ điểm của node: khung vùng / màn hình, đổi điểm lưu (neo) ↔ toạ độ màn hình, điểm đầu / cuối của bước chuột.
// Hàm đọc state, không đổi state – dùng chung cho nodes / points / multi và overlay, xem trước.
import { anchorPoint, insideBounds, nearestAnchor, resolvePoint } from '../anchor';
import { flowOrder } from '../flow/model';
import { Point, Rect, Script, Step, StepPoint } from '../types';
import { getEditor } from './editor';

// Khung tính toạ độ điểm: vùng, không có thì màn hình chính (bounds trong anchor.go).
export const boundsOf = (region: Rect | null): Rect => region ?? getEditor().screen;

// Điểm đã lưu (neo theo vùng) → toạ độ màn hình thật.
export const toScreen = (p: StepPoint, region: Rect | null): Point => resolvePoint(p, boundsOf(region));

// Toạ độ màn hình → điểm lưu: like = điểm cũ để giữ neo (kéo marker, chọn lại), không có thì neo theo cài đặt kịch bản
// (góc gần nhất / luôn trên trái). null nếu nằm ngoài vùng.
export const toRelative = (p: Point, r: Rect | null, like?: StepPoint): StepPoint | null => {
  const b = boundsOf(r);
  if (!insideBounds(p, b)) return null;
  if (like) return anchorPoint(p, b, like.ax, like.ay);
  return getEditor().script.anchorMode === 'start' ? anchorPoint(p, b, 'start', 'start') : nearestAnchor(p, b);
};

// Điểm giữa vùng / màn hình (điểm mặc định của chờ pixel, cuộn tại điểm).
export const centerPoint = (region: Rect | null): StepPoint => {
  const b = boundsOf(region);
  return nearestAnchor({ x: b.x + Math.round(b.w / 2), y: b.y + Math.round(b.h / 2) }, b);
};

export const pointKey = (id: string, index: number) => `${id}:${index}`;

// Điểm chọn được của một node: A "từ điểm trước" tự tính nên không chọn / biến đổi.
export const selectableIndexes = (s: Step) =>
  s.points.length === 0 || s.type === 'wait' || s.type === 'type' || s.type === 'key' || s.type === 'open' ? []
    : s.type === 'drag' ? (s.dragFromPrev ? [1] : [0, 1]) : [0];

export type Ends = { start: Point | null; end: Point | null; prevStep: string | null };

// Toạ độ màn hình (danh nghĩa, chưa tính ngẫu nhiên) điểm đầu / cuối của các node chuột đang bật. "Điểm trước" = chỗ chuột
// dừng ở node chuột gần nhất phía trước theo thứ tự hiển thị (flowOrder – nhánh chính); lúc chạy thật theo đường đi thực tế.
// prevStep: id node cung cấp điểm trước cho node kéo "từ điểm trước".
export function endpoints(sc: Pick<Script, 'steps' | 'start' | 'edges'>, region: Rect | null): Map<string, Ends> {
  const out = new Map<string, Ends>();
  const byId = new Map(sc.steps.map((s) => [s.id, s]));
  let last: Point | null = null;
  let lastStep: string | null = null;
  for (const id of flowOrder(sc)) {
    const s = byId.get(id)!;
    if (!s.enabled || !s.points.length) continue;
    if (s.type === 'waitPixel') { // có điểm (marker) nhưng không di chuột
      const p = toScreen(s.points[0], region);
      out.set(s.id, { start: p, end: p, prevStep: null });
    } else if (s.type === 'click' || s.type === 'double' || s.type === 'hold' || s.type === 'scroll') {
      const p = toScreen(s.points[0], region);
      out.set(s.id, { start: p, end: p, prevStep: null });
      last = p;
      lastStep = s.id;
    } else if (s.type === 'drag' && s.points.length >= 2) {
      const start: Point | null = s.dragFromPrev ? last : toScreen(s.points[0], region);
      const d = s.points[1];
      const end: Point | null = s.dragRelative ? (start && { x: start.x + d.x, y: start.y + d.y }) : toScreen(d, region);
      out.set(s.id, { start, end, prevStep: s.dragFromPrev ? lastStep : null });
      if (end) {
        last = end;
        lastStep = s.id;
      }
    }
  }
  return out;
}

// Mã node để hiện trong lời báo (n2), không thấy thì "?"
export const refOf = (id: string) => getEditor().script.steps.find((s) => s.id === id)?.ref || '?';
