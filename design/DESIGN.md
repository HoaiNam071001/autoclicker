# Auto Clicker – Design v2

Canvas gốc (xem trực tiếp, chỉnh Tweaks sáng/tối, màu nhấn): https://claude.ai/artifact/5WdDyFZdB3WKxEXwkticnw
Xuất ảnh / PDF: mở canvas → Share › Export.

## Nguyên tắc

- **Nhẹ, ít màu.** Nền và bề mặt là màu đặc cố định, trung tính. Chỉ một màu nhấn (`accent`).
- **Màu nhấn chỉ cho trạng thái:** bước đang chọn, nhãn "Mới", công tắc bật, gạch chân tab, nút chính.
- **Gradient chỉ ở 4 chỗ:** nút chính, viền mục đang chọn, viền ô nhập đang focus, vòng tiến độ. Gradient luôn là một màu nhạt dần (accent → accent trong suốt), không phối nhiều màu.
- **Opacity thay vì màu mới:** chip, nền icon, nút Dừng dùng màu pha độ trong suốt thấp.
- **Số, toạ độ, thời gian** dùng font mono.
- Bố cục và nút cửa sổ kiểu Windows 11 (thu nhỏ / phóng to / đóng bên phải).

## Token màu

> **Theme tối đã đổi sang bảng "Slate"** (`design/newdarkmode.md`, 2026-10-16) – cột *Tối* dưới đây là bản cũ; giá trị đang dùng nằm trong `theme.css` / `frontend/src/style.css`. Thêm token `on-accent` (chữ trên nền accent: trắng ở theme sáng, `#0B1B33` ở theme tối).

Khai báo trong `theme.css`. Class Tailwind: `bg-surface`, `text-muted`, `border-line`, `bg-accent-soft`…

| Token | Sáng | Tối | Dùng cho |
|---|---|---|---|
| `bg` | `#F6F7F9` | `#141518` | Nền vùng nội dung |
| `chrome` | `#FBFBFC` | `#18191C` | Thanh tiêu đề, toolbar, status bar |
| `side` | `#F3F4F7` | `#18191C` | Thanh bên thư viện |
| `surface` | `#FFFFFF` | `#1E1F23` | Thẻ bước, inspector, popup, thanh mini |
| `raised` | `#FFFFFF` | `#2A2C31` | Lựa chọn đang bật trong segmented |
| `field` | `#FFFFFF` | `#17181B` | Ô nhập, select |
| `soft` | ink 4.5% | trắng 5% | Chip thời gian, nền ô icon, nền segmented |
| `hover` | ink 5% | trắng 6% | Mục đang chọn ở thanh bên, hover |
| `group` | ink 1.5% | trắng 2% | Nền khung nhóm bước |
| `line` | `#E6E8EC` | trắng 8% | Viền thường |
| `line-strong` | `#D3D6DC` | trắng 16% | Viền nét đứt (nhóm, nút "+ Thêm bước") |
| `ink` | `#16181D` | `#ECEDF0` | Chữ chính |
| `muted` | `#636A76` | `#A0A5B0` | Chữ phụ, icon, nhãn |
| `accent` | `#4F6BED` | `#4F6BED` | Màu nhấn (đổi được) |
| `accent-soft` | accent 10% | accent 16% | Nền nhãn "Mới", ô icon bước đang chọn |
| `accent-text` | `#3B55CF` | `#B4C1FF` | Chữ trên `accent-soft` |
| `warn` | `#C8893A` | `#D9A35B` | Chấm "Chưa lưu" |
| `ok` | `#3E9B5C` | `#5DBB7A` | Chấm "Sẵn sàng" |
| `danger` / `danger-text` / `danger-soft` | `#D6343A` / `#B4252B` / 8% | — / `#FF9A9D` / 14% | Nút Dừng |
| `hud`, `hud-text`, `hud-line`, `dim` | luôn tối | | Thanh công cụ overlay, lớp phủ mờ |

Màu nhấn thay thế đã thử trên canvas: `#2F8F83` (xanh ngọc), `#8B6CE0` (tím nhạt), `#5B6472` (xám, gần đơn sắc).

## Chữ

- **Sans:** Segoe UI Variable (Windows) → Plus Jakarta Sans (cài `@fontsource/plus-jakarta-sans` để chạy offline, có tiếng Việt).
- **Mono:** Cascadia Mono (Windows) → JetBrains Mono (`@fontsource/jetbrains-mono`).
- Cỡ: tên kịch bản 19/600, tiêu đề vùng 16/600, tên bước 14/600, chữ thường 14/400, phụ 12–13, mono 12–13.

## Bo góc & khoảng cách

