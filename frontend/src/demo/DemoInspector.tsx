// Inspector của màn #demo: tab Bước (nội dung theo loại bước mẫu) / Input (InputsPanel).
// Mọi thay đổi chỉ là state cục bộ của màn mẫu.
import { ReactNode, useState } from 'react';
import { StepIcon } from '../components/icons';
import { Switch } from '../components/RunSettings';
import NumInput from '../components/NumInput';
import PinButton, { KeyChip } from '../components/inputs/PinButton';
import VarTextInput from '../components/inputs/VarTextInput';
import InputsPanel from '../components/inputs/InputsPanel';
import { ScriptInput } from '../components/inputs/model';
import { DemoGroup, DemoStep } from './demoScript';
import {
  ClickImageFields, FocusWindowFields, GroupFields, OnTimeoutFields, OpenFields, ScrollFields, WaitImageFields, WaitPixelFields,
} from '../components/steps/StepFields';
import AnchorField from '../components/steps/AnchorField';
import {
  AnchoredPoint, ClickImageData, FocusWindowData, OnTimeout, OpenData, ScrollData, WaitImageData, WaitPixelData,
} from '../components/steps/model';

type Props = {
  step: DemoStep | DemoGroup;
  n: string;
  inputs: ScriptInput[];
  setInputs: (inputs: ScriptInput[]) => void;
};

export default function DemoInspector({ step, n, inputs, setInputs }: Props) {
  const [tab, setTab] = useState<'step' | 'inputs'>('step');
  const [enabled, setEnabled] = useState(step.type === 'group' || step.enabled !== false);

  return (
    <aside aria-label={`Thuộc tính bước ${n} (mẫu)`} className="flex w-[340px] flex-none flex-col gap-4 overflow-y-auto border-l border-line bg-surface p-4 max-[1100px]:w-[300px]">
      <div className="flex items-center gap-3">
        <span className="grid h-9 w-9 flex-none place-items-center rounded-[9px] bg-accent-soft text-accent-text"><StepIcon type={step.type} className="h-[18px] w-[18px]" /></span>
        <div className="min-w-0 flex-1">
          <div className="truncate text-[15px] font-semibold">{step.type === 'group' ? `Nhóm: ${step.name}` : step.title}</div>
          <div className="text-xs text-muted">
            {step.type === 'group' ? `${step.children.length} bước, lặp ×${step.repeat}` : `Bước ${n}${inGroup(step) ? ', trong nhóm Gửi tin' : ''}`}
          </div>
        </div>
        <Switch on={enabled} onChange={setEnabled} label="Bật bước" />
      </div>

      <div className="-mt-1 flex gap-[18px] border-b border-line" role="tablist">
        <TabBtn active={tab === 'step'} onClick={() => setTab('step')}>Bước</TabBtn>
        <TabBtn active={tab === 'inputs'} onClick={() => setTab('inputs')}>
          Input <span className="rounded-full bg-accent-soft px-1.5 text-[11px] font-semibold text-accent-text">{inputs.length}</span>
        </TabBtn>
      </div>

      {tab === 'inputs'
        ? <InputsPanel inputs={inputs} onChange={setInputs} onAdd={() => setInputs([...inputs, newInput(inputs)])} />
        : <StepTab step={step} inputs={inputs} setInputs={setInputs} />}
    </aside>
  );
}

const inGroup = (s: DemoStep) => ['d3', 'd4', 'd5', 'd6'].includes(s.id);

const newInput = (inputs: ScriptInput[]): ScriptInput =>
  ({ key: `bien_${inputs.length + 1}`, label: 'Input mới', type: 'text', usedBy: [] });

function TabBtn({ active, onClick, children }: { active: boolean; onClick: () => void; children: ReactNode }) {
  return (
    <button role="tab" aria-selected={active} onClick={onClick}
      className={`relative flex h-[34px] items-center gap-1.5 px-0.5 ${active ? 'font-semibold' : 'text-muted hover:text-ink'}`}>
      {children}
      {active && <span className="absolute inset-x-[15%] -bottom-px h-[3px] rounded-sm bg-accent" />}
    </button>
  );
}

// Nội dung tab Bước theo loại (giá trị mẫu, sửa được trong màn #demo)
function StepTab({ step, inputs, setInputs }: Omit<Props, 'n'>) {
  if (step.type === 'group') return <DemoGroupTab group={step} />;
  return (
    <>
      <StepBody step={step} inputs={inputs} setInputs={setInputs} />
      <Footer step={step} />
    </>
  );
}

