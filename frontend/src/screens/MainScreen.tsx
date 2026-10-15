// Phương án A – danh sách bước + form cấu hình
import { useEffect, useState } from 'react';
import { Drag, Logo, Play, Region, Target } from '../components/icons';
import { Modifier, Step, STEP_LABEL, StepType, fmtMs, fmtPos, makeStep, sampleSteps } from '../types';
import { RunScript, SpikeClick, StopScript } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { EventsOn, Quit, WindowMinimise } from '../../wailsjs/runtime/runtime';
import { inWails, setWindowMode } from '../windowMode';
import { setLastError, STOP_REASON, useEngine } from '../engineStore';

const btn = 'flex h-11 items-center gap-2 rounded-lg border border-line bg-white px-3.5';
const field = 'h-10 rounded-lg border border-line px-2.5 font-mono';
const types: StepType[] = ['click', 'double', 'hold', 'drag', 'wait', 'type', 'key'];
const modifiers: Modifier[] = ['cmd', 'ctrl', 'alt', 'shift'];

// Thử engine (bước 7): mở TextEdit / Notepad sao cho phủ điểm (500, 500) rồi bấm "Chạy thử engine".
const cmdKey: Modifier = navigator.platform.startsWith('Mac') ? 'cmd' : 'ctrl';
const engineDemo = main.Script.createFrom({
  name: 'Demo engine', region: null, loops: 3, maxMinutes: 0,
  steps: [
    makeStep({ id: 'd1', type: 'click', points: [{ x: 500, y: 500 }] }),
    makeStep({ id: 'd2', type: 'type', text: 'Xin chào – gõ tự động 123' }),
    makeStep({ id: 'd3', type: 'key', key: 'enter' }),
    makeStep({ id: 'd4', type: 'type', text: 'dòng thứ hai', delayAfter: 800 }),
    makeStep({ id: 'd5', type: 'double', points: [{ x: 500, y: 500 }], delayAfter: 800 }),
    makeStep({ id: 'd6', type: 'drag', points: [{ x: 420, y: 500 }, { x: 650, y: 500 }], dragMs: 600, delayAfter: 800 }),
    makeStep({ id: 'd7', type: 'key', key: 'a', modifiers: [cmdKey] }),
  ],
});

