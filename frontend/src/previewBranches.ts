// Nhánh xem trước đã chọn (menu Nhánh trên thanh Xem trước – bước 20 F6): nhớ trong phiên để mở lại overlay không mất,
// không lưu vào file kịch bản (chỉ là giả sử lúc xem trước).
import { useSyncExternalStore } from 'react';
import type { BranchChoices } from './flow/preview';

let choices: BranchChoices = {};
const listeners = new Set<() => void>();
const set = (c: BranchChoices) => {
  choices = c;
  listeners.forEach((l) => l());
};

export const chooseBranch = (id: string, port: string) => set({ ...choices, [id]: port });
export const resetBranches = () => set({});

export const useBranchChoices = () => useSyncExternalStore((l) => {
  listeners.add(l);
  return () => listeners.delete(l);
}, () => choices);
