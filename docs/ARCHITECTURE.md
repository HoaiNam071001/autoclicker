# Kiến trúc & luồng chính

> Đọc file này để nắm hệ thống mà không phải mở code. Tìm file theo việc: §6; danh sách file đầy đủ: §8.
> Tính năng đang có / chưa có: bảng *Tính năng của hệ thống* trong `docs/tasks/README.md` (không chép sang đây).
> Đổi luồng / thêm API / thêm event → sửa file này cùng commit.

## 1. Bức tranh tổng

```
┌──────────── Frontend (React/TS, WebView) ────────────┐        ┌──────── Go: main (root – chỉ Wails) ─────────┐
│ screens/  MainScreen · OverlayScreen · MiniBar ·     │ gọi    │ App (app.go) – method public = API bind,     │
│           PreviewScreen · Demo* (hash URL chọn màn)  │──────► │   mỗi method vài dòng gọi xuống internal/    │
│ components/ (chỉ hiển thị + gọi hàm store)           │ Wails  │ runner.go · storage.go · imageapi.go ·       │
│ store/ (kịch bản đang soạn, hoàn tác)                │ bind   │ window.go · hotkey.go (gohook F6/F8/F10)     │
│ engineStore · runStore · runForm · scriptFile ·      │◄────── ├──────────────── internal/ ───────────────────┤
│ imageSteps · coordMap · windowMode                   │ events │ engine  chạy flow (node, biến, thử lại)      │
│ flow/ anchor.ts transform.ts (logic thuần, có test)  │        │ script  đọc / kiểm tra kịch bản, input, neo  │
└──────────────────────────────────────────────────────┘        │ store   cấu hình, file, nháp, lịch sử chạy   │
                                                                │ device  chuột / phím, chụp, OCR (interface)  │
                                                                │ vision  tìm ảnh / chữ   launch  mở URL / app │
                                                                │ model   kiểu dữ liệu (mọi package dùng)      │
                                                                └──────────────────────────────────────────────┘
```

- **Một cửa sổ Wails** (frameless, nền trong suốt). Overlay / mini bar = đổi cỡ + vị trí + always-on-top cửa sổ đó (`SetWindowMode`), frontend đổi `location.hash` → `App.tsx` chọn màn. Các màn **không cùng mount** → state dùng chung nằm trong store module (`useSyncExternalStore`), không trong React state.
- **Go là nguồn chuẩn model**; frontend gọi Go qua `frontend/wailsjs/go/main/App` (tự sinh), Go đẩy trạng thái về qua `runtime.EventsEmit`.
- **Trình duyệt không có Go** (`inWails = false`): store nạp dữ liệu mẫu, nút Chạy tắt – dùng để xem UI.

## 2. Model (`internal/model`: model.go, flow.go)

- `Script` = `{version(4), name, region(nil = màn chính), loops(0 = vô hạn), maxMinutes, steps[], start, edges[], anchorMode, inputs[]}`.
- `Step` = **node**. Định danh bằng `id` (cạnh, `start`, `imageFrom` trỏ theo id); `ref` (`n2`) chỉ để hiển thị + biến. Thứ tự mảng `steps` **không** quyết định thứ tự chạy.
- `Edge{from, port, to}` – mỗi cổng tối đa một cạnh. Cổng theo loại (`portsOf`, Go + TS cùng công thức):

| Loại | Cổng ra |
|---|---|
| hành động (click, double, hold, drag, wait, type, key, scroll, open, waitImage, clickImage, waitPixel) | `success`, `fail` |
| `if` | `true`, `false` |
| `switch` / `switchScreen` | `case:<key>`… + `default` / `none` |
| `counter` | `next`, `done` |
| `setVar` | `success` |
| `end` | (không có) |

