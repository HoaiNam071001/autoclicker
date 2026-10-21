# Bộ Tesseract kèm trong exe (bước 40)

`tess.zip` (không commit – `.gitignore`) do `scripts/fetch-tesseract.sh` tạo: DLL Tesseract 5 đã strip + `vie.traineddata` (tessdata_fast). `ocr_tessbundle_windows.go` nhúng thư mục này (`go:embed`); thiếu `tess.zip` thì app vẫn chạy, đọc chữ bằng OCR Windows.

Giấy phép: Tesseract, tessdata – Apache-2.0; leptonica – BSD-2; các thư viện kèm (libarchive, zlib, libpng, libjpeg-turbo, libtiff, libwebp, openjpeg, OpenSSL libcrypto, libstdc++ / libgcc – GPL kèm ngoại lệ runtime, …) theo giấy phép của từng gói – xem `THIRD_PARTY.md` ở gốc repo.
