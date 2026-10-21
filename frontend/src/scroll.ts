// Cuộn (node Cuộn chuột): mũi tên hướng, quãng một nấc, độ dịch trang – dùng ở inspector, overlay, xem trước.
import type { Point, ScrollDirection, Step } from './types';

export const SCROLL_ARROW: Record<ScrollDirection, string> = { up: '↑', down: '↓', left: '←', right: '→' };
// Quãng ước tính một nấc cuộn (khớp scrollNotchPx – input.go): macOS gửi đúng 100 px, Windows một nấc con lăn
// (Chrome / Edge ≈ 100 px, app khác có thể khác).
export const SCROLL_NOTCH_PX = 100;
// Trang dịch đi bao nhiêu (px, theo hướng cuộn): xuống = +y, phải = +x.
export const scrollDelta = (s: Pick<Step, 'direction' | 'amount'>): Point => {
  const d = (s.amount ?? 0) * SCROLL_NOTCH_PX;
  return { down: { x: 0, y: d }, up: { x: 0, y: -d }, right: { x: d, y: 0 }, left: { x: -d, y: 0 } }[s.direction ?? 'down'];
};
