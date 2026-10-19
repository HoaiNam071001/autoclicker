// Công tắc 40×20 (design v2): bật nền accent, núm trắng 14px.
export function Switch({ on, onChange, label }: { on: boolean; onChange: (on: boolean) => void; label: string }) {
  return (
    <button role="switch" aria-checked={on} aria-label={label} onClick={() => onChange(!on)}
      className={`relative h-5 w-10 flex-none rounded-full transition-colors ${on ? 'bg-accent' : 'bg-line-strong'}`}>
      <span className={`absolute top-[3px] h-3.5 w-3.5 rounded-full bg-white shadow-sm transition-all ${on ? 'left-[23px]' : 'left-[3px]'}`} />
    </button>
  );
}
