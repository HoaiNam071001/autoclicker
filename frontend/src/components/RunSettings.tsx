// Cài đặt chạy (bước 16): lặp theo số vòng / vô hạn / theo thời gian, đếm ngược trước khi chạy
import NumInput from './NumInput';
import { setCountdown, updateScript, useEditor } from '../scriptStore';
import { Script } from '../types';

type Mode = 'loops' | 'infinite' | 'time';
const MODES: { id: Mode; label: string }[] = [
  { id: 'loops', label: 'Số vòng' },
  { id: 'infinite', label: 'Vô hạn' },
  { id: 'time', label: 'Theo thời gian' },
];

// Script chỉ có loops (0 = vô hạn) và maxMinutes (0 = không giới hạn) – suy ra chế độ từ hai số này.
const modeOf = (sc: Script): Mode => (sc.loops > 0 ? 'loops' : sc.maxMinutes > 0 ? 'time' : 'infinite');

// Nhớ số vòng / số phút đã nhập để chuyển chế độ qua lại không mất.
let lastLoops = 10;
let lastMinutes = 5;

export default function RunSettings() {
  const { script, countdown } = useEditor();
  const mode = modeOf(script);

  const setMode = (m: Mode) => {
    if (script.loops > 0) lastLoops = script.loops;
    if (script.maxMinutes > 0) lastMinutes = script.maxMinutes;
    updateScript(
      m === 'loops' ? { loops: lastLoops, maxMinutes: 0 }
        : m === 'time' ? { loops: 0, maxMinutes: lastMinutes }
        : { loops: 0, maxMinutes: 0 },
    );
  };

  return (
    <div className="mt-auto flex flex-wrap items-center gap-x-5 gap-y-3 border-t border-line pt-4">
      <span className="font-semibold">Lặp kịch bản</span>
      <div className="flex overflow-hidden rounded-lg border border-line" role="radiogroup" aria-label="Chế độ lặp">
        {MODES.map((m, i) => (
          <button key={m.id} role="radio" aria-checked={mode === m.id} onClick={() => setMode(m.id)}
            className={`h-10 px-3.5 ${i > 0 ? 'border-l border-line' : ''} ${mode === m.id ? 'bg-ink text-white' : 'bg-white hover:bg-ground'}`}>
            {m.label}
          </button>
        ))}
      </div>
      {mode === 'loops' && (
        <span className="w-[120px]" aria-label="Số vòng">
          <NumInput value={script.loops} min={1} unit="vòng" onChange={(loops) => updateScript({ loops })} />
        </span>
      )}
      {mode === 'time' && (
        <span className="w-[120px]" aria-label="Số phút">
          <NumInput value={script.maxMinutes} min={1} unit="phút" onChange={(maxMinutes) => updateScript({ maxMinutes })} />
        </span>
      )}
      {mode === 'infinite' && <span className="text-[13px] text-muted">Chạy tới khi bấm Dừng / F8</span>}
      <label className="flex items-center gap-2 text-[#3F4857]">
        <input type="checkbox" checked={countdown} onChange={(e) => setCountdown(e.target.checked)} className="h-[18px] w-[18px] accent-accent" />
        Đếm ngược 3 giây
      </label>
    </div>
  );
}
