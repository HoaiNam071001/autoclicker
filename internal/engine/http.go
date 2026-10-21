package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ghosthand/internal/model"
)

// Node HTTP (bước 35 N3): gọi API, 2xx → success; mã khác → fail http_status (output vẫn có – nhánh Thất bại đọc
// {{n2.output.status}} / body); không gửi được → fail network. Lời báo lỗi không lộ header / tham số bí mật.

const (
	httpReadMax   = 5 << 20  // đọc tối đa 5 MB phản hồi
	httpOutputMax = 64 << 10 // output body tối đa 64 KB – nhật ký chạy lưu output
	httpTimeoutMs = 15000    // TimeoutMs = 0
)

type httpStatusErr struct{ msg string }

func (e *httpStatusErr) Error() string { return e.msg }

type networkErr struct{ msg string }

func (e *networkErr) Error() string { return e.msg }

func (e *Engine) httpCall(ctx context.Context, s model.Step) (map[string]any, error) {
	req, err := httpRequest(ctx, s)
	if err != nil {
		return nil, err
	}
	timeout := s.TimeoutMs
	if timeout <= 0 {
		timeout = httpTimeoutMs
	}
	cctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
	defer cancel()
	if err := e.gate.wait(ctx); err != nil {
		return nil, err
	}
	resp, err := e.httpc.Do(req.WithContext(cctx))
	where := maskURL(req.URL)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err() // người dùng dừng
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &networkErr{fmt.Sprintf("%s %s: hết %s chưa có phản hồi", req.Method, where, fmtSeconds(timeout))}
		}
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err // url.Error kèm URL đầy đủ (có thể có token) – chỉ lấy lý do
		}
		return nil, &networkErr{fmt.Sprintf("%s %s: không gọi được – %v", req.Method, where, err)}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, httpReadMax))
	if err != nil {
		return nil, &networkErr{fmt.Sprintf("%s %s: đứt kết nối khi đọc phản hồi – %v", req.Method, where, err)}
	}
	out := httpOutput(resp, body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		snippet := strings.TrimSpace(string(body))
		if r := []rune(snippet); len(r) > 200 {
			snippet = string(r[:200]) + "…"
		}
		msg := fmt.Sprintf("%s %s: HTTP %s", req.Method, where, resp.Status)
		if snippet != "" {
			msg += " – " + snippet
		}
		return out, &httpStatusErr{msg}
	}
	return out, nil
}

// httpRequest dựng request từ node (biến đã thay). Lỗi dữ liệu → invalid_input.
func httpRequest(ctx context.Context, s model.Step) (*http.Request, error) {
	method := s.Method
	if method == "" {
		method = http.MethodGet
	}
	u, err := url.Parse(strings.TrimSpace(s.URL))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, &invalidErr{fmt.Sprintf("URL “%s” sai dạng – cần http:// hoặc https://", strings.TrimSpace(s.URL))}
	}
	if len(s.Query) > 0 {
		q := u.Query()
		for _, kv := range s.Query {
			if k := strings.TrimSpace(kv.Key); k != "" {
				q.Add(k, kv.Value)
			}
		}
		u.RawQuery = q.Encode()
	}
	var body io.Reader
	ctype := ""
	switch s.BodyType {
	case model.BodyJSON:
		if strings.TrimSpace(s.Body) != "" && !json.Valid([]byte(s.Body)) {
			return nil, &invalidErr{"thân JSON sai cú pháp sau khi thay biến – biến chữ đặt trong dấu nháy \"{{…}}\""}
		}
		body, ctype = strings.NewReader(s.Body), "application/json"
	case model.BodyText:
		body, ctype = strings.NewReader(s.Body), "text/plain; charset=utf-8"
	case model.BodyForm:
		f := url.Values{}
		for _, kv := range s.Form {
			if k := strings.TrimSpace(kv.Key); k != "" {
				f.Add(k, kv.Value)
			}
		}
		body, ctype = strings.NewReader(f.Encode()), "application/x-www-form-urlencoded"
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, &invalidErr{fmt.Sprintf("không tạo được request: %v", err)}
	}
	req.Header.Set("User-Agent", "Ghosthand")
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	for _, kv := range s.Headers {
		if k := strings.TrimSpace(kv.Key); k != "" {
			req.Header.Set(k, kv.Value) // header người dùng đặt thắng mặc định (vd Content-Type)
		}
	}
	return req, nil
}

// httpOutput – status, body (cắt ở httpOutputMax), headers (tên chữ thường), json khi phản hồi là JSON.
func httpOutput(resp *http.Response, body []byte) map[string]any {
	headers := map[string]any{}
	for k, vs := range resp.Header {
		k = strings.ToLower(k)
		v := strings.Join(vs, ", ")
		if k == "set-cookie" {
			v = "***"
		}
		headers[k] = v
	}
	text := string(body)
	if len(body) > httpOutputMax {
		text = strings.ToValidUTF8(string(body[:httpOutputMax]), "") + "…"
	}
	out := map[string]any{"status": resp.StatusCode, "body": text, "headers": headers}
	if t := bytes.TrimSpace(body); len(t) > 0 && (t[0] == '{' || t[0] == '[' || strings.Contains(resp.Header.Get("Content-Type"), "json")) {
		var v any
		if json.Unmarshal(t, &v) == nil {
			out["json"] = normalizeJSON(v)
		}
	}
	return out
}

// normalizeJSON – số nguyên trong JSON (float64 khi giải mã) thành int để biến in ra "42" thay vì "42.0" / ô số đọc được.
func normalizeJSON(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			x[k] = normalizeJSON(e)
		}
	case []any:
		for i, e := range x {
			x[i] = normalizeJSON(e)
		}
	case float64:
		if x == float64(int(x)) && x < 1e15 && x > -1e15 {
			return int(x)
		}
	}
	return v
}

// maskURL – URL cho lời báo lỗi: tham số tên bí mật (token, key…) và mật khẩu trong URL thành ***.
func maskURL(u *url.URL) string {
	c := *u
	if c.User != nil {
		c.User = url.User(c.User.Username())
	}
	q := c.Query()
	for k := range q {
		if model.IsSecretKey(k) {
			q.Set(k, "***")
		}
	}
	c.RawQuery = q.Encode()
	return strings.ReplaceAll(c.String(), "%2A%2A%2A", "***")
}
