// Node (Step) → dữ liệu hiển thị thẻ trên canvas (FlowNodeView): tên, dòng chi tiết, chip, cổng ra có nhãn, cảnh báo,
// trạng thái chạy. Hàm thuần – FlowEditor (app) gọi với dữ liệu thật; test ở view.test.ts.
import type { FlowNodeView, PortView, RunState } from '../components/flow/model';
import { fmtMs, fmtPos, Script, Step, stepLabel } from '../types';
import { reachable } from './edit';
import { PORT_FAIL, PORT_SUCCESS, portsOf } from './model';

// Nhãn + màu cổng ra. Node điều khiển (F3) thêm cổng của mình ở đây.
const PORTS: Record<string, Omit<PortView, 'id'>> = {
  [PORT_SUCCESS]: { label: 'Thành công', kind: 'ok' },
  [PORT_FAIL]: { label: 'Thất bại', kind: 'fail' },
};
export const portView = (id: string): PortView => ({ id, ...(PORTS[id] ?? { label: id, kind: 'neutral' }) });

export const nodeTitle = (s: Step) => s.name || `${stepLabel(s)}${s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}`;

export type ViewOptions = {
  inputs?: (id: string) => number; // số input ghim vào node
  run?: Map<string, { run: RunState; visits?: number }>; // trạng thái lần chạy đang xem
};

export function nodeViews(sc: Pick<Script, 'steps' | 'start' | 'edges'>, opt: ViewOptions = {}): FlowNodeView[] {
  const ok = reachable(sc.edges, sc.start);
  const refs = new Map(sc.steps.map((s) => [s.id, s.ref]));
  const refOf = (id: string) => refs.get(id) ?? '?';
  return sc.steps.map((s) => {
    const n = opt.inputs?.(s.id) ?? 0;
    const r = opt.run?.get(s.id);
    return {
      id: s.id,
      ref: s.ref,
      icon: s.type,
      title: nodeTitle(s),
      detail: `${fmtPos(s, refOf)} · ${fmtMs(s.delayAfter)}`,
      ports: portsOf(s).map(portView),
      enabled: s.enabled,
      chips: [...(n ? [`${n} input`] : []), ...(s.retry?.count ? [`thử lại ${s.retry.count}`] : [])],
      warn: warningOf(s, ok.has(s.id), refs),
      run: r?.run,
      visits: r?.visits,
    };
  });
}

// Cảnh báo khi soạn (không chặn lưu; chạy thì Go kiểm tra đủ): node không đi tới được, click theo ảnh thiếu nguồn.
function warningOf(s: Step, reachable: boolean, refs: Map<string, string>): string | undefined {
  if (!reachable) return 'Không có đường nào đi tới node này – sẽ không bao giờ chạy';
  if (s.type === 'clickImage' && s.imageSource !== 'search' && !refs.has(s.imageFrom ?? '')) return 'Chưa chọn node lấy chỗ thấy (Chờ ảnh / Chờ chữ)';
  return undefined;
}

// Output node trả về khi thành công – khớp Outputs trong steptypes.go. Node sau đọc bằng {{ref.output.key}}.
export type OutputField = { key: string; desc: string };
const FOUND: OutputField[] = [
  { key: 'x', desc: 'tâm chỗ thấy – X' }, { key: 'y', desc: 'tâm chỗ thấy – Y' }, { key: 'w', desc: 'rộng' }, { key: 'h', desc: 'cao' },
  { key: 'score', desc: '% khớp' }, { key: 'text', desc: 'chữ tìm thấy (tìm theo chữ)' },
];
export function outputsOf(s: Step): OutputField[] {
  switch (s.type) {
    case 'click': case 'double': case 'hold': case 'drag':
      return [{ key: 'x', desc: 'X chuột dừng' }, { key: 'y', desc: 'Y chuột dừng' }];
    case 'type': return [{ key: 'text', desc: 'chữ đã gõ (sau khi thay biến)' }];
    case 'open': return [{ key: 'target', desc: 'URL / app đã mở' }];
    case 'waitImage': return FOUND;
    case 'clickImage': return [{ key: 'clickX', desc: 'X đã click' }, { key: 'clickY', desc: 'Y đã click' }, ...(s.imageSource === 'search' ? FOUND : [])];
    case 'waitPixel': return [{ key: 'color', desc: 'màu đọc được' }];
    default: return [];
  }
}
