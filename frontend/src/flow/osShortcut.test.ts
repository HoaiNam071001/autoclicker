import { expect, it } from 'vitest';
import { fmtOSShortcut } from './osShortcut';

it('fmtOSShortcut', () => {
  expect(fmtOSShortcut({})).toBe('(chưa chọn)');
  expect(fmtOSShortcut({ shortcutName: 'Đổi tên' })).toBe('Đổi tên');
  expect(fmtOSShortcut({ shortcutName: 'Đổi tên', text: 'a.png' })).toBe('Đổi tên · có input');
});
