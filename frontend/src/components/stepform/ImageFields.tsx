// Node chờ ảnh / chữ, click theo ảnh / chữ, chờ màu pixel (bước 25, 25e): ảnh mẫu / chữ, thử tìm, nguồn chỗ thấy, màu.
import { useEffect, useState } from 'react';
import NumInput, { fieldBox } from '../NumInput';
import Segmented from '../Segmented';
import { Switch } from '../RunSettings';
import VarTextInput from '../inputs/VarTextInput';
import AnchorField from '../steps/AnchorField';
import { setAnchor, updateStep, useScript } from '../../store';
import { findsImage } from '../../flow/model';
import { addInput } from '../../scriptInputs';
import { checkCapture, pickColor, requestCapture, startCapture, templateUrl, testFind, useImageState } from '../../imageSteps';
import { openOverlay } from '../../windowMode';
import { nodeTitle } from '../../flow/view';
import { Step } from '../../types';
import { ButtonField, Field, fmtSec, hint, NumField, PickButton, selectBox, SetPoint, useNodeVars } from './shared';
import PointSource from './PointSource';
import { usesVar } from '../../anchor';

export const isImageStep = (s: Step) => s.type === 'waitImage' || s.type === 'clickImage' || s.type === 'waitPixel';

export // macOS: chưa có quyền ghi màn hình thì chụp chỉ ra hình nền – báo ngay trong inspector.
function CapturePermission() {
  const { allowed } = useImageState();
  useEffect(() => {
    checkCapture();
  }, []);
  if (allowed !== false) return null;
  return (
    <div className="col-span-2 flex flex-col gap-1.5 rounded-lg bg-danger-soft px-3 py-2 text-xs text-danger-text">
      Chưa có quyền ghi màn hình (Screen Recording) – app không chụp được màn hình để tìm ảnh / đọc màu.
      <span>System Settings → Privacy &amp; Security → Screen Recording → bật cho app (khi dev: terminal / VS Code), rồi mở lại app.</span>
      <button onClick={requestCapture} className="self-start rounded-control border border-danger/40 px-2.5 py-1 font-semibold hover:bg-danger/10">Xin quyền</button>
    </div>
  );
}


