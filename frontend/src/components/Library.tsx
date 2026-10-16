// Thanh bên thư viện kịch bản (artboard Main): tìm, danh sách kịch bản (chấm "chưa lưu", giờ hẹn), mở gần đây,
// lịch sử chạy, Kịch bản mới / Mở / Lưu. Thu gọn / mở rộng bằng nút ở đầu cột hoặc Ctrl/Cmd+B (nhớ bằng localStorage).
// Giao diện dựng ở bước 27b; lưu / mở thật ở bước 21, lịch sử ở 25c, hẹn giờ ở 22 – app thật hiện tắt + "Sắp có".
import { ReactNode, useEffect, useState } from 'react';
import { useScript } from '../scriptStore';
import { SidebarToggle } from './icons';
import { Soon } from './StepCard';

export type LibraryScript = { name: string; unsaved?: boolean; schedule?: string };
export type LibraryRun = { script: string; when: string; result: string; ok: boolean };
export type LibraryProps = {
  scripts: LibraryScript[];
  active: number;
  recent: string[];
  history: LibraryRun[];
  ready: boolean; // false: ô tìm / nút chưa hoạt động (app thật trước bước 21)
};

const KEY = 'library-collapsed';
const isMac = navigator.platform.startsWith('Mac');
const mod = isMac ? '⌘' : 'Ctrl+';

const loadCollapsed = () => {
  try {
    return localStorage.getItem(KEY) === '1';
  } catch {
    return false;
  }
};

// App thật: chỉ có kịch bản đang soạn, chưa lưu được.
export default function Library() {
  const { name } = useScript();
  return <LibraryView scripts={[{ name, unsaved: true }]} active={0} recent={[]} history={[]} ready={false} />;
}

