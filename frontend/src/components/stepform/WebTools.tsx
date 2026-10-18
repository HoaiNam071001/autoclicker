// Công cụ của node web (bước 25d D3) – chỉ hiển thị: trạng thái Chrome của app, Chọn phần tử trên trang, Thử tìm.
// Nối store ở WebFields (webSteps.ts).
const btn = 'h-8 rounded-control border border-line bg-surface px-3 text-ink enabled:hover:bg-hover disabled:opacity-50';

export interface WebToolsProps {
  status: string;
  open: boolean;
  enabled: boolean; // có Go (không phải trình duyệt / #demo)
  picking: boolean;
  finding: boolean;
  canFind: boolean; // đã nhập cách tìm hợp lệ
  result: { ok: boolean; text: string; elText?: string } | null;
  error: string;
  onOpen: () => void;
  onRefresh: () => void;
  onPick: () => void;
  onCancelPick: () => void;
  onFind: () => void;
  onView: () => void;
}

export function WebTools(p: WebToolsProps) {
  return (
    <div className="col-span-2 flex flex-col gap-2 text-[13px]">
      <div className="flex items-center gap-2 text-xs text-muted">
        <span className={`h-2 w-2 shrink-0 rounded-full ${p.open ? 'bg-ok' : 'bg-line-strong'}`} />
        <span className="min-w-0 flex-1 truncate" title={p.status}>{p.status}</span>
        {p.enabled && !p.open && <button onClick={p.onOpen} className="text-ink underline">Mở</button>}
        {p.enabled && <button onClick={p.onRefresh} className="underline">Làm mới</button>}
      </div>
      <div className="grid grid-cols-2 gap-2">
        {p.picking
          ? <button onClick={p.onCancelPick} className={btn} title="Hoặc bấm Esc trong trang">Huỷ chọn</button>
          : <button onClick={p.onPick} disabled={!p.enabled} className={btn}
            title="Đưa Chrome của app lên trước – bấm vào phần tử cần tìm, app tự điền cách tìm">Chọn phần tử trên trang</button>}
        <button onClick={p.onFind} disabled={!p.enabled || !p.canFind || p.finding || p.picking} className={btn}
          title="Tìm một lần trên tab hiện tại (lần đầu chuột đi vào trang để đo toạ độ)">
          {p.finding ? 'Đang tìm…' : 'Thử tìm'}
        </button>
      </div>
      {p.picking && <span className="text-xs text-accent-text">Bấm vào phần tử trong Chrome của app (Esc để huỷ)…</span>}
      {p.result && (
        <span className={`flex flex-col gap-0.5 text-xs ${p.result.ok ? 'text-ok' : 'text-warn'}`}>
          <span>{p.result.text}{' '}
            {p.result.ok && <button onClick={p.onView} className="text-muted underline">xem trên màn hình</button>}
          </span>
          {p.result.elText && <span className="truncate text-muted">Chữ: {p.result.elText}</span>}
        </span>
      )}
      {p.error && <span className="text-xs text-danger-text">{p.error}</span>}
    </div>
  );
}
