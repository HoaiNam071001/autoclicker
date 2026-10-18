import { model } from '../wailsjs/go/models';
import { fmtAxis, fmtPoint } from './anchor';

// Neo của một trục (anchor.ts, khớp internal/script/anchor.go): bỏ trống = start (cách cạnh trái / trên).
export type Anchor = 'start' | 'center' | 'end' | 'ratio';
// Toạ độ màn hình / cửa sổ.
export type Point = { x: number; y: number };
// Điểm lưu trong bước: neo theo cạnh vùng từng trục (khớp Point trong model.go).
export type StepPoint = Point & { ax?: Anchor; ay?: Anchor };
export type Rect = { x: number; y: number; w: number; h: number };
// Node thất bại (không thấy ảnh / chữ / màu, lỗi): thử lại thêm count lần, nghỉ delayMs; hết lượt → cổng fail. Khớp Retry trong model.go.
export type Retry = { count?: number; delayMs?: number };
// Node hay "thất bại" theo nghĩa người dùng hiểu – mới hiện phần Khi thất bại (thử lại)
export const CAN_FAIL: StepType[] = ['waitImage', 'clickImage', 'waitPixel', 'open'];

// ---- Flow (bước 20, flow.go) ----
// Cạnh: node from đi ra cổng port thì chạy node to. Mỗi cổng tối đa một cạnh. Khớp Edge trong model.go.
export interface Edge { id: string; from: string; port: string; to: string }
export type StepStatus = 'success' | 'fail' | 'skipped' | 'cancelled';
export type ErrorCode = 'not_found' | 'timeout' | 'permission' | 'invalid_input' | 'cancelled' | 'end' | 'internal';
// Kết quả một lần chạy node (event engine:node) – khớp StepResult trong flow.go.
export interface StepResult {
  nodeId: string;
  ref: string;
  status: StepStatus;
  port?: string;
  output?: Record<string, unknown>;
  error?: { code: ErrorCode; message: string };
  logs?: { atMs: number; level: 'info' | 'warn' | 'error'; msg: string }[];
  attempts: number;
  loop: number; // vòng thứ mấy, từ 1
  shot?: string; // ảnh lúc thất bại (lịch sử chạy – RunShot)
  startedAt: string;
  durationMs: number;
}
const edgeMatchesGo: Same<keyof Edge, keyof Fields<model.Edge>> = true;
const resultMatchesGo: Same<keyof StepResult, keyof Fields<model.StepResult>> = true;
void edgeMatchesGo, resultMatchesGo;

export type ActionType = 'click' | 'double' | 'hold' | 'drag' | 'wait' | 'type' | 'key' | 'scroll' | 'open' | 'waitImage' | 'clickImage' | 'waitPixel';
// Node điều khiển (bước 20 F3, internal/engine/nodes.go): chỉ chọn cổng ra, inspector riêng (components/flow/nodes/ControlFields)
export type ControlType = 'if' | 'switch' | 'switchScreen' | 'counter' | 'end' | 'setVar';
export type StepType = ActionType | ControlType;
export const CONTROL_TYPES: ControlType[] = ['if', 'switch', 'switchScreen', 'counter', 'end', 'setVar'];
export const isControl = (t: StepType): t is ControlType => (CONTROL_TYPES as StepType[]).includes(t);

// Điều kiện của node Nếu / case của Switch theo màn hình – khớp Cond trong model.go (internal/engine/cond.go).
export type CondKind = 'image' | 'text' | 'pixel' | 'value' | 'all' | 'any';
export type CondOp = 'eq' | 'ne' | 'contains' | 'gt' | 'lt' | 'empty';
export interface Cond {
  kind: CondKind;
  not?: boolean; // đảo kết quả
  template?: string; templateScale?: number; threshold?: number; searchArea?: string; // image (như Chờ ảnh)
  findText?: string; textExact?: boolean; textNoAccent?: boolean; occurrence?: number; // text, nhận {{biến}}
  points?: StepPoint[]; color?: string; tolerance?: number; // pixel: points[0]
  left?: string; op?: CondOp; right?: string; // value: so 2 vế, nhận {{biến}}
  conds?: Cond[]; // all (và) / any (hoặc)
}
// Nhánh của switch / switchScreen – cổng ra "case:" + key. Khớp Case trong model.go.
export type CaseOp = 'eq' | 'contains' | 'regex';
export interface Case {
  key: string;
  label?: string; // tên trên cổng, bỏ trống = key
  op?: CaseOp; // switch
  match?: string; // switch, nhận {{biến}}
  cond?: Cond; // switchScreen
}
export type ScrollDirection = 'up' | 'down' | 'left' | 'right';
export type MouseButton = 'left' | 'right' | 'center';
export type Modifier = 'cmd' | 'ctrl' | 'alt' | 'shift';
// chord: nhấn giữ cùng lúc (Ctrl+C+T) · sequence: giữ phím kèm, gõ lần lượt từng phím chính
export type KeyMode = 'chord' | 'sequence';

