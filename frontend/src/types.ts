import { model } from '../wailsjs/go/models';
import type { ReadFrom, ScrollMode } from './flow/web';
import type { ShotArea } from './flow/shot';
import type { ClipOp } from './flow/clipboard';
import type { BodyType, HttpMethod } from './flow/http';
import type { FileOp } from './flow/file';
import type { SysOp } from './flow/system';
import type { IfPassed } from './flow/waitUntil';

// Neo của một trục (anchor.ts, khớp internal/script/anchor.go): bỏ trống = start (cách cạnh trái / trên).
export type Anchor = 'start' | 'center' | 'end' | 'ratio';
// Toạ độ màn hình / cửa sổ.
export type Point = { x: number; y: number };
// Điểm lưu trong bước: neo theo cạnh vùng từng trục (khớp Point trong model.go). xVar / yVar (F5b): toạ độ màn hình từ biến
// ("{{n2.output.x}} + 20") – có thì bỏ qua x / y + neo (anchor.ts usesVar).
export type StepPoint = Point & { ax?: Anchor; ay?: Anchor; xVar?: string; yVar?: string };
export type Rect = { x: number; y: number; w: number; h: number };
// Một dòng bảng query / header / form của node HTTP (khớp KV trong internal/model/http.go).
export type KV = { key: string; value: string };
// Node thất bại (không thấy ảnh / chữ / màu, lỗi): thử lại thêm count lần, nghỉ delayMs; hết lượt → cổng fail. Khớp Retry trong model.go.
export type Retry = { count?: number; delayMs?: number };
// Node hay "thất bại" theo nghĩa người dùng hiểu – mới hiện phần Khi thất bại (thử lại)
export const CAN_FAIL: StepType[] = ['waitImage', 'clickImage', 'waitPixel', 'screenshot', 'clipboard', 'http', 'file', 'system', 'waitUntil', 'webWait', 'webClick', 'webType', 'webRead', 'webScroll', 'open'];

// ---- Flow (bước 20, flow.go) ----
// Cạnh: node from đi ra cổng port thì chạy node to. Mỗi cổng tối đa một cạnh. Khớp Edge trong model.go.
export interface Edge { id: string; from: string; port: string; to: string }
export type StepStatus = 'success' | 'fail' | 'skipped' | 'cancelled';
export type ErrorCode = 'not_found' | 'timeout' | 'permission' | 'invalid_input' | 'cancelled' | 'end' | 'browser' | 'http_status' | 'network' | 'internal';
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

export type ActionType = 'click' | 'double' | 'hold' | 'drag' | 'wait' | 'type' | 'key' | 'scroll' | 'open' | 'waitImage' | 'clickImage' | 'waitPixel' | 'screenshot' | 'clipboard' | 'http' | 'file' | 'system' | 'waitUntil'
  | 'webWait' | 'webClick' | 'webType' | 'webRead' | 'webScroll';
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

