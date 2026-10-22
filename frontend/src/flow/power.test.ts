import { describe, expect, it } from 'vitest';
import { fmtPower, powerNeedsConsent } from './power';

describe('fmtPower', () => {
  it('không đếm ngược với khoá màn, có với hành động mất phiên', () => {
    expect(fmtPower({})).toBe('Khoá màn hình');
    expect(fmtPower({ powerOp: 'displaySleep', countdownSec: 5 })).toBe('Tắt màn hình sau 5 giây');
    expect(fmtPower({ powerOp: 'shutdown' })).toBe('Tắt máy sau 30 giây');
    expect(fmtPower({ powerOp: 'restart', countdownSec: 90 })).toBe('Khởi động lại sau 90 giây');
  });
  it('hành động cần đồng ý', () => {
    expect(powerNeedsConsent('lock')).toBe(false);
    expect(powerNeedsConsent('sleep')).toBe(true);
  });
});
