// Thẻ node trên canvas flow (bước 20). Phần trên giống thẻ bước (ô icon, mã n2, tên, chi tiết mono);
// phần dưới là các cổng ra, mỗi cổng một dòng có nhãn + chấm màu (thành công ok, thất bại danger, nhánh khác trung tính).
// Cỡ cố định theo model.ts để xếp tự động không cần đo DOM.
import { Handle, Node, NodeProps, Position } from '@xyflow/react';
import { StepIcon } from '../icons';
import { tileClass } from '../StepCard';
import { FlowNodeView, inputY, NODE_HEAD, NODE_W, PORT_ROW, PORTS_PAD } from './model';

export type CardData = { view: FlowNodeView; start: boolean; linked: string[] };
export type CardNode = Node<CardData, 'card'>;

export default function NodeCard({ data, selected }: NodeProps<CardNode>) {
  const { view: v, start, linked } = data;
  const off = v.enabled === false;
  const outline = v.run === 'running' ? 'outline-2 outline-offset-2 outline-accent' : v.run === 'fail' ? 'outline-2 outline-offset-2 outline-danger' : '';
  return (
    <div style={{ width: NODE_W }} title={v.warn}
      className={`relative rounded-card ${selected ? 'ring-selected' : 'border border-line bg-surface shadow-[0_1px_3px_rgb(0_0_0/.05)]'} ${outline} ${v.run === 'skipped' ? 'opacity-60' : ''}`}>
      {start && (
        <span className="absolute -top-[23px] left-0 flex h-[19px] items-center gap-1 rounded-full bg-accent-soft px-2 text-[11px] font-semibold text-accent-text">
          <svg viewBox="0 0 24 24" className="h-2 w-2" fill="currentColor" aria-hidden><path d="M6 4l14 8-14 8z" /></svg>Bắt đầu
        </span>
      )}
      <Badge v={v} />
      <Handle type="target" position={Position.Left} style={{ top: inputY }} className="flow-handle flow-handle-in" />

      <div style={{ height: NODE_HEAD }} className="flex items-center gap-2.5 px-3">
        <span className={`${tileClass(!!selected)} ${off ? 'opacity-50' : ''}`}><StepIcon type={v.icon} /></span>
        <span className={`min-w-0 flex-1 ${off ? 'opacity-50' : ''}`}>
          <span className="flex items-center gap-1.5">
            {v.breakpoint && <span className="h-2 w-2 flex-none rounded-full bg-danger" title="Điểm dừng (debug)" />}
            <span className="font-mono text-[11px] text-muted">{v.ref}</span>
            <span className="truncate font-semibold">{v.title}</span>
          </span>
          <span className="block truncate font-mono text-xs text-muted">{v.detail}</span>
        </span>
        {off && <span className="flex-none rounded-full bg-soft px-[7px] py-px text-[11px] text-muted">Tắt</span>}
      </div>

      {v.ports.length > 0 && (
        <div style={{ paddingBlock: PORTS_PAD }} className="relative border-t border-line">
          {/* chip phụ nằm bên trái vùng cổng (bên phải là nhãn cổng) */}
          {!!v.chips?.length && (
            <span style={{ height: PORT_ROW, top: PORTS_PAD }} className="absolute left-3 flex items-center gap-1">
              {v.chips.map((c) => <span key={c} className="rounded-[5px] border border-line px-1.5 text-[11px] leading-[17px] text-muted">{c}</span>)}
            </span>
          )}
          {v.ports.map((p) => {
            const free = !linked.includes(p.id);
            return (
              <div key={p.id} style={{ height: PORT_ROW }} className="relative flex items-center justify-end gap-1.5 pr-4 text-xs">
                <span className={p.kind === 'fail' ? 'text-danger-text' : 'text-muted'}>{p.label}</span>
                {free && (
                  <span className="font-mono text-[11px] text-muted opacity-70" title={p.kind === 'fail' ? 'Chưa nối: lỗi ở đây thì dừng kịch bản' : 'Chưa nối: đi ra cổng này thì kết thúc kịch bản'}>
                    → {p.kind === 'fail' ? 'dừng' : 'hết'}
                  </span>
                )}
                <Handle type="source" id={p.id} position={Position.Right} className={`flow-handle flow-port-${p.kind}`} />
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

// Huy hiệu góc phải trên: trạng thái chạy (đang chạy / xong / lỗi + số lần đi qua), hoặc cảnh báo khi đang soạn.
function Badge({ v }: { v: FlowNodeView }) {
  const visits = v.visits && v.visits > 1 ? <span className="font-mono text-[11px]">×{v.visits}</span> : null;
  const box = 'absolute -top-2.5 -right-2.5 flex h-5 min-w-5 items-center justify-center gap-0.5 rounded-full border bg-surface px-1';
  if (v.run === 'running') {
    return <span className={`${box} border-accent text-accent-text`} title="Đang chạy"><span className="h-2.5 w-2.5 animate-spin rounded-full border-2 border-accent border-t-transparent" />{visits}</span>;
  }
  if (v.run === 'ok') return <span className={`${box} border-ok text-ok`} title="Thành công"><Check />{visits}</span>;
  if (v.run === 'fail') return <span className={`${box} border-danger text-danger-text`} title="Thất bại"><Cross />{visits}</span>;
  if (v.run === 'skipped') return <span className={`${box} border-line text-muted`} title="Bỏ qua (node tắt)"><span className="font-mono text-[11px]">–</span></span>;
  if (v.warn) return <span className={`${box} border-warn font-mono text-[11px] font-bold text-warn`} title={v.warn}>!</span>;
  return null;
}

const Check = () => <svg viewBox="0 0 24 24" className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={3} strokeLinecap="round" strokeLinejoin="round" aria-hidden><path d="M5 12.5l4.5 4.5L19 7.5" /></svg>;
const Cross = () => <svg viewBox="0 0 24 24" className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth={3} strokeLinecap="round" aria-hidden><path d="M6 6l12 12M18 6L6 18" /></svg>;
