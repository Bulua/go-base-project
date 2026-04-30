package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	WriteJSON(w, http.StatusOK, Body{
		Code:    0,
		Message: "ok",
		Data:    data,
		TraceID: TraceID(r),
	})
}

func Error(w http.ResponseWriter, r *http.Request, status int, code int, message string) {
	WriteJSON(w, status, Body{
		Code:    code,
		Message: message,
		TraceID: TraceID(r),
	})
}

func WriteJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func ReadJSON(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func AllowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	WriteJSON(w, http.StatusMethodNotAllowed, Body{
		Code:    405,
		Message: "method not allowed",
		TraceID: TraceID(r),
	})
	return false
}

func TraceID(r *http.Request) string {
	if value := strings.TrimSpace(r.Header.Get("X-Trace-ID")); value != "" {
		return value
	}
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
