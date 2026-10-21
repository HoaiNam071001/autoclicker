import { describe, expect, it } from 'vitest';
import { makeStep, type Step } from '../types';
import { chain } from './edit';
import { coordSources, nodeVarOptions, nodeViews, outputsOf } from './view';

describe('nodeViews', () => {
  const steps = [
    makeStep({ id: 'a', ref: 'n1', type: 'waitImage', findBy: 'text', findText: 'Gửi', timeoutMs: 3000, retry: { count: 2 } }),
    makeStep({ id: 'b', ref: 'n2', type: 'clickImage', imageSource: 'node', imageFrom: 'a' }),
    makeStep({ id: 'c', ref: 'n3', type: 'clickImage', imageSource: 'node' }),
  ];
  const sc = { steps, start: 'a', edges: chain(['a', 'b']) };

  it('cổng có nhãn + màu, chip thử lại / input, nguồn chỗ thấy theo mã node', () => {
    const [a, b] = nodeViews(sc, { inputs: (id) => (id === 'a' ? 1 : 0) });
    expect(a.ports).toEqual([{ id: 'success', label: 'Thành công', kind: 'ok' }, { id: 'fail', label: 'Thất bại', kind: 'fail' }]);
    expect(a.chips).toEqual(['1 input', 'thử lại 2']);
    expect(a.title).toBe('Chờ chữ');
    expect(b.detail).toContain('chỗ n1 thấy');
  });

  it('cảnh báo node không đi tới được', () => {
    const c = nodeViews(sc)[2];
    expect(c.warn).toContain('Không có đường nào');
    expect(nodeViews({ ...sc, edges: chain(['a', 'c']) })[2].warn).toContain('Chưa chọn node');
  });
});

describe('nodeVarOptions', () => {
  it('output + status của node khác theo thứ tự flow, bỏ chính nó', () => {
    const steps = [makeStep({ id: 'b', ref: 'n2', type: 'type', text: 'x' }), makeStep({ id: 'a', ref: 'n1', type: 'waitImage' })];
    const keys = nodeVarOptions({ steps, start: 'a', edges: chain(['a', 'b']) }, 'b').map((v) => v.key);
    expect(keys[0]).toBe('n1.output.x');
    expect(keys).toContain('n1.status');
    expect(keys.some((k) => k.startsWith('n2.'))).toBe(false);
  });
});

describe('coordSources', () => {
  it('node có output x, y theo thứ tự flow – bỏ chính nó và node không có toạ độ', () => {
    const steps = [
      makeStep({ id: 'c', ref: 'n3', type: 'click', points: [{ x: 1, y: 1 }] }),
      makeStep({ id: 'a', ref: 'n1', type: 'waitImage', findBy: 'text', findText: 'Gửi' }),
      makeStep({ id: 't', ref: 'n2', type: 'type', text: 'x' }),
    ];
    const got = coordSources({ steps, start: 'a', edges: chain(['a', 't', 'c']) }, 'c');
    expect(got).toEqual([{ ref: 'n1', label: 'n1 · Chờ chữ' }]);
  });
});

describe('node điều khiển', () => {
  it('cổng case theo nhãn case, cổng cuối, chi tiết + cảnh báo thiếu dữ liệu', () => {
    const sw = makeStep({ id: 's', ref: 'n1', type: 'switchScreen', timeoutMs: 3000, cases: [{ key: 'a', label: 'Đăng nhập', cond: { kind: 'text', findText: 'x' } }, { key: 'b', cond: { kind: 'image' } }] });
    const iff = makeStep({ id: 'i', ref: 'n2', type: 'if', cond: { kind: 'value', left: '{{n1.output.case}}', op: 'eq', right: 'a' } });
    const end = makeStep({ id: 'e', ref: 'n3', type: 'end', endStatus: 'fail', message: 'lỗi' });
    const [a, b, c] = nodeViews({ steps: [sw, iff, end], start: 's', edges: [...chain(['s', 'i']).map((e) => ({ ...e, port: 'case:a' })), { id: 'x', from: 'i', port: 'true', to: 'e' }] });
    expect(a.ports.map((p) => p.label)).toEqual(['Đăng nhập', 'b', 'Không thấy']);
    expect(a.detail).toBe('chờ tối đa 3 s · 2 case');
    expect(a.warn).toContain('Case b: chưa có ảnh mẫu');
    expect(b.ports.map((p) => p.id)).toEqual(['true', 'false']);
    expect(b.detail).toBe('{{n1.output.case}} = a');
    expect(c.ports).toEqual([]);
    expect(c.detail).toBe('thất bại · lỗi');
  });

  it('Lặp qua danh sách: cổng, chi tiết theo nguồn, output, cảnh báo trống', () => {
    const v = (p: Partial<Step>) => nodeViews({ steps: [makeStep({ id: 'f', ref: 'n1', type: 'forEach', ...p })], start: 'f', edges: [] })[0];
    const lines = v({ value: 'a\n\nb\n' });
    expect(lines.ports.map((p) => p.label)).toEqual(['Mỗi phần tử', 'Hết danh sách']);
    expect(lines.detail).toBe('2 dòng');
    expect(lines.warn).toBeUndefined();
    expect(v({ value: '{{n3.output.json.items}}' }).detail).toBe('{{n3.output.json.items}}');
    expect(v({ listFrom: 'file', filePath: 'C:\\data\\khach.csv' }).detail).toBe('khach.csv');
    expect(v({ value: ' ' }).warn).toBe('Danh sách trống');
    expect(v({ listFrom: 'file' }).warn).toBe('Chưa chọn file CSV / Excel');
    expect(outputsOf(makeStep({ id: 'f', type: 'forEach' })).map((o) => o.key).slice(0, 3)).toEqual(['item', 'index', 'count']);
  });
});

describe('outputsOf (F5c)', () => {
  it('mọi node có startedAt / durationMs; node phím / chờ / cuộn có output riêng', () => {
    const keys = (s: Parameters<typeof outputsOf>[0]) => outputsOf(s).map((o) => o.key);
    expect(keys(makeStep({ id: 'k', type: 'key' }))).toEqual(['keys', 'window', 'windowApp', 'startedAt', 'durationMs']);
    expect(keys(makeStep({ id: 'o', type: 'open' }))).toEqual(['target', 'app', 'window', 'windowApp', 'newWindow', 'startedAt', 'durationMs']);
    expect(keys(makeStep({ id: 'w', type: 'wait' }))).toEqual(['ms', 'startedAt', 'durationMs']);
    expect(keys(makeStep({ id: 's', type: 'scroll' }))).toEqual(['x', 'y', 'direction', 'amount', 'startedAt', 'durationMs']);
  });
  it('thẻ node: số từ biến hiện biểu thức', () => {
    const s = makeStep({ id: 'a', ref: 'n1', type: 'scroll', direction: 'down', amount: 3, numVars: [{ field: 'amount', expr: '{{vars.n}}' }, { field: 'delayAfter', expr: '{{n2.output.ms}}' }] });
    expect(nodeViews({ steps: [s], start: 'a', edges: [] })[0].detail).toBe('↓ {{vars.n}} nấc tại chuột · {{n2.output.ms}}');
  });
});
