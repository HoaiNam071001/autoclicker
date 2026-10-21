// Công tắc AI 4 mức (bước 31b): Tắt = tắt MCP; Chỉ soạn / Hỏi / Tự chạy = MCPSettings.run (khớp store.AIRun* trong Go).
// Dùng chung Cài đặt (AiRunSettings) và chip AI trên thanh mini.
export type RunMode = 'edit' | 'ask' | 'auto';
export type AiLevel = 'off' | RunMode;
type Mcp = { enabled: boolean; run?: string };

export const AI_LEVELS: { id: AiLevel; label: string }[] = [
  { id: 'off', label: 'Tắt' },
  { id: 'edit', label: 'Chỉ soạn' },
  { id: 'ask', label: 'Hỏi' },
  { id: 'auto', label: 'Tự chạy' },
];

// runMode – mức chạy khi MCP bật; giá trị lạ / trống = hỏi (mặc định của Go).
export const runMode = (m: Mcp): RunMode => (m.run === 'edit' || m.run === 'auto' ? m.run : 'ask');

export const aiLevel = (m: Mcp): AiLevel => (m.enabled ? runMode(m) : 'off');

// withLevel – cài đặt MCP sau khi chọn mức l (Tắt giữ nguyên run để bật lại về mức cũ).
export const withLevel = <T extends Mcp>(m: T, l: AiLevel): T => (l === 'off' ? { ...m, enabled: false } : { ...m, enabled: true, run: l });

export const levelLabel = (l: AiLevel) => AI_LEVELS.find((x) => x.id === l)!.label;

// Ưu tiên node khi AI soạn (bước 44, khớp store.NodePref* trong Go): giá trị lạ / trống = ảnh / chữ.
export type NodePref = 'vision' | 'web';
export const NODE_PREFS: { id: NodePref; label: string }[] = [
  { id: 'vision', label: 'Ảnh / chữ' },
  { id: 'web', label: 'Trình duyệt app' },
];
export const nodePref = (m: { nodePref?: string }): NodePref => (m.nodePref === 'web' ? 'web' : 'vision');
export const nodePrefLabel = (p: string) => NODE_PREFS.find((x) => x.id === nodePref({ nodePref: p }))!.label;
