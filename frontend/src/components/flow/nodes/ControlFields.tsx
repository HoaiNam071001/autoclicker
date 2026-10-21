// Inspector node điều khiển của flow (bước 20 F3): Nếu, Switch theo giá trị, Switch theo màn hình, Đếm, Kết thúc,
// Đặt biến, Tính toán, xử lý chữ. Chỉ hiển thị: nhận node + set(patch) – app thật nối store (stepform/ControlSection), màn mẫu nối state cục bộ.
import Select from '../../Select';
import NumInput, { fieldBox } from '../../NumInput';
import Segmented from '../../Segmented';
import { errText, Field, hint } from '../../stepform/shared';
import { Trash } from '../../icons';
import VarTextInput from '../../inputs/VarTextInput';
import CondEditor, { CondCtx } from './CondEditor';
import ForEachFields from './ForEachFields';
import CalcFields from './CalcFields';
import TextNodeFields from './TextNodeFields';
import { isTextNode } from '../../../flow/textNodes';
import { newCase } from '../../../flow/cond';
import type { Case, CaseOp, Step } from '../../../types';

export type { CondCtx, CondTools, VarProps } from './CondEditor';

const textBox = `${fieldBox} w-full text-ink outline-none`;
const sec = (ms?: number) => `${((ms ?? 0) / 1000).toLocaleString('vi-VN', { maximumFractionDigits: 1 })} s`;
const Grid = ({ children }: { children: React.ReactNode }) => <div className="grid grid-cols-2 gap-3">{children}</div>;

export default function ControlFields({ ctx }: { ctx: CondCtx }) {
  switch (ctx.step.type) {
    case 'if': return <IfFields ctx={ctx} />;
    case 'switch': return <SwitchFields ctx={ctx} />;
    case 'switchScreen': return <SwitchScreenFields ctx={ctx} />;
    case 'counter': return <CounterFields ctx={ctx} />;
    case 'forEach': return <ForEachFields ctx={ctx} />;
    case 'end': return <EndFields ctx={ctx} />;
    case 'setVar': return <SetVarFields ctx={ctx} />;
    case 'calc': return <CalcFields ctx={ctx} />;
    default: return isTextNode(ctx.step.type) ? <TextNodeFields ctx={ctx} type={ctx.step.type} /> : null;
  }
}

// ---- Nếu ----
function IfFields({ ctx }: { ctx: CondCtx }) {
  const { step: s, set } = ctx;
  return (
    <Grid>
      <CondEditor ctx={ctx} at={{ path: [] }} />
      <Field label="Chờ tối đa" help={s.timeoutMs ? `Xem lại tới khi đúng, quá ${sec(s.timeoutMs)} thì là sai` : 'Xem một lần (0) – không đợi'}>
        <NumInput value={s.timeoutMs ?? 0} unit="ms" onChange={(timeoutMs) => set({ timeoutMs: timeoutMs || undefined, intervalMs: s.intervalMs || 250 })} />
      </Field>
    </Grid>
  );
}

