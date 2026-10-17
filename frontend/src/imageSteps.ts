// Chờ ảnh / click theo ảnh / chờ pixel (bước 25) – phía giao diện: khoanh ảnh mẫu trên overlay, thử tìm trên màn hình,
// lấy màu pixel, quyền chụp màn hình (macOS). Engine + so khớp ở Go (imagestep.go, match.go, imageapi.go).
import { useSyncExternalStore } from 'react';
import { CaptureTemplate, FindImage, PixelColorAt, RequestScreenCapture, ScreenCaptureAllowed, TestCond } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime/runtime';
import { main } from '../wailsjs/go/models';
import { addStep, getEditor, selectStep, setNotice, stepById, toScreen, updateStep } from './scriptStore';
import { condAt, updateCondAt } from './stepTree';
import { errorText } from './engineStore';
import { inWails, openOverlay, setWindowMode } from './windowMode';
import { Cond, Rect } from './types';

// fromRun: chỗ giống nhất của lần chạy hết giờ không thấy (engine:failed), không phải từ nút Thử tìm
export type FindResult = { found: boolean; score: number; x: number; y: number; w: number; h: number; area: Rect; ms: number; read?: string[]; fromRun?: boolean };

// Kết quả "Thử điều kiện" (condtest.go)
export type CondPart = { path: number[]; kind: string; ok: boolean; known: boolean; detail: string; hit?: Rect };
export type CondTestResult = { ok: boolean; known: boolean; parts: CondPart[]; ms: number };

type ImageState = {
  condTests: Record<string, CondTestResult>; // theo id bước
  allowed: boolean | null; // quyền chụp màn hình (null = chưa kiểm tra)
  finds: Record<string, FindResult>; // kết quả "Thử tìm" theo id bước – overlay vẽ khung
  busy: string; // id bước đang thử tìm / lấy màu
  error: string;
};
let state: ImageState = { condTests: {}, allowed: null, finds: {}, busy: '', error: '' };
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
let captureCond: number[] | null = null; // ảnh mẫu cho điều kiện của khối (đường dẫn điều kiện con), null = của bước
export const captureTargetStep = () => captureTarget;

export function startCapture(stepId: string | null, condPath: number[] | null = null) {
  captureTarget = stepId;
  captureCond = condPath;
  return openOverlay('capture');
}

// Sửa điều kiện con tại path của khối id.
export function updateCond(id: string, path: number[], f: (c: Cond) => Cond) {
  const s = stepById(id);
  if (s?.cond) updateStep(id, { cond: updateCondAt(s.cond, path, f) });
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
    if (captureCond && captureTarget && stepById(captureTarget)) {
      updateCond(captureTarget, captureCond, (c) => ({ ...c, template: shot.png, templateScale: shot.scale }));
      const { [captureTarget]: _, ...condTests } = state.condTests;
      set({ condTests });
      setNotice(`Đã lưu ảnh mẫu ${shot.w} × ${shot.h} cho điều kiện`);
      captureTarget = captureCond = null;
      await setWindowMode('main');
      return null;
    }
    const target = captureTarget && stepById(captureTarget) ? captureTarget : null;
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

// "Thử điều kiện": xét điều kiện trên màn hình lúc này; chỗ thấy đầu tiên vẽ khung trên overlay (như Thử tìm).
export async function testCond(stepId: string) {
  if (!inWails) return;
  set({ busy: stepId, error: '' });
  try {
    const r = (await TestCond(main.Script.createFrom(getEditor().script), stepId, [])) as CondTestResult;
    const hit = r.parts.find((p) => p.hit && p.ok)?.hit ?? r.parts.find((p) => p.hit)?.hit;
    const finds = { ...state.finds };
    if (hit) finds[stepId] = { found: r.ok, score: 100, ...hit, area: hit, ms: r.ms };
    else delete finds[stepId];
    set({ condTests: { ...state.condTests, [stepId]: r }, finds, busy: '' });
  } catch (e) {
    set({ busy: '', error: errorText(e) });
  }
}

// "Lấy màu" cho điều kiện màu pixel của khối.
export async function pickCondColor(stepId: string, path: number[]) {
  const { script } = getEditor();
  const c = stepById(stepId)?.cond && condAt(stepById(stepId)!.cond!, path);
  if (!c?.points?.[0] || !inWails) return;
  set({ busy: stepId, error: '' });
  try {
    const p = toScreen(c.points[0], script.region);
    const color = await PixelColorAt(p.x, p.y);
    updateCond(stepId, path, (x) => ({ ...x, color }));
    setNotice(`Màu tại X ${p.x} Y ${p.y}: ${color}`);
    set({ busy: '' });
  } catch (e) {
    set({ busy: '', error: errorText(e) });
  }
}

export async function pickColor(stepId: string) {
  const { script } = getEditor();
  const s = stepById(stepId);
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
