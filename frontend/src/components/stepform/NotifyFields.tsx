// Node Thông báo hệ thống (bước 43 S9): tiêu đề + nội dung.
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { Field, hint, useNodeVars } from './shared';

export default function NotifyFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const varProps = { inputs: script.inputs ?? [], nodeVars, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  return (
    <>
      <Field label="Tiêu đề" wide group help="Trống = Ghosthand">
        <VarTextInput {...varProps} rows={1} compact value={s.notifyTitle ?? ''} onChange={(v) => updateStep(s.id, { notifyTitle: v || undefined })} placeholder="Ghosthand" />
      </Field>
      <Field label="Nội dung" wide group>
        <VarTextInput {...varProps} rows={3} value={s.text ?? ''} onChange={(v) => updateStep(s.id, { text: v })} placeholder="Đã xử lý xong {{n1.output.count}} file" />
      </Field>
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Hiện một lần rồi đi tiếp, không chờ bấm. macOS: cấp quyền Thông báo cho Ghosthand (hoặc terminal chạy wails dev) và tắt Không làm phiền; Windows: toast hiện dưới tên PowerShell.
      </span>
    </>
  );
}
