// Hẹn giờ + thông báo (giao diện cho bước 22): chạy lúc HH:mm một lần / hằng ngày / các ngày trong tuần, với bộ giá
// trị input đã chọn; thông báo Telegram / Discord / webhook khi xong hoặc lỗi. Mở từ popup "N vòng" (ScheduleRow).
import { useState } from 'react';
import Segmented from './Segmented';
import { Switch } from './RunSettings';
import { fieldBox } from './NumInput';

export type Schedule = {
  time: string; // HH:mm
  repeat: 'once' | 'daily' | 'weekdays';
  days: number[]; // 1 = T2 … 7 = CN (repeat = weekdays)
  preset: string;
  tray: boolean; // thu xuống khay khi chờ tới giờ
  notify: 'none' | 'telegram' | 'discord' | 'webhook';
  notifyTarget: string;
  notifyOn: 'done' | 'error' | 'both';
  screenshot: boolean;
};

export const DAY_LABEL = ['T2', 'T3', 'T4', 'T5', 'T6', 'T7', 'CN'];
export const scheduleSummary = (s: Schedule | null) =>
  !s ? 'Chưa hẹn'
    : `${s.time} · ${s.repeat === 'once' ? 'một lần' : s.repeat === 'daily' ? 'hằng ngày' : s.days.map((d) => DAY_LABEL[d - 1]).join(', ')}`;

// Dòng trong popup cài đặt chạy
export function ScheduleRow({ value, onEdit, soon }: { value: Schedule | null; onEdit?: () => void; soon?: boolean }) {
  return (
    <div className="flex items-center justify-between border-t border-line pt-3">
      <span className="flex flex-col">
        Hẹn giờ chạy
        <span className="font-mono text-xs text-muted">{scheduleSummary(value)}</span>
      </span>
      <button onClick={onEdit} disabled={soon} title={soon ? 'Sắp có (bước 22)' : undefined}
        className="h-8 rounded-control border border-line px-2.5 text-[13px] enabled:hover:bg-hover disabled:opacity-50">
        {value ? 'Sửa…' : 'Đặt giờ…'}
      </button>
    </div>
  );
}

const NOTIFY_HINT = { none: '', telegram: 'Bot token và chat id, vd 123:ABC@456789', discord: 'Webhook URL của kênh Discord', webhook: 'URL nhận POST JSON (kết quả, ảnh chụp)' };

export default function ScheduleDialog({ value, presets, onSave, onRemove, onClose }: {
  value: Schedule | null; presets: string[]; onSave: (s: Schedule) => void; onRemove?: () => void; onClose: () => void;
}) {
  const [s, setS] = useState<Schedule>(value ?? {
    time: '08:00', repeat: 'daily', days: [1, 2, 3, 4, 5], preset: presets[0] ?? '', tray: true,
    notify: 'none', notifyTarget: '', notifyOn: 'both', screenshot: true,
  });
  const set = (p: Partial<Schedule>) => setS({ ...s, ...p });
  const box = `${fieldBox} text-ink outline-none`;
  return (
    <div className="fixed inset-0 z-40 flex items-start justify-center bg-black/25 px-4 pt-[8vh]" onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div role="dialog" aria-label="Hẹn giờ chạy" className="flex w-[440px] max-w-full flex-col overflow-hidden rounded-panel border border-line bg-surface text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
        <div className="flex flex-col gap-4 p-[18px]">
          <div>
            <div className="text-base font-semibold">Hẹn giờ chạy</div>
            <div className="text-[13px] text-muted">App cần đang mở (có thể thu xuống khay) và máy không ngủ.</div>
          </div>
          <div className="flex items-end gap-3">
            <label className="flex flex-col gap-1.5 text-muted">Lúc
              <input type="time" value={s.time} onChange={(e) => set({ time: e.target.value })} className={`${box} w-[110px] font-mono`} />
            </label>
            <div className="flex flex-1 flex-col gap-1.5 text-muted">Lặp lại
              <Segmented value={s.repeat} onChange={(repeat) => set({ repeat })} full
                options={[{ id: 'once', label: 'Một lần' }, { id: 'daily', label: 'Hằng ngày' }, { id: 'weekdays', label: 'Chọn ngày' }]} />
            </div>
          </div>
          {s.repeat === 'weekdays' && (
            <div className="flex gap-1">
              {DAY_LABEL.map((d, i) => {
                const on = s.days.includes(i + 1);
                return (
                  <button key={d} aria-pressed={on} onClick={() => set({ days: on ? s.days.filter((x) => x !== i + 1) : [...s.days, i + 1].sort() })}
                    className={`h-8 flex-1 rounded-control border text-[13px] ${on ? 'border-accent bg-accent-soft font-semibold text-accent-text' : 'border-line hover:bg-hover'}`}>{d}</button>
                );
              })}
            </div>
          )}
          {presets.length > 0 && (
            <label className="flex flex-col gap-1.5 text-muted">Bộ giá trị input
              <select value={s.preset} onChange={(e) => set({ preset: e.target.value })} className={box}>{presets.map((p) => <option key={p}>{p}</option>)}</select>
            </label>
          )}
          <label className="flex items-center justify-between">Thu xuống khay khi chờ tới giờ<Switch on={s.tray} onChange={(tray) => set({ tray })} label="Thu xuống khay" /></label>

          <div className="flex flex-col gap-2 border-t border-line pt-3">
            <span className="font-semibold">Thông báo</span>
            <Segmented value={s.notify} onChange={(notify) => set({ notify })} full
              options={[{ id: 'none', label: 'Không' }, { id: 'telegram', label: 'Telegram' }, { id: 'discord', label: 'Discord' }, { id: 'webhook', label: 'Webhook' }]} />
            {s.notify !== 'none' && (
              <>
                <input value={s.notifyTarget} onChange={(e) => set({ notifyTarget: e.target.value })} placeholder={NOTIFY_HINT[s.notify]} className={`${box} font-mono text-[13px]`} />
                <div className="flex items-center gap-2 text-muted">Gửi khi
                  <Segmented value={s.notifyOn} onChange={(notifyOn) => set({ notifyOn })}
                    options={[{ id: 'done', label: 'Xong' }, { id: 'error', label: 'Lỗi' }, { id: 'both', label: 'Cả hai' }]} />
                </div>
                <label className="flex items-center justify-between">Kèm ảnh chụp màn hình<Switch on={s.screenshot} onChange={(screenshot) => set({ screenshot })} label="Kèm ảnh chụp" /></label>
              </>
            )}
          </div>
        </div>
        <div className="flex items-center gap-2 border-t border-line bg-chrome px-[18px] py-3">
          {value && onRemove && <button onClick={onRemove} className="h-9 rounded-control px-3 text-danger-text hover:bg-danger-soft">Bỏ hẹn giờ</button>}
          <button onClick={onClose} className="ml-auto h-9 rounded-control border border-line bg-surface px-3.5 hover:bg-hover">Huỷ</button>
          <button onClick={() => onSave(s)} disabled={s.repeat === 'weekdays' && s.days.length === 0}
            className="btn-primary h-9 rounded-control px-4 font-semibold disabled:opacity-50">Lưu hẹn giờ</button>
        </div>
      </div>
    </div>
  );
}
