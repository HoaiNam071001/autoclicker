// Inspector node Lặp qua danh sách (bước 42 G2): nguồn chữ / biến hoặc file CSV / xlsx, nút Đọc thử (app thật) hiện
// số phần tử + cột → khoá biến. Chỉ hiển thị – nối store ở stepform/ControlSection.
import { useEffect, useState } from 'react';
import Segmented from '../../Segmented';
import VarTextInput from '../../inputs/VarTextInput';
import { errText, Field, hint } from '../../stepform/shared';
import type { CondCtx, ListPreview } from './CondEditor';

type Peek = { list?: ListPreview; error?: string; busy?: boolean };

export default function ForEachFields({ ctx: { step: s, set, vars, tools } }: { ctx: CondCtx }) {
  const file = s.listFrom === 'file';
  const path = s.filePath?.trim() ?? '';
  const badExt = file && !!path && !path.includes('{{') && !/\.(csv|xlsx)$/i.test(path);
  const [peek, setPeek] = useState<Peek>({});
  useEffect(() => setPeek({}), [s.id, s.listFrom, s.filePath, s.sheet, s.value]);
  const read = tools?.previewList && (() => {
    setPeek({ busy: true });
    tools.previewList!().then((list) => setPeek({ list }), (e) => setPeek({ error: String(e) }));
  });
  return (
    <div className="flex flex-col gap-3">
      <Field label="Danh sách từ" group>
        <Segmented value={file ? 'file' : 'text'} onChange={(v) => set({ listFrom: v === 'file' ? 'file' : undefined })} full label="Danh sách từ"
          options={[{ id: 'text', label: 'Chữ / biến' }, { id: 'file', label: 'File CSV / Excel' }]} />
      </Field>
      {file ? <>
        <Field label="File" group help="~ = thư mục người dùng · dòng đầu là tên cột">
          <VarTextInput {...vars} rows={1} compact value={s.filePath ?? ''} onChange={(v) => set({ filePath: v || undefined })} placeholder="~/Documents/khach-hang.csv" />
          {(!path || badExt) && <span className={errText}>{!path ? 'Chưa nhập file – chạy sẽ báo lỗi' : 'Chỉ đọc được file .csv hoặc .xlsx'}</span>}
        </Field>
        {/\.xlsx$/i.test(path) && (
          <Field label="Sheet" group>
            <VarTextInput {...vars} rows={1} compact value={s.sheet ?? ''} onChange={(v) => set({ sheet: v || undefined })} placeholder="trống = sheet đầu" />
          </Field>
        )}
      </> : (
        <Field label="Danh sách" group help="Mỗi dòng một phần tử, hoặc mảng JSON / biến trỏ output mảng (vd json.items của node HTTP)">
          <VarTextInput {...vars} rows={4} value={s.value ?? ''} onChange={(value) => set({ value })} placeholder={'Táo\nCam\nXoài'} />
          {!s.value?.trim() && <span className={errText}>Danh sách trống – chạy sẽ báo lỗi</span>}
        </Field>
      )}
      {read && (
        <div className="flex flex-col gap-1.5">
          <button onClick={read} disabled={peek.busy} className="h-8 w-fit rounded-control border border-line bg-surface px-3 text-[13px] text-ink enabled:hover:bg-hover disabled:opacity-50">
            {peek.busy ? 'Đang đọc…' : 'Đọc thử'}
          </button>
          {peek.error && <span className={errText}>{peek.error}</span>}
          {peek.list && <PeekResult list={peek.list} ref_={s.ref} />}
        </div>
      )}
      <span className={hint}>
        Mỗi lần đi vào lấy một phần tử ra cổng <b className="text-ink">Mỗi phần tử</b> – nối cạnh từ cuối đoạn cần lặp quay lại node này; hết thì ra <b className="text-ink">Hết danh sách</b>.
        Node sau dùng <span className="var-token">{`${s.ref}.output.item`}</span>, <span className="var-token">{`${s.ref}.output.index`}</span> (từ 1).
      </span>
    </div>
  );
}

function PeekResult({ list, ref_ }: { list: ListPreview; ref_: string }) {
  const first = typeof list.first === 'string' ? list.first : list.first === undefined ? '' : JSON.stringify(list.first);
  return (
    <div className="flex flex-col gap-1 rounded-card bg-soft px-3 py-2 text-xs text-muted">
      <span><b className="text-ink">{list.count}</b> phần tử{!list.columns?.length && first ? <> · đầu tiên: <span className="font-mono text-ink">{first.length > 60 ? `${first.slice(0, 59)}…` : first}</span></> : null}</span>
      {list.columns?.map((c) => (
        <span key={c.key} className="flex items-center gap-1.5">
          <span className="min-w-0 truncate text-ink">{c.name || '(không tên)'}</span>→<span className="var-token">{`${ref_}.output.item.${c.key}`}</span>
        </span>
      ))}
    </div>
  );
}
