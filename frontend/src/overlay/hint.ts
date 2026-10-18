// Chữ gợi ý một dòng trên overlay theo công cụ / trạng thái đang làm.
import type { OverlayTool } from '../windowMode';
import { STEP_LABEL, Step, StepType } from '../types';

export type HintState = {
  focusStep: Step | null; // xem riêng một node
  repicking: boolean; // đang "Chọn lại trên màn hình" khi xem riêng
  dragA: boolean; // kéo thả: đã có A, chờ B
  tool: OverlayTool;
  captured: boolean; // khoanh ảnh: đã có khung, chờ Dùng ảnh này
  captureRef: string | null; // mã node sẽ nhận ảnh mẫu, null = tạo node Chờ ảnh mới
  selected: number; // số điểm đang chọn (Chọn nhiều)
  picking: boolean; // đang chọn lại điểm của node có sẵn
  placing: StepType; // loại node sẽ thêm khi đặt điểm
  modKey: string; // ⌘ / Ctrl
};

export function overlayHint(h: HintState): string {
  const f = h.focusStep;
  if (f) {
    if (!h.repicking) return 'Kéo marker để sửa · Esc quay lại';
    if (f.type === 'drag' && !f.dragFromPrev) return h.dragA ? 'Click điểm B (nơi thả)' : 'Click điểm A (nơi bắt đầu kéo)';
    return 'Click (hoặc fn+F6) để đặt lại điểm';
  }
  switch (h.tool) {
    case 'preview':
      return 'Con trỏ giả, không click thật';
    case 'region':
      return 'Kéo để khoanh vùng';
    case 'capture':
      return h.captured ? 'Enter: dùng ảnh này · kéo lại hoặc Esc: khoanh lại'
        : `Kéo khung sát quanh nút / chữ cần chờ${h.captureRef ? ` (${h.captureRef})` : ' – tạo bước Chờ ảnh mới'}`;
    case 'select':
      return h.selected
        ? `${h.selected} điểm · kéo khung / núm để biến đổi · ${h.modKey}+núm: xiên, phối cảnh`
        : `Click, Shift+click hoặc kéo khung để chọn · ${h.modKey}A tất cả`;
    case 'drag':
      return h.dragA ? 'Click điểm B · Esc huỷ A' : `Click điểm A${h.picking ? ' (chọn lại)' : ''}`;
  }
  return h.picking ? 'Click để đặt lại điểm' : `Click / F6: thêm ${STEP_LABEL[h.placing].toLowerCase()} · kéo marker để dời`;
}
