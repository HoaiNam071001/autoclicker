// Bảng thêm bước (design v2, artboard AddStep) – mở bằng "+ Thêm bước" hoặc Ctrl/Cmd+K.
// Ô tìm (không dấu cũng được), 4 nhóm, phím mũi tên + Enter; bước mới chèn sau bước đang chọn.
// Loại chưa làm hiện nhãn "Sắp có" và không chọn được (logic ở các bước sau).
import { KeyboardEvent, useEffect, useMemo, useRef, useState } from 'react';
import { StepIcon, StepIconName } from './icons';
import { addPointStep, addStep, getEditor, insertTarget, pathOf, wrapStep } from '../scriptStore';
import { addImageStep, addTextStep } from '../imageSteps';

type Item = { icon: StepIconName; label: string; keywords: string; run?: () => unknown; soon?: string };

const GROUPS: { name: string; items: Item[] }[] = [
  {
    name: 'Chuột', items: [
      { icon: 'click', label: 'Click', keywords: 'bam nhan chuot trai phai', run: () => addPointStep('click') },
      { icon: 'double', label: 'Double click', keywords: 'nhay dup hai lan', run: () => addPointStep('double') },
      { icon: 'hold', label: 'Giữ chuột', keywords: 'hold nhan giu', run: () => addPointStep('hold') },
      { icon: 'drag', label: 'Kéo thả', keywords: 'drag keo tha a b', run: () => addPointStep('drag') },
      { icon: 'scroll', label: 'Cuộn', keywords: 'scroll lan chuot', run: () => addStep('scroll') },
      { icon: 'clickImage', label: 'Click theo ảnh', keywords: 'image anh mau tim nut', run: () => addImageStep('clickImage') },
      { icon: 'clickImage', label: 'Click theo chữ', keywords: 'text chu ocr doc nut tim', run: () => addTextStep('clickImage') },
    ],
  },
  {
    name: 'Bàn phím', items: [
      { icon: 'type', label: 'Gõ chữ', keywords: 'type text nhap van ban', run: () => addStep('type') },
      { icon: 'key', label: 'Phím tắt', keywords: 'key hotkey enter ctrl to hop', run: () => addStep('key') },
    ],
  },
  {
    name: 'Chờ & điều kiện', items: [
      { icon: 'wait', label: 'Chờ', keywords: 'wait delay nghi doi', run: () => addStep('wait') },
      { icon: 'waitImage', label: 'Chờ ảnh', keywords: 'wait image anh mau xuat hien khoanh', run: () => addImageStep('waitImage') },
      { icon: 'waitImage', label: 'Chờ chữ', keywords: 'wait text chu ocr doc xuat hien', run: () => addTextStep('waitImage') },
      { icon: 'waitPixel', label: 'Chờ màu pixel', keywords: 'pixel mau color diem', run: () => addImageStep('waitPixel') },
    ],
  },
  {
    name: 'Ứng dụng & cấu trúc', items: [
      { icon: 'open', label: 'Mở URL / app', keywords: 'url web link trinh duyet ung dung open youtube', run: () => addStep('open') },
      { icon: 'focusWindow', label: 'Focus cửa sổ', keywords: 'window cua so ung dung', soon: 'bước 26' },
    ],
  },
  {
    name: 'Điều kiện & lặp', items: [
      { icon: 'if', label: 'Nếu… / Nếu không', keywords: 'if else dieu kien re nhanh thay khong thay', run: () => addStep('if') },
      { icon: 'group', label: 'Nhóm / Lặp', keywords: 'group loop lap lai repeat den khi', run: () => addStep('group') },
      { icon: 'break', label: 'Thoát khối lặp', keywords: 'break thoat dung lap', run: () => addStep('break') },
      { icon: 'stop', label: 'Dừng kịch bản', keywords: 'stop dung ket thuc bao loi', run: () => addStep('stop') },
      { icon: 'group', label: 'Bọc bước đang chọn trong Nhóm', keywords: 'wrap boc gom nhom lap', run: () => wrapSelected('group') },
      { icon: 'if', label: 'Bọc bước đang chọn trong Nếu', keywords: 'wrap boc dieu kien if', run: () => wrapSelected('if') },
    ],
  },
];

const wrapSelected = (type: 'group' | 'if') => {
  const id = getEditor().selected;
  if (id) wrapStep(id, type);
};
const COLS = 3;

// So khớp không dấu: "cuon" tìm được "Cuộn".
const fold = (t: string) => t.normalize('NFD').replace(/[̀-ͯ]/g, '').replace(/đ/g, 'd').replace(/Đ/g, 'D').toLowerCase();

// Bước mới vào đâu: cuối khối (nút "+ Thêm bước" trong khối) hoặc sau bước đang chọn.
function where() {
  const t = insertTarget();
  if (t) return `thêm vào cuối ${t.branch === 'else' ? 'nhánh “Nếu không” của ' : ''}khối ${pathOf(t.parent)}`;
  const sel = getEditor().selected;
  return sel ? `thêm sau bước ${pathOf(sel)}` : 'thêm vào cuối';
}

