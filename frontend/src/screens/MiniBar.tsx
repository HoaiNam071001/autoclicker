// Thanh mini khi đang chạy – cửa sổ nhỏ always-on-top
import { Pause, Stop } from '../components/icons';
import { setWindowMode } from '../windowMode';

type Props = { step: number; total: number; loop: number; loops: number; clicks: number; elapsed: string };

export default function MiniBar({ step = 3, total = 6, loop = 4, loops = 10, clicks = 142, elapsed = '01:26' }: Partial<Props>) {
  const progress = ((loop - 1) * total + step) / (loops * total);
  return (
    <div className="drag relative flex h-[76px] w-[460px] items-center gap-3.5 overflow-hidden rounded-[14px] bg-dark-bg pl-4 pr-3 text-dark-text">
      <span className="h-2.5 w-2.5 flex-none rounded-full bg-ok" />
      <div className="flex min-w-0 flex-1 flex-col gap-0.5">
        <span className="font-semibold">Đang chạy bước {step} / {total}</span>
        <span className="whitespace-pre font-mono text-xs text-dark-muted">Vòng {loop} / {loops}   {clicks} click   {elapsed}</span>
      </div>
      <button aria-label="Tạm dừng" className="grid h-11 w-11 place-items-center rounded-[10px] bg-dark-field"><Pause /></button>
      <button onClick={() => setWindowMode('main')} className="flex h-11 items-center gap-2 rounded-[10px] bg-danger px-3.5 font-semibold text-white">
        <Stop />Dừng <span className="font-mono text-xs opacity-85">F8</span>
      </button>
      <div className="absolute inset-x-0 bottom-0 h-[3px] bg-dark-field">
        <div className="h-full bg-[#5B82F0]" style={{ width: `${progress * 100}%` }} />
      </div>
    </div>
  );
}
