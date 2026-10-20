// Nút xoá ở đầu hộp Lịch sử chạy (bước 41 L3): xoá lần đang chọn ngay; "Xoá hết" hỏi lại ngay trong thanh.
// Chỉ hiển thị – RunsDialog gọi runStore.deleteRuns.
import { useState } from 'react';
import { Trash } from '../icons';

export default function RunsDelete({ count, canDelete, onDelete, onClear }: {
  count: number; canDelete: boolean; onDelete: () => void; onClear: () => void;
}) {
  const [asking, setAsking] = useState(false);
  const btn = 'h-8 rounded-control border border-line px-2.5 text-[13px] enabled:hover:bg-hover disabled:opacity-50';
  if (asking) {
    return (
      <span className="flex items-center gap-1.5 rounded-control bg-soft px-2 py-1 text-[13px]">
        Xoá cả {count} lần chạy đã lưu (kèm ảnh)?
        <button onClick={() => { setAsking(false); onClear(); }} className="btn-danger-soft h-7 rounded-control px-2.5 font-semibold">Xoá hết</button>
        <button onClick={() => setAsking(false)} className="h-7 rounded-control border border-line px-2.5 hover:bg-hover">Huỷ</button>
      </span>
    );
  }
  return (
    <>
      <button disabled={!canDelete} onClick={onDelete} title="Xoá lần chạy đang chọn (nhật ký + ảnh lúc lỗi)" className={`${btn} flex items-center gap-1.5`}>
        <Trash className="h-3.5 w-3.5" />Xoá lần này
      </button>
      <button disabled={!count} onClick={() => setAsking(true)} className={btn}>Xoá hết…</button>
    </>
  );
}
