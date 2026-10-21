// Đổi loại node từ inspector: cùng bảng với Thêm node (có hoạt hoạ minh hoạ) nhưng chỉ các loại node thuần và bấm là đổi node đang chọn.
import { useMemo } from 'react';
import { PaletteView } from '../PaletteView';
import { GROUPS } from './groups';
import { changeStepType } from '../../store';
import type { Step } from '../../types';

export default function ChangeTypePalette({ step, onClose }: { step: Step; onClose: () => void }) {
  const groups = useMemo(() => GROUPS
    .map((g) => ({
      ...g, note: undefined,
      items: g.items.filter((it) => it.type).map((it) => ({ ...it, run: () => changeStepType(step.id, it.type!) })),
    }))
    .filter((g) => g.items.length > 0), [step.id]);
  return (
    <PaletteView groups={groups} onClose={onClose} label="Đổi loại node"
      context={<>Đổi loại của <span className="font-mono">{step.ref}</span> – cài đặt chung (vị trí, thời gian, nối cạnh) được giữ</>}
      footer={<span><span className="font-mono">Enter</span> đổi sang loại này</span>} />
  );
}
