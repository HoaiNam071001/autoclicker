import { describe, expect, it } from 'vitest';
import { fmtProcess } from './process';

describe('fmtProcess', () => {
  it('đọc', () => {
    expect(fmtProcess({})).toBe('mọi tiến trình → đọc CPU / RAM');
    expect(fmtProcess({ procName: 'chrome' })).toBe('chrome → đọc CPU / RAM');
  });
  it('chờ', () => {
    expect(fmtProcess({ procOp: 'waitStart' })).toBe('chưa nhập tên → chờ chạy');
    expect(fmtProcess({ procOp: 'waitQuit', procName: 'zalo' })).toBe('zalo → chờ thoát');
    expect(fmtProcess({ procOp: 'waitIdle', procName: 'ffmpeg' })).toBe('ffmpeg → chờ hết bận (CPU < 5 %)');
    expect(fmtProcess({ procOp: 'waitIdle', procName: 'ffmpeg', cpuBelow: 20 })).toBe('ffmpeg → chờ hết bận (CPU < 20 %)');
  });
});
