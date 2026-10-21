package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"ghosthand/internal/store"
)

func TestParse(t *testing.T) {
	a, err := Parse([]string{"run", "--preset", "Gửi cho Lan", "a.json", "--set", "M1=Lan", "--set=M2=a=b", "--countdown=3", "--no-wait"})
	if err != nil {
		t.Fatal(err)
	}
	if a.Path != "a.json" || a.Preset != "Gửi cho Lan" || a.Values["M1"] != "Lan" || a.Values["M2"] != "a=b" || a.Countdown != 3 || !a.NoWait {
		t.Fatalf("%+v", a)
	}
	for _, bad := range [][]string{{"run"}, {"run", "a", "b"}, {"run", "a", "--set", "x"}, {"run", "a", "--countdown", "-1"}, {"run", "a", "--lạ"}, {"run", "a", "--preset"}} {
		if _, err := Parse(bad); err == nil || errors.Is(err, errHelp) {
			t.Errorf("%v phải lỗi", bad)
		}
	}
	if _, err := Parse([]string{"help"}); !errors.Is(err, errHelp) {
		t.Error("help")
	}
	if !IsCommand([]string{"run"}) || IsCommand(nil) || IsCommand([]string{"-psn_0_123"}) {
		t.Error("IsCommand")
	}
}

// fakeApp – app giả: /health, /run (ghi lại thân), /events trả vài event SSE.
func fakeApp(t *testing.T, end string) (*httptest.Server, *map[string]any) {
	t.Helper()
	var got map[string]any
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte(`{"ok":true}`)) })
	mux.HandleFunc("POST /run", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			w.WriteHeader(401)
			return
		}
		json.NewDecoder(r.Body).Decode(&got)
		if got["preset"] == "khong-co" {
			w.WriteHeader(400)
			w.Write([]byte(`{"error":"không có bộ giá trị","code":"invalid"}`))
			return
		}
		w.Write([]byte(`{"runId":"r1"}`))
	})
	mux.HandleFunc("GET /events", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "event: status\ndata: {\"runId\":\"r1\",\"state\":\"running\"}\n\n")
		fmt.Fprint(w, "event: engine:node\ndata: {\"ref\":\"n1\",\"status\":\"success\",\"durationMs\":12}\n\n")
		fmt.Fprintf(w, "event: run:ended\ndata: %s\n\n", end)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, &got
}

func clientFor(srv *httptest.Server) *Client {
	var port int
	fmt.Sscanf(srv.URL[strings.LastIndex(srv.URL, ":")+1:], "%d", &port)
	return &Client{Info: func() (store.APIInfo, error) { return store.APIInfo{Port: port, Token: "tok"}, nil }, HTTP: srv.Client(),
		Launch: func() error { return errors.New("không được mở") }}
}

func TestRunFollow(t *testing.T) {
	srv, got := fakeApp(t, `{"runId":"r1","state":"success","durationMs":1500}`)
	var out, errOut bytes.Buffer
	code := clientFor(srv).Run(Args{Path: "/a.json", Values: map[string]any{"M1": "Lan"}}, &out, &errOut)
	if code != ExitOK || !strings.Contains(out.String(), "✓ n1 · 12 ms") || !strings.Contains(out.String(), "Xong sau 1.5 s") {
		t.Fatalf("code %d\n%s%s", code, out.String(), errOut.String())
	}
	if (*got)["source"] != "cli" || (*got)["path"] != "/a.json" {
		t.Fatalf("thân /run: %v", *got)
	}
}

func TestRunFailAndErrors(t *testing.T) {
	srv, _ := fakeApp(t, `{"runId":"r1","state":"fail","error":"n2: không thấy ảnh"}`)
	var out, errOut bytes.Buffer
	if code := clientFor(srv).Run(Args{Path: "/a.json"}, &out, &errOut); code != ExitFail || !strings.Contains(out.String(), "không thấy ảnh") {
		t.Fatalf("lỗi chạy: %d %s", code, out.String())
	}
	if code := clientFor(srv).Run(Args{Path: "/a.json", Preset: "khong-co"}, &out, &errOut); code != ExitFail || !strings.Contains(errOut.String(), "không có bộ giá trị") {
		t.Fatalf("app từ chối: %d %s", code, errOut.String())
	}
	c := &Client{Info: func() (store.APIInfo, error) { return store.APIInfo{}, os.ErrNotExist }, HTTP: http.DefaultClient,
		Launch: func() error { return nil }, WaitApp: 50 * time.Millisecond, poll: 10 * time.Millisecond}
	if code := c.Run(Args{Path: "/a.json"}, &out, &errOut); code != ExitConnect {
		t.Fatalf("không có app phải mã 3: %d", code)
	}
}
