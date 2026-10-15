// Phương án B – xem trước màn hình là trung tâm (nền tối)
import { Logo, Play } from '../components/icons';

type Pt = { label: string; x: number; y: number; name?: string; delay?: string; drag?: boolean };
const pts: Pt[] = [
  { label: '1', x: 26, y: 24, name: 'Click', delay: '500 ms' },
  { label: '2', x: 40, y: 24, name: 'Double', delay: '300 ms' },
  { label: '3', x: 51, y: 48, name: 'Click ×5', delay: '800 ms' },
  { label: '4A', x: 19, y: 67, name: 'Kéo thả', delay: '1 s', drag: true },
  { label: '4B', x: 56, y: 67, drag: true },
  { label: '5', x: 69, y: 58, name: 'Click phải', delay: '500 ms' },
];
const SELECTED = '3';
const input = 'h-10 rounded-lg border border-dark-border bg-dark-bg px-2.5 font-mono text-dark-text';

export default function PreviewScreen() {
  return (
    <div className="flex min-h-screen flex-col bg-dark-bg text-dark-text">
      <header className="flex flex-wrap items-center gap-2.5 border-b border-dark-line px-4 py-2.5">
        <span className="text-accent"><Logo /></span>
        <span className="mr-4 text-base font-bold">Kịch bản 1</span>
        <div className="flex gap-1 rounded-[10px] bg-dark-panel p-1">
          {['Chọn điểm', 'Kéo thả', 'Chọn vùng'].map((m, i) => (
            <button key={m} className={`h-10 rounded-[7px] px-3.5 ${i === 0 ? 'bg-dark-border text-white' : 'text-[#C3CAD5]'}`}>{m}</button>
          ))}
        </div>
        <div className="ml-auto flex items-center gap-2.5">
          <span className="text-dark-muted">10 vòng</span>
          <button className="flex h-11 items-center gap-2 rounded-lg bg-accent px-5 font-semibold text-white"><Play />Chạy</button>
        </div>
      </header>

      <div className="flex flex-1 flex-wrap gap-4 p-4">
        <section aria-label="Xem trước màn hình" className="flex min-w-0 flex-[999_1_640px] flex-col gap-2.5">
          {/* Thay nền này bằng ảnh chụp màn hình thật */}
          <div className="relative aspect-video w-full overflow-hidden rounded-[10px] border border-dark-line bg-[#2A303B]">
            <div className="absolute left-[12%] top-[14%] h-[66%] w-[66%] rounded border-2 border-dashed border-accent bg-accent/10" />
            <span className="absolute left-[12%] top-[calc(14%-26px)] rounded bg-accent px-2 py-0.5 font-mono text-xs text-white">Vùng 1056 × 594</span>
            <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="absolute inset-0 h-full w-full" aria-hidden>
              <polyline points="26,24 40,24 51,48 19,67" fill="none" stroke="#8FA6EE" strokeWidth={1.5} strokeDasharray="4 4" vectorEffect="non-scaling-stroke" />
              <line x1="19" y1="67" x2="56" y2="67" stroke="#F0913A" strokeWidth={2.5} vectorEffect="non-scaling-stroke" />
              <polyline points="56,67 69,58" fill="none" stroke="#8FA6EE" strokeWidth={1.5} strokeDasharray="4 4" vectorEffect="non-scaling-stroke" />
            </svg>
            {pts.map((p) => {
              const sel = p.label === SELECTED;
              return (
                <span key={p.label}
                  style={{ left: `${p.x}%`, top: `${p.y}%` }}
                  className={`absolute grid -translate-x-1/2 -translate-y-1/2 place-items-center rounded-full border-2 text-xs font-bold shadow-[0_0_0_4px_rgba(28,33,43,.6)]
                    ${sel ? 'h-[34px] w-[34px] bg-accent' : 'h-7 w-7 bg-dark-bg'} ${p.drag ? 'border-drag-light' : 'border-accent'}`}>
                  {p.label}
                </span>
              );
            })}
          </div>
          <div className="flex gap-4 text-[13px] text-dark-muted">
            <span>Kéo điểm để sửa vị trí</span><span>Click chuột phải để xoá</span>
            <span className="ml-auto font-mono">X 812 Y 436</span>
          </div>
        </section>

        <aside aria-label="Thuộc tính bước" className="flex max-w-[360px] flex-[1_1_300px] flex-col gap-4 rounded-[10px] border border-dark-line bg-dark-panel p-[18px]">
          <div className="flex items-center gap-2.5">
            <span className="grid h-[30px] w-[30px] place-items-center rounded-full bg-accent font-bold">3</span>
            <h2 className="text-base font-semibold">Click trái</h2>
          </div>
          <div className="grid grid-cols-2 gap-3">
            {[['X', 812], ['Y', 436], ['Số click', 5], ['Giữa click (ms)', 120]].map(([l, v]) => (
              <label key={l} className="flex flex-col gap-1.5 text-dark-muted">{l}<input defaultValue={v} className={input} /></label>
            ))}
          </div>
          <Range label="Chờ trước bước tiếp theo" max={3000} value={800} unit="ms" />
          <Range label="Ngẫu nhiên vị trí" max={20} value={3} unit="px" prefix="±" />
        </aside>
      </div>

      <section aria-label="Thứ tự bước" className="px-4 pb-4">
        <div className="overflow-x-auto">
          <ol className="flex min-w-max items-center gap-2 rounded-[10px] border border-dark-line bg-dark-panel p-3">
            {pts.filter((p) => p.name).map((p) => (
              <li key={p.label} className="flex items-center gap-2">
                <span className={`flex h-11 items-center gap-2 rounded-lg border px-3 ${p.label === SELECTED ? 'border-accent bg-[#2F3A55]' : 'border-dark-border bg-dark-bg'}`}>
                  <span className={`grid h-[22px] w-[22px] place-items-center rounded-full border-2 text-[11px] font-bold ${p.drag ? 'border-drag-light' : 'border-accent'}`}>{p.label}</span>
                  {p.name}
                  <span className="font-mono text-xs text-dark-muted">{p.delay}</span>
                </span>
                <span className="text-[#6B7484]" aria-hidden>›</span>
              </li>
            ))}
            <li><button className="h-11 rounded-lg border border-dashed border-[#4A5464] px-3.5 text-[#C3CAD5]">+ Thêm</button></li>
          </ol>
        </div>
      </section>
    </div>
  );
}

function Range({ label, max, value, unit, prefix = '' }: { label: string; max: number; value: number; unit: string; prefix?: string }) {
  return (
    <label className="flex flex-col gap-1.5 text-dark-muted">{label}
      <span className="flex items-center gap-3">
        <input type="range" min={0} max={max} defaultValue={value} className="flex-1 accent-accent" />
        <span className="font-mono text-dark-text">{prefix}{value} {unit}</span>
      </span>
    </label>
  );
}
