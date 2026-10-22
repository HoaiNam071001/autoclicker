import { describe, expect, it } from 'vitest';
import { fmtInputSource } from './inputSource';

describe('fmtInputSource', () => {
  it('đọc / đổi', () => {
    expect(fmtInputSource({})).toBe('đọc bộ gõ đang dùng');
    expect(fmtInputSource({ inputOp: 'set', inputLayout: ' ABC ' })).toBe('đổi sang ABC');
    expect(fmtInputSource({ inputOp: 'set' })).toBe('đổi sang chưa nhập bộ gõ');
    expect(fmtInputSource({ inputOp: 'toggleIME' })).toBe('bấm ctrl+shift (đảo bật / tắt tiếng Việt)');
    expect(fmtInputSource({ inputOp: 'toggleIME', inputHotkey: ' alt+z ' })).toBe('bấm alt+z (đảo bật / tắt tiếng Việt)');
  });
});
