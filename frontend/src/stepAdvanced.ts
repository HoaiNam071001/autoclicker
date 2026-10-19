// Trường riêng của các node nâng cao (bước 35) – Step kế thừa (types.ts), khớp struct Step trong internal/model/model.go
// (tsc kiểm qua stepMatchesGo). Tách khỏi types.ts cho file không phình theo số loại node.
import type { ShotArea } from './flow/shot';
import type { ClipOp } from './flow/clipboard';
import type { BodyType, HttpMethod } from './flow/http';
import type { FileOp } from './flow/file';
import type { SysOp } from './flow/system';
import type { IfPassed } from './flow/waitUntil';
import type { Shape } from './flow/shape';
import type { WinMatch, WinOp, WinPreset } from './flow/window';
import type { KV, Rect } from './types';

export interface AdvancedFields {
  // chụp màn hình (bước 35 N1, internal/engine/shot.go)
  shotArea?: ShotArea; // bỏ trống = vùng của kịch bản
  display?: number; // shotArea screen: 0 = màn chứa vùng kịch bản, n = màn thứ n
  shotRect?: Rect; // shotArea rect: toạ độ màn hình
  saveDir?: string; // trống = Pictures/autoclicker; nhận {{biến}}
  fileName?: string; // trống = autoclicker-<ngày>-<giờ>; nhận {{biến}}
  toClipboard?: boolean;
  // clipboard (bước 35 N2, internal/engine/clip.go)
  clipOp?: ClipOp; // bỏ trống = đọc chữ
  filePath?: string; // clipOp image: file PNG; nhận {{biến}}
  // HTTP (bước 35 N3, internal/engine/http.go) – mọi ô chữ nhận {{biến}}
  method?: HttpMethod; // bỏ trống = GET
  url?: string; // biến chèn nguyên văn – giá trị cần mã hoá đặt trong query
  query?: KV[];
  headers?: KV[];
  bodyType?: BodyType; // bỏ trống = không body
  body?: string; // json / text
  form?: KV[]; // bodyType form
  // tệp / thư mục (bước 35 N4, internal/engine/file.go) – filePath + text dùng chung
  fileOp?: FileOp; // bỏ trống = ghi chữ
  fileAppend?: boolean; // ghi chữ: thêm vào cuối
  cells?: string[]; // row: ô của dòng mới; nhận {{biến}}
  sheet?: string; // row xlsx: tên sheet, trống = sheet đầu
  filePattern?: string; // trash: mẫu tên file, vd *.png
  // hệ thống (bước 35 N5, internal/engine/system.go)
  sysOp?: SysOp;
  volume?: number; // volume: mức 0–100; volumeUp / volumeDown: mỗi lần (0 = 10)
  // chờ đến giờ (bước 35 N6, internal/engine/waituntil.go)
  at?: string; // HH:mm[:ss[.mmm]] (kèm ngày nếu cần); nhận {{biến}}
  ifPassed?: IfPassed;
  useNtp?: boolean;
  // vẽ đường (bước 35 N7, internal/engine/draw.go) – points, button, dragMs dùng chung
  shape?: Shape; // bỏ trống = qua các điểm
  shapeW?: number;
  shapeH?: number; // 0 = bằng shapeW
  closed?: boolean; // qua các điểm: nối điểm cuối về đầu
  // cửa sổ (bước 26, internal/engine/window.go) – timeoutMs: chờ cửa sổ khớp; display: màn của winPreset (0 = màn chứa cửa sổ)
  winOp?: WinOp; // bỏ trống = focus
  winApp?: string; // tên app / tên exe / bundle id; nhận {{biến}}
  winTitle?: string; // tiêu đề; nhận {{biến}}
  winMatch?: WinMatch; // bỏ trống = chứa (không phân biệt hoa thường, dấu)
  winPreset?: WinPreset; // setBounds: bỏ trống = đúng winRect
  winRect?: Rect; // setBounds: khung đích (preset trái / giữa / phải: chỉ lấy w, h)
}
