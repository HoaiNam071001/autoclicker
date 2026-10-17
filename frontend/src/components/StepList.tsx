// Danh sách bước (design v2: thẻ bước): chọn, kéo ô icon để sắp xếp, bật/tắt, nhân bản, xoá, thêm bước (bước 14).
// Dạng cây (bước 20b): bước trong khối Nếu / Nhóm thụt vào, có đường kẻ nét đứt bên trái; khối gập / mở được;
// kéo bước vào / ra khối – danh sách phẳng có thụt lề (stepTree.ts rows / dropLoc), thả xuống thì vị trí theo dòng phía trên.
import { CSSProperties, useEffect, useState } from 'react';
import {
  closestCenter, DndContext, DragEndEvent, DragStartEvent, KeyboardSensor, PointerSensor, useSensor, useSensors,
} from '@dnd-kit/core';
import {
  arrayMove, SortableContext, sortableKeyboardCoordinates, useSortable, verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Copy, Power, Region, StepIcon, Trash } from './icons';
import {
  duplicateStep, moveStepTo, removeStep, selectStep, setInsertInto, setMulti, setRegion, toggleMultiStep, toggleStep,
  undoDelete, useEditor,
} from '../scriptStore';
import { findStep, fmtMs, fmtPos, OnFail, Script, Step, STEP_LABEL } from '../types';
import { Branch, dropLoc, flatOf, isBlock, Row, rows as buildRows } from '../stepTree';
import { stepInputCount } from '../scriptInputs';
import { templateUrl } from '../imageSteps';
import { useEngine } from '../engineStore';
import { openOverlay } from '../windowMode';

import AddStepPalette from './AddStepPalette';
import { CardBody, cardClass, tileClass } from './StepCard';

const isMac = navigator.platform.startsWith('Mac');
const INDENT = 22; // px mỗi tầng khối

