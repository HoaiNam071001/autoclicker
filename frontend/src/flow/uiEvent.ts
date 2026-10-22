// Node Chờ sự kiện giao diện (bước 43 E2, internal/engine/waituievent.go) – nhãn sự kiện, dòng chi tiết thẻ.
import { fmtMs, type Step } from '../types';
import { firstSelector } from './element';

// Khớp model.UIEvents ('' = newWindow). window: lọc theo cửa sổ (winApp / winTitle); element: theo selectors.
export type UiEvent = 'newWindow' | 'focus' | 'newElement' | 'valueChange';
export const UI_EVENTS: { id: UiEvent; label: string; target: 'window' | 'element' }[] = [
  { id: 'newWindow', label: 'Cửa sổ mới mở', target: 'window' },
  { id: 'focus', label: 'Đổi cửa sổ trên cùng', target: 'window' },
  { id: 'newElement', label: 'Phần tử mới hiện (menu, hộp thoại)', target: 'element' },
  { id: 'valueChange', label: 'Giá trị phần tử đổi', target: 'element' },
];

export function uiEventOf(s: Step): (typeof UI_EVENTS)[number] {
  return UI_EVENTS.find((e) => e.id === (s.uiEvent || 'newWindow')) ?? UI_EVENTS[0];
}

// Dòng chi tiết: cửa sổ mới Notes · “Lưu thành” · tối đa 10 s  |  phần tử mới “role=dialog” · Safari · tối đa 10 s
export function fmtWaitUIEvent(s: Step): string {
  const ev = uiEventOf(s);
  const wait = `tối đa ${fmtMs(s.timeoutMs ?? 0)}`;
  if (ev.target === 'element') {
    return `${ev.id === 'newElement' ? 'phần tử mới' : 'giá trị đổi'} ${firstSelector(s.selectors)} · ${s.elementApp?.trim() || 'app trên cùng'} · ${wait}`;
  }
  const what = [s.winApp?.trim(), s.winTitle?.trim() && `“${s.winTitle.trim()}”`].filter(Boolean).join(' · ') || 'mọi cửa sổ';
  return `${ev.id === 'newWindow' ? 'cửa sổ mới' : 'đổi sang'} ${what} · ${wait}`;
}
