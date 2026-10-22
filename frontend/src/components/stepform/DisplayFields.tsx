// Node Màn hình / âm thanh (bước 43 S7): chế độ sáng / tối, độ sáng, loa / micro mặc định, phát âm báo.
import NumInput from '../NumInput';
import Segmented from '../Segmented';
import Select from '../Select';
import { Switch } from '../Switch';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { DISPLAY_OPS, SYSTEM_SOUNDS, THEMES, type DisplayOp } from '../../flow/display';
import { Step } from '../../types';
import { Field, hint, selectBox, useNodeVars } from './shared';

const OP_HINT: Record<DisplayOp, string> = {
  '': 'Đọc một lần, không đổi gì – dùng theme / brightness / outputName với node Nếu.',
  theme: 'Đổi chế độ sáng / tối của hệ điều hành (macOS lần đầu hỏi quyền điều khiển System Events).',
  brightness: 'Chỉ màn hình laptop; màn ngoài / máy bàn → Thất bại.',
  audioOut: 'Tên hoặc mã thiết bị, chứa là khớp. Không có → lỗi kèm danh sách thiết bị máy có.',
  audioIn: 'Tên hoặc mã thiết bị, chứa là khớp. Không có → lỗi kèm danh sách thiết bị máy có.',
  sound: 'macOS: file aiff / wav / mp3 hoặc tên âm (Glass, Ping…) · Windows: file .wav hoặc Asterisk, Beep, Exclamation… Âm lượng → node Hệ thống.',
};

export default function DisplayFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const op = s.displayOp ?? '';
  const value = s.displayValue ?? '';
  const varProps = { inputs: script.inputs ?? [], nodeVars, rows: 1, compact: true, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  const text = (placeholder: string) => <VarTextInput {...varProps} value={value} onChange={(v) => set({ displayValue: v || undefined })} placeholder={placeholder} />;
  return (
    <>
      <Field label="Làm gì" wide group>
        <Select value={op} onChange={(v) => set({ displayOp: (v || undefined) as DisplayOp | undefined, displayValue: undefined })} className={selectBox} aria-label="Thao tác màn hình / âm thanh"
          options={DISPLAY_OPS.map((o) => ({ value: o.id, label: o.label }))} />
      </Field>
      {op === 'theme' && (
        <Field label="Chế độ" wide group>
          <Segmented value={value} full label="Chế độ sáng / tối" onChange={(v) => set({ displayValue: v })} options={THEMES} />
        </Field>
      )}
      {op === 'brightness' && (
        <Field label="Độ sáng" help="0–100 %">
          <NumInput value={Number(value) || 0} unit="%" min={0} onChange={(n) => set({ displayValue: String(Math.min(100, n)) })} />
        </Field>
      )}
      {(op === 'audioOut' || op === 'audioIn') && (
        <Field label={op === 'audioOut' ? 'Loa' : 'Micro'} wide group>{text(op === 'audioOut' ? 'AirPods' : 'MacBook Pro Microphone')}</Field>
      )}
      {op === 'sound' && (
        <>
          <Field label="File hoặc tên âm" wide group help={`Âm hệ thống macOS: ${SYSTEM_SOUNDS.join(', ')}`}>{text('Glass')}</Field>
          <Field label="Chờ phát xong" help="Bật: node chỉ đi tiếp khi âm kết thúc">
            <Switch on={!!s.soundWait} onChange={(soundWait) => set({ soundWait: soundWait || undefined })} label="Chờ phát xong" />
          </Field>
        </>
      )}
      <span className={`col-span-2 -mt-1 ${hint}`}>{OP_HINT[op]}</span>
    </>
  );
}
