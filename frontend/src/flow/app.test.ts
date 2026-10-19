import { describe, expect, it } from 'vitest';
import { appChoices, appStatus, fmtApp, matchApps } from './app';
// Cùng dữ liệu với internal/model/app_test.go
import data from '../../../testdata/apps.json';

describe('app', () => {
  it.each(data.cases)('matchApps($app, $match)', (c) => {
    const run = () => matchApps(data.apps, c.app, c.match);
    if (c.want === 'error') expect(run).toThrow();
    else expect(run()).toEqual(c.want);
  });
  it('appStatus', () => {
    expect(appStatus(data.apps, 'chrome', '')).toEqual({ running: true, windows: 2 });
    expect(appStatus(data.apps, 'telegram', '')).toEqual({ running: false, windows: 0 });
    expect(appStatus(data.apps, '{{ten}}', '')).toBeNull();
    expect(appStatus(data.apps, '(', 'regex')).toBeNull();
  });
  it('appChoices – gộp tiến trình cùng tên, app có cửa sổ trước', () => {
    expect(appChoices(data.apps).map((a) => `${a.name}:${a.windows}`)).toEqual(['chrome:2', 'Notepad:1', 'TextEdit:1', 'Zalo:0']);
  });
  it('fmtApp', () => {
    expect(fmtApp({ appName: 'Zalo' })).toBe('Zalo → kiểm tra đang chạy');
    expect(fmtApp({ appName: 'notepad.exe', appOp: 'quit', appForce: true })).toBe('notepad.exe → thoát, quá giờ thì buộc thoát');
    expect(fmtApp({ appOp: 'waitStart' })).toBe('chưa nhập app → chờ app chạy');
  });
});
