// Ô Tab của node web + Phím trên trang web: chọn tab Chrome của app theo URL / tiêu đề / id tab; nhận {{biến}}.
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { Step } from '../../types';
import { Field, useNodeVars } from './shared';

export function WebTabField({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  return (
    <Field label="Tab" wide group help="Tuỳ chọn: URL hoặc tiêu đề tab chứa chữ này (hoặc id tab AI đọc được – đổi khi Chrome mở lại). Trang phải mở bằng node Mở với “Chrome của app”.">
      <VarTextInput value={s.tab ?? ''} onChange={(tab) => updateStep(s.id, { tab: tab || undefined })} inputs={script.inputs ?? []} nodeVars={nodeVars}
        onCreateInput={(key) => addInput({ key, label: key, type: 'text' })} rows={1} compact placeholder="trống = tab đang hiện" />
    </Field>
  );
}
