// Node Chờ đến giờ (bước 35 N6): giờ hẹn (nhận {{biến}}) + xem trước còn bao lâu, giờ đã qua thì làm gì, canh giờ chuẩn.
import { useEffect, useState } from 'react';
import Segmented from '../Segmented';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { PASSED_MODES, untilText, type IfPassed } from '../../flow/waitUntil';
import { Step } from '../../types';
import { Field, SwitchField, useNodeVars } from './shared';

export function WaitUntilFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const [now, setNow] = useState(() => new Date());
  useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 500);
    return () => clearInterval(t);
  }, []);
  const preview = untilText(s, now);
  return (
    <>
      <Field label="Giờ hẹn" wide group help="HH:mm[:ss[.mmm]], thêm ngày phía trước nếu cần">
        <VarTextInput inputs={script.inputs ?? []} nodeVars={nodeVars} rows={1} compact value={s.at ?? ''} onChange={(at) => set({ at })}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} placeholder="20:00:00.000" />
        <span className={`text-xs ${preview?.bad ? 'text-danger-text' : 'text-muted'}`}>
          {preview ? preview.text : 'Giờ lấy từ biến – biết khi chạy'}
        </span>
      </Field>
      <Field label="Nếu giờ đã qua" wide group help={s.ifPassed === 'tomorrow' ? 'Hợp với vòng lặp chạy mỗi ngày – giờ có ghi ngày thì đi tiếp ngay.' : s.ifPassed === 'fail' ? 'Đi cổng Thất bại (timeout).' : 'Tới node khi đã quá giờ thì chạy tiếp luôn.'}>
        <Segmented value={s.ifPassed ?? ''} options={PASSED_MODES} onChange={(v: IfPassed) => set({ ifPassed: v || undefined })} full label="Nếu giờ đã qua" />
      </Field>
      <SwitchField label="Giờ chuẩn" help="Canh theo giờ chuẩn internet: đồng hồ máy lệch vẫn đúng giờ – nên bật khi săn sale. Không lấy được giờ chuẩn thì dùng giờ máy."
        on={!!s.useNtp} onChange={(on) => set({ useNtp: on || undefined })} />
    </>
  );
}
