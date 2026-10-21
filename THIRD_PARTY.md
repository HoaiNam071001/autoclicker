# Phần mềm bên thứ ba kèm trong bản phát hành

## Bản Windows – đọc chữ tiếng Việt (bước 40)

`ghosthand.exe` nhúng bộ Tesseract (`internal/device/tess/tess.zip`, tạo bởi `scripts/fetch-tesseract.sh`), giải nén vào thư mục cache khi đọc chữ lần đầu. DLL lấy từ bản build Windows của UB Mannheim (v5.4.0.20240606), chỉ strip ký hiệu gỡ lỗi, không sửa.

| Thành phần | Giấy phép | Nguồn |
|---|---|---|
| Tesseract OCR 5.4 (`libtesseract-5.dll`) | Apache-2.0 | github.com/tesseract-ocr/tesseract |
| Mô hình `vie.traineddata` (tessdata_fast) | Apache-2.0 | github.com/tesseract-ocr/tessdata_fast |
| Leptonica (`libleptonica-6.dll`) | BSD-2-Clause | github.com/DanBloomberg/leptonica |
| libarchive | BSD-2-Clause | libarchive.org |
| OpenSSL libcrypto 3 | Apache-2.0 | openssl.org |
| zlib, libpng | zlib / libpng | zlib.net, libpng.org |
| libjpeg-turbo | IJG + BSD-3-Clause | libjpeg-turbo.org |
| libtiff | libtiff (BSD-style) | libtiff.gitlab.io |
| libwebp, libwebpmux, libsharpyuv | BSD-3-Clause | chromium.googlesource.com/webm/libwebp |
| OpenJPEG | BSD-2-Clause | openjpeg.org |
| giflib | MIT | giflib.sourceforge.net |
| LERC | Apache-2.0 | github.com/Esri/lerc |
| JBIG-KIT | GPL-2.0-or-later | cl.cam.ac.uk/~mgk25/jbigkit |
| libdeflate | MIT | github.com/ebiggers/libdeflate |
| zstd, lz4, BLAKE2 (`libb2`) | BSD-3 / BSD-2 / CC0 | facebook/zstd, lz4/lz4, BLAKE2/libb2 |
| xz (`liblzma`), bzip2 | 0BSD / bzip2 | tukaani.org/xz, sourceware.org/bzip2 |
| expat | MIT | libexpat.github.io |
| GNU libiconv | LGPL-2.1-or-later | gnu.org/software/libiconv |
| libgcc, libstdc++ (GCC runtime) | GPL-3.0 kèm GCC Runtime Library Exception | gcc.gnu.org |
| winpthreads (MinGW-w64) | MIT / BSD | mingw-w64.org |
