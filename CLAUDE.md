# Auto Clicker

Desktop app tự động click theo kịch bản (chọn vùng, chọn điểm, kéo thả A→B, lặp N vòng). Wails v2 + Go + React/TS.

## Bắt đầu mỗi session

1. Đọc **`docs/tasks/README.md`**: *Trạng thái hiện tại* (đang làm bước nào), bảng thứ tự các bước + link file từng bước. Chỉ mở file của bước đang làm (`docs/tasks/phase3/<số>-<tên>.md`); `decisions.md` là quyết định đã chốt (không bàn lại), `checks.md` là việc chờ người dùng kiểm tra tay, `log.md` là nhật ký, `overview.md` là stack / data model. Thiết kế UI v2 (token, thành phần, icon, artboard) ở **`design/DESIGN.md`** + `design/theme.css`.
2. Xong một bước (hoặc chuyển sang bước mới) → cập nhật **trước khi commit / kết thúc session**: tick `[x]` trong file của bước, sửa *Trạng thái hiện tại* ở `README.md`, thêm 1 dòng vào `log.md` (ngày tuyệt đối, kết quả kiểm tra). Quy ước:
   - Mục con làm xong là tick ngay; việc mới phát sinh (sửa theo người dùng báo, việc để sau) thêm thành mục con `[ ]` / `[x]` trong file của bước, không chỉ ghi ở chat.
   - Mục cha tick khi hết phần code. Việc *Kiểm tra tay* của phần đã xong code ghi vào **`checks.md`** (mới nhất ở trên, nhóm theo bước), file bước chỉ để dòng link; người dùng thử xong thì tick ở `checks.md` + ghi `log.md`.
   - Bước mới → file mới `docs/tasks/phase3/<số>-<tên>.md` + một dòng trong bảng ở `README.md`. File quá dài (~200 dòng) thì tách theo phần con.
   - Trước khi kết thúc, rà lại các mục cha có mọi mục con code đã `[x]` mà chưa tick.
   - Tính năng đổi trạng thái (xong code, kiểm tra tay xong, thêm tính năng mới) → sửa bảng *Tính năng của hệ thống* ở đầu `README.md` (viết tay, theo nhóm).
   - Tick / thêm bước / sửa `checks.md` xong → chạy `python3 scripts/tasks-index.py` để cập nhật bảng *Tổng quan* trong `README.md` (bảng tự sinh, không sửa tay).
3. Người dùng trao đổi bằng **tiếng Việt** – trả lời, comment UI và tài liệu bằng tiếng Việt; tên biến/hàm bằng tiếng Anh.

## Lệnh

```bash
wails dev                      # chạy app (hot reload). wails nằm ở ~/go/bin
wails build                    # build ra build/bin/
cd frontend && npm run build   # type-check (tsc) + build frontend – chạy sau mỗi thay đổi TS
cd frontend && npm test        # test đơn vị frontend (vitest, *.test.ts)
```

Xem riêng từng màn trong `wails dev` / trình duyệt: thêm hash `#preview`, `#overlay`, `#minibar` (mặc định là MainScreen). Trên trình duyệt không có Go nên nút Chạy bị tắt.
Màn mẫu giao diện đầy đủ (kể cả tính năng chưa có logic, dữ liệu giả trong `frontend/src/demo/`): `#demo`, `#demo-overlay`, `#demo-minibar`, `#demo-flow` (canvas flow bước 20; `?mode=running` / `?mode=failed`).

## Cấu trúc

