import { describe, expect, it } from 'vitest';
import { aiLevel, levelLabel, nodePref, nodePrefLabel, runMode, withLevel } from './aiLevel';

describe('aiLevel', () => {
  it('đọc mức từ cài đặt MCP', () => {
    expect(aiLevel({ enabled: false, run: 'auto' })).toBe('off');
    expect(aiLevel({ enabled: true, run: 'auto' })).toBe('auto');
    expect(aiLevel({ enabled: true, run: '' })).toBe('ask');
    expect(runMode({ enabled: true, run: 'lạ' })).toBe('ask');
  });
  it('chọn mức: Tắt chỉ tắt MCP, giữ run; mức khác bật lại', () => {
    const m = { enabled: true, run: 'auto', maxMinutes: 5 };
    expect(withLevel(m, 'off')).toEqual({ enabled: false, run: 'auto', maxMinutes: 5 });
    expect(withLevel({ ...m, enabled: false }, 'edit')).toEqual({ enabled: true, run: 'edit', maxMinutes: 5 });
    expect(levelLabel('ask')).toBe('Hỏi');
  });
});

describe('nodePref', () => {
  it('trống / lạ = ảnh / chữ', () => {
    expect(nodePref({})).toBe('vision');
    expect(nodePref({ nodePref: 'dom' })).toBe('vision');
    expect(nodePref({ nodePref: 'web' })).toBe('web');
  });
  it('nhãn', () => expect(nodePrefLabel('web')).toBe('Trình duyệt app'));
});
