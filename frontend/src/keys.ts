// Tên phím robotgo (gợi ý / kiểm tra trong form Phím tắt) và chữ hiển thị tổ hợp phím.
import type { Step } from './types';

// Tên phím robotgo hỗ trợ (ngoài các ký tự đơn như a, 1, /) – gợi ý trong form và kiểm tra.
export const KEY_NAMES = [
  'enter', 'tab', 'esc', 'space', 'backspace', 'delete', 'up', 'down', 'left', 'right',
  'home', 'end', 'pageup', 'pagedown', 'insert', 'capslock', 'printscreen', 'menu',
  ...Array.from({ length: 24 }, (_, i) => `f${i + 1}`),
  'num0', 'num1', 'num2', 'num3', 'num4', 'num5', 'num6', 'num7', 'num8', 'num9', 'num_enter',
  'audio_mute', 'audio_vol_down', 'audio_vol_up', 'audio_play', 'audio_next', 'audio_prev',
];
export const isValidKey = (k: string) => [...k].length === 1 || KEY_NAMES.includes(k);

// Mô tả tiếng Việt cho gợi ý tên phím trong form.
const FIXED_LABEL: Record<string, string> = {
  enter: 'Enter – xuống dòng / xác nhận', tab: 'Tab – sang ô kế', esc: 'Esc – thoát / huỷ', space: 'Phím cách',
  backspace: 'Backspace – xoá lùi', delete: 'Delete – xoá tới', up: '↑ Mũi tên lên', down: '↓ Mũi tên xuống',
  left: '← Mũi tên trái', right: '→ Mũi tên phải', home: 'Home – đầu dòng', end: 'End – cuối dòng',
  pageup: 'Page Up – lên một trang', pagedown: 'Page Down – xuống một trang', insert: 'Insert', capslock: 'Caps Lock – bật / tắt chữ hoa',
  printscreen: 'Print Screen – chụp màn hình', menu: 'Menu ngữ cảnh', num_enter: 'Enter bàn phím số',
  audio_mute: 'Tắt / bật tiếng', audio_vol_down: 'Giảm âm lượng', audio_vol_up: 'Tăng âm lượng',
  audio_play: 'Phát / tạm dừng nhạc', audio_next: 'Bài kế tiếp', audio_prev: 'Bài trước',
};
export const keyLabel = (k: string) =>
  FIXED_LABEL[k] ?? (/^f\d+$/.test(k) ? `Phím chức năng ${k.toUpperCase()}` : /^num\d$/.test(k) ? `Số ${k.slice(3)} bàn phím số` : '');

const cap = (k: string) => k[0].toUpperCase() + k.slice(1);
// chord: "Ctrl + C + T" · sequence: "Ctrl + C → T"
export const fmtKey = (s: Step) => {
  const keys = s.keys.map(cap).join(s.keyMode === 'sequence' ? ' → ' : ' + ');
  return [...s.modifiers.map(cap), keys].filter(Boolean).join(' + ') || '—';
};