- `main.go`, `app.go` – Wails app; method public của `App` được bind sang frontend (`frontend/wailsjs/go/main/App`, file tự sinh – không sửa tay).
- `model.go` – model kịch bản (`Step` = node, `Script` có `start` + `edges`, `Retry`); `flow.go` – flow (bước 20): `StepResult`, cổng ra `portsOf`, `flowOrder` (thứ tự hiển thị – `frontend/src/flow/model.ts` cùng công thức, test chung `testdata/flow_order.json`); test viết kịch bản dạng danh sách thì bọc `linear(sc)` (`flowtest_test.go`); `flowcheck.go` – `validateFlow`; `flowrun.go` – engine chạy theo cạnh (`runLoop`, `execNode`, thử lại, event `engine:node`); `vars.go` – biến `{{n2.output.x}}` lúc chạy, `varFields` (mọi ô chữ nhận biến – input, biến node, kiểm tra dùng chung); `nodes.go` – node điều khiển (if / switch / switchScreen / counter / end / setVar: `runControl`, cổng, kiểm tra); `cond.go` – điều kiện (`evalCond`, `snapshot` chụp một lần mỗi lần xét); `nodetest.go` – nút Thử node (`TestNode`); `window.go` – chế độ cửa sổ main / overlay / minibar; `hotkey.go` – phím toàn cục (gohook); `engine.go` – chạy kịch bản qua interface `Input` (`input.go` = robotgo, test dùng input giả – `go test ./...`); `runner.go` – `RunScript` / `StopScript` / `PauseScript` / `ResumeScript`; `failsafe.go` – dừng khẩn cấp khi chuột chạm góc màn hình. `open.go` – bước Mở URL / app (`openCommand` theo hệ điều hành). `match.go` – tìm ảnh mẫu (NCC, Go thuần); `capture.go` / `capture_darwin.go` – chụp màn hình qua interface `Screen`, quyền Screen Recording macOS; `imagestep.go` – engine chờ ảnh / click theo ảnh / chờ pixel; `imageapi.go` – `CaptureTemplate`, `FindImage`, `PixelColorAt` cho giao diện (ẩn cửa sổ lúc chụp). `ocr.go` / `ocr_darwin.go` – tìm theo chữ (đọc chữ bằng Vision, so chữ `findText`; Windows chưa có). `steptypes.go` – mô tả + ví dụ JSON từng loại bước (cho AI / MCP; thêm loại bước mới phải thêm dòng, có test). `inputs.go` – input của kịch bản (`ScriptInput`, `resolveInputs` áp giá trị vào bản sao trước khi chạy, `ScriptError` lỗi có cấu trúc, giá trị lần trước / preset trong `inputs.json`). `anchor.go` – neo điểm theo cạnh vùng (`resolvePoint`; `frontend/src/anchor.ts` phải cùng công thức). `frontend/src/transform.ts` – ma trận biến đổi nhiều điểm (bước 30, có test); `components/TransformTool.tsx` – công cụ Chọn nhiều trên overlay (khung biến đổi, thanh số, hộp ánh xạ vùng).
- `frontend/src/store/` – store kịch bản đang soạn (dùng chung các màn), import qua `store/index.ts`: `editor` (state, hoàn tác / làm lại, thay kịch bản), `geometry` (khung, neo ↔ toạ độ màn hình, `endpoints` – chỉ đọc), `nodes` (thêm / sửa / xoá / nhân bản node, nối cạnh, node bắt đầu, vị trí, chỗ nối node mới `insertAt`), `control` (thêm / sửa node điều khiển, sửa điều kiện lồng nhau), `points` (đặt / chọn lại điểm trên overlay, F6, kéo thả), `multi` (chọn nhiều + biến đổi, đổi vùng). `frontend/src/flow/` – logic flow thuần, có test: `model` (cổng ra, `flowOrder`, `orderedSteps`), `edit` (nối / cắt cạnh, chèn / xoá node, `nextRef`, `pruneEdges`), `layout` (elkjs nạp lười, `placeNode`), `view` (node → dữ liệu thẻ, nhãn cổng, `outputsOf`), `cond` (điều kiện lồng nhau theo đường dẫn `CondAt`, case mới). `scriptInputs.ts` – ghim trường / sửa input / bước dùng input (bước 21b, có test); `runForm.ts` – bấm Chạy / F8 → form Chạy khi có input, preset, chạy luôn; `imageSteps.ts` – khoanh ảnh mẫu (overlay tool `capture`), thử tìm, lấy màu, quyền chụp màn hình; `components/MockPage.tsx` – trang mẫu phủ nền khi Xem trước cuộn; `coordMap.ts` – đổi toạ độ overlay ↔ robotgo (tỷ lệ + độ lệch đo từ chuột); `engineStore.ts` – store trạng thái engine (event `engine:state` / `engine:progress`); `windowMode.ts` – `setWindowMode` đổi cửa sổ + màn; `theme.ts` – theme sáng / tối / theo hệ thống.
- `settings.go` – cấu hình app (theme, mở gần đây) trong `os.UserConfigDir()/autoclicker/`; `storage.go` – lưu / mở file kịch bản `.json` (`ScriptVersion`, `checkScript`), bản nháp `draft.json`, hỏi lưu khi đóng app. `frontend/src/scriptFile.ts` – Lưu / Mở / Mới / mở gần đây, tự ghi nháp, cờ chưa lưu; `components/SaveChangesDialog.tsx` – hỏi Lưu / Không lưu / Huỷ.
- `frontend/src/screens/` – các màn UI; `components/TitleBar.tsx` – thanh tiêu đề kiểu Windows 11 + chọn theme; `components/Library.tsx` – thanh bên thư viện; `components/flow/FlowEditor.tsx` – canvas flow của màn chính (nối store); `components/AddStepPalette.tsx` – bảng thêm node Ctrl/Cmd+K (`PaletteView` dùng chung); `components/StepForm.tsx` – khung inspector, trường theo loại ở `components/stepform/` (`shared`, `MouseFields`, `ScrollFields`, `KeyFields`, `TextFields`, `ImageFields`, `FlowSections`, `InputsTab`); `components/RunSettings.tsx` – lặp / đếm ngược (popup nút "N vòng"); `components/Simulation.tsx` – xem trước bằng con trỏ giả; `components/icons.tsx` – icon SVG; `types.ts` – model `Step` (phải khớp struct Go trong `model.go`; tsc tự kiểm tra).
- `frontend/src/style.css` – Tailwind v4 + token thiết kế trong `@theme`.
- `components/inputs/` – ghim (`PinMenu`), ô chữ `{{biến}}`, tab Input, form Chạy (đã nối logic ở 21b). Giao diện dựng sẵn cho tính năng chưa có logic (bước 27b, app thật hiện tắt + *Sắp có*): `components/steps/` (inspector loại bước mới, hết giờ, nhóm, neo ±), `ScheduleDialog.tsx`, `AppChrome.tsx` (`Toolbar`, `StatusBar`, nút debug / ghi). Kiểu dữ liệu UI tạm trong `steps/model.ts` – chuyển vào `types.ts` + `model.go` khi làm logic bước đó.
- `components/flow/` – phần hiển thị flow (chỉ nhận props, app thật và màn mẫu dùng chung): `FlowCanvas` (`@xyflow/react`), `NodeCard`, `FlowEdge`, `PortField` (Cổng ra / Khi thất bại / Output / Lần chạy gần nhất), `RunLog`, `nodes/ControlFields` + `nodes/CondEditor` + `nodes/NodeTestPanel` (inspector node điều khiển; app nối qua `stepform/ControlSection`), kiểu hiển thị `model.ts` (cỡ thẻ cố định). Màn mẫu `screens/DemoFlow.tsx` + `demo/flowDemo.ts` (`?sel=<id>` mở sẵn inspector node).
- `autoclicker-ui/` – bản thiết kế gốc (Tailwind v3, chỉ để tham khảo, không build) – bị thay bởi design v2 từ bước 27.
- `design/` – design v2: `DESIGN.md` (token, thành phần, icon, map artboard → file), `theme.css` (token Tailwind v4 sáng / tối), `canvas/` (mã artboard, chỉ tham khảo – cần runtime canvas, không chạy riêng).

