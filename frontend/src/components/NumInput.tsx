import { useEffect, useState } from 'react';

// Ô nhập design v2: cao 32, nền field, viền line, focus viền gradient (ring-focus).
export const fieldBox = 'flex h-8 items-center rounded-control border border-line bg-field px-2.5 focus-within:ring-focus';

// Ô số nguyên ≥ min (signed: cho phép số âm, vd độ lệch ΔX). Giữ bản nháp khi đang gõ (cho phép xoá trắng), chỉ ghi vào store khi hợp lệ;
// rời ô thì hiện lại giá trị thật.
export default function NumInput({ value, onChange, min = 0, unit, prefix, className = '', signed = false }: {
  value: number; onChange: (v: number) => void; min?: number; unit?: string; prefix?: string; className?: string;
  signed?: boolean;
}) {
  const [draft, setDraft] = useState(String(value));
  useEffect(() => setDraft(String(value)), [value]);
  return (
    <span className={`${fieldBox} gap-1.5 text-muted ${className}`}>
      {prefix && <span>{prefix}</span>}
      <input inputMode="numeric" value={draft} className="w-full min-w-0 flex-1 bg-transparent font-mono text-[13px] text-ink outline-none"
        onChange={(e) => {
          const v = e.target.value;
          setDraft(v);
          if (signed ? /^-?\d+$/.test(v) : /^\d+$/.test(v)) onChange(signed ? Number(v) : Math.max(min, Number(v)));
        }}
        onBlur={() => setDraft(String(value))} />
      {unit && <span>{unit}</span>}
    </span>
  );
}
