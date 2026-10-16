import { useEffect, useState } from 'react';

// Ô nhập design v2: cao 32, nền field, viền line, focus viền gradient (ring-focus).
export const fieldBox = 'flex h-8 items-center rounded-control border border-line bg-field px-2.5 focus-within:ring-focus';
// Trường đã ghim làm input (design v2: field-pinned)
const pinnedBox = 'flex h-8 items-center rounded-control border border-dashed border-accent bg-field px-2.5';

// Ô số nguyên ≥ min (signed: cho phép số âm, vd độ lệch ΔX). Giữ bản nháp khi đang gõ (cho phép xoá trắng), chỉ ghi vào store khi hợp lệ;
// rời ô thì hiện lại giá trị thật. scale > 1: hiện value / scale, cho gõ số lẻ (vd neo % lưu theo 1/10000, scale 100 → 33.33%).
export default function NumInput({ value, onChange, min = 0, unit, prefix, className = '', signed = false, pinned = false, scale = 1 }: {
  value: number; onChange: (v: number) => void; min?: number; unit?: string; prefix?: string; className?: string;
  signed?: boolean; pinned?: boolean; scale?: number;
}) {
  const fmt = (v: number) => (Object.is(v, -0) ? '-0' : String(scale === 1 ? v : Number((v / scale).toFixed(2))));
  const valid = scale === 1 ? (signed ? /^-?\d+$/ : /^\d+$/) : (signed ? /^-?\d+([.,]\d*)?$/ : /^\d+([.,]\d*)?$/);
  const [draft, setDraft] = useState(fmt(value));
  useEffect(() => setDraft(fmt(value)), [value, scale]); // eslint-disable-line react-hooks/exhaustive-deps
  return (
    <span className={`${pinned ? pinnedBox : fieldBox} gap-1.5 text-muted ${className}`}>
      {prefix && <span>{prefix}</span>}
      <input inputMode="numeric" value={draft} className="w-full min-w-0 flex-1 bg-transparent font-mono text-[13px] text-ink outline-none"
        onChange={(e) => {
          const v = e.target.value;
          setDraft(v);
          if (!valid.test(v)) return;
          const n = Math.round(Number(v.replace(',', '.')) * scale);
          onChange(signed ? n : Math.max(min, n));
        }}
        onBlur={() => setDraft(fmt(value))} />
      {unit && <span>{unit}</span>}
    </span>
  );
}
