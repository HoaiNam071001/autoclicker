import { describe, expect, it } from 'vitest';
import { portsSummary, retrySummary, timingSummary } from './summary';
import type { PortView } from '../components/flow/model';

describe('tóm tắt khối inspector', () => {
  it('thời gian', () => {
    expect(timingSummary(500)).toBe('chờ 500 ms, cố định');
    expect(timingSummary(500, 100)).toBe('chờ 400–600 ms');
    expect(timingSummary(1500, 500)).toBe('chờ 1 s–2 s');
    expect(timingSummary(0)).toBe('không chờ');
    expect(timingSummary('{{gap}}')).toBe('chờ {{gap}}');
  });

  it('thử lại', () => {
    expect(retrySummary(undefined)).toBe('không thử lại');
    expect(retrySummary({ count: 2, delayMs: 500 })).toBe('thử lại 2 lần · nghỉ 500 ms');
  });

  it('cổng ra', () => {
    const ok: PortView = { id: 'success', label: 'Thành công', kind: 'ok' };
    const fail: PortView = { id: 'fail', label: 'Thất bại', kind: 'fail' };
    const edges = [{ id: 'e1', from: 'a', port: 'success', to: 'b' }];
    const ref = (id: string) => ({ b: 'n4' })[id] ?? '?';
    expect(portsSummary('a', [ok, fail], edges, ref)).toBe('ok → n4 · lỗi → dừng');
    expect(portsSummary('a', [], edges, ref)).toBe('không có – kết thúc ở đây');
    const n: PortView = { id: 'none', label: 'Không thấy', kind: 'neutral' };
    expect(portsSummary('a', [ok, fail, n], edges, ref)).toBe('3 cổng · 1 đã nối');
  });
});
