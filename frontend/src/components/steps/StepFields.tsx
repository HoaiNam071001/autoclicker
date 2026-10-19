// Phần inspector cho các loại bước chưa làm logic (bước 27b) – mỗi khối nhận value / onChange để bước làm logic
// nối thẳng vào store: cuộn, mở URL / app (23), chờ ảnh, click theo ảnh, chờ pixel (25), khi hết giờ (25b),
// focus cửa sổ (26), nhóm bước (20).
import { ReactNode } from 'react';
import NumInput, { fieldBox } from '../NumInput';
import Segmented from '../Segmented';
import { Target } from '../icons';
import { KeyChip } from '../inputs/PinButton';
import { Field as RowField } from '../stepform/shared';
import {
  ClickImageData, FocusWindowData, GroupData, OnTimeout, OpenData, SAMPLE_WINDOWS, ScrollData, WaitImageData, WaitPixelData,
} from './model';

type Ctl<T> = { value: T; onChange: (v: T) => void };
const textBox = `${fieldBox} w-full text-ink outline-none`;
const btn = 'h-8 flex-1 rounded-control border border-line bg-surface text-ink hover:bg-hover';

// Hàng trường dùng chung với inspector thật (stepform/shared – design v3); pinKey: chip key input đã ghim dưới nhãn.
export function Field({ label, children, wide = false, pinKey, group = false, help }: {
  label: string; children: ReactNode; wide?: boolean; pinKey?: string; group?: boolean; help?: string;
}) {
  return <RowField label={label} wide={wide} group={group} help={help} extra={pinKey && <KeyChip k={pinKey} />}>{children}</RowField>;
}

// Toạ độ X / Y + nút chọn trên màn hình
function PointInput({ value, onChange }: Ctl<{ x: number; y: number }>) {
  return (
    <div className="flex gap-2">
      {(['x', 'y'] as const).map((k) => (
        <NumInput key={k} prefix={k.toUpperCase()} value={value[k]} onChange={(v) => onChange({ ...value, [k]: v })} className="flex-1" />
      ))}
      <button title="Chọn trên màn hình" aria-label="Chọn trên màn hình"
        className="grid h-8 w-9 flex-none place-items-center rounded-control border border-line text-muted hover:bg-hover hover:text-ink"><Target /></button>
    </div>
  );
}

const Grid = ({ children }: { children: ReactNode }) => <div className="grid grid-cols-2 gap-3">{children}</div>;

// ---- Cuộn (bước 23) ----
export function ScrollFields({ value: v, onChange }: Ctl<ScrollData>) {
  return (
    <Grid>
      <Field label="Hướng" wide group>
        <Segmented value={v.direction} onChange={(direction) => onChange({ ...v, direction })} full label="Hướng cuộn"
          options={[{ id: 'up', label: '↑ Lên' }, { id: 'down', label: '↓ Xuống' }, { id: 'left', label: '← Trái' }, { id: 'right', label: '→ Phải' }]} />
      </Field>
      <Field label="Số nấc"><NumInput value={v.amount} min={1} unit="nấc" onChange={(amount) => onChange({ ...v, amount })} /></Field>
      <Field label="Cuộn tại" group>
        <Segmented value={v.at ? 'point' : 'mouse'} onChange={(m) => onChange({ ...v, at: m === 'point' ? v.at ?? { x: 0, y: 0 } : null })} full
          options={[{ id: 'mouse', label: 'Chuột' }, { id: 'point', label: 'Điểm' }]} />
      </Field>
      {v.at && <Field label="Vị trí" wide group><PointInput value={v.at} onChange={(at) => onChange({ ...v, at })} /></Field>}
    </Grid>
  );
}

// ---- Mở URL / app (bước 23) ----
export function OpenFields({ value: v, onChange, pinKey }: Ctl<OpenData> & { pinKey?: string }) {
  return (
    <Grid>
      <Field label="URL hoặc đường dẫn app" wide pinKey={pinKey}>
        <input value={v.target} onChange={(e) => onChange({ ...v, target: e.target.value })} className={`${textBox} font-mono text-[13px]`} />
        <span className="text-xs">Nhận biến, vd <span className="font-mono">facebook.com/{'{{profile}}'}</span></span>
      </Field>
      <Field label="Mở bằng" wide>
        <select value={v.app} onChange={(e) => onChange({ ...v, app: e.target.value })} className={textBox}>
          <option value="">Mặc định của máy</option><option>Google Chrome</option><option>Microsoft Edge</option><option>Firefox</option>
        </select>
      </Field>
    </Grid>
  );
}

