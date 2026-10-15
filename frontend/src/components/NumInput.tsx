import { useEffect, useState } from 'react';

const box = 'flex h-10 items-center rounded-lg border border-line px-2.5 font-normal focus-within:border-accent';

// Ô số nguyên ≥ min. Giữ bản nháp khi đang gõ (cho phép xoá trắng), chỉ ghi vào store khi hợp lệ;
// rời ô thì hiện lại giá trị thật.
export default function NumInput({ value, onChange, min = 0, unit, prefix, className = '' }: {
  value: number; onChange: (v: number) => void; min?: number; unit?: string; prefix?: string; className?: string;
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
          if (/^\d+$/.test(v)) onChange(Math.max(min, Number(v)));
        }}
        onBlur={() => setDraft(String(value))} />
      {unit && <span className="text-muted">{unit}</span>}
    </span>
  );
}
