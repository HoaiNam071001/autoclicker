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
