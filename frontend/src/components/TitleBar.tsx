// Thanh tiêu đề tự vẽ kiểu Windows 11 (cửa sổ frameless): kéo được cửa sổ, nút thu nhỏ / phóng to / đóng bên phải.
import { Quit, WindowMinimise, WindowToggleMaximise } from '../../wailsjs/runtime/runtime';
import { inWails } from '../windowMode';
import { ReactNode } from 'react';
import { Logo, Monitor, Moon, Sun } from './icons';
import { setTheme, Theme, useTheme } from '../theme';

const winBase = 'grid h-10 w-[46px] place-items-center text-ink';
const winBtn = `${winBase} hover:bg-hover`;
// Nút đóng: nền đỏ khi rê chuột (không gộp với hover:bg-hover – hai class hover cùng lúc thì bg-hover thắng)
const closeBtn = `${winBase} hover:bg-danger hover:text-white`;

export default function TitleBar() {
  return (
    // double click vùng kéo = phóng to / thu về, giống cửa sổ thường
    <div className="drag flex h-10 flex-none items-center gap-2.5 bg-chrome pl-3.5" 
      onDoubleClick={(e) => inWails && !(e.target as HTMLElement).closest('button') && WindowToggleMaximise()}>
      <Logo className="h-[18px] w-[18px] flex-none" />
      <span className="text-xs text-muted">Auto Clicker</span>
      <ThemeSwitch />
      {inWails && (
        <div className="flex">
          <button aria-label="Thu nhỏ" onClick={WindowMinimise} className={winBtn}>
            <svg width="10" height="10" viewBox="0 0 10 10" aria-hidden><path d="M0 5h10" stroke="currentColor" /></svg>
          </button>
          <button aria-label="Phóng to" onClick={WindowToggleMaximise} className={winBtn}>
            <svg width="10" height="10" viewBox="0 0 10 10" aria-hidden><rect x=".5" y=".5" width="9" height="9" rx="1.5" fill="none" stroke="currentColor" /></svg>
          </button>
          <button aria-label="Đóng" onClick={Quit} className={closeBtn}>
            <svg width="10" height="10" viewBox="0 0 10 10" aria-hidden><path d="M0 0l10 10M10 0L0 10" stroke="currentColor" /></svg>
          </button>
        </div>
      )}
    </div>
  );
}

const THEMES: { id: Theme; label: string; icon: ReactNode }[] = [
  { id: 'light', label: 'Sáng', icon: <Sun /> },
  { id: 'dark', label: 'Tối', icon: <Moon /> },
  { id: 'system', label: 'Theo hệ thống', icon: <Monitor /> },
];

// Chọn theme: Sáng / Tối / Theo hệ thống – lưu vào cấu hình app.
function ThemeSwitch() {
  const theme = useTheme();
  return (
    <div role="radiogroup" aria-label="Giao diện" className="mr-1 ml-auto flex rounded-md bg-soft p-0.5">
      {THEMES.map((t) => (
        <button key={t.id} role="radio" aria-checked={theme === t.id} aria-label={t.label} title={`Giao diện: ${t.label}`}
          onClick={() => setTheme(t.id)}
          className={`grid h-6 w-7 place-items-center rounded ${theme === t.id ? 'bg-raised text-ink shadow-[var(--c-raise)]' : 'text-muted hover:text-ink'}`}>
          {t.icon}
        </button>
      ))}
    </div>
  );
}
