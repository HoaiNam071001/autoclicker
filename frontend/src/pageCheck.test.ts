import { describe, expect, it } from 'vitest';
import { normalizeUrl, verdictInfo } from './pageCheck';

describe('verdictInfo', () => {
  it('theo kết luận', () => {
    expect(verdictInfo('ok').tone).toBe('ok');
    expect(verdictInfo('captcha').title).toBe('Có captcha');
    expect(verdictInfo('blocked').tone).toBe('danger');
    expect(verdictInfo('loading').title).toBe('Trang chưa tải xong');
  });
  it('giá trị lạ = trống', () => expect(verdictInfo('???')).toBe(verdictInfo('empty')));
});

describe('normalizeUrl', () => {
  it('thêm https:// khi thiếu', () => expect(normalizeUrl(' game-a.vn/play ')).toBe('https://game-a.vn/play'));
  it('giữ giao thức có sẵn', () => {
    expect(normalizeUrl('http://localhost:8080')).toBe('http://localhost:8080');
    expect(normalizeUrl('file:///C:/a.html')).toBe('file:///C:/a.html');
  });
  it('trống', () => expect(normalizeUrl('  ')).toBe(''));
});
