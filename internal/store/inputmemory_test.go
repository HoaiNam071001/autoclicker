package store

import "testing"

func TestInputMemoryRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)            // macOS / Linux: UserConfigDir theo HOME
	t.Setenv("XDG_CONFIG_HOME", dir) // Linux
	t.Setenv("AppData", dir)         // Windows
	m := InputMemory{Last: map[string]any{"M1": "Lan"}, Presets: []InputPreset{{Name: "Gửi cho Lan", Values: map[string]any{"M1": "Lan"}}}, AutoRun: true}
	if err := SaveInputMemory("/k.json", m); err != nil {
		t.Fatal(err)
	}
	got := InputMemoryFor("/k.json")
	if got.Last["M1"] != "Lan" || len(got.Presets) != 1 || !got.AutoRun {
		t.Fatalf("%+v", got)
	}
	if other := InputMemoryFor("/khac.json"); other.Last != nil {
		t.Fatalf("kịch bản khác không được dùng chung: %+v", other)
	}
}
