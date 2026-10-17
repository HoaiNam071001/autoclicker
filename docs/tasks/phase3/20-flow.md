# 20. Flow kiểu n8n – mỗi bước là một module, rẽ nhánh theo kết quả

- [ ] **20. Flow kiểu n8n – mỗi bước là một module, rẽ nhánh theo kết quả** *(người dùng chốt 2026-10-17, thay bản "cây khối" 20a–20c đã revert)* – Kịch bản thành **đồ thị node + cạnh**: mỗi node chạy xong trả về một `StepResult` chung (status, cổng ra, output, lỗi, log); cạnh nối từ **cổng ra** của node này sang node kế → luồng thành công / thất bại, if / else, switch case, vòng lặp đều là cạnh. Soạn trên **canvas kéo nối như n8n**. Gồm cả debug + lịch sử chạy (25c cũ, ở F5).
  - **Thứ tự làm:** F0 → **F4a** (thiết kế + `#demo-flow`, dữ liệu giả) → F1 (model) → F2 (engine) → **F4b** (nối canvas vào app, bỏ danh sách bước cũ) → F3 (node điều khiển: engine + inspector cùng lúc) → F5 (chạy + nhật ký + debug) → F6 (xem trước, thanh mini) → F7 (MCP). Lý do: duyệt giao diện sớm khi sửa còn rẻ; canvas thay danh sách ngay sau khi engine chạy được để app không bị "nửa cũ nửa mới" lâu
  - **File cũ:** bỏ, không chuyển đổi *(người dùng chốt 2026-10-17)* – `ScriptVersion` 4, file version < 4 và `draft.json` cũ báo "Kịch bản tạo bằng bản cũ, không mở được" (nháp cũ xoá luôn). Bỏ hẳn `onFail.action` / `goto` (thay bằng cạnh `fail`), không giữ code tương thích
  - **Nguyên tắc:**
    - `Step` giữ nguyên làm **node** (mọi trường bước hiện có dùng lại), `Script.steps` = danh sách node (thứ tự trong mảng không có nghĩa chạy); thêm `Script.start` + `Script.edges`
    - Cổng ra không nối: `fail` → **dừng kịch bản, báo lỗi**; `success` / cổng khác → **kết thúc thành công**. Node tắt → `skipped`, đi theo cổng `success`
    - Vòng lặp = cạnh quay lại; chặn lặp vô hạn bằng node *Đếm* + giới hạn 10 000 lần chạy node mỗi vòng
    - Node sau đọc kết quả node trước bằng biến: `{{n2.output.x}}`, `{{n2.status}}`, `{{n2.error.message}}`, `{{last.output.text}}`, `{{vars.ten}}` – `n2` là mã ngắn của node (`ref`, tự đặt n1, n2…, sửa được, duy nhất)
    - *Click theo kết quả bước trước* → chọn node nguồn (`imageSource: "node"`, `imageFrom` = **id** node – cạnh / start / imageFrom đều trỏ theo id, `ref` chỉ dùng trong biến); *Kéo từ chỗ chuột dừng* lấy vị trí chuột lúc chạy (chưa node chuột nào chạy thì lỗi)
    - Lặp N vòng cả kịch bản, tạm dừng / dừng / F8, failsafe giữ nguyên
  - **Chia module** (mỗi file một việc; phần thuần có test riêng; Go và TS cùng tên khái niệm):

    | Module | Go | Frontend | Việc |
    |---|---|---|---|
    | Model flow | `flow.go` | `types.ts` (+ `flow.ts`) | `Edge`, `StepResult`, `StepError`, `portsOf(step)`, hằng cổng / status / mã lỗi |
    | Thứ tự hiển thị | `flow.go` `flowOrder` | `flow.ts` `flowOrder` | đánh số node (marker, Xem trước, form Chạy) – hai bên cùng công thức, test cùng dữ liệu |
    | Kiểm tra | `flowcheck.go` | – | `validateFlow` (cạnh, cổng, `ref`, biến, `imageFrom`, node không tới được) |
    | Engine | `flowrun.go` | – | `runFlow`, `execNode`, thử lại, giới hạn lặp, event |
    | Biến lúc chạy | `vars.go` | `flowVars.ts` | thay `{{nX.output…}}` / `{{vars…}}`, gợi ý biến cho ô nhập |
    | Node điều khiển | `nodes.go` (+ `cond.go`) | `components/flow/nodes/*` | `if`, `switch`, `switchScreen`, `counter`, `end`, `setVar`; Thử node |
    | Sửa đồ thị | – | `flowEdit.ts` | thêm / xoá / nối / cắt cạnh, chèn node giữa cạnh, copy / paste, đặt bắt đầu (thuần, có test) |
    | Xếp tự động | – | `flowLayout.ts` | `elkjs` → vị trí node |
    | Canvas | – | `components/flow/FlowCanvas.tsx`, `NodeCard.tsx`, `FlowEdge.tsx`, `PortField.tsx` | `@xyflow/react`, thẻ node, cạnh có nhãn cổng, mục Cổng ra trong inspector |
    | Nhật ký chạy | `runlog.go` | `components/flow/RunLog.tsx`, `runStore.ts` | lưu / đọc lịch sử chạy, bảng nhật ký |
    | Store | – | `scriptStore.ts` | gọi `flowEdit` + hoàn tác / làm lại (không tự sửa cạnh trong store) |
    | Mẫu | – | `demo/flowDemo.ts`, màn `#demo-flow` | dữ liệu giả đủ mọi trạng thái |

  - **Ví dụ JSON** (đích của F1–F3):
    ```json
    { "version": 4, "start": "a1",
      "steps": [
        { "id": "a1", "ref": "n1", "type": "open", "target": "facebook.com" },
        { "id": "a2", "ref": "n2", "type": "waitImage", "name": "Tìm nút Nhắn tin", "findBy": "text", "findText": "Nhắn tin", "timeoutMs": 3000 },
        { "id": "a3", "ref": "n3", "type": "clickImage", "imageSource": "node", "imageFrom": "a2" },
        { "id": "a4", "ref": "n4", "type": "scroll", "direction": "down", "amount": 3 },
        { "id": "a5", "ref": "n5", "type": "counter", "max": 10 },
        { "id": "a6", "ref": "n6", "type": "end", "endStatus": "fail", "message": "Không thấy nút Nhắn tin" } ],
      "edges": [
        { "from": "a1", "port": "success", "to": "a2" }, { "from": "a2", "port": "success", "to": "a3" },
        { "from": "a2", "port": "fail", "to": "a5" },    { "from": "a5", "port": "next", "to": "a4" },
        { "from": "a5", "port": "done", "to": "a6" },    { "from": "a4", "port": "success", "to": "a2" } ] }
    ```
  - **Loại node và cổng ra:**

    | Nhóm | Node | Cổng ra | Output |
    |---|---|---|---|
    | Hành động | các loại bước hiện có (click, kéo, gõ, phím, cuộn, mở, chờ ảnh / chữ / pixel, click theo ảnh / chữ) | `success` / `fail` | chờ ảnh / chữ: `{x, y, w, h, score, text}`; chờ pixel: `{color}`; click: `{x, y}`; mở: `{target}` |
    | Rẽ nhánh | `if` **Nếu** – điều kiện ảnh / chữ / màu / so giá trị / và-hoặc / đảo | `true` / `false` | từng điều kiện lá đúng / sai + chi tiết |
    | | `switch` **Switch theo giá trị** – `{{biến}}` so với từng case (bằng / chứa / khớp mẫu) | `case:<key>` / `default` | `{value, case}` |
    | | `switchScreen` **Switch theo màn hình** – chờ tối đa X s, cái nào trong N ảnh / chữ hiện trước | `case:<key>` / `none` | `{case, x, y, w, h, score}` |
    | Luồng | `counter` **Đếm** – đếm số lần đi qua (đặt lại mỗi vòng kịch bản) | `next` (còn < max) / `done` | `{count}` |
    | | `end` **Kết thúc** – thành công / thất bại + lời nhắn | – | – |
    | | `setVar` **Đặt biến** – gán / cộng / lấy từ output node khác | `success` | `{name, value}` |
    | Sau này | `subflow` – gọi một kịch bản khác như một node (thay cho Nhóm) | `success` / `fail` | output của node `end` bên trong |

  - [x] **F0. Gỡ bước 20 cũ** *(2026-10-17)*: revert 3 commit 20a / 20b / 20c (`60dcfe2`, `956ec76`, `d5882a7`), giữ 2 commit sửa lỗi sau đó (overlay theo màn chứa vùng, `stepLabel`). Code điều kiện (`evalCond`, `compareValues`, `condtest.go` / `TestCond`, `CondEditor`) **lấy lại từ lịch sử git** ở F3, không viết lại
  - [x] **F4a. Thiết kế canvas + màn mẫu `#demo-flow`** (dữ liệu giả, chưa nối Go)
    - [x] `design/DESIGN.md` thêm phần *Canvas flow*: nền canvas (lưới chấm), thẻ node (cỡ, ô icon, tên, mã `n2`, chip tóm tắt), cổng vào / ra + nhãn cổng (`success` màu `ok`, `fail` màu `danger`, case trung tính), cạnh (cong, mũi tên, nhãn, nút × khi rê), node bắt đầu / tắt / đang chọn / cảnh báo / đang chạy / xong / lỗi, minimap + nút zoom, cả 2 theme; token mới (nếu có) thêm cả `theme.css` + `style.css` – không cần token mới; thêm 6 icon node điều khiển
    - [x] Cài `@xyflow/react` + `elkjs`; `components/flow/NodeCard.tsx`, `FlowEdge.tsx`, `FlowCanvas.tsx` nhận dữ liệu qua props (dùng lại được ở F4b). Thêm `components/flow/PortField.tsx` (Cổng ra, Khi thất bại, Output, Lần chạy gần nhất), `RunLog.tsx`, `nodes/ControlFields.tsx`; `flowEdit.ts` (có test) + `flowLayout.ts` (elkjs nạp lười) làm sớm vì màn mẫu dùng luôn; `AddStepPalette` tách `PaletteView`
    - [x] `demo/flowDemo.ts`: ví dụ Facebook ở trên + 1 node If, 1 Switch theo màn hình, 1 Đếm, 1 Kết thúc, node tắt, node cảnh báo; công tắc xem trạng thái "đang chạy" (node chạy / xong / lỗi, cạnh đã đi) và bảng Nhật ký mẫu (`?mode=running` / `?mode=failed`)
    - [x] Màn `#demo-flow`: bố cục 3 cột (thư viện | canvas | inspector có mục Cổng ra), Ctrl+K khi kéo cổng ra chỗ trống (sửa được thật trong state cục bộ: kéo node, nối / xoá cạnh, thả cổng ở chỗ trống → bảng thêm node, inspector)
    - [x] Chụp headless cả 2 theme (soạn / đang chạy / thất bại); `npm run build`, `npm test` (31) sạch
    - [x] Người dùng duyệt `#demo-flow` (OK 2026-10-17)
  - [x] **F1. Model flow + kết quả** (Go + TS, chưa UI – sửa bằng JSON / test)
    - [x] `model.go`: `Script.start`, `Script.edges []Edge {id, from, port, to}`; `Step.ref`, `Step.name`, `Step.pos {x, y}`, `Step.imageFrom`; `onFail` thành `retry {count, delayMs}`
    - [x] `flow.go`: `StepResult {nodeId, ref, status: success | fail | skipped | cancelled, port, output, error {code, message}, logs [{t, level, msg}], attempts, startedAt, durationMs}`; mã lỗi `not_found | timeout | permission | invalid_input | cancelled | internal`; `portsOf`; `flowOrder` (duyệt từ `start`, ưu tiên `success` / `true` / case đầu)
    - [x] `storage.go`: `ScriptVersion` 4, file < 4 báo không mở được; `NewScript` tạo sẵn `start`
    - [x] `flowcheck.go` `validateFlow`: có `start`; cạnh trỏ node có thật, cổng đúng loại node, mỗi cổng tối đa 1 cạnh; `ref` duy nhất; `{{nX…}}` trỏ node có thật; `imageFrom` là node chờ ảnh / chữ; node không tới được → cảnh báo (không chặn chạy). Lỗi trả `nodeId` + `ref`
    - [x] `types.ts` khớp Go (tsc kiểm tra), `flow.ts` (`portsOf`, `flowOrder`); `steptypes.go` ghi cổng ra + output mỗi loại, ví dụ qua được `validateFlow`
    - [x] Test: `validateFlow` từng lỗi, `flowOrder` có vòng lặp (Go + TS cùng kết quả – `testdata/flow_order.json`), file cũ bị từ chối, sửa cạnh hỏng khi mở file
    - Ghi chú F1: `ScriptError` thêm `ref` (lỗi ghi "n2: …"); `NewStep` id thêm bộ đếm (gọi liên tiếp từng trùng id); `App.FlowGuide()` (cấu trúc flow + kịch bản mẫu cho AI, có test chạy qua `validateScript`); lỗi thiếu quyền chụp màn hình là `errNoPermission` (mã `permission`)
    - **Cầu nối tạm** `linearToFlow` (flow.go): danh sách bước cũ chưa có start / cạnh → nối thành chuỗi success theo thứ tự mảng, đặt ref, click theo ảnh lấy node chờ ảnh gần nhất phía trước – để màn danh sách cũ vẫn chạy được tới F4b. **Xoá ở F4b**
  - [x] **F2. Engine chạy theo đồ thị** (`flowrun.go`, `vars.go`)
    - [x] `runFlow`: từ `start` → `execNode` trả `StepResult` → cạnh theo `port` → node kế; hết cạnh kết thúc theo quy tắc. `runStep` bọc lại để trả `StepResult`; thử lại theo `retry`, mỗi lần ghi log; thay `runWithRetry` / `goto` cũ
    - [x] Ngữ cảnh chạy: `results[nodeId]`, `vars`, `last`; thay `{{nX.output…}}` / `{{vars…}}` lúc chạy từng node (input `{{key}}` vẫn qua `resolveInputs`); click theo ảnh lấy chỗ thấy từ output node `imageFrom`
    - [x] Giới hạn 10 000 lần chạy node mỗi vòng → lỗi "có thể lặp vô hạn ở n2 → n5 → n4"
    - [x] Event: `engine:progress` thêm `nodeId`, `ref`, `port`; event mới `engine:node` gửi `StepResult`
    - [x] Test bằng input giả: `fail` có / không nối, thử lại rồi qua `fail`, vòng lặp, node tắt, biến từ output, dừng / tạm dừng, giới hạn 10 000 (`flowrun_test.go`; test cũ chạy qua cầu nối)
    - Ghi chú F2: làm cùng F1 (đổi model là phải đổi vòng chạy). `runStep` trả output theo loại (chuột: x, y; gõ: text; mở: target; chờ ảnh / chữ: x, y, w, h, score, text; click theo ảnh: clickX, clickY; chờ pixel: color); `delayAfter` chỉ chờ khi node thành công; mã lỗi `not_found` / `permission` / `invalid_input` / `internal`; UI danh sách cũ: mục *Khi không thấy* chỉ còn Thử lại N lần + nghỉ (rẽ nhánh làm trên canvas)
  - [x] **F4b. Nối canvas vào app** (bỏ danh sách bước cũ)
    - [x] `MainScreen`: cột giữa là `components/flow/FlowEditor.tsx` (nối store → `FlowCanvas`); bỏ `StepList.tsx`; cầu nối `linearToFlow` chỉ còn là helper test `linear()` (`flowtest_test.go`). `@dnd-kit` vẫn giữ cho kéo sắp xếp input (`InputsPanel`)
    - [x] Tổ chức lại *(người dùng yêu cầu 2026-10-17: tách module rõ ràng, quy tắc ghi vào CLAUDE.md – mục Quy tắc code)*: `src/flow/` (`model`, `edit`, `layout`, `view` – thuần, có test); `scriptStore.ts` (606 dòng) tách thành `src/store/` (`editor`, `geometry`, `nodes`, `points`, `multi`, `index`); `StepForm.tsx` (751 dòng) tách thành khung + `components/stepform/` (`shared`, `MouseFields`, `ScrollFields`, `KeyFields`, `TextFields`, `ImageFields`, `FlowSections`, `InputsTab`)
    - [x] Store: `ref` cho node mới / nhân bản (`nextRef`), node đầu tiên thành `start`, node mới đặt cạnh node nguồn (`placeNode`); chỗ nối node mới `insertAt` (thả cổng ở chỗ trống / Ctrl+K sau cổng success của node đang chọn / đặt nhiều điểm liên tiếp trên overlay nối tiếp nhau / F6); đổi loại node bỏ cạnh của cổng không còn (`pruneEdges`); xoá node kèm cạnh (hoàn tác bằng ⌘Z thay cho nút "Hoàn tác" cũ)
    - [x] `flow/edit.ts` (có test): nối / cắt cạnh, xoá node kèm cạnh, chèn node giữa cạnh, `nextRef`, `insertNode`, `pruneEdges`, `chain`; `flow/layout.ts`: xếp elkjs + `placeNode`
    - [x] Thao tác: kéo node (1 lần hoàn tác), kéo cổng sang node = nối, thả cổng ở chỗ trống → bảng thêm node (dòng "Nối từ n2 · cổng Thất bại"), Ctrl+K / nút "+ Thêm node", Delete xoá node / cạnh, Shift / ⌘ + click node = chọn nhiều điểm (bước 30), vừa màn hình / zoom / minimap / Sắp xếp lại; kịch bản mở ra chưa có vị trí thì tự xếp (không tính là sửa)
    - [x] Inspector: tên node sửa tại chỗ, mã node sửa được (báo trùng / sai dạng), *Khi thất bại* dùng chung `RetryPanel`, **Cổng ra** (chọn node đích), **Output** (bấm chép `{{n2.output.x}}`, `outputsOf` khớp `steptypes.go`), Đặt làm bắt đầu / Xoá node; click theo ảnh chọn node nguồn trong các node Chờ ảnh / Chờ chữ
    - [x] Overlay / xem trước / thanh mini / form Chạy dùng mã node (n2) và thứ tự `flowOrder` thay cho số thứ tự mảng; `endpoints` (điểm trước của kéo thả) theo `flowOrder`
    - [x] Thẻ node: cảnh báo không đi tới được / click theo ảnh thiếu nguồn (`flow/view.ts`), node đang chạy / node làm lần chạy trước dừng
    - [x] Copy / paste node + cạnh bên trong (Ctrl+C / V, dán lệch dần, đổi `imageFrom` / biến `{{nX.…}}` trỏ node đã chép sang bản sao – `flow/edit.ts` `copyNodes` / `pasteNodes`, nhân bản dùng chung); Shift + kéo khung / Shift + click chọn nhiều node (`picked` trong store, kèm điểm cho bước 30), Ctrl+A chọn hết, Delete xoá cả nhóm một lần hoàn tác
    - [x] `VarTextInput` gợi ý cả output các node (`{{n2.output.x}}`, status, error.message – `flow/view.ts` `nodeVarOptions`), tô màu biến node, mã node không có tô đỏ
    - Kiểm tra tay: xem [checks.md](../checks.md)
  - [x] **F3. Node điều khiển** (`nodes.go`, `cond.go`, `components/flow/nodes/*`) *(2026-10-17)*
    - [x] Lấy lại `Cond`, `evalCond`, `compareValues` (từ `60dcfe2`) và `CondEditor` (từ `d5882a7`) cho `if`; điều kiện `value` (so 2 vế, nhận biến node) thay `input` / `result` cũ; `if` chờ theo `timeoutMs` của node (xét cả cây điều kiện mỗi lần, hết giờ = sai) thay `waitMs` từng điều kiện
    - [x] `switch` (eq / contains / regex), `switchScreen` (mỗi lần xét chụp + đọc chữ một lần cho mọi case – `snapshot`), `counter` (lần 1..max → next, lần max+1 → done rồi đếm lại), `end` (fail = lỗi mã `end`), `setVar` (gán / cộng) – mỗi loại: engine + `checkStep` + `validateScript` + `portsOf` (Go + TS) + `steptypes.go` + inspector + thẻ node + test
    - [x] **Thử node** (`nodetest.go` `TestNode`, thay `condtest.go`): `if` / `switch` / `switchScreen` lúc này → cổng sẽ đi + chi tiết từng phần, biến node → "chỉ biết lúc chạy", khung chỗ thấy trên overlay
    - [x] Ctrl+K nhóm *Rẽ nhánh & luồng*; Ctrl+K / đặt điểm liên tiếp nối vào cổng đầu của node (true / next…) thay vì luôn success
    - [x] Màn mẫu `#demo-flow` dùng chung inspector thật (`components/flow/nodes/ControlFields`, `CondEditor`, `NodeTestPanel` – chỉ nhận props); thêm `?sel=<id>` để mở sẵn inspector node
    - Ghi chú F3: `varFields` (vars.go) – một chỗ liệt kê ô chữ nhận biến (input, biến node, kiểm tra biến), TS `mapVarText` (`scriptInputs.ts`) cùng danh sách; `cloneStep` chép sâu điều kiện / case trước khi thay biến; Click theo ảnh lấy được chỗ thấy từ node Nếu / Switch theo màn hình (`findsImage`); `flowExample` cho AI dùng Đếm + Kết thúc; test `nodes_test.go` (12), `flow/cond.test.ts`, `testdata/flow_order.json` thêm case node điều khiển
    - File thử tay: `examples/thu-node-dieu-khien.json` (Google tìm từ khoá – đủ 6 node điều khiển, qua `validateScript`)
    - Kiểm tra tay: xem [checks.md](../checks.md)
  - [ ] **F5. Lúc chạy + nhật ký + debug** *(thay bước 25c cũ)* (`runlog.go`, `RunLog.tsx`, `runStore.ts`)
    - [ ] Canvas khi chạy: node đang chạy viền `accent` + vòng tiến độ, xong tô `ok` / `danger`, cạnh vừa đi sáng lên, số lần đi qua trên node
    - [ ] Lời nhắn của node Kết thúc **thành công** hiện ra (thanh trạng thái / thanh mini / nhật ký) – hiện chỉ ghi vào log node, chưa chỗ nào đọc (thất bại thì đã hiện ở "Lỗi lần chạy trước")
    - [ ] Bảng **Nhật ký** dưới canvas: mỗi dòng một node đã chạy (giờ, `n2`, tên, status, ms, cổng đi) → bấm mở output / lỗi / log / ảnh lúc lỗi; bấm dòng → chọn node
    - [ ] Lưu mỗi lần chạy `UserConfigDir/autoclicker/runs/<kịch bản>/<thời gian>.json` (`StepResult[]`, input đã dùng, nguồn: người dùng / hẹn giờ / CLI / AI, kết quả cuối) + PNG khi node lỗi / hết giờ; giữ 50 lần gần nhất mỗi kịch bản; mục *Lịch sử* trong thư viện mở lại (tô đường đã đi)
    - [ ] Debug: chạy từ node đang chọn, chạy từng node (Tiếp / F10), breakpoint; `RunNode` là hàm Go riêng (Thử node và MCP dùng chung)
  - [ ] **F6. Xem trước, thanh mini**
    - [ ] Xem trước (`Simulation.tsx`) đi theo đồ thị: mặc định `success` / `true` / case đầu; node rẽ nhánh có công tắc "giả sử thất bại" / chọn case; Đếm chạy 1 vòng rồi `done`; đường đang mô phỏng tô trên canvas
    - [ ] Thanh mini: "n2 · Tìm nút Nhắn tin · thất bại → n5 Đếm (3 / 10)"
  - [ ] **F7. Chuẩn bị cho MCP (31):** JSON flow là định dạng duy nhất; `get_step_types` trả cổng ra + output từng loại; `validateFlow` lỗi có `nodeId` / `ref`; kết quả chạy trả `StepResult[]` để AI đọc lỗi và tự sửa kịch bản
  - Sau này (28): đọc chữ trong vùng / clipboard vào biến (`setVar` từ OCR), lặp theo từng dòng CSV, `subflow`
  - Xong khi: ví dụ Facebook ở trên (thấy "Nhắn tin" → click; không thấy → cuộn 3 nấc, lặp tối đa 10 lần rồi kết thúc báo lỗi) soạn hết bằng canvas, chạy thật được, xem được nhật ký từng node; `go test ./...`, `npm test`, `npm run build` sạch
