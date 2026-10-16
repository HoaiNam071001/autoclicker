// Danh sách bước (design v2: thẻ bước): chọn, kéo ô icon để sắp xếp, bật/tắt, nhân bản, xoá, thêm bước (bước 14)
import { useEffect, useState } from 'react';
import {
  closestCenter, DndContext, DragEndEvent, KeyboardSensor, PointerSensor, useSensor, useSensors,
} from '@dnd-kit/core';
import {
  SortableContext, sortableKeyboardCoordinates, useSortable, verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Copy, Power, Region, StepIcon, Trash } from './icons';
import {
  duplicateStep, moveStep, removeStep, selectStep, setRegion, toggleStep, undoDelete, useEditor,
} from '../scriptStore';
import { fmtMs, fmtPos, Step, STEP_LABEL } from '../types';
import { openOverlay } from '../windowMode';

import AddStepPalette from './AddStepPalette';

const isMac = navigator.platform.startsWith('Mac');

export default function StepList() {
  const { script: { steps, region }, selected, deleted } = useEditor();
  const [palette, setPalette] = useState(false);
  // Ctrl+K (macOS: Cmd+K) mở / đóng bảng thêm bước
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key.toLowerCase() !== 'k' || !(isMac ? e.metaKey : e.ctrlKey)) return;
      e.preventDefault();
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

  const onDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return;
    moveStep(steps.findIndex((s) => s.id === active.id), steps.findIndex((s) => s.id === over.id));
  };

  return (
    <main aria-label="Danh sách bước" className="flex min-w-0 flex-1 flex-col gap-3 overflow-y-auto px-[22px] py-4">
      <div className="flex flex-wrap items-center gap-2.5">
        <h1 className="text-base font-semibold">Các bước</h1>
        <span className="text-muted">{steps.length} bước</span>
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
      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={onDragEnd}>
        <SortableContext items={steps.map((s) => s.id)} strategy={verticalListSortingStrategy}>
          <ol className="flex flex-col gap-1.5">
            {steps.map((s, i) => <Card key={s.id} step={s} index={i} selected={s.id === selected} />)}
          </ol>
        </SortableContext>
      </DndContext>

      {deleted && (
        <div className="flex items-center justify-between rounded-card bg-soft px-3 py-2 text-[13px]">
          <span>Đã xoá bước {deleted.index + 1} ({STEP_LABEL[deleted.step.type]})</span>
          <button onClick={undoDelete} className="font-semibold text-accent-text">Hoàn tác</button>
        </div>
      )}

      <button onClick={() => setPalette(true)} aria-haspopup="dialog"
        className="flex h-[42px] flex-none items-center justify-center gap-2.5 rounded-card border border-dashed border-line-strong text-muted hover:bg-hover hover:text-ink">
        + Thêm bước <span className="font-mono text-[11px]">{isMac ? '⌘K' : 'Ctrl+K'}</span>
      </button>
      {palette && <AddStepPalette onClose={() => setPalette(false)} />}
    </main>
  );
}

function Card({ step: s, index: i, selected }: { step: Step; index: number; selected: boolean }) {
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id: s.id });
  const iconBtn = 'grid h-7 w-7 place-items-center rounded-md text-muted hover:bg-hover hover:text-ink';
  const n = String(i + 1).padStart(2, '0');

  return (
    <li ref={setNodeRef} style={{ transform: CSS.Translate.toString(transform), transition }} onClick={() => selectStep(s.id)}
      className={`group relative flex cursor-default items-center gap-3 rounded-card px-3 py-[9px] ${selected ? 'ring-selected' : 'border border-line bg-surface'} ${isDragging ? 'z-10 shadow-[0_6px_20px_rgb(0_0_0/.18)]' : ''}`}>
      {/* ô icon = tay kéo để sắp xếp */}
      <button ref={setActivatorNodeRef} {...listeners} {...attributes} aria-label={`Kéo để sắp xếp bước ${i + 1}`} title="Kéo để sắp xếp"
        className={`grid h-8 w-8 flex-none cursor-grab touch-none place-items-center rounded-lg active:cursor-grabbing ${selected ? 'bg-accent-soft text-accent-text' : 'bg-soft text-muted'} ${s.enabled ? '' : 'opacity-50'}`}>
        <StepIcon type={s.type} />
      </button>
      <span className={`min-w-0 flex-1 ${s.enabled ? '' : 'opacity-50'}`}>
        <span className="flex items-center gap-2">
          <span className="font-mono text-xs text-muted">{n}</span>
          <span className="truncate font-semibold">{STEP_LABEL[s.type]}{s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}</span>
          {!s.enabled && <span className="rounded-full bg-soft px-[7px] py-px text-[11px] text-muted">Tắt</span>}
        </span>
        <span className="block truncate font-mono text-xs text-muted">{fmtPos(s)}</span>
      </span>
      {/* invisible (không phải hidden) để giữ nguyên chiều rộng – rê chuột không làm thẻ nhảy */}
      <span className="rounded-[5px] bg-soft px-2 py-[3px] font-mono text-xs whitespace-nowrap text-muted group-hover:invisible">
        {fmtMs(s.delayAfter)}{s.jitterMs > 0 && ` ±${s.jitterMs}`}
      </span>
      {/* nút nổi ở góc phải khi rê chuột – không đẩy bố cục */}
      <span onClick={(e) => e.stopPropagation()} className="absolute inset-y-0 right-2 hidden items-center gap-0.5 group-hover:flex">
        <button onClick={() => toggleStep(s.id)} aria-label={`${s.enabled ? 'Tắt' : 'Bật'} bước ${i + 1}`} title={s.enabled ? 'Tắt bước' : 'Bật bước'}
          className={`${iconBtn} ${s.enabled ? '' : 'text-accent-text'}`}><Power /></button>
        <button onClick={() => duplicateStep(s.id)} aria-label={`Nhân bản bước ${i + 1}`} title="Nhân bản" className={iconBtn}><Copy /></button>
        <button onClick={() => removeStep(s.id)} aria-label={`Xoá bước ${i + 1}`} title="Xoá (Delete)" className={`${iconBtn} hover:text-danger-text`}><Trash /></button>
      </span>
    </li>
  );
}
