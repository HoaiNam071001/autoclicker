// Thẻ một kịch bản mẫu: tên, độ khó, nhãn web thật, mô tả, các loại bước được test, nút mở trang / dùng mẫu.
import { StepIcon, StepIconName, STEP_ICON } from '../icons';
import { levelText, Sample, stepName } from '../../samples/filter';
import { categoryOfType } from '../../flow/category';

const LEVEL_TONE = ['', 'text-ok', 'text-warn', 'text-danger'];

export default function SampleCard({ sample, ready, onPick, onOpenPage }: {
  sample: Sample; ready: boolean; onPick: () => void; onOpenPage: () => void;
}) {
  const off = ready ? undefined : 'Cần chạy trong app (wails dev)';
  return (
    <div className="flex flex-col gap-2 rounded-[10px] border border-line bg-surface p-3">
      <div className="flex items-start gap-2">
        <span className="flex-1 font-semibold">{sample.name}</span>
        {sample.real && <span className="rounded-full bg-soft px-2 py-0.5 text-[11px] text-warn" title="Chạy trên trang thật – có thể hỏng khi trang đổi giao diện">web thật</span>}
        <span className={`text-xs font-medium ${LEVEL_TONE[sample.level] ?? ''}`}>{levelText(sample.level)}</span>
      </div>
      <p className="text-[13px] leading-snug text-muted">{sample.description}</p>
      {sample.prep && <p className="text-xs text-muted">Chuẩn bị: {sample.prep}</p>}
      <div className="flex flex-wrap gap-1">
        {sample.steps.map((t) => (
          <span key={t} className={`${categoryOfType(t) ? `cat-${categoryOfType(t)}` : ''} flex items-center gap-1 rounded-md bg-soft px-1.5 py-0.5 text-[11px] text-muted`}>
            {t in STEP_ICON && <StepIcon type={t as StepIconName} className="h-3 w-3 text-cat" />}{stepName(t)}
          </span>
        ))}
      </div>
      <div className="mt-auto flex items-center gap-2 pt-1">
        <span className="flex-1 truncate font-mono text-[11px] text-muted" title={sample.page}>{sample.page}</span>
        {sample.page && <button disabled={!ready} title={off ?? 'Mở trang đích trong trình duyệt'} onClick={onOpenPage}
          className="h-8 rounded-control px-2.5 text-[13px] text-muted enabled:hover:bg-hover enabled:hover:text-ink disabled:opacity-60">Mở trang</button>}
        <button disabled={!ready} title={off ?? 'Mở bản sao trong trình soạn (chưa lưu)'} onClick={onPick}
          className="btn-primary h-8 rounded-control px-3 text-[13px] font-medium disabled:opacity-60">Dùng mẫu</button>
      </div>
    </div>
  );
}
