import { expect, it } from 'vitest';
import { fmtShell, shellPreview } from './shell';

it('shellPreview lấy dòng đầu không trống, cắt ngắn', () => {
  expect(shellPreview('\n  echo hi \nls')).toBe('echo hi');
  expect(shellPreview('x'.repeat(60))).toBe(`${'x'.repeat(47)}…`);
  expect(shellPreview(undefined)).toBe('');
});

it('fmtShell', () => {
  expect(fmtShell({ script: 'echo hi' })).toBe('echo hi');
  expect(fmtShell({ shellLang: 'applescript', script: 'beep' })).toBe('AppleScript: beep');
  expect(fmtShell({})).toBe('(chưa có lệnh)');
});
