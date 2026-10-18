// Hộp thoại Mẫu thử (bước 33): kịch bản mẫu đơn giản để test từng loại bước, lọc theo loại bước / kiểu app + ô tìm.
// Chỉ hiển thị – nạp danh sách / mở mẫu ở Samples.tsx (app thật) hoặc dữ liệu giả (#demo).
import { ReactNode, useEffect, useMemo, useState } from 'react';
import { StepIcon, StepIconName, STEP_ICON } from '../icons';
import { fieldBox } from '../NumInput';
import { ALL, appsOf, filterSamples, Sample, SampleFilter, stepCounts, stepName } from '../../samples/filter';
import SampleCard from './SampleCard';

export type SamplesDialogProps = {
  samples: Sample[] | null; // null = đang nạp
  error?: string;
  ready: boolean; // false: trình duyệt không có Go – không mở được
  onPick: (s: Sample) => void;
  onOpenPage: (page: string) => void; // '' = danh sách trang sân thử
  onClose: () => void;
};

export default function SamplesDialog({ samples, error, ready, onPick, onOpenPage, onClose }: SamplesDialogProps) {
  const [f, setF] = useState<SampleFilter>(ALL);
  const list = samples ?? [];
  const shown = useMemo(() => filterSamples(list, f), [list, f]);
  const steps = useMemo(() => stepCounts(list), [list]);
  const apps = useMemo(() => appsOf(list), [list]);
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [onClose]);
  const toggle = (k: 'step' | 'app', v: string) => setF({ ...f, [k]: f[k] === v ? '' : v });
  return (
    <div className="fixed inset-0 z-40 flex items-start justify-center bg-black/25 px-4 pt-[6vh]" onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div role="dialog" aria-label="Mẫu thử" className="flex max-h-[86vh] w-[820px] max-w-full flex-col overflow-hidden rounded-panel border border-line bg-surface text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
        <div className="flex flex-col gap-3 border-b border-line p-[18px]">
          <div className="flex items-start gap-3">
            <div className="flex-1">
              <div className="text-base font-semibold">Mẫu thử</div>
              <div className="text-[13px] text-muted">Kịch bản đơn giản để test từng loại bước, chạy trên sân thử cục bộ (trang giả lập offline, tự hiện kết quả). Dùng mẫu = mở bản sao chưa lưu.</div>
            </div>
            <button disabled={!ready} onClick={() => onOpenPage('')} className="h-8 shrink-0 rounded-control border border-line px-2.5 text-[13px] enabled:hover:bg-hover disabled:opacity-60">Mở sân thử</button>
            <button onClick={onClose} aria-label="Đóng" className="h-8 w-8 shrink-0 rounded-control text-muted hover:bg-hover hover:text-ink">✕</button>
          </div>
          <label className={fieldBox}>
            <input autoFocus value={f.q} onChange={(e) => setF({ ...f, q: e.target.value })} placeholder="Tìm mẫu (không dấu cũng được)…" className="w-full bg-transparent text-ink outline-none" />
          </label>
          <ChipRow label="Kiểu">
            {apps.map((a) => <Chip key={a} on={f.app === a} onClick={() => toggle('app', a)}>{a}</Chip>)}
          </ChipRow>
          <ChipRow label="Bước">
            {steps.map(({ type, count }) => (
              <Chip key={type} on={f.step === type} onClick={() => toggle('step', type)}>
                {type in STEP_ICON && <StepIcon type={type as StepIconName} className="h-3.5 w-3.5" />}
                {stepName(type)} <span className="font-mono text-[11px] opacity-70">{count}</span>
              </Chip>
            ))}
          </ChipRow>
        </div>
        <div className="grid flex-1 grid-cols-1 content-start gap-3 overflow-y-auto bg-bg p-[18px] sm:grid-cols-2">
          {error ? <Empty>Không nạp được danh sách mẫu: {error}</Empty>
            : !samples ? <Empty>Đang nạp…</Empty>
            : !shown.length ? <Empty>Không có mẫu khớp bộ lọc.</Empty>
            : shown.map((s) => <SampleCard key={s.id} sample={s} ready={ready} onPick={() => onPick(s)} onOpenPage={() => onOpenPage(s.page)} />)}
        </div>
        <div className="border-t border-line px-[18px] py-2 text-xs text-muted">{shown.length} / {list.length} mẫu</div>
      </div>
    </div>
  );
}

function ChipRow({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex items-start gap-2">
      <span className="w-10 shrink-0 pt-1 text-xs text-muted">{label}</span>
      <div className="flex flex-wrap gap-1">{children}</div>
    </div>
  );
}

function Chip({ on, onClick, children }: { on: boolean; onClick: () => void; children: ReactNode }) {
  return (
    <button onClick={onClick} aria-pressed={on}
      className={`flex h-7 items-center gap-1 rounded-full border px-2.5 text-xs ${on ? 'border-accent bg-accent-soft text-accent-text' : 'border-line text-muted hover:bg-hover hover:text-ink'}`}>
      {children}
    </button>
  );
}

const Empty = ({ children }: { children: ReactNode }) => <div className="col-span-full py-10 text-center text-muted">{children}</div>;
