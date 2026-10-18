import { describe, expect, it } from 'vitest';
import { ALL, appsOf, filterSamples, Sample, stepCounts } from './filter';

const s = (id: string, app: string, steps: string[], name = id): Sample =>
  ({ id, name, description: '', app, level: 1, page: 'click.html', real: false, steps });

const list = [
  s('a', 'Cơ bản', ['open', 'click'], 'Click nhiều lần'),
  s('b', 'Cơ bản', ['open', 'drag'], 'Kéo thả'),
  s('c', 'Điều khiển', ['open', 'if', 'click']),
];

describe('filterSamples', () => {
  it('không lọc gì → giữ cả danh sách', () => {
    expect(filterSamples(list, ALL)).toHaveLength(3);
  });
  it('lọc theo loại bước và kiểu app', () => {
    expect(filterSamples(list, { ...ALL, step: 'click' }).map((x) => x.id)).toEqual(['a', 'c']);
    expect(filterSamples(list, { ...ALL, step: 'click', app: 'Điều khiển' }).map((x) => x.id)).toEqual(['c']);
  });
  it('tìm không dấu, mọi từ phải khớp, tìm cả tên loại bước', () => {
    expect(filterSamples(list, { ...ALL, q: 'keo' }).map((x) => x.id)).toEqual(['b']);
    expect(filterSamples(list, { ...ALL, q: 'click lan' }).map((x) => x.id)).toEqual(['a']);
    expect(filterSamples(list, { ...ALL, q: 'nếu' }).map((x) => x.id)).toEqual(['c']);
  });
});

describe('stepCounts / appsOf', () => {
  it('đếm số mẫu theo loại bước, theo thứ tự bảng bước', () => {
    expect(stepCounts(list)).toEqual([
      { type: 'click', count: 2 }, { type: 'drag', count: 1 }, { type: 'open', count: 3 }, { type: 'if', count: 1 },
    ]);
  });
  it('kiểu app không trùng, giữ thứ tự', () => {
    expect(appsOf(list)).toEqual(['Cơ bản', 'Điều khiển']);
  });
});
