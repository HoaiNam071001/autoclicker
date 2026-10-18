import { describe, expect, it } from 'vitest';
import { fmtSelectors, parseSelector, Selector } from './web';
import data from '../../../testdata/selectors.json';

// Cùng dữ liệu với TestParseSelector (internal/model/web_test.go) – Go và TS phải đọc selector giống nhau.
const cases = data as { in: string; want?: Selector; bad?: boolean }[];

describe('parseSelector (chung với Go)', () => {
  it('có dữ liệu', () => expect(cases.length).toBeGreaterThan(0));
  for (const c of cases) it(JSON.stringify(c.in), () => expect(parseSelector(c.in)).toEqual(c.bad ? null : c.want));
});

describe('fmtSelectors', () => {
  it('cách đầu + số dự phòng, bỏ dòng trống', () => {
    expect(fmtSelectors(['text=Xem lịch sử', '', 'css=#h'])).toBe('chữ “Xem lịch sử” (+1 dự phòng)');
    expect(fmtSelectors(['css=#go'])).toBe('CSS #go');
    expect(fmtSelectors(['role=button[name="Gửi"]'])).toBe('vai trò button “Gửi”');
    expect(fmtSelectors(['role=link'])).toBe('vai trò link');
  });
  it('chưa nhập', () => {
    expect(fmtSelectors(undefined)).toBe('chưa nhập cách tìm phần tử');
    expect(fmtSelectors(['text='])).toBe('chưa nhập cách tìm phần tử');
  });
});