- `Point{x, y, ax, ay}` – toạ độ **trong vùng**, mỗi trục neo `start|end|center|ratio` (`script/anchor.go` ↔ `anchor.ts`). Đổi ra toạ độ màn hình lúc chạy bằng `script.ResolvePoint`.
- `StepResult{nodeId, ref, status, port, output, error{code,message}, logs, attempts, loop, shot}` – kết quả mỗi lần chạy node. Mã lỗi: `not_found`, `timeout`, `permission`, `invalid_input`, `cancelled`, `end`, `internal`.
- Output theo loại (đầy đủ ở `internal/model/steptypes.go`): chờ / click ảnh·chữ → `x, y, w, h, score(, text)`; waitPixel → `color`; if → `result` + chỗ thấy; end thành công → `message`.

## 3. Luồng soạn kịch bản

```
component ──gọi──► store/nodes|control|points|multi ──► flow/edit.ts (logic thuần) ──► updateScript()
                                                                                      │ push past[], bật dirty
                                                                                      ▼
                                       listeners → mọi màn re-render;  scriptFile.ts ghi nháp (giãn 800 ms) → SaveDraft
```

- Mọi sửa kịch bản đi qua `updateScript` (`store/editor.ts`) → có hoàn tác / làm lại + cờ chưa lưu (`SetDirty` sang Go để hỏi khi đóng).
- Thêm node: bảng Ctrl/Cmd+K (`AddStepPalette`) → `addStep` (gọi Go `NewStep` lấy mặc định) / `addControl` → nối vào node đang chọn (`insertAt`).
- Canvas: `FlowEditor` (nối store) → `FlowCanvas` / `NodeCard` / `FlowEdge` (chỉ props). Tên, chi tiết, nhãn cổng, output gợi ý lấy từ `flow/view.ts`. Tự xếp: `flow/layout.ts` (elkjs nạp lười).
- Inspector: `StepForm` → trường theo loại ở `components/stepform/*`, node điều khiển ở `components/flow/nodes/*` (qua `stepform/ControlSection`).

### Đặt điểm / vùng trên overlay

1. `openOverlay(tool, focus)` (`windowMode.ts`) → Go `OverlayDisplayFor` chọn màn chứa vùng → `SetWindowMode('overlay')` → hash `#overlay`.
2. `OverlayScreen` vẽ vùng + marker; click → toạ độ cửa sổ → **`coordMap.ts`** đổi sang toạ độ robotgo (ước lượng `k, b` từ mẫu `MousePos`, xử lý DPI / màn phụ) → `store/points` đổi sang điểm neo theo vùng (`store/geometry`).
3. F6 (Go `hotkey.go`) → event `hotkey:f6 {x,y}` (toạ độ robotgo sẵn) → `store/points` thêm điểm.
4. Công cụ: `point`, `drag` (A rồi B), `region`, `preview` (con trỏ giả – `Simulation.tsx`), `select` (biến đổi nhiều điểm – `TransformTool` + `transform.ts`), `capture` (khoanh ảnh mẫu – `imageSteps.ts` → Go `CaptureTemplate`, ẩn cửa sổ lúc chụp).

### Lưu / mở (storage.go → `store` / `script` ↔ scriptFile.ts)

- File `.json` (`ScriptVersion` = 4; cũ hơn từ chối, không chuyển đổi). Mở: `store.ReadScriptFile` → `script.DecodeScript` → `checkScript` (chỉ sửa cấu trúc: ref trùng, cạnh hỏng…) – **không** `validateScript` (bản soạn dở vẫn mở được).
- Nháp `draft.json` (kịch bản + path + dirty) – mở app nạp lại. Đóng app khi dirty: `beforeClose` chặn → event `app:confirmClose` → `SaveChangesDialog` → `ConfirmQuit`.
- Thư mục cấu hình `os.UserConfigDir()/autoclicker/`: `settings.json` (theme, mở gần đây), `draft.json`, `inputs.json` (giá trị input lần trước / preset theo kịch bản), `runs/<tên>-<hash path>/` (lịch sử chạy), `sandbox/` (trang sân thử của Mẫu thử, ghi đè mỗi lần mở mẫu).
- Mẫu thử (bước 33): `OpenSample` → `samples.Open` (qua `DecodeScript` như file) → `scriptFile.openSample` mở như kịch bản mới chưa có file (hỏi lưu nếu đang dirty).

## 4. Luồng chạy