export // Ảnh mẫu (hoặc chữ) + thử tìm + độ khớp / chờ tối đa / vùng tìm – dùng cho Chờ ảnh và Click theo ảnh (tự tìm).
function ImageSearchFields({ step: s }: { step: Step }) {
  const { finds, busy, error } = useImageState();
  const find = finds[s.id];
  const [size, setSize] = useState<{ w: number; h: number } | null>(null);
  const k = s.templateScale || 1;
  const byText = s.findBy === 'text';
  return (
    <>
      <Field label="Tìm theo" wide group>
        <Segmented value={byText ? 'text' : 'image'} full label="Tìm theo"
          onChange={(v) => updateStep(s.id, { findBy: v === 'text' ? 'text' : undefined })}
          options={[{ id: 'image', label: 'Ảnh mẫu' }, { id: 'text', label: 'Chữ' }]} />
        <span className={hint}>
          {byText ? 'Đọc chữ trên màn hình – không bị ảnh hưởng bởi theme sáng / tối, zoom, đổi màu khi rê chuột. Dùng khi nút / nhãn có chữ.'
            : 'So từng điểm ảnh – dùng cho icon, nút không có chữ. Nút có chữ nên chọn “Chữ”.'}
        </span>
      </Field>
      {byText ? <TextSearchFields step={s} /> : <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span>Ảnh mẫu</span>
        <div className="relative grid min-h-[92px] place-items-center overflow-hidden rounded-card border border-line bg-soft p-2 text-[13px]">
          {s.template ? (
            <>
              <img src={templateUrl(s.template)} alt="Ảnh mẫu" className="max-h-[140px] max-w-full object-contain"
                style={size ? { width: Math.min(size.w / k, 300) } : undefined}
                onLoad={(e) => setSize({ w: e.currentTarget.naturalWidth, h: e.currentTarget.naturalHeight })} />
              {size && (
                <span className="absolute right-1.5 bottom-1.5 rounded bg-hud/80 px-1.5 font-mono text-[11px] text-hud-text">
                  {Math.round(size.w / k)} × {Math.round(size.h / k)}
                </span>
              )}
            </>
          ) : 'Chưa có ảnh mẫu – bấm “Khoanh ảnh mẫu”, kéo khung quanh nút / chữ cần chờ'}
        </div>
        <div className="flex gap-1.5">
          <button onClick={() => startCapture(s.id)} className="h-8 flex-1 rounded-control border border-line bg-surface text-ink hover:bg-hover">
            {s.template ? 'Khoanh lại' : 'Khoanh ảnh mẫu'}
          </button>
          <button onClick={() => testFind(s.id)} disabled={!s.template || busy === s.id}
            title="Tìm ảnh mẫu trên màn hình hiện tại (app ẩn đi một chút) – để biết độ khớp thực tế"
            className="h-8 flex-1 rounded-control border border-line bg-surface text-ink enabled:hover:bg-hover disabled:opacity-50">
            {busy === s.id ? 'Đang tìm…' : 'Thử tìm'}
          </button>
        </div>
        {find && (
          <span className={`flex flex-wrap items-center gap-x-2 text-xs ${find.found ? 'text-ok' : 'text-warn'}`}>
            {find.found
              ? <>Thấy tại X {find.x + Math.round(find.w / 2)} Y {find.y + Math.round(find.h / 2)} · khớp {find.score}% · {find.ms} ms</>
              : find.fromRun
                ? <>Lần chạy trước hết giờ không thấy – chỗ giống nhất khớp {find.score}% (cần {s.threshold}%)</>
                : <>Không thấy trên màn hình hiện tại – khớp nhất {find.score}% (cần {s.threshold}%)</>}
            <button onClick={() => openOverlay('point', s.id)} className="text-muted underline">xem trên màn hình</button>
          </span>
        )}
        {error && busy === '' && <span className="text-xs text-danger-text">{error}</span>}
      </div>}
      {!byText && <NumField label="Độ khớp tối thiểu" step={s} field="threshold" unit="%" min={1} />}
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 250} min={20} unit="ms" onChange={(intervalMs) => updateStep(s.id, { intervalMs })} />
      </Field>
      <span className={`self-end pb-2 ${hint}`}>Không thấy sau {fmtSec(s.timeoutMs)} → cổng Thất bại</span>
      {!byText && (s.threshold ?? 0) < 70 && (
        <span className="col-span-2 -mt-1 text-xs text-warn">
          Độ khớp {s.threshold}% rất thấp – gần như chỗ nào cũng “khớp”, dễ click nhầm. Ảnh khoanh đúng thường khớp 95–100%;
          nên để 80–90% và bấm Thử tìm để xem độ khớp thật.
        </span>
      )}
      <Field label="Tìm trong" wide group>
        <Segmented value={s.searchArea ?? 'region'} onChange={(v) => updateStep(s.id, { searchArea: v === 'screen' ? 'screen' : undefined })} full
          options={[{ id: 'region', label: 'Vùng của kịch bản' }, { id: 'screen', label: 'Cả màn hình' }]} />
        <span className={hint}>
          {byText ? 'Vùng nhỏ đọc nhanh hơn (cả màn ~1 s mỗi lần) và ít nhầm chữ trùng ở chỗ khác.'
            : 'Vùng nhỏ tìm nhanh hơn và ít nhầm hơn. Ảnh mẫu chụp ở màn Retina vẫn tìm được trên màn thường.'}
          {' '}Cả màn hình = cả màn chứa vùng.
        </span>
      </Field>
    </>
  );
}


