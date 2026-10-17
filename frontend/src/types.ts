import { main } from '../wailsjs/go/models';
import { fmtAxis, fmtPoint } from './anchor';

// Neo của một trục (anchor.ts, khớp anchor.go): bỏ trống = start (cách cạnh trái / trên).
export type Anchor = 'start' | 'center' | 'end' | 'ratio';
// Toạ độ màn hình / cửa sổ.
export type Point = { x: number; y: number };
// Điểm lưu trong bước: neo theo cạnh vùng từng trục (khớp Point trong model.go).
export type StepPoint = Point & { ax?: Anchor; ay?: Anchor };
export type Rect = { x: number; y: number; w: number; h: number };
export type StepType = 'click' | 'double' | 'hold' | 'drag' | 'wait' | 'type' | 'key' | 'scroll' | 'open' | 'waitImage' | 'clickImage' | 'waitPixel';
export type ScrollDirection = 'up' | 'down' | 'left' | 'right';
export type MouseButton = 'left' | 'right' | 'center';
export type Modifier = 'cmd' | 'ctrl' | 'alt' | 'shift';
// chord: nhấn giữ cùng lúc (Ctrl+C+T) · sequence: giữ phím kèm, gõ lần lượt từng phím chính
export type KeyMode = 'chord' | 'sequence';

export interface Step {
  id: string;
  type: StepType;
  points: StepPoint[]; // drag: points[0] -> points[1]; B lệch so với A (dragRelative) không dùng neo
  button: MouseButton;
  clicks: number;
  clickInterval: number; // ms
  holdMs: number;
  dragMs: number;
  dragFromPrev: boolean; // drag: bắt đầu từ chỗ chuột dừng ở bước chuột trước (bỏ qua points[0])
  dragRelative: boolean; // drag: points[1] là độ lệch (dx, dy) so với điểm bắt đầu
  text: string; // type: chuỗi cần gõ
  keys: string[]; // key: phím chính theo tên robotgo (enter, tab, a, f5…)
  keyMode: KeyMode; // key: chỉ khác nhau khi có ≥ 2 phím chính
  modifiers: Modifier[]; // key: tổ hợp, vd ['cmd'] + 'v'
  delayAfter: number; // ms
  jitterMs: number; // ngẫu nhiên ± ms cho delayAfter và clickInterval
  jitterPx: number; // ngẫu nhiên ± px cho mỗi lần di chuột tới điểm
  enabled: boolean;
  direction?: ScrollDirection; // scroll
  amount?: number; // scroll: số nấc
  target?: string; // open: URL / đường dẫn / tên app, nhận {{biến}}
  app?: string; // open: mở bằng app này, '' = mặc định của máy
  // chờ ảnh / click theo ảnh / chờ pixel (bước 25, imagestep.go)
  template?: string; // ảnh mẫu PNG base64
  templateScale?: number; // px ảnh mẫu / 1 đơn vị toạ độ lúc chụp (Retina 2)
  threshold?: number; // % khớp tối thiểu
  timeoutMs?: number; // chờ tối đa
  intervalMs?: number; // kiểm tra mỗi
  searchArea?: 'screen'; // bỏ trống = vùng của kịch bản
  imageSource?: 'prev' | 'search'; // clickImage
  dx?: number; // clickImage: lệch so với tâm ảnh
  dy?: number;
  findBy?: 'text'; // waitImage / clickImage tự tìm: bỏ trống = ảnh mẫu, text = tìm chữ (ocr.go)
  findText?: string; // chữ cần tìm, nhận {{biến}}
  textExact?: boolean; // cả dòng đúng chuỗi (mặc định: chứa chuỗi)
  textNoAccent?: boolean; // bỏ qua dấu
  occurrence?: number; // chỗ thứ mấy theo thứ tự đọc, 0 = chỗ đầu
  color?: string; // waitPixel: #RRGGBB
  tolerance?: number; // waitPixel: sai số mỗi kênh màu
}

// Kiểm tra lúc type-check: Step phải khớp struct Go (model.go) – thiếu/thừa/sai kiểu field là tsc báo lỗi.
type Fields<T> = { [K in keyof T as T[K] extends Function ? never : K]: T[K] };
type Same<A, B> = [A] extends [B] ? ([B] extends [A] ? true : false) : false;
const stepMatchesGo: [Step extends Fields<main.Step> ? true : false, Same<keyof Step, keyof Fields<main.Step>>] = [true, true];
void stepMatchesGo;

// Kịch bản – khớp struct Script trong model.go. Toạ độ điểm tương đối theo region (null = toàn màn hình).
export interface Script {
  version: number; // định dạng file (storage.go) – Go ghi khi lưu
  name: string;
  region: Rect | null; // toạ độ màn hình thật (robotgo)
  loops: number; // 0 = vô hạn
  maxMinutes: number; // 0 = không giới hạn
  steps: Step[];
  anchorMode?: AnchorMode; // điểm mới đặt neo theo: bỏ trống = góc gần nhất, start = luôn trên trái
  inputs?: ScriptInput[]; // giá trị hỏi khi bấm Chạy (bước 21b, inputs.go)
}
export type AnchorMode = 'start';
const scriptMatchesGo: Same<keyof Script, keyof Fields<main.Script>> = true;
void scriptMatchesGo;

