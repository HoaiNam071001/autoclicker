package notify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"autoclicker/internal/model"
)

type hit struct {
	path, ctype, body string
	form              map[string]string
	file              []byte
}

func recorder(t *testing.T, code int) (*httptest.Server, *[]hit) {
	t.Helper()
	var hits []hit
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := hit{path: r.URL.Path, ctype: r.Header.Get("Content-Type"), form: map[string]string{}}
		if strings.HasPrefix(h.ctype, "multipart/") {
			r.ParseMultipartForm(1 << 20)
			for k, v := range r.MultipartForm.Value {
				h.form[k] = v[0]
			}
			for _, fs := range r.MultipartForm.File {
				f, _ := fs[0].Open()
				h.file, _ = io.ReadAll(f)
			}
		} else {
			b, _ := io.ReadAll(r.Body)
			h.body = string(b)
		}
		hits = append(hits, h)
		w.WriteHeader(code)
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

func TestTelegram(t *testing.T) {
	srv, hits := recorder(t, 200)
	telegramAPI = srv.URL
	defer func() { telegramAPI = "https://api.telegram.org" }()
	ctx := context.Background()
	if err := Send(ctx, model.NotifyTelegram, "123:ABC@456", Message{Text: "xong"}); err != nil {
		t.Fatal(err)
	}
	if err := Send(ctx, model.NotifyTelegram, "123:ABC@456", Message{Text: "lỗi", Image: []byte("png")}); err != nil {
		t.Fatal(err)
	}
	h := *hits
	if h[0].path != "/bot123:ABC/sendMessage" || !strings.Contains(h[0].body, `"chat_id":"456"`) {
		t.Errorf("sendMessage: %+v", h[0])
	}
	if h[1].path != "/bot123:ABC/sendPhoto" || h[1].form["caption"] != "lỗi" || string(h[1].file) != "png" {
		t.Errorf("sendPhoto: %+v", h[1])
	}
}

func TestDiscordWebhook(t *testing.T) {
	srv, hits := recorder(t, 204)
	ctx := context.Background()
	if err := Send(ctx, model.NotifyDiscord, srv.URL+"/d", Message{Text: "xong", Image: []byte("png")}); err != nil {
		t.Fatal(err)
	}
	if err := Send(ctx, model.NotifyWebhook, srv.URL+"/w", Message{Text: "xong", Image: []byte("png"), Fields: map[string]any{"status": "success"}}); err != nil {
		t.Fatal(err)
	}
	h := *hits
	if !strings.Contains(h[0].form["payload_json"], `"content":"xong"`) || string(h[0].file) != "png" {
		t.Errorf("discord: %+v", h[0])
	}
	var v map[string]any
	json.Unmarshal([]byte(h[1].body), &v)
	if v["status"] != "success" || v["text"] != "xong" || !strings.HasPrefix(v["image"].(string), "data:image/png;base64,") {
		t.Errorf("webhook: %v", v)
	}
	if Send(ctx, model.NotifyNone, "", Message{}) != nil || len(*hits) != 2 {
		t.Error("none không gửi gì")
	}
}

func TestSendError(t *testing.T) {
	srv, _ := recorder(t, 400)
	if err := Send(context.Background(), model.NotifyWebhook, srv.URL, Message{Text: "x"}); err == nil || !strings.Contains(err.Error(), "400") {
		t.Fatalf("mã lỗi HTTP phải báo: %v", err)
	}
}
