import { SetWindowMode } from '../wailsjs/go/main/App';

export type WindowMode = 'main' | 'overlay' | 'minibar';

// false khi xem trên trình duyệt (không có Go / Wails runtime)
export const inWails = 'runtime' in window;

// Đổi kích thước/vị trí cửa sổ (Go) rồi đổi màn qua hash URL – App.tsx chọn màn theo hash.
export async function setWindowMode(mode: WindowMode) {
  if (inWails) await SetWindowMode(mode);
  location.hash = mode === 'main' ? '' : mode;
}
