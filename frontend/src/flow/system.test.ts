import { expect, it } from 'vitest';
import { fmtSystem, sysHasLevel } from './system';

it('fmtSystem', () => {
  expect(fmtSystem({ sysOp: 'volume', volume: 30 })).toBe('âm lượng = 30 %');
  expect(fmtSystem({ sysOp: 'volumeUp' })).toBe('âm lượng +10 %');
  expect(fmtSystem({ sysOp: 'volumeDown', volume: 5 })).toBe('âm lượng −5 %');
  expect(fmtSystem({ sysOp: 'mute' })).toBe('tắt tiếng');
  expect(sysHasLevel('volumeUp') && !sysHasLevel('mute')).toBe(true);
});