// ---- Input của kịch bản (bước 21b) – khớp ScriptInput trong inputs.go ----
export type InputType = 'text' | 'multiline' | 'number' | 'ms' | 'bool' | 'select';
export type InputValue = string | number | boolean;
// Trường ghim được (pinFields trong inputs.go): số → input number / ms, enabled → bool.
export type PinField = 'clicks' | 'clickInterval' | 'holdMs' | 'dragMs' | 'delayAfter' | 'jitterMs' | 'jitterPx' | 'amount' | 'timeoutMs' | 'threshold' | 'enabled';
export interface InputTarget { stepId: string; field: PinField }
export interface ScriptInput {
  key: string; // chữ không dấu, số, _ – dùng trong {{key}}
  label: string;
  type: InputType;
  default?: InputValue; // trống = không có mặc định
  required?: boolean;
  min?: number;
  max?: number;
  options?: string[]; // select
  hint?: string; // gợi ý dưới ô trong form Chạy
  targets?: InputTarget[]; // trường được ghim
}
const inputMatchesGo: Same<keyof ScriptInput, keyof Fields<main.ScriptInput>> = true;
void inputMatchesGo;

// Tạo bước với giá trị trống, ghi đè bằng p (dùng cho dữ liệu mẫu / demo).
export const makeStep = (p: Partial<Step> & Pick<Step, 'id' | 'type'>): Step => ({
  points: [], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false,
  text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true, ...p,
});

export const STEP_LABEL: Record<StepType, string> = {
  click: 'Click',
  double: 'Double click',
  hold: 'Giữ chuột',
  drag: 'Kéo thả',
  wait: 'Chờ',
  type: 'Gõ chữ',
  key: 'Phím tắt',
  scroll: 'Cuộn',
  open: 'Mở URL / app',
  waitImage: 'Chờ ảnh',
  clickImage: 'Click theo ảnh',
  waitPixel: 'Chờ màu pixel',
};

export const SCROLL_ARROW: Record<ScrollDirection, string> = { up: '↑', down: '↓', left: '←', right: '→' };
// Quãng ước tính một nấc cuộn (khớp scrollNotchPx – input.go): macOS gửi đúng 100 px, Windows một nấc con lăn
// (Chrome / Edge ≈ 100 px, app khác có thể khác).
export const SCROLL_NOTCH_PX = 100;
// Trang dịch đi bao nhiêu (px, theo hướng cuộn): xuống = +y, phải = +x.
export const scrollDelta = (s: Pick<Step, 'direction' | 'amount'>): Point => {
  const d = (s.amount ?? 0) * SCROLL_NOTCH_PX;
  return { down: { x: 0, y: d }, up: { x: 0, y: -d }, right: { x: d, y: 0 }, left: { x: -d, y: 0 } }[s.direction ?? 'down'];
};

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

export const sampleSteps: Step[] = [
  { id: '1', type: 'click', points: [{ x: 412, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '2', type: 'double', points: [{ x: 640, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '3', type: 'click', points: [{ x: 812, y: 436 }], button: 'left', clicks: 5, clickInterval: 120, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 800, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '4', type: 'drag', points: [{ x: 300, y: 600 }, { x: 900, y: 600 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 400, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 1000, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '5', type: 'wait', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 2000, jitterMs: 0, jitterPx: 0, enabled: false },
  { id: '6', type: 'click', points: [{ x: 1100, y: 520 }], button: 'right', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '7', type: 'type', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: 'xin chao 123', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '8', type: 'key', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: ['enter'], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
];

export const fmtMs = (ms: number) => (ms >= 1000 ? `${ms / 1000} s` : `${ms} ms`);
const fmtSigned = (v: number) => (v > 0 ? `+${v}` : `${v}`);
const fmtPair = (p: StepPoint) => `${fmtAxis(p.x, p.ax)}, ${fmtAxis(p.y, p.ay)}`;
export const fmtPos = (s: Step) =>
  s.type === 'drag'
    ? `${s.dragFromPrev ? 'Từ điểm trước' : `A ${fmtPair(s.points[0])}`} → ${
      s.dragRelative ? `Δ ${fmtSigned(s.points[1].x)},${fmtSigned(s.points[1].y)}` : `B ${fmtPair(s.points[1])}`}`
    : s.type === 'type' ? `"${s.text}"`
    : s.type === 'key' ? fmtKey(s)
    : s.type === 'scroll' ? `${SCROLL_ARROW[s.direction ?? 'down']} ${s.amount ?? 0} nấc ≈ ${(s.amount ?? 0) * SCROLL_NOTCH_PX} px ${s.points[0] ? `tại ${fmtPoint(s.points[0])}` : 'tại chuột'}`
    : s.type === 'open' ? `${s.target || '—'}${s.app ? ` · ${s.app}` : ''}`
    : s.type === 'waitImage' || (s.type === 'clickImage' && s.imageSource === 'search')
      ? `${s.type === 'clickImage' ? 'tự tìm · ' : ''}${s.findBy === 'text'
        ? (s.findText?.trim() ? `chữ “${s.findText}”${(s.occurrence ?? 0) > 1 ? ` (chỗ thứ ${s.occurrence})` : ''} · tối đa ${fmtMs(s.timeoutMs ?? 0)}` : 'chưa nhập chữ cần tìm')
        : s.template ? `khớp ≥ ${s.threshold}% · tối đa ${fmtMs(s.timeoutMs ?? 0)}` : 'chưa có ảnh mẫu'}`
    : s.type === 'clickImage' ? `chỗ bước trước thấy ảnh${s.dx || s.dy ? ` ${fmtSigned(s.dx ?? 0)},${fmtSigned(s.dy ?? 0)}` : ''}`
    : s.type === 'waitPixel' && s.points[0] ? `${fmtPoint(s.points[0])} = ${s.color} ±${s.tolerance} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`
    : s.points[0] ? fmtPoint(s.points[0]) : '—';
