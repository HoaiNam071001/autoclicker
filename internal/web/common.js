// Helper dùng chung cho find.js / probe.js (find.go ghép vào trước mỗi script, chạy trong trang chính).
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
