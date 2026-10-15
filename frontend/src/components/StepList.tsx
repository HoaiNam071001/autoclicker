// Danh sách bước: chọn, kéo ⋮⋮ để sắp xếp, bật/tắt, nhân bản, xoá, thêm bước (bước 14)
import { useState } from 'react';
import {
  closestCenter, DndContext, DragEndEvent, KeyboardSensor, PointerSensor, useSensor, useSensors,
} from '@dnd-kit/core';
import {
  SortableContext, sortableKeyboardCoordinates, useSortable, verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Copy, Trash } from './icons';
import {
  addStep, duplicateStep, moveStep, removeStep, selectStep, toggleStep, undoDelete, useEditor,
} from '../scriptStore';
import { fmtMs, fmtPos, Step, STEP_LABEL } from '../types';
import { openOverlay } from '../windowMode';

const ADD_MENU: { label: string; run: () => unknown }[] = [
  { label: 'Click – chọn trên màn hình', run: () => openOverlay('point') },
  { label: 'Kéo thả – chọn A → B', run: () => openOverlay('drag') },
  { label: 'Chờ', run: () => addStep('wait') },
  { label: 'Gõ chữ', run: () => addStep('type') },
  { label: 'Phím tắt', run: () => addStep('key') },
];

export default function StepList() {
  const { script: { steps }, selected, deleted } = useEditor();
  const [menu, setMenu] = useState(false);
  // distance: click vào ⋮⋮ không kéo thì không tính là kéo
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 4 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  const onDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return;
    moveStep(steps.findIndex((s) => s.id === active.id), steps.findIndex((s) => s.id === over.id));
  };

  return (
    <section aria-label="Danh sách bước" className="flex max-w-[420px] flex-[1_1_340px] flex-col rounded-[10px] border border-line bg-white">
      <div className="flex items-center justify-between border-b border-line px-4 py-3.5">
        <h2 className="text-[15px] font-semibold">Các bước</h2>
        <span className="text-muted">{steps.length} bước</span>
      </div>

      {steps.length === 0 && <p className="px-4 py-6 text-center text-muted">Chưa có bước nào – bấm “Thêm điểm” hoặc “+ Thêm bước”.</p>}
      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={onDragEnd}>
        <SortableContext items={steps.map((s) => s.id)} strategy={verticalListSortingStrategy}>
          <ol className="flex flex-col gap-0.5 p-1.5">
            {steps.map((s, i) => <Row key={s.id} step={s} index={i} selected={s.id === selected} />)}
          </ol>
        </SortableContext>
      </DndContext>

      {deleted && (
        <div className="mx-3 mb-2 flex items-center justify-between rounded-lg bg-ground px-3 py-2 text-[13px]">
          <span>Đã xoá bước {deleted.index + 1} ({STEP_LABEL[deleted.step.type]})</span>
          <button onClick={undoDelete} className="font-semibold text-accent">Hoàn tác</button>
        </div>
      )}

      <div className="relative mx-3 mb-3 mt-auto">
        {menu && (
          <ul className="absolute inset-x-0 bottom-full mb-1.5 overflow-hidden rounded-lg border border-line bg-white shadow-[0_6px_24px_rgba(0,0,0,.12)]">
            {ADD_MENU.map((m) => (
              <li key={m.label}>
                <button onClick={() => { setMenu(false); m.run(); }} className="w-full px-3.5 py-2.5 text-left hover:bg-ground">{m.label}</button>
              </li>
            ))}
          </ul>
        )}
        <button onClick={() => setMenu(!menu)} aria-expanded={menu}
          className="h-11 w-full rounded-lg border border-dashed border-[#B8C0CC] text-[#3F4857]">+ Thêm bước</button>
      </div>
    </section>
  );
}

function Row({ step: s, index: i, selected }: { step: Step; index: number; selected: boolean }) {
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id: s.id });
  const ring = s.type === 'drag' ? 'border-drag text-drag' : 'border-accent text-accent';
  const iconBtn = 'grid h-7 w-7 place-items-center rounded-md text-muted hover:bg-white hover:text-ink';

  return (
    <li ref={setNodeRef} style={{ transform: CSS.Translate.toString(transform), transition }}
      className={`group relative flex items-center gap-2 rounded-lg p-2 ${selected ? 'bg-accent-soft' : 'hover:bg-ground'} ${isDragging ? 'z-10 shadow-[0_4px_16px_rgba(0,0,0,.15)]' : ''}`}>
      <button ref={setActivatorNodeRef} {...listeners} {...attributes} aria-label={`Kéo để sắp xếp bước ${i + 1}`}
        className="cursor-grab touch-none px-0.5 text-[#9AA2AE] active:cursor-grabbing">⋮⋮</button>
      <input type="checkbox" checked={s.enabled} onChange={() => toggleStep(s.id)} aria-label={`Bật bước ${i + 1}`}
        className="h-4 w-4 flex-none accent-accent" />
      <button onClick={() => selectStep(s.id)} className={`flex min-w-0 flex-1 items-center gap-2.5 text-left ${s.enabled ? '' : 'opacity-50'}`}>
        <span className={`grid h-7 w-7 flex-none place-items-center rounded-full border-2 text-[13px] font-semibold ${ring}`}>{i + 1}</span>
        <span className="min-w-0 flex-1">
          <span className="block font-semibold">{STEP_LABEL[s.type]}{s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}</span>
          <span className="block truncate font-mono text-xs text-muted">{fmtPos(s)}</span>
        </span>
        <span className="font-mono text-xs text-muted group-hover:hidden">{fmtMs(s.delayAfter)}</span>
      </button>
      {/* hiện khi rê chuột vào dòng (che chỗ thời gian chờ) */}
      <span className="hidden gap-0.5 group-hover:flex">
        <button onClick={() => duplicateStep(s.id)} aria-label={`Nhân bản bước ${i + 1}`} title="Nhân bản" className={iconBtn}><Copy /></button>
        <button onClick={() => removeStep(s.id)} aria-label={`Xoá bước ${i + 1}`} title="Xoá (Delete)" className={`${iconBtn} hover:text-danger`}><Trash /></button>
      </span>
    </li>
  );
}
