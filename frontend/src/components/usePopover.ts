// Popup / menu gắn vào một nút: bấm ra ngoài khung hoặc Esc thì đóng. Dùng cho nút vòng lặp, menu Chạy, menu theme.
import { useEffect, useRef, useState } from 'react';

// escOnly: Esc chỉ đóng popup, không tới lắng nghe khác (overlay – Esc thoát overlay; nghe ở pha capture).
export function usePopover({ escOnly = false } = {}) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!open) return;
    const close = (e: MouseEvent | KeyboardEvent) => {
      if (e instanceof KeyboardEvent ? e.key !== 'Escape' : ref.current?.contains(e.target as Node)) return;
      if (escOnly) e.stopPropagation();
      setOpen(false);
    };
    window.addEventListener('mousedown', close);
    window.addEventListener('keydown', close, escOnly);
    return () => {
      window.removeEventListener('mousedown', close);
      window.removeEventListener('keydown', close, escOnly);
    };
  }, [open, escOnly]);
  return { open, setOpen, ref };
}
