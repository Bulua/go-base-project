package apperror_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gobaseproject/server/internal/apperror"
	usermodel "gobaseproject/server/internal/model/user"
)

func TestWriteUsesRegisteredBusinessErrorMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", nil)
	req.Header.Set("X-Trace-ID", "trace-001")
	rec := httptest.NewRecorder()

	apperror.Write(rec, req, usermodel.ErrLoginNameTaken)

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		TraceID string `json:"trace_id"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected http status 409, got %d", rec.Code)
	}
	if body.Code != 409201 {
		t.Fatalf("expected code 409201, got %d", body.Code)
	}
	if body.Message != "登录账号已存在" {
		t.Fatalf("expected Chinese error message, got %q", body.Message)
	}
	if body.TraceID != "trace-001" {
		t.Fatalf("expected trace id trace-001, got %q", body.TraceID)
	}
}

func TestWriteHidesUnregisteredInternalError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()

	apperror.Write(rec, req, errors.New("sql: password leaked in driver message"))

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected http status 500, got %d", rec.Code)
	}
	if body.Code != 500001 {
		t.Fatalf("expected code 500001, got %d", body.Code)
	}
	if body.Message != "系统繁忙，请稍后重试" {
		t.Fatalf("expected safe internal error message, got %q", body.Message)
	}
}

func TestAllowMethodWritesUnifiedChineseError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/login", nil)
	rec := httptest.NewRecorder()

	ok := apperror.AllowMethod(rec, req, http.MethodPost)

	if ok {
		t.Fatal("expected method to be rejected")
	}
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected http status 405, got %d", rec.Code)
	}
	if rec.Header().Get("Allow") != http.MethodPost {
		t.Fatalf("expected Allow header POST, got %q", rec.Header().Get("Allow"))
	}
	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != 405001 {
		t.Fatalf("expected code 405001, got %d", body.Code)
	}
	if body.Message != "请求方法不允许" {
		t.Fatalf("expected Chinese method error message, got %q", body.Message)
	}
}
