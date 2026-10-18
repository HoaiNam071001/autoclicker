// Node Clipboard (bước 35 N2): đọc chữ / ghi chữ / ghi ảnh PNG từ file – ô chữ nhận {{biến}}.
import Segmented from '../Segmented';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { CLIP_OPS, type ClipOp } from '../../flow/clipboard';
import { Step } from '../../types';
import { Field, hint, useNodeVars } from './shared';

export function ClipFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const op = s.clipOp ?? '';
  const varProps = { inputs: script.inputs ?? [], nodeVars, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  // node Chụp màn hình gần nhất – gợi ý điền sẵn đường dẫn ảnh của nó
  const shot = script.steps.find((x) => x.type === 'screenshot' && x.id !== s.id);
  return (
    <>
      <Field label="Làm gì" wide group>
        <Segmented value={op} options={CLIP_OPS} onChange={(clipOp: ClipOp) => set({ clipOp: clipOp || undefined })} full label="Phép clipboard" />
        <span className={hint}>{op === 'write' ? 'Thay nội dung clipboard bằng chữ dưới – để trống = xoá clipboard'
          : op === 'image' ? 'Chép ảnh PNG vào clipboard để dán (Cmd/Ctrl+V) vào chat, email…'
          : <>Đọc chữ đang có – node sau dùng <span className="font-mono">{`{{${s.ref}.output.text}}`}</span>; clipboard chỉ có ảnh thì ra rỗng</>}</span>
      </Field>
      {op === 'write' && (
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">Nội dung
          <VarTextInput {...varProps} value={s.text} onChange={(text) => set({ text })} placeholder="Mã đơn: {{n2.output.text}}" />
        </div>
      )}
      {op === 'image' && (
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">File ảnh PNG
          <VarTextInput {...varProps} rows={1} compact value={s.filePath ?? ''} onChange={(v) => set({ filePath: v || undefined })}
            placeholder={shot ? `{{${shot.ref}.output.path}}` : '~/Pictures/anh.png'} />
          <span className={`text-xs ${s.filePath?.trim() ? 'text-muted' : 'text-danger-text'}`}>
            {s.filePath?.trim() ? '~ / đường dẫn tương đối tính từ thư mục người dùng' : 'Chưa nhập file – chạy sẽ báo lỗi'}
            {shot && !s.filePath && <> · <button className="text-accent-text underline" onClick={() => set({ filePath: `{{${shot.ref}.output.path}}` })}>
              dùng ảnh của {shot.ref}</button></>}
          </span>
        </div>
      )}
    </>
  );
}