export function LibraryView({ scripts, active, recent, history, ready }: LibraryProps) {
  const [collapsed, setCollapsed] = useState(loadCollapsed);
  const [q, setQ] = useState('');
  useEffect(() => {
    try {
      localStorage.setItem(KEY, collapsed ? '1' : '0');
    } catch { /* không lưu được thì thôi – chỉ là tiện ích */ }
  }, [collapsed]);
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key.toLowerCase() !== 'b' || !(isMac ? e.metaKey : e.ctrlKey)) return;
      e.preventDefault();
      setCollapsed((c) => !c);
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);

  const toggle = (
    <button onClick={() => setCollapsed(!collapsed)} aria-expanded={!collapsed}
      aria-label={collapsed ? 'Mở rộng thư viện' : 'Thu gọn thư viện'} title={`${collapsed ? 'Mở rộng' : 'Thu gọn'} thư viện (${mod}B)`}
      className="grid h-7 w-7 flex-none place-items-center rounded-md text-muted hover:bg-hover hover:text-ink">
      <SidebarToggle open={!collapsed} />
    </button>
  );

  if (collapsed) {
    return (
      <nav aria-label="Thư viện kịch bản (thu gọn)" className="flex w-11 flex-none flex-col items-center border-r border-line bg-side py-3">
        {toggle}
      </nav>
    );
  }

  const shown = scripts.map((s, i) => ({ s, i })).filter(({ s }) => s.name.toLowerCase().includes(q.trim().toLowerCase()));
  const soon = (step: string) => (ready ? undefined : `Sắp có (bước ${step})`);
  return (
    <nav aria-label="Thư viện kịch bản" className="flex w-[236px] flex-none flex-col gap-3 overflow-y-auto border-r border-line bg-side px-2 py-3 max-[1100px]:w-[200px]">
      <div className="flex items-center gap-1">
        <label className={`mx-1 flex h-[34px] min-w-0 flex-1 items-center gap-2 rounded-control border border-line bg-surface px-2.5 text-muted focus-within:ring-focus ${ready ? '' : 'opacity-60'}`}
          title={soon('21')}>
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden><circle cx="11" cy="11" r="7" /><path d="M20 20l-3.5-3.5" /></svg>
          <input value={q} onChange={(e) => setQ(e.target.value)} disabled={!ready} placeholder="Tìm kịch bản"
            className="min-w-0 flex-1 bg-transparent text-ink outline-none" />
        </label>
        {toggle}
      </div>

      <Section title="Kịch bản">
        {shown.map(({ s, i }) => (
          <a key={i} href="#" onClick={(e) => e.preventDefault()} aria-current={i === active}
            className={`relative flex h-9 items-center gap-2.5 rounded-md px-3 ${i === active ? 'bg-hover font-semibold' : 'hover:bg-hover'}`}>
            {i === active && <span className="absolute top-2.5 left-0 h-4 w-[3px] rounded-sm bg-accent" />}
            <span className="min-w-0 flex-1 truncate">{s.name}</span>
            {s.schedule && (
              <span className="flex items-center gap-1 font-mono text-[11px] font-normal text-muted" title={`Hẹn chạy lúc ${s.schedule}`}>
                <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4" aria-hidden><circle cx="12" cy="12" r="9" /><path d="M12 7v5l3 2" /></svg>
                {s.schedule}
              </span>
            )}
            {s.unsaved && <span className="h-1.5 w-1.5 flex-none rounded-full bg-warn" aria-label="Chưa lưu" title="Chưa lưu" />}
          </a>
        ))}
        {shown.length === 0 && <span className="px-3 py-1.5 text-[13px] text-muted">Không có kịch bản nào khớp</span>}
      </Section>

      <Section title="Mở gần đây" soonStep={ready ? undefined : '21'}>
        {recent.map((r) => (
          <a key={r} href="#" onClick={(e) => e.preventDefault()} className="truncate rounded-md px-3 py-[7px] text-[13px] text-muted hover:bg-hover hover:text-ink">{r}</a>
        ))}
        {recent.length === 0 && <span className="px-3 py-1 text-[13px] text-muted">Chưa có</span>}
      </Section>

      <Section title="Lịch sử chạy" soonStep={ready ? undefined : '25c'}>
        {history.map((h, i) => (
          <a key={i} href="#" onClick={(e) => e.preventDefault()} className="flex items-center gap-2 rounded-md px-3 py-[7px] text-[13px] hover:bg-hover">
            <span className={`h-1.5 w-1.5 flex-none rounded-full ${h.ok ? 'bg-ok' : 'bg-danger'}`} />
            <span className="min-w-0 flex-1 truncate">{h.script}</span>
            <span className="font-mono text-[11px] text-muted" title={h.result}>{h.when}</span>
          </a>
        ))}
        {history.length === 0 && <span className="px-3 py-1 text-[13px] text-muted">Chưa có lần chạy nào</span>}
      </Section>

      <div className="mt-auto flex flex-col gap-1 px-1">
        <button disabled={!ready} title={soon('21')} className="h-9 rounded-control border border-line bg-surface hover:bg-hover disabled:opacity-60">
          + Kịch bản mới <span className="font-mono text-[11px] text-muted">{mod}N</span>
        </button>
        <div className="flex gap-1">
          <button disabled={!ready} title={soon('21')} className="h-[34px] flex-1 rounded-md text-muted hover:bg-hover hover:text-ink disabled:opacity-60">Mở…</button>
          <button disabled={!ready} title={soon('21')} className="h-[34px] flex-1 rounded-md text-muted hover:bg-hover hover:text-ink disabled:opacity-60">
            Lưu <span className="font-mono text-[11px]">{mod}S</span>
          </button>
        </div>
      </div>
    </nav>
  );
}

function Section({ title, soonStep, children }: { title: string; soonStep?: string; children: ReactNode }) {
  return (
    <div className="flex flex-col gap-0.5">
      <span className="flex items-center gap-1.5 px-3 py-1 text-xs font-semibold text-muted">{title}{soonStep && <Soon step={soonStep} />}</span>
      {children}
    </div>
  );
}
