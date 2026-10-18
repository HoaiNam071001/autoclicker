#!/usr/bin/env python3
"""Cập nhật bảng tổng quan trong docs/tasks/README.md (giữa <!-- index:start --> và <!-- index:end -->):
mỗi bước một dòng – trạng thái, tiến độ mục con, số việc chờ kiểm tra tay, link file bước.
Chạy lại sau khi tick: python3 scripts/tasks-index.py
Bước dài: mục con đã xong có thể nằm ở phase3/<tên>-xong.md – vẫn được đếm vào bước."""
import glob
import os
import re

ROOT = os.path.join(os.path.dirname(__file__), '..', 'docs', 'tasks')
README = os.path.join(ROOT, 'README.md')


def steps(rel):
    out, cur = [], None
    for line in open(os.path.join(ROOT, rel), encoding='utf-8'):
        m = re.match(r'- \[(.)\] \*\*(\d+[a-z]?)\. (.+?)\*\*', line)
        if m:
            cur = dict(done=m.group(1) == 'x', num=m.group(2), title=m.group(3).rstrip('.'), x=0, n=0, file=rel)
            out.append(cur)
            continue
        m = re.match(r'\s+- \[(.)\]', line)
        if m and cur:
            cur['n'] += 1
            cur['x'] += m.group(1) == 'x'
    # Phần đã xong của bước dài tách sang <tên>-xong.md: mục con trong đó cộng vào bước của file chính.
    done = os.path.join(ROOT, rel[:-3] + '-xong.md')
    if out and os.path.exists(done):
        for line in open(done, encoding='utf-8'):
            m = re.match(r'\s+- \[(.)\]', line)
            if m:
                out[-1]['n'] += 1
                out[-1]['x'] += m.group(1) == 'x'
    return out


def pending_checks():
    """Số việc kiểm tra tay chưa tick theo bước (tiêu đề ## <số> … trong checks.md)."""
    out, num = {}, None
    for line in open(os.path.join(ROOT, 'checks.md'), encoding='utf-8'):
        m = re.match(r'## (\d+[a-z]?)', line)
        if m:
            num = m.group(1)
        elif line.startswith('## '):
            num = None
        elif num and line.startswith('- [ ]'):
            out[num] = out.get(num, 0) + 1
    return out


def row(s, checks):
    c = checks.get(s['num'], 0)
    icon = '🧪' if s['done'] and c else '✅' if s['done'] else '🟡' if s['x'] else '⬜'
    state = {'🧪': 'chờ kiểm tra tay', '✅': 'xong', '🟡': 'đang làm', '⬜': 'chưa làm'}[icon]
    prog = f"{s['x']} / {s['n']}" if s['n'] else '–'
    chk = f"[{c} việc](checks.md)" if c else ''
    return f"| {icon} {state} | [{s['num']}. {s['title']}]({s['file']}) | {prog} | {chk} |"


def main():
    text = open(README, encoding='utf-8').read()
    checks = pending_checks()
    all_steps = {}
    for f in ['phase0.md', 'phase1-mvp.md', 'phase2.md'] + sorted(os.path.relpath(p, ROOT) for p in glob.glob(os.path.join(ROOT, 'phase3', '*.md')) if not p.endswith('-xong.md')):
        for s in steps(f):
            all_steps[s['num']] = s
    # Phase 3 theo thứ tự làm: thứ tự các link phase3/… trong bảng "Thứ tự làm"
    order = []
    for f in re.findall(r'\]\((phase3/[^)]+\.md)\)', text.split('<!-- index:end -->')[-1]):
        for s in all_steps.values():
            if s['file'] == f and s['num'] not in order:
                order.append(s['num'])
    head = '| Trạng thái | Bước | Mục con xong | Kiểm tra tay |\n|---|---|---|---|'
    phase = lambda files: [s for s in all_steps.values() if s['file'] in files]
    p3 = [all_steps[n] for n in order] + [s for s in all_steps.values() if s['file'].startswith('phase3/') and s['num'] not in order]
    total = sum(checks.values())
    parts = [
        f"**{sum(s['done'] for s in all_steps.values())} / {len(all_steps)} bước xong code** · "
        f"🧪 **[{total} việc chờ bạn kiểm tra tay](checks.md)** · ✅ xong · 🧪 xong code, chờ kiểm tra tay · 🟡 đang làm · ⬜ chưa làm",
        '', '### Phase 3 – Design v2 & tự động hoá (theo thứ tự làm)', '', head, *[row(s, checks) for s in p3],
        '', '### Phase 0 – 2', '', head, *[row(s, checks) for s in phase({'phase0.md', 'phase1-mvp.md', 'phase2.md'})],
    ]
    new = re.sub(r'(<!-- index:start -->\n).*?(<!-- index:end -->)', lambda m: m.group(1) + '\n'.join(parts) + '\n' + m.group(2), text, flags=re.S)
    open(README, 'w', encoding='utf-8').write(new)


if __name__ == '__main__':
    main()