// ---- Chờ ảnh (bước 25) – artboard Main ----
export function WaitImageFields({ value: v, onChange, timeoutPin }: Ctl<WaitImageData> & { timeoutPin?: string }) {
  return (
    <>
      <section className="flex flex-col gap-2">
        <span className="font-semibold">Ảnh mẫu</span>
        <TemplatePreview size={v.template} />
        <div className="flex gap-1.5">
          <button className={btn}>{v.template ? 'Chụp lại' : 'Khoanh ảnh mẫu'}</button>
          <button className={btn} disabled={!v.template}>Xem trên màn hình</button>
        </div>
        <label className="flex flex-col gap-1.5 text-muted">Độ khớp tối thiểu
          <span className="flex items-center gap-2.5">
            <input type="range" min={50} max={100} value={v.threshold} onChange={(e) => onChange({ ...v, threshold: +e.target.value })}
              className="flex-1 accent-accent" />
            <span className="w-10 text-right font-mono text-ink">{v.threshold}%</span>
          </span>
        </label>
      </section>
      <Grid>
        <Field label="Chờ tối đa" pinKey={timeoutPin}>
          <NumInput value={v.timeoutMs / 1000} min={1} unit="s" pinned={!!timeoutPin} onChange={(s) => onChange({ ...v, timeoutMs: s * 1000 })} />
        </Field>
        <Field label="Kiểm tra mỗi"><NumInput value={v.intervalMs} min={50} unit="ms" onChange={(intervalMs) => onChange({ ...v, intervalMs })} /></Field>
        <Field label="Tìm trong" wide group>
          <Segmented value={v.area} onChange={(area) => onChange({ ...v, area })} full
            options={[{ id: 'region', label: 'Vùng của kịch bản' }, { id: 'screen', label: 'Toàn màn hình' }]} />
        </Field>
      </Grid>
    </>
  );
}

// Khung ảnh mẫu (mẫu: vẽ nút "Nhắn tin" bằng CSS thay cho ảnh chụp thật)
export function TemplatePreview({ size }: { size: { w: number; h: number } | null }) {
  return (
    <div className="relative grid h-[92px] place-items-center overflow-hidden rounded-card border border-line bg-soft text-[13px] text-muted">
      {size ? (
        <>
          <span className="grid h-11 w-[132px] place-items-center rounded-lg bg-[#DCE6FF] font-semibold text-[#1E3FAF]">Nhắn tin</span>
          <span className="absolute right-1.5 bottom-1.5 rounded bg-hud/80 px-1.5 font-mono text-[11px] text-hud-text">{size.w} × {size.h}</span>
        </>
      ) : 'Chưa có ảnh mẫu – bấm “Khoanh ảnh mẫu”'}
    </div>
  );
}

// ---- Click theo ảnh (bước 25) ----
export function ClickImageFields({ value: v, onChange, prevStep }: Ctl<ClickImageData> & { prevStep: string }) {
  return (
    <Grid>
      <Field label="Vị trí ảnh" wide group>
        <Segmented value={v.source} onChange={(source) => onChange({ ...v, source })} full
          options={[{ id: 'prev', label: `Kết quả bước ${prevStep}` }, { id: 'search', label: 'Tự tìm ảnh mẫu' }]} />
        <span className="text-xs">{v.source === 'prev' ? `Click vào chỗ bước ${prevStep} (Chờ ảnh) tìm thấy` : 'Chụp ảnh mẫu riêng, tìm lúc chạy (dùng khi hết giờ như Chờ ảnh)'}</span>
      </Field>
      <Field label="Lệch so với tâm ảnh" wide group>
        <div className="flex gap-2">
          <NumInput signed prefix="ΔX" value={v.dx} onChange={(dx) => onChange({ ...v, dx })} className="flex-1" />
          <NumInput signed prefix="ΔY" value={v.dy} onChange={(dy) => onChange({ ...v, dy })} className="flex-1" />
        </div>
      </Field>
      <Field label="Nút chuột" wide group>
        <Segmented value={v.button} onChange={(button) => onChange({ ...v, button })} full
          options={[{ id: 'left', label: 'Trái' }, { id: 'right', label: 'Phải' }, { id: 'center', label: 'Giữa' }]} />
      </Field>
    </Grid>
  );
}

// ---- Chờ màu pixel (bước 25) ----
export function WaitPixelFields({ value: v, onChange }: Ctl<WaitPixelData>) {
  return (
    <Grid>
      <Field label="Điểm" wide group><PointInput value={v.at} onChange={(at) => onChange({ ...v, at })} /></Field>
      <Field label="Màu" wide group>
        <div className="flex gap-2">
          <span className={`${fieldBox} flex-1 gap-2`}>
            <span className="h-4 w-4 flex-none rounded border border-line" style={{ background: v.color }} />
            <input value={v.color} onChange={(e) => onChange({ ...v, color: e.target.value })} className="min-w-0 flex-1 bg-transparent font-mono text-[13px] text-ink outline-none" />
          </span>
          <button title="Lấy màu đang hiện tại điểm trên màn hình" className="h-8 rounded-control border border-line px-2.5 text-[13px] whitespace-nowrap text-ink hover:bg-hover">Lấy màu</button>
        </div>
      </Field>
      <Field label="Sai số màu"><NumInput prefix="±" value={v.tolerance} onChange={(tolerance) => onChange({ ...v, tolerance })} /></Field>
      <Field label="Chờ tối đa"><NumInput value={v.timeoutMs / 1000} min={1} unit="s" onChange={(s) => onChange({ ...v, timeoutMs: s * 1000 })} /></Field>
    </Grid>
  );
}

