import { describe, expect, it } from 'vitest';
import { lastRowOf, runHead, runRows, runTrace } from './run';
import type { Edge, StepResult } from '../types';

const edges: Edge[] = [
  { id: 'e1', from: 'w', port: 'success', to: 'c' },
  { id: 'e2', from: 'w', port: 'fail', to: 's' },
  { id: 'e3', from: 's', port: 'success', to: 'w' },
];
const steps = [{ id: 'w', ref: 'n1' }, { id: 'c', ref: 'n2' }, { id: 's', ref: 'n3' }];
const res = (nodeId: string, status: StepResult['status'], port: string, loop = 1): StepResult =>
  ({ nodeId, ref: '', status, port, attempts: 1, loop, startedAt: '2026-10-17T14:32:05+07:00', durationMs: 1200 });

// w thất bại → cuộn → w thấy → c
const results = [res('w', 'fail', 'fail'), res('s', 'success', 'success'), res('w', 'success', 'success'), res('c', 'success', 'success')];

describe('runTrace', () => {
  it('trạng thái lần cuối, số lần đi qua, cạnh đã đi', () => {
    const t = runTrace(results, edges);
    expect(t.nodes.get('w')).toEqual({ run: 'ok', visits: 2 });
    expect(t.nodes.get('c')).toEqual({ run: 'ok', visits: 1 });
    expect([...t.taken].sort()).toEqual(['e1', 'e2', 'e3']);
    expect(t.active).toBeUndefined();
  });
  it('đang chạy: node hiện đang chạy, cạnh vừa đi sáng', () => {
    const t = runTrace(results.slice(0, 2), edges, 'w');
    expect(t.nodes.get('w')).toEqual({ run: 'running', visits: 1 });
    expect(t.active).toBe('e3');
  });
  it('sang vòng mới không tính là đi cạnh', () => {
    const t = runTrace([res('c', 'success', 'success', 1), res('w', 'success', 'success', 2)], [{ id: 'x', from: 'c', port: 'success', to: 'w' }]);
    expect(t.taken.size).toBe(0);
  });
});

describe('runRows', () => {
  it('cột đi tiếp + log', () => {
    const rows = runRows([...results.slice(0, 1), { ...res('s', 'fail', 'fail'), error: { code: 'internal', message: 'x' }, logs: [{ atMs: 1500, level: 'error', msg: 'x' }] }], steps, edges);
    expect(rows[0].next).toBe('n3');
    expect(rows[1].next).toBe('dừng');
    expect(rows[1].logs).toEqual([{ t: '+1,5', level: 'error', msg: 'x' }]);
    expect(rows[0].at).toMatch(/^\d\d:\d\d:05$/);
    expect(runRows([res('c', 'success', 'success')], steps, edges)[0].next).toBe('hết');
    expect(lastRowOf(runRows(results, steps, edges), 'w')?.status).toBe('success');
  });
});

describe('runHead', () => {
  it('thành công có lời nhắn Kết thúc / thất bại / đang chạy', () => {
    const base = { durationMs: 3200, nodes: 4, startedAt: '2026-10-17T14:32:05+07:00' };
    expect(runHead(results, { ...base, status: 'success', message: 'Đã gửi' }, true, false)).toMatchObject({ summary: 'Thành công · Đã gửi · 4 node · 3,2 s', tone: 'ok' });
    expect(runHead(results, { ...base, status: 'fail', error: 'n1: không thấy' }, false, false)).toMatchObject({ summary: 'Thất bại · n1: không thấy', tone: 'fail' });
    expect(runHead(results, null, true, true).tone).toBe('running');
  });
});
