// Inspector cho node điều khiển của flow (bước 20): Nếu, Switch theo giá trị, Switch theo màn hình, Đếm, Kết thúc, Đặt biến.
// Mỗi khối nhận value / onChange (như steps/StepFields) để F3 nối thẳng vào store. Kiểu dữ liệu tạm ở đây –
// chuyển vào types.ts + model.go khi làm logic (F3).
import NumInput, { fieldBox } from '../../NumInput';
import Segmented from '../../Segmented';
import { Field } from '../../steps/StepFields';
import { Switch } from '../../RunSettings';

type Ctl<T> = { value: T; onChange: (v: T) => void };
const textBox = `${fieldBox} w-full text-ink outline-none`;
const Grid = ({ children }: { children: React.ReactNode }) => <div className="grid grid-cols-2 gap-3">{children}</div>;
const TestBtn = ({ label }: { label: string }) => (
  <button disabled title="Sắp có (bước 20 F3)" className="h-8 rounded-control border border-line bg-surface text-ink disabled:opacity-50">{label}</button>
);

// ---- Nếu ----
export type IfData = {
  kind: 'text' | 'image' | 'pixel' | 'value';
  text: string; waitMs: number; // thấy chữ / ảnh / màu
  left: string; op: 'eq' | 'ne' | 'contains' | 'gt' | 'lt' | 'empty'; right: string; // so giá trị
  negate: boolean;
};
const OPS: { id: IfData['op']; label: string }[] = [
  { id: 'eq', label: 'bằng' }, { id: 'ne', label: 'khác' }, { id: 'contains', label: 'chứa' },
  { id: 'gt', label: 'lớn hơn' }, { id: 'lt', label: 'nhỏ hơn' }, { id: 'empty', label: 'trống' },
];
export function IfFields({ value: v, onChange }: Ctl<IfData>) {
  return (
    <Grid>
      <Field label="Điều kiện" wide group>
        <Segmented value={v.kind} onChange={(kind) => onChange({ ...v, kind })} full label="Loại điều kiện"
          options={[{ id: 'text', label: 'Thấy chữ' }, { id: 'image', label: 'Thấy ảnh' }, { id: 'pixel', label: 'Màu' }, { id: 'value', label: 'Giá trị' }]} />
      </Field>
      {v.kind === 'value' ? (
        <>
          <Field label="Giá trị" wide><input value={v.left} onChange={(e) => onChange({ ...v, left: e.target.value })} className={`${textBox} font-mono text-[13px]`} placeholder="{{n2.output.text}}" /></Field>
          <Field label="Phép so"><select value={v.op} onChange={(e) => onChange({ ...v, op: e.target.value as IfData['op'] })} className={textBox}>{OPS.map((o) => <option key={o.id} value={o.id}>{o.label}</option>)}</select></Field>
          {v.op !== 'empty' && <Field label="Với"><input value={v.right} onChange={(e) => onChange({ ...v, right: e.target.value })} className={`${textBox} font-mono text-[13px]`} /></Field>}
        </>
      ) : (
        <>
          {v.kind === 'text' && <Field label="Chữ cần thấy" wide><input value={v.text} onChange={(e) => onChange({ ...v, text: e.target.value })} className={textBox} /></Field>}
          {v.kind !== 'text' && <p className="col-span-2 text-[13px] text-muted">{v.kind === 'image' ? 'Khoanh ảnh mẫu như bước Chờ ảnh.' : 'Chọn điểm + màu như bước Chờ màu pixel.'}</p>}
          <Field label="Chờ tối đa"><NumInput value={v.waitMs} unit="ms" onChange={(waitMs) => onChange({ ...v, waitMs })} /></Field>
          <span className="self-end pb-2 text-xs text-muted">0 = xem một lần</span>
        </>
      )}
      <label className="col-span-2 flex items-center justify-between text-muted">
        Đảo điều kiện (không thấy / không khớp)<Switch on={v.negate} onChange={(negate) => onChange({ ...v, negate })} label="Đảo điều kiện" />
      </label>
      <div className="col-span-2 grid"><TestBtn label="Thử điều kiện trên màn hình" /></div>
    </Grid>
  );
}

// ---- Switch theo màn hình: chờ, cái nào hiện trước thì đi cổng đó ----
export type ScreenCase = { key: string; label: string; by: 'text' | 'image'; text: string };
export type SwitchScreenData = { cases: ScreenCase[]; timeoutMs: number };
export function SwitchScreenFields({ value: v, onChange }: Ctl<SwitchScreenData>) {
  const set = (i: number, c: Partial<ScreenCase>) => onChange({ ...v, cases: v.cases.map((x, j) => (j === i ? { ...x, ...c } : x)) });
  return (
    <div className="flex flex-col gap-3">
      <p className="text-[13px] text-muted">Chụp màn hình, xem cái nào trong các case hiện trước thì đi ra cổng của case đó.</p>
      {v.cases.map((c, i) => (
        <div key={c.key} className="flex flex-col gap-2 rounded-card border border-line p-2.5">
          <div className="flex items-center gap-2">
            <span className="h-2 w-2 rounded-full bg-muted" />
            <input value={c.label} onChange={(e) => set(i, { label: e.target.value })} aria-label="Tên case" className="min-w-0 flex-1 bg-transparent font-semibold outline-none" />
            <button onClick={() => onChange({ ...v, cases: v.cases.filter((_, j) => j !== i) })} aria-label="Xoá case" className="text-muted hover:text-danger-text">×</button>
          </div>
          <Segmented value={c.by} onChange={(by) => set(i, { by })} full options={[{ id: 'text', label: 'Thấy chữ' }, { id: 'image', label: 'Thấy ảnh' }]} />
          {c.by === 'text' ? <input value={c.text} onChange={(e) => set(i, { text: e.target.value })} className={textBox} placeholder="Chữ cần thấy" />
            : <span className="text-xs text-muted">Khoanh ảnh mẫu như bước Chờ ảnh.</span>}
        </div>
      ))}
      <button onClick={() => onChange({ ...v, cases: [...v.cases, { key: `c${v.cases.length + 1}`, label: `Case ${v.cases.length + 1}`, by: 'text', text: '' }] })}
        className="h-8 rounded-control border border-dashed border-line-strong text-muted hover:text-ink">+ Thêm case</button>
      <Grid>
        <Field label="Chờ tối đa"><NumInput value={v.timeoutMs} unit="ms" onChange={(timeoutMs) => onChange({ ...v, timeoutMs })} /></Field>
        <span className="self-end pb-2 text-xs text-muted">hết giờ → cổng Không thấy</span>
      </Grid>
      <TestBtn label="Thử trên màn hình" />
    </div>
  );
}

