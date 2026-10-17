// Kết quả nút "Thử node" (bước 20 F3, nodetest.go): cổng sẽ đi + từng điều kiện đúng / sai và vì sao. Chỉ hiển thị.
import { portView } from '../../../flow/view';
import type { Step } from '../../../types';

export type NodeTestView = {
  port: string; known: boolean; ms: number; value?: string;
  parts: { path: number[]; case?: string; ok: boolean; known: boolean; detail: string; hit?: unknown }[];
};

type Props = {
  step: Pick<Step, 'type' | 'cases'>;
  result?: NodeTestView;
  busy?: boolean;
  error?: string;
  onRun?: () => void; // không có = màn mẫu / trình duyệt (nút tắt)
  onShow?: () => void; // xem chỗ thấy trên màn hình
};

export default function NodeTestPanel({ step: s, result: r, busy, error, onRun, onShow }: Props) {
  const screen = s.type !== 'switch';
  const caseLabel = (key?: string) => (key ? portView(`case:${key}`, s).label : '');
  return (
    <div className="flex flex-col gap-1.5">
      <button onClick={onRun} disabled={!onRun || busy}
        title={screen ? 'Xét trên màn hình hiện tại (app ẩn đi một chút) – một lần, không chờ' : 'Xét giá trị lúc này (input lấy giá trị mặc định)'}
        className="h-8 rounded-control border border-line bg-surface text-ink enabled:hover:bg-hover disabled:opacity-50">
        {busy ? 'Đang xét…' : screen ? 'Thử trên màn hình' : 'Thử với giá trị lúc này'}
      </button>
      {r && (
        <div className="flex flex-col gap-1 text-xs">
          <span className={`font-semibold ${!r.known ? 'text-warn' : 'text-ink'}`}>
            {!r.known ? 'Chưa biết – có biến node, chỉ biết lúc chạy' : <>Lúc này → cổng <b className="text-accent-text">{portView(r.port, s).label}</b></>}
            <span className="font-normal text-muted"> · {r.ms} ms</span>
          </span>
          {r.value !== undefined && <span className="text-muted">Giá trị: <span className="font-mono text-ink">“{r.value}”</span></span>}
          {r.parts.map((p, i) => (
            <span key={i} className="flex gap-1.5 text-muted">
              <span className={!p.known ? 'text-warn' : p.ok ? 'text-ok' : 'text-danger-text'}>{!p.known ? '?' : p.ok ? '✓' : '✗'}</span>
              <span>
                {p.case && <b className="font-semibold">{caseLabel(p.case)}: </b>}
                {!p.case && p.path.length > 0 && <b className="font-semibold">Điều kiện {p.path.map((n) => n + 1).join('.')}: </b>}
                {p.detail}
              </span>
            </span>
          ))}
          {onShow && r.parts.some((p) => p.hit) && <button onClick={onShow} className="self-start text-muted underline">xem trên màn hình</button>}
        </div>
      )}
      {error && !busy && <span className="text-xs text-danger-text">{error}</span>}
    </div>
  );
}
