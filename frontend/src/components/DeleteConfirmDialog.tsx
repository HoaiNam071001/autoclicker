// Xác nhận thao tác xoá trước lần chạy đầu (bước 35 N4) – liệt kê từng node xoá gì. Enter = Chạy, Esc = Huỷ.
import { useEffect } from 'react';
import { answerDelete, useDeleteAsk } from '../deleteConfirm';
import { deleteOpText } from '../flow/file';
import { Trash } from './icons';

export default function DeleteConfirmDialog() {
  const ask = useDeleteAsk();
  useEffect(() => {
    if (!ask) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') answerDelete(false);
      else if (e.key === 'Enter') answerDelete(true);
      else return;
      e.preventDefault();
      e.stopPropagation();
    };
    window.addEventListener('keydown', onKey, true);
    return () => window.removeEventListener('keydown', onKey, true);
  }, [ask]);
  if (!ask) return null;
  const permanent = ask.ops.some((o) => o.op === 'emptyTrash');
  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center bg-black/25 px-4 pt-[14vh]"
      onMouseDown={(e) => e.target === e.currentTarget && answerDelete(false)}>
      <div role="alertdialog" aria-labelledby="del-title" aria-describedby="del-desc"
        className="flex w-[460px] max-w-full flex-col overflow-hidden rounded-panel border border-line bg-surface text-ink shadow-[0_16px_48px_rgb(0_0_0/.22)]">
        <div className="flex gap-3 p-[18px]">
          <span className="grid h-9 w-9 flex-none place-items-center rounded-[9px] bg-danger-soft text-danger-text"><Trash className="h-[18px] w-[18px]" /></span>
          <div className="flex min-w-0 flex-col gap-2">
            <div id="del-title" className="text-base font-semibold">Kịch bản này sẽ xoá file</div>
            <ul id="del-desc" className="flex flex-col gap-1.5 text-[13px]">
              {ask.ops.map((o, i) => (
                <li key={i} className="flex gap-2"><span className="flex-none rounded bg-soft px-1.5 font-mono text-xs leading-5 text-muted">{o.ref}</span>
                  <span className={o.op === 'emptyTrash' ? 'text-danger-text' : ''}>{deleteOpText(o)}</span></li>
              ))}
            </ul>
            <span className="text-xs text-muted">Đồng ý một lần – lần sau chạy không hỏi lại, trừ khi đổi thư mục / mẫu tên.{permanent ? ' Dọn Thùng rác xoá cả file không phải do kịch bản tạo.' : ''}</span>
          </div>
        </div>
        <div className="flex items-center gap-2 border-t border-line bg-chrome px-[18px] py-3">
          <button onClick={() => answerDelete(false)} className="ml-auto h-9 rounded-control border border-line bg-surface px-3.5 hover:bg-hover">Huỷ <span className="font-mono text-[11px] text-muted">Esc</span></button>
          <button autoFocus onClick={() => answerDelete(true)} className="btn-danger-soft h-9 rounded-control px-4 font-semibold">Cho phép xoá và chạy</button>
        </div>
      </div>
    </div>
  );
}
