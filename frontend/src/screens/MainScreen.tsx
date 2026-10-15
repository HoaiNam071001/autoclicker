// Phương án A – danh sách bước + form cấu hình
import { useEffect, useState } from 'react';
import { Drag, Eye, Logo, Play, Region, Stop, Target } from '../components/icons';
import { MousePos, StopScript } from '../../wailsjs/go/main/App';
import { Quit, WindowMinimise } from '../../wailsjs/runtime/runtime';
import { inWails, openOverlay } from '../windowMode';
import { removeStep, runCurrent, setRegion, useEditor } from '../scriptStore';
import StepList from '../components/StepList';
import StepForm from '../components/StepForm';
import RunSettings from '../components/RunSettings';
import { STOP_REASON, useEngine } from '../engineStore';
import { Point } from '../types';

const btn = 'flex h-11 items-center gap-2 rounded-lg border border-line bg-white px-3.5';
const MOUSE_POLL_MS = 250;

export default function MainScreen() {
  const { script: { steps, region }, selected, notice } = useEditor();
  const cur = steps.find((s) => s.id === selected) ?? steps[0]; // undefined khi chưa có bước nào
  const idx = steps.indexOf(cur) + 1;
  const { running, lastError, stopReason, hotkeyError } = useEngine();

  // Toạ độ chuột thật ở thanh dưới – để đọc nhanh toạ độ khi nhập tay.
  const [mouse, setMouse] = useState<Point | null>(null);
  useEffect(() => {
    if (!inWails) return;
    const id = setInterval(() => MousePos().then(setMouse), MOUSE_POLL_MS);
    return () => clearInterval(id);
  }, []);

  // Delete / Backspace xoá bước đang chọn (trừ khi đang gõ trong ô nhập).
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const el = e.target as HTMLElement;
      if ((e.key !== 'Delete' && e.key !== 'Backspace') || el.closest('input, textarea, select')) return;
      if (cur) removeStep(cur.id);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [cur]);

  return (
    <div className="flex min-h-screen flex-col bg-ground">
      <header className="drag flex flex-wrap items-center gap-3 border-b border-line bg-white px-5 py-3">
        <div className="mr-3 flex items-center gap-2.5 text-accent">
          <Logo /><span className="text-base font-bold text-ink">Auto Clicker</span>
        </div>
        <div className="flex h-9 items-center gap-1.5 rounded-lg border border-line px-3">
          <span className="text-muted">Kịch bản</span><span className="font-semibold">Kịch bản 1</span>
        </div>
        <div className="flex flex-1 flex-wrap gap-2">
          <button onClick={() => openOverlay('region')} className={btn}><Region />Chọn vùng</button>
          <button onClick={() => openOverlay('point')} className={btn}><Target />Thêm điểm
            <kbd className="rounded border border-line px-1.5 font-mono text-xs text-muted">F6</kbd>
          </button>
          <button onClick={() => openOverlay('drag')} className={btn}><Drag className="h-4 w-4 text-drag" />Kéo thả</button>
        </div>
        <button onClick={() => openOverlay('preview')} className={btn} title="Xem đường đi các bước bằng con trỏ giả, không click thật">
          <Eye />Xem trước
        </button>
        <button onClick={() => (running ? StopScript() : runCurrent())} disabled={!inWails}
          title={inWails ? undefined : 'Chỉ chạy được trong app'}
          className={`flex h-11 items-center gap-2.5 rounded-lg px-5 font-semibold text-white disabled:opacity-50 ${running ? 'bg-danger' : 'bg-accent'}`}>
          {running ? <><Stop />Dừng</> : <><Play />Chạy</>} <span className="font-mono text-xs opacity-85">F8</span>
        </button>
        {inWails && (
          <div className="flex gap-1">
            <button aria-label="Thu nhỏ" onClick={WindowMinimise} className="grid h-8 w-8 place-items-center rounded-md text-muted hover:bg-ground">–</button>
            <button aria-label="Đóng" onClick={Quit} className="grid h-8 w-8 place-items-center rounded-md text-muted hover:bg-danger hover:text-white">✕</button>
          </div>
        )}
      </header>

      <div className="flex flex-1 flex-wrap gap-4 px-5 py-4">
        <StepList />

        {!cur ? (
          <section className="grid flex-[999_1_560px] place-items-center rounded-[10px] border border-line bg-white text-muted">
            Chọn hoặc thêm một bước để cấu hình
          </section>
        ) : (<>
        <StepForm step={cur} index={idx}>
          <RunSettings />
        </StepForm>
        </>)}
      </div>

      <footer className="flex flex-wrap gap-5 border-t border-line bg-white px-5 py-2 text-[13px] text-muted">
        <span>Vùng: {region
          ? <><span className="font-mono">X {region.x} Y {region.y} · {region.w} × {region.h}</span>
              <button onClick={() => setRegion(null)} className="ml-2 underline">Toàn màn hình</button></>
          : 'toàn màn hình'}</span>
        <span>Chuột: <span className="font-mono">{mouse ? `X ${mouse.x} Y ${mouse.y}` : '—'}</span></span>
        {hotkeyError && <span className="text-danger">{hotkeyError}</span>}
        {lastError && <span className="text-danger">Lỗi lần chạy trước: {lastError}</span>}
        {stopReason && <span>{STOP_REASON[stopReason] ?? stopReason}</span>}
        {notice && <span className="text-accent">{notice}</span>}
        <span className="ml-auto">Dừng: F8 hoặc rê chuột vào góc màn hình</span>
      </footer>
    </div>
  );
}
