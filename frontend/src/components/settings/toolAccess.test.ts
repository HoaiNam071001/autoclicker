import { describe, expect, it } from 'vitest';
import { clearOverride, clientChoices, hasOverride, offFor, setToolOn } from './toolAccess';

const base = { offTools: ['screenshot'], clients: { cursor: ['run_script'] } };

describe('mcpTools', () => {
  it('ghi đè riêng thay hẳn mặc định', () => {
    expect(offFor(base, '')).toEqual(['screenshot']);
    expect(offFor(base, 'cursor')).toEqual(['run_script']);
    expect(offFor(base, 'claude-code')).toEqual(['screenshot']); // chưa ghi đè → mặc định
    expect(hasOverride(base, 'cursor')).toBe(true);
    expect(hasOverride(base, '')).toBe(false);
  });
  it('sửa mặc định / sửa client', () => {
    expect(setToolOn(base, '', 'run_script', false)).toEqual({ offTools: ['screenshot', 'run_script'] });
    expect(setToolOn(base, '', 'screenshot', true)).toEqual({ offTools: [] });
    // client chưa ghi đè: tạo từ mặc định
    expect(setToolOn(base, 'codex', 'update_script', false).clients).toEqual({ cursor: ['run_script'], codex: ['screenshot', 'update_script'] });
  });
  it('bỏ ghi đè + danh sách client', () => {
    expect(clearOverride(base, 'cursor')).toEqual({ clients: {} });
    expect(clientChoices(base, ['codex', 'cursor'])).toEqual(['codex', 'cursor']);
  });
});