// Tìm theo chữ: ô chữ (nhận {{biến}}), cách so, chỗ thứ mấy, thử tìm.
function TextSearchFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const { finds, busy, error } = useImageState();
  const find = finds[s.id];
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const empty = !s.findText?.trim();
  return (
    <div className="col-span-2 flex flex-col gap-1.5 text-muted">
      <span className="flex items-center justify-between">Chữ cần tìm
        <span className="text-xs">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
      </span>
      <VarTextInput value={s.findText ?? ''} onChange={(findText) => set({ findText })} inputs={script.inputs ?? []} nodeVars={nodeVars}
        placeholder="Xem lịch sử" onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} />
      <span className={`text-xs ${empty ? 'text-danger-text' : ''}`}>
        {empty ? 'Chưa nhập chữ – chạy sẽ báo lỗi. ' : ''}Không phân biệt hoa / thường, bỏ qua khoảng trắng. Click vào giữa chữ.
      </span>
      <div className="flex flex-wrap items-center gap-x-4 gap-y-2 py-1 text-[13px] text-ink">
        <label className="flex items-center gap-2"><Switch on={!!s.textExact} onChange={(textExact) => set({ textExact: textExact || undefined })} label="Cả dòng" />Cả dòng đúng chữ này</label>
        <label className="flex items-center gap-2"><Switch on={!!s.textNoAccent} onChange={(textNoAccent) => set({ textNoAccent: textNoAccent || undefined })} label="Bỏ qua dấu" />Bỏ qua dấu</label>
        <span className="flex items-center gap-2">Chỗ thứ
          <NumInput value={Math.max(s.occurrence ?? 1, 1)} min={1} onChange={(n) => set({ occurrence: n > 1 ? n : undefined })} className="w-20" />
        </span>
      </div>
      <span className={hint}>
        “Cả dòng”: chỉ khớp ô / nút chỉ có đúng chữ này (tránh khớp trong câu dài). “Bỏ qua dấu”: khi máy đọc sai dấu.
        Chữ hiện nhiều lần → chỗ thứ mấy tính từ trên xuống, trái sang phải.
      </span>
      <button onClick={() => testFind(s.id)} disabled={empty || busy === s.id}
        title="Đọc chữ trên màn hình hiện tại và tìm (app ẩn đi một chút)"
        className="h-8 rounded-control border border-line bg-surface text-ink enabled:hover:bg-hover disabled:opacity-50">
        {busy === s.id ? 'Đang đọc chữ…' : 'Thử tìm'}
      </button>
      {find && (
        <span className={`flex flex-col gap-1 text-xs ${find.found ? 'text-ok' : 'text-warn'}`}>
          <span className="flex flex-wrap items-center gap-x-2">
            {find.found ? <>Thấy tại X {find.x + Math.round(find.w / 2)} Y {find.y + Math.round(find.h / 2)} · {find.ms} ms</>
              : find.fromRun ? <>Lần chạy trước hết giờ không thấy chữ này</>
              : <>Không thấy trên màn hình hiện tại{find.ms ? ` · ${find.ms} ms` : ''}</>}
            {find.found && <button onClick={() => openOverlay('point', s.id)} className="text-muted underline">xem trên màn hình</button>}
          </span>
          {!find.found && (
            <span className="text-muted">
              {find.read?.length ? <>Đọc được: {find.read.slice(0, 12).map((t, i) => (
                <button key={i} onClick={() => set({ findText: t })} title="Dùng chữ này"
                  className="mr-1 mb-1 inline-block rounded border border-line px-1.5 text-ink hover:bg-hover">{t}</button>
              ))}</> : 'Không đọc được chữ nào trong vùng tìm.'}
            </span>
          )}
        </span>
      )}
      {error && busy === '' && <span className="text-xs text-danger-text">{error}</span>}
    </div>
  );
}