// ---- Switch theo giá trị ----
export type ValueCase = { key: string; label: string; op: 'eq' | 'contains' | 'regex'; match: string };
export type SwitchData = { value: string; cases: ValueCase[] };
export function SwitchFields({ value: v, onChange }: Ctl<SwitchData>) {
  const set = (i: number, c: Partial<ValueCase>) => onChange({ ...v, cases: v.cases.map((x, j) => (j === i ? { ...x, ...c } : x)) });
  return (
    <div className="flex flex-col gap-3">
      <Field label="Giá trị cần xét"><input value={v.value} onChange={(e) => onChange({ ...v, value: e.target.value })} className={`${textBox} font-mono text-[13px]`} placeholder="{{n2.output.text}}" /></Field>
      {v.cases.map((c, i) => (
        <div key={c.key} className="grid grid-cols-[minmax(0,1fr)_96px] gap-2 rounded-card border border-line p-2.5">
          <input value={c.label} onChange={(e) => set(i, { label: e.target.value })} aria-label="Tên case" className="min-w-0 bg-transparent font-semibold outline-none" />
          <select value={c.op} onChange={(e) => set(i, { op: e.target.value as ValueCase['op'] })} className={textBox}>
            <option value="eq">bằng</option><option value="contains">chứa</option><option value="regex">khớp mẫu</option>
          </select>
          <input value={c.match} onChange={(e) => set(i, { match: e.target.value })} className={`${textBox} col-span-2 font-mono text-[13px]`} />
        </div>
      ))}
      <button onClick={() => onChange({ ...v, cases: [...v.cases, { key: `c${v.cases.length + 1}`, label: `Case ${v.cases.length + 1}`, op: 'eq', match: '' }] })}
        className="h-8 rounded-control border border-dashed border-line-strong text-muted hover:text-ink">+ Thêm case</button>
      <span className="text-xs text-muted">Xét lần lượt từ trên xuống; không case nào khớp → cổng Khác.</span>
    </div>
  );
}

// ---- Đếm ----
export type CounterData = { max: number };
export function CounterFields({ value: v, onChange }: Ctl<CounterData>) {
  return (
    <div className="flex flex-col gap-2">
      <Grid><Field label="Tối đa"><NumInput value={v.max} min={1} unit="lần" onChange={(max) => onChange({ max })} /></Field></Grid>
      <span className="text-xs text-muted">Mỗi lần đi qua cộng 1: còn dưới {v.max} → cổng <b>Còn lượt</b>, đủ {v.max} → <b>Hết lượt</b>. Đặt lại mỗi vòng kịch bản. Dùng để chặn vòng lặp vô hạn.</span>
    </div>
  );
}

// ---- Kết thúc ----
export type EndData = { status: 'success' | 'fail'; message: string };
export function EndFields({ value: v, onChange }: Ctl<EndData>) {
  return (
    <Grid>
      <Field label="Kết quả" wide group>
        <Segmented value={v.status} onChange={(status) => onChange({ ...v, status })} full options={[{ id: 'success', label: 'Thành công' }, { id: 'fail', label: 'Thất bại' }]} />
      </Field>
      <Field label="Lời nhắn" wide><input value={v.message} onChange={(e) => onChange({ ...v, message: e.target.value })} className={textBox} /></Field>
      <span className="col-span-2 text-xs text-muted">{v.status === 'fail' ? 'Dừng kịch bản và báo lỗi kèm lời nhắn (ghi vào nhật ký, thanh mini).' : 'Kết thúc vòng này bình thường – còn vòng lặp thì chạy vòng sau.'}</span>
    </Grid>
  );
}

// ---- Đặt biến ----
export type SetVarData = { name: string; op: 'set' | 'add'; value: string };
export function SetVarFields({ value: v, onChange }: Ctl<SetVarData>) {
  return (
    <Grid>
      <Field label="Tên biến"><input value={v.name} onChange={(e) => onChange({ ...v, name: e.target.value })} className={`${textBox} font-mono text-[13px]`} /></Field>
      <Field label="Cách đặt" group>
        <Segmented value={v.op} onChange={(op) => onChange({ ...v, op })} full options={[{ id: 'set', label: 'Gán' }, { id: 'add', label: 'Cộng' }]} />
      </Field>
      <Field label="Giá trị" wide><input value={v.value} onChange={(e) => onChange({ ...v, value: e.target.value })} className={`${textBox} font-mono text-[13px]`} /></Field>
      <span className="col-span-2 text-xs text-muted">Dùng ở node sau: <span className="var-token">{`vars.${v.name || 'ten'}`}</span></span>
    </Grid>
  );
}