- Nút, ô nhập: 6–7 px. Thẻ bước: 9 px. Khung nhóm: 11 px. Popup: 12 px. Thanh mini: 14 px.
- Thẻ bước: padding 9×12, cách nhau 6 px, ô icon 32×32 bo 8.
- Thanh bên rộng 236, inspector 320–360, nội dung giữa co giãn.

## Gradient (4 chỗ duy nhất)

| Chỗ | Utility | Mô tả |
|---|---|---|
| Nút chính (Chạy, Dùng ảnh này, Xong) | `btn-primary` | Dọc: accent 82% → accent, viền accent 90%, viền sáng trên 25%, bóng accent 22% |
| Mục đang chọn | `ring-selected` | Viền 1px chéo 135°: accent → accent 25% |
| Ô nhập đang focus | `ring-focus` | Viền 1px ngang: accent → accent 20% |
| Vòng tiến độ thanh mini | SVG | `<linearGradient>` accent → accent 35%, nét 3.5, đầu tròn |

```html
<svg width="50" height="50" viewBox="0 0 50 50">
  <defs><linearGradient id="ring" x1="0" y1="0" x2="1" y2="1">
    <stop offset="0" stop-color="var(--c-accent)"/><stop offset="1" stop-color="var(--c-accent)" stop-opacity=".35"/>
  </linearGradient></defs>
  <circle cx="25" cy="25" r="22" fill="none" stroke="var(--c-soft)" stroke-width="3.5"/>
  <circle cx="25" cy="25" r="22" fill="none" stroke="url(#ring)" stroke-width="3.5" stroke-linecap="round"
          stroke-dasharray="{progress * 138.2} 138.2" transform="rotate(-90 25 25)"/>
</svg>
```

## Thành phần

- **Thẻ bước:** `bg-surface border border-line rounded-card`. Ô icon `bg-soft text-muted`; khi chọn: `ring-selected`, ô icon `bg-accent-soft text-accent-text`. Phải: chip thời gian `bg-soft text-muted font-mono`. Bước tắt: `opacity-50`.
- **Khung nhóm:** `border border-dashed border-line-strong bg-group rounded-group p-2`; đầu khung: icon, tên nhóm, chip `lặp ×N`.
- **Nhãn "Mới" / trạng thái:** `bg-accent-soft text-accent-text rounded-full text-[11px] font-semibold`.
- **Tab Soạn / Xem trước:** chữ, mục đang chọn có gạch dưới 3px `bg-accent` rộng 50%.
- **Thanh bên:** mục đang chọn `bg-hover` + vạch trái 3×16 `bg-accent`.
- **Segmented (Khi hết giờ):** nền `bg-soft`, lựa chọn bật `bg-raised shadow-[var(--c-raise)]`.
- **Công tắc:** 40×20, bật `bg-accent`, núm trắng 14px.
- **Nút phụ:** `bg-surface border border-line rounded-control`.
- **Overlay:** lớp phủ `bg-dim`, khung chọn nét đứt trắng, 4 núm vuông viền accent, nhãn kích thước nền accent. Thanh công cụ dưới cùng `bg-hud text-hud-text border border-hud-line rounded-xl`, công cụ đang bật nền trắng 14%.
- **Thanh mini:** `bg-surface border border-line rounded-hud`, vòng tiến độ, 2 dòng mono (`whitespace-pre`), nút tạm dừng phụ, nút Dừng `btn-danger-soft`.

## Input của kịch bản

- **Tab Bước / Input (N)** ở đầu inspector; số N là chip `accent-soft`.
- **Nút ghim** cạnh mỗi trường số / chữ: 32×32, viền `line`, icon `muted`; đã ghim: nền `accent-soft`, icon `accent-text`, viền accent 45%. Trường đã ghim: `field-pinned` + chip key (`var-token` có icon ghim) cạnh nhãn.
- **Biến trong ô chữ:** `{{M1}}` hiển thị bằng `var-token`. Gõ `{{` mở popup gợi ý (nền `surface`, viền `line`, bóng nhẹ): key mono `accent-text` + nhãn `muted`, mục đang chọn `bg-hover`, cuối danh sách "+ Tạo input mới…". Dòng phụ dưới ô: "N biến, tự thêm vào tab Input".
- **Thẻ bước** có input: chip viền `line` "1 input" / "2 input".
- **Tab Input:** mỗi input một dòng (⋮⋮ kéo, nhãn + `*` nếu bắt buộc, `key · mặc định` mono, chip kiểu `bg-soft`, chip bước dùng viền `line`). Input đang sửa mở rộng với `ring-selected`: nhãn, kiểu, mặc định, bắt buộc (công tắc), gợi ý. Cuối: "+ Thêm input", ghi chú giá trị lần trước lưu trên máy.
- **Form Chạy kịch bản** (popup 560, bo 12):
  - Đầu: icon Chạy `accent-soft`, tên kịch bản + số input, nút đóng.
  - Hàng "Bộ giá trị": select preset + "Lưu bộ mới".
  - Thân: nhóm theo bước (nhãn nhỏ `muted` có số bước mono), ô cao 34. Ô đã khác mặc định dùng `ring-focus` + dòng "đã đổi, mặc định …" màu `accent-text` + "khôi phục". Bắt buộc: `*` màu `warn`. Key mono nhỏ cạnh nhãn.
  - Khung "Sẽ gõ" `bg-soft` xem trước câu sau khi thay biến, phần thay thế đậm `accent-text`.
  - Chân `chrome`: ô tick "Lần sau chạy luôn với giá trị này", "Đếm ngược 3 giây", nút Huỷ (phụ) và Chạy (`btn-primary`, F8).

