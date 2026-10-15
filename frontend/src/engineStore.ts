import { useSyncExternalStore } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { inWails } from './windowMode';

// Khớp struct Progress trong engine.go (chỉ đi qua event nên Wails không sinh type).
export type Progress = {
  loop: number;
  loops: number; // 0 = vô hạn
  step: number;
  total: number;
  clicks: number;
  elapsedMs: number; // không tính thời gian tạm dừng
  paused: boolean;
};

// lastError: lỗi của lần chạy gần nhất – giữ ở store vì màn gọi RunScript bị gỡ khi chuyển sang mini bar.
// stopReason: lý do lần dừng gần nhất ('f8' | 'corner'), '' nếu chạy hết / dừng bằng nút.
type EngineState = {
  running: boolean;
  progress: Progress | null;
  receivedAt: number;
  lastError: string;
  stopReason: string;
};

// Store chung nghe event từ Go ngay khi load, để màn mở sau (mini bar) vẫn có dữ liệu mới nhất.
let state: EngineState = { running: false, progress: null, receivedAt: 0, lastError: '', stopReason: '' };
const listeners = new Set<() => void>();
const set = (patch: Partial<EngineState>) => {
  state = { ...state, ...patch };
  listeners.forEach((l) => l());
};

if (inWails) {
  EventsOn('engine:state', (running: boolean) => set(running ? { running, progress: null, lastError: '', stopReason: '' } : { running }));
  EventsOn('engine:stopped', (stopReason: string) => set({ stopReason }));
  EventsOn('engine:progress', (progress: Progress) => set({ progress, receivedAt: Date.now() }));
}

export const setLastError = (lastError: string) => set({ lastError });

const subscribe = (l: () => void) => {
  listeners.add(l);
  return () => listeners.delete(l);
};

export const useEngine = () => useSyncExternalStore(subscribe, () => state);

// Thời gian đã chạy tính tới hiện tại (Go chỉ gửi khi có thay đổi).
export const elapsedNow = ({ progress, receivedAt }: EngineState) =>
  progress ? progress.elapsedMs + (progress.paused ? 0 : Date.now() - receivedAt) : 0;

export const fmtElapsed = (ms: number) => {
  const s = Math.floor(ms / 1000);
  const mm = String(Math.floor(s / 60) % 60).padStart(2, '0');
  const ss = String(s % 60).padStart(2, '0');
  return s >= 3600 ? `${Math.floor(s / 3600)}:${mm}:${ss}` : `${mm}:${ss}`;
};

export const STOP_REASON: Record<string, string> = {
  f8: 'Đã dừng bằng F8',
  corner: 'Dừng khẩn cấp: chuột chạm góc màn hình',
};