export default function AddStepPalette({ onClose }: { onClose: () => void }) {
  const [q, setQ] = useState('');
  const groups = useMemo(() => {
    const words = fold(q).split(/\s+/).filter(Boolean);
    return GROUPS
      .map((g) => ({ ...g, items: g.items.filter((it) => words.every((w) => fold(`${it.label} ${it.keywords} ${g.name}`).includes(w))) }))
      .filter((g) => g.items.length > 0);
  }, [q]);
  const flat = groups.flatMap((g) => g.items);
  const [sel, setSel] = useState(0);
  useEffect(() => setSel(Math.max(flat.findIndex((it) => it.run), 0)), [q]); // eslint-disable-line react-hooks/exhaustive-deps
  const panel = useRef<HTMLDivElement>(null);

  const pick = (it: Item | undefined) => {
    if (!it?.run) return;
    onClose();
    it.run();
  };

  // ←→ đi từng mục, ↑↓ đi theo cột (3 mục một hàng, tính trong danh sách đang hiện).
  const onKey = (e: KeyboardEvent) => {
    const move = { ArrowLeft: -1, ArrowRight: 1, ArrowUp: -COLS, ArrowDown: COLS }[e.key];
    if (move !== undefined) {
      e.preventDefault();
      setSel((i) => Math.min(Math.max(i + move, 0), flat.length - 1));
    } else if (e.key === 'Enter') {
      e.preventDefault();
      pick(flat[sel]);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
    }
  };

  let n = -1; // chỉ số trong flat khi vẽ
  return (
    <div className="fixed inset-0 z-30 flex items-start justify-center bg-black/20 px-4 pt-[10vh]"
      onMouseDown={(e) => { if (!panel.current?.contains(e.target as Node)) onClose(); }}>
      <div ref={panel} role="dialog" aria-label="Thêm bước" onKeyDown={onKey}
        className="flex max-h-[80vh] w-[680px] max-w-full flex-col overflow-hidden rounded-panel border border-line bg-surface text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
        <label className="mx-3.5 mt-3.5 mb-1.5 flex h-10 flex-none items-center gap-3 rounded-lg px-3 ring-focus">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-muted" aria-hidden>
            <circle cx="11" cy="11" r="7" /><path d="M20 20l-3.5-3.5" />
          </svg>
          <input autoFocus value={q} onChange={(e) => setQ(e.target.value)} placeholder="Tìm loại bước… (vd: chờ, cuộn, url)"
            className="min-w-0 flex-1 bg-transparent text-[15px] outline-none" />
          <span className="font-mono text-[11px] text-muted">Esc</span>
        </label>

        <div className="flex min-h-0 flex-1 flex-col gap-3 overflow-y-auto px-3.5 py-2">
          {groups.length === 0 && <p className="py-6 text-center text-muted">Không có loại bước nào khớp “{q}”</p>}
          {groups.map((g) => (
            <section key={g.name} className="flex flex-col gap-1.5">
              <span className="px-1 text-xs font-semibold text-muted">{g.name}</span>
              <div className="grid grid-cols-3 gap-1.5">
                {g.items.map((it) => {
                  const i = ++n;
                  const on = i === sel;
                  return (
                    <button key={it.label} onClick={() => pick(it)} onMouseEnter={() => setSel(i)} aria-disabled={!it.run}
                      title={it.soon ? `Sắp có (${it.soon})` : undefined}
                      className={`flex min-h-[46px] items-center gap-2.5 rounded-card px-2.5 py-1.5 text-left ${on ? 'ring-selected' : 'border border-line bg-surface'} ${it.run ? '' : 'cursor-default'}`}>
                      <span className={`grid h-[30px] w-[30px] flex-none place-items-center rounded-[7px] ${on && it.run ? 'bg-accent-soft text-accent-text' : 'bg-soft text-muted'} ${it.run ? '' : 'opacity-60'}`}>
                        <StepIcon type={it.icon} className="h-[15px] w-[15px]" />
                      </span>
                      <span className={`min-w-0 flex-1 font-medium ${it.run ? '' : 'text-muted'}`}>{it.label}</span>
                      {it.soon && <span className="rounded-full bg-soft px-1.5 py-px text-[10px] font-semibold whitespace-nowrap text-muted">Sắp có</span>}
                    </button>
                  );
                })}
              </div>
            </section>
          ))}
        </div>

        <div className="flex flex-none flex-wrap items-center gap-4 border-t border-line bg-chrome px-4 py-2.5 text-xs text-muted">
          <span><span className="font-mono">↑↓←→</span> chọn</span>
          <span><span className="font-mono">Enter</span> {where()}</span>
          <span className="ml-auto">Click / Double / Giữ / Kéo thả mở overlay chọn điểm</span>
        </div>
      </div>
    </div>
  );
}