## Icon bước (SVG 24×24, nét 1.8, đầu tròn)

| Loại | path `d` |
|---|---|
| click | `M5 3l13 7-6 1.5L9 18z` |
| double | `M5 3l13 7-6 1.5L9 18zM15 3l2-2M19 6l2-1` |
| hold | `M5 3l13 7-6 1.5L9 18zM17 15a3 3 0 1 0 0 .01` |
| drag | `M5 19a2 2 0 1 0 0-.01M7 17L18 6M13 6h5v5` |
| scroll | `M12 4v16M8 8l4-4 4 4M8 16l4 4 4-4` |
| type | `M4 7V5h12v2M10 5v14M7 19h6M15 13h5M17.5 13v6` |
| key | `M5 9h14v10H5zM8 13h1M11 13h1M14 13h2M8 16h8` |
| wait | `M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zM12 7v5l3 2` |
| waitImage | `M4 8V5h3M17 5h3v3M20 16v3h-3M7 19H4v-3M9 12a3 3 0 1 0 6 0 3 3 0 0 0-6 0` |
| clickImage | `M4 8V5h3M17 5h3v3M7 19H4v-3M10 9l9 5-4 1-2 4z` |
| waitPixel | `M4 4h7v7H4zM13 13h7v7h-7zM13 4h7v7h-7z` |
| open | `M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18zM3 12h18M12 3c2.5 2.7 2.5 15.3 0 18M12 3c-2.5 2.7-2.5 15.3 0 18` |
| focusWindow | `M4 5h16v14H4zM4 9h16` |
| group | `M3 3h7v7H3zM14 3h7v7h-7zM3 14h7v7H3zM14 14h7v7h-7z` |
| ghim input | `M9 4h6l-1 5 3 3v2H7v-2l3-3zM12 14v6` |

## Màn hình trên canvas → file trong repo

| Artboard | Nội dung | File trong repo |
|---|---|---|
| `Main` / `MainDark` | Soạn kịch bản: thư viện, danh sách bước + nhóm, inspector "Chờ ảnh" | `screens/MainScreen.tsx`, `components/StepList.tsx`, `StepForm.tsx`, `RunSettings.tsx`; thêm `components/Library.tsx` |
| `AddStep` / `AddStepDark` | Bảng thêm bước (Ctrl+K) | Thay menu "+ Thêm bước" trong `StepList.tsx` bằng `components/AddStepPalette.tsx` |
| `Overlay` | Chế độ chụp ảnh mẫu | `screens/OverlayScreen.tsx`, thêm tool `capture` |
| `MiniBar` / `MiniBarDark` | Đang chạy, đang chờ ảnh | `screens/MiniBar.tsx` |
| `StepInputs` | Bước gõ chữ: biến `{{…}}`, gợi ý, nút ghim | `components/StepForm.tsx`; thêm `components/VarTextInput.tsx`, `components/PinButton.tsx` |
| `InputsTab` | Tab Input của kịch bản | thêm `components/InputsPanel.tsx` |
| `RunDialog` / `RunDialogDark` | Form Chạy kịch bản | thêm `components/RunDialog.tsx`; nút Chạy / F8 mở form khi kịch bản có input |

## Đổi tên token từ bản cũ

| Cũ | Mới |
|---|---|
| `ground` | `bg` |
| `ink`, `muted`, `line`, `accent`, `accent-soft`, `danger`, `ok` | giữ tên (giá trị mới) |
| `accent-strong` | `accent-text` |
| `accent-light`, `drag`, `drag-light`, `drag-marker` | bỏ; kéo thả dùng chung icon/màu trung tính, marker overlay dùng `accent` |
| `dark-*` | bỏ; dùng token chung + `data-theme="dark"` (MiniBar, PreviewScreen) |

`design/canvas/` chứa mã nguồn các artboard (`.dc.html`) để tham khảo; chúng cần runtime của canvas, không chạy riêng được.
