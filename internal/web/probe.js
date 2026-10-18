// Chạy trong trang (find.go, sau common.js): gắn bộ nghe pointermove (capture, trước mọi handler của trang) ở trang
// chính + mọi iframe cùng nguồn (gắn lại mỗi lần gọi để bắt frame mới), ghi lần cuối theo toạ độ client trang chính.
// Trả số lần + toạ độ (đo toạ độ – calibrate.go) và phần tử dưới chuột có thuộc window.__acTarget không (thử hover).
() => {
  const install = w => {
    if (!w.__acMvOn) {
      w.__acMvOn = true;
      w.addEventListener('pointermove', e => {
        const off = frameOffset(w), m = window.__acMv;
        window.__acMv = { X: e.clientX + off.x, Y: e.clientY + off.y, N: m.N + 1 };
        window.__acMvT = e.target;
      }, true);
    }
    for (const c of childDocs(w.document)) install(c.doc.defaultView);
  };
  if (!window.__acMv) window.__acMv = { X: 0, Y: 0, N: 0 };
  install(window);
  const t = window.__acTarget, at = window.__acMvT;
  const hit = !!(t && at && at.nodeType === Node.ELEMENT_NODE && (t === at || t.contains(at) || at.contains(t)));
  return { ...window.__acMv, Hit: hit, At: hit || !at ? '' : describe(at) };
}
