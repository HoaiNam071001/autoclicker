// Node Đọc chữ trang (bước 43 B4): phạm vi (cả cửa sổ | trong một phần tử) + gồm phần ngoài khung nhìn.
import NumInput from '../NumInput';
import { updateStep } from '../../store';
import { Step } from '../../types';
import { fmtSec, Field, hint, NumField, SwitchField } from './shared';
import { ElementQueryFields } from './ElementFields';

export default function PageTextFields({ step: s }: { step: Step }) {
  return (
    <>
      <ElementQueryFields step={s} optional />
      <SwitchField label="Gồm phần ngoài khung nhìn" help="Mặc định chỉ đọc chữ đang hiện trên màn (phần cuộn xuống chưa tới bị bỏ). Bật để lấy cả trang dài – cây của trình duyệt vẫn giữ các khối chữ chưa cuộn tới."
        on={!!s.pageAll} onChange={(pageAll) => updateStep(s.id, { pageAll: pageAll || undefined })} />
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 250} min={20} unit="ms" onChange={(intervalMs) => updateStep(s.id, { intervalMs })} />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Không di chuột, không cuộn, không OCR. Node sau dùng {`{{${s.ref}.output.text}}`} (cả khối) hoặc {`{{${s.ref}.output.lines}}`} nối vào Lặp qua danh sách.
        Chưa có chữ nào sau {fmtSec(s.timeoutMs)} → cổng Thất bại. Cần quyền Accessibility.
      </span>
    </>
  );
}
