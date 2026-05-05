package authhandler

import (
	"context"
	"net/http"
	"strings"

	"gobaseproject/server/internal/apperror"
	authmodel "gobaseproject/server/internal/model/auth"
	authservice "gobaseproject/server/internal/service/auth"
	"gobaseproject/server/pkg/response"
	"gobaseproject/server/pkg/routereg"
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
	mux.HandleFunc("/api/v1/auth/login", h.login)     // skip rule – not in api_resources
	mux.HandleFunc("/api/v1/auth/refresh", h.refresh) // skip rule – not in api_resources
	mux.HandleFunc("/api/v1/auth/logout", h.withAuth(h.logout))
	mux.HandleFunc("/api/v1/auth/profile", h.withAuth(h.profile))
	mux.HandleFunc("/api/v1/auth/routes", h.withAuth(h.routes))
	mux.HandleFunc("/api/v1/auth/actions", h.withAuth(h.actions))

	routereg.Add("POST", "/api/v1/auth/logout",  "auth", "退出登录")
	routereg.Add("GET",  "/api/v1/auth/profile", "auth", "当前用户信息")
	routereg.Add("GET",  "/api/v1/auth/routes",  "auth", "菜单路由")
	routereg.Add("GET",  "/api/v1/auth/actions", "auth", "按钮权限")
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if !apperror.AllowMethod(w, r, http.MethodPost) {
		return
	}
	var req authmodel.LoginRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	session, err := h.service.Login(r.Context(), req, requestMeta(r))
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, session)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request, current Context) {
	if !apperror.AllowMethod(w, r, http.MethodPost) {
		return
	}
	if err := h.service.Logout(r.Context(), current.Token, "logout"); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) {
	if !apperror.AllowMethod(w, r, http.MethodPost) {
		return
	}
	var req authmodel.RefreshRequest
	_ = response.ReadJSON(r, &req)
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		refreshToken = bearerToken(r)
	}
	if refreshToken == "" {
		apperror.WriteDefinition(w, r, apperror.MissingRefreshToken)
		return
	}
	session, err := h.service.Refresh(r.Context(), refreshToken)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, session)
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request, current Context) {
	if !apperror.AllowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := h.service.Profile(r.Context(), current.Claims.UserID, current.Claims.LoginName)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, user)
}

func (h *Handler) routes(w http.ResponseWriter, r *http.Request, current Context) {
	if !apperror.AllowMethod(w, r, http.MethodGet) {
		return
	}
	menus, err := h.service.Routes(r.Context(), current.Claims.UserID)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, menus)
}

func (h *Handler) actions(w http.ResponseWriter, r *http.Request, current Context) {
	if !apperror.AllowMethod(w, r, http.MethodGet) {
		return
	}
	actions, err := h.service.Actions(r.Context(), current.Claims.UserID)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, actions)
}

func (h *Handler) withAuth(next func(http.ResponseWriter, *http.Request, Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			apperror.WriteDefinition(w, r, apperror.MissingAuthToken)
			return
		}
		claims, err := h.service.ParseAccessToken(r.Context(), token)
		if err != nil {
			apperror.Write(w, r, err)
			return
		}
		current := Context{Token: token, Claims: claims}
		ctx := context.WithValue(r.Context(), contextKey{}, current)
		next(w, r.WithContext(ctx), current)
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
