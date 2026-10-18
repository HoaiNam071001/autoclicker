// Bảng thêm bước (design v2, artboard AddStep) – mở bằng "+ Thêm bước" hoặc Ctrl/Cmd+K.
// Danh sách nhóm / loại node + nối store; phần hiển thị ở PaletteView. Node mới nối theo chỗ đã đặt (store insertAt).
import { PaletteGroup, PaletteView } from './PaletteView';
import { addControl, addPointStep, addStep, getEditor, insertTarget } from '../store';
import { portView } from '../flow/view';
import { WEB_BASIC_NOTE } from '../flow/web';
import { NODE_INFO } from '../flow/nodeInfo';
import { addImageStep, addTextStep } from '../imageSteps';

const GROUPS: PaletteGroup[] = [
  {
    name: 'Chuột', items: [
      { icon: 'click', label: 'Click', info: NODE_INFO.click, keywords: 'bam nhan chuot trai phai', run: () => addPointStep('click') },
      { icon: 'double', label: 'Double click', info: NODE_INFO.double, keywords: 'nhay dup hai lan', run: () => addPointStep('double') },
      { icon: 'hold', label: 'Giữ chuột', info: NODE_INFO.hold, keywords: 'hold nhan giu', run: () => addPointStep('hold') },
      { icon: 'drag', label: 'Kéo thả', info: NODE_INFO.drag, keywords: 'drag keo tha a b', run: () => addPointStep('drag') },
      { icon: 'scroll', label: 'Cuộn', info: NODE_INFO.scroll, keywords: 'scroll lan chuot', run: () => addStep('scroll') },
      { icon: 'clickImage', label: 'Click theo ảnh', info: NODE_INFO.clickImage, keywords: 'image anh mau tim nut', run: () => addImageStep('clickImage') },
      { icon: 'clickImage', label: 'Click theo chữ', info: NODE_INFO.clickText, keywords: 'text chu ocr doc nut tim', run: () => addTextStep('clickImage') },
    ],
  },
  {
    name: 'Bàn phím', items: [
      { icon: 'type', label: 'Gõ chữ', info: NODE_INFO.type, keywords: 'type text nhap van ban', run: () => addStep('type') },
      { icon: 'key', label: 'Phím tắt', info: NODE_INFO.key, keywords: 'key hotkey enter ctrl to hop', run: () => addStep('key') },
    ],
  },
  {
    name: 'Chờ & điều kiện', items: [
      { icon: 'wait', label: 'Chờ', info: NODE_INFO.wait, keywords: 'wait delay nghi doi', run: () => addStep('wait') },
      { icon: 'waitImage', label: 'Chờ ảnh', info: NODE_INFO.waitImage, keywords: 'wait image anh mau xuat hien khoanh', run: () => addImageStep('waitImage') },
      { icon: 'waitImage', label: 'Chờ chữ', info: NODE_INFO.waitText, keywords: 'wait text chu ocr doc xuat hien', run: () => addTextStep('waitImage') },
      { icon: 'waitPixel', label: 'Chờ màu pixel', info: NODE_INFO.waitPixel, keywords: 'pixel mau color diem', run: () => addImageStep('waitPixel') },
    ],
  },
  {
    name: 'Ứng dụng & cấu trúc', items: [
      { icon: 'open', label: 'Mở URL / app', info: NODE_INFO.open, keywords: 'url web link trinh duyet ung dung open youtube', run: () => addStep('open') },
      { icon: 'focusWindow', label: 'Focus cửa sổ', info: NODE_INFO.focusWindow, keywords: 'window cua so ung dung', soon: 'bước 26' },
    ],
  },
  {
    name: 'Rẽ nhánh & luồng', items: [
      { icon: 'if', label: 'Nếu', info: NODE_INFO.if, keywords: 'if dieu kien re nhanh else thay chu anh mau so sanh', run: () => addControl('if') },
      { icon: 'switchScreen', label: 'Switch theo màn hình', info: NODE_INFO.switchScreen, keywords: 'switch case man hinh thay', run: () => addControl('switchScreen') },
      { icon: 'switch', label: 'Switch theo giá trị', info: NODE_INFO.switch, keywords: 'switch case gia tri bien', run: () => addControl('switch') },
      { icon: 'counter', label: 'Đếm', info: NODE_INFO.counter, keywords: 'counter lap loop dem gioi han so lan thu', run: () => addControl('counter') },
      { icon: 'end', label: 'Kết thúc', info: NODE_INFO.end, keywords: 'end stop dung ket thuc', run: () => addControl('end') },
      { icon: 'setVar', label: 'Đặt biến', info: NODE_INFO.setVar, keywords: 'set var bien gan', run: () => addControl('setVar') },
    ],
  },
  {
    // bước 25d: tìm phần tử bằng HTML qua Chrome của app, bấm bằng chuột thật – để cuối, chỉ cho trang cơ bản
    name: 'Web (trình duyệt)', note: WEB_BASIC_NOTE, items: [
      { icon: 'open', label: 'Mở trang (Chrome của app)', info: NODE_INFO.openWeb, keywords: 'web url trinh duyet chrome dieu khien dom html', run: () => addStep('open', { webBrowser: true }) },
      { icon: 'webWait', label: 'Chờ phần tử web', info: NODE_INFO.webWait, keywords: 'web dom html selector css xpath aria role cho xuat hien mat', run: () => addStep('webWait') },
      { icon: 'webClick', label: 'Click phần tử web', info: NODE_INFO.webClick, keywords: 'web dom html selector css xpath aria role nut bam click', run: () => addStep('webClick') },
      { icon: 'webType', label: 'Gõ vào ô web', info: NODE_INFO.webType, keywords: 'web dom html o nhap input go chu dien form dang nhap', run: () => addStep('webType') },
      { icon: 'webRead', label: 'Đọc phần tử web', info: NODE_INFO.webRead, keywords: 'web dom html doc chu gia tri thuoc tinh href bien lay du lieu', run: () => addStep('webRead') },
      { icon: 'webScroll', label: 'Cuộn trang web', info: NODE_INFO.webScroll, keywords: 'web dom html cuon lan scroll feed toi khi thay', run: () => addStep('webScroll') },
    ],
  },
];

// onCancel: đóng mà không chọn (Esc, bấm ra ngoài) – bỏ chỗ nối đã đặt (store insertAt).
export default function AddStepPalette({ onClose, onCancel }: { onClose: () => void; onCancel?: () => void }) {
  const at = insertTarget();
  const from = at && getEditor().script.steps.find((s) => s.id === at.from);
  return (
    <PaletteView groups={GROUPS} onClose={onClose} onCancel={onCancel}
      context={from ? <>Nối từ <span className="font-mono">{from.ref}</span> · cổng {portView(at.port, from).label}</> : 'Node mới đứng riêng – kéo từ cổng của node khác để nối'}
      footer={<><span><span className="font-mono">Enter</span> thêm node</span><span className="ml-auto">Click / Double / Giữ / Kéo thả mở overlay chọn điểm</span></>} />
  );
}
