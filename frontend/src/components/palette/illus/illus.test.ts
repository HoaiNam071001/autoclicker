import { describe, expect, it } from 'vitest';
import { ILLUS } from '.';
import { NODE_INFO } from '../../../flow/nodeInfo';

describe('hoạt hoạ minh hoạ node', () => {
  it('mọi loại node có mô tả đều có hoạt hoạ', () => {
    for (const k of Object.keys(NODE_INFO)) expect(ILLUS[k as keyof typeof NODE_INFO], k).toBeTypeOf('function');
  });
});
