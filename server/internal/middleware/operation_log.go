package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	auditmodel "gobaseproject/server/internal/model/audit"
	authmodel "gobaseproject/server/internal/model/auth"
)

type logWriter interface {
	InsertOperationLog(ctx context.Context, log auditmodel.OperationLog) error
}

// OperationLog records all HTTP requests into gbp_operation_audit_logs.
// It wraps the inner handler, captures response status/body, and writes async.
func OperationLog(repo logWriter, tokens jwtParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldRecord(r.Method) {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()

			// read + restore request body
			var reqBody string
			if r.Body != nil {
				raw, err := io.ReadAll(io.LimitReader(r.Body, 64*1024))
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(raw))
				if err == nil {
					reqBody = redactBody(raw)
				}
			}

			rw := &captureWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)
			costMs := time.Since(start).Milliseconds()

			userID := extractUserID(r, tokens)
			log := auditmodel.OperationLog{
				UserID:        userID,
				SourceIP:      clientIP(r),
				RequestMethod: r.Method,
				RequestPath:   r.URL.Path,
				StatusCode:    rw.status,
				CostMs:        costMs,
				UserAgent:     r.UserAgent(),
				RequestBody:   reqBody,
				ResponseBody:  rw.body.String(),
			}
			go repo.InsertOperationLog(context.Background(), log) //nolint:errcheck
		})
	}
}

func shouldRecord(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	}
	return false
}

// redactBody masks sensitive fields (password, token, secret, key).
func redactBody(raw []byte) string {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return string(raw)
	}
	for k := range m {
		lower := strings.ToLower(k)
		if strings.Contains(lower, "password") ||
			strings.Contains(lower, "token") ||
			strings.Contains(lower, "secret") ||
			strings.Contains(lower, "key") {
			m[k] = "***"
		}
	}
	out, _ := json.Marshal(m)
	return string(out)
}

func extractUserID(r *http.Request, tokens jwtParser) *uint64 {
	token := bearerToken(r)
	if token == "" {
		return nil
	}
	claims, err := tokens.Parse(token, authmodel.TokenTypeAccess)
	if err != nil {
		return nil
	}
	return &claims.UserID
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.Split(fwd, ",")[0]
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i != -1 {
		return host[:i]
	}
	return host
}

type captureWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (cw *captureWriter) WriteHeader(status int) {
	cw.status = status
	cw.ResponseWriter.WriteHeader(status)
}

func (cw *captureWriter) Write(b []byte) (int, error) {
	if cw.body.Len() < 64*1024 {
		cw.body.Write(b)
	}
	return cw.ResponseWriter.Write(b)
}
