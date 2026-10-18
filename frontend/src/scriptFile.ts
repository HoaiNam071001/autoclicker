// Lưu / mở kịch bản (bước 21) – nối scriptStore với storage.go: Lưu / Lưu thành / Mở / Kịch bản mới / Mở gần đây,
// bản nháp tự lưu (mở lại app vẫn còn) và hỏi lưu trước khi bỏ thay đổi (đổi kịch bản, đóng app).
import { useSyncExternalStore } from 'react';
import { EventsOn } from '../wailsjs/runtime/runtime';
import {
  ClearDraft, ConfirmQuit, LoadDraft, NewScript, OpenSample, OpenScript, OpenScriptPath, RecentScripts, SaveDraft, SaveScript, SaveScriptAs, SetDirty,
} from '../wailsjs/go/main/App';
import { model, store } from '../wailsjs/go/models';
import { getEditor, markSaved, replaceScript, subscribe } from './store';
import { inWails, setWindowMode } from './windowMode';
import { Script } from './types';

const DRAFT_DELAY_MS = 800;

// Việc đang chờ người dùng trả lời "Lưu thay đổi?" – chạy sau khi Lưu / Không lưu.
type Pending = { action: 'new' } | { action: 'open' } | { action: 'recent'; path: string } | { action: 'sample'; id: string; name: string }
  | { action: 'quit' };

type FileState = {
  recent: string[]; // mở gần đây, mới nhất trước
  draftAt: number | null; // lần cuối ghi bản nháp (ms)
  error: string; // lỗi lưu / mở gần nhất
  pending: Pending | null; // đang hiện hộp thoại hỏi lưu
};

let fs: FileState = { recent: [], draftAt: null, error: '', pending: null };
const listeners = new Set<() => void>();
const setFs = (patch: Partial<FileState>) => {
  fs = { ...fs, ...patch };
  listeners.forEach((l) => l());
};
export const useFile = () =>
  useSyncExternalStore((l) => {
    listeners.add(l);
    return () => listeners.delete(l);
  }, () => fs);

// Struct Go → Script của frontend (class Wails có region?: Rect, anchorMode?: string).
const fromGo = (sc: model.Script) => ({ ...sc, region: sc.region ?? null }) as unknown as Script;
const toGo = (sc: Script) => model.Script.createFrom(sc);
const errText = (e: unknown) => (e instanceof Error ? e.message : String(e));

export const fileName = (path: string) => path.split(/[\\/]/).pop() ?? path;

const refreshRecent = () => RecentScripts().then((recent) => setFs({ recent })).catch(() => {});

// ---- Lưu ----

// Trả false nếu người dùng huỷ hộp thoại hoặc lỗi – việc đang chờ (đổi kịch bản, đóng app) không chạy tiếp.
export async function save(as = false): Promise<boolean> {
  if (!inWails) return false;
  const { script, path } = getEditor();
  try {
    const res = await (as ? SaveScriptAs : SaveScript)(toGo(script), path ?? '');
    if (!res.path) return false;
    // còn sửa trong lúc hộp thoại mở thì vẫn "chưa lưu" phần sửa đó
    markSaved(res.path, res.script.name, `Đã lưu ${fileName(res.path)}`, getEditor().script !== script);
    setFs({ error: '' });
    refreshRecent();
    await flushDraft(); // nháp khớp file ngay – đóng app liền sau đó cũng không bị báo "chưa lưu"
    return true;
  } catch (e) {
    setFs({ error: errText(e) });
    return false;
  }
}

// ---- Đổi kịch bản: hỏi lưu nếu có thay đổi ----

// discard: đóng app chọn "Không lưu" thì bỏ luôn bản nháp – lần sau mở kịch bản trống.
async function run(p: Pending, discard = false) {
  try {
    if (p.action === 'quit') {
      if (discard) await ClearDraft().catch(() => {});
      return ConfirmQuit();
    }
    if (p.action === 'new') return replaceScript(fromGo(await NewScript()), null, false, 'Kịch bản mới');
    // mẫu thử (bước 33): bản sao chưa có file – Lưu sẽ hỏi chỗ lưu, mẫu gốc không đổi
    if (p.action === 'sample') return replaceScript(fromGo(await OpenSample(p.id)), null, false, `Đã mở mẫu “${p.name}”`);
    const res = p.action === 'open' ? await OpenScript() : await OpenScriptPath(p.path);
    if (!res.path) return; // huỷ hộp thoại
    replaceScript(fromGo(res.script), res.path, false, `Đã mở ${fileName(res.path)}`);
    setFs({ error: '' });
  } catch (e) {
    setFs({ error: errText(e) });
  } finally {
    if (p.action !== 'quit') refreshRecent();
  }
}

function guard(p: Pending) {
  if (!inWails) return;
  if (getEditor().dirty) setFs({ pending: p });
  else run(p);
}

export const newScript = () => guard({ action: 'new' });
export const openScript = () => guard({ action: 'open' });
export const openRecent = (path: string) => guard({ action: 'recent', path });
export const openSample = (id: string, name: string) => guard({ action: 'sample', id, name });

// Trả lời hộp thoại hỏi lưu.
export async function answerPending(choice: 'save' | 'discard' | 'cancel') {
  const p = fs.pending;
  setFs({ pending: null });
  if (!p || choice === 'cancel') return;
  if (choice === 'save' && !(await save())) return;
  run(p, choice === 'discard');
}

// ---- Bản nháp tự lưu + báo Go có thay đổi chưa lưu ----

let draftTimer: ReturnType<typeof setTimeout> | undefined;
let lastScript: Script | null = null;
let lastDirty = false;

async function flushDraft() {
  clearTimeout(draftTimer);
  const e = getEditor();
  try {
    await SaveDraft(store.Draft.createFrom({ script: e.script, path: e.path ?? '', dirty: e.dirty, savedAt: 0 }));
    setFs({ draftAt: Date.now() });
  } catch (err) {
    setFs({ error: `Không ghi được bản nháp: ${errText(err)}` });
  }
}

function watch() {
  const { script, path, dirty } = getEditor();
  if (dirty !== lastDirty) {
    lastDirty = dirty;
    SetDirty(dirty).catch(() => {});
  }
  if (script === lastScript) return;
  lastScript = script;
  clearTimeout(draftTimer);
  draftTimer = setTimeout(flushDraft, DRAFT_DELAY_MS);
}

// Mở app: khôi phục bản nháp lần trước, không có thì kịch bản trống. Chỉ bắt đầu ghi nháp sau khi nạp xong
// để kịch bản tạm lúc khởi động không đè lên nháp cũ.
async function init() {
  try {
    const d = await LoadDraft();
    if (d) {
      replaceScript(fromGo(d.script), d.path || null, d.dirty, d.dirty ? 'Đã khôi phục bản nháp chưa lưu' : '');
      setFs({ draftAt: d.savedAt || null });
    } else {
      replaceScript(fromGo(await NewScript()), null, false);
    }
  } catch (e) {
    setFs({ error: errText(e) });
  }
  lastScript = getEditor().script;
  lastDirty = !getEditor().dirty; // gửi trạng thái ban đầu cho Go
  watch();
  subscribe(watch);
  refreshRecent();
}

if (inWails) {
  init();
  // Đóng app (nút ✕, Cmd+Q) khi còn thay đổi chưa lưu: Go giữ app lại và nhờ hỏi.
  EventsOn('app:confirmClose', () => {
    if (location.hash) setWindowMode('main');
    setFs({ pending: { action: 'quit' } });
  });
}
