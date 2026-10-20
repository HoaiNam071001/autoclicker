// Chữ cho hộp hỏi AI chạy kịch bản (bước 31b) – logic thuần, aiAsk.ts / AIRunDialog dùng.

export type AIAsk = { id: string; script: string; path: string; steps: number; from?: string; stopAt?: string; nodes?: string[]; resume?: boolean };

// Mô tả phần sẽ chạy: cả kịch bản / từ n3 / đến n5 / cụm n2, n3.
export function aiRunScope(a: AIAsk): string {
  const parts = [a.nodes?.length ? `chỉ ${a.nodes.join(', ')}` : '', a.from ? `từ ${a.from}` : '', a.stopAt ? `đến ${a.stopAt}` : ''].filter(Boolean);
  const scope = parts.length ? parts.join(' ') : `cả kịch bản (${a.steps} node)`;
  return a.resume ? `${scope}, dùng kết quả lần chạy trước` : scope;
}
