// Ô chữ có biến {{key}} (artboard StepInputs): tô biến bằng var-token, gõ "{{" mở gợi ý input (↑↓ Enter / Tab chèn,
// Esc đóng), cuối danh sách "+ Tạo input mới…". Cách làm: textarea chữ trong suốt đặt trên một lớp vẽ cùng nội dung
// (cùng font mono, cùng padding) – biến chỉ đổi màu nền / chữ, không đổi độ rộng nên con trỏ luôn khớp.
import { KeyboardEvent, useLayoutEffect, useMemo, useRef, useState } from 'react';
import { isValidKey, ScriptInput, splitVars, varsIn } from './model';

type Props = {
  value: string;
  onChange: (v: string) => void;
  inputs: Pick<ScriptInput, 'key' | 'label'>[];
  onCreateInput?: (key: string) => void; // biến mới chưa có trong tab Input
  rows?: number;
  placeholder?: string;
};

// Đoạn "{{abc" ngay trước con trỏ (chưa đóng) → đang gõ biến, trả phần đã gõ.
const openVar = (text: string, caret: number) => {
  const m = /\{\{([A-Za-z0-9_]*)$/.exec(text.slice(0, caret));
  return m && text[m.index - 1] !== '\\' ? { start: m.index, typed: m[1] } : null;
};

export default function VarTextInput({ value, onChange, inputs, onCreateInput, rows = 3, placeholder }: Props) {
  const ta = useRef<HTMLTextAreaElement>(null);
  const mirror = useRef<HTMLDivElement>(null);
  const [caret, setCaret] = useState(0);
  const [focused, setFocused] = useState(false);
  const [sel, setSel] = useState(0);
  const [dismissed, setDismissed] = useState(-1); // Esc: ẩn gợi ý cho tới khi con trỏ sang chỗ khác
  // Vị trí con trỏ cần đặt sau khi chèn biến – đặt ngay khi React ghi giá trị mới vào textarea
  // (đặt trễ bằng requestAnimationFrame thì chữ gõ tiếp kịp chen vào chỗ cũ).
  const pendingCaret = useRef<number | null>(null);
  useLayoutEffect(() => {
    const pos = pendingCaret.current;
    if (pos === null || !ta.current) return;
    pendingCaret.current = null;
    ta.current.setSelectionRange(pos, pos);
    setCaret(pos);
  }, [value]);

  const open = focused ? openVar(value, caret) : null;
  const matches = useMemo(
    () => (open ? inputs.filter((i) => i.key.toLowerCase().startsWith(open.typed.toLowerCase())) : []),
    [open?.typed, inputs], // eslint-disable-line react-hooks/exhaustive-deps
  );
  const canCreate = !!onCreateInput && !!open && isValidKey(open.typed) && !inputs.some((i) => i.key === open.typed);
  const options = matches.length + (canCreate || (onCreateInput && open) ? 1 : 0);
  const showPopup = !!open && dismissed !== open.start && options > 0;

  const used = varsIn(value);
  const unknown = used.filter((k) => !inputs.some((i) => i.key === k));

  const insert = (key: string) => {
    if (!open) return;
    const end = value.slice(caret).startsWith('}}') ? caret + 2 : caret;
    const next = `${value.slice(0, open.start)}{{${key}}}${value.slice(end)}`;
    onChange(next);
    pendingCaret.current = open.start + key.length + 4;
  };
  const choose = (i: number) => {
    if (i < matches.length) return insert(matches[i].key);
    if (!open) return;
    const key = isValidKey(open.typed) ? open.typed : `bien_${inputs.length + 1}`;
    onCreateInput?.(key);
    insert(key);
  };

  const onKey = (e: KeyboardEvent) => {
    if (!showPopup) return;
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
      e.preventDefault();
      setSel((s) => (s + (e.key === 'ArrowDown' ? 1 : options - 1)) % options);
    } else if (e.key === 'Enter' || e.key === 'Tab') {
      e.preventDefault();
      choose(sel);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      e.stopPropagation();
      setDismissed(open!.start);
    }
  };

  const syncCaret = () => {
    const c = ta.current?.selectionStart ?? 0;
    setCaret(c);
    setSel(0);
    if (mirror.current && ta.current) mirror.current.scrollTop = ta.current.scrollTop;
  };

  const box = 'col-start-1 row-start-1 whitespace-pre-wrap break-words px-3 py-2.5 font-mono text-[13px] leading-[1.9]';
  return (
    <div className="flex flex-col gap-1.5">
      <div className={`relative grid rounded-lg border bg-field ${focused ? 'ring-focus' : 'border-line'}`}>
        {/* lớp vẽ: cùng chữ, biến tô màu; dòng cuối thêm khoảng trắng để cao bằng textarea */}
        <div ref={mirror} aria-hidden className={`${box} overflow-hidden text-ink`}>
          {splitVars(value).map((s, i) => s.varKey
            ? <span key={i} className={`rounded-[4px] ${unknown.includes(s.varKey) ? 'bg-danger-soft text-danger-text' : 'bg-accent-soft text-accent-text'} shadow-[0_0_0_2px_var(--c-accent-soft)]`}>{s.text}</span>
            : <span key={i}>{s.text}</span>)}
          {'​'}
        </div>
        <textarea ref={ta} value={value} rows={rows} placeholder={placeholder} spellCheck={false}
          onChange={(e) => { onChange(e.target.value); setCaret(e.target.selectionStart); setSel(0); }}
          onSelect={syncCaret} onScroll={syncCaret} onKeyDown={onKey}
          onFocus={() => setFocused(true)} onBlur={() => setFocused(false)}
          className={`${box} resize-none bg-transparent text-transparent caret-[var(--c-ink)] outline-none placeholder:text-muted`} />

        {showPopup && (
          <div role="listbox" className="absolute top-full left-0 z-20 mt-1 flex w-[260px] flex-col gap-0.5 rounded-[10px] border border-line bg-surface p-1.5 shadow-[0_10px_30px_rgb(0_0_0/.14)]"
            onMouseDown={(e) => e.preventDefault() /* giữ focus ở textarea */}>
            {matches.map((m, i) => (
              <button key={m.key} role="option" aria-selected={sel === i} onClick={() => choose(i)} onMouseEnter={() => setSel(i)}
                className={`flex h-[38px] items-center gap-2.5 rounded-[7px] px-2.5 text-left ${sel === i ? 'bg-hover' : ''}`}>
                <span className="min-w-14 font-mono text-xs text-accent-text">{m.key}</span>
                <span className="flex-1 truncate text-[13px] text-muted">{m.label}</span>
              </button>
            ))}
            {onCreateInput && (
              <>
                {matches.length > 0 && <div className="mx-1 my-0.5 h-px bg-line" />}
                <button role="option" aria-selected={sel === matches.length} onClick={() => choose(matches.length)} onMouseEnter={() => setSel(matches.length)}
                  className={`h-9 rounded-[7px] px-2.5 text-left text-[13px] ${sel === matches.length ? 'bg-hover' : ''}`}>
                  + Tạo input mới{canCreate ? <> <span className="font-mono text-accent-text">{open!.typed}</span></> : '…'}
                </button>
              </>
            )}
          </div>
        )}
      </div>
      <span className="flex items-center gap-1.5 text-xs text-muted">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden><circle cx="12" cy="12" r="9" /><path d="M12 8h.01M11 12h1v4h1" /></svg>
        {used.length === 0 ? 'Gõ {{ để chèn biến'
          : unknown.length ? <span className="text-danger-text">Biến chưa có trong tab Input: {unknown.join(', ')}</span>
          : `${used.length} biến, tự thêm vào tab Input`}
      </span>
    </div>
  );
}
