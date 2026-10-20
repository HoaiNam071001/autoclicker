package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"autoclicker/internal/session"
)

const sseKeepAlive = 15 * time.Second

// events – GET /events (SSE). Gửi trước trạng thái hiện tại (event "status"), rồi mọi event của app. Có runId:
// đóng sau run:ended của phiên đó (phiên đã xong từ trước thì đóng ngay sau "status").
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, 500, "internal", "không stream được")
		return
	}
	id := r.URL.Query().Get("runId")
	ch, stop := s.hub.Subscribe() // nghe trước khi đọc trạng thái – không lỡ run:ended xảy ra giữa hai việc
	defer stop()
	st, err := s.back.Status(id)
	if err != nil && id != "" {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(200)
	if err == nil {
		sendSSE(w, "status", st)
	}
	fl.Flush()
	if id != "" && st.Done() {
		return
	}
	tick := time.NewTicker(sseKeepAlive)
	defer tick.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-tick.C:
			fmt.Fprint(w, ": ping\n\n")
			fl.Flush()
		case ev := <-ch:
			sendSSE(w, ev.Name, ev.Data)
			fl.Flush()
			if end, ok := ev.Data.(session.Status); ok && ev.Name == "run:ended" && id != "" && end.RunID == id {
				return
			}
		}
	}
}

func sendSSE(w http.ResponseWriter, name string, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", name, b)
}
