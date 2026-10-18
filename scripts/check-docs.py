#!/usr/bin/env python3
"""Kiểm tra tài liệu / quy tắc lệch với code. Chạy: python3 scripts/check-docs.py (thoát 1 nếu có lỗi).
- Mọi method public của App và mọi event EventsEmit phải có trong docs/ARCHITECTURE.md.
- Mọi loại node trong stepTypes (model.go) phải có nhãn (types.ts), icon (icons.tsx), mô tả (steptypes.go).
- Cảnh báo (không lỗi): file code > 400 dòng (quy tắc *Một file một việc*)."""
import glob
import os
import re
import sys

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..')
read = lambda p: open(os.path.join(ROOT, p), encoding='utf-8').read()

go_files = [f for f in glob.glob(os.path.join(ROOT, '*.go')) if not f.endswith('_test.go')]
go_src = {os.path.basename(f): open(f, encoding='utf-8').read() for f in go_files}
arch = read('docs/ARCHITECTURE.md')
errors, warns = [], []

# Method public của App (bind sang frontend). Bỏ Greet (mẫu Wails).
for name, src in go_src.items():
    for m in re.finditer(r'^func \(a \*App\) ([A-Z]\w*)\(', src, re.M):
        if m.group(1) != 'Greet' and f'`{m.group(1)}' not in arch and m.group(1) not in arch:
            errors.append(f'ARCHITECTURE.md thiếu method App.{m.group(1)} ({name})')

# Event Go → frontend.
for name, src in go_src.items():
    for m in re.finditer(r'EventsEmit\([^,]+,\s*"([^"]+)"', src):
        if m.group(1) not in arch:
            errors.append(f'ARCHITECTURE.md thiếu event "{m.group(1)}" ({name})')

# Loại node: hằng trong stepTypes → giá trị chuỗi.
model = go_src['model.go']
consts = dict(re.findall(r'^\s*(Step\w+)\s*=\s*"(\w+)"', model, re.M))
listed = re.search(r'var stepTypes = \[\]string\{(.+?)\}', model, re.S)
types = [consts[c] for c in re.findall(r'Step\w+', listed.group(1)) if c in consts] if listed else []
for rel in ('frontend/src/types.ts', 'frontend/src/components/icons.tsx', 'steptypes.go'):
    src = read(rel)
    for t in types:
        if not re.search(rf"(^|[\s{{,'\"])'?{t}'?\s*:|Step{t[0].upper() + t[1:]}\b|\"{t}\"", src, re.M):
            errors.append(f'{rel} thiếu loại node "{t}" (nhãn / icon / mô tả – checklist Thêm loại node)')

# File dài.
code = go_files + glob.glob(os.path.join(ROOT, 'frontend/src/**/*.ts*'), recursive=True)
for f in sorted(code):
    if '/demo/' in f or f.endswith('.test.ts'):
        continue
    n = sum(1 for _ in open(f, encoding='utf-8'))
    if n > 400:
        warns.append(f'{os.path.relpath(f, ROOT)}: {n} dòng (> 400) – cân nhắc tách khi đang sửa file này')

for w in warns:
    print('cảnh báo:', w)
for e in errors:
    print('LỖI:', e)
print(f'{len(errors)} lỗi, {len(warns)} cảnh báo')
sys.exit(1 if errors else 0)
