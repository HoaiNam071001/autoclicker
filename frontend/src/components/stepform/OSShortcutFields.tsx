// Node Shortcuts hệ điều hành (bước 43 S2, macOS): tên Shortcut (gõ hoặc chọn từ app Phím tắt) + input chữ + hết giờ.
import { useEffect, useState } from 'react';
import { ListShortcuts } from '../../../wailsjs/go/main/App';
import Select from '../Select';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { errorText } from '../../engineStore';
import { inWails } from '../../windowMode';
import { Step } from '../../types';
import { Field, hint, NumField, selectBox, useNodeVars } from './shared';

export default function OSShortcutFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const varProps = { inputs: script.inputs ?? [], nodeVars, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  const [names, setNames] = useState<string[]>([]);
  const [error, setError] = useState<string>();
  useEffect(() => {
    if (!inWails) return;
    ListShortcuts().then((l) => setNames(l ?? [])).catch((e) => setError(errorText(e)));
  }, []);
  return (
    <>
      <Field label="Tên Shortcut" wide group help="Đúng như trong ứng dụng Phím tắt · chỉ chạy trên macOS">
        <VarTextInput {...varProps} rows={1} compact value={s.shortcutName ?? ''} onChange={(v) => updateStep(s.id, { shortcutName: v || undefined })} placeholder="Đổi tên ảnh" />
        {names.length > 0 && (
          <Select value="" onChange={(v) => v && updateStep(s.id, { shortcutName: v })} className={selectBox} aria-label="Chọn Shortcut có sẵn"
            options={[{ value: '', label: 'Chọn từ danh sách…' }, ...names.map((n) => ({ value: n, label: n }))]} />
        )}
        {error && <span className="text-xs text-danger-text">{error}</span>}
      </Field>
      <Field label="Input (chữ)" wide group help="Trống = chạy không input">
        <VarTextInput {...varProps} rows={3} value={s.text ?? ''} onChange={(v) => updateStep(s.id, { text: v || undefined })} placeholder="{{n1.output.text}}" />
      </Field>
      <NumField label="Chờ tối đa" info="Quá giờ thì dừng Shortcut và đi cổng Thất bại" step={s} field="timeoutMs" unit="ms" min={1} />
      <span className={`col-span-2 -mt-1 ${hint}`}>
        Node sau dùng {`{{${s.ref}.output.text}}`} (chữ Shortcut trả về). Không có Shortcut tên đó → cổng Thất bại.
      </span>
    </>
  );
}
