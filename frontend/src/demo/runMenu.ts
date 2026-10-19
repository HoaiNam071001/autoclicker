// Menu Chạy của màn mẫu (#demo, #demo-flow): mục bấm được, chỉ hiện thông báo.
import { RUN_ICON, RunMenuItem } from '../components/RunMenu';

export function demoRunItems(notify: (msg: string) => void, recording?: { on: boolean; toggle: () => void }): RunMenuItem[] {
  return [
    { label: 'Chạy từ node đang chọn', icon: RUN_ICON.from, onClick: () => notify('Màn mẫu – chạy từ node đang chọn') },
    { label: 'Chạy riêng node đang chọn', icon: RUN_ICON.only, onClick: () => notify('Màn mẫu – chạy riêng node đang chọn') },
    { label: 'Chạy từng node', icon: RUN_ICON.step, hotkey: 'F10', onClick: () => notify('Màn mẫu – chạy từng node, F10 để tiếp') },
    { label: 'Xem trước trên màn hình', icon: RUN_ICON.preview, onClick: () => notify('Màn mẫu – xem trước') },
    { label: recording?.on ? 'Dừng ghi' : 'Ghi thao tác', icon: RUN_ICON.record, hotkey: 'F9', danger: true,
      disabled: !recording, title: recording ? undefined : 'Ghi macro – sắp có (bước 24)', onClick: recording?.toggle },
    { label: 'Sửa input trước khi chạy', minor: true, onClick: () => notify('Màn mẫu – mở form input') },
  ];
}
