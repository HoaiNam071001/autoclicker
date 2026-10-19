// Dữ liệu cho phần giao diện các loại bước chưa làm logic (bước 27b). Khi làm logic, từng kiểu chuyển vào
// types.ts + model.go: cuộn / mở (23), chờ ảnh / click theo ảnh / chờ pixel (25), hết giờ (25b), focus cửa sổ (26),
// nhóm (20), neo điểm (29). Tới lúc đó chỉ màn #demo dùng.
import { Point } from '../../types';

export type ScrollData = { direction: 'up' | 'down' | 'left' | 'right'; amount: number; at: Point | null }; // at null = vị trí chuột hiện tại
export type OpenData = { target: string; app: string }; // app '' = mặc định của hệ điều hành
export type WaitImageData = {
  template: { w: number; h: number } | null; // ảnh mẫu đã khoanh (null = chưa có)
  threshold: number; // % khớp tối thiểu
  timeoutMs: number;
  intervalMs: number;
  area: 'region' | 'screen'; // tìm trong vùng của kịch bản / toàn màn hình
};
export type ClickImageData = { source: 'prev' | 'search'; dx: number; dy: number; button: 'left' | 'right' | 'center' };
export type WaitPixelData = { at: Point; color: string; tolerance: number; timeoutMs: number };
export type OnTimeout = { action: 'retry' | 'skip' | 'stop'; retries: number; then: 'skipGroup' | 'stop' };
export type GroupData = { name: string; repeat: number };

// Cửa sổ đang mở (mẫu) cho ô chọn của Focus cửa sổ
