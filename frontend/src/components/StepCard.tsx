// Phần hiển thị thẻ bước + khung nhóm (design v2) – dùng chung cho danh sách thật (StepList) và màn #demo.
import { ReactNode } from 'react';
import { StepIcon, StepIconName } from './icons';

export const cardClass = (selected: boolean) =>
  `group relative flex cursor-default items-center gap-3 rounded-card px-3 py-[9px] ${selected ? 'ring-selected' : 'border border-line bg-surface'}`;

export const tileClass = (selected: boolean) =>
  `grid h-8 w-8 flex-none place-items-center rounded-lg ${selected ? 'bg-accent-soft text-accent-text' : 'bg-soft text-muted'}`;

// Phần giữa + phải của thẻ: "01  Tên  [chip]" / chi tiết mono, chip phụ (viền), chip thời gian.
export function CardBody({ n, title, detail, delay, enabled = true, tags, chips, breakpoint, media }: {
  n: string; title: string; detail: string; delay: string; enabled?: boolean;
  tags?: ReactNode; // cạnh tên: "Tắt", "Sắp có"
  chips?: string[]; // viền line: "1 input", "Hết giờ: thử lại 2 lần"
  breakpoint?: boolean; // chấm đỏ – điểm dừng khi debug (bước 25c)
  media?: ReactNode; // hình nhỏ bên phải (ảnh mẫu của bước chờ ảnh)
}) {
  return (
    <>
      <span className={`min-w-0 flex-1 ${enabled ? '' : 'opacity-50'}`}>
        <span className="flex items-center gap-2">
          {breakpoint && <span className="h-2 w-2 flex-none rounded-full bg-danger" title="Điểm dừng (debug)" />}
          <span className="font-mono text-xs text-muted">{n}</span>
          <span className="truncate font-semibold">{title}</span>
          {tags}
        </span>
        <span className="block truncate font-mono text-xs text-muted">{detail}</span>
      </span>
      {media && <span className="flex-none group-hover:invisible">{media}</span>}
      {chips?.map((c) => (
        <span key={c} className="rounded-[5px] border border-line px-2 py-0.5 text-xs whitespace-nowrap text-muted group-hover:invisible">{c}</span>
      ))}
      {/* invisible (không phải hidden) để giữ nguyên chiều rộng – rê chuột không làm thẻ nhảy */}
      {delay && <span className="rounded-[5px] bg-soft px-2 py-[3px] font-mono text-xs whitespace-nowrap text-muted group-hover:invisible">{delay}</span>}
    </>
  );
}

// Thẻ tĩnh (không kéo được) – màn #demo
export function StaticCard({ icon, selected, onSelect, ...body }: Parameters<typeof CardBody>[0] & {
  icon: StepIconName; selected: boolean; onSelect: () => void;
}) {
  return (
    <li onClick={onSelect} className={cardClass(selected)}>
      <span className={`${tileClass(selected)} ${body.enabled === false ? 'opacity-50' : ''}`}><StepIcon type={icon} /></span>
      <CardBody {...body} />
    </li>
  );
}

// Khung nhóm bước (bước 20): nét đứt, tên + chip "lặp ×N", các thẻ con bên trong.
export function GroupFrame({ name, repeat, selected, onSelect, children }: {
  name: string; repeat: number; selected?: boolean; onSelect?: () => void; children: ReactNode;
}) {
  return (
    <li className={`flex flex-col gap-1.5 rounded-group border border-dashed bg-group p-2 ${selected ? 'border-accent' : 'border-line-strong'}`}>
      <div onClick={onSelect} className="flex cursor-default items-center gap-2.5 px-1 py-0.5">
        <StepIcon type="group" className="h-[15px] w-[15px] text-muted" />
        <span className="font-semibold">Nhóm: {name}</span>
        <span className="ml-auto rounded-[5px] bg-soft px-2 py-[3px] font-mono text-xs text-muted">lặp ×{repeat}</span>
      </div>
      <ol className="flex flex-col gap-1.5">{children}</ol>
    </li>
  );
}

// Nhãn "Sắp có" – tính năng đã có giao diện, logic làm ở bước sau (27b).
export const Soon = ({ step }: { step?: string }) => (
  <span title={step ? `Sắp có (bước ${step})` : 'Sắp có'}
    className="rounded-full bg-soft px-1.5 py-px text-[10px] font-semibold whitespace-nowrap text-muted">Sắp có</span>
);
