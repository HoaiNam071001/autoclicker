// Nhật ký lần chạy đang xem (bước 20 F5): kết quả từng node của lần đang chạy / vừa chạy (event engine:node) hoặc
// một lần trong lịch sử (runlog.go). Canvas tô đường đã đi, bảng Nhật ký, mục "Lần chạy gần nhất" đều đọc từ đây.
import { useSyncExternalStore } from 'react';
import { GetRun, ListRuns, RunShot } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { getEditor, subscribe } from './store';
import { inWails } from './windowMode';
import type { StepResult } from './types';

// Khớp RunSummary trong runlog.go.
export type RunSummary = {
  id: string; script: string; source: string; startedAt: string; durationMs: number;
  status: 'success' | 'fail' | 'stopped'; error?: string; message?: string; nodes: number;
};

const MAX_RESULTS = 2000; // khớp maxRunResults (runlog.go)

type RunLogState = {
  results: StepResult[];
  summary: RunSummary | null; // null = đang chạy / chưa chạy
  live: boolean; // true: lần chạy trong phiên này; false: mở từ lịch sử
  open: boolean; // bảng Nhật ký đang hiện
  history: RunSummary[]; // lịch sử của kịch bản đang mở, mới nhất trước
  // kịch bản (tên + file) của lần đang xem – để đọc ảnh lúc lỗi đúng thư mục lịch sử
  owner: { name: string; path: string };
};

let state: RunLogState = { results: [], summary: null, live: true, open: false, history: [], owner: { name: '', path: '' } };
const listeners = new Set<() => void>();
const set = (patch: Partial<RunLogState>) => {
  state = { ...state, ...patch };
  listeners.forEach((l) => l());
};
export const useRunLog = () => useSyncExternalStore((l) => {
  listeners.add(l);
  return () => listeners.delete(l);
}, () => state);

const ownerNow = () => {
  const { script, path } = getEditor();
  return { name: script.name, path: path ?? '' };
};

// Lịch sử chạy của kịch bản đang mở (gọi khi mở kịch bản khác / chạy xong).
export async function loadHistory() {
  if (!inWails) return;
  const { name, path } = ownerNow();
  const history = ((await ListRuns(name, path).catch(() => [])) ?? []) as unknown as RunSummary[];
  set({ history });
}

export async function openRun(id: string) {
  const owner = ownerNow();
  const rec = await GetRun(owner.name, owner.path, id).catch(() => null);
  if (!rec) return;
  const { results, inputs: _, truncated: __, ...summary } = rec as unknown as RunSummary & { results: StepResult[]; inputs: unknown; truncated: number };
  set({ results: results ?? [], summary: summary as RunSummary, live: false, open: true, owner });
}

export const closeRunLog = () => set({ open: false });

export const shotUrl = (shot: string) =>
  state.summary ? RunShot(state.owner.name, state.owner.path, state.summary.id, shot).catch(() => '') : Promise.resolve('');

// Mở kịch bản khác: bỏ nhật ký của kịch bản cũ (id node không còn nghĩa), nạp lịch sử của kịch bản mới.
// Lưu thành file khác / đổi tên: lịch sử theo file mới.
let seen = { opened: -1, key: '' };
function watchEditor() {
  const { opened, path, script } = getEditor();
  const key = `${path ?? ''}\n${script.name}`;
  if (opened === seen.opened && key === seen.key) return;
  if (opened !== seen.opened) set({ results: [], summary: null, live: true, open: false });
  seen = { opened, key };
  loadHistory();
}

if (inWails) {
  subscribe(watchEditor);
  EventsOn('engine:state', (running: boolean) => {
    if (running) set({ results: [], summary: null, live: true, open: true, owner: ownerNow() });
  });
  EventsOn('engine:node', (r: StepResult) => {
    const results = state.results.length >= MAX_RESULTS ? state.results.slice(1) : state.results.slice();
    results.push(r);
    set({ results });
  });
  EventsOn('engine:done', (summary: RunSummary) => {
    set({ summary });
    loadHistory();
  });
}
