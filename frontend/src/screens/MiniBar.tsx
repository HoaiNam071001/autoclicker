// Thanh mini khi đang chạy – cửa sổ nhỏ always-on-top
import { useEffect, useReducer } from 'react';
import { PauseScript, ResumeScript, StopScript } from '../../wailsjs/go/main/App';
import { Pause, Play, Stop } from '../components/icons';
import { elapsedNow, fmtElapsed, Progress, useEngine } from '../engineStore';
import { inWails, setWindowMode } from '../windowMode';

// Số mẫu khi xem trên trình duyệt (#minibar)
const sample: Progress = { loop: 4, loops: 10, step: 3, total: 6, clicks: 142, elapsedMs: 86_000, paused: false };

export default function MiniBar() {
  const engine = useEngine();
  const p = engine.progress ?? (inWails ? null : sample);
  const running = engine.running || !inWails;
  const paused = !!p?.paused;

  // Vẽ lại mỗi 0,5 giây để đồng hồ chạy giữa các lần Go báo tiến độ.
  const [, tick] = useReducer((n: number) => n + 1, 0);
  useEffect(() => {
    const id = setInterval(tick, 500);
    return () => clearInterval(id);
  }, []);

  const elapsed = inWails ? elapsedNow(engine) : sample.elapsedMs;
  // Vô hạn vòng: thanh tiến độ chỉ theo bước trong vòng hiện tại.
  const progress = !p || p.total === 0 ? 0
    : p.loops > 0 ? ((p.loop - 1) * p.total + p.step) / (p.loops * p.total)
    : p.step / p.total;

  const stop = () => {
    StopScript();
    setWindowMode('main');
  };

  return (
    <div className="drag relative flex h-[76px] w-[460px] items-center gap-3.5 overflow-hidden rounded-[14px] bg-dark-bg pl-4 pr-3 text-dark-text">
      <span className={`h-2.5 w-2.5 flex-none rounded-full ${!running ? 'bg-dark-muted' : paused || engine.countdown > 0 ? 'bg-drag-light' : 'bg-ok'}`} />
      <div className="flex min-w-0 flex-1 flex-col gap-0.5">
        <span className="font-semibold">
          {!running ? 'Chưa chạy'
            : engine.countdown > 0 ? `Bắt đầu sau ${engine.countdown} giây…`
            : !p ? 'Đang bắt đầu…'
            : `${paused ? 'Tạm dừng ở' : 'Đang chạy'} bước ${p.step} / ${p.total}`}
        </span>
        <span className="whitespace-pre font-mono text-xs text-dark-muted">
          Vòng {p?.loop ?? 0} / {p && p.loops > 0 ? p.loops : '∞'}   {p?.clicks ?? 0} click   {fmtElapsed(elapsed)}
        </span>
      </div>
      <button aria-label={paused ? 'Tiếp tục' : 'Tạm dừng'} disabled={!running || engine.countdown > 0}
        onClick={() => (paused ? ResumeScript() : PauseScript())}
        className="grid h-11 w-11 place-items-center rounded-[10px] bg-dark-field disabled:opacity-40">
        {paused ? <Play /> : <Pause />}
      </button>
      <button onClick={stop} className="flex h-11 items-center gap-2 rounded-[10px] bg-danger px-3.5 font-semibold text-white">
        <Stop />Dừng <span className="font-mono text-xs opacity-85">F8</span>
      </button>
      <div className="absolute inset-x-0 bottom-0 h-[3px] bg-dark-field">
        <div className="h-full bg-[#5B82F0]" style={{ width: `${progress * 100}%` }} />
      </div>
    </div>
  );
}
