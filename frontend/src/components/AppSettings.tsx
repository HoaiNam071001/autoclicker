// Cài đặt app (nút bánh răng trên thanh tiêu đề) – API local cho chương trình khác (bước 22): bật / tắt, cổng, token.
// CLI `autoclicker run` luôn dùng được (kênh riêng của app – api.json), không phụ thuộc công tắc này.
import { useState } from 'react';
import { GetSettings, SetAPISettings } from '../../wailsjs/go/main/App';
import { store } from '../../wailsjs/go/models';
import { errorText } from '../engineStore';
import NumInput from './NumInput';
import { Switch } from './Switch';
import { usePopover } from './usePopover';
import { Copy, Gear } from './icons';

const CLI_HINT = 'autoclicker run kich-ban.json --preset "Tên bộ" --set M1=Lan';

export default function AppSettings() {
  const { open, setOpen, ref } = usePopover();
  const [api, setApi] = useState<store.APISettings | null>(null);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  const toggle = () => {
    if (!open) GetSettings().then((s) => setApi(s.api)).catch((e) => setError(errorText(e)));
    setOpen(!open);
  };
  const save = (patch: Partial<store.APISettings>) => {
    if (!api) return;
    const next = store.APISettings.createFrom({ ...api, ...patch });
    setApi(next);
    SetAPISettings(next).then((s) => { setApi(s); setError(''); }).catch((e) => setError(errorText(e)));
  };
  const copy = (text: string) => navigator.clipboard.writeText(text).then(() => { setCopied(true); setTimeout(() => setCopied(false), 1200); });

  return (
    <div ref={ref} className="relative">
      <button onClick={toggle} aria-label="Cài đặt" title="Cài đặt" aria-haspopup="dialog" aria-expanded={open}
        className="grid h-9 w-9 place-items-center text-muted hover:bg-hover hover:text-ink"><Gear /></button>
      {open && api && (
        <div role="dialog" aria-label="Cài đặt" className="absolute right-0 top-full z-30 mt-1 flex w-[340px] flex-col gap-3 rounded-[10px] border border-line bg-surface p-3.5 text-[13px] shadow-[0_12px_32px_rgb(0_0_0/.14)]">
          <span className="font-semibold">Chạy từ ngoài app</span>
          <div className="flex flex-col gap-1 text-muted">
            Dòng lệnh (luôn bật – app chưa mở thì tự mở):
            <code className="break-all rounded-control bg-field px-2 py-1.5 font-mono text-xs text-ink">{CLI_HINT}</code>
          </div>
          <label className="flex items-center justify-between gap-3 border-t border-line pt-3">
            <span className="flex flex-col">API cho chương trình khác
              <span className="text-xs text-muted">n8n, script… gọi http://127.0.0.1 – chỉ máy này</span>
            </span>
            <Switch on={api.enabled} onChange={(enabled) => save({ enabled })} label="Bật API" />
          </label>
          {api.enabled && (
            <>
              <label className="flex items-center justify-between text-muted">Cổng
                <NumInput value={api.port} min={1024} onChange={(port) => save({ port })} className="w-[110px]" />
              </label>
              <div className="flex flex-col gap-1 text-muted">Token (gửi kèm Authorization: Bearer …)
                <div className="flex gap-1.5">
                  <code className="min-w-0 flex-1 truncate rounded-control bg-field px-2 py-1.5 font-mono text-xs text-ink">{api.token}</code>
                  <button onClick={() => copy(api.token ?? '')} aria-label="Chép token" title={copied ? 'Đã chép' : 'Chép token'}
                    className="grid h-8 w-8 place-items-center rounded-control border border-line hover:bg-hover"><Copy className="h-3.5 w-3.5" /></button>
                </div>
                <button onClick={() => save({ token: '' })} className="self-start text-accent-text hover:underline">Tạo token mới (token cũ hết dùng được)</button>
              </div>
              <span className="text-xs text-muted">POST /run {'{'}"path": "…json", "preset": "…"{'}'} · POST /stop · GET /status · GET /events (SSE)</span>
            </>
          )}
          {error && <span className="text-danger-text">{error}</span>}
        </div>
      )}
    </div>
  );
}
