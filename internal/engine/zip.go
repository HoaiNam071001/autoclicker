package engine

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"ghosthand/internal/model"
)

// Node Nén / giải nén (bước 43 S8): .zip chuẩn của thư viện Go, không gọi lệnh ngoài nên chạy giống nhau ở macOS +
// Windows. Giải nén chặn mọi đường dẫn thoát ra ngoài thư mục đích (zip slip) và liên kết tượng trưng.

// zipMaxBytes – tổng dung lượng giải nén tối đa (chặn zip bomb).
const zipMaxBytes = 8 << 30

func (e *Engine) zipStep(ctx context.Context, s model.Step) (map[string]any, error) {
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	src, err := userPath(s.FilePath)
	if err != nil {
		return nil, err
	}
	if s.ZipOp == model.ZipExtract {
		return unzipTo(ctx, src, s.ZipTo)
	}
	return zipFrom(ctx, src, strings.TrimSpace(s.FilePattern), s.ZipTo)
}

// zipFrom nén src (file, hoặc thư mục – giữ tên thư mục làm gốc trong zip) ra dest; pattern lọc tên file khi nén thư mục.
func zipFrom(ctx context.Context, src, pattern, dest string) (map[string]any, error) {
	info, err := os.Stat(src)
	if err != nil {
		return nil, &notFoundErr{msg: fmt.Sprintf("không thấy %s", src)}
	}
	if dest = strings.TrimSpace(dest); dest == "" {
		dest = src + ".zip"
	} else if dest, err = userPath(dest); err != nil {
		return nil, err
	}
	if err := guardPath(dest); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return nil, fmt.Errorf("không tạo được thư mục chứa file nén: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".tmp-*.zip")
	if err != nil {
		return nil, fmt.Errorf("không tạo được file nén: %w", err)
	}
	defer os.Remove(tmp.Name()) // đã đổi tên thì thành no-op
	zw := zip.NewWriter(tmp)
	var count, size int
	add := func(path, name string, fi fs.FileInfo) error {
		if path == dest {
			return nil
		}
		n, err := addZipFile(ctx, zw, path, name, fi)
		count, size = count+1, size+int(n)
		return err
	}
	if info.IsDir() {
		root := filepath.Dir(filepath.Clean(src))
		err = filepath.WalkDir(src, func(path string, d fs.DirEntry, werr error) error {
			if werr != nil || d.IsDir() || !d.Type().IsRegular() {
				return werr
			}
			if ok, _ := filepath.Match(pattern, d.Name()); pattern != "" && !ok {
				return nil
			}
			fi, err := d.Info()
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, path)
			return add(path, filepath.ToSlash(rel), fi)
		})
	} else {
		err = add(src, info.Name(), info)
	}
	if err == nil && count == 0 {
		err = &notFoundErr{msg: fmt.Sprintf("không có file nào để nén trong %s", src)}
	}
	if err == nil {
		err = zw.Close()
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(tmp.Name(), dest)
	}
	if err != nil {
		return nil, err
	}
	return map[string]any{"path": dest, "count": count, "size": size}, nil
}

func addZipFile(ctx context.Context, zw *zip.Writer, path, name string, fi fs.FileInfo) (int64, error) {
	hdr, err := zip.FileInfoHeader(fi)
	if err != nil {
		return 0, err
	}
	hdr.Name, hdr.Method = name, zip.Deflate
	w, err := zw.CreateHeader(hdr)
	if err != nil {
		return 0, err
	}
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("không đọc được %s: %w", path, err)
	}
	defer f.Close()
	return copyCtx(ctx, w, f, -1)
}

// unzipTo giải nén zipPath vào dest (trống = thư mục cạnh file zip, tên bỏ đuôi .zip).
func unzipTo(ctx context.Context, zipPath, dest string) (map[string]any, error) {
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &notFoundErr{msg: fmt.Sprintf("không thấy %s", zipPath)}
		}
		return nil, &invalidErr{fmt.Sprintf("không mở được %s như file zip: %v", filepath.Base(zipPath), err)}
	}
	defer zr.Close()
	if dest = strings.TrimSpace(dest); dest == "" {
		dest = strings.TrimSuffix(zipPath, filepath.Ext(zipPath))
	} else if dest, err = userPath(dest); err != nil {
		return nil, err
	}
	if err := guardPath(dest); err != nil {
		return nil, err
	}
	if fi, err := os.Stat(dest); err == nil && !fi.IsDir() {
		return nil, &invalidErr{fmt.Sprintf("%s là một file, không phải thư mục – đặt “Giải nén vào” sang thư mục khác", dest)}
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return nil, fmt.Errorf("không tạo được thư mục đích: %w", err)
	}
	root := realPath(dest)
	budget := int64(zipMaxBytes)
	var count int
	var size int64
	for _, f := range zr.File {
		target, err := zipTarget(root, f)
		if err != nil {
			return nil, err
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return nil, fmt.Errorf("không tạo được thư mục %s: %w", f.Name, err)
			}
			continue
		}
		n, err := extractZipFile(ctx, f, target, budget)
		if err != nil {
			return nil, err
		}
		budget -= n
		count, size = count+1, size+n
	}
	return map[string]any{"path": dest, "count": count, "size": size}, nil
}

// zipTarget – đường dẫn thật của một mục zip; mục thoát khỏi root, tuyệt đối hoặc là liên kết thì bị từ chối.
func zipTarget(root string, f *zip.File) (string, error) {
	bad := func(why string) error {
		return &invalidErr{fmt.Sprintf("file zip không an toàn – mục %q %s", f.Name, why)}
	}
	if f.Mode()&fs.ModeSymlink != 0 {
		return "", bad("là liên kết tượng trưng")
	}
	name := strings.ReplaceAll(f.Name, `\`, "/") // zip tạo trên Windows đôi khi dùng \
	if strings.HasPrefix(name, "/") || filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return "", bad("có đường dẫn tuyệt đối")
	}
	target := filepath.Join(root, filepath.FromSlash(name))
	if target == root || !within(realPath(target), root) {
		return "", bad("thoát ra ngoài thư mục đích")
	}
	return target, nil
}

func extractZipFile(ctx context.Context, f *zip.File, target string, budget int64) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0, fmt.Errorf("không tạo được thư mục chứa %s: %w", f.Name, err)
	}
	rc, err := f.Open()
	if err != nil {
		return 0, fmt.Errorf("không đọc được %s trong file zip: %w", f.Name, err)
	}
	defer rc.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode().Perm()|0o600)
	if err != nil {
		return 0, fmt.Errorf("không ghi được %s: %w", target, err)
	}
	n, err := copyCtx(ctx, out, rc, budget)
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	return n, err
}

// copyCtx chép từng khối để Dừng có tác dụng giữa file lớn; limit ≥ 0: vượt thì báo lỗi (zip bomb).
func copyCtx(ctx context.Context, dst io.Writer, src io.Reader, limit int64) (int64, error) {
	buf := make([]byte, 256<<10)
	var total int64
	for {
		if err := ctx.Err(); err != nil {
			return total, err
		}
		n, rerr := src.Read(buf)
		if n > 0 {
			if total += int64(n); limit >= 0 && total > limit {
				return total, &invalidErr{"giải nén vượt 8 GB – dừng để tránh đầy ổ đĩa"}
			}
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return total, fmt.Errorf("không ghi được dữ liệu: %w", werr)
			}
		}
		if rerr == io.EOF {
			return total, nil
		}
		if rerr != nil {
			return total, rerr
		}
	}
}
