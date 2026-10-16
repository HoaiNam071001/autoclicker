// Nút ghim cạnh trường số / chữ (artboard StepInputs): ghim = cho phép đổi giá trị khi bấm Chạy (input của kịch bản).
import { Pin } from '../icons';

export default function PinButton({ pinned, onToggle, disabled, title }: {
  pinned: boolean; onToggle?: () => void; disabled?: boolean; title?: string;
}) {
  return (
    <button type="button" aria-pressed={pinned} onClick={onToggle} disabled={disabled}
      aria-label={pinned ? 'Bỏ ghim input' : 'Đặt làm input'}
      title={title ?? (pinned ? 'Bỏ ghim – không hỏi khi chạy nữa' : 'Đặt làm input – hỏi giá trị khi bấm Chạy')}
      className={`grid h-8 w-8 flex-none place-items-center rounded-control border disabled:opacity-50 ${pinned
        ? 'border-accent/45 bg-accent-soft text-accent-text'
        : 'border-line text-muted enabled:hover:bg-hover enabled:hover:text-ink'}`}>
      <Pin />
    </button>
  );
}

// Chip key cạnh nhãn của trường đã ghim
export const KeyChip = ({ k }: { k: string }) => (
  <span className="var-token inline-flex items-center gap-1 py-0! text-[11px]"><Pin className="h-2.5 w-2.5" />{k}</span>
);
