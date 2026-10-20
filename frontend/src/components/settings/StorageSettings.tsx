// Cài đặt → Lưu trữ (bước 41 L2): thư mục kịch bản (lưu tay / AI tạo) và ảnh chụp. Ô trống = mặc định (hiện mờ).
// Chỉ hiển thị – AppSettings chọn / mở thư mục qua Go và lưu.
import { store } from '../../../wailsjs/go/models';

export type StorageKey = keyof store.StorageSettings;

const ROWS: { key: StorageKey; title: string; hint: string }[] = [
  { key: 'scriptsDir', title: 'Kịch bản', hint: 'Lưu lần đầu và Mở kịch bản mở sẵn ở đây' },
  { key: 'aiScriptsDir', title: 'Kịch bản AI tạo', hint: 'AI tạo kịch bản mới qua MCP vào đây' },
  { key: 'shotsDir', title: 'Ảnh chụp màn hình', hint: 'Node Chụp màn hình để trống thư mục lưu vào đây' },
];

export default function StorageSettings({ storage, defaults, onPick, onOpen, onReset }: {
  storage: store.StorageSettings; defaults: store.StorageDefaults | null;
  onPick: (key: StorageKey, from: string) => void; onOpen: (dir: string) => void; onReset: (key: StorageKey) => void;
}) {
  const fallback = (key: StorageKey) =>
    key === 'aiScriptsDir' ? storage.scriptsDir || defaults?.scriptsDir || '' : key === 'shotsDir' ? defaults?.shotsDir ?? '' : defaults?.scriptsDir ?? '';
  const btn = 'h-7 rounded-control border border-line px-2.5 text-xs hover:bg-hover';
  return (
    <>
      {ROWS.map(({ key, title, hint }, i) => {
        const custom = storage[key] ?? '';
        const dir = custom || fallback(key);
        return (
          <div key={key} className={`flex flex-col gap-1.5 ${i ? 'border-t border-line pt-3' : ''}`}>
            <span className="flex flex-col">{title}<span className="text-xs text-muted">{hint}</span></span>
            <code title={dir} className={`truncate rounded-control bg-field px-2 py-1.5 font-mono text-xs ${custom ? 'text-ink' : 'text-muted'}`}>
              {custom || `Mặc định${key === 'aiScriptsDir' ? ' (theo Kịch bản)' : ''}: ${dir}`}
            </code>
            <span className="flex gap-1.5">
              <button onClick={() => onPick(key, dir)} className={btn}>Chọn…</button>
              <button onClick={() => onOpen(dir)} disabled={!dir} className={btn} title="Mở trong Finder / Explorer">Mở</button>
              {custom && <button onClick={() => onReset(key)} className={btn}>Về mặc định</button>}
            </span>
          </div>
        );
      })}
      <span className="text-xs text-muted">Kịch bản đã có file thì Lưu vẫn ghi tại chỗ cũ; node Chụp màn hình đã đặt thư mục riêng thì giữ.</span>
    </>
  );
}
