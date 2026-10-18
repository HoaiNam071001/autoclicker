package engine

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"autoclicker/internal/device/devicetest"
	"autoclicker/internal/model"
	"autoclicker/internal/model/modeltest"
	"autoclicker/internal/script/scripttest"
)

func httpStep(f func(*model.Step)) model.Step { return modeltest.Step(model.StepHTTP, f) }

// runHTTP chạy kịch bản; trả kết quả từng node theo mã.
func runHTTP(t *testing.T, sc model.Script, clip *devicetest.FakeClipboard) (map[string]model.StepResult, error) {
	t.Helper()
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(sc), nil)
	e.clip = clip
	res := map[string]model.StepResult{}
	e.onNode = func(r model.StepResult) { res[r.Ref] = r }
	return res, e.Run(context.Background())
}

func TestHTTPJSONAndDeepVars(t *testing.T) {
	var got struct {
		method, ctype, auth, q string
		body              map[string]any
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.method, got.ctype, got.auth, got.q = r.Method, r.Header.Get("Content-Type"), r.Header.Get("Authorization"), r.URL.Query().Get("tim")
		b, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(b, &got.body); err != nil {
			t.Errorf("server nhận JSON hỏng: %s", b)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Set-Cookie", "sid=bimat")
		io.WriteString(w, `{"data":{"items":[{"id":7},{"id":42,"ten":"Hà"}]},"count":2}`)
	}))
	defer srv.Close()
	clip := &devicetest.FakeClipboard{Text: `Anh "Tâm"` + "\nxuống dòng"}
	sc := model.Script{Loops: 1, Steps: []model.Step{
		modeltest.Step(model.StepClipboard, func(s *model.Step) { s.Ref = "n1" }),
		httpStep(func(s *model.Step) {
			s.Ref, s.Method, s.URL = "n2", "POST", srv.URL+"/api"
			s.Query = []model.KV{{Key: "tim", Value: "cà phê & trà"}}
			s.Headers = []model.KV{{Key: "Authorization", Value: "Bearer abc"}}
			s.BodyType, s.Body = model.BodyJSON, `{"ten": "{{n1.output.text}}", "so": 3}`
		}),
		modeltest.Step(model.StepClipboard, func(s *model.Step) {
			s.Ref, s.ClipOp, s.Text = "n3", model.ClipWrite, "{{n2.output.json.data.items.1.id}} / {{n2.output.json.count}} / {{n2.output.json.data.items.1}}"
		}),
	}}
	res, err := runHTTP(t, sc, clip)
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.ctype != "application/json" || got.auth != "Bearer abc" || got.q != "cà phê & trà" {
		t.Errorf("request: %+v", got)
	}
	if got.body["ten"] != `Anh "Tâm"`+"\nxuống dòng" || got.body["so"] != float64(3) {
		t.Errorf("body JSON: %v – biến phải được thoát trong chuỗi", got.body)
	}
	if want := `42 / 2 / {"id":42,"ten":"Hà"}`; clip.Text != want {
		t.Errorf("đọc sâu: %q, muốn %q", clip.Text, want)
	}
	out := res["n2"].Output
	if out["status"] != 200 || out["headers"].(map[string]any)["set-cookie"] != "***" {
		t.Errorf("output %v", out)
	}
}

func TestHTTPErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/slow":
			time.Sleep(300 * time.Millisecond)
		case "/form":
			r.ParseForm()
			io.WriteString(w, r.PostForm.Get("ten"))
		default:
			w.WriteHeader(404)
			io.WriteString(w, "không có")
		}
	}))
	defer srv.Close()
	one := func(f func(*model.Step)) model.StepResult {
		res, _ := runHTTP(t, model.Script{Loops: 1, Steps: []model.Step{httpStep(func(s *model.Step) { s.Ref = "n1"; f(s) })}}, &devicetest.FakeClipboard{})
		return res["n1"]
	}
	r := one(func(s *model.Step) { s.URL, s.Query = srv.URL+"/x?token=bimat", []model.KV{{Key: "api_key", Value: "k123"}} })
	if r.Status != model.StatusFail || r.Error.Code != model.ErrHTTPStatus || r.Output["status"] != 404 || r.Output["body"] != "không có" {
		t.Errorf("404: %+v", r)
	}
	if strings.Contains(r.Error.Message, "bimat") || strings.Contains(r.Error.Message, "k123") || !strings.Contains(r.Error.Message, "404") {
		t.Errorf("lời báo lộ bí mật / thiếu mã: %s", r.Error.Message)
	}
	r = one(func(s *model.Step) { s.URL, s.TimeoutMs = srv.URL+"/slow", 50 })
	if r.Error == nil || r.Error.Code != model.ErrNetwork || !strings.Contains(r.Error.Message, "hết") {
		t.Errorf("hết giờ: %+v", r.Error)
	}
	r = one(func(s *model.Step) { s.URL = "http://127.0.0.1:1/?password=xx" })
	if r.Error == nil || r.Error.Code != model.ErrNetwork || strings.Contains(r.Error.Message, "xx") {
		t.Errorf("không kết nối được: %+v", r.Error)
	}
	r = one(func(s *model.Step) { s.Method, s.URL, s.BodyType, s.Form = "POST", srv.URL+"/form", model.BodyForm, []model.KV{{Key: "ten", Value: "Lan & Hà"}} })
	if r.Status != model.StatusSuccess || r.Output["body"] != "Lan & Hà" {
		t.Errorf("form: %+v", r)
	}
	// URL sai dạng: kiểm tra chặn trước khi chạy
	if _, err := runHTTP(t, model.Script{Loops: 1, Steps: []model.Step{httpStep(func(s *model.Step) { s.URL = "ftp://x" })}}, &devicetest.FakeClipboard{}); err == nil || !strings.Contains(err.Error(), "sai dạng") {
		t.Errorf("URL sai: %v", err)
	}
}

func TestHTTPCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer srv.Close()
	e := newEngine(&devicetest.FakeInput{}, scripttest.Linear(model.Script{Loops: 1, Steps: []model.Step{httpStep(func(s *model.Step) { s.URL = srv.URL })}}), nil)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	e.Run(ctx)
	if d := time.Since(start); d > time.Second {
		t.Errorf("Dừng giữa lúc chờ phản hồi mất %v", d)
	}
}
