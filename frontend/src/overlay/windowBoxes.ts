// Khung đích của node Cửa sổ "Đặt vị trí + cỡ" (bước 26 W3b) để vẽ trên overlay – cùng cách engine tính
// (windowTarget, internal/engine/window.go): toạ độ → đúng winRect; preset → windowPreset trong vùng làm việc của màn.
import { windowPreset, windowTarget } from '../flow/window';
import type { Rect, Step } from '../types';

export type ScreenInfo = { rect: Rect; work: Rect; current: boolean };
// guess: màn 0 (màn đang chứa cửa sổ) chỉ biết lúc chạy – vẽ tạm trên màn đang mở overlay.
export type WindowBox = { id: string; ref: string; rect: Rect; who: string; guess: boolean };

export function windowBoxes(steps: Step[], screens: ScreenInfo[]): WindowBox[] {
  const out: WindowBox[] = [];
  for (const s of steps) {
    if (s.type !== 'window' || s.winOp !== 'setBounds') continue;
    const want = s.winRect ?? { x: 0, y: 0, w: 0, h: 0 };
    const base = { id: s.id, ref: s.ref, who: windowTarget(s), guess: false };
    if (!s.winPreset) {
      if (want.w > 0 && want.h > 0) out.push({ ...base, rect: want });
      continue;
    }
    const scr = s.display ? screens[s.display - 1] : screens.find((d) => d.current) ?? screens[0];
    if (!scr) continue; // không có màn (màn số lớn hơn số màn máy này) – chạy thật sẽ báo lỗi
    out.push({ ...base, rect: windowPreset(s.winPreset, scr.work, want), guess: !s.display });
  }
  return out;
}