function StepBody({ step, inputs, setInputs }: { step: DemoStep } & Omit<Props, 'step' | 'n'>) {
  const [scroll, setScroll] = useState<ScrollData>({ direction: 'down', amount: 5, at: { x: 760, y: 520 } });
  const [open, setOpen] = useState<OpenData>({ target: 'facebook.com/{{profile}}', app: 'Google Chrome' });
  const [wait, setWait] = useState<WaitImageData>({ template: { w: 164, h: 64 }, threshold: 90, timeoutMs: 10000, intervalMs: 250, area: 'region' });
  const [click, setClick] = useState<ClickImageData>({ source: 'prev', dx: 0, dy: 0, button: 'left' });
  const [pixel, setPixel] = useState<WaitPixelData>({ at: { x: 640, y: 300 }, color: '#1877F2', tolerance: 10, timeoutMs: 5000 });
  const [win, setWin] = useState<FocusWindowData>({ title: 'Google Chrome', match: 'contains' });
  const [timeout, setTimeoutCfg] = useState<OnTimeout>({ action: 'retry', retries: 2, then: 'skipGroup' });
  const [pt, setPt] = useState<AnchoredPoint>({ x: 300, y: 300, ax: 'end', ay: 'start' });
  const onTimeout = <OnTimeoutFields value={timeout} onChange={setTimeoutCfg} inGroup={inGroup(step)} />;

  switch (step.type) {
    case 'type': return <TypeStep inputs={inputs} setInputs={setInputs} />;
    case 'scroll': return <ScrollFields value={scroll} onChange={setScroll} />;
    case 'open': return <OpenFields value={open} onChange={setOpen} pinKey="profile" />;
    case 'waitImage': return <><WaitImageFields value={wait} onChange={setWait} timeoutPin="wait_s" />{onTimeout}</>;
    case 'clickImage': return <ClickImageFields value={click} onChange={setClick} prevStep="03" />;
    case 'waitPixel': return <><WaitPixelFields value={pixel} onChange={setPixel} />{onTimeout}</>;
    case 'focusWindow': return <><FocusWindowFields value={win} onChange={setWin} />{onTimeout}</>;
    case 'click': return (
      <div className="flex flex-col gap-1.5 text-muted">
        <span>Vị trí <span className="text-xs">(neo theo góc vùng)</span></span>
        <AnchorField value={pt} onChange={setPt} />
        <span className="mt-1 text-xs">Vùng đổi cỡ thì điểm vẫn cách cạnh phải và cạnh trên đúng khoảng này. Gõ số âm để neo cạnh phải / dưới.</span>
      </div>
    );
    default: return <p className="text-[13px] text-muted">Loại bước này dùng inspector thật (xem màn chính).</p>;
  }
}

// Cuối inspector: chờ sau bước + điểm dừng (debug, bước 25c)
function Footer({ step }: { step: DemoStep }) {
  const [delay, setDelay] = useState(parseInt(step.delay) * (step.delay.includes(' s') ? 1000 : 1) || 300);
  const [bp, setBp] = useState(!!step.breakpoint);
  return (
    <section className="mt-auto flex flex-col gap-3 border-t border-line pt-3">
      <div className="grid grid-cols-2 gap-3">
        <label className="flex flex-col gap-1.5 text-muted">Chờ sau bước<NumInput value={delay} unit="ms" onChange={setDelay} /></label>
        <label className="flex flex-col gap-1.5 text-muted">Ngẫu nhiên<NumInput prefix="±" value={step.delay.includes('±') ? 50 : 0} unit="ms" onChange={() => {}} /></label>
      </div>
      <label className="flex items-center justify-between">
        <span className="flex items-center gap-2"><span className="h-2 w-2 rounded-full bg-danger" />Điểm dừng khi chạy từng bước</span>
        <Switch on={bp} onChange={setBp} label="Điểm dừng" />
      </label>
    </section>
  );
}

function DemoGroupTab({ group }: { group: DemoGroup }) {
  const [g, setG] = useState({ name: group.name, repeat: group.repeat });
  return <GroupFields value={g} onChange={setG} count={group.children.length} repeatPin="send_times" />;
}

function TypeStep({ inputs, setInputs }: Omit<Props, 'step' | 'n'>) {
  const [text, setText] = useState('xin chào {{M1}} và {{M2}}, ');
  const [vals, setVals] = useState({ charMs: 30, jitter: 10, delay: 200 });
  const [pins, setPins] = useState<Record<string, string>>({}); // trường → key input
  const fields: { id: keyof typeof vals; label: string; unit: string; prefix?: string; key: string }[] = [
    { id: 'charMs', label: 'Cách nhau giữa các ký tự', unit: 'ms', key: 'char_ms' },
    { id: 'jitter', label: 'Ngẫu nhiên thời gian', unit: 'ms', prefix: '±', key: 'jitter_ms' },
    { id: 'delay', label: 'Chờ trước bước tiếp', unit: 'ms', key: 'type_delay' },
  ];
  const togglePin = (f: (typeof fields)[number]) => {
    if (pins[f.id]) {
      setPins(({ [f.id]: _, ...rest }) => rest);
      setInputs(inputs.filter((i) => i.key !== pins[f.id]));
    } else {
      setPins({ ...pins, [f.id]: f.key });
      setInputs([...inputs, { key: f.key, label: f.label, type: 'ms', default: vals[f.id], unit: 'ms', usedBy: [{ n: '05', title: 'Gõ chữ' }] }]);
    }
  };

  return (
    <>
      <section className="flex flex-col gap-2">
        <div className="flex items-center justify-between">
          <span className="font-semibold">Nội dung</span>
          <span className="text-xs text-muted">gõ <span className="font-mono">{'{{'}</span> để chèn biến</span>
        </div>
        <VarTextInput value={text} onChange={setText} inputs={inputs}
          onCreateInput={(key) => setInputs([...inputs, { key, label: key, type: 'text', usedBy: [{ n: '05', title: 'Gõ chữ' }] }])} />
      </section>
      <section className="flex flex-col gap-3">
        {fields.map((f) => (
          <div key={f.id} className="flex flex-col gap-1.5 text-muted">
            <span className="flex items-center gap-2">{f.label}{pins[f.id] && <KeyChip k={pins[f.id]} />}</span>
            <div className="flex gap-1.5">
              <NumInput value={vals[f.id]} unit={f.unit} prefix={f.prefix} pinned={!!pins[f.id]} className="flex-1"
                onChange={(v) => setVals({ ...vals, [f.id]: v })} />
              <PinButton pinned={!!pins[f.id]} onToggle={() => togglePin(f)} />
            </div>
          </div>
        ))}
        <span className="text-xs text-muted">Bấm ghim để cho phép đổi giá trị khi chạy. Giá trị hiện tại thành mặc định.</span>
      </section>
    </>
  );
}
