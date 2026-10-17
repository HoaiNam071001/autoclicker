// Các mục của inspector khi soạn flow (bước 20):
// - Cổng ra: mỗi cổng một ô chọn node đích – cách nối không cần kéo, dùng được bằng bàn phím
// - Khi thất bại: thử lại N lần, nghỉ giữa các lần (node hành động)
// - Output: các giá trị node trả về, chèn vào ô chữ của node sau bằng {{n2.output.x}}
// - Lần chạy gần nhất: status, lỗi, log của node (StepResult)
import { useState } from 'react';
import NumInput, { fieldBox } from '../NumInput';
import { FlowEdgeView, FlowNodeView, PortKind, RunRow } from './model';

const dot: Record<PortKind, string> = { ok: 'bg-ok', fail: 'bg-danger', neutral: 'bg-muted' };
const section = 'flex flex-col gap-2.5 border-t border-line pt-3';
const head = 'flex items-center justify-between font-semibold';

export function PortsPanel({ node, nodes, edges, onConnect, onDisconnect, onJump }: {
  node: FlowNodeView; nodes: FlowNodeView[]; edges: FlowEdgeView[];
  onConnect: (port: string, to: string) => void; onDisconnect: (port: string) => void; onJump: (id: string) => void;
}) {
  if (!node.ports.length) {
    return <section className={section}><span className={head}>Cổng ra</span><span className="text-[13px] text-muted">Node Kết thúc không có cổng ra – kịch bản dừng ở đây.</span></section>;
  }
  const others = nodes.filter((n) => n.id !== node.id);
  return (
    <section className={section}>
      <span className={head}>Cổng ra <span className="text-xs font-normal text-muted">hoặc kéo từ chấm trên thẻ</span></span>
      {node.ports.map((p) => {
        const to = edges.find((e) => e.from === node.id && e.port === p.id)?.to ?? '';
        return (
          <div key={p.id} className="flex flex-col gap-1.5">
            <span className="flex items-center gap-2 text-muted"><span className={`h-2 w-2 rounded-full ${dot[p.kind]}`} />{p.label}</span>
            <div className="flex gap-1.5">
              <select value={to} onChange={(e) => (e.target.value ? onConnect(p.id, e.target.value) : onDisconnect(p.id))}
                className={`${fieldBox} min-w-0 flex-1 text-[13px] text-ink outline-none`}>
                <option value="">{p.kind === 'fail' ? '— Không nối: dừng kịch bản, báo lỗi' : '— Không nối: kết thúc kịch bản'}</option>
                {others.map((n) => <option key={n.id} value={n.id}>{n.ref} · {n.title}</option>)}
              </select>
              <button disabled={!to} onClick={() => onJump(to)} title="Chọn node đích" aria-label="Chọn node đích"
                className="grid h-8 w-8 flex-none place-items-center rounded-control border border-line text-muted enabled:hover:bg-hover enabled:hover:text-ink disabled:opacity-40">→</button>
            </div>
          </div>
        );
      })}
    </section>
  );
}

export type Retry = { count: number; delayMs: number };
export function RetryPanel({ value, onChange }: { value: Retry; onChange: (v: Retry) => void }) {
  return (
    <section className={section}>
      <span className={head}>Khi thất bại</span>
      <div className="grid grid-cols-2 gap-3">
        <label className="flex flex-col gap-1.5 text-muted">Thử lại<NumInput value={value.count} unit="lần" onChange={(count) => onChange({ ...value, count })} /></label>
        <label className="flex flex-col gap-1.5 text-muted">Nghỉ giữa các lần<NumInput value={value.delayMs} unit="ms" onChange={(delayMs) => onChange({ ...value, delayMs })} /></label>
      </div>
      <span className="text-xs text-muted">Thử hết vẫn không được thì đi ra cổng <span className="text-danger-text">Thất bại</span>.</span>
    </section>
  );
}

import type { OutputField } from '../../flow/view';
export type { OutputField };
export function OutputsPanel({ refName, fields }: { refName: string; fields: OutputField[] }) {
  const [copied, setCopied] = useState('');
  if (!fields.length) return null;
  const copy = (k: string) => {
    const t = `{{${refName}.output.${k}}}`;
    navigator.clipboard?.writeText(t).catch(() => {});
    setCopied(k);
    setTimeout(() => setCopied(''), 1200);
  };
  return (
    <section className={section}>
      <span className={head}>Output <span className="text-xs font-normal text-muted">bấm để chép biến</span></span>
      <div className="flex flex-col gap-1">
        {fields.map((f) => (
          <button key={f.key} onClick={() => copy(f.key)} className="flex items-center gap-2 rounded-md px-1 py-0.5 text-left hover:bg-hover">
            <span className="var-token">{`${refName}.output.${f.key}`}</span>
            <span className="min-w-0 flex-1 truncate text-xs text-muted">{copied === f.key ? 'đã chép' : f.desc}</span>
          </button>
        ))}
      </div>
      <span className="text-xs text-muted">Còn có <span className="font-mono">{`{{${refName}.status}}`}</span>, <span className="font-mono">{`{{${refName}.error.message}}`}</span>.</span>
    </section>
  );
}

export function LastRunPanel({ row }: { row: RunRow | undefined }) {
  if (!row) return null;
  return (
    <section className={section}>
      <span className={head}>Lần chạy gần nhất <StatusChip status={row.status} /></span>
      <div className="flex flex-wrap gap-x-4 gap-y-1 font-mono text-xs text-muted">
        <span>{row.at}</span>
        {row.ms !== undefined && <span>{fmtMs(row.ms)}</span>}
        {row.attempts && row.attempts > 1 && <span>{row.attempts} lần thử</span>}
        {row.port && <span>→ {row.port}{row.next ? ` · ${row.next}` : ''}</span>}
      </div>
      {row.error && (
        <div className="rounded-control bg-danger-soft px-2.5 py-2 text-[13px] text-danger-text">
          <span className="font-mono text-xs">{row.error.code}</span> · {row.error.message}
        </div>
      )}
      {row.output && <pre className="overflow-x-auto rounded-control bg-soft px-2.5 py-2 font-mono text-xs text-ink">{JSON.stringify(row.output, null, 1)}</pre>}
      {!!row.logs?.length && <LogLines logs={row.logs} />}
    </section>
  );
}

export function LogLines({ logs }: { logs: NonNullable<RunRow['logs']> }) {
  return (
    <ol className="flex flex-col gap-0.5 font-mono text-xs">
      {logs.map((l, i) => (
        <li key={i} className="flex gap-2">
          <span className="text-muted">{l.t}</span>
          <span className={l.level === 'error' ? 'text-danger-text' : l.level === 'warn' ? 'text-warn' : 'text-ink'}>{l.msg}</span>
        </li>
      ))}
    </ol>
  );
}

const STATUS: Record<RunRow['status'], [string, string]> = {
  success: ['Thành công', 'bg-soft text-ok'],
  fail: ['Thất bại', 'bg-danger-soft text-danger-text'],
  skipped: ['Bỏ qua', 'bg-soft text-muted'],
  running: ['Đang chạy', 'bg-accent-soft text-accent-text'],
};
export const StatusChip = ({ status }: { status: RunRow['status'] }) => (
  <span className={`rounded-full px-2 py-px text-[11px] font-semibold whitespace-nowrap ${STATUS[status][1]}`}>{STATUS[status][0]}</span>
);

export const fmtMs = (ms: number) => (ms >= 1000 ? `${(ms / 1000).toFixed(1).replace('.', ',')} s` : `${ms} ms`);
