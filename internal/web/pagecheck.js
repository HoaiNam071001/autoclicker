// Chạy trong trang (pagecheck.go, sau common.js): gom dấu hiệu captcha / chặn bot / game canvas để Go phân loại.
// Chỉ đọc, không đụng trang.
() => {
  const vw = innerWidth, vh = innerHeight;
  const shown = el => {
    const r = el.getBoundingClientRect(), st = getComputedStyle(el);
    return r.width > 0 && r.height > 0 && st.visibility !== 'hidden' && st.display !== 'none' && +st.opacity !== 0;
  };
  const frames = [...document.querySelectorAll('iframe')].map(f => {
    const r = f.getBoundingClientRect();
    return { src: f.src || '', title: f.title || '', w: Math.round(r.width), h: Math.round(r.height), shown: shown(f) };
  });
  const marks = [];
  for (const sel of ['.g-recaptcha', '.h-captcha', '.cf-turnstile', '#challenge-form', '#cf-challenge-running', '#px-captcha']) {
    const el = document.querySelector(sel);
    if (el && shown(el)) marks.push(sel);
  }
  let canvas = 0; // phần khung nhìn canvas lớn nhất che (0–1)
  for (const c of document.querySelectorAll('canvas')) {
    if (!shown(c)) continue;
    const r = c.getBoundingClientRect();
    const w = Math.max(0, Math.min(r.right, vw) - Math.max(r.left, 0)), h = Math.max(0, Math.min(r.bottom, vh) - Math.max(r.top, 0));
    canvas = Math.max(canvas, (w * h) / (vw * vh || 1));
  }
  const text = (document.body ? document.body.innerText : '').slice(0, 3000);
  return { url: location.href, title: document.title, text, frames, marks, canvas };
}
