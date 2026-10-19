// Dòng đầu canvas flow (design v3): "Flow · N node" bên trái, nút (vùng, + Thêm node…) bên phải – nằm đè trên canvas,
// nền trong suốt; chỉ nút nhận chuột để kéo / chọn trên canvas không bị chặn. Dùng chung app (FlowEditor) và #demo-flow.
import { ReactNode } from 'react';

export default function FlowHeader({ count, children }: { count: number; children: ReactNode }) {
  return (
    <div className="pointer-events-none absolute inset-x-4 top-3 z-10 flex flex-wrap items-center gap-2.5">
      <h1 className="font-semibold">Flow</h1>
      <span className="text-[12.5px] whitespace-nowrap text-muted">{count} node</span>
      <span className="pointer-events-auto ml-auto flex items-center gap-1.5">{children}</span>
    </div>
  );
}

// Nút trên dòng đầu canvas
export const headerBtn = 'flex h-8 items-center gap-1.5 whitespace-nowrap rounded-control border border-line bg-surface px-2.5 text-[12.5px] hover:bg-hover';
