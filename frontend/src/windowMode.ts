import { SetWindowMode } from '../wailsjs/go/main/App';

export type WindowMode = 'main' | 'overlay' | 'minibar';

// false khi xem trên trình duyệt (không có Go / Wails runtime)
export const inWails = 'runtime' in window;

// Đổi kích thước/vị trí cửa sổ (Go) rồi đổi màn qua hash URL – App.tsx chọn màn theo hash.
export async function setWindowMode(mode: WindowMode) {
  if (inWails) await SetWindowMode(mode);
  location.hash = mode === 'main' ? '' : mode;
}

// preview: xem trước kịch bản (con trỏ giả), không sửa điểm
export type OverlayTool = 'point' | 'drag' | 'region' | 'preview';

// Công cụ được chọn sẵn khi mở overlay (OverlayScreen đọc lúc mount).
export let overlayTool: OverlayTool = 'region';

export function openOverlay(tool: OverlayTool = 'point') {
  overlayTool = tool;
  return setWindowMode('overlay');
}
