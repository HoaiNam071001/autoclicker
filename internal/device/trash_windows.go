package device

import (
	"fmt"
	"strings"
)

// Recycle Bin qua Microsoft.VisualBasic (SendToRecycleBin) – khôi phục được. Danh sách đường dẫn qua stdin, mỗi dòng một file.
func (SystemTrash) MoveToTrash(paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	script := "Add-Type -AssemblyName Microsoft.VisualBasic;" +
		"foreach($p in [Console]::In.ReadToEnd() -split \"`n\"){ if($p){ [Microsoft.VisualBasic.FileIO.FileSystem]::DeleteFile($p.Trim(),'OnlyErrorDialogs','SendToRecycleBin') } }"
	if _, err := run(ps(script), strings.Join(paths, "\n")); err != nil {
		return fmt.Errorf("không chuyển file vào Thùng rác được: %w", err)
	}
	return nil
}

func (SystemTrash) EmptyTrash() error {
	if _, err := run(ps("Clear-RecycleBin -Force -ErrorAction Stop"), ""); err != nil {
		return fmt.Errorf("không dọn Thùng rác được: %w", err)
	}
	return nil
}
