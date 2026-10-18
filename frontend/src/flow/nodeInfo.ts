// Mô tả từng loại node cho người dùng (bước 35 N0) – bảng thêm node (khung dưới) + inspector (đầu form).
// Chữ cho người dùng, khác internal/model/steptypes.go (cho AI / MCP, đủ trường). Thêm loại node → thêm một dòng ở đây (test bắt thiếu).
import type { Step, StepType } from '../types';

// desc: làm gì; use: dùng khi nào / ví dụ.
export type NodeInfo = { desc: string; use: string };
// Biến thể hiện như loại riêng trên UI: Chờ chữ / Click theo chữ (findBy text), Mở trang bằng Chrome của app.
export type InfoKey = StepType | 'waitText' | 'clickText' | 'openWeb' | 'focusWindow';

export const NODE_INFO: Record<InfoKey, NodeInfo> = {
  click: { desc: 'Bấm chuột tại một điểm trên màn hình, một hoặc nhiều lần.', use: 'Bấm nút, mở mục, chọn ô – khi chỗ cần bấm luôn nằm cố định.' },
  double: { desc: 'Nhấp đúp chuột tại một điểm.', use: 'Mở file / thư mục, chọn cả một từ, phóng to cửa sổ.' },
  hold: { desc: 'Nhấn giữ chuột tại một điểm trong một khoảng thời gian rồi nhả.', use: 'Nút phải giữ mới chạy (giữ để xác nhận, giữ để ghi âm).' },
  drag: { desc: 'Nhấn chuột ở điểm A, kéo tới điểm B rồi nhả.', use: 'Kéo thẻ kanban, kéo file vào thư mục, kéo thanh trượt, bôi đen chữ.' },
  wait: { desc: 'Đứng yên một khoảng thời gian rồi chạy tiếp.', use: 'Đợi trang / app tải khi không có ảnh hay chữ nào để chờ.' },
  type: { desc: 'Gõ một đoạn chữ vào ô đang được chọn, nhận {{biến}}.', use: 'Điền form, ô tìm kiếm, tin nhắn – thường đặt sau node Click vào ô.' },
  key: { desc: 'Bấm một phím hoặc tổ hợp phím (Enter, Ctrl+C, Cmd+Space…).', use: 'Gửi bằng Enter, copy / paste, chuyển tab, mở ô tìm app của máy.' },
  scroll: { desc: 'Lăn chuột lên / xuống / sang ngang một số nấc, tại chỗ chuột hoặc một điểm.', use: 'Xem thêm nội dung trong trang, danh sách dài, feed.' },
  open: { desc: 'Mở trang web, file hoặc app trên máy; URL nhận {{biến}}.', use: 'Đầu kịch bản: mở YouTube với từ khoá, bật Spotify, mở file Excel.' },
  openWeb: { desc: 'Mở trang web bằng Chrome riêng của app để các node Web tìm được phần tử trong trang.', use: 'Đặt trước các node Chờ / Click / Gõ / Đọc phần tử web.' },
  waitImage: { desc: 'Chờ tới khi ảnh mẫu (nút, biểu tượng) hiện trên màn hình, cho biết vị trí tìm thấy.', use: 'Đợi trang tải xong, đợi hộp thoại hiện, kiểm tra đã đăng nhập chưa.' },
  waitText: { desc: 'Chờ tới khi một dòng chữ hiện trên màn hình (đọc chữ bằng OCR), cho biết vị trí.', use: 'Đợi thông báo “Thành công”, đợi tên nút hiện – không cần chụp ảnh mẫu.' },
  clickImage: { desc: 'Tìm ảnh mẫu trên màn hình rồi bấm vào đó.', use: 'Nút không nằm cố định (đổi chỗ theo cửa sổ, theo trang đã cuộn).' },
  clickText: { desc: 'Tìm một dòng chữ trên màn hình rồi bấm vào đó.', use: 'Bấm nút / liên kết theo tên hiện trên màn, vd “Đăng nhập”, “Xem thêm”.' },
  waitPixel: { desc: 'Chờ tới khi một điểm trên màn hình có đúng màu cho trước.', use: 'Đèn trạng thái đổi màu, thanh tiến độ đầy, nút chuyển từ xám sang xanh.' },
  screenshot: { desc: 'Chụp vùng kịch bản, cả màn, vùng cố định hoặc đúng chỗ node trước thấy, lưu PNG; chép vào clipboard được.', use: 'Lưu bằng chứng sau khi đặt hàng / chuyển tiền, chụp kết quả để gửi đi.' },
  clipboard: { desc: 'Đọc chữ đang có trong clipboard vào biến, ghi chữ (nhận {{biến}}) hoặc ghi ảnh PNG từ file.', use: 'Copy / paste: ghép với Phím tắt Cmd/Ctrl+C, Cmd/Ctrl+V; chép ảnh vừa chụp để dán vào chat.' },
  http: { desc: 'Gọi API (GET / POST / PUT / PATCH / DELETE) kèm tham số, header, thân JSON / form; đọc kết quả JSON vào biến.', use: 'Lấy dữ liệu từ server, gửi đơn / tin nhắn qua API, báo kết quả chạy về Telegram / webhook.' },
  file: { desc: 'Ghi chữ vào file, thêm dòng vào bảng .csv / .xlsx, tạo thư mục, xoá file vào Thùng rác hoặc dọn Thùng rác.', use: 'Lưu nhật ký / kết quả từng vòng vào Excel, dọn ảnh chụp cũ trong thư mục.' },
  system: { desc: 'Chỉnh âm lượng loa của máy: đặt mức, tăng / giảm, tắt / bật tiếng.', use: 'Bật nhạc nền vừa đủ nghe, tắt tiếng trước khi họp, giảm âm lượng ban đêm.' },
  waitUntil: { desc: 'Chờ tới đúng giờ hẹn (tới ms, có thể kèm ngày) rồi chạy tiếp; canh được theo giờ chuẩn internet.', use: 'Săn sale lúc 20:00:00, gửi tin nhắn đúng giờ, chạy việc mỗi ngày lúc 9:00.' },
  path: { desc: 'Giữ chuột kéo một nét liền qua nhiều điểm, hoặc vẽ hình tròn / chữ nhật / sao quanh một tâm.', use: 'Vẽ trong app vẽ, ký tên, mở khoá bằng hình vẽ, khoanh vùng trên bản đồ.' },
  webWait: { desc: 'Chờ phần tử trong trang web (theo chữ, nhãn, CSS…) hiện ra hoặc biến mất.', use: 'Đợi trang tải xong, đợi kết quả tìm kiếm – nhanh và chắc hơn chờ ảnh.' },
  webClick: { desc: 'Tìm phần tử trong trang web theo chữ / nhãn rồi bấm vào bằng chuột.', use: 'Trang hay đổi bố cục / theme khiến ảnh mẫu không còn khớp.' },
  webType: { desc: 'Tìm ô nhập trong trang web, bấm vào rồi gõ chữ (xoá chữ cũ nếu muốn).', use: 'Điền form đăng ký, ô tìm kiếm, bình luận.' },
  elementRead: { desc: 'Đọc chữ / giá trị của phần tử trong cửa sổ app đang mở (Chrome thường, app khác) – chữ chính xác, không qua OCR.', use: 'Lấy giá, số dư, nội dung ô nhập trên trang đã đăng nhập sẵn (macOS).' },
  webRead: { desc: 'Đọc chữ, giá trị hoặc thuộc tính của phần tử trong trang web vào biến.', use: 'Lấy giá, mã đơn, số dư, link để node sau dùng hoặc so sánh.' },
  webScroll: { desc: 'Cuộn trang web tới phần tử, cuộn trong một khung, hoặc cuộn tới khi thấy phần tử.', use: 'Đưa nút cuối trang vào tầm nhìn, lăn feed tới bài cần tìm.' },
  focusWindow: { desc: 'Đưa cửa sổ của một app lên trước để thao tác tiếp vào đúng app đó.', use: 'Chuyển qua lại giữa nhiều app trong một kịch bản.' },
  if: { desc: 'Kiểm tra điều kiện (thấy chữ / ảnh / màu, so giá trị) rồi đi nhánh Đúng hoặc Sai.', use: 'Đã đăng nhập thì bỏ qua bước đăng nhập; hết hàng thì dừng.' },
  switch: { desc: 'So một giá trị (input, kết quả node trước) với nhiều trường hợp, mỗi trường hợp một nhánh.', use: 'Rẽ theo trạng thái đơn hàng, theo lựa chọn nhập lúc Chạy.' },
  switchScreen: { desc: 'Xem màn hình đang hiện gì (theo chữ / ảnh) rồi đi nhánh tương ứng.', use: 'Trang có thể ra nhiều kiểu: màn đăng nhập, captcha, trang chính…' },
  counter: { desc: 'Đếm số lần đi qua: còn lượt thì ra cổng Còn lượt, đủ số lần thì ra Hết lượt.', use: 'Lặp một đoạn N lần, giới hạn số lần thử lại.' },
  end: { desc: 'Kết thúc vòng chạy, ghi thành công hoặc thất bại kèm lời nhắn.', use: 'Cuối một nhánh: báo “Đã đặt hàng”, hoặc dừng với lỗi “Hết hàng”.' },
  setVar: { desc: 'Đặt hoặc cộng một biến để node sau đọc bằng {{vars.tên}}.', use: 'Đếm số đơn đã xử lý, ghi lại kết quả giữa các nhánh.' },
};

// Mô tả cho một node cụ thể (biến thể theo cách tìm / trình duyệt).
export function infoKey(s: Pick<Step, 'type' | 'findBy' | 'imageSource' | 'webBrowser'>): InfoKey {
  if (s.findBy === 'text' && s.type === 'waitImage') return 'waitText';
  if (s.findBy === 'text' && s.type === 'clickImage' && s.imageSource === 'search') return 'clickText';
  if (s.type === 'open' && s.webBrowser) return 'openWeb';
  return s.type;
}

export const nodeInfo = (s: Pick<Step, 'type' | 'findBy' | 'imageSource' | 'webBrowser'>): NodeInfo => NODE_INFO[infoKey(s)];
