import { describe, expect, it } from 'vitest';
import { flowOrder } from './flow';
import type { Script } from './types';
import data from '../../testdata/flow_order.json';

// Cùng dữ liệu với TestFlowOrderShared (flowrun_test.go) – Go và TS phải đánh số node giống nhau.
const cases = data as unknown as { name: string; script: Pick<Script, 'start' | 'steps' | 'edges'>; order: string[] }[];

describe('flowOrder (chung với Go)', () => {
  it('có dữ liệu', () => expect(cases.length).toBeGreaterThan(0));
  for (const c of cases) it(c.name, () => expect(flowOrder(c.script)).toEqual(c.order));
});
