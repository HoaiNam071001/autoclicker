// Bảng thêm bước (design v2, artboard AddStep) – mở bằng "+ Thêm bước" hoặc Ctrl/Cmd+K.
// Danh sách nhóm / loại node + nối store; phần hiển thị ở PaletteView. Node mới nối theo chỗ đã đặt (store insertAt).
import { PaletteGroup, PaletteView } from './PaletteView';
import { addControl, addPointStep, addStep, getEditor, insertTarget } from '../store';
import { portView } from '../flow/view';
import { WEB_BASIC_NOTE } from '../flow/web';
import { NODE_INFO } from '../flow/nodeInfo';
import { addElementStep, addImageStep, addTextStep } from '../imageSteps';

// Nhóm theo flow/category (cùng màu thẻ node) – tên nhóm lấy từ CATEGORIES
const GROUPS: PaletteGroup[] = [
  {
    cat: 'mouse', items: [
      { icon: 'click', label: 'Click', info: NODE_INFO.click, keywords: 'bam nhan chuot trai phai', run: () => addPointStep('click') },
      { icon: 'double', label: 'Double click', info: NODE_INFO.double, keywords: 'nhay dup hai lan', run: () => addPointStep('double') },
      { icon: 'hold', label: 'Giữ chuột', info: NODE_INFO.hold, keywords: 'hold nhan giu', run: () => addPointStep('hold') },
      { icon: 'drag', label: 'Kéo thả', info: NODE_INFO.drag, keywords: 'drag keo tha a b', run: () => addPointStep('drag') },
      { icon: 'path', label: 'Vẽ đường', info: NODE_INFO.path, keywords: 've duong net hinh tron chu nhat sao ky ten draw path shape', run: () => addStep('path') },
      { icon: 'scroll', label: 'Cuộn', info: NODE_INFO.scroll, keywords: 'scroll lan chuot', run: () => addStep('scroll') },
      { icon: 'clickImage', label: 'Click theo ảnh', info: NODE_INFO.clickImage, keywords: 'image anh mau tim nut', run: () => addImageStep('clickImage') },
      { icon: 'clickImage', label: 'Click theo chữ', info: NODE_INFO.clickText, keywords: 'text chu ocr doc nut tim', run: () => addTextStep('clickImage') },
      { icon: 'clickImage', label: 'Click theo phần tử', info: NODE_INFO.clickElement, keywords: 'phan tu accessibility nut chrome thuong app aria role bam', run: () => addElementStep('clickImage') },
    ],
  },
  {
    cat: 'kb', items: [
      { icon: 'type', label: 'Gõ chữ', info: NODE_INFO.type, keywords: 'type text nhap van ban', run: () => addStep('type') },
      { icon: 'key', label: 'Phím tắt', info: NODE_INFO.key, keywords: 'key hotkey enter ctrl to hop', run: () => addStep('key') },
    ],
  },
  {
    cat: 'wait', items: [
      { icon: 'wait', label: 'Chờ', info: NODE_INFO.wait, keywords: 'wait delay nghi doi', run: () => addStep('wait') },
      { icon: 'waitImage', label: 'Chờ ảnh', info: NODE_INFO.waitImage, keywords: 'wait image anh mau xuat hien khoanh', run: () => addImageStep('waitImage') },
      { icon: 'waitImage', label: 'Chờ chữ', info: NODE_INFO.waitText, keywords: 'wait text chu ocr doc xuat hien', run: () => addTextStep('waitImage') },
      { icon: 'waitImage', label: 'Chờ phần tử', info: NODE_INFO.waitElement, keywords: 'phan tu accessibility chrome thuong app aria role xuat hien', run: () => addElementStep('waitImage') },
      { icon: 'findAll', label: 'Tìm tất cả', info: NODE_INFO.findAll, keywords: 'tim tat ca moi cho nhieu dem danh sach find all nut lap tung', run: () => addStep('findAll') },
      { icon: 'waitScreen', label: 'Chờ màn hình', info: NODE_INFO.waitScreen, keywords: 'cho man hinh dung yen thay doi tai xong load spinner loading on dinh doi', run: () => addStep('waitScreen') },
      { icon: 'waitUntil', label: 'Chờ đến giờ', info: NODE_INFO.waitUntil, keywords: 'cho den gio hen gio dong ho time schedule san sale flash ntp', run: () => addStep('waitUntil') },
      { icon: 'waitPixel', label: 'Chờ màu pixel', info: NODE_INFO.waitPixel, keywords: 'pixel mau color diem', run: () => addImageStep('waitPixel') },
      { icon: 'screenshot', label: 'Chụp màn hình', info: NODE_INFO.screenshot, keywords: 'screenshot chup anh man hinh luu png clipboard bang chung', run: () => addStep('screenshot') },
    ],
  },
  {
    // bước 35: node nâng cao – hệ thống, tệp, API
    cat: 'data', items: [
      { icon: 'clipboard', label: 'Clipboard', info: NODE_INFO.clipboard, keywords: 'clipboard copy paste dan chep bo nho tam doc ghi', run: () => addStep('clipboard') },
      { icon: 'http', label: 'HTTP / API', info: NODE_INFO.http, keywords: 'http api rest get post put delete webhook json goi server request', run: () => addStep('http') },
      { icon: 'file', label: 'Tệp / thư mục', info: NODE_INFO.file, keywords: 'file tep thu muc folder ghi luu csv excel xlsx xoa thung rac don dep', run: () => addStep('file') },
      { icon: 'system', label: 'Âm lượng', info: NODE_INFO.system, keywords: 'he thong am luong volume loa tat tieng mute sound', run: () => addStep('system') },
      { icon: 'calc', label: 'Tính toán', info: NODE_INFO.calc, keywords: 'tinh toan bieu thuc cong tru nhan chia lam tron so sanh calc math round', run: () => addControl('calc') },
      { icon: 'textNumber', label: 'Số & ngày', info: NODE_INFO.textNumber, keywords: 'lay so gia tien dinh dang so ngay gio thang nam number date format nghin', run: () => addControl('textNumber') },
      { icon: 'textMatch', label: 'Lọc mẫu', info: NODE_INFO.textMatch, keywords: 'loc mau regex khop trich lay ma don email so dien thoai match', run: () => addControl('textMatch') },
      { icon: 'textSplit', label: 'Tách chữ', info: NODE_INFO.textSplit, keywords: 'tach chu chia danh sach dau phay dong split', run: () => addControl('textSplit') },
      { icon: 'textClean', label: 'Làm sạch chữ', info: NODE_INFO.textClean, keywords: 'lam sach chu thay the bo dau hoa thuong khoang trang trim replace lower upper', run: () => addControl('textClean') },
      { icon: 'setVar', label: 'Đặt biến', info: NODE_INFO.setVar, keywords: 'set var bien gan', run: () => addControl('setVar') },
    ],
  },
  {
    cat: 'app', items: [
      { icon: 'open', label: 'Mở URL / app', info: NODE_INFO.open, keywords: 'url web link trinh duyet ung dung open youtube', run: () => addStep('open') },
      { icon: 'window', label: 'Cửa sổ', info: NODE_INFO.window, keywords: 'window cua so ung dung focus dua len truoc dat vi tri co kich thuoc thu nho phong to dong cho mo', run: () => addStep('window') },
      { icon: 'app', label: 'Ứng dụng', info: NODE_INFO.app, keywords: 'app ung dung chuong trinh dang chay kiem tra cho mo thoat tat buoc thoat kill quit process tien trinh', run: () => addStep('app') },
      { icon: 'browserTab', label: 'Tab trình duyệt', info: NODE_INFO.browserTab, keywords: 'tab trinh duyet chrome edge brave coc coc dong tab mo tab moi chuyen tab tai lai quay lai url tieu de browser', run: () => addStep('browserTab') },
      { icon: 'browserGo', label: 'Đi tới URL', info: NODE_INFO.browserGo, keywords: 'di toi url doi trang dia chi trong tab hien tai chrome trinh duyet go navigate tai trang', run: () => addStep('browserGo') },
      { icon: 'browserNorm', label: 'Chuẩn hoá trình duyệt', info: NODE_INFO.browserNorm, keywords: 'chuan hoa trinh duyet zoom 100 toan man hinh fullscreen dat co cua so popup dich trang chrome toa do lech', run: () => addStep('browserNorm') },
      { icon: 'pageText', label: 'Đọc chữ trang', info: NODE_INFO.pageText, keywords: 'doc chu trang toan bo noi dung van ban danh sach bang gia accessibility chrome thuong app lay het chu text', run: () => addStep('pageText') },
      { icon: 'elementRead', label: 'Đọc phần tử', info: NODE_INFO.elementRead, keywords: 'doc phan tu accessibility chrome thuong app gia tri o nhap gia tien lay du lieu', run: () => addStep('elementRead') },
    ],
  },
  {
    cat: 'flow', items: [
      { icon: 'if', label: 'Nếu', info: NODE_INFO.if, keywords: 'if dieu kien re nhanh else thay chu anh mau so sanh', run: () => addControl('if') },
      { icon: 'switchScreen', label: 'Switch theo màn hình', info: NODE_INFO.switchScreen, keywords: 'switch case man hinh thay', run: () => addControl('switchScreen') },
      { icon: 'switch', label: 'Switch theo giá trị', info: NODE_INFO.switch, keywords: 'switch case gia tri bien', run: () => addControl('switch') },
      { icon: 'forEach', label: 'Lặp qua danh sách', info: NODE_INFO.forEach, keywords: 'foreach lap qua danh sach moi phan tu dong csv excel xlsx mang json loop list', run: () => addControl('forEach') },
      { icon: 'counter', label: 'Đếm', info: NODE_INFO.counter, keywords: 'counter lap loop dem gioi han so lan thu', run: () => addControl('counter') },
      { icon: 'end', label: 'Kết thúc', info: NODE_INFO.end, keywords: 'end stop dung ket thuc', run: () => addControl('end') },
    ],
  },
  {
    // bước 25d: tìm phần tử bằng HTML qua Chrome của app, bấm bằng chuột thật – để cuối, chỉ cho trang cơ bản
    cat: 'web', note: WEB_BASIC_NOTE, items: [
      { icon: 'open', label: 'Mở trang (Chrome của app)', info: NODE_INFO.openWeb, keywords: 'web url trinh duyet chrome dieu khien dom html', run: () => addStep('open', { webBrowser: true }) },
      { icon: 'webWait', label: 'Chờ phần tử web', info: NODE_INFO.webWait, keywords: 'web dom html selector css xpath aria role cho xuat hien mat', run: () => addStep('webWait') },
      { icon: 'webClick', label: 'Click phần tử web', info: NODE_INFO.webClick, keywords: 'web dom html selector css xpath aria role nut bam click', run: () => addStep('webClick') },
      { icon: 'webType', label: 'Gõ vào ô web', info: NODE_INFO.webType, keywords: 'web dom html o nhap input go chu dien form dang nhap', run: () => addStep('webType') },
      { icon: 'webRead', label: 'Đọc phần tử web', info: NODE_INFO.webRead, keywords: 'web dom html doc chu gia tri thuoc tinh href bien lay du lieu', run: () => addStep('webRead') },
      { icon: 'webKey', label: 'Phím trên trang web', info: NODE_INFO.webKey, keywords: 'web dom phim key enter tab esc chrome nen', run: () => addStep('webKey') },
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
