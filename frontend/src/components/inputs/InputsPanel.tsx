// Tab Input của inspector (artboard InputsTab): các giá trị hỏi khi bấm Chạy, theo thứ tự – kéo ⋮⋮ để sắp xếp,
// bấm một dòng để mở rộng sửa nhãn / kiểu / mặc định / bắt buộc / gợi ý, xoá (bỏ ghim) input.
import { useState } from 'react';
import { closestCenter, DndContext, DragEndEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { arrayMove, SortableContext, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Trash } from '../icons';
import { Switch } from '../RunSettings';
import { fieldBox } from '../NumInput';
import { fmtValue, INPUT_TYPE_LABEL, InputType, ScriptInput } from './model';

type Props = { inputs: ScriptInput[]; onChange: (inputs: ScriptInput[]) => void; onAdd?: () => void };

export default function InputsPanel({ inputs, onChange, onAdd }: Props) {
  const [open, setOpen] = useState<string | null>(null);
  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 4 } }));
  const onDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return;
    onChange(arrayMove(inputs, inputs.findIndex((i) => i.key === active.id), inputs.findIndex((i) => i.key === over.id)));
  };
  const update = (key: string, patch: Partial<ScriptInput>) => onChange(inputs.map((i) => (i.key === key ? { ...i, ...patch } : i)));

  return (
    <div className="flex flex-1 flex-col gap-3.5">
      <p className="text-xs text-muted">Các giá trị được hỏi khi bấm Chạy, theo thứ tự dưới đây.</p>
      <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={onDragEnd}>
        <SortableContext items={inputs.map((i) => i.key)} strategy={verticalListSortingStrategy}>
          <div className="flex flex-col gap-1.5">
            {inputs.map((inp) => (
              <Row key={inp.key} inp={inp} open={open === inp.key} onToggle={() => setOpen(open === inp.key ? null : inp.key)}
                onChange={(p) => update(inp.key, p)} onRemove={() => onChange(inputs.filter((i) => i.key !== inp.key))} />
            ))}
          </div>
        </SortableContext>
      </DndContext>
      {inputs.length === 0 && <p className="py-4 text-center text-[13px] text-muted">Chưa có input – bấm ghim cạnh một trường hoặc gõ {'{{'} trong ô chữ.</p>}
      <button onClick={onAdd} disabled={!onAdd}
        className="h-[38px] rounded-lg border border-dashed border-line-strong text-muted enabled:hover:bg-hover enabled:hover:text-ink">+ Thêm input</button>
      <p className="mt-auto text-xs text-muted">Giá trị lần chạy trước lưu trên máy này, không ghi vào file kịch bản.</p>
    </div>
  );
}

function Row({ inp, open, onToggle, onChange, onRemove }: {
  inp: ScriptInput; open: boolean; onToggle: () => void; onChange: (p: Partial<ScriptInput>) => void; onRemove: () => void;
}) {
  const { attributes, listeners, setNodeRef, setActivatorNodeRef, transform, transition, isDragging } = useSortable({ id: inp.key });
  const handle = (
    <span ref={setActivatorNodeRef} {...listeners} {...attributes} aria-label={`Kéo để sắp xếp ${inp.label}`}
      className="cursor-grab touch-none text-muted active:cursor-grabbing" onClick={(e) => e.stopPropagation()}>⋮⋮</span>
  );
  const used = inp.usedBy.map((u) => `B${u.n}`).join(', ');
  const style = { transform: CSS.Translate.toString(transform), transition };
  const lift = isDragging ? 'relative z-10 shadow-[0_6px_20px_rgb(0_0_0/.18)]' : '';

  if (!open) {
    return (
      <div ref={setNodeRef} style={style} onClick={onToggle}
        className={`flex cursor-default items-center gap-2.5 rounded-card border border-line bg-surface px-2.5 py-[9px] hover:bg-hover ${lift}`}>
        {handle}
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-1.5 font-semibold"><span className="truncate">{inp.label}</span>{inp.required && <span className="text-warn">*</span>}</div>
          <div className="truncate font-mono text-xs text-muted">{inp.key}  ·  {inp.required && inp.default === undefined ? 'bắt buộc' : `mặc định ${fmtValue(inp.default, inp.unit)}`}</div>
        </div>
        <span className="rounded-[5px] bg-soft px-[7px] py-0.5 text-[11px] whitespace-nowrap text-muted">{INPUT_TYPE_LABEL[inp.type]}</span>
        {used && <span className="rounded-[5px] border border-line px-1.5 py-px font-mono text-[11px] text-muted">{used}</span>}
      </div>
    );
  }

  const label = 'flex flex-col gap-[5px] text-[13px] text-muted';
  const input = `${fieldBox} text-sm text-ink outline-none`;
  return (
    <div ref={setNodeRef} style={style} className={`ring-selected flex flex-col gap-2.5 rounded-card p-3 ${lift}`}>
      <div className="flex items-center gap-2.5" onClick={onToggle}>
        {handle}
        <span className="var-token text-[13px]">{inp.key}</span>
        <span className="truncate text-xs text-muted">dùng ở {inp.usedBy.map((u) => `B${u.n} ${u.title}`).join(', ') || 'chưa bước nào'}</span>
        <button onClick={(e) => { e.stopPropagation(); onRemove(); }} aria-label="Xoá input" title="Xoá input (bỏ ghim, biến {{…}} trong ô chữ sẽ báo lỗi)"
          className="ml-auto grid h-7 w-7 flex-none place-items-center rounded-md text-muted hover:bg-hover hover:text-danger-text"><Trash className="h-3.5 w-3.5" /></button>
      </div>
      <label className={label}>Nhãn
        <input value={inp.label} onChange={(e) => onChange({ label: e.target.value })} className={input} />
      </label>
      <div className="grid grid-cols-2 gap-2">
        <label className={label}>Kiểu
          <select value={inp.type} onChange={(e) => onChange({ type: e.target.value as InputType })} className={input}>
            {Object.entries(INPUT_TYPE_LABEL).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
          </select>
        </label>
        <label className={label}>Mặc định
          <input value={inp.default === undefined ? '' : String(inp.default)} placeholder="(trống)"
            onChange={(e) => onChange({ default: e.target.value === '' ? undefined : e.target.value })} className={`${input} font-mono`} />
        </label>
      </div>
      <label className="flex items-center justify-between text-[13px]">Bắt buộc nhập
        <Switch on={!!inp.required} onChange={(required) => onChange({ required })} label="Bắt buộc nhập" />
      </label>
      <label className={label}>Gợi ý hiện dưới ô
        <input value={inp.hint ?? ''} onChange={(e) => onChange({ hint: e.target.value })} className={input} />
      </label>
    </div>
  );
}
