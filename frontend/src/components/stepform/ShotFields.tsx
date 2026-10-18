// Node Chụp màn hình (bước 35 N1): vùng chụp, thư mục / tên file (nhận {{biến}}), chép vào clipboard.
import { useEffect, useState } from 'react';
import { Displays } from '../../../wailsjs/go/main/App';
import NumInput from '../NumInput';
import { Region } from '../icons';
import { Switch } from '../RunSettings';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { startCapture } from '../../imageSteps';
import { inWails } from '../../windowMode';
import { nodeTitle } from '../../flow/view';
import { DEFAULT_SHOT_DIR, hasFoundBox, SHOT_AREAS, type ShotArea } from '../../flow/shot';
import { Rect, Step } from '../../types';
import { Field, hint, selectBox, useNodeVars } from './shared';

type Display = { index: number; rect: Rect; main: boolean };

export function ShotFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const area = s.shotArea ?? '';
  const varProps = { inputs: script.inputs ?? [], nodeVars, rows: 1, compact: true, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  return (
    <>
      <Field label="Chụp" wide group>
        <select value={area} className={selectBox} aria-label="Vùng chụp"
          onChange={(e) => set({ shotArea: (e.target.value || undefined) as ShotArea | undefined })}>
          {SHOT_AREAS.map((a) => <option key={a.id} value={a.id}>{a.label}</option>)}
        </select>
        {area === '' && <span className={hint}>{script.region
          ? `Vùng của kịch bản: ${script.region.w} × ${script.region.h} tại ${script.region.x}, ${script.region.y}`
          : 'Kịch bản chưa đặt vùng – chụp cả màn chính'}</span>}
        {area === 'screen' && <DisplayPicker value={s.display ?? 0} onChange={(display) => set({ display: display || undefined })} />}
        {area === 'rect' && <RectPicker step={s} onChange={(shotRect) => set({ shotRect })} />}
        {area === 'node' && <NodePicker step={s} onChange={(imageFrom) => set({ imageFrom })} />}
      </Field>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">Lưu vào thư mục
        <VarTextInput {...varProps} value={s.saveDir ?? ''} onChange={(v) => set({ saveDir: v || undefined })} placeholder={DEFAULT_SHOT_DIR} />
        <span className={hint}>Trống = <span className="font-mono">~/{DEFAULT_SHOT_DIR}</span> · chưa có thì tự tạo</span>
      </div>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">Tên file
        <VarTextInput {...varProps} value={s.fileName ?? ''} onChange={(v) => set({ fileName: v || undefined })} placeholder="autoclicker-<ngày>-<giờ>" />
        <span className={hint}>Tự thêm <span className="font-mono">.png</span> · trùng tên thì thêm (2), không ghi đè</span>
      </div>
      <div className="col-span-2 flex flex-col gap-1">
        <label className="flex items-center gap-2 text-ink">
          <Switch on={!!s.toClipboard} onChange={(on) => set({ toClipboard: on || undefined })} label="Chép ảnh vào clipboard" />
          Chép ảnh vào clipboard
        </label>
        <span className={`pl-11 ${hint}`}>Dán (Cmd/Ctrl+V) thẳng vào chat, email, Paint…</span>
      </div>
    </>
  );
}

function DisplayPicker({ value, onChange }: { value: number; onChange: (n: number) => void }) {
  const [ds, setDs] = useState<Display[]>([]);
  useEffect(() => { if (inWails) Displays().then(setDs).catch(() => {}); }, []);
  const known = value === 0 || value <= ds.length;
  return (
    <>
      <select value={value} onChange={(e) => onChange(Number(e.target.value))} className={selectBox} aria-label="Màn hình">
        <option value={0}>Màn chứa vùng kịch bản</option>
        {ds.map((d) => <option key={d.index} value={d.index + 1}>Màn {d.index + 1}{d.main ? ' (chính)' : ''} · {d.rect.w} × {d.rect.h}</option>)}
        {!known && <option value={value}>Màn {value} (không thấy trên máy này)</option>}
      </select>
      {!known && <span className="text-xs text-danger-text">Máy này không có màn {value} – chạy sẽ báo lỗi</span>}
    </>
  );
}

function RectPicker({ step: s, onChange }: { step: Step; onChange: (r: Rect) => void }) {
  const r = s.shotRect ?? { x: 0, y: 0, w: 0, h: 0 };
  const num = (k: keyof Rect, prefix: string, min?: number) => (
    <NumInput prefix={prefix} value={r[k]} min={min} signed={min === undefined} onChange={(v) => onChange({ ...r, [k]: v })} className="flex-1" />
  );
  return (
    <>
      <div className="flex gap-2">{num('x', 'X')}{num('y', 'Y')}</div>
      <div className="flex gap-2">{num('w', 'W', 1)}{num('h', 'H', 1)}</div>
      <button onClick={() => startCapture(s.id)} className="flex h-8 items-center justify-center gap-2 rounded-control border border-line text-ink hover:bg-hover">
        <Region className="h-3.5 w-3.5" />Khoanh trên màn hình…
      </button>
      <span className={`text-xs ${r.w > 0 && r.h > 0 ? 'text-muted' : 'text-danger-text'}`}>
        {r.w > 0 && r.h > 0 ? 'Toạ độ màn hình – không đổi theo vùng của kịch bản' : 'Chưa chọn vùng – chạy sẽ báo lỗi'}
      </span>
    </>
  );
}

function NodePicker({ step: s, onChange }: { step: Step; onChange: (id: string | undefined) => void }) {
  const { steps } = useScript();
  const sources = steps.filter((x) => x.id !== s.id && hasFoundBox(x));
  const src = sources.find((x) => x.id === s.imageFrom);
  return (
    <>
      <select value={src ? src.id : ''} onChange={(e) => onChange(e.target.value || undefined)} aria-label="Node lấy khung chụp"
        className={`${selectBox} ${src ? '' : 'border-danger'}`}>
        <option value="">— chọn node có chỗ thấy —</option>
        {sources.map((x) => <option key={x.id} value={x.id}>{x.ref} · {nodeTitle(x)}</option>)}
      </select>
      <span className={`text-xs ${src ? 'text-muted' : 'text-danger-text'}`}>
        {src ? `Chụp đúng khung ${src.ref} thấy lần gần nhất trong vòng – ${src.ref} chưa thấy thì node này thất bại`
          : sources.length ? 'Chọn node để lấy khung chụp' : 'Chưa có node Chờ ảnh / chữ hay phần tử web nào'}
      </span>
    </>
  );
}
