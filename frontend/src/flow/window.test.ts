import { describe, expect, it } from 'vitest';
import { fmtWindow, matchWindow, windowHasBox, windowPreset } from './window';
// Cùng dữ liệu với internal/model/window_test.go
import wins from '../../../testdata/windows.json';
import presets from '../../../testdata/window_presets.json';

describe('window', () => {
  it.each(wins.cases)('matchWindow($app, $title, $match)', (c) => {
    const run = () => matchWindow(wins.wins, c.app, c.title, c.match);
    if (c.want === 'error') expect(run).toThrow();
    else expect(run()).toBe(c.want);
  });
  it.each(presets.cases)('windowPreset($preset)', (c) => {
    expect(windowPreset(c.preset, c.work, c.want)).toEqual(c.out);
  });
  it('fmtWindow', () => {
    expect(fmtWindow({ winApp: 'Zalo' })).toBe('Zalo → đưa lên trước (focus)');
    expect(fmtWindow({ winApp: 'chrome', winTitle: 'Sân thử', winOp: 'setBounds', winPreset: 'left', winRect: { x: 0, y: 0, w: 1280, h: 800 } }))
      .toBe('chrome · “Sân thử” → trái 1280 × 800');
    expect(fmtWindow({ winTitle: 'a', winOp: 'setBounds', winRect: { x: 10, y: 20, w: 800, h: 600 }, display: 2 })).toBe('“a” → 10, 20 800 × 600 · màn 2');
    expect(fmtWindow({ winOp: 'setBounds', winPreset: 'full' })).toBe('chưa chọn cửa sổ → đầy màn');
  });
  it('windowHasBox – khớp HasFoundBox (flowcheck.go)', () => {
    expect(windowHasBox('')).toBe(true);
    expect(windowHasBox('close')).toBe(false);
    expect(windowHasBox('waitClose')).toBe(false);
  });
});
