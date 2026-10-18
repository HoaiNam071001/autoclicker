// Hỏi xác nhận thao tác xoá trước khi chạy (bước 35 N4): runScript gọi askDelete, DeleteConfirmDialog hiện + trả lời.
// Go nhớ đồng ý theo kịch bản (deleteapi.go) – lần sau không hỏi, đổi thư mục / mẫu tên thì hỏi lại.
import { useSyncExternalStore } from 'react';
import type { DeleteOp } from './flow/file';

type Ask = { ops: DeleteOp[]; resolve: (ok: boolean) => void } | null;
let ask: Ask = null;
const listeners = new Set<() => void>();
const set = (a: Ask) => {
  ask = a;
  listeners.forEach((l) => l());
};

export const askDelete = (ops: DeleteOp[]) => new Promise<boolean>((resolve) => set({ ops, resolve }));

export function answerDelete(ok: boolean) {
  ask?.resolve(ok);
  set(null);
}

export const useDeleteAsk = () => useSyncExternalStore((l) => {
  listeners.add(l);
  return () => listeners.delete(l);
}, () => ask);
