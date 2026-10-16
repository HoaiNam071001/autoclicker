// Màn soạn kịch bản (design v2, artboard Main): thanh tiêu đề, thanh công cụ, thư viện | danh sách bước | inspector, status bar
import { useEffect, useRef, useState } from 'react';
import { Loop, Play, Stop } from '../components/icons';
import { MousePos, StopScript } from '../../wailsjs/go/main/App';
import { inWails, openOverlay } from '../windowMode';
import { removeStep, runCurrent, useEditor } from '../scriptStore';
import TitleBar from '../components/TitleBar';
import Library from '../components/Library';
import StepList from '../components/StepList';
import StepForm from '../components/StepForm';
import RunSettings, { loopSummary } from '../components/RunSettings';
import { STOP_REASON, useEngine } from '../engineStore';
import { Point } from '../types';

const MOUSE_POLL_MS = 250;

export default function MainScreen() {
  const { script, selected, notice } = useEditor();
  const { steps } = script;
  const cur = steps.find((s) => s.id === selected) ?? steps[0]; // undefined khi chưa có bước nào
  const idx = steps.indexOf(cur) + 1;
  const { running, lastError, stopReason, hotkeyError } = useEngine();

  // Toạ độ chuột thật ở status bar – để đọc nhanh toạ độ khi nhập tay.
  const [mouse, setMouse] = useState<Point | null>(null);
  useEffect(() => {
    if (!inWails) return;
    const id = setInterval(() => MousePos().then(setMouse), MOUSE_POLL_MS);
    return () => clearInterval(id);
  }, []);

  // Delete / Backspace xoá bước đang chọn (trừ khi đang gõ trong ô nhập).
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const el = e.target as HTMLElement;
      if ((e.key !== 'Delete' && e.key !== 'Backspace') || el.closest('input, textarea, select')) return;
      if (cur) removeStep(cur.id);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [cur]);

  return (
    // h-screen + từng cột tự cuộn: thanh tiêu đề và status bar luôn nằm trong cửa sổ.
    <div className="flex h-screen flex-col overflow-hidden bg-bg text-ink">
      <TitleBar />

      <div className="drag flex flex-none flex-wrap items-center gap-3.5 bg-chrome px-[18px] pt-1.5 pb-3">
        <span className="max-w-[280px] truncate text-[19px] font-semibold" title="Kịch bản">{script.name}</span>
        <div className="ml-3 flex gap-[18px]" role="tablist" aria-label="Chế độ">
          <Tab active>Soạn</Tab>
          <Tab onClick={() => openOverlay('preview')} title="Xem đường đi các bước bằng con trỏ giả trên màn hình, không click thật">Xem trước</Tab>
        </div>
        <div className="ml-auto flex gap-2">
          <LoopButton summary={loopSummary(script)} />
          <button onClick={() => (running ? StopScript() : runCurrent())} disabled={!inWails}
            title={inWails ? undefined : 'Chỉ chạy được trong app'}
            className={`flex h-9 items-center gap-2.5 rounded-control px-4 font-semibold disabled:opacity-50 ${running ? 'btn-danger-soft' : 'btn-primary'}`}>
            {running ? <><Stop className="h-2.5 w-2.5" />Dừng</> : <><Play className="h-[11px] w-[11px]" />Chạy</>}
            <span className="font-mono text-[11px] font-normal opacity-80">F8</span>
          </button>
        </div>
      </div>

      <div className="flex min-h-0 flex-1 border-t border-line">
        <Library />
        <StepList />
        {cur ? (
          <StepForm step={cur} index={idx} />
        ) : (
          <aside className="grid w-[340px] flex-none place-items-center border-l border-line bg-surface p-4 text-center text-muted">
            Chọn hoặc thêm một bước để cấu hình
          </aside>
        )}
      </div>

      <footer className="flex flex-none flex-wrap items-center gap-x-[18px] gap-y-1 border-t border-line bg-chrome px-4 py-1.5 text-xs text-muted">
        <span className="flex items-center gap-1.5">
          <span className={`h-[7px] w-[7px] rounded-full ${running ? 'bg-accent' : 'bg-ok'}`} />{running ? 'Đang chạy' : 'Sẵn sàng'}
        </span>
        <span className="font-mono">{mouse ? `X ${mouse.x} Y ${mouse.y}` : 'X — Y —'}</span>
        {hotkeyError && <span className="text-danger-text">{hotkeyError}</span>}
        {lastError && <span className="text-danger-text">Lỗi lần chạy trước: {lastError}</span>}
        {stopReason && <span>{STOP_REASON[stopReason] ?? stopReason}</span>}
        {notice && <span className="text-accent-text">{notice}</span>}
        <span className="ml-auto">Dừng khẩn cấp: F8 hoặc rê chuột vào góc màn hình</span>
      </footer>
    </div>
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
function LoopButton({ summary }: { summary: string }) {
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
          <RunSettings />
        </div>
      )}
    </div>
  );
}
