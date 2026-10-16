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
export const Loop = ({ className = 'h-[15px] w-[15px]' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base} strokeWidth={1.8}><path d="M17 2l4 4-4 4M3 11V9a3 3 0 0 1 3-3h15M7 22l-4-4 4-4M21 13v2a3 3 0 0 1-3 3H3" /></svg>
);
export const Power = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><path d="M12 3v8M6.3 6.3a8 8 0 1 0 11.4 0" /></svg>
);
export const Sun = ({ className = 'h-3.5 w-3.5' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base} strokeLinecap="round"><circle cx="12" cy="12" r="4" /><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" /></svg>
);
export const Moon = ({ className = 'h-3.5 w-3.5' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base} strokeLinejoin="round"><path d="M20 14.5A8 8 0 0 1 9.5 4a8 8 0 1 0 10.5 10.5z" /></svg>
);
export const Monitor = ({ className = 'h-3.5 w-3.5' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><rect x="3" y="4" width="18" height="12" rx="2" /><path d="M8 20h8M12 16v4" /></svg>
);
// Thu gọn / mở rộng thanh bên: khung có cột trái, mũi tên theo hướng
export const SidebarToggle = ({ className = 'h-4 w-4', open = true }: P & { open?: boolean }) => (
  <svg viewBox="0 0 24 24" className={className} {...base} strokeWidth={1.8} strokeLinecap="round" strokeLinejoin="round">
    <rect x="3" y="4" width="18" height="16" rx="2" /><path d="M9 4v16" /><path d={open ? 'M15.5 10l-2 2 2 2' : 'M13.5 10l2 2-2 2'} />
  </svg>
);
// Chọn nhiều (khung chọn nét đứt + con trỏ)
export const SelectIcon = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base} strokeWidth={1.8} strokeLinecap="round" strokeLinejoin="round">
    <path d="M3 9V5a2 2 0 0 1 2-2h4M15 3h4a2 2 0 0 1 2 2v4M3 15v4a2 2 0 0 0 2 2h4" /><path d="M13 13l8 3.5-3.5 1-1 3.5z" />
  </svg>
);
export const Copy = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><rect x="8" y="8" width="12" height="12" rx="2" /><path d="M16 8V5a1 1 0 0 0-1-1H5a1 1 0 0 0-1 1v10a1 1 0 0 0 1 1h3" /></svg>
);
export const Trash = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><path d="M4 7h16M10 11v6M14 11v6M6 7l1 13h10l1-13M9 7V4h6v3" /></svg>
);
export const Eye = ({ className = 'h-4 w-4' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base}><path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7S2 12 2 12z" /><circle cx="12" cy="12" r="3" /></svg>
);

// Icon loại bước – design v2 (design/DESIGN.md, "Icon bước"): SVG 24×24, nét 1.8, đầu tròn.
// Có cả loại chưa làm (scroll, open, waitImage…) để dựng sẵn UI (bước 27b).
export const STEP_ICON = {
  click: 'M5 3l13 7-6 1.5L9 18z',
  double: 'M5 3l13 7-6 1.5L9 18zM15 3l2-2M19 6l2-1',
  hold: 'M5 3l13 7-6 1.5L9 18zM17 15a3 3 0 1 0 0 .01',
  drag: 'M5 19a2 2 0 1 0 0-.01M7 17L18 6M13 6h5v5',
  scroll: 'M12 4v16M8 8l4-4 4 4M8 16l4 4 4-4',
  type: 'M4 7V5h12v2M10 5v14M7 19h6M15 13h5M17.5 13v6',
  key: 'M5 9h14v10H5zM8 13h1M11 13h1M14 13h2M8 16h8',
  wait: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zM12 7v5l3 2',
  waitImage: 'M4 8V5h3M17 5h3v3M20 16v3h-3M7 19H4v-3M9 12a3 3 0 1 0 6 0 3 3 0 0 0-6 0',
  clickImage: 'M4 8V5h3M17 5h3v3M7 19H4v-3M10 9l9 5-4 1-2 4z',
  waitPixel: 'M4 4h7v7H4zM13 13h7v7h-7zM13 4h7v7h-7z',
  open: 'M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zM3 12h18M12 3c2.5 2.7 2.5 15.3 0 18M12 3c-2.5 2.7-2.5 15.3 0 18',
  focusWindow: 'M4 5h16v14H4zM4 9h16',
  group: 'M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3zM14 14h7v7h-7z',
} as const;
export type StepIconName = keyof typeof STEP_ICON;

export const StepIcon = ({ type, className = 'h-4 w-4' }: P & { type: StepIconName }) => (
  <svg viewBox="0 0 24 24" className={className} fill="none" stroke="currentColor" strokeWidth={1.8}
    strokeLinecap="round" strokeLinejoin="round" aria-hidden><path d={STEP_ICON[type]} /></svg>
);
export const Pin = ({ className = 'h-3.5 w-3.5' }: P) => (
  <svg viewBox="0 0 24 24" className={className} {...base} strokeLinecap="round" strokeLinejoin="round"><path d="M9 4h6l-1 5 3 3v2H7v-2l3-3zM12 14v6" /></svg>
);
