// Node Thao tác phần tử (bước 43 E1, internal/engine/elementaction.go) – nhãn thao tác, dòng chi tiết thẻ.
import { fmtMs, type Step } from '../types';
import { firstSelector } from './element';

// Khớp model.ElementActs ('' = press)
export type ElementAct = 'press' | 'setValue' | 'check' | 'uncheck' | 'expand' | 'collapse' | 'select' | 'scroll';
export const ELEMENT_ACTS: { id: ElementAct; label: string; verb: string }[] = [
  { id: 'press', label: 'Bấm', verb: 'bấm' }, { id: 'setValue', label: 'Đặt giá trị', verb: 'đặt giá trị' },
  { id: 'check', label: 'Tick', verb: 'tick' }, { id: 'uncheck', label: 'Bỏ tick', verb: 'bỏ tick' },
  { id: 'expand', label: 'Mở rộng', verb: 'mở rộng' }, { id: 'collapse', label: 'Thu gọn', verb: 'thu gọn' },
  { id: 'select', label: 'Chọn mục', verb: 'chọn' }, { id: 'scroll', label: 'Cuộn tới', verb: 'cuộn tới' },
];

// Dòng chi tiết: bấm “role=button[name=Gửi]” · Safari · tối đa 10 s
export function fmtElementAction(s: Step): string {
  const verb = ELEMENT_ACTS.find((a) => a.id === (s.elementAct || 'press'))?.verb ?? 'bấm';
  const value = s.elementAct === 'setValue' ? ` = “${s.text ?? ''}”` : '';
  return `${verb} ${firstSelector(s.selectors)}${value} · ${s.elementApp?.trim() || 'app trên cùng'} · tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
}