export default function MainScreen() {
  const [steps] = useState<Step[]>(sampleSteps);
  const [selected, setSelected] = useState('3');
  const cur = steps.find((s) => s.id === selected)!;
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
          <button className={btn}><Region />Chọn vùng</button>
          <button className={btn}><Target />Thêm điểm
            <kbd className="rounded border border-line px-1.5 font-mono text-xs text-muted">F6</kbd>
          </button>
          <button className={btn}><Drag className="h-4 w-4 text-drag" />Kéo thả</button>
        </div>
        <button onClick={runSpike} className={btn} title="Spike bước 3">Thử click 500,500</button>
        <button onClick={toggleEngine} className={btn} title="Thử engine bước 7">{running ? 'Dừng engine' : 'Chạy thử engine'}</button>
        <button onClick={() => setWindowMode('overlay')} className={btn} title="Spike bước 5">Thử overlay</button>
        <button onClick={() => setWindowMode('minibar')} className={btn} title="Spike bước 5">Thử mini bar</button>
        {spikeMsg && <span className="font-mono text-xs text-muted">{spikeMsg}</span>}
        {lastError && <span className="font-mono text-xs text-danger">Lỗi engine: {lastError}</span>}
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
        <section aria-label="Danh sách bước" className="flex max-w-[420px] flex-[1_1_340px] flex-col rounded-[10px] border border-line bg-white">
          <div className="flex items-center justify-between border-b border-line px-4 py-3.5">
            <h2 className="text-[15px] font-semibold">Các bước</h2>
            <span className="text-muted">{steps.length} bước</span>
          </div>
          <ol className="flex flex-col gap-0.5 p-1.5">
            {steps.map((s, i) => {
              const ring = s.type === 'drag' ? 'border-drag text-drag' : 'border-accent text-accent';
              return (
                <li key={s.id}>
                  <button
                    onClick={() => setSelected(s.id)}
                    className={`flex w-full items-center gap-3 rounded-lg p-2.5 text-left ${s.id === selected ? 'bg-accent-soft' : 'hover:bg-ground'} ${s.enabled ? '' : 'opacity-50'}`}
                  >
                    <span className="cursor-grab text-[#9AA2AE]" aria-hidden>⋮⋮</span>
                    <span className={`grid h-7 w-7 flex-none place-items-center rounded-full border-2 text-[13px] font-semibold ${ring}`}>{i + 1}</span>
                    <span className="min-w-0 flex-1">
                      <span className="block font-semibold">{STEP_LABEL[s.type]}{s.clicks > 1 ? ` ×${s.clicks}` : ''}</span>
                      <span className="block font-mono text-xs text-muted">{fmtPos(s)}</span>
                    </span>
                    <span className="font-mono text-xs text-muted">{fmtMs(s.delayAfter)}</span>
                  </button>
                </li>
              );
            })}
          </ol>
          <button className="mx-3 mb-3 mt-auto h-11 rounded-lg border border-dashed border-[#B8C0CC] text-[#3F4857]">+ Thêm bước</button>
        </section>

        <section aria-label={`Cấu hình bước ${idx}`} className="flex min-w-0 flex-[999_1_560px] flex-col gap-5 rounded-[10px] border border-line bg-white px-6 py-5">
          <div className="flex items-center gap-3">
            <span className="grid h-8 w-8 place-items-center rounded-full bg-accent font-bold text-white">{idx}</span>
            <h2 className="text-lg font-semibold">Bước {idx}</h2>
            <label className="ml-auto flex items-center gap-2 text-[#3F4857]">
              <input type="checkbox" defaultChecked={cur.enabled} className="h-[18px] w-[18px] accent-accent" />Bật
            </label>
          </div>

          <fieldset className="flex flex-col gap-2">
            <legend className="mb-2 font-semibold">Hành động</legend>
            <div className="flex flex-wrap gap-1.5">
              {types.map((t) => (
                <button key={t} className={`h-10 rounded-lg px-4 ${t === cur.type ? 'border-2 border-accent bg-accent-soft font-semibold text-accent-strong' : 'border border-line bg-white'}`}>
                  {STEP_LABEL[t]}
                </button>
              ))}
            </div>
          </fieldset>

          <div className="grid grid-cols-2 gap-x-5 gap-y-4">
            {cur.type === 'type' && (
              <label className="col-span-2 flex flex-col gap-1.5 font-semibold">Nội dung gõ
                <textarea defaultValue={cur.text} rows={3} className="rounded-lg border border-line px-2.5 py-2 font-mono font-normal" />
                <span className="text-xs font-normal text-muted">Gõ vào ô đang focus – thường đặt sau một bước click. Tắt bộ gõ tiếng Việt khi chạy.</span>
              </label>
            )}
            {cur.type === 'key' && (
              <>
                <label className="flex flex-col gap-1.5 font-semibold">Phím
                  <input defaultValue={cur.key} placeholder="enter, tab, a, f5…" className={`${field} font-normal`} />
                </label>
                <fieldset className="flex flex-col gap-1.5">
                  <legend className="mb-1.5 font-semibold">Kèm phím</legend>
                  <div className="flex flex-wrap gap-1.5">
                    {modifiers.map((m) => (
                      <label key={m} className="flex h-10 items-center gap-1.5 rounded-lg border border-line px-3">
                        <input type="checkbox" defaultChecked={cur.modifiers.includes(m)} className="accent-accent" />
                        <span className="font-mono capitalize">{m}</span>
                      </label>
                    ))}
                  </div>
                </fieldset>
              </>
            )}
            {cur.type !== 'type' && cur.type !== 'key' && (<>
            <label className="flex flex-col gap-1.5 font-semibold">Nút chuột
              <select defaultValue={cur.button} className="h-10 rounded-lg border border-line bg-white px-2.5 font-normal">
                <option value="left">Trái</option><option value="right">Phải</option><option value="center">Giữa</option>
              </select>
            </label>
            <div className="flex flex-col gap-1.5">
              <span className="font-semibold">Vị trí</span>
              <div className="flex gap-2">
                {(['x', 'y'] as const).map((k) => (
                  <label key={k} className="flex h-10 flex-1 items-center gap-1.5 rounded-lg border border-line px-2.5">
                    <span className="uppercase text-muted">{k}</span>
                    <input defaultValue={cur.points[0]?.[k]} className="w-full font-mono outline-none" />
                  </label>
                ))}
                <button aria-label="Chọn lại trên màn hình" className="grid h-10 w-11 place-items-center rounded-lg border border-line"><Target /></button>
              </div>
            </div>
            <label className="flex flex-col gap-1.5 font-semibold">Số click
              <input defaultValue={cur.clicks} className={`${field} font-normal`} />
            </label>
            <MsField label="Giữa các click" value={cur.clickInterval} />
            </>)}
            <MsField label="Chờ trước bước tiếp theo" value={cur.delayAfter} />
          </div>

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
        </section>
      </div>

      <footer className="flex flex-wrap gap-5 border-t border-line bg-white px-5 py-2 text-[13px] text-muted">
        <span>Vùng: toàn màn hình <span className="font-mono">1920 × 1080</span></span>
        <span>Chuột: <span className="font-mono">X 1204 Y 388</span></span>
        <span className="ml-auto">Dừng khẩn cấp: rê chuột vào góc màn hình</span>
      </footer>
    </div>
  );
}

function MsField({ label, value }: { label: string; value: number }) {
  return (
    <label className="flex flex-col gap-1.5 font-semibold">{label}
      <span className="flex h-10 items-center rounded-lg border border-line px-2.5 font-normal">
        <input defaultValue={value} className="flex-1 font-mono outline-none" /><span className="text-muted">ms</span>
      </span>
    </label>
  );
}
