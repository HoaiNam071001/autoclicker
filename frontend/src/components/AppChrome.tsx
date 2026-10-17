// Khung chung của màn soạn (design v2): thanh công cụ dưới thanh tiêu đề + status bar.
// Dùng cho MainScreen (dữ liệu thật) và DemoScreen (#demo, dữ liệu mẫu).
import { ReactNode, useEffect, useRef, useState } from 'react';
import { Loop, Play, Stop } from './icons';

export function Toolbar({ name, nameTitle, unsaved, unsavedTitle, loop, loopPanel, running, runDisabled, runTitle, onRun, onPreview, extra }: {
  name: string; nameTitle?: string;
  unsaved?: boolean; unsavedTitle?: string;
  loop: string; loopPanel: ReactNode; // chữ trên nút "N vòng" + nội dung popup
  running: boolean; runDisabled?: boolean; runTitle?: string; onRun: () => void;
  onPreview: () => void;
  extra?: ReactNode; // nút thêm cạnh nút Chạy (debug, ghi macro…)
}) {
  return (
    <div className="drag flex flex-none flex-wrap items-center gap-3.5 bg-chrome px-[18px] pt-1.5 pb-3">
      <span className="max-w-[280px] truncate text-[19px] font-semibold" title={nameTitle ?? 'Kịch bản'}>{name}</span>
      {unsaved && (
        <span title={unsavedTitle} className="flex h-[22px] items-center gap-1.5 rounded-full bg-soft px-[9px] text-xs text-muted">
          <span className="h-1.5 w-1.5 rounded-full bg-warn" />Chưa lưu
        </span>
      )}
      <div className="ml-3 flex gap-[18px]" role="tablist" aria-label="Chế độ">
        <Tab active>Soạn</Tab>
        <Tab onClick={onPreview} title="Xem đường đi các bước bằng con trỏ giả trên màn hình, không click thật">Xem trước</Tab>
      </div>
      <div className="ml-auto flex gap-2">
        {extra}
        <LoopButton summary={loop}>{loopPanel}</LoopButton>
        <button onClick={onRun} disabled={runDisabled} title={runTitle}
          className={`flex h-9 items-center gap-2.5 rounded-control px-4 font-semibold disabled:opacity-50 ${running ? 'btn-danger-soft' : 'btn-primary'}`}>
          {running ? <><Stop className="h-2.5 w-2.5" />Dừng</> : <><Play className="h-[11px] w-[11px]" />Chạy</>}
          <span className="font-mono text-[11px] font-normal opacity-80">F8</span>
        </button>
      </div>
    </div>
  );
}

// Nút debug (bước 25c): chạy từ bước đang chọn, chạy từng bước (Tiếp = F10). soon: giao diện có, logic chưa.
export function DebugButtons({ soon, onFrom, onStep }: { soon?: boolean; onFrom?: () => void; onStep?: () => void }) {
  const cls = 'flex h-9 items-center gap-2 rounded-control border border-line bg-surface px-3 text-[13px] enabled:hover:bg-hover disabled:opacity-50';
  const title = (t: string) => (soon ? `${t} – sắp có (bước 25c)` : t);
  return (
    <div className="flex gap-1">
      <button disabled={soon} onClick={onFrom} title={title('Chạy từ bước đang chọn tới hết')} className={cls}>
        <Play className="h-2.5 w-2.5" />Từ bước này
      </button>
      <button disabled={soon} onClick={onStep} title={title('Chạy từng bước: dừng sau mỗi bước, bấm Tiếp / F10 để chạy bước kế')} className={cls}>
        <StepOver />Từng bước <span className="font-mono text-[11px] text-muted">F10</span>
      </button>
    </div>
  );
}

const StepOver = () => (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.9" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
    <path d="M4 14a8 8 0 0 1 15-4M19 4v6h-6" /><circle cx="12" cy="19" r="1.6" />
  </svg>
);

// Nút ghi macro (bước 24): ghi click / kéo / gõ thật thành bước, F9 dừng. soon: giao diện có, logic chưa.
export function RecordButton({ soon, recording, onClick }: { soon?: boolean; recording?: string | null; onClick?: () => void }) {
  return (
    <button onClick={onClick} disabled={soon} title={soon ? 'Ghi macro – sắp có (bước 24)' : recording ? 'Dừng ghi (F9)' : 'Ghi thao tác chuột / phím thành bước (F9 dừng)'}
      className={`flex h-9 items-center gap-2 rounded-control border px-3 text-[13px] disabled:opacity-50 ${recording ? 'btn-danger-soft' : 'border-line bg-surface enabled:hover:bg-hover'}`}>
      <span className={`h-2.5 w-2.5 rounded-full bg-danger ${recording ? 'animate-pulse' : ''}`} />
      {recording ? <>Đang ghi <span className="font-mono">{recording}</span></> : 'Ghi'}
      <span className="font-mono text-[11px] text-muted">F9</span>
    </button>
  );
}

// Status bar: chấm trạng thái, toạ độ chuột, các thông báo (children), gợi ý dừng khẩn cấp bên phải.
export function StatusBar({ running, mouse, children }: { running: boolean; mouse: string; children?: ReactNode }) {
  return (
    <footer className="flex flex-none flex-wrap items-center gap-x-[18px] gap-y-1 border-t border-line bg-chrome px-4 py-1.5 text-xs text-muted">
      <span className="flex items-center gap-1.5">
        <span className={`h-[7px] w-[7px] rounded-full ${running ? 'bg-accent' : 'bg-ok'}`} />{running ? 'Đang chạy' : 'Sẵn sàng'}
      </span>
      <span className="font-mono">{mouse}</span>
      {children}
      <span className="ml-auto">Dừng khẩn cấp: F8 hoặc rê chuột vào góc màn hình</span>
    </footer>
  );
}

function Tab({ active, children, ...rest }: { active?: boolean; children: string; onClick?: () => void; title?: string }) {
  return (
    <button role="tab" aria-selected={!!active} {...rest}
      className={`relative h-[34px] px-0.5 ${active ? 'font-semibold text-ink' : 'text-muted hover:text-ink'}`}>
      {children}
      {active && <span className="absolute inset-x-1/4 bottom-px h-[3px] rounded-sm bg-accent" />}
    </button>
  );
}

// Nút "10 vòng": mở popup cài đặt chạy (lặp, đếm ngược).
function LoopButton({ summary, children }: { summary: string; children: ReactNode }) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!open) return;
    const close = (e: MouseEvent | KeyboardEvent) => {
      if (e instanceof KeyboardEvent ? e.key === 'Escape' : !ref.current?.contains(e.target as Node)) setOpen(false);
    };
    window.addEventListener('mousedown', close);
    window.addEventListener('keydown', close);
    return () => {
      window.removeEventListener('mousedown', close);
      window.removeEventListener('keydown', close);
    };
  }, [open]);
  return (
    <div ref={ref} className="relative">
      <button onClick={() => setOpen(!open)} aria-expanded={open} title="Cài đặt lặp / đếm ngược"
        className="flex h-9 items-center gap-2 rounded-control border border-line bg-surface px-3 hover:bg-hover">
        <Loop /><span className="font-mono text-[13px]">{summary}</span>
      </button>
      {open && (
        <div className="absolute right-0 top-full z-20 mt-1.5 w-[340px] rounded-panel border border-line bg-surface p-4 shadow-[0_8px_28px_rgb(0_0_0/.16)]">
          {children}
        </div>
      )}
    </div>
  );
}
