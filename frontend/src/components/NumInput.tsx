import { useEffect, useState } from 'react';

const box = 'flex h-10 items-center rounded-lg border border-line px-2.5 font-normal focus-within:border-accent';

// Ô số nguyên ≥ min (signed: cho phép số âm, vd độ lệch ΔX). Giữ bản nháp khi đang gõ (cho phép xoá trắng), chỉ ghi vào store khi hợp lệ;
// rời ô thì hiện lại giá trị thật.
export default function NumInput({ value, onChange, min = 0, unit, prefix, className = '', signed = false }: {
  value: number; onChange: (v: number) => void; min?: number; unit?: string; prefix?: string; className?: string;
  signed?: boolean;
}) {
  const [draft, setDraft] = useState(String(value));
  useEffect(() => setDraft(String(value)), [value]);
  return (
    <span className={`${box} gap-1.5 ${className}`}>
      {prefix && <span className="text-muted">{prefix}</span>}
      <input inputMode="numeric" value={draft} className="w-full min-w-0 flex-1 font-mono outline-none"
        onChange={(e) => {
          const v = e.target.value;
          setDraft(v);
          if (signed ? /^-?\d+$/.test(v) : /^\d+$/.test(v)) onChange(signed ? Number(v) : Math.max(min, Number(v)));
        }}
        onBlur={() => setDraft(String(value))} />
      {unit && <span className="text-muted">{unit}</span>}
    </span>
  );
}