```
Bấm Chạy / F8 (hotkey:f8start)
  └► runForm.startRun: có input & chưa "chạy luôn" → RunDialog (preset, giá trị lần trước) ; không thì chạy ngay
      └► engineStore.runScript
           1. Go ValidateScript(sc, values)   — lỗi: ở lại màn chính, hiện lastError
           2. setWindowMode('minibar')
           3. Go RunScript(sc, values, RunOptions{countdown, path, source, from, step})   (chặn tới khi xong)
           4. setWindowMode('main')
```

**Bên Go – `RunScript` (runner.go):**
1. `script.ResolveInputs` – thay `{{key}}` input vào **bản sao** kịch bản; `script.ValidateScript` (flow: `validateFlow`, rồi từng node).
2. `engine.New(device.TrackingInput{RobotgoInput}, sc, engine.Options{…})`; đếm ngược (`engine:countdown`); bật `device.WatchCorners` (chuột chạm góc → dừng, `engine:stopped "corner"`).
3. `Engine.Run` (`internal/engine/flowrun.go`) → mỗi vòng `runLoop`: từ `start`, lặp:
   - `breakBefore` (debug: breakpoint / chạy từng node → chờ Tiếp / F10)
   - `execNode`: node tắt → `skipped` đi `success`; `resolveNodeVars` (biến `{{n2.output.x}}`, `{{last.…}}`, `{{vars.x}}` – thiếu biến → `invalid_input`); node điều khiển → `runControl` (`nodes.go`, `cond.go` xét điều kiện, chụp một lần mỗi lần xét); node hành động → `runStep` (`actions.go` `doStep`, ảnh / pixel / chữ ở `imagestep.go` + `internal/vision`); lỗi thì thử lại theo `Retry` (node hành động).
   - Thất bại → chụp màn (`shoot`, `shots.go`). Gửi `engine:node` (StepResult). Đi theo `edgeOut(port)`; `fail` không nối → dừng báo lỗi (`engine:failed`); cổng khác không nối → xong vòng. Thành công → chờ `delayAfter ± jitter`.
   - Quá 10 000 node / vòng → lỗi lặp vô hạn. Kết quả + biến đặt lại mỗi vòng.
4. `finishRun` → `store.SaveRun`: `runs/…/<id>.json` + ảnh lỗi (giữ 50 lần) → `engine:done` (RunSummary). Cuối cùng `engine:state false`.

Dừng / tạm dừng: `StopScript` (huỷ ctx), `PauseScript` / `ResumeScript` / `StepScript` (F10) qua `pauseGate` (`engine/gate.go`). F8 lúc đang chạy = dừng (`engine:stopped "f8"`).

**Frontend nhận:** `engineStore` (state, countdown, progress, failed, stopped, hotkey:error) → MiniBar / thanh trạng thái; `runStore` (engine:node, engine:done) → canvas tô trạng thái node + cạnh đã đi (`flow/run.ts`) + `RunLog`; `imageSteps` (engine:failed kèm chỗ giống nhất) để xem lại ảnh mẫu.

### Các đường chạy phụ

- **Thử node** (`TestNode` → `Engine.TestNode`, `engine/nodetest.go`): xét if / switch / switchScreen trên màn hình lúc này, không chạy cả flow → `NodeTestPanel` + khung chỗ thấy trên overlay.
- **Thử tìm ảnh / chữ** (`FindImage` → `Engine.FindOnce`), **lấy màu** (`PixelColorAt`), **quyền chụp** (`ScreenCaptureAllowed` / `RequestScreenCapture`, macOS Screen Recording).
- **Chạy từ node**: `RunOptions.from` → một vòng bắt đầu ở node đó (biến node trước trống).
- **Lịch sử**: `ListRuns` / `GetRun` / `RunShot` – mở từ thư viện, hiện như lần chạy đang xem.

## 5. API Go (method của App) & event

