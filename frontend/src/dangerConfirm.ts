// Hỏi xác nhận thao tác nguy hiểm trước khi chạy (xoá file, buộc thoát app, đóng cửa sổ – flow/danger.ts): runScript gọi
// askDanger, DangerConfirmDialog hiện + trả lời. Go nhớ đồng ý theo kịch bản (dangerapi.go) – lần sau không hỏi, đổi
// đích (thư mục, mẫu tên, app…) thì hỏi lại.
import { useSyncExternalStore } from 'react';
import type { DangerOp } from './flow/danger';

type Ask = { ops: DangerOp[]; resolve: (ok: boolean) => void } | null;
let ask: Ask = null;
const listeners = new Set<() => void>();
const set = (a: Ask) => {
  ask = a;
  listeners.forEach((l) => l());
};

export const askDanger = (ops: DangerOp[]) => new Promise<boolean>((resolve) => set({ ops, resolve }));

export function answerDanger(ok: boolean) {
  ask?.resolve(ok);
  set(null);
}

export const useDangerAsk = () => useSyncExternalStore((l) => {
  listeners.add(l);
  return () => listeners.delete(l);
}, () => ask);
