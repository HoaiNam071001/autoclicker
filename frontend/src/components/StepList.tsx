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
  duplicateStep, moveStep, removeStep, selectStep, setMulti, setRegion, toggleMultiStep, toggleStep, undoDelete, useEditor,
} from '../scriptStore';
import { fmtMs, fmtPos, Script, Step, STEP_LABEL } from '../types';
import { stepInputCount } from '../scriptInputs';
import { templateUrl } from '../imageSteps';
import { useEngine } from '../engineStore';
import { openOverlay } from '../windowMode';

import AddStepPalette from './AddStepPalette';
import { CardBody, cardClass, tileClass } from './StepCard';

const isMac = navigator.platform.startsWith('Mac');

export default function StepList() {
  const { script, selected, deleted, multi } = useEditor();
  const { steps, region } = script;
  const multiSteps = new Set(multi.map((k) => k.split(':')[0]));
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
            {steps.map((s, i) => <Card key={s.id} step={s} index={i} selected={s.id === selected} inMulti={multiSteps.has(s.id)} script={script} />)}
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

// inMulti: có điểm đang chọn nhiều (bước 30) – Shift / ⌘ / Ctrl + click thẻ thêm / bỏ điểm của bước; click thường bỏ chọn nhiều.
function Card({ step: s, index: i, selected, inMulti, script }: { step: Step; index: number; selected: boolean; inMulti: boolean; script: Script }) {
  const inputs = stepInputCount(script, s.id);
  const failed = useEngine().failedStep === s.id;
  // ảnh mẫu thu nhỏ: nhìn danh sách là biết bước nào chờ ảnh nào
  const tpl = s.template && (s.type === 'waitImage' || (s.type === 'clickImage' && s.imageSource === 'search'));
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id: s.id });
  const iconBtn = 'grid h-7 w-7 place-items-center rounded-md text-muted hover:bg-hover hover:text-ink';

  return (
    <li ref={setNodeRef} style={{ transform: CSS.Translate.toString(transform), transition }} onClick={(e) => {
        if (e.shiftKey || e.metaKey || e.ctrlKey) return toggleMultiStep(s.id);
        selectStep(s.id);
        setMulti([]);
      }}
      title={failed ? 'Lần chạy trước dừng ở bước này' : inMulti ? 'Có điểm đang chọn nhiều – Shift / ⌘ + click để bỏ' : undefined}
      className={`${cardClass(selected)} ${inMulti ? 'outline-2 outline-offset-1 outline-accent' : ''} ${failed ? 'outline-2 outline-offset-1 outline-danger' : ''} ${isDragging ? 'z-10 shadow-[0_6px_20px_rgb(0_0_0/.18)]' : ''}`}>
      {/* ô icon = tay kéo để sắp xếp */}
      <button ref={setActivatorNodeRef} {...listeners} {...attributes} aria-label={`Kéo để sắp xếp bước ${i + 1}`} title="Kéo để sắp xếp"
        className={`${tileClass(selected)} cursor-grab touch-none active:cursor-grabbing ${s.enabled ? '' : 'opacity-50'}`}>
        <StepIcon type={s.type} />
      </button>
      <CardBody n={String(i + 1).padStart(2, '0')} enabled={s.enabled}
        title={`${STEP_LABEL[s.type]}${s.type === 'click' && s.clicks > 1 ? ` ×${s.clicks}` : ''}`}
        chips={inputs ? [`${inputs} input`] : undefined}
        tags={<>
          {!s.enabled && <span className="rounded-full bg-soft px-[7px] py-px text-[11px] text-muted">Tắt</span>}
          {failed && <span className="rounded-full bg-danger-soft px-[7px] py-px text-[11px] font-semibold text-danger-text">Dừng ở đây</span>}
        </>}
        media={tpl ? <img src={templateUrl(s.template)} alt="Ảnh mẫu" className="h-7 max-w-[72px] rounded border border-line object-contain" /> : undefined}
        detail={fmtPos(s)} delay={`${fmtMs(s.delayAfter)}${s.jitterMs > 0 ? ` ±${s.jitterMs}` : ''}`} />
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
