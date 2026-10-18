import { describe, expect, it } from 'vitest';
import { controls, curveAt, fitts, HUMAN_PROFILES, humanLegs } from './humanPath';
import type { HumanLevel } from '../types';
// Cùng dữ liệu với internal/humanize/shared_test.go – xem trước vẽ đúng công thức đường chuột thật.
import data from '../../../testdata/human_path.json';

// nguồn ngẫu nhiên lặp lại được
const seeded = (seed: number) => () => ((seed = (seed * 1103515245 + 12345) % 2 ** 31) / 2 ** 31);

describe('humanPath – dữ liệu chung với Go', () => {
  it('thông số mức khớp profile.go', () => {
    expect(HUMAN_PROFILES).toEqual(data.profiles);
  });
  it.each(data.fitts)('fitts($d, $size) = $ms', (f) => {
    expect(fitts(f.d, f.size)).toBeCloseTo(f.ms, 2);
  });
  it.each(data.curves)('controls + curveAt %#', (c) => {
    const [c1, c2] = controls(c.a, c.b, c.o1, c.o2);
    expect(c1.x).toBeCloseTo(c.c1.x, 2); expect(c1.y).toBeCloseTo(c.c1.y, 2);
    expect(c2.x).toBeCloseTo(c.c2.x, 2); expect(c2.y).toBeCloseTo(c.c2.y, 2);
    for (const a of c.at) {
      const p = curveAt(c.a, c1, c2, c.b, a.t);
      expect(p.x).toBeCloseTo(a.p.x, 2); expect(p.y).toBeCloseTo(a.p.y, 2);
    }
  });
});

describe('humanLegs', () => {
  const from = { x: 100, y: 100 }, to = { x: 900, y: 500 };
  it.each(['low', 'medium', 'high'] as HumanLevel[])('mức %s: nối liền, tới đúng đích, thời gian hợp lý', (lv) => {
    for (let seed = 1; seed < 200; seed++) {
      const legs = humanLegs(from, to, HUMAN_PROFILES[lv]!, {}, seeded(seed));
      expect(legs[0].from).toEqual(from);
      expect(legs.at(-1)!.to).toEqual(to);
      legs.slice(1).forEach((l, i) => expect(l.from).toEqual(legs[i].to));
      const ms = legs.reduce((n, l) => n + l.ms + l.pauseMs, 0);
      expect(ms).toBeGreaterThan(fitts(894, 0) * HUMAN_PROFILES[lv]!.speed * 0.8);
      expect(ms).toBeLessThan(fitts(894, 0) * HUMAN_PROFILES[lv]!.speed * 1.2 + 650);
    }
  });
  it('có lúc đi lố rồi chỉnh; kéo thả (durMs) và quãng ngắn thì không', () => {
    const p = HUMAN_PROFILES.high!;
    const counts = Array.from({ length: 300 }, (_, i) => humanLegs(from, to, p, {}, seeded(i + 1)).length);
    expect(counts).toContain(2);
    for (let i = 1; i < 100; i++) {
      expect(humanLegs(from, to, p, { durMs: 500 }, seeded(i))).toMatchObject([{ ms: 500 }]);
      expect(humanLegs(from, { x: 150, y: 120 }, p, {}, seeded(i))).toHaveLength(1);
    }
  });
  it('đường cong thật, không thẳng', () => {
    const [l] = humanLegs(from, to, HUMAN_PROFILES.medium!, { durMs: 400 }, () => 0.95);
    const mid = curveAt(l.from, l.c1, l.c2, l.to, 0.5);
    const off = Math.abs((to.x - from.x) * (mid.y - from.y) - (to.y - from.y) * (mid.x - from.x)) / Math.hypot(to.x - from.x, to.y - from.y);
    expect(off).toBeGreaterThan(20);
  });
  it('quãng < 2 px: nhảy luôn', () => {
    expect(humanLegs(from, { x: 101, y: 100 }, HUMAN_PROFILES.medium!)).toMatchObject([{ ms: 0, to: { x: 101, y: 100 } }]);
  });
});
