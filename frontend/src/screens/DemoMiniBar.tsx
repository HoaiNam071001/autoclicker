// Màn #demo-minibar (bước 27b): các trạng thái thanh mini, gồm trạng thái chưa có logic –
// đang chờ ảnh (vòng theo thời gian chờ, bước 25), thử lại (25b), tên bộ giá trị input (21b).
import { MiniBarView } from './MiniBar';

const noop = () => {};
const STATES: { note: string; props: Parameters<typeof MiniBarView>[0] }[] = [
  {
    note: 'Đang chờ ảnh – vòng theo thời gian chờ (bước 25)',
    props: { title: 'Bước 3 / 8: Chờ ảnh nút Nhắn tin', lines: ['đang tìm 4.2 / 10 s   thử lại 0 / 2', 'vòng 2 / 10   38 click   01:26'], icon: 'waitImage', progress: 0.42, paused: false, countdown: 0, canPause: true, onPause: noop, onStop: noop },
  },
  {
    note: 'Hết giờ, đang thử lại (bước 25b) + tên bộ giá trị input (bước 21b)',
    props: { title: 'Bước 3 / 8: Chờ ảnh nút Nhắn tin', lines: ['thử lại 1 / 2   đang tìm 1.0 / 10 s', 'Gửi cho Lan · vòng 2 / 10   01:38'], icon: 'waitImage', progress: 0.1, paused: false, countdown: 0, canPause: true, onPause: noop, onStop: noop },
  },
  {
    note: 'Đang chạy (đã có)',
    props: { title: 'Bước 5 / 8: Gõ chữ', lines: ['vòng 4 / 10   142 click   02:05'], icon: 'type', progress: 0.36, paused: false, countdown: 0, canPause: true, onPause: noop, onStop: noop },
  },
  {
    note: 'Tạm dừng (đã có)',
    props: { title: 'Tạm dừng · Bước 5 / 8: Gõ chữ', lines: ['vòng 4 / 10   142 click   02:05'], icon: 'type', progress: 0.36, paused: true, countdown: 0, canPause: true, onPause: noop, onStop: noop },
  },
  {
    note: 'Đếm ngược (đã có)',
    props: { title: 'Bắt đầu sau 2 giây…', lines: ['vòng 0 / 10   0 click   00:00'], icon: 'click', progress: 0, paused: false, countdown: 2, canPause: false, onPause: noop, onStop: noop },
  },
];

export default function DemoMiniBar() {
  return (
    <div className="flex min-h-screen flex-col items-center gap-5 bg-bg p-8 text-ink">
      <span className="rounded-full bg-accent-soft px-2.5 py-0.5 text-xs font-semibold text-accent-text">Màn mẫu #demo-minibar</span>
      {STATES.map((s) => (
        <div key={s.note} className="flex flex-col gap-1.5">
          <span className="text-xs text-muted">{s.note}</span>
          <MiniBarView {...s.props} />
        </div>
      ))}
    </div>
  );
}
