import { expect, it } from 'vitest';
import { fmtNotify } from './notify';

it('fmtNotify', () => {
  expect(fmtNotify({})).toBe('(chưa nhập nội dung)');
  expect(fmtNotify({ text: ' Xong ' })).toBe('Xong');
  expect(fmtNotify({ notifyTitle: 'Báo cáo', text: 'Đã nén 3 file' })).toBe('Báo cáo · Đã nén 3 file');
});
