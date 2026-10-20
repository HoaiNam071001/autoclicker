// Cách so ảnh mẫu + kiểm tra màu (bước 40) – chỉ hiển thị, dùng chung cho node Chờ / Click theo ảnh và điều kiện thấy ảnh.
import NumInput from '../NumInput';
import Segmented from '../Segmented';
import { Switch } from '../Switch';
import { MatchMode } from '../../types';
import { Field } from './shared';

// Mức lệch màu mặc định khi bật: nút xám / xanh cùng hình lệch ~30%, cùng nút khác theme sáng nhẹ lệch vài %.
export const DEFAULT_COLOR_TOL = 15;

export type MatchValue = { matchMode?: MatchMode; colorTol?: number };

export default function MatchFields({ value: v, onChange }: { value: MatchValue; onChange: (patch: MatchValue) => void }) {
  const colorOn = (v.colorTol ?? 0) > 0;
  return (
    <>
      <Field label="So theo" wide group
        help="Độ sáng: giống từng điểm ảnh – chắc nhất khi giao diện không đổi. Đường viền: so hình dáng nét – ảnh khoanh ở theme sáng vẫn tìm được ở theme tối / đổi màu nền, nhưng dễ nhầm hơn với ảnh ít nét.">
        <Segmented value={v.matchMode || 'gray'} full label="So theo"
          onChange={(m) => onChange({ matchMode: m === 'edge' ? 'edge' : undefined })}
          options={[{ id: 'gray', label: 'Độ sáng' }, { id: 'edge', label: 'Đường viền' }]} />
      </Field>
      <Field label="Kiểm tra màu" group
        help="Sau khi thấy hình giống, so màu trung bình với ảnh mẫu – không nhận nhầm nút cùng hình khác màu (xám / xanh). Không dùng cùng Đường viền khi theme đổi màu.">
        <Switch on={colorOn} label="Kiểm tra màu" onChange={(on) => onChange({ colorTol: on ? DEFAULT_COLOR_TOL : undefined })} />
      </Field>
      {colorOn && (
        <Field label="Lệch màu tối đa" help="% lệch màu trung bình cho phép. Thử tìm để xem độ lệch thật.">
          <NumInput value={v.colorTol ?? DEFAULT_COLOR_TOL} min={1} unit="%" onChange={(colorTol) => onChange({ colorTol: Math.min(colorTol, 100) })} />
        </Field>
      )}
    </>
  );
}
