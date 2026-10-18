// Hộp thoại Mẫu thử với dữ liệu giả cho màn #demo (trình duyệt không có Go) – một phần danh sách thật ở internal/samples.
import SamplesDialog from '../components/samples/SamplesDialog';
import { Sample } from '../samples/filter';

const s = (id: string, name: string, app: string, level: number, page: string, steps: string[], description: string, extra: Partial<Sample> = {}): Sample =>
  ({ id, name, app, level, page, steps, description, real: false, ...extra });

const DEMO_SAMPLES: Sample[] = [
  s('click-nhieu-lan', 'Click nhiều lần', 'Cơ bản', 1, 'click.html', ['open', 'waitImage', 'click'],
    'Mở Bảng click, click trái 10 lần cách nhau 150 ms. Đúng: ô Click trái = 10, nhịp ≈ 150 ms.',
    { prep: 'Trình duyệt mở to gần toàn màn hình (điểm click neo theo giữa màn hình).' }),
  s('keo-tha', 'Kéo thả', 'Cơ bản', 1, 'draw.html', ['open', 'waitImage', 'drag'],
    'Vẽ hình vuông trên Bảng vẽ: kéo tuyệt đối, kéo tương đối và kéo tiếp từ chỗ chuột dừng.'),
  s('phim-tat', 'Phím tắt', 'Cơ bản', 1, 'type.html', ['open', 'waitImage', 'click', 'type', 'key'],
    'Gõ hai dòng bằng Enter, bấm Tab, Shift+A, rồi gõ lần lượt X Y Z.'),
  s('cho-mau-pixel', 'Chờ màu pixel', 'Chờ / tìm', 1, 'signal.html', ['open', 'waitPixel', 'end'],
    'Chờ điểm dưới giữa màn hình đổi sang xanh #16a34a.'),
  s('cho-anh-logo', 'Chờ ảnh + click theo ảnh', 'Chờ / tìm', 2, 'image.html', ['open', 'waitImage', 'clickImage', 'counter', 'end'],
    'Tìm logo ở vị trí ngẫu nhiên, click vào nó; logo nhảy chỗ khác, lặp 3 lần.', { prep: 'Trình duyệt để zoom 100%.' }),
  s('switch-man-hinh', 'Switch theo màn hình', 'Điều khiển', 2, 'popup.html', ['open', 'switchScreen', 'clickImage', 'counter', 'end'],
    'Màn nào hiện trước thì đi nhánh đó: quảng cáo → đóng, đăng nhập → Tiếp tục, tới khi thấy Trang chủ.'),
  s('wiki-tim-bai', 'Wikipedia – tìm bài', 'Web thật', 2, 'https://vi.wikipedia.org', ['open', 'clickImage', 'type', 'key', 'waitImage'],
    'Mở Wikipedia tiếng Việt, gõ từ khoá vào ô tìm, Enter, chờ tiêu đề bài.', { real: true }),
];

export default function DemoSamples({ onClose }: { onClose: () => void }) {
  return <SamplesDialog samples={DEMO_SAMPLES} ready={false} onClose={onClose} onPick={() => {}} onOpenPage={() => {}} />;
}