// ---- Focus cửa sổ (bước 26) ----
export function FocusWindowFields({ value: v, onChange }: Ctl<FocusWindowData>) {
  return (
    <Grid>
      <Field label="Tên cửa sổ / ứng dụng" wide>
        <input value={v.title} list="open-windows" onChange={(e) => onChange({ ...v, title: e.target.value })} className={textBox} />
        <datalist id="open-windows">{SAMPLE_WINDOWS.map((w) => <option key={w} value={w} />)}</datalist>
      </Field>
      <Field label="So khớp" wide group>
        <Segmented value={v.match} onChange={(match) => onChange({ ...v, match })} full
          options={[{ id: 'contains', label: 'Chứa chữ này' }, { id: 'exact', label: 'Đúng tên' }]} />
      </Field>
      <div className="col-span-2 flex flex-col gap-1 text-xs text-muted">
        <span>Đang mở:</span>
        <div className="flex flex-wrap gap-1">
          {SAMPLE_WINDOWS.map((w) => (
            <button key={w} onClick={() => onChange({ ...v, title: w })}
              className={`rounded-full border px-2 py-0.5 ${v.title === w ? 'border-accent bg-accent-soft text-accent-text' : 'border-line hover:bg-hover hover:text-ink'}`}>{w}</button>
          ))}
        </div>
      </div>
    </Grid>
  );
}

// ---- Khi hết giờ (bước 25b) – artboard Main ----
export function OnTimeoutFields({ value: v, onChange, inGroup }: Ctl<OnTimeout> & { inGroup: boolean }) {
  return (
    <section className="flex flex-col gap-2">
      <span className="font-semibold">Khi hết giờ</span>
      <Segmented value={v.action} onChange={(action) => onChange({ ...v, action })} full label="Khi hết giờ"
        options={[{ id: 'retry', label: 'Thử lại' }, { id: 'skip', label: 'Bỏ qua' }, { id: 'stop', label: 'Dừng' }]} />
      {v.action === 'retry' && (
        <>
          <label className="flex items-center justify-between text-muted">Số lần thử lại
            <NumInput value={v.retries} min={1} onChange={(retries) => onChange({ ...v, retries })} className="w-[72px]" />
          </label>
          <label className="flex items-center justify-between text-muted">Thử hết vẫn không thấy
            <select value={v.then} onChange={(e) => onChange({ ...v, then: e.target.value as OnTimeout['then'] })} className={`${fieldBox} w-[150px] text-ink`}>
              {inGroup && <option value="skipGroup">Bỏ qua nhóm</option>}
              <option value="stop">Dừng kịch bản</option>
            </select>
          </label>
        </>
      )}
      <span className="text-xs text-muted">
        {v.action === 'retry' ? `Thử thêm ${v.retries} lần (mỗi lần chờ lại từ đầu)`
          : v.action === 'skip' ? 'Không thấy thì bỏ qua bước này, chạy tiếp bước sau'
          : 'Không thấy thì dừng kịch bản, ghi lý do vào lịch sử chạy'}
      </span>
    </section>
  );
}

// ---- Nhóm bước (bước 20) ----
export function GroupFields({ value: v, onChange, count, repeatPin }: Ctl<GroupData> & { count: number; repeatPin?: string }) {
  return (
    <Grid>
      <Field label="Tên nhóm" wide>
        <input value={v.name} onChange={(e) => onChange({ ...v, name: e.target.value })} className={textBox} />
      </Field>
      <Field label="Lặp" pinKey={repeatPin}>
        <NumInput value={v.repeat} min={1} unit="lần" pinned={!!repeatPin} onChange={(repeat) => onChange({ ...v, repeat })} />
      </Field>
      <Field label="Bước trong nhóm"><span className="flex h-8 items-center font-mono text-ink">{count}</span></Field>
      <span className="col-span-2 text-xs text-muted">
        Cả nhóm chạy {v.repeat} lần liên tiếp rồi mới sang bước sau. Kéo thẻ bước vào / ra khung để thêm / bớt.
        Input trong nhóm dùng cùng giá trị ở mọi lần lặp.
      </span>
    </Grid>
  );
}
