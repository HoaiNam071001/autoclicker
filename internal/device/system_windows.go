package device

import (
	"fmt"
	"strconv"
	"strings"
)

// Core Audio (IAudioEndpointVolume của loa mặc định) qua C# trong PowerShell – Windows không có lệnh sẵn để chỉnh âm lượng.
const audioCS = `Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
[Guid("5CDF2C82-841E-4546-9722-0CF74078229A"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IAudioEndpointVolume {
  int f(); int g(); int h(); int i();
  int SetMasterVolumeLevelScalar(float fLevel, Guid ctx);
  int j();
  int GetMasterVolumeLevelScalar(out float pfLevel);
  int k(); int l(); int m(); int n();
  int SetMute([MarshalAs(UnmanagedType.Bool)] bool bMute, Guid ctx);
  int GetMute(out bool pbMute);
}
[Guid("D666063F-1587-4E43-81F1-B948E807363F"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IMMDevice { int Activate(ref Guid id, int clsCtx, IntPtr p, out IAudioEndpointVolume aev); }
[Guid("A95664D2-9614-4F35-A746-DE8DB63617E6"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IMMDeviceEnumerator { int f(); int GetDefaultAudioEndpoint(int dataFlow, int role, out IMMDevice ep); }
[ComImport, Guid("BCDE0395-E52F-467C-8E3D-C4579291692E")] class MMDeviceEnumerator { }
public static class AcAudio {
  static IAudioEndpointVolume Vol() {
    var en = (IMMDeviceEnumerator)(new MMDeviceEnumerator());
    IMMDevice dev; Marshal.ThrowExceptionForHR(en.GetDefaultAudioEndpoint(0, 1, out dev));
    IAudioEndpointVolume v; var id = typeof(IAudioEndpointVolume).GUID;
    Marshal.ThrowExceptionForHR(dev.Activate(ref id, 23, IntPtr.Zero, out v));
    return v;
  }
  public static string Get() { float l; bool m; var v = Vol(); Marshal.ThrowExceptionForHR(v.GetMasterVolumeLevelScalar(out l)); Marshal.ThrowExceptionForHR(v.GetMute(out m)); return ((int)Math.Round(l * 100)) + "," + (m ? "true" : "false"); }
  public static void Set(int level) { var v = Vol(); Marshal.ThrowExceptionForHR(v.SetMasterVolumeLevelScalar(level / 100f, Guid.Empty)); Marshal.ThrowExceptionForHR(v.SetMute(false, Guid.Empty)); }
  public static void Mute(bool m) { Marshal.ThrowExceptionForHR(Vol().SetMute(m, Guid.Empty)); }
}
'@;`

func (SystemAudio) Volume() (int, bool, error) {
	out, err := run(ps(audioCS+"[AcAudio]::Get()"), "")
	if err != nil {
		return 0, false, fmt.Errorf("không đọc được âm lượng: %w", err)
	}
	v, m, _ := strings.Cut(strings.TrimSpace(out), ",")
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false, fmt.Errorf("không đọc được âm lượng (%q)", strings.TrimSpace(out))
	}
	return n, m == "true", nil
}

func (SystemAudio) SetVolume(level int) error {
	if _, err := run(ps(fmt.Sprintf("%s[AcAudio]::Set(%d)", audioCS, level)), ""); err != nil {
		return fmt.Errorf("không đặt được âm lượng: %w", err)
	}
	return nil
}

func (SystemAudio) SetMuted(muted bool) error {
	if _, err := run(ps(fmt.Sprintf("%s[AcAudio]::Mute($%v)", audioCS, muted)), ""); err != nil {
		return fmt.Errorf("không đổi được tắt tiếng: %w", err)
	}
	return nil
}
