// Bảng Nhật ký dưới canvas (bước 20 F5): mỗi dòng một lần chạy node (StepResult). Bấm dòng → chọn node trên canvas
// và mở chi tiết (output / lỗi / log). Thu gọn được còn một dòng tiêu đề.
import { useState } from 'react';
import { StepIcon } from '../icons';
import { FlowNodeView, RunRow } from './model';
import { fmtMs, LogLines, StatusChip } from './PortField';

export default function RunLog({ title, summary, ok, rows, nodes, selected, onSelect }: {
  title: string; summary: string; ok: boolean | null; // null = đang chạy
  rows: RunRow[]; nodes: FlowNodeView[]; selected: string | null; onSelect: (id: string) => void;
}) {
  const [open, setOpen] = useState(true);
  const [row, setRow] = useState<number | null>(null);
  const byId = new Map(nodes.map((n) => [n.id, n]));
  return (
    <section aria-label="Nhật ký chạy" className="flex flex-none flex-col border-t border-line bg-surface" style={{ height: open ? 220 : 'auto' }}>
      <button onClick={() => setOpen(!open)} aria-expanded={open} className="flex h-9 flex-none items-center gap-3 px-4 text-left hover:bg-hover">
        <svg viewBox="0 0 24 24" className={`h-3.5 w-3.5 text-muted transition-transform ${open ? '' : '-rotate-90'}`} fill="none" stroke="currentColor" strokeWidth={2} aria-hidden><path d="M6 9l6 6 6-6" /></svg>
        <span className="font-semibold">Nhật ký</span>
        <span className="text-muted">{title}</span>
        <span className={`ml-auto text-xs ${ok === false ? 'text-danger-text' : ok ? 'text-ok' : 'text-accent-text'}`}>{summary}</span>
      </button>
      {open && (
        <div className="flex min-h-0 flex-1">
          <ol className="min-w-0 flex-1 overflow-y-auto border-t border-line">
            {rows.map((r, i) => {
              const n = byId.get(r.nodeId);
              const on = row === i;
              return (
                <li key={i}>
                  <button onClick={() => { setRow(i); onSelect(r.nodeId); }}
                    className={`grid w-full grid-cols-[64px_34px_minmax(0,1fr)_88px_64px_minmax(0,160px)] items-center gap-2 px-4 py-1 text-left text-[13px] ${on ? 'bg-accent-soft' : selected === r.nodeId ? 'bg-hover' : 'hover:bg-hover'}`}>
                    <span className="font-mono text-xs text-muted">{r.at}</span>
                    <span className="font-mono text-xs text-muted">{n?.ref}</span>
                    <span className="flex min-w-0 items-center gap-2">
                      {n && <StepIcon type={n.icon} className="h-3.5 w-3.5 flex-none text-muted" />}
                      <span className="truncate">{n?.title}</span>
                      {r.attempts && r.attempts > 1 && <span className="font-mono text-[11px] text-muted">{r.attempts} lần thử</span>}
                    </span>
                    <StatusChip status={r.status} />
                    <span className="font-mono text-xs text-muted">{r.ms !== undefined ? fmtMs(r.ms) : ''}</span>
                    <span className="truncate font-mono text-xs text-muted">{r.port ? `→ ${r.port}${r.next ? ` · ${r.next}` : ''}` : ''}</span>
                  </button>
                </li>
              );
            })}
          </ol>
          {row !== null && rows[row] && <Detail r={rows[row]} />}
        </div>
      )}
    </section>
  );
}

function Detail({ r }: { r: RunRow }) {
  return (
    <div className="flex w-[340px] flex-none flex-col gap-2 overflow-y-auto border-t border-l border-line p-3 max-[1100px]:w-[280px]">
      {r.error && <div className="rounded-control bg-danger-soft px-2.5 py-2 text-[13px] text-danger-text"><span className="font-mono text-xs">{r.error.code}</span> · {r.error.message}</div>}
      {r.output && <pre className="overflow-x-auto rounded-control bg-soft px-2.5 py-2 font-mono text-xs">{JSON.stringify(r.output, null, 1)}</pre>}
      {r.logs?.length ? <LogLines logs={r.logs} /> : <span className="text-xs text-muted">Không có log.</span>}
    </div>
  );
}
