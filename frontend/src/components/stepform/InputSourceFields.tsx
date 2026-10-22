// Node Bộ gõ (bước 43 S6): đọc hoặc đổi bộ gõ của hệ điều hành.
import Segmented from '../Segmented';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { INPUT_OPS, type InputOp } from '../../flow/inputSource';
import { Step } from '../../types';
import { Field, hint, useNodeVars } from './shared';

export default function InputSourceFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const varProps = { inputs: script.inputs ?? [], nodeVars, rows: 1, compact: true, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  return (
    <>
      <Field label="Làm gì" wide group>
        <Segmented value={s.inputOp ?? ''} full label="Thao tác bộ gõ" onChange={(v) => set({ inputOp: (v || undefined) as InputOp | undefined })} options={INPUT_OPS} />
      </Field>
      {s.inputOp === 'set' && (
        <Field label="Bộ gõ cần đổi sang" wide group help="Tên hoặc mã bộ gõ, chứa là khớp, không phân biệt hoa thường. macOS: ABC, U.S., Telex · Windows: en-US, vi-VN. Không có → lỗi kèm danh sách bộ gõ máy có.">
          <VarTextInput {...varProps} value={s.inputLayout ?? ''} onChange={(v) => set({ inputLayout: v || undefined })} placeholder="ABC" />
        </Field>
      )}
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Đặt trước node Gõ chữ để dấu không bị méo. Node sau dùng {`{{${s.ref}.output.name}}`}, {`{{${s.ref}.output.vietnamese}}`}, {`{{${s.ref}.output.imeApp}}`} (Unikey / EVKey đang chạy – chỉ báo, chưa bật / tắt được)…
      </span>
    </>
  );
}
