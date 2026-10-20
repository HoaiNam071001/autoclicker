// Chọn tool MCP cho từng loại AI (bước 31b2) – phần tính thuần, khớp store.MCPSettings.OffToolsFor (Go).
import type { store } from '../../../wailsjs/go/models';

type Mcp = Pick<store.MCPSettings, 'offTools' | 'clients'>;
type Patch = { offTools?: string[]; clients?: Record<string, string[]> };

/** Tool đang tắt cho client ('' = mặc định chung): ghi đè riêng nếu có, không thì mặc định. */
export function offFor(m: Mcp, client: string): string[] {
  return (client && m.clients?.[client]) || m.offTools || [];
}

/** Có ghi đè riêng cho client không. */
export const hasOverride = (m: Mcp, client: string): boolean => !!client && !!m.clients && client in m.clients;

/** Bật / tắt một tool cho client; lần đầu sửa client thì tạo ghi đè từ mặc định hiện tại. */
export function setToolOn(m: Mcp, client: string, tool: string, on: boolean): Patch {
  const next = offFor(m, client).filter((t) => t !== tool);
  if (!on) next.push(tool);
  return client ? { clients: { ...m.clients, [client]: next } } : { offTools: next };
}

/** Bỏ ghi đè của client – quay về mặc định chung. */
export function clearOverride(m: Mcp, client: string): Patch {
  const { [client]: _gone, ...rest } = m.clients ?? {};
  return { clients: rest };
}

/** Client để chọn: đã kết nối + đã có ghi đè (không trùng), theo chữ cái. */
export const clientChoices = (m: Mcp, seen: string[]): string[] =>
  [...new Set([...seen, ...Object.keys(m.clients ?? {})])].sort();
