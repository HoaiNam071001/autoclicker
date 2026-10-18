// Mô tả loại node ở đầu inspector (bước 35 N0) – chữ từ flow/nodeInfo; ẩn / hiện bằng nút ⓘ ở đầu inspector, nhớ cho mọi node.
import { useEffect, useState } from 'react';
import type { NodeInfo } from '../../flow/nodeInfo';

const KEY = 'node-info-hidden';

const loadHidden = () => {
  try {
    return localStorage.getItem(KEY) === '1';
  } catch {
    return false;
  }
};

// [đang hiện, bật / tắt]
export function useNodeInfoShown(): [boolean, () => void] {
  const [hidden, setHidden] = useState(loadHidden);
  useEffect(() => {
    try {
      localStorage.setItem(KEY, hidden ? '1' : '0');
    } catch { /* không lưu được thì thôi – chỉ là tiện ích */ }
  }, [hidden]);
  return [!hidden, () => setHidden((h) => !h)];
}

export function NodeInfoCard({ info, onHide }: { info: NodeInfo; onHide: () => void }) {
  return (
    <div className="relative flex flex-col gap-1 rounded-card border border-line bg-soft py-2.5 pr-8 pl-3">
      <span className="text-[13px] leading-snug text-ink">{info.desc}</span>
      <span className="text-xs leading-snug text-muted"><span className="font-semibold">Dùng khi:</span> {info.use}</span>
      <button onClick={onHide} aria-label="Ẩn mô tả" title="Ẩn mô tả – bấm ⓘ ở đầu inspector để hiện lại"
        className="absolute top-1.5 right-1.5 grid h-6 w-6 place-items-center rounded-control text-muted hover:bg-hover hover:text-ink">
        <svg viewBox="0 0 24 24" className="h-3.5 w-3.5" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" aria-hidden><path d="M6 6l12 12M18 6L6 18" /></svg>
      </button>
    </div>
  );
}