| Nhóm | Method |
|---|---|
| Model | `NewStep(type)`, `NewScript()`, `StepTypes()`, `FlowGuide()` (mô tả cho AI) |
| File | `SaveScript`, `SaveScriptAs`, `OpenScript`, `OpenScriptPath`, `RecentScripts`, `SaveDraft`, `LoadDraft`, `ClearDraft`, `SetDirty`, `ConfirmQuit` |
| Chạy | `ValidateScript`, `RunScript`, `StopScript`, `PauseScript`, `ResumeScript`, `StepScript`, `TestNode` |
| Lịch sử | `ListRuns`, `GetRun`, `RunShot` |
| Input | `GetInputMemory`, `SaveInputMemory` |
| Ảnh / màn hình | `CaptureTemplate`, `FindImage`, `PixelColorAt`, `ScreenCaptureAllowed`, `RequestScreenCapture`, `ScreenBounds`, `MousePos` |
| Cửa sổ | `SetWindowMode`, `Displays`, `OverlayDisplayFor`, `SetOverlayDisplay` |
| Cấu hình | `GetSettings`, `SetTheme` |
| Mẫu thử | `Samples`, `OpenSample(id)` (ghi sân thử ra `ConfigDir/sandbox`, trả kịch bản), `OpenSandboxPage(page)` |

| Event (Go → frontend) | Payload | Nghe ở |
|---|---|---|
| `engine:state` | bool | engineStore, runStore |
| `engine:countdown` | số giây | engineStore |
| `engine:progress` | `Progress` (vòng, node, thử lại, tạm dừng, thời gian) | engineStore |
| `engine:node` | `StepResult` | runStore |
| `engine:done` | `RunSummary` | runStore |
| `engine:failed` | `{stepId, miss?}` | engineStore, imageSteps |
| `engine:stopped` | `"f8"` \| `"corner"` | engineStore |
| `hotkey:f8start` / `hotkey:f6` / `hotkey:error` | – / `Point` / chuỗi | runForm / store/points / engineStore |
| `app:confirmClose` | – | scriptFile |

## 6. Muốn làm X → sửa ở đâu

| Việc | Chỗ |
|---|---|
| Thêm loại node | Skill `add-node-type` (Go model + check + run + PortsOf + VarFields + steptypes; TS view + form + palette + test) |
| Thêm trường cho node | `internal/model/model.go` → `wails generate module` → `types.ts` → `script` (`checkStep` / `ValidateScript`) → form trong `stepform/` |
| Đổi cách chạy / rẽ nhánh | `internal/engine`: `flowrun.go` (vòng, thử lại, cạnh), `nodes.go` (node điều khiển), `cond.go` (điều kiện) |
| Ô chữ mới nhận `{{biến}}` | `script.VarFields` (`varfields.go`) + `mapVarText` (`scriptInputs.ts`) |
| Tìm ảnh / chữ / màu | `engine/imagestep.go`, `internal/vision` (`match.go` NCC, `ocr.go` tìm chữ), `internal/device` (`ocr_*.go` đọc chữ, `capture*.go` chụp); UI `imageSteps.ts`, `stepform/ImageFields` |
| Toạ độ / neo / DPI | `script/anchor.go` ↔ `anchor.ts` (cùng công thức), `coordMap.ts` (overlay ↔ robotgo), `window.go` (màn hình) |
| Overlay (công cụ, gợi ý, vẽ marker) | Tính: `overlay/geometry.ts` / `hint.ts`; vẽ: `components/overlay/`; chuột / phím + state: `OverlayScreen` |
| Thao tác canvas | `flow/edit.ts` (thuần) → `store/nodes.ts` → `FlowEditor` |
| Hiển thị khi chạy | `flow/run.ts` (thuần), `runStore.ts`, `components/flow/RunLog.tsx` |
| Định dạng file | `model.ScriptVersion`, `script` (`codec.go`, `check.go`), `store/scriptfile.go` |
| UI chưa có logic | `components/steps/`, `AppChrome.tsx`, dữ liệu giả `src/demo/` (màn `#demo*`) |
| Mô tả cho AI / MCP | `internal/model/steptypes.go` (có test bắt đủ loại node) |
| Kịch bản mẫu / trang sân thử | `internal/samples/scripts/*.json` (JSON gọn `{meta, script}`), `pages/*.html`; test bắt mỗi loại node ≥ 1 mẫu; UI `components/samples/`, lọc `samples/filter.ts` |

