import { describe, expect, it } from 'vitest';
import { addConnected, connect, disconnect, Edge, reachable, removeNode } from './flowEdit';

const e = (id: string, from: string, port: string, to: string): Edge => ({ id, from, port, to });

describe('flowEdit', () => {
  it('connect: mỗi cổng một cạnh, nối lại thì thay, giữ id', () => {
    let edges = connect([], 'a', 'success', 'b');
    expect(edges).toHaveLength(1);
    const id = edges[0].id;
    edges = connect(edges, 'a', 'success', 'c');
    expect(edges).toEqual([{ id, from: 'a', port: 'success', to: 'c' }]);
    edges = connect(edges, 'a', 'fail', 'b');
    expect(edges).toHaveLength(2);
  });

  it('connect: không nối node vào chính nó', () => {
    expect(connect([], 'a', 'success', 'a')).toEqual([]);
  });

  it('disconnect chỉ bỏ cạnh của đúng cổng', () => {
    const edges = [e('1', 'a', 'success', 'b'), e('2', 'a', 'fail', 'c')];
    expect(disconnect(edges, 'a', 'fail')).toEqual([edges[0]]);
  });

  it('removeNode bỏ cạnh vào / ra, chuyển start', () => {
    const g = { nodes: [{ id: 'a' }, { id: 'b' }, { id: 'c' }], edges: [e('1', 'a', 'success', 'b'), e('2', 'b', 'success', 'c')], start: 'a' };
    const r = removeNode(g, 'b');
    expect(r.edges).toEqual([]);
    expect(removeNode(g, 'a').start).toBe('b');
  });

  it('addConnected chèn vào giữa cạnh có sẵn', () => {
    const g = { nodes: [{ id: 'a' }, { id: 'b' }], edges: [e('1', 'a', 'success', 'b')] };
    const r = addConnected(g, { id: 'n' }, 'a', 'success', 'success');
    expect(r.edges.map((x) => `${x.from}.${x.port}>${x.to}`).sort()).toEqual(['a.success>n', 'n.success>b']);
    // cổng chưa nối: chỉ thêm một cạnh
    const r2 = addConnected(g, { id: 'm' }, 'b', 'fail', 'success');
    expect(r2.edges.map((x) => `${x.from}.${x.port}>${x.to}`).sort()).toEqual(['a.success>b', 'b.fail>m']);
  });

  it('reachable đi theo cả vòng lặp', () => {
    const edges = [e('1', 'a', 'success', 'b'), e('2', 'b', 'fail', 'a'), e('3', 'c', 'success', 'a')];
    expect([...reachable(edges, 'a')].sort()).toEqual(['a', 'b']);
  });
});