export default function StepList() {
  const { script, selected, deleted, multi } = useEditor();
  const { steps, region } = script;
  const multiSteps = new Set(multi.map((k) => k.split(':')[0]));
  const [palette, setPalette] = useState(false);
  const [collapsed, setCollapsed] = useState<Set<string>>(new Set());
  const [dragging, setDragging] = useState<string | null>(null);
  // Bảng thêm bước: từ nút trong khối thì bước mới vào cuối khối đó, còn lại chèn sau bước đang chọn.
  const openPalette = (into: { parent: string; branch: Branch } | null) => {
    setInsertInto(into);
    setPalette(true);
  };
  // Ctrl+K (macOS: Cmd+K) mở / đóng bảng thêm bước
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key.toLowerCase() !== 'k' || !(isMac ? e.metaKey : e.ctrlKey)) return;
      e.preventDefault();
      setInsertInto(null);
      setPalette((o) => !o);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);
  // distance: click vào ô icon không kéo thì không tính là kéo
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 4 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  // đang kéo một khối: ẩn bước con của nó (kéo cả khối, không thả được vào chính nó)
  const list = buildRows(steps, collapsed, dragging);
  const total = list.filter((r) => r.kind === 'step').length;

  const onDragEnd = ({ active, over }: DragEndEvent) => {
    setDragging(null);
    if (!over || active.id === over.id) return;
    const ids = list.map((r) => r.id);
    const from = ids.indexOf(String(active.id));
    const to = ids.indexOf(String(over.id));
    if (from < 0 || to < 0) return;
    const moved = arrayMove(ids, from, to);
    const at = moved.indexOf(String(active.id));
    const above = at > 0 ? list.find((r) => r.id === moved[at - 1]) ?? null : null;
    moveStepTo(String(active.id), dropLoc(steps, String(active.id), above));
  };
  const toggleFold = (id: string) => setCollapsed((c) => {
    const n = new Set(c);
    if (!n.delete(id)) n.add(id);
    return n;
  });

  return (
    <main aria-label="Danh sách bước" className="flex min-w-0 flex-1 flex-col gap-3 overflow-y-auto px-[22px] py-4">
      <div className="flex flex-wrap items-center gap-2.5">
        <h1 className="text-base font-semibold">Các bước</h1>
        <span className="text-muted">{total} bước</span>
        <span className="ml-auto flex items-center gap-1">
          <button onClick={() => openOverlay('region')} title="Chọn vùng trên màn hình – toạ độ điểm tính theo vùng"
            className="flex h-7 items-center gap-1.5 rounded-md border border-line bg-surface px-2.5 text-[13px] text-muted hover:bg-hover hover:text-ink">
            <Region className="h-3.5 w-3.5" />
            {region ? <span className="font-mono">Vùng {region.w} × {region.h}</span> : 'Toàn màn hình'}
          </button>
          {region && (
            <button onClick={() => setRegion(null)} aria-label="Về toàn màn hình" title="Về toàn màn hình"
              className="grid h-7 w-7 place-items-center rounded-md text-muted hover:bg-hover hover:text-ink">✕</button>
          )}
        </span>
      </div>

      {steps.length === 0 && <p className="py-6 text-center text-muted">Chưa có bước nào – bấm “+ Thêm bước” hoặc F6 để lưu vị trí chuột.</p>}
      <DndContext sensors={sensors} collisionDetection={closestCenter}
        onDragStart={({ active }: DragStartEvent) => setDragging(String(active.id))} onDragCancel={() => setDragging(null)} onDragEnd={onDragEnd}>
        <SortableContext items={list.map((r) => r.id)} strategy={verticalListSortingStrategy}>
          <ol className="flex flex-col gap-1.5">
            {list.map((r) => r.kind === 'step'
              ? <Card key={r.id} row={r} selected={r.id === selected} inMulti={multiSteps.has(r.id)} script={script}
                  folded={collapsed.has(r.id)} onFold={() => toggleFold(r.id)} />
              : <ExtraRow key={r.id} row={r} steps={steps} onAdd={(b) => openPalette({ parent: r.parent, branch: b })} />)}
          </ol>
        </SortableContext>
      </DndContext>

      {deleted && (
        <div className="flex items-center justify-between rounded-card bg-soft px-3 py-2 text-[13px]">
          <span>Đã xoá bước {deleted.path} ({STEP_LABEL[deleted.step.type]}{deleted.step.children?.length ? ` và ${countAll(deleted.step) - 1} bước bên trong` : ''})</span>
          <button onClick={undoDelete} className="font-semibold text-accent-text">Hoàn tác</button>
        </div>
      )}

      <button onClick={() => openPalette(null)} aria-haspopup="dialog"
        className="flex h-[42px] flex-none items-center justify-center gap-2.5 rounded-card border border-dashed border-line-strong text-muted hover:bg-hover hover:text-ink">
        + Thêm bước <span className="font-mono text-[11px]">{isMac ? '⌘K' : 'Ctrl+K'}</span>
      </button>
      {palette && <AddStepPalette onClose={() => setPalette(false)} />}
    </main>
  );
}

const countAll = (s: Step): number => 1 + [...(s.children ?? []), ...(s.else ?? [])].reduce((n, c) => n + countAll(c), 0);

// Đường kẻ nét đứt bên trái cho mỗi tầng khối chứa dòng – nhìn như khung khối bao các bước con.
function Rails({ depth }: { depth: number }) {
  return (
    <>
      {Array.from({ length: depth }, (_, k) => (
        <span key={k} aria-hidden className="pointer-events-none absolute -top-[3px] -bottom-[3px] border-l border-dashed border-line-strong"
          style={{ left: k * INDENT + 9 }} />
      ))}
    </>
  );
}

