// Kịch bản mẫu khi mở giao diện trên trình duyệt (không có Go) – xem canvas / overlay / inspector (store/editor.ts).
// Cuối có vòng rẽ nhánh (Chờ chữ thất bại → Đếm → Cuộn → tìm lại) để xem menu Nhánh của Xem trước.
import { chain } from '../flow/edit';
import { Edge, makeStep, Step } from '../types';

export const browserSteps: Step[] = [
  makeStep({ id: '1', ref: 'n1', type: 'click', points: [{ x: 412, y: 220 }], delayAfter: 500 }),
  makeStep({ id: '2', ref: 'n2', type: 'double', points: [{ x: 640, y: 220 }] }),
  makeStep({ id: '3', ref: 'n3', type: 'click', points: [{ x: 812, y: 436 }], clicks: 5, clickInterval: 120, delayAfter: 800 }),
  makeStep({ id: '4', ref: 'n4', type: 'drag', points: [{ x: 300, y: 600 }, { x: 900, y: 600 }], dragMs: 400, delayAfter: 1000 }),
  makeStep({ id: '5', ref: 'n5', type: 'wait', clicks: 0, delayAfter: 2000, enabled: false }),
  makeStep({ id: '6', ref: 'n6', type: 'click', points: [{ x: 1100, y: 520 }], button: 'right', delayAfter: 500 }),
  makeStep({ id: '7', ref: 'n7', type: 'type', clicks: 0, text: 'xin chao 123' }),
  makeStep({ id: '8', ref: 'n8', type: 'key', clicks: 0, keys: ['enter'], delayAfter: 500 }),
  makeStep({ id: '9', ref: 'n9', type: 'waitImage', name: 'Tìm nút Nhắn tin', findBy: 'text', findText: 'Nhắn tin', timeoutMs: 3000 }),
  makeStep({ id: '10', ref: 'n10', type: 'counter', max: 5 }),
  makeStep({ id: '11', ref: 'n11', type: 'scroll', amount: 3 }),
];

export const browserEdges: Edge[] = [
  ...chain(browserSteps.slice(0, 9).map((s) => s.id)),
  { id: 'b1', from: '9', port: 'fail', to: '10' },
  { id: 'b2', from: '10', port: 'next', to: '11' },
  { id: 'b3', from: '11', port: 'success', to: '9' },
];
