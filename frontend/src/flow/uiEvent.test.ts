import { expect, it } from 'vitest';
import { fmtWaitUIEvent, uiEventOf } from './uiEvent';
import type { Step } from '../types';

const step = (p: Partial<Step>) => ({ id: 's1', ref: 'n1', type: 'waitUIEvent', timeoutMs: 10000, ...p }) as Step;

it('fmtWaitUIEvent', () => {
  expect(fmtWaitUIEvent(step({}))).toBe('cửa sổ mới mọi cửa sổ · tối đa 10 s');
  expect(fmtWaitUIEvent(step({ uiEvent: 'focus', winApp: 'Mail', winTitle: 'Hộp thư' }))).toBe('đổi sang Mail · “Hộp thư” · tối đa 10 s');
  expect(fmtWaitUIEvent(step({ uiEvent: 'newElement', selectors: ['role=dialog'], elementApp: 'Safari' }))).toBe('phần tử mới “role=dialog” · Safari · tối đa 10 s');
  expect(fmtWaitUIEvent(step({ uiEvent: 'valueChange', selectors: ['text=Tiến độ'] }))).toContain('giá trị đổi');
});

it('uiEventOf mặc định newWindow', () => {
  expect(uiEventOf(step({})).id).toBe('newWindow');
  expect(uiEventOf(step({ uiEvent: 'newElement' })).target).toBe('element');
});
