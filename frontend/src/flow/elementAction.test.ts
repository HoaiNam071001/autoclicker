import { expect, it } from 'vitest';
import { fmtElementAction } from './elementAction';
import type { Step } from '../types';

const step = (p: Partial<Step>) => ({ timeoutMs: 10000, selectors: ['role=button[name="Gửi"]'], ...p }) as Step;

it('fmtElementAction', () => {
  expect(fmtElementAction(step({}))).toBe('bấm “role=button[name="Gửi"]” · app trên cùng · tối đa 10 s');
  expect(fmtElementAction(step({ elementAct: 'setValue', text: 'An', elementApp: 'Safari' }))).toBe('đặt giá trị “role=button[name="Gửi"]” = “An” · Safari · tối đa 10 s');
  expect(fmtElementAction(step({ elementAct: 'uncheck' }))).toContain('bỏ tick');
});
