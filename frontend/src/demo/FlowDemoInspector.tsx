// Inspector của màn #demo-flow (bước 20 F4a): đầu (icon, mã node, tên, bật / tắt), tab Bước / Input.
// Tab Bước: nội dung theo loại node → Khi thất bại (node hành động) → Cổng ra → Output → Lần chạy gần nhất → chờ sau bước.
// Mọi thay đổi chỉ là state của màn mẫu.
import { ReactNode, useState } from 'react';
import { StepIcon } from '../components/icons';
import { Switch } from '../components/RunSettings';
import NumInput from '../components/NumInput';
import InputsPanel from '../components/inputs/InputsPanel';
import { InputView } from '../components/inputs/model';
import { FlowEdgeView, FlowNodeView, RunRow } from '../components/flow/model';
import { LastRunPanel, OutputsPanel, PortsPanel, RetryPanel } from '../components/flow/PortField';
import {
  CounterFields, EndFields, IfFields, SetVarFields, SwitchFields, SwitchScreenFields,
} from '../components/flow/nodes/ControlFields';
import { Footer, StepBody } from './DemoInspector';
import { DemoFlowNode, NodeBody, outputsOf } from './flowDemo';

type Props = {
  node: DemoFlowNode; view: FlowNodeView; views: FlowNodeView[]; edges: FlowEdgeView[];
  start: boolean; lastRun?: RunRow;
  inputs: InputView[]; setInputs: (v: InputView[]) => void;
  onChange: (patch: Partial<DemoFlowNode>) => void;
  onConnect: (port: string, to: string) => void; onDisconnect: (port: string) => void;
  onJump: (id: string) => void; onSetStart: () => void; onDelete: () => void;
};

const KIND_LABEL: Record<NodeBody['kind'], string> = {
  step: 'Hành động', wait: 'Hành động', if: 'Rẽ nhánh', switch: 'Rẽ nhánh', switchScreen: 'Rẽ nhánh', counter: 'Luồng', end: 'Luồng', setVar: 'Luồng',
};

export default function FlowDemoInspector(p: Props) {
  const { node, view, inputs, setInputs } = p;
  const [tab, setTab] = useState<'step' | 'inputs'>('step');
  const action = node.body.kind === 'step' || node.body.kind === 'wait';
  return (
    <aside aria-label={`Thuộc tính node ${node.ref} (mẫu)`} className="flex w-[340px] flex-none flex-col gap-4 overflow-y-auto border-l border-line bg-surface p-4 max-[1100px]:w-[300px]">
      <div className="flex items-center gap-3">
        <span className="grid h-9 w-9 flex-none place-items-center rounded-[9px] bg-accent-soft text-accent-text"><StepIcon type={node.icon} className="h-[18px] w-[18px]" /></span>
        <div className="min-w-0 flex-1">
          <input value={node.title} onChange={(e) => p.onChange({ title: e.target.value })} aria-label="Tên node"
            className="w-full truncate bg-transparent text-[15px] font-semibold outline-none" />
          <div className="flex items-center gap-1.5 text-xs text-muted">
            <input value={node.ref} onChange={(e) => p.onChange({ ref: e.target.value.replace(/[^a-z0-9_]/gi, '') })} aria-label="Mã node"
              title="Mã node – dùng trong biến {{n2.output.x}}" className="w-12 rounded bg-soft px-1 font-mono outline-none" />
            {KIND_LABEL[node.body.kind]}
            {p.start && <span className="rounded-full bg-accent-soft px-1.5 font-semibold text-accent-text">Bắt đầu</span>}
          </div>
        </div>
        <Switch on={node.enabled !== false} onChange={(on) => p.onChange({ enabled: on })} label="Bật node" />
      </div>

      <div className="-mt-1 flex gap-[18px] border-b border-line" role="tablist">
        <TabBtn active={tab === 'step'} onClick={() => setTab('step')}>Node</TabBtn>
        <TabBtn active={tab === 'inputs'} onClick={() => setTab('inputs')}>
          Input <span className="rounded-full bg-accent-soft px-1.5 text-[11px] font-semibold text-accent-text">{inputs.length}</span>
        </TabBtn>
      </div>

      {tab === 'inputs' ? (
        <InputsPanel inputs={inputs} onAdd={() => {}} onMove={() => {}}
          onUpdate={(key, x) => setInputs(inputs.map((i) => (i.key === key ? { ...i, ...x } : i)))}
          onRemove={(key) => setInputs(inputs.filter((i) => i.key !== key))} />
      ) : (
        <>
          {view.warn && <div className="rounded-control border border-warn/40 bg-warn/10 px-2.5 py-2 text-[13px] text-ink"><b className="text-warn">!</b> {view.warn}</div>}
          <Body node={node} inputs={inputs} setInputs={setInputs} onChange={p.onChange} />
          {action && <RetryPanel value={node.retry ?? { count: 0, delayMs: 300 }} onChange={(retry) => p.onChange({ retry, chips: retry.count ? [`thử lại ${retry.count}`] : [] })} />}
          <PortsPanel node={view} nodes={p.views} edges={p.edges} onConnect={p.onConnect} onDisconnect={p.onDisconnect} onJump={p.onJump} />
          <OutputsPanel refName={node.ref} fields={outputsOf(node.body)} />
          <LastRunPanel row={p.lastRun} />
          {action && node.body.kind === 'step' && <Footer step={node.body.step} />}
          <div className="flex gap-2 border-t border-line pt-3">
            <button disabled={p.start} onClick={p.onSetStart} className="h-8 flex-1 rounded-control border border-line bg-surface enabled:hover:bg-hover disabled:opacity-50">
              {p.start ? 'Đang là node bắt đầu' : 'Đặt làm bắt đầu'}
            </button>
            <button onClick={p.onDelete} className="h-8 rounded-control border border-line bg-surface px-3 text-danger-text hover:bg-danger-soft">Xoá node</button>
          </div>
        </>
      )}
    </aside>
  );
}

