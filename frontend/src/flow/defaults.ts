// Giá trị mặc định khi thêm node theo loại – giống NewStep (internal/model/newstep.go). Dùng khi thêm node (store/nodes,
// store/points) và cho mục "Về mặc định" của ô nhập (bước 38 V3b). Phần phụ thuộc kịch bản (điểm giữa vùng, node lấy chỗ
// thấy gần nhất) do store ghép thêm.
import { makeStep, NumVarField, Step, StepType } from '../types';

export const STEP_DEFAULTS: Partial<Record<StepType, Partial<Step>>> = {
  click: { delayAfter: 500 },
  double: { delayAfter: 500 },
  hold: { holdMs: 1000, delayAfter: 500 },
  drag: { dragMs: 400, delayAfter: 500 },
  wait: { clicks: 0, delayAfter: 1000 },
  type: { clicks: 0, delayAfter: 500 },
  key: { clicks: 0, keys: ['enter'], delayAfter: 500 },
  scroll: { clicks: 0, direction: 'down', amount: 3, delayAfter: 500 },
  open: { clicks: 0, target: 'https://', delayAfter: 2000 },
  waitImage: { clicks: 0, threshold: 90, timeoutMs: 10000, intervalMs: 250, delayAfter: 300 },
  clickImage: { imageSource: 'node', threshold: 90, timeoutMs: 10000, intervalMs: 250, delayAfter: 500 },
  waitPixel: { clicks: 0, color: '#000000', tolerance: 10, timeoutMs: 5000, intervalMs: 100, delayAfter: 300 },
  screenshot: { clicks: 0, delayAfter: 300 },
  clipboard: { clicks: 0, delayAfter: 300 },
  path: { dragMs: 1500, delayAfter: 300 },
  waitUntil: { clicks: 0, at: '12:00:00', delayAfter: 0 },
  system: { clicks: 0, sysOp: 'volume', volume: 50, delayAfter: 300 },
  window: { clicks: 0, timeoutMs: 10000, delayAfter: 300 },
  app: { clicks: 0, timeoutMs: 10000, delayAfter: 300 },
  file: { clicks: 0, filePath: '~/Documents/ghosthand/nhat-ky.txt', delayAfter: 300 },
  http: { clicks: 0, method: 'GET', url: 'https://', timeoutMs: 15000, delayAfter: 0 },
  webWait: { clicks: 0, selectors: ['text='], timeoutMs: 10000, intervalMs: 200, delayAfter: 300 },
  webClick: { selectors: ['text='], timeoutMs: 10000, intervalMs: 200, delayAfter: 500 },
  webType: { clicks: 0, selectors: ['text='], webClear: true, timeoutMs: 10000, intervalMs: 200, delayAfter: 300 },
  webRead: { clicks: 0, selectors: ['text='], timeoutMs: 10000, intervalMs: 200, delayAfter: 0 },
  elementRead: { clicks: 0, selectors: ['text='], timeoutMs: 10000, intervalMs: 250, delayAfter: 0 },
  webScroll: { clicks: 0, selectors: ['text='], direction: 'down', amount: 3, timeoutMs: 10000, intervalMs: 200, delayAfter: 300 },
};

// Mặc định của một ô số theo loại node; undefined = loại này không có mặc định cho ô đó (ẩn "Về mặc định").
export function defaultOf(type: StepType, field: NumVarField): number | undefined {
  const v = { ...makeStep({ id: '', type }), ...STEP_DEFAULTS[type] }[field];
  return typeof v === 'number' ? v : undefined;
}
