// Cài đặt app (nút bánh răng trên thanh tiêu đề): chạy từ ngoài app – dòng lệnh, API local (bước 22, settings/ApiSettings)
// và AI qua MCP (bước 31, settings/McpSettings).
import { useState } from 'react';
import { GetSettings, SetAPISettings, SetMCPSettings } from '../../wailsjs/go/main/App';
import { store } from '../../wailsjs/go/models';
import { errorText } from '../engineStore';
import { usePopover } from './usePopover';
import { Gear } from './icons';
import ApiSettings from './settings/ApiSettings';
import McpSettings from './settings/McpSettings';

export default function AppSettings() {
  const { open, setOpen, ref } = usePopover();
  const [settings, setSettings] = useState<store.Settings | null>(null);
  const [error, setError] = useState('');

  const toggle = () => {
    if (!open) GetSettings().then(setSettings).catch((e) => setError(errorText(e)));
    setOpen(!open);
  };
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
    <div ref={ref} className="relative">
      <button onClick={toggle} aria-label="Cài đặt" title="Cài đặt" aria-haspopup="dialog" aria-expanded={open}
        className="grid h-9 w-9 place-items-center text-muted hover:bg-hover hover:text-ink"><Gear /></button>
      {open && settings && (
        <div role="dialog" aria-label="Cài đặt" className="absolute right-0 top-full z-30 mt-1 flex max-h-[calc(100vh-64px)] w-[360px] flex-col gap-3 overflow-y-auto rounded-[10px] border border-line bg-surface p-3.5 text-[13px] shadow-[0_12px_32px_rgb(0_0_0/.14)]">
          <ApiSettings api={settings.api} onChange={saveApi} />
          <McpSettings mcp={settings.mcp} onChange={saveMcp} />
          {error && <span className="text-danger-text">{error}</span>}
        </div>
      )}
    </div>
  );
}
