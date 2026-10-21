// Hộp Thử trang trên Chrome của app (bước 44) – chỉ hiển thị: ô URL, nút Thử, kết luận + dấu hiệu. PageCheckButton giữ
// state + gọi Go; màn mẫu #demo dùng lại với dữ liệu giả.
import { web } from '../../../wailsjs/go/models';
import { Tone, verdictInfo } from '../../pageCheck';

const TONE: Record<Tone, string> = { ok: 'text-ok', warn: 'text-warn', danger: 'text-danger-text' };
const ICON: Record<Tone, string> = { ok: '✓', warn: '⚠', danger: '✕' };

type Props = {
  url: string;
  onUrl: (u: string) => void;
  busy: boolean;
  result: web.PageCheck | null;
  error: string;
  onCheck: () => void;
  onRecheck: () => void; // dò lại tab đang mở, không tải lại
  onClose: () => void;
};

export default function PageCheckDialog({ url, onUrl, busy, result, error, onCheck, onRecheck, onClose }: Props) {
  return (
    <div className="fixed inset-0 z-40 flex items-start justify-center bg-black/25 px-4 pt-[6vh]" onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div role="dialog" aria-modal="true" aria-label="Thử trang" className="flex max-h-[88vh] w-[440px] max-w-full flex-col gap-3 overflow-y-auto rounded-panel border border-line bg-surface p-[18px] text-[13px] text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
        <div className="flex items-center justify-between">
          <div className="text-base font-semibold">Thử trang trên Chrome của app</div>
          <button onClick={onClose} aria-label="Đóng" className="rounded-md px-2 py-1 text-muted hover:bg-hover hover:text-ink">✕</button>
        </div>
        <span className="text-xs text-muted">Xem trang có captcha / chặn bot không – không có thì node web dùng được (chính xác nhất). Chỉ để thử: không đăng nhập tài khoản thật ở đây.</span>
        <form className="flex gap-1.5" onSubmit={(e) => { e.preventDefault(); if (!busy) onCheck(); }}>
          <input value={url} onChange={(e) => onUrl(e.target.value)} placeholder="vd game-a.vn/choi" aria-label="URL trang" autoFocus
            className="h-8 flex-1 rounded-control border border-line bg-field px-2.5 font-mono text-xs outline-none focus:ring-focus" />
          <button type="submit" disabled={busy || !url.trim()} className="btn-primary rounded-control px-3.5 text-on-accent disabled:opacity-50">
            {busy ? 'Đang thử…' : 'Thử'}
          </button>
        </form>
        {busy && <span className="text-xs text-muted">Đang chờ trang hiện xong rồi dò (tối đa ~30 giây)…</span>}
        {result && !busy && <PageCheckResult result={result} />}
        {error && <span className="text-danger-text">{error}</span>}
        {(result || error) && !busy && (
          <button onClick={onRecheck} title="Dò lại trang đang mở trên Chrome của app, không tải lại – cho trang tải lâu"
            className="h-8 self-start rounded-control border border-line bg-surface px-3 hover:bg-hover">Dò lại</button>
        )}
      </div>
    </div>
  );
}

function PageCheckResult({ result }: { result: web.PageCheck }) {
  const v = verdictInfo(result.verdict);
  return (
    <div className="flex flex-col gap-1.5 rounded-control bg-soft p-2.5">
      <span className={`font-semibold ${TONE[v.tone]}`}>{ICON[v.tone]} {v.title}</span>
      <span className="text-xs text-muted">{v.hint}</span>
      <span className="truncate text-xs" title={result.url}>{result.title || result.url}
        <span className="text-muted"> · <span className="font-mono">{result.elements}</span> phần tử đọc được</span>
      </span>
      {result.signs.length > 0 && (
        <ul className="flex list-disc flex-col gap-0.5 pl-4 text-xs text-muted">
          {result.signs.map((s, i) => <li key={i}>{s}</li>)}
        </ul>
      )}
    </div>
  );
}
