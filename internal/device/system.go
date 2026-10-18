package device

// Audio – âm lượng loa mặc định của máy (node Hệ thống, bước 35 N5). Mức 0–100. Qua interface để test bằng bản giả
// (devicetest.FakeAudio) – test không đổi âm lượng thật.
type Audio interface {
	Volume() (level int, muted bool, err error)
	SetVolume(level int) error
	SetMuted(muted bool) error
}

// SystemAudio – macOS qua osascript (volume settings – không cần quyền), Windows qua Core Audio (PowerShell) – system_<os>.go.
type SystemAudio struct{}
