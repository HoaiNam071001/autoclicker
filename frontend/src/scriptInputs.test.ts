import { beforeEach, describe, expect, it, vi } from 'vitest';
import { makeStep, Script } from './types';
import { chain } from './flow/edit';

// store cần window (trình duyệt) – thay bằng store giả chỉ giữ kịch bản.
const store = vi.hoisted(() => ({ script: null as unknown as Script }));
vi.mock('./store', () => ({
  getEditor: () => ({ script: store.script }),
  updateScript: (p: Partial<Script>) => {
    store.script = { ...store.script, ...p };
  },
}));

const {
  addInput, inputUsage, parseValue, pinField, pinnedInput, pinToInput, runSections, setPinnedValue, stepInputCount, unpinField, updateInput,
} = await import('./scriptInputs');

const base = (): Script => ({
  version: 4, name: 'K', region: null, loops: 1, maxMinutes: 0, start: 'a', edges: chain(['a', 'b', 'c']),
  steps: [
    makeStep({ id: 'a', ref: 'n1', type: 'click', points: [{ x: 1, y: 1 }], delayAfter: 500 }),
    makeStep({ id: 'b', ref: 'n2', type: 'type', text: 'xin chào {{M1}}' }),
    makeStep({ id: 'c', ref: 'n3', type: 'click', points: [{ x: 2, y: 2 }], delayAfter: 800 }),
  ],
});

beforeEach(() => {
  store.script = base();
});

describe('ghim trường', () => {
  it('tạo input với giá trị hiện tại làm mặc định', () => {
    const key = pinField('a', 'delayAfter');
    expect(key).toBe('delay_1');
    expect(pinnedInput(store.script, 'a', 'delayAfter')).toMatchObject({ type: 'ms', default: 500, targets: [{ stepId: 'a', field: 'delayAfter' }] });
  });

  it('dùng chung một input cho nhiều bước – trường mới nhận mặc định, sửa một chỗ đổi cả hai', () => {
    pinField('a', 'delayAfter');
    pinToInput('c', 'delayAfter', 'delay_1');
    expect(store.script.steps[2].delayAfter).toBe(500);
    setPinnedValue('c', 'delayAfter', 1000);
    expect(store.script.steps.map((s) => s.delayAfter)).toEqual([1000, 300, 1000]);
    expect(store.script.inputs?.[0].default).toBe(1000);
  });

  it('bỏ ghim trường cuối thì xoá input', () => {
    pinField('a', 'delayAfter');
    unpinField('a', 'delayAfter');
    expect(store.script.inputs).toBeUndefined();
  });
});

describe('sửa input', () => {
  it('đổi key → đổi {{key}} trong ô chữ', () => {
    addInput({ key: 'M1', type: 'text' });
    expect(updateInput('M1', { key: 'ten' })).toBeNull();
    expect(store.script.steps[1].text).toBe('xin chào {{ten}}');
  });

  it('key sai / trùng báo lỗi, không đổi', () => {
    addInput({ key: 'M1' });
    addInput({ key: 'M2' });
    expect(updateInput('M1', { key: 'tên' })).toMatch(/không dấu/);
    expect(updateInput('M1', { key: 'M2' })).toMatch(/Đã có/);
    expect(store.script.steps[1].text).toBe('xin chào {{M1}}');
  });

  it('mặc định từ ô nhập (chuỗi) đổi sang số và ghi vào trường ghim', () => {
    pinField('a', 'delayAfter');
    updateInput('delay_1', { default: '750' });
    expect(store.script.inputs?.[0].default).toBe(750);
    expect(store.script.steps[0].delayAfter).toBe(750);
  });

  it('đổi sang kiểu chữ thì bỏ ghim vào trường số', () => {
    pinField('a', 'delayAfter');
    updateInput('delay_1', { type: 'text' });
    expect(store.script.inputs?.[0].targets).toEqual([]);
  });
});

describe('bước dùng input', () => {
  it('đếm cả ghim lẫn {{biến}}, bỏ \\{{', () => {
    addInput({ key: 'M1' });
    pinField('a', 'delayAfter');
    store.script.steps[2] = makeStep({ id: 'c', type: 'type', text: '\\{{M1}} không phải biến' });
    const usage = inputUsage(store.script);
    expect(usage.get('M1')?.map((s) => s.id)).toEqual(['b']);
    expect(stepInputCount(store.script, 'a')).toBe(1);
    expect(stepInputCount(store.script, 'c')).toBe(0);
  });

  it('form Chạy nhóm theo bước, có câu mẫu cho bước gõ chữ', () => {
    addInput({ key: 'M1' });
    pinField('a', 'delayAfter');
    pinToInput('c', 'delayAfter', 'delay_1');
    addInput({ key: 'thua' });
    expect(runSections(store.script)).toEqual([
      { steps: 'n1, n3', title: 'Click', keys: ['delay_1'], template: undefined },
      { steps: 'n2', title: 'Gõ chữ', keys: ['M1'], template: 'xin chào {{M1}}' },
      { steps: '—', title: 'Chưa bước nào dùng', keys: ['thua'], template: undefined },
    ]);
  });
});

it('parseValue', () => {
  expect(parseValue('ms', '120')).toBe(120);
  expect(parseValue('ms', 'abc')).toBeUndefined();
  expect(parseValue('text', '')).toBeUndefined();
  expect(parseValue('bool', 'true')).toBe(true);
  expect(parseValue('text', 5)).toBe('5');
});

it('biến trong URL của bước Mở: tính là dùng, đổi key đổi theo', () => {
  store.script.steps.push(makeStep({ id: 'o', type: 'open', target: 'youtube.com/results?search_query={{bai_hat}}' }));
  addInput({ key: 'M1' });
  addInput({ key: 'bai_hat' });
  expect(inputUsage(store.script).get('bai_hat')?.map((s) => s.id)).toEqual(['o']);
  updateInput('bai_hat', { key: 'song' });
  expect(store.script.steps[3].target).toBe('youtube.com/results?search_query={{song}}');
});
