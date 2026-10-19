// Thanh công cụ giữa phía trên overlay (chỉ hiển thị; đáy bị Dock / taskbar che) – tự ẩn, hiện khi chuột sát mép trên.
// Ba dạng: xem riêng một node, xem trước, các công cụ đặt điểm / vùng.
import { ReactNode } from 'react';
import { main } from '../../../wailsjs/go/models';
import { PageKindPicker } from '../MockPage';
import { Play, Region, SelectIcon, StepIcon, Stop, Target } from '../icons';
import { Step, stepLabel } from '../../types';
import type { OverlayTool } from '../../windowMode';
import type { BranchChoices, BranchNode } from '../../flow/preview';
import { BranchMenu } from './BranchMenu';

export const TOOLS: { id: OverlayTool; label: string; soon?: string; icon: ReactNode }[] = [
  { id: 'point', label: 'Đặt điểm', icon: <StepIcon type="click" /> },
  { id: 'drag', label: 'Kéo thả A → B', icon: <StepIcon type="drag" /> },
  { id: 'region', label: 'Chọn vùng', icon: <Region /> },
  { id: 'select', label: 'Chọn nhiều', icon: <SelectIcon /> },
  { id: 'capture', label: 'Khoanh ảnh', icon: <StepIcon type="waitImage" /> },
];
const SPEEDS = [1, 2, 4];
const toolBtn = 'flex h-11 items-center gap-2 rounded-[7px] px-3.5 font-semibold';
const on = (active: boolean) => (active ? 'bg-white/14 text-white' : 'hover:bg-white/6');
const Sep = () => <span className="mx-1.5 h-[26px] w-px bg-white/16" />;

export type ToolbarProps = {
  open: boolean; onEnter: () => void; onLeave: () => void; onDone: () => void;
  // xem riêng một node
  focusStep: Step | null; repicking: boolean; onToggleRepick: () => void; onShowAll: () => void;
  // xem trước
  preview: boolean; simRunning: boolean; onToggleSim: () => void; speed: number; setSpeed: (v: number) => void; onEditPoints: () => void;
  branches: BranchNode[]; choices: BranchChoices; onChooseBranch: (id: string, port: string) => void; onResetBranches: () => void;
  // công cụ; fullRegion: nút đặt vùng = cả màn (null = đang là cả màn)
  tool: OverlayTool; setTool: (t: OverlayTool) => void; fullRegion: { label: string; onClick: () => void } | null;
  displays: main.DisplayInfo[]; onPickDisplay: (i: number) => void;
};

export function OverlayToolbar(p: ToolbarProps) {
  return (
    <div data-ui onMouseEnter={p.onEnter} onMouseLeave={p.onLeave}
      className={`absolute left-1/2 top-2 flex -translate-x-1/2 transition-[translate,opacity] duration-150 ${p.open ? '' : 'pointer-events-none -translate-y-[calc(100%+16px)] opacity-0'} items-center gap-1 whitespace-nowrap rounded-xl border border-hud-line bg-hud p-1.5 text-hud-text shadow-[0_12px_32px_rgb(0_0_0/.4)]`}>
      {p.focusStep ? <FocusTools {...p} focusStep={p.focusStep} /> : p.preview ? <PreviewTools {...p} /> : <MainTools {...p} />}
      {p.displays.length > 1 && <><Sep /><DisplayPicker displays={p.displays} onPick={p.onPickDisplay} /></>}
      <button onClick={p.onDone} className={`${toolBtn} btn-primary ml-1 font-bold`}>Xong <span className="font-mono text-[11px] font-normal opacity-70">Esc</span></button>
    </div>
  );
}

function FocusTools({ focusStep, repicking, onToggleRepick, onShowAll }: ToolbarProps & { focusStep: Step }) {
  return (
    <>
      <span className="flex items-center gap-2 px-2.5 font-semibold">
        <span className="grid h-7 w-7 place-items-center rounded-full bg-accent text-[13px] text-on-accent">{focusStep.ref}</span>
        {focusStep.ref} · {stepLabel(focusStep)}
      </span>
      <button onClick={onToggleRepick} className={`${toolBtn} ${on(repicking)}`}>
        <Target />{repicking ? 'Huỷ chọn lại' : 'Chọn lại trên màn hình'}
      </button>
      <Sep />
      <button onClick={onShowAll} className={`${toolBtn} font-normal hover:bg-white/6`}>Xem tất cả</button>
    </>
  );
}