// Dòng phụ: tiêu đề "Nếu không" / "Khi thất bại", cuối khối "+ Thêm bước vào khối". Không kéo được, nhưng thả lên được.
function ExtraRow({ row, steps, onAdd }: { row: Exclude<Row, { kind: 'step' }>; steps: Step[]; onAdd: (b: Branch) => void }) {
  const { setNodeRef, transform, transition } = useSortable({ id: row.id, disabled: { draggable: true, droppable: false } });
  const style: CSSProperties = { transform: CSS.Translate.toString(transform), transition, paddingLeft: row.depth * INDENT };
  const parent = flatOf(steps, row.parent);
  if (row.kind === 'else') {
    return (
      <li ref={setNodeRef} style={style} className="relative flex items-center gap-2 py-0.5 text-[13px] font-semibold text-muted">
        <Rails depth={row.depth} />
        <span className={row.depth ? 'pl-2' : ''}>{row.label}</span>
        <span className="font-normal">{row.label === 'Nếu không'
          ? `– điều kiện bước ${parent?.path} sai thì chạy các bước dưới`
          : `– bước ${parent?.path} thất bại (sau khi thử lại) thì chạy các bước dưới rồi đi tiếp`}</span>
      </li>
    );
  }
  return (
    <li ref={setNodeRef} style={style} className="relative">
      <Rails depth={row.depth} />
      <button onClick={() => onAdd(row.branch)}
        className="ml-2 flex h-7 items-center gap-1.5 rounded-md px-2 text-[12px] text-muted hover:bg-hover hover:text-ink">
        + Thêm bước vào {row.branch === 'else' ? `nhánh “${parent?.step.type === 'if' ? 'Nếu không' : 'Khi thất bại'}”` : `khối ${parent?.path}`}
      </button>
    </li>
  );
}

// Chip "khi thất bại" trên thẻ bước: ↻2 → bỏ qua / → bước 5 / dừng
function failChip(f: OnFail, script: Script) {
  const retry = f.retries ? `↻${f.retries} ` : '';
  const target = f.goto ? flatOf(script.steps, f.goto) : undefined;
  return `${retry}${f.action === 'next' ? '→ bỏ qua' : f.action === 'goto' ? (target ? `→ bước ${target.path}` : '→ ?') : f.action === 'else' ? '→ chạy bước khác' : 'không thấy → dừng'}`;
}

