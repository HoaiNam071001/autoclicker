// Mục flow của inspector node (bước 20): thử lại khi thất bại, cổng ra (chọn node đích), output, node bắt đầu / xoá.
// Phần hiển thị dùng chung với màn mẫu (components/flow/PortField); ở đây chỉ nối vào store.
import { LastRunPanel, OutputsPanel, PortsPanel, RetryPanel } from '../flow/PortField';
import { nodeViews, outputsOf } from '../../flow/view';
import { runRows } from '../../flow/run';
import { connectPort, disconnectPort, removeStep, selectStep, setStart, updateStep, useScript } from '../../store';
import { useRunLog } from '../../runStore';
import { CAN_FAIL, Step } from '../../types';

export default function FlowSections({ step: s }: { step: Step }) {
  const script = useScript();
  const { results } = useRunLog();
  const views = nodeViews(script);
  const view = views.find((v) => v.id === s.id)!;
  const isStart = script.start === s.id;
  const i = results.findLastIndex((r) => r.nodeId === s.id);
  const lastRun = i < 0 ? undefined : runRows(results.slice(i, i + 2), script.steps, script.edges)[0];
  return (
    <>
      {CAN_FAIL.includes(s.type) && (
        <RetryPanel value={{ count: s.retry?.count ?? 0, delayMs: s.retry?.delayMs ?? 0 }}
          onChange={(r) => updateStep(s.id, { retry: r.count ? r : undefined })} />
      )}
      <PortsPanel node={view} nodes={views} edges={script.edges}
        onConnect={(port, to) => connectPort(s.id, port, to)} onDisconnect={(port) => disconnectPort(s.id, port)} onJump={selectStep} />
      <OutputsPanel refName={s.ref} fields={outputsOf(s)} />
      <LastRunPanel row={lastRun} />
      <div className="flex gap-2 border-t border-line pt-3">
        <button disabled={isStart} onClick={() => setStart(s.id)}
          className="h-8 flex-1 rounded-control border border-line bg-surface enabled:hover:bg-hover disabled:opacity-50">
          {isStart ? 'Đang là node bắt đầu' : 'Đặt làm bắt đầu'}
        </button>
        <button onClick={() => updateStep(s.id, { breakpoint: !s.breakpoint || undefined })} aria-pressed={!!s.breakpoint}
          title="Điểm dừng (debug): chạy tới node này thì tạm dừng trước khi chạy nó – Tiếp / F10 để đi tiếp"
          className={`flex h-8 items-center gap-1.5 rounded-control border px-2.5 ${s.breakpoint ? 'border-danger bg-danger-soft text-danger-text' : 'border-line bg-surface text-muted hover:bg-hover hover:text-ink'}`}>
          <span className={`h-2 w-2 rounded-full ${s.breakpoint ? 'bg-danger' : 'border border-current'}`} />Điểm dừng
        </button>
        <button onClick={() => removeStep(s.id)} className="h-8 rounded-control border border-line bg-surface px-3 text-danger-text hover:bg-danger-soft">Xoá node</button>
      </div>
    </>
  );
}
