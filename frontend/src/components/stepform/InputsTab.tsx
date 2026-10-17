import { useEffect } from 'react';
import InputsPanel from '../inputs/InputsPanel';
import { useScript } from '../../store';
import { addInput, inputViews, moveInput, removeInput, updateInput } from '../../scriptInputs';
import { loadMemory, stopAutoRun, useRunForm } from '../../runForm';

// Tab Input: input của cả kịch bản (không riêng bước đang chọn).
export default function InputsTab() {
  const script = useScript();
  const { memory } = useRunForm();
  useEffect(() => {
    loadMemory();
  }, [script.name]);
  return (
    <InputsPanel inputs={inputViews(script)} onUpdate={updateInput} onRemove={removeInput} onMove={moveInput}
      onAdd={() => addInput({ label: 'Input mới' })}
      footer={memory.autoRun && !!script.inputs?.length && (
        <div className="flex items-center gap-2 rounded-lg bg-soft px-3 py-2 text-xs text-muted">
          <span className="flex-1">Bấm Chạy đang chạy luôn với giá trị lần trước, không hỏi (Shift + bấm Chạy để sửa)</span>
          <button onClick={stopAutoRun} className="font-semibold text-accent-text">Hỏi lại</button>
        </div>
      )} />
  );
}