function Body({ node, inputs, setInputs, onChange }: { node: DemoFlowNode; inputs: InputView[]; setInputs: (v: InputView[]) => void; onChange: Props['onChange'] }) {
  const b = node.body;
  const set = (body: NodeBody) => onChange({ body });
  switch (b.kind) {
    case 'step': return <StepBody step={b.step} inputs={inputs} setInputs={setInputs} flow />;
    case 'wait': return <label className="flex flex-col gap-1.5 text-muted">Chờ<NumInput value={b.ms} unit="ms" onChange={(ms) => onChange({ body: { kind: 'wait', ms }, detail: `${ms} ms` })} /></label>;
    case 'if': return <IfFields value={b.data} onChange={(data) => set({ kind: 'if', data })} />;
    case 'switch': return <SwitchFields value={b.data} onChange={(data) => set({ kind: 'switch', data })} />;
    case 'switchScreen': return <SwitchScreenFields value={b.data} onChange={(data) => onChange({ body: { kind: 'switchScreen', data }, detail: `chờ tối đa ${data.timeoutMs / 1000} s · ${data.cases.length} case` })} />;
    case 'counter': return <CounterFields value={b.data} onChange={(data) => onChange({ body: { kind: 'counter', data }, detail: `tối đa ${data.max} lần` })} />;
    case 'end': return <EndFields value={b.data} onChange={(data) => set({ kind: 'end', data })} />;
    case 'setVar': return <SetVarFields value={b.data} onChange={(data) => set({ kind: 'setVar', data })} />;
  }
}

function TabBtn({ active, onClick, children }: { active: boolean; onClick: () => void; children: ReactNode }) {
  return (
    <button role="tab" aria-selected={active} onClick={onClick}
      className={`relative flex h-[34px] items-center gap-1.5 px-0.5 ${active ? 'font-semibold' : 'text-muted hover:text-ink'}`}>
      {children}
      {active && <span className="absolute inset-x-[15%] -bottom-px h-[3px] rounded-sm bg-accent" />}
    </button>
  );
}
