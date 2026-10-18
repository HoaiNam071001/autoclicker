// Chạy trong trang (CDP Runtime.evaluate, find.go, sau common.js): tìm phần tử theo danh sách selector
// {by, value, exact, name}, cái đầu thấy phần tử hiện thì dùng – cả trong iframe cùng nguồn. scroll = phần tử ngoài
// viewport thì cuộn vào giữa. Trả khung theo toạ độ client của trang chính (CSS px), phần tử che tâm (elementFromPoint
// từng tầng frame) và hình học cửa sổ để Go biết khi nào phải đo lại toạ độ. Phần tử thấy được ghi vào
// window.__acTarget để probe.js kiểm khi thử hover.
(sels, scroll) => {
  const norm = s => clean(s).toLowerCase();
  const matches = (got, value, exact) => {
    const a = norm(got), b = norm(value);
    return a !== '' && (exact ? a === b : a.includes(b));
  };
  const SKIP = new Set(['SCRIPT', 'STYLE', 'NOSCRIPT', 'TEMPLATE', 'HEAD']);
  const all = doc => Array.from(doc.body ? doc.body.querySelectorAll('*') : []).filter(el => !SKIP.has(tag(el)));
  // Phần tử sâu nhất chứa chữ: bỏ phần tử có con cũng khớp (chữ trong <button><span>Gửi</span></button> → span).
  const deepest = list => { const set = new Set(list); return list.filter(el => !Array.from(el.children).some(c => set.has(c))); };
  const inDocOrder = list => Array.from(new Set(list)).sort((a, b) => a.compareDocumentPosition(b) & Node.DOCUMENT_POSITION_FOLLOWING ? -1 : 1);
  // Nhanh: text node chứa chữ → phần tử cha (+ ô nhập theo value / placeholder). Chữ bị tách qua nhiều thẻ
  // ("Xem <b>lịch</b> sử") thì mới quét mọi phần tử (chậm trên trang lớn).
  const byText = (sel, doc) => {
    const hits = [];
    const w = doc.createTreeWalker(doc.body || doc.documentElement, NodeFilter.SHOW_TEXT);
    for (let n; (n = w.nextNode());) {
      const p = n.parentElement;
      if (p && !SKIP.has(tag(p)) && matches(n.data, sel.value, false) && matches(ownText(p), sel.value, sel.exact)) hits.push(p);
    }
    for (const el of doc.querySelectorAll('input, textarea')) if (matches(ownText(el), sel.value, sel.exact)) hits.push(el);
    if (hits.length) return inDocOrder(hits);
    return deepest(all(doc).filter(el => matches(el.textContent, sel.value, false) && matches(ownText(el), sel.value, sel.exact)));
  };
  const byAria = (sel, doc) => all(doc).filter(el => (el.hasAttribute('aria-label') || el.hasAttribute('aria-labelledby')) && matches(ariaLabel(el), sel.value, sel.exact));
  const byRole = (sel, doc) => {
    const explicit = all(doc).filter(el => el.hasAttribute('role') && el.getAttribute('role').trim().split(/\s+/)[0].toLowerCase() === sel.value);
    const implicit = IMPLICIT[sel.value] ? Array.from(doc.querySelectorAll(IMPLICIT[sel.value])).filter(el => !el.hasAttribute('role')) : [];
    const list = inDocOrder([...explicit, ...implicit]);
    return sel.name ? list.filter(el => matches(accName(el), sel.name, false)) : list;
  };
  const byXPath = (v, doc) => {
    const r = doc.evaluate(v, doc, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
    return Array.from({ length: r.snapshotLength }, (_, i) => r.snapshotItem(i)).filter(n => n.nodeType === Node.ELEMENT_NODE);
  };
  const findIn = (sel, doc) => {
    switch (sel.by) {
      case 'text': return byText(sel, doc);
      case 'aria': return byAria(sel, doc);
      case 'role': return byRole(sel, doc);
      case 'css': return Array.from(doc.querySelectorAll(sel.value));
      case 'xpath': return byXPath(sel.value, doc);
    }
    throw new Error('cách tìm không hợp lệ: ' + sel.by);
  };
  // document chính trước, rồi các iframe cùng nguồn (lồng nhau) theo thứ tự trong trang
  const docs = [];
  const walk = doc => { docs.push(doc); for (const c of childDocs(doc)) walk(c.doc); };
  walk(document);
  const shown = el => {
    const r = el.getBoundingClientRect();
    return r.width >= 1 && r.height >= 1 && el.checkVisibility({ opacityProperty: true, visibilityProperty: true });
  };
  // Khung theo toạ độ trang chính + có nằm trong vùng nhìn thấy (viewport chính và vùng nội dung mọi frame chứa nó).
  const place = el => {
    const r = el.getBoundingClientRect(), off = frameOffset(el.ownerDocument.defaultView);
    const box = { x: r.left + off.x, y: r.top + off.y, w: r.width, h: r.height };
    const cx = box.x + box.w / 2, cy = box.y + box.h / 2;
    let inView = cx >= 0 && cy >= 0 && cx < innerWidth && cy < innerHeight;
    for (let f = el.ownerDocument.defaultView.frameElement; f && inView; f = f.ownerDocument.defaultView.frameElement) {
      const c = frameContent(f), o = frameOffset(f.ownerDocument.defaultView);
      inView = cx >= c.x + o.x && cy >= c.y + o.y && cx < c.x + o.x + c.w && cy < c.y + o.y + c.h;
    }
    return { box, cx, cy, inView };
  };
  // Phần tử che tâm: xét từ trang chính xuống, mỗi tầng điểm đó phải rơi vào frame chứa đích (tầng cuối: chính đích).
  const coverOf = (el, cx, cy) => {
    const chain = [el];
    for (let f = el.ownerDocument.defaultView.frameElement; f; f = f.ownerDocument.defaultView.frameElement) chain.unshift(f);
    for (const want of chain) {
      const doc = want.ownerDocument, off = frameOffset(doc.defaultView);
      const top = doc.elementFromPoint(cx - off.x, cy - off.y);
      if (top && top !== want && !want.contains(top) && !top.contains(want)) return describe(top);
    }
    return '';
  };
  const geom = { screenX, screenY, outerW: outerWidth, outerH: outerHeight, innerW: innerWidth, innerH: innerHeight, dpr: devicePixelRatio };
  // đang tải (tab mới mở, chưa có body): không thấy chưa có nghĩa là đã mất
  const loading = document.readyState === 'loading' || location.href === 'about:blank' || !document.body;
  const out = { sel: -1, hidden: 0, geom, loading, focus: document.hasFocus(), visible: document.visibilityState === 'visible' };
  window.__acTarget = null;
  for (let i = 0; i < sels.length; i++) {
    const found = docs.flatMap(doc => findIn(sels[i], doc));
    const el = found.find(shown);
    if (!el) { out.hidden = Math.max(out.hidden, found.length); continue; }
    let p = place(el);
    if (scroll && !p.inView) {
      el.scrollIntoView({ block: 'center', inline: 'center', behavior: 'instant' });
      p = place(el);
    }
    window.__acTarget = el;
    Object.assign(out, { sel: i, x: p.box.x, y: p.box.y, w: p.box.w, h: p.box.h, text: clean(ownText(el)).slice(0, 200), inView: p.inView });
    // tâm ngoài vùng nhìn thấy (chưa cuộn) → chưa bấm được, nhưng webWait vẫn tính là thấy
    if (p.inView) out.covered = coverOf(el, p.cx, p.cy);
    return out;
  }
  return out;
}
