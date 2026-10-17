// Mục flow của inspector node (bước 20): thử lại khi thất bại, cổng ra (chọn node đích), output, node bắt đầu / xoá.
// Phần hiển thị dùng chung với màn mẫu (components/flow/PortField); ở đây chỉ nối vào store.
import { OutputsPanel, PortsPanel, RetryPanel } from '../flow/PortField';
import { nodeViews, outputsOf } from '../../flow/view';
import { connectPort, disconnectPort, removeStep, selectStep, setStart, updateStep, useScript } from '../../store';
import { CAN_FAIL, Step } from '../../types';

export default function FlowSections({ step: s }: { step: Step }) {
  const script = useScript();
  const views = nodeViews(script);
  const view = views.find((v) => v.id === s.id)!;
  const isStart = script.start === s.id;
  return (
    <>
      {CAN_FAIL.includes(s.type) && (
        <RetryPanel value={{ count: s.retry?.count ?? 0, delayMs: s.retry?.delayMs ?? 0 }}
          onChange={(r) => updateStep(s.id, { retry: r.count ? r : undefined })} />
      )}
      <PortsPanel node={view} nodes={views} edges={script.edges}
        onConnect={(port, to) => connectPort(s.id, port, to)} onDisconnect={(port) => disconnectPort(s.id, port)} onJump={selectStep} />
      <OutputsPanel refName={s.ref} fields={outputsOf(s)} />
      <div className="flex gap-2 border-t border-line pt-3">
        <button disabled={isStart} onClick={() => setStart(s.id)}
          className="h-8 flex-1 rounded-control border border-line bg-surface enabled:hover:bg-hover disabled:opacity-50">
          {isStart ? 'Đang là node bắt đầu' : 'Đặt làm bắt đầu'}
        </button>
        <button onClick={() => removeStep(s.id)} className="h-8 rounded-control border border-line bg-surface px-3 text-danger-text hover:bg-danger-soft">Xoá node</button>
      </div>
    </>
  );
}
