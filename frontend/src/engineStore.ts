import { useSyncExternalStore } from 'react';
import { RunScript, ValidateScript } from '../wailsjs/go/main/App';
import { main } from '../wailsjs/go/models';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { Script } from './types';
import { inWails, setWindowMode } from './windowMode';

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
  countdown: number; // giây còn lại trước khi chạy, 0 = không đếm ngược
  hotkeyError: string; // Go không nghe được phím toàn cục (thiếu quyền)
  preset: string; // bộ giá trị input của lần chạy này ('' = tự điền / không có input) – hiện trên thanh mini
};

// Store chung nghe event từ Go ngay khi load, để màn mở sau (mini bar) vẫn có dữ liệu mới nhất.
let state: EngineState = {
  running: false, progress: null, receivedAt: 0, lastError: '', stopReason: '', countdown: 0, hotkeyError: '', preset: '',
};
const listeners = new Set<() => void>();
const set = (patch: Partial<EngineState>) => {
  state = { ...state, ...patch };
  listeners.forEach((l) => l());
};

if (inWails) {
  EventsOn('engine:state', (running: boolean) => set(running ? { running, progress: null, lastError: '', stopReason: '' } : { running, countdown: 0 }));
  EventsOn('engine:countdown', (countdown: number) => set({ countdown }));
  EventsOn('hotkey:error', (hotkeyError: string) => set({ hotkeyError }));
  EventsOn('engine:stopped', (stopReason: string) => set({ stopReason }));
  EventsOn('engine:progress', (progress: Progress) => set({ progress, receivedAt: Date.now() }));
}

export const setLastError = (lastError: string) => set({ lastError });

const COUNTDOWN_SEC = 3;

// Chạy kịch bản với giá trị input (Go áp vào bản sao – resolveInputs): kiểm tra trước (lỗi thì ở lại màn chính),
// thu cửa sổ thành mini bar, chạy tới khi xong / dừng, rồi về màn chính. Form Chạy: runForm.ts.
export async function runScript(script: Script, countdown: boolean, values: Record<string, unknown> = {}, preset = '') {
  if (state.running || !inWails) return;
  const sc = main.Script.createFrom(script);
  set({ preset });
  try {
    await ValidateScript(sc, values);
  } catch (e) {
    set({ lastError: String(e), stopReason: '' });
    return;
  }
  await setWindowMode('minibar');
  try {
    await RunScript(sc, values, countdown ? COUNTDOWN_SEC : 0);
  } catch (e) {
    set({ lastError: String(e) });
  }
  await setWindowMode('main');
}

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
