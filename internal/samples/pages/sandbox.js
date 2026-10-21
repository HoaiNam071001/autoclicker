// Sân thử Ghosthand – hàm dùng chung: ô đếm, nhật ký, dựng header.
(function () {
  const sb = {};
  // header(title): thanh trên cùng có tên trang, link về danh sách, nút Làm lại.
  sb.header = function (title) {
    const h = document.createElement('header');
    h.className = 'sb';
    h.innerHTML = '<a href="index.html">← Sân thử</a><h1></h1><span class="sp"></span><button type="button">Làm lại</button>';
    h.querySelector('h1').textContent = title;
    h.querySelector('button').onclick = () => location.reload();
    document.body.prepend(h);
    document.title = title + ' – Sân thử';
  };
  // stat(id, label): thêm ô đếm vào .stats, trả hàm set(value).
  sb.stat = function (label, value) {
    const box = document.querySelector('.stats');
    const el = document.createElement('div');
    el.className = 'stat';
    el.innerHTML = '<b></b><span></span>';
    el.querySelector('span').textContent = label;
    const b = el.querySelector('b');
    b.textContent = value ?? 0;
    box.appendChild(el);
    return (v) => { b.textContent = v; };
  };
  // log(text): thêm dòng (mới nhất trên cùng) vào .log, kèm giờ:phút:giây.mili.
  sb.log = function (text) {
    const box = document.querySelector('.log');
    if (!box) return;
    const d = new Date();
    const t = d.toTimeString().slice(0, 8) + '.' + String(d.getMilliseconds()).padStart(3, '0');
    const line = document.createElement('div');
    line.textContent = t + '  ' + text;
    box.prepend(line);
    while (box.children.length > 60) box.lastChild.remove();
  };
  sb.rand = (min, max) => min + Math.random() * (max - min);
  window.sb = sb;
})();
