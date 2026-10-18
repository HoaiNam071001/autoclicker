import { expect, it } from 'vitest';
import { fmtClip } from './clipboard';

it('fmtClip', () => {
  expect(fmtClip({ text: '' })).toBe('đọc chữ vào output.text');
  expect(fmtClip({ clipOp: 'write', text: 'Mã {{n2.output.text}}' })).toBe('ghi “Mã {{n2.output.text}}”');
  expect(fmtClip({ clipOp: 'write', text: '' })).toBe('xoá clipboard');
  expect(fmtClip({ clipOp: 'write', text: 'x'.repeat(50) })).toBe(`ghi “${'x'.repeat(40)}…”`);
  expect(fmtClip({ clipOp: 'image', text: '', filePath: '{{n3.output.path}}' })).toBe('ảnh {{n3.output.path}}');
  expect(fmtClip({ clipOp: 'image', text: '' })).toBe('chưa nhập file ảnh');
});
