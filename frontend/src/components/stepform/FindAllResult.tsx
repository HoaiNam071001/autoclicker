// Kết quả Thử tìm của node Tìm tất cả (42 G3) – thay dòng "Thấy tại X Y" trong ô tìm ảnh / chữ / phần tử: số chỗ + xem khung đánh số.
import { findAllSummary } from '../../flow/findAll';
import { openOverlay } from '../../windowMode';
import type { FindResult } from '../../imageSteps';
import { Step } from '../../types';

export function FindAllResult({ step: s, find }: { step: Step; find: FindResult }) {
  return (
    <span className={`flex flex-col gap-1 text-xs ${find.found ? 'text-ok' : 'text-warn'}`}>
      <span className="flex flex-wrap items-center gap-x-2">
        {findAllSummary(find, s)}
        {!!find.all?.length && <button onClick={() => openOverlay('point', s.id)} className="text-muted underline">xem trên màn hình</button>}
      </span>
      {!find.count && !!find.read?.length && <span className="text-muted">Đọc được: {find.read.slice(0, 12).join(' · ')}{find.read.length > 12 ? ' …' : ''}</span>}
    </span>
  );
}
