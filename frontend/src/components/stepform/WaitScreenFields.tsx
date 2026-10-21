// Node Chờ màn hình (bước 42 G1): cách chờ (đứng yên / thay đổi), vùng (như Chụp màn hình), ngưỡng, thời gian.
import Segmented from '../Segmented';
import NumInput from '../NumInput';
import { updateStep } from '../../store';
import { WAIT_MODES, type WaitMode } from '../../flow/waitScreen';
import { Step } from '../../types';
import { Field, fmtSec, NumField } from './shared';
import { ShotAreaField } from './ShotFields';

export function WaitScreenFields({ step: s }: { step: Step }) {
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const change = s.waitMode === 'change';
  return (
    <>
      <Field label="Chờ tới khi" wide group help={change
        ? 'Màn hình khác lúc bắt đầu node – đặt ngay sau node click để chờ trang phản hồi (mở hộp thoại, chuyển trang).'
        : 'Hai lần chụp liền nhau gần như y hệt trong suốt thời gian đứng yên – trang tải xong, spinner tắt, hiệu ứng dừng.'}>
        <Segmented value={s.waitMode ?? ''} options={WAIT_MODES} onChange={(v: WaitMode) => set({ waitMode: v || undefined })} full label="Chờ tới khi" />
      </Field>
      <ShotAreaField step={s} label="Vùng theo dõi" />
      {!change && (
        <Field label="Đứng yên" help={`Không đổi liên tục ${fmtSec(s.stableMs)} mới đi tiếp`}>
          <NumInput value={s.stableMs ?? 0} unit="ms" onChange={(stableMs) => set({ stableMs })} />
        </Field>
      )}
      <Field label="Ngưỡng đổi" help={`% điểm ảnh đổi rõ giữa hai lần so – ${change ? 'vượt ngưỡng là đã đổi' : 'dưới ngưỡng là đứng yên'}. Con trỏ chữ nhấp nháy chỉ ~0,01 %; spinner nhỏ trong vùng lớn có thể < 1 % – thu hẹp vùng hoặc đặt 0.`}>
        <NumInput value={s.diffPct ?? 0} unit="%" onChange={(diffPct) => set({ diffPct: Math.min(diffPct, 100) })} />
      </Field>
      <NumField label="Chờ tối đa" info={`Chưa ${change ? 'đổi' : 'đứng yên'} sau ${fmtSec(s.timeoutMs)} → cổng Thất bại`} step={s} field="timeoutMs" unit="ms" min={1} />
      <Field label="Kiểm tra mỗi">
        <NumInput value={s.intervalMs ?? 200} min={20} unit="ms" onChange={(intervalMs) => set({ intervalMs })} />
      </Field>
    </>
  );
}
