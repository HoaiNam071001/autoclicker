import { describe, expect, it } from 'vitest';
import { Edge, makeStep, Step } from '../types';
import { branchNodes, BranchChoices, hopText, newWalk, previewHop, routeText, stopText } from './preview';

const e = (from: string, port: string, to: string): Edge => ({ id: `${from}-${port}`, from, port, to });

// Đi hết một lượt xem trước, trả đường đã đi + chữ lúc dừng.
function walk(steps: Step[], edges: Edge[], choices: BranchChoices = {}) {
  const sc = { steps, start: steps[0].id, edges };
  const w = newWalk();
  const route: string[] = [];
  let cur = sc.start;
  let last;
  while (cur) {
    const s = steps.find((x) => x.id === cur)!;
    last = { h: previewHop(sc, s, choices, w), s };
    route.push(hopText(last.h, s));
    cur = last.h.next;
  }
  return { route, stop: stopText(last?.h, last?.s) };
}

// Ví dụ Facebook: tìm nút; không thấy → Đếm → cuộn → tìm lại; hết lượt → Kết thúc báo lỗi.
const fb = [
  makeStep({ id: 'w', ref: 'n1', type: 'waitImage' }),
  makeStep({ id: 'c', ref: 'n2', type: 'click' }),
  makeStep({ id: 'k', ref: 'n3', type: 'counter', max: 10 }),
  makeStep({ id: 's', ref: 'n4', type: 'scroll' }),
  makeStep({ id: 'x', ref: 'n5', type: 'end', endStatus: 'fail' }),
];
const fbEdges = [e('w', 'success', 'c'), e('w', 'fail', 'k'), e('k', 'next', 's'), e('k', 'done', 'x'), e('s', 'success', 'w')];

describe('previewHop', () => {
  it('mặc định đi nhánh chính', () => {
    expect(walk(fb, fbEdges)).toEqual({ route: ['n1', 'n2'], stop: 'Xong 1 vòng' });
  });

  it('chọn Thất bại → Đếm chạy 1 lượt rồi hết lượt', () => {
    expect(walk(fb, fbEdges, { w: 'fail' })).toEqual({
      route: ['n1 ✗', 'n3 (1/10)', 'n4', 'n1 ✗', 'n3 hết lượt', 'n5'], stop: 'Kết thúc ở n5 – báo lỗi',
    });
  });

  it('Lặp qua danh sách: chọn Mỗi phần tử → 1 lượt rồi hết; chọn Hết danh sách → ra luôn', () => {
    const steps = [
      makeStep({ id: 'f', ref: 'n1', type: 'forEach', value: 'a' }),
      makeStep({ id: 't', ref: 'n2', type: 'type' }),
      makeStep({ id: 'x', ref: 'n3', type: 'end' }),
    ];
    const edges = [e('f', 'item', 't'), e('t', 'success', 'f'), e('f', 'done', 'x')];
    expect(walk(steps, edges).route).toEqual(['n1 (phần tử 1)', 'n2', 'n1 hết danh sách', 'n3']);
    expect(walk(steps, edges, { f: 'done' }).route).toEqual(['n1 hết danh sách', 'n3']);
    expect(branchNodes({ steps, start: 'f', edges }).map((b) => b.ports.map((p) => p.label))).toEqual([['Mỗi phần tử', 'Hết danh sách']]);
  });

  it('thất bại không có nhánh → dừng báo lỗi; cổng chọn không còn → về cổng đầu', () => {
    expect(walk(fb.slice(0, 2), [e('w', 'success', 'c')], { w: 'fail' }).stop).toBe('n1 thất bại, không có nhánh Thất bại – dừng báo lỗi');
    expect(walk(fb, fbEdges, { w: 'case:x' }).route).toEqual(['n1', 'n2']);
  });

  it('vòng lặp không có Đếm → dừng sau 2 lần', () => {
    const r = walk(fb.slice(0, 2), [e('w', 'fail', 'c'), e('c', 'success', 'w')], { w: 'fail' });
    expect(r.route).toEqual(['n1 ✗', 'n2', 'n1 ✗', 'n2', 'n1 ✗']);
    expect(r.stop).toBe('Dừng xem trước: lặp lại ở n1 mà không có node Đếm');
  });

  it('Nếu / Switch theo cổng chọn; node tắt đi success', () => {
    const steps = [
      makeStep({ id: 'i', ref: 'n1', type: 'if' }),
      makeStep({ id: 'v', ref: 'n2', type: 'switch', cases: [{ key: 'a', label: 'Có tin' }, { key: 'b', label: 'Hết' }] }),
      makeStep({ id: 'd', ref: 'n3', type: 'waitImage', enabled: false }),
    ];
    const edges = [e('i', 'false', 'v'), e('v', 'case:b', 'd')];
    expect(walk(steps, edges).route).toEqual(['n1 Đúng']);
    expect(walk(steps, edges, { i: 'false', v: 'case:b' }).route).toEqual(['n1 Sai', 'n2 Hết', 'n3']);
  });
});

describe('branchNodes', () => {
  it('Nếu / Switch + node nối cổng khác cổng đầu; bỏ Đếm, node tắt, node chỉ nối success', () => {
    const steps = [...fb, makeStep({ id: 'i', ref: 'n6', type: 'if' }), makeStep({ id: 'z', ref: 'n7', type: 'waitImage', enabled: false })];
    const got = branchNodes({ steps, start: 'w', edges: [...fbEdges, e('z', 'success', 'w'), e('z', 'fail', 'k')] });
    expect(got.map((b) => b.ref)).toEqual(['n1', 'n6']);
    expect(got[0].ports.map((p) => p.label)).toEqual(['Thành công', 'Thất bại']);
  });
});

it('routeText giữ 6 chặng cuối', () => {
  expect(routeText(['a', 'b'])).toBe('a → b');
  expect(routeText(['1', '2', '3', '4', '5', '6', '7'])).toBe('… → 2 → 3 → 4 → 5 → 6 → 7');
});
