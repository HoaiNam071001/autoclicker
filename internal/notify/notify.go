// Package notify – gửi thông báo sau lần chạy hẹn giờ (bước 22): Telegram (bot), Discord (webhook), webhook JSON.
// Ảnh chụp màn hình (PNG) kèm theo nếu có.
package notify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"autoclicker/internal/model"
)

// Message – nội dung thông báo. Fields thêm vào JSON của webhook (status, runId…).
type Message struct {
	Text   string
	Image  []byte // PNG, nil = không kèm
	Fields map[string]any
}

// telegramAPI – đổi trong test.
var telegramAPI = "https://api.telegram.org"

const timeout = 20 * time.Second

// Send gửi m theo kind (model.Notify*) tới target. kind none / "" → không làm gì.
func Send(ctx context.Context, kind, target string, m Message) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	switch kind {
	case "", model.NotifyNone:
		return nil
	case model.NotifyTelegram:
		return telegram(ctx, target, m)
	case model.NotifyDiscord:
		return discord(ctx, target, m)
	case model.NotifyWebhook:
		return webhook(ctx, target, m)
	}
	return fmt.Errorf("kiểu thông báo %q không hợp lệ", kind)
}

func telegram(ctx context.Context, target string, m Message) error {
	token, chat, ok := strings.Cut(target, "@")
	if !ok {
		return fmt.Errorf("Telegram cần dạng botToken@chatId")
	}
	base := telegramAPI + "/bot" + token
	if m.Image == nil {
		b, _ := json.Marshal(map[string]string{"chat_id": chat, "text": m.Text})
		return post(ctx, base+"/sendMessage", "application/json", bytes.NewReader(b))
	}
	body, ctype, err := form(map[string]string{"chat_id": chat, "caption": m.Text}, "photo", m.Image)
	if err != nil {
		return err
	}
	return post(ctx, base+"/sendPhoto", ctype, body)
}

func discord(ctx context.Context, url string, m Message) error {
	payload, _ := json.Marshal(map[string]string{"content": m.Text})
	if m.Image == nil {
		return post(ctx, url, "application/json", bytes.NewReader(payload))
	}
	body, ctype, err := form(map[string]string{"payload_json": string(payload)}, "files[0]", m.Image)
	if err != nil {
		return err
	}
	return post(ctx, url, ctype, body)
}

func webhook(ctx context.Context, url string, m Message) error {
	v := map[string]any{"text": m.Text}
	for k, x := range m.Fields {
		v[k] = x
	}
	if m.Image != nil {
		v["image"] = "data:image/png;base64," + base64.StdEncoding.EncodeToString(m.Image)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return post(ctx, url, "application/json", bytes.NewReader(b))
}

// form – multipart với các trường chữ + một file PNG.
func form(fields map[string]string, fileField string, png []byte) (io.Reader, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}
	f, err := w.CreateFormFile(fileField, "man-hinh.png")
	if err != nil {
		return nil, "", err
	}
	if _, err := f.Write(png); err != nil {
		return nil, "", err
	}
	return &buf, w.FormDataContentType(), w.Close()
}

func post(ctx context.Context, url, ctype string, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", ctype)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("gửi thông báo: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		msg, _ := io.ReadAll(io.LimitReader(res.Body, 300))
		return fmt.Errorf("gửi thông báo: %s %s", res.Status, strings.TrimSpace(string(msg)))
	}
	return nil
}