function PreviewTools({ simRunning, onToggleSim, speed, setSpeed, onEditPoints, branches, choices, onChooseBranch, onResetBranches }: ToolbarProps) {
  return (
    <>
      <button onClick={onToggleSim} className={`${toolBtn} ${simRunning ? 'btn-danger-soft' : 'btn-primary'}`}>
        {simRunning ? <><Stop className="h-2.5 w-2.5" />Dừng</> : <><Play className="h-3 w-3" />Chạy thử</>}
      </button>
      <div className="flex rounded-lg bg-white/5 p-[3px]" role="radiogroup" aria-label="Tốc độ xem trước">
        {SPEEDS.map((v) => (
          <button key={v} role="radio" aria-checked={speed === v} onClick={() => setSpeed(v)}
            className={`h-[38px] rounded-md px-3 font-mono ${on(speed === v)}`}>
            {v}×
          </button>
        ))}
      </div>
      {branches.length > 0 && <BranchMenu branches={branches} choices={choices} onChoose={onChooseBranch} onReset={onResetBranches} className={`${toolBtn} font-normal`} />}
      <PageKindPicker className={`${toolBtn} font-normal`} />
      <Sep />
      <button onClick={onEditPoints} title="Chọn nhiều điểm để di chuyển / co giãn / xoay, rồi Xem trước lại" className={`${toolBtn} font-normal hover:bg-white/6`}>Sửa điểm</button>
    </>
  );
}

function MainTools({ tool, setTool, fullRegion }: ToolbarProps) {
  return (
    <>
      {TOOLS.map((t) => (
        <button key={t.id} disabled={!!t.soon} title={t.soon && `Sắp có (bước ${t.soon})`} onClick={() => setTool(t.id)}
          className={`${toolBtn} disabled:opacity-40 ${on(tool === t.id)}`}>
          {t.icon}{t.label}
        </button>
      ))}
      {tool === 'region' && fullRegion && (
        <button onClick={fullRegion.onClick} className={`${toolBtn} font-normal hover:bg-white/6`}>{fullRegion.label}</button>
      )}
      <Sep />
      <span className="flex items-center gap-1.5 px-1.5 text-[13px] text-hud-text/60">
        <kbd className="rounded border border-hud-line px-1.5 font-mono text-hud-text">F6</kbd>lưu vị trí chuột
      </span>
      <button onClick={() => setTool('preview')} className={`${toolBtn} font-normal hover:bg-white/6`}>Xem trước</button>
    </>
  );
}

function DisplayPicker({ displays, onPick }: { displays: main.DisplayInfo[]; onPick: (i: number) => void }) {
  return (
    <div className="flex rounded-lg bg-white/5 p-[3px]" role="radiogroup" aria-label="Màn hình">
      {displays.map((d) => (
        <button key={d.index} role="radio" aria-checked={d.current} onClick={() => onPick(d.index)}
          title={`Màn ${d.index + 1}${d.main ? ' (chính)' : ''} · ${d.rect.w} × ${d.rect.h} tại X ${d.rect.x} Y ${d.rect.y}`}
          className={`flex h-[38px] items-center gap-1.5 rounded-md px-2.5 ${on(d.current)}`}>
          <svg width="14" height="12" viewBox="0 0 14 12" aria-hidden><rect x="1" y="1" width="12" height="8" rx="1.5" fill="none" stroke="currentColor" strokeWidth="1.4" /><path d="M5 11h4" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" /></svg>
          <span className="font-mono text-[13px]">{d.index + 1}</span>
          {d.main && <span className="text-[11px] text-hud-text/60">chính</span>}
        </button>
      ))}
    </div>
  );
}
