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
// select (chọn nhiều điểm để biến đổi – bước 30) và capture (khoanh ảnh mẫu – bước 25): giao diện ở #demo-overlay, chưa bật
export type OverlayTool = 'point' | 'drag' | 'region' | 'preview' | 'select' | 'capture';

// Công cụ được chọn sẵn khi mở overlay và bước đang xem riêng (null = xem tất cả) – OverlayScreen đọc lúc mount.
export let overlayTool: OverlayTool = 'region';
export let overlayFocus: string | null = null;

export function openOverlay(tool: OverlayTool = 'point', focus: string | null = null) {
  overlayTool = tool;
  overlayFocus = focus;
  return setWindowMode('overlay');
}
