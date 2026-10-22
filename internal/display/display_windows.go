package display

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Windows: mọi việc qua PowerShell (chữ truyền bằng biến môi trường, không ghép vào lệnh). Viết nhưng chưa chạy thử
// trên máy Windows thật. Loa / micro mặc định đổi bằng IPolicyConfig (COM không công bố của Windows, cách mọi công cụ
// đổi thiết bị đang dùng).

// System – Backend thật cho Windows.
type System struct{}

func ps(ctx context.Context, script string, env ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

const personalize = `HKCU:\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

func (System) Dark() (bool, error) {
	out, err := ps(context.Background(), `(Get-ItemProperty '`+personalize+`').AppsUseLightTheme`)
	if err != nil {
		return false, err
	}
	return out == "0", nil
}

func (System) SetDark(dark bool) error {
	light := "1"
	if dark {
		light = "0"
	}
	_, err := ps(context.Background(), `Set-ItemProperty '`+personalize+`' AppsUseLightTheme `+light+`; Set-ItemProperty '`+personalize+`' SystemUsesLightTheme `+light)
	return err
}

func (System) Brightness() (int, error) {
	out, err := ps(context.Background(), `(Get-CimInstance -Namespace root/WMI -ClassName WmiMonitorBrightness).CurrentBrightness`)
	if err != nil || out == "" {
		return 0, ErrNoBrightness
	}
	n, err := strconv.Atoi(strings.Fields(out)[0])
	if err != nil {
		return 0, ErrNoBrightness
	}
	return n, nil
}

func (System) SetBrightness(percent int) error {
	_, err := ps(context.Background(), `Invoke-CimMethod -Namespace root/WMI -ClassName WmiMonitorBrightnessMethods -MethodName WmiSetBrightness -Arguments @{Timeout=1;Brightness=`+strconv.Itoa(percent)+`}`)
	if err != nil {
		return ErrNoBrightness
	}
	return nil
}

const mmdevices = `HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\MMDevices\Audio\`

func (System) AudioDevices() ([]Audio, error) {
	var out []Audio
	for _, kind := range []string{"Render", "Capture"} {
		text, err := ps(context.Background(), `Get-ChildItem '`+mmdevices+kind+`' | ForEach-Object { `+
			`if ((Get-ItemProperty $_.PSPath).DeviceState -eq 1) { $p = Get-ItemProperty ($_.PSPath + '\Properties'); `+
			`$_.PSChildName + [char]9 + $p.'{a45c254e-df1c-4efd-8020-67d146a850e0},2' } }`)
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(text, "\n") {
			if f := strings.SplitN(strings.TrimSpace(line), "\t", 2); len(f) == 2 {
				out = append(out, Audio{ID: f[0], Name: f[1], Input: kind == "Capture"})
			}
		}
	}
	return out, nil
}

const policyConfig = `
Add-Type -TypeDefinition @'
using System; using System.Runtime.InteropServices;
[Guid("F8679F50-850A-41CF-9C72-430F290290C8"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
public interface IPolicyConfig {
 int GetMixFormat([MarshalAs(UnmanagedType.LPWStr)] string d, IntPtr f);
 int GetDeviceFormat([MarshalAs(UnmanagedType.LPWStr)] string d, int b, IntPtr f);
 int ResetDeviceFormat([MarshalAs(UnmanagedType.LPWStr)] string d);
 int SetDeviceFormat([MarshalAs(UnmanagedType.LPWStr)] string d, IntPtr a, IntPtr b);
 int GetProcessingPeriod([MarshalAs(UnmanagedType.LPWStr)] string d, int b, IntPtr a, IntPtr c);
 int SetProcessingPeriod([MarshalAs(UnmanagedType.LPWStr)] string d, IntPtr a);
 int GetShareMode([MarshalAs(UnmanagedType.LPWStr)] string d, IntPtr m);
 int SetShareMode([MarshalAs(UnmanagedType.LPWStr)] string d, IntPtr m);
 int GetPropertyValue([MarshalAs(UnmanagedType.LPWStr)] string d, int b, IntPtr k, IntPtr v);
 int SetPropertyValue([MarshalAs(UnmanagedType.LPWStr)] string d, int b, IntPtr k, IntPtr v);
 int SetDefaultEndpoint([MarshalAs(UnmanagedType.LPWStr)] string d, int role);
 int SetEndpointVisibility([MarshalAs(UnmanagedType.LPWStr)] string d, int v);
}
[ComImport, Guid("870AF99C-171D-4F9E-AF0D-E63DF40C2BC9")] public class PolicyConfigClient {}
'@
$c = [IPolicyConfig](New-Object PolicyConfigClient)
0..2 | ForEach-Object { [void]$c.SetDefaultEndpoint($env:GH_AUDIO_ENDPOINT, $_) }
`

func (System) SetDefaultAudio(id string, input bool) error {
	prefix := "{0.0.0.00000000}."
	if input {
		prefix = "{0.0.1.00000000}."
	}
	_, err := ps(context.Background(), policyConfig, "GH_AUDIO_ENDPOINT="+prefix+id)
	return err
}

const playScript = `
$t = $env:GH_SOUND
if (Test-Path -LiteralPath $t) { $p = New-Object System.Media.SoundPlayer $t; if ($env:GH_WAIT -eq '1') { $p.PlaySync() } else { $p.Play(); Start-Sleep -Milliseconds 200 } }
else { $s = [System.Media.SystemSounds]::$t; if (-not $s) { exit 3 }; $s.Play(); Start-Sleep -Milliseconds 800 }
`

func (System) PlaySound(ctx context.Context, target string, wait bool) error {
	w := "0"
	if wait {
		w = "1"
	}
	if filepath.Ext(target) != "" && strings.ToLower(filepath.Ext(target)) != ".wav" {
		return fmt.Errorf("Windows chỉ phát được file .wav (hoặc tên âm hệ thống: Asterisk, Beep, Exclamation, Hand, Question)")
	}
	if _, err := ps(ctx, playScript, "GH_SOUND="+target, "GH_WAIT="+w); err != nil {
		return fmt.Errorf("không phát được %q: %w", target, err)
	}
	return nil
}
