// Input của kịch bản – kiểu dữ liệu cho phần giao diện (bước 27b). Bước 21b chuyển vào types.ts + model.go
// (Script.inputs) và thêm resolveInputs bên Go; tới lúc đó chỉ màn #demo dùng.

export type InputType = 'text' | 'multiline' | 'number' | 'ms' | 'bool' | 'select' | 'point';
export type InputValue = string | number | boolean;

export const INPUT_TYPE_LABEL: Record<InputType, string> = {
  text: 'Chữ',
  multiline: 'Nhiều dòng',
  number: 'Số',
  ms: 'ms',
  bool: 'Bật / tắt',
  select: 'Danh sách chọn',
  point: 'Toạ độ',
};

export type ScriptInput = {
  key: string; // chữ không dấu, số, _ – dùng trong {{key}}
  label: string;
  type: InputType;
  default?: InputValue;
  required?: boolean;
  hint?: string;
  unit?: string; // hiện bên phải ô (s, ms)
  prefix?: string; // hiện mờ bên trái ô (facebook.com/)
  options?: string[]; // type = select
  usedBy: { n: string; title: string }[]; // bước dùng input này (ghim hoặc {{key}})
};

export const isValidKey = (k: string) => /^[A-Za-z_][A-Za-z0-9_]*$/.test(k);

// Tách chuỗi theo {{key}} – \{{ là ký tự {{ thật, không phải biến.
export type Segment = { text: string; varKey?: string };
export function splitVars(text: string): Segment[] {
  const out: Segment[] = [];
  const re = /\\\{\{|\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}/g;
  let last = 0;
  for (let m = re.exec(text); m; m = re.exec(text)) {
    if (m.index > last) out.push({ text: text.slice(last, m.index) });
    out.push(m[1] ? { text: m[0], varKey: m[1] } : { text: m[0] });
    last = m.index + m[0].length;
  }
  if (last < text.length) out.push({ text: text.slice(last) });
  return out;
}

export const varsIn = (text: string) => [...new Set(splitVars(text).flatMap((s) => (s.varKey ? [s.varKey] : [])))];

export const fmtValue = (v: InputValue | undefined, unit?: string) =>
  v === undefined || v === '' ? 'trống' : typeof v === 'boolean' ? (v ? 'bật' : 'tắt') : `${v}${unit ? ` ${unit}` : ''}`;
