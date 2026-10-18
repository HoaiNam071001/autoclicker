#!/bin/sh
# Hook SessionStart (.claude/settings.json): in trạng thái dự án để Claude khỏi đọc cả docs/tasks/README.md.
cd "$(dirname "$0")/.." || exit 0
echo "# Trạng thái dự án (tự chèn – scripts/session-status.sh)"
sed -n '/^## Trạng thái hiện tại/,/^## /p' docs/tasks/README.md | sed '1d;$d' | grep -v '^> ' | grep -v '^$' | sed 's|(phase3/|(docs/tasks/phase3/|g; s|(checks.md)|(docs/tasks/checks.md)|g'
n=$(grep -c '^- \[ \]' docs/tasks/checks.md 2>/dev/null)
echo "- **Chờ người dùng kiểm tra tay:** ${n:-0} việc (docs/tasks/checks.md)"
echo "- **Git:** nhánh $(git branch --show-current 2>/dev/null), $(git status --porcelain 2>/dev/null | wc -l | tr -d ' ') file chưa commit; commit gần nhất: $(git log -1 --format='%h %s' 2>/dev/null)"
