// Inspector – cấu hình node đang chọn (design v2, cột phải), sửa trực tiếp vào store.
// Khung: đầu (loại, tên, mã node, bật / tắt), tab Node / Input; trường theo loại nằm ở components/stepform/.
import { ReactNode, useState } from 'react';
import { Eye, Info, StepIcon } from './icons';
import { Switch } from './RunSettings';
import { openOverlay } from '../windowMode';
import { changeStepType, toggleStep, updateStep, useScript } from '../store';
import { isControl, STEP_LABEL, Step, stepLabel, StepPoint, StepType } from '../types';
import { hint, NumField, selectBox, TabBtn } from './stepform/shared';
import { isMouse, JitterField, MouseFields } from './stepform/MouseFields';
import { ScrollFields } from './stepform/ScrollFields';
import { KeyFields } from './stepform/KeyFields';
import { OpenFields, TypeFields } from './stepform/TextFields';
import { CapturePermission, ClickImageFields, ImageSearchFields, isImageStep, WaitPixelFields } from './stepform/ImageFields';
import { WebFields } from './stepform/WebFields';
import { ShotFields } from './stepform/ShotFields';
import { ClipFields } from './stepform/ClipFields';
import { HttpFields } from './stepform/HttpFields';
import { FileFields } from './stepform/FileFields';
import { isWebStep } from '../flow/web';
import FlowSections from './stepform/FlowSections';
import ControlSection from './stepform/ControlSection';
import InputsTab from './stepform/InputsTab';
import { HumanNodeField } from './HumanFields';
import { usesHuman } from '../flow/human';
import { nodeInfo } from '../flow/nodeInfo';
import { NodeInfoCard, useNodeInfoShown } from './stepform/NodeInfoCard';

// Loại node hành động đổi qua lại được trong ô Hành động (node điều khiển có trường riêng, không đổi loại)
const TYPES: StepType[] = ['click', 'double', 'hold', 'drag', 'scroll', 'wait', 'type', 'key', 'open', 'waitImage', 'clickImage', 'waitPixel', 'screenshot', 'clipboard', 'http', 'file', 'webWait', 'webClick',
  'webType', 'webRead', 'webScroll'];
const REF_RE = /^[A-Za-z][A-Za-z0-9_]{0,31}$/; // khớp validRef (flow.go)

export default function StepForm({ step: s }: { step: Step }) {
  const script = useScript();
  const [tab, setTab] = useState<'step' | 'inputs'>('step');
  const [infoShown, toggleInfo] = useNodeInfoShown();
  const setPoint = (i: number, p: Partial<StepPoint>) => updateStep(s.id, { points: s.points.map((q, j) => (j === i ? { ...q, ...p } : q)) });
  const atPoint = isMouse(s) || s.type === 'waitPixel' || (s.type === 'scroll' && s.points.length > 0); // có vị trí trên màn hình

  return (
    <aside aria-label={`Thuộc tính node ${s.ref}`} className="flex w-[340px] flex-none flex-col gap-4 overflow-y-auto border-l border-line bg-surface p-4 max-[1100px]:w-[300px]">
      <div className="flex items-center gap-3">
        <span className="grid h-9 w-9 flex-none place-items-center rounded-[9px] bg-accent-soft text-accent-text"><StepIcon type={s.type} className="h-[18px] w-[18px]" /></span>
        <div className="min-w-0 flex-1">
          <input value={s.name ?? ''} placeholder={stepLabel(s)} onChange={(e) => updateStep(s.id, { name: e.target.value || undefined })}
            aria-label="Tên node" title="Tên hiện trên thẻ – bỏ trống = tên theo loại"
            className="w-full truncate bg-transparent text-[15px] font-semibold outline-none placeholder:text-ink" />
          <RefInput step={s} />
        </div>
        <button onClick={toggleInfo} aria-label="Mô tả loại node" aria-pressed={infoShown} title={infoShown ? 'Ẩn mô tả loại node' : 'Hiện mô tả loại node'}
          className={`grid h-8 w-8 flex-none place-items-center rounded-control border border-line hover:bg-hover hover:text-ink ${infoShown ? 'bg-soft text-ink' : 'text-muted'}`}><Info /></button>
        {atPoint && (
          <button onClick={() => openOverlay('point', s.id)} aria-label="Xem trên màn hình"
            title="Xem / sửa vị trí riêng node này trên màn hình – các node khác mờ đi"
            className="grid h-8 w-8 flex-none place-items-center rounded-control border border-line text-muted hover:bg-hover hover:text-ink"><Eye /></button>
        )}
        <Switch on={s.enabled} onChange={() => toggleStep(s.id)} label="Bật node" />
      </div>

      <div className="-mt-1 flex gap-[18px] border-b border-line" role="tablist">
        <TabBtn active={tab === 'step'} onClick={() => setTab('step')}>Node</TabBtn>
        <TabBtn active={tab === 'inputs'} onClick={() => setTab('inputs')} title="Các giá trị hỏi khi bấm Chạy">
          Input {!!script.inputs?.length && <span className="rounded-full bg-accent-soft px-1.5 text-[11px] font-semibold text-accent-text">{script.inputs.length}</span>}
        </TabBtn>
      </div>

      {/* mô tả loại node ngay dưới tab, trên các trường */}
      {infoShown && tab === 'step' && <NodeInfoCard info={nodeInfo(s)} onHide={toggleInfo} />}

      {tab === 'inputs' ? <InputsTab /> : isControl(s.type) ? <>
        <ControlSection step={s} />
        <FlowSections step={s} />
      </> : <>
        <label className="flex flex-col gap-1.5 text-muted">Hành động
          <select value={s.type} onChange={(e) => changeStepType(s.id, e.target.value as StepType)} className={selectBox}>
            {TYPES.map((t) => <option key={t} value={t}>{t === s.type ? stepLabel(s) : STEP_LABEL[t]}</option>)}
          </select>
        </label>

        <div className="grid grid-cols-2 gap-3">
          <TypeBody step={s} setPoint={setPoint} />
        </div>

        <FlowSections step={s} />

        {/* Chờ sau node + ngẫu nhiên thời gian (bước 19) ở cuối, giống artboard */}
        <section className="mt-auto grid grid-cols-2 gap-3 border-t border-line pt-3">
          <NumField label={s.type === 'wait' ? 'Thời gian chờ' : 'Chờ sau khi xong'} step={s} field="delayAfter" unit="ms" />
          <NumField label="Ngẫu nhiên" step={s} field="jitterMs" unit="ms" prefix="±" />
          <span className={`col-span-2 -mt-1 ${hint}`}>
            {s.jitterMs > 0
              ? `Chờ ${Math.max(s.delayAfter - s.jitterMs, 0)}–${s.delayAfter + s.jitterMs} ms${s.type === 'click' && s.clicks > 1 ? ', giữa các click cũng lệch' : ''}`
              : 'Thời gian chờ cố định (ngẫu nhiên = 0)'}
            {' '}· chỉ chờ khi node thành công
          </span>
        </section>
      </>}
    </aside>
  );
}

