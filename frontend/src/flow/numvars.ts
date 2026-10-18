// Trường số lấy từ biến (bước 20 F5c, Step.numVars) – hàm thuần sửa node; khớp script/numvars.go.
import type { NumVarField, Step } from '../types';

// Đặt biểu thức của trường field (giữ vị trí dòng cũ); expr undefined = bỏ, trường dùng lại số cố định.
export function withNumVar(s: Step, field: NumVarField, expr: string | undefined): Step {
  const list = s.numVars ?? [];
  const i = list.findIndex((v) => v.field === field);
  const numVars = expr === undefined ? list.filter((v) => v.field !== field)
    : i < 0 ? [...list, { field, expr }] : list.map((v, j) => (j === i ? { field, expr } : v));
  return { ...s, numVars: numVars.length ? numVars : undefined };
}

// Bật: điền số hiện tại để sửa thành biến. Tắt: ô đang là số nguyên không âm thì giữ số đó làm giá trị cố định.
export function toggleNumVar(s: Step, field: NumVarField, on: boolean): Step {
  if (on) return withNumVar(s, field, String(s[field] ?? 0));
  const expr = s.numVars?.find((v) => v.field === field)?.expr.trim() ?? '';
  const out = withNumVar(s, field, undefined);
  return /^\d+$/.test(expr) ? { ...out, [field]: Number(expr) } : out;
}