// Node của flow (bước 20) – khớp Step trong model.go.
export interface Step {
  id: string;
  type: StepType;
  ref: string; // mã ngắn n1, n2… (duy nhất) – biến {{n2.output.x}}
  name?: string; // tên trên thẻ, bỏ trống = tên theo loại
  pos?: Point; // vị trí thẻ trên canvas
  breakpoint?: boolean; // debug: chạy từ app thì tạm dừng trước node này
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
  // chờ ảnh / click theo ảnh / chờ pixel (bước 25, internal/engine/imagestep.go)
  template?: string; // ảnh mẫu PNG base64
  templateScale?: number; // px ảnh mẫu / 1 đơn vị toạ độ lúc chụp (Retina 2)
  threshold?: number; // % khớp tối thiểu
  timeoutMs?: number; // chờ tối đa
  intervalMs?: number; // kiểm tra mỗi
  searchArea?: 'screen'; // bỏ trống = vùng của kịch bản
  imageSource?: 'node' | 'search'; // clickImage: chỗ node imageFrom thấy / tự tìm
  imageFrom?: string; // clickImage imageSource node: id node lấy chỗ thấy
  dx?: number; // clickImage: lệch so với tâm ảnh
  dy?: number;
  findBy?: 'text'; // waitImage / clickImage tự tìm: bỏ trống = ảnh mẫu, text = tìm chữ (ocr.go)
  findText?: string; // chữ cần tìm, nhận {{biến}}
  textExact?: boolean; // cả dòng đúng chuỗi (mặc định: chứa chuỗi)
  textNoAccent?: boolean; // bỏ qua dấu
  occurrence?: number; // chỗ thứ mấy theo thứ tự đọc, 0 = chỗ đầu
  retry?: Retry; // thất bại thì thử lại; hết lượt → cổng fail (internal/engine/flowrun.go execNode)
  color?: string; // waitPixel: #RRGGBB
  tolerance?: number; // waitPixel: sai số mỗi kênh màu
  // node điều khiển (bước 20 F3)
  cond?: Cond; // if
  cases?: Case[]; // switch / switchScreen
  value?: string; // switch: giá trị cần xét; setVar: giá trị gán / cộng
  max?: number; // counter: số lần đi cổng next tối đa
  endStatus?: 'success' | 'fail'; // end
  message?: string; // end: lời nhắn / lời báo lỗi
  varName?: string; // setVar
  varOp?: 'set' | 'add'; // setVar
}

// Kiểm tra lúc type-check: Step phải khớp struct Go (model.go) – thiếu/thừa/sai kiểu field là tsc báo lỗi.
// bỏ hàm (convertValues của class sinh từ Go) ở mọi tầng – Step có điều kiện / case lồng nhau
type Fields<T> = { [K in keyof T as T[K] extends Function ? never : K]: DeepFields<T[K]> };
type DeepFields<T> = T extends (infer U)[] ? DeepFields<U>[] : T extends object ? Fields<T> : T;
type Same<A, B> = [A] extends [B] ? ([B] extends [A] ? true : false) : false;
const stepMatchesGo: [Step extends Fields<model.Step> ? true : false, Same<keyof Step, keyof Fields<model.Step>>] = [true, true];
void stepMatchesGo;
const condMatchesGo: [Same<keyof Cond, keyof Fields<model.Cond>>, Same<keyof Case, keyof Fields<model.Case>>] = [true, true];
void condMatchesGo;

