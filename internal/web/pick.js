// Chạy trong trang (pick.go, sau common.js): "Chọn phần tử trên trang" (D3). op = 'start': tô viền phần tử dưới chuột
// (trang chính + iframe cùng nguồn), click → chặn click của trang, ghi phần tử vào window.__acPicked + danh sách
// selector ứng viên theo thứ tự độ bền (chữ → aria → role → id / data-* → đường CSS); Esc → huỷ. op = 'read': trạng thái.
// op = 'stop': gỡ hết. Go kiểm lại từng ứng viên bằng find.js (chỉ giữ cái tìm ra đúng phần tử này).
(op) => {
  const st = window.__acPick;
  if (op === 'read') return st ? { state: st.state, cands: st.cands, text: st.text } : { state: 'gone' };
  if (st) st.stop();
  if (op === 'stop') return { state: 'stopped' };

  // Phần tử đáng chọn: phần tử tương tác gần nhất chứa điểm bấm (span trong nút → nút).
  const ACTIVE = 'a[href], button, input, select, textarea, summary, label, [role], [onclick], [tabindex]';
  const targetOf = el => (el.closest && el.closest(ACTIVE)) || el;
  const candidates = el => {
    const out = [];
    const text = clean(ownText(el) || '');
    if (usable(text)) out.push('text=' + text, 'text="' + text + '"');
    const aria = clean(ariaLabel(el));
    if (usable(aria)) out.push('aria=' + aria);
    const role = roleOf(el), name = clean(accName(el));
    if (role) out.push(usable(name) ? `role=${role}[name="${name}"]` : 'role=' + role);
    if (el.id && !autoLike(el.id)) out.push('css=#' + CSS.escape(el.id));
    for (const a of ['data-testid', 'data-test', 'data-qa', 'name']) {
      const v = el.getAttribute(a);
      if (v && usable(v)) out.push(`css=${tag(el).toLowerCase()}[${a}="${v}"]`);
    }
    out.push('css=' + cssPath(el));
    return out;
  };

  const box = document.createElement('div');
  box.style.cssText = 'position:fixed;z-index:2147483647;pointer-events:none;border:2px solid #f59e0b;background:rgba(245,158,11,.15);border-radius:3px;display:none';
  const bar = document.createElement('div');
  bar.textContent = 'Ghosthand: bấm vào phần tử cần chọn · Esc để huỷ';
  bar.style.cssText = 'position:fixed;z-index:2147483647;top:8px;left:50%;transform:translateX(-50%);pointer-events:none;padding:6px 12px;border-radius:8px;background:#111827;color:#f9fafb;font:13px system-ui,sans-serif;box-shadow:0 4px 16px rgba(0,0,0,.3)';
  document.documentElement.append(box, bar);

  const s = { state: 'picking', cands: [], text: '', hover: null };
  const show = el => {
    const r = el.getBoundingClientRect(), off = frameOffset(el.ownerDocument.defaultView);
    Object.assign(box.style, { display: 'block', left: r.left + off.x + 'px', top: r.top + off.y + 'px', width: r.width + 'px', height: r.height + 'px' });
  };
  const block = e => { e.preventDefault(); e.stopImmediatePropagation(); };
  const handlers = {
    pointermove: e => { s.hover = targetOf(e.target); show(s.hover); },
    pointerdown: block, mousedown: block, pointerup: block, mouseup: block,
    click: e => {
      block(e);
      const el = targetOf(e.target);
      window.__acPicked = el;
      Object.assign(s, { state: 'done', cands: candidates(el), text: clean(ownText(el) || '').slice(0, 200) });
      s.stop(true);
    },
    keydown: e => { if (e.key === 'Escape') { block(e); s.state = 'cancelled'; s.stop(true); } },
  };
  const wins = [];
  const install = w => {
    wins.push(w);
    for (const [k, f] of Object.entries(handlers)) w.addEventListener(k, f, true);
    for (const c of childDocs(w.document)) install(c.doc.defaultView);
  };
  install(window);
  // keep = giữ trạng thái cho Go đọc (chỉ gỡ giao diện + bộ nghe)
  s.stop = keep => {
    for (const w of wins) for (const [k, f] of Object.entries(handlers)) w.removeEventListener(k, f, true);
    box.remove(); bar.remove();
    if (!keep) window.__acPick = null;
  };
  window.__acPick = s;
  return { state: s.state };
}
