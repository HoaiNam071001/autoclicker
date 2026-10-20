import { describe, expect, it } from 'vitest';
import { windowBoxes, type ScreenInfo } from './windowBoxes';
import type { Step } from '../types';

const screens: ScreenInfo[] = [
  { rect: { x: 0, y: 0, w: 1920, h: 1080 }, work: { x: 0, y: 0, w: 1920, h: 1040 }, current: false },
  { rect: { x: 1920, y: 0, w: 1280, h: 800 }, work: { x: 1920, y: 0, w: 1280, h: 760 }, current: true },
];
const win = (p: Partial<Step>) => ({ id: 's1', ref: 'n1', type: 'window', winOp: 'setBounds', winApp: 'chrome', ...p }) as Step;

describe('windowBoxes', () => {
  it('toạ độ: đúng winRect, thiếu cỡ thì bỏ', () => {
    expect(windowBoxes([win({ winRect: { x: 10, y: 20, w: 800, h: 600 } })], screens))
      .toEqual([{ id: 's1', ref: 'n1', who: 'chrome', rect: { x: 10, y: 20, w: 800, h: 600 }, guess: false }]);
    expect(windowBoxes([win({ winRect: { x: 10, y: 20, w: 0, h: 600 } })], screens)).toEqual([]);
  });
  it('preset theo vùng làm việc của màn đã chọn', () => {
    const [b] = windowBoxes([win({ winPreset: 'right', winRect: { x: 0, y: 0, w: 1000, h: 0 }, display: 1 })], screens);
    expect(b.rect).toEqual({ x: 920, y: 0, w: 1000, h: 1040 });
    expect(b.guess).toBe(false);
  });
  it('màn 0: vẽ tạm trên màn đang mở overlay', () => {
    const [b] = windowBoxes([win({ winPreset: 'full' })], screens);
    expect(b.rect).toEqual(screens[1].work);
    expect(b.guess).toBe(true);
  });
  it('bỏ node khác và màn không có', () => {
    expect(windowBoxes([win({ winOp: '' }), win({ winPreset: 'full', display: 3 }), { ...win({}), type: 'click' } as Step], screens)).toEqual([]);
  });
});