// ---- Danh sách case (dùng chung 2 loại Switch) ----
function CaseList({ ctx, children, last }: { ctx: CondCtx; children: (c: Case, i: number) => React.ReactNode; last: string }) {
  const { step: s, set } = ctx;
  const cases = s.cases ?? [];
  const setCase = (i: number, p: Partial<Case>) => set({ cases: cases.map((c, j) => (j === i ? { ...c, ...p } : c)) });
  const move = (i: number, d: -1 | 1) => {
    const cs = [...cases];
    [cs[i], cs[i + d]] = [cs[i + d], cs[i]];
    set({ cases: cs });
  };
  const icon = 'grid h-6 w-6 place-items-center rounded text-muted enabled:hover:bg-hover enabled:hover:text-ink disabled:opacity-30';
  return (
    <div className="flex flex-col gap-2.5">
      {cases.map((c, i) => (
        <div key={c.key} className="flex flex-col gap-2 rounded-card border border-line p-2.5">
          <div className="flex items-center gap-1.5">
            <span className="h-2 w-2 flex-none rounded-full bg-muted" />
            <input value={c.label ?? ''} placeholder={c.key} onChange={(e) => setCase(i, { label: e.target.value || undefined })} aria-label="Tên case"
              title={`Tên case – hiện trên cổng ra (mã cổng case:${c.key})`} className="min-w-0 flex-1 bg-transparent font-semibold outline-none" />
            <button disabled={i === 0} onClick={() => move(i, -1)} aria-label="Lên" title="Xét trước" className={icon}>↑</button>
            <button disabled={i === cases.length - 1} onClick={() => move(i, 1)} aria-label="Xuống" title="Xét sau" className={icon}>↓</button>
            <button onClick={() => set({ cases: cases.filter((_, j) => j !== i) })} aria-label="Xoá case" title="Xoá case (bỏ cạnh của cổng này)"
              className={`${icon} hover:text-danger-text`}><Trash className="h-3.5 w-3.5" /></button>
          </div>
          {children(c, i)}
        </div>
      ))}
      <button onClick={() => set({ cases: [...cases, newCase(cases, s.type === 'switchScreen')] })}
        className="h-8 rounded-control border border-dashed border-line-strong text-[13px] text-muted hover:bg-hover hover:text-ink">+ Thêm case</button>
      <span className={hint}>Xét lần lượt từ trên xuống; {last}.</span>
    </div>
  );
}

// ---- Switch theo giá trị ----
const CASE_OPS: { id: CaseOp; label: string }[] = [{ id: 'eq', label: 'bằng' }, { id: 'contains', label: 'chứa' }, { id: 'regex', label: 'khớp mẫu' }];
function SwitchFields({ ctx }: { ctx: CondCtx }) {
  const { step: s, set, vars } = ctx;
  const cases = s.cases ?? [];
  const setCase = (i: number, p: Partial<Case>) => set({ cases: cases.map((c, j) => (j === i ? { ...c, ...p } : c)) });
  return (
    <div className="flex flex-col gap-3">
      <Field label="Giá trị cần xét" wide group help="Bằng / chứa không phân biệt hoa thường; khớp mẫu là biểu thức chính quy (thêm (?i) để bỏ qua hoa thường).">
        <VarTextInput value={s.value ?? ''} onChange={(value) => set({ value })} rows={1} placeholder="{{n2.output.text}}" {...vars} />
      </Field>
      <CaseList ctx={ctx} last="không case nào khớp → cổng Khác">
        {(c, i) => (
          <div className="grid grid-cols-[96px_minmax(0,1fr)] gap-2">
            <Select value={c.op ?? 'eq'} onChange={(op) => setCase(i, { op })} aria-label="Cách so" className={textBox}
              options={CASE_OPS.map((o) => ({ value: o.id, label: o.label }))} />
            <VarTextInput value={c.match ?? ''} onChange={(match) => setCase(i, { match })} rows={1} placeholder={c.op === 'regex' ? '(?i)lỗi|error' : 'giá trị'} {...vars} />
          </div>
        )}
      </CaseList>
    </div>
  );
}

// ---- Switch theo màn hình ----
function SwitchScreenFields({ ctx }: { ctx: CondCtx }) {
  const { step: s, set } = ctx;
  return (
    <div className="flex flex-col gap-3">
      <p className="text-[13px] text-muted">Chụp màn hình, case nào trong danh sách hiện trước thì đi ra cổng của case đó.</p>
      <CaseList ctx={ctx} last={`hết ${sec(s.timeoutMs)} không case nào hiện → cổng Không thấy`}>
        {(_, i) => <div className="grid grid-cols-2 gap-3"><CondEditor ctx={ctx} at={{ case: i, path: [] }} screenOnly /></div>}
      </CaseList>
      <Field label="Chờ tối đa"><NumInput value={s.timeoutMs ?? 0} unit="ms" onChange={(timeoutMs) => set({ timeoutMs })} /></Field>
      <Field label="Xem lại mỗi"><NumInput value={s.intervalMs ?? 250} min={10} unit="ms" onChange={(intervalMs) => set({ intervalMs })} /></Field>
    </div>
  );
}

