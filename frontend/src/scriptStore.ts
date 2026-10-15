import { useSyncExternalStore } from 'react';
import { Point, Rect, Script, sampleSteps } from './types';

// Store kịch bản đang soạn – dùng chung giữa MainScreen và OverlayScreen (hai màn không cùng mount).
// Chưa lưu xuống đĩa (bước 21).
let script: Script = { name: 'Kịch bản 1', region: null, loops: 10, maxMinutes: 0, steps: sampleSteps };
const listeners = new Set<() => void>();

const subscribe = (l: () => void) => {
  listeners.add(l);
  return () => listeners.delete(l);
};

export const useScript = () => useSyncExternalStore(subscribe, () => script);

export const updateScript = (patch: Partial<Script>) => {
  script = { ...script, ...patch };
  listeners.forEach((l) => l());
};

// Đổi vùng; toạ độ điểm vẫn tương đối nên các điểm đi theo vùng mới. null = toàn màn hình.
export const setRegion = (region: Rect | null) => updateScript({ region });

// Điểm tương đối theo vùng → toạ độ màn hình thật.
export const toScreen = (p: Point, region: Rect | null): Point =>
  region ? { x: region.x + p.x, y: region.y + p.y } : p;
