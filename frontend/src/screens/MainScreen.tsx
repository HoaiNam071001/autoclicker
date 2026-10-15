// Phương án A – danh sách bước + form cấu hình
import { useEffect, useState } from 'react';
import { Drag, Logo, Play, Region, Target } from '../components/icons';
import { Modifier, makeStep } from '../types';
import { RunScript, SpikeClick, StopScript } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { EventsOn, Quit, WindowMinimise } from '../../wailsjs/runtime/runtime';
import { inWails, openOverlay, setWindowMode } from '../windowMode';
import { removeStep, setRegion, useEditor } from '../scriptStore';
import StepList from '../components/StepList';
import StepForm from '../components/StepForm';
import { setLastError, STOP_REASON, useEngine } from '../engineStore';

const btn = 'flex h-11 items-center gap-2 rounded-lg border border-line bg-white px-3.5';
const field = 'h-10 rounded-lg border border-line px-2.5 font-mono';

// Thử engine (bước 7): mở TextEdit / Notepad sao cho phủ điểm (500, 500) rồi bấm "Chạy thử engine".
const cmdKey: Modifier = navigator.platform.startsWith('Mac') ? 'cmd' : 'ctrl';
const engineDemo = main.Script.createFrom({
  name: 'Demo engine', region: null, loops: 3, maxMinutes: 0,
  steps: [
    makeStep({ id: 'd1', type: 'click', points: [{ x: 500, y: 500 }] }),
    makeStep({ id: 'd2', type: 'type', text: 'Xin chào – gõ tự động 123' }),
    makeStep({ id: 'd3', type: 'key', keys: ['enter'] }),
    makeStep({ id: 'd4', type: 'type', text: 'dòng thứ hai', delayAfter: 800 }),
    makeStep({ id: 'd5', type: 'double', points: [{ x: 500, y: 500 }], delayAfter: 800 }),
    makeStep({ id: 'd6', type: 'drag', points: [{ x: 420, y: 500 }, { x: 650, y: 500 }], dragMs: 600, delayAfter: 800 }),
    makeStep({ id: 'd7', type: 'key', keys: ['a'], modifiers: [cmdKey] }),
  ],
});

export default function MainScreen() {
  const { script: { steps, region }, selected, notice } = useEditor();
  const cur = steps.find((s) => s.id === selected) ?? steps[0]; // undefined khi chưa có bước nào
  const idx = steps.indexOf(cur) + 1;
  const [spikeMsg, setSpikeMsg] = useState('');
  const { running, lastError, stopReason } = useEngine();

  useEffect(() => {
    if (!inWails) return;
    const offs = [
      EventsOn('hotkey:error', (msg: string) => setSpikeMsg(`Lỗi: ${msg}`)),
    ];
    return () => offs.forEach((off) => off());
  }, []);

  // Spike bước 3: đếm ngược 3 giây (Go chờ) rồi click (500, 500).
  const runSpike = async () => {
    setSpikeMsg('Click sau 3 giây…');
    try {
      const pos = await SpikeClick(500, 500, 3000);
      setSpikeMsg(`Đã click tại ${pos.x}, ${pos.y}`);
    } catch (e) {
      setSpikeMsg(`Lỗi: ${e}`);
    }
  };

  // Thử engine + mini bar (bước 7–8): thu cửa sổ thành mini bar khi chạy, về màn chính khi xong / dừng.
  const toggleEngine = async () => {
    if (running) return StopScript();
    setSpikeMsg('');
    await setWindowMode('minibar');
    try {
      await RunScript(engineDemo);
    } catch (e) {
      setLastError(String(e));
    }
    await setWindowMode('main');
  };

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
        <button onClick={runSpike} className={btn} title="Spike bước 3">Thử click 500,500</button>
        <button onClick={toggleEngine} className={btn} title="Thử engine bước 7">{running ? 'Dừng engine' : 'Chạy thử engine'}</button>
        <button onClick={() => setWindowMode('minibar')} className={btn} title="Spike bước 5">Thử mini bar</button>
        {spikeMsg && <span className="font-mono text-xs text-muted">{spikeMsg}</span>}
        {lastError && <span className="font-mono text-xs text-danger">Lỗi engine: {lastError}</span>}
        {notice && <span className="font-mono text-xs text-accent">{notice}</span>}
        {stopReason && <span className="font-mono text-xs text-muted">{STOP_REASON[stopReason] ?? stopReason}</span>}
        <button className="flex h-11 items-center gap-2.5 rounded-lg bg-accent px-5 font-semibold text-white">
          <Play />Chạy <span className="font-mono text-xs opacity-85">F8</span>
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
          {/* Lặp kịch bản – chưa nối (bước 16) */}
          <div className="mt-auto flex flex-wrap items-center gap-5 border-t border-line pt-4">
            <span className="font-semibold">Lặp kịch bản</span>
            <div className="flex overflow-hidden rounded-lg border border-line">
              {['Số vòng', 'Vô hạn', 'Theo thời gian'].map((m, i) => (
                <button key={m} className={`h-10 px-3.5 ${i === 0 ? 'bg-ink text-white' : 'border-l border-line bg-white'}`}>{m}</button>
              ))}
            </div>
            <input aria-label="Số vòng" defaultValue={10} className={`${field} w-[72px]`} />
            <label className="flex items-center gap-2 text-[#3F4857]">
              <input type="checkbox" defaultChecked className="h-[18px] w-[18px] accent-accent" />Đếm ngược 3 giây
            </label>
          </div>
        </StepForm>
        </>)}
      </div>

      <footer className="flex flex-wrap gap-5 border-t border-line bg-white px-5 py-2 text-[13px] text-muted">
        <span>Vùng: {region
          ? <><span className="font-mono">X {region.x} Y {region.y} · {region.w} × {region.h}</span>
              <button onClick={() => setRegion(null)} className="ml-2 underline">Toàn màn hình</button></>
          : 'toàn màn hình'}</span>
        <span>Chuột: <span className="font-mono">X 1204 Y 388</span></span>
        <span className="ml-auto">Dừng khẩn cấp: rê chuột vào góc màn hình</span>
      </footer>
    </div>
  );
}
