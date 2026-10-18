package device

// Trash – Thùng rác của hệ điều hành (node Tệp / thư mục, bước 35 N4): chuyển file vào (khôi phục được), dọn hẳn.
// Qua interface để test engine bằng bản giả (devicetest.FakeTrash) – test không bao giờ xoá file thật.
type Trash interface {
	MoveToTrash(paths []string) error
	EmptyTrash() error
}

// SystemTrash – macOS: đổi tên vào ~/.Trash, dọn qua Finder (lần đầu hệ điều hành hỏi quyền điều khiển Finder);
// Windows: Recycle Bin (PowerShell) – trash_<os>.go.
type SystemTrash struct{}
