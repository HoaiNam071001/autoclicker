import { expect, it } from 'vitest';
import { fmtDisplay } from './display';

it('fmtDisplay', () => {
  expect(fmtDisplay({})).toBe('đọc sáng / tối, độ sáng, loa, micro');
  expect(fmtDisplay({ displayOp: 'theme', displayValue: 'toggle' })).toBe('chế độ đảo sáng ↔ tối');
  expect(fmtDisplay({ displayOp: 'brightness', displayValue: '70' })).toBe('độ sáng 70 %');
  expect(fmtDisplay({ displayOp: 'audioOut', displayValue: ' AirPods ' })).toBe('loa → AirPods');
  expect(fmtDisplay({ displayOp: 'audioIn' })).toBe('micro → chưa nhập');
  expect(fmtDisplay({ displayOp: 'sound', displayValue: 'Glass' })).toBe('phát Glass');
});
