package store

import "golang.org/x/sys/windows"

// knownFolders – Documents / Pictures thật qua SHGetKnownFolderPath: máy có OneDrive chuyển thư mục vẫn ra đúng chỗ
// (ghép home + "Documents" sẽ trỏ vào thư mục cũ đã trống).
func knownFolders(home string) (docs, pics string) {
	docs, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		docs = home + `\Documents`
	}
	pics, err = windows.KnownFolderPath(windows.FOLDERID_Pictures, 0)
	if err != nil {
		pics = home + `\Pictures`
	}
	return docs, pics
}
