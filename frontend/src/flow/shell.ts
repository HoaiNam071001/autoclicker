// Node Chạy lệnh (bước 43 S1) – ngôn ngữ script, dòng chi tiết thẻ. Chạy ở internal/engine/shell.go; nội dung lệnh cần
// người dùng đồng ý (flow/danger.ts – script.DangerOps).
import type { Step } from '../types';

// '' = shell của máy (khớp model.ShellLangs)
export type ShellLang = '' | 'powershell' | 'applescript' | 'jxa';
export const SHELL_LANGS: { id: ShellLang; label: string; mac?: boolean }[] = [
  { id: '', label: 'Shell của máy (zsh / cmd)' }, { id: 'powershell', label: 'PowerShell' },
  { id: 'applescript', label: 'AppleScript (macOS)', mac: true }, { id: 'jxa', label: 'JXA – JavaScript (macOS)', mac: true },
];

// Dòng đầu của lệnh, cắt ngắn: `ls -1 "$DIR" | wc -l`
export function shellPreview(script?: string): string {
  const line = (script ?? '').split('\n').map((l) => l.trim()).find(Boolean) ?? '';
  return line.length > 48 ? `${line.slice(0, 47)}…` : line;
}

export function fmtShell(s: Pick<Step, 'shellLang' | 'script'>): string {
  const lang = SHELL_LANGS.find((l) => l.id === (s.shellLang ?? ''));
  const cmd = shellPreview(s.script) || '(chưa có lệnh)';
  return s.shellLang ? `${lang?.label.replace(/ \(.*/, '')}: ${cmd}` : cmd;
}