// inMulti: có điểm đang chọn nhiều (bước 30) – Shift / ⌘ / Ctrl + click thẻ thêm / bỏ điểm của bước; click thường bỏ chọn nhiều.
function Card({ row, selected, inMulti, script, folded, onFold }: {
  row: Extract<Row, { kind: 'step' }>; selected: boolean; inMulti: boolean; script: Script; folded: boolean; onFold: () => void;
}) {
  const { step: s, path, depth, active } = row.f;
  const inputs = stepInputCount(script, s.id);
  const failedId = useEngine().failedStep;
  // lỗi ở bước con mà khối đang gập → viền khối
  const failed = !!failedId && (failedId === s.id || (folded && (!!findStep(s.children ?? [], failedId) || !!findStep(s.else ?? [], failedId))));
  // ảnh mẫu thu nhỏ: nhìn danh sách là biết bước nào chờ ảnh nào
  const tpl = s.template && s.findBy !== 'text' && (s.type === 'waitImage' || (s.type === 'clickImage' && s.imageSource === 'search'));
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id: s.id });
  const iconBtn = 'grid h-7 w-7 place-items-center rounded-md text-muted hover:bg-hover hover:text-ink';
  const block = isBlock(s);
  const foldable = block || s.onFail?.action === 'else';
  const inner = (s.children?.length ?? 0) + (s.else?.length ?? 0);

  return (
    <li ref={setNodeRef} style={{ transform: CSS.Translate.toString(transform), transition, paddingLeft: depth * INDENT }}
      className={`relative ${isDragging ? 'z-10' : ''}`}>
      <Rails depth={depth} />
      <div onClick={(e) => {
          if (e.shiftKey || e.metaKey || e.ctrlKey) return toggleMultiStep(s.id);
          selectStep(s.id);
          setMulti([]);
        }}
        title={failed ? 'Lần chạy trước dừng ở bước này' : inMulti ? 'Có điểm đang chọn nhiều – Shift / ⌘ + click để bỏ' : undefined}
        className={`${cardClass(selected)} ${block && !selected ? 'border-dashed border-line-strong! bg-group!' : ''} ${inMulti ? 'outline-2 outline-offset-1 outline-accent' : ''} ${failed ? 'outline-2 outline-offset-1 outline-danger' : ''} ${isDragging ? 'shadow-[0_6px_20px_rgb(0_0_0/.18)]' : ''}`}>
        {/* ô icon = tay kéo để sắp xếp */}
        <button ref={setActivatorNodeRef} {...listeners} {...attributes} aria-label={`Kéo để sắp xếp bước ${path}`} title="Kéo để sắp xếp – thả vào trong khối để đưa bước vào khối"
          className={`${tileClass(selected)} cursor-grab touch-none active:cursor-grabbing ${active ? '' : 'opacity-50'}`}>
          <StepIcon type={s.type} />
        </button>
        <CardBody n={depth ? path : path.padStart(2, '0')} enabled={active}
          title={`${STEP_LABEL[s.type]}${s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}`}
          chips={[...(inputs ? [`${inputs} input`] : []), ...(s.onFail ? [failChip(s.onFail, script)] : []), ...(foldable && folded ? [`${inner} bước bên trong`] : [])]}
          tags={<>
            {!s.enabled && <span className="rounded-full bg-soft px-[7px] py-px text-[11px] text-muted">Tắt</span>}
            {s.enabled && !active && <span className="rounded-full bg-soft px-[7px] py-px text-[11px] text-muted" title="Khối chứa bước này đang tắt">Khối tắt</span>}
            {failed && <span className="rounded-full bg-danger-soft px-[7px] py-px text-[11px] font-semibold text-danger-text">Dừng ở đây</span>}
          </>}
          media={tpl ? <img src={templateUrl(s.template)} alt="Ảnh mẫu" className="h-7 max-w-[72px] rounded border border-line object-contain" /> : undefined}
          detail={fmtPos(s)} delay={(block || s.type === 'break' || s.type === 'stop') && !s.delayAfter ? '' : `${fmtMs(s.delayAfter)}${s.jitterMs > 0 ? ` ±${s.jitterMs}` : ''}`} />
        {/* nút nổi ở góc phải khi rê chuột – không đẩy bố cục */}
        <span onClick={(e) => e.stopPropagation()} className="absolute inset-y-0 right-2 hidden items-center gap-0.5 group-hover:flex">
          {foldable && (
            <button onClick={onFold} aria-label={`${folded ? 'Mở' : 'Gập'} khối ${path}`} title={folded ? 'Mở khối' : 'Gập khối'} className={iconBtn}>
              <svg viewBox="0 0 24 24" className={`h-4 w-4 transition-transform ${folded ? '-rotate-90' : ''}`} fill="none" stroke="currentColor" strokeWidth={1.8} strokeLinecap="round" strokeLinejoin="round"><path d="M6 9l6 6 6-6" /></svg>
            </button>
          )}
          <button onClick={() => toggleStep(s.id)} aria-label={`${s.enabled ? 'Tắt' : 'Bật'} bước ${path}`} title={s.enabled ? 'Tắt bước' : 'Bật bước'}
            className={`${iconBtn} ${s.enabled ? '' : 'text-accent-text'}`}><Power /></button>
          <button onClick={() => duplicateStep(s.id)} aria-label={`Nhân bản bước ${path}`} title={block ? 'Nhân bản (cả bước bên trong)' : 'Nhân bản'} className={iconBtn}><Copy /></button>
          <button onClick={() => removeStep(s.id)} aria-label={`Xoá bước ${path}`} title={block ? 'Xoá (cả bước bên trong) – Delete' : 'Xoá (Delete)'} className={`${iconBtn} hover:text-danger-text`}><Trash /></button>
        </span>
      </div>
    </li>
  );
}

