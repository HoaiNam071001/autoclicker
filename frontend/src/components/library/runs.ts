// Dòng lịch sử chạy hiển thị ở thư viện + hộp thoại Lịch sử chạy (dùng chung): nguồn chạy, giờ, kết quả.
import type { RunSummary } from '../../runStore';

export type LibraryRun = { id?: string; script: string; when: string; result: string; ok: boolean; muted?: boolean; active?: boolean };

const SOURCE: Record<string, string> = { user: 'Chạy tay', schedule: 'Hẹn giờ', cli: 'Dòng lệnh', ai: 'AI' };
const two = (n: number) => String(n).padStart(2, '0');

// "14:32" hôm nay, "14:32 16/10" ngày khác
export const whenOf = (iso: string) => {
  const d = new Date(iso);
  if (isNaN(d.getTime())) return '';
  const t = `${two(d.getHours())}:${two(d.getMinutes())}`;
  return d.toDateString() === new Date().toDateString() ? t : `${t} ${two(d.getDate())}/${two(d.getMonth() + 1)}`;
};

export const libraryRun = (h: RunSummary, active = false): LibraryRun => ({
  id: h.id, script: SOURCE[h.source] ?? h.source, when: whenOf(h.startedAt), ok: h.status !== 'fail', muted: h.status === 'stopped',
  result: h.status === 'fail' ? `Thất bại: ${h.error ?? ''}` : h.status === 'stopped' ? 'Đã dừng' : `Thành công${h.message ? `: ${h.message}` : ''}`,
  active,
});

export const runDot = (h: LibraryRun) => (h.muted ? 'bg-muted' : h.ok ? 'bg-ok' : 'bg-danger');
