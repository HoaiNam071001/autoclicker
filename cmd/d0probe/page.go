//go:build probe

package main

// Trang thử: 5 nút 20×20 ở góc + giữa, 1 nút phải cuộn mới thấy.
// Ghi pointermove cuối (để đo) và pointerdown / click (để chấm sai số).
const probePage = `<!doctype html><html><head><meta charset="utf-8"><title>D0 probe</title><style>
html,body{margin:0;font:14px sans-serif}
body{height:3000px;background:#f4f4f4;position:relative}
.t{position:absolute;width:20px;height:20px;background:#2563eb;border:0;padding:0;cursor:pointer}
</style></head><body>
<button class="t" id="tl" style="left:40px;top:40px"></button>
<button class="t" id="tr" style="right:40px;top:40px"></button>
<button class="t" id="bl" style="left:40px;top:calc(100vh - 60px)"></button>
<button class="t" id="br" style="right:40px;top:calc(100vh - 60px)"></button>
<button class="t" id="c" style="left:calc(50% - 10px);top:calc(50vh - 10px)"></button>
<button class="t" id="far" style="left:300px;top:2500px"></button>
<script>
window.__mv = {x:0, y:0, n:0};
window.__downs = [];
addEventListener('pointermove', e => { __mv = {x:e.clientX, y:e.clientY, n:__mv.n+1}; }, true);
addEventListener('pointerdown', e => { __downs.push({id:e.target.id||'', x:e.clientX, y:e.clientY, trusted:e.isTrusted}); }, true);
addEventListener('click', e => e.preventDefault(), true);
</script></body></html>`

// targetIDs – thứ tự bấm; "far" cuộn vào giữa trước khi đo rect.
var targetIDs = []string{"tl", "tr", "c", "bl", "br", "far"}