// Mức thao tác giống người – khớp HumanLevels (model.go).
export type HumanLevel = 'off' | 'low' | 'medium' | 'high';

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
  text: string; // type / webType: chuỗi cần gõ
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
  webBrowser?: boolean; // open: mở bằng Chrome của app (profile riêng) – node web điều khiển được; bỏ qua app
  // web qua DOM (bước 25d, internal/engine/webstep.go) – cách tìm ở flow/web.ts
  selectors?: string[]; // node web: mỗi dòng một cách tìm (text= / aria= / role= / css= / xpath=), dòng đầu ưu tiên; nhận {{biến}}
  tab?: string; // tab có URL / tiêu đề chứa chuỗi này, bỏ trống = tab đang hiện; nhận {{biến}}
  webGone?: boolean; // webWait: chờ phần tử mất
  webClear?: boolean; // webType: xoá chữ cũ trước khi gõ
  webEnter?: boolean; // webType: Enter sau khi gõ
  webBackground?: boolean; // webWait / webClick / webType / webScroll: chạy nền qua CDP, không đụng chuột / phím thật (flow/web.ts canBackground)
  readFrom?: ReadFrom; // webRead: bỏ trống = chữ hiển thị
  attr?: string; // webRead readFrom attr: tên thuộc tính
  scrollMode?: ScrollMode; // webScroll: bỏ trống = tới phần tử
  scrollIn?: string[]; // webScroll until: khung cuộn (như selectors), trống = giữa trang; nhận {{biến}}
  // chờ ảnh / click theo ảnh / chờ pixel (bước 25, internal/engine/imagestep.go)
  template?: string; // ảnh mẫu PNG base64
  templateScale?: number; // px ảnh mẫu / 1 đơn vị toạ độ lúc chụp (Retina 2)
  threshold?: number; // % khớp tối thiểu
  timeoutMs?: number; // chờ tối đa
  intervalMs?: number; // kiểm tra mỗi
  searchArea?: 'screen'; // bỏ trống = vùng của kịch bản
  imageSource?: 'node' | 'search'; // clickImage: chỗ node imageFrom thấy / tự tìm
  imageFrom?: string; // clickImage imageSource node / screenshot shotArea node: id node lấy chỗ thấy
  dx?: number; // clickImage / webClick: lệch so với tâm ảnh / phần tử
  dy?: number;
  findBy?: 'text'; // waitImage / clickImage tự tìm: bỏ trống = ảnh mẫu, text = tìm chữ (ocr.go)
  findText?: string; // chữ cần tìm, nhận {{biến}}
  textExact?: boolean; // cả dòng đúng chuỗi (mặc định: chứa chuỗi)
  textNoAccent?: boolean; // bỏ qua dấu
  occurrence?: number; // chỗ thứ mấy theo thứ tự đọc, 0 = chỗ đầu
  retry?: Retry; // thất bại thì thử lại; hết lượt → cổng fail (internal/engine/flowrun.go execNode)
  numVars?: NumVar[]; // trường số lấy từ biến (F5c) – giá trị lúc chạy thay cho trường đó
  human?: HumanLevel; // mức giống người riêng của node – trống = theo kịch bản (bước 34)
  // chụp màn hình (bước 35 N1, internal/engine/shot.go)
  shotArea?: ShotArea; // bỏ trống = vùng của kịch bản
  display?: number; // shotArea screen: 0 = màn chứa vùng kịch bản, n = màn thứ n
  shotRect?: Rect; // shotArea rect: toạ độ màn hình
  saveDir?: string; // trống = Pictures/autoclicker; nhận {{biến}}
  fileName?: string; // trống = autoclicker-<ngày>-<giờ>; nhận {{biến}}
  toClipboard?: boolean;
  // clipboard (bước 35 N2, internal/engine/clip.go)
  clipOp?: ClipOp; // bỏ trống = đọc chữ
  filePath?: string; // clipOp image: file PNG; nhận {{biến}}
  // HTTP (bước 35 N3, internal/engine/http.go) – mọi ô chữ nhận {{biến}}
  method?: HttpMethod; // bỏ trống = GET
  url?: string; // biến chèn nguyên văn – giá trị cần mã hoá đặt trong query
  query?: KV[];
  headers?: KV[];
  bodyType?: BodyType; // bỏ trống = không body
  body?: string; // json / text
  form?: KV[]; // bodyType form
  // tệp / thư mục (bước 35 N4, internal/engine/file.go) – filePath + text dùng chung
  fileOp?: FileOp; // bỏ trống = ghi chữ
  fileAppend?: boolean; // ghi chữ: thêm vào cuối
  cells?: string[]; // row: ô của dòng mới; nhận {{biến}}
  sheet?: string; // row xlsx: tên sheet, trống = sheet đầu
  filePattern?: string; // trash: mẫu tên file, vd *.png
  // hệ thống (bước 35 N5, internal/engine/system.go)
  sysOp?: SysOp;
  volume?: number; // volume: mức 0–100; volumeUp / volumeDown: mỗi lần (0 = 10)
  // chờ đến giờ (bước 35 N6, internal/engine/waituntil.go)
  at?: string; // HH:mm[:ss[.mmm]] (kèm ngày nếu cần); nhận {{biến}}
  ifPassed?: IfPassed;
  useNtp?: boolean;
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
  human?: HumanLevel; // thao tác giống người (bước 34, internal/humanize): trống / off = tắt
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
// Trường số lấy từ biến (F5c, khớp NumVar trong model.go): expr là số / biến cộng trừ, vd "{{n2.output.count}} + 1".
export type NumVarField = Exclude<PinField, 'enabled'>;
export interface NumVar { field: NumVarField; expr: string }
const numVarMatchesGo: Same<keyof NumVar, keyof Fields<model.NumVar>> = true;
void numVarMatchesGo;
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
  screenshot: 'Chụp màn hình',
  clipboard: 'Clipboard',
  http: 'HTTP',
  file: 'Tệp / thư mục',
  system: 'Hệ thống',
  waitUntil: 'Chờ đến giờ',
  webWait: 'Chờ phần tử web',
  webClick: 'Click phần tử web',
  webType: 'Gõ vào ô web',
  webRead: 'Đọc phần tử web',
  webScroll: 'Cuộn trang web',
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

export const fmtMs = (ms: number) => (ms >= 1000 ? `${ms / 1000} s` : `${ms} ms`);
// Biểu thức biến của trường số (F5c) – undefined = trường dùng số cố định.
export const numVarOf = (s: Pick<Step, 'numVars'>, field: NumVarField) => s.numVars?.find((v) => v.field === field)?.expr;
