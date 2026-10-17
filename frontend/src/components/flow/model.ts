// Kiểu dữ liệu hiển thị của canvas flow (bước 20, F4a). Thẻ node / cạnh / nhật ký chỉ nhận các kiểu này qua props –
// màn mẫu #demo-flow tạo từ dữ liệu giả, app thật (F4b) tạo từ Script (model F1).
import type { StepIconName } from '../icons';

// Màu cổng ra: ok = thành công, fail = thất bại, neutral = nhánh khác (đúng / sai, case, còn lượt…)
export type PortKind = 'ok' | 'fail' | 'neutral';
export type PortView = {
  id: string; // success | fail | true | false | case:<key> | none | next | done
  label: string;
  kind: PortKind;
};

// Trạng thái node trong lần chạy (đang chạy / vừa xem lại nhật ký)
export type RunState = 'running' | 'ok' | 'fail' | 'skipped';

export type FlowNodeView = {
  id: string;
  ref: string; // mã ngắn n1, n2… – dùng trong biến {{n2.output.x}}
  icon: StepIconName;
  title: string;
  detail: string; // dòng mono dưới tên
  ports: PortView[]; // cổng ra theo thứ tự trên thẻ (trên → dưới)
  enabled?: boolean;
  warn?: string; // cảnh báo kiểm tra (thiếu dữ liệu, không đi tới được…)
  chips?: string[];
  breakpoint?: boolean;
  run?: RunState;
  visits?: number; // số lần đã đi qua trong lần chạy
};

export type FlowEdgeView = {
  id: string;
  from: string;
  port: string;
  to: string;
};

// Một dòng nhật ký = StepResult của một lần chạy node (model Go ở F1)
export type LogLine = { t: string; level: 'info' | 'warn' | 'error'; msg: string };
export type RunRow = {
  nodeId: string;
  at: string; // giờ bắt đầu, mono
  loop?: number; // vòng thứ mấy của kịch bản
  status: 'success' | 'fail' | 'skipped' | 'cancelled' | 'running';
  port?: string;
  next?: string; // ref node kế (hoặc "dừng" / "kết thúc")
  ms?: number;
  attempts?: number;
  output?: Record<string, unknown>;
  error?: { code: string; message: string };
  logs?: LogLine[];
  shot?: string; // ảnh màn hình lúc thất bại (file trong lịch sử chạy)
};

// Cỡ thẻ node – cố định để xếp tự động (flow/layout.ts) không cần đo DOM. NodeCard dùng đúng các số này.
export const NODE_W = 248;
export const NODE_HEAD = 56; // phần icon + tên + chi tiết
export const PORT_ROW = 24; // mỗi cổng ra một dòng
export const PORTS_PAD = 4; // đệm trên / dưới vùng cổng

export const nodeHeight = (ports: number) => 2 + NODE_HEAD + (ports ? 1 + 2 * PORTS_PAD + ports * PORT_ROW : 0);
// Toạ độ y (tính từ mép trên thẻ) của cổng ra thứ i và cổng vào
export const portY = (i: number) => 1 + NODE_HEAD + 1 + PORTS_PAD + i * PORT_ROW + PORT_ROW / 2;
export const inputY = 1 + NODE_HEAD / 2;