## Quy tắc code

Mục tiêu: thêm tính năng sau chỉ phải đụng đúng chỗ. Code mới theo các quy tắc dưới; gặp code cũ lệch thì sửa khi đang làm phần đó.

- **Một file một việc.** File quá ~400 dòng hoặc gom nhiều nhóm việc → tách theo nhóm (vd `components/stepform/`, `store/`). Tên file nói đúng việc nó làm.
- **Ba lớp ở frontend**, phụ thuộc một chiều logic thuần ← store ← component:
  1. *Logic thuần* (`src/flow/`, `anchor.ts`, `transform.ts`, phần tính trong `scriptInputs.ts`): không React, không Wails, không đọc store; có test `*.test.ts` cạnh file.
  2. *Store* (`src/store/`): ghép logic thuần với state, hoàn tác, thông báo. Chỉ module trong `store/` gọi `setEditor`; mọi sửa kịch bản đi qua `updateScript` (để có hoàn tác + cờ chưa lưu).
  3. *Component*: hiển thị + gọi hàm nghiệp vụ của store (`connectPort`, `updateStep`…), không tự dựng `script` mới. Tách thành phần **chỉ hiển thị** (nhận props: `FlowCanvas`, `NodeCard`, `PortField`) khỏi phần **nối store** (`FlowEditor`, `stepform/FlowSections`) để màn mẫu `#demo*` và app thật dùng chung một giao diện.
