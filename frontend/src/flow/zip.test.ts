import { describe, expect, it } from 'vitest';
import { fmtZip } from './zip';

describe('fmtZip', () => {
  it('nén', () => {
    expect(fmtZip({})).toBe('chưa nhập nguồn → nén');
    expect(fmtZip({ filePath: '~/Documents/bao-cao/' })).toBe('bao-cao → nén');
    expect(fmtZip({ filePath: 'C:\\data\\bao-cao', filePattern: ' *.pdf ' })).toBe('bao-cao → nén (*.pdf)');
  });
  it('giải nén', () => {
    expect(fmtZip({ zipOp: 'extract', filePath: '~/Downloads/pack.zip' })).toBe('pack.zip → giải nén');
  });
});
