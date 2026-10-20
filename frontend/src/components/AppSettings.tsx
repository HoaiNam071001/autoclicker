// Cài đặt app (nút bánh răng trên thanh tiêu đề): chạy từ ngoài app – dòng lệnh, API local (bước 22, settings/ApiSettings)
// và AI qua MCP (bước 31, settings/McpSettings).
import { useEffect, useState } from 'react';
import { GetSettings, SetAPISettings, SetMCPSettings } from '../../wailsjs/go/main/App';
import { store } from '../../wailsjs/go/models';
import { errorText } from '../engineStore';
import { Gear } from './icons';
import ApiSettings from './settings/ApiSettings';
import McpSettings from './settings/McpSettings';

export default function AppSettings() {
  const [open, setOpen] = useState(false);
  const [settings, setSettings] = useState<store.Settings | null>(null);
  const [error, setError] = useState('');

  const toggle = () => {
    if (!open) GetSettings().then(setSettings).catch((e) => setError(errorText(e)));
    setOpen(!open);
  };
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && setOpen(false);
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [open]);
  const saveApi = (patch: Partial<store.APISettings>) => {
    if (!settings) return;
    const api = store.APISettings.createFrom({ ...settings.api, ...patch });
    setSettings(store.Settings.createFrom({ ...settings, api }));
    SetAPISettings(api).then((saved) => { setSettings((s) => s && store.Settings.createFrom({ ...s, api: saved })); setError(''); })
      .catch((e) => setError(errorText(e)));
  };
  const saveMcp = (patch: Partial<store.MCPSettings>) => {
    if (!settings) return;
    const mcp = store.MCPSettings.createFrom({ ...settings.mcp, ...patch });
    setSettings(store.Settings.createFrom({ ...settings, mcp }));
    SetMCPSettings(mcp).then(() => setError('')).catch((e) => setError(errorText(e)));
  };

  return (
    <>
      <button onClick={toggle} aria-label="Cài đặt" title="Cài đặt" aria-haspopup="dialog" aria-expanded={open}
        className="grid h-9 w-9 place-items-center text-muted hover:bg-hover hover:text-ink"><Gear /></button>
      {open && settings && (
        <div className="fixed inset-0 z-40 flex items-start justify-center bg-black/25 px-4 pt-[6vh]" onMouseDown={(e) => e.target === e.currentTarget && setOpen(false)}>
          <div role="dialog" aria-modal="true" aria-label="Cài đặt" className="flex max-h-[88vh] w-[440px] max-w-full flex-col gap-3 overflow-y-auto rounded-panel border border-line bg-surface p-[18px] text-[13px] text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
            <div className="flex items-center justify-between">
              <div className="text-base font-semibold">Cài đặt</div>
              <button onClick={() => setOpen(false)} aria-label="Đóng" className="rounded-md px-2 py-1 text-muted hover:bg-hover hover:text-ink">✕</button>
            </div>
            <ApiSettings api={settings.api} onChange={saveApi} />
            <McpSettings mcp={settings.mcp} onChange={saveMcp} />
            {error && <span className="text-danger-text">{error}</span>}
          </div>
        </div>
      )}
    </>
  );
}
