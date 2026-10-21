// Thử trang trên Chrome của app (bước 44) – chữ + tông màu theo kết luận web.PageCheck.verdict (khớp web.Page* trong Go).
export type Verdict = 'ok' | 'captcha' | 'blocked' | 'canvas' | 'empty';
export type Tone = 'ok' | 'warn' | 'danger';

const VERDICTS: Record<Verdict, { title: string; hint: string; tone: Tone }> = {
  ok: { title: 'Dùng được node web', hint: 'Không thấy captcha – node web (Chờ / Click / Gõ phần tử web) chính xác nhất cho trang này', tone: 'ok' },
  captcha: { title: 'Có captcha', hint: 'Trang kiểm tra “bạn là người” – nên dùng node ảnh / chữ trên trình duyệt thường', tone: 'warn' },
  blocked: { title: 'Trang chặn truy cập', hint: 'Trang không cho trình duyệt của app vào – dùng node ảnh / chữ trên trình duyệt thường', tone: 'danger' },
  canvas: { title: 'Nội dung vẽ trong canvas', hint: 'Game / trang vẽ bằng canvas – node web không thấy nút bên trong, dùng node ảnh / chữ', tone: 'warn' },
  empty: { title: 'Không đọc được gì', hint: 'Trang trống hoặc chưa tải xong – thử lại, hoặc dùng node ảnh / chữ', tone: 'warn' },
};

/** Chữ + tông của kết luận; giá trị lạ coi như trống. */
export const verdictInfo = (v: string) => VERDICTS[v as Verdict] ?? VERDICTS.empty;

/** URL người dùng gõ → URL mở được (thiếu giao thức thì thêm https://), '' khi trống. */
export function normalizeUrl(raw: string): string {
  const u = raw.trim();
  if (!u) return '';
  return /^[a-z][a-z0-9+.-]*:/i.test(u) ? u : `https://${u}`;
}