## 7. Test

- Go: `go test ./...` – engine chạy bằng `Input` / `Screen` / `TextReader` giả; kịch bản danh sách bọc `scripttest.Linear(sc)`; helper test dùng chung ở `modeltest` / `scripttest` / `visiontest` / `devicetest`; `testdata/flow_order.json` dùng chung Go + TS.
- TS: `cd frontend && npm test` (vitest, logic thuần `flow/`, `anchor`, `transform`, `scriptInputs`); `npm run build` để type-check.
- Ví dụ kịch bản thật: `examples/*.json`.

## 8. Bản đồ file

**Go** – package theo miền trong `internal/` (bước 32), phụ thuộc một chiều: `model` ← `vision` ← `device`; `model`, `vision`, `launch` ← `script` ← `store`; `model`, `vision`, `device`, `launch`, `script` ← `engine`; tất cả ← `main`. Package con không import `main` hay import ngược.

`internal/model` – kiểu dữ liệu thuần, không thiết bị / Wails:
- `model.go` – `Step` = node, `Script` có `start` + `edges`, `Retry`, `ScriptVersion`, `NewStep` / `NewScript` (mặc định). `flow.go` – `StepResult`, `StepError` + mã lỗi, cổng `PortsOf`, `FlowOrder`, `EdgeOut` (test chung `testdata/flow_order.json` với `flow/model.ts`). `nodes.go` – cổng node điều khiển, `EndStatus` / `VarOp`. `cond.go` – `Cond.Kind`, `Probe`. `inputs.go` – `ScriptInput`, `ScriptError`. `steptypes.go` – mô tả + ví dụ JSON từng loại node cho AI / MCP (`StepTypeDocs`, `Guide`). `run.go` – `RunRecord` / `RunSummary` (một lần chạy, `Add` giữ tối đa `MaxRunResults`).

`internal/model/modeltest` – dựng node / cạnh / điều kiện cho test (`Node`, `Edge`, `Step`, `TextCond`). `internal/script/scripttest` – `Linear` (kịch bản danh sách → flow), `TextScript`.

`internal/vision` – tìm ảnh / chữ thuần trên ảnh đã chụp: `match.go` (NCC, `MatchScaled`, `ImageScale`), `ocr.go` (`TextLine`, `FindText`, `ReadingOrder`), `image.go` (PNG base64, màu `#RRGGBB`, `Shrink`). `visiontest/` – ảnh giả dùng chung cho test (`Scene`, `Button`, `Crop`).

`internal/device` – thiết bị thật qua interface: `input.go` (`Input` + robotgo, `Displays`, `MainScreen`), `capture.go` / `capture_darwin.go` / `capture_other.go` (`Screen`, quyền Screen Recording, chụp không lẫn cửa sổ app), `ocr.go` / `ocr_darwin.go` / `ocr_other.go` (`TextReader`; Vision, Windows chưa có), `failsafe.go` (dừng khi chuột chạm góc, `TrackingInput`). `devicetest/` – `FakeInput` cho test.

`cmd/d0probe` – code thử bước 25d D0 (build tag `probe`, ngoài build thường): đo toạ độ trang → màn hình qua CDP + robotgo; xoá khi D1 có test thật.

`internal/samples` – kịch bản mẫu để test (bước 33): `samples.go` (`List`, `Open` – điền mặc định `NewStep`, nối cạnh theo thứ tự, thay `{{sandbox}}` / ảnh mẫu `@file`), `sandbox.go` (`WriteSandbox`, `FileURL`), `scripts/*.json`, `pages/` (trang sân thử HTML offline + `logo.png`). `samples_real_test.go` – chạy mẫu bằng engine trên máy thật (`AUTOCLICKER_REAL=1`).

`internal/launch` – `open.go`: mở URL / file / app (`OpenTarget`, `openCommand` theo hệ điều hành).

