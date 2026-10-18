import { describe, expect, it } from 'vitest';
import { fmtCountdown, fmtDuration, parseAt, untilText } from './waitUntil';
// Cùng dữ liệu với internal/script/at_test.go
import data from '../../../testdata/wait_until.json';

describe('waitUntil', () => {
  const now = new Date(data.now.y, data.now.m - 1, data.now.d, 15, 0, 0);
  it.each(data.cases)('parseAt($at)', (c) => {
    const p = parseAt(c.at, now);
    expect(!!p).toBe(c.ok);
    if (p && c.ok) {
      expect(p.t.getTime()).toBe(new Date(c.y!, c.m! - 1, c.d!, c.h!, c.mi!, c.s!, c.ms!).getTime());
      expect(p.dated).toBe(c.dated);
    }
  });
  it('fmtDuration / fmtCountdown', () => {
    expect(fmtDuration(12500)).toBe('12,5 giây');
    expect(fmtDuration(200000)).toBe('3 phút 20 giây');
    expect(fmtDuration(3900000)).toBe('1 giờ 5 phút');
    expect(fmtCountdown(7170000)).toBe('01:59:30');
    expect(fmtCountdown(245200)).toBe('04:05.2');
  });
  it('untilText', () => {
    expect(untilText({ at: '15:00:30' }, now)).toEqual({ text: 'Còn 30 giây', bad: false });
    expect(untilText({ at: '14:00', ifPassed: 'tomorrow' }, now)?.text).toContain('ngày mai (còn 23 giờ 0 phút)');
    expect(untilText({ at: '14:00', ifPassed: 'fail' }, now)?.bad).toBe(true);
    expect(untilText({ at: '{{gio}}' }, now)).toBeNull();
    expect(untilText({ at: '12h' }, now)?.bad).toBe(true);
  });
});
