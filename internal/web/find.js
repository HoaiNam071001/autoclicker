// Chạy trong trang (CDP Runtime.evaluate, find.go): tìm phần tử theo danh sách selector {by, value, exact}, cái đầu thấy
// phần tử hiện thì dùng. scroll = phần tử ngoài viewport thì cuộn vào giữa. Trả khung theo toạ độ client (CSS px),
// phần tử che tâm (elementFromPoint) và hình học cửa sổ để Go biết khi nào phải đo lại toạ độ.
(sels, scroll) => {
  const clean = s => (s || '').replace(/\s+/g, ' ').trim();
  const norm = s => clean(s).toLowerCase();
  const matches = (got, sel) => {
    const a = norm(got), b = norm(sel.value);
    return a !== '' && (sel.exact ? a === b : a.includes(b));
  };
  const SKIP = new Set(['SCRIPT', 'STYLE', 'NOSCRIPT', 'TEMPLATE', 'HEAD']);
  const all = () => Array.from(document.body ? document.body.querySelectorAll('*') : []).filter(el => !SKIP.has(el.tagName));
  const ownText = el => {
    if (el instanceof HTMLInputElement) return ['button', 'submit', 'reset'].includes(el.type) ? el.value : el.placeholder;
    if (el instanceof HTMLTextAreaElement) return el.placeholder;
    return el.innerText;
  };
  // Phần tử sâu nhất chứa chữ: bỏ phần tử có con cũng khớp (chữ trong <button><span>Gửi</span></button> → span).
  const deepest = list => { const set = new Set(list); return list.filter(el => !Array.from(el.children).some(c => set.has(c))); };
  const inDocOrder = list => Array.from(new Set(list)).sort((a, b) => a.compareDocumentPosition(b) & Node.DOCUMENT_POSITION_FOLLOWING ? -1 : 1);
  // Nhanh: text node chứa chữ → phần tử cha (+ ô nhập theo value / placeholder). Chữ bị tách qua nhiều thẻ
  // ("Xem <b>lịch</b> sử") thì mới quét mọi phần tử (chậm trên trang lớn).
  const byText = sel => {
    const loose = { value: sel.value, exact: false }, hits = [];
    const w = document.createTreeWalker(document.body || document.documentElement, NodeFilter.SHOW_TEXT);
    for (let n; (n = w.nextNode());) {
      const p = n.parentElement;
      if (p && !SKIP.has(p.tagName) && matches(n.data, loose) && matches(ownText(p), sel)) hits.push(p);
    }
    for (const el of document.querySelectorAll('input, textarea')) if (matches(ownText(el), sel)) hits.push(el);
    if (hits.length) return inDocOrder(hits);
    return deepest(all().filter(el => matches(el.textContent, loose) && matches(ownText(el), sel)));
  };
  const label = el => {
    const ids = (el.getAttribute('aria-labelledby') || '').split(/\s+/).filter(Boolean);
    if (ids.length) return ids.map(id => document.getElementById(id)?.innerText || '').join(' ');
    return el.getAttribute('aria-label') || '';
  };
  const byAria = sel => all().filter(el => (el.hasAttribute('aria-label') || el.hasAttribute('aria-labelledby')) && matches(label(el), sel));
  const byXPath = v => {
    const r = document.evaluate(v, document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
    return Array.from({ length: r.snapshotLength }, (_, i) => r.snapshotItem(i)).filter(n => n instanceof Element);
  };
  const find = sel => {
    switch (sel.by) {
      case 'text': return byText(sel);
      case 'aria': return byAria(sel);
      case 'css': return Array.from(document.querySelectorAll(sel.value));
      case 'xpath': return byXPath(sel.value);
    }
    throw new Error('cách tìm không hợp lệ: ' + sel.by);
  };
  const shown = el => {
    const r = el.getBoundingClientRect();
    return r.width >= 1 && r.height >= 1 && el.checkVisibility({ opacityProperty: true, visibilityProperty: true });
  };
  const describe = el => {
    let s = el.tagName.toLowerCase();
    if (el.id) s += '#' + el.id;
    const t = clean(el.innerText).slice(0, 40);
    return t ? s + ' “' + t + '”' : s;
  };
  const geom = { screenX, screenY, outerW: outerWidth, outerH: outerHeight, innerW: innerWidth, innerH: innerHeight, dpr: devicePixelRatio };
  // đang tải (tab mới mở, chưa có body): không thấy chưa có nghĩa là đã mất
  const loading = document.readyState === 'loading' || location.href === 'about:blank' || !document.body;
  const out = { sel: -1, hidden: 0, geom, loading, focus: document.hasFocus(), visible: document.visibilityState === 'visible' };
  for (let i = 0; i < sels.length; i++) {
    const found = find(sels[i]);
    const el = found.find(shown);
    if (!el) { out.hidden = Math.max(out.hidden, found.length); continue; }
    let r = el.getBoundingClientRect();
    if (scroll && (r.top < 0 || r.left < 0 || r.bottom > innerHeight || r.right > innerWidth)) {
      el.scrollIntoView({ block: 'center', inline: 'center', behavior: 'instant' });
      r = el.getBoundingClientRect();
    }
    const cx = r.left + r.width / 2, cy = r.top + r.height / 2;
    const top = document.elementFromPoint(cx, cy);
    Object.assign(out, { sel: i, x: r.left, y: r.top, w: r.width, h: r.height, text: clean(ownText(el)).slice(0, 200) });
    // tâm ngoài viewport (chưa cuộn) → elementFromPoint null: chưa bấm được, nhưng webWait vẫn tính là thấy
    out.inView = cx >= 0 && cy >= 0 && cx < innerWidth && cy < innerHeight;
    if (out.inView && top && top !== el && !el.contains(top) && !top.contains(el)) out.covered = describe(top);
    return out;
  }
  return out;
}
