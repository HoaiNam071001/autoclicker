package model

// Node Nén / giải nén (bước 43 S8, engine/zip.go).

// Thao tác của node Nén / giải nén (Step.ZipOp).
const (
	ZipCompress = ""        // nén FilePath (file / thư mục) ra file .zip
	ZipExtract  = "extract" // giải nén file .zip vào thư mục
)

var ZipOps = []string{ZipCompress, ZipExtract}
