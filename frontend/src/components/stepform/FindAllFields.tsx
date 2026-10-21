// Node Tìm tất cả (bước 42 G3): cùng ô tìm như Chờ ảnh / chữ / phần tử (ImageSearchFields) + số chỗ cần thấy.
import NumInput from '../NumInput';
import { updateStep } from '../../store';
import { countRange } from '../../flow/findAll';
import { Step } from '../../types';
import { Field, hint } from './shared';
import { ImageSearchFields } from './ImageFields';

export function FindAllFields({ step: s }: { step: Step }) {
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  return (
    <>
      <ImageSearchFields step={s} />
      <Field label="Số chỗ cần thấy" wide group
        help="Chưa đủ tối thiểu (hoặc nhiều hơn tối đa) thì nhìn lại tới hết thời gian chờ rồi đi cổng Thất bại. Tối đa 0 = không giới hạn.">
        <div className="flex gap-2">
          <NumInput prefix="ít nhất" value={Math.max(s.minCount ?? 0, 1)} min={1} onChange={(n) => set({ minCount: n })} className="flex-1" />
          <NumInput prefix="tối đa" value={s.maxCount ?? 0} onChange={(n) => set({ maxCount: n || undefined })} className="flex-1" />
        </div>
        <span className={hint}>
          Cần {countRange(s)} chỗ. Lặp từng chỗ: nối vào Lặp qua danh sách với danh sách {`{{${s.ref}.output.items}}`}, rồi click theo item.x / item.y.
        </span>
      </Field>
    </>
  );
}