- **Flow:** thao tác node / cạnh qua `flow/edit.ts`; thứ tự node qua `flow/model.ts` (`flowOrder`, `orderedSteps`) – không dựa vào thứ tự mảng `script.steps`. Tham chiếu node bằng `id` (cạnh, `start`, `imageFrom`); `ref` (n2) chỉ để hiển thị và dùng trong biến. Tên / dòng chi tiết / nhãn cổng / output của node chỉ lấy từ `flow/view.ts`.
- **Go là nguồn chuẩn của model:** thêm trường → `model.go` → `~/go/bin/wails generate module` → `types.ts` (tsc tự báo lệch). Công thức chạy ở cả hai bên (neo, `flowOrder`, cổng ra) giữ cùng tên hàm và test chung dữ liệu trong `testdata/`.
- **Thêm loại node** – đủ các chỗ: hằng + mặc định (`model.go` `NewStep` + `stepTypes`, `store/nodes.ts` `addStep` / `store/control.ts` `addControl`), `checkStep`, `validateScript`, `runStep` / `runControl` (output), `portsOf` (Go + TS), ô chữ nhận biến (`varFields` Go + `mapVarText` TS), `steptypes.go` (mô tả, output, ví dụ), `flow/view.ts` (nhãn cổng, chi tiết, `outputsOf`), trường trong `components/stepform/` hoặc `components/flow/nodes/`, mục trong bảng thêm node, test Go + TS.
- **Dùng lại trước khi viết mới:** `NumInput`, `Segmented`, `Switch`, `Field` / `NumField` / `PickButton` (`stepform/shared`), `PaletteView`, `PortField`, `StepIcon`. Hai chỗ cần cùng một thứ → đưa vào chỗ dùng chung, không chép.
- **Không giữ code chết / tương thích ngược** khi không cần (file kịch bản cũ đã bỏ). Cầu nối tạm thời phải ghi bước sẽ xoá trong file của bước (`docs/tasks/`). Helper chỉ cho test để trong `*_test.go` / `*.test.ts`.
- **Go:** một file một nhóm (`flow.go` kiểu + thứ tự, `flowcheck.go` kiểm tra, `flowrun.go` engine, `vars.go` biến…). Thiết bị (chuột, màn hình, OCR, mở app) qua interface để test bằng bản giả. Lỗi cho người dùng / AI có cấu trúc (`ScriptError`, `StepError` + mã).
- Dữ liệu giả chỉ ở `src/demo/`. Comment tiếng Việt, ngắn, nói *vì sao* / ràng buộc; không comment lại điều code đã nói.

## Quy ước UI

Theo **design v2** – chi tiết token, thành phần, icon trong `design/DESIGN.md`.

- Tailwind **v4**: token khai báo trong `style.css` (chép từ `design/theme.css` – sửa token thì sửa cả hai), không tạo `tailwind.config.js`. Body không tô nền (overlay / mini bar cần trong suốt) – mỗi màn tự tô `bg-bg` / `bg-surface`.
- Token: `bg` nền, `chrome` thanh tiêu đề / status bar, `side` thanh bên, `surface` thẻ / inspector / popup, `field` ô nhập, `soft` / `hover` nền pha nhẹ, `line` / `line-strong` viền, `ink` / `muted` chữ, `accent` / `accent-soft` / `accent-text` (chỉ cho trạng thái), `warn`, `ok`, `danger*`; `hud*` / `dim` cho overlay (luôn tối).
- Theme tối dùng bảng **Slate** trong `design/newdarkmode.md` (lớp càng cao càng sáng: `bg` → `surface` → `field` / `soft` → `raised`). Chữ trên nền `accent` dùng `text-on-accent`, không dùng `text-white`.
- Theme: `data-theme="light|dark"` trên `<html>` (`theme.ts`, lưu qua `settings.go`), không đặt = theo hệ điều hành; overlay đặt `data-theme="dark"` riêng. Xem thử trên trình duyệt: `?theme=dark`.
- Gradient chỉ ở 4 chỗ: `btn-primary`, `ring-selected`, `ring-focus`, vòng tiến độ thanh mini. Nút Dừng `btn-danger-soft`.
- Thành phần dùng chung: `NumInput` (+ `fieldBox`), `Segmented`, `Switch` (trong `RunSettings.tsx`), `StepIcon` (`icons.tsx`, đủ icon cả loại bước chưa làm).
- Font `font-sans` (Segoe UI Variable trên Windows → Plus Jakarta Sans) cho chữ, `font-mono` (Cascadia Mono → JetBrains Mono) cho số/toạ độ/thời gian. Font cài local qua `@fontsource`, không dùng Google Fonts.

## Lưu ý nền tảng

- Đang dev trên macOS; mục tiêu phát hành đầu tiên là Windows (robotgo cần cgo + MinGW-w64).
- macOS: cấp quyền Accessibility (System Settings → Privacy & Security) cho terminal / app thì robotgo mới click được.
- Wails v2 chỉ có một cửa sổ: overlay và mini bar là chế độ của cửa sổ chính.
