package model

// Giá trị hợp lệ của các trường liệt kê trong Step / Script (script/check.go kiểm tra lúc mở file).

// Nguồn vị trí của bước click theo ảnh (Step.ImageSource).
const (
	ImageFromNode   = "node"   // chỗ node ImageFrom (chờ ảnh / chữ, click theo ảnh tự tìm) tìm thấy lần chạy gần nhất
	ImageFromSearch = "search" // tự tìm ảnh mẫu riêng của bước
)

// Vùng chụp của node Chụp màn hình (Step.ShotArea).
const (
	ShotRegion = ""       // vùng của kịch bản (không có vùng: màn chính)
	ShotScreen = "screen" // cả một màn: Display (0 = màn chứa vùng kịch bản / màn chính)
	ShotRect   = "rect"   // vùng cố định ShotRect (toạ độ màn hình)
	ShotNode   = "node"   // khung chỗ thấy của node ImageFrom (output x, y = tâm, w, h)
)

var ShotAreas = []string{ShotRegion, ShotScreen, ShotRect, ShotNode}

// Mức thao tác giống người (Script.Human, Step.Human – bước 34, internal/humanize). Kịch bản "" = tắt; node "" = theo
// kịch bản. Kịch bản mới (NewScript) bật sẵn mức vừa.
const (
	HumanOff    = "off"
	HumanLow    = "low"
	HumanMedium = "medium"
	HumanHigh   = "high"
)

var HumanLevels = []string{HumanOff, HumanLow, HumanMedium, HumanHigh}

// Hướng cuộn (Step.Direction).
var ScrollDirections = []string{"up", "down", "left", "right"}

// Cách bấm khi bước phím có nhiều phím chính (Step.KeyMode).
const (
	KeyChord    = "chord"    // nhấn giữ phím kèm + mọi phím chính cùng lúc rồi nhả (vd Ctrl+C+T)
	KeySequence = "sequence" // giữ phím kèm, gõ lần lượt từng phím chính rồi nhả phím kèm
)

// Phép của node Clipboard (Step.ClipOp).
const (
	ClipRead  = ""      // đọc chữ đang có vào output text
	ClipWrite = "write" // ghi Text (rỗng = xoá clipboard)
	ClipImage = "image" // ghi ảnh PNG từ FilePath
)

var ClipOps = []string{ClipRead, ClipWrite, ClipImage}

// Thao tác của node Tệp / thư mục (Step.FileOp). trash / emptyTrash là thao tác xoá – chạy phải được người dùng xác nhận
// trước (script.DeleteOps, RunOptions.ConfirmDelete).
const (
	FileWrite      = ""           // ghi Text vào FilePath (FileAppend: thêm vào cuối)
	FileRow        = "row"        // thêm một dòng Cells vào file .csv / .xlsx (chưa có thì tạo)
	FileMkdir      = "mkdir"      // tạo thư mục FilePath (cả thư mục cha)
	FileTrash      = "trash"      // chuyển file khớp FilePattern trong thư mục FilePath vào Thùng rác
	FileEmptyTrash = "emptyTrash" // dọn Thùng rác (xoá hẳn)
)

var FileOps = []string{FileWrite, FileRow, FileMkdir, FileTrash, FileEmptyTrash}

// Thao tác của node Hệ thống (Step.SysOp).
const (
	SysVolume     = "volume"     // đặt âm lượng = Volume %
	SysVolumeUp   = "volumeUp"   // tăng Volume % (0 = 10)
	SysVolumeDown = "volumeDown" // giảm Volume % (0 = 10)
	SysMute       = "mute"
	SysUnmute     = "unmute"
)

var SysOps = []string{SysVolume, SysVolumeUp, SysVolumeDown, SysMute, SysUnmute}

// Giờ hẹn đã qua lúc tới node Chờ đến giờ (Step.IfPassed).
const (
	PassedGo       = ""         // đi tiếp ngay (output lateMs > 0)
	PassedTomorrow = "tomorrow" // chờ tới giờ đó ngày mai – giờ không ghi ngày
	PassedFail     = "fail"     // node thất bại (code timeout)
)

var PassedModes = []string{PassedGo, PassedTomorrow, PassedFail}
