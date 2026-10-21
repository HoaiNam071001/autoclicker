#!/usr/bin/env bash
# Tạo internal/device/tess/tess.zip – bộ Tesseract kèm trong exe Windows (bước 40 O2): DLL đã strip + vie.traineddata.
# Chạy một lần trước `wails build` / `go test` trên Windows (máy dev Windows: Git Bash + 7-Zip; CI Linux: 7z + MinGW).
# Nguồn cố định phiên bản + sha256; đã có zip thì bỏ qua (xoá zip để tạo lại).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/internal/device/tess/tess.zip"
[ -f "$OUT" ] && { echo "đã có $OUT"; exit 0; }

SETUP_URL="https://github.com/UB-Mannheim/tesseract/releases/download/v5.4.0.20240606/tesseract-ocr-w64-setup-5.4.0.20240606.exe"
SETUP_SHA="c885fff6998e0608ba4bb8ab51436e1c6775c2bafc2559a19b423e18678b60c9"
VIE_URL="https://github.com/tesseract-ocr/tessdata_fast/raw/923915d4ced2a7235221788285785a29c4a42d4a/vie.traineddata"
VIE_SHA="79df64caf7bcfb2a27df5042ecb6121e196eada34da774956995747636d5bfa1"
# libtesseract-5.dll và các DLL nó cần (đệ quy, trừ DLL hệ thống) – dò bằng objdump -p ở 40 O0
DLLS="libtesseract-5.dll libleptonica-6.dll libarchive-13.dll libb2-1.dll libbz2-1.dll libcrypto-3-x64.dll
libexpat-1.dll libiconv-2.dll liblz4.dll liblzma-5.dll zlib1.dll libzstd.dll libgcc_s_seh-1.dll libstdc++-6.dll
libwinpthread-1.dll libgif-7.dll libjpeg-8.dll libopenjp2-7.dll libpng16-16.dll libtiff-6.dll libwebp-7.dll
libwebpmux-3.dll libsharpyuv-0.dll libdeflate.dll libjbig-0.dll libLerc.dll"

SEVENZ="$(command -v 7z || command -v 7za || echo "/c/Program Files/7-Zip/7z.exe")"
STRIP="$(command -v x86_64-w64-mingw32-strip || command -v strip)"
[ -x "$SEVENZ" ] || { echo "thiếu 7-Zip (Windows: cài 7-Zip; Linux: apt install p7zip-full)"; exit 1; }

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT
fetch() { # url sha file [file có sẵn]
  if [ -n "${4:-}" ]; then cp "$4" "$3"; else curl -fsSL --retry 3 -o "$3" "$1"; fi
  echo "$2  $3" | sha256sum -c --quiet - || { echo "sai sha256: $1"; exit 1; }
}
# TESS_SETUP / TESS_VIE: dùng file đã tải sẵn (mạng chậm) – vẫn kiểm tra sha256
fetch "$SETUP_URL" "$SETUP_SHA" "$WORK/setup.exe" "${TESS_SETUP:-}"
fetch "$VIE_URL" "$VIE_SHA" "$WORK/vie.traineddata" "${TESS_VIE:-}"

mkdir -p "$WORK/pkg"
# chỉ giải nén DLL cần (bộ cài NSIS còn trình gỡ cài đặt 7-Zip đọc lỗi – không cần)
(cd "$WORK/pkg" && "$SEVENZ" e -y "../setup.exe" $DLLS >/dev/null)
for d in $DLLS; do
  [ -f "$WORK/pkg/$d" ] || { echo "bộ cài thiếu $d"; exit 1; }
  "$STRIP" "$WORK/pkg/$d"
done
cp "$WORK/vie.traineddata" "$WORK/pkg/"
rm -f "$OUT"
(cd "$WORK/pkg" && "$SEVENZ" a -tzip -mx=9 "$WORK/tess.zip" ./* >/dev/null)
mv "$WORK/tess.zip" "$OUT"
echo "xong: $OUT ($(du -h "$OUT" | cut -f1))"
