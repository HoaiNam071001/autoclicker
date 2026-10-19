// Khung chung của màn soạn (design v3): thanh công cụ dưới thanh tiêu đề + status bar.
// Dùng cho MainScreen (dữ liệu thật) và màn mẫu (#demo, #demo-flow).
import { ReactNode } from 'react';
import { Loop } from './icons';
import { RunButton, RunMenuItem } from './RunMenu';
import { usePopover } from './usePopover';

export function Toolbar({ name, nameTitle, unsaved, unsavedTitle, loop, loopPanel, running, runDisabled, runTitle, onRun, runItems, onPreview, extra }: {
  name: string; nameTitle?: string;
  unsaved?: boolean; unsavedTitle?: string;
  loop: string; loopPanel: ReactNode; // chữ trên nút "N vòng" + nội dung popup
  running: boolean; runDisabled?: boolean; runTitle?: string; onRun: (e: { shiftKey: boolean }) => void;
  runItems: RunMenuItem[]; // menu ▾ cạnh nút Chạy
  onPreview: () => void;
  extra?: ReactNode; // trạng thái thêm trước nút vòng lặp (vd "Đang ghi 0:07")
}) {
  return (
    <div className="drag flex flex-none flex-wrap items-center gap-3.5 bg-chrome px-4 pt-1 pb-2.5">
      <span className="max-w-[280px] truncate text-lg font-semibold" title={nameTitle ?? 'Kịch bản'}>{name}</span>
      {unsaved && (
        <span title={unsavedTitle} className="flex h-[22px] items-center gap-1.5 rounded-full bg-soft px-[9px] text-xs text-muted">
          <span className="h-1.5 w-1.5 rounded-full bg-warn" />Chưa lưu
        </span>
      )}
      <div className="ml-2 flex gap-4" role="tablist" aria-label="Chế độ">
        <Tab active>Soạn</Tab>
        <Tab onClick={onPreview} title="Xem đường đi các bước bằng con trỏ giả trên màn hình, không click thật">Xem trước</Tab>
      </div>
      <div className="ml-auto flex items-center gap-2">
        {extra}
        <LoopButton summary={loop}>{loopPanel}</LoopButton>
        <RunButton running={running} disabled={runDisabled} title={runTitle} onRun={onRun} items={runItems} />
      </div>
    </div>
  );
}

// Status bar: chấm trạng thái, toạ độ chuột, các thông báo (children), gợi ý dừng khẩn cấp bên phải.
export function StatusBar({ running, mouse, children }: { running: boolean; mouse: string; children?: ReactNode }) {
  return (
    <footer className="flex flex-none flex-wrap items-center gap-x-4 gap-y-1 border-t border-line bg-chrome px-4 py-[5px] text-xs text-muted">
      <span className="flex items-center gap-1.5">
        <span className={`h-[7px] w-[7px] rounded-full ${running ? 'bg-accent' : 'bg-ok'}`} />{running ? 'Đang chạy' : 'Sẵn sàng'}
      </span>
      <span className="font-mono">{mouse}</span>
      {children}
      <span className="ml-auto" title="Bấm F8 hoặc rê chuột vào một góc màn hình để dừng kịch bản đang chạy">Dừng khẩn cấp: F8 / góc màn hình</span>
    </footer>
  );
}

function Tab({ active, children, ...rest }: { active?: boolean; children: string; onClick?: () => void; title?: string }) {
  return (
    <button role="tab" aria-selected={!!active} {...rest}
      className={`relative h-8 px-0.5 ${active ? 'font-semibold text-ink' : 'text-muted hover:text-ink'}`}>
      {children}
      {active && <span className="absolute inset-x-[22%] bottom-0 h-[3px] rounded-sm bg-accent" />}
    </button>
  );
}

// Nút "10 vòng": mở popup cài đặt chạy (lặp, đếm ngược).
function LoopButton({ summary, children }: { summary: string; children: ReactNode }) {
  const { open, setOpen, ref } = usePopover();
  return (
    <div ref={ref} className="relative">
      <button onClick={() => setOpen(!open)} aria-expanded={open} title="Cài đặt lặp / đếm ngược"
        className="flex h-[34px] items-center gap-[7px] rounded-control border border-line bg-surface px-3 hover:bg-hover">
        <Loop className="h-3.5 w-3.5" /><span className="font-mono text-[12.5px]">{summary}</span>
      </button>
      {open && (
        <div className="absolute right-0 top-full z-20 mt-1.5 w-[340px] rounded-panel border border-line bg-surface p-4 shadow-[0_8px_28px_rgb(0_0_0/.16)]">
          {children}
        </div>
      )}
    </div>
  );
}