// ---- Đếm ----
function CounterFields({ ctx: { step: s, set } }: { ctx: CondCtx }) {
  const max = s.max ?? 0;
  return (
    <Field label="Tối đa" help={`Mỗi lần đi qua cộng 1: lần 1–${max} → cổng Còn lượt, lần ${max + 1} → Hết lượt (rồi đếm lại từ đầu). Đặt lại mỗi vòng kịch bản. Dùng để chặn vòng lặp vô hạn.`}>
      <NumInput value={max} min={1} unit="lần" stepper onChange={(m) => set({ max: m })} />
    </Field>
  );
}

// ---- Kết thúc ----
function EndFields({ ctx: { step: s, set, vars } }: { ctx: CondCtx }) {
  const fail = s.endStatus === 'fail';
  return (
    <Grid>
      <Field label="Kết quả" wide group help={fail ? 'Dừng kịch bản và báo lỗi kèm lời nhắn (nhật ký, thanh mini).' : 'Kết thúc vòng này bình thường – còn vòng lặp thì chạy vòng sau.'}>
        <Segmented value={fail ? 'fail' : 'success'} onChange={(endStatus) => set({ endStatus })} full label="Kết quả"
          options={[{ id: 'success', label: 'Thành công' }, { id: 'fail', label: 'Thất bại' }]} />
      </Field>
      <Field label="Lời nhắn" wide group>
        <VarTextInput value={s.message ?? ''} onChange={(message) => set({ message })} rows={2} placeholder={fail ? 'Không thấy nút Nhắn tin' : ''} {...vars} />
      </Field>
    </Grid>
  );
}

// ---- Đặt biến ----
function SetVarFields({ ctx: { step: s, set, vars } }: { ctx: CondCtx }) {
  const ok = /^[A-Za-z_][A-Za-z0-9_]{0,31}$/.test(s.varName ?? '');
  return (
    <Grid>
      <Field label="Tên biến">
        <input value={s.varName ?? ''} onChange={(e) => set({ varName: e.target.value.trim() })} spellCheck={false}
          className={`${textBox} font-mono text-[13px] ${ok ? '' : 'border-danger'}`} />
      </Field>
      <Field label="Cách đặt" group>
        <Segmented value={s.varOp === 'add' ? 'add' : 'set'} onChange={(varOp) => set({ varOp })} full label="Cách đặt"
          options={[{ id: 'set', label: 'Gán' }, { id: 'add', label: 'Cộng' }]} />
      </Field>
      <Field label={s.varOp === 'add' ? 'Cộng thêm (số)' : 'Giá trị'} wide group>
        <VarTextInput value={s.value ?? ''} onChange={(value) => set({ value })} rows={1} placeholder={s.varOp === 'add' ? '1' : '{{n2.output.text}}'} {...vars} />
      </Field>
      {ok ? <span className={`col-span-2 ${hint}`}>Dùng ở node sau: <span className="var-token">{`vars.${s.varName}`}</span>{s.varOp === 'add' ? ' – chưa có thì tính từ 0' : ''}. Đặt lại mỗi vòng.</span>
        : <span className={`col-span-2 ${errText}`}>Tên biến: chữ cái / _ đầu, sau đó chữ / số / _</span>}
    </Grid>
  );
}

// Node có nút Thử (xét trên màn hình / giá trị lúc này).
export const canTest = (s: Pick<Step, 'type'>) => s.type === 'if' || s.type === 'switch' || s.type === 'switchScreen';
