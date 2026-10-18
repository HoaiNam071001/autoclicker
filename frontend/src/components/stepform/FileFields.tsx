// Node Tệp / thư mục (bước 35 N4): ghi chữ, thêm dòng csv / xlsx, tạo thư mục, xoá vào / dọn Thùng rác. Ô chữ nhận {{biến}}.
import { Switch } from '../RunSettings';
import VarTextInput from '../inputs/VarTextInput';
import { updateStep, useScript } from '../../store';
import { addInput } from '../../scriptInputs';
import { FILE_OPS, folderOf, isDeleteOp, type FileOp } from '../../flow/file';
import { Step } from '../../types';
import { Field, hint, selectBox, useNodeVars } from './shared';

export function FileFields({ step: s }: { step: Step }) {
  const script = useScript();
  const nodeVars = useNodeVars(s.id);
  const set = (patch: Partial<Step>) => updateStep(s.id, patch);
  const op = s.fileOp ?? '';
  const varProps = { inputs: script.inputs ?? [], nodeVars, onCreateInput: (key: string) => addInput({ key, label: key, type: 'text' }) };
  const folder = op === 'mkdir' || op === 'trash';
  const xlsx = /\.xlsx$/i.test(s.filePath ?? '');
  const badTable = op === 'row' && !!s.filePath?.trim() && !s.filePath.includes('{{') && !/\.(csv|xlsx)$/i.test(s.filePath.trim());
  const setOp = (fileOp: FileOp) => set({
    fileOp: fileOp || undefined,
    ...(fileOp === 'row' && !s.cells?.length ? { cells: [''] } : {}),
    // sang thao tác cần thư mục: ô đang là file thì giữ thư mục chứa nó
    ...((fileOp === 'mkdir' || fileOp === 'trash') && s.filePath ? { filePath: folderOf(s.filePath) } : {}),
  });
  return (
    <>
      <Field label="Làm gì" wide>
        <select value={op} onChange={(e) => setOp(e.target.value as FileOp)} className={selectBox}>
          {FILE_OPS.map((o) => <option key={o.id} value={o.id}>{o.label}</option>)}
        </select>
      </Field>
      {isDeleteOp(op) && (
        <div className={`col-span-2 rounded-card px-3 py-2.5 text-xs ${op === 'emptyTrash' ? 'bg-danger-soft text-danger-text' : 'bg-soft text-muted'}`}>
          {op === 'emptyTrash'
            ? <><b>Xoá hẳn</b> mọi thứ trong Thùng rác – kể cả file không phải do kịch bản tạo, không khôi phục được.</>
            : <>Chuyển vào Thùng rác – <b className="text-ink">khôi phục được</b>. Chỉ file nằm ngay trong thư mục, không gồm thư mục con.</>}
          {' '}Lần chạy đầu app hỏi xác nhận; đổi thư mục / mẫu tên thì hỏi lại.
        </div>
      )}
      {op !== 'emptyTrash' && (
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">{folder ? 'Thư mục' : 'File'}
          <VarTextInput {...varProps} rows={1} compact value={s.filePath ?? ''} onChange={(v) => set({ filePath: v || undefined })}
            placeholder={folder ? '~/Downloads' : op === 'row' ? '~/Documents/don-hang.xlsx' : '~/Documents/nhat-ky.txt'} />
          <span className={`text-xs ${s.filePath?.trim() && !badTable ? 'text-muted' : 'text-danger-text'}`}>
            {!s.filePath?.trim() ? `Chưa nhập ${folder ? 'thư mục' : 'file'} – chạy sẽ báo lỗi`
              : badTable ? 'Thêm dòng chỉ cho file .csv hoặc .xlsx – đổi đuôi file'
              : '~ = thư mục người dùng · thư mục hệ thống, ~/Library bị chặn'}
          </span>
        </div>
      )}
      {op === '' && <>
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">Nội dung
          <VarTextInput {...varProps} rows={4} value={s.text} onChange={(text) => set({ text })} placeholder={'Chạy lúc {{n1.output.startedAt}}\n'} />
        </div>
        <label className="col-span-2 flex items-center gap-2 text-ink">
          <Switch on={!!s.fileAppend} onChange={(on) => set({ fileAppend: on || undefined })} label="Thêm vào cuối file" />
          Thêm vào cuối file <span className="text-xs text-muted">(tắt = ghi đè)</span>
        </label>
      </>}
      {op === 'row' && <CellsField step={s} varProps={varProps} xlsx={xlsx} />}
      {op === 'trash' && (
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">Mẫu tên file
          <VarTextInput {...varProps} rows={1} compact value={s.filePattern ?? ''} onChange={(v) => set({ filePattern: v || undefined })} placeholder="*.png" />
          <span className={`text-xs ${s.filePattern?.trim() ? 'text-muted' : 'text-danger-text'}`}>
            {s.filePattern?.trim() ? <><span className="font-mono">*</span> = mọi ký tự, <span className="font-mono">?</span> = một ký tự · vd <span className="font-mono">bao-cao-*.csv</span> · file ẩn bỏ qua</> : 'Chưa nhập mẫu tên – chạy sẽ báo lỗi'}
          </span>
        </div>
      )}
    </>
  );
}

function CellsField({ step: s, varProps, xlsx }: { step: Step; varProps: Omit<Parameters<typeof VarTextInput>[0], 'value' | 'onChange'>; xlsx: boolean }) {
  const cells = s.cells ?? [];
  const set = (next: string[]) => updateStep(s.id, { cells: next });
  return (
    <>
      <div className="col-span-2 flex flex-col gap-1.5 text-muted">
        <span className="flex items-center justify-between">Các ô của dòng mới<span className="text-xs">{cells.length} ô</span></span>
        {cells.map((c, i) => (
          <div key={i} className="flex items-start gap-1.5">
            <span className="grid h-8 w-6 flex-none place-items-center font-mono text-xs">{String.fromCharCode(65 + (i % 26))}</span>
            <div className="min-w-0 flex-1">
              <VarTextInput {...varProps} rows={1} compact value={c} onChange={(v) => set(cells.map((x, j) => (j === i ? v : x)))} placeholder={i === 0 ? '{{ten}}' : ''} />
            </div>
            <button onClick={() => set(cells.filter((_, j) => j !== i))} aria-label={`Xoá ô ${i + 1}`} disabled={cells.length === 1}
              className="grid h-8 w-7 flex-none place-items-center rounded-control text-muted hover:bg-hover hover:text-ink disabled:opacity-30">
              <svg viewBox="0 0 24 24" className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" aria-hidden><path d="M6 6l12 12M18 6L6 18" /></svg>
            </button>
          </div>
        ))}
        <button onClick={() => set([...cells, ''])} className="h-8 w-fit rounded-control border border-dashed border-line-strong px-3 text-[13px] text-ink hover:bg-hover">+ Thêm ô</button>
        <span className={hint}>{xlsx ? 'Ô trông như số ghi thành số (Excel cộng được). Đang mở file trong Excel thì không ghi được.' : 'CSV: ô có dấu phẩy / xuống dòng tự đặt trong nháy; file mới có BOM để Excel hiện đúng tiếng Việt.'}</span>
      </div>
      {xlsx && (
        <div className="col-span-2 flex flex-col gap-1.5 text-muted">Sheet
          <VarTextInput {...varProps} rows={1} compact value={s.sheet ?? ''} onChange={(v) => updateStep(s.id, { sheet: v || undefined })} placeholder="bỏ trống = sheet đầu tiên" />
        </div>
      )}
    </>
  );
}
