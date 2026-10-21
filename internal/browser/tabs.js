// Chạy bằng osascript -l JavaScript (tabs_darwin.go). argv = [lệnh, trình duyệt, tham số…] – không ghép vào chuỗi lệnh.
const IDS = { chrome: 'com.google.Chrome', edge: 'com.microsoft.edgemac', brave: 'com.brave.Browser', coccoc: 'com.coccoc.Coccoc' };

function app(b) {
  if (!IDS[b]) throw new Error('trình duyệt không hỗ trợ: ' + b);
  return Application(IDS[b]);
}

// running – không mở app (chưa cài → Application ném lỗi).
function running(b) {
  try { return Application(IDS[b]).running(); } catch (e) { return false; }
}

// list – mỗi trình duyệt đang chạy: id cửa sổ (trên cùng trước), tab đang hiện, id / tiêu đề / URL tab (mảng theo cửa sổ).
function list(names) {
  const out = [];
  for (const b of names) {
    if (!running(b)) continue;
    const a = app(b);
    out.push({
      browser: b, front: a.frontmost(), windows: a.windows.id().map(Number), active: a.windows.activeTabIndex(),
      ids: a.windows.tabs.id().map((w) => w.map(Number)), titles: a.windows.tabs.title(), urls: a.windows.tabs.url(), loading: a.windows.tabs.loading(),
    });
  }
  return out;
}

function tab(b, w, t) { return app(b).windows.byId(Number(w)).tabs.byId(Number(t)); }

function newTab(b, url) {
  const a = app(b);
  let w;
  if (a.windows.length === 0) {
    a.Window().make(); // cửa sổ mới có sẵn một tab trống – dùng luôn tab đó
    w = a.windows[0];
    if (url) w.activeTab.url = url;
  } else {
    w = a.windows[0];
    w.tabs.push(a.Tab(url ? { url } : {}));
  }
  a.activate();
  const t = w.activeTab;
  return { window: Number(w.id()), id: Number(t.id()), index: w.activeTabIndex(), title: t.title(), url: t.url() };
}

// Phím và thuộc tính cửa sổ đi qua System Events (cần quyền Accessibility – app đã xin cho click). Tab phải đang ở trước.
function winOf(b) { return Application('System Events').processes.whose({ bundleIdentifier: IDS[b] })[0].windows[0]; }

function isFullscreen(b) {
  try { return winOf(b).attributes.byName('AXFullScreen').value() === true; } catch (e) { return false; }
}

function frame(b, w) {
  const r = app(b).windows.byId(Number(w)).bounds();
  return { box: { x: r.x, y: r.y, w: r.width, h: r.height }, fullscreen: isFullscreen(b) };
}

function normalize(b, w, o) {
  let was = false;
  if (o.fullscreen && isFullscreen(b)) {
    was = true;
    winOf(b).attributes.byName('AXFullScreen').value = false;
    delay(1.2); // chờ hiệu ứng thoát toàn màn hình
  }
  if (o.bounds) app(b).windows.byId(Number(w)).bounds = { x: o.bounds.x, y: o.bounds.y, width: o.bounds.w, height: o.bounds.h };
  const se = Application('System Events');
  if (o.zoom) se.keystroke('0', { using: 'command down' });
  if (o.popups) se.keyCode(53); // Esc
  return { wasFullscreen: was };
}

function run(argv) {
  const [cmd, b, w, t, url] = argv;
  switch (cmd) {
    case 'list':
      return JSON.stringify(list(b ? [b] : Object.keys(IDS)));
    case 'new':
      return JSON.stringify(newTab(b, w));
    case 'activate': {
      const a = app(b);
      const win = a.windows.byId(Number(w));
      const k = win.tabs.id().map(Number).indexOf(Number(t));
      if (k < 0) throw new Error('tab gone (-1728)');
      win.activeTabIndex = k + 1;
      win.index = 1;
      a.activate();
      return '';
    }
    case 'go': tab(b, w, t).url = url; return '';
    case 'frame': return JSON.stringify(frame(b, w));
    case 'normalize': return JSON.stringify(normalize(b, w, JSON.parse(url))); // url = tuỳ chọn dạng JSON
    case 'close': tab(b, w, t).close(); return '';
    case 'reload': tab(b, w, t).reload(); return '';
    case 'back': tab(b, w, t).goBack(); return '';
    case 'forward': tab(b, w, t).goForward(); return '';
  }
  throw new Error('lệnh không hỗ trợ: ' + cmd);
}
