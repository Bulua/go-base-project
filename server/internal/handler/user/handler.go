package userhandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	usermodel "gobaseproject/server/internal/model/user"
	authservice "gobaseproject/server/internal/service/auth"
	userservice "gobaseproject/server/internal/service/user"
	"gobaseproject/server/pkg/response"
)

type TokenParser interface {
	ParseAccessToken(ctx context.Context, token string) (*authservice.Claims, error)
}

type Handler struct {
	service *userservice.Service
	tokens  TokenParser
}

func NewHandler(service *userservice.Service, tokens TokenParser) *Handler {
	return &Handler{service: service, tokens: tokens}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/users", h.withAuth(h.list))
	mux.HandleFunc("POST /api/v1/users", h.withAuth(h.create))
	mux.HandleFunc("GET /api/v1/users/role-options", h.withAuth(h.roleOptions))
	mux.HandleFunc("GET /api/v1/users/{id}", h.withAuth(h.get))
	mux.HandleFunc("PUT /api/v1/users/{id}", h.withAuth(h.update))
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.withAuth(h.delete))
	mux.HandleFunc("PUT /api/v1/users/{id}/status", h.withAuth(h.updateStatus))
	mux.HandleFunc("PUT /api/v1/users/{id}/password", h.withAuth(h.resetPassword))
	mux.HandleFunc("PUT /api/v1/users/{id}/roles", h.withAuth(h.assignRoles))
}

// ── Endpoints ─────────────────────────────────────────────────────────────

func (h *Handler) list(w http.ResponseWriter, r *http.Request, _ usermodel.ActorContext) {
	q := usermodel.ListQuery{
		Page:       parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:   parseIntDefault(r.URL.Query().Get("page_size"), 20),
		Keyword:    r.URL.Query().Get("keyword"),
		UserStatus: parseIntDefault(r.URL.Query().Get("user_status"), 0),
		RoleID:     parseUint64(r.URL.Query().Get("role_id")),
	}
	result, err := h.service.List(r.Context(), q)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, 500, "list users failed: "+err.Error())
		return
	}
	response.OK(w, r, result)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request, _ usermodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, user)
}

func (h *Handler) roleOptions(w http.ResponseWriter, r *http.Request, _ usermodel.ActorContext) {
	roles, err := h.service.ListRoles(r.Context())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, 500, "list roles failed: "+err.Error())
		return
	}
	response.OK(w, r, roles)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request, actor usermodel.ActorContext) {
	var req usermodel.CreateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
		return
	}
	user, err := h.service.Create(r.Context(), req, actor)
	if err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, user)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, actor usermodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	var req usermodel.UpdateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
		return
	}
	user, err := h.service.Update(r.Context(), id, req, actor)
	if err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, user)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, actor usermodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id, actor); err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request, actor usermodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	var req usermodel.UpdateStatusRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
		return
	}
	if err := h.service.UpdateStatus(r.Context(), id, req.UserStatus, actor); err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request, actor usermodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	var req usermodel.ResetPasswordRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
		return
	}
	if err := h.service.ResetPassword(r.Context(), id, req, actor); err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) assignRoles(w http.ResponseWriter, r *http.Request, actor usermodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	var req usermodel.AssignRolesRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
		return
	}
	if err := h.service.AssignRoles(r.Context(), id, req.RoleIDs, actor); err != nil {
		writeUserError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// ── Middleware ────────────────────────────────────────────────────────────

func (h *Handler) withAuth(next func(http.ResponseWriter, *http.Request, usermodel.ActorContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			response.Error(w, r, http.StatusUnauthorized, 401, "missing authorization token")
			return
		}
		claims, err := h.tokens.ParseAccessToken(r.Context(), token)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, 401, "invalid or expired token")
			return
		}
		actor := usermodel.ActorContext{
			UserID:    claims.UserID,
			LoginName: claims.LoginName,
			SourceIP:  clientIP(r),
			UserAgent: r.UserAgent(),
		}
		next(w, r, actor)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────

func parseIDFromPath(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	raw := r.PathValue("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid user id")
		return 0, false
	}
	return id, true
}

func parseIntDefault(raw string, def int) int {
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}

func parseUint64(raw string) uint64 {
	if raw == "" {
		return 0
	}
	v, _ := strconv.ParseUint(raw, 10, 64)
	return v
}

func writeUserError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, usermodel.ErrUserNotFound):
		response.Error(w, r, http.StatusNotFound, 404, "user not found")
	case errors.Is(err, usermodel.ErrLoginNameTaken):
		response.Error(w, r, http.StatusConflict, 409, "login name already exists")
	case errors.Is(err, usermodel.ErrLoginNameInvalid):
		response.Error(w, r, http.StatusBadRequest, 400, "login name must start with a letter and be 3-32 characters of letters/digits/_.-")
	case errors.Is(err, usermodel.ErrPasswordWeak):
		response.Error(w, r, http.StatusBadRequest, 400, "password must be at least 6 characters")
	case errors.Is(err, usermodel.ErrInvalidStatus):
		response.Error(w, r, http.StatusBadRequest, 400, "invalid user status")
	case errors.Is(err, usermodel.ErrAdminProtected):
		response.Error(w, r, http.StatusForbidden, 403, "built-in admin user is protected")
	default:
		response.Error(w, r, http.StatusInternalServerError, 500, "internal server error: "+err.Error())
	}
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	prefix := "Bearer "
	if !strings.HasPrefix(strings.ToLower(header), strings.ToLower(prefix)) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx > -1 {
		return host[:idx]
	}
	return host
}
