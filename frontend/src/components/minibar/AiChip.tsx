// Chip AI trên thanh mini (bước 31b): lượt chạy do AI khởi động – chấm nhấp nháy = AI đang điều khiển. Bấm để mở hàng
// chọn mức (AiLevelRow). Chỉ hiển thị.
import { AI_LEVELS, AiLevel, levelLabel } from '../../aiLevel';
import Segmented from '../Segmented';

export function AiChip({ level, open, onToggle }: { level: AiLevel; open: boolean; onToggle: () => void }) {
  return (
    <button onClick={onToggle} aria-expanded={open} title="AI đang điều khiển – bấm để đổi mức AI"
      className="flex h-[22px] flex-none items-center gap-1.5 rounded-full bg-accent-soft px-2 text-[11.5px] font-semibold text-accent-text hover:brightness-95">
      <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-accent" aria-hidden />
      AI · {levelLabel(level)}
      <span className="font-mono text-[10px] opacity-70" aria-hidden>{open ? '▴' : '▾'}</span>
    </button>
  );
}

// Hàng chọn mức thay cho các dòng chữ khi chip mở. Tắt / Chỉ soạn dừng ngay lượt AI đang chạy (Go – stopAIRun).
export function AiLevelRow({ level, onLevel }: { level: AiLevel; onLevel: (l: AiLevel) => void }) {
  return <Segmented value={level} options={AI_LEVELS} onChange={onLevel} label="Mức AI" />;
}
