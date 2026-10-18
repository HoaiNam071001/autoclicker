// Helper dùng chung cho find.js / probe.js / pick.js (find.go ghép vào trước mỗi script, chạy trong trang chính).
const clean = s => (s || '').replace(/\s+/g, ' ').trim();
const describe = el => {
  if (!el || !el.tagName) return 'phần tử khác';
  let s = el.tagName.toLowerCase();
  if (el.id) s += '#' + el.id;
  const t = clean(el.innerText).slice(0, 40);
  return t ? s + ' “' + t + '”' : s;
};
// Các document con cùng nguồn (iframe / frame) – khác nguồn thì contentDocument ném lỗi / null, bỏ qua.
const childDocs = doc => {
  const out = [];
  for (const f of doc.querySelectorAll('iframe, frame')) {
    try { if (f.contentDocument && f.contentDocument.documentElement) out.push({ frame: f, doc: f.contentDocument }); } catch (e) { /* khác nguồn */ }
  }
  return out;
};
// Góc trên trái vùng nội dung của frame trong viewport của document chứa nó (cộng viền + padding).
const frameContent = f => {
  const r = f.getBoundingClientRect(), cs = f.ownerDocument.defaultView.getComputedStyle(f);
  const x = r.left + f.clientLeft + parseFloat(cs.paddingLeft), y = r.top + f.clientTop + parseFloat(cs.paddingTop);
  return { x, y, w: f.clientWidth - parseFloat(cs.paddingLeft) - parseFloat(cs.paddingRight), h: f.clientHeight - parseFloat(cs.paddingTop) - parseFloat(cs.paddingBottom) };
};
// Độ lệch toạ độ client của window w so với trang chính (cộng dồn qua các frame lồng nhau).
const frameOffset = w => {
  let x = 0, y = 0;
  for (let f = w.frameElement; f; f = f.ownerDocument.defaultView.frameElement) {
    const c = frameContent(f);
    x += c.x; y += c.y;
  }
  return { x, y };
};
// so tagName thay vì instanceof: phần tử trong iframe thuộc realm khác
const tag = el => el.tagName;
// Chữ của phần tử theo cách người dùng thấy: nút input → value, ô nhập → placeholder.
const ownText = el => {
  if (tag(el) === 'INPUT') return ['button', 'submit', 'reset'].includes(el.type) ? el.value : el.placeholder;
  if (tag(el) === 'TEXTAREA') return el.placeholder;
  return el.innerText;
};
const ariaLabel = el => {
  const ids = (el.getAttribute('aria-labelledby') || '').split(/\s+/).filter(Boolean);
  if (ids.length) return ids.map(id => el.ownerDocument.getElementById(id)?.innerText || '').join(' ');
  return el.getAttribute('aria-label') || '';
};
// Vai trò ngầm định của thẻ (rút gọn theo HTML-ARIA; ô mật khẩu tính là textbox cho tiện viết form đăng nhập).
const IMPLICIT = {
  button: 'button, input[type=button], input[type=submit], input[type=reset], input[type=image], summary',
  link: 'a[href], area[href]',
  textbox: 'input:not([type]), input[type=text], input[type=email], input[type=tel], input[type=url], input[type=password], textarea',
  searchbox: 'input[type=search]', checkbox: 'input[type=checkbox]', radio: 'input[type=radio]', slider: 'input[type=range]',
  spinbutton: 'input[type=number]', combobox: 'select', option: 'option', heading: 'h1, h2, h3, h4, h5, h6',
  img: 'img[alt]:not([alt=""])', list: 'ul, ol', listitem: 'li', navigation: 'nav', main: 'main', dialog: 'dialog',
  form: 'form', article: 'article', table: 'table', row: 'tr', cell: 'td', columnheader: 'th',
};
// Tên truy cập rút gọn: aria → nhãn <label> → alt / title / placeholder → chữ.
const accName = el => {
  const a = ariaLabel(el);
  if (a) return a;
  if (el.labels && el.labels.length) return Array.from(el.labels).map(l => l.innerText).join(' ');
  return el.getAttribute('alt') || el.getAttribute('title') || el.getAttribute('placeholder') || ownText(el) || '';
};
// Vai trò của phần tử: role ghi rõ (từ đầu) → vai trò ngầm định của thẻ.
const roleOf = el => {
  // có thuộc tính role thì find.js không xét vai trò ngầm định (kể cả role rỗng)
  if (el.hasAttribute('role')) return el.getAttribute('role').trim().split(/\s+/)[0].toLowerCase();
  return Object.keys(IMPLICIT).find(k => el.matches(IMPLICIT[k])) || '';
};
// Giá trị dùng được trong selector sinh tự động (ngắn, không dấu nháy / biến).
const usable = s => s && s.length <= 60 && !s.includes('"') && !s.includes('{{');
// id / giá trị do framework sinh (":r1:", "x1n2onr6", "ember123") – đổi mỗi lần dựng trang, không dùng
const autoLike = s => /[:]|\d{3,}|^[a-z]{1,3}\d|[a-z0-9]{8,}/.test(s) && !/^[a-z]+(-[a-z]+)*$/.test(s);
// Đường CSS tới phần tử trong document của nó (dừng ở tổ tiên có id bền).
const cssPath = el => {
  const parts = [];
  for (let e = el; e && e.nodeType === Node.ELEMENT_NODE && tag(e) !== 'HTML'; e = e.parentElement) {
    if (e.id && !autoLike(e.id)) { parts.unshift('#' + CSS.escape(e.id)); break; }
    const same = Array.from(e.parentElement ? e.parentElement.children : []).filter(c => tag(c) === tag(e));
    const t = tag(e).toLowerCase();
    parts.unshift(same.length > 1 ? `${t}:nth-of-type(${same.indexOf(e) + 1})` : t);
  }
  return parts.join(' > ');
};
// Phần tử đang hiện (có cỡ, không ẩn / trong suốt) – find.js chỉ dùng phần tử thấy được.
const shown = el => {
  const r = el.getBoundingClientRect();
  return r.width >= 1 && r.height >= 1 && el.checkVisibility({ opacityProperty: true, visibilityProperty: true });
};
// document chính trước, rồi các iframe cùng nguồn (lồng nhau) theo thứ tự trong trang – thứ tự find.js tìm.
const allDocs = () => {
  const docs = [];
  const walk = doc => { docs.push(doc); for (const c of childDocs(doc)) walk(c.doc); };
  walk(document);
  return docs;
};
// Tâm (x, y – toạ độ client của frame chứa el) bị khung cuộn tổ tiên (overflow khác visible) cắt mất – vd cột trong
// dải cuộn ngang: elementFromPoint ra chính khung nên không bắt được là bị che.
const clipped = (el, x, y) => {
  const win = el.ownerDocument.defaultView;
  for (let p = el.parentElement; p && p !== el.ownerDocument.body && p !== el.ownerDocument.documentElement; p = p.parentElement) {
    const cs = win.getComputedStyle(p);
    if (cs.overflowX === 'visible' && cs.overflowY === 'visible') continue;
    const r = p.getBoundingClientRect();
    if (x < r.left || y < r.top || x >= r.right || y >= r.bottom) return true;
  }
  return false;
};
// Khung theo toạ độ trang chính + có nằm trong vùng nhìn thấy (viewport chính và vùng nội dung mọi frame chứa nó).
const place = el => {
  const r = el.getBoundingClientRect(), off = frameOffset(el.ownerDocument.defaultView);
  const box = { x: r.left + off.x, y: r.top + off.y, w: r.width, h: r.height };
  const cx = box.x + box.w / 2, cy = box.y + box.h / 2;
  let inView = cx >= 0 && cy >= 0 && cx < innerWidth && cy < innerHeight && !clipped(el, cx - off.x, cy - off.y);
  for (let f = el.ownerDocument.defaultView.frameElement; f && inView; f = f.ownerDocument.defaultView.frameElement) {
    const c = frameContent(f), o = frameOffset(f.ownerDocument.defaultView);
    inView = cx >= c.x + o.x && cy >= c.y + o.y && cx < c.x + o.x + c.w && cy < c.y + o.y + c.h;
  }
  return { box, cx, cy, inView };
};
