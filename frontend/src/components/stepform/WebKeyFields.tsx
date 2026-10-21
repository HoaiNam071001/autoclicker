// Node Phím trên trang web (46 B3): phím như node Phím tắt + tab của Chrome của app – gửi qua trình duyệt, không đụng
// bàn phím thật.
import { Step } from '../../types';
import { KeyFields } from './KeyFields';
import { WebTabField } from './WebTabField';
import { hint } from './shared';

export function WebKeyFields({ step: s }: { step: Step }) {
  return <>
    <KeyFields step={s} />
    <WebTabField step={s} />
    <span className={`col-span-2 -mt-1 ${hint}`}>
      Phím vào ô đang chọn của trang trên Chrome của app – Chrome bị che vẫn chạy, bàn phím thật không gõ gì. Dùng được: enter, tab, esc, backspace, delete, space, mũi tên, home, end, pageup, pagedown, f1–f12, chữ / số.
    </span>
  </>;
}
