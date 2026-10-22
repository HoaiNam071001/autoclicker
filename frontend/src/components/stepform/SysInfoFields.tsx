// Node Thông tin máy (bước 43 S4): chỉ có ô chọn ổ / thư mục đo chỗ trống; mọi thông số khác luôn được đọc.
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { Field, hint, useNodeVars } from './shared';

export default function SysInfoFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const varProps = { inputs: script.inputs ?? [], nodeVars, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  return (
    <>
      <Field label="Ổ / thư mục đo chỗ trống" wide group help="Trống = thư mục người dùng · ví dụ / hoặc D:\">
        <VarTextInput {...varProps} rows={1} compact value={s.filePath ?? ''} onChange={(v) => updateStep(s.id, { filePath: v || undefined })} placeholder="/" />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Node sau dùng {`{{${s.ref}.output.batteryPercent}}`}, {`{{${s.ref}.output.diskFreeGB}}`}, {`{{${s.ref}.output.online}}`}… Rẽ nhánh bằng node Nếu (vd pin &lt; 20).
      </span>
    </>
  );
}