// Click theo ảnh / chữ: chỗ một node chờ ảnh / chữ đã thấy (chọn node nguồn), hoặc tự tìm.
export function ClickImageFields({ step: s }: { step: Step }) {
  const { steps } = useScript();
  const sources = steps.filter((x) => x.id !== s.id && findsImage(x));
  const src = sources.find((x) => x.id === s.imageFrom);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const byNode = s.imageSource !== 'search';
  return (
    <>
      <Field label="Click vào" wide group>
        <Segmented value={byNode ? 'node' : 'search'} onChange={(imageSource: 'node' | 'search') => set({ imageSource })} full
          options={[{ id: 'node', label: 'Chỗ node khác thấy' }, { id: 'search', label: 'Tự tìm (ảnh / chữ)' }]} label="Click vào" />
        {byNode && (
          <select value={src ? src.id : ''} onChange={(e) => set({ imageFrom: e.target.value || undefined })}
            className={`${selectBox} ${src ? '' : 'border-danger'}`}>
            <option value="">— chọn node tìm ảnh / chữ —</option>
            {sources.map((x) => <option key={x.id} value={x.id}>{x.ref} · {nodeTitle(x)}</option>)}
          </select>
        )}
        <span className={`text-xs ${byNode && !src ? 'text-danger-text' : 'text-muted'}`}>
          {!byNode ? `Tự tìm ${s.findBy === 'text' ? 'chữ' : 'ảnh mẫu'} lúc chạy (chờ như ${s.findBy === 'text' ? 'Chờ chữ' : 'Chờ ảnh'}) rồi click vào giữa`
            : src ? `Click vào chỗ ${src.ref} thấy lần gần nhất trong vòng – ${src.ref} chưa thấy (chưa chạy / thất bại) thì node này thất bại`
            : sources.length ? 'Chọn node tìm ảnh / chữ để lấy chỗ thấy' : 'Chưa có node Chờ ảnh / Chờ chữ nào – thêm một node, hoặc chọn “Tự tìm”'}
        </span>
      </Field>
      {!byNode && <ImageSearchFields step={s} />}
      <Field label={`Lệch so với tâm ${!byNode && s.findBy === 'text' ? 'chữ' : 'chỗ thấy'}`} wide group>
        <div className="flex gap-2">
          <NumInput signed prefix="ΔX" value={s.dx ?? 0} onChange={(dx) => set({ dx })} className="flex-1" />
          <NumInput signed prefix="ΔY" value={s.dy ?? 0} onChange={(dy) => set({ dy })} className="flex-1" />
        </div>
        <span className={hint}>Vd tìm nhãn “Họ tên”, ô nhập nằm bên phải 150 px → ΔX 150.</span>
      </Field>
      <ButtonField step={s} />
    </>
  );
}

export function WaitPixelFields({ step: s, setPoint }: { step: Step; setPoint: SetPoint }) {
  const { region } = useScript();
  const { busy, error } = useImageState();
  const valid = /^#[0-9a-fA-F]{6}$/.test(s.color ?? '');
  return (
    <>
      {s.points[0] && (
        <PointSource step={s} index={0} setPoint={setPoint} label="Điểm">
          {(menu) => <AnchorField value={s.points[0]} onChange={(p) => setPoint(0, p)} onAnchor={(ax, ay) => setAnchor(s.id, 0, ax, ay)}
            pick={<PickButton step={s} />} menu={menu} />}
        </PointSource>
      )}
      <Field label="Màu cần chờ" wide group>
        <div className="flex gap-2">
          <span className={`${fieldBox} flex-1 gap-2 ${valid ? '' : 'border-danger'}`}>
            <span className="h-4 w-4 flex-none rounded border border-line" style={{ background: valid ? s.color : 'transparent' }} />
            <input value={s.color ?? ''} onChange={(e) => updateStep(s.id, { color: e.target.value.trim() })} spellCheck={false}
              className="min-w-0 flex-1 bg-transparent font-mono text-[13px] text-ink outline-none" />
          </span>
          <button onClick={() => pickColor(s.id)} disabled={busy === s.id || (!!s.points[0] && usesVar(s.points[0]))}
            title="Lấy màu đang hiện tại điểm này trên màn hình (app ẩn đi một chút) – mở sẵn màn hình có màu cần chờ"
            className="h-8 rounded-control border border-line px-2.5 text-[13px] whitespace-nowrap text-ink enabled:hover:bg-hover disabled:opacity-50">
            {busy === s.id ? 'Đang lấy…' : 'Lấy màu'}
          </button>
        </div>
        {!valid && <span className="text-xs text-danger-text">Màu dạng #RRGGBB, vd #1877F2</span>}
        {error && busy === '' && <span className="text-xs text-danger-text">{error}</span>}
      </Field>
      <Field label="Sai số màu">
        <NumInput prefix="±" value={s.tolerance ?? 0} onChange={(tolerance) => updateStep(s.id, { tolerance: Math.min(tolerance, 255) })} />
      </Field>
      <NumField label="Chờ tối đa" step={s} field="timeoutMs" unit="ms" min={1} />
      <span className={`col-span-2 ${hint}`}>
        Mỗi {s.intervalMs ?? 100} ms đọc màu tại điểm{region ? ' (trong vùng)' : ''}, mỗi kênh R / G / B lệch ≤ {s.tolerance} là được.
        Không đổi màu sau {fmtSec(s.timeoutMs)} → cổng Thất bại. Không di chuột.
      </span>
    </>
  );
}
