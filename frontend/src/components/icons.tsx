type P = { className?: string };
const base = { fill: 'none', stroke: 'currentColor', strokeWidth: 2, 'aria-hidden': true } as const;

export const Target = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><circle cx="12" cy="12" r="5" /><path d="M12 3v3M12 18v3M3 12h3M18 12h3" /></svg>
);
export const Logo = ({ className = 'h-[22px] w-[22px]' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><circle cx="12" cy="12" r="7" /><path d="M12 1v5M12 18v5M1 12h5M18 12h5" /></svg>
);
export const Region = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><rect x="3" y="5" width="18" height="14" rx="1" strokeDasharray="3 3" /></svg>
);
export const Drag = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><circle cx="5" cy="18" r="2" /><path d="M7 16 18 6M13 6h5v5" /></svg>
);
export const Play = ({ className = 'h-3.5 w-3.5' }: P) => (
  <svg viewBox="0 0 24 24" className={className} fill="currentColor" aria-hidden><path d="M6 4l14 8-14 8z" /></svg>
);
export const Pause = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} fill="currentColor" aria-hidden><rect x="6" y="4" width="4" height="16" rx="1" /><rect x="14" y="4" width="4" height="16" rx="1" /></svg>
);
export const Stop = ({ className = 'h-3 w-3' }: P) => (
  <svg viewBox="0 0 24 24" className={className} fill="currentColor" aria-hidden><rect x="3" y="3" width="18" height="18" rx="2" /></svg>
);
