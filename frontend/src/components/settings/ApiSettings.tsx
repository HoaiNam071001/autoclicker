// Cài đặt – chạy từ ngoài app (bước 22): dòng lệnh (luôn bật) + API local cho chương trình khác (bật / tắt, cổng, token).
import { store } from '../../../wailsjs/go/models';
import NumInput from '../NumInput';
import { Switch } from '../Switch';
import { CopyButton, SettingCode } from './parts';

const CLI_HINT = 'ghosthand run kich-ban.json --preset "Tên bộ" --set M1=Lan';

export default function ApiSettings({ api, onChange }: { api: store.APISettings; onChange: (patch: Partial<store.APISettings>) => void }) {
  return (
    <>
      <span className="font-semibold">Chạy từ ngoài app</span>
      <div className="flex flex-col gap-1 text-muted">
        Dòng lệnh (luôn bật – app chưa mở thì tự mở):
        <SettingCode>{CLI_HINT}</SettingCode>
      </div>
      <label className="flex items-center justify-between gap-3 border-t border-line pt-3">
        <span className="flex flex-col">API cho chương trình khác
          <span className="text-xs text-muted">n8n, script… gọi http://127.0.0.1 – chỉ máy này</span>
        </span>
        <Switch on={api.enabled} onChange={(enabled) => onChange({ enabled })} label="Bật API" />
      </label>
      {api.enabled && (
        <>
          <label className="flex items-center justify-between text-muted">Cổng
            <NumInput value={api.port} min={1024} onChange={(port) => onChange({ port })} className="w-[110px]" />
          </label>
          <div className="flex flex-col gap-1 text-muted">Token (gửi kèm Authorization: Bearer …)
            <div className="flex gap-1.5">
              <code className="min-w-0 flex-1 truncate rounded-control bg-field px-2 py-1.5 font-mono text-xs text-ink">{api.token}</code>
              <CopyButton text={api.token ?? ''} label="Chép token" />
            </div>
            <button onClick={() => onChange({ token: '' })} className="self-start text-accent-text hover:underline">Tạo token mới (token cũ hết dùng được)</button>
          </div>
          <span className="text-xs text-muted">POST /run {'{'}"path": "…json", "preset": "…"{'}'} · POST /stop · GET /status · GET /events (SSE) · /mcp (AI)</span>
        </>
      )}
    </>
  );
}