// Kịch bản – khớp struct Script trong model.go. Toạ độ điểm tương đối theo region (null = toàn màn hình).
export interface Script {
  version: number; // định dạng file (storage.go) – Go ghi khi lưu
  name: string;
  region: Rect | null; // toạ độ màn hình thật (robotgo)
  loops: number; // 0 = vô hạn
  maxMinutes: number; // 0 = không giới hạn
  steps: Step[]; // node của flow – thứ tự mảng không quyết định thứ tự chạy
  start: string; // id node bắt đầu
  edges: Edge[];
  anchorMode?: AnchorMode; // điểm mới đặt neo theo: bỏ trống = góc gần nhất, start = luôn trên trái
  inputs?: ScriptInput[]; // giá trị hỏi khi bấm Chạy (bước 21b, internal/script/inputs.go)
}
export type AnchorMode = 'start';
const scriptMatchesGo: Same<keyof Script, keyof Fields<model.Script>> = true;
void scriptMatchesGo;

// ---- Input của kịch bản (bước 21b) – khớp ScriptInput trong internal/model/inputs.go ----
export type InputType = 'text' | 'multiline' | 'number' | 'ms' | 'bool' | 'select';
export type InputValue = string | number | boolean;
// Trường ghim được (pinFields trong internal/script/inputvalue.go): số → input number / ms, enabled → bool.
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
const inputMatchesGo: Same<keyof ScriptInput, keyof Fields<model.ScriptInput>> = true;
void inputMatchesGo;

// Tạo bước với giá trị trống, ghi đè bằng p (dùng cho dữ liệu mẫu / demo).
export const makeStep = (p: Partial<Step> & Pick<Step, 'id' | 'type'>): Step => ({
  ref: '', points: [], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false,
  text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true, ...p,
});

// Tên hiện cho một bước cụ thể: Chờ ảnh / Click theo ảnh đang tìm theo chữ thì là "Chờ chữ" / "Click theo chữ".
export function stepLabel(s: Step): string {
  if (s.findBy === 'text' && (s.type === 'waitImage' || (s.type === 'clickImage' && s.imageSource === 'search'))) {
    return s.type === 'waitImage' ? 'Chờ chữ' : 'Click theo chữ';
  }
  return STEP_LABEL[s.type];
}

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
  if: 'Nếu',
  switch: 'Switch theo giá trị',
  switchScreen: 'Switch theo màn hình',
  counter: 'Đếm',
  end: 'Kết thúc',
  setVar: 'Đặt biến',
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
  { id: '1', ref: 'n1', type: 'click', points: [{ x: 412, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '2', ref: 'n2', type: 'double', points: [{ x: 640, y: 220 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '3', ref: 'n3', type: 'click', points: [{ x: 812, y: 436 }], button: 'left', clicks: 5, clickInterval: 120, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 800, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '4', ref: 'n4', type: 'drag', points: [{ x: 300, y: 600 }, { x: 900, y: 600 }], button: 'left', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 400, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 1000, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '5', ref: 'n5', type: 'wait', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 2000, jitterMs: 0, jitterPx: 0, enabled: false },
  { id: '6', ref: 'n6', type: 'click', points: [{ x: 1100, y: 520 }], button: 'right', clicks: 1, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '7', ref: 'n7', type: 'type', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: 'xin chao 123', keys: [], keyMode: 'chord', modifiers: [], delayAfter: 300, jitterMs: 0, jitterPx: 0, enabled: true },
  { id: '8', ref: 'n8', type: 'key', points: [], button: 'left', clicks: 0, clickInterval: 0, holdMs: 0, dragMs: 0, dragFromPrev: false, dragRelative: false, text: '', keys: ['enter'], keyMode: 'chord', modifiers: [], delayAfter: 500, jitterMs: 0, jitterPx: 0, enabled: true },
];

export const fmtMs = (ms: number) => (ms >= 1000 ? `${ms / 1000} s` : `${ms} ms`);
const fmtSigned = (v: number) => (v > 0 ? `+${v}` : `${v}`);
const fmtPair = (p: StepPoint) => `${fmtAxis(p.x, p.ax)}, ${fmtAxis(p.y, p.ay)}`;
// refOf: id node → mã (n2) để ghi nguồn chỗ thấy của click theo ảnh.
export const fmtPos = (s: Step, refOf: (id: string) => string = () => '?') =>
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
    : s.type === 'clickImage' ? `chỗ ${s.imageFrom ? refOf(s.imageFrom) : '(chưa chọn node)'} thấy${s.dx || s.dy ? ` ${fmtSigned(s.dx ?? 0)},${fmtSigned(s.dy ?? 0)}` : ''}`
    : s.type === 'waitPixel' && s.points[0] ? `${fmtPoint(s.points[0])} = ${s.color} ±${s.tolerance} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`
    : s.points[0] ? fmtPoint(s.points[0]) : '—';
