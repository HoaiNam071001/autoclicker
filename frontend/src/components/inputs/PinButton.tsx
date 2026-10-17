// Nút ghim cạnh trường số / chữ (artboard StepInputs): ghim = cho phép đổi giá trị khi bấm Chạy (input của kịch bản).
import { CSSProperties, useEffect, useRef, useState } from 'react';
import { Pin } from '../icons';

export default function PinButton({ pinned, onToggle, disabled, title }: {
  pinned: boolean; onToggle?: () => void; disabled?: boolean; title?: string;
}) {
  return (
    <button type="button" aria-pressed={pinned} onClick={onToggle} disabled={disabled}
      aria-label={pinned ? 'Bỏ ghim input' : 'Đặt làm input'}
      title={title ?? (pinned ? 'Bỏ ghim – không hỏi khi chạy nữa' : 'Đặt làm input – hỏi giá trị khi bấm Chạy')}
      className={`grid h-8 w-8 flex-none place-items-center rounded-control border disabled:opacity-50 ${pinned
        ? 'border-accent/45 bg-accent-soft text-accent-text'
        : 'border-line text-muted enabled:hover:bg-hover enabled:hover:text-ink'}`}>
      <Pin />
    </button>
  );
}

// Chip key cạnh nhãn của trường đã ghim
export const KeyChip = ({ k }: { k: string }) => (
  <span className="var-token inline-flex items-center gap-1 py-0! text-[11px]"><Pin className="h-2.5 w-2.5" />{k}</span>
);

// Nút ghim của trường trong inspector: chưa ghim → tạo input mới (giá trị hiện tại làm mặc định), hoặc nếu đã có input
// cùng kiểu thì mở menu chọn "dùng chung" (một input cho nhiều bước); đã ghim → bỏ ghim.
export function PinMenu({ pinnedKey, options, onPinNew, onPinTo, onUnpin }: {
  pinnedKey?: string;
  options: { key: string; label: string }[];
  onPinNew: () => void; onPinTo: (key: string) => void; onUnpin: () => void;
}) {
  const [open, setOpen] = useState(false);
  // fixed theo vị trí nút: inspector cuộn (overflow) sẽ cắt menu nếu đặt absolute; gần đáy cửa sổ thì mở lên trên
  const [pos, setPos] = useState<CSSProperties>({});
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
  if (pinnedKey) return <PinButton pinned onToggle={onUnpin} title={`Bỏ ghim – không hỏi ${pinnedKey} khi chạy nữa`} />;
  const choose = (f: () => void) => {
    setOpen(false);
    f();
  };
  const item = 'flex h-9 w-full items-center gap-2.5 rounded-[7px] px-2.5 text-left text-[13px] hover:bg-hover';
  return (
    <div ref={ref} className="relative flex-none">
      <PinButton pinned={false} onToggle={() => {
        if (!options.length) return onPinNew();
        const r = ref.current!.getBoundingClientRect();
        const up = r.bottom + 120 + options.length * 36 > window.innerHeight;
        setPos({ right: Math.max(window.innerWidth - r.right, 8), ...(up ? { bottom: window.innerHeight - r.top + 4 } : { top: r.bottom + 4 }) });
        setOpen(!open);
      }} />
      {open && (
        <div role="menu" style={pos} className={`fixed z-30 flex w-[250px] flex-col gap-0.5 rounded-[10px] border border-line bg-surface p-1.5 shadow-[0_10px_30px_rgb(0_0_0/.14)]`}>
          <button role="menuitem" onClick={() => choose(onPinNew)} className="flex w-full flex-col items-start rounded-[7px] px-2.5 py-1.5 text-left text-[13px] text-ink hover:bg-hover">
            + Input mới<span className="text-xs text-muted">giá trị hiện tại làm mặc định</span>
          </button>
          <div className="mx-1 my-0.5 h-px bg-line" />
          <span className="px-2.5 py-1 text-xs text-muted">Dùng chung input có sẵn</span>
          {options.map((o) => (
            <button key={o.key} role="menuitem" onClick={() => choose(() => onPinTo(o.key))} className={item}>
              <span className="font-mono text-xs text-accent-text">{o.key}</span>
              <span className="min-w-0 flex-1 truncate text-muted">{o.label}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
