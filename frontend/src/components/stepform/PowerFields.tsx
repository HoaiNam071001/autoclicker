// Node Nguồn (bước 43 S3): hành động + đếm ngược. Hành động mất phiên làm việc hỏi đồng ý lúc bấm Chạy.
import NumInput from '../NumInput';
import Select from '../Select';
import { updateStep } from '../../store';
import { POWER_COUNTDOWN_SEC, POWER_OPS, powerNeedsConsent, type PowerOp } from '../../flow/power';
import { Step } from '../../types';
import { Field, hint, selectBox } from './shared';

export default function PowerFields({ step: s }: { step: Step }) {
  const op: PowerOp = s.powerOp ?? 'lock';
  const danger = powerNeedsConsent(op);
  return (
    <>
      <Field label="Làm gì" wide group>
        <Select value={op} onChange={(v) => updateStep(s.id, { powerOp: v as PowerOp })} className={selectBox} aria-label="Hành động nguồn"
          options={POWER_OPS.map((o) => ({ value: o.id, label: o.label }))} />
      </Field>
      <Field label="Đếm ngược" help={danger ? `Đợi trước khi làm để kịp bấm Dừng. Trống / 0 = ${POWER_COUNTDOWN_SEC} giây` : 'Tuỳ chọn – 0 = làm ngay'}>
        <NumInput value={s.countdownSec ?? 0} unit="giây" min={0} onChange={(countdownSec) => updateStep(s.id, { countdownSec })} />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        {danger
          ? 'Mất phiên làm việc: lần chạy đầu app hỏi xác nhận. Trong lúc đếm ngược bấm Dừng (thanh mini, F8) để huỷ. Đặt node này ở CUỐI kịch bản.'
          : op === 'lock' ? 'Khoá màn hình (macOS bấm Ctrl+Cmd+Q – cần quyền Accessibility). Không cần xác nhận.' : 'Tắt màn hình, máy vẫn chạy. Không cần xác nhận.'}
        {' '}Muốn máy không ngủ trong lúc chạy → Cài đặt chạy → Giữ máy thức.
      </span>
    </>
  );
}