`internal/script` – kịch bản như dữ liệu (không thiết bị, không đĩa): `codec.go` (`EncodeScript` / `DecodeScript`, `SafeFileName`), `check.go` (kiểm tra lúc mở: `checkScript`, `checkStep`, `FixRefs`), `validate.go` (`ValidateScript` – đủ dữ liệu để chạy), `flowcheck.go` (`validateFlow`, `UnreachableNodes`), `control.go` (kiểm tra node điều khiển + điều kiện), `inputs.go` / `inputvalue.go` (`ResolveInputs`, ép kiểu, ghim trường), `varfields.go` (`VarFields` – mọi ô chữ nhận `{{biến}}`), `anchor.go` (`ResolvePoint`, `Bounds`), `errors.go` (`stepErr` → `model.ScriptError`).

`internal/store` – dữ liệu app trên máy (`ConfigDir`, `WriteFileAtomic`): `scriptfile.go` (đọc / ghi file kịch bản), `draft.go` (nháp tự lưu), `settings.go` (theme), `recent.go` (mở gần đây), `inputmemory.go` (`inputs.json` – giá trị input lần trước / preset), `runs.go` (lịch sử chạy: `SaveRun`, `RunsOf`, `GetRun`, `RunShot`, `safeName`).

`internal/engine` – chạy kịch bản: `engineapi.go` (`New` + `Options`, `Failed`, `Shots`, `FindOnce`, `FindResult`), `engine.go` (`Engine`, `Progress`), `flowrun.go` (`Run`, `runLoop`, `execNode`, thử lại), `actions.go` (`runStep` / `doStep` – chuột, phím, cuộn, mở), `imagestep.go` (chờ ảnh / chữ / pixel, click theo ảnh), `nodes.go` (`runControl`), `cond.go` (`evalCond`, `snapshot` chụp một lần mỗi lần xét), `vars.go` (`resolveNodeVars`), `gate.go` (tạm dừng / debug, `Sleep`), `nodetest.go` (`TestNode`, `NodeTest`), `shots.go` (ảnh lúc node lỗi).

`internal/winexclude/` – Windows: loại cửa sổ app khỏi ảnh chụp.

**Root (package main)** – chỉ Wails: method `App` bind sang frontend, mỗi method vài dòng gọi xuống `internal/`.

- `main.go`, `app.go` – Wails app, `App`; method mỏng: `NewStep` / `NewScript`, mô tả loại node, theme, input đã nhớ (`frontend/wailsjs/` tự sinh – không sửa tay, không cần đọc; namespace TS theo package Go: `model.Step`, `store.Settings`, `engine.FindResult`, `main.RunOptions`).
- `runner.go` – `RunScript` / Stop / Pause / Resume / Step, `RunOptions`, lưu lịch sử chạy, method xem lịch sử. `samplesapi.go` – Mẫu thử (`Samples`, `OpenSample`, `OpenSandboxPage`). `imageapi.go` – chụp ảnh mẫu, Thử tìm, Thử node, lấy màu, quyền chụp (ẩn cửa sổ lúc chụp – `withHidden`).
- `storage.go` – hộp thoại lưu / mở, nháp, hỏi lưu khi đóng. `window.go` – chế độ cửa sổ, màn hình, `ScreenBounds`. `hotkey.go` – F6 / F8 / F10 (gohook).

**Frontend (`frontend/src/`)**

