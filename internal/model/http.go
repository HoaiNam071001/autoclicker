package model

import "strings"

// Node HTTP (bước 35 N3, engine/http.go): gọi API – method, URL, query, header, body; mọi ô chữ nhận {{biến}}.

// KV – một dòng của bảng query / header / form.
type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var HTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// Kiểu body (Step.BodyType).
const (
	BodyNone = ""     // không gửi body
	BodyJSON = "json" // Body là JSON (Content-Type application/json); biến trong Body thoát theo chuỗi JSON
	BodyForm = "form" // Form – application/x-www-form-urlencoded
	BodyText = "text" // Body gửi nguyên văn (Content-Type text/plain nếu header không đặt)
)

var BodyTypes = []string{BodyNone, BodyJSON, BodyForm, BodyText}

// secretWords – tên header / tham số có các chữ này coi là bí mật (token, mật khẩu…): che trong nhật ký chạy.
var secretWords = []string{"auth", "token", "secret", "pass", "apikey", "api-key", "api_key", "cookie", "session", "signature", "credential"}

// IsSecretKey – tên header / tham số chứa bí mật. Cùng công thức isSecretKey (frontend/src/flow/http.ts).
func IsSecretKey(name string) bool {
	n := strings.ToLower(name)
	for _, w := range secretWords {
		if strings.Contains(n, w) {
			return true
		}
	}
	return false
}
