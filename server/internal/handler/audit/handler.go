package audithandler

import (
	"context"
	"net/http"
	"strconv"

	"gobaseproject/server/internal/apperror"
	auditmodel "gobaseproject/server/internal/model/audit"
	"gobaseproject/server/pkg/response"
)

type auditService interface {
	ListLoginLogs(ctx context.Context, q auditmodel.LoginLogQuery) (*auditmodel.LoginLogResult, error)
	CleanupLoginLogs(ctx context.Context, days int) (*auditmodel.CleanupResult, error)
	ListOperationLogs(ctx context.Context, q auditmodel.OperationLogQuery) (*auditmodel.OperationLogResult, error)
	GetOperationLogByID(ctx context.Context, id uint64) (*auditmodel.OperationLogRecord, error)
	CleanupOperationLogs(ctx context.Context, days int) (*auditmodel.CleanupResult, error)
}

type Handler struct {
	svc auditService
}

func NewHandler(svc auditService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// specific paths before parameterized
	mux.HandleFunc("DELETE /api/v1/audit/login-logs/cleanup", h.cleanupLoginLogs)
	mux.HandleFunc("DELETE /api/v1/audit/operation-logs/cleanup", h.cleanupOperationLogs)
	mux.HandleFunc("GET /api/v1/audit/login-logs", h.listLoginLogs)
	mux.HandleFunc("GET /api/v1/audit/operation-logs", h.listOperationLogs)
	mux.HandleFunc("GET /api/v1/audit/operation-logs/{id}", h.getOperationLog)
}

// GET /api/v1/audit/login-logs
func (h *Handler) listLoginLogs(w http.ResponseWriter, r *http.Request) {
	q := auditmodel.LoginLogQuery{
		Page:         intParam(r, "page", 1),
		PageSize:     intParam(r, "page_size", 20),
		Keyword:      r.URL.Query().Get("keyword"),
		LoginSuccess: intParam(r, "login_success", 0),
		StartDate:    r.URL.Query().Get("start_date"),
		EndDate:      r.URL.Query().Get("end_date"),
	}
	result, err := h.svc.ListLoginLogs(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

// DELETE /api/v1/audit/login-logs/cleanup
func (h *Handler) cleanupLoginLogs(w http.ResponseWriter, r *http.Request) {
	days := intParam(r, "days", 90)
	result, err := h.svc.CleanupLoginLogs(r.Context(), days)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

// GET /api/v1/audit/operation-logs
func (h *Handler) listOperationLogs(w http.ResponseWriter, r *http.Request) {
	q := auditmodel.OperationLogQuery{
		Page:       intParam(r, "page", 1),
		PageSize:   intParam(r, "page_size", 20),
		Keyword:    r.URL.Query().Get("keyword"),
		Method:     r.URL.Query().Get("method"),
		StatusCode: intParam(r, "status_code", 0),
		StartDate:  r.URL.Query().Get("start_date"),
		EndDate:    r.URL.Query().Get("end_date"),
	}
	result, err := h.svc.ListOperationLogs(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

// GET /api/v1/audit/operation-logs/{id}
func (h *Handler) getOperationLog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidParams)
		return
	}
	rec, err := h.svc.GetOperationLogByID(r.Context(), id)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	if rec == nil {
		apperror.WriteDefinition(w, r, apperror.NotFound)
		return
	}
	response.OK(w, r, rec)
}

// DELETE /api/v1/audit/operation-logs/cleanup
func (h *Handler) cleanupOperationLogs(w http.ResponseWriter, r *http.Request) {
	days := intParam(r, "days", 90)
	result, err := h.svc.CleanupOperationLogs(r.Context(), days)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

func intParam(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