- `App.tsx` chọn màn theo hash; `screens/` – `MainScreen`, `OverlayScreen` (state + chuột / phím; vẽ ở `components/overlay/`), `MiniBar`, `PreviewScreen`, màn mẫu `Demo*`.
- Overlay: `overlay/` (thuần, có test) – `geometry` (marker, đường neo, mũi tên cuộn, chỗ đặt nút khoanh ảnh), `hint` (chữ gợi ý theo công cụ); `components/overlay/` (chỉ props) – `OverlayToolbar`, `Hud` (gợi ý, crosshair, ô toạ độ), `MarkerLayer` (đường nối, ô ngẫu nhiên, marker), `Shapes` (vùng, khung Thử tìm, khung khoanh ảnh), `hooks` (`useRefState`, thanh tự ẩn, chọn màn hình, `useCapture` khoanh ảnh).
- `types.ts` – model TS (khớp `internal/model/model.go`, tsc kiểm), nhãn loại node `STEP_LABEL`.
- Logic thuần (có test): `flow/` – `model` (cổng, `flowOrder`, `orderedSteps`), `edit` (nối / cắt cạnh, chèn / xoá node, `nextRef`, `pruneEdges`), `layout` (elkjs, `placeNode`), `view` (thẻ node, nhãn cổng, `outputsOf`), `cond` (điều kiện lồng nhau `CondAt`), `run` (kết quả chạy → trạng thái thẻ, nhật ký); `anchor.ts`; `transform.ts` (biến đổi nhiều điểm); `overlay/` (hình học + gợi ý overlay); `scriptInputs.ts` (ghim trường, input, `mapVarText`); `samples/filter.ts` (lọc Mẫu thử); `fold.ts` (so khớp không dấu).
- `store/` (import qua `store/index.ts`): `editor` (state, hoàn tác, `updateScript`, `replaceScript`), `geometry` (khung, neo ↔ màn hình), `nodes` (thêm / sửa / xoá / nhân bản, nối cạnh, `insertAt`), `control` (node điều khiển, điều kiện), `points` (điểm trên overlay, F6, kéo thả), `multi` (chọn nhiều + biến đổi, phím công cụ Chọn nhiều `selectToolKey`, đổi vùng).
- Store khác: `engineStore.ts` (trạng thái engine, `runScript`), `runStore.ts` (nhật ký lần chạy + lịch sử), `runForm.ts` (Chạy / F8 → form Chạy, preset), `scriptFile.ts` (lưu / mở / nháp / cờ chưa lưu), `imageSteps.ts` (khoanh ảnh mẫu, thử tìm, lấy màu, quyền chụp), `coordMap.ts` (overlay ↔ robotgo), `windowMode.ts` (`setWindowMode`, `openOverlay`), `theme.ts`.
- `components/`: `TitleBar`, `Library` (thư viện + lịch sử), `AddStepPalette` (Ctrl/Cmd+K, `PaletteView`), `StepForm` (inspector) + `stepform/` (`shared`, `MouseFields`, `ScrollFields`, `KeyFields`, `TextFields`, `ImageFields`, `FlowSections`, `ControlSection`, `InputsTab`), `RunSettings` (lặp / đếm ngược, `Switch`), `Simulation` (xem trước con trỏ giả), `MockPage`, `TransformTool`, `SaveChangesDialog`, `NumInput`, `Segmented`, `icons` (`StepIcon`), `StepCard`.
- `components/flow/` (chỉ props, app + màn mẫu dùng chung): `FlowCanvas` (`@xyflow/react`), `NodeCard`, `FlowEdge`, `PortField`, `RunLog`, `model.ts` (cỡ thẻ), `nodes/ControlFields` + `CondEditor` + `NodeTestPanel`; `FlowEditor` = bản nối store.
- `components/samples/` – hộp thoại Mẫu thử: `SamplesDialog` + `SampleCard` (chỉ props), `Samples` (nối Go + `scriptFile.openSample`).
- `components/inputs/` – ghim, ô chữ `{{biến}}` (`VarTextInput`), tab Input, form Chạy (`RunDialog`).
- UI dựng sẵn chưa có logic (bước 27b, *Sắp có*): `components/steps/` (kiểu tạm `steps/model.ts`), `ScheduleDialog`, `AppChrome` (`Toolbar`, `StatusBar`). Dữ liệu giả: `demo/` (màn `#demo*`, `?sel=<id>`, `?mode=running|failed`).
- `style.css` – Tailwind v4 + token (`@theme`, chép từ `design/theme.css`).

**Khác:** `design/` (`DESIGN.md` token / thành phần / icon, `theme.css`, `newdarkmode.md` bảng Slate), `docs/tasks/` (plan), `examples/` (kịch bản mẫu), `testdata/`, `scripts/` (`tasks-index.py`, `check-docs.py`, `session-status.sh`), `build/` (icon, plist, manifest cho `wails build` – `build/bin` là bản build, bỏ qua).
