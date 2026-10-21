// Inspector node Tính toán (bước 42 G4a): ô biểu thức + kết quả xem trước (không còn biến) / lỗi cú pháp, bảng cú pháp
// ngắn. Chỉ hiển thị – nối store ở stepform/ControlSection.
import { Fragment } from 'react';
import VarTextInput from '../../inputs/VarTextInput';
import { errText, Field, hint } from '../../stepform/shared';
import { calcPreview } from '../../../flow/calc';
import type { CondCtx } from './CondEditor';

const SYNTAX: [string, string][] = [
  ['+ - * / %', 'cộng trừ nhân chia, chia lấy dư'],
  ['== != < <= > >=', 'so sánh → true / false'],
  ['and  or  not', 'và / hoặc / không (hoặc && || !)'],
  ['round(x, 2)  floor  ceil', 'làm tròn (2 chữ số), xuống, lên'],
  ['min(a, b)  max  abs', 'nhỏ nhất, lớn nhất, trị tuyệt đối'],
  ['"Hết hàng"', 'chữ trong ngoặc kép – biến chữ tự thành chuỗi'],
];

export default function CalcFields({ ctx: { step: s, set, vars } }: { ctx: CondCtx }) {
  const p = calcPreview(s.value ?? '');
  return (
    <div className="flex flex-col gap-3">
      <Field label="Biểu thức" group help="Số thập phân dùng dấu chấm: 1.5">
        <VarTextInput {...vars} rows={2} value={s.value ?? ''} onChange={(value) => set({ value })} placeholder={'{{n2.output.count}} * 2 + 1'} />
        {p.error ? <span className={errText}>{p.error}</span>
          : p.value !== undefined ? <span className={hint}>Kết quả: <b className="font-mono text-ink">{p.value}</b></span>
          : p.needsRun && <span className={hint}>Kết quả biết khi chạy (có biến)</span>}
      </Field>
      <div className="grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 rounded-card bg-soft px-3 py-2 text-xs text-muted">
        {SYNTAX.map(([code, desc]) => <Fragment key={code}><span className="font-mono text-ink">{code}</span><span>{desc}</span></Fragment>)}
      </div>
      <span className={hint}>
        Node sau dùng <span className="var-token">{`${s.ref}.output.value`}</span> – rẽ nhánh theo kết quả bằng node <b className="text-ink">Nếu</b>.
        Số kiểu <span className="font-mono">1.234.000 đ</span> là chữ: lấy số trước bằng node <b className="text-ink">Số &amp; ngày</b>.
      </span>
    </div>
  );
}
