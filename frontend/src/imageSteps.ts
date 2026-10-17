// Chờ ảnh / click theo ảnh / chờ pixel (bước 25) – phía giao diện: khoanh ảnh mẫu trên overlay, thử tìm trên màn hình,
// lấy màu pixel, quyền chụp màn hình (macOS). Engine + so khớp ở Go (imagestep.go, match.go, imageapi.go).
import { useSyncExternalStore } from 'react';
import { CaptureTemplate, FindImage, PixelColorAt, RequestScreenCapture, ScreenCaptureAllowed } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { main } from '../wailsjs/go/models';
import { addStep, getEditor, selectStep, setNotice, toScreen, updateStep } from './store';
import { errorText } from './engineStore';
import { inWails, openOverlay, setWindowMode } from './windowMode';
import { Rect } from './types';

// fromRun: chỗ giống nhất của lần chạy hết giờ không thấy (engine:failed), không phải từ nút Thử tìm
export type FindResult = { found: boolean; score: number; x: number; y: number; w: number; h: number; area: Rect; ms: number; read?: string[]; fromRun?: boolean };

type ImageState = {
  allowed: boolean | null; // quyền chụp màn hình (null = chưa kiểm tra)
  finds: Record<string, FindResult>; // kết quả "Thử tìm" theo id bước – overlay vẽ khung
  busy: string; // id bước đang thử tìm / lấy màu
  error: string;
};
let state: ImageState = { allowed: null, finds: {}, busy: '', error: '' };
const listeners = new Set<() => void>();
const set = (p: Partial<ImageState>) => {
  state = { ...state, ...p };
  listeners.forEach((l) => l());
};
export const useImageState = () => useSyncExternalStore((l) => {
  listeners.add(l);
  return () => listeners.delete(l);
}, () => state);

export async function checkCapture() {
  if (!inWails) return set({ allowed: true });
  set({ allowed: await ScreenCaptureAllowed().catch(() => true) });
}
export async function requestCapture() {
  if (inWails) set({ allowed: await RequestScreenCapture().catch(() => false) });
}

// ---- Khoanh ảnh mẫu ----
let captureTarget: string | null = null; // bước nhận ảnh mẫu; null = tạo bước Chờ ảnh mới
export const captureTargetStep = () => captureTarget;

export function startCapture(stepId: string | null) {
  captureTarget = stepId;
  return openOverlay('capture');
}

// Bảng thêm bước: Chờ ảnh → thêm bước rồi khoanh ảnh luôn.
export function addImageStep(type: 'waitImage' | 'clickImage' | 'waitPixel') {
  addStep(type);
  if (type === 'waitImage') startCapture(getEditor().selected);
}

// Bảng thêm bước: Chờ chữ / Click theo chữ = bước chờ ảnh / click theo ảnh (tự tìm) tìm theo chữ – không khoanh ảnh.
export function addTextStep(type: 'waitImage' | 'clickImage') {
  addStep(type);
  updateStep(getEditor().selected, { findBy: 'text', findText: '', ...(type === 'clickImage' ? { imageSource: 'search' as const } : {}) });
}

// "Dùng ảnh này" trên overlay: chụp vùng r (toạ độ màn hình thật) – Go ẩn cửa sổ lúc chụp – ghi vào bước, về màn chính.
export async function applyCapture(r: Rect): Promise<string | null> {
  if (!inWails) return 'Chỉ chụp được trong app';
  try {
    const shot = await CaptureTemplate(main.Rect.createFrom(r));
    const target = captureTarget && getEditor().script.steps.some((s) => s.id === captureTarget) ? captureTarget : null;
    if (!target) addStep('waitImage');
    const id = target ?? getEditor().selected;
    updateStep(id, { template: shot.png, templateScale: shot.scale });
    const { [id]: _, ...finds } = state.finds;
    set({ finds });
    setNotice(`Đã lưu ảnh mẫu ${shot.w} × ${shot.h}`);
    captureTarget = null;
    await setWindowMode('main');
    return null;
  } catch (e) {
    return errorText(e);
  }
}

// ---- Thử tìm / lấy màu ----
export async function testFind(stepId: string) {
  if (!inWails) return;
  set({ busy: stepId, error: '' });
  try {
    const r = await FindImage(main.Script.createFrom(getEditor().script), stepId);
    set({ finds: { ...state.finds, [stepId]: r as FindResult }, busy: '' });
  } catch (e) {
    set({ busy: '', error: errorText(e) });
  }
}

export async function pickColor(stepId: string) {
  const { script } = getEditor();
  const s = script.steps.find((x) => x.id === stepId);
  if (!s?.points[0] || !inWails) return;
  set({ busy: stepId, error: '' });
  try {
    const p = toScreen(s.points[0], script.region);
    const color = await PixelColorAt(p.x, p.y);
    updateStep(stepId, { color });
    setNotice(`Màu tại X ${p.x} Y ${p.y}: ${color}`);
    set({ busy: '' });
  } catch (e) {
    set({ busy: '', error: errorText(e) });
  }
}

// Ảnh mẫu → URL hiển thị trong <img>
export const templateUrl = (png?: string) => (png ? `data:image/png;base64,${png}` : '');

// Lần chạy dừng vì lỗi: chọn bước lỗi để người dùng thấy ngay; chờ ảnh hết giờ thì nhớ chỗ giống nhất (xem trên màn hình).
if (inWails) {
  EventsOn('engine:failed', (f: { stepId: string; miss?: FindResult }) => {
    selectStep(f.stepId);
    if (f.miss) set({ finds: { ...state.finds, [f.stepId]: { ...f.miss, found: false, fromRun: true } } });
  });
}
