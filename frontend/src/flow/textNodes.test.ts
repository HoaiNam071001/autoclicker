import { describe, expect, it } from 'vitest';
import { makeStep } from '../types';
import { textDetail, textOutputs, textPortLabel, textWarning } from './textNodes';

describe('textNodes', () => {
  it('dòng chi tiết', () => {
    expect(textDetail(makeStep({ id: 'a', type: 'textMatch', pattern: '(\\d+)', matchAll: true }))).toBe('/(\\d+)/ · mọi chỗ');
    expect(textDetail(makeStep({ id: 'a', type: 'textSplit' }))).toBe('tách theo mỗi dòng');
    expect(textDetail(makeStep({ id: 'a', type: 'textClean', stripAccents: true, caseMode: 'lower', trim: true }))).toBe('bỏ dấu · thường · gọn khoảng trắng');
    expect(textDetail(makeStep({ id: 'a', type: 'textClean' }))).toBe('chưa chọn việc');
    expect(textDetail(makeStep({ id: 'a', type: 'textNumber', numOp: 'format', decimals: 2 }))).toBe('định dạng 1.234,5 · 2 số lẻ');
    expect(textDetail(makeStep({ id: 'a', type: 'textNumber', numOp: 'date' }))).toBe('ngày → dd/MM/yyyy');
  });
  it('cảnh báo: thiếu chữ vào (trừ ngày = bây giờ), thiếu mẫu', () => {
    expect(textWarning(makeStep({ id: 'a', type: 'textSplit' }))).toBe('Chưa có chữ cần xử lý');
    expect(textWarning(makeStep({ id: 'a', type: 'textNumber', numOp: 'date' }))).toBeUndefined();
    expect(textWarning(makeStep({ id: 'a', type: 'textMatch', value: '{{n1.output.text}}' }))).toBe('Chưa nhập mẫu');
  });
  it('cổng + output', () => {
    expect(textPortLabel('textMatch', 'fail')).toBe('Không khớp');
    expect(textPortLabel('textClean', 'success')).toBe('Tiếp');
    expect(textOutputs(makeStep({ id: 'a', type: 'textNumber' }))[0].key).toBe('value');
    expect(textOutputs(makeStep({ id: 'a', type: 'textNumber', numOp: 'format' }))[0].key).toBe('text');
  });
});