// Trường riêng theo loại node (trong lưới 2 cột).
function TypeBody({ step: s, setPoint }: { step: Step; setPoint: (i: number, p: Partial<StepPoint>) => void }): ReactNode {
  const script = useScript();
  return (
    <>
      {isMouse(s) && <MouseFields step={s} setPoint={setPoint} />}
      {s.type === 'scroll' && <ScrollFields step={s} setPoint={setPoint} />}
      {(isImageStep(s) || s.type === 'screenshot') && <CapturePermission />}
      {s.type === 'waitImage' && <ImageSearchFields step={s} />}
      {s.type === 'clickImage' && <ClickImageFields step={s} />}
      {s.type === 'waitPixel' && <WaitPixelFields step={s} setPoint={setPoint} />}
      {s.type === 'screenshot' && <ShotFields step={s} />}
      {s.type === 'clipboard' && <ClipFields step={s} />}
      {s.type === 'http' && <HttpFields step={s} />}
      {s.type === 'file' && <FileFields step={s} />}
      {isWebStep(s.type) && <WebFields step={s} />}
      {s.type === 'open' && <OpenFields step={s} />}
      {s.type === 'type' && <TypeFields step={s} />}
      {s.type === 'key' && <KeyFields step={s} />}
      {isMouse(s) && <JitterField step={s} />}
      {usesHuman(s) && <HumanNodeField value={s.human} scriptLevel={script.human} className="col-span-2"
        onChange={(human) => updateStep(s.id, { human })} />}
    </>
  );
}

// Mã node (n2): dùng trong biến {{n2.output.x}}. Giữ bản nháp khi gõ, chỉ ghi khi hợp lệ và không trùng.
function RefInput({ step: s }: { step: Step }) {
  const { steps } = useScript();
  const [draft, setDraft] = useState(s.ref);
  const dup = steps.some((x) => x.id !== s.id && x.ref === draft);
  const bad = !REF_RE.test(draft) || draft === 'vars' || draft === 'last';
  return (
    <div className="flex items-center gap-1.5 text-xs text-muted">
      <input value={draft} aria-label="Mã node" spellCheck={false}
        title="Mã node – dùng trong biến {{mã.output.x}}. Chữ cái đầu, sau đó chữ / số / _"
        onChange={(e) => {
          const v = e.target.value.trim();
          setDraft(v);
          if (REF_RE.test(v) && v !== 'vars' && v !== 'last' && !steps.some((x) => x.id !== s.id && x.ref === v)) updateStep(s.id, { ref: v });
        }}
        onBlur={() => setDraft(s.ref)}
        className={`w-14 rounded bg-soft px-1 font-mono outline-none ${dup || bad ? 'text-danger-text' : 'text-ink'}`} />
      {dup ? 'trùng mã node khác' : bad ? 'mã không hợp lệ' : <>{STEP_LABEL[s.type]}{s.enabled ? '' : ' · đang tắt'}</>}
    </div>
  );
}
