// Node Chờ đến giờ (bước 35 N6): giờ hẹn (nhận {{biến}}) + xem trước còn bao lâu, giờ đã qua thì làm gì, canh giờ chuẩn.
import { useEffect, useState } from 'react';
import Segmented from '../Segmented';
import { Switch } from '../RunSettings';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { PASSED_MODES, untilText, type IfPassed } from '../../flow/waitUntil';
import { Step } from '../../types';
import { Field, hint, useNodeVars } from './shared';

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
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">Giờ hẹn
        <VarTextInput inputs={script.inputs ?? []} nodeVars={nodeVars} rows={1} compact value={s.at ?? ''} onChange={(at) => set({ at })}
          onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} placeholder="20:00:00.000" />
        <span className={`text-xs ${preview?.bad ? 'text-danger-text' : 'text-muted'}`}>
          {preview ? preview.text : 'Giờ lấy từ biến – biết khi chạy'}
          <span className="text-muted"> · HH:mm[:ss[.mmm]], thêm ngày phía trước nếu cần</span>
        </span>
      </div>
      <Field label="Nếu giờ đã qua" wide group>
        <Segmented value={s.ifPassed ?? ''} options={PASSED_MODES} onChange={(v: IfPassed) => set({ ifPassed: v || undefined })} full label="Nếu giờ đã qua" />
        <span className={hint}>{s.ifPassed === 'tomorrow' ? 'Hợp với vòng lặp chạy mỗi ngày – giờ có ghi ngày thì đi tiếp ngay.' : s.ifPassed === 'fail' ? 'Đi cổng Thất bại (timeout).' : 'Tới node khi đã quá giờ thì chạy tiếp luôn.'}</span>
      </Field>
      <div className="col-span-2 flex flex-col gap-1">
        <label className="flex items-center gap-2 text-ink">
          <Switch on={!!s.useNtp} onChange={(on) => set({ useNtp: on || undefined })} label="Canh giờ chuẩn internet" />
          Canh giờ chuẩn internet (NTP)
        </label>
        <span className={`pl-11 ${hint}`}>Đồng hồ máy lệch vẫn đúng giờ – nên bật khi săn sale. Không lấy được giờ chuẩn thì dùng giờ máy.</span>
      </div>
    </>
  );
}
