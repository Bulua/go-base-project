package authhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	authmodel "gobaseproject/server/internal/model/auth"
	authservice "gobaseproject/server/internal/service/auth"
	"gobaseproject/server/pkg/response"
)

type contextKey struct{}

type Context struct {
	Token  string
	Claims *authservice.Claims
}

type Handler struct {
	service *authservice.Service
}

func NewHandler(service *authservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/auth/login", h.login)
	mux.HandleFunc("/api/v1/auth/logout", h.withAuth(h.logout))
	mux.HandleFunc("/api/v1/auth/refresh", h.refresh)
	mux.HandleFunc("/api/v1/auth/profile", h.withAuth(h.profile))
	mux.HandleFunc("/api/v1/auth/routes", h.withAuth(h.routes))
	mux.HandleFunc("/api/v1/auth/actions", h.withAuth(h.actions))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if !response.AllowMethod(w, r, http.MethodPost) {
		return
	}
	var req authmodel.LoginRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
		return
	}
	session, err := h.service.Login(r.Context(), req, requestMeta(r))
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	response.OK(w, r, session)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request, current Context) {
	if !response.AllowMethod(w, r, http.MethodPost) {
		return
	}
	if err := h.service.Logout(r.Context(), current.Token, "logout"); err != nil {
		writeAuthError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if !response.AllowMethod(w, r, http.MethodPost) {
		return
	}
	var req authmodel.RefreshRequest
	_ = response.ReadJSON(r, &req)
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		refreshToken = bearerToken(r)
	}
	if refreshToken == "" {
		response.Error(w, r, http.StatusUnauthorized, 401, "missing refresh token")
		return
	}
	session, err := h.service.Refresh(r.Context(), refreshToken)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	response.OK(w, r, session)
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request, current Context) {
	if !response.AllowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := h.service.Profile(r.Context(), current.Claims.UserID, current.Claims.LoginName)
	if err != nil {
		writeAuthError(w, r, err)
		return
	}
	response.OK(w, r, user)
}

func (h *Handler) routes(w http.ResponseWriter, r *http.Request, current Context) {
	if !response.AllowMethod(w, r, http.MethodGet) {
		return
	}
	menus, err := h.service.Routes(r.Context(), current.Claims.UserID)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, 500, "query routes failed")
		return
	}
	response.OK(w, r, menus)
}

func (h *Handler) actions(w http.ResponseWriter, r *http.Request, current Context) {
	if !response.AllowMethod(w, r, http.MethodGet) {
		return
	}
	actions, err := h.service.Actions(r.Context(), current.Claims.UserID)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, 500, "query actions failed")
		return
	}
	response.OK(w, r, actions)
}

func (h *Handler) withAuth(next func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			response.Error(w, r, http.StatusUnauthorized, 401, "missing authorization token")
			return
		}
		claims, err := h.service.ParseAccessToken(r.Context(), token)
		if err != nil {
			writeAuthError(w, r, err)
			return
		}
		current := Context{Token: token, Claims: claims}
		ctx := context.WithValue(r.Context(), contextKey{}, current)
		next(w, r.WithContext(ctx), current)
	}
}

func writeAuthError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, authmodel.ErrInvalidCredentials):
		response.Error(w, r, http.StatusUnauthorized, 401, "invalid login name or password")
	case errors.Is(err, authmodel.ErrUserDisabled):
		response.Error(w, r, http.StatusForbidden, 403, "user disabled")
	case errors.Is(err, authmodel.ErrInvalidToken), errors.Is(err, authmodel.ErrTokenBlocked):
		response.Error(w, r, http.StatusUnauthorized, 401, "invalid or expired token")
	default:
		response.Error(w, r, http.StatusInternalServerError, 500, "internal server error")
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

func requestMeta(r *http.Request) authmodel.RequestMeta {
	return authmodel.RequestMeta{
		SourceIP:  clientIP(r),
		UserAgent: r.UserAgent(),
	}
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
