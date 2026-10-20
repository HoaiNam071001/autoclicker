package api

import "sync"

// Event – một event của app gửi ra SSE (cùng tên / dữ liệu với event Wails).
type Event struct {
	Name string
	Data any
}

// Hub phát event tới các kết nối SSE đang mở. Người nghe chậm thì bỏ event (không chặn engine).
type Hub struct {
	mu   sync.Mutex
	subs map[chan Event]struct{}
}

const subBuffer = 256

func NewHub() *Hub { return &Hub{subs: map[chan Event]struct{}{}} }

// Publish gửi event cho mọi người nghe; nil Hub thì bỏ qua.
func (h *Hub) Publish(name string, data any) {
	if h == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.subs {
		select {
		case c <- Event{name, data}:
		default:
		}
	}
}

// Subscribe – kênh nhận event; gọi hàm trả về để thôi nghe.
func (h *Hub) Subscribe() (<-chan Event, func()) {
	c := make(chan Event, subBuffer)
	h.mu.Lock()
	h.subs[c] = struct{}{}
	h.mu.Unlock()
	return c, func() {
		h.mu.Lock()
		delete(h.subs, c)
		h.mu.Unlock()
	}
}
