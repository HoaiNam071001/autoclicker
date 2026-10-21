import { describe, expect, it } from 'vitest';
import { calcPreview, calcQuote, checkCalc, evalCalc } from './calc';
// Cùng dữ liệu với internal/script/calc_test.go
import data from '../../../testdata/calc.json';

describe('calc', () => {
  it.each(data.eval)('evalCalc($expr)', (c) => {
    if (c.error) expect(() => evalCalc(c.expr)).toThrow();
    else expect(evalCalc(c.expr)).toBe(c.value);
  });
  it.each(data.quote)('calcQuote($in)', (c) => expect(calcQuote(c.in)).toBe(c.out));
  it('checkCalc: chỉ báo lỗi cú pháp, biến coi là số', () => {
    expect(checkCalc('{{n2.output.count}} * 2 + 1')).toBeUndefined();
    expect(checkCalc('10 / {{vars.n}}')).toBeUndefined();
    expect(checkCalc('{{n2.output.count}} *')).toBeTruthy();
    expect(checkCalc('round({{price}}, 2, 3)')).toBeTruthy();
  });
  it('calcPreview', () => {
    expect(calcPreview('2 * 3')).toEqual({ value: '6' });
    expect(calcPreview('{{n1.output.value}} > 3')).toEqual({ needsRun: true });
    expect(calcPreview('1 / 0').error).toBe('chia cho 0');
  });
});
