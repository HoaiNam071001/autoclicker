# Dark mode · xám ám xanh

Có 3 biến thể dùng chung một bộ token. Đổi giữa chúng bằng `data-theme` trên `<html>`.

| Theme | `data-theme` | Đặc điểm |
|---|---|---|
| Slate | `slate` | Tối nhất, xanh rất nhẹ, gọn và hiện đại |
| Dimmed | `dimmed` | Sáng và dịu nhất (giống GitHub Dark dimmed) |
| Nord | `nord` | Xanh rõ nhất, accent xanh băng, cảm giác mát |

## Bảng màu nền

| Token | Vai trò | Slate | Dimmed | Nord |
|---|---|---|---|---|
| `--bg` | Nền trang | `#171b22` | `#22272e` | `#1e2430` |
| `--surface-1` | Card, widget | `#222833` | `#2d333b` | `#2a3240` |
| `--surface-2` | Popover, hover, input, track | `#2e3644` | `#373e47` | `#364052` |
| `--border` | Viền, divider | `#343d4d` | `#444c56` | `#3d4759` |
| `--text` | Chữ chính | `#f1f5f9` | `#e6edf3` | `#eceff4` |
| `--text-muted` | Chữ phụ, label | `#9aa6b8` | `#a8b1bc` | `#a3adc2` |
| `--accent` | Nút chính, link, focus | `#60a5fa` | `#539bf5` | `#88c0d0` |
| `--on-accent` | Chữ trên accent | `#0b1b33` | `#0d1a2b` | `#10262c` |

## Màu ngữ nghĩa (dùng chung cho cả 3 theme)

| Token | Màu chữ / số | Nền tag (alpha) | Chữ trên tag |
|---|---|---|---|
| `success` (thu, tăng) | `#4ade80` | `rgba(74,222,128,.16)` | `#86efac` |
| `danger` (chi, giảm) | `#f87171` | `rgba(248,113,113,.16)` | `#fca5a5` |
| `warning` | `#facc15` | `rgba(250,204,21,.16)` | `#fde047` |
| `info` | `#60a5fa` | `rgba(96,165,250,.18)` | `#93c5fd` |

## Màu biểu đồ / danh mục

| Token | Màu | Gợi ý |
|---|---|---|
| `--chart-1` | `#60a5fa` | Xanh dương |
| `--chart-2` | `#fb923c` | Cam |
| `--chart-3` | `#a78bfa` | Tím |
| `--chart-4` | `#2dd4bf` | Teal |
| `--chart-5` | `#f472b6` | Hồng |
| `--chart-6` | `#facc15` | Vàng |

Thanh progress dùng nền track là `--surface-2`, phần đã lấp đầy dùng màu chart.

## CSS

```css
:root,
[data-theme="slate"] {
  --bg: #171b22;
  --surface-1: #222833;
  --surface-2: #2e3644;
  --border: #343d4d;
  --text: #f1f5f9;
  --text-muted: #9aa6b8;
  --accent: #60a5fa;
  --on-accent: #0b1b33;
}

[data-theme="dimmed"] {
  --bg: #22272e;
  --surface-1: #2d333b;
  --surface-2: #373e47;
  --border: #444c56;
  --text: #e6edf3;
  --text-muted: #a8b1bc;
  --accent: #539bf5;
  --on-accent: #0d1a2b;
}

[data-theme="nord"] {
  --bg: #1e2430;
  --surface-1: #2a3240;
  --surface-2: #364052;
  --border: #3d4759;
  --text: #eceff4;
  --text-muted: #a3adc2;
  --accent: #88c0d0;
  --on-accent: #10262c;
}

:root {
  --success: #4ade80;
  --success-bg: rgba(74, 222, 128, 0.16);
  --success-fg: #86efac;
  --danger: #f87171;
  --danger-bg: rgba(248, 113, 113, 0.16);
  --danger-fg: #fca5a5;
  --warning: #facc15;
  --warning-bg: rgba(250, 204, 21, 0.16);
  --warning-fg: #fde047;
  --info: #60a5fa;
  --info-bg: rgba(96, 165, 250, 0.18);
  --info-fg: #93c5fd;

  --chart-1: #60a5fa;
  --chart-2: #fb923c;
  --chart-3: #a78bfa;
  --chart-4: #2dd4bf;
  --chart-5: #f472b6;
  --chart-6: #facc15;
}
```

## Tailwind v4 + shadcn/ui

shadcn đã có sẵn biến `--accent` (dùng cho màu hover nhạt). Vì vậy ở đây token `--accent` của mình được map vào `--primary`, còn `--accent` của shadcn trỏ sang `--surface-2`:

```css
[data-theme], :root {
  --background: var(--bg);
  --foreground: var(--text);
  --card: var(--surface-1);
  --card-foreground: var(--text);
  --popover: var(--surface-2);
  --popover-foreground: var(--text);
  --primary: var(--accent);
  --primary-foreground: var(--on-accent);
  --secondary: var(--surface-2);
  --secondary-foreground: var(--text);
  --muted: var(--surface-2);
  --muted-foreground: var(--text-muted);
  --destructive: var(--danger);
  --input: var(--border);
  --ring: var(--accent);
}

@theme inline {
  --color-surface-1: var(--surface-1);
  --color-surface-2: var(--surface-2);
  --color-success: var(--success);
  --color-success-bg: var(--success-bg);
  --color-danger: var(--danger);
  --color-danger-bg: var(--danger-bg);
  --color-chart-1: var(--chart-1);
  --color-chart-2: var(--chart-2);
  --color-chart-3: var(--chart-3);
  --color-chart-4: var(--chart-4);
  --color-chart-5: var(--chart-5);
  --color-chart-6: var(--chart-6);
}
```

> Biến `--accent` của shadcn phải trỏ về `--surface-2`. Vì khối trên đã dùng tên `--accent` cho màu chính, hãy đổi token của mình thành `--brand` / `--on-brand` trước khi dán vào dự án, sau đó thêm dòng `--accent: var(--surface-2);` vào khối map.

Ví dụ dùng:

```tsx
<div className="rounded-xl border bg-card p-4">
  <p className="text-sm text-muted-foreground">Thu</p>
  <p className="text-2xl font-medium text-success">8.200.000đ</p>
  <span className="rounded-full bg-success-bg px-2 py-0.5 text-xs text-success">
    Lương, freelance
  </span>
</div>
```

## Quy tắc dùng

- Lớp càng cao thì càng sáng: `bg` → `surface-1` → `surface-2`. Không dùng shadow để tạo độ sâu.
- Mỗi màn hình chỉ có một nút dùng `--accent`, các nút còn lại dùng `surface-2` kèm `border`.
- Tag và badge dùng nền alpha 16–18% với chữ cùng màu ở mức sáng hơn, tránh dùng nền màu đặc.
- Kiểm tra tương phản cho cặp `text-muted` trên `surface-2`, vì đây là cặp dễ thiếu nhất. Ở cả 3 theme cặp này vẫn ≥ 4.5:1.
