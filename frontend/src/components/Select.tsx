// Ô chọn tự vẽ thay <select> của hệ điều hành: danh sách gốc của macOS / Windows không theo theme – ở theme tối chữ
// bị trùng màu nền. Menu đặt fixed theo ô (inspector cuộn sẽ cắt nếu absolute); cuộn trang thì đóng.
import { CSSProperties, KeyboardEvent, ReactNode, useEffect, useRef, useState } from 'react';
import { usePopover } from './usePopover';

export type SelectOption<T extends string = string> = {
  value: T;
  label: string;
  icon?: ReactNode;
  sub?: string; // chữ phụ bên phải (số lượng, phím tắt…)
  disabled?: boolean;
};

type Props<T extends string> = {
  value: T;
  options: SelectOption<T>[];
  onChange: (v: T) => void;
  className?: string; // khung ô (fieldBox / selectBox…)
  placeholder?: string; // value không có trong options
  trigger?: ReactNode; // nội dung nút thay cho nhãn + ▾ (vd chữ loại node dưới tên)
  menuWidth?: number;
  disabled?: boolean;
  title?: string;
  'aria-label'?: string;
};

const ITEM_H = 32;
const Chevron = () => (
  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4" className="flex-none text-muted" aria-hidden><path d="M6 9l6 6 6-6" /></svg>
);

export default function Select<T extends string>({ value, options, onChange, className = '', placeholder = '', trigger, menuWidth, disabled, title, ...rest }: Props<T>) {
  const { open, setOpen, ref } = usePopover();
  const [pos, setPos] = useState<CSSProperties>({});
  const [active, setActive] = useState(-1);
  const list = useRef<HTMLDivElement>(null);
  const cur = options.find((o) => o.value === value);
  useEffect(() => {
    if (!open) return;
    const close = (e: Event) => !list.current?.contains(e.target as Node) && setOpen(false);
    window.addEventListener('scroll', close, true);
    window.addEventListener('resize', close);
    list.current?.querySelector('[aria-selected="true"]')?.scrollIntoView({ block: 'nearest' });
    return () => {
      window.removeEventListener('scroll', close, true);
      window.removeEventListener('resize', close);
    };
  }, [open, setOpen]);
  const show = () => {
    const r = ref.current!.getBoundingClientRect();
    const h = Math.min(options.length * ITEM_H + 10, 300);
    const up = r.bottom + h + 8 > window.innerHeight && r.top - h - 8 > 0;
    const w = Math.max(r.width, menuWidth ?? 180);
    setPos({ left: Math.max(Math.min(r.left, window.innerWidth - w - 8), 8), width: w, maxHeight: h, ...(up ? { bottom: window.innerHeight - r.top + 4 } : { top: r.bottom + 4 }) });
    setActive(options.findIndex((o) => o.value === value));
    setOpen(true);
  };
  const pick = (o?: SelectOption<T>) => {
    if (!o || o.disabled) return;
    setOpen(false);
    if (o.value !== value) onChange(o.value);
  };
  const step = (d: number) => {
    for (let i = 1; i <= options.length; i++) {
      const k = (active + d * i + options.length * 2) % options.length;
      if (!options[k].disabled) return setActive(k);
    }
  };
  const onKey = (e: KeyboardEvent) => {
    if (!open) {
      if (['ArrowDown', 'ArrowUp', 'Enter', ' '].includes(e.key)) { e.preventDefault(); show(); }
      return;
    }
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') { e.preventDefault(); step(e.key === 'ArrowDown' ? 1 : -1); }
    else if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); pick(options[active]); }
    else if (e.key === 'Tab') setOpen(false);
  };
  return (
    <div ref={ref} className={trigger ? 'relative inline-flex min-w-0' : 'min-w-0'}>
      <button type="button" disabled={disabled} title={title} aria-label={rest['aria-label']} aria-haspopup="listbox" aria-expanded={open}
        onClick={() => (open ? setOpen(false) : show())} onKeyDown={onKey}
        className={trigger ? 'inline-flex min-w-0 items-center text-left' : `${className} w-full min-w-0 gap-1.5 text-left outline-none disabled:opacity-60`}>
        {trigger ?? (
          <>
            <span className={`min-w-0 flex-1 truncate ${cur ? '' : 'text-muted'}`}>{cur?.label ?? placeholder}</span>
            <Chevron />
          </>
        )}
      </button>
      {open && (
        <div ref={list} role="listbox" style={pos}
          className="fixed z-50 flex flex-col gap-px overflow-y-auto rounded-[9px] border border-line bg-surface p-[5px] text-[13px] text-ink shadow-[0_10px_28px_rgb(0_0_0/.18)]">
          {options.map((o, i) => (
            <button key={o.value} type="button" role="option" aria-selected={o.value === value} disabled={o.disabled}
              onMouseEnter={() => setActive(i)} onClick={() => pick(o)}
              className={`flex min-h-[30px] flex-none items-center gap-2 rounded-md px-2 text-left disabled:opacity-50 ${i === active ? 'bg-hover' : ''} ${o.value === value ? 'font-semibold text-accent-text' : ''}`}>
              {o.icon}
              <span className="min-w-0 flex-1 truncate">{o.label}</span>
              {o.sub && <span className="flex-none font-mono text-[11px] text-muted">{o.sub}</span>}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
