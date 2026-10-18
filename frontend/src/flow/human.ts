// Thao tác giống người (bước 34) – nhãn mức, chữ giải thích, node nào áp dụng. Công thức chạy ở internal/humanize (Go).
import { isWebStep } from './web';
import type { HumanLevel, Step } from '../types';

export const HUMAN_LEVELS: { id: HumanLevel; label: string }[] = [
  { id: 'off', label: 'Tắt' }, { id: 'low', label: 'Thấp' }, { id: 'medium', label: 'Vừa' }, { id: 'high', label: 'Cao' },
];

export const humanLabel = (lv: HumanLevel) => HUMAN_LEVELS.find((l) => l.id === lv)?.label ?? lv;

// Mức kịch bản: trống = tắt (file cũ); mức của node: trống = theo kịch bản.
export const scriptHuman = (lv?: HumanLevel): HumanLevel => lv || 'off';
export const effectiveHuman = (script?: HumanLevel, node?: HumanLevel): HumanLevel => node || scriptHuman(script);

// Mô tả mỗi mức (khớp profiles trong internal/humanize/profile.go).
export function humanHint(lv: HumanLevel): string {
  switch (lv) {
    case 'low': return 'Chuột đi đường hơi cong, nhanh; click giữ 40–80 ms; gõ / cuộn nhịp hơi lệch';
    case 'medium': return 'Chuột đi đường cong, tăng rồi giảm tốc, đôi khi đi lố rồi chỉnh; click giữ 60–150 ms, bấm lệch quanh tâm; gõ / cuộn nhịp không đều';
    case 'high': return 'Như Vừa nhưng chậm, cong hơn, hay ngừng; đôi khi gõ nhầm rồi xoá';
    default: return 'Chuột nhảy thẳng tới đích, click / gõ / cuộn nhanh nhất';
  }
}

// Node điều khiển chuột / phím thật – có ô ghi đè mức (node web chạy nền không đụng chuột / phím).
export function usesHuman(s: Pick<Step, 'type' | 'webBackground' | 'scrollMode'>): boolean {
  if (isWebStep(s.type)) {
    if (s.type === 'webWait' || s.type === 'webRead' || s.webBackground) return false;
    return s.type !== 'webScroll' || s.scrollMode === 'in' || s.scrollMode === 'until';
  }
  return ['click', 'double', 'hold', 'drag', 'type', 'key', 'scroll', 'clickImage'].includes(s.type);
}
