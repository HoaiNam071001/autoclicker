import { describe, expect, it } from 'vitest';
import { HintState, overlayHint } from './hint';
import { makeStep } from '../types';

const base: HintState = {
  focusStep: null, repicking: false, dragA: false, tool: 'point', captured: false, captureRef: null, selected: 0, picking: false,
  placing: 'click', modKey: '⌘',
};

describe('overlayHint', () => {
  it('xem riêng node kéo thả đang chọn lại: chờ A rồi B', () => {
    const f = makeStep({ id: 'd', type: 'drag' });
    expect(overlayHint({ ...base, focusStep: f, repicking: true })).toContain('điểm A');
    expect(overlayHint({ ...base, focusStep: f, repicking: true, dragA: true })).toContain('điểm B');
    expect(overlayHint({ ...base, focusStep: f })).toBe('Kéo marker để sửa · Esc quay lại');
  });
  it('khoanh ảnh: ghi node nhận ảnh hoặc tạo node mới', () => {
    expect(overlayHint({ ...base, tool: 'capture', captureRef: 'n4' })).toContain('(n4)');
    expect(overlayHint({ ...base, tool: 'capture' })).toContain('tạo bước Chờ ảnh mới');
    expect(overlayHint({ ...base, tool: 'capture', captured: true })).toMatch(/^Enter/);
  });
  it('chọn nhiều: số điểm đã chọn', () => {
    expect(overlayHint({ ...base, tool: 'select', selected: 3 })).toMatch(/^3 điểm/);
  });
  it('đặt điểm: chọn lại hoặc thêm node theo loại', () => {
    expect(overlayHint({ ...base, picking: true })).toBe('Click để đặt lại điểm');
    expect(overlayHint(base)).toMatch(/^Click \/ F6: thêm /);
  });
});
